package web

import (
	"net/http"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"

	"github.com/gin-gonic/gin"
)

// PingController has the ping endpoint.
type PingController struct {
	App chainlink.Application
}

// Show returns pong.
func (eic *PingController) Show(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
