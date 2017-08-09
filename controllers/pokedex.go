package controllers

import (
	"net/http"
	"../models"
	"github.com/asaskevich/govalidator"
	"strings"
	"../helpers"
)

func Pokedex(w http.ResponseWriter, r *http.Request) {
	helpers.Render(w, "pokedex", nil)
}
func Pokedex_create(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		//p := &models.Pokedex{}
		//p.Errors = make(map[string]string)
		//p.Errors["pokemon_name"] = getMsg(w, r, "pokemon_name")
		//p.Errors["pokemon_images"] = getMsg(w, r, "pokemon_images")
		//p.Errors["elements"] = getMsg(w, r, "elements")

		helpers.Render(w, "pokedex_create", nil)

	case "POST":
		p := &models.Pokedex{
			Uuid: models.Uuid(),
			Name: r.FormValue("pokemon_name"),
			Images: r.FormValue("pokemon_images"),
			Elements: r.FormValue("elements"),
		}
		result, err := govalidator.ValidateStruct(p)
		if err != nil {
			e := err.Error()
			if re := strings.Contains(e, "Name"); re {
				helpers.SetMsg(w, "pokemon_name", "please enter pokemon name!")
			}
			if re := strings.Contains(e, "Images"); re {
				helpers.SetMsg(w, "pokemon_images", "please enter pokemon images!")
			}
			if re := strings.Contains(e, "Elements"); re {
				helpers.SetMsg(w, "elements", "please enter pokemon elements!")
			}
		}
		// ready to insert
		if result {
			models.InsertPokedex(p)
			http.Redirect(w, r, "/pokedex", 302)
			return
		}
		http.Redirect(w, r, "/pokedex/create", 302)
	}

}

