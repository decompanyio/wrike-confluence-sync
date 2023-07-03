package confluence

import (
	goconfluence "github.com/virtomize/confluence-go-api"
	"net/url"
)

type Client struct {
	Client  *goconfluence.API // confluence client
	spaceId string            // confluence space key
}

// NewClient confluence client 생성
func NewClient(domain string, username string, token string, spaceId string) (*Client, error) {
	// domain validation
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
