package models

import (
	"bytes"
	"fmt"

	//"fmt"
	//"gopkg.in/gomail.v2"
	"html/template"
	"log"
)

func (c *Client) SendMail() {
	c.ParseTemplate(false)
	fmt.Println(c.TemplateReceive)
	//msg := gomail.NewMessage()
	//msg.SetHeader("From", fmt.Sprintf("%s <%s>", c.Alias, c.Sender))
	//msg.SetHeader("To", "<jhonatan.quihuiri@gmail.com>")
	//msg.SetHeader("Subject", "Nuevo Mensaje de Cliente")
	//msg.SetBody("text/html", c.TemplateReceive)
	//n := gomail.NewDialer(c.Host, c.Port, c.Sender, c.Password)
	//
	//// Send the email
	//if err := n.DialAndSend(msg); err != nil {
	//panic(err)
	//}
}

func (c *Client) ParseTemplate(answer bool) {
	var (
		tmp *template.Template
		err error
	)
	switch answer {
	case false:
		tmp, err = template.New("Template").Parse(c.TemplateReceive)
	case true:
		tmp, err = template.New("Template").Parse(c.TemplateSend)
	}
	if err != nil {
		log.Fatal(err)
	}
	temp := new(bytes.Buffer)
	tmp.ExecuteTemplate(temp, "Template", "LOGICIELAPPLAB")

	switch answer {
	case false:
		c.TemplateReceive = temp.String()
	case true:
		c.TemplateSend = temp.String()
	}
}
