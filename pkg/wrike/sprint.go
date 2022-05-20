package wrike

import (
	"fmt"
	"github.com/cloudflare/ahocorasick"
	"sort"
	"strings"
	"sync"
)

type Sprint struct {
	AuthorName string `json:"authorName"`
	Tasks      []Task `json:"tasks"`
	SprintGoal string `json:"sprintGoal"`
}

type SprintWeekly struct {
	Title   string `json:"title"`
	Sprints []Sprint
}

// 스프린트 조회 - 스프린트 이름으로 필터
// 파라미터 예시: "2022.03.SP1"
func (w *WrikeClient) Sprints(spMonth string, sprintRootLink string, outputDomains []string) []SprintWeekly {
	// wrike 스프린트 루트 폴더 (Sprint)
	rootProject := w.ProjectsByLink(sprintRootLink, nil)

	// API 호출 제한 때문에 전체를 가져와서 필터링
	var userAll = w.UserAll()
	findUserById := func(userId string) User {
		return userAll[userId]
	}

	var attachmentAll = w.AttachmentAll()
	findAttachmentByTaskId := func(taskId string) []Attachment {
		return attachmentAll[taskId]
	}

	var folderAll = w.FolderAll()
	findFolderByIds := func(folderIds []string) []Project {
		var projectTemp []Project
		for _, id := range folderIds {
			projectTemp = append(projectTemp, folderAll[id])
		}
		return projectTemp
	}

	// 스프린트 2 뎁스 조회 - 2022년 04월
	projectsD2 := findFolderByIds(rootProject.Data[0].ChildIds)
	projectD2 := Project{}
	for _, p := range projectsD2 {
		if p.Title == spMonth {
			projectD2 = p
			break
		}
	}

	// 산출물 도메인 필터
	m := ahocorasick.NewStringMatcher(outputDomains)
	outputFilter := func(url string) bool {
		return len(m.Match([]byte(url))) > 0
	}

	tasksAllPerParentId, tasksAllPerTaskId := w.TaskAll(projectD2.ID)
	findTaskByIds := func(parentId string, authorName string) []Task {
		var taskTemp = tasksAllPerParentId[parentId]
		// 하위 작업 조회
		for _, task := range taskTemp {
			for _, subTaskId := range task.SubTaskIds {
				taskTemp = append(taskTemp, tasksAllPerTaskId[subTaskId])
			}
		}
		for i, task := range taskTemp {
			// 본인 제외 협업담당자
			for _, responsibleId := range task.ResponsibleIds {
				author := findUserById(responsibleId)
				if strings.ToLower(authorName) != strings.ToLower(author.FirstName) {
					taskTemp[i].Coworkers = append(taskTemp[i].Coworkers, author)
				}
			}
			// 기한이 이상한 날짜 형식으로 와서 자르기
			if len(task.Dates.Due) > 0 {
				taskTemp[i].Dates.Due = task.Dates.Due[0:10]
			}
			// 첨부파일 조회
			if task.HasAttachments {
				attachments := findAttachmentByTaskId(task.ID)
				for _, attachment := range attachments {
					// 성능을 위해 ahocorasick 알고리즘 사용
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

	// 하위 폴더 조회 - 2022.04.SPX
	projectsD3 := w.ProjectsByIds(projectD2.ChildIds)

	var sprintWeekly []SprintWeekly
	// 비동기 처리
	var wg sync.WaitGroup
	var mutex sync.Mutex
	wg.Add(len(projectsD3.Data))

	convertToSprint := func(p Project) {
		// 팀원 별 폴더 조회 - 2022.04.SP1.anthony
		foldersPerMember := findFolderByIds(p.ChildIds)

		// 비동기 처리
		var wgChild sync.WaitGroup
		var mutexChild sync.Mutex
		wgChild.Add(len(foldersPerMember))

		// 팀원 별 프로젝트 하위 작업 조회
		var sprints []Sprint
		for _, pMember := range foldersPerMember {
			go func(pMember Project) {
				authorName := strings.Split(pMember.Title, ".")[3]
				mutexChild.Lock()
				sprints = append(sprints, Sprint{
					AuthorName: authorName,
					Tasks:      findTaskByIds(pMember.ID, authorName),
					SprintGoal: pMember.Description,
				})
				mutexChild.Unlock()
				wgChild.Done()
			}(pMember)
		}
		wgChild.Wait()

		// 이름 순으로 정렬
		sort.Slice(sprints, func(i, j int) bool { return sprints[i].AuthorName < sprints[j].AuthorName })

		// 1주치 Sprint 구조체 생성
		mutex.Lock()
		sprintWeekly = append(sprintWeekly, SprintWeekly{
			Title:   p.Title,
			Sprints: sprints,
		})
		mutex.Unlock()
		wg.Done()
	}
	for _, folders := range projectsD3.Data {
		fmt.Printf("동기화할 Wrike의 Sprint ==> %s\n", folders.Title)
		go convertToSprint(folders)
	}
	wg.Wait()

	return sprintWeekly
}
