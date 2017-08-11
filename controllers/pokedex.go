package controllers

import (
	"net/http"
	"../models"
	"github.com/asaskevich/govalidator"
	"strings"
	"../helpers"
	"os"
	"io"
)


func Pokedex(w http.ResponseWriter, r *http.Request) {
	p := models.GetPokedex()

	uuid := helpers.GetUuid(r)
	u := models.GetUserFromUuid(uuid)

	m := map[string]interface{}{
		"Pokedex": p,
		"User":   u,
	}

	helpers.Render(w, "pokedex", m)
}


func Pokedex_create(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		p := &models.Pokedex{}
		p.Errors = make(map[string]string)
		p.Errors["pname"] = helpers.GetMsg(w, r, "pname")
		p.Errors["pimages"] = helpers.GetMsg(w, r, "pimages")
		p.Errors["elements"] = helpers.GetMsg(w, r, "elements")

		uuid := helpers.GetUuid(r)
		us := models.GetUserFromUuid(uuid)

		m := map[string]interface{}{
			"User": us,
			"Info": p,
		}

		helpers.Render(w, "pokedex_create", m)

	case "POST":
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("pimages")
		if err != nil {
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			helpers.SetMsg(w, "pimages", "please enter pokemon images!")
			http.Redirect(w, r, "/pokedex/create", 302)
			return
		}
		defer file.Close()

		p := &models.Pokedex{
			Uuid: models.Uuid(),
			Name: r.FormValue("pname"),
			Images : handler.Filename,
			Elements: r.FormValue("elements"),
		}
		result, err := govalidator.ValidateStruct(p)
		if err != nil {
			e := err.Error()
			if re := strings.Contains(e, "Name"); re {
				helpers.SetMsg(w, "pname", "please enter pokemon name!")
			}
			if re := strings.Contains(e, "Elements"); re {
				helpers.SetMsg(w, "elements", "please enter pokemon elements!")
			}
		}
		// ready to insert
		if result {
			// insert to folder
			f, err := os.OpenFile("./statics/pokedex/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			defer f.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer f.Close()
			io.Copy(f, file)

			models.InsertPokedex(p)
			http.Redirect(w, r, "/pokedex", 302)
			return
		}
		http.Redirect(w, r, "/pokedex/create", 302)
	}

}

func Pokedex_destroy(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Path[len("/pokedex/delete/"):]
	models.DeletePokedex(uuid)
	http.Redirect(w, r, "/pokedex", 302)
}

func Pokedex_edit(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Path[len("/pokedex/edit/"):]

	switch r.Method {
	case "GET":
		p := models.GetOnePokedex(uuid)
		uuid := helpers.GetUuid(r)
		us := models.GetUserFromUuid(uuid)

		m := map[string]interface{}{
			"User": us,
			"Info": p,
		}

		helpers.Render(w, "pokedex_edit", m)

	case "POST":
		p := &models.Pokedex{
			Uuid: uuid,
			Name: r.FormValue("pname"),
			Elements: r.FormValue("elements"),
		}
		models.UpdatePokedex(p)
		http.Redirect(w, r, "/pokedex", 302)
		return
	}

}

