package main

import (
	"log"
	"os"

	"online-music-library/handlers"
	"online-music-library/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем конфигурацию из .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Инициализация базы данных
	models.InitDB()

	// Создаем роутер
	router := gin.Default()
	
	// Регистрируем маршруты
	handlers.RegisterRoutes(router)

	// Запускаем сервер
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // Значение по умолчанию
	}
	log.Printf("Starting server on port %s...", port)
	router.Run(":" + port)
}
