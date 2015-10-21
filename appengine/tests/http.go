package tests

import (
	"errors"

	"app/users" 
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




