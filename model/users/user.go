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

type UserTag struct {
	Id     int64 `json:",string" datastore:"-"`
	UserId int64
	Tag    string
}

func (ut UserTag) ID() int64 {
	return ut.Id
}

func (ut *UserTag) SetID(id int64) {
	ut.Id = id
}

type UserTagBuffer []*UserTag

func NewUserTagBuffer() UserTagBuffer {
	return make([]*UserTag, 0)
}

func (v UserTagBuffer) At(i int) data.DataItem {
	return data.DataItem(v[i])
}

func (v UserTagBuffer) Set(i int, t data.DataItem) {
	v[i] = t.(*UserTag)
}

func (v UserTagBuffer) Len() int {
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

func putUser(wr srv.WrapperRequest, nu users.NUser) error {

	if err := nu.IsValid(); err != nil {
		return err
	}

	nu.TimeStamp = time.Now()

	_, err := getUserByMail(wr, nu.Mail)
	if err == nil {
		return errors.New(users.ERR_DUPLICATEDUSER)
	}

	q := data.NewConn(wr, "users")
	q.Put(&nu)

	// Add a UserTags entry for each tag for this user
	addUserTags(wr, nu)

	return nil
}

func updateUser(wr srv.WrapperRequest, nu users.NUser) error {

	if err := nu.IsValid(); err != nil {
		return err
	}

	old, err := getUserById(wr, fmt.Sprintf("%d", nu.Id))
	if err != nil {
		return errors.New(users.ERR_USERNOTFOUND)
	}

	nu.Mail = old.Mail
	nu.Id = old.Id

	q := data.NewConn(wr, "users")
	q.Put(&nu)

	// Delete all users-tags
	err = deleteUserTags(wr, &nu)
	if err != nil {
		srv.AppWarning(wr, err.Error())
	}
	// Add a UserTags entry for each tag for this user
	addUserTags(wr, nu)

	return nil
}

func deleteUser(wr srv.WrapperRequest, nu *users.NUser) error {
	if err := nu.IsValid(); err != nil {
		return err
	}

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

func getUsersByTags(wr srv.WrapperRequest, tags []string) ([]*users.NUser, error) {
	nus := NewNUserBuffer()
	uTagsAll := NewUserTagBuffer()

	q := data.NewConn(wr, "users-tags")
	if q.GetMany(&uTagsAll) != nil {
		return nus, errors.New(users.ERR_USERNOTFOUND)
	}

	// After recover all UserTags it makes a homemade filtering
	// because dinamically filteres based on tags array are not
	// allowed in GAE datastore

	filtered := make(map[int64]int)
	for _, tag := range tags {
		for _, ut := range uTagsAll {
			if ut.Tag == tag {
				if _, ok := filtered[ut.UserId]; !ok {
					filtered[ut.UserId] = 1
				} else {
					filtered[ut.UserId]++
				}
			}
		}
	}

	for id, _ := range filtered {
		if filtered[id] == len(tags) {
			nu, err := getUserById(wr, fmt.Sprintf("%d", id))
			if err != nil {
				return nus, err
			}
			nus = append(nus, nu)
		}
	}

	return nus, nil
}

func getUserTags(wr srv.WrapperRequest, nu *users.NUser) ([]string, error) {
	tags := make([]string, 0)
	userTags := NewUserTagBuffer()

	q := data.NewConn(wr, "users-tags")
	q.AddFilter("UserId =", nu.Id)
	if q.GetMany(&userTags) != nil {
		return tags, errors.New(users.ERR_USERNOTFOUND)
	}

	tags = make([]string, 0)
	for _, utag := range userTags {
		tags = append(tags, utag.Tag)
	}

	return tags, nil

}

func addUserTags(wr srv.WrapperRequest, nu users.NUser) error {
	q := data.NewConn(wr, "users-tags")
	for _, tag := range nu.Tags {
		ut := &UserTag{UserId: nu.Id, Tag: tag}
		err := q.Put(ut)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteUserTags(wr srv.WrapperRequest, nu *users.NUser) error {
	userTags := NewUserTagBuffer()

	q := data.NewConn(wr, "users-tags")
	q.AddFilter("UserId =", nu.Id)
	q.GetMany(&userTags)

	for _, utag := range userTags {
		srv.AppWarning(wr, fmt.Sprintf("Borramos utag con id %d", utag.Id))
		err := q.Delete(utag)
		if err != nil {
			return err
		}
	}
	return nil
}

func getAllUserTags(wr srv.WrapperRequest) ([]string, error) {
	tagsMap := make(map[string]int, 0)
	userTags := NewUserTagBuffer()
	tags := make([]string, 0)

	q := data.NewConn(wr, "users-tags")
	err := q.GetMany(&userTags)
	if err != nil {
		return tags, err
	}

	tags = make([]string, len(userTags))
	for _, utag := range userTags {
		if _, ok := tagsMap[utag.Tag]; !ok {
			tagsMap[utag.Tag] = 1
			tags = append(tags, utag.Tag)
		}
	}

	return tags, nil
}
