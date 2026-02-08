package main

import (
	"log"

	"health-check-service/config"
	"health-check-service/infrastructure/discord"
	"health-check-service/infrastructure/telegram"
	"health-check-service/infrastructure/waha"
)

func main() {
	cfg := config.Load()

	// Override with manual tokens if needed for local test, or rely on .env
	// cfg.DiscordBotToken = "..."
	// cfg.TelegramBotToken = "..."

	// Initialize clients
	wahaClient := waha.NewClient(cfg)
	telegramClient := telegram.NewClient(cfg)
	discordClient := discord.NewClient(cfg)

	// Prevent unused variable error
	_ = wahaClient

	log.Println("Starting verification...")

	msg := "Test Notification from Health Service"

	// Test Discord
	if cfg.DiscordBotToken != "" {
		err := discordClient.SendMessage("REPLACE_WITH_CHANNEL_ID", msg)
		if err != nil {
			log.Printf("Discord failed: %v", err)
		} else {
			log.Println("Discord success")
		}
	} else {
		log.Println("Discord token missing")
	}

	// Test Telegram
	if cfg.TelegramBotToken != "" {
		err := telegramClient.SendMessage("REPLACE_WITH_CHAT_ID", msg)
		if err != nil {
			log.Printf("Telegram failed: %v", err)
		} else {
			log.Println("Telegram success")
		}
	} else {
		log.Println("Telegram token missing")
	}
}
