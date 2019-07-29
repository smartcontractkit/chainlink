package web

import (
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services"
)

// PingController has the ping endpoint.
type PingController struct {
	App services.Application
}

// Show returns pong.
func (eic *PingController) Show(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
