package routes

import (
	"shopflow/application/handlers"
	"shopflow/application/middleware"
	"shopflow/application/services"

	"github.com/gin-gonic/gin"
)

func RegisterApplicationRoutes(r *gin.Engine, appSvc *services.ApplicationService, publisher *services.NotificationService) {
	api := r.Group("/api")
	{
		appGroup := api.Group("/applications")
		appGroup.Use(middleware.AuthMiddleware())
		{
			appGroup.POST("/store", handlers.CreateApplicationHandler(appSvc, publisher))
			appGroup.GET("", handlers.GetApplicationHandler(appSvc))
		}
	}
}
