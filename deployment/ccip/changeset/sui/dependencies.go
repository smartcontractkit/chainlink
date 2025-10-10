package sui

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

type Deps struct {
	AB               *cldf.AddressBookMap
	SuiChain         sui_ops.OpTxDeps
	CCIPOnChainState stateview.CCIPOnChainState
}
