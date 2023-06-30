package main

import (
	"github.com/rs/zerolog/log"
	"os"
	"sync"
	"time"
	"wrike-confluence-sync/pkg/confluence"
	"wrike-confluence-sync/pkg/wrike"
)

func main() {
	// 매핑할 산출물 도메인 설정
	outputDomains := []string{
		os.Getenv("CONFLUENCE_DOMAIN"),
		"https://www.polarissharetech.net",
		"https://www.figma.com",
		"https://www.polarisoffice.com",
		"https://github.com/decompanyio",
	}

	// 클라이언트 생성
	cfClient, err := confluence.NewConfluenceClient(
		os.Getenv("CONFLUENCE_DOMAIN"),
		os.Getenv("CONFLUENCE_USER"),
		os.Getenv("CONFLUENCE_TOKEN"),
		os.Getenv("CONFLUENCE_SPACEID"),
	)
	if err != nil {
		log.Err(err).Msg("failed to create confluence client")
		return
	}

	wrikeClient, err := wrike.NewWrikeClient(
		os.Getenv("WRIKE_BASE_URL"),
		os.Getenv("WRIKE_TOKEN"),
		os.Getenv("WRIKE_SPACE_ID"),
	)
	if err != nil {
		log.Err(err).Msg("failed to create wrike client:")
		return
	}

	// 현재 날짜 구하기 (yyyy년 M월)
	loc, _ := time.LoadLocation("Asia/Seoul")
	now := time.Now()
	// 일을 항상 1일로 설정 (8월 31일에 월 +1 하니까 9월이 아니라 10월이 되어버림)
	now = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var spMonths []time.Time
	spMonths = append(spMonths, now.In(loc).AddDate(0, -1, 0)) // 저번달
	spMonths = append(spMonths, now.In(loc).AddDate(0, 0, 0))  // 이번달
	spMonths = append(spMonths, now.In(loc).AddDate(0, 1, 0))  // 다음달

	data := wrike.AllData{
		UserAll:       wrikeClient.UserAll(),
		AttachmentAll: wrikeClient.AttachmentAll(),
		ProjectAll:    wrikeClient.ProjectAll(),
	}

	var wg sync.WaitGroup
	done := make(chan struct{})
	errCh := make(chan error)

	for _, spMonth := range spMonths {
		wg.Add(1)
		go func(date time.Time) {
			defer wg.Done()
			// 동기화 실행
			sprintProjects, err := wrikeClient.FindSprintProjects(data.ProjectAll, os.Getenv("WRIKE_SPRINT_ROOT_URL"), date.Format("2006.01"))
			if err != nil {
				errCh <- err
				return
			}

			for _, sprintProject := range sprintProjects {
				sprint, errChild := wrikeClient.Sprint(sprintProject, outputDomains, data)
				if errChild != nil {
					return
				}

				syncConfig := confluence.SyncConfig{
					Date:             date,
					AncestorId:       os.Getenv("CONFLUENCE_ANCESTOR_ID"),
					OutputDomains:    outputDomains,
					ConfluenceDomain: os.Getenv("CONFLUENCE_DOMAIN"),
				}
				errSync := cfClient.SyncContent(sprint, syncConfig)
				if errSync != nil {
					errCh <- errSync
					return
				}
			}
			done <- struct{}{}
		}(spMonth)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	for i := 0; i < len(spMonths); i++ {
		select {
		case <-done:
		case e := <-errCh:
			log.Err(e).Msg("error occurred")
		case <-time.After(30 * time.Second):
			log.Error().Msg("Root goroutine timeout for 30s")
		}
	}
}
