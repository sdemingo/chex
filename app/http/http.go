package http

import (
	"net/http"

	"model/users"

	"appengine/srv"
)

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/logout", logout)
}

func root(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/welcome", http.StatusMovedPermanently)
}

func logout(w http.ResponseWriter, r *http.Request) {
	srv.RedirectUserLogin(w, r)
}

func RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	srv.RedirectUserLogin(w, r)
}

var adminTmpl = "app/tmpl/adminWelcome.html"
var studentTmpl = "app/tmpl/studentWelcome.html"
var teacherTmpl = "app/tmpl/teacherWelcome.html"
var helpTmpl = "app/tmpl/help.html"

func Welcome(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	// err := srv.CheckPerm(wr, users.OP_VIEWER)
	// if err != nil {
	// 	return "", errors.New(users.ERR_NOTOPERATIONALLOWED)
	// }

	// añadir más información a tc

	//if wr.NU.IsAdmin() {
	if wr.NU.GetRole() == users.ROLE_ADMIN {
		return adminTmpl, nil
	}

	//if wr.NU.IsTeacher() {
	if wr.NU.GetRole() == users.ROLE_TEACHER {
		return teacherTmpl, nil
	}

	return studentTmpl, nil
}

func Help(wr srv.WrapperRequest, tc map[string]interface{}) (string, error) {
	return helpTmpl, nil
}
