package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"health-check-service/config"

	"github.com/sirupsen/logrus"
)

// Client represents the Telegram Bot API client
type Client struct {
	botToken   string
	httpClient *http.Client
}

// SendMessageRequest represents the request to send a message
type SendMessageRequest struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

// NewClient creates a new Telegram client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		botToken: cfg.TelegramBotToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IsConfigured returns true if Telegram is configured
func (c *Client) IsConfigured() bool {
	return c.botToken != ""
}

// SendMessage sends a text message to the specified chat ID
func (c *Client) SendMessage(chatID, message string) error {
	if !c.IsConfigured() {
		return fmt.Errorf("telegram bot token not configured")
	}

	payload := SendMessageRequest{
		ChatID:    chatID,
		Text:      message,
		ParseMode: "Markdown",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.botToken)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("telegram API error: status code %d", resp.StatusCode)
	}

	logrus.Infof("Message sent to Telegram chat %s", chatID)
	return nil
}
