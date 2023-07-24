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

	response = c.SendMessage(c.Sender, "Nuevo mensaje de cliente", "")
	if c.Answer == 1 {
		c.SendMessage(req.Mail, "Hemos recibido tu correo", "")
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
		for _, toName := range []string{"Correo", "Mail", "Email"} {
			_, ok := val[toName]
			if ok {
				toColumn = toName
				break
			}
		}
		go func(client Client, req map[string]string, toColumn string, failedMail chan string) {
			fmt.Println(req)
			client.ParseTemplate(req)
			mailBack := mailStore.MailStore{Mail: req[toColumn]}
			mailBack.AddMail(req)
			fmt.Println("1." + req[toColumn])
			fmt.Println(client.TemplateReceive)
			response = client.SendMessage(req[toColumn], "HOLA", "")
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
	temp, err := os.ReadFile("templates/reportMailTemplate.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	c.TemplateReceive = string(temp)
	c.SendMessage(c.Sender, "Reporte de correos masivos enviados", nf.Name())
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

func (c *Client) SendMessage(tos, subject, attachFile string) request.RequestResponse {
	msg := gomail.NewMessage()
	msg.SetHeader("From", fmt.Sprintf("%s <%s>", c.Alias, c.Sender))
	msg.SetHeader("To", tos)
	msg.SetHeader("Subject", subject)
	if tos == c.Sender {
		msg.SetBody("text/html", c.TemplateReceive)
	} else {
		msg.SetBody("text/html", c.TemplateSend)
	}

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
