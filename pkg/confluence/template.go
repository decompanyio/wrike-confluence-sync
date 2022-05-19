package confluence

import (
	"bytes"
	"html/template"
	"strings"
)

func NewTemplate(dataParam interface{}, confluenceDomain string) string {
	// html 템플릿 생성
	html := `
{{- range . -}}
<h3><strong>{{ .AuthorName }}</strong></h3>
{{- .SprintGoal -}}
<table data-layout="wide">
    <colgroup>
        <col style="width: 250.0px;" />
        <col style="width: 75.0px;" />
        <col style="width: 70.0px;" />
        <col style="width: 85.0px;" />
        <col style="width: 85.0px;" />
        <col style="width: 150.0px;" />
    </colgroup>
    <tbody>
        <tr>
            <th><p><strong>제목</strong></p></th>
            <th><p><strong>협업담당자</strong></p></th>
            <th><p><strong>중요</strong></p></th>
            <th><p><strong>작성일</strong></p></th>
            <th><p><strong>기한</strong></p></th>
            <th><p><strong>산출물</strong></p></th>
        </tr>
        {{- range .Tasks -}}
        <tr>
            <td>
                <p><a href="{{ .Permalink }}">
                    {{- if eq .Status "Completed" -}}
                    <del>{{ .Title }}</del>
                    {{- else -}}
                    {{ .Title }}
                    {{- end -}}
                </a></p>
            </td>
            <td><p>{{- range .Coworkers -}}{{ printf "%s " .FirstName }}{{- end -}}</p></td>
            <td><p>{{ .Importance }}</p></td>
            <td><p>{{- if .CreatedDate -}}{{ .CreatedDate.Format "2006-01-02" }}{{- end -}}</p></td>
            <td><p>{{- if .Dates.Due -}}{{ .Dates.Due }}{{- end -}}</p></td>
            <td>
				{{- range .Attachments -}}
                    {{- if .IsDomain "` + confluenceDomain + `" -}}
						<p><a href="{{ .Url }}" data-card-appearance="inline">{{ .Url }}</a></p>
                    {{- else -}}
						<p><a href="{{ .Url }}">{{ .Name }}</a></p>
					{{- end -}}
                {{- end -}}
			</td>
        </tr>
        {{- end -}}
    </tbody>
</table>
<p/>
{{- end -}}`

	// html 템플릿 로드
	tmpl := template.New("confluence-template")
	var err error

	tmpl, err = tmpl.Parse(html)
	errHandler(err)

	// 데이터를 기반으로 html 템플릿 동적 생성
	var tmplString bytes.Buffer
	err = tmpl.Execute(&tmplString, dataParam)
	errHandler(err)

	// 태그 사이 공백 제거
	// 예시: </th>  <th>  ==> </th><th>
	result := strings.ReplaceAll(tmplString.String(), `/\>\s+\</m`, `><`)
	result = escapeSpecialHTML(result)
	return result
}

func escapeSpecialHTML(str string) string {
	str = strings.Replace(str, `&lt;`, `<`, -1)
	str = strings.Replace(str, `&gt;`, `>`, -1)
	str = strings.Replace(str, `&amp;`, `&`, -1)
	return str
}
