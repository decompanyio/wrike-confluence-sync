package confluence

import (
	goconfluence "github.com/virtomize/confluence-go-api"
	"log"
)

type ConfluenceClient struct {
	Client  *goconfluence.API
	spaceId string
}

func NewConfluenceClient(domain string, username string, token string, spaceId string) *ConfluenceClient {
	client, err := goconfluence.NewAPI(domain+"/wiki/rest/api", username, token)
	errHandler(err)
	return &ConfluenceClient{
		Client:  client,
		spaceId: spaceId,
	}
}

func errHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
