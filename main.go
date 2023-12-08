package main

import (
	"fmt"
	"os"
	"time"
	"os/signal"
	"syscall"
	"path/filepath"
	utils "github.com/0187773933/ReStreamer/v1/utils"
	server "github.com/0187773933/ReStreamer/v1/server"
	bolt_api "github.com/boltdb/bolt"
)

var s server.Server

func SetupCloseHandler() {
	c := make( chan os.Signal )
	signal.Notify( c , os.Interrupt , syscall.SIGTERM , syscall.SIGINT )
	go func() {
		<-c
		fmt.Println( "\r- Ctrl+C pressed in Terminal" )
		fmt.Println( "Shutting Down Re-Streaming Server" )
		s.FiberApp.Shutdown()
		os.Exit( 0 )
	}()
}

func main() {

	config_file_path := "./config.yaml"
	if len( os.Args ) > 1 { config_file_path , _ = filepath.Abs( os.Args[ 1 ] ) }
	config := utils.ParseConfig( config_file_path )
	fmt.Printf( "Loaded Config File From : %s\n" , config_file_path )

	// 1.) Setup StreamDeck
	db , _ := bolt_api.Open( config.BoltDBPath , 0600 , &bolt_api.Options{ Timeout: ( 3 * time.Second ) } )
	fmt.Println( db )

	// 2.) Start Server
	SetupCloseHandler()
	// utils.GenerateNewKeys()
	s = server.New( config )
	s.Start()

}
