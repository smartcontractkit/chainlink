package web

import (
	"chainlink/core/services"

	"github.com/gin-gonic/gin"
)

// TxAttemptsController lists TxAttempts requests.
type TxAttemptsController struct {
	App services.Application
}

// Index returns paginated transaction attempts
func (tac *TxAttemptsController) Index(c *gin.Context, size, page, offset int) {
	tas, count, err := tac.App.GetStore().TxAttempts(offset, size)
	paginatedResponse(c, "TxAttempts", size, page, tas, count, err)
}
