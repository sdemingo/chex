package questions

import (
	"encoding/json"
	"errors"

	"app/users"
	"model/answers"

	"appengine/srv"
)

// Templates

var newTmpl = "model/questions/tmpl/edit.html"
var viewTmpl = "model/questions/tmpl/view.html"
var infoTmpl = "model/questions/tmpl/info.html"
var mainTmpl = "model/questions/tmpl/main.html"

func Main(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	return mainTmpl, nil
}

func GetList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	err := srv.CheckPerm(wr, users.OP_COMMITTER)
	if err != nil {
		return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	wr.R.ParseForm()
	qs, err := getQuestions(wr, wr.R.Form)
	if err != nil {
		return infoTmpl, err
	}

	tc["Content"] = qs

	return infoTmpl, nil
}

func GetOne(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {

	// only teacher must entry to the question throught this
	// handler. A student or other should use test handlers
	err := srv.CheckPerm(wr, users.OP_COMMITTER)
	if err != nil {
		return viewTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	wr.R.ParseForm()
	qs, err := getQuestions(wr, wr.R.Form)
	if len(qs) == 0 || err != nil {
		return viewTmpl, errors.New(ERR_QUESTNOTFOUND)
	}
	q := qs[0]

	// A question only can be viewed by and admin or by their writer
	if !wr.NU.IsAdmin() && q.AuthorId != wr.NU.Id {
		return viewTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	// if the question hasn't got a answer to render. It makes a
	// blank anwser based on Atype of the question to render it
	if !q.IsSolved() {
		q.SolutionId = -1
		q.Solution, err = answers.NewAnswerWithBody(-1, -1, q.AType)
		if err != nil {
			return viewTmpl, err
		}
	}

	unSolved, solved, err := q.Solution.Body.GetHTML(q.Options)
	tc["OptionsSolved"] = solved
	tc["OptionsUnSolved"] = unSolved
	tc["UnSolvedQuestion"] = !q.IsSolved()
	tc["Content"] = q

	return viewTmpl, nil
}

func GetTagsList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	err := srv.CheckPerm(wr, users.OP_VIEWER)
	if err != nil {
		return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	var tags []string
	if wr.NU.IsAdmin() {
		tags, err = getAllQuestionsTags(wr)
	} else {
		tags, err = getQuestionsTagsFromUser(wr, wr.NU.Id)
	}

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
