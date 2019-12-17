package foursquare

import (
	"fmt"
	// "io/ioutil"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"
	// "github.com/sjsafranek/logger"
)

var client = &http.Client{
	Timeout: time.Second * 10,
	Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	},
}

const (
	FoursquareApiUrl     = "https://api.foursquare.com/v2"
	FoursquareApiVersion = "20191212"
)

type Client struct {
	ClientID     string
	ClientSecret string
}

func (self *Client) newRequest(url string) (*http.Request, error) {
	// create new http request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return req, err
	}

	// set basic headers
	req.Header.Add("Accept", "application/json")

	// automatically include required parmeters
	q := req.URL.Query()
	q.Add("v", FoursquareApiVersion)
	q.Add("client_id", self.ClientID)
	q.Add("client_secret", self.ClientSecret)
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func (self *Client) SearchVenues(longitude, latitude float64, categories []string) ([]Venue, error) {
	url := fmt.Sprintf("%v/venues/search", FoursquareApiUrl)

	req, err := self.newRequest(url)
	if err != nil {
		return []Venue{}, err
	}

	q := req.URL.Query()
	q.Add("ll", fmt.Sprintf("%v,%v", latitude, longitude)) // latlng formatting
	q.Add("limit", "50")                                   // 50 max
	q.Add("radius", "2000")                                // meters
	q.Add("categoryId", strings.Join(categories, ","))     // list of category ids
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return []Venue{}, err
	}

	defer resp.Body.Close()
	var response ApiResponse
	err = json.NewDecoder(resp.Body).Decode(&response)

	return response.Response.Venues, nil
}

func (self *Client) Categories() ([]Category, error) {
	url := fmt.Sprintf("%v/venues/categories", FoursquareApiUrl)

	req, err := self.newRequest(url)
	if err != nil {
		return []Category{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return []Category{}, err
	}

	defer resp.Body.Close()
	var response ApiResponse
	err = json.NewDecoder(resp.Body).Decode(&response)

	return response.Response.Categories, nil
}
