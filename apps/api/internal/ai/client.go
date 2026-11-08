package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ProviderConfig struct {
	Name    string
	BaseURL string
	Model   string
	APIKey  string
	Timeout time.Duration
}

type Client struct {
	cfg  ProviderConfig
	http *http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

type ContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type UpstreamError struct {
	Status int
	Body   string
}

func (e *UpstreamError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("AI upstream returned HTTP %d", e.Status)
	}
	return fmt.Sprintf("AI upstream returned HTTP %d: %s", e.Status, e.Body)
}

var (
	ErrInvalidOutput = errors.New("AI returned invalid structured output")
	ErrNotDetected   = errors.New("aksara Bali tidak terdeteksi pada gambar")
)

func NewClient(cfg ProviderConfig) *Client {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 60 * time.Second
	}
	return &Client{
		cfg:  cfg,
		http: &http.Client{Timeout: cfg.Timeout},
	}
}

func (c *Client) Model() string {
	return c.cfg.Model
}

func (c *Client) Complete(ctx context.Context, messages []Message) (string, error) {
	return c.CompleteWithTokens(ctx, messages, 2400)
}

func (c *Client) CompleteWithTokens(ctx context.Context, messages []Message, maxTokens int) (string, error) {
	body := map[string]any{
		"model":       c.cfg.Model,
		"messages":    messages,
		"stream":      false,
		"temperature": 0,
		"top_p":       1,
		"max_tokens":  maxTokens,
		"response_format": map[string]string{
			"type": "json_object",
		},
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("encode AI request: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			delay := time.Duration(attempt) * time.Second
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(delay):
			}
		}

		content, err := c.do(ctx, payload)
		if err == nil {
			return content, nil
		}
		lastErr = err
		if !retryable(err) {
			break
		}
	}
	return "", lastErr
}

func (c *Client) do(ctx context.Context, payload []byte) (string, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		strings.TrimRight(c.cfg.BaseURL, "/")+"/chat/completions",
		bytes.NewReader(payload),
	)
	if err != nil {
		return "", fmt.Errorf("create AI request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	if c.cfg.APIKey != "" {
		request.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	}

	response, err := c.http.Do(request)
	if err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		return "", fmt.Errorf("AI upstream request failed: %w", err)
	}
	defer response.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("read AI response: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return "", &UpstreamError{
			Status: response.StatusCode,
			Body:   upstreamMessage(raw),
		}
	}

	var decoded struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return "", &UpstreamError{Status: response.StatusCode, Body: "respons AI bukan JSON"}
	}
	if decoded.Error != nil && decoded.Error.Message != "" {
		return "", &UpstreamError{Status: response.StatusCode, Body: decoded.Error.Message}
	}
	if len(decoded.Choices) == 0 || strings.TrimSpace(decoded.Choices[0].Message.Content) == "" {
		return "", &UpstreamError{Status: response.StatusCode, Body: "respons AI kosong"}
	}
	return decoded.Choices[0].Message.Content, nil
}

func (c *Client) Health(ctx context.Context) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		strings.TrimRight(c.cfg.BaseURL, "/")+"/models",
		nil,
	)
	if err != nil {
		return err
	}
	if c.cfg.APIKey != "" {
		request.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	}
	response, err := c.http.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("9router health HTTP %d", response.StatusCode)
	}
	return nil
}

func retryable(err error) bool {
	var upstream *UpstreamError
	if errors.As(err, &upstream) {
		return upstream.Status == http.StatusTooManyRequests || upstream.Status >= 500
	}
	var network net.Error
	if errors.As(err, &network) && network.Timeout() {
		return true
	}
	var urlErr *url.Error
	return errors.As(err, &urlErr)
}

func shouldFallback(err error) bool {
	return retryable(err)
}

func upstreamMessage(raw []byte) string {
	var decoded struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if json.Unmarshal(raw, &decoded) == nil && decoded.Error.Message != "" {
		return truncate(decoded.Error.Message, 400)
	}
	return truncate(strings.TrimSpace(string(raw)), 400)
}

func truncate(value string, limit int) string {
	if len(value) <= limit {
		return value
	}
	return value[:limit] + "…"
}
