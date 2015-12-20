package answers

import (
	"encoding/json"
	//"fmt"

	"appengine/srv"
)

// Templates

var infoTmpl = ""

// It must be rename to solve and move to exercises/http.go
func Add(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	var err error
	// err := srv.CheckPerm(wr, users.OP_MAKER)
	// if err != nil {
	// 	return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	// }

	var a *Answer

	decoder := json.NewDecoder(wr.R.Body)
	err = decoder.Decode(&a)
	if err != nil {
		return infoTmpl, err
	}

	a.BuildBody()
	a.AuthorId = wr.NU.ID()

	err = putAnswer(wr, a)
	if err != nil {
		return infoTmpl, err
	}

	tc["Content"] = a
	return infoTmpl, nil
}
