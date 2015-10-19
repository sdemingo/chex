package users


import (
	"fmt"
	"net/http"
	"errors"
	"html/template"

	"app"
)


/*

// Routes

func init() {
	http.HandleFunc("/users/logout", app.AppLogout)
	http.HandleFunc("/users/get", getOneUser)
	http.HandleFunc("/users/list", getListUsers)
}

*/


// Templates

var listTmpl = "appengine/users/tmpl/list.html"
var newTmpl  = "appengine/users/tmpl/new.html"
var viewTmpl = "appengine/users/tmpl/view.html"



func getOneUser (w http.ResponseWriter, r *http.Request) {
	wr:=app.NewWrapperRequest(r)
	user,err:=GetCurrentUser(wr)
	if err!=nil{
		app.RedirectUserLogin(w,wr.R)
		return
	}
	err=user.CheckPerm(wr, OP_ADMIN)
	if err!=nil{
		app.AppError(wr,w,err)
		return
	}


	wr.R.ParseForm()
	nus,err:=Get(wr,wr.R.Form)
	if len(nus)==0 || err!=nil{
		app.AppError(wr,w,errors.New("Usuario no encontrado"))
		return
	}


	
	tc := make(map[string]interface{})
	tc["User"] = user
	tc["Content"] = nus[0]

	if err := viewTmpl.Execute(w, tc); err != nil {
		app.AppError(wr,w,err)
		return
	}
}


func getListUsers (w http.ResponseWriter, r *http.Request) {
	wr:=app.NewWrapperRequest(r)
	user,err:=GetCurrentUser(wr)
	if err!=nil{
		app.RedirectUserLogin(w,wr.R)
		return
	}
	err=user.CheckPerm(wr, OP_ADMIN)
	if err!=nil{
		app.AppError(wr,w,err)
		return
	}


	filters:=make(map[string][]string)
	filters["role"]=[]string{fmt.Sprintf("%d",ROLE_ADMIN)}
	admins,err:=Get(wr,filters)

	filters["role"]=[]string{fmt.Sprintf("%d",ROLE_TEACHER)}
	teachers,err:=Get(wr,filters)

	filters["role"]=[]string{fmt.Sprintf("%d",ROLE_STUDENT)}
	students,err:=Get(wr,filters)


	tc := make(map[string]interface{})
	tc["User"] = user
	tc["Admins"] = admins
	tc["Teachers"] = teachers
	tc["Students"] = students

	if err := listTmpl.Execute(w, tc); err != nil {
		app.AppError(wr,w,err)
		return
	}
}
*/




















/*
func addUser(w http.ResponseWriter, r *http.Request) {
	var err error

	c := appengine.NewContext(r)
	if err:=CheckPerm(w,r,OP_ADMIN); err!=nil{
		return
	}

	var nu NUser 

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&nu)
	if err != nil {
		app.ServeError(c,w,err)
		return
	}

	if err = nu.IsValid(); err!=nil{
		app.ServeError(c,w,err)
		return
	}

	_, err = GetUserByMail(c,nu.Mail)
	if err==nil{
		app.ServeError(c,w,errors.New("User duplicated"))
		return
	}

	key := datastore.NewKey(c, "users", "", 0, nil)
	key, err = datastore.Put(c, key, &nu)
	if err!=nil{
		app.ServeError(c,w,err)
		return
	}
	
	nu.Id = key.IntID()

	jbody,err:=json.Marshal(nu)
	if err != nil {
		app.ServeError(c,w,err)
		return
	}

	fmt.Fprintf(w, "%s", string(jbody[:len(jbody)]))
}






func allUsersForm(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	if err:=CheckPerm(w,r,OP_ADMIN); err!=nil{
		return
	}

	user,err:=GetCurrentUser(c)
	if err!=nil{
		app.ServeError(c,w,err)
		return
	}
	admins,err:=GetUsersByRole(c,fmt.Sprintf("%d",ROLE_ADMIN))
	teachers,err:=GetUsersByRole(c,fmt.Sprintf("%d",ROLE_TEACHER))
	students,err:=GetUsersByRole(c,fmt.Sprintf("%d",ROLE_STUDENT))


	tc := make(map[string]interface{})
	tc["User"] = user
	tc["Admins"] = admins
	tc["Teachers"] = teachers
	tc["Students"] = students

	if err := listTmpl.Execute(w, tc); err != nil {
		app.ServeError(c,w,err)
		return
	}
}


func newUserForm(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if err:=CheckPerm(w,r,OP_ADMIN); err!=nil{
		return
	}

	user,err:=GetCurrentUser(c)
	if err!=nil{
		app.ServeError(c,w,err)
		return
	}

	tc := make(map[string]interface{})
	tc["User"] = user

	if err := newTmpl.Execute(w, tc); err != nil {
		app.ServeError(c,w,err)
		return
	}
}




*/




