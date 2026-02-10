package telegram

import (
	"fmt"
	"time"

	"health-check-service/config"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// Client represents the Telegram Bot API client
type Client struct {
	botToken string
	client   *resty.Client
}

// SendMessageRequest represents the request to send a message
type SendMessageRequest struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

// NewClient creates a new Telegram client
func NewClient(cfg *config.Config) *Client {
	client := resty.New().
		SetTimeout(30*time.Second).
		SetHeader("Content-Type", "application/json")

	return &Client{
		botToken: cfg.TelegramBotToken,
		client:   client,
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

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.botToken)

	logrus.Info("Sending message to Telegram chat", chatID)
	logrus.Info("Payload: ", payload)
	resp, err := c.client.R().
		SetBody(payload).
		Post(url)

	if err != nil {
		return err
	}

	if resp.IsError() {
		logrus.Errorf("Telegram API error: %s", resp.String())
		return fmt.Errorf("telegram API error: %s", resp.Status())
	}

	logrus.Infof("Message sent to Telegram chat %s", chatID)
	return nil
}
