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

// Add the users allowed in the test to the database
func addUsersAllowed(wr srv.WrapperRequest, t *Test, err error) error {
	if err != nil {
		return err
	}
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

// Delete allowed users of the test
func deleteUsersAllowed(wr srv.WrapperRequest, t *Test, err error) error {
	if err != nil {
		return err
	}

	us := NewTestUserBuffer()
	q := data.NewConn(wr, "tests-users")
	q.AddFilter("TestId =", t.Id)
	q.GetMany(&us)

	for _, u := range us {
		err := q.Delete(u)
		if err != nil {
			return err
		}
	}

	return nil
}

// Fill the test UList array (allowed users ids lists)
func loadAllowed(wr srv.WrapperRequest, t *Test, err error) error {
	if err != nil {
		return err
	}
	tus := NewTestUserBuffer()
	t.UList = make([]int64, 0)

	q := data.NewConn(wr, "tests-users")
	q.AddFilter("TestId=", t.Id)
	err = q.GetMany(&tus)
	if err != nil {
		return err
	}

	for i := range tus {
		t.UList = append(t.UList, tus[i].UserId)
	}
	return err
}

// Return the users allowed for this tests
func getUsersAllowed(wr srv.WrapperRequest, t *Test) ([]*users.NUser, error) {
	nus := make([]*users.NUser, 0)
	qu := data.NewConn(wr, "users")
	for i := range t.UList {
		us := new(users.NUser)
		us.SetID(t.UList[i])
		err := qu.Get(us)
		if err != nil {
			continue
		}
		nus = append(nus, us)
	}

	return nus, nil
}

// Return all tests allowed for a user
func getTestAllowedForUser(wr srv.WrapperRequest, userId int64) ([]*Test, error) {
	tus := NewTestUserBuffer()
	tst := make([]*Test, 0)

	q := data.NewConn(wr, "tests-users")
	q.AddFilter("UserId=", userId)
	err := q.GetMany(&tus)
	if err != nil {
		return tst, err
	}

	for i := range tus {
		t, err := getTestById(wr, tus[i].TestId)
		if err != nil {
			continue
		}
		tst = append(tst, t)
	}

	return tst, nil
}
