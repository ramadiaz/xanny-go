package dto

type EmailRequest struct {
	Email   string
	Subject string
	Body    string
}

type EmailExample struct {
	Message string
	Email   string
	Subject string
	Body    string
}
