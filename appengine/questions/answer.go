package questions

import (
	"app/users"
)

const (
	TYPE_TESTSINGLE   = iota
	TYPE_TESTMULTIPLE = iota

	ERR_BADANSWERMARK = "Mark not valid"
)

type Answer struct {
	Id       int64 `json:",string" datastore:"-"`
	QuestId  int64 `json:",string"`
	Quest    *Question
	AuthorId int64 `json:",string"`
	Author   *users.NUser

	Body    AnswerBody
	Comment string
	// Timestamp??
}

type AnswerBodyType int

type AnswerBody interface {
	GetId() int64
	GetType() int
	Equals(master AnswerBody) bool
}

func NewAnswer(body AnswerBody) Answer {
	a := new(Answer)
	a.Id = -1
	a.QuestId = -1
	a.AuthorId = -1
	a.Body = body
	a.Comment = ""

	return *a
}
