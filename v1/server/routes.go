package server

import (
	// "fmt"
	fiber "github.com/gofiber/fiber/v2"
	// strconv "strconv"
	// types "github.com/0187773933/ReStreamer/v1/types"
	// bolt_api "github.com/boltdb/bolt"
)

// var ui_html_pages = map[ string ]string {
// 	"/": "./v1/server/html/admin.html" ,
// 	"/users": "./v1/server/html/admin_view_users.html" ,
// 	// "/user/new": "./v1/server/html/admin_user_new.html" ,
// 	"/user/new/handoff/:uuid": "./v1/server/html/admin_user_new_handoff.html" ,
// 	"/user/checkin": "./v1/server/html/admin_user_checkin.html" ,
// 	"/user/checkin/:uuid": "./v1/server/html/admin_user_checkin.html" ,
// 	"/user/checkin/:uuid/edit": "./v1/server/html/admin_user_checkin.html" ,
// 	"/user/checkin/new": "./v1/server/html/admin_user_checkin.html" ,
// 	"/user/edit/:uuid": "./v1/server/html/admin_user_edit.html" ,
// 	"/checkins": "./v1/server/html/admin_view_total_checkins.html" ,
// 	"/emails": "./v1/server/html/admin_view_all_emails.html" ,
// 	"/phone-numbers": "./v1/server/html/admin_view_all_phone_numbers.html" ,
// 	"/barcodes": "./v1/server/html/admin_view_all_barcodes.html" ,
// 	"/sms": "./v1/server/html/admin_sms_all_users.html" ,
// 	"/email": "./v1/server/html/admin_email_all_users.html" ,
// }

// func ( s *Server ) PressButton( context *fiber.Ctx ) ( error ) {
// 	if validate_admin( context ) == false { return serve_failed_attempt( context ) }

// 	// s.UI.ActivePageID = "spotify-triple"
// 	// s.UI.Render()

// 	button_type := ""
// 	button_int , button_int_err := strconv.Atoi( context.Params( "button" ) )
// 	if button_int_err == nil {
// 		// Check if the integer is within the uint8 range (0-255)
// 		if button_int >= 0 && button_int <= 255 {
// 			button_type = "number"
// 		}
// 	} else {
// 		button_type = "string"
// 	}

// 	switch button_type {
// 		case "number":
// 			s.UI.SingleClickNumber( uint8( button_int ) )
// 			break;
// 		case "string":
// 			s.UI.SingleClickId( context.Params( "button" ) )
// 			break;
// 	}

// 	return context.JSON( fiber.Map{
// 		"route": "/:button" ,
// 		"type": button_type ,
// 		"id": context.Params( "button" ) ,
// 		"result": "success" ,
// 	})
// }

// func ( s *Server ) RenderPage( context *fiber.Ctx ) ( error ) {
// 	if validate_admin( context ) == false { return serve_failed_attempt( context ) }
// 	page_id := context.Params( "id" )
// 	s.UI.Clear()

// 	// s.UI.SetActivePageID( string( page_id ) )
// 	// race condition ????
// 	s.UI.DB.Update( func( tx *bolt_api.Tx ) error {
// 		tmp2_bucket , _ := tx.CreateBucketIfNotExists( []byte( "tmp2" ) )
// 		tmp2_bucket.Put( []byte( "active-page-id" ) , []byte( page_id ) )
// 		return nil
// 	})
// 	s.UI.Render()
// 	return context.JSON( fiber.Map{
// 		"route": "/page/:id" ,
// 		"id": page_id ,
// 		"result": "success" ,
// 	})
// }

func ( s *Server ) Home( context *fiber.Ctx ) ( error ) {
	return context.JSON( fiber.Map{
		"route": "/" ,
		"source": "https://github.com/0187773933/ReStreamer" ,
		"result": "success" ,
	})
}

func ( s *Server ) SetupRoutes() {

	// admin_route_group := s.FiberApp.Group( "/admin" )

	// // HTML UI Pages
	// admin_route_group.Get( "/login" , ServeLoginPage )
	// for url , _ := range ui_html_pages {
	// 	admin_route_group.Get( url , ServeAuthenticatedPage )
	// }

	s.FiberApp.Get( "/" , s.Home )
	// s.FiberApp.Get( "/page/:id" , s.RenderPage )

}