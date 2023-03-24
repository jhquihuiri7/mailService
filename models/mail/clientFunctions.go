package mail

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
)

func (c *Client) ParseClient(req *gin.Context) {
	err := json.NewDecoder(req.Request.Body).Decode(&c)
	if err != nil {
		log.Fatal(err)
	}
}

// create
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

	var wg sync.WaitGroup
	go c.Templates.AddTemplates(db, response, wg)
	go c.Auth.AddAuth(db, response, wg)
	wg.Wait()
}
func (t *Templates) AddTemplates(db *sql.DB, id int, wg sync.WaitGroup) {
	wg.Add(1)
	_, err := db.Exec(`
	INSERT INTO templates(id,templateReceive, templateSend) VALUES (?,?,?);
	`, id, t.TemplateReceive, t.TemplateSend)
	if err != nil {
		fmt.Println(err)
	}
	wg.Done()
}
func (a *Auth) AddAuth(db *sql.DB, id int, wg sync.WaitGroup) {
	wg.Add(1)
	_, err := db.Exec(`
	INSERT INTO auth(id,sender, alias, password, host, port) VALUES (?,?,?,?,?,?);
	`, id, a.Sender, a.Alias, a.Password, a.Host, a.Port)
	if err != nil {
		fmt.Println(err)
	}
	wg.Done()
}

// list
func (c *Client) GetClient(db *sql.DB) {
	//get ID
	var id int
	idRow := db.QueryRow("SELECT id, answer FROM clients WHERE name = ?", c.Name)
	err := idRow.Scan(&id, &c.Answer)
	if err != nil {
		log.Fatal(err)
	}

	//get DATA
	rowTemp := db.QueryRow("SELECT templateReceive, templateSend FROM templates WHERE id = ?", id)
	err = rowTemp.Scan(&c.TemplateReceive, &c.TemplateSend)
	if err != nil {
		log.Fatal(err)
	}

	rowAuth := db.QueryRow("SELECT * FROM auth WHERE id = ?", id)
	err = rowAuth.Scan(&id, &c.Sender, &c.Alias, &c.Password, &c.Host, &c.Port)
	if err != nil {
		log.Fatal(err)
	}
}
func (cs *Clients) ListClients(db *sql.DB) {
	rows, err := db.Query(`
		SELECT clients.id,clients.name, clients.answer,templates.templateReceive,templates.templateSend,auth.sender,auth.alias,auth.host,auth.port
		FROM clients
		INNER JOIN templates
		    ON clients.id = templates.id
		INNER JOIN auth
		    ON clients.id = auth.id
	`)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		client := Client{Auth: Auth{Password: "XXXXXXXX"}}
		rows.Scan(
			&client.Id, &client.Name, &client.Answer, &client.TemplateReceive,
			&client.TemplateSend, &client.Sender, &client.Alias, &client.Host, &client.Port)
		cs.List = append(cs.List, client)
	}
}

// delete
func (c *Client) DeleteClient(db *sql.DB) {
	_, err := db.Exec(
		`DELETE FROM clients WHERE id=?;
				DELETE FROM templates WHERE id=?;
				DELETE FROM auth WHERE id=?;`, c.Id, c.Id, c.Id)
	if err != nil {
		log.Fatal(err)
	}

}

func (c *Client) UpdateClient(db *sql.DB) {
	update := ""
	if c.TemplateReceive != "" {
		update += "templateReceive = " + c.TemplateReceive
	}
	if c.TemplateSend != "" {
		if update == "" {
			update += "templateSend = " + c.TemplateSend
		} else {
			update += ", templateSend = " + c.TemplateSend
		}
	}
	fmt.Println(update)
}
