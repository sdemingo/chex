package srv

import (
	"errors"
	"net/http"

	"app/users"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type ErrorResponse struct {
	Error string
}

type WrapperRequest struct {
	R            *http.Request
	C            appengine.Context
	U            *user.User
	NU           *users.NUser
	JsonResponse bool
}

func NewWrapperRequest(r *http.Request) WrapperRequest {
	c := appengine.NewContext(r)
	return WrapperRequest{r, c, user.Current(c), nil, false}
}

func (wr WrapperRequest) IsAdminRequest() bool {
	return user.IsAdmin(wr.C)
}

func (wr *WrapperRequest) GetCurrentUser() error {

	// Busco información del usuario de sesión

	var nu users.NUser
	u := wr.U
	if u == nil {
		return errors.New("No user session founded")
	}

	q := datastore.NewQuery("users").Filter("Mail =", u.Email)
	var nusers []users.NUser
	keys, _ := q.GetAll(wr.C, &nusers)

	if len(nusers) <= 0 {
		// El usuario de sesión no esta en el datastore
		// Usamos el admin de la app aunque no este en el datastore
		if wr.IsAdminRequest() {
			nu = users.New(u.Email, "Administrador", users.ROLE_ADMIN)
		} else {
			return errors.New("No user id found")
		}
	} else {
		nu = nusers[0]
		nu.Id = keys[0].IntID()
	}

	wr.NU = &nu
	return nil
}
