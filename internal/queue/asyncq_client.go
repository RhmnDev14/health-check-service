package queue

import (
	"health-check-service/config"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

var Client *asynq.Client

// InitClient initializes the AsyncQ client
func InitClient(cfg *config.Config) {
	Client = asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})
	logrus.Info("AsyncQ client initialized")
}

// CloseClient closes the AsyncQ client
func CloseClient() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// EnqueueNotification enqueues a WhatsApp notification task
func EnqueueNotification(payload NotificationPayload) error {
	task, err := NewNotificationTask(payload)
	if err != nil {
		return err
	}

	info, err := Client.Enqueue(task,
		asynq.MaxRetry(3),
		asynq.Queue("notifications"),
	)
	if err != nil {
		return err
	}

	logrus.Infof("Enqueued notification task: %s", info.ID)
	return nil
}
