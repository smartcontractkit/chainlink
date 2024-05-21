package txmgr

import (
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
)

const (
	TxUnstarted               = txmgrtypes.TxState("unstarted")
	TxInProgress              = txmgrtypes.TxState("in_progress")
	TxFatalError              = txmgrtypes.TxState("fatal_error")
	TxUnconfirmed             = txmgrtypes.TxState("unconfirmed")
	TxConfirmed               = txmgrtypes.TxState("confirmed")
	TxConfirmedMissingReceipt = txmgrtypes.TxState("confirmed_missing_receipt")
)
