package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	appusers "app/users"
	"model/users"

	"appengine/data"
	"appengine/srv"
)

const (
	JSON_ACCEPT_HEADER = "application/json"
)

type ErrorResponse struct {
	Error string
}

type WrapperHandler func(wr srv.WrapperRequest, tc map[string]interface{}) (string, error)

func AppHandler(w http.ResponseWriter, r *http.Request, whandler WrapperHandler) {
	wr := srv.NewWrapperRequest(r)
	err := getCurrentUser(&wr)
	if err != nil {
		RedirectToLogin(w, wr.R)
		return
	}

	rformat := r.Header.Get("Accept")
	wr.JsonResponse = (strings.Index(rformat, JSON_ACCEPT_HEADER) >= 0)

	// Perform the Handler
	rdata := make(map[string]interface{})
	rdata["User"] = wr.NU
	tmplName, err := whandler(wr, rdata)
	if err != nil {
		errorResponse(wr, w, err)
		return
	}

	if wr.JsonResponse {
		// Json Response
		jbody, err := json.Marshal(rdata["Content"])
		if err != nil {
			errorResponse(wr, w, err)
			return
		}
		fmt.Fprintf(w, "%s", string(jbody[:len(jbody)]))

	} else {
		// HTML Response
		// Template funcions
		funcMap := template.FuncMap{
			"inc": func(i int) int {
				return i + 1
			},
		}

		tmpl := template.Must(template.New("base.html").Funcs(funcMap).ParseFiles("app/tmpl/base.html", tmplName))
		if err := tmpl.Execute(w, rdata); err != nil {
			errorResponse(wr, w, err)
			return
		}
	}

}

func errorResponse(wr srv.WrapperRequest, w http.ResponseWriter, err error) {
	wr.C.Errorf("%v", err)

	if wr.JsonResponse {
		// Json Response
		jbody, err := json.Marshal(ErrorResponse{err.Error()})
		if err != nil {
			srv.Log(wr, err.Error())
			return
		}
		fmt.Fprintf(w, "%s", string(jbody[:len(jbody)]))
	} else {
		// HTML Response
		errorTmpl := template.Must(template.ParseFiles("app/tmpl/error.html"))
		if err := errorTmpl.Execute(w, err.Error()); err != nil {
			wr.C.Errorf("%v", err)
			return
		}
	}
}

func getCurrentUser(wr *srv.WrapperRequest) error {

	//var nu nusers.NUser
	var nu appusers.AppUser

	nus := users.NewNUserBuffer()
	u := wr.U
	if u == nil {
		return errors.New("No user session founded")
	}

	q := data.NewConn(*wr, "users")
	q.AddFilter("Mail =", u.Email)
	err := q.GetMany(&nus)

	if len(nus) <= 0 {
		// If the session users not in the datastore we use de admin
		// users of the app
		if wr.IsAdminRequest() {
			//nu = nusers.New(u.Email, "Administrador", nusers.ROLE_ADMIN)
			nu = GetDefaultUser(u.Email)
		} else {
			return errors.New("No user id found")
		}
	} else {
		nu = nus[0]
	}

	wr.NU = nu
	return err
}

func GetDefaultUser(email string) appusers.AppUser {
	n := new(users.NUser)
	n.Id = -1
	n.Mail = email
	n.Name = "Administrador"
	n.Role = users.ROLE_ADMIN

	return n
}
