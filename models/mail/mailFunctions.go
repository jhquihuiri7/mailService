package mail

import (
	"bytes"
	"fmt"
	"gopkg.in/gomail.v2"
	"mailService/models/request"
	"os"
	"sync"
	//"fmt"
	//"gopkg.in/gomail.v2"
	"html/template"
	"log"
)

func (c *Client) SendStandardMail(req request.RequestStandard) request.RequestResponse {
	var response request.RequestResponse
	c.ParseTemplate(req)

	response = c.SendMessage(c.Sender, "")
	if c.Answer == 1 {
		c.SendMessage(req.Mail, "")
	}
	return response
}

func (c *Client) SendBulkMail(req request.RequestBulk) request.RequestResponse {
	var response request.RequestResponse
	nf, err := os.Create("./logs/bulk.txt")
	defer nf.Close()
	if err != nil {
		log.Fatal(err)
	}
	nf.WriteString(fmt.Sprintf("%s %s\n%s %d al %d %s\n\n%s\n",
		"Reporte de correos masivos enviador por", c.Name,
		"------------ Correos entre los números", req.Limits[0], req.Limits[1], "------------",
		"Correos enviados fallidos:"))
	var wg sync.WaitGroup
	c.TemplateReceive = req.Template
	wg.Add(len(req.Tos[req.Limits[0]-1 : req.Limits[1]]))
	for _, v := range req.Tos[req.Limits[0]-1 : req.Limits[1]] {
		go func(req request.RequestStandard) {
			c.ParseTemplate(v)
			response = c.SendMessage(req.Mail, "")
			if response.Error != "" {
				nf.WriteString(fmt.Sprintf("- %s\n", req.Mail))
			}
			fmt.Println(req)
			wg.Done()
		}(v)
	}
	wg.Wait()
	fmt.Println(nf.Name())
	c.SendMessage(c.Sender, nf.Name())
	return response
}

func (c *Client) ParseTemplate(data interface{}) {
	var (
		tmpReceive, tmpSend   *template.Template
		err                   error
		tempReceive, tempSend bytes.Buffer
	)
	switch c.Answer {
	case 1:
		tmpReceive, err = template.New("TemplateReceive").Parse(c.TemplateReceive)
		tmpSend, err = template.New("TemplateSend").Parse(c.TemplateSend)

		tmpReceive.ExecuteTemplate(&tempReceive, "TemplateReceive", data)
		tmpSend.ExecuteTemplate(&tempSend, "TemplateSend", data)

		c.TemplateReceive = tempReceive.String()
		c.TemplateSend = tempSend.String()
	case 0:
		tmpReceive, err = template.New("TemplateReceive").Parse(c.TemplateReceive)

		err = tmpReceive.ExecuteTemplate(&tempReceive, "TemplateReceive", data)

		c.TemplateReceive = tempReceive.String()
	}
	fmt.Println(tempReceive.String())
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Client) SendMessage(tos string, attachFile string) request.RequestResponse {
	msg := gomail.NewMessage()
	msg.SetHeader("From", fmt.Sprintf("%s <%s>", c.Alias, c.Sender))
	msg.SetHeader("To", tos)
	msg.SetHeader("Subject", "Nuevo Mensaje de Cliente")
	msg.SetBody("text/html", c.TemplateReceive)
	if attachFile != "" {
		msg.Attach(attachFile)
	}
	n := gomail.NewDialer(c.Host, c.Port, c.Sender, c.Password)

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		return request.RequestResponse{Error: "Error al enviar correo"}
	} else {
		return request.RequestResponse{Success: "Correo enviado correctamente"}
	}
}
