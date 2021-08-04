package synchronization

import (
	"encoding/json"

	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	null "gopkg.in/guregu/null.v4"
)

func formatEthereumReceipt(str string) (*syncReceiptPresenter, error) {
	var receipt types.Receipt
	err := json.Unmarshal([]byte(str), &receipt)
	if err != nil {
		return nil, err
	}
	return &syncReceiptPresenter{
		Hash:   receipt.TxHash,
		Status: runLogStatusPresenter(receipt),
	}, nil
}

type syncReceiptPresenter struct {
	Hash   common.Hash `json:"transactionHash"`
	Status TxStatus    `json:"transactionStatus"`
}

// TxStatus indicates if a transaction is fulfilled or not
type TxStatus string

const (
	// StatusFulfilledRunLog indicates that a ChainlinkFulfilled event was
	// detected in the transaction receipt.
	StatusFulfilledRunLog TxStatus = "fulfilledRunLog"
	// StatusNoFulfilledRunLog indicates that no ChainlinkFulfilled events were
	// detected in the transaction receipt.
	StatusNoFulfilledRunLog = "noFulfilledRunLog"
)

func runLogStatusPresenter(receipt types.Receipt) TxStatus {
	if models.ReceiptIndicatesRunLogFulfillment(receipt) {
		return StatusFulfilledRunLog
	}
	return StatusNoFulfilledRunLog
}

type syncInitiatorPresenter struct {
	Type      string               `json:"type"`
	RequestID *common.Hash         `json:"requestId,omitempty"`
	TxHash    *common.Hash         `json:"txHash,omitempty"`
	Requester *ethkey.EIP55Address `json:"requester,omitempty"`
}

type syncTaskRunPresenter struct {
	Index                            int           `json:"index"`
	Type                             string        `json:"type"`
	Status                           string        `json:"status"`
	Error                            null.String   `json:"error"`
	Result                           interface{}   `json:"result,omitempty"`
	ObservedIncomingConfirmations    clnull.Uint32 `json:"confirmations"`
	MinRequiredIncomingConfirmations clnull.Uint32 `json:"minimumConfirmations"`
}
