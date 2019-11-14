package web

import (
	"net/http"

	"chainlink/core/services"

	"github.com/gin-gonic/gin"
)

// PingController has the ping endpoint.
type PingController struct {
	App services.Application
}

// Show returns pong.
func (eic *PingController) Show(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
