package main

import (
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"mailService/functions"
	"mailService/models"
	"net/http"
	"os"
	"time"
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
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Access-Control-Allow-Origin", "Content-Length", "Content-type"},
		ExposeHeaders:    []string{"Content-Length", "Content-type"},
		AllowCredentials: true,

		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))
	router.POST("/api/standardMail", StandardMail)
	router.POST("/api/createStandardClient", CreateClient)
	port := os.Getenv("PORT")
	if err = http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
func StandardMail(c *gin.Context) {
	var request models.RequestMessage
	request.ParseRequestData(c)
	response := functions.SendMail(db, request)
	c.Writer.WriteString(response.Marshal())
}
func CreateClient(c *gin.Context) {
	var newClient models.Client
	newClient.ParseClient(c)
	newClient.CreateClient(db)
}
