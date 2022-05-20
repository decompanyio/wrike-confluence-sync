package confluence

import (
	goconfluence "github.com/virtomize/confluence-go-api"
	"log"
	"net/url"
)

type Client struct {
	Client  *goconfluence.API
	spaceId string
}

func NewConfluenceClient(domain string, username string, token string, spaceId string) *Client {
	domainValid, err := url.ParseRequestURI(domain)
	errHandler(err)

	client, err := goconfluence.NewAPI(domainValid.String()+"/wiki/rest/api", username, token)
	errHandler(err)

	return &Client{
		Client:  client,
		spaceId: spaceId,
	}
}

func errHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
