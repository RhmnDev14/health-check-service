package domain

import "time"

// Service represents a service to be monitored
type Service struct {
	ID            string    `bson:"_id" json:"id"`
	Name          string    `bson:"name" json:"name"`
	URL           string    `bson:"url" json:"url"`
	CheckInterval int       `bson:"check_interval" json:"check_interval"` // seconds
	RetryInterval int       `bson:"retry_interval" json:"retry_interval"` // seconds for notification retry
	IsActive      bool      `bson:"is_active" json:"is_active"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}

// CreateServiceRequest represents request to create a service
type CreateServiceRequest struct {
	ID            string `json:"id" binding:"required"`
	Name          string `json:"name" binding:"required"`
	URL           string `json:"url" binding:"required"`
	CheckInterval int    `json:"check_interval"`
	RetryInterval int    `json:"retry_interval"`
}

// UpdateServiceRequest represents request to update a service
type UpdateServiceRequest struct {
	Name          string `json:"name"`
	URL           string `json:"url"`
	CheckInterval int    `json:"check_interval"`
	RetryInterval int    `json:"retry_interval"`
	IsActive      *bool  `json:"is_active"`
}
