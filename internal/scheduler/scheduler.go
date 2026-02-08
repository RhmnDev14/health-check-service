package scheduler

import (
	"context"
	"sync"
	"time"

	"health-check-service/internal/domain"
	"health-check-service/internal/usecase"

	"github.com/sirupsen/logrus"
)

// Scheduler manages health check scheduling for all services
type Scheduler struct {
	healthChecker *usecase.HealthChecker
	tickers       map[string]*time.Ticker
	stopChan      map[string]chan bool
	mu            sync.RWMutex
	running       bool
}

// NewScheduler creates a new Scheduler
func NewScheduler(healthChecker *usecase.HealthChecker) *Scheduler {
	return &Scheduler{
		healthChecker: healthChecker,
		tickers:       make(map[string]*time.Ticker),
		stopChan:      make(map[string]chan bool),
	}
}

// Start starts the scheduler
func (s *Scheduler) Start(ctx context.Context) error {
	s.mu.Lock()
	s.running = true
	s.mu.Unlock()

	logrus.Info("Scheduler started")

	// Initial check for all services
	go s.healthChecker.CheckAllServices(ctx)

	return nil
}

// ScheduleService starts a ticker for a specific service
func (s *Scheduler) ScheduleService(service *domain.Service) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Stop existing ticker if any
	if stop, exists := s.stopChan[service.ID]; exists {
		close(stop)
		delete(s.stopChan, service.ID)
		delete(s.tickers, service.ID)
	}

	if !service.IsActive {
		logrus.Debugf("Service %s is not active, skipping scheduling", service.Name)
		return
	}

	interval := time.Duration(service.CheckInterval) * time.Second
	ticker := time.NewTicker(interval)
	stop := make(chan bool)

	s.tickers[service.ID] = ticker
	s.stopChan[service.ID] = stop

	go func(svc domain.Service) {
		logrus.Infof("Started health check ticker for %s (interval: %ds)", svc.Name, svc.CheckInterval)
		for {
			select {
			case <-ticker.C:
				ctx := context.Background()
				s.healthChecker.CheckService(ctx, &svc)
			case <-stop:
				ticker.Stop()
				logrus.Infof("Stopped health check ticker for %s", svc.Name)
				return
			}
		}
	}(*service)
}

// UnscheduleService stops the ticker for a specific service
func (s *Scheduler) UnscheduleService(serviceID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if stop, exists := s.stopChan[serviceID]; exists {
		close(stop)
		delete(s.stopChan, serviceID)
		delete(s.tickers, serviceID)
	}
}

// Stop stops all tickers
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, stop := range s.stopChan {
		close(stop)
		delete(s.stopChan, id)
		delete(s.tickers, id)
	}

	s.running = false
	logrus.Info("Scheduler stopped")
}

// IsRunning returns whether the scheduler is running
func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}
