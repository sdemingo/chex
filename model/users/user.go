package users

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"appengine/data"
	"appengine/srv"
)

const (
	ROLE_GUEST   = iota
	ROLE_STUDENT = iota
	ROLE_TEACHER = iota
	ROLE_ADMIN   = iota

	MAXSZUSERNAME = 100

	ERR_NOTVALIDUSER        = "Usuario no valido"
	ERR_DUPLICATEDUSER      = "Usuario duplicado"
	ERR_USERNOTFOUND        = "Usuario no encontrado"
	ERR_NOTOPERATIONALLOWED = "Operación no permitida"
)

var roleNames = []string{
	ROLE_GUEST:   "Invitado",
	ROLE_STUDENT: "Estudiante",
	ROLE_TEACHER: "Profesor",
	ROLE_ADMIN:   "Administrador"}

type NUser struct {
	Id        int64 `json:",string" datastore:"-"`
	Mail      string
	Name      string
	Role      int8      `json:",string"`
	Tags      []string  `datastore:"-"`
	TimeStamp time.Time `json:"`
}

func (n *NUser) GetRole() int8 {
	return n.Role
}

func (n *NUser) GetEmail() string {
	return n.Mail
}

func (n *NUser) GetInfo() map[string]string {
	info := make(map[string]string)

	info["Username"] = n.Name
	if int(n.Role) < len(roleNames) {
		info["RoleName"] = roleNames[n.Role]
	} else {
		info["RoleName"] = ""
	}
	info["Tags"] = strings.Join(n.Tags, ",")
	info["TimeStamp"] = n.TimeStamp.Format("02/01/2006")

	return info
}

func (n *NUser) ID() int64 {
	return n.Id
}

func (n *NUser) SetID(id int64) {
	n.Id = id
}

type NUserBuffer []*NUser

func NewNUserBuffer() NUserBuffer {
	return make([]*NUser, 0)
}

func (v NUserBuffer) At(i int) data.DataItem {
	return data.DataItem(v[i])
}

func (v NUserBuffer) Set(i int, t data.DataItem) {
	v[i] = t.(*NUser)
}

func (v NUserBuffer) Len() int {
	return len(v)
}

func getUsers(wr srv.WrapperRequest, filters map[string][]string) (nus []*NUser, err error) {

	if filters["id"] != nil {
		nu, err := getUserById(wr, filters["id"][0])
		nus := make([]*NUser, 1)
		nus[0] = nu
		return nus, err
	}

	if filters["mail"] != nil {
		nu, err := getUserByMail(wr, filters["mail"][0])
		nus := make([]*NUser, 1)
		nus[0] = nu
		return nus, err
	}

	if filters["tags"] != nil {
		nus, err := getUsersByTags(wr, strings.Split(filters["tags"][0], ","))
		return nus, err
	}

	return
}

func putUser(wr srv.WrapperRequest, nu *NUser) error {

	/*if err := nu.IsValid(); err != nil {
		return err
	}*/

	nu.TimeStamp = time.Now()

	_, err := getUserByMail(wr, nu.Mail)
	if err == nil {
		return errors.New(ERR_DUPLICATEDUSER)
	}

	q := data.NewConn(wr, "users")
	q.Put(nu)

	// Add a UserTags entry for each tag for this user
	addUserTags(wr, nu)

	return nil
}

func updateUser(wr srv.WrapperRequest, nu *NUser) error {

	/*if err := nu.IsValid(); err != nil {
		return err
	}*/

	old, err := getUserById(wr, fmt.Sprintf("%d", nu.Id))
	if err != nil {
		return errors.New(ERR_USERNOTFOUND)
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

func deleteUser(wr srv.WrapperRequest, nu *NUser) error {
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

func getUserByMail(wr srv.WrapperRequest, email string) (*NUser, error) {
	nus := NewNUserBuffer()
	nu := new(NUser)

	q := data.NewConn(wr, "users")
	q.AddFilter("Mail =", email)
	q.GetMany(&nus)
	if len(nus) == 0 {
		return nu, errors.New(ERR_USERNOTFOUND)
	}
	nu = nus[0]
	nu.Tags, _ = getUserTags(wr, nu)

	return nu, nil
}

func getUserById(wr srv.WrapperRequest, s_id string) (*NUser, error) {
	nu := new(NUser)

	id, err := strconv.ParseInt(s_id, 10, 64)
	if err != nil {
		return nu, errors.New(ERR_USERNOTFOUND)
	}

	nu.Id = id
	q := data.NewConn(wr, "users")
	q.Get(nu)
	nu.Tags, _ = getUserTags(wr, nu)

	return nu, nil
}
