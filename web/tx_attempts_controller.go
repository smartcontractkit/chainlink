package web

import (
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
)

// TxAttemptsController lists TxAttempts requests.
type TxAttemptsController struct {
	App services.Application
}

// Index returns paginated transaction attempts
func (tac *TxAttemptsController) Index(c *gin.Context) {
	size, page, offset, err := ParsePaginatedRequest(c.Query("size"), c.Query("page"))
	if err != nil {
		c.AbortWithError(422, err)
		return
	}

	tas, count, err := tac.App.GetStore().TxAttempts(offset, size)
	paginatedResponse(c, "TxAttempts", size, page, tas, count, err)
}
