package confluence

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	goconfluence "github.com/virtomize/confluence-go-api"
	"os"
	"testing"
	"wrike-confluence-sync/pkg/wrike"
)

var cf *Client

const sprintRootLink = "https://app-us2.wrike.com/open.htm?id=850512856"

func init() {
	godotenv.Load()
	cf = NewConfluenceClient(
		os.Getenv("CONFLUENCE_DOMAIN"),
		os.Getenv("CONFLUENCE_USER"),
		os.Getenv("CONFLUENCE_TOKEN"),
		os.Getenv("CONFLUENCE_SPACEID"),
	)
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

// wrike 데이터로 html 동적 생성
func TestNewTemplate(t *testing.T) {
	// wrike 데이터 조회
	wrikeAPI := wrike.NewWrikeClient(
		os.Getenv("WRIKE_BASE_URL"),
		os.Getenv("WRIKE_TOKEN"),
		os.Getenv("WRIKE_SPACE_ID"),
		nil)
	sprintWeekly, err := wrikeAPI.Sprints("2022년 04월", sprintRootLink, []string{"https://google.com"})
	assert.NoError(t, err)

	for _, weekly := range sprintWeekly {
		data := NewTemplate(weekly.Sprints, os.Getenv("DOMAIN"))
		fmt.Println(data)
		assert.NotEqual(t, data, nil)
	}
}

func TestCreateContent(t *testing.T) {
	syncConfig := SyncConfig{
		SpMonth:        "2022년 4월",
		SprintRootLink: sprintRootLink,
		WrikeBaseUrl:   os.Getenv("WRIKE_BASE_URL"),
		WrikeToken:     os.Getenv("WRIKE_TOKEN"),
		AncestorId:     os.Getenv("ANCESTORID"),
	}

	cf.SyncContent(syncConfig)

	assert.NotEqual(t, 1, nil)
}
