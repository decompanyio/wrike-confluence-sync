package wrike

import (
	"time"
)

type Tasks struct {
	Kind string `json:"kind"`
	Data []Task `json:"data"`
}

type Task struct {
	ID               string    `json:"id"`
	AccountID        string    `json:"accountId"`
	Coworkers        []User    `json:"coworkers"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	BriefDescription string    `json:"briefDescription"`
	ParentIds        []string  `json:"parentIds"`
	SuperParentIds   []string  `json:"superParentIds"`
	SharedIds        []string  `json:"sharedIds"`
	ResponsibleIds   []string  `json:"responsibleIds"`
	Status           string    `json:"status"`
	Importance       string    `json:"importance"`
	CreatedDate      time.Time `json:"createdDate"`
	UpdatedDate      time.Time `json:"updatedDate"`
	Dates            struct {
		Type     string `json:"type"`
		Duration int    `json:"duration"`
		Start    string `json:"start"`
		Due      string `json:"due"`
	} `json:"dates"`
	Scope          string        `json:"scope"`
	AuthorIds      []string      `json:"authorIds"`
	CustomStatusID string        `json:"customStatusId"`
	HasAttachments bool          `json:"hasAttachments"`
	Attachments    []Attachment  `json:"attachments"`
	Permalink      string        `json:"permalink"`
	Priority       string        `json:"priority"`
	FollowedByMe   bool          `json:"followedByMe"`
	FollowerIds    []string      `json:"followerIds"`
	SuperTaskIds   []string      `json:"superTaskIds"`
	SubTaskIds     []interface{} `json:"subTaskIds"`
	DependencyIds  []string      `json:"dependencyIds"`
	Metadata       []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"metadata"`
	CustomFields []struct {
		ID    string `json:"id"`
		Value string `json:"value"`
	} `json:"customFields"`
}

type AllTaskMap map[string][]Task

// 모든 작업 조회
func (w *WrikeClient) TaskAll(rootFolderId string) AllTaskMap {
	tasks := Tasks{}
	urlQuery := map[string]string{
		"status":      `["Active","Completed"]`,
		"fields":      `["authorIds","responsibleIds","hasAttachments","parentIds"]`,
		"descendants": "true",
		"subTasks":    "true",
		"sortField":   `DueDate`,
	}
	w.newAPI("/folders/"+rootFolderId+"/tasks", urlQuery, &tasks)

	taskAll := AllTaskMap{}
	for _, task := range tasks.Data {
		for _, id := range task.ParentIds {
			taskAll[id] = append(taskAll[id], task)
		}
	}

	return taskAll
}
