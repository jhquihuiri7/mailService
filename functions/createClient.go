package functions

import (
	"database/sql"
	"mailService/models"
)

func CreateClient(db *sql.DB) {
	templates := models.Templates{TemplateSend: "HOLA", TemplateReceive: "MUNDO"}
	auth := models.Auth{
		Sender:   "logicielapplab@gmail.com",
		Alias:    "Logiciel Applab",
		Password: "cgmjfmgwjlsnqqku",
		Host:     "smtp.gmail.com",
		Port:     587,
	}
	newClient := models.Client{
		Name:      "Logiciel Applab",
		Answer:    1,
		Templates: templates,
		Auth:      auth,
	}
	newClient.CreateClient(db)
}
