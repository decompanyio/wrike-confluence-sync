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

type AllAttachmentMap map[string][]Attachment

// AttachmentAll 모든 첨부파일 조회 후 parentId를 키인 map 반환
func (w *Client) AttachmentAll() AllAttachmentMap {
	urlQuery := map[string]string{
		"withUrls": `true`,
	}
	attachments := Attachments{}
	w.newAPI("/attachments", urlQuery, &attachments)

	attachmentAll := AllAttachmentMap{}
	for _, attachment := range attachments.Data {
		attachmentAll[attachment.TaskID] = append(attachmentAll[attachment.TaskID], attachment)
	}

	return attachmentAll
}

func (a *Attachment) IsDomain(domain string) bool {
	return strings.Index(a.Url, domain) > -1
}

func (aam *AllAttachmentMap) findByTaskId(taskId string) []Attachment {
	return (*aam)[taskId]
}
