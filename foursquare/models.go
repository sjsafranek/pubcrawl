package foursquare

import (
	"strings"
)

type ApiResponse struct {
	Meta     Meta     `json:"meta"`
	Response Response `json:"response"`
}

type Meta struct {
	Code int `json:"code"`
	// RequestId string `json:"requestId"`
}

type Response struct {
	// Confident bool    `json:"confident"`
	Venues []Venue `json:"venues"`
}

type Venue struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Location Location `json:"location"`
	// "categories": [{
	// 		"id": "5370f356bcbc57f1066c94c2",
	// 		"name": "Beer Store",
	// 		"pluralName": "Beer Stores",
	// 		"shortName": "Beer Store",
	// 		"icon": {
	// 			"prefix": "https:\/\/ss3.4sqi.net\/img\/categories_v2\/nightlife\/beergarden_",
	// 			"suffix": ".png"
	// 		},
	// 		"primary": true
	// 	}]
}

type Location struct {
	Address          string   `json:"address"`
	CrossStreet      string   `json:"crossStreet"`
	Lat              float64  `json:"lat"`
	Lng              float64  `json:"lng"`
	PostalCode       string   `json:"postalCode"`
	CC               string   `json:"cc"`
	City             string   `json:"city"`
	State            string   `json:"state"`
	Country          string   `json:"country"`
	FormattedAddress []string `json:"formattedAddress"`
}

func (self *Location) ToString() string {
	return strings.Join(self.FormattedAddress, ",")
}
