package wrike

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"log/slog"
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

	// wrike API를 사용하기 위한 설정
	// - BaseURL: API 호출을 위한 기본 URL
	// - SetAuthToken: API 호출을 위한 인증 토큰
	httpClient := resty.New()
	httpClient.BaseURL = parsedHost.String()
	httpClient.SetAuthToken(bearer)

	return &Client{
		host:       parsedHost,
		bearer:     bearer,
		spaceId:    spaceId,
		httpClient: httpClient,
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

type AllData struct {
	UserAll       AllUserMap
	AttachmentAll AllAttachmentMap
	ProjectAll    AllProjectMap
}

func (c *Client) GetAllData() (AllData, error) {
	users, err := c.GetAllUsers()
	if err != nil {
		slog.Error("failed to get all users", slog.String("error", err.Error()))
		return AllData{}, err
	}
	attachments, err := c.GetAllAttachments()
	if err != nil {
		slog.Error("failed to get all attachments", slog.String("error", err.Error()))
		return AllData{}, err
	}

	projects, err := c.GetAllProjects()
	if err != nil {
		slog.Error("failed to get all projects", slog.String("error", err.Error()))
		return AllData{}, err
	}

	return AllData{
		UserAll:       users,
		AttachmentAll: attachments,
		ProjectAll:    projects,
	}, nil
}
