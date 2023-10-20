package wrike

import (
	"fmt"
	"github.com/cloudflare/ahocorasick"
	"log/slog"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	ImportanceHigh   = "High"
	ImportanceNormal = "Normal"
	StatusActive     = "Active"
	StatusCompleted  = "Completed"
)

type Sprint struct {
	AuthorName string `json:"authorName"`
	Tasks      []Task `json:"tasks"`
	SprintGoal string `json:"sprintGoal"`
}

// SprintWeekly 스프린트 회차 단위 데이터
type SprintWeekly struct {
	Title                string `json:"title"`
	Sprints              []Sprint
	ImportanceStatistics map[string]*ImportanceStatistic
}

// ImportanceStatistic 중요도 별 진행률 통계
type ImportanceStatistic struct {
	Active          int
	Completed       int
	CompletePercent int
	Total           int
	TaskMap         map[string]Task
}

// analyzeImportance 스프린트 데이터의 중요도 별 진행률 분석
func (sw *SprintWeekly) analyzeImportance() SprintWeekly {
	sw.ImportanceStatistics[ImportanceHigh] = &ImportanceStatistic{TaskMap: map[string]Task{}}
	sw.ImportanceStatistics[ImportanceNormal] = &ImportanceStatistic{TaskMap: map[string]Task{}}

	// 중요도 별 작업 분류
	for _, sprint := range sw.Sprints {
		for _, task := range sprint.Tasks {
			switch task.Importance {
			case ImportanceHigh:
				sw.ImportanceStatistics[ImportanceHigh].TaskMap[task.ID] = task
			case ImportanceNormal:
				sw.ImportanceStatistics[ImportanceNormal].TaskMap[task.ID] = task
			}
		}
	}

	// 중요도 별 진행률 계산
	for _, is := range sw.ImportanceStatistics {
		for _, task := range is.TaskMap {
			switch task.Status {
			case StatusActive:
				is.Active++
			case StatusCompleted:
				is.Completed++
			}
		}
		is.Total = is.Active + is.Completed
		if is.Total > 0 {
			is.CompletePercent = is.Completed * 100 / is.Total
		}
	}
	return *sw
}

// Sprint 스프린트 데이터 조회
func (w *Client) Sprint(sprintProject Project, outputDomains []string, data AllData) (SprintWeekly, error) {
	fmt.Printf("동기화할 Wrike의 Sprint ==> %s\n", sprintProject.Title)

	// 산출물 도메인 필터
	m := ahocorasick.NewStringMatcher(outputDomains)
	outputFilter := func(url string) bool {
		return len(m.Match([]byte(url))) > 0
	}

	// 작업 조회 및 데이터 가공을 위한 익명 함수
	tasksAllPerParentId, tasksAllPerTaskId, err := w.TaskAll(sprintProject.ID)
	if err != nil {
		return SprintWeekly{}, err
	}

	findTaskByIds := func(parentId string, authorName string) []Task {
		var taskTemp = tasksAllPerParentId[parentId]
		// sprint에 폴더 형태로 등록했을 경우, 하위 작업을 조회하여 포함한다
		for _, task := range taskTemp {
			for _, subTaskId := range task.SubTaskIds {
				taskTemp = append(taskTemp, tasksAllPerTaskId[subTaskId])
			}
		}
		for i, task := range taskTemp {
			// 본인을 제외한 협업담당자를 설정한다
			for _, responsibleId := range task.ResponsibleIds {
				author := data.UserAll.findUser(responsibleId)
				if strings.ToLower(authorName) != strings.ToLower(author.FirstName) {
					taskTemp[i].Coworkers = append(taskTemp[i].Coworkers, author)
				}
			}
			// 기한 날짜 파싱해서 yyyy-MM-dd 포맷으로 변경한다
			if len(task.Dates.Due) > 0 {
				parse, _ := time.Parse("2006-01-02T15:04:05", task.Dates.Due)
				taskTemp[i].Dates.Due = parse.Format("2006-01-02")
			}
			// 각 작업의 첨부파일을 조회하여 산출물로 설정한다
			if task.HasAttachments {
				attachments := data.AttachmentAll.findByTaskId(task.ID)
				for _, attachment := range attachments {
					// 도메인 필터
					if outputFilter(attachment.Url) {
						taskTemp[i].Attachments = append(taskTemp[i].Attachments, attachment)
					}
				}
			}
		}
		// 작업을 기한(오름차순), 이름(오름차순) 순으로 정렬
		sort.Slice(taskTemp, func(i, j int) bool {
			switch strings.Compare(taskTemp[i].Dates.Due, taskTemp[j].Dates.Due) {
			case -1:
				return true
			case 1:
				return false
			}
			return taskTemp[i].Title < taskTemp[j].Title
		})
		return taskTemp
	}

	// sprint 데이터를 가공한다
	// 팀원 별 프로젝트 하위 작업 조회
	var wg sync.WaitGroup
	var mutex sync.Mutex
	done := make(chan struct{})

	var sprints []Sprint
	foldersPerMember := data.ProjectAll.GetProjectsByIds(sprintProject.ChildIds)
	for _, pMember := range foldersPerMember {
		wg.Add(1)
		go func(pMember Project) {
			defer wg.Done()

			sprintTitleSlice := strings.Split(pMember.Title, ".")
			authorName := sprintTitleSlice[len(sprintTitleSlice)-1]

			mutex.Lock()
			sprints = append(sprints, Sprint{
				AuthorName: authorName,
				Tasks:      findTaskByIds(pMember.ID, authorName),
				SprintGoal: pMember.Description,
			})
			mutex.Unlock()

			done <- struct{}{}
		}(pMember)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	for i := 0; i < len(foldersPerMember); i++ {
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			slog.Error("[wrike] spring api timeout for 5s")
		}
	}

	// 이름 순으로 정렬
	sort.Slice(sprints, func(i, j int) bool { return sprints[i].AuthorName < sprints[j].AuthorName })

	sprintWeekly := SprintWeekly{
		Title:                sprintProject.Title,
		Sprints:              sprints,
		ImportanceStatistics: map[string]*ImportanceStatistic{},
	}
	return sprintWeekly.analyzeImportance(), nil
}

// FindSprintProjects 동기화 할 월별 폴더 조회 (ex. month: "2022.10.SP1")
// return []Project, error (Project의 Title은 2023.04.SP1 형식)
func (w *Client) FindSprintProjects(fa AllProjectMap, rootProjectLink string, month string) ([]Project, error) {
	var result []Project

	rootProject, err := w.GetProjectsByLink(rootProjectLink, nil)
	if err != nil {
		return nil, err
	}

	projects := fa.GetProjectsByIds(rootProject.Data[0].ChildIds)

	for _, p := range projects {
		if strings.HasPrefix(p.Title, month) {
			result = append(result, p)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("wrike에 [%s] sprint 폴더가 존재하지 않아요\n", month)
	}
	return result, nil
}
