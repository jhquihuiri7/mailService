package functions

import (
	"database/sql"
	"log"
	"mailService/models"
)

func SendMail(db *sql.DB, message models.RequestMessage) models.RequestResponse {
	//createClient
	newClient := models.Client{
		Name: message.ClientName,
	}

	//get ID
	var id int
	idRow := db.QueryRow("SELECT id, answer FROM clients WHERE name = ?", newClient.Name)
	err := idRow.Scan(&id, &newClient.Answer)
	if err != nil {
		log.Fatal(err)
	}

	//get DATA
	rowTemp := db.QueryRow("SELECT templateReceive, templateSend FROM templates WHERE id = ?", id)
	err = rowTemp.Scan(&newClient.TemplateReceive, &newClient.TemplateSend)
	if err != nil {
		log.Fatal(err)
	}

	rowAuth := db.QueryRow("SELECT * FROM auth WHERE id = ?", id)
	err = rowAuth.Scan(&id, &newClient.Sender, &newClient.Alias, &newClient.Password, &newClient.Host, &newClient.Port)
	if err != nil {
		log.Fatal(err)
	}
	return newClient.SendMail()
}
