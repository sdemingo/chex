package users

import (
	"fmt"
	"net/http"
	"app"
	"errors"
	"encoding/json"
	//"strconv"
	"html/template"

	"appengine"
	"appengine/user"
	"appengine/datastore"
)

// Routes

func init() {
	http.HandleFunc("/users/logout", logout)
	http.HandleFunc("/users/get", getUser)
	http.HandleFunc("/users/edit", addUser)
	http.HandleFunc("/users/all",allUsers)

	http.HandleFunc("/users/init", initUsers)

}


// Templates

var listTmpl = template.Must(template.ParseFiles("app/tmpl/base.html",
	"appengine/users/tmpl/list.html"))

//var newTmpl = template.Must(template.ParseFiles("app/tmpl/base.html",
//	"appengine/tests/tmpl/new.html"))




func getUser (w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	CheckPerm(w,r,OP_VIEW)

	r.ParseForm()

	var nu NUser
	var err error

	if r.Form["id"]!=nil{
		nu,err=GetUserById(c,r.Form["id"][0])
	}

	if r.Form["mail"]!=nil{
		nu,err=GetUserById(c,r.Form["mail"][0])
	}

	if err!=nil{
		app.ServeError(c,w,err)
		return
	}

	jbody,err:=json.Marshal(nu)
	if err != nil {
		app.ServeError(c,w,err)
		return
	}

	fmt.Fprintf(w, "%s", string(jbody[:len(jbody)]))
}



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

	key := datastore.NewKey(c, "users", "", 0, nil)
	key, err = datastore.Put(c, key, &nu)
	nu.Id = key.IntID()

	jbody,err:=json.Marshal(nu)
	if err != nil {
		app.ServeError(c,w,err)
		return
	}

	fmt.Fprintf(w, "%s", string(jbody[:len(jbody)]))
}



func allUsers(w http.ResponseWriter, r *http.Request) {

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





func logout (w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	if u == nil {
		RedirectUserLogin(w,r)
		return
	}

	url, err := user.LogoutURL(c, r.URL.String())
	if err != nil {
		app.ServeError(c,w,err)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
}





func RedirectUserLogin(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	url, err := user.LoginURL(c, r.URL.String())
	if err != nil {
		app.ServeError(c,w,err)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
}





func CheckPerm(w http.ResponseWriter, r *http.Request, op byte)(error) {

	c := appengine.NewContext(r)
	u := user.Current(c)

	if u == nil {
		c.Infof("Not user session founded to check perm")
		RedirectUserLogin(w,r)
		return errors.New("user not exits")
	}

	if (!user.IsAdmin(c)){
		// Si no es admin, deberiamos buscarlo en nuestra base
		// de datos de usuarios permitidos y comprobar si 
		// con su rol puede hacer dicha operación
		// De esa busqueda calculamos la variable perm y la comparamos
		// con op

		nu,err:=GetCurrentUser(c)
		if err!=nil{
			RedirectUserLogin(w,r)
			return err
		}

		if !IsAllowed(nu.Role,op){
			c.Infof("Perm:"+fmt.Sprintf("%b",nu.Role)+" "+fmt.Sprintf("%b",op))
	                c.Infof("User "+nu.Mail+" failed allowed access")
			RedirectUserLogin(w,r)
			return err
		}

		//RedirectUserLogin(w,r)
		//return errors.New("user has not perm for the operation")
	}

	// Si es admin puede cualquier cosa
	return nil
}






/***************************************************************************/

/*******                     Debug!!!!                               *******/  

/***************************************************************************/



func initUsers(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	nu:=new(NUser)
	nu.Mail="admin@chex"
	nu.Name="Administrador"
	nu.Role=ROLE_ADMIN

	key := datastore.NewKey(c, "users", "", 0, nil)
	key, _ = datastore.Put(c, key, nu)
	nu.Id = key.IntID()

	nu=new(NUser)
	nu.Mail="teacher@chex"
	nu.Name="Profesor"
	nu.Role=ROLE_TEACHER

	key = datastore.NewKey(c, "users", "", 0, nil)
	key, _ = datastore.Put(c, key, nu)
	nu.Id = key.IntID()


	nu=new(NUser)
	nu.Mail="student@chex"
	nu.Name="Estudiante"
	nu.Role=ROLE_STUDENT

	key = datastore.NewKey(c, "users", "", 0, nil)
	key, _ = datastore.Put(c, key, nu)
	nu.Id = key.IntID()

}





