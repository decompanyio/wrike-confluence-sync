package confluence

import (
	"bytes"
	"html/template"
	"os"
	"strings"
	"wrike-confluence-sync/pkg/wrike"
)

func NewTemplate() string {
	// html 파일 로드
	tmpl, err := template.ParseFiles("sprint-template.html")
	errHandler(err)

	// wrike 데이터 조회
	wrikeAPI := wrike.NewWrikeClient(os.Getenv("WRIKE_TOKEN"), nil)
	sprints := wrikeAPI.Sprints("2022.03.SP1")

	// wrike 데이터를 기반으로 html 템플릿 동적 생성
	var tmplString bytes.Buffer
	err = tmpl.Execute(&tmplString, sprints)
	errHandler(err)

	// 태그 사이 공백 제거
	// 예시: </th>  <th>  ==> </th><th>
	result := strings.ReplaceAll(tmplString.String(), `/\>\s+\</m`, `><`)
	return result
}
