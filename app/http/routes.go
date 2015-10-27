package http

import (
	"net/http"

	"appengine/questions"
	"appengine/srv"
	//"appengine/tests"
	"appengine/users"
)

func init() {
	http.HandleFunc("/welcome", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w, r, Welcome)
	})

	http.HandleFunc("/logout", srv.AppLogout)

	// Questions
	http.HandleFunc("/questions/list", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w, r, questions.GetList)
	})
	http.HandleFunc("/questions/new", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w, r, questions.New)
	})

	// Users routes

	http.HandleFunc("/users/get", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w, r, users.GetOne)
	})
	http.HandleFunc("/users/list", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w, r, users.GetList)
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
	http.HandleFunc("/users/tags/list", func(w http.ResponseWriter, r *http.Request) {
		srv.AppHandler(w, r, users.GetTagsList)
	})

}
