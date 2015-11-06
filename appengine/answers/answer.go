package answers

import (
	"errors"
	//"fmt"
	"html/template"
	"strconv"
	"time"

	//"app/users"
	"appengine/datastore"
	"appengine/srv"
)

const (
	TYPE_TESTSINGLE   = 1
	TYPE_TESTMULTIPLE = iota

	ERR_ANSWERNOTFOUND    = "Respuesta no encontrada"
	ERR_BADRENDEREDANSWER = "Respuesta renderizada erroneamente"
	ERR_ANSWERWITHOUTBODY = "Respuesta sin cuerpo definido"
	ERR_BADANSWERTYPE     = "Respuesta con tipo de cuerpo desconocido"
)

var bodiesTable = []string{"",
	"testsingles-bodies",
	"testmultiples-bodies"}

type Answer struct {
	Id          int64  `json:",string" datastore:"-"`
	RawSolution string `datastore:"-"`

	QuestId   int64 `json:",string"`
	AuthorId  int64 `json:",string"`
	TimeStamp time.Time
	//Author    *users.NUser

	BodyType AnswerBodyType `json:",string"`
	Body     AnswerBody     `datastore:"-"`
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

// Return an answer wihtout a body
func NewAnswer(questionId int64, authorId int64) *Answer {
	a := new(Answer)
	a.Id = 0
	a.QuestId = questionId
	a.AuthorId = authorId
	a.Comment = ""
	a.BodyType = -1
	a.BodyId = -1

	return a
}

// Return an answer with a blank body of bodyType
func NewAnswerWithBody(questionId int64, authorId int64, bodyType AnswerBodyType) (*Answer, error) {
	a := NewAnswer(questionId, authorId)

	var abody AnswerBody

	switch bodyType {
	case TYPE_TESTSINGLE:
		abody = NewTestSingleAnswer(-1)
	default:
		return nil, errors.New(ERR_BADANSWERTYPE)
	}

	a.SetBody(abody)

	return a, nil
}

func (a *Answer) SetBody(abody AnswerBody) {
	a.Body = abody
	a.BodyId = abody.GetId()
	a.BodyType = abody.GetType()
}

// Create or update an answer
func putAnswer(wr srv.WrapperRequest, a *Answer) error {

	if a.BodyType < 0 {
		return errors.New(ERR_ANSWERWITHOUTBODY)
	}

	key := datastore.NewKey(wr.C, "answers", "", a.Id, nil)
	a.Id = key.IntID()

	a.TimeStamp = time.Now()
	a.AuthorId = wr.NU.Id
	a.BodyType = a.Body.GetType()

	key, err := datastore.Put(wr.C, key, a)
	if err != nil {
		return err
	}

	/*
		// metemos en answer body también
		bodyTable := bodiesTable[a.BodyType]
		key = datastore.NewKey(wr.C, bodyTable, "", 0, nil)
		srv.AppWarning(wr, fmt.Sprintf("%s  %d", bodyTable, a.Body.GetType()))
		key, err = datastore.Put(wr.C, key, &a.Body)
		if err != nil {
			return err
		}
		a.BodyId = key.IntID()
	*/
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
