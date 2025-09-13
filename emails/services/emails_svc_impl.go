package emails

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"
	"xanny-go/emails/dto"
	"xanny-go/pkg/config"
	"xanny-go/pkg/exceptions"

	"gopkg.in/gomail.v2"
)

func SendEmail(data dto.EmailRequest) *exceptions.Exception {
	email := config.GetSMTPEmail()
	password := config.GetSMTPPassword()
	server := config.GetSMTPServer()
	smtpPort := config.GetSMTPPort()

	i, err := strconv.Atoi(smtpPort)
	if err != nil {
		return exceptions.NewException(http.StatusInternalServerError, exceptions.ErrInternalServer)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", email)
	m.SetHeader("To", data.Email)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", data.Body)

	d := gomail.NewDialer(server, i, email, password)

	if err := d.DialAndSend(m); err != nil {
		return exceptions.NewException(http.StatusBadGateway, err.Error())
	}

	return nil
}

func ExampleEmail(data dto.EmailExample) *exceptions.Exception {
	tmpl, exc := template.ParseFiles("emails/templates/example.html")
	if exc != nil {
		return exceptions.NewException(http.StatusInternalServerError, exc.Error())
	}

	var body bytes.Buffer
	if exc := tmpl.Execute(&body, data); exc != nil {
		return exceptions.NewException(http.StatusInternalServerError, exc.Error())
	}

	emailData := dto.EmailRequest{
		Email:   data.Email,
		Subject: data.Subject,
		Body:    body.String(),
	}

	err := SendEmail(emailData)
	if err != nil {
		return err
	}

	return nil
}

func VerificationEmail(data dto.EmailVerification) *exceptions.Exception {
	tmpl, exc := template.ParseFiles("emails/templates/verification.html")
	if exc != nil {
		return exceptions.NewException(http.StatusInternalServerError, exc.Error())
	}

	var body bytes.Buffer
	if exc := tmpl.Execute(&body, data); exc != nil {
		return exceptions.NewException(http.StatusInternalServerError, exc.Error())
	}

	emailData := dto.EmailRequest{
		Email:   data.Email,
		Subject: "[Xanware] Email Verification for Account Activation",
		Body:    body.String(),
	}

	err := SendEmail(emailData)
	if err != nil {
		return err
	}

	return nil
}
