package questions

import (
	"encoding/json"
	"errors"

	"model/answers"
	"model/users"

	"appengine/srv"
)

// Templates

var newTmpl = "model/questions/tmpl/edit.html"
var viewTmpl = "model/questions/tmpl/view.html"
var infoTmpl = "model/questions/tmpl/info.html"
var mainTmpl = "model/questions/tmpl/main.html"

func Main(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	return mainTmpl, nil
}

func GetList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	var err error
	// err := srv.CheckPerm(wr, users.OP_COMMITTER)
	// if err != nil {
	// 	return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	// }

	wr.R.ParseForm()
	qs, err := getQuestions(wr, wr.R.Form)
	if err != nil {
		return infoTmpl, err
	}

	tc["Content"] = qs

	return infoTmpl, nil
}

func GetOne(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	var err error
	// only teacher must entry to the question throught this
	// handler. A student or other should use test handlers
	// err := srv.CheckPerm(wr, users.OP_COMMITTER)
	// if err != nil {
	// 	return viewTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	// }

	wr.R.ParseForm()
	qs, err := getQuestions(wr, wr.R.Form)
	if len(qs) == 0 || err != nil {
		return viewTmpl, errors.New(ERR_QUESTNOTFOUND)
	}
	q := qs[0]

	// A question only can be viewed by and admin or by their writer
	// TODO

	tc["Content"] = q

	return viewTmpl, nil
}

func GetTagsList(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	var err error
	// err := srv.CheckPerm(wr, users.OP_VIEWER)
	// if err != nil {
	// 	return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	// }

	var tags []string
	if wr.NU.GetRole() == users.ROLE_ADMIN {
		tags, err = getAllQuestionsTags(wr)
	} else {
		tags, err = getQuestionsTagsFromUser(wr, wr.NU.ID())
	}

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
	var err error
	// err := srv.CheckPerm(wr, users.OP_COMMITTER)
	// if err != nil {
	// 	return newTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	// }

	wr.R.ParseForm()
	qs, err := getQuestions(wr, wr.R.Form)
	if len(qs) == 0 || err != nil {
		return viewTmpl, errors.New(ERR_QUESTNOTFOUND)
	}
	q := qs[0]

	// Chequear tambien la autoria de la pregunta
	// solo puede actualizarla el autor

	// Chequear que no haya sido añadida a examenes o ya
	//contestada. En ese caso no se puede actualizar

	tc["Content"] = q

	return newTmpl, nil
}

func Add(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	var err error
	// err := srv.CheckPerm(wr, users.OP_COMMITTER)
	// if err != nil {
	// 	return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	// }

	var q Question

	decoder := json.NewDecoder(wr.R.Body)
	err = decoder.Decode(&q)
	if err != nil {
		return infoTmpl, err
	}

	err = putQuestion(wr, &q)
	if err != nil {
		return infoTmpl, err
	}

	tc["Content"] = q

	return infoTmpl, nil
}

func Solve(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	var err error
	// err := srv.CheckPerm(wr, users.OP_COMMITTER)
	// if err != nil {
	// 	return infoTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
	// }

	var a *answers.Answer

	decoder := json.NewDecoder(wr.R.Body)
	err = decoder.Decode(&a)
	if err != nil {
		return infoTmpl, err
	}

	a.BuildBody()
	a.AuthorId = wr.NU.ID()

	quest, err := GetQuestById(wr, a.QuestId)
	if err != nil {
		return infoTmpl, err
	}

	// Only que teacher author of the quest can add the solution
	if quest.AuthorId != a.AuthorId {
		return infoTmpl, errors.New(answers.ERR_AUTHORID)
	}

	err = answers.PutSolutionAnswer(wr, a)
	if err != nil {
		return infoTmpl, err
	}

	tc["Content"] = a
	return infoTmpl, nil
}
