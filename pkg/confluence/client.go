package confluence

import (
	goconfluence "github.com/virtomize/confluence-go-api"
	"log"
)

type confluence struct {
	client  *goconfluence.API
	spaceId string
}

func NewConfluence(domain string, username string, token string, spaceId string) *confluence {
	client, err := goconfluence.NewAPI(domain+"/wiki/rest/api", username, token)
	errHandler(err)
	return &confluence{
		client:  client,
		spaceId: spaceId,
	}
}

func errHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
