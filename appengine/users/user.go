package users

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"app/users"
	"appengine/datastore"
	"appengine/srv"
)

type UserTag struct {
	Id     int64 `json:",string" datastore:"-"`
	UserId int64
	Tag    string
}

func getUsers(wr srv.WrapperRequest, filters map[string][]string) (nus []users.NUser, err error) {

	if filters["id"] != nil {
		nu, err := getUserById(wr, filters["id"][0])
		nus := make([]users.NUser, 1)
		nus[0] = nu
		return nus, err
	}

	if filters["mail"] != nil {
		nu, err := getUserByMail(wr, filters["mail"][0])
		nus := make([]users.NUser, 1)
		nus[0] = nu
		return nus, err
	}

	if filters["role"] != nil {
		nus, err := getUsersByRole(wr, filters["role"][0])
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
		return errors.New("Usuario duplicado")
	}

	key := datastore.NewKey(wr.C, "users", "", 0, nil)
	key, err = datastore.Put(wr.C, key, &nu)
	if err != nil {
		return err
	}
	nu.Id = key.IntID()

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
		return errors.New("Usuario no encontrado")
	}

	nu.Mail = old.Mail
	key := datastore.NewKey(wr.C, "users", "", old.Id, nil)
	key, err = datastore.Put(wr.C, key, &nu)
	if err != nil {
		return err
	}

	// Delete all users-tags
	deleteUserTags(wr, nu)
	// Add a UserTags entry for each tag for this user
	addUserTags(wr, nu)

	return nil
}

func deleteUser(wr srv.WrapperRequest, nu users.NUser) error {
	if err := nu.IsValid(); err != nil {
		return err
	}

	// Delete all users-tags
	err := deleteUserTags(wr, nu)
	if err != nil {
		return err
	}

	key := datastore.NewKey(wr.C, "users", "", nu.Id, nil)
	return datastore.Delete(wr.C, key)
}

func getUserByMail(wr srv.WrapperRequest, email string) (users.NUser, error) {
	var nus []users.NUser
	var nu users.NUser

	q := datastore.NewQuery("users").Filter("Mail =", email)

	keys, err := q.GetAll(wr.C, &nus)
	if (len(keys) == 0) || err != nil {
		return nu, errors.New("User not found. Bad mail")
	}
	nu = nus[0]
	nu.Id = keys[0].IntID()
	nu.Tags, _ = getUserTags(wr, nu)

	return nu, nil
}

func getUserById(wr srv.WrapperRequest, s_id string) (users.NUser, error) {

	var nu users.NUser

	id, err := strconv.ParseInt(s_id, 10, 64)
	if err != nil {
		return nu, errors.New("User not found. Bad ID")
	}

	if id != 0 {
		k := datastore.NewKey(wr.C, "users", "", id, nil)
		datastore.Get(wr.C, k, &nu)
	} else {
		return nu, errors.New("User not found. Bad ID")
	}

	nu.Id = id
	nu.Tags, _ = getUserTags(wr, nu)

	return nu, nil
}

func getUsersByRole(wr srv.WrapperRequest, s_role string) ([]users.NUser, error) {
	var nus []users.NUser

	role, err := strconv.ParseInt(s_role, 10, 64)
	if err != nil {
		return nus, errors.New("User role bad formatted")
	}

	q := datastore.NewQuery("users").Filter("Role =", role)

	keys, err := q.GetAll(wr.C, &nus)
	if (len(keys) == 0) || err != nil {
		return nus, errors.New("User not found. Bad role")
	}

	for i := 0; i < len(nus); i++ {
		nus[i].Id = keys[i].IntID()
		nus[i].Tags, _ = getUserTags(wr, nus[i])
	}

	return nus, nil
}

func getUsersByTags(wr srv.WrapperRequest, tags []string) ([]users.NUser, error) {
	var nus []users.NUser
	var uTagsAll []UserTag

	q := datastore.NewQuery("users-tags")
	_, err := q.GetAll(wr.C, &uTagsAll)
	if err != nil {
		return nus, errors.New("User not found. Bad tags")
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

func getUserTags(wr srv.WrapperRequest, nu users.NUser) ([]string, error) {

	var tags []string
	var userTags []UserTag

	q := datastore.NewQuery("users-tags").Filter("UserId =", nu.Id)
	_, err := q.GetAll(wr.C, &userTags)
	if err != nil {
		return tags, err
	}
	tags = make([]string, 0)
	for _, utag := range userTags {
		tags = append(tags, utag.Tag)
	}

	return tags, nil

}

func addUserTags(wr srv.WrapperRequest, nu users.NUser) error {
	for _, tag := range nu.Tags {
		key := datastore.NewKey(wr.C, "users-tags", "", 0, nil)
		ut := UserTag{UserId: nu.Id, Tag: tag}
		key, err := datastore.Put(wr.C, key, &ut)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteUserTags(wr srv.WrapperRequest, nu users.NUser) error {
	var userTags []UserTag

	q := datastore.NewQuery("users-tags").Filter("UserId =", nu.Id)
	keys, err := q.GetAll(wr.C, &userTags)
	if err != nil {
		return err
	}
	for i := 0; i < len(userTags); i++ {
		userTags[i].Id = keys[i].IntID()
	}

	for _, utag := range userTags {
		key := datastore.NewKey(wr.C, "users-tags", "", utag.Id, nil)
		err = datastore.Delete(wr.C, key)
		if err != nil {
			return err
		}
	}
	return nil
}

func getAllUserTags(wr srv.WrapperRequest) ([]string, error) {
	var tagsMap = make(map[string]int, 0)
	var userTags []UserTag
	var tags = make([]string, 0)

	q := datastore.NewQuery("users-tags")
	_, err := q.GetAll(wr.C, &userTags)
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
