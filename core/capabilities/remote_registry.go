package capabilities

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// remoteRegistry contains a local cache of the CapabilityRegistry deployed
// on-chain. It is updated by the registrySyncer and is otherwise read-only.
type remoteRegistry struct {
	address      common.Address
	capabilities []Capability
	lggr         logger.Logger
}

// NewRemoteRegistry creates a new remote capability registry
func NewRemoteRegistry(remoteRegistryAddress string, lggr logger.Logger) *remoteRegistry {
	onchainCapabilityRegistryAddress, err := evmtypes.NewEIP55Address(remoteRegistryAddress)
	if err != nil {
		panic(fmt.Sprintf("failed to remote capability registry address. Received address: %v. Err: %v", remoteRegistryAddress, err))
	}

	return &remoteRegistry{
		address: onchainCapabilityRegistryAddress.Address(),
		lggr:    lggr,
	}
}
