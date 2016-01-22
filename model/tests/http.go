package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"model/answers"
	"model/questions"
	"model/users"

	"appengine/srv"
)

// Templates

//var listTmpl = "appengine/tests/tmpl/list.html"
var newTmpl = "model/tests/tmpl/edit.html"
var viewTmpl = "model/tests/tmpl/view.html"
var doTmpl = "model/tests/tmpl/do.html"
var infoTmpl = "model/tests/tmpl/info.html"

func GetList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	if wr.NU.GetRole() < users.ROLE_STUDENT {
		return viewTmpl, fmt.Errorf("tests: getlist: %v", users.ERR_NOTOPERATIONALLOWED)
	}

	wr.R.ParseForm()

	tests, err := getTests(wr, wr.R.Form)
	if err != nil {
		return infoTmpl, fmt.Errorf("tests: getlist: %v", err)
	}

	tc["Content"] = tests

	return infoTmpl, nil
}

func GetOne(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	if wr.NU.GetRole() < users.ROLE_TEACHER {
		return viewTmpl, fmt.Errorf("tests: getone: %v", users.ERR_NOTOPERATIONALLOWED)
	}
	wr.R.ParseForm()

	ts, err := getTests(wr, wr.R.Form)
	if err != nil {
		return viewTmpl, fmt.Errorf("tests: getone: %v", err)
	}
	t := ts[0]
	if wr.NU.ID() != t.AuthorId {
		return viewTmpl, fmt.Errorf("tests: getone: %s", users.ERR_NOTOPERATIONALLOWED)
	}

	tc["Content"] = t

	return viewTmpl, nil
}

func DoOne(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	if wr.NU.GetRole() < users.ROLE_STUDENT {
		return viewTmpl, fmt.Errorf("tests: doone: %v", users.ERR_NOTOPERATIONALLOWED)
	}
	wr.R.ParseForm()

	ts, err := getTests(wr, wr.R.Form)
	if err != nil {
		return viewTmpl, fmt.Errorf("tests: doone: %v", err)
	}
	t := ts[0]

	// TODO:
	// check if the users is allowed for this exam

	tc["Content"] = t

	return doTmpl, nil
}

func GetTagsList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	if wr.NU.GetRole() < users.ROLE_STUDENT {
		return viewTmpl, fmt.Errorf("tests: gettaglist: %v", users.ERR_NOTOPERATIONALLOWED)
	}

	var tags []string
	var err error
	if wr.NU.GetRole() == users.ROLE_ADMIN {
		tags, err = getAllTestsTags(wr)
	} else {
		tags, err = getTestTagsFromUser(wr, wr.NU.ID())
	}

	if err != nil {
		return infoTmpl, fmt.Errorf("tests: gettaglist: %v", err)
	}

	tc["Content"] = tags

	return infoTmpl, nil
}

func GetUsersList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	if wr.NU.GetRole() < users.ROLE_TEACHER {
		return viewTmpl, fmt.Errorf("tests: getuserlist: %v", users.ERR_NOTOPERATIONALLOWED)
	}
	wr.R.ParseForm()

	ts, err := getTests(wr, wr.R.Form)
	if err != nil {
		return viewTmpl, fmt.Errorf("tests: getuserlist: %v", err)
	}
	t := ts[0]
	if wr.NU.ID() != t.AuthorId {
		return viewTmpl, fmt.Errorf("tests: getuserlist: %s", users.ERR_NOTOPERATIONALLOWED)
	}

	users, err := getUsersAllowed(wr, t)
	if err != nil {
		return infoTmpl, fmt.Errorf("tests: getuserlist: %v", err)
	}

	tc["Content"] = users

	return infoTmpl, nil
}

func GetExercisesList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	if wr.NU.GetRole() < users.ROLE_TEACHER {
		return viewTmpl, fmt.Errorf("tests: getexerciseslist: %v", users.ERR_NOTOPERATIONALLOWED)
	}
	wr.R.ParseForm()

	ts, err := getTests(wr, wr.R.Form)
	if err != nil {
		return viewTmpl, fmt.Errorf("tests: getexerciseslist: %v", err)
	}

	t := ts[0]
	if wr.NU.ID() != t.AuthorId {
		return viewTmpl, fmt.Errorf("tests: getexerciseslist: %s", users.ERR_NOTOPERATIONALLOWED)
	}

	tc["Content"] = t.Exercises

	return infoTmpl, nil
}

func GetExercisesAnswers(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	if wr.NU.GetRole() < users.ROLE_STUDENT {
		return viewTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	ts, err := getTests(wr, wr.R.Form)
	if err != nil {
		return viewTmpl, fmt.Errorf("tests: getexercisesanwsers: %v", err)
	}
	t := ts[0]

	authorId := wr.NU.ID()
	as, err := getExerciseAnswered(wr, t, authorId)
	if err != nil {
		return viewTmpl, fmt.Errorf("tests: getexercisesanwsers: %v", err)
	}

	tc["Content"] = as
	return infoTmpl, nil

}

func DoExercise(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	if wr.NU.GetRole() < users.ROLE_STUDENT {
		return viewTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	var a *answers.Answer

	decoder := json.NewDecoder(wr.R.Body)
	err := decoder.Decode(&a)
	if err != nil {
		return infoTmpl, err
	}

	// Check if the a.AuthorId is an allowed user for
	ex, err := getExerciseById(wr, a.ExerciseId)
	if err != nil {
		return infoTmpl, err
	}
	if ok := IsTestAllowedUser(wr, ex.TestId, a.AuthorId); !ok {
		return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	a.BuildBody()
	a.AuthorId = wr.NU.ID()

	_, err = questions.GetQuestById(wr, a.QuestId)
	if err != nil {
		return infoTmpl, err
	}

	err = answers.PutAnswer(wr, a)
	if err != nil {
		return infoTmpl, err
	}

	tc["Content"] = a
	return infoTmpl, nil
}

func New(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	return newTmpl, nil
}

func Edit(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	if wr.NU.GetRole() < users.ROLE_TEACHER {
		return viewTmpl, fmt.Errorf("tests: edit: %v", users.ERR_NOTOPERATIONALLOWED)
	}

	wr.R.ParseForm()
	ts, err := getTests(wr, wr.R.Form)
	if err != nil {
		return viewTmpl, fmt.Errorf("tests: edit: %v", err)
	}

	t := ts[0]
	if wr.NU.ID() != t.AuthorId {
		return viewTmpl, fmt.Errorf("tests: edit: %s", users.ERR_NOTOPERATIONALLOWED)
	}

	tc["Content"] = t
	tc["FromEditHandler"] = true

	return newTmpl, nil
}

func Add(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	if wr.NU.GetRole() < users.ROLE_TEACHER {
		return viewTmpl, fmt.Errorf("tests: add: %v", users.ERR_NOTOPERATIONALLOWED)
	}

	var t Test
	decoder := json.NewDecoder(wr.R.Body)
	err := decoder.Decode(&t)
	if err != nil {
		return infoTmpl, fmt.Errorf("tests: add: %v", err)
	}

	// clean fields
	t.Desc = strings.Trim(t.Desc, " \t\n")

	err = putTest(wr, &t)
	if err != nil {
		return infoTmpl, fmt.Errorf("tests: add: %v", err)
	}

	tc["Content"] = t

	return infoTmpl, nil
}

func Update(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	if wr.NU.GetRole() < users.ROLE_TEACHER {
		return viewTmpl, fmt.Errorf("tests: update: %s", users.ERR_NOTOPERATIONALLOWED)
	}

	var t Test

	decoder := json.NewDecoder(wr.R.Body)
	err := decoder.Decode(&t)
	if err != nil {
		return infoTmpl, fmt.Errorf("tests: update: %v", err)
	}

	// clean fields
	t.Desc = strings.Trim(t.Desc, " \t\n")

	err = updateTest(wr, &t)
	if err != nil {
		return infoTmpl, fmt.Errorf("tests: update: %v", err)
	}

	tc["Content"] = t

	return infoTmpl, nil
}

func Delete(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	if wr.NU.GetRole() < users.ROLE_TEACHER {
		return viewTmpl, fmt.Errorf("tests: delete: %s", users.ERR_NOTOPERATIONALLOWED)
	}

	wr.R.ParseForm()
	ts, err := getTests(wr, wr.R.Form)
	if err != nil {
		return infoTmpl, fmt.Errorf("tests: delete: %v", err)
	}

	t := ts[0]
	if wr.NU.ID() != t.AuthorId {
		return viewTmpl, fmt.Errorf("tests: delete: %s", users.ERR_NOTOPERATIONALLOWED)
	}

	err = deleteTest(wr, t)
	if err != nil {
		return infoTmpl, fmt.Errorf("tests: delete: %v", err)
	}

	tc["Content"] = "Test borrado con éxito"

	return infoTmpl, nil
}
