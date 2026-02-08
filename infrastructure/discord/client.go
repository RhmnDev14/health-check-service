package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"health-check-service/config"

	"github.com/sirupsen/logrus"
)

// Client represents the Discord API client
type Client struct {
	botToken   string
	httpClient *http.Client
}

// SendMessageRequest represents the request to send a message
type SendMessageRequest struct {
	Content string `json:"content"`
}

// NewClient creates a new Discord client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		botToken: cfg.DiscordBotToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IsConfigured returns true if Discord is configured
func (c *Client) IsConfigured() bool {
	return c.botToken != ""
}

// SendMessage sends a text message to the specified channel ID
func (c *Client) SendMessage(channelID, message string) error {
	if !c.IsConfigured() {
		return fmt.Errorf("discord bot token not configured")
	}

	payload := SendMessageRequest{
		Content: message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", channelID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", c.botToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("discord API error: status code %d", resp.StatusCode)
	}

	logrus.Infof("Message sent to Discord channel %s", channelID)
	return nil
}
