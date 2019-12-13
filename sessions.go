package main

import (
	"net/http"

	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/facebook"
	"github.com/dghubble/sessions"
	"golang.org/x/oauth2"
	facebookOAuth2 "golang.org/x/oauth2/facebook"

	// "github.com/sjsafranek/pubcrawl/foursquare"
	"github.com/sjsafranek/lemur/middleware"
)

const (
	sessionName   = "chocolate-ship"
	sessionSecret = "cookies"
)

// sessionStore encodes and decodes session data stored in signed cookies
var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

// New returns a new ServeMux with app routes.
func New(config *Config) *http.ServeMux {

	// build app
	mux := http.NewServeMux()

	mux.Handle("/", middleware.Attach(http.HandlerFunc(welcomeHandler)))
	mux.Handle("/profile", middleware.Attach(requireLogin(http.HandlerFunc(profileHandler))))
	mux.Handle("/logout", middleware.Attach(http.HandlerFunc(logoutHandler)))

	// 1. Register Login and Callback handlers
	oauth2Config := &oauth2.Config{
		ClientID:     config.FacebookClientID,
		ClientSecret: config.FacebookClientSecret,
		RedirectURL:  "http://localhost:8080/facebook/callback",
		Endpoint:     facebookOAuth2.Endpoint,
		Scopes:       []string{"email"},
	}

	// state param cookies require HTTPS by default; disable for localhost development
	stateConfig := gologin.DebugOnlyCookieConfig
	mux.Handle("/facebook/login", middleware.Attach(facebook.StateHandler(stateConfig, facebook.LoginHandler(oauth2Config, nil))))
	mux.Handle("/facebook/callback", middleware.Attach(facebook.StateHandler(stateConfig, facebook.CallbackHandler(oauth2Config, issueSession(), nil))))
	return mux
}

// issueSession issues a cookie session after successful Facebook login
func issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		facebookUser, err := facebook.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 2. Implement a success handler to issue some form of session
		session := sessionStore.New(sessionName)
		session.Values["facebook-id"] = facebookUser.ID
		session.Values["facebook-name"] = facebookUser.Name
		session.Save(w)
		http.Redirect(w, req, "/profile", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

// isAuthenticated returns true if the user has a signed session cookie.
func isAuthenticated(req *http.Request) bool {
	if _, err := sessionStore.Get(req, sessionName); err == nil {
		return true
	}
	return false
}
