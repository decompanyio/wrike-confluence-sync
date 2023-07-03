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

// GetAllAttachments 모든 첨부파일 조회 후 parentId가 키인 map 반환
func (w *Client) GetAllAttachments() AllAttachmentMap {
	urlQuery := map[string]string{
		"withUrls": "true",
	}

	var attachments Attachments
	w.callAPI("/attachments", urlQuery, &attachments)

	attachmentMap := make(AllAttachmentMap)
	for _, attachment := range attachments.Data {
		attachmentMap[attachment.TaskID] = append(attachmentMap[attachment.TaskID], attachment)
	}

	return attachmentMap
}

func (a *Attachment) IsDomain(domain string) bool {
	return strings.Contains(a.Url, domain)
}

func (aam *AllAttachmentMap) findByTaskId(taskId string) []Attachment {
	return (*aam)[taskId]
}
