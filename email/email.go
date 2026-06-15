package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func SendOTP(to, code string) error {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		log.Printf("[EMAIL] OTP for %s: %s", to, code)
		return nil
	}

	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = user
	}

	body := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Your URL Shortener Verification Code\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\nYour verification code is: %s\r\n\r\nThis code expires in 15 minutes.\r\n", from, to, code)

	addr := host + ":" + port
	auth := smtp.PlainAuth("", user, pass, host)

	return smtp.SendMail(addr, auth, from, []string{to}, []byte(body))
}
