package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/sjsafranek/logger"
	"github.com/sjsafranek/pubcrawl/lib/api"
)

// func apiHandler(w http.ResponseWriter, r *http.Request) {
// 	var data string
//
// 	val, _ := sessionManager.Get(r)
// 	useremail := val.Values["useremail"].(string)
// 	if 0 == len(useremail) {
// 		apiBasicResponse(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
// 		return
// 	}
//
// 	status_code, err := func() (int, error) {
// 		switch r.Method {
// 		case "POST":
// 			body, err := ioutil.ReadAll(r.Body)
// 			if err != nil {
// 				return http.StatusBadRequest, err
// 			}
// 			r.Body.Close()
//
// 			var request api.Request
//
// 			// TODO
// 			//  - use request.Unmarshal
// 			json.Unmarshal(body, &request)
//
// 			// protect against request hijacking!
// 			request.Params.Username = useremail
// 			request.Params.Apikey = ""
//
// 			logger.Info(string(body))
// 			// logger.Info(request)
//
// 			if !conf.Api.IsPublicMethod(request.Method) {
// 				logger.Warnf("Not a public api method: %v", request.Method)
// 				return http.StatusMethodNotAllowed, errors.New(http.StatusText(http.StatusMethodNotAllowed))
// 			}
//
// 			// run api request
// 			response, err := rpcApi.Do(&request)
//
// 			if nil != err {
// 				results, _ := response.Marshal()
// 				data = results
// 				return http.StatusBadRequest, nil
// 			}
//
// 			results, _ := response.Marshal()
// 			data = results
// 			return http.StatusOK, nil
//
// 		default:
// 			return http.StatusMethodNotAllowed, errors.New(http.StatusText(http.StatusMethodNotAllowed))
// 		}
// 	}()
//
// 	if nil != err {
// 		apiBasicResponse(w, status_code, err)
// 		return
// 	}
//
// 	apiJSONResponse(w, []byte(data), status_code)
// }


func getApiRequestFromHttpRequest(r *http.Request) (*api.Request, error) {
	var request api.Request

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &request, err
	}
	r.Body.Close()

	// TODO
	//  - use request.Unmarshal
	return &request, json.Unmarshal(body, &request)
}

func execApiRequest(request *api.Request) (string, int, error) {
	// run api request
	response, err := rpcApi.Do(request)
	results, _ := response.Marshal()

	if nil != err {
		return results, http.StatusBadRequest, nil
	}

	return results, http.StatusOK, nil
}

func apiHandler(w http.ResponseWriter, r *http.Request) {

	var data string

	status_code, err := func() (int, error) {
		switch r.Method {
		case "POST":
			// get api request
			request, err := getApiRequestFromHttpRequest(r)
			if err != nil {
				return http.StatusBadRequest, err
			}

			userResult, err := rpcApi.Do(&api.Request{Method: "get_user", Params: &api.RequestParams{Apikey: request.Params.Apikey}})
			if nil != err {
				return http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized))
			}
			user := userResult.Data.User

			// check against allowed methods
			if !user.IsSuperuser && !rpcApi.IsPublicMethod(request.Method) {
				logger.Warnf("Not a public api method: %v", request.Method)
				return http.StatusMethodNotAllowed, errors.New(http.StatusText(http.StatusMethodNotAllowed))
			}

			// run api request
			results, statusCode, err := execApiRequest(request)
			data = results
			return statusCode, err

		default:
			return http.StatusMethodNotAllowed, errors.New(http.StatusText(http.StatusMethodNotAllowed))
		}
	}()

	if nil != err {
		apiBasicResponse(w, status_code, err)
		return
	}

	apiJSONResponse(w, []byte(data), status_code)
}

func apiWithSessionHandler(w http.ResponseWriter, r *http.Request) {

	var data string

	val, _ := sessionManager.Get(r)
	useremail := val.Values["useremail"].(string)
	if 0 == len(useremail) {
		apiBasicResponse(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	userResult, err := rpcApi.Do(&api.Request{Method: "get_user", Params: &api.RequestParams{Username: useremail}})
	if nil != err {
		apiBasicResponse(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	user := userResult.Data.User

	status_code, err := func() (int, error) {
		switch r.Method {
		case "POST":
			// get api request
			request, err := getApiRequestFromHttpRequest(r)
			if err != nil {
				return http.StatusBadRequest, err
			}

			// WARNING - protect against request hijacking!
			if nil == request.Params {
				request.Params = &api.RequestParams{}
			}
			request.Params.Apikey = user.Apikey

			// check against allowed methods
			if !user.IsSuperuser && !rpcApi.IsPublicMethod(request.Method) {
				logger.Warnf("Not a public api method: %v", request.Method)
				return http.StatusMethodNotAllowed, errors.New(http.StatusText(http.StatusMethodNotAllowed))
			}

			// run api request
			results, statusCode, err := execApiRequest(request)
			data = results
			return statusCode, err

		default:
			return http.StatusMethodNotAllowed, errors.New(http.StatusText(http.StatusMethodNotAllowed))
		}
	}()

	if nil != err {
		apiBasicResponse(w, status_code, err)
		return
	}

	apiJSONResponse(w, []byte(data), status_code)
}
