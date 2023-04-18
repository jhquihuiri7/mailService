package mail

import (
	"bytes"
	"fmt"
	"gopkg.in/gomail.v2"
	"html/template"
	"log"
	"mailService/models/mailStore"
	"mailService/models/pdfReport"
	"mailService/models/request"
	"os"
	"sync"
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
	report := pdfReport.PdfReport{
		ClientName: req.ClientName,
		Limits:     req.Limits,
	}
	failedMail := make(chan string, 500)

	var wg sync.WaitGroup
	c.TemplateReceive = req.Template
	wg.Add(len(req.Tos[req.Limits[0]-1 : req.Limits[1]]))
	for _, val := range req.Tos[req.Limits[0]-1 : req.Limits[1]] {
		var toColumn string
		for _, toName := range []string{"correo", "mail", "email"} {
			_, ok := val[toName]
			if ok {
				toColumn = toName
				break
			}
		}
		go func(client Client, req map[string]string, toColumn string, failedMail chan string) {
			client.ParseTemplate(req)
			mailBack := mailStore.MailStore{Mail: req[toColumn]}
			mailBack.AddMail(req)
			response = client.SendMessage(req[toColumn], "")
			if response.Error != "" {
				failedMail <- req[toColumn]
			}
			wg.Done()
		}(*c, val, toColumn, failedMail)
	}
	wg.Wait()
	close(failedMail)
	for chanVal := range failedMail {
		report.ErrorMail = append(report.ErrorMail, chanVal)
		report.ErrorCount++
	}
	nf := report.GenerateBulkReport()
	defer os.Remove(nf.Name())
	defer nf.Close()
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
