package http


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
	http.HandleFunc("/users/get", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w,r,users.GetOne)
	})
	http.HandleFunc("/users/new", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w,r,users.New)
	})
	http.HandleFunc("/users/add", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w,r,users.Add)
	})


	http.HandleFunc("/logout", srv.AppLogout)
}
