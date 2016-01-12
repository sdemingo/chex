package tests

import (
	"encoding/json"
	"fmt"
	"strings"

	"model/users"

	"appengine/srv"
)

// Templates

//var listTmpl = "appengine/tests/tmpl/list.html"
var newTmpl = "model/tests/tmpl/edit.html"
var viewTmpl = "model/tests/tmpl/view.html"
var infoTmpl = "model/tests/tmpl/info.html"

func GetList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	if wr.NU.GetRole() < users.ROLE_TEACHER {
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

	tc["Content"] = t

	return viewTmpl, nil
}

func GetTagsList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	if wr.NU.GetRole() < users.ROLE_STUDENT {
		return viewTmpl, fmt.Errorf("tests: gettaglist: %v", users.ERR_NOTOPERATIONALLOWED)
	}

	tags, err := getAllTestsTags(wr)
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

	tc["Content"] = t.Exercises

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
	tc["Content"] = ts[0]
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
		return infoTmpl, fmt.Errorf("tests: add: decode: %v", err)
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
		return infoTmpl, fmt.Errorf("tests: update: decode: %v", err)
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
