package web

import (
	"net/http"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/store/presenters"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// TransactionsController displays Ethereum transactions requests.
type TransactionsController struct {
	App chainlink.Application
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

	txAttempt, err := tc.App.GetStore().FindTxAttempt(hash)
	if errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("Transaction not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewTxFromAttempt(*txAttempt), "transaction")
}
