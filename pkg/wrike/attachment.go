package wrike

import (
	"fmt"
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

// GetAllAttachments 모든 첨부파일 조회 후 해당 첨부파일의 TaskID가 키인 map 반환
func (w *Client) GetAllAttachments() (AllAttachmentMap, error) {

	// Fetch data from Wrike
	var attachments Attachments

	urlQuery := map[string]string{
		"withUrls": "true",
	}

	resp, err := w.httpClient.R().SetQueryParams(urlQuery).
		SetResult(&attachments).
		Get("/attachments")

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode(), resp.String())
	}
	if err != nil {
		panic(err)

	}

	// TaskID를 키로 하는 map 생성
	attachmentMap := make(AllAttachmentMap)
	for _, attachment := range attachments.Data {
		attachmentMap[attachment.TaskID] = append(attachmentMap[attachment.TaskID], attachment)
	}

	return attachmentMap, nil
}

func (a *Attachment) IsDomain(domain string) bool {
	return strings.Contains(a.Url, domain)
}

func (aam *AllAttachmentMap) findByTaskId(taskId string) []Attachment {
	return (*aam)[taskId]
}
