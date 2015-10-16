package app


import (
	"net/http"
	"appengine"
	"appengine/user"
	"encoding/json"
	"fmt"
)


type ErrorResponse struct{
	Error string
}



type WrapperRequest struct{
	R *http.Request
	C appengine.Context
	U *user.User
}


func NewWrapperRequest(r *http.Request) (WrapperRequest) {
	c:=appengine.NewContext(r)
	return WrapperRequest{r, c, user.Current(c)}
}

func (wr WrapperRequest) IsAdminRequest()(bool){
	return user.IsAdmin(wr.C)
}



func init() {
	http.HandleFunc("/", root)
}


func AppError(wr WrapperRequest, w http.ResponseWriter, err error) {
	wr.C.Errorf("%v", err)
	er := &ErrorResponse{err.Error()}
	js,err:=json.Marshal(er)
	fmt.Fprintf(w, "%s", string(js[:len(js)]))
}


func AppWarning(wr WrapperRequest, msg string) {
	wr.C.Infof("%s", msg)
}



func AppLogout (w http.ResponseWriter, r *http.Request) {
	wr:=NewWrapperRequest(r)

	url, err := user.LogoutURL(wr.C, wr.R.URL.String())
	if err != nil {
		AppError(wr,w,err)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
}


func RedirectUserLogin(w http.ResponseWriter, r *http.Request){
	wr:=NewWrapperRequest(r)
	url, err := user.LoginURL(wr.C, wr.R.URL.String())
	if err != nil {
		AppError(wr,w,err)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
}








func root(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w,r,"/test/all",http.StatusMovedPermanently)
}




