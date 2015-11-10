package tests

import (
	"time"

	"app/users"
	"appengine/answers"
	"appengine/questions"
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
