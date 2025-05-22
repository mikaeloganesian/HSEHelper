package handlers

import (
	"fmt"
	"gateway/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AnalyzeFile godoc
// @Summary      Analyze file content
// @Description  Analyze the content of a file and return the analysis result
// @Tags         File Analysis
// @Accept       json
// @Produce      json
// @Param        analyzeRequest body services.AnalyzeRequest true "Analyze request"
// @Success      200 {object} services.AnalyzeResponse "Analysis successful"
// @Failure      400 {object} gin.H "Bad request"
// @Failure      500 {object} gin.H "Internal server error"
// @Router       /analyze [post]

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
