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

// 작업 전체 조회
func (w *WrikeClient) Tasks() Tasks {
	tasks := Tasks{}
	w.newAPI("/tasks", nil, &tasks)

	return tasks
}

// 특정 프로젝트의 작업 조회
func (w *WrikeClient) TasksInProject(folderId string) Tasks {
	tasks := Tasks{}
	urlQuery := map[string]string{
		"fields":    `["authorIds","responsibleIds"]`,
		"sortField": `["DueDate"]`,
	}
	w.newAPI("/folders/"+folderId+"/tasks", urlQuery, &tasks)
	for i, data := range tasks.Data {
		// 본인 제외 협업담당자
		for _, responsibleId := range data.ResponsibleIds {
			if data.AuthorIds[0] != responsibleId {
				tasks.Data[i].Coworkers = append(tasks.Data[i].Coworkers, w.User(responsibleId))
			}
		}
		// 기한이 이상한 날짜 형식으로 와서 자르기
		if len(data.Dates.Due) > 0 {
			tasks.Data[i].Dates.Due = data.Dates.Due[0:10]
		}
	}
	return tasks
}
