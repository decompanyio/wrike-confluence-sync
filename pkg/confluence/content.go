package confluence

import (
	"fmt"
	goconfluence "github.com/virtomize/confluence-go-api"
	"sync"
	"wrike-confluence-sync/pkg/wrike"
)

type SyncConfig struct {
	SpMonth          string
	SprintRootLink   string
	WrikeBaseUrl     string
	WrikeToken       string
	WrikeSpaceId     string
	AncestorId       string
	OutputDomains    []string
	ConfluenceDomain string
}

func (c *Client) checkContentExist(title string) (bool, goconfluence.ContentSearch) {
	contentSearch, err := c.Client.GetContent(goconfluence.ContentQuery{
		Title:    title,
		Type:     "page",
		Expand:   []string{"version"},
		SpaceKey: c.spaceId,
	})
	errHandler(err)

	return contentSearch.Size > 0, *contentSearch
}

func (c *Client) NewContent(ancestorId string, title string, body string, contentSearch goconfluence.ContentSearch) *goconfluence.Content {
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
			Number:    contentSearch.Results[0].Version.Number + 1,
			MinorEdit: true, // 관찰자에게 알리지 않음
		}

		contentResult, err = c.Client.UpdateContent(content)
	} else {
		contentResult, err = c.Client.CreateContent(content)
	}
	errHandler(err)

	return contentResult
}

func (c *Client) SyncContent(syncConfig SyncConfig) {
	// Root 페이지 하위에 이미 sprint 페이지가 있는지 조회, 없으면 생성
	// 페이지명: yyyy년 MM월 Sprint
	// parentId 페이지 하위에 동기화
	searchTitle := syncConfig.SpMonth + " Sprint"
	var parentId string
	isExist, parentContent := c.checkContentExist(searchTitle)
	if isExist {
		parentId = parentContent.Results[0].ID
	} else {
		parentContentNew := c.NewContent(syncConfig.AncestorId, searchTitle, "", parentContent)
		parentId = parentContentNew.ID
	}

	// wrike 데이터 조회
	wrikeAPI := wrike.NewWrikeClient(
		syncConfig.WrikeBaseUrl,
		syncConfig.WrikeToken,
		syncConfig.WrikeSpaceId,
		nil)
	sprintWeekly := wrikeAPI.Sprints(syncConfig.SpMonth,
		syncConfig.SprintRootLink,
		syncConfig.OutputDomains)

	// 각 주차마다 비동기로 빠르게 처리
	var wg sync.WaitGroup
	wg.Add(len(sprintWeekly))

	for _, weekly := range sprintWeekly {
		go func(weekly wrike.SprintWeekly) {
			var content *goconfluence.Content
			title := weekly.Title
			body := NewTemplate(weekly, syncConfig.ConfluenceDomain)

			// 이미 존재하는 페이지인지 title로 조회
			contentSearch, err := c.Client.GetContent(goconfluence.ContentQuery{
				Title:    title,
				Type:     "page",
				Expand:   []string{"version"},
				SpaceKey: c.spaceId,
			})
			errHandler(err)

			// 페이지 생성/수정
			content = &goconfluence.Content{}
			content = c.NewContent(parentId, title, body, *contentSearch)
			fmt.Printf("동기화된 컨플 페이지 링크 ==> %s (%s)\n", weekly.Title, content.Links.Base+content.Links.TinyUI)
			wg.Done()
		}(weekly)
	}
	wg.Wait()
}
