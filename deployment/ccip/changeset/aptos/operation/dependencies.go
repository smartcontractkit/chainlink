package operation

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

type AptosDeps struct {
	AB         *cldf.AddressBookMap
	AptosChain deployment.AptosChain
	// TODO: Refactor this?
	OnChainState     changeset.AptosCCIPChainState
	CCIPOnChainState changeset.CCIPOnChainState
}
