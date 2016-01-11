package http

import (
	"net/http"

	"model/answers"
	"model/questions"
	"model/tests"
	"model/users"
)

func init() {
	http.HandleFunc("/welcome", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, Welcome)
	})
	http.HandleFunc("/help", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, Help)
	})

	// Questions
	http.HandleFunc("/questions/main", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, questions.Main)
	})
	http.HandleFunc("/questions/list", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, questions.GetList)
	})
	http.HandleFunc("/questions/get", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, questions.GetOne)
	})
	http.HandleFunc("/questions/new", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, questions.New)
	})
	http.HandleFunc("/questions/add", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, questions.Add)
	})
	http.HandleFunc("/questions/tags/list", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, questions.GetTagsList)
	})
	http.HandleFunc("/questions/solve", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, questions.Solve)
	})

	// Answers
	http.HandleFunc("/answers/add", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, answers.Add)
	})

	// Tests
	http.HandleFunc("/tests/new", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, tests.New)
	})
	http.HandleFunc("/tests/edit", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, tests.Edit)
	})
	http.HandleFunc("/tests/update", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, tests.Update)
	})
	http.HandleFunc("/tests/add", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, tests.Add)
	})
	http.HandleFunc("/tests/list", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, tests.GetList)
	})
	http.HandleFunc("/tests/get", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, tests.GetOne)
	})
	http.HandleFunc("/tests/users/list", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, tests.GetUsersList)
	})
	http.HandleFunc("/tests/exercises/list", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, tests.GetExercisesList)
	})
	http.HandleFunc("/tests/tags/list", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, tests.GetTagsList)
	})

	// Users routes

	http.HandleFunc("/users/main", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, users.Main)
	})
	http.HandleFunc("/users/get", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, users.GetOne)
	})
	http.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, users.GetOne)
	})
	http.HandleFunc("/users/list", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, users.GetList)
	})
	http.HandleFunc("/users/new", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, users.New)
	})
	http.HandleFunc("/users/add", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, users.Add)
	})
	http.HandleFunc("/users/edit", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, users.Edit)
	})
	http.HandleFunc("/users/update", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, users.Update)
	})
	http.HandleFunc("/users/delete", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, users.Delete)
	})
	http.HandleFunc("/users/import", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, users.Import)
	})
	http.HandleFunc("/users/tags/list", func(w http.ResponseWriter, r *http.Request) {
		AppHandler(w, r, users.GetTagsList)
	})
}
