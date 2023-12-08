package utils

import (
	"fmt"
	"time"
	"strings"
	"unicode"
	tz "4d63.com/tz"
	"encoding/json"
	// "strings"
	"io/ioutil"
	yaml "gopkg.in/yaml.v2"
	// hid "github.com/dh1tw/hid"
	types "github.com/0187773933/ReStreamer/v1/types"
	fiber_cookie "github.com/gofiber/fiber/v2/middleware/encryptcookie"
	encryption "github.com/0187773933/encryption/v1/encryption"
)

type Device struct {
	Name      string
	VendorID  uint16
	ProductID uint16
	Serial string
}
// func GetAllHIDDevices() ( result []Device ) {
// 	devices := hid.Enumerate( 0 , 0 )
// 	devices_map := make( map[ Device ]bool )
// 	for _ , device := range devices {
// 		// fmt.Println( i , device )
// 		// fmt.Printf( "%d === %s %s === %d === %d\n" , i , device.Manufacturer , device.Product , device.VendorID , device.ProductID )
// 		name := strings.TrimSpace( device.Manufacturer + " " + device.Product )
// 		if name == "" { continue }
// 		d := Device{
// 			Name: name ,
// 			VendorID: device.VendorID ,
// 			ProductID: device.ProductID ,
// 			Serial: device.Serial ,
// 		}
// 		devices_map[ d ] = true
// 	}
// 	result = make( []Device , 0 , len( devices_map ) )
// 	for device := range devices_map { result = append( result , device ) }
// 	return
// }

// func GetHIDDevice( serial string ) ( result *hid.Device ) {
// 	devices := hid.Enumerate( 0 , 0 )
// 	for _ , device := range devices {
// 		if device.Serial == serial {
// 			// result = device
// 			// return
// 			result , _ = device.Open()
// 			return
// 		}
// 	}
// 	return
// }

func GetFormattedTimeString() ( result string ) {
	// location , _ := time.LoadLocation( "America/New_York" )
	location , _ := tz.LoadLocation( "America/New_York" )
	time_object := time.Now().In( location )
	month_name := strings.ToUpper( time_object.Format( "Jan" ) )
	milliseconds := time_object.Format( ".000" )
	date_part := fmt.Sprintf( "%02d%s%d" , time_object.Day() , month_name , time_object.Year() )
	time_part := fmt.Sprintf( "%02d:%02d:%02d%s" , time_object.Hour() , time_object.Minute() , time_object.Second() , milliseconds )
	result = fmt.Sprintf( "%s === %s" , date_part , time_part )
	return
}

func RemoveNonASCII( input string ) ( result string ) {
	for _ , i := range input {
		if i > unicode.MaxASCII { continue }
		result += string( i )
	}
	return
}

const SanitizedStringSizeLimit = 100
func SanitizeInputString( input string ) ( result string ) {
	trimmed := strings.TrimSpace( input )
    if len( trimmed ) > SanitizedStringSizeLimit { trimmed = strings.TrimSpace( trimmed[ 0 : SanitizedStringSizeLimit ] ) }
	result = RemoveNonASCII( trimmed )
	return
}


func PrintDevices( devices []Device ) {
	jd , _ := json.MarshalIndent( devices , "" , "  " )
	fmt.Println( string( jd ) )
}

func WriteJSON( filePath string , data interface{} ) {
	file, _ := json.MarshalIndent( data , "" , " " )
	_ = ioutil.WriteFile( filePath , file , 0644 )
}

func ParseConfig( file_path string ) ( result types.ConfigFile ) {
	// file_data , _ := ioutil.ReadFile( file_path )
	// err := json.Unmarshal( file_data , &result )
	// if err != nil { fmt.Println( err ) }
	// return
	file , _ := ioutil.ReadFile( file_path )
	error := yaml.Unmarshal( file , &result )
	if error != nil { panic( error ) }
	return
}

func GenerateNewKeys() {
	fiber_cookie_key := fiber_cookie.GenerateKey()
	bolt_db_key := encryption.GenerateRandomString( 32 )
	server_api_key := encryption.GenerateRandomString( 16 )
	admin_username := encryption.GenerateRandomString( 16 )
	admin_password := encryption.GenerateRandomString( 16 )
	fmt.Println( "Generated New Keys :" )
	fmt.Printf( "\tFiber Cookie Key === %s\n" , fiber_cookie_key )
	fmt.Printf( "\tBolt DB Key === %s\n" , bolt_db_key )
	fmt.Printf( "\tServer API Key === %s\n" , server_api_key )
	fmt.Printf( "\tAdmin Username === %s\n" , admin_username )
	fmt.Printf( "\tAdmin Password === %s\n\n" , admin_password )
}