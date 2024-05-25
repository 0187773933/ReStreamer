package server

import (
	"fmt"
	"time"
	fiber "github.com/gofiber/fiber/v2"
	bcrypt "golang.org/x/crypto/bcrypt"
	encryption "github.com/0187773933/encryption/v1/encryption"
)

func validate_login_credentials( context *fiber.Ctx ) ( result bool ) {
	result = false
	uploaded_username := context.FormValue( "username" )
	if uploaded_username == "" { fmt.Println( "username empty" ); return }
	if uploaded_username != GlobalConfig.AdminUsername { fmt.Println( "username not correct" ); return }
	uploaded_password := context.FormValue( "password" )
	if uploaded_password == "" { fmt.Println( "password empty" ); return }
	password_matches := bcrypt.CompareHashAndPassword( []byte( uploaded_password ) , []byte( GlobalConfig.AdminPassword ) )
	if password_matches != nil { fmt.Println( "bcrypted password doesn't match" ); return }
	result = true
	return
}

func validate_inernal_auth_mw( context *fiber.Ctx ) ( error ) {
	k_param := context.Params( "k" )
	if k_param == "" {
		return context.Status( fiber.StatusUnauthorized ).SendString( "why" )
	}
	if len( k_param ) > 50 {
		return context.Status( fiber.StatusUnauthorized ).SendString( "why" )
	}
	if k_param != GlobalConfig.ServerAPIKey {
		return context.Status( fiber.StatusUnauthorized ).SendString( "why" )
	}
	return context.Next()
}

func Logout( context *fiber.Ctx ) ( error ) {
	context.Cookie( &fiber.Cookie{
		Name: GlobalConfig.ServerCookieName ,
		Value: "" ,
		Expires: time.Now().Add( -time.Hour ) , // set the expiration to the past
		HTTPOnly: true ,
		Secure: true ,
	})
	context.Set( "Content-Type" , "text/html" )
	return context.SendString( "<h1>Logged Out</h1>" )
}

// POST http://localhost:5950/admin/login
func Login( context *fiber.Ctx ) ( error ) {
	valid_login := validate_login_credentials( context )
	if valid_login == false { return serve_failed_attempt( context ) }
	context.Cookie(
		&fiber.Cookie{
			Name: GlobalConfig.ServerCookieName ,
			Value: encryption.SecretBoxEncrypt( GlobalConfig.BoltDBEncryptionKey , GlobalConfig.ServerCookieAdminSecretMessage ) ,
			Secure: true ,
			Path: "/" ,
			// Domain: "blah.ngrok.io" , // probably should set this for webkit
			HTTPOnly: true ,
			SameSite: "Lax" ,
			Expires: time.Now().AddDate( 10 , 0 , 0 ) , // aka 10 years from now
		} ,
	)
	return context.Redirect( "/" )
}

func validate_admin_cookie( context *fiber.Ctx ) ( result bool ) {
	result = false
	admin_cookie := context.Cookies( GlobalConfig.ServerCookieName )
	if admin_cookie == "" { fmt.Println( "admin cookie was blank" ); return }
	admin_cookie_value := encryption.SecretBoxDecrypt( GlobalConfig.BoltDBEncryptionKey , admin_cookie )
	if admin_cookie_value != GlobalConfig.ServerCookieAdminSecretMessage { fmt.Println( "admin cookie secret message was not equal" ); return }
	result = true
	return
}

func validate_admin( context *fiber.Ctx ) ( result bool ) {
	result = false
	admin_cookie := context.Cookies( GlobalConfig.ServerCookieName )
	if admin_cookie != "" {
		admin_cookie_value := encryption.SecretBoxDecrypt( GlobalConfig.BoltDBEncryptionKey , admin_cookie )
		if admin_cookie_value == GlobalConfig.ServerCookieAdminSecretMessage {
			result = true
			return
		}
	}
	admin_api_key_header := context.Get( "key" )
	if admin_api_key_header != "" {
		if admin_api_key_header == GlobalConfig.ServerAPIKey {
			result = true
			return
		}
	}
	admin_api_key_query := context.Query( "k" )
	if admin_api_key_query != "" {
		if admin_api_key_query == GlobalConfig.ServerAPIKey {
			result = true
			return
		}
	}
	return
}

func RenderLoginPage( context *fiber.Ctx ) ( error ) {
	context.Set( "Content-Type" , "text/html" )
	return context.SendFile( "./v1/server/html/login.html" )
}

// need somthing different for hls .ts parts
// too slow
func validate_admin_mw( context *fiber.Ctx ) ( error ) {
	admin_cookie := context.Cookies( GlobalConfig.ServerCookieName )
	if admin_cookie != "" {
		admin_cookie_value := encryption.SecretBoxDecrypt( GlobalConfig.BoltDBEncryptionKey , admin_cookie )
		if admin_cookie_value == GlobalConfig.ServerCookieAdminSecretMessage {
			return context.Next()
		}
	}
	admin_api_key_header := context.Get( "key" )
	if admin_api_key_header != "" {
		if admin_api_key_header == GlobalConfig.ServerAPIKey {
			return context.Next()
		}
	}
	admin_api_key_query := context.Query( "k" )
	if admin_api_key_query != "" {
		if admin_api_key_query == GlobalConfig.ServerAPIKey {
			return context.Next()
		}
	}
	return context.Status( fiber.StatusUnauthorized ).SendString( "why" )
}