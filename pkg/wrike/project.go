package wrike

import (
	"strings"
	"time"
)

type Projects struct {
	Kind string    `json:"kind"`
	Data []Project `json:"data"`
}

type Project struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Children    []string `json:"children"`
	ChildIds    []string `json:"childIds"`
	Scope       string   `json:"scope"`
	Project     struct {
		AuthorID       string    `json:"authorId"`
		OwnerIds       []string  `json:"ownerIds"`
		CustomStatusID string    `json:"customStatusId"`
		StartDate      string    `json:"startDate"`
		EndDate        string    `json:"endDate"`
		CreatedDate    time.Time `json:"createdDate"`
	} `json:"project,omitempty"`
}

type AllFolderMap map[string]Project

func (afm *AllFolderMap) findFolderByIds(folderIds []string) []Project {
	var projectTemp []Project
	for _, id := range folderIds {
		projectTemp = append(projectTemp, (*afm)[id])
	}
	return projectTemp
}

// FolderAll 모든 폴더 조회(프로젝트 제외) 후 folderId가 키인 map 반환
func (w *Client) FolderAll() AllFolderMap {
	var folders Projects
	urlQuery := map[string]string{
		"deleted": "false",
		"project": "false",
		"fields":  `["description"]`,
	}
	w.newAPI("/spaces/"+w.spaceId+"/folders", urlQuery, &folders)

	allFolderMap := AllFolderMap{}
	for _, folder := range folders.Data {
		allFolderMap[folder.ID] = folder
	}

	return allFolderMap
}

// ProjectsByLink 프로젝트 & 폴더 - 고정 링크로 필터
func (w *Client) ProjectsByLink(link string, urlQuery map[string]string) Projects {
	if urlQuery == nil {
		urlQuery = map[string]string{
			"deleted": "false",
		}
	}
	if len(link) > 0 {
		urlQuery["permalink"] = link
	}

	projects := Projects{}
	w.newAPI("/folders", urlQuery, &projects)

	return projects
}

// ProjectsByIds 프로젝트 & 폴더 - ID로 필터
func (w *Client) ProjectsByIds(ids []string) Projects {
	projects := Projects{}
	if len(ids) > 0 {
		w.newAPI("/folders/"+strings.Join(ids, ","), nil, &projects)
	}

	return projects
}
