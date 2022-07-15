package mail

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"
	"os"
)

type Mail struct {
	from    string
	to      []string
	subject string
	message string
}

func NewMail(to []string, subject string) *Mail {
	return &Mail{
		to:      to,
		subject: subject,
	}
}
func (m *Mail) parseTemplate(filename string, data interface{}) error {
	t, err := template.ParseFiles(filename)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err := t.Execute(buffer, data); err != nil {
		return err
	}
	m.message = buffer.String()
	return nil
}

func (m *Mail) sendMail() bool {
	body := "To: " + m.to[0] + "\r\nSubject: " + m.subject + "\r\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n\r\n" + m.message
	auth := smtp.PlainAuth("", os.Getenv("MAIL_USER"), os.Getenv("MAIL_PASS"), os.Getenv("MAIL_HOST"))
	if err := smtp.SendMail(os.Getenv("MAIL_HOST")+":"+os.Getenv("MAIL_PORT"), auth, m.from, m.to, []byte(body)); err != nil {
		return false
	}
	return true
}

func (m *Mail) Send(filename string, data interface{}) {

	if err := m.parseTemplate(filename, data); err != nil {
		log.Fatal(err)
	}
	if ok := m.sendMail(); ok {
		log.Printf("Email has been sent to %s\n", m.to)
	} else {
		log.Fatalf("Failed to send the email to %s\n", m.to)
	}
}
