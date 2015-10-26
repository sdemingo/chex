package tags

import (
	//"encoding/json"
	"errors"
	//"fmt"

	"app/users"
	"appengine/srv"
)

// Templates

var listTmpl = ""
var newTmpl = ""
var viewTmpl = ""
var infoTmpl = ""

func GetAll(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	err := srv.CheckPerm(wr, users.OP_VIEW)
	if err != nil {
		return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	tags, err := getAllTags(wr)
	if err != nil {
		return infoTmpl, errors.New("Etiquetas no encontradas")
	}

	tc["Content"] = tags

	return listTmpl, nil
}

func Delete(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {

	return infoTmpl, nil
}

func Update(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {

	return "", nil
}

func Add(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {

	return "", nil
}
