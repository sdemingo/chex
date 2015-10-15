package app


import (
	"net/http"
	"appengine"
	"encoding/json"
	"fmt"

)


type ErrorResponse struct{
	Error string
}


func init() {
	http.HandleFunc("/", root)
}


func ServeError(c appengine.Context, w http.ResponseWriter, err error) {
	c.Errorf("%v", err)
	er := &ErrorResponse{err.Error()}
	js,err:=json.Marshal(er)
	fmt.Fprintf(w, "%s", string(js[:len(js)]))
}



func root(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w,r,"/test/all",http.StatusMovedPermanently)
}

