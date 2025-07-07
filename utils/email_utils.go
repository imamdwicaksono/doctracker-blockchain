package utils

import (
	"log"
	"net/smtp"
	"os"
)

func SendEmailOTP(email, otp string) error {
	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")
	host := os.Getenv("EMAIL_HOST")
	if host == "" {
		host = "smtp.gmail.com"
	}
	port := os.Getenv("EMAIL_PORT")
	if port == "" {
		port = "587"
	}
	if username == "" || password == "" || host == "" {
		log.Fatal("EMAIL_USERNAME or EMAIL_PASSWORD not set in environment")
	}
	auth := smtp.PlainAuth("", username, password, host)
	msg := []byte("Subject: Your OTP Code\n\nYour OTP is: " + otp)
	return smtp.SendMail(host+":"+port, auth, username, []string{email}, msg)
}
