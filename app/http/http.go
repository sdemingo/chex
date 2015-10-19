package http


import (
	"net/http"
)




func init() {
	http.HandleFunc("/", root)
}



func root(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w,r,"/test/all",http.StatusMovedPermanently)
}




