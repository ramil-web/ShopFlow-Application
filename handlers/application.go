package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"shopflow/application/models"
	_ "shopflow/application/publisher"
	"shopflow/application/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
	AppSvc    *services.ApplicationService
	Publisher *services.NotificationService
}

type ApplicationHandler struct {
	*BaseHandler
}

// CreateApplication godoc
// @Summary Create Application
// @Security BearerAuth
// @Tags UserApplication
// @Accept json
// @Produce json
// @Param input body models.CreateApplicationRequest true "Create Application"
// @Success 201 {object} models.Application "Application successfully created"
// @Failure 400 {object} map[string]string
// @Router /api/applications [post]
func (h *ApplicationHandler) CreateApplication(c *gin.Context) {
	var req models.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	email := c.GetString("email")

	app, err := h.AppSvc.CreateApplication(userID, req.Text, req.FileURL, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Формируем сообщение для уведомления
	msg := services.ApplicationCreatedMessage{
		ID:     app.ID,
		UserID: app.UserID,
		Text:   app.Text,
		File:   app.FileURL,
		Email:  email,
	}

	// Публикуем асинхронно
	go func(m services.ApplicationCreatedMessage) {
		if err := h.Publisher.PublishApplicationCreated(m); err != nil {
			log.Println("[ERROR] failed to publish application_created event:", err)
		}
	}(msg)

	c.JSON(http.StatusCreated, app)
}

// GetApplications godoc
// @Summary Gets all the Applications
// @Security BearerAuth
// @Tags UserApplication
// @Accept json
// @Produce json
// @Param user_id query int false "Filter by user ID"
// @Success 200 {object} models.Application
// @Failure 400  {object} map[string]string
// @Router /api/applications [get]
func (h *ApplicationHandler) GetApplications(c *gin.Context) {
	userIDStr := c.Query("user_id")
	var userID uint = 0
	if userIDStr != "" {
		if uid, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			userID = uint(uid)
		}
	}

	apps, err := h.AppSvc.GetAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, apps)
}

// GetApplicationById godoc
// @Summary Gets Application by ID
// @Tags UserApplication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Application ID"
// @Success 200 {object} models.Application "Данные одной заявки"
// @Failure 400 {object} map[string]string
// @Router /api/applications/{id} [get]
func (h *ApplicationHandler) GetApplicationById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		// любой нечисловой ID возвращаем 404
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	app, err := h.AppSvc.GetApplicationById(uint(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, app)
}

// DeleteApplication godoc
// @Summary Delete Application by ID
// @Tags UserApplication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Application ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/applications/{id} [delete]
func (h *ApplicationHandler) DeleteApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		// любой нечисловой ID возвращаем 404
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	err = h.AppSvc.DeleteApplication(uint(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// UpdateApplication  godoc
// @Summary Update Application
// @Tags UserApplication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Application ID"
// @Param input body models.UpdateApplicationRequest true "Application request with user_id, text и status"
// @Success 200 {object} models.Application
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/applications/{id} [patch]
func (h *ApplicationHandler) UpdateApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
	}

	var req models.UpdateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	app, err := h.AppSvc.UpdateApplication(req, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if app == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}
	c.JSON(http.StatusOK, app)
}
