package users

import (
	"strings"
	"time"
)

const (
	ROLE_GUEST   = iota
	ROLE_STUDENT = iota
	ROLE_TEACHER = iota
	ROLE_ADMIN   = iota

	//deprecated:
	OP_VIEWER    = 1 // ver tests
	OP_MAKER     = 2 // hacer tests
	OP_COMMITTER = 4 // crear tests
	OP_ADMIN     = 8 // labores de administración (crear usuarios, ....)

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

type AppUser interface {
	ID() int64
	GetInfo() map[string]string
	GetRole() int8
	GetEmail() string
}

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

func GetDefaultUser(email string) AppUser {
	n := new(NUser)
	n.Id = -1
	n.Mail = email
	n.Name = "Administrador"
	n.Role = ROLE_ADMIN

	return n
}

/*
















 */

/*

func IsAllowed(userPerm int8, opMask byte) bool {
	return opMask == byte(userPerm)&opMask
}

func New(mail string, name string, role int8) NUser {
	nu := NUser{-1, mail, name, role, make([]string, 10), time.Now()}
	return nu
}

func (n NUser) IsAdmin() bool {
	return n.Role == ROLE_ADMIN
}

func (n NUser) IsTeacher() bool {
	return n.Role == ROLE_TEACHER
}

func (n NUser) IsStudent() bool {
	return n.Role == ROLE_STUDENT
}

func (n NUser) GetStringTags() string {
	return strings.Join(n.Tags, ",")
}

func (n NUser) GetStringRole() string {
	switch n.Role {
	case ROLE_ADMIN:
		return "Administrador"
	case ROLE_TEACHER:
		return "Profesor"
	case ROLE_STUDENT:
		return "Estudiante"
	}
	return ""
}

func (n NUser) GetStringTimeStamp() string {
	return n.TimeStamp.Format("02/01/2006")
}

func (n NUser) IsValid() (err error) {

	if len(n.Name) < 0 || len(n.Name) > MAXSZUSERNAME {
		return errors.New(ERR_NOTVALIDUSER)
	}

	if len(n.Mail) < 0 || len(n.Mail) > MAXSZUSERNAME {
		return errors.New(ERR_NOTVALIDUSER)
	}

	if n.Role != ROLE_GUEST && n.Role != ROLE_TEACHER &&
		n.Role != ROLE_STUDENT && n.Role != ROLE_ADMIN {
		return errors.New(ERR_NOTVALIDUSER)
	}

	return
}
*/
