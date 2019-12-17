package api

import (
	"encoding/json"

	"github.com/sjsafranek/pubcrawl/lib/database"
)

type Request struct {
	Method    string  `json:"method,omitempty"`
	Email     string  `json:"email,omitempty"`
	Username  string  `json:"username,omitempty"`
	Password  string  `json:"password,omitempty"`
	Apikey    string  `json:"apikey,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Name      string  `json:"name,omitempty"`
	CrawlId   string  `json:"crawl_id,omitempty"`
	// Type        string           `json:"type,omitempty"`
	// PlaceName   string           `json:"placename,omitempty"`
	// PlaceId     int              `json:"place_id,ompitempty"`
	// Place       *database.Place  `json:"place,omitempty"`
	// Limit       int              `json:"limit,omitempty"`
	// Filter      *database.Filter `json:"filter,ompitempty"`
	Callback string `json:"callback,ompitempty"`
	// PlaceStatus string           `json:"place_status,ompitempty"`
	// Polygon   *geojson.Geometry `json:"polygon,omitempty"`
}

func (self *Request) Unmarshal(data string) error {
	return json.Unmarshal([]byte(data), self)
}

type ResponseData struct {
	Users []*database.User `json:"users,omitempty"`
	User  *database.User   `json:"user,omitempty"`
	// Places   *geojson.FeatureCollection `json:"places,omitempty"`
	// Location *gogeocoder.Location       `json:"location,omitempty"`
	// Place    *database.Place            `json:"place,omitempty"`
	Crawl  *database.Crawl   `json:"crawl,omitempty"`
	Crawls []*database.Crawl `json:"crawls,omitempty"`
}

type Response struct {
	Status   string       `json:"status"`
	Message  string       `json:"message,omitempty"`
	Error    string       `json:"error,omitempty"`
	Data     ResponseData `json:"data,omitempty"`
	Callback string       `json:"callback,omitempty"`
	// CaptchaID string       `json:"captchaID,omitempty"`
}

func (self *Response) Marshal() (string, error) {
	b, err := json.Marshal(self)
	if nil != err {
		return "", err
	}
	return string(b), err
}

func (self *Response) SetError(err error) {
	self.Status = "error"
	self.Error = err.Error()
}
