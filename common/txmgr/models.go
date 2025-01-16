package txmgr

import (
	"github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
)

const (
	TxUnstarted               = types.TxState("unstarted")
	TxInProgress              = types.TxState("in_progress")
	TxFatalError              = types.TxState("fatal_error")
	TxUnconfirmed             = types.TxState("unconfirmed")
	TxConfirmed               = types.TxState("confirmed")
	TxConfirmedMissingReceipt = types.TxState("confirmed_missing_receipt")
	TxFinalized               = types.TxState("finalized")
)
