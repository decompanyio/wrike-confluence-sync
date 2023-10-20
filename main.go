package main

import (
	"log/slog"
	"os"
	"sync"
	"time"
	"wrike-confluence-sync/pkg/confluence"
	"wrike-confluence-sync/pkg/wrike"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	slog.SetDefault(logger)

	// 동기화 설정
	outputDomains := []string{
		os.Getenv("CONFLUENCE_DOMAIN"),
		"https://www.polarissharetech.net",
		"https://www.figma.com",
		"https://www.polarisoffice.com",
		"https://github.com/decompanyio",
	}

	cfClient, wrikeClient, err := createClients()
	if err != nil {
		slog.Error("failed to create clients", slog.String("error", err.Error()))
		return
	}

	spMonths, err := getSprintMonths()
	if err != nil {
		slog.Error("failed to generate sprint months", slog.String("error", err.Error()))
		return
	}

	// 성능을 위해 모든 데이터를 한번에 조회
	// - wrike API 호출 횟수를 줄이기 위함 (API 호출 횟수 제한이 있음)
	data, err := wrikeClient.GetAllData()
	if err != nil {
		slog.Error("failed to get all data", slog.String("error", err.Error()))
		return
	}

	// 월 별 wrike 프로젝트 데이터를 조회, 가공하여 confluence에 동기화 (비동기)
	var wg sync.WaitGroup
	done := make(chan struct{})
	errCh := make(chan error)

	for _, spMonth := range spMonths {
		wg.Add(1)
		go processSprintProject(spMonth, cfClient, wrikeClient, data, outputDomains, &wg, done, errCh)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	for i := 0; i < len(spMonths); i++ {
		select {
		case <-done:
		case e := <-errCh:
			slog.Error("error occurred", slog.String("error", e.Error()))
		case <-time.After(30 * time.Second):
			slog.Error("root goroutine timeout for 30s")
		}
	}
}

func createClients() (*confluence.Client, *wrike.Client, error) {
	cfClient, err := confluence.NewClient(
		os.Getenv("CONFLUENCE_DOMAIN"),
		os.Getenv("CONFLUENCE_USER"),
		os.Getenv("CONFLUENCE_TOKEN"),
		os.Getenv("CONFLUENCE_SPACEID"),
	)
	if err != nil {
		return nil, nil, err
	}

	wrikeClient, err := wrike.NewClient(
		os.Getenv("WRIKE_BASE_URL"),
		os.Getenv("WRIKE_TOKEN"),
		os.Getenv("WRIKE_SPACE_ID"),
	)
	if err != nil {
		return nil, nil, err
	}

	return cfClient, wrikeClient, nil
}

// getSprintMonths 현재 월을 기준으로 지난달, 이번달, 다음달을 반환
func getSprintMonths() ([]time.Time, error) {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return nil, err
	}

	now := time.Now().In(loc)
	now = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	return []time.Time{
		now.AddDate(0, -1, 0), // last month
		now,                   // this month
		now.AddDate(0, 1, 0),  // next month
	}, nil
}

// processSprintProject wrike 프로젝트를 조회하여 confluence에 동기화
func processSprintProject(date time.Time, cfClient *confluence.Client, wrikeClient *wrike.Client, data wrike.AllData, outputDomains []string, wg *sync.WaitGroup, done chan struct{}, errCh chan error) {
	defer wg.Done()

	sprintProjects, err := wrikeClient.FindSprintProjects(data.ProjectAll, os.Getenv("WRIKE_SPRINT_ROOT_URL"), date.Format("2006.01"))
	if err != nil {
		errCh <- err
		return
	}

	for _, sprintProject := range sprintProjects {
		sprint, err := wrikeClient.Sprint(sprintProject, outputDomains, data)
		if err != nil {
			errCh <- err
			return
		}

		syncConfig := confluence.SyncConfig{
			Date:             date,
			AncestorId:       os.Getenv("CONFLUENCE_ANCESTOR_ID"),
			OutputDomains:    outputDomains,
			ConfluenceDomain: os.Getenv("CONFLUENCE_DOMAIN"),
		}
		err = cfClient.SyncContent(sprint, syncConfig)
		if err != nil {
			errCh <- err
			return
		}
	}

	done <- struct{}{}
}
