package internal

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

func Send(to, subject, body string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("<%s>", os.Getenv("EMAIL_ADDRESS"))
	e.To = []string{os.Getenv("RECIPIENT")}
	e.Subject = subject
	e.HTML = []byte(body)
	return e.Send(fmt.Sprintf("%s:%s", os.Getenv("EMAIL_HOST"), os.Getenv("EMAIL_PORT")), smtp.PlainAuth("", os.Getenv("EMAIL_ADDRESS"), os.Getenv("EMAIL_PASSWORD"), os.Getenv("EMAIL_HOST")))
}

func SendWithAttachment(to, subject, body, attachmentFile string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("<%s>", os.Getenv("EMAIL_ADDRESS"))
	e.To = []string{os.Getenv("RECIPIENT")}
	e.Subject = subject
	e.HTML = []byte(body)
	_, err := e.AttachFile(attachmentFile)
	if err != nil {
		return err
	}
	return e.Send(fmt.Sprintf("%s:%s", os.Getenv("EMAIL_HOST"), os.Getenv("EMAIL_PORT")), smtp.PlainAuth("", os.Getenv("EMAIL_ADDRESS"), os.Getenv("EMAIL_PASSWORD"), os.Getenv("EMAIL_HOST")))
}
