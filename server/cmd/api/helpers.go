package main

import (
	"bytes"
	"crypto/rand"
	"embed"
	"fmt"
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"os"
	"time"
)

//go:embed templates
var emailTemplateFS embed.FS
var (
	EmailFrom = os.Getenv("EmailFrom")
	SMTPPass  = os.Getenv("SMTPPass")
	SMTPUser  = os.Getenv("SMTPUser")
	SMTPHost  = os.Getenv("SMTPHost")
	SMTPPort  = os.Getenv("SMTPPort")
)

func (app *application) GenerateUserId() (string, error) {
	p, err := rand.Prime(rand.Reader, 64)
	if err != nil {
		return "", err
	}
	return p.String(), nil
}

func (app *application) SendEmail(from, to, subject, tmpl string, data interface{}) error {
	templateToRender := fmt.Sprintf("templates/%s.html.tmpl", tmpl)

	t, err := template.New("email.html").ParseFS(emailTemplateFS, templateToRender)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "body", data); err != nil {
		app.errorLog.Println(err)
		return err
	}

	formattedMsg := tpl.String()

	templateToRender = fmt.Sprintf("templates/%s.plain.tmpl", tmpl)
	t, err = template.New("email-plain").ParseFS(emailTemplateFS, templateToRender)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	err = t.ExecuteTemplate(&tpl, "body", data)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	plainMsg := tpl.String()
	fmt.Println(formattedMsg, plainMsg)

	// send the email
	server := mail.NewSMTPClient()
	server.Host = app.cfg.smtp.host
	server.Port = app.cfg.smtp.port
	server.Username = app.cfg.smtp.username
	server.Password = app.cfg.smtp.password
	server.Encryption = mail.EncryptionTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	email := mail.NewMSG()
	email.
		SetFrom(from).
		AddTo(to).
		SetSubject(subject)

	email.SetBody(mail.TextHTML, formattedMsg)
	email.AddAlternative(mail.TextPlain, plainMsg)

	err = email.Send(smtpClient)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}
	fmt.Println("message send")
	return nil
}
