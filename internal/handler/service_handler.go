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
// @Summary Get all services
// @Description Get all monitored services
// @Tags services
// @Produce json
// @Success 200 {array} domain.Service
// @Router /services [get]
func (h *ServiceHandler) GetAll(c *gin.Context) {
	services, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}

// GetByID returns a service by ID
// @Summary Get service by ID
// @Description Get a service by its ID
// @Tags services
// @Produce json
// @Param id path string true "Service ID"
// @Success 200 {object} domain.Service
// @Router /services/{id} [get]
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
// @Summary Create a new service
// @Description Create a new monitored service
// @Tags services
// @Accept json
// @Produce json
// @Param service body domain.CreateServiceRequest true "Service data"
// @Success 201 {object} domain.Service
// @Router /services [post]
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
// @Summary Update a service
// @Description Update an existing service
// @Tags services
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Param service body domain.UpdateServiceRequest true "Service data"
// @Success 200 {object} map[string]string
// @Router /services/{id} [put]
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
// @Summary Delete a service
// @Description Delete a monitored service
// @Tags services
// @Produce json
// @Param id path string true "Service ID"
// @Success 200 {object} map[string]string
// @Router /services/{id} [delete]
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
