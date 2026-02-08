package waha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"health-check-service/config"

	"github.com/sirupsen/logrus"
)

// Client represents the WAHA API client
type Client struct {
	baseURL    string
	session    string
	apiKey     string
	httpClient *http.Client
}

// SendTextRequest represents the request to send text message
type SendTextRequest struct {
	ChatID string `json:"chatId"`
	Text   string `json:"text"`
}

// NewClient creates a new WAHA client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		baseURL: cfg.WahaAPIURL,
		session: cfg.WahaSession,
		apiKey:  cfg.WahaAPIKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendText sends a text message to the specified phone number
func (c *Client) SendText(phone, message string) error {
	// Format chatId for WhatsApp
	chatID := fmt.Sprintf("%s@c.us", phone)

	payload := SendTextRequest{
		ChatID: chatID,
		Text:   message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/api/sendText", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("X-Api-Key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("WAHA API error: status code %d", resp.StatusCode)
	}

	logrus.Infof("Message sent to %s via WAHA", phone)
	return nil
}

// SendTextWithSession sends a text message using specific session
func (c *Client) SendTextWithSession(session, phone, message string) error {
	chatID := fmt.Sprintf("%s@c.us", phone)

	payload := SendTextRequest{
		ChatID: chatID,
		Text:   message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/api/sessions/%s/sendText", c.baseURL, session)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("X-Api-Key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("WAHA API error: status code %d", resp.StatusCode)
	}

	logrus.Infof("Message sent to %s via WAHA session %s", phone, session)
	return nil
}
