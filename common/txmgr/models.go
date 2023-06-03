package txmgr

import (
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
)

const (
	// TODO: change Eth prefix: https://smartcontract-it.atlassian.net/browse/BCI-1198
	EthTxUnstarted               = txmgrtypes.TxState("unstarted")
	EthTxInProgress              = txmgrtypes.TxState("in_progress")
	EthTxFatalError              = txmgrtypes.TxState("fatal_error")
	EthTxUnconfirmed             = txmgrtypes.TxState("unconfirmed")
	EthTxConfirmed               = txmgrtypes.TxState("confirmed")
	EthTxConfirmedMissingReceipt = txmgrtypes.TxState("confirmed_missing_receipt")
)
