package handlers

import (
	"gateway/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetReport godoc
// @Summary      Get analysis report by ID
// @Description  Fetch an analysis report from file-analysis by its ID
// @Tags         Report Get
// @Accept       json
// @Produce      json
// @Param        id path string true "Report ID"
// @Success      200 {object} services.Report "Report fetched successfully"
// @Failure      400 {object} gin.H "Bad request"
// @Failure      500 {object} gin.H "Internal server error"
// @Router       /reports/{id} [get]
func GetReport(c *gin.Context) {
	id := c.Param("id")

	report, err := services.GetReportByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch report" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}
