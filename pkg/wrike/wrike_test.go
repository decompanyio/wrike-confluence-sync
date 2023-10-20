package wrike

import (
	"github.com/joho/godotenv"
	"github.com/k0kubun/pp/v3"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	wrikeClient *Client
)

func init() {
	godotenv.Load()
	wrikeClient, _ = NewClient(
		os.Getenv("WRIKE_BASE_URL"),
		os.Getenv("WRIKE_TOKEN"),
		os.Getenv("WRIKE_SPACE_ID"))
}

// 유저 조회
func TestUsers(t *testing.T) {
	users, err := wrikeClient.GetAllUsers()

	assert.NoError(t, err)
	assert.Greater(t, len(users), 0)
	pp.Println(users)
}

// 모든 폴더 조회
func TestFoldersAll(t *testing.T) {
	folders, err := wrikeClient.GetAllProjects()

	assert.NoError(t, err)
	assert.Greater(t, len(folders), 0)
	pp.Println(folders)
}

// 모든 작업 조회
func TestTaskAll(t *testing.T) {
	tasksParentId, tasksTaskId, err := wrikeClient.TaskAll("IEACTJ64I42PUE7V")

	assert.NoError(t, err)
	assert.Greater(t, len(tasksParentId), 0)
	assert.Greater(t, len(tasksTaskId), 0)
}

// 모든 첨부파일 조회
func TestAttachmentAll(t *testing.T) {
	attachments, err := wrikeClient.GetAllAttachments()

	assert.NoError(t, err)
	assert.Greater(t, len(attachments), 0)
}

// 특정 프로젝트 조회 (링크)
func TestProjectsByLink(t *testing.T) {
	projects, err := wrikeClient.GetProjectsByLink("https://app-us2.wrike.com/open.htm?id=897180682", nil)

	assert.NoError(t, err)
	assert.Greater(t, len(projects.Data), 0)
	pp.Println(projects.Data)
}

func TestProjectsByLinkV2(t *testing.T) {
	projects, err := wrikeClient.GetProjectsByLink("https://app-us2.wrike.com/open.htm?id=1084138983", nil)

	assert.NoError(t, err)
	assert.Greater(t, len(projects.Data), 0)

	pp.Println(projects.Data)
}
