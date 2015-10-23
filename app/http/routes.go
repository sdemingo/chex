package http

import (
	"net/http"

	"appengine/srv"
	"appengine/tests"
	"appengine/users"
)

func init() {
	http.HandleFunc("/welcome", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w, r, Welcome)
	})

	http.HandleFunc("/logout", srv.AppLogout)

	// Test routes
	http.HandleFunc("/test/all", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w, r, tests.GetAll)
	})

	// Users routes
	http.HandleFunc("/users/all", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w, r, users.GetAll)
	})
	http.HandleFunc("/users/get", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w, r, users.GetOne)
	})
	http.HandleFunc("/users/new", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w, r, users.New)
	})
	http.HandleFunc("/users/add", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w, r, users.Add)
	})
	http.HandleFunc("/users/edit", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w, r, users.Edit)
	})
	http.HandleFunc("/users/update", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w, r, users.Update)
	})
	http.HandleFunc("/users/delete", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w, r, users.Delete)
	})

}
