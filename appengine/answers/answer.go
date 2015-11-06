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
	Id      int64  `json:",string" datastore:"-"`
	RawBody string `datastore:"-"`

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
	a.BodyType = bodyType
	err := a.BuildBody()
	return a, err
}

// Try build a body of  BodyType from the RawBody property
func (a *Answer) BuildBody() error {
	var abody AnswerBody

	if a.BodyType < 0 {
		return errors.New(ERR_ANSWERWITHOUTBODY)
	}

	switch a.BodyType {
	case TYPE_TESTSINGLE:
		sol, err := strconv.ParseInt(a.RawBody, 10, 32)
		if err != nil {
			abody = NewTestSingleAnswer(-1)
		} else {
			abody = NewTestSingleAnswer(int(sol))
		}
	default:
		return errors.New(ERR_BADANSWERTYPE)
	}

	a.SetBody(abody)
	return nil
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

	var key *datastore.Key
	var bkey *datastore.Key

	a2, err := getAnswer(wr, a.AuthorId, a.QuestId)
	if err != nil { // answer not found. Is new answer
		key = datastore.NewKey(wr.C, "answers", "", 0, nil)
		a.Id = key.IntID()
		a.AuthorId = wr.NU.Id
		a.BodyType = a.Body.GetType()
		a.TimeStamp = time.Now()

		// actualizamos primero el body
		bodyTable := bodiesTable[a.BodyType]
		bkey = datastore.NewKey(wr.C, bodyTable, "", 0, nil)

		// debo meter un tipo concreto no una interfaz
		switch body := a.Body.(type) {
		case TestSingleBody:
			bkey, err = datastore.Put(wr.C, bkey, &body)

		}
		a.BodyId = bkey.IntID()
		if err != nil {
			return err
		}

		// actualizo ahora la respuesta
		_, err := datastore.Put(wr.C, key, a)
		if err != nil {
			return err
		}
		a.Id = key.IntID()

	} else { // answer found. Updated
		a.BodyId = a2.BodyId
		key = datastore.NewKey(wr.C, "answers", "", a2.Id, nil)
		a.TimeStamp = time.Now()

		// actualizamos primero el body
		bodyTable := bodiesTable[a.BodyType]
		bkey = datastore.NewKey(wr.C, bodyTable, "", a.BodyId, nil)

		// debo meter un tipo concreto no una interfaz
		switch body := a.Body.(type) {
		case TestSingleBody:
			bkey, err = datastore.Put(wr.C, bkey, &body)
		}
		//a.BodyId = bkey.IntID()
		if err != nil {
			return err
		}

		// actualizo ahora la respuesta
		_, err := datastore.Put(wr.C, key, a)
		if err != nil {
			return err
		}
		a.Id = key.IntID()
	}

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

func getAnswer(wr srv.WrapperRequest, authorId int64, questId int64) (*Answer, error) {

	var as []Answer
	var a Answer

	q := datastore.NewQuery("answers").Filter("AuthorId =", authorId).Filter("QuestId =", questId)

	keys, err := q.GetAll(wr.C, &as)
	if (len(keys) == 0) || err != nil {
		return nil, errors.New(ERR_ANSWERNOTFOUND)
	}
	a = as[0]
	a.Id = keys[0].IntID()

	return &a, nil
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
