package wrike

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
)

type Client struct {
	host       string
	bearer     string
	spaceId    string
	httpClient *resty.Client
}

// NewWrikeClient wrike client 생성
func NewWrikeClient(host string, bearer string, spaceId string) (*Client, error) {
	hostValid, err := url.ParseRequestURI(host)
	if err != nil {
		return nil, errors.New("failed to create wrike client")
	}

	if len(bearer) == 0 {
		return nil, errors.New("failed to create wrike client")
	}

	return &Client{
		host:       hostValid.String(),
		bearer:     bearer,
		spaceId:    spaceId,
		httpClient: resty.New(),
	}, nil
}

// newAPI wrike api 호출
func (w *Client) newAPI(uri string, urlQuery map[string]string, target interface{}) {
	// request 생성
	req, err := http.NewRequest("GET", w.host+uri, nil)
	if err != nil {
		log.Err(err).Msg("failed to call wrike api")
		panic(err)
	}

	// query string
	if urlQuery != nil {
		q := req.URL.Query()
		for k, v := range urlQuery {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	// api 호출
	resp, err := w.httpClient.R().
		SetHeader("Authorization", "Bearer "+w.bearer).
		SetResult(target).
		Get(req.URL.String())

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Err(err).Msg("failed to call wrike api")
		panic(err)
	}
	if resp.StatusCode() != http.StatusOK {
		log.Error().Caller().Msg(string(resp.Body()))
		panic(err)
	}
}
