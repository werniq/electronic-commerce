package utils

import (
	"new-e-commerce/models"
	"os"
)

type EmailData struct {
	URL       string
	Firstname string
	Subject   string
}

var (
	from     = os.Getenv("EmailFrom")
	smtpPass = os.Getenv("SMTPPass")
	smtpUser = os.Getenv("SMTPUser")
	smtpHost = os.Getenv("SMTPHost")
	smtpPort = os.Getenv("SMTPPort")
)

func SendEmail(user *models.User, data *EmailData, templateName string) error {
	to := user.Email

}
