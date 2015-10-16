package tests

import (
	//"fmt"
	"net/http"
	"html/template"

	"app"
	"appengine/users"  //my users
)

// http://stackoverflow.com/questions/9573644/go-appengine-how-to-structure-templates-for-application

// Routes

func init() {
	http.HandleFunc("/test/all", getAllTest)
	//http.HandleFunc("/test/new", newTest)
}





// Templates

var listTmpl = template.Must(template.ParseFiles("app/tmpl/base.html",
	"appengine/tests/tmpl/list.html"))

var newTmpl = template.Must(template.ParseFiles("app/tmpl/base.html",
	"appengine/tests/tmpl/new.html"))



// Handlers

func getAllTest (w http.ResponseWriter, r *http.Request) {
	wr:=app.NewWrapperRequest(r)

	nu,err:=users.GetCurrentUser(wr)
	if err!=nil{
		app.RedirectUserLogin(w,r)
		return
	}
	
	if err:=users.CheckPerm(w,wr,users.OP_VIEW); err!=nil{
		return
	}

	// MOcks:
	tests := make([]Test,10)
	tests[0].Title="Test de prueba1"
	tests[0].Author="Sergio de Mingo"

	tests[1].Title="Test de prueba2"
	tests[1].Author="Sergio de Mingo"

	tests[2].Title="Test de prueba3"
	tests[2].Author="Sergio de Mingo"

	tests[3].Title="Test de prueba4"
	tests[3].Author="Sergio de Mingo"
	

	tc := make(map[string]interface{})
	tc["User"] = nu
	tc["Content"] = tests

	if err := listTmpl.Execute(w, tc); err != nil {
		app.AppError(wr,w,err)
		return
	}
}


/*
func newTest (w http.ResponseWriter, r *http.Request) {
	
	c := appengine.NewContext(r)

	nu,err:=users.GetCurrentUser(c)
	if err!=nil{
		app.ServeError(c,w,err)
		return
	}

	if err:=users.CheckPerm(w,r,users.OP_UPDATE); err!=nil{
		return
	}

	test := new(Test)
	test.Author=nu.Mail

	tc := make(map[string]interface{})
	tc["User"] = nu
	tc["Content"] = test

	if err := newTmpl.Execute(w, tc); err != nil {
		app.ServeError(c,w,err)
		return
	}
}

*/




