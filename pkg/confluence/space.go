package confluence

import goconfluence "github.com/virtomize/confluence-go-api"

func (c *confluence) Spaces() *goconfluence.AllSpaces {
	spaces, err := c.client.GetAllSpaces(goconfluence.AllSpacesQuery{})
	errHandler(err)
	return spaces
}
