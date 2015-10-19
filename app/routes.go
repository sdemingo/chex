package app


import (
	"net/http"

	"appengine/srv"
	"appengine/tests"
	"appengine/users"
)



func init() {
	http.HandleFunc("/test/all", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w,r,tests.GetAll)
	})

	http.HandleFunc("/users/all", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w,r,users.GetAll)
	})

	http.HandleFunc("/logout", srv.AppLogout)
}

