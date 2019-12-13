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
	Venues     []Venue    `json:"venues"`
	Categories []Category `json:"categories"`
}

type Venue struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Location   Location   `json:"location"`
	Categories []Category `json:"categories"`
}

type Category struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Categories []Category `json:"categories"` // handles sub categories
	// 		"pluralName": "Beer Stores",
	// 		"shortName": "Beer Store",
	// 		"icon": {
	// 			"prefix": "https:\/\/ss3.4sqi.net\/img\/categories_v2\/nightlife\/beergarden_",
	// 			"suffix": ".png"
	// 		},
	// 		"primary": true

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
