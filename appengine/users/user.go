package users

import (
	"errors"
	"strconv"

	"appengine"
	"appengine/datastore"
	"appengine/user"
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


func (n NUser)IsAdmin()(bool){
	return n.Role == ROLE_ADMIN
}

func (n NUser)IsTeacher()(bool){
	return n.Role == ROLE_TEACHER || n.Role == ROLE_ADMIN
}

func (n NUser)IsStudent()(bool){
	return n.Role == ROLE_STUDENT || n.Role == ROLE_ADMIN || n.Role == ROLE_TEACHER
}


func (n NUser)IsValid()(err error){

	if len(n.Name)<0 || len(n.Name)>MAXSZUSERNAME{
		return errors.New(ERR_NOTVALIDUSER)
	}

	if len(n.Mail)<0 || len(n.Mail)>MAXSZUSERNAME{
		return errors.New(ERR_NOTVALIDUSER)
	}

	return
}





func GetUserByMail(c appengine.Context, email string)(NUser, error){
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


func GetUserById(c appengine.Context, s_id string)(NUser, error){

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



func GetCurrentUser(c appengine.Context)(NUser, error){

	// Busco información del usuario de sesión

	var nu NUser

	u := user.Current(c)
	if u == nil {
		return nu, errors.New("Session user not found")
	}

	q := datastore.NewQuery("users").Filter("Mail =",u.Email)
	var nusers []NUser
	keys,_:= q.GetAll(c, &nusers)
	
	if (len(nusers)<=0){
		// El usuario de sesión no esta en el datastore
		// Usamos el admin de la app aunque no este en el datastore
		if user.IsAdmin(c){
			nu = NUser{-1,u.Email,"Administrador",ROLE_ADMIN}
		}else{
			return nu, errors.New("No user id found")
		}
	}else{
		nu=nusers[0]
		nu.Id=keys[0].IntID()
	}
	return nu,nil
}



func GetUsersByRole(c appengine.Context, s_role string)([]NUser, error){
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
	//nu = nus[0]
	//nu.Id =  keys[0].IntID()

	return nus,nil

}






