package marks

import (
//"fmt"
)

type Mark struct {
	Id       int64 `json:",string" datastore:"-"`
	AuthorId int64 `json:",string"`
	TestId   int64 `json:",string"`
	Value    int
}
