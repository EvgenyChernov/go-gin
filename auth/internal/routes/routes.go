package routes

import (
	"auth/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter настраивает все маршруты приложения
// Принимает обработчик, который содержит логику для работы с пользователями
// Возвращает настроенный роутер
func SetupRouter(h *handler.Handler) *gin.Engine {
	router := gin.Default()

	// Группа маршрутов для аутентификации
	auth := router.Group("/auth")
	{
		// Публичные endpoints
		auth.POST("/register", h.RegisterUser)
		auth.POST("/login", h.LoginUser)
		// Защищенные endpoints
		auth.GET("/user", h.GetUserInfo)
		auth.PUT("/user", h.UpdateUser)
		auth.DELETE("/user", h.DeleteUser)
	}

	return router
}
