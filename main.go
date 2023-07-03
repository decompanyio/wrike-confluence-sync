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
	outputDomains := []string{
		"https://www.polarissharetech.net",
		"https://www.figma.com",
		"https://www.polarisoffice.com",
		"https://github.com/decompanyio",
	}

	cfClient, wrikeClient, err := createClients()
	if err != nil {
		log.Err(err).Msg("failed to create clients")
		return
	}

	spMonths, err := getSprintMonths()
	if err != nil {
		log.Err(err).Msg("failed to generate sprint months")
		return
	}

	data := wrike.AllData{
		UserAll:       wrikeClient.GetAllUsers(),
		AttachmentAll: wrikeClient.GetAllAttachments(),
		ProjectAll:    wrikeClient.GetAllProjects(),
	}

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
			log.Err(e).Msg("error occurred")
		case <-time.After(30 * time.Second):
			log.Error().Msg("Root goroutine timeout for 30s")
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
