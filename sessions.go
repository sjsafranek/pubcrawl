package main

import (
	"fmt"
	"net/http"

	"github.com/sjsafranek/gosocialsessions"

	"github.com/sjsafranek/lemur/middleware"
	"github.com/sjsafranek/pubcrawl/lib/clients/websockets"
	"github.com/sjsafranek/pubcrawl/lib/config"
)

var sessionManager = gosocialsessions.New("chocolate-ship", "cookies")

// New returns a new ServeMux with app routes.
func New(conf *config.Config) *http.ServeMux {

	// build app
	mux := http.NewServeMux()

	mux.Handle("/", middleware.Attach(http.HandlerFunc(indexHandler)))
	mux.Handle("/profile", middleware.Attach(sessionManager.RequireLogin(http.HandlerFunc(profileHandler))))
	mux.Handle("/api", middleware.Attach(sessionManager.RequireLogin(http.HandlerFunc(apiWithSessionHandler))))
	mux.Handle("/logout", middleware.Attach(http.HandlerFunc(sessionManager.LogoutHandler)))

	// Static files
	fsvr := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fsvr))

	// Enable FaceBook login
	if conf.OAuth2.HasFacebook() {
		// get facebook login handlers
		loginHandler, callbackHandler := sessionManager.GetFacebookLoginHandlers(
			conf.OAuth2.Facebook.ClientID,
			conf.OAuth2.Facebook.ClientSecret,
			fmt.Sprintf("%v/facebook/callback", conf.Server.GetURLString()))

		// attach facebook login handlers to mux
		mux.Handle("/facebook/login", middleware.Attach(loginHandler))
		mux.Handle("/facebook/callback", middleware.Attach(callbackHandler))
	}

	if conf.OAuth2.HasGoogle() {
		// get facebook login handlers
		loginHandler, callbackHandler := sessionManager.GetGoogleLoginHandlers(
			conf.OAuth2.Google.ClientID,
			conf.OAuth2.Google.ClientSecret,
			fmt.Sprintf("%v/google/callback", conf.Server.GetURLString()))

		// attach facebook login handlers to mux
		mux.Handle("/google/login", middleware.Attach(loginHandler))
		mux.Handle("/google/callback", middleware.Attach(callbackHandler))
	}

	// websockets
	hub, _ := websockets.New(rpcApi)
	mux.Handle("/ws", middleware.Attach(http.HandlerFunc(hub.WebSocketHandler)))

	return mux
}
