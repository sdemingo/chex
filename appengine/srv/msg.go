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
