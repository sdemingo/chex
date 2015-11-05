package answers

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
	ERR_ANSWERNOTFOUND    = "Respuesta no encontrada"
	ERR_BADRENDEREDANSWER = "Respuesta renderizada erroneamente"
	ERR_ANSWERWITHOUTBODY = "Respuesta sin cuerpo definido"
)

var bodiesTable = []string{"",
	"testsingles-bodies",
	"testmultiples-bodies"}

type Answer struct {
	Id        int64 `json:",string" datastore:"-"`
	QuestId   int64 `json:",string"`
	AuthorId  int64 `json:",string"`
	TimeStamp time.Time
	Author    *users.NUser

	BodyType AnswerBodyType `json:",string"`
	Body     AnswerBody
	BodyId   int64
	Comment  string
}

type AnswerBodyType int

type AnswerBody interface {
	GetId() int64
	GetType() AnswerBodyType
	Equals(master AnswerBody) bool
	GetHTML(options []string) template.HTML
	IsUnsolved() bool
}

func NewAnswer(questionId int64, authorId int64) *Answer {
	a := new(Answer)
	a.Id = -1
	a.QuestId = questionId
	a.AuthorId = authorId
	a.Comment = ""
	a.BodyType = -1
	a.BodyId = -1

	return a
}

func (a *Answer) SetBody(abody AnswerBody) {
	a.Body = abody
	a.BodyId = abody.GetId()
	a.BodyType = abody.GetType()
}

func putAnswer(wr srv.WrapperRequest, a *Answer) error {

	if a.BodyId < 0 || a.BodyType < 0 {
		return errors.New(ERR_ANSWERWITHOUTBODY)
	}
	a.TimeStamp = time.Now()
	a.AuthorId = wr.NU.Id
	a.BodyType = a.Body.GetType()

	key := datastore.NewKey(wr.C, "answers", "", 0, nil)
	key, err := datastore.Put(wr.C, key, a)
	if err != nil {
		return err
	}
	a.Id = key.IntID()

	// metemos en answer body también
	bodyTable := bodiesTable[a.BodyType]
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
	bodyTable := bodiesTable[a.BodyType]
	k := datastore.NewKey(wr.C, bodyTable, "", id, nil)
	datastore.Get(wr.C, k, &b)
	a.Body = b
	return a, nil
}
