package srv

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"app/users"

	"appengine/user"
)

const (
	JSON_ACCEPT_HEADER = "application/json"
)

type WrapperHandler func(wr WrapperRequest, tc map[string]interface{}) (string, error)

func AppHandler(w http.ResponseWriter, r *http.Request, whandler WrapperHandler) {
	wr := NewWrapperRequest(r)
	err := wr.GetCurrentUser()
	if err != nil {
		RedirectUserLogin(w, wr.R)
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
		tmpl := template.Must(template.ParseFiles("app/tmpl/base.html",
			tmplName))

		if err := tmpl.Execute(w, rdata); err != nil {
			errorResponse(wr, w, err)
			return
		}
	}

}

func errorResponse(wr WrapperRequest, w http.ResponseWriter, err error) {
	wr.C.Errorf("%v", err)

	if wr.JsonResponse {
		// Json Response
		jbody, err := json.Marshal(ErrorResponse{err.Error()})
		if err != nil {
			wr.C.Errorf("%v", err)
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

func Log(wr WrapperRequest, msg string) {
	wr.C.Infof("%s", msg)
}

func RedirectUserLogin(w http.ResponseWriter, r *http.Request) {
	wr := NewWrapperRequest(r)
	var url string
	var err error
	if wr.U != nil {
		url, err = user.LogoutURL(wr.C, "/")
	} else {
		url, err = user.LoginURL(wr.C, wr.R.URL.String())
	}
	if err != nil {
		errorResponse(wr, w, err)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
}

func CheckPerm(wr WrapperRequest, op byte) error {

	n := wr.NU
	if !n.IsAdmin() {
		// Si no es admin, deberiamos buscarlo en nuestra base
		// de datos de usuarios permitidos y comprobar si
		// con su rol puede hacer dicha operación
		// De esa busqueda calculamos la variable perm y la comparamos
		// con op

		if !users.IsAllowed(n.Role, op) {
			Log(wr, "Perm:"+fmt.Sprintf("%b", n.Role)+" "+fmt.Sprintf("%b", op))
			Log(wr, "User "+n.Mail+" failed allowed access")
			return errors.New("Operation not allowed")
		}
	}

	// Si es admin puede cualquier cosa
	return nil
}
