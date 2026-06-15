package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type resendPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

type resendError struct {
	StatusCode int    `json:"statusCode"`
	Name       string `json:"name"`
	Message    string `json:"message"`
}

func SendOTP(to, code string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		log.Printf("[EMAIL] OTP for %s: %s", to, code)
		return nil
	}

	from := os.Getenv("EMAIL_FROM")
	if from == "" {
		from = "onboarding@resend.dev"
	}

	body := resendPayload{
		From:    from,
		To:      to,
		Subject: "Your URL Shortener Verification Code",
		Text:    fmt.Sprintf("Your verification code is: %s\n\nThis code expires in 15 minutes.", code),
	}

	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("resend request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resend API error (%d): %s", resp.StatusCode, string(respBody))
	}
	return nil
}
