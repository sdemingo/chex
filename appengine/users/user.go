package users

import (
	"errors"
	"fmt"
	"strconv"

	"app/users"
	"appengine/datastore"
	"appengine/srv"
)

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
		nus, err := getUsersByTags(wr, filters["tags"])
		return nus, err
	}

	return
}

func getAllTags(wr srv.WrapperRequest) (tags []string, err error) {
	var nus []users.NUser
	found := make(map[string]int)

	q := datastore.NewQuery("users").Order("Tags")

	_, err = q.GetAll(wr.C, &nus)
	if err != nil {
		return tags, errors.New("User tags not found")
	}

	for i := 0; i < len(nus); i++ {
		if len(nus[i].Tags) > 0 {
			for j := 0; j < len(nus[i].Tags); j++ {
				if _, ok := found[nus[i].Tags[j]]; !ok {
					found[nus[i].Tags[j]] = 0
					tags = append(tags, nus[i].Tags[j])
				}
			}
		}
	}

	return tags, err
}

func putUser(wr srv.WrapperRequest, nu users.NUser) error {

	if err := nu.IsValid(); err != nil {
		return err
	}

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

	// Store the user-tags relations
	for _, tag := range nu.Tags {
		tkey := datastore.NewKey(wr.C, "users-tags", "", 0, nil)
		_, err = datastore.Put(wr.C, tkey, tag)
		if err != nil {
			return err
		}
	}
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

	return nil
}

func deleteUser(wr srv.WrapperRequest, nu users.NUser) error {
	if err := nu.IsValid(); err != nil {
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
	}

	return nus, nil
}

func getUsersByTags(wr srv.WrapperRequest, tags []string) ([]users.NUser, error) {
	var nus []users.NUser

	q := datastore.NewQuery("users").Filter("Tags =", tags)

	keys, err := q.GetAll(wr.C, &nus)
	if (len(keys) == 0) || err != nil {
		srv.AppWarning(wr, fmt.Sprintf("%s", err))
		return nus, errors.New("User not found. Bad tags")
	}

	for i := 0; i < len(nus); i++ {
		nus[i].Id = keys[i].IntID()
	}

	return nus, nil
}
