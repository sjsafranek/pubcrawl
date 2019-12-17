package config

import (
	"fmt"
)

// Config configures the app
type Config struct {
	Facebook   Facebook
	Foursquare Foursquare
	Server     Server
	Database   Database
}

type Server struct {
	HttpPort int
}

type Facebook struct {
	ClientID     string
	ClientSecret string
}

type Foursquare struct {
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
