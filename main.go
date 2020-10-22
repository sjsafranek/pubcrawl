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
	"github.com/sjsafranek/pubcrawl/lib/api"
	"github.com/sjsafranek/pubcrawl/lib/config"
)

const (
	PROJECT                   string = "PubCrawl"
	VERSION                   string = "0.0.1"
	DEFAULT_HTTP_PORT         int    = 8080
	DEFAULT_HTTP_HOST string = "localhost"
	DEFAULT_HTTP_PROTOCOL string = "http"
	DEFAULT_DATABASE_ENGINE   string = "postgres"
	DEFAULT_DATABASE_DATABASE string = "crawldb"
	DEFAULT_DATABASE_PASSWORD string = "dev"
	DEFAULT_DATABASE_USERNAME string = "crawluser"
	DEFAULT_DATABASE_HOST     string = "localhost"
	DEFAULT_DATABASE_PORT     int64  = 5432
)

var (
	FACEBOOK_CLIENT_ID     string = os.Getenv("FACEBOOK_CLIENT_ID")
	FACEBOOK_CLIENT_SECRET string = os.Getenv("FACEBOOK_CLIENT_SECRET")
	GOOGLE_CLIENT_ID       string = os.Getenv("GOOGLE_CLIENT_ID")
	GOOGLE_CLIENT_SECRET   string = os.Getenv("GOOGLE_CLIENT_SECRET")
	FOURSQUARE_CLIENT_ID     string = os.Getenv("FOURSQUARE_CLIENT_ID")
	FOURSQUARE_CLIENT_SECRET string = os.Getenv("FOURSQUARE_CLIENT_SECRET")
	DATABASE_ENGINE          string = DEFAULT_DATABASE_ENGINE
	DATABASE_DATABASE        string = DEFAULT_DATABASE_DATABASE
	DATABASE_PASSWORD        string = DEFAULT_DATABASE_PASSWORD
	DATABASE_USERNAME        string = DEFAULT_DATABASE_USERNAME
	DATABASE_HOST            string = DEFAULT_DATABASE_HOST
	DATABASE_PORT            int64  = DEFAULT_DATABASE_PORT
	API_REQUEST              string = ""
	rpcApi                   *api.Api
	conf                     *config.Config
)

// main creates and starts a Server listening.
func main() {
	// read credentials from environment variables if available
	conf = &config.Config{
		Api: config.Api{
			PublicMethods: []string{
				"create_crawl",
				"get_crawl",
				"get_crawls",
				"delete_crawl",
				"get_venues",
				"up_vote",
				"down_vote",
			},
		},
		Server: config.Server{
			HttpPort: DEFAULT_HTTP_PORT,
			HttpHost: DEFAULT_HTTP_HOST,
			HttpProtocol: DEFAULT_HTTP_PROTOCOL,
		},
		Foursquare: config.Foursquare{
			ClientID:     FOURSQUARE_CLIENT_ID,
			ClientSecret: FOURSQUARE_CLIENT_SECRET,
		},
		OAuth2: config.OAuth2{
			Facebook: config.SocialOAuth2{
				ClientID:     FACEBOOK_CLIENT_ID,
				ClientSecret: FACEBOOK_CLIENT_SECRET,
			},
			Google: config.SocialOAuth2{
				ClientID:     GOOGLE_CLIENT_ID,
				ClientSecret: GOOGLE_CLIENT_SECRET,
			},
		},
		Database: config.Database{
			DatabaseEngine: DATABASE_ENGINE,
			DatabaseHost:   DEFAULT_DATABASE_HOST,
			DatabaseName:   DEFAULT_DATABASE_DATABASE,
			DatabasePass:   DEFAULT_DATABASE_PASSWORD,
			DatabaseUser:   DEFAULT_DATABASE_USERNAME,
			DatabasePort:   DEFAULT_DATABASE_PORT,
		},
	}

	// allow consumer credential flags to override config fields
	var printVersion bool
	flag.StringVar(&conf.OAuth2.Facebook.ClientID, "facebook-client-id", FACEBOOK_CLIENT_ID, "Facebook Client ID")
	flag.StringVar(&conf.OAuth2.Facebook.ClientSecret, "facebook-client-secret", FACEBOOK_CLIENT_SECRET, "Facebook Client Secret")
	flag.StringVar(&conf.OAuth2.Google.ClientID, "gmail-client-id", GOOGLE_CLIENT_ID, "Google Client ID")
	flag.StringVar(&conf.OAuth2.Google.ClientSecret, "gmail-client-secret", GOOGLE_CLIENT_SECRET, "Google Client Secret")
	flag.StringVar(&conf.Foursquare.ClientID, "foursquare-client-id", FOURSQUARE_CLIENT_ID, "Foursquare Client ID")
	flag.StringVar(&conf.Foursquare.ClientSecret, "foursquare-client-secret", FOURSQUARE_CLIENT_SECRET, "Foursquare Client Secret")
	flag.BoolVar(&printVersion, "V", false, "Print version and exit")
	flag.IntVar(&conf.Server.HttpPort, "httpport", conf.Server.HttpPort, "Server port")
	flag.StringVar(&conf.Database.DatabaseHost, "dbhost", DEFAULT_DATABASE_HOST, "database host")
	flag.StringVar(&conf.Database.DatabaseName, "dbname", DEFAULT_DATABASE_DATABASE, "database name")
	flag.StringVar(&conf.Database.DatabasePass, "dbpass", DEFAULT_DATABASE_PASSWORD, "database password")
	flag.StringVar(&conf.Database.DatabaseUser, "dbuser", DEFAULT_DATABASE_USERNAME, "database username")
	flag.Int64Var(&conf.Database.DatabasePort, "dbport", DEFAULT_DATABASE_PORT, "Database port")
	flag.StringVar(&API_REQUEST, "query", "", "Api query to execute")
	flag.Parse()

	rpcApi = api.New(conf)

	if "" != API_REQUEST {
		response, err := rpcApi.DoJSON(API_REQUEST)
		if nil != err {
			fmt.Println(err)
		}

		results, err := response.Marshal()
		if nil != err {
			panic(err)
		}
		fmt.Println(results)
		return
	}

	logger.Debug("GOOS: ", runtime.GOOS)
	logger.Debug("CPUS: ", runtime.NumCPU())
	logger.Debug("PID: ", os.Getpid())
	logger.Debug("Go Version: ", runtime.Version())
	logger.Debug("Go Arch: ", runtime.GOARCH)
	logger.Debug("Go Compiler: ", runtime.Compiler)
	logger.Debug("NumGoroutine: ", runtime.NumGoroutine())

	resp, err := rpcApi.DoJSON(`{"method":"get_database_version"}`)
	if nil != err {
		panic(err)
	}
	logger.Debugf("Database version: %v", resp.Message)
	//
	address := fmt.Sprintf(":%v", conf.Server.HttpPort)
	logger.Infof("The magic happens on %s\n", address)
	err = http.ListenAndServe(address, New(conf))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
