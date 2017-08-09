package helpers

import (
	"../models"
	"net/http"
)

// set user session
func SetSession(u *models.User, w http.ResponseWriter) {

	value := map[string]string {
		"uuid" : u.Uuid,
	}

	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name: "session",
			Value: encoded,
			Path: "/",
		}
		http.SetCookie(w, cookie)
	}

}

// clear user session
func ClearSession(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name: name,
		Value: "",
		Path: "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

// get user uuid
func GetUuid(r *http.Request) (uuid string) {
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			uuid = cookieValue["uuid"]
		}
	}
	return uuid
}

// set cookie to store messages
func SetMsg(w http.ResponseWriter, name string, msg string) {
	value := map[string]string {
		name: msg,
	}
	if encoded, err := cookieHandler.Encode(name, value); err == nil {
		cookie := &http.Cookie{
			Name: name,
			Value: encoded,
			Path: "/",
		}
		http.SetCookie(w, cookie)
	}
}

// get cookie messages
func GetMsg(w http.ResponseWriter, r *http.Request, name string) (msg string) {
	if cookie, err := r.Cookie(name); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode(name, cookie.Value, &cookieValue); err == nil {
			msg = cookieValue[name]
			ClearSession(w, name)
		}
	}

	return msg
}
