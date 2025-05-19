package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"file-analysis/models"
	"file-analysis/repository"
	"file-analysis/services"
)

func AnalyzeHandler(c *gin.Context) {
	var req models.AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileName := req.FileName
	if fileName == "" {
		fileName = "unknown.txt"
	}

	hash := services.CalculateHash(req.Text)
	existingReport, err := repository.FindReportByHash(hash)
	if err == nil && existingReport != nil {
		c.JSON(http.StatusOK, models.AnalyzeResponse{
			Paragraphs:    existingReport.Paragraphs,
			Words:         existingReport.Words,
			Characters:    existingReport.Characters,
			IsPlagiarized: true,
		})
		return
	}

	result := services.AnalyzeText(req.Text)
	report := &models.Report{
		FileName:   fileName,
		Paragraphs: result.Paragraphs,
		Words:      result.Words,
		Characters: result.Characters,
		Hash:       hash,
	}

	err = repository.InsertReport(report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result.IsPlagiarized = false
	c.JSON(http.StatusOK, result)
}

func GetReportByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	report, err := repository.FindReportByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "report not found"})
		return
	}

	c.JSON(http.StatusOK, report)
}
