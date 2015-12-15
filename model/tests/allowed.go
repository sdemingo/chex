package tests

import (
	"appengine/data"
	"appengine/srv"
)

type TestUser struct {
	Id     int64 `json:",string" datastore:"-"`
	TestId int64
	UserId int64
}

func (t TestUser) ID() int64 {
	return t.Id
}

func (t *TestUser) SetID(id int64) {
	t.Id = id
}

type TestUserBuffer []*TestUser

func NewTestUserBuffer() TestUserBuffer {
	return make([]*TestUser, 0)
}

func (v TestUserBuffer) At(i int) data.DataItem {
	return data.DataItem(v[i])
}

func (v TestUserBuffer) Set(i int, t data.DataItem) {
	v[i] = t.(*TestUser)
}

func (v TestUserBuffer) Len() int {
	return len(v)
}

func addUsersAllowed(wr srv.WrapperRequest, t *Test) error {
	q := data.NewConn(wr, "tests-users")
	for _, ex := range t.Exercises {
		err := q.Put(ex)
		if err != nil {
			return err
		}
	}
	return nil
}
