package operation

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	aptosstate "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/aptos"
)

type AptosDeps struct {
	AB               *cldf.AddressBookMap
	AptosChain       cldf.AptosChain
	CCIPOnChainState stateview.CCIPOnChainState
	OnChainState     aptosstate.CCIPChainState
}
