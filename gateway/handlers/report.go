package handlers

import (
	"gateway/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /reports/:id
func GetReport(c *gin.Context) {
	id := c.Param("id")

	report, err := services.GetReportByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch report" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}
