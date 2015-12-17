package tests

import (
	"encoding/json"
	"errors"

	"app/users"

	"appengine/srv"
)

// Templates

//var listTmpl = "appengine/tests/tmpl/list.html"
var newTmpl = "model/tests/tmpl/edit.html"
var viewTmpl = "model/tests/tmpl/view.html"
var infoTmpl = "model/tests/tmpl/info.html"

func GetList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	err := srv.CheckPerm(wr, users.OP_COMMITTER)
	if err != nil {
		return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	wr.R.ParseForm()

	tests, err := getTests(wr, wr.R.Form)
	if err != nil {
		return infoTmpl, err
	}

	tc["Content"] = tests

	return infoTmpl, nil
}

func GetOne(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {

	err := srv.CheckPerm(wr, users.OP_COMMITTER)
	if err != nil {
		return viewTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	wr.R.ParseForm()
	/*
	 qs, err = getQuestions(wr, wr.R.Form)
	 if len(qs) == 0 || err != nil {
	 return viewTmpl, errors.New(ERR_QUESTNOTFOUND)
	 }
	 q := &qs[0]

	 // if the question hasn't got a answer to render. It makes a
	 // blank anwser based on Atype of the question to render it
	 if q.Solution == nil {
	 q.Solution, err = answers.NewAnswerWithBody(-1, -1, q.AType)
	 if err != nil {
	 return viewTmpl, err
	 }
	 }

	 unSolved, solved, err := q.Solution.Body.GetHTML(q.Options)
	 tc["OptionsSolved"] = solved
	 tc["OptionsUnSolved"] = unSolved
	 tc["Content"] = q
	*/
	return viewTmpl, nil
}

func GetTagsList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	err := srv.CheckPerm(wr, users.OP_VIEWER)
	if err != nil {
		return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	tags, err := getAllTestsTags(wr)
	if err != nil {
		return infoTmpl, err
	}

	tc["Content"] = tags

	return infoTmpl, nil
}

func New(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	return newTmpl, nil
}

func Add(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	err := srv.CheckPerm(wr, users.OP_COMMITTER)
	if err != nil {
		return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	var t Test

	decoder := json.NewDecoder(wr.R.Body)
	err = decoder.Decode(&t)
	if err != nil {
		return infoTmpl, err
	}

	err = putTest(wr, &t)
	if err != nil {
		return infoTmpl, err
	}

	tc["Content"] = t

	return infoTmpl, nil
}
