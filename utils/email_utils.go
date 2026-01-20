package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

func getAccessToken() (string, error) {
	data := url.Values{}
	data.Set("client_id", os.Getenv("AZURE_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("AZURE_CLIENT_SECRET"))
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "https://graph.microsoft.com/.default")

	req, _ := http.NewRequest(
		"POST",
		"https://login.microsoftonline.com/"+os.Getenv("AZURE_TENANT_ID")+"/oauth2/v2.0/token",
		strings.NewReader(data.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var r struct {
		AccessToken string `json:"access_token"`
	}
	json.NewDecoder(resp.Body).Decode(&r)

	return r.AccessToken, nil
}

func SendEmailOTP(to, otp string, templateLocation string) error {
	token, err := getAccessToken()
	if err != nil {
		return err
	}

	// 🔹 Load HTML template
	htmlBody, err := LoadOTPTemplate(templateLocation, OTPData{
		AppName: os.Getenv("APP_NAME"),
		Code:    otp,
		Minutes: 5,
	})
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"message": map[string]interface{}{
			"subject": "Your OTP Code",
			"body": map[string]string{
				"contentType": "HTML",
				"content":     htmlBody,
			},
			"toRecipients": []map[string]interface{}{
				{
					"emailAddress": map[string]string{
						"address": to,
					},
				},
			},
		},
		"saveToSentItems": false,
	}

	b, _ := json.Marshal(payload)

	req, _ := http.NewRequest(
		"POST",
		"https://graph.microsoft.com/v1.0/users/"+os.Getenv("EMAIL_FROM")+"/sendMail",
		bytes.NewBuffer(b),
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("send mail failed: %s", resp.Status)
	}

	return nil
}

type OTPData struct {
	AppName string
	Code    string
	Minutes int
}

func LoadOTPTemplate(path string, data OTPData) (string, error) {
	tpl, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// func SendEmailOTP(email, otp string) error {
// 	username := os.Getenv("EMAIL_USERNAME")
// 	password := os.Getenv("EMAIL_PASSWORD")
// 	host := os.Getenv("EMAIL_HOST")
// 	if host == "" {
// 		host = "smtp.gmail.com"
// 	}
// 	port := os.Getenv("EMAIL_PORT")
// 	if port == "" {
// 		port = "587"
// 	}
// 	if username == "" || password == "" || host == "" {
// 		log.Fatal("EMAIL_USERNAME or EMAIL_PASSWORD not set in environment")
// 	}
// 	auth := smtp.PlainAuth("", username, password, host)
// 	msg := []byte("Subject: Your OTP Code\n\nYour OTP is: " + otp)
// 	return smtp.SendMail(host+":"+port, auth, username, []string{email}, msg)
// }
