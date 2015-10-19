package users

import (
	"errors"
	"strconv"
	//"fmt"

	"users"
	"appengine/srv"
	"appengine/datastore"
)





func getUsers(wr srv.WrapperRequest, filters map[string][]string)(nus []users.NUser,err error){
	if filters["id"]!=nil{
		nu,err:=getUserById(wr,filters["id"][0])
		nus:=make([]users.NUser,1)
		nus[0]=nu
		return nus,err
	}

	if filters["mail"]!=nil{
		nu,err:=getUserById(wr,filters["mail"][0])
		nus:=make([]users.NUser,1)
		nus[0]=nu
		return nus,err
	}

	if filters["role"]!=nil{
		nus,err:=getUsersByRole(wr,filters["role"][0])
		return nus,err
	}

	return
}








func getUserByMail(wr srv.WrapperRequest, email string)(users.NUser, error){
	var nus []users.NUser
	var nu users.NUser
	
	q := datastore.NewQuery("users").Filter("Mail =", email)

	keys, err := q.GetAll(wr.C, &nus)
	if (len(keys)==0) || err!=nil{
		return nu, errors.New("User not found. Bad mail")
	}
	nu = nus[0]
	nu.Id =  keys[0].IntID()

	return nu,nil
}


func getUserById(wr srv.WrapperRequest, s_id string)(users.NUser, error){

	var nu users.NUser

	id,err:=strconv.ParseInt(s_id,10,64)
	if err!=nil{
		return nu, errors.New("User not found. Bad ID")
	}

	if id!=0{
		k := datastore.NewKey(wr.C, "users", "", id, nil)
		datastore.Get(wr.C, k, &nu)
	}else{
		return nu, errors.New("User not found. Bad ID")
	}

	return nu,nil
}


func getUsersByRole(wr srv.WrapperRequest, s_role string)([]users.NUser, error){
	var nus []users.NUser

	role,err:=strconv.ParseInt(s_role,10,64)
	if err!=nil{
		return nus, errors.New("User role bad formatted")
	}
	
	q := datastore.NewQuery("users").Filter("Role =", role)

	keys, err := q.GetAll(wr.C, &nus)
	if (len(keys)==0) || err!=nil{
		return nus, errors.New("User not found. Bad role")
	}

	for i:=0;i<len(nus);i++{
		nus[i].Id=keys[i].IntID()
	}

	return nus,nil

}

