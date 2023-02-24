package internal

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

func Send(path string) error {
	e := email.NewEmail()
	e.From = "Tony Mannino <anthonyjosephmannino@gmail.com>"
	e.To = []string{os.Getenv("RECIPIENT")}
	e.Subject = "sage test email"
	e.Text = []byte("Text Body is, of course, supported!")
	e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	e.AttachFile(path)
	return e.Send(fmt.Sprintf("%s:%s", os.Getenv("EMAIL_HOST"), os.Getenv("EMAIL_PORT")), smtp.PlainAuth("", os.Getenv("EMAIL_ADDRESS"), os.Getenv("EMAIL_PASSWORD"), os.Getenv("EMAIL_HOST")))
}

// SendMail builds and sends a plaintext email to a recipient
func SendMail(to, subject, body string) error {
	email := os.Getenv("EMAIL_ADDRESS")
	if email == "" {
		return errors.New("no smtp email in env")
	}
	password := os.Getenv("EMAIL_PASSWORD")
	if password == "" {
		return errors.New("no smtp password in env")
	}
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		// default gmail for myself
		host = "smtp.gmail.com"
	}
	port := os.Getenv("SMTP_PORT")
	if port == "" {
		// default gmail port form myself
		port = "587"
	}

	message := fmt.Sprintf("Subject: %s\n%s", subject, body)
	auth := smtp.PlainAuth("", email, password, host)
	err := smtp.SendMail(fmt.Sprintf("%s:%s", host, port), auth, email, []string{to}, []byte(message))
	if err != nil {
		return err
	}

	return nil
}
