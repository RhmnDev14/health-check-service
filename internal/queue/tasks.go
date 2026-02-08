package queue

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeNotificationWhatsApp = "notification:whatsapp"
)

// NotificationPayload represents the payload for WhatsApp notification task
type NotificationPayload struct {
	ServiceID    string `json:"service_id"`
	ServiceName  string `json:"service_name"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
	LogID        string `json:"log_id"`
}

// NewNotificationTask creates a new WhatsApp notification task
func NewNotificationTask(payload NotificationPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeNotificationWhatsApp, data), nil
}
