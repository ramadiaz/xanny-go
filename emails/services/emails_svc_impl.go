package emails

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"xanny-go-template/emails/dto"
	"xanny-go-template/pkg/exceptions"

	"gopkg.in/gomail.v2"
)

func SendEmail(data dto.EmailRequest) *exceptions.Exception {
	email := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	server := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")

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
	body := fmt.Sprintf(
		`<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8" />
				<meta name="viewport" content="width=device-width, initial-scale=1.0" />
				<link rel="icon" type="img" href="/tixchain-logo.png" />
				<title>Email Example</title>
			</head>
			<body>
				<p>%s</p>
			</body>
			</html>
		`, data.Body)

	emailData := dto.EmailRequest{
		Email:   data.Email,
		Subject: data.Subject,
		Body:    body,
	}

	err := SendEmail(emailData)
	if err != nil {
		return err
	}

	return nil
}