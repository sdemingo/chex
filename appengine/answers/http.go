package answers

import (
	"encoding/json"
	"errors"
	//"fmt"

	"app/users"
	"appengine/srv"
)

// Templates

var infoTmpl = ""

func Add(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	err := srv.CheckPerm(wr, users.OP_MAKER)
	if err != nil {
		return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	var a *Answer

	decoder := json.NewDecoder(wr.R.Body)
	err = decoder.Decode(&a)
	if err != nil {
		return infoTmpl, err
	}

	a.BuildBody()
	a.AuthorId = wr.NU.Id

	err = putAnswer(wr, a)
	if err != nil {
		return infoTmpl, err
	}

	tc["Content"] = a
	return infoTmpl, nil
}
