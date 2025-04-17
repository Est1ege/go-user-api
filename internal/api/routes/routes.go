package routes

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/Est1ege/go-user-api/internal/api/handlers"
	"github.com/Est1ege/go-user-api/internal/api/middleware"
)

// SetupRouter настраивает маршруты API и веб-интерфейса
func SetupRouter(userHandler *handlers.UserHandler, webHandler *handlers.WebHandler) *gin.Engine {
	router := gin.Default()
	
	// Добавляем middleware для логирования
	router.Use(middleware.Logger())
	
	// Настройка сессий
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("user-api-session", store))

	// Обработка флеш-сообщений
	router.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		if flashes := session.Flashes(); len(flashes) > 0 {
			c.Set("success", flashes[0])
			session.Save()
		}
		c.Next()
	})
	
	// Загрузка шаблонов
	router.LoadHTMLGlob("templates/**/*")
	
	// API v1
	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", userHandler.Create)
			users.GET("/:id", userHandler.GetByID)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}
	}
	
	// Веб-интерфейс
	web := router.Group("/web")
	{
		users := web.Group("/users")
		{
			users.GET("", webHandler.Index)
			users.POST("", webHandler.Create)
			users.POST("/:id", webHandler.Update)
			users.POST("/:id/delete", webHandler.Delete)
		}
	}
	
	// Редирект с корня на веб-интерфейс
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/web/users")
	})
	
	return router
}