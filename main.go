package main

import (
	"os"
	"time"
	"wrike-confluence-sync/pkg/confluence"
)

var (
	cfClient *confluence.Client
)

func main() {
	// 환경변수 설정
	configure()

	// 현재 날짜 구하기 (yyyy년 M월)
	loc, _ := time.LoadLocation("Asia/Seoul")
	now := time.Now().In(loc).Format("2006년 1월")

	// 클라이언트 생성
	cfClient = confluence.NewConfluenceClient(
		CONFLUENCE_DOMAIN,
		CONFLUENCE_USER,
		CONFLUENCE_TOKEN,
		CONFLUENCE_SPACEID,
	)

	// 매핑할 산출물 도메인 설정
	outputDomains := []string{
		CONFLUENCE_DOMAIN,
		"https://www.polarissharetech.net",
		"https://www.figma.com",
		"https://www.polarisoffice.com",
	}

	// 동기화 실행
	syncConfig := confluence.SyncConfig{
		SpMonth:          now,
		SprintRootLink:   WRIKE_SPRINT_ROOT_URL,
		WrikeBaseUrl:     WRIKE_BASE_URL,
		WrikeToken:       WRIKE_TOKEN,
		WrikeSpaceId:     WRIKE_SPACE_ID,
		AncestorId:       CONFLUENCE_ANCESTOR_ID,
		OutputDomains:    outputDomains,
		ConfluenceDomain: CONFLUENCE_DOMAIN,
	}
	cfClient.SyncContent(syncConfig)
}

var (
	CONFLUENCE_DOMAIN      string
	CONFLUENCE_USER        string
	CONFLUENCE_TOKEN       string
	CONFLUENCE_SPACEID     string
	CONFLUENCE_ANCESTOR_ID string
	WRIKE_BASE_URL         string
	WRIKE_TOKEN            string
	WRIKE_SPACE_ID         string
	WRIKE_SPRINT_ROOT_URL  string
)

func configure() {
	CONFLUENCE_DOMAIN = os.Getenv("CONFLUENCE_DOMAIN")
	CONFLUENCE_USER = os.Getenv("CONFLUENCE_USER")
	CONFLUENCE_TOKEN = os.Getenv("CONFLUENCE_TOKEN")
	CONFLUENCE_SPACEID = os.Getenv("CONFLUENCE_SPACEID")
	CONFLUENCE_ANCESTOR_ID = os.Getenv("CONFLUENCE_ANCESTOR_ID")
	WRIKE_BASE_URL = os.Getenv("WRIKE_BASE_URL")
	WRIKE_TOKEN = os.Getenv("WRIKE_TOKEN")
	WRIKE_SPACE_ID = os.Getenv("WRIKE_SPACE_ID")
	WRIKE_SPRINT_ROOT_URL = os.Getenv("WRIKE_SPRINT_ROOT_URL")
}
