package confluence

import (
	goconfluence "github.com/virtomize/confluence-go-api"
	"log"
)

const (
	spaceId = "~166200948"
)

type confluence struct {
	client *goconfluence.API
}

func NewConfluence(domain string, username string, token string) *confluence {
	client, err := goconfluence.NewAPI(domain+"/wiki/rest/api", username, token)
	errHandler(err)
	return &confluence{
		client: client,
	}
}

func errHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
