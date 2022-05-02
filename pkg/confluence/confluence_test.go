package confluence

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	goconfluence "github.com/virtomize/confluence-go-api"
	"os"
	"testing"
)

var cf *confluence

func init() {
	godotenv.Load()
	cf = NewConfluence(
		os.Getenv("DOMAIN"),
		os.Getenv("USER"),
		os.Getenv("TOKEN"),
		os.Getenv("SPACEID"),
	)
}

// 스페이스 리스트 조회
func TestSpace(t *testing.T) {
	spaces, err := cf.client.GetAllSpaces(goconfluence.AllSpacesQuery{})
	errHandler(err)

	for i, space := range spaces.Results {
		fmt.Printf("%d번째 : %s\n", i, space.Name)
	}

	assert.NotEqual(t, spaces, nil)
}

// wrike 데이터로 html 동적 생성
func TestNewTemplate(t *testing.T) {
	data := NewTemplate()
	fmt.Println(data)
	assert.NotEqual(t, data, nil)
}

func TestCreateContent(t *testing.T) {
	title := "Sprint Temp Page"
	ancestorId := os.Getenv("ANCESTORID")

	// 이미 존재하는 페이지인지 title로 조회
	contentSearch, err := cf.client.GetContent(goconfluence.ContentQuery{
		Title:  title,
		Type:   "page",
		Expand: []string{"version"},
	})
	errHandler(err)

	content := &goconfluence.Content{}
	content = cf.newContent(ancestorId, *contentSearch)

	fmt.Println(content.Status)
	fmt.Println(content.Links.Base + content.Links.TinyUI)

	assert.NotEqual(t, content, nil)
}
