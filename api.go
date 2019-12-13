package main

import (
	"github.com/sjsafranek/logger"
	// "html/template"
	"net/http"
)

// profileHandler shows protected user content.
func createPubCrawl(w http.ResponseWriter, req *http.Request) {
	// val, _ := sessionStore.Get(req, sessionName)
	val, _ := sessionManager.Get(req)
	facebookName := val.Values["username"]

	args := make(map[string]interface{})
	args["username"] = facebookName

	// # /api/v1/crawl?longitude=-123.088000&latitude=44.046174
	// CATEGORIES = {
	//     'Brewery': '50327c8591d4c4b30a586d5d',
	//     'Bar': '4bf58dd8d48988d116941735'
	// }

	venues, err := apiClient.SearchVenues(-123.088000, 44.046174, []string{"50327c8591d4c4b30a586d5d", "4bf58dd8d48988d116941735"})
	logger.Info(venues, err)

	categories, err := apiClient.Categories()
	logger.Info(categories, err)
}
