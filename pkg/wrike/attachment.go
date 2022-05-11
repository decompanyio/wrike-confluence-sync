package wrike

import (
	"strings"
	"time"
)

type Attachments struct {
	Kind string       `json:"kind"`
	Data []Attachment `json:"data"`
}

type Attachment struct {
	ID          string    `json:"id"`
	AuthorID    string    `json:"authorId"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"createdDate"`
	Version     int       `json:"version"`
	Type        string    `json:"type"`
	ContentType string    `json:"contentType"`
	Size        int       `json:"size"`
	TaskID      string    `json:"taskId"`
	Url         string    `json:"url,omitempty"`
}

// 작업 ID로 첨부파일 조회
func (w *WrikeClient) AttachmentsByTask(taskId string) Attachments {
	urlQuery := map[string]string{
		"withUrls": `true`,
	}
	attachments := Attachments{}
	w.newAPI("/tasks/"+taskId+"/attachments", urlQuery, &attachments)

	return attachments
}

func (a *Attachment) IsDomain(domain string) bool {
	return strings.Index(a.Url, domain) > -1
}
