package waha

import (
	"fmt"
	"time"

	"health-check-service/config"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// Client represents the WAHA API client
type Client struct {
	baseURL string
	session string
	apiKey  string
	client  *resty.Client
}

// SendTextRequest represents the request to send text message
type SendTextRequest struct {
	ChatID string `json:"chatId"`
	Text   string `json:"text"`
}

// NewClient creates a new WAHA client
func NewClient(cfg *config.Config) *Client {
	client := resty.New().
		SetBaseURL(cfg.WahaAPIURL).
		SetTimeout(30*time.Second).
		SetHeader("Content-Type", "application/json")

	if cfg.WahaAPIKey != "" {
		client.SetHeader("X-Api-Key", cfg.WahaAPIKey)
	}

	logrus.Infof("WAHA client initialized")
	logrus.Infof("WAHA API URL: %s", cfg.WahaAPIURL)
	logrus.Infof("WAHA API Key: %s", cfg.WahaAPIKey)
	logrus.Infof("WAHA Session: %s", cfg.WahaSession)

	return &Client{
		baseURL: cfg.WahaAPIURL,
		session: cfg.WahaSession,
		apiKey:  cfg.WahaAPIKey,
		client:  client,
	}
}

// SendText sends a text message
func (c *Client) SendText(phone, message string) error {
	if c.session != "" {
		return c.SendTextWithSession(c.session, phone, message)
	}

	chatID := fmt.Sprintf("%s@c.us", phone)

	payload := SendTextRequest{
		ChatID: chatID,
		Text:   message,
	}

	logrus.Infof("Sending message to %s via WAHA", phone)

	resp, err := c.client.R().
		SetBody(payload).
		Post("/api/sendText")

	if err != nil {
		logrus.Errorf("Failed to send request: %v", err)
		return err
	}

	if resp.IsError() {
		logrus.Errorf("WAHA API error: %s", resp.String())
		return fmt.Errorf("WAHA API error: %s", resp.Status())
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

	logrus.Infof("Sending message to %s via WAHA session %s", phone, session)

	resp, err := c.client.R().
		SetBody(payload).
		Post(fmt.Sprintf("/api/sessions/%s/sendText", session))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("WAHA API error: %s", resp.Status())
	}

	logrus.Infof("Message sent to %s via WAHA session %s", phone, session)
	return nil
}
