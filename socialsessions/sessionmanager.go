package socialsessions

import (
	"net/http"

	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/facebook"
	"github.com/dghubble/sessions"
	"golang.org/x/oauth2"
	facebookOAuth2 "golang.org/x/oauth2/facebook"
)

type SessionManager struct {
	sessionName   string
	sessionSecret string
	sessionStore  *sessions.CookieStore
}

func New(sessionName, sessionSecret string) *SessionManager {
	// sessionStore encodes and decodes session data stored in signed cookies
	return &SessionManager{
		sessionName:   sessionName,
		sessionSecret: sessionSecret,
		sessionStore:  sessions.NewCookieStore([]byte(sessionSecret), nil),
	}
}

func (self *SessionManager) Get(req *http.Request) (*sessions.Session, error) {
	return self.sessionStore.Get(req, self.sessionName)
}

func (self *SessionManager) issueSession() *sessions.Session {
	return self.sessionStore.New(self.sessionName)
}

func (self *SessionManager) GetFacebookLoginHandlers(clientID, clientSecret, callbackUrl string) (http.Handler, http.Handler) {
	// 1. Register Login and Callback handlers
	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  callbackUrl,
		Endpoint:     facebookOAuth2.Endpoint,
		Scopes:       []string{"email"},
	}

	// state param cookies require HTTPS by default; disable for localhost development
	stateConfig := gologin.DebugOnlyCookieConfig
	loginHandler := facebook.StateHandler(stateConfig, facebook.LoginHandler(oauth2Config, nil))
	callbackHandler := facebook.StateHandler(stateConfig, facebook.CallbackHandler(oauth2Config, self.issueFacebookSession(), nil))
	return loginHandler, callbackHandler
}

// issueSession issues a cookie session after successful Facebook login
func (self *SessionManager) issueFacebookSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		facebookUser, err := facebook.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 2. Implement a success handler to issue some form of session
		session := self.issueSession()
		session.Values["userid"] = facebookUser.ID
		session.Values["username"] = facebookUser.Name
		session.Values["usertype"] = "facebook"
		session.Save(w)
		http.Redirect(w, req, "/profile", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

func (self *SessionManager) destroySession(w http.ResponseWriter) {
	self.sessionStore.Destroy(w, self.sessionName)
}

// logoutHandler destroys the session on POSTs and redirects to home.
func (self *SessionManager) LogoutHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		self.destroySession(w)
	}
	http.Redirect(w, req, "/", http.StatusFound)
}

// requireLogin redirects unauthenticated users to the login route.
func (self *SessionManager) RequireLogin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if !self.IsAuthenticated(req) {
			http.Redirect(w, req, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

// isAuthenticated returns true if the user has a signed session cookie.
func (self *SessionManager) IsAuthenticated(req *http.Request) bool {
	if _, err := self.sessionStore.Get(req, self.sessionName); err == nil {
		return true
	}
	return false
}
