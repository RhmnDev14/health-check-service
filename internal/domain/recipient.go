package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Notification channel constants
const (
	ChannelWhatsApp = "whatsapp"
	ChannelTelegram = "telegram"
	ChannelDiscord  = "discord"
)

// Recipient represents a notification recipient
type Recipient struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	Channel    string             `bson:"channel" json:"channel"`                             // whatsapp, telegram, email
	Phone      string             `bson:"phone,omitempty" json:"phone,omitempty"`             // WhatsApp
	TelegramID string             `bson:"telegram_id,omitempty" json:"telegram_id,omitempty"` // Telegram chat ID
	DiscordID  string             `bson:"discord_id,omitempty" json:"discord_id,omitempty"`   // Discord channel/user ID
	IsActive   bool               `bson:"is_active" json:"is_active"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

// CreateRecipientRequest represents request to create a recipient
type CreateRecipientRequest struct {
	Name       string `json:"name" binding:"required"`
	Channel    string `json:"channel" binding:"required,oneof=whatsapp telegram discord"`
	Phone      string `json:"phone,omitempty"`
	TelegramID string `json:"telegram_id,omitempty"`
	DiscordID  string `json:"discord_id,omitempty"`
}

// UpdateRecipientRequest represents request to update a recipient
type UpdateRecipientRequest struct {
	Name       string `json:"name,omitempty"`
	Channel    string `json:"channel,omitempty"`
	Phone      string `json:"phone,omitempty"`
	TelegramID string `json:"telegram_id,omitempty"`
	DiscordID  string `json:"discord_id,omitempty"`
	IsActive   *bool  `json:"is_active,omitempty"`
}

// ValidChannels returns list of valid channels
func ValidChannels() []string {
	return []string{ChannelWhatsApp, ChannelTelegram, ChannelDiscord}
}

// IsValidChannel checks if a channel is valid
func IsValidChannel(channel string) bool {
	for _, c := range ValidChannels() {
		if c == channel {
			return true
		}
	}
	return false
}
