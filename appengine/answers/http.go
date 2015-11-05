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

	var a Answer

	decoder := json.NewDecoder(wr.R.Body)
	err = decoder.Decode(&a)
	if err != nil {
		return infoTmpl, err
	}

	// Del cliente la respuesta ha llegado con los campos (QuestId, AuthorId, AType)
	// Ahora necesitamos crear el body y asociarselo

	// err = putQuestion(wr, &q)
	// if err != nil {
	// 	return infoTmpl, err
	// }

	// tc["Content"] = q

	return infoTmpl, nil
}
