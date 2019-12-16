package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	// "github.com/sjsafranek/lemur"
	"github.com/sjsafranek/logger"
	"github.com/sjsafranek/pubcrawl/database"
	"github.com/sjsafranek/pubcrawl/foursquare"
)

const (
	PROJECT                   string = "PubCrawl"
	VERSION                   string = "0.0.1"
	DEFAULT_HTTP_PORT         int    = 8080
	DEFAULT_DATABASE_ENGINE   string = "postgres"
	DEFAULT_DATABASE_DATABASE string = "crawldb"
	DEFAULT_DATABASE_PASSWORD string = "dev"
	DEFAULT_DATABASE_USERNAME string = "crawluser"
	DEFAULT_DATABASE_HOST     string = "localhost"
	DEFAULT_DATABASE_PORT     int64  = 5432
)

var (
	FACEBOOK_CLIENT_ID       string = os.Getenv("FACEBOOK_CLIENT_ID")
	FACEBOOK_CLIENT_SECRET   string = os.Getenv("FACEBOOK_CLIENT_SECRET")
	FOURSQUARE_CLIENT_ID     string = os.Getenv("FOURSQUARE_CLIENT_ID")
	FOURSQUARE_CLIENT_SECRET string = os.Getenv("FOURSQUARE_CLIENT_SECRET")
	HTTP_PORT                int    = DEFAULT_HTTP_PORT
	DATABASE_ENGINE          string = DEFAULT_DATABASE_ENGINE
	DATABASE_DATABASE        string = DEFAULT_DATABASE_DATABASE
	DATABASE_PASSWORD        string = DEFAULT_DATABASE_PASSWORD
	DATABASE_USERNAME        string = DEFAULT_DATABASE_USERNAME
	DATABASE_HOST            string = DEFAULT_DATABASE_HOST
	DATABASE_PORT            int64  = DEFAULT_DATABASE_PORT

	apiClient *foursquare.Client
	db        *database.Database
)

// main creates and starts a Server listening.
func main() {
	const address = "localhost:8080"

	// read credentials from environment variables if available
	config := &Config{
		FacebookClientID:       FACEBOOK_CLIENT_ID,
		FacebookClientSecret:   FACEBOOK_CLIENT_SECRET,
		FoursquareClientID:     FOURSQUARE_CLIENT_ID,
		FoursquareClientSecret: FOURSQUARE_CLIENT_SECRET,
	}

	// allow consumer credential flags to override config fields
	var printVersion bool
	flag.StringVar(&config.FacebookClientID, "facebook-client-id", FACEBOOK_CLIENT_ID, "Facebook Client ID")
	flag.StringVar(&config.FacebookClientSecret, "facebook-client-secret", FACEBOOK_CLIENT_SECRET, "Facebook Client Secret")
	flag.StringVar(&config.FoursquareClientID, "foursquare-client-id", FOURSQUARE_CLIENT_ID, "Foursquare Client ID")
	flag.StringVar(&config.FoursquareClientSecret, "foursquare-client-secret", FOURSQUARE_CLIENT_SECRET, "Foursquare Client Secret")
	flag.BoolVar(&printVersion, "V", false, "Print version and exit")
	flag.IntVar(&HTTP_PORT, "httpport", DEFAULT_HTTP_PORT, "Server port")
	flag.StringVar(&DATABASE_HOST, "dbhost", DEFAULT_DATABASE_HOST, "database host")
	flag.StringVar(&DATABASE_DATABASE, "dbname", DEFAULT_DATABASE_DATABASE, "database name")
	flag.StringVar(&DATABASE_PASSWORD, "dbpass", DEFAULT_DATABASE_PASSWORD, "database password")
	flag.StringVar(&DATABASE_USERNAME, "dbuser", DEFAULT_DATABASE_USERNAME, "database username")
	flag.Int64Var(&DATABASE_PORT, "dbport", DEFAULT_DATABASE_PORT, "Database port")
	flag.Parse()

	if config.FacebookClientID == "" {
		log.Fatal("Missing Facebook Client ID")
	}
	if config.FacebookClientSecret == "" {
		log.Fatal("Missing Facebook Client Secret")
	}

	logger.Debug("GOOS: ", runtime.GOOS)
	logger.Debug("CPUS: ", runtime.NumCPU())
	logger.Debug("PID: ", os.Getpid())
	logger.Debug("Go Version: ", runtime.Version())
	logger.Debug("Go Arch: ", runtime.GOARCH)
	logger.Debug("Go Compiler: ", runtime.Compiler)
	logger.Debug("NumGoroutine: ", runtime.NumGoroutine())

	dbConnectionString := fmt.Sprintf("%v://%v:%v@%v:%v/%v?sslmode=disable", DATABASE_ENGINE, DATABASE_USERNAME, DATABASE_PASSWORD, DATABASE_HOST, DATABASE_PORT, DATABASE_DATABASE)
	db = database.New(dbConnectionString)

	apiClient = foursquare.New(config.FoursquareClientID, config.FoursquareClientSecret)

	logger.Infof("The magic happens on %s\n", address)
	err := http.ListenAndServe(address, New(config))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
