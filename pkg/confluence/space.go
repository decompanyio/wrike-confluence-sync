package confluence

import goconfluence "github.com/virtomize/confluence-go-api"

// Spaces 스페이스 리스트 조회
func (c *Client) Spaces() (*goconfluence.AllSpaces, error) {
	return c.Client.GetAllSpaces(goconfluence.AllSpacesQuery{})

}
