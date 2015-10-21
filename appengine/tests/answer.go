package tests



type Answer struct{
	Id int64           `json:",string" datastore:"-"`
	QuestionId uint64  `json:",string"`
	AuthorId uint64    `json:",string"`

	Solutions []int
	Comment string   //for non-test question
	// Timestamp??
}