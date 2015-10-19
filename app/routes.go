package app


import (
	"net/http"

	"appengine/srv"
	"appengine/tests"
)



func init() {
	http.HandleFunc("/test/all", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w,r,tests.GetAll)
	})

	http.HandleFunc("/logout", srv.AppLogout)
}

