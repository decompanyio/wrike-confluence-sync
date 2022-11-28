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

func NewConfluenceClient(domain string, username string, token string, spaceId string) (*Client, error) {
	domainValid, err := url.ParseRequestURI(domain)
	if err != nil {
		return nil, err
	}

	client, err := goconfluence.NewAPI(domainValid.String()+"/wiki/rest/api", username, token)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client:  client,
		spaceId: spaceId,
	}, nil
}

func errHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
