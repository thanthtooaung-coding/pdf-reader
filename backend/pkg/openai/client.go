package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.openai.com/v1"
const defaultModel = "gpt-4o-mini"

type Client interface {
	Enabled() bool
	Summarize(ctx context.Context, documentText string) (string, error)
	Translate(ctx context.Context, documentText, targetLanguage string) (string, error)
}

type Config struct {
	APIKey  string
	Model   string
	BaseURL string
	Timeout time.Duration
}

type client struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

func NewClient(cfg Config) Client {
	if strings.TrimSpace(cfg.APIKey) == "" {
		return &stubClient{}
	}

	model := cfg.Model
	if model == "" {
		model = defaultModel
	}
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 120 * time.Second
	}

	return &client{
		apiKey:  cfg.APIKey,
		model:   model,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

type stubClient struct{}

func (s *stubClient) Enabled() bool { return false }

func (s *stubClient) Summarize(_ context.Context, _ string) (string, error) {
	return "[stub] Summary of the PDF content. Set OPENAI_API_KEY to use OpenAI.", nil
}

func (s *stubClient) Translate(_ context.Context, _ string, targetLanguage string) (string, error) {
	lang := targetLanguage
	if lang == "" {
		lang = "en"
	}
	return fmt.Sprintf("[stub] Translated document to %s. Set OPENAI_API_KEY to use OpenAI.", lang), nil
}

func (c *client) Enabled() bool { return true }

func (c *client) Summarize(ctx context.Context, documentText string) (string, error) {
	system := "You are a helpful assistant that summarizes documents clearly and concisely."
	user := "Summarize the following document. Use bullet points when helpful.\n\n" + documentText
	return c.chat(ctx, system, user)
}

func (c *client) Translate(ctx context.Context, documentText, targetLanguage string) (string, error) {
	lang := strings.TrimSpace(targetLanguage)
	if lang == "" {
		lang = "English"
	}
	system := "You are a professional translator. Preserve meaning and tone."
	user := fmt.Sprintf("Translate the following document to %s:\n\n%s", lang, documentText)
	return c.chat(ctx, system, user)
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (c *client) chat(ctx context.Context, system, user string) (string, error) {
	payload, err := json.Marshal(chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("openai request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}

	var out chatResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("decode openai response: %w", err)
	}
	if out.Error != nil {
		return "", errors.New(out.Error.Message)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("openai api error (%d): %s", resp.StatusCode, string(body))
	}
	if len(out.Choices) == 0 || strings.TrimSpace(out.Choices[0].Message.Content) == "" {
		return "", errors.New("openai returned empty response")
	}
	return strings.TrimSpace(out.Choices[0].Message.Content), nil
}
