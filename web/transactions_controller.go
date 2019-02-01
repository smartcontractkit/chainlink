package web

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/orm"
	"github.com/smartcontractkit/chainlink/store/presenters"
)

// TransactionsController displays Ethereum transactions requests.
type TransactionsController struct {
	App services.Application
}

// Index returns paginated transaction attempts
func (tc *TransactionsController) Index(c *gin.Context) {
	size, page, offset, err := ParsePaginatedRequest(c.Query("size"), c.Query("page"))
	if err != nil {
		c.AbortWithError(422, err)
		return
	}

	txs, count, err := tc.App.GetStore().Transactions(offset, size)
	if err == orm.ErrorNotFound {
		c.Data(404, MediaType, emptyJSON)
	} else if err != nil {
		c.AbortWithError(500, fmt.Errorf("error getting paged Transactions: %+v", err))
	} else if buffer, err := NewPaginatedResponse(*c.Request.URL, size, page, count, txs); err != nil {
		c.AbortWithError(500, err)
	} else {
		c.Data(200, MediaType, buffer)
	}
}

// Show returns the details of a Ethereum Transasction details.
// Example:
//  "<application>/transactions/:TxHash"
func (tc *TransactionsController) Show(c *gin.Context) {
	hash := common.HexToHash(c.Param("TxHash"))
	if tx, err := tc.App.GetStore().FindTxByAttempt(hash); err == orm.ErrorNotFound {
		c.AbortWithError(404, errors.New("Transaction not found"))
	} else if err != nil {
		c.AbortWithError(500, err)
	} else if doc, err := jsonapi.Marshal(presenters.NewTx(tx)); err != nil {
		c.AbortWithError(500, err)
	} else {
		c.Data(200, MediaType, doc)
	}
}
