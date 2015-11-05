// Single Test Answer
package answers

import (
	"bytes"
	"fmt"
	"html/template"
)

type TestSingleBody struct {
	Id       int64          `json:",string" datastore:"-"`
	atype    AnswerBodyType `json:",string"`
	Solution int
}

func NewTestSingleAnswer(sol int) TestSingleBody {
	return TestSingleBody{-1, TYPE_TESTSINGLE, sol}
}

func (a TestSingleBody) GetId() int64 {
	return a.Id
}

func (a TestSingleBody) GetType() AnswerBodyType {
	return TYPE_TESTSINGLE
}

func (a TestSingleBody) GetHTML(options []string) template.HTML {

	var doc bytes.Buffer
	unsolvedTmpl := `
	<ul>{{range $index, $item := .}}
        <li><input type="radio" name="solution" value="{{$index}}" /><label>{{ $item }}</label></li>
        {{end}}</ul>
`

	solvedTmpl := `
        <ul class="list-group">{{range $index, $item := .}}
        {{if eq $index ` + fmt.Sprintf("%d", a.Solution) + `}}
        <li class="list-group-item list-group-item-success">{{ $item }}</li>
        {{else}}
        <li class="list-group-item">{{ $item }}</li>
        {{end}}
        {{end}}</ul>
`
	var t *template.Template
	if a.Solution < 0 {
		t, _ = template.New("options").Parse(unsolvedTmpl)
	} else {
		t, _ = template.New("options").Parse(solvedTmpl)
	}
	err := t.Execute(&doc, options)
	if err != nil {
		return template.HTML(ERR_BADRENDEREDANSWER)
	} else {
		return template.HTML(doc.String())
	}
}

func (a TestSingleBody) Equals(master AnswerBody) bool {
	if master.GetType() == TYPE_TESTSINGLE {
		sol := master.(TestSingleBody)
		return a.Solution == sol.Solution
	}
	return false
}

func (a TestSingleBody) IsUnsolved() bool {
	return a.Solution == -1
}
