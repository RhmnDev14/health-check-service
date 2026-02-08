package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"health-check-service/config"
)

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
	From User   `json:"from"`
}

type Chat struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

func main() {
	cfg := config.Load()
	token := cfg.TelegramBotToken

	if token == "" {
		log.Fatal("Telegram Bot Token not found in .env")
	}

	log.Printf("Checking updates for bot token: %s...", token[:5]+"...")
	log.Println("Please send a message to your bot now if you haven't already.")

	url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("Failed to get updates: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var updateResp UpdateResponse
	if err := json.Unmarshal(body, &updateResp); err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	if !updateResp.Ok {
		log.Fatalf("Telegram API error: %s", string(body))
	}

	if len(updateResp.Result) == 0 {
		log.Println("No updates found. Send a message to your bot and try again.")
		return
	}

	log.Println("--- Found Chats ---")
	seen := make(map[int64]bool)
	for _, update := range updateResp.Result {
		chatID := update.Message.Chat.ID
		if !seen[chatID] {
			log.Printf("Chat ID: %d | User: %s (@%s) | Message: %s\n",
				chatID,
				update.Message.From.FirstName,
				update.Message.From.Username,
				update.Message.Text,
			)
			seen[chatID] = true
		}
	}
	log.Println("-------------------")
	log.Println("Copy the Chat ID above and update your recipient configuration.")
}
