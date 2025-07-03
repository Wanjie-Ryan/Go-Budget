package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/gomail.v2"
)

// the line below tells the go compiler to include all the entire templates/ folder and its files inside the final binary.
//
//go:embed templates
var templateFS embed.FS

// templateFS is now an in memory filesystem you can read from at runtime via template.ParseFS
type Mailer struct {
	dailer *gomail.Dialer
	// sender *gomail.Sender // no-sender@Wanjie-Ryan.com
	sender string
}

type EmailData struct {
	AppName string
	Subject string
	Meta    interface{}
}

func NewMailer() Mailer {
	// anytime you are getting details from an env, they coming as string, eg. the port, hence conver to int
	// strconv.Atoi is used to convert string to integer
	mailPort, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		panic(err)
	}
	mailHost := os.Getenv("MAIL_HOST")
	mailUsername := os.Getenv("MAIL_USERNAME")
	mailPassword := os.Getenv("MAIL_PASSWORD")
	mailSender := os.Getenv("MAIL_SENDER")

	dailer := gomail.NewDialer(mailHost, mailPort, mailUsername, mailPassword)

	return Mailer{dailer: dailer, sender: mailSender}

}

// will be used to send emails
func (mail *Mailer) Send(recipient string, templateFile string, data EmailData) error {
	// generate an absolute path like templates/hello.html
	absolutePath := filepath.Join("templates", templateFile)
	// reads the file out of the embedded FS and parses any {{}} directives in it.
	tmpl, err := template.ParseFS(templateFS, absolutePath)

	if err != nil {
		fmt.Println("template parse err", err)
		return err
	}
	data.AppName = os.Getenv("APP_NAME")
	// the subjectBuf acts like a variable where the subject is stored
	subjectBuf := new(bytes.Buffer)
	// the code below runs the {{define "subject"}} block and writes its output into subjectBuf
	err = tmpl.ExecuteTemplate(subjectBuf, "subject", data)
	if err != nil {
		fmt.Println("subject error", err)
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		fmt.Println("html body error", err)
		return err
	}

	//  after getting the html body and the subject, then create the actual message to send
	message := gomail.NewMessage()
	// who is the recipient
	message.SetHeader("To", recipient)
	// who is the sender
	message.SetHeader("From", mail.sender)
	// the subject of the email
	message.SetHeader("Subject", subjectBuf.String())
	// the body of the email
	message.SetBody("text/html", htmlBody.String())

	// dials first to test the connection, then sends the message
	err = mail.dailer.DialAndSend(message)
	if err != nil {
		fmt.Println("dial and send error", err)
		return err
	}
	return nil

}
