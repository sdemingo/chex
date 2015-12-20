package users

import (
	"app/users"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"appengine/data"
	"appengine/srv"
)

type NUserBuffer []*users.NUser

func NewNUserBuffer() NUserBuffer {
	return make([]*users.NUser, 0)
}

func (v NUserBuffer) At(i int) data.DataItem {
	return data.DataItem(v[i])
}

func (v NUserBuffer) Set(i int, t data.DataItem) {
	v[i] = t.(*users.NUser)
}

func (v NUserBuffer) Len() int {
	return len(v)
}

func getUsers(wr srv.WrapperRequest, filters map[string][]string) (nus []*users.NUser, err error) {

	if filters["id"] != nil {
		nu, err := getUserById(wr, filters["id"][0])
		nus := make([]*users.NUser, 1)
		nus[0] = nu
		return nus, err
	}

	if filters["mail"] != nil {
		nu, err := getUserByMail(wr, filters["mail"][0])
		nus := make([]*users.NUser, 1)
		nus[0] = nu
		return nus, err
	}

	if filters["tags"] != nil {
		nus, err := getUsersByTags(wr, strings.Split(filters["tags"][0], ","))
		return nus, err
	}

	return
}

func putUser(wr srv.WrapperRequest, nu *users.NUser) error {

	/*if err := nu.IsValid(); err != nil {
		return err
	}*/

	nu.TimeStamp = time.Now()

	_, err := getUserByMail(wr, nu.Mail)
	if err == nil {
		return errors.New(users.ERR_DUPLICATEDUSER)
	}

	q := data.NewConn(wr, "users")
	q.Put(nu)

	// Add a UserTags entry for each tag for this user
	addUserTags(wr, nu)

	return nil
}

func updateUser(wr srv.WrapperRequest, nu *users.NUser) error {

	/*if err := nu.IsValid(); err != nil {
		return err
	}*/

	old, err := getUserById(wr, fmt.Sprintf("%d", nu.Id))
	if err != nil {
		return errors.New(users.ERR_USERNOTFOUND)
	}

	// invariant fields
	nu.Mail = old.Mail
	nu.Id = old.Id
	nu.TimeStamp = old.TimeStamp

	q := data.NewConn(wr, "users")
	q.Put(nu)

	// Delete all users-tags
	err = deleteUserTags(wr, nu)
	if err != nil {
		srv.Log(wr, err.Error())
	}
	// Add a UserTags entry for each tag for this user
	addUserTags(wr, nu)

	return nil
}

func deleteUser(wr srv.WrapperRequest, nu *users.NUser) error {
	/*if err := nu.IsValid(); err != nil {
		return err
	}*/

	// Delete all users-tags
	err := deleteUserTags(wr, nu)
	if err != nil {
		return err
	}

	q := data.NewConn(wr, "users")
	return q.Delete(nu)
}

func getUserByMail(wr srv.WrapperRequest, email string) (*users.NUser, error) {
	nus := NewNUserBuffer()
	nu := new(users.NUser)

	q := data.NewConn(wr, "users")
	q.AddFilter("Mail =", email)
	q.GetMany(&nus)
	if len(nus) == 0 {
		return nu, errors.New(users.ERR_USERNOTFOUND)
	}
	nu = nus[0]
	nu.Tags, _ = getUserTags(wr, nu)

	return nu, nil
}

func getUserById(wr srv.WrapperRequest, s_id string) (*users.NUser, error) {
	nu := new(users.NUser)

	id, err := strconv.ParseInt(s_id, 10, 64)
	if err != nil {
		return nu, errors.New(users.ERR_USERNOTFOUND)
	}

	nu.Id = id
	q := data.NewConn(wr, "users")
	q.Get(nu)
	nu.Tags, _ = getUserTags(wr, nu)

	return nu, nil
}
