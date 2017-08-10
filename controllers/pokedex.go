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

func images(w http.ResponseWriter, r *http.Request, name string) string {
	var filename string
	err := r.ParseMultipartForm(100000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	m := r.MultipartForm
	files := m.File[name]
	for i := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}
		f, err := os.Create("./statics/pokedex/"+files[i].Filename)
		defer f.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}
		if _, err := io.Copy(f, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}
		filename = files[i].Filename
	}

	return filename
}

func Pokedex(w http.ResponseWriter, r *http.Request) {
	helpers.Render(w, "pokedex", nil)
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
		p := &models.Pokedex{
			Uuid: models.Uuid(),
			Name: r.FormValue("pname"),
			Images : "hello", //images(w, r, "pimages")
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
			models.InsertPokedex(p)
			http.Redirect(w, r, "/pokedex", 302)
			return
		}
		http.Redirect(w, r, "/pokedex/create", 302)
	}

}

