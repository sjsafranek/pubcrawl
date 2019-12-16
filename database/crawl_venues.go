package database

import (
	"encoding/json"

	"github.com/sjsafranek/pubcrawl/foursquare"
)

type Venue foursquare.Venue

type CrawlVenues struct {
	Next      []string `json:"next"`
	Visited   []*Venue `json:"visited"`
	Unvisited []*Venue `json:"unvisited"`
}

func (self *CrawlVenues) Marshal() (string, error) {
	b, err := json.Marshal(self)
	if nil != err {
		return "", err
	}
	return string(b), err
}

func (self *CrawlVenues) Unmarshal(data string) error {
	return json.Unmarshal([]byte(data), self)
}
