package confluence

import (
	"bytes"
	"html/template"
	"strings"
)

// escapeSpecialHTML 특수 문자 변환
func escapeSpecialHTML(str string) string {
	str = strings.Replace(str, `&lt;`, `<`, -1)
	str = strings.Replace(str, `&gt;`, `>`, -1)
	str = strings.Replace(str, `&amp;`, `&`, -1)
	return str
}

// NewTemplate 템플릿 생성
func NewTemplate(dataParam interface{}, confluenceDomain string) string {
	html := `<ac:structured-macro ac:name="toc" ac:schema-version="1" data-layout="default"><ac:parameter ac:name="minLevel">1</ac:parameter><ac:parameter ac:name="maxLevel">7</ac:parameter><ac:parameter ac:name="type">flat</ac:parameter></ac:structured-macro>

<ac:layout>
	<ac:layout-section ac:type="two_equal" ac:breakout-mode="default">
	{{- range $key, $obj := .ImportanceStatistics -}}
		<ac:layout-cell>
			<ac:structured-macro ac:name="info" ac:schema-version="1">
				<ac:rich-text-body>
					<p>{{- if eq $key "High" -}}<ac:emoticon ac:name="blue-star" ac:emoji-shortname=":exclamation:" ac:emoji-fallback="❗" />{{- end -}}[중요도 {{ $key }}] 진행률: {{ $obj.CompletePercent }}% ({{ $obj.Completed }}/{{ $obj.Total }} 완료)</p>
				</ac:rich-text-body>
			</ac:structured-macro>
			<ac:structured-macro ac:name="expand" ac:schema-version="1">
				<ac:rich-text-body>
					<ul>
					{{- range $taskId, $task := $obj.TaskMap -}}
					<li>
						<p><a href="{{ .Permalink }}">
							{{- if eq $task.Status "Completed" -}}
								<del>{{ $task.Title }}</del>
							{{- else -}}
								{{ $task.Title }}
							{{- end -}}
						</a></p>
					</li>
					{{- end -}}
					</ul>
				</ac:rich-text-body>
			</ac:structured-macro>
		</ac:layout-cell>
	{{- end -}}
	</ac:layout-section>
</ac:layout>

{{- range .Sprints -}}
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
            <td><p>
				{{- if .CreatedDate -}}
					{{ .CreatedDate.Format "2006-01-02" }}
				{{- end -}}
			</p></td>
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
