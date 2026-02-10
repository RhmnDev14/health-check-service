package discord

import (
	"fmt"
	"time"

	"health-check-service/config"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// Client represents the Discord API client
type Client struct {
	botToken string
	client   *resty.Client
}

// SendMessageRequest represents the request to send a message
type SendMessageRequest struct {
	Content string `json:"content"`
}

// NewClient creates a new Discord client
func NewClient(cfg *config.Config) *Client {
	client := resty.New().
		SetBaseURL("https://discord.com/api/v10").
		SetTimeout(30*time.Second).
		SetHeader("Content-Type", "application/json")

	if cfg.DiscordBotToken != "" {
		client.SetHeader("Authorization", fmt.Sprintf("Bot %s", cfg.DiscordBotToken))
	}

	return &Client{
		botToken: cfg.DiscordBotToken,
		client:   client,
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

	resp, err := c.client.R().
		SetBody(payload).
		Post(fmt.Sprintf("/channels/%s/messages", channelID))

	if err != nil {
		return err
	}

	if resp.IsError() {
		logrus.Errorf("Discord API error: %s", resp.String())
		return fmt.Errorf("discord API error: %s", resp.Status())
	}

	logrus.Infof("Message sent to Discord channel %s", channelID)
	return nil
}
