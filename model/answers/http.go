package answers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"model/users"

	"appengine/srv"
)

// Templates

var infoTmpl = ""

func Add(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	if wr.NU.GetRole() < users.ROLE_STUDENT {
		return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	}

	var a *Answer

	decoder := json.NewDecoder(wr.R.Body)
	err := decoder.Decode(&a)
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

func GetList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	if wr.NU.GetRole() < users.ROLE_STUDENT {
		return infoTmpl, fmt.Errorf("answers: getlist: %v", users.ERR_NOTOPERATIONALLOWED)
	}

	wr.R.ParseForm()
	var err error
	filters := wr.R.Form
	exerciseId := int64(0)
	questId := int64(0)
	authorId := wr.NU.ID()

	if filters["exercise"] != nil {
		exerciseId, err = strconv.ParseInt(filters["exercise"][0], 10, 64)
		if err != nil {
			return infoTmpl, fmt.Errorf("answers: getlist: %v", err)
		}
	}

	if filters["quest"] != nil {
		questId, err = strconv.ParseInt(filters["quest"][0], 10, 64)
		if err != nil {
			return infoTmpl, fmt.Errorf("answers: getlist: %v", err)
		}
	}

	all, err := GetAnswers(wr, authorId, questId, exerciseId)
	if err != nil {
		return infoTmpl, fmt.Errorf("answers: getlist: %v", err)
	}

	tc["Content"] = all

	return infoTmpl, nil
}
