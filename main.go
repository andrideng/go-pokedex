package main

import (
	"net/http"
	"os"
	"./helpers"
	"./models"
	"./controllers"
	"./logs"
	"github.com/asaskevich/govalidator"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	uuid := helpers.GetUuid(r)
	u := models.GetUserFromUuid(uuid)

	m := map[string]interface{}{
		"User":   u,
	}

	helpers.Render(w, "main", m)
}


func main() {
	// set env variable
	os.Setenv("port", ":8000")
	os.Setenv("dbname", "cache/pokedex_db.sqlite3")

	// handle static folder
	http.Handle("/statics/", http.StripPrefix("/statics/", http.FileServer(http.Dir("statics"))))
	http.Handle("/statics/pokedex", http.StripPrefix("/statics/pokedex/", http.FileServer(http.Dir("statics/pokedex"))))

	// set validator
	govalidator.SetFieldsRequiredByDefault(true)

	// routes
	http.HandleFunc("/", mainPage)

		// auth routes
	http.HandleFunc("/login", helpers.BeforeLogin(controllers.Login))
	http.HandleFunc("/register", helpers.BeforeLogin(controllers.Register))
	http.HandleFunc("/logout",  helpers.AfterLogin(controllers.Logout))

		// pokedex routes
	http.HandleFunc("/pokedex", controllers.Pokedex)
	http.HandleFunc("/pokedex/create", helpers.AfterLogin(controllers.Pokedex_create))
	http.HandleFunc("/pokedex/delete/", helpers.AfterLogin(controllers.Pokedex_destroy))
	http.HandleFunc("/pokedex/edit/", helpers.AfterLogin(controllers.Pokedex_edit))

	// listen and serve
	logs.Logger.Info("start serve at:%v", os.Getenv("port"))
	err := http.ListenAndServe(os.Getenv("port"), nil)
	logs.Logger.Critical("server err: %v", err)
}