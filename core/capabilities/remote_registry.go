package capabilities

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type Capability struct {
	ID [32]byte
	// The `Name` is a partially qualified ID for the capability.
	// Validation: ^[a-z0-9_\-:]{1,32}$
	Name string
	// Semver, e.g., "1.2.3"
	Version      string
	ResponseType CapabilityResponseType
	// An address to the capability configuration contract. Having this defined
	// on a capability enforces consistent configuration across DON instances
	// serving the same capability.
	//
	// The main use cases are:
	// 1) Sharing capability configuration across DON instances
	// 2) Inspect and modify on-chain configuration without off-chain
	// capability code.
	ConfigurationContract common.Address
}

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
