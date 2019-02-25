package docs

import (
	"github.com/pivotal/indicator-protocol/pkg/indicator"
	"log"
	"strings"

	"bytes"
	"fmt"
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
)

var indicatorTmpl = template.Must(template.New("Indicator").Parse(`
<table>
    <tr>
        <th width="25%">Description</th>
        <td>{{.Description}}</td>
    </tr>
    <tr>
        <th>PromQL</th>
        <td>
			<code>{{.PromQL}}</code>
		</td>
    </tr>
    <tr>
		{{if .Thresholds}}
        <th>Thresholds</th>
        <td>
            {{range .Thresholds}} <em>{{.Level}}</em>: {{.Operator}} {{.Value}}<br/> {{end}}
			{{if ne .ThresholdNote ""}}
				{{.ThresholdNote}}
			{{- end}}
        </td>
		{{- end}}
    </tr>
	{{range $key, $value :=.OtherDocumentationFields}}
    	<tr>
    	    <th>{{$key}}</th>
    	    <td>
    	        {{$value}}
    	    </td>
    	</tr>
	{{- end}}
</table>`))

type indicatorPresenter struct {
	indicator.Indicator
}

func NewIndicatorPresenter(i indicator.Indicator) indicatorPresenter {
	return indicatorPresenter{i}
}

func (p *indicatorPresenter) HTML() template.HTML {
	buffer := bytes.NewBuffer(nil)
	err := indicatorTmpl.Execute(buffer, p)

	if err != nil {
		log.Fatalf("could not render indicator: %s", err.Error())
	}

	return template.HTML(buffer.String())
}

func (p indicatorPresenter) PromQL() template.HTML {
	return template.HTML(p.Indicator.PromQL)
}

func (p indicatorPresenter) Title() string {
	t, found := p.Documentation["title"]
	if !found {
		return p.Name
	}

	return t
}

func (p indicatorPresenter) Description() template.HTML {
	return p.markdownDocumentationField("description")
}

func (p indicatorPresenter) ThresholdNote() template.HTML {
	return p.markdownDocumentationField("thresholdNote")
}

func (p indicatorPresenter) OtherDocumentationFields() map[string]template.HTML {
	fields := make(map[string]template.HTML, 0)

	for k, v := range p.Documentation {
		if isUnusedDocumentationField(k) {
			title := strings.Title(getHumanReadableTitle(k))
			fields[title] = template.HTML(blackfriday.Run([]byte(v)))
		}
	}

	return fields
}

func getHumanReadableTitle(titleKey string) string {
	switch titleKey {
	case "recommendedResponse":
		return "Recommended Response"
	default:
		return titleKey
	}
}

func isUnusedDocumentationField(fieldName string) bool {
	return fieldName != "title" && fieldName != "description" && fieldName != "thresholdNote"
}

func (p indicatorPresenter) markdownDocumentationField(field string) template.HTML {
	d, found := p.Documentation[field]
	if !found {
		return ""
	}

	return template.HTML(blackfriday.Run([]byte(d)))
}

type thresholdPresenter struct {
	threshold indicator.Threshold
}

func (p indicatorPresenter) Thresholds() []thresholdPresenter {
	var tp []thresholdPresenter
	for _, t := range p.Indicator.Thresholds {
		tp = append(tp, thresholdPresenter{t})
	}
	return tp
}

func (t thresholdPresenter) Level() string {
	switch t.threshold.Level {
	case "warning":
		return "Yellow warning"
	case "critical":
		return "Red critical"
	default:
		return t.threshold.Level
	}
}

func (t thresholdPresenter) Operator() string {
	return t.threshold.GetComparator()
}

func (t thresholdPresenter) Value() string {
	return fmt.Sprintf("%v", t.threshold.Value)
}
