package wrike

import (
	"sort"
	"strings"
	"sync"
)

type Sprint struct {
	AuthorName string `json:"authorName"`
	Tasks      Tasks  `json:"tasks"`
}

type SprintWeekly struct {
	Title   string `json:"title"`
	Sprints []Sprint
}

// 스프린트 조회 - 스프린트 이름으로 필터
// 파라미터 예시: "2022.03.SP1"
func (w *WrikeClient) Sprints(spMonth string, sprintRootLink string, outputDomains []string) []SprintWeekly {
	// wrike 스프린트 루트 폴더 (Sprint)
	projectsD1 := w.ProjectsByLink(sprintRootLink, nil)

	// 스프린트 2 뎁스 조회 - 2022년 04월 프로젝트
	projectsD2 := w.ProjectsByIds(projectsD1.Data[0].ChildIds)
	projectD2 := Project{}

	for _, p := range projectsD2.Data {
		if p.Title == spMonth {
			projectD2 = p
			break
		}
	}

	// 하위 폴더 조회 - 2022.04.SPX
	projectsD3 := w.ProjectsByIds(projectD2.ChildIds)

	var sprintWeekly []SprintWeekly

	// 비동기 처리
	var wg sync.WaitGroup
	wg.Add(len(projectsD3.Data))

	convertToSprint := func(p Project) {
		// 팀원 별 프로젝트 조회 - 2022.04.SP1.anthony
		foldersPerMember := w.ProjectsByIds(p.ChildIds)
		// 비동기 처리
		var wgChild sync.WaitGroup
		wgChild.Add(len(foldersPerMember.Data))

		// 팀원 별 프로젝트 하위 작업 조회
		var sprints []Sprint
		for _, pMember := range foldersPerMember.Data {
			go func(pMember Project) {
				sprints = append(sprints, Sprint{
					AuthorName: strings.Split(pMember.Title, ".")[3],
					Tasks:      w.TasksInProject(pMember.ID, outputDomains),
				})
				wgChild.Done()
			}(pMember)
		}
		wgChild.Wait()

		// 이름 순으로 정렬
		sort.Slice(sprints, func(i, j int) bool { return sprints[i].AuthorName < sprints[j].AuthorName })

		// 1주치 Sprint 구조체 생성
		sprintWeekly = append(sprintWeekly, SprintWeekly{
			Title:   p.Title,
			Sprints: sprints,
		})
		wg.Done()
	}

	for _, folders := range projectsD3.Data {
		go convertToSprint(folders)
	}

	wg.Wait()

	return sprintWeekly
}
