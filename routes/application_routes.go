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

		// создаём один экземпляр хендлера
		h := &handlers.ApplicationHandler{
			BaseHandler: &handlers.BaseHandler{
				AppSvc:    appSvc,
				Publisher: publisher,
			},
		}

		// регистрируем методы структуры как обработчики
		appGroup.POST("", h.CreateApplication)
		appGroup.GET("", h.GetApplications)
		appGroup.GET("/:id", h.GetApplicationById)
		appGroup.DELETE("/:id", h.DeleteApplication)
		appGroup.PATCH("/:id", h.UpdateApplication)
	}
}
