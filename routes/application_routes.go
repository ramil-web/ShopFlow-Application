package routes

import (
	"shopflow/application/handlers"
	"shopflow/application/middleware"
	"shopflow/application/services"

	"github.com/gin-gonic/gin"
)

// RegisterApplicationRoutes регистрирует маршруты для Application сервиса
func RegisterApplicationRoutes(r *gin.Engine, appSvc *services.ApplicationService, publisher *services.NotificationService) {
	api := r.Group("/api")
	{
		appGroup := api.Group("/applications")
		appGroup.Use(middleware.AuthMiddleware()) // JWT или любая другая авторизация

		// создаём один экземпляр хендлера с DI
		h := &handlers.ApplicationHandler{
			BaseHandler: &handlers.BaseHandler{
				AppSvc:    appSvc,
				Publisher: publisher,
			},
		}

		// маршруты
		appGroup.POST("", h.CreateApplication)       // создание заявки (с gRPC Auth проверкой)
		appGroup.GET("", h.GetApplications)          // получить все заявки текущего пользователя
		appGroup.GET("/:id", h.GetApplicationById)   // получить заявку по ID
		appGroup.DELETE("/:id", h.DeleteApplication) // удалить заявку
		appGroup.PATCH("/:id", h.UpdateApplication)  // обновить заявку
	}
}
