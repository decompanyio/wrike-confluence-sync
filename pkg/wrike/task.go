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
	Scope          string       `json:"scope"`
	AuthorIds      []string     `json:"authorIds"`
	CustomStatusID string       `json:"customStatusId"`
	HasAttachments bool         `json:"hasAttachments"`
	Attachments    []Attachment `json:"attachments"`
	Permalink      string       `json:"permalink"`
	Priority       string       `json:"priority"`
	FollowedByMe   bool         `json:"followedByMe"`
	FollowerIds    []string     `json:"followerIds"`
	SuperTaskIds   []string     `json:"superTaskIds"`
	SubTaskIds     []string     `json:"subTaskIds"`
	DependencyIds  []string     `json:"dependencyIds"`
	Metadata       []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"metadata"`
	CustomFields []struct {
		ID    string `json:"id"`
		Value string `json:"value"`
	} `json:"customFields"`
}

// TaskMapByParentIdAll key: task의 ParentId
type TaskMapByParentIdAll map[string][]Task

// TaskMapByTaskIdAll key: task의 Id
type TaskMapByTaskIdAll map[string]Task

// TaskAll 모든 작업 조회 후 parentId를 키로 하는 map과 taskId를 키로 하는 map 2개를 반환
func (w *Client) TaskAll(rootFolderId string) (TaskMapByParentIdAll, TaskMapByTaskIdAll) {
	tasks := Tasks{}
	urlQuery := map[string]string{
		"status":      `["Active","Completed"]`,
		"fields":      `["authorIds","responsibleIds","hasAttachments","parentIds","subTaskIds"]`,
		"descendants": "true",
		"subTasks":    "true",
		"sortField":   `DueDate`,
	}
	w.newAPI("/folders/"+rootFolderId+"/tasks", urlQuery, &tasks)

	// map key: 작업의 부모 ID
	taskMapByParentIdAll := TaskMapByParentIdAll{}
	for _, task := range tasks.Data {
		for _, parentId := range task.ParentIds {
			taskMapByParentIdAll[parentId] = append(taskMapByParentIdAll[parentId], task)
		}
	}

	// map key: 작업의 ID
	taskMapByTaskIdAll := TaskMapByTaskIdAll{}
	for _, task := range tasks.Data {
		taskMapByTaskIdAll[task.ID] = task
	}

	return taskMapByParentIdAll, taskMapByTaskIdAll
}
