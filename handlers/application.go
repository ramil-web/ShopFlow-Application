package handlers

import (
	"log"
	"net/http"
	"shopflow/application/models"
	_ "shopflow/application/publisher"
	"shopflow/application/services"

	"github.com/gin-gonic/gin"
)

// CreateApplicationHandler godoc
// @Summary Create User Application
// @Security BearerAuth
// @Tags UserApplication
// @Accept json
// @Produce json
// @Param input body models.CreateApplicationRequest true "Create Application"
// @Success 201 {object} models.Application "Application successfully created"
// @Failure 400 {object} map[string]string
// @Router /api/applications/store [post]
func CreateApplicationHandler(appSvc *services.ApplicationService, publisher *services.NotificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateApplicationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID := c.GetUint("user_id")
		email := c.GetString("email")

		app, err := appSvc.CreateApplication(userID, req.Text, req.FileURL, req.Status)
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
			if err := publisher.PublishApplicationCreated(m); err != nil {
				log.Println("[ERROR] failed to publish application_created event:", err)
			}
		}(msg)

		c.JSON(http.StatusCreated, app)
	}
}

// GetApplicationHandler  godoc
// @Summary Gets all the applications
// @Security BearerAuth
// @Tags UserApplication
// @Accept json
// @Produce json
// @Param user_id query int false "Filter by user ID"
// @Success 200 {array} models.Application
// @Failure 400  {object} map[string]string
// @Router /api/applications [get]
func GetApplicationHandler(appSvc *services.ApplicationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		apps, err := appSvc.GetAll(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, apps)
	}
}
