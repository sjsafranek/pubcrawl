package main

import (
	"errors"

		"github.com/sjsafranek/pubcrawl/lib/api"
	"github.com/sjsafranek/logger"
	"html/template"
	"net/http"
)
//
// var (
// 	LOGIN_TEMPLATE   *template.Template = template.Must(template.ParseFiles("tmpl/global_header.html", "tmpl/global_footer.html", "tmpl/login.html"))
// 	PROFILE_TEMPLATE *template.Template = template.Must(template.ParseFiles("tmpl/global_header.html", "tmpl/global_footer.html", "tmpl/navbar.html", "tmpl/profile.html"))
// )
//
// // welcomeHandler shows a welcome message and login button.
// func welcomeHandler(w http.ResponseWriter, req *http.Request) {
// 	if req.URL.Path != "/" {
// 		http.NotFound(w, req)
// 		return
// 	}
//
// 	if sessionManager.IsAuthenticated(req) {
// 		http.Redirect(w, req, "/profile", http.StatusFound)
// 		return
// 	}
//
// 	err := LOGIN_TEMPLATE.ExecuteTemplate(w, "login", nil)
// 	if nil != err {
// 		logger.Error(err)
// 		apiBasicResponse(w, http.StatusInternalServerError, err)
// 	}
// }
//
// // profileHandler shows protected user content.
// func profileHandler(w http.ResponseWriter, req *http.Request) {
// 	val, _ := sessionManager.Get(req)
// 	username := val.Values["username"].(string)
// 	usertype := val.Values["usertype"].(string)
// 	userid := val.Values["userid"].(string)
// 	useremail := val.Values["useremail"].(string)
//
// 	args := make(map[string]interface{})
// 	args["username"] = username
//
// 	user, err := rpcApi.GetDatabase().CreateUserIfNotExists(useremail, useremail)
// 	if nil != err {
// 		logger.Error(err)
// 		apiBasicResponse(w, http.StatusInternalServerError, err)
// 	}
// 	user.CreateSocialAccountIfNotExists(userid, username, usertype)
//
// 	err = PROFILE_TEMPLATE.ExecuteTemplate(w, "profile", args)
// 	if nil != err {
// 		logger.Error(err)
// 		apiBasicResponse(w, http.StatusInternalServerError, err)
// 	}
// }



var (
	LOGIN_TEMPLATE   *template.Template = template.Must(template.ParseFiles("tmpl/global_header.html", "tmpl/global_footer.html", "tmpl/login.html"))

	PROFILE_TEMPLATE *template.Template = template.Must(template.ParseFiles("tmpl/global_header.html", "tmpl/global_footer.html", "tmpl/navbar.html", "tmpl/profile.html"))
)

func  executeLoginTemplate(w http.ResponseWriter, options map[string]interface{}) {
	logger.Info(options)
	err := LOGIN_TEMPLATE.ExecuteTemplate(w, "login", options)
	if nil != err {
		logger.Error(err)
		apiBasicResponse(w, http.StatusInternalServerError, err)
	}
}

func  getHandlerOptions(r *http.Request) map[string]interface{} {
	options := make(map[string]interface{})

	oauth2Options := make(map[string]bool)
	oauth2Options["facebook"] = conf.OAuth2.HasFacebook()
	oauth2Options["google"] = conf.OAuth2.HasGoogle()
	options["oauth2"] = oauth2Options

	val, _ := sessionManager.Get(r)
	if nil != val {
		useremail := val.Values["useremail"]
		username := useremail.(string)
		userOptions := make(map[string]string)
		userOptions["username"] = username
		results, err := rpcApi.Do(&api.Request{Method: "get_user", Params: &api.RequestParams{Username: username}})
		if nil == err {
			userOptions["apikey"] = results.Data.User.Apikey
		}

		options["user"] = userOptions
	}

	return options
}

// welcomeHandler shows a welcome message and login button.
func  indexHandler(w http.ResponseWriter, r *http.Request) {

	if sessionManager.IsAuthenticated(r) {
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}

	options := getHandlerOptions(r)

	if "POST" == r.Method {

		username := r.FormValue("username")
		password := r.FormValue("password")
		if "" == username && "" == password {
			usr, psw, ok := r.BasicAuth()
			if !ok {
				err := errors.New("Unable to get credentials")
				options["error"] = err.Error()
				executeLoginTemplate(w, options)
				return
			}
			username = usr
			password = psw
		}

		results, err := rpcApi.Do(&api.Request{Method: "get_user", Params: &api.RequestParams{Username: username}})
		if nil != err {
			options["error"] = err.Error()
			executeLoginTemplate(w, options)
			return
		}

		is_password, _ := results.Data.User.IsPassword(password)
		if !is_password {
			err = errors.New("Incorrect password")
			options["error"] = err.Error()
			executeLoginTemplate(w, options)
			return
		}

		session := sessionManager.IssueSession()
		session.Values["userid"] = ""
		session.Values["username"] = results.Data.User.Username
		session.Values["useremail"] = results.Data.User.Email
		session.Values["usertype"] = "find5"
		session.Save(w)
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}

	executeLoginTemplate(w, options)
}

// profileHandler shows protected user content.
func  profileHandler(w http.ResponseWriter, r *http.Request) {
	val, _ := sessionManager.Get(r)
	username := val.Values["username"].(string)
	usertype := val.Values["usertype"].(string)
	userid := val.Values["userid"].(string)
	useremail := val.Values["useremail"].(string)

	// options := make(map[string]interface{})
	// options["username"] = username
	options := getHandlerOptions(r)

	user, err := rpcApi.GetDatabase().CreateUserIfNotExists(useremail, useremail)
	if nil != err {
		logger.Error(err)
		apiBasicResponse(w, http.StatusInternalServerError, err)
	}
	user.CreateSocialAccountIfNotExists(userid, username, usertype)

	err = PROFILE_TEMPLATE.ExecuteTemplate(w, "profile", options)
	if nil != err {
		logger.Error(err)
		apiBasicResponse(w, http.StatusInternalServerError, err)
	}
}
