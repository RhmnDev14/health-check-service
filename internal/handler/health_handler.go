package handler

import (
	"net/http"
	"strconv"

	"health-check-service/internal/repository/mongodb"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	healthLogRepo *mongodb.HealthLogRepository
	serviceRepo   *mongodb.ServiceRepository
}

func NewHealthHandler(healthLogRepo *mongodb.HealthLogRepository, serviceRepo *mongodb.ServiceRepository) *HealthHandler {
	return &HealthHandler{
		healthLogRepo: healthLogRepo,
		serviceRepo:   serviceRepo,
	}
}

// GetRecentLogs returns recent health check logs
func (h *HealthHandler) GetRecentLogs(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 50
	}

	logs, err := h.healthLogRepo.GetRecent(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// GetLogsByServiceID returns logs for a specific service
func (h *HealthHandler) GetLogsByServiceID(c *gin.Context) {
	serviceID := c.Param("serviceId")
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 50
	}

	logs, err := h.healthLogRepo.GetByServiceID(c.Request.Context(), serviceID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// GetDashboardStats returns dashboard statistics
func (h *HealthHandler) GetDashboardStats(c *gin.Context) {
	ctx := c.Request.Context()

	// Get all services
	services, err := h.serviceRepo.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stats := make([]gin.H, 0)
	upCount := 0
	downCount := 0

	for _, service := range services {
		latestLog, err := h.healthLogRepo.GetLatestByServiceID(ctx, service.ID)
		status := "UNKNOWN"
		var responseTime int64 = 0
		var errorMsg string

		if err == nil && latestLog != nil {
			status = latestLog.Status
			responseTime = latestLog.ResponseTime
			errorMsg = latestLog.ErrorMessage

			if status == "UP" {
				upCount++
			} else if status == "DOWN" {
				downCount++
			}
		}

		stats = append(stats, gin.H{
			"service_id":    service.ID,
			"service_name":  service.Name,
			"url":           service.URL,
			"status":        status,
			"response_time": responseTime,
			"error_message": errorMsg,
			"is_active":     service.IsActive,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total_services": len(services),
		"up_count":       upCount,
		"down_count":     downCount,
		"services":       stats,
	})
}
