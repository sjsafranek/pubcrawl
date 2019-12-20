package database

import (
	"time"
	"encoding/json"

	"github.com/sjsafranek/pubcrawl/lib/foursquare"
)

type Venue foursquare.Venue

type CrawlVenue struct {
	VenueId   string    `json:"venue_id"`
	Visited   bool      `json:"visited"`
	UpVotes   int       `json:"up_votes"`
	DownVotes int       `json:"down_votes"`
	Venue    	*Venue  `json:"venue"`
	CreatedAt time.Time `json:"created_at,string"`
	UpdatedAt time.Time `json:"updated_at,string"`
}

type CrawlVenues []*CrawlVenue

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
