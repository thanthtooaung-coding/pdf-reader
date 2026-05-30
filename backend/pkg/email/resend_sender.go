package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultResendBaseURL = "https://api.resend.com"

type ResendConfig struct {
	APIKey  string
	BaseURL string
	From    string
}

type resendSender struct {
	cfg        ResendConfig
	httpClient *http.Client
}

func NewResendSender(cfg ResendConfig) Sender {
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultResendBaseURL
	}
	return &resendSender{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *resendSender) Enabled() bool {
	return s.cfg.APIKey != "" && s.cfg.From != ""
}

type resendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Text    string   `json:"text,omitempty"`
	HTML    string   `json:"html,omitempty"`
}

func (s *resendSender) SendOTP(to, otp string) error {
	if !s.Enabled() {
		return fmt.Errorf("resend is not configured")
	}

	subject := "Your PDF Reader verification code"
	text := fmt.Sprintf("Your verification code is: %s\n\nIt expires in 10 minutes.", otp)
	html := fmt.Sprintf(
		`<p>Your verification code is: <strong>%s</strong></p><p>It expires in 10 minutes.</p>`,
		otp,
	)

	payload, err := json.Marshal(resendEmailRequest{
		From:    s.cfg.From,
		To:      []string{to},
		Subject: subject,
		Text:    text,
		HTML:    html,
	})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.BaseURL+"/emails", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("resend api error (%d): %s", resp.StatusCode, string(body))
	}
	return nil
}
