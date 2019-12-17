package database

import (
	// "database/sql"
	"encoding/json"
	"time"

	"github.com/sjsafranek/pubcrawl/lib/foursquare"
)

type Crawl struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	Owner           string      `json:"owner"`
	MaxVotesPerUser int         `json:"max_votes_per_user"`
	IsDeleted       bool        `json:"is_deleted"`
	CreatedAt       time.Time   `json:"created_at,string"`
	UpdatedAt       time.Time   `json:"updated_at,string"`
	Venues          CrawlVenues `json:"venues"`
	db              *Database   `json:"-"`
}

func (self *Crawl) Marshal() (string, error) {
	b, err := json.Marshal(self)
	if nil != err {
		return "", err
	}
	return string(b), err
}

func (self *Crawl) Unmarshal(data string) error {
	return json.Unmarshal([]byte(data), self)
}

// Delete deletes user
func (self *Crawl) Delete() error {
	self.IsDeleted = true
	return self.Update()
}

// Update updates user data in database
func (self *Crawl) Update() error {
	return self.db.Insert(`
		UPDATE crawls
			SET
				is_deleted=$2
			WHERE id=$1;`, self.ID, self.IsDeleted)
}

func (self *Crawl) AddVenues(venues []foursquare.Venue) error {
	for _, venue := range venues {
		data, err := venue.Marshal()
		if nil != err {
			return err
		}

		err = self.db.Insert(`
			INSERT INTO venues(id, crawl_id, data) VALUES($1, $2, $3);
		`, venue.ID, self.ID, data)
	}

	return nil
}

func (self *Crawl) GetVenues() (*CrawlVenues, error) {

	if 0 != len(self.Venues.Unvisited) || 0 != len(self.Venues.Visited) {
		return &self.Venues, nil
	}

	var venues CrawlVenues

	query := `
		SELECT
			venues_json
		FROM venues_view
		WHERE
			crawl_id = $1;
	`

	var temp string
	err := self.db.QueryRow(query, self.ID).Scan(&temp)
	if nil != err {
		return &venues, err
	}
	venues.Unmarshal(temp)

	self.Venues = venues

	return &venues, nil
}
