package models

import (
	"encoding/json"
	"io"
	"log"
)

type RequestMessage struct {
	ClientName string `json:"clientName"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Mail       string `json:"mail"`
	Message    string `json:"message"`
}

func (r *RequestMessage) ParseRequestData(body io.ReadCloser) {
	err := json.NewDecoder(body).Decode(&r)
	if err != nil {
		log.Fatal(err)
	}
}
