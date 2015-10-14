package app


import (
	"net/http"
	"appengine"

)


func init() {
	http.HandleFunc("/", root)
}


func ServeError(c appengine.Context, w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	c.Errorf("%v", err)
}



func root(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w,r,"/test/all",http.StatusMovedPermanently)
}

