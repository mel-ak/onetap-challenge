package notification

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
)

type EmailNotifier struct {
	from     string
	to       string
	host     string
	port     string
	username string
	password string
}

func NewEmailNotifier(from, to, host, port, username, password string) *EmailNotifier {
	return &EmailNotifier{
		from:     from,
		to:       to,
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (n *EmailNotifier) NotifyAdmin(ctx context.Context, message string, severity string) error {
	subject := fmt.Sprintf("[%s] System Notification", severity)
	body := fmt.Sprintf("Message: %s\nSeverity: %s", message, severity)
	return n.sendEmail(subject, body)
}

func (n *EmailNotifier) NotifyError(ctx context.Context, err error, context string) error {
	subject := "[ERROR] System Error"
	body := fmt.Sprintf("Context: %s\nError: %v", context, err)
	return n.sendEmail(subject, body)
}

func (n *EmailNotifier) sendEmail(subject, body string) error {
	auth := smtp.PlainAuth("", n.username, n.password, n.host)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		n.from, n.to, subject, body)

	addr := fmt.Sprintf("%s:%s", n.host, n.port)
	if err := smtp.SendMail(addr, auth, n.from, []string{n.to}, []byte(msg)); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}
	return nil
}
