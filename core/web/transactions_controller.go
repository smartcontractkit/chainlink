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
	txs, count, err := tc.App.GetStore().EthTransactionsWithOrderedAttempts(offset, size)
	ptxs := make([]presenters.EthTx, len(txs))
	for i, tx := range txs {
		if len(tx.EthTxAttempts) > 0 {
			lastAttempt := len(tx.EthTxAttempts) - 1
			ptxs[i] = presenters.NewEthTxWithAttempt(&tx, &tx.EthTxAttempts[lastAttempt])
		} else {
			ptxs[i] = presenters.NewEthTx(&tx)
		}
	}
	paginatedResponse(c, "eth_transactions", size, page, ptxs, count, err)
}

// Show returns the details of a Ethereum Transasction details.
// Example:
//  "<application>/transactions/:TxHash"
func (tc *TransactionsController) Show(c *gin.Context) {
	hash := common.HexToHash(c.Param("TxHash"))

	ethTxAttempt, err := tc.App.GetStore().FindEthTxAttempt(hash)
	if errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("Transaction not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewEthTxWithAttempt(&ethTxAttempt.EthTx, ethTxAttempt), "transaction")
}
