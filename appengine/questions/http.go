package questions

import (
	"encoding/json"
	"errors"
	//"fmt"

	"app/users"
	"appengine/answers"
	"appengine/srv"
)

// Templates

var listTmpl = "appengine/questions/tmpl/list.html"
var newTmpl = "appengine/questions/tmpl/edit.html"
var viewTmpl = "appengine/questions/tmpl/view.html"
var infoTmpl = "appengine/questions/tmpl/info.html"

func GetList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	err := srv.CheckPerm(wr, users.OP_COMMITTER)
	if err != nil {
		return listTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	wr.R.ParseForm()
	qs, err := getQuestions(wr, wr.R.Form)
	if err != nil {
		return listTmpl, err
	}

	tc["Content"] = qs

	return listTmpl, nil
}

func GetOne(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	var qs []Question

	err := srv.CheckPerm(wr, users.OP_MAKER)
	if err != nil {
		return viewTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	wr.R.ParseForm()
	qs, err = getQuestions(wr, wr.R.Form)
	if len(qs) == 0 || err != nil {
		return viewTmpl, errors.New(ERR_QUESTNOTFOUND)
	}
	q := &qs[0]

	// if the question hasn't got a answer to render. It makes a
	// blank anwser based on Atype of the question to render it

	var blankAnswer *answers.Answer
	if q.Solution == nil {
		blankAnswer = answers.NewAnswerWithBody(-1, -1, q.AType)
	}

	tc["Options"] = blankAnswer.Body.GetHTML(q.Options)
	tc["Content"] = q

	return viewTmpl, nil
}

func GetTagsList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	err := srv.CheckPerm(wr, users.OP_VIEWER)
	if err != nil {
		return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	tags, err := getAllQuestionsTags(wr)
	if err != nil {
		return infoTmpl, err
	}

	tc["Content"] = tags

	return infoTmpl, nil
}

func New(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	tc["ImportForm"] = true
	return newTmpl, nil
}

func Add(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	err := srv.CheckPerm(wr, users.OP_COMMITTER)
	if err != nil {
		return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	var q Question

	decoder := json.NewDecoder(wr.R.Body)
	err = decoder.Decode(&q)
	if err != nil {
		return infoTmpl, err
	}

	err = putQuestion(wr, &q)
	if err != nil {
		return infoTmpl, err
	}

	tc["Content"] = q

	return infoTmpl, nil
}
