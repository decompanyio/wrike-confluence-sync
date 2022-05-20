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
	wrikeClient   *Client
	outputDomains []string
)

func init() {
	godotenv.Load()
	wrikeClient = NewWrikeClient(
		os.Getenv("WRIKE_BASE_URL"),
		os.Getenv("WRIKE_TOKEN"),
		os.Getenv("WRIKE_SPACE_ID"),
		nil)
	outputDomains = []string{os.Getenv("CONFLUENCE_DOMAIN"), "https://www.polarissharetech.net"}
}

// 유저 조회
func TestUsers(t *testing.T) {
	users := wrikeClient.UserAll()

	fmt.Println(prettyPrint(users))
	assert.Greater(t, len(users), 0)
}

// 모든 폴더 조회
func TestFoldersAll(t *testing.T) {
	folders := wrikeClient.FolderAll()

	fmt.Println(prettyPrint(folders))
	assert.Greater(t, len(folders), 0)
}

// 모든 작업 조회
func TestTaskAll(t *testing.T) {
	tasksParentId, tasksTaskId := wrikeClient.TaskAll("IEACTJ64I42PUE7V")

	assert.Greater(t, len(tasksParentId), 0)
	assert.Greater(t, len(tasksTaskId), 0)
}

// 모든 첨부파일 조회
func TestAttachmentAll(t *testing.T) {
	attachments := wrikeClient.AttachmentAll()

	fmt.Println(prettyPrint(attachments))
	assert.Greater(t, len(attachments), 0)
}

// 특정 프로젝트 조회 (링크)
func TestProjectsByLink(t *testing.T) {
	projects := wrikeClient.ProjectsByLink("https://app-us2.wrike.com/open.htm?id=897180682", nil)
	println(prettyPrint(projects))
	//projects := wrikeClient.ProjectsByLink("https://www.wrike.com/open.htm?id=865199939", nil)
	assert.Greater(t, len(projects.Data), 0)
}

// ID로 프로젝트 조회
func TestProjectsByIds(t *testing.T) {
	projects := wrikeClient.ProjectsByLink("https://www.wrike.com/open.htm?id=897180682", nil)
	projectsSearch := wrikeClient.ProjectsByIds(projects.Data[0].ChildIds)
	fmt.Println(len(projectsSearch.Data))
	fmt.Printf("%+v\n", projectsSearch.Data)
	assert.Greater(t, len(projectsSearch.Data), 0)
}

// "2022.03.SP1"로 특정 스프린트 하위 폴더 조회
func TestSprints(t *testing.T) {
	sprintWeekly := wrikeClient.Sprints("2022년 5월", "https://app-us2.wrike.com/open.htm?id=850512856", outputDomains)
	//sprintWeekly := wrikeClient.Sprints("2022년 3월", "https://www.wrike.com/open.htm?id=865199939")

	fmt.Println(len(sprintWeekly))
	for _, v := range sprintWeekly {
		fmt.Println(prettyPrint(v.ImportanceStatistics["Normal"]))
	}

	assert.Greater(t, len(sprintWeekly), 0)
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
