package handler

import (
	"net/http"
	"time"

	"health-check-service/config"
	"health-check-service/internal/domain"
	"health-check-service/internal/repository/mongodb"
	"health-check-service/internal/scheduler"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	repo      *mongodb.ServiceRepository
	scheduler *scheduler.Scheduler
	cfg       *config.Config
}

func NewServiceHandler(repo *mongodb.ServiceRepository, scheduler *scheduler.Scheduler, cfg *config.Config) *ServiceHandler {
	return &ServiceHandler{
		repo:      repo,
		scheduler: scheduler,
		cfg:       cfg,
	}
}

// GetAll returns all services
func (h *ServiceHandler) GetAll(c *gin.Context) {
	services, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}

// GetByID returns a service by ID
func (h *ServiceHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	service, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}
	c.JSON(http.StatusOK, service)
}

// Create creates a new service
func (h *ServiceHandler) Create(c *gin.Context) {
	var req domain.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default intervals if not provided
	if req.CheckInterval == 0 {
		req.CheckInterval = h.cfg.DefaultCheckInterval
	}
	if req.RetryInterval == 0 {
		req.RetryInterval = h.cfg.DefaultRetryInterval
	}

	service := &domain.Service{
		ID:            req.ID,
		Name:          req.Name,
		URL:           req.URL,
		CheckInterval: req.CheckInterval,
		RetryInterval: req.RetryInterval,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := h.repo.Create(c.Request.Context(), service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Schedule the new service
	h.scheduler.ScheduleService(service)

	c.JSON(http.StatusCreated, service)
}

// Update updates a service
func (h *ServiceHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req domain.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Update(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Re-fetch and re-schedule
	service, err := h.repo.GetByID(c.Request.Context(), id)
	if err == nil {
		h.scheduler.ScheduleService(service)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service updated"})
}

// Delete deletes a service
func (h *ServiceHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	// Unschedule first
	h.scheduler.UnscheduleService(id)

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted"})
}
