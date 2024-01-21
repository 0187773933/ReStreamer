package server

import (
	"fmt"
	"time"
	"os"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	rate_limiter "github.com/gofiber/fiber/v2/middleware/limiter"
	"net/http"
	"os/exec"
	"strings"
)

var public_limiter = rate_limiter.New( rate_limiter.Config{
	Max: 1 ,
	Expiration: 1 * time.Second ,
	KeyGenerator: func( c *fiber.Ctx ) string {
		return c.Get( "x-forwarded-for" )
	} ,
	LimitReached: func( c *fiber.Ctx ) error {
		ip_address := c.IP()
		log_message := fmt.Sprintf( "%s === %s === %s === PUBLIC RATE LIMIT REACHED !!!" , ip_address , c.Method() , c.Path() );
		fmt.Println( log_message )
		c.Set( "Content-Type" , "text/html" )
		return c.SendString( "<html><h1>loading ...</h1><script>setTimeout(function(){ window.location.reload(1); }, 6);</script></html>" )
	} ,
})

func ( s *Server ) Home( context *fiber.Ctx ) ( error ) {
	return context.JSON( fiber.Map{
		"route": "/" ,
		"source": "https://github.com/0187773933/ReStreamer" ,
		"result": "success" ,
	})
}

func ( s *Server ) Que( context *fiber.Ctx ) ( error ) {
	api_key := context.Query( "k" )
	if api_key != s.Config.ServerAPIKey {
		return context.Status( fiber.StatusUnauthorized ).SendString( "why" )
	}
	x_url := context.Params( "*" )
	fmt.Printf( "ReStreamURL( %s )\n" , x_url )

	fmt.Println( "Killing yt-dlp" )
	kill_ytdlp := exec.Command( "pkill" , "yt-dlp" )
	kill_ytdlp.Run()

	fmt.Println( "Killing ffmpeg" )
	kill_ffmpeg := exec.Command( "pkill" , "ffmpeg" )
	kill_ffmpeg.Run()

	time.Sleep( 500 * time.Millisecond )
	os.MkdirAll( "./hls-files" , os.ModePerm )
	fmt.Println( "Removing Existing HLS Files" )
	rm_existing := exec.Command( "bash" , "-c" , "rm -rf ./hls-files/*" )
	rm_existing.Run()

	fmt.Println( "getting live url" )
	var live_url_cmd *exec.Cmd
	if s.Config.CookiesFilePath == "" {
		live_url_cmd = exec.Command( "yt-dlp" , "--cookies" , s.Config.CookiesFilePath , "-g" , x_url )
	} else {
		live_url_cmd = exec.Command( "yt-dlp" , "-g" , x_url )
	}
	live_url_cmd_output , _ := live_url_cmd.Output()
	live_url_cmd_output_string := strings.TrimSpace( string( live_url_cmd_output ) )
	fmt.Println( live_url_cmd_output_string )

	cmd_string := "ffmpeg -i \"" + live_url_cmd_output_string + "\" -c:v libx264 -preset ultrafast -tune zerolatency -c:a copy -f hls -hls_time 4 -hls_list_size 10 -hls_segment_filename \"./hls-files/stream%03d.ts\" -hls_flags delete_segments ./hls-files/stream.m3u8"
	fmt.Println( cmd_string )
	cmd := exec.Command( "bash" , "-c" , cmd_string )
	go cmd.Run()
	stream_url := fmt.Sprintf( "/%s/stream.m3u8" , s.Config.HLSURLPrefix )
	return context.JSON( fiber.Map{
		"url": "/que/url/*" ,
		"param_url": x_url ,
		"stream_url": stream_url ,
		"result": true ,
	})
}

func ( s *Server ) Stop( context *fiber.Ctx ) ( error ) {
	api_key := context.Query( "k" )
	if api_key != s.Config.ServerAPIKey {
		return context.Status( fiber.StatusUnauthorized ).SendString( "why" )
	}
	fmt.Println( "Killing yt-dlp" )
	kill_ytdlp := exec.Command( "pkill" , "yt-dlp" )
	kill_ytdlp.Run()
	fmt.Println( "Killing ffmpeg" )
	kill_ffmpeg := exec.Command( "pkill" , "ffmpeg" )
	kill_ffmpeg.Run()
	time.Sleep( 500 * time.Millisecond )
	fmt.Println( "Removing Existing HLS Files" )
	rm_existing := exec.Command( "bash" , "-c" , "rm -rf ./hls-files/*" )
	rm_existing.Run()
	return context.JSON( fiber.Map{
		"url": "/stop" ,
		"result": true ,
	})
}

func ( s *Server ) SetupRoutes() {
	s.FiberApp.Get( "/" , public_limiter , s.Home )
	s.FiberApp.Get( "/que/url/*" , public_limiter , s.Que )
	s.FiberApp.Get( "/stop" , public_limiter , s.Stop )
	s.FiberApp.Use( fmt.Sprintf( "/%s" , s.Config.HLSURLPrefix ) , filesystem.New( filesystem.Config{
		Root: http.Dir( "./hls-files" ) ,
		Browse: false ,
		Index: "" ,
		MaxAge: 3600 ,
		PathPrefix: "" ,
	}))
}