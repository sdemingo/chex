package questions

import (
	"errors"
	"html/template"
	"strconv"
	"time"

	"app/users"
	"appengine/datastore"
	"appengine/srv"
)

const (
	TYPE_TESTSINGLE   = 1
	TYPE_TESTMULTIPLE = iota

	//ERR_BADANSWERMARK  = "Mark not valid"
	ERR_ANSWERNOTFOUND = "Respuesta no encontrada"
)

var bodiesTable = []string{"",
	"testsingles-bodies",
	"testmultiples-bodies"}

type Answer struct {
	Id        int64 `json:",string" datastore:"-"`
	QuestId   int64 `json:",string"`
	Quest     *Question
	AuthorId  int64          `json:",string"`
	TimeStamp time.Time      `json:"`
	AType     AnswerBodyType `json:"`
	Author    *users.NUser

	Body    AnswerBody
	BodyId  int64 `json:",string"`
	Comment string
}

type AnswerBodyType int

type AnswerBody interface {
	GetId() int64
	GetType() AnswerBodyType
	Equals(master AnswerBody) bool
	GetHTML(options []string) template.HTML
	IsUnsolved() bool
}

func NewAnswer(questionId int64, authorId int64, body AnswerBody) *Answer {
	a := new(Answer)
	a.Id = -1
	a.QuestId = questionId
	a.AuthorId = authorId
	a.Body = body
	a.BodyId = body.GetId()
	a.Comment = ""

	return a
}

func putAnswer(wr srv.WrapperRequest, a *Answer) error {

	a.TimeStamp = time.Now()
	a.AuthorId = wr.NU.Id
	a.AType = a.Body.GetType()

	key := datastore.NewKey(wr.C, "answers", "", 0, nil)
	key, err := datastore.Put(wr.C, key, a)
	if err != nil {
		return err
	}
	a.Id = key.IntID()

	// metemos en answer body también
	bodyTable := bodiesTable[a.AType]
	key = datastore.NewKey(wr.C, bodyTable, "", 0, nil)
	key, err = datastore.Put(wr.C, key, a.Body)
	if err != nil {
		return err
	}
	a.BodyId = key.IntID()

	return nil
}

func getAnswers(wr srv.WrapperRequest, filters map[string][]string) ([]Answer, error) {

	var as []Answer
	var err error

	if filters["id"] != nil {
		a, err := getAnswersById(wr, filters["id"][0])
		as := make([]Answer, 1)
		as[0] = a
		return as, err
	}

	return as, err
}

func getAnswersById(wr srv.WrapperRequest, s_id string) (Answer, error) {

	var a Answer

	id, err := strconv.ParseInt(s_id, 10, 64)
	if err != nil {
		return a, errors.New(ERR_ANSWERNOTFOUND)
	}

	if id != 0 {
		k := datastore.NewKey(wr.C, "answers", "", id, nil)
		datastore.Get(wr.C, k, &a)
	} else {
		return a, errors.New(ERR_ANSWERNOTFOUND)
	}

	a.Id = id

	// falta el answer body
	var b AnswerBody
	bodyTable := bodiesTable[a.AType]
	k := datastore.NewKey(wr.C, bodyTable, "", id, nil)
	datastore.Get(wr.C, k, &b)
	a.Body = b
	return a, nil
}
