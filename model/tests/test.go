package tests

import (
	"errors"
	"strconv"
	"time"

	"app/users"

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
	Course    string
	Desc      string
	State     int
	Exercises []*Exercise `datastore:"-"` // all exercises
	UList     []int64     `datastore:"-"` // users allowed
	Tags      []string    `datastore:"-"`
}

func NewTest() *Test {
	t := new(Test)
	t.Exercises = make([]*Exercise, 0)
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

func putTest(wr srv.WrapperRequest, t *Test) error {
	if err := t.IsValid(); err != nil {
		return err
	}

	t.TimeStamp = time.Now()
	t.AuthorId = wr.NU.Id

	q := data.NewConn(wr, "tests")
	q.Put(t)

	err := addExercises(wr, t)
	err = addUsersAllowed(wr, t)
	err = addTestTags(wr, t)

	return err
}

func getTests(wr srv.WrapperRequest, filters map[string][]string) (TestBuffer, error) {
	ts := NewTestBuffer()
	var err error

	if filters["id"] != nil {
		// q, err := getQuestById(wr, filters["id"][0])
		// qs = append(qs, q)
		return ts, err
	}

	if filters["tags"] != nil {
		//qs, err := getQuestByTags(wr, strings.Split(filters["tags"][0], ","))
		return ts, err
	}

	if filters["author"] != nil {
		ts, err = getTestsByAuthor(wr, filters["author"][0])
		return ts, err
	}

	return ts, err
}

func getTestsByAuthor(wr srv.WrapperRequest, authorId string) (TestBuffer, error) {
	ts := NewTestBuffer()

	id, err := strconv.ParseInt(authorId, 10, 64)
	if err != nil {
		return ts, errors.New(ERR_TESTNOTFOUND)
	}

	qry := data.NewConn(wr, "tests")
	qry.AddFilter("AuthorId =", id)
	qry.GetMany(&ts)

	// for i := range qs {
	// 	qs[i].Tags, _ = getQuestTags(wr, qs[i])
	// }

	return ts, nil
}
