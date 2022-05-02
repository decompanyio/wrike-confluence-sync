package confluence

import goconfluence "github.com/virtomize/confluence-go-api"

func (c *confluence) CreateContent() *goconfluence.Content {
	content := c.newContent("1015809", "", 0)

	content, err := c.client.CreateContent(content)
	errHandler(err)

	return content
}

func (c *confluence) UpdateContent(contentSearch goconfluence.Content) *goconfluence.Content {
	content := c.newContent("1015809", contentSearch.ID, contentSearch.Version.Number)

	content, err := c.client.UpdateContent(content)
	errHandler(err)

	return content
}

func (c confluence) newContent(ancestorId string, contentId string, version int) *goconfluence.Content {
	content := &goconfluence.Content{
		Title: "api-test-02",
		Type:  "page",
		Body: goconfluence.Body{
			Storage: goconfluence.Storage{
				Value:          NewTemplate(),
				Representation: "editor2",
			},
		},
		Ancestors: []goconfluence.Ancestor{
			goconfluence.Ancestor{
				ID: ancestorId,
			},
		},
		Space: goconfluence.Space{Key: spaceId},
	}
	if len(contentId) > 0 {
		content.ID = contentId
	}
	// 페이지의 버전 증가
	if version > 0 {
		content.Version = &goconfluence.Version{
			Number: version + 1,
		}
	}

	return content
}
