package tests

import (
	"appengine/data"
	"appengine/srv"
)

type TestTag struct {
	Id     int64 `json:",string" datastore:"-"`
	TestId int64
	Tag    string
}

func (t TestTag) ID() int64 {
	return t.Id
}

func (t *TestTag) SetID(id int64) {
	t.Id = id
}

type TestTagBuffer []*TestTag

func NewTestTagBuffer() TestTagBuffer {
	return make([]*TestTag, 0)
}

func (v TestTagBuffer) At(i int) data.DataItem {
	return data.DataItem(v[i])
}

func (v TestTagBuffer) Set(i int, t data.DataItem) {
	v[i] = t.(*TestTag)
}

func (v TestTagBuffer) Len() int {
	return len(v)
}

func addTestTags(wr srv.WrapperRequest, t *Test) error {
	q := data.NewConn(wr, "tests-tags")
	for _, tag := range t.Tags {
		tt := &TestTag{TestId: t.Id, Tag: tag}
		err := q.Put(tt)
		if err != nil {
			return err
		}
	}
	return nil
}
