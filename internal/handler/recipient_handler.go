package handler

import (
	"net/http"
	"time"

	"health-check-service/internal/domain"
	"health-check-service/internal/repository/mongodb"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecipientHandler struct {
	repo *mongodb.RecipientRepository
}

func NewRecipientHandler(repo *mongodb.RecipientRepository) *RecipientHandler {
	return &RecipientHandler{
		repo: repo,
	}
}

// GetAll returns all recipients
// @Summary Get all recipients
// @Description Get all notification recipients
// @Tags recipients
// @Produce json
// @Success 200 {array} domain.Recipient
// @Router /recipients [get]
func (h *RecipientHandler) GetAll(c *gin.Context) {
	recipients, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipients)
}

// GetByID returns a recipient by ID
// @Summary Get recipient by ID
// @Description Get a notification recipient by ID
// @Tags recipients
// @Produce json
// @Param id path string true "Recipient ID"
// @Success 200 {object} domain.Recipient
// @Router /recipients/{id} [get]
func (h *RecipientHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	recipient, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipient not found"})
		return
	}
	c.JSON(http.StatusOK, recipient)
}

// Create creates a new recipient
// @Summary Create a new recipient
// @Description Create a new notification recipient
// @Tags recipients
// @Accept json
// @Produce json
// @Param recipient body domain.CreateRecipientRequest true "Recipient data"
// @Success 201 {object} domain.Recipient
// @Router /recipients [post]
func (h *RecipientHandler) Create(c *gin.Context) {
	var req domain.CreateRecipientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipient := &domain.Recipient{
		Name:       req.Name,
		Phone:      req.Phone,
		Channel:    req.Channel,
		TelegramID: req.TelegramID,
		DiscordID:  req.DiscordID,
		IsActive:   true,
		CreatedAt:  time.Now(),
	}

	logrus.Info("Recipient created: ", recipient)

	if err := h.repo.Create(c.Request.Context(), recipient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, recipient)
}

// Update updates a recipient
// @Summary Update a recipient
// @Description Update an existing notification recipient
// @Tags recipients
// @Accept json
// @Produce json
// @Param id path string true "Recipient ID"
// @Param recipient body domain.UpdateRecipientRequest true "Recipient data"
// @Success 200 {object} map[string]string
// @Router /recipients/{id} [put]
func (h *RecipientHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req domain.UpdateRecipientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Update(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipient updated"})
}

// Delete deletes a recipient
// @Summary Delete a recipient
// @Description Delete a notification recipient
// @Tags recipients
// @Produce json
// @Param id path string true "Recipient ID"
// @Success 200 {object} map[string]string
// @Router /recipients/{id} [delete]
func (h *RecipientHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipient deleted"})
}
