package wrike

import (
	"errors"
	"fmt"
	"github.com/cloudflare/ahocorasick"
	"sort"
	"strings"
	"sync"
	"time"
)

type Sprint struct {
	AuthorName string `json:"authorName"`
	Tasks      []Task `json:"tasks"`
	SprintGoal string `json:"sprintGoal"`
}

type SprintWeekly struct {
	Title                string `json:"title"`
	Sprints              []Sprint
	ImportanceStatistics map[string]*ImportanceStatistic
}

type ImportanceStatistic struct {
	Active          int
	Completed       int
	CompletePercent int
	Total           int
	TaskMap         map[string]Task
}

// analyzeImportance 스프린트 데이터의 중요도 별 진행률 분석
func (sw *SprintWeekly) analyzeImportance() {
	sw.ImportanceStatistics["High"] = &ImportanceStatistic{TaskMap: map[string]Task{}}
	sw.ImportanceStatistics["Normal"] = &ImportanceStatistic{TaskMap: map[string]Task{}}

	for _, sprint := range sw.Sprints {
		for _, task := range sprint.Tasks {
			switch task.Importance {
			case "High":
				sw.ImportanceStatistics["High"].TaskMap[task.ID] = task
			case "Normal":
				sw.ImportanceStatistics["Normal"].TaskMap[task.ID] = task
			}
		}
	}

	for _, is := range sw.ImportanceStatistics {
		for _, task := range is.TaskMap {
			switch task.Status {
			case "Active":
				is.Active++
			case "Completed":
				is.Completed++
			}
		}
		is.Total = is.Active + is.Completed
		if is.Total > 0 {
			is.CompletePercent = is.Completed * 100 / is.Total
		}
	}
}

// Sprints 스프린트 데이터 조회
func (w *Client) Sprints(spMonth string, sprintRootLink string, outputDomains []string) ([]SprintWeekly, error) {
	rootProject := w.ProjectsByLink(sprintRootLink, nil)

	// API 호출 제한 때문에 전체를 가져와서 필터링
	var userAll = w.UserAll()
	var attachmentAll = w.AttachmentAll()
	var folderAll = w.FolderAll()

	monthProject, err := findSprintMonthProject(&folderAll, &rootProject, spMonth)
	if err != nil {
		return nil, err
	}

	// 산출물 도메인 필터
	m := ahocorasick.NewStringMatcher(outputDomains)
	outputFilter := func(url string) bool {
		return len(m.Match([]byte(url))) > 0
	}

	// 작업 조회 및 데이터 가공
	tasksAllPerParentId, tasksAllPerTaskId := w.TaskAll(monthProject.ID)
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
				author := userAll.findUser(responsibleId)
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
				attachments := attachmentAll.findByTaskId(task.ID)
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

	// 해당 월의 각 sprint 회차 폴더를 조회한다 (ex. "2022.11.SP1")
	projectsD3 := w.ProjectsByIds(monthProject.ChildIds)

	var sprintWeeklyList []SprintWeekly

	var wg sync.WaitGroup
	var mutex sync.Mutex
	wg.Add(len(projectsD3.Data))

	// sprint 데이터를 가공한다
	for _, folder := range projectsD3.Data {
		fmt.Printf("동기화할 Wrike의 Sprint ==> %s\n", folder.Title)
		go func(p Project) {
			defer wg.Done()

			// 팀원 별 폴더 조회 - 2022.04.SP1.anthony
			foldersPerMember := folderAll.findFolderByIds(p.ChildIds)

			var wgChild sync.WaitGroup
			var mutexChild sync.Mutex
			wgChild.Add(len(foldersPerMember))

			// 팀원 별 프로젝트 하위 작업 조회
			var sprints []Sprint
			for _, pMember := range foldersPerMember {
				go func(pMember Project) {
					defer wgChild.Done()
					mutexChild.Lock()
					defer mutexChild.Unlock()

					authorName := strings.Split(pMember.Title, ".")[3]

					sprints = append(sprints, Sprint{
						AuthorName: authorName,
						Tasks:      findTaskByIds(pMember.ID, authorName),
						SprintGoal: pMember.Description,
					})
				}(pMember)
			}
			wgChild.Wait()

			// 이름 순으로 정렬
			sort.Slice(sprints, func(i, j int) bool { return sprints[i].AuthorName < sprints[j].AuthorName })

			// 1주치 Sprint 구조체 생성
			mutex.Lock()
			defer mutex.Unlock()

			sprintWeekly := SprintWeekly{
				Title:                p.Title,
				Sprints:              sprints,
				ImportanceStatistics: map[string]*ImportanceStatistic{},
			}
			sprintWeekly.analyzeImportance()
			sprintWeeklyList = append(sprintWeeklyList, sprintWeekly)
		}(folder)
	}
	wg.Wait()
	return sprintWeeklyList, nil
}

// findSprintMonthProject 동기화 할 월별 폴더 조회 (ex. "2022년 10월")
func findSprintMonthProject(fa *AllFolderMap, rootProject *Projects, spMonth string) (Project, error) {
	projectsD2 := fa.findFolderByIds(rootProject.Data[0].ChildIds)
	result := Project{}

	for _, p := range projectsD2 {
		if p.Title == spMonth {
			result = p
			break
		}
	}
	if len(result.Title) == 0 {
		msg := fmt.Sprintf("wrike에 [%s] sprint 폴더가 존재하지 않아요\n", spMonth)
		return Project{}, errors.New(msg)
	}
	return result, nil
}
