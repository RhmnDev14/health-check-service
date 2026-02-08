package queue

import (
	"context"
	"encoding/json"

	"health-check-service/config"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

// NotificationSender interface for sending notifications
type NotificationSender interface {
	SendNotification(ctx context.Context, payload NotificationPayload) error
}

// Server represents the AsyncQ server/worker
type Server struct {
	server   *asynq.Server
	notifier NotificationSender
}

// NewServer creates a new AsyncQ server
func NewServer(cfg *config.Config, notifier NotificationSender) *Server {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"notifications": 10,
				"default":       5,
			},
		},
	)

	return &Server{
		server:   srv,
		notifier: notifier,
	}
}

// Start starts the AsyncQ server
func (s *Server) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeNotificationWhatsApp, s.handleNotification)

	logrus.Info("Starting AsyncQ server...")
	return s.server.Start(mux)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() {
	s.server.Shutdown()
}

// handleNotification processes WhatsApp notification tasks
func (s *Server) handleNotification(ctx context.Context, task *asynq.Task) error {
	var payload NotificationPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		logrus.Errorf("Failed to unmarshal notification payload: %v", err)
		return err
	}

	logrus.Infof("Processing notification for service: %s (status: %s)", payload.ServiceName, payload.Status)

	// Send notification to all active recipients
	err := s.notifier.SendNotification(ctx, payload)
	if err != nil {
		logrus.Errorf("Failed to send notification: %v", err)
		return err
	}

	return nil
}
