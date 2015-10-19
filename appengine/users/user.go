package users

import (
	"errors"
	"strconv"
	//"fmt"

	"appengine"
	"appengine/datastore"
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
)


type NUser struct{
	Id int64     `datastore:"-"`  // ignored by datastore
	Mail string  
	Name string
	Role int8    `json:",string"`
}



func IsAllowed(userPerm int8, opMask byte)(bool){
	return opMask == byte(userPerm) & opMask
}


func Get(c appengine.Context, filters map[string][]string)(nus []NUser,err error){
	if filters["id"]!=nil{
		nu,err:=getUserById(c,filters["id"][0])
		nus:=make([]NUser,1)
		nus[0]=nu
		return nus,err
	}

	if filters["mail"]!=nil{
		nu,err:=getUserById(c,filters["mail"][0])
		nus:=make([]NUser,1)
		nus[0]=nu
		return nus,err
	}

	if filters["role"]!=nil{
		nus,err:=getUsersByRole(c,filters["role"][0])
		return nus,err
	}

	return
}


func (n NUser) IsAdmin()(bool){
	return n.Role == ROLE_ADMIN
}

func (n NUser) IsTeacher()(bool){
	return n.Role == ROLE_TEACHER || n.Role == ROLE_ADMIN
}

func (n NUser) IsStudent()(bool){
	return n.Role == ROLE_STUDENT || n.Role == ROLE_ADMIN || n.Role == ROLE_TEACHER
}

func (n NUser) GetEmail()(string){
	return n.Mail
}


func (n NUser) IsValid()(err error){

	if len(n.Name)<0 || len(n.Name)>MAXSZUSERNAME{
		return errors.New(ERR_NOTVALIDUSER)
	}

	if len(n.Mail)<0 || len(n.Mail)>MAXSZUSERNAME{
		return errors.New(ERR_NOTVALIDUSER)
	}

	return
}










func getUserByMail(c appengine.Context, email string)(NUser, error){
	var nus []NUser
	var nu NUser
	
	q := datastore.NewQuery("users").Filter("Mail =", email)

	keys, err := q.GetAll(c, &nus)
	if (len(keys)==0) || err!=nil{
		return nu, errors.New("User not found. Bad mail")
	}
	nu = nus[0]
	nu.Id =  keys[0].IntID()

	return nu,nil
}


func getUserById(c appengine.Context, s_id string)(NUser, error){

	var nu NUser

	id,err:=strconv.ParseInt(s_id,10,64)
	if err!=nil{
		return nu, errors.New("User not found. Bad ID")
	}

	if id!=0{
		k := datastore.NewKey(c, "users", "", id, nil)
		datastore.Get(c, k, &nu)
	}else{
		return nu, errors.New("User not found. Bad ID")
	}

	return nu,nil
}


func getUsersByRole(c appengine.Context, s_role string)([]NUser, error){
	var nus []NUser

	role,err:=strconv.ParseInt(s_role,10,64)
	if err!=nil{
		return nus, errors.New("User role bad formatted")
	}
	
	q := datastore.NewQuery("users").Filter("Role =", role)

	keys, err := q.GetAll(c, &nus)
	if (len(keys)==0) || err!=nil{
		return nus, errors.New("User not found. Bad role")
	}

	for i:=0;i<len(nus);i++{
		nus[i].Id=keys[0].IntID()
	}

	return nus,nil

}

