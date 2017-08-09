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
