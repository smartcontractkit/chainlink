package web

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"chainlink/core/services"
	"chainlink/core/store/orm"
	"chainlink/core/store/presenters"
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
		txp := presenters.NewTx(&tx)
		ptxs[i] = txp
	}
	paginatedResponse(c, "Transactions", size, page, ptxs, count, err)
}

// Show returns the details of a Ethereum Transasction details.
// Example:
//  "<application>/transactions/:TxHash"
func (tc *TransactionsController) Show(c *gin.Context) {
	hash := common.HexToHash(c.Param("TxHash"))
	if txAttempt, err := tc.App.GetStore().FindTxAttempt(hash); errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("Transaction not found"))
	} else if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponse(c, presenters.NewTxFromAttempt(*txAttempt), "transaction")
	}
}
