package confluence

import (
	"fmt"
	"github.com/rs/zerolog/log"
	goconfluence "github.com/virtomize/confluence-go-api"
	"time"
	"wrike-confluence-sync/pkg/wrike"
)

type SyncConfig struct {
	Date             time.Time
	AncestorId       string
	OutputDomains    []string
	ConfluenceDomain string
}

// getPageContent fetches and returns the content of a Confluence page.
func (c *Client) getPageContent(title string) (*goconfluence.ContentSearch, error) {
	contentSearch, err := c.Client.GetContent(goconfluence.ContentQuery{
		Title:    title,
		Type:     "page",
		Expand:   []string{"version"},
		SpaceKey: c.spaceId,
	})
	if err != nil {
		log.Err(err).Msgf("GetContent failed : %s", title)
		return nil, err
	}
	return contentSearch, nil
}

// createOrUpdateContent creates or updates a Confluence page and returns the result.
func (c *Client) createOrUpdateContent(ancestorId, title, body string, existingContent *goconfluence.ContentSearch) (*goconfluence.Content, error) {
	content := c.prepareContent(ancestorId, title, body)

	if existingContent.Size > 0 {
		return c.updateContent(content, existingContent)
	}
	return c.createContent(content)
}

// prepareContent prepares a Confluence page content.
func (c *Client) prepareContent(ancestorId, title, body string) *goconfluence.Content {
	return &goconfluence.Content{
		Title: title,
		Type:  "page",
		Body: goconfluence.Body{
			Storage: goconfluence.Storage{
				Value:          body,
				Representation: "editor2",
			},
		},
		Ancestors: []goconfluence.Ancestor{{ID: ancestorId}},
		Space:     &goconfluence.Space{Key: c.spaceId},
	}
}

// updateContent updates a Confluence page and returns the result.
func (c *Client) updateContent(content *goconfluence.Content, existingContent *goconfluence.ContentSearch) (*goconfluence.Content, error) {
	content.ID = existingContent.Results[0].ID
	content.Version = &goconfluence.Version{
		Number:    existingContent.Results[0].Version.Number + 1,
		MinorEdit: true,
	}
	return c.Client.UpdateContent(content)
}

// createContent creates a Confluence page and returns the result.
func (c *Client) createContent(content *goconfluence.Content) (*goconfluence.Content, error) {
	return c.Client.CreateContent(content)
}

// ensureParentContent ensures the parent content exists and returns its ID.
func (c *Client) ensureParentContent(searchTitle string, ancestorId string) (string, error) {
	parentContent, err := c.getPageContent(searchTitle)
	if err != nil {
		return "", err
	}

	if parentContent.Size > 0 {
		return parentContent.Results[0].ID, nil
	}

	newParentContent, err := c.createOrUpdateContent(ancestorId, searchTitle, "", parentContent)
	if err != nil {
		return "", err
	}

	return newParentContent.ID, nil
}

// SyncContent syncs Wrike's sprint data to a Confluence page.
func (c *Client) SyncContent(sprint wrike.SprintWeekly, syncConfig SyncConfig) error {
	searchTitle := syncConfig.Date.Format("2006년 1월") + " Sprint"

	parentId, err := c.ensureParentContent(searchTitle, syncConfig.AncestorId)
	if err != nil {
		return err
	}

	sprintContent, err := c.getPageContent(sprint.Title)
	if err != nil {
		return err
	}

	body := NewTemplate(sprint, syncConfig.ConfluenceDomain)
	updatedContent, err := c.createOrUpdateContent(parentId, sprint.Title, body, sprintContent)
	if err != nil {
		return err
	}

	fmt.Printf("동기화된 컨플 페이지 ==> %s (%s)\n", sprint.Title, updatedContent.Links.Base+updatedContent.Links.TinyUI)
	return nil
}
