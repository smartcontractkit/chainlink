package operation

import (
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

type AptosDeps struct {
	AB               *cldf.AddressBookMap
	AptosChain       cldf_aptos.Chain
	CCIPOnChainState stateview.CCIPOnChainState
}
