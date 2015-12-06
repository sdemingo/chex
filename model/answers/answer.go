package answers

import (
	"errors"
	"html/template"
	"strconv"
	"time"

	"appengine/data"
	"appengine/srv"
)

const (
	TYPE_TESTSINGLE   = iota + 1
	TYPE_TESTMULTIPLE = iota + 1

	ERR_ANSWERNOTFOUND    = "Respuesta no encontrada"
	ERR_BADRENDEREDANSWER = "Respuesta renderizada erroneamente"
	ERR_ANSWERWITHOUTBODY = "Respuesta sin cuerpo definido"
	ERR_BADANSWERTYPE     = "Respuesta con tipo de cuerpo desconocido"
)

var bodiesTable = []string{
	TYPE_TESTSINGLE:   "testsingles-bodies",
	TYPE_TESTMULTIPLE: "testmultiples-bodies"}

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
	GetType() AnswerBodyType
	Equals(master AnswerBody) bool
	GetHTML(options []string) (template.HTML, template.HTML, error)
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
	a.BodyType = abody.GetType()
}

func (a Answer) ID() int64 {
	return a.Id
}

func (a *Answer) SetID(id int64) {
	a.Id = id
}

type AnswerBuffer []*Answer

func NewAnswerBuffer() AnswerBuffer {
	return make([]*Answer, 0)
}

func (v AnswerBuffer) At(i int) data.DataItem {
	return data.DataItem(v[i])
}

func (v AnswerBuffer) Set(i int, t data.DataItem) {
	v[i] = t.(*Answer)
}

func (v AnswerBuffer) Len() int {
	return len(v)
}

func GetAnswer(wr srv.WrapperRequest, authorId int64, questId int64) (*Answer, error) {
	as := NewAnswerBuffer()
	a := new(Answer)

	q := data.NewConn(wr, "answers")
	q.AddFilter("QuestId =", questId)
	err := q.GetMany(&as)
	if err != nil || len(as) == 0 {
		return nil, errors.New(ERR_ANSWERNOTFOUND)
	}

	err = getAnswerBody(wr, a)
	if err != nil {
		return nil, errors.New(ERR_ANSWERNOTFOUND)
	}

	return a, err
}

// Create or update an answer
func putAnswer(wr srv.WrapperRequest, a *Answer) error {
	if a.BodyType < 0 {
		return errors.New(ERR_ANSWERWITHOUTBODY)
	}

	a2, err := GetAnswer(wr, a.AuthorId, a.QuestId)
	qry := data.NewConn(wr, "answers")

	if err != nil { // New
		a.AuthorId = wr.NU.Id
		a.TimeStamp = time.Now()

		err = putAnswerBody(wr, a)
		if err != nil {
			return err
		}

	} else { // Updated
		a.BodyId = a2.BodyId
		a.TimeStamp = time.Now()
		a.Id = a2.Id

		err = putAnswerBody(wr, a)
		if err != nil {
			return err
		}
	}

	qry.Put(a)

	return nil
}

func putAnswerBody(wr srv.WrapperRequest, a *Answer) error {

	bodyTable := bodiesTable[a.BodyType]
	q := data.NewConn(wr, bodyTable)
	var err error

	switch a.Body.GetType() {
	case TYPE_TESTSINGLE:
		tbody := a.Body.(*TestSingleBody)
		q.Put(tbody)
		a.BodyId = tbody.ID()
	default:
		err = errors.New(ERR_ANSWERWITHOUTBODY)
	}
	return err
}

func getAnswerBody(wr srv.WrapperRequest, a *Answer) error {

	bodyTable := bodiesTable[a.BodyType]

	q := data.NewConn(wr, bodyTable)
	var err error
	switch a.BodyType {
	case TYPE_TESTSINGLE:
		body := NewTestSingleAnswer(-1)
		body.Id = a.BodyId
		err = q.Get(body)
		a.Body = body
	}
	return err
}

func getAnswers(wr srv.WrapperRequest, filters map[string][]string) (AnswerBuffer, error) {
	as := NewAnswerBuffer()
	var err error

	if filters["id"] != nil {
		a, err := getAnswersById(wr, filters["id"][0])
		as[0] = a
		return as, err
	}

	return as, err
}

func getAnswersById(wr srv.WrapperRequest, s_id string) (*Answer, error) {
	a := NewAnswer(-1, -1)

	id, err := strconv.ParseInt(s_id, 10, 64)
	if err != nil {
		return a, errors.New(ERR_ANSWERNOTFOUND)
	}

	qry := data.NewConn(wr, "answers")
	a.Id = id
	if id != 0 {
		qry.Get(a)

	} else {
		return a, errors.New(ERR_ANSWERNOTFOUND)
	}

	// falta el answer body
	getAnswerBody(wr, a)

	return a, err
}
