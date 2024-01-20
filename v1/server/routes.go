package server

import (
	"fmt"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"net/http"
	"os/exec"
	"strings"
)

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

	fmt.Println( "Removing Existing HLS Files" )
	rm_existing := exec.Command( "bash" , "-c" , "rm -rf ./hls-files/*" )
	rm_existing.Run()

	fmt.Println( "Killing yt-dlp" )
	kill_ytdlp := exec.Command( "pkill" , "yt-dlp" )
	kill_ytdlp.Run()

	fmt.Println( "Killing ffmpeg" )
	kill_ffmpeg := exec.Command( "pkill" , "ffmpeg" )
	kill_ffmpeg.Run()

	fmt.Println( "getting live url" )
	cookie_file_path := "/Users/morpheous/Library/CloudStorage/Dropbox/Misc/Cookies/twitch_youtube.txt"
	// live_url_cmd_string := "yt-dlp --cookies=/Users/morpheous/Library/CloudStorage/Dropbox/Misc/Cookies/twitch_youtube.txt -q -o - \"" + x_url + "\""
	// live_url_cmd := exec.Command( "bash" , "-c" , live_url_cmd_string )
	live_url_cmd := exec.Command( "yt-dlp" , "--cookies" , cookie_file_path , "-g" , x_url )
	live_url_cmd_output , _ := live_url_cmd.Output()
	live_url_cmd_output_string := strings.TrimSpace( string( live_url_cmd_output ) )
	fmt.Println( live_url_cmd_output_string )

	// cmdString := "yt-dlp --cookies=/Users/morpheous/Library/CloudStorage/Dropbox/Misc/Cookies/twitch_youtube.txt -q -o - \"" + x_url + "\" | ffmpeg -i - -c:v libx264 -preset ultrafast -tune zerolatency -c:a copy -f hls -hls_time 4 -hls_list_size 10 -hls_segment_filename \"./hls-files/stream%03d.ts\" -hls_flags delete_segments ./hls-files/stream.m3u8"
	// fmt.Println( cmdString )

	cmd_string := "ffmpeg -i \"" + live_url_cmd_output_string + "\" -c:v libx264 -preset ultrafast -tune zerolatency -c:a copy -f hls -hls_time 4 -hls_list_size 10 -hls_segment_filename \"./hls-files/stream%03d.ts\" -hls_flags delete_segments ./hls-files/stream.m3u8"
	fmt.Println( cmd_string )
	cmd := exec.Command( "bash" , "-c" , cmd_string )
	go cmd.Run()
	return context.Redirect( "/hls/stream.m3u8" )
}

func ( s *Server ) SetupRoutes() {
	s.FiberApp.Get( "/" , s.Home )
	s.FiberApp.Get( "/que/url/*" , s.Que )
	s.FiberApp.Use( "/hls" , filesystem.New( filesystem.Config{
		Root: http.Dir( "./hls-files" ) ,
		Browse: false ,
		Index: "" ,
		MaxAge: 3600 ,
		PathPrefix: "" ,
	}))
}