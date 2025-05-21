package handlers

import (
	"fmt"
	"gateway/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AnalyzeFile(c *gin.Context) {
	var req services.AnalyzeRequest
	fmt.Print(req)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body" + err.Error() + req.Text + req.FileName})
		return
	}

	result, err := services.AnalyzeText(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze text" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
