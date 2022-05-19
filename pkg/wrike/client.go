package wrike

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

const Host string = "https://www.wrike.com/api/v4"

type WrikeClient struct {
	host       string
	bearer     string
	spaceId    string
	httpClient *http.Client
}

// WrikeClient 생성자
func NewWrikeClient(host string, bearer string, spaceId string, httpClient *http.Client) *WrikeClient {
	if len(bearer) == 0 {
		log.Fatal("토큰이 없습니다.")
	}

	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	return &WrikeClient{
		host:       host,
		bearer:     bearer,
		spaceId:    spaceId,
		httpClient: httpClient,
	}
}

// API 공통 모듈 (internal)
func (w *WrikeClient) newAPI(uri string, urlQuery map[string]string, target interface{}) {
	req, err := http.NewRequest("GET", w.host+uri, nil)
	errorHandler(err)

	req.Header.Add("Authorization", "Bearer "+w.bearer)

	if urlQuery != nil {
		q := req.URL.Query()
		for k, v := range urlQuery {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := w.httpClient.Do(req)
	errorHandler(err)

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(target)
		errorHandler(err)
	} else {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(string(bodyBytes))
	}
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
