package usecase

import (
	"context"
	"net/http"
	"time"

	"health-check-service/config"
	"health-check-service/internal/domain"
	"health-check-service/internal/queue"
	"health-check-service/internal/repository/mongodb"

	"github.com/sirupsen/logrus"
)

// HealthChecker performs health checks on services
type HealthChecker struct {
	serviceRepo   *mongodb.ServiceRepository
	healthLogRepo *mongodb.HealthLogRepository
	httpClient    *http.Client
	cfg           *config.Config
}

// NewHealthChecker creates a new HealthChecker
func NewHealthChecker(
	serviceRepo *mongodb.ServiceRepository,
	healthLogRepo *mongodb.HealthLogRepository,
	cfg *config.Config,
) *HealthChecker {
	return &HealthChecker{
		serviceRepo:   serviceRepo,
		healthLogRepo: healthLogRepo,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.HTTPTimeout) * time.Second,
		},
		cfg: cfg,
	}
}

// CheckService performs a health check on a single service
func (h *HealthChecker) CheckService(ctx context.Context, service *domain.Service) (*domain.HealthLog, error) {
	startTime := time.Now()

	// Perform HTTP GET request
	resp, err := h.httpClient.Get(service.URL)
	responseTime := time.Since(startTime).Milliseconds()

	log := &domain.HealthLog{
		ServiceID:        service.ID,
		ServiceName:      service.Name,
		ResponseTime:     responseTime,
		NotificationSent: false,
	}

	if err != nil {
		log.Status = domain.StatusDOWN
		log.ErrorMessage = err.Error()
		logrus.Warnf("Service %s is DOWN: %v", service.Name, err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			log.Status = domain.StatusUP
			logrus.Debugf("Service %s is UP (status: %d, response: %dms)", service.Name, resp.StatusCode, responseTime)
		} else {
			log.Status = domain.StatusDOWN
			log.ErrorMessage = http.StatusText(resp.StatusCode)
			logrus.Warnf("Service %s is DOWN (status: %d)", service.Name, resp.StatusCode)
		}
	}

	// Save health log
	if err := h.healthLogRepo.Create(ctx, log); err != nil {
		logrus.Errorf("Failed to save health log: %v", err)
		return log, err
	}

	// If service is DOWN, check if we need to send notification
	if log.Status == domain.StatusDOWN {
		shouldNotify, _ := h.shouldSendNotification(ctx, service)
		if shouldNotify {
			// Enqueue notification
			payload := queue.NotificationPayload{
				ServiceID:    service.ID,
				ServiceName:  service.Name,
				Status:       log.Status,
				ErrorMessage: log.ErrorMessage,
				LogID:        log.ID.Hex(),
			}
			if err := queue.EnqueueNotification(payload); err != nil {
				logrus.Errorf("Failed to enqueue notification: %v", err)
			}
		}
	}

	return log, nil
}

// shouldSendNotification determines if a notification should be sent
func (h *HealthChecker) shouldSendNotification(ctx context.Context, service *domain.Service) (bool, error) {
	// Get the latest log for this service
	latestLog, err := h.healthLogRepo.GetLatestByServiceID(ctx, service.ID)
	if err != nil {
		// No previous log, send notification
		return true, nil
	}

	// If notification was already sent and retry interval hasn't passed, don't send
	if latestLog.NotificationSent {
		timeSinceCheck := time.Since(latestLog.CheckedAt)
		retryInterval := time.Duration(service.RetryInterval) * time.Second
		if timeSinceCheck < retryInterval {
			logrus.Debugf("Skipping notification for %s, retry interval not passed (%.0fs < %ds)",
				service.Name, timeSinceCheck.Seconds(), service.RetryInterval)
			return false, nil
		}
	}

	return true, nil
}

// CheckAllServices performs health checks on all active services
func (h *HealthChecker) CheckAllServices(ctx context.Context) error {
	services, err := h.serviceRepo.GetActive(ctx)
	if err != nil {
		return err
	}

	logrus.Infof("Checking %d active services", len(services))

	for _, service := range services {
		h.CheckService(ctx, &service)
	}

	return nil
}
