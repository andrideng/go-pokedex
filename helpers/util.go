package helpers

import (
	"net/http"
	"html/template"
	"github.com/gorilla/securecookie"
)

var (
	cookieHandler = securecookie.New(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32),
	)
)

func Render(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := template.ParseGlob("views/*.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	tmpl.ExecuteTemplate(w, name, data)
}

func BeforeLogin(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r* http.Request) {
		uuid := GetUuid(r)
		if uuid == "" {
			fn(w, r)
			return
		}
		http.Redirect(w, r, "/", 302)
	}
}


func AfterLogin(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r* http.Request) {
		uuid := GetUuid(r)
		if uuid != "" {
			fn(w, r)
			return
		}
		http.Redirect(w, r, "/login", 302)
	}
}