package wrike

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/url"
)

type Client struct {
	host       *url.URL
	bearer     string
	spaceId    string
	httpClient *resty.Client
}

func NewClient(host, bearer, spaceId string) (*Client, error) {
	if host == "" || bearer == "" {
		return nil, fmt.Errorf("both host and bearer are required to create a Wrike client")
	}

	parsedHost, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("invalid host: %v", err)
	}

	return &Client{
		host:       parsedHost,
		bearer:     bearer,
		spaceId:    spaceId,
		httpClient: resty.New(),
	}, nil
}

func (c *Client) callAPI(endpoint string, queryParams map[string]string, result interface{}) error {
	resp, err := c.httpClient.R().
		SetHeader("Authorization", "Bearer "+c.bearer).
		SetQueryParams(queryParams).
		SetResult(result).
		Get(c.host.String() + endpoint)

	if err != nil {
		return fmt.Errorf("failed to call Wrike API: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode(), resp.String())
	}

	return nil
}
