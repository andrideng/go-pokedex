package main

import (
	"net/http"
	"os"
	"./helpers"
	"./models"
	"./controllers"
	"github.com/asaskevich/govalidator"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	uuid := helpers.GetUuid(r)
	u := models.GetUserFromUuid(uuid)
	helpers.Render(w, "main", u)
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
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/logout", controllers.Logout)
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/pokedex", controllers.Pokedex)
	http.HandleFunc("/pokedex/create", controllers.Pokedex_create)

	// listen and serve
	http.ListenAndServe(os.Getenv("port"), nil)
}