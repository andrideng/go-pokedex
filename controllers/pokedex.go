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
	helpers.Render(w, "pokedex", p)
}


func Pokedex_create(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		p := &models.Pokedex{}
		p.Errors = make(map[string]string)
		p.Errors["pname"] = helpers.GetMsg(w, r, "pname")
		p.Errors["elements"] = helpers.GetMsg(w, r, "elements")

		helpers.Render(w, "pokedex_create", nil)

	case "POST":
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("pimages")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

func Pokedex_del(w http.ResponseWriter, r *http.Request) {
	uuid := strings.TrimPrefix(r.URL.Path, "/pokedex/delete")
	println(uuid)
}

