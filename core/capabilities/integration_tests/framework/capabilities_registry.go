package framework

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	registrysyncerv2 "github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer/v2"
)

type CapabilitiesRegistry struct {
	t              *testing.T
	backend        *EthBlockchain
	contract       *kcr.CapabilitiesRegistry
	addr           common.Address
	nodeOperatorID uint32
	nextDonID      int
}

func NewCapabilitiesRegistry(ctx context.Context, t *testing.T, backend *EthBlockchain) *CapabilitiesRegistry {
	// Tests routinely run DONs with a single node and F=0, which the registry rejects unless explicitly allowed.
	addr, _, contract, err := kcr.DeployCapabilitiesRegistry(backend.transactionOpts, backend.Client(),
		kcr.CapabilitiesRegistryConstructorParams{CanAddOneNodeDONs: true})
	require.NoError(t, err)
	backend.Commit()

	_, err = contract.AddNodeOperators(backend.transactionOpts, []kcr.CapabilitiesRegistryNodeOperatorParams{
		{
			Admin: backend.transactionOpts.From,
			Name:  "TEST_NODE_OPERATOR",
		},
	})
	require.NoError(t, err)
	blockHash := backend.Commit()

	logs, err := backend.Client().FilterLogs(ctx, ethereum.FilterQuery{
		BlockHash: &blockHash,
		FromBlock: nil,
		ToBlock:   nil,
		Addresses: nil,
		Topics:    nil,
	})

	require.NoError(t, err)

	recLog, err := contract.ParseNodeOperatorAdded(logs[0])
	require.NoError(t, err)

	nopID := recLog.NodeOperatorId

	return &CapabilitiesRegistry{t: t, addr: addr, contract: contract, backend: backend, nodeOperatorID: nopID}
}

func (r *CapabilitiesRegistry) getAddress() common.Address {
	return r.addr
}

// NewRegistryCapability builds the registry entry for a capability. The v2 registry identifies
// capabilities by their "<name>@<version>" ID and carries the capability and response types as JSON
// metadata rather than as dedicated fields, so it is encoded here using the struct the registry
// syncer decodes it with.
func NewRegistryCapability(t *testing.T, capabilityID string, capabilityType registrysyncerv2.ContractCapabilityType,
	responseType uint8) kcr.CapabilitiesRegistryCapability {
	metadata, err := json.Marshal(registrysyncerv2.CapabilityMetadata{
		CapabilityType: uint8(capabilityType),
		ResponseType:   responseType,
	})
	require.NoError(t, err)

	return kcr.CapabilitiesRegistryCapability{
		CapabilityId: capabilityID,
		Metadata:     metadata,
	}
}

type capability struct {
	donCapabilityConfig *pb.CapabilityConfig
	registryConfig      kcr.CapabilitiesRegistryCapability
	// internalOnly is true if the capability is published in the registry but not made available outside the DON in which it runs
	internalOnly bool
}

// SetupDON sets up a new DON with the given capabilities and returns the DON ID
func (r *CapabilitiesRegistry) setupDON(donInfo DonConfiguration, capabilities []capability) int {
	registryCapabilities := make([]kcr.CapabilitiesRegistryCapability, 0, len(capabilities))
	capabilityIDs := make([]string, 0, len(capabilities))
	for _, c := range capabilities {
		registryCapabilities = append(registryCapabilities, c.registryConfig)
		capabilityIDs = append(capabilityIDs, c.registryConfig.CapabilityId)
	}

	_, err := r.contract.AddCapabilities(r.backend.transactionOpts, registryCapabilities)
	require.NoError(r.t, err)

	r.backend.Commit()

	peerIDs := make([][32]byte, 0, len(donInfo.p2pKeys))
	nodes := make([]kcr.CapabilitiesRegistryNodeParams, 0, len(donInfo.p2pKeys))
	for i, p2pkey := range donInfo.p2pKeys {
		signer, innerErr := getSignerStringFromOCRKeyBundle(donInfo.KeyBundles[i])
		require.NoError(r.t, innerErr)
		peer := peerIDAndOCRSigner{PeerID: p2ptypes.PeerID(p2pkey.PeerID()), Signer: signer}
		peerIDs = append(peerIDs, p2pkey.PeerID())
		n, innerErr := peerToNode(r.nodeOperatorID, peer)
		require.NoError(r.t, innerErr)

		n.CapabilityIds = capabilityIDs
		nodes = append(nodes, n)
	}

	_, err = r.contract.AddNodes(r.backend.transactionOpts, nodes)
	require.NoError(r.t, err)
	r.backend.Commit()

	capabilityConfigurations := make([]kcr.CapabilitiesRegistryCapabilityConfiguration, 0, len(capabilities))
	for _, c := range capabilities {
		configBinary, err2 := proto.Marshal(c.donCapabilityConfig)
		require.NoError(r.t, err2)

		capabilityConfigurations = append(capabilityConfigurations, kcr.CapabilitiesRegistryCapabilityConfiguration{
			CapabilityId: c.registryConfig.CapabilityId,
			Config:       configBinary,
		})
	}

	_, err = r.contract.AddDONs(r.backend.transactionOpts, []kcr.CapabilitiesRegistryNewDONParams{
		{
			Name:                     donInfo.name,
			DonFamilies:              []string{donInfo.name},
			CapabilityConfigurations: capabilityConfigurations,
			Nodes:                    peerIDs,
			F:                        donInfo.F,
			IsPublic:                 true,
			AcceptsWorkflows:         donInfo.AcceptsWorkflows,
		},
	})
	require.NoError(r.t, err)
	r.backend.Commit()

	r.nextDonID++
	return r.nextDonID
}
