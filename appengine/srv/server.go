package srv

import (
	"net/http"

	//"encoding/json"
	"html/template"
	"fmt"
	"errors"

	"appengine"
	"appengine/user"
	//"appengine/users"
	"users"
	"appengine/datastore"
)




type ErrorResponse struct{
	Error string
}



type WrapperRequest struct{
	R *http.Request
	C appengine.Context
	U *user.User
	NU *users.NUser
}


type WrapperHandler func(wr WrapperRequest, tc map[string]interface{}) (string, error)


func NewWrapperRequest(r *http.Request) (WrapperRequest) {
	c:=appengine.NewContext(r)
	return WrapperRequest{r, c, user.Current(c),nil}
}

func (wr WrapperRequest) IsAdminRequest()(bool){
	return user.IsAdmin(wr.C)
}




func AppHandler(w http.ResponseWriter, r *http.Request, whandler WrapperHandler) {
	wr:=NewWrapperRequest(r)
	err:=GetCurrentUser(&wr)
	if err!=nil{
		RedirectUserLogin(w,wr.R)
		return
	}

	rdata := make(map[string]interface{})
	rdata["User"] = wr.NU
	tmplName,err:=whandler(wr, rdata)

	if err!=nil{
		AppError(wr,w,err)
		return
	}

	tmpl:= template.Must(template.ParseFiles("app/tmpl/base.html",
		tmplName))

	if err := tmpl.Execute(w, rdata); err != nil {
		AppError(wr,w,err)
		return
	}

}



func AppError(wr WrapperRequest, w http.ResponseWriter, err error) {
	wr.C.Errorf("%v", err)

	errorTmpl := template.Must(template.ParseFiles("app/tmpl/error.html"))
	if err := errorTmpl.Execute(w, err.Error()); err != nil {
		wr.C.Errorf("%v", err)
		return
	}
}


func AppWarning(wr WrapperRequest, msg string) {
	wr.C.Infof("%s", msg)
}



func AppLogout (w http.ResponseWriter, r *http.Request) {
	wr:=NewWrapperRequest(r)
	RedirectUserLogin(w,wr.R)
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




func GetCurrentUser(wr *WrapperRequest)(error){

	// Busco información del usuario de sesión

	var nu users.NUser
	u := wr.U
	if u==nil{
		return errors.New("No user session founded")
	}

	q := datastore.NewQuery("users").Filter("Mail =",u.Email)
	var nusers []users.NUser
	keys,_:= q.GetAll(wr.C, &nusers)
	
	if (len(nusers)<=0){
		// El usuario de sesión no esta en el datastore
		// Usamos el admin de la app aunque no este en el datastore
		if wr.IsAdminRequest(){
			nu = users.NUser{-1,u.Email,"Administrador",users.ROLE_ADMIN}
		}else{
			return errors.New("No user id found")
		}
	}else{
		nu=nusers[0]
		nu.Id=keys[0].IntID()
	}

	wr.NU = &nu
	return nil
}


func CheckPerm(wr WrapperRequest, op byte)(error) {

	n:= wr.NU
	if (!n.IsAdmin()){
		// Si no es admin, deberiamos buscarlo en nuestra base
		// de datos de usuarios permitidos y comprobar si 
		// con su rol puede hacer dicha operación
		// De esa busqueda calculamos la variable perm y la comparamos
		// con op

		if !users.IsAllowed(n.Role,op){
			AppWarning(wr,"Perm:"+fmt.Sprintf("%b",n.Role)+" "+fmt.Sprintf("%b",op))
	                AppWarning(wr,"User "+n.Mail+" failed allowed access")
			return errors.New("Operation not allowed")
		}
	}

	// Si es admin puede cualquier cosa
	return nil
}
