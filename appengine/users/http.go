package users

import (
	"fmt"
	"net/http"
	"errors"
	//"encoding/json"
	//"strconv"
	"html/template"

	"app"
)



// Routes

func init() {
	http.HandleFunc("/users/logout", app.AppLogout)
	http.HandleFunc("/users/get", getOneUser)
	http.HandleFunc("/users/list", getListUsers)
	/*http.HandleFunc("/users/add", addUser)
	http.HandleFunc("/users/new", newUserForm)
	http.HandleFunc("/users/list",allUsersForm)*/



}


// Templates

var listTmpl = template.Must(template.ParseFiles("app/tmpl/base.html",
	"appengine/users/tmpl/list.html"))

var newTmpl = template.Must(template.ParseFiles("app/tmpl/base.html",
	"appengine/users/tmpl/new.html"))

var viewTmpl = template.Must(template.ParseFiles("app/tmpl/base.html",
	"appengine/users/tmpl/view.html"))



func getOneUser (w http.ResponseWriter, r *http.Request) {
	wr:=app.NewWrapperRequest(r)
	CheckPerm(w,wr,OP_VIEW)

	wr.R.ParseForm()

	nus,err:=Get(wr,wr.R.Form)
	if len(nus)==0 || err!=nil{
		app.AppError(wr,w,errors.New("Usuario no encontrado"))
		return
	}

	user,err:=GetCurrentUser(wr)
	if err!=nil{
		app.AppError(wr,w,err)
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
	if err:=CheckPerm(w,wr,OP_ADMIN); err!=nil{
		return
	}

	user,err:=GetCurrentUser(wr)
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

/*
	admins,err:=GetUsersByRole(wr,fmt.Sprintf("%d",ROLE_ADMIN))
	teachers,err:=GetUsersByRole(wr,fmt.Sprintf("%d",ROLE_TEACHER))
	students,err:=GetUsersByRole(wr,fmt.Sprintf("%d",ROLE_STUDENT))
*/

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







func CheckPerm(w http.ResponseWriter, wr app.WrapperRequest, op byte)(error) {

	if wr.U == nil {
		wr.C.Infof("Not user session founded to check perm")
		app.RedirectUserLogin(w,wr.R)
		return errors.New("user not exits")
	}

	if (!wr.IsAdminRequest()){
		// Si no es admin, deberiamos buscarlo en nuestra base
		// de datos de usuarios permitidos y comprobar si 
		// con su rol puede hacer dicha operación
		// De esa busqueda calculamos la variable perm y la comparamos
		// con op

		nu,err:=GetCurrentUser(wr)
		if err!=nil{
			app.RedirectUserLogin(w,wr.R)
			return err
		}

		if !IsAllowed(nu.Role,op){
			wr.C.Infof("Perm:"+fmt.Sprintf("%b",nu.Role)+" "+fmt.Sprintf("%b",op))
	                wr.C.Infof("User "+nu.Mail+" failed allowed access")
			app.RedirectUserLogin(w,wr.R)
			return err
		}

		//RedirectUserLogin(w,r)
		//return errors.New("user has not perm for the operation")
	}

	// Si es admin puede cualquier cosa
	return nil
}



