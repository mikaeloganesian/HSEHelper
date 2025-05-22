// @title HSE Helper Gateway API
// @version 1.0
// @description Gateway для взаимодействия с file-storing и file-analysis
// @host localhost:8080
// @BasePath /

// @schemes http
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"gateway/config"
	"gateway/handlers"
)

func main() {

	// Загрузка .env
	config.LoadEnv()

	// Инициализация роутера
	router := gin.Default()

	// Роуты
	router.POST("/upload", handlers.UploadAndAnalyze)
	router.GET("/files", handlers.ListFiles)
	router.GET("/files/:id", handlers.GetFile)
	router.POST("/analyze", handlers.AnalyzeFile)
	router.GET("/reports/:id", handlers.GetReport)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Gateway listening on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
