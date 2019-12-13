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

func (self *Client) SearchVenues(longitude, latitude float64, categories []string) ([]Venue, error) {
	url := fmt.Sprintf("%v/venues/search", FoursquareApiUrl)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []Venue{}, err
	}

	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	q.Add("ll", fmt.Sprintf("%v,%v", latitude, longitude))
	// q.Add("limit", "50")
	q.Add("limit", "2")
	q.Add("radius", "2000") // meters
	q.Add("categoryId", strings.Join(categories, ","))
	q.Add("v", FoursquareApiVersion)
	q.Add("client_id", self.ClientID)
	q.Add("client_secret", self.ClientSecret)
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
