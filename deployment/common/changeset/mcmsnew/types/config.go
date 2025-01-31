package mcmsnew

import (
	"math/big"

	mcmsTypes "github.com/smartcontractkit/mcms/types"
)

// MCMSWithTimelockConfig holds the configuration for an MCMS with timelock.
// Note that this type already exists in types.go, but this one is using the new lib version.
type MCMSWithTimelockConfig struct {
	Canceller        mcmsTypes.Config
	Bypasser         mcmsTypes.Config
	Proposer         mcmsTypes.Config
	TimelockMinDelay *big.Int
}
