package models

import (
	"bytes"
	"fmt"
	"gopkg.in/gomail.v2"

	//"fmt"
	//"gopkg.in/gomail.v2"
	"html/template"
	"log"
)

func (c *Client) SendMail() RequestResponse {
	c.ParseTemplate()
	return c.SendMessage()
}

func (c *Client) ParseTemplate() {
	var (
		tmpReceive, tmpSend   *template.Template
		err                   error
		tempReceive, tempSend bytes.Buffer
	)
	switch c.Answer {
	case 1:
		tmpReceive, err = template.New("TemplateReceive").Parse(c.TemplateReceive)
		tmpSend, err = template.New("TemplateSend").Parse(c.TemplateSend)

		tmpReceive.ExecuteTemplate(&tempReceive, "TemplateReceive", "LOGICIELAPPLAB")
		tmpSend.ExecuteTemplate(&tempSend, "TemplateSend", "LOGICIELAPPLAB")

		c.TemplateReceive = tempReceive.String()
		c.TemplateSend = tempSend.String()
	case 0:
		tmpReceive, err = template.New("TemplateReceive").Parse(c.TemplateReceive)

		tmpReceive.ExecuteTemplate(&tempReceive, "TemplateReceive", "LOGICIELAPPLAB")

		c.TemplateReceive = tempReceive.String()
	}
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Client) SendMessage() RequestResponse {
	msg := gomail.NewMessage()
	msg.SetHeader("From", fmt.Sprintf("%s <%s>", c.Alias, c.Sender))
	msg.SetHeader("To", "<logicielapplab@gmail.com>")
	msg.SetHeader("Subject", "Nuevo Mensaje de Cliente")
	msg.SetBody("text/html", c.TemplateReceive)
	n := gomail.NewDialer(c.Host, c.Port, c.Sender, c.Password)

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		panic(err)
	} else {
		if c.Answer == 1 {
			msg.SetHeader("To", "<jhonatan.quihuiri@gmail.com>")
			msg.SetBody("text/html", c.TemplateSend)
			n.DialAndSend(msg)
		}
	}
	return RequestResponse{Success: "Correo enviado correctamente"}
}
