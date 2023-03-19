package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func (c *Client) ParseClient(req *gin.Context) {
	err := json.NewDecoder(req.Request.Body).Decode(&c)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c)
}
func (c *Client) CreateClient(db *sql.DB) {
	_, err := db.Exec(`
	INSERT INTO clients(name,answer) VALUES (?,?);
	`, c.Name, c.Answer)
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := db.Prepare("SELECT id FROM clients WHERE name = ?")
	if err != nil {
		log.Fatal(err)
	}
	res := stmt.QueryRow(c.Name)
	var response int
	res.Scan(&response)
	c.Templates.AddTemplates(db, response)
	c.Auth.AddAuth(db, response)
}
func (t *Templates) AddTemplates(db *sql.DB, id int) {
	_, err := db.Exec(`
	INSERT INTO templates(id,templateReceive, templateSend) VALUES (?,?,?);
	`, id, t.TemplateReceive, t.TemplateSend)
	if err != nil {
		log.Fatal(err)
	}
}
func (a *Auth) AddAuth(db *sql.DB, id int) {
	_, err := db.Exec(`
	INSERT INTO auth(id,sender, alias, password, host, port) VALUES (?,?,?,?,?,?);
	`, id, a.Sender, a.Alias, a.Password, a.Host, a.Port)
	if err != nil {
		log.Fatal(err)
	}
}
