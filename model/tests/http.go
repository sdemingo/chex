package tests

import (
	"encoding/json"
	"errors"
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
		return viewTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
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
	if wr.NU.GetRole() < users.ROLE_TEACHER {
		return viewTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}
	wr.R.ParseForm()

	ts, err := getTests(wr, wr.R.Form)
	if len(ts) == 0 || err != nil {
		return viewTmpl, errors.New(ERR_TESTNOTFOUND)
	}
	t := ts[0]

	tc["Content"] = t

	return viewTmpl, nil
}

func GetTagsList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	if wr.NU.GetRole() < users.ROLE_STUDENT {
		return viewTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	tags, err := getAllTestsTags(wr)
	if err != nil {
		return infoTmpl, err
	}

	tc["Content"] = tags

	return infoTmpl, nil
}

func GetUsersList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	if wr.NU.GetRole() < users.ROLE_TEACHER {
		return viewTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}
	wr.R.ParseForm()

	ts, err := getTests(wr, wr.R.Form)
	if len(ts) == 0 || err != nil {
		return viewTmpl, errors.New(ERR_TESTNOTFOUND)
	}
	t := ts[0]

	users, err := getUsersAllowed(wr, t)
	if err != nil {
		return infoTmpl, err
	}

	tc["Content"] = users

	return infoTmpl, nil
}

func New(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	return newTmpl, nil
}

func Add(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	if wr.NU.GetRole() < users.ROLE_TEACHER {
		return viewTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	var t Test

	decoder := json.NewDecoder(wr.R.Body)
	err := decoder.Decode(&t)
	if err != nil {
		return infoTmpl, err
	}

	// clean fields
	t.Desc = strings.Trim(t.Desc, " \t\n")

	err = putTest(wr, &t)
	if err != nil {
		return infoTmpl, err
	}

	tc["Content"] = t

	return infoTmpl, nil
}
