package srv

import (
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
