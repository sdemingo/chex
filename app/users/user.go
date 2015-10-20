package users

import(
	"errors"
	"strings"
)

const (
	ROLE_ADMIN = 15
	ROLE_TEACHER = 7   
	ROLE_STUDENT = 3
	ROLE_GUEST = 1
	
	OP_VIEW = 1    // ver tests
	OP_CHECKIN = 2 // hacer tests
	OP_UPDATE = 4  // crear tests
	OP_ADMIN = 8   // labores de administración (crear usuarios, ....)

	MAXSZUSERNAME = 100

	ERR_NOTVALIDUSER = "Usuario no valido"
	ERR_NOTOPERATIONALLOWED = "Operación no permitida"
)


type NUser struct{
	Id int64     `json:",string" datastore:"-"`
	Mail string  
	Name string
	Role int8    `json:",string"`
	Tags []string
}



func IsAllowed(userPerm int8, opMask byte)(bool){
	return opMask == byte(userPerm) & opMask
}

func New(mail string, name string, role int8)(NUser) {
	nu := NUser{-1,mail,name,role,make([]string,10)}
	return nu
}


func (n  NUser) IsAdmin()(bool){
	return n.Role == ROLE_ADMIN
}

func (n  NUser) IsTeacher()(bool){
	return n.Role == ROLE_TEACHER 
}

func (n  NUser) IsStudent()(bool){
	return n.Role == ROLE_STUDENT
}

func (n  NUser) GetEmail()(string){
	return n.Mail
}

func (n NUser) GetStringTags()(string){
	return strings.Join(n.Tags,",")
}


func (n  NUser) IsValid()(err error){

	if len(n.Name)<0 || len(n.Name)>MAXSZUSERNAME{
		return errors.New(ERR_NOTVALIDUSER)
	}

	if len(n.Mail)<0 || len(n.Mail)>MAXSZUSERNAME{
		return errors.New(ERR_NOTVALIDUSER)
	}

	return
}
