package main

import (
	"net/http"

	"github.com/sjsafranek/lemur/middleware"
	"github.com/sjsafranek/pubcrawl/lib/config"
	"github.com/sjsafranek/pubcrawl/lib/socialsessions"
)

var sessionManager = socialsessions.New("chocolate-ship", "cookies")

// New returns a new ServeMux with app routes.
func New(conf *config.Config) *http.ServeMux {

	// build app
	mux := http.NewServeMux()

	mux.Handle("/", middleware.Attach(http.HandlerFunc(welcomeHandler)))
	mux.Handle("/profile", middleware.Attach(sessionManager.RequireLogin(http.HandlerFunc(profileHandler))))
	mux.Handle("/api/v1/new_crawl", middleware.Attach(sessionManager.RequireLogin(http.HandlerFunc(newPubCrawlHandler))))
	mux.Handle("/api/v1/get_crawl", middleware.Attach(sessionManager.RequireLogin(http.HandlerFunc(getPubCrawlHandler))))
	mux.Handle("/logout", middleware.Attach(http.HandlerFunc(sessionManager.LogoutHandler)))

	// Static files
	fsvr := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fsvr))

	// get facebook login handlers
	loginHandler, callbackHandler := sessionManager.GetFacebookLoginHandlers(
		conf.Facebook.ClientID,
		conf.Facebook.ClientSecret,
		"http://localhost:8080/facebook/callback")

	// attach facebook login handlers to mux
	mux.Handle("/facebook/login", middleware.Attach(loginHandler))
	mux.Handle("/facebook/callback", middleware.Attach(callbackHandler))

	return mux
}
