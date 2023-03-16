package models

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
)

type RequestMessage struct {
	ClientName string `json:"clientName"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Mail       string `json:"mail"`
	Message    string `json:"message"`
}
type RequestResponse struct {
	Success string `json:"success"`
	Error   string `json:"error"`
}

func (r *RequestMessage) ParseRequestData(c *gin.Context) {
	err := json.NewDecoder(c.Request.Body).Decode(&r)
	if err != nil {
		log.Fatal(err)
	}
}

func (resp *RequestResponse) Marshal() string {
	JSONresponse, _ := json.Marshal(resp)
	return string(JSONresponse)
}
