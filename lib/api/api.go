package api

import (
	// "encoding/json"
	"errors"
	// "fmt"
	// "strings"
	"time"

	"github.com/karlseguin/ccache"
	"github.com/sjsafranek/logger"
	"github.com/sjsafranek/pubcrawl/lib/config"
	"github.com/sjsafranek/pubcrawl/lib/database"
	"github.com/sjsafranek/pubcrawl/lib/foursquare"
)

func New(conf *config.Config) *Api {
	dbConnStr := conf.Database.GetDatabaseConnection()
	return &Api{
		config:     conf,
		db:         database.New(dbConnStr),
		cache:      ccache.Layered(ccache.Configure()),
		foursquare: foursquare.New(conf.Foursquare.ClientID, conf.Foursquare.ClientSecret),
	}
}

type Api struct {
	config     *config.Config
	db         *database.Database
	cache      *ccache.LayeredCache
	foursquare *foursquare.Client
}

func (self *Api) IsPublicMethod(method string) bool {
	return self.config.Api.IsPublicMethod(method)
}

func (self *Api) GetDatabase() *database.Database {
	return self.db
}

func (self *Api) fetchUser(request *Request, clbk func(*database.User) error) error {
	var user *database.User
	var err error
	if "" != request.Params.Apikey {
		user, err = self.getUserByApikey(request.Params.Apikey)
	} else if "" != request.Params.Username {
		user, err = self.getUserByUsername(request.Params.Username)
	} else {
		err = errors.New("Missing parameters")
	}
	if nil != err {
		return err
	}
	return clbk(user)
}

// CreateUser
func (self *Api) createUser(email, username, password string) (*database.User, error) {
	user, err := self.db.CreateUser(email, username)
	if nil == err {
		// cache apikey user pair
		err = user.SetPassword(password)
		if nil == err {
			self.cache.Set("user", user.Apikey, user, 5*time.Minute)
		}
	}
	return user, err
}

// GetUserByUserName
func (self *Api) getUserByUsername(username string) (*database.User, error) {
	return self.db.GetUserByUsername(username)
}

// GetUserByApikey fetches user via apikey. This method uses an inmemory LRU cache to
// decrease the number of database transactions.
func (self *Api) getUserByApikey(apikey string) (*database.User, error) {
	// check cache for apikey user pair
	item := self.cache.Get("user", apikey)
	if nil != item {
		return item.Value().(*database.User), nil
	}

	user, err := self.db.GetUserByApikey(apikey)
	if nil == err {
		// cache apikey user pair
		self.cache.Set("user", apikey, user, 5*time.Minute)
	}
	return user, err
}

func (self *Api) DoJSON(jdata string) (*Response, error) {
	var request Request
	// err := json.Unmarshal([]byte(jdata), &request)
	err := request.Unmarshal(jdata)
	if nil != err {
		response := &Response{Status: "err"}
		response.SetError(err)
		return response, err
	}
	return self.Do(&request)
}

func (self *Api) Do(request *Request) (*Response, error) {
	var response Response

	// TODO HANDLE API VERSIONS
	response.Version = request.Version
	if "" == request.Version {
		response.Version = VERSION
		request.Version = VERSION
	}

	response.Status = "ok"
	response.Id = request.Id

	err := func() error {
		switch request.Method {

		case "get_database_version":
			// {"method":"get_database_version"}
			version, err := self.db.GetVersion()
			if nil != err {
				return err
			}
			response.Message = version
			return nil

		case "ping":
			// {"method":"ping"}
			response.Message = "pong"
			return nil

		case "create_user":
			// {"method":"create_user","params":{"username": "admin_user" "email":"admin@email.com","password":"1234"}+
			if "" == request.Params.Username {
				return errors.New("missing parameters")
			}

			user, err := self.createUser(request.Params.Email, request.Params.Username, request.Params.Password)
			if nil != err {
				return err
			}

			response.Data.User = user
			return nil

		case "get_users":
			// {"method":"get_users"}
			users, err := self.db.GetUsers()
			if nil != err {
				return err
			}
			response.Data.Users = users
			return nil

		case "get_user":
			// {"method":"get_user","params":{"username":"admin_user"}}
			// {"method":"get_user","params":{"apikey":"<apikey>"}}
			return self.fetchUser(request, func(user *database.User) error {
				response.Data.User = user
				return nil
			})

		case "delete_user":
			// {"method":"delete_user","params":{"username":"admin_user"}}
			// {"method":"delete_user","params":{"apikey":"<apikey>"}}
			return self.fetchUser(request, func(user *database.User) error {
				self.cache.Delete("user", user.Apikey)
				return user.Delete()
			})

		case "activate_user":
			// {"method":"activate_user","params":{"username":"admin_user"}}
			// {"method":"activate_user","params":{"apikey":"<apikey>"}}
			return self.fetchUser(request, func(user *database.User) error {
				self.cache.Delete("user", user.Apikey)
				return user.Activate()
			})

		case "deactivate_user":
			// {"method":"deactivate_user","params":{"username":"admin_user"}}
			// {"method":"deactivate_user","params":{"apikey":"<apikey>"}}
			return self.fetchUser(request, func(user *database.User) error {
				self.cache.Delete("user", user.Apikey)
				return user.Deactivate()
			})

		case "set_password":
			// {"method":"set_password","params":{"username":"admin_user","password":"1234"}}
			// {"method":"set_password","params":{"apikey":"<apikey>","password":"1234"}}
			return self.fetchUser(request, func(user *database.User) error {
				return user.SetPassword(request.Params.Password)
			})

		case "create_crawl":
			// {"method":"create_crawl","params":{"username":"sjsafranek@gmail.com","name":"12 bars","longitude":-123.088000,"latitude":44.046174}}
			logger.Info(request)
			return self.fetchUser(request, func(user *database.User) error {
				venues, err := self.foursquare.SearchVenues(request.Params.Longitude, request.Params.Latitude, searchCategeories)
				if nil != err {
					return err
				}

				crawl, err := user.CreateCrawl(request.Params.Name)
				if nil != err {
					return err
				}

				err = crawl.AddVenues(venues)
				if nil != err {
					return err
				}

				crawl, err = user.GetCrawl(crawl.ID)
				if nil != err {
					return err
				}

				response.Data.Crawl = crawl

				return nil
			})

		case "get_crawl":
			// {"method":"get_crawl","params":{"username":"sjsafranek@gmail.com","crawl_id":"62b6eacf-bc9a-1201-ad99-70e35fb00b10"}}
			return self.fetchUser(request, func(user *database.User) error {
				crawl, err := user.GetCrawl(request.Params.CrawlId)
				if nil != err {
					return err
				}
				response.Data.Crawl = crawl
				return nil
			})

		case "get_crawls":
			// {"method":"get_crawls","params":{"username":"sjsafranek@gmail.com"}}
			return self.fetchUser(request, func(user *database.User) error {
				crawls, err := user.GetCrawls()
				if nil != err {
					return err
				}
				response.Data.Crawls = crawls
				return nil
			})

		case "delete_crawl":
			// {"method":"delete_crawl","params":{"username":"sjsafranek@gmail.com","crawl_id":"62b6eacf-bc9a-1201-ad99-70e35fb00b10"}}
			return self.fetchUser(request, func(user *database.User) error {
				crawl, err := user.GetCrawl(request.Params.CrawlId)
				if nil != err {
					return err
				}
				return crawl.Delete()
			})

		case "up_vote":
			// {"method":"up_vote","params":{"username":"sjsafranek@gmail.com","crawl_id":"62b6eacf-bc9a-1201-ad99-70e35fb00b10", "venue_id":""}}
			return self.fetchUser(request, func(user *database.User) error {
				return user.UpVoteVenue(request.Params.CrawlId, request.Params.VenueId)
			})

		case "down_vote":
			// {"method":"down_vote","params":{"username":"sjsafranek@gmail.com","crawl_id":"62b6eacf-bc9a-1201-ad99-70e35fb00b10", "venue_id":""}}
			return self.fetchUser(request, func(user *database.User) error {
				return user.DownVoteVenue(request.Params.CrawlId, request.Params.VenueId)
			})

		default:
			return errors.New("method not found")

		}
	}()

	if nil != err {
		response.SetError(err)
	}

	return &response, err
}
