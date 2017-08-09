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
