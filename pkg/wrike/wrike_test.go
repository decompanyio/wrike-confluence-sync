package wrike

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"runtime"
	"testing"
)

var (
	wrikeClient      *WrikeClient
	confluenceDomain string
)

func init() {
	godotenv.Load()
	wrikeClient = NewWrikeClient(os.Getenv("WRIKE_BASE_URL"), os.Getenv("WRIKE_TOKEN"), nil)
	confluenceDomain = os.Getenv("CONFLUENCE_DOMAIN")
}

// 프로젝트 리스트 조회
func TestProject(t *testing.T) {
	projects := wrikeClient.Projects(nil)

	fmt.Println(len(projects.Data))
	assert.NotEqual(t, projects, nil)
}

// 특정 프로젝트 조회 (링크)
func TestProjectsByLink(t *testing.T) {
	projects := wrikeClient.ProjectsByLink("https://app-us2.wrike.com/open.htm?id=850512856", nil)
	//projects := wrikeClient.ProjectsByLink("https://www.wrike.com/open.htm?id=865199939", nil)

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

// 폴더 ID로 TASK 조회
func TestTasksInProject(t *testing.T) {
	tasks := wrikeClient.TasksInProject("IEACTJ64I42AQ7PZ", confluenceDomain)

	fmt.Println(prettyPrint(tasks.Data))
	assert.NotEqual(t, tasks, nil)
}

// "2022.03.SP1"로 특정 스프린트 하위 폴더 조회
func TestSprints(t *testing.T) {
	// CPU 최대로 사용
	runtime.GOMAXPROCS(runtime.NumCPU())

	sprintWeekly := wrikeClient.Sprints("2022년 4월", "https://app-us2.wrike.com/open.htm?id=850512856", confluenceDomain)
	//sprintWeekly := wrikeClient.Sprints("2022년 3월", "https://www.wrike.com/open.htm?id=865199939")

	fmt.Println(len(sprintWeekly))
	for _, v := range sprintWeekly {
		fmt.Println(v.Title)
		fmt.Println(prettyPrint(v.Sprints))
	}

	assert.NotEqual(t, sprintWeekly, nil)
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
