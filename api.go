package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/sjsafranek/logger"
)

// # /api/v1/crawl?longitude=-123.088000&latitude=44.046174

const (
	CategoryBrewery string = "brewery"
	CategoryBar     string = "bar"
)

var categoryName = map[string]string{
	CategoryBrewery: "50327c8591d4c4b30a586d5d",
	CategoryBar:     "4bf58dd8d48988d116941735",
}

func CategoryCode(name string) string {
	return categoryName[name]
}

var searchCategeories = []string{"50327c8591d4c4b30a586d5d", "4bf58dd8d48988d116941735"}

// profileHandler shows protected user content.
func newPubCrawlHandler(w http.ResponseWriter, req *http.Request) {

	response := Response{}

	val, _ := sessionManager.Get(req)
	username := val.Values["username"].(string)
	userid := val.Values["userid"].(string)
	usertype := val.Values["usertype"].(string)
	useremail := val.Values["useremail"].(string)

	// if req.Method == "POST" {
	// }

	statusCode, err := func() (int, error) {

		longitudeString, ok := req.URL.Query()["longitude"]
		if !ok {
			return http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest))
		}

		latitudeString, ok := req.URL.Query()["latitude"]
		if !ok {
			return http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest))
		}

		longitude, err := strconv.ParseFloat(longitudeString[0], 64)
		latitude, err := strconv.ParseFloat(latitudeString[0], 64)

		crawlName, _ := req.URL.Query()["name"]

		venues, err := apiClient.SearchVenues(longitude, latitude, searchCategeories)
		if nil != err {
			return http.StatusInternalServerError, err
		}

		user, err := db.CreateUserIfNotExists(useremail, useremail)
		if nil != err {
			return http.StatusInternalServerError, err
		}
		user.CreateSocialAccountIfNotExists(userid, username, usertype)

		crawl, err := user.CreateCrawl(crawlName[0])
		if nil != err {
			return http.StatusInternalServerError, err

		}

		err = crawl.AddVenues(venues)
		if nil != err {
			return http.StatusInternalServerError, err
		}
		// crawl.GetVenues()
		crawl, err = user.GetCrawl(crawl.ID)
		if nil != err {
			return http.StatusInternalServerError, err
		}

		response.Data.PubCrawl = crawl

		return http.StatusOK, nil

	}()

	if nil != err {
		logger.Error(err)
		apiBasicResponse(w, statusCode, err)
		return
	}

	results, _ := response.Marshal()

	apiJSONResponse(w, []byte(results), http.StatusOK)
}

func getPubCrawlHandler(w http.ResponseWriter, req *http.Request) {

	response := Response{}

	val, _ := sessionManager.Get(req)
	// username := val.Values["username"].(string)
	// userid := val.Values["userid"].(string)
	// usertype := val.Values["usertype"].(string)
	useremail := val.Values["useremail"].(string)
	logger.Info(useremail, response)
}
