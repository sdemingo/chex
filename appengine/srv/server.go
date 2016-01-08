package srv

import (
	"net/http"

	"app/users"

	"appengine"
	"appengine/user"
)

type WrapperRequest struct {
	R  *http.Request
	C  appengine.Context
	U  *user.User
	NU users.AppUser
	//NU           *users.NUser

	JsonResponse bool
}

func NewWrapperRequest(r *http.Request) WrapperRequest {
	c := appengine.NewContext(r)
	return WrapperRequest{r, c, user.Current(c), nil, false}
}

func (wr WrapperRequest) IsAdminRequest() bool {
	return user.IsAdmin(wr.C)
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
		//errorResponse(wr, w, err)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
}
