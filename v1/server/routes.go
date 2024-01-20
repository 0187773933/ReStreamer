package server

import (
	"fmt"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"net/http"
	// "path/filepath"
	"os/exec"
	// strconv "strconv"
	// types "github.com/0187773933/ReStreamer/v1/types"
	// bolt_api "github.com/boltdb/bolt"
)

func ( s *Server ) Home( context *fiber.Ctx ) ( error ) {
	return context.JSON( fiber.Map{
		"route": "/" ,
		"source": "https://github.com/0187773933/ReStreamer" ,
		"result": "success" ,
	})
}

func ( s *Server ) Test( context *fiber.Ctx ) ( error ) {
    // tikTokURL := "https://pull-hls-f16-va01.tiktokcdn.com/stage/stream-2996525957852168265_or4/index.m3u8" // Replace with the actual URL
    // tikTokURL := "https://www.tiktok.com/t/ZPRvTSNvH"

    kill_ytdlp := exec.Command( "pkill" , "yt-dlp" )
    kill_ytdlp.Run()

    kill_ffmpeg := exec.Command( "pkill" , "ffmpeg" )
    kill_ffmpeg.Run()

    // Construct the command to use yt-dlp and FFmpeg for HLS
    cmdString := "yt-dlp --cookies=/Users/morpheous/Library/CloudStorage/Dropbox/Misc/Cookies/twitch_youtube.txt -q -o - \"https://www.tiktok.com/t/ZPRcc7YkJ\" | ffmpeg -i - -c:v libx264 -preset ultrafast -tune zerolatency -c:a copy -f hls -hls_time 4 -hls_list_size 10 -hls_segment_filename \"./hls-files/stream%03d.ts\" -hls_flags delete_segments ./hls-files/stream.m3u8"
    fmt.Println( cmdString )

    // Execute the command
    cmd := exec.Command( "bash" , "-c" , cmdString )
    if err := cmd.Run(); err != nil {
        return err
    }

    // Redirect to the HLS playlist
    return context.Redirect( "/hls/stream.m3u8" )
}

func ( s *Server ) Que( context *fiber.Ctx ) ( error ) {
    // tikTokURL := "https://pull-hls-f16-va01.tiktokcdn.com/stage/stream-2996525957852168265_or4/index.m3u8" // Replace with the actual URL
    // tikTokURL := "https://www.tiktok.com/t/ZPRvTSNvH"
	api_key := context.Query( "k" )
	if api_key != s.Config.ServerAPIKey {
		return context.Status( fiber.StatusUnauthorized ).SendString( "why" )
	}
	x_url := context.Params( "*" )
	fmt.Println( api_key )
	fmt.Sprintf( fmt.Sprintf( "ReStreamURL( %s )" , x_url ) )

	rm_existing := exec.Command( "rm" , "-rf" , "./hls-files/*" )
	rm_existing.Run()

    kill_ytdlp := exec.Command( "pkill" , "yt-dlp" )
    kill_ytdlp.Run()

    kill_ffmpeg := exec.Command( "pkill" , "ffmpeg" )
    kill_ffmpeg.Run()

    // Construct the command to use yt-dlp and FFmpeg for HLS
    cmdString := "yt-dlp --cookies=/Users/morpheous/Library/CloudStorage/Dropbox/Misc/Cookies/twitch_youtube.txt -q -o - \"" + x_url + "\" | ffmpeg -i - -c:v libx264 -preset ultrafast -tune zerolatency -c:a copy -f hls -hls_time 4 -hls_list_size 10 -hls_segment_filename \"./hls-files/stream%03d.ts\" -hls_flags delete_segments ./hls-files/stream.m3u8"
    fmt.Println( cmdString )

    // Execute the command
    cmd := exec.Command( "bash" , "-c" , cmdString )
    if err := cmd.Run(); err != nil {
        return err
    }

    // Redirect to the HLS playlist
    return context.Redirect( "/hls/stream.m3u8" )
}

func ( s *Server ) SetupRoutes() {

	// admin_route_group := s.FiberApp.Group( "/admin" )

	// // HTML UI Pages
	// admin_route_group.Get( "/login" , ServeLoginPage )
	// for url , _ := range ui_html_pages {
	// 	admin_route_group.Get( url , ServeAuthenticatedPage )
	// }

	s.FiberApp.Get( "/" , s.Home )
	s.FiberApp.Get( "/test" , s.Test )
	s.FiberApp.Get( "/que/url/*" , s.Que )
	s.FiberApp.Use( "/hls" , filesystem.New( filesystem.Config{
		Root: http.Dir( "./hls-files" ) ,
		Browse: false ,
		Index: "" ,
		MaxAge: 3600 ,
		PathPrefix: "" ,
	}))
}