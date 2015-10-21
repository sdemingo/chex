package tests




type Test struct{
	Id int64           `json:",string" datastore:"-"`
	AuthorId uint64    `json:",string"`

	Title string
	Description string
	Questions []Question
	Tags []string
}

