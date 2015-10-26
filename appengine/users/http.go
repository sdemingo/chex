package users

import (
	"encoding/json"
	"errors"
	"fmt"

	"app/users"
	"appengine/srv"
)

// Templates

var listTmpl = "appengine/users/tmpl/list.html"
var newTmpl = "appengine/users/tmpl/edit.html"
var viewTmpl = "appengine/users/tmpl/view.html"
var infoTmpl = "appengine/users/tmpl/info.html"

// it would be deprecated!
func GetAll(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {

	err := srv.CheckPerm(wr, users.OP_ADMIN)
	if err != nil {
		return listTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	filters := make(map[string][]string)
	filters["role"] = []string{fmt.Sprintf("%d", users.ROLE_ADMIN)}
	admins, err := getUsers(wr, filters)

	filters["role"] = []string{fmt.Sprintf("%d", users.ROLE_TEACHER)}
	teachers, err := getUsers(wr, filters)

	filters["role"] = []string{fmt.Sprintf("%d", users.ROLE_STUDENT)}
	students, err := getUsers(wr, filters)

	tc["Admins"] = admins
	tc["Teachers"] = teachers
	tc["Students"] = students

	return listTmpl, nil
}

func GetList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	err := srv.CheckPerm(wr, users.OP_ADMIN)
	if err != nil {
		return listTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	wr.R.ParseForm()
	srv.AppWarning(wr, fmt.Sprintf("%s", wr.R.Form["tags"]))
	nus, err := getUsers(wr, wr.R.Form)
	if err != nil {
		return listTmpl, err
	}

	tc["Content"] = nus

	return listTmpl, nil
}

func GetOne(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {

	err := srv.CheckPerm(wr, users.OP_ADMIN)
	if err != nil {
		return viewTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	wr.R.ParseForm()
	nus, err := getUsers(wr, wr.R.Form)
	if len(nus) == 0 || err != nil {
		return viewTmpl, errors.New("Usuario no encontrado")
	}

	tc["Content"] = nus[0]

	return viewTmpl, nil
}

func GetTagsList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	err := srv.CheckPerm(wr, users.OP_VIEW)
	if err != nil {
		return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	tags, err := getAllUserTags(wr)
	if err != nil {
		return infoTmpl, err
	}

	tc["Content"] = tags

	return infoTmpl, nil
}

func New(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	return newTmpl, nil
}

func Edit(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {

	err := srv.CheckPerm(wr, users.OP_ADMIN)
	if err != nil {
		return newTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	wr.R.ParseForm()
	nus, err := getUsers(wr, wr.R.Form)
	if len(nus) == 0 || err != nil {
		return viewTmpl, errors.New("Usuario no encontrado")
	}

	tc["Content"] = nus[0]

	return newTmpl, nil
}

func Delete(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	err := srv.CheckPerm(wr, users.OP_ADMIN)
	if err != nil {
		return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	wr.R.ParseForm()
	nus, err := getUsers(wr, wr.R.Form)
	if len(nus) == 0 || err != nil {
		return infoTmpl, errors.New("Usuario no encontrado")
	}
	err = deleteUser(wr, nus[0])
	if err != nil {
		return infoTmpl, errors.New("Usuario no encontrado")
	}

	tc["Content"] = "Usuario borrado con éxito"

	return infoTmpl, nil
}

func Update(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {

	err := srv.CheckPerm(wr, users.OP_ADMIN)
	if err != nil {
		return "", errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	var nu users.NUser

	decoder := json.NewDecoder(wr.R.Body)
	err = decoder.Decode(&nu)
	if err != nil {
		return "", err
	}

	err = updateUser(wr, nu)
	if err != nil {
		return "", err
	}

	tc["Content"] = nu

	return "", nil
}

func Add(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	err := srv.CheckPerm(wr, users.OP_ADMIN)
	if err != nil {
		return "", errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	var nu users.NUser

	decoder := json.NewDecoder(wr.R.Body)
	err = decoder.Decode(&nu)
	if err != nil {
		return "", err
	}

	err = putUser(wr, nu)
	if err != nil {
		return "", err
	}

	tc["Content"] = nu

	return "", nil
}
