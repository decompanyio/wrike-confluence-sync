package confluence

import goconfluence "github.com/virtomize/confluence-go-api"

// Spaces 스페이스 리스트 조회
func (c *Client) Spaces() *goconfluence.AllSpaces {
	spaces, err := c.Client.GetAllSpaces(goconfluence.AllSpacesQuery{})
	errHandler(err)
	return spaces
}
