package users

import (
	"app/users"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"appengine/srv"
)

// Templates

var newTmpl = "model/users/tmpl/edit.html"
var viewTmpl = "model/users/tmpl/view.html"
var infoTmpl = "model/users/tmpl/info.html"
var mainTmpl = "model/users/tmpl/main.html"

func Main(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	return mainTmpl, nil
}

func GetList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	err := srv.CheckPerm(wr, users.OP_COMMITTER)
	if err != nil {
		return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	wr.R.ParseForm()
	nus, err := getUsers(wr, wr.R.Form)
	if err != nil {
		return infoTmpl, err
	}

	tc["Content"] = nus

	return infoTmpl, nil
}

func GetOne(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {

	var nus []*users.NUser

	if strings.HasSuffix(wr.R.URL.Path, "/me") {
		if wr.NU.Id > 0 {
			filters := map[string][]string{"id": []string{fmt.Sprintf("%d", wr.NU.Id)}}
			nus, err := getUsers(wr, filters)
			if len(nus) == 0 || err != nil {
				return viewTmpl, errors.New(users.ERR_USERNOTFOUND)
			}
			tc["Content"] = nus[0]
		} else {
			tc["Content"] = wr.NU
		}

	} else {
		err := srv.CheckPerm(wr, users.OP_ADMIN)
		if err != nil {
			return viewTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
		}

		wr.R.ParseForm()
		nus, err = getUsers(wr, wr.R.Form)
		if len(nus) == 0 || err != nil {
			return viewTmpl, errors.New(users.ERR_USERNOTFOUND)
		}
		tc["Content"] = nus[0]
		tc["UserProfile"] = true
	}

	return viewTmpl, nil
}

func GetTagsList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	err := srv.CheckPerm(wr, users.OP_VIEWER)
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
	tc["ImportForm"] = true
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
		return viewTmpl, errors.New(users.ERR_USERNOTFOUND)
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
		return infoTmpl, errors.New(users.ERR_USERNOTFOUND)
	}
	err = deleteUser(wr, nus[0])
	if err != nil {
		return infoTmpl, errors.New(users.ERR_USERNOTFOUND)
	}

	tc["Content"] = "Usuario borrado con éxito"

	return infoTmpl, nil
}

func Update(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {

	err := srv.CheckPerm(wr, users.OP_ADMIN)
	if err != nil {
		return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	nu := new(users.NUser)

	decoder := json.NewDecoder(wr.R.Body)
	err = decoder.Decode(nu)
	if err != nil {
		return infoTmpl, err
	}

	err = updateUser(wr, nu)
	if err != nil {
		return infoTmpl, err
	}

	tc["Content"] = nu

	return infoTmpl, nil
}

func Add(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	err := srv.CheckPerm(wr, users.OP_ADMIN)
	if err != nil {
		return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	nu := new(users.NUser)

	decoder := json.NewDecoder(wr.R.Body)
	err = decoder.Decode(nu)
	if err != nil {
		return infoTmpl, err
	}

	err = putUser(wr, nu)
	if err != nil {
		return infoTmpl, err
	}

	tc["Content"] = nu

	return infoTmpl, nil
}

func Import(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	err := srv.CheckPerm(wr, users.OP_ADMIN)
	if err != nil {
		return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	file, _, err := wr.R.FormFile("importFile")
	if err != nil {
		return infoTmpl, err
	}

	//var nus []users.NUser
	nus := NewNUserBuffer()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&nus)
	if err != nil {
		return infoTmpl, err
	}

	for _, nu := range nus {
		err = nu.IsValid()
		if err != nil {
			return infoTmpl, err
		}
	}

	for _, nu := range nus {
		err = putUser(wr, nu)
		if err != nil {
			return infoTmpl, err
		}
	}

	tc["Content"] = fmt.Sprintf("Se han creado en la base de datos %d usuarios", len(nus))

	return infoTmpl, nil
}
