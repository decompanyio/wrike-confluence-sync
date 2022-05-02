package wrike

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var wrikeClient *WrikeClient

func init() {
	godotenv.Load()
	wrikeClient = NewWrikeClient(os.Getenv("WRIKE_TOKEN"), nil)
}

// 프로젝트 리스트 조회
func TestProject(t *testing.T) {
	projects := wrikeClient.Projects(nil)

	fmt.Println(len(projects.Data))
	assert.NotEqual(t, projects, nil)
}

// 특정 프로젝트 조회 (링크)
func TestProjectsByLink(t *testing.T) {
	projects := wrikeClient.ProjectsByLink("https://www.wrike.com/open.htm?id=865199939", nil)

	fmt.Println(len(projects.Data))
	assert.NotEqual(t, projects, nil)
}

// ID로 프로젝트 조회
func TestProjectsByIds(t *testing.T) {
	projects := wrikeClient.ProjectsByLink("https://www.wrike.com/open.htm?id=865199939", nil)
	projectsSearch := wrikeClient.ProjectsByIds(projects.Data[0].ChildIds)

	fmt.Println(len(projectsSearch.Data))
	fmt.Printf("%+v\n", projectsSearch.Data)
	assert.NotEqual(t, projectsSearch, nil)
}

// "2022.03.SP1"로 특정 스프린트 하위 폴더 조회
func TestSprints(t *testing.T) {
	sprints := wrikeClient.Sprints("2022.03.SP1")

	fmt.Println(len(sprints))
	for _, v := range sprints {
		fmt.Printf("%+v\n\n", v)
	}

	assert.NotEqual(t, sprints, nil)
}
