package main

import (
	"fmt"
	"log"
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
		log.Fatalln("failed to create confluence client:", err.Error())
	}

	wrikeClient, err := wrike.NewWrikeClient(
		os.Getenv("WRIKE_BASE_URL"),
		os.Getenv("WRIKE_TOKEN"),
		os.Getenv("WRIKE_SPACE_ID"),
		nil,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	// 현재 날짜 구하기 (yyyy년 M월)
	loc, _ := time.LoadLocation("Asia/Seoul")
	now := time.Now()
	// 일을 항상 1일로 설정 (8월 31일에 월 +1 하니까 9월이 아니라 10월이 되어버림)
	now = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var wg sync.WaitGroup
	var spMonths []string

	wg.Add(3)
	spMonths = append(spMonths, now.In(loc).AddDate(0, -1, 0).Format("2006년 1월")) // 저번달
	spMonths = append(spMonths, now.In(loc).AddDate(0, 0, 0).Format("2006년 1월"))  // 이번달
	spMonths = append(spMonths, now.In(loc).AddDate(0, 1, 0).Format("2006년 1월"))  // 다음달

	for _, spMonth := range spMonths {
		go func(date string) {
			defer wg.Done()
			// 동기화 실행
			syncConfig := confluence.SyncConfig{
				SpMonth:          date,
				SprintRootLink:   os.Getenv("WRIKE_SPRINT_ROOT_URL"),
				AncestorId:       os.Getenv("CONFLUENCE_ANCESTOR_ID"),
				OutputDomains:    outputDomains,
				ConfluenceDomain: os.Getenv("CONFLUENCE_DOMAIN"),
			}
			errSync := cfClient.SyncContent(syncConfig, wrikeClient)
			if errSync != nil {
				fmt.Println(errSync.Error())
			}
		}(spMonth)
	}
	wg.Wait()
}
