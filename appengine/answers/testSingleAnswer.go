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

func (a TestSingleBody) GetHTML(options []string) (template.HTML, template.HTML, error) {

	var doc1, doc2 bytes.Buffer
	unsolvedTmpl := `
	<ul>{{range $index, $item := .}}
        <li><input type="radio" name="RawBody" value="{{$index}}" /><label>{{ $item }}</label></li>
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

	tu, err := template.New("options").Parse(unsolvedTmpl)
	err = tu.Execute(&doc1, options)

	ts, err := template.New("options").Parse(solvedTmpl)
	err = ts.Execute(&doc2, options)

	return template.HTML(doc1.String()), template.HTML(doc2.String()), err
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
