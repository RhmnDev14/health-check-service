package handler

import (
	"net/http"
	"time"

	"health-check-service/internal/domain"
	"health-check-service/internal/repository/mongodb"

	"github.com/gin-gonic/gin"
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
func (h *RecipientHandler) GetAll(c *gin.Context) {
	recipients, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipients)
}

// GetByID returns a recipient by ID
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
func (h *RecipientHandler) Create(c *gin.Context) {
	var req domain.CreateRecipientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipient := &domain.Recipient{
		Name:      req.Name,
		Phone:     req.Phone,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	if err := h.repo.Create(c.Request.Context(), recipient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, recipient)
}

// Update updates a recipient
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
