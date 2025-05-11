package handlers

import (
	"net/http"

	"github.com/QuorumMind/transaction-service/internal/models"
	"github.com/QuorumMind/transaction-service/internal/service"
	"github.com/gin-gonic/gin"
)

func TransferHandler(c *gin.Context) {
	var req models.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := service.ProcessTransfer(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "transfer accepted"})
}
