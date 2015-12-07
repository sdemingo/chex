package srv

import (
	"app/users"
	"errors"
	"fmt"
	"net/http"

	"appengine/user"
)

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
		//errorResponse(wr, w, err)
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
