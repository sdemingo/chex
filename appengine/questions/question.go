package questions

import (
	"app/users"
)

type Question struct {
	Id         int64        `json:",string" datastore:"-"`
	AuthorId   int64        `json:",string"`
	Author     *users.NUser `datastore:"-"`
	SolutionId int64        `json:",string"`
	Solution   *Answer      `datastore:"-"`

	AType   AnswerBodyType
	Text    string
	Hint    string
	Options []string
	Tags    []string
}

func NewQuestion(text string, options []string, tags []string) Question {
	q := new(Question)
	q.Id = -1
	q.AuthorId = -1
	q.SolutionId = -1

	q.Text = text
	q.Hint = ""
	q.Options = options
	q.Tags = tags

	return *q
}

func (q Question) SetSolution(sol *Answer) {
	if sol != nil {
		q.Solution = sol
		q.SolutionId = sol.Id
	}
}

func (q Question) SetAuthor(author *users.NUser) {
	if author != nil {
		q.Author = author
		q.AuthorId = author.Id
	}
}
