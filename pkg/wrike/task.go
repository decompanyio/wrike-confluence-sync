package wrike

import (
	"github.com/cloudflare/ahocorasick"
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

// 작업 전체 조회
func (w *WrikeClient) Tasks() Tasks {
	tasks := Tasks{}
	w.newAPI("/tasks", nil, &tasks)

	return tasks
}

// 특정 프로젝트의 작업 조회
func (w *WrikeClient) TasksInProject(folderId string, outputDomains []string) Tasks {
	tasks := Tasks{}
	urlQuery := map[string]string{
		"fields":    `["authorIds","responsibleIds","hasAttachments"]`,
		"sortField": `DueDate`,
	}
	w.newAPI("/folders/"+folderId+"/tasks", urlQuery, &tasks)

	// 산출물 도메인 필터
	m := ahocorasick.NewStringMatcher(outputDomains)
	outputFilter := func(url string) bool {
		return len(m.Match([]byte(url))) > 0
	}

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
		// 첨부파일 조회
		if data.HasAttachments {
			attachments := w.AttachmentsByTask(data.ID)
			for _, attachment := range attachments.Data {
				// 성능을 위해 ahocorasick 알고리즘 사용
				if outputFilter(attachment.Url) {
					tasks.Data[i].Attachments = append(tasks.Data[i].Attachments, attachment)
				}
			}
		}
	}
	return tasks
}
