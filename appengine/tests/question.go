package tests




type Question struct{
	Id int64           `json:",string" datastore:"-"`
	AuthorId uint64    `json:",string"`
	SolutionId uint64  `json:",string"`

	Text string
	Hint string
	Options []string
	Tags []string
}
