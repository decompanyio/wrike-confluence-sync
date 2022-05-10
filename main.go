package main

import (
	"log"
	"os"
	"time"
	"wrike-confluence-sync/pkg/confluence"
	"wrike-confluence-sync/pkg/wrike"
)

var (
	CONFLUENCE_DOMAIN      string
	CONFLUENCE_USER        string
	CONFLUENCE_TOKEN       string
	CONFLUENCE_SPACEID     string
	CONFLUENCE_ANCESTOR_ID string
	WRIKE_BASE_URL         string
	WRIKE_TOKEN            string
	WRIKE_SPRINT_ROOT_URL  string
)

var (
	cfClient    *confluence.ConfluenceClient
	wrikeClient *wrike.WrikeClient
)

func main() {
	// 환경변수 설정
	configure()

	// 현재 날짜 구하기 (yyyy년 M월)
	loc, _ := time.LoadLocation("Asia/Seoul")
	now := time.Now().In(loc).Format("2006년 1월")

	// 클라이언트 생성
	cfClient = confluence.NewConfluenceClient(CONFLUENCE_DOMAIN, CONFLUENCE_USER, CONFLUENCE_TOKEN, CONFLUENCE_SPACEID)
	wrikeClient = wrike.NewWrikeClient(WRIKE_BASE_URL, WRIKE_TOKEN, nil)

	// 싱크 실행
	syncConfig := confluence.SyncConfig{
		SpMonth:          now,
		SprintRootLink:   WRIKE_SPRINT_ROOT_URL,
		WrikeBaseUrl:     WRIKE_BASE_URL,
		WrikeToken:       WRIKE_TOKEN,
		AncestorId:       CONFLUENCE_ANCESTOR_ID,
		ConfluenceDomain: CONFLUENCE_DOMAIN,
	}
	cfClient.SyncContent(syncConfig)
}

func configure() {
	CONFLUENCE_DOMAIN = os.Getenv("CONFLUENCE_DOMAIN")
	CONFLUENCE_USER = os.Getenv("CONFLUENCE_USER")
	CONFLUENCE_TOKEN = os.Getenv("CONFLUENCE_TOKEN")
	CONFLUENCE_DOMAIN = os.Getenv("CONFLUENCE_DOMAIN")
	CONFLUENCE_SPACEID = os.Getenv("CONFLUENCE_SPACEID")
	CONFLUENCE_ANCESTOR_ID = os.Getenv("CONFLUENCE_ANCESTOR_ID")
	WRIKE_BASE_URL = os.Getenv("WRIKE_BASE_URL")
	WRIKE_TOKEN = os.Getenv("WRIKE_TOKEN")
	WRIKE_SPRINT_ROOT_URL = os.Getenv("WRIKE_SPRINT_ROOT_URL")
}

func errHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
