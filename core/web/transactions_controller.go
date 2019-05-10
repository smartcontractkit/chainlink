package web

import (
	"errors"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
)

// TransactionsController displays Ethereum transactions requests.
type TransactionsController struct {
	App services.Application
}

// Index returns paginated transaction attempts
func (tc *TransactionsController) Index(c *gin.Context, size, page, offset int) {
	txs, count, err := tc.App.GetStore().Transactions(offset, size)
	ptxs := make([]presenters.Tx, len(txs))
	for i, tx := range txs {
		ptxs[i] = presenters.NewTx(&tx)
	}
	paginatedResponse(c, "Transactions", size, page, ptxs, count, err)
}

// Show returns the details of a Ethereum Transasction details.
// Example:
//  "<application>/transactions/:TxHash"
func (tc *TransactionsController) Show(c *gin.Context) {
	hash := common.HexToHash(c.Param("TxHash"))
	if tx, err := tc.App.GetStore().FindTxByAttempt(hash); err == orm.ErrorNotFound {
		publicError(c, http.StatusNotFound, errors.New("Transaction not found"))
	} else if err != nil {
		publicError(c, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponse(c, presenters.NewTx(tx), "transaction")
	}
}
