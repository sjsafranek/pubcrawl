package database

import (
	"encoding/json"

	"github.com/sjsafranek/pubcrawl/lib/foursquare"
)

type Venue foursquare.Venue

type VenueVotes struct {
	VenueId   string `json:"venue_id"`
	UpVotes   int    `json:"up_votes"`
	DownVotes int    `json:"down_votes"`
	Votes     int    `json:"votes"`
}

type CrawlVenues struct {
	Next      []string      `json:"next"`
	Visited   []*Venue      `json:"visited"`
	Unvisited []*Venue      `json:"unvisited"`
	Votes     []*VenueVotes `json:"votes"`
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
