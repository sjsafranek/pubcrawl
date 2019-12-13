package foursquare

func New(clientID, clientSecret string) *Client {
	return &Client{ClientID: clientID, ClientSecret: clientSecret}
}
