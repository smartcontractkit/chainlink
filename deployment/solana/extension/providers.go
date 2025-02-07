package deployment_solana

import (
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
)

// Providers are custom implementations of Solana dependencies.

// No need to reexport this. If signature with Dependency doesn't match, we would implement the provider here
var SendAndConfirm = solCommonUtil.SendAndConfirm
