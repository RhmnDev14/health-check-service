package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"health-check-service/config"
	"health-check-service/infrastructure/discord"
	"health-check-service/infrastructure/telegram"
	"health-check-service/infrastructure/waha"
	"health-check-service/internal/api"
	"health-check-service/internal/queue"
	"health-check-service/internal/repository/mongodb"
	"health-check-service/internal/scheduler"
	"health-check-service/internal/usecase"

	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg := config.Load()
	logrus.Info("Configuration loaded")

	// Connect to MongoDB
	if err := mongodb.Connect(cfg); err != nil {
		logrus.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongodb.Disconnect()

	// Initialize repositories
	serviceRepo := mongodb.NewServiceRepository()
	recipientRepo := mongodb.NewRecipientRepository()
	healthLogRepo := mongodb.NewHealthLogRepository()

	// Ensure indexes
	ctx := context.Background()
	serviceRepo.EnsureIndexes(ctx)
	recipientRepo.EnsureIndexes(ctx)
	healthLogRepo.EnsureIndexes(ctx)

	// Initialize AsyncQ client
	queue.InitClient(cfg)
	defer queue.CloseClient()

	// Initialize WAHA client
	wahaClient := waha.NewClient(cfg)

	// Initialize Telegram client
	telegramClient := telegram.NewClient(cfg)

	// Initialize Discord client
	discordClient := discord.NewClient(cfg)

	// Initialize usecases
	notifier := usecase.NewNotifier(wahaClient, telegramClient, discordClient, recipientRepo, healthLogRepo)
	healthChecker := usecase.NewHealthChecker(serviceRepo, healthLogRepo, cfg)

	// Initialize and start AsyncQ server
	asyncqServer := queue.NewServer(cfg, notifier)
	go func() {
		if err := asyncqServer.Start(); err != nil {
			logrus.Errorf("AsyncQ server error: %v", err)
		}
	}()
	defer asyncqServer.Shutdown()

	// Initialize scheduler
	sched := scheduler.NewScheduler(healthChecker)
	sched.Start(ctx)
	defer sched.Stop()

	// Schedule all active services
	services, err := serviceRepo.GetActive(ctx)
	if err == nil {
		for _, service := range services {
			sched.ScheduleService(&service)
		}
	}

	// Initialize and start API router
	router := api.NewRouter(cfg, sched)
	router.Setup(serviceRepo, recipientRepo, healthLogRepo)

	// Start server in goroutine
	go func() {
		logrus.Infof("Starting HTTP server on port %s", cfg.AppPort)
		if err := router.Run(); err != nil {
			logrus.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")
}
