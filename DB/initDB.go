package DB

import (
	"database/sql"
	"log"
)

var (
	SQliteDB *sql.DB
	err      error
)

func InitSQLite() {
	SQliteDB, err = sql.Open("sqlite3", "./data/clients.db")
	if err != nil {
		log.Fatal(err)
	}
}
