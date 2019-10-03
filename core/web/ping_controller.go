package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"chainlink/core/services"
)

// PingController has the ping endpoint.
type PingController struct {
	App services.Application
}

// Show returns pong.
func (eic *PingController) Show(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
