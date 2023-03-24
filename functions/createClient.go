package functions

import (
	"database/sql"
	"mailService/models/mail"
)

func CreateClient(db *sql.DB) {
	templates := mail.Templates{TemplateSend: "HOLA", TemplateReceive: "MUNDO"}
	auth := mail.Auth{
		Sender:   "logicielapplab@gmail.com",
		Alias:    "Logiciel Applab",
		Password: "cgmjfmgwjlsnqqku",
		Host:     "smtp.gmail.com",
		Port:     587,
	}
	newClient := mail.Client{
		Name:      "Logiciel Applab",
		Answer:    1,
		Templates: templates,
		Auth:      auth,
	}
	newClient.CreateClient(db)
}
