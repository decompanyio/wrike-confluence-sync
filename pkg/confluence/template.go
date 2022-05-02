package confluence

import (
	"bytes"
	"html/template"
	"strings"
)

func NewTemplate(dataParam interface{}) string {
	// html 파일 로드
	tmpl, err := template.ParseFiles("sprint-template.html")
	errHandler(err)

	// 데이터를 기반으로 html 템플릿 동적 생성
	var tmplString bytes.Buffer
	err = tmpl.Execute(&tmplString, dataParam)
	errHandler(err)

	// 태그 사이 공백 제거
	// 예시: </th>  <th>  ==> </th><th>
	result := strings.ReplaceAll(tmplString.String(), `/\>\s+\</m`, `><`)
	return result
}
