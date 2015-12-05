package tests

import (
	"errors"
	"time"

	"app/users"
	"model/answers"
	"model/questions"

	"appengine/data"
	"appengine/srv"
)

const (
	ST_TESTOPEN   = iota + 1
	ST_TESTCLOSED = iota + 1

	ERR_NOTVALIDTEST = "Test no valido"
	ERR_TESTNOTFOUND = "Test no encontrado"
)

type Test struct {
	Id        int64        `json:",string" datastore:"-"`
	AuthorId  int64        `json:",string"`
	Author    *users.NUser `datastore:"-"`
	TimeStamp time.Time    `json:"`

	Title     string
	Desc      string
	Alias     string
	State     int
	Exercises []int64 // all exercises
	UList     []int64 // users allowed
	Tags      []string
}

func NewTest() *Test {
	t := new(Test)
	t.Exercises = make([]int64, 0)
	t.UList = make([]int64, 0)
	t.Tags = make([]string, 0)
	return t
}

func (t *Test) IsValid() error {
	if t != nil && t.Title == "" {
		return errors.New(ERR_NOTVALIDTEST)
	}
	return nil
}

func (t Test) ID() int64 {
	return t.Id
}

func (t *Test) SetID(id int64) {
	t.Id = id
}

type TestBuffer []*Test

func NewTestBuffer() TestBuffer {
	return make([]*Test, 0)
}

func (v TestBuffer) At(i int) data.DataItem {
	return data.DataItem(v[i])
}

func (v TestBuffer) Set(i int, t data.DataItem) {
	v[i] = t.(*Test)
}

func (v TestBuffer) Len() int {
	return len(v)
}

type TestTag struct {
	Id     int64 `json:",string" datastore:"-"`
	TestId int64
	Tag    string
}

func (t TestTag) ID() int64 {
	return t.Id
}

func (t *TestTag) SetID(id int64) {
	t.Id = id
}

type TestTagBuffer []*TestTag

func NewTestTagBuffer() TestTagBuffer {
	return make([]*TestTag, 0)
}

func (v TestTagBuffer) At(i int) data.DataItem {
	return data.DataItem(v[i])
}

func (v TestTagBuffer) Set(i int, t data.DataItem) {
	v[i] = t.(*TestTag)
}

func (v TestTagBuffer) Len() int {
	return len(v)
}

type Exercise struct {
	Id        int64              `json:",string" datastore:"-"`
	QuestId   int64              `json:",string"`
	Quest     questions.Question `datastore:"-"`
	BadPoint  int
	GoodPoint int
	AnswersId []int64           // list of answers
	Answers   []*answers.Answer `datastore:"-"`
}

func NewExercise() *Exercise {
	e := new(Exercise)
	e.AnswersId = make([]int64, 0)
	e.Answers = make([]*answers.Answer, 0)
	return e
}

func (e Exercise) ID() int64 {
	return e.Id
}

func (e *Exercise) SetID(id int64) {
	e.Id = id
}

type ExerciseBuffer []*Exercise

func NewExerciseBuffer() ExerciseBuffer {
	return make([]*Exercise, 0)
}

func (v ExerciseBuffer) At(i int) data.DataItem {
	return data.DataItem(v[i])
}

func (v ExerciseBuffer) Set(i int, t data.DataItem) {
	v[i] = t.(*Exercise)
}

func (v ExerciseBuffer) Len() int {
	return len(v)
}

func putTest(wr srv.WrapperRequest, t *Test) error {
	if err := t.IsValid(); err != nil {
		return err
	}

	t.TimeStamp = time.Now()
	t.AuthorId = wr.NU.Id

	q := data.NewConn(wr, "tests")
	q.Put(t)

	// Now, all questions  must be add. Checksums must be check
	// to avoid insert duplicated questions

	// Add a TestsTags entry for each tag for this questions
	addTestTags(wr, t)

	return nil
}

func addTestTags(wr srv.WrapperRequest, t *Test) error {
	q := data.NewConn(wr, "tests-tags")
	for _, tag := range t.Tags {
		tt := &TestTag{TestId: t.Id, Tag: tag}
		err := q.Put(tt)
		if err != nil {
			return err
		}
	}
	return nil
}
