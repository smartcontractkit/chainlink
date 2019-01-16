package web

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/orm"
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

	attempts, count, err := tac.App.GetStore().TxAttempts(offset, size)
	if err == orm.ErrorNotFound {
		c.Data(404, MediaType, emptyJSON)
	} else if err != nil {
		c.AbortWithError(500, fmt.Errorf("error getting paged TxAttempts: %+v", err))
	} else if buffer, err := NewPaginatedResponse(*c.Request.URL, size, page, count, attempts); err != nil {
		c.AbortWithError(500, err)
	} else {
		c.Data(200, MediaType, buffer)
	}
}

// Show returns the details of a Transasction Attempt.
// Example:
//  "<application>/txattempts/:TxHash"
func (tac *TxAttemptsController) Show(c *gin.Context) {
	hash := common.HexToHash(c.Param("TxHash"))
	if tx, err := tac.App.GetStore().FindFullTxAttempt(hash); err == orm.ErrorNotFound {
		c.AbortWithError(404, errors.New("Transaction Attempt not found"))
	} else if err != nil {
		c.AbortWithError(500, err)
	} else if doc, err := jsonapi.Marshal(tx); err != nil {
		c.AbortWithError(500, err)
	} else {
		c.Data(200, MediaType, doc)
	}
}
