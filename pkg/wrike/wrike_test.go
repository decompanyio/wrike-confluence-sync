package wrike

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var (
	wrikeClient   *Client
	outputDomains []string
)

func init() {
	godotenv.Load()
	wrikeClient, _ = NewWrikeClient(
		os.Getenv("WRIKE_BASE_URL"),
		os.Getenv("WRIKE_TOKEN"),
		os.Getenv("WRIKE_SPACE_ID"))
	outputDomains = []string{
		os.Getenv("CONFLUENCE_DOMAIN"),
		"https://www.polarissharetech.net",
		"https://www.figma.com",
		"https://www.polarisoffice.com",
		"https://github.com/decompanyio",
	}
}

// 유저 조회
func TestUsers(t *testing.T) {
	users := wrikeClient.UserAll()

	fmt.Println(prettyPrint(users))
	assert.Greater(t, len(users), 0)
}

// 모든 폴더 조회
func TestFoldersAll(t *testing.T) {
	folders := wrikeClient.ProjectAll()

	fmt.Println(prettyPrint(folders))
	assert.Greater(t, len(folders), 0)
}

// 모든 작업 조회
func TestTaskAll(t *testing.T) {
	tasksParentId, tasksTaskId := wrikeClient.TaskAll("IEACTJ64I42PUE7V")

	fmt.Println(prettyPrint(tasksParentId))
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
	projects := wrikeClient.ProjectsByLink("https://www.wrike.com/open.htm?id=850512856", nil)
	projectsSearch := wrikeClient.ProjectsByIds(projects.Data[0].ChildIds)
	fmt.Println(len(projectsSearch.Data))
	fmt.Printf("%+v\n", projectsSearch.Data)
	assert.Greater(t, len(projectsSearch.Data), 0)
}

// "2022.03.SP1"로 특정 스프린트 하위 폴더 조회
func TestSprints(t *testing.T) {
	rootLink := "https://app-us2.wrike.com/open.htm?id=1084138983"

	data := AllData{
		UserAll:       wrikeClient.UserAll(),
		AttachmentAll: wrikeClient.AttachmentAll(),
		ProjectAll:    wrikeClient.ProjectAll(),
	}

	date := time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC)
	sprint, err := wrikeClient.Sprint(date, rootLink, outputDomains, data)
	assert.NoError(t, err)
	assert.NotNil(t, sprint)

	fmt.Println(sprint)
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func TestProjectsByLinkV2(t *testing.T) {
	projects := wrikeClient.ProjectsByLink("https://app-us2.wrike.com/open.htm?id=1084138983", nil)
	println(prettyPrint(projects))
	assert.Greater(t, len(projects.Data), 0)
}
