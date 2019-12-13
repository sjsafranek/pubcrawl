package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"

	// "github.com/sjsafranek/lemur"
	"github.com/sjsafranek/logger"
	"github.com/sjsafranek/pubcrawl/foursquare"
)

var (
	FACEBOOK_CLIENT_ID       string = os.Getenv("FACEBOOK_CLIENT_ID")
	FACEBOOK_CLIENT_SECRET   string = os.Getenv("FACEBOOK_CLIENT_SECRET")
	FOURSQUARE_CLIENT_ID     string = os.Getenv("FOURSQUARE_CLIENT_ID")
	FOURSQUARE_CLIENT_SECRET string = os.Getenv("FOURSQUARE_CLIENT_SECRET")

	apiClient *foursquare.Client
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
	flag.StringVar(&config.FacebookClientID, "facebook-client-id", FACEBOOK_CLIENT_ID, "Facebook Client ID")
	flag.StringVar(&config.FacebookClientSecret, "facebook-client-secret", FACEBOOK_CLIENT_SECRET, "Facebook Client Secret")
	flag.StringVar(&config.FoursquareClientID, "foursquare-client-id", FOURSQUARE_CLIENT_ID, "Foursquare Client ID")
	flag.StringVar(&config.FoursquareClientSecret, "foursquare-client-secret", FOURSQUARE_CLIENT_SECRET, "Foursquare Client Secret")
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

	apiClient = foursquare.New(config.FoursquareClientID, config.FoursquareClientSecret)

	logger.Infof("The magic happens on %s\n", address)
	err := http.ListenAndServe(address, New(config))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
