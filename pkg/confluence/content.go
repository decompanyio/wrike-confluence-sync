package confluence

import goconfluence "github.com/virtomize/confluence-go-api"

func (c confluence) newContent(ancestorId string, title string, body string, contentSearch goconfluence.ContentSearch) *goconfluence.Content {
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

		contentResult, err = c.client.UpdateContent(content)
	} else {
		contentResult, err = c.client.CreateContent(content)
	}
	errHandler(err)

	return contentResult
}
