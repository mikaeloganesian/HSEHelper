package main

import (
	"log"

	"file-storing/config"
	"file-storing/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	defer config.CloseDB()

	r := gin.Default()

	r.POST("/upload", handlers.UploadFile)
	r.GET("/files/:id", handlers.GetFileByID)
	r.GET("/files", handlers.ListFiles)

	log.Println("file-storing service is running on port 8082")
	r.Run(":8082")
}
