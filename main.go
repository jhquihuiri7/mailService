package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"mailService/functions"
	"mailService/models"
)

func main() {
	db, err := sql.Open("sqlite3", "./data/clients.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//_, err = db.Exec("CREATE TABLE auth(id integer not null primary key, sender text, alias text, password text, host text, port integer)")
	//if err != nil {
	//	log.Fatal(err)
	//}

	//functions.CreateClient(db)
	request := models.RequestMessage{Name: "Logiciel Applab"}
	functions.SendMail(db, request)
}
