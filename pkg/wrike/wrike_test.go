package wrike

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	wrikeClient      *WrikeClient
	outputDomains    []string
	confluenceDomain string
)

func init() {
	godotenv.Load()
	wrikeClient = NewWrikeClient(
		os.Getenv("WRIKE_BASE_URL"),
		os.Getenv("WRIKE_TOKEN"),
		os.Getenv("WRIKE_SPACE_ID"),
		nil)
	outputDomains = []string{os.Getenv("CONFLUENCE_DOMAIN"), "https://www.polarissharetech.net"}
	confluenceDomain = os.Getenv("CONFLUENCE_DOMAIN")
}

// 프로젝트 리스트 조회
func TestProject(t *testing.T) {
	projects := wrikeClient.Projects(nil)

	fmt.Println(len(projects.Data))
	assert.NotEqual(t, projects, nil)
// 모든 첨부파일 조회
func TestAttachmentAll(t *testing.T) {
	attachments := wrikeClient.AttachmentAll()

	fmt.Println(prettyPrint(attachments))
	assert.NotEqual(t, len(attachments), 0)
}

// 특정 프로젝트 조회 (링크)
func TestProjectsByLink(t *testing.T) {
	projects := wrikeClient.ProjectsByLink("https://app-us2.wrike.com/open.htm?id=897180682", nil)
	println(prettyPrint(projects))
	//projects := wrikeClient.ProjectsByLink("https://www.wrike.com/open.htm?id=865199939", nil)
	assert.NotEqual(t, projects, nil)
}

// ID로 프로젝트 조회
func TestProjectsByIds(t *testing.T) {
	projects := wrikeClient.ProjectsByLink("https://www.wrike.com/open.htm?id=897180682", nil)
	projectsSearch := wrikeClient.ProjectsByIds(projects.Data[0].ChildIds)
	fmt.Println(len(projectsSearch.Data))
	fmt.Printf("%+v\n", projectsSearch.Data)
	assert.NotEqual(t, projectsSearch, nil)
}

// 폴더 ID로 TASK 조회
func TestTasksInProject(t *testing.T) {
	tasks := wrikeClient.TasksInProject("IEACTJ64I42AQ7PZ", outputDomains)

	fmt.Println(prettyPrint(tasks.Data))
	assert.NotEqual(t, tasks.Data, nil)
}

// "2022.03.SP1"로 특정 스프린트 하위 폴더 조회
func TestSprints(t *testing.T) {
	sprintWeekly := wrikeClient.Sprints("2022년 5월", "https://app-us2.wrike.com/open.htm?id=850512856", outputDomains)
	//sprintWeekly := wrikeClient.Sprints("2022년 3월", "https://www.wrike.com/open.htm?id=865199939")

	fmt.Println(len(sprintWeekly))
	for _, v := range sprintWeekly {
		fmt.Println(v.Title)
		fmt.Println(prettyPrint(v.Sprints))
	}

	assert.Greater(t, len(sprintWeekly), 0)
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
