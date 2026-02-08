package usecase

import (
	"context"
	"fmt"
	"time"

	"health-check-service/infrastructure/discord"
	"health-check-service/infrastructure/telegram"
	"health-check-service/infrastructure/waha"
	"health-check-service/internal/domain"
	"health-check-service/internal/queue"
	"health-check-service/internal/repository/mongodb"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Notifier handles sending notifications
type Notifier struct {
	wahaClient     *waha.Client
	telegramClient *telegram.Client
	discordClient  *discord.Client
	recipientRepo  *mongodb.RecipientRepository
	healthLogRepo  *mongodb.HealthLogRepository
}

// NewNotifier creates a new Notifier
func NewNotifier(
	wahaClient *waha.Client,
	telegramClient *telegram.Client,
	discordClient *discord.Client,
	recipientRepo *mongodb.RecipientRepository,
	healthLogRepo *mongodb.HealthLogRepository,
) *Notifier {
	return &Notifier{
		wahaClient:     wahaClient,
		telegramClient: telegramClient,
		discordClient:  discordClient,
		recipientRepo:  recipientRepo,
		healthLogRepo:  healthLogRepo,
	}
}

// SendNotification sends notification to all active recipients
// Implements queue.NotificationSender interface
func (n *Notifier) SendNotification(ctx context.Context, payload queue.NotificationPayload) error {
	// Get all active recipients
	recipients, err := n.recipientRepo.GetActive(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active recipients: %w", err)
	}

	if len(recipients) == 0 {
		logrus.Warn("No active recipients found, skipping notification")
		return nil
	}

	// Build message
	message := n.buildMessage(payload)

	// Send to all recipients
	var lastErr error
	successCount := 0
	for _, recipient := range recipients {
		var err error
		switch recipient.Channel {
		case domain.ChannelWhatsApp:
			err = n.wahaClient.SendText(recipient.Phone, message)
		case domain.ChannelTelegram:
			err = n.telegramClient.SendMessage(recipient.TelegramID, message)
		case domain.ChannelDiscord:
			err = n.discordClient.SendMessage(recipient.DiscordID, message)
		default:
			logrus.Warnf("Unknown channel %s for recipient %s", recipient.Channel, recipient.Name)
			continue
		}

		if err != nil {
			logrus.Errorf("Failed to send notification to %s via %s: %v", recipient.Name, recipient.Channel, err)
			lastErr = err
		} else {
			successCount++
		}
	}

	logrus.Infof("Notification sent to %d/%d recipients", successCount, len(recipients))

	// Update health log notification_sent status
	if payload.LogID != "" {
		logID, err := primitive.ObjectIDFromHex(payload.LogID)
		if err == nil {
			n.healthLogRepo.UpdateNotificationSent(ctx, logID, true)
		}
	}

	return lastErr
}

// buildMessage builds the notification message
func (n *Notifier) buildMessage(payload queue.NotificationPayload) string {
	var statusEmoji string
	if payload.Status == domain.StatusDOWN {
		statusEmoji = "🔴"
	} else {
		statusEmoji = "🟢"
	}

	message := fmt.Sprintf(
		"%s *Health Check Alert*\n\n"+
			"*Service:* %s\n"+
			"*Status:* %s %s\n"+
			"*Time:* %s",
		statusEmoji,
		payload.ServiceName,
		statusEmoji,
		payload.Status,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	if payload.ErrorMessage != "" {
		message += fmt.Sprintf("\n*Error:* %s", payload.ErrorMessage)
	}

	return message
}
