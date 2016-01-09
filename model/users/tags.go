package users

import (
	"fmt"

	"appengine/data"
	"appengine/srv"
)

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

func getUsersByTags(wr srv.WrapperRequest, tags []string) ([]*NUser, error) {
	nus := NewNUserBuffer()
	uTagsAll := NewUserTagBuffer()

	q := data.NewConn(wr, "users-tags")
	err := q.GetMany(&uTagsAll)
	if err != nil {
		return nus, fmt.Errorf("%v: %s", err, ERR_USERNOTFOUND)

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

func getUserTags(wr srv.WrapperRequest, nu *NUser) ([]string, error) {
	tags := make([]string, 0)
	userTags := NewUserTagBuffer()

	q := data.NewConn(wr, "users-tags")
	q.AddFilter("UserId =", nu.Id)
	err := q.GetMany(&userTags)
	if err != nil {
		return tags, fmt.Errorf("%v: %s", err, ERR_USERNOTFOUND)
	}

	tags = make([]string, 0)
	for _, utag := range userTags {
		tags = append(tags, utag.Tag)
	}

	return tags, nil

}

func addUserTags(wr srv.WrapperRequest, nu *NUser) error {
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

func deleteUserTags(wr srv.WrapperRequest, nu *NUser) error {
	userTags := NewUserTagBuffer()

	q := data.NewConn(wr, "users-tags")
	q.AddFilter("UserId =", nu.Id)
	q.GetMany(&userTags)

	for _, utag := range userTags {
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
