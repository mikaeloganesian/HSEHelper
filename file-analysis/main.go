package main

import (
	"file-analysis/config"
	"file-analysis/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	defer config.CloseDB()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/analyze", handlers.AnalyzeHandler)
	r.GET("/reports/:id", handlers.GetReportByIDHandler)

	r.Run(":8081")
}
