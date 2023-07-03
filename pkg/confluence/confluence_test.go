package confluence

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	goconfluence "github.com/virtomize/confluence-go-api"
	"os"
	"testing"
)

var cf *Client

const sprintRootLink = "https://app-us2.wrike.com/open.htm?id=850512856"

func init() {
	godotenv.Load()

	var err error
	cf, err = NewClient(
		os.Getenv("CONFLUENCE_DOMAIN"),
		os.Getenv("CONFLUENCE_USER"),
		os.Getenv("CONFLUENCE_TOKEN"),
		os.Getenv("CONFLUENCE_SPACEID"),
	)
	if err != nil {
		panic(err)
	}
}

// 스페이스 리스트 조회
func TestSpace(t *testing.T) {
	spaces, err := cf.Client.GetAllSpaces(goconfluence.AllSpacesQuery{})
	assert.NoError(t, err)

	for i, space := range spaces.Results {
		fmt.Printf("%d번째 : %s\n", i, space.Name)
	}

	assert.NotEqual(t, spaces, nil)
}
