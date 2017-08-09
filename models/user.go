package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Uuid 		string `valid:"required,uuidv4"`
	Username 	string `valid:"required,alphanum"`
	Password 	string `valid:"required"`
	Fullname 	string `valid:"required"`
	Errors		map[string]string `valid:"-"`
}

// register user
func Register(u *User) error{
	var db, _ = sql.Open("sqlite3", os.Getenv("dbname"))
	defer db.Close()
	db.Exec("create table if not exists users (uuid text not null unique, username text not null unique, password text not null, fullname text not null, primary key (uuid))")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into users(uuid, username, password, fullname) values(?, ?, ?, ?)")
	_, err := stmt.Exec(u.Uuid, u.Username, u.Password, u.Fullname)
	tx.Commit()
	return err
}

// check exists username for register
func CheckUser(username string) bool {
	var db, _ = sql.Open("sqlite3", os.Getenv("dbname"))
	defer db.Close()
	var un string
	q, err := db.Query("select username from users where username='"+username+"'")
	if err != nil {
		return false
	}
	for q.Next() {
		q.Scan(&un)
	}
	if un == username {
		return true
	}
	return false
}

// check username for login
func UserExists(u *User) (bool, string) {
	var db, _ = sql.Open("sqlite3", os.Getenv("dbname"))
	defer db.Close()
	var pass, uuid string
	q, err := db.Query("select uuid, password from users where username = '"+u.Username+"'")
	if err != nil {
		return false, ""
	}
	for q.Next() {
		q.Scan(&uuid, &pass)
	}
	pw := bcrypt.CompareHashAndPassword([]byte(pass), []byte(u.Password))
	if uuid != "" && pw == nil {
		return true, uuid
	}
	return false, ""
}

func GetUserFromUuid(uuid string) *User {
	var db, _ = sql.Open("sqlite3", os.Getenv("dbname"))
	defer db.Close()
	var uu, username, fullname string
	q, err := db.Query("select uuid, username, fullname from users where uuid = '"+uuid+"'")
	if err != nil {
		return &User{}
	}
	if q.Next() {
		q.Scan(&uu, &username, &fullname)
	}
	return &User{Uuid:uu, Fullname:fullname, Username:username}
}

// utilities for users model
func EncryptPass(password string) string {
	pass := []byte(password)
	hashpw, _ := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	return string(hashpw)
}

func Uuid() (string) {
	id := uuid.NewV4().String()
	return id
}

