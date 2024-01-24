package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"mailService/DB"
	"mailService/middleware"
	"mailService/models/mail"
	"mailService/models/request"
	"net/http"
)

var (
	err      error
	listBulk request.ListBulk
)

func init() {
	DB.InitSQLite()
}
func main() {
	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.POST("/api/standardMail", StandardMail)
	router.POST("/api/bulkMail", BulkMail)
	router.POST("/api/createStandardClient", CreateClient)
	router.GET("/api/listClients", ListClients)
	router.POST("/api/deleteClient", DeleteClient)
	router.POST("/api/updateClient", UpdateClient)
	router.POST("/api/validateDataInput", ValidateDataInput)
	router.POST("/api/validateBulkTemplate", ValidateBulkTemplate)
	//port := os.Getenv("PORT")
	if err = http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
	defer DB.SQliteDB.Close()
}

func StandardMail(c *gin.Context) {
	var request request.RequestStandard
	request.ParseRequestStandardData(c)
	newClient := mail.Client{
		Name: request.ClientName,
	}
	newClient.GetClient(DB.SQliteDB)
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
	newClient.GetClient(DB.SQliteDB)
	response = newClient.SendBulkMail(mailRequest)
	var tempList request.ListBulk
	for _, v := range listBulk.List {
		if v.ClientName != req.ClientName {
			tempList.List = append(tempList.List, v)
		}
	}
	listBulk = tempList
	c.Writer.WriteString(response.Marshal())
}
func CreateClient(c *gin.Context) {
	var newClient mail.Client
	newClient.ParseClient(c)
	newClient.CreateClient(DB.SQliteDB)
}
func ListClients(c *gin.Context) {
	var clients mail.Clients
	clients.ListClients(DB.SQliteDB)
	fmt.Fprintln(c.Writer, clients.List)
}
func DeleteClient(c *gin.Context) {
	var client mail.Client
	client.ParseClient(c)
	client.DeleteClient(DB.SQliteDB)
}
func UpdateClient(c *gin.Context) {
	var client mail.Client
	client.ParseClient(c)
	client.UpdateClient(DB.SQliteDB)
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
