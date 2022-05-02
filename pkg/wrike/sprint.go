package wrike

import (
	"fmt"
	"strconv"
	"strings"
)

type Sprint struct {
	AuthorName string `json:"title"`
	Tasks      Tasks  `json:"tasks"`
}

// 스프린트 조회 - 스프린트 이름으로 필터
// 파라미터 예시: "2022.03.SP1"
func (w *WrikeClient) Sprints(spName string) []Sprint {
	// wrike 스프린트 루트 프로젝트 (Sprint)
	projectsD1 := w.ProjectsByLink("https://www.wrike.com/open.htm?id=865199939", nil)

	// 스프린트 2 뎁스 조회 (xxxx년 xx월)
	projectsD2 := w.ProjectsByIds(projectsD1.Data[0].ChildIds)
	projectD2 := Project{}

	// 파라미터 split
	spNameSlice := strings.Split(spName, ".")
	year, err := strconv.Atoi(spNameSlice[0])
	errorHandler(err)
	month, err := strconv.Atoi(spNameSlice[1])
	errorHandler(err)

	// 스프린트 이름으로 필터 (xxxx년 xx월)
	for _, p := range projectsD2.Data {
		if p.Title == fmt.Sprintf("%d년 %d월", year, month) {
			projectD2 = p
			break
		}
	}

	// 하위 폴더 조회 (xxxx년 xx월 / xxxx.xx.SPx)
	projectsD3 := w.ProjectsByIds(projectD2.ChildIds)
	projectD3 := Project{}
	for _, p := range projectsD3.Data {
		if p.Title == spName {
			projectD3 = p
			break
		}
	}

	// 팀원 별 프로젝트 조회
	projectsD4 := w.ProjectsByIds(projectD3.ChildIds)

	// 팀원 별 프로젝트 하위 작업 조회
	sprints := []Sprint{}

	for _, p := range projectsD4.Data {
		sprints = append(sprints, Sprint{
			AuthorName: strings.Split(p.Title, ".")[3],
			Tasks:      w.TasksInProject(p.ID, true),
		})
	}

	return sprints
}
