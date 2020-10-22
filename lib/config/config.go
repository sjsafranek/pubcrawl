package config

import (
	"fmt"
)

// Config configures the app
type Config struct {
	Foursquare Foursquare
	OAuth2   OAuth2
	Server     Server
	Database   Database
	Api        Api
}

type Api struct {
	PublicMethods []string
}

func (self *Api) IsPublicMethod(method string) bool {
	for _, publicMethod := range self.PublicMethods {
		if publicMethod == method {
			return true
		}
	}
	return false
}

type Server struct {
	HttpPort int
}

// type Facebook struct {
// 	ClientID     string
// 	ClientSecret string
// }
//
type Foursquare struct {
	ClientID     string
	ClientSecret string
}


type OAuth2 struct {
	Facebook SocialOAuth2
	Google SocialOAuth2
}

func (self *OAuth2) HasFacebook() bool {
	return "" != self.Facebook.ClientID && "" != self.Facebook.ClientSecret
}

func (self *OAuth2) HasGoogle() bool {
	return "" != self.Google.ClientID && "" != self.Google.ClientSecret
}


type SocialOAuth2 struct {
	ClientID     string
	ClientSecret string
}


type Database struct {
	DatabaseEngine string
	DatabaseName   string
	DatabasePass   string
	DatabaseUser   string
	DatabaseHost   string
	DatabasePort   int64
}

func (self *Database) GetDatabaseConnection() string {
	return fmt.Sprintf("%v://%v:%v@%v:%v/%v?sslmode=disable",
		self.DatabaseEngine,
		self.DatabaseUser,
		self.DatabasePass,
		self.DatabaseHost,
		self.DatabasePort,
		self.DatabaseName)
}
