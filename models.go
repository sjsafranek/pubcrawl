package main

//
// import (
// 	"encoding/json"
// 	// "errors"
//
// 	// "github.com/google/uuid"
// 	"github.com/sjsafranek/pubcrawl/lib/database"
// )
//
// //
// // func NewPubCrawl(owner_id string, venues []foursquare.Venue) *PubCrawl {
// // 	pc := PubCrawl{
// // 		ID:      uuid.New().String(),
// // 		Venues:  venues,
// // 		Visited: make(map[string]bool),
// // 	}
// // 	for _, venue := range pc.Venues {
// // 		pc.Visited[venue.ID] = false
// // 	}
// // 	return &pc
// // }
// //
// // type User struct {
// // 	UserID   string `json:"userid"`
// // 	UserName string `json:"username"`
// // 	UserType string `json:"usertype"`
// // }
// //
// // type UserVote struct {
// // 	UpVote   int `json:"up_vote"`
// // 	DownVote int `json:"down_vote"`
// // }
// //
// // type UserRule struct {
// // 	UpVoteWeight   int `json:"up_vote_weight"`
// // 	DownVoteWeight int `json:"down_vote_weight"`
// // }
// //
// // type PubCrawl struct {
// // 	ID        string                          `json:"id"`
// // 	OwnerID   string                          `json:"owner_id"`
// // 	Venues    []foursquare.Venue              `json:"venues"`
// // 	Visited   map[string]bool                 `json:"visited"`
// // 	UserVotes map[string]map[string]*UserVote `json:"user_votes"`
// // 	UserRules map[string]UserRule             `json:"user_rules"`
// // 	// Users     map[string]User                 `json:"users"`
// // }
//
// // func (self *PubCrawl) Has(venue_id string) bool {
// // 	_, err := self.Get(venue_id)
// // 	return nil == err
// // }
// //
// // func (self *PubCrawl) Get(venue_id string) (*foursquare.Venue, error) {
// // 	for _, venue := range self.Venues {
// // 		if venue_id == venue.ID {
// // 			return &venue, nil
// // 		}
// // 	}
// // 	return &foursquare.Venue{}, errors.New("Not Found")
// // }
//
// //
// // func (self *PubCrawl) UpVote(user_id, venue_id string) error {
// // 	if self.Has(venue_id) {
// // 		if nil == self.UserVotes {
// // 			self.UserVotes = make(map[string]map[string]*UserVote)
// // 		}
// // 		if nil == self.UserVotes[venue_id] {
// // 			self.UserVotes[venue_id] = make(map[string]*UserVote)
// // 		}
// // 		_, ok := self.UserVotes[venue_id][user_id]
// // 		if !ok {
// // 			self.UserVotes[venue_id][user_id] = &UserVote{}
// // 		}
// // 		self.UserVotes[venue_id][user_id].UpVote++
// // 	}
// // 	return errors.New("Not Found")
// // }
// //
// // func (self *PubCrawl) DownVote(user_id, venue_id string) error {
// // 	if self.Has(venue_id) {
// // 		if nil == self.UserVotes {
// // 			self.UserVotes = make(map[string]map[string]*UserVote)
// // 		}
// // 		if nil == self.UserVotes[venue_id] {
// // 			self.UserVotes[venue_id] = make(map[string]*UserVote)
// // 		}
// // 		_, ok := self.UserVotes[venue_id][user_id]
// // 		if !ok {
// // 			self.UserVotes[venue_id][user_id] = &UserVote{}
// // 		}
// // 		self.UserVotes[venue_id][user_id].DownVote++
// // 	}
// // 	return errors.New("Not Found")
// // }
// //
// // func (self *PubCrawl) Visit(venue_id string) error {
// // 	if self.Has(venue_id) {
// // 		if nil == self.Visited {
// // 			self.Visited = make(map[string]bool)
// // 		}
// // 	}
// // 	return errors.New("Not Found")
// // }
//
// type ResponseData struct {
// 	PubCrawl *database.Crawl `json:"pub_crawl"`
// }
//
// type Response struct {
// 	Status  string       `json:"status"`
// 	Message string       `json:"message,omitempty"`
// 	Error   string       `json:"error,omitempty"`
// 	Data    ResponseData `json:"data,omitempty"`
// }
//
// func (self *Response) Marshal() (string, error) {
// 	b, err := json.Marshal(self)
// 	if nil != err {
// 		return "", err
// 	}
// 	return string(b), err
// }
//
// func (self *Response) SetError(err error) {
// 	self.Status = "error"
// 	self.Error = err.Error()
// }
