package tests

import (
	"model/users"

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
	for _, uid := range t.UList {
		tu := &TestUser{Id: 0, TestId: t.Id, UserId: uid}
		err := q.Put(tu)
		if err != nil {
			return err
		}
	}
	return nil
}

func getUsersAllowed(wr srv.WrapperRequest, t *Test) ([]*users.NUser, error) {
	tus := NewTestUserBuffer()
	nus := make([]*users.NUser, 0)

	q := data.NewConn(wr, "tests-users")
	q.AddFilter("TestId=", t.Id)
	err := q.GetMany(tus)

	qu := data.NewConn(wr, "users")
	for i := range tus {
		tu := tus.At(i)
		us := new(users.NUser)
		us.SetID(tu.ID())
		qu.Get(us)
		nus = append(nus, us)
	}

	return nus, err
}
