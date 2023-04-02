package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"mailService/models/mail"
	"mailService/models/request"
	"net/http"
	"os"
	"time"
)

var (
	db       *sql.DB
	err      error
	listBulk request.ListBulk
)

func init() {
	db, err = sql.Open("sqlite3", "./data/clients.db")
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
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
	router.POST("/api/bulkMail", BulkMail)
	router.POST("/api/createStandardClient", CreateClient)
	router.GET("/api/listClients", ListClients)
	router.POST("/api/deleteClient", DeleteClient)
	router.POST("/api/updateClient", UpdateClient)
	router.POST("/api/validateDataInput", ValidateDataInput)
	router.POST("/api/validateBulkTemplate", ValidateBulkTemplate)

	port := os.Getenv("PORT")
	if err = http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
func StandardMail(c *gin.Context) {
	var request request.RequestStandard
	request.ParseRequestStandardData(c)
	newClient := mail.Client{
		Name: request.ClientName,
	}
	newClient.GetClient(db)
	response := newClient.SendStandardMail(request)
	c.Writer.WriteString(response.Marshal())
}
func BulkMail(c *gin.Context) {
	var req request.RequestBulk
	var response request.RequestResponse
	req.ParseRequestBulkData(c)
	mailRequest, response := listBulk.GetRequestItemLimits(req)
	newClient := mail.Client{
		Name: req.ClientName,
	}
	newClient.GetClient(db)
	response = newClient.SendBulkMail(mailRequest)
	c.Writer.WriteString(response.Marshal())
}
func CreateClient(c *gin.Context) {
	var newClient mail.Client
	newClient.ParseClient(c)
	newClient.CreateClient(db)
}
func ListClients(c *gin.Context) {
	var clients mail.Clients
	clients.ListClients(db)
	fmt.Fprintln(c.Writer, clients.List)
}
func DeleteClient(c *gin.Context) {
	var client mail.Client
	client.ParseClient(c)
	client.DeleteClient(db)
}
func UpdateClient(c *gin.Context) {
	var client mail.Client
	client.ParseClient(c)
	client.UpdateClient(db)
}
func ValidateDataInput(c *gin.Context) {
	var request request.RequestBulk
	request.ClientName = c.Request.URL.Query()["clientName"][0]
	response := request.ValidateDataInput(c)
	if response.Error == "" {
		listBulk.List = append(listBulk.List, request)
	}
	c.Writer.WriteString(response.Marshal())
}
func ValidateBulkTemplate(c *gin.Context) {
	var request request.RequestTemplate
	request.ParseRequestBulkTemplate(c)
	response := request.ValidateTemplate()
	if response.Error == "" {
		response = listBulk.GetRequestItemTmp(request)
	}
	c.Writer.WriteString(response.Marshal())
}
