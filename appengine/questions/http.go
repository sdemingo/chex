package questions

import (
	//"encoding/json"
	"errors"
	//"fmt"

	"app/users"
	"appengine/srv"
)

// Templates

var listTmpl = "appengine/questions/tmpl/list.html"
var newTmpl = "appengine/questions/tmpl/edit.html"
var viewTmpl = "appengine/questions/tmpl/view.html"
var infoTmpl = "appengine/questions/tmpl/info.html"

func GetList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	err := srv.CheckPerm(wr, users.OP_COMMIT)
	if err != nil {
		return listTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	wr.R.ParseForm()
	/*
		nus, err := getUsers(wr, wr.R.Form)
		if err != nil {
			return listTmpl, err
		}

		tc["Content"] = nus
	*/

	return listTmpl, nil
}

func New(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	return newTmpl, nil
}
