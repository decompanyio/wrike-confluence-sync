package confluence

import (
	"fmt"
	goconfluence "github.com/virtomize/confluence-go-api"
	"sync"
	"wrike-confluence-sync/pkg/wrike"
)

func (c ConfluenceClient) NewContent(ancestorId string, title string, body string, contentSearch goconfluence.ContentSearch) *goconfluence.Content {
	// 컨플 컨텐트 구조체 생성
	content := &goconfluence.Content{
		Title: title,
		Type:  "page",
		Body: goconfluence.Body{
			Storage: goconfluence.Storage{
				Value:          body,
				Representation: "editor2",
			},
		},
		Ancestors: []goconfluence.Ancestor{
			{
				ID: ancestorId,
			},
		},
		Space: goconfluence.Space{Key: c.spaceId},
	}

	// 컨플 페이지 등록 또는 수정
	var contentResult *goconfluence.Content
	var err error
	if contentSearch.Size > 0 {
		content.ID = contentSearch.Results[0].ID
		content.Version = &goconfluence.Version{
			Number: contentSearch.Results[0].Version.Number + 1,
		}

		contentResult, err = c.Client.UpdateContent(content)
	} else {
		contentResult, err = c.Client.CreateContent(content)
	}
	errHandler(err)

	return contentResult
}

type SyncConfig struct {
	SpMonth          string
	SprintRootLink   string
	WrikeBaseUrl     string
	WrikeToken       string
	AncestorId       string
	OutputDomains    []string
	ConfluenceDomain string
}

func (c *ConfluenceClient) SyncContent(syncConfig SyncConfig) {
	// wrike 데이터 조회
	wrikeAPI := wrike.NewWrikeClient(syncConfig.WrikeBaseUrl, syncConfig.WrikeToken, nil)
	sprintWeekly := wrikeAPI.Sprints(syncConfig.SpMonth, syncConfig.SprintRootLink, syncConfig.OutputDomains)

	// 각 주차마다 비동기로 빠르게 처리
	var wg sync.WaitGroup
	wg.Add(len(sprintWeekly))

	// 익명 함수
	newContent := func(weekly wrike.SprintWeekly) {
		var content *goconfluence.Content
		title := weekly.Title
		body := NewTemplate(weekly.Sprints, syncConfig.ConfluenceDomain)

		// 이미 존재하는 페이지인지 title로 조회
		contentSearch, err := c.Client.GetContent(goconfluence.ContentQuery{
			Title:  title,
			Type:   "page",
			Expand: []string{"version"},
		})
		errHandler(err)

		content = &goconfluence.Content{}
		content = c.NewContent(syncConfig.AncestorId, title, body, *contentSearch)
		fmt.Printf("동기화된 컨플 페이지 링크 ==> %s (%s)\n", weekly.Title, content.Links.Base+content.Links.TinyUI)
		wg.Done()
	}

	for _, weekly := range sprintWeekly {
		go newContent(weekly)
	}
	wg.Wait()
	fmt.Println("Done")
}
