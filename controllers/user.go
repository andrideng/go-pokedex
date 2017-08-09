package controllers

import (
	"../helpers"
	"strings"
	"net/http"
	"../models"
	"github.com/asaskevich/govalidator"
)

// user controller
func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		u := &models.User{}
		u.Errors = make(map[string]string)
		u.Errors["username"] = helpers.GetMsg(w, r, "username")
		u.Errors["password"] = helpers.GetMsg(w, r, "password")
		u.Errors["message"] = helpers.GetMsg(w, r, "message")

		helpers.Render(w, "login", u)

	case "POST":
		username := r.FormValue("username")
		password := r.FormValue("password")
		u := &models.User{
			Username: username,
			Password: password,
		}
		if username == "" {
			helpers.SetMsg(w, "username", "username must be filled!")
		}
		if password == "" {
			helpers.SetMsg(w, "password", "password must be filled!")
		}
		if username != "" && password != "" {
			if auth, uuid := models.UserExists(u); auth {
				helpers.SetSession(&models.User{Uuid: uuid}, w)
				http.Redirect(w, r, "/", 302)
				return
			} else {
				helpers.SetMsg(w, "message", "Please register or enter valid username and password!")
				http.Redirect(w, r, "/login", 302)
				return
			}
		}

		http.Redirect(w, r, "/login", 302)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		u := &models.User{}
		u.Errors = make(map[string]string)
		u.Errors["username"] = helpers.GetMsg(w, r, "username")
		u.Errors["fullname"] = helpers.GetMsg(w, r, "fullname")
		u.Errors["password"] = helpers.GetMsg(w, r, "password")
		u.Errors["cpassword"] = helpers.GetMsg(w, r, "cpassword")

		helpers.Render(w, "register", u)

	case "POST":
		if n := models.CheckUser(r.FormValue("username")); n == true && r.FormValue("username") != "" {
			helpers.SetMsg(w, "username", "username already exists, please enter a unique username!")
			http.Redirect(w, r, "/register", 302)
			return
		}
		u := &models.User{
			Uuid: models.Uuid(),
			Fullname: r.FormValue("fullname"),
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		}

		result, err := govalidator.ValidateStruct(u)
		if err != nil {
			e := err.Error()
			if re := strings.Contains(e, "Username"); re == true {
				helpers.SetMsg(w, "username", "Please enter valid username!")
			}
			if re := strings.Contains(e, "Fullname"); re == true {
				helpers.SetMsg(w, "fullname", "Please enter valid fullname!")
			}
			if re := strings.Contains(e, "Password"); re == true {
				helpers.SetMsg(w, "password", "Please enter a password!")
			}
		}
		if r.FormValue("cpassword") == "" {
			helpers.SetMsg(w, "cpassword", "Please enter a confirm password!")
			http.Redirect(w, r, "/register", 302)
			return
		}
		if r.FormValue("cpassword") != r.FormValue("password") {
			helpers.SetMsg(w, "cpassword", "Confirm Password not match!")
			http.Redirect(w, r, "/register", 302)
			return
		}
		// save user to db
		if result == true {
			u.Password = models.EncryptPass(u.Password)
			models.Register(u)
			http.Redirect(w, r, "/login", 302)
			return
		}
		http.Redirect(w, r, "/register", 302)

	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	helpers.ClearSession(w, "session")
	http.Redirect(w, r, "/", 302)
}