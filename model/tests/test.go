package tests

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"model/users"

	"appengine/data"
	"appengine/srv"
)

const (
	ST_TESTOPEN   = iota + 1
	ST_TESTCLOSED = iota + 1

	ERR_TESTNOTFOUND          = "Test no encontrado"
	ERR_TESTEXERCISESNOTVALID = "Un Test no puede contener ejercicicios sin solucionar. Eliminelos para poder guardar"
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

func (t *Test) GetStringTags() string {
	return strings.Join(t.Tags, ", ")
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

func getTests(wr srv.WrapperRequest, filters map[string][]string) (TestBuffer, error) {
	ts := NewTestBuffer()
	var err error

	if filters["id"] != nil {
		id, err := strconv.ParseInt(filters["id"][0], 10, 64)
		if err != nil {
			return ts, fmt.Errorf("%v: %s", err, ERR_TESTNOTFOUND)
		}
		t, err := getTestById(wr, id)
		ts = append(ts, t)
		return ts, err
	}

	if filters["tags"] != nil {
		ts, err = getTestsByTags(wr, strings.Split(filters["tags"][0], ","))
		return ts, err
	}

	if filters["foruser"] != nil {
		id, err := strconv.ParseInt(filters["foruser"][0], 10, 64)
		if err != nil {
			return ts, fmt.Errorf("%v: %s", err, ERR_TESTNOTFOUND)
		}
		ts, err = GetTestAllowedForUser(wr, id)
		return ts, err
	}

	if filters["author"] != nil {
		id, err := strconv.ParseInt(filters["author"][0], 10, 64)
		if err != nil {
			return ts, fmt.Errorf("%v: %s", err, ERR_TESTNOTFOUND)
		}
		ts, err = getTestsByAuthor(wr, id)
		return ts, err
	}

	return ts, err
}

// Return all tests from authorId
func getTestsByAuthor(wr srv.WrapperRequest, authorId int64) (TestBuffer, error) {
	ts := NewTestBuffer()

	qry := data.NewConn(wr, "tests")
	qry.AddFilter("AuthorId =", authorId)
	err := qry.GetMany(&ts)

	for i := range ts {
		loadTestTags(wr, ts[i], err)
	}

	if err != nil {
		return ts, fmt.Errorf("gettestbyauthor: %v", err)
	}

	return ts, nil
}

// Return the test with id
func getTestById(wr srv.WrapperRequest, id int64) (*Test, error) {
	t := NewTest()

	t.Id = id
	qry := data.NewConn(wr, "tests")
	err := qry.Get(t)

	err = loadTestTags(wr, t, err)
	err = loadExercises(wr, t, err)
	err = loadAllowed(wr, t, err)

	if err != nil {
		return nil, fmt.Errorf("gettestbyid: %v", err)
	}

	return t, nil
}

// Write new test in the database
func putTest(wr srv.WrapperRequest, t *Test) error {

	t.TimeStamp = time.Now()
	t.AuthorId = wr.NU.ID()

	q := data.NewConn(wr, "tests")
	err := q.Put(t)

	err = checkExercises(wr, t, err)
	err = addExercises(wr, t, err)
	err = addUsersAllowed(wr, t, err)
	err = addTestTags(wr, t, err)

	if err != nil {
		return fmt.Errorf("puttest: %v", err)
	}

	return nil
}

func updateTest(wr srv.WrapperRequest, t *Test) error {

	old, err := getTestById(wr, t.Id)
	if err != nil {
		return fmt.Errorf("updatetest: %v", err)
	}

	if wr.NU.ID() != old.AuthorId {
		return fmt.Errorf("updatetest: %s", users.ERR_NOTOPERATIONALLOWED)
	}

	// invariant fields
	t.TimeStamp = old.TimeStamp
	t.AuthorId = old.AuthorId

	q := data.NewConn(wr, "tests")
	err = q.Put(t)
	err = checkExercises(wr, t, err)
	err = deleteExercises(wr, t, err)
	err = addExercises(wr, t, err)
	err = deleteUsersAllowed(wr, t, err)
	err = addUsersAllowed(wr, t, err)
	err = deleteTestTags(wr, t, err)
	err = addTestTags(wr, t, err)

	if err != nil {
		return fmt.Errorf("updatetest: %v", err)
	}

	return nil
}

func deleteTest(wr srv.WrapperRequest, t *Test) error {

	q := data.NewConn(wr, "tests")
	err := q.Delete(t)
	err = deleteExercises(wr, t, err)
	err = deleteUsersAllowed(wr, t, err)
	err = deleteTestTags(wr, t, err)

	if err != nil {
		return fmt.Errorf("deletetest: %v", err)
	}
	return nil
}
