// Single Test Answer
package answers

import (
	"bytes"
	"html/template"
)

type TestSingleBody struct {
	Id       int64          `json:",string" datastore:"-"`
	atype    AnswerBodyType `json:",string"`
	Solution int
}

func NewTestSingleAnswer(sol int) *TestSingleBody {
	t := new(TestSingleBody)
	t.Id = -1
	t.atype = TYPE_TESTSINGLE
	t.Solution = sol
	//return &TestSingleBody{-1, TYPE_TESTSINGLE, sol}
	return t
}

func (a TestSingleBody) ID() int64 {
	return a.Id
}

func (a *TestSingleBody) SetID(id int64) {
	a.Id = id
}

func (a TestSingleBody) GetType() AnswerBodyType {
	return TYPE_TESTSINGLE
}

func (a TestSingleBody) GetHTML(options []string) (template.HTML, error) {

	var doc bytes.Buffer
	unsolvedTmpl := `
	<ul>{{range $index, $item := .}}
        <li><input type="radio" name="RawBody" value="{{$index}}" /><label>{{ $item }}</label></li>
        {{end}}</ul>
`

	tu, err := template.New("options").Parse(unsolvedTmpl)
	err = tu.Execute(&doc, options)

	return template.HTML(doc.String()), err
}

func (a TestSingleBody) Equals(master AnswerBody) bool {
	if master.GetType() == TYPE_TESTSINGLE {
		sol := master.(*TestSingleBody)
		return a.Solution == sol.Solution
	}
	return false
}

func (a TestSingleBody) IsUnsolved() bool {
	return a.Solution == -1
}
