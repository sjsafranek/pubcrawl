package api

import (
	"encoding/json"

	"github.com/sjsafranek/pubcrawl/lib/database"
)

/*
Format for JSON RPC

https://en.wikipedia.org/wiki/JSON-RPC
*/

const VERSION string = "0.0.1"

type Request struct {
	Method  string `json:"method"`
	Version string `json:"version"`
	Params  Params `json:"params"`
	Id      string `json:"id,ompitempty"`
}

type Params struct {
	Email     string  `json:"email,omitempty"`
	Username  string  `json:"username,omitempty"`
	Password  string  `json:"password,omitempty"`
	Apikey    string  `json:"apikey,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Name      string  `json:"name,omitempty"`
	CrawlId   string  `json:"crawl_id,omitempty"`
	VenueId   string  `json:"venue_id,omitempty"`
}

func (self *Request) Unmarshal(data string) error {
	return json.Unmarshal([]byte(data), self)
}

type ResponseData struct {
	Users  []*database.User  `json:"users,omitempty"`
	User   *database.User    `json:"user,omitempty"`
	Crawl  *database.Crawl   `json:"crawl,omitempty"`
	Crawls []*database.Crawl `json:"crawls,omitempty"`
}

type Response struct {
	Status  string       `json:"status"`
	Version string       `json:"version,omitempty"`
	Message string       `json:"message,omitempty"`
	Error   string       `json:"error,omitempty"`
	Data    ResponseData `json:"data,omitempty"`
	Id      string       `json:"id,omitempty"`
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
