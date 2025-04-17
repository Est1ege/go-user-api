package main

import (
	"log"

	"github.com/Est1ege/go-user-api/internal/api/handlers"
	"github.com/Est1ege/go-user-api/internal/api/routes"
	"github.com/Est1ege/go-user-api/internal/config"
	"github.com/Est1ege/go-user-api/internal/repository/postgres"
	"github.com/Est1ege/go-user-api/internal/service"
	"github.com/Est1ege/go-user-api/pkg/database"
	"github.com/Est1ege/go-user-api/pkg/validator"
)

func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Настройка валидатора
	validator.SetupValidator()

	// Подключение к базе данных
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err.Error())
	}

	// Инициализация репозиториев
	userRepo := postgres.NewUserRepository(db)

	// Инициализация сервисов
	userService := service.NewUserService(userRepo)

	// Инициализация обработчиков
	userHandler := handlers.NewUserHandler(userService)
	webHandler := handlers.NewWebHandler(userService)

	// Настройка маршрутов
	router := routes.SetupRouter(userHandler, webHandler)

	// Запуск сервера
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Printf("Web interface available at http://localhost:%s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}
}