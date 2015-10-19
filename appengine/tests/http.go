package tests

import (
	"errors"

	"users" 
	"appengine/srv" 
)




// Templates

var listTmpl = "appengine/tests/tmpl/list.html"
var newTmpl  = "appengine/tests/tmpl/new.html"



// Handlers

func GetAll (wr srv.WrapperRequest, tc map[string]interface{}) (string, error){

	err:=srv.CheckPerm(wr, users.OP_VIEW)
	if err!=nil{
		return listTmpl, errors.New(users.ERR_NOTOPERATIONALLOWED)
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
	
	tc["Content"] = tests

	return listTmpl, nil
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




