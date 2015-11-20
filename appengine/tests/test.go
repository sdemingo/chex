package tests

import (
	"errors"
	"time"

	"app/users"
	"appengine/answers"
	"appengine/datastore"
	"appengine/questions"
	"appengine/srv"
)

const (
	ST_TESTOPEN   = 1
	ST_TESTCLOSED = iota

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

func (t *Test) IsValid() error {
	if t != nil && t.Title == "" {
		return errors.New(ERR_NOTVALIDTEST)
	}
	return nil
}

type TestTag struct {
	Id     int64 `json:",string" datastore:"-"`
	TestId int64
	Tag    string
}

type Exercise struct {
	Id        int64              `json:",string" datastore:"-"`
	QuestId   int64              `json:",string"`
	Quest     questions.Question `datastore:"-"`
	BadPoint  int
	GoodPoint int
	AnswersId []int64          // list of answers
	Answers   []answers.Answer `datastore:"-"`
}

func putTest(wr srv.WrapperRequest, t *Test) error {
	if err := t.IsValid(); err != nil {
		return err
	}

	t.TimeStamp = time.Now()
	t.AuthorId = wr.NU.Id

	key := datastore.NewKey(wr.C, "tests", "", 0, nil)
	key, err := datastore.Put(wr.C, key, t)
	if err != nil {
		return err
	}
	t.Id = key.IntID()

	// Now, all questions  must be add. Checksums must be check
	// to avoid insert duplicated questions

	// Add a TestsTags entry for each tag for this questions
	addTestTags(wr, t)

	return nil
}

func addTestTags(wr srv.WrapperRequest, t *Test) error {
	for _, tag := range t.Tags {
		key := datastore.NewKey(wr.C, "tests-tags", "", 0, nil)
		tt := TestTag{TestId: t.Id, Tag: tag}
		key, err := datastore.Put(wr.C, key, &tt)
		if err != nil {
			return err
		}
	}
	return nil
}
