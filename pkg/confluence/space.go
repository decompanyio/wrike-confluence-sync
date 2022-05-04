package confluence

import goconfluence "github.com/virtomize/confluence-go-api"

func (c *ConfluenceClient) Spaces() *goconfluence.AllSpaces {
	spaces, err := c.Client.GetAllSpaces(goconfluence.AllSpacesQuery{})
	errHandler(err)
	return spaces
}
