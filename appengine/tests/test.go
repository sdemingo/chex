package tests

import (
	"appengine/questions"
)

type Test struct {
	Id       int64  `json:",string" datastore:"-"`
	AuthorId uint64 `json:",string"`

	Title       string
	Description string
	Questions   []questions.Question
	Tags        []string
}
