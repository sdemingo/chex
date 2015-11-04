// Single Test Answer
package questions

import (
	"bytes"
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
	tmpl := `
	<ul>{{range $index, $item := .}}
        <li><input type="radio" name="solution" value="{{$index}}" /><label>{{ $item }}</label></li>
        {{end}}</ul>
`
	t, err := template.New("options").Parse(tmpl)
	err = t.Execute(&doc, options)
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
