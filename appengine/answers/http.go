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

	// Del cliente la respuesta ha llegado con los campos (QuestId, AuthorId, AType)
	// Ahora necesitamos crear el body y asociarselo

	var abody AnswerBody
	switch a.BodyType {
	case TYPE_TESTSINGLE:
		abody = NewTestSingleAnswer(-1)
	}

	a.AuthorId = wr.NU.Id
	a.SetBody(abody)

	err = putAnswer(wr, a)

	tc["Content"] = a

	return infoTmpl, nil
}
