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
	AncestorId       string
	OutputDomains    []string
	ConfluenceDomain string
}

// checkContentExist 컨플 페이지가 존재하는지 확인한다
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

// NewContent 컨플 페이지를 생성한다
//
// @param ancestorId 부모 컨플의 ID. 부모 컨플 하위 페이지로 생성
// @param title 컨플 페이지 제목
// @param body 컨플 본문
// @param contentSearch 해당 페이지가 존재하는지에 대한 검색 결과
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
		Space: &goconfluence.Space{Key: c.spaceId},
	}

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

// SyncContent wrike의 sprint 데이터를 조회하여 컨플 페이지로 등록한다
func (c *Client) SyncContent(syncConfig SyncConfig, wrikeClient *wrike.Client) error {
	searchTitle := syncConfig.SpMonth + " Sprint"
	var parentId string

	isExist, parentContent := c.checkContentExist(searchTitle)
	if isExist {
		parentId = parentContent.Results[0].ID
	} else {
		parentContentNew := c.NewContent(syncConfig.AncestorId, searchTitle, "", parentContent)
		parentId = parentContentNew.ID
	}

	sprintWeekly, err := wrikeClient.Sprints(
		syncConfig.SpMonth,
		syncConfig.SprintRootLink,
		syncConfig.OutputDomains)

	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(sprintWeekly))

	for _, weekly := range sprintWeekly {
		go func(weekly wrike.SprintWeekly) {
			defer wg.Done()

			title := weekly.Title
			body := NewTemplate(weekly, syncConfig.ConfluenceDomain)

			// 이미 존재하는 페이지인지 확인
			contentSearch, errChild := c.Client.GetContent(goconfluence.ContentQuery{
				Title:    title,
				Type:     "page",
				Expand:   []string{"version"},
				SpaceKey: c.spaceId,
			})
			errHandler(errChild)

			// 페이지 생성/수정
			content := c.NewContent(parentId, title, body, *contentSearch)
			fmt.Printf("동기화된 컨플 페이지 링크 ==> %s (%s)\n", weekly.Title, content.Links.Base+content.Links.TinyUI)
		}(weekly)
	}
	wg.Wait()
	return nil
}
