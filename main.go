package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"mailService/functions"
	"mailService/models"
	"net/http"
	"os"
)

var (
	db  *sql.DB
	err error
)

func init() {
	db, err = sql.Open("sqlite3", "./data/clients.db")
	if err != nil {
		log.Fatal(err)
	}
}
func main() {

	//_, err = db.Exec("CREATE TABLE auth(id integer not null primary key, sender text, alias text, password text, host text, port integer)")
	//if err != nil {
	//	log.Fatal(err)
	//}

	//functions.CreateClient(db)
	//request := models.RequestMessage{Name: "Logiciel Applab"}
	//functions.SendMail(db, request)
	router := gin.New()
	router.POST("/api/standarMail", StandarMail)
	port := os.Getenv("PORT")
	if err = http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
func StandarMail(c *gin.Context) {
	var request models.RequestMessage
	request.ParseRequestData(c)
	response := functions.SendMail(db, request)
	c.Writer.WriteString(response.Marshal())

}
