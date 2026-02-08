package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HealthLog represents a health check log entry
type HealthLog struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ServiceID        string             `bson:"service_id" json:"service_id"`
	ServiceName      string             `bson:"service_name" json:"service_name"`
	Status           string             `bson:"status" json:"status"` // UP, DOWN
	ResponseTime     int64              `bson:"response_time_ms" json:"response_time_ms"`
	ErrorMessage     string             `bson:"error_message,omitempty" json:"error_message,omitempty"`
	CheckedAt        time.Time          `bson:"checked_at" json:"checked_at"`
	NotificationSent bool               `bson:"notification_sent" json:"notification_sent"`
}

const (
	StatusUP   = "UP"
	StatusDOWN = "DOWN"
)
