package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type Pokedex struct {
	Uuid 		string	`valid:"required,uuidv4"`
	Name 		string	`valid:"required"`
	Images 		string	`valid:"required"`
	Elements	string	`valid:"required"`
	Errors 		map[string]string	`valid:"-"`
}

type PokedexList []Pokedex

func InsertPokedex(p *Pokedex) error {
	var db, _ = sql.Open("sqlite3", os.Getenv("dbname"))
	defer db.Close()
	db.Exec("create table if not exists pokedex (uuid text not null unique, name text not null, images text not null, elements text not null, primary key (uuid))")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into pokedex(uuid, name, images, elements) values(?,?,?,?)")
	_, err := stmt.Exec(p.Uuid, p.Name, p.Images, p.Elements)
	tx.Commit()
	return err
}

func GetPokedex() *PokedexList{
	var db, _ = sql.Open("sqlite3", os.Getenv("dbname"))
	defer db.Close()
	q, err := db.Query("select uuid, name, images, elements from pokedex")
	if err != nil {
		return &PokedexList{}
	}

	list := PokedexList{}

	for q.Next() {
		var uuid, name, images, elements string
		q.Scan(&uuid, &name, &images, &elements)
		p := Pokedex{Uuid:uuid, Name: name, Images: images, Elements: elements}

		list = append(list, p)
	}
	return &list
}

func UpdatePokedex(p *Pokedex) (error, bool) {
	var db, _ = sql.Open("sqlite3", os.Getenv("dbname"))
	defer db.Close()
	_, err := db.Query("update pokedex set name='"+p.Name+"', images='"+p.Images+"', elements='"+p.Elements+"' where uuid='"+p.Uuid+"'")
	return err, true
}

func DeletePokedex(uuid string) bool{
	var db, _ = sql.Open("sqlite3", os.Getenv("dbname"))
	defer db.Close()
	_, err := db.Query("delete from pokedex where uuid = '"+uuid+"'")
	if err != nil {
		return false
	}
	return true
}