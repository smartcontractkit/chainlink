package framework

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"

	"testing"

	"github.com/stretchr/testify/require"
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
	addr, _, contract, err := kcr.DeployCapabilitiesRegistry(backend.transactionOpts, backend)
	require.NoError(t, err)
	backend.Commit()

	_, err = contract.AddNodeOperators(backend.transactionOpts, []kcr.CapabilitiesRegistryNodeOperator{
		{
			Admin: backend.transactionOpts.From,
			Name:  "TEST_NODE_OPERATOR",
		},
	})
	require.NoError(t, err)
	blockHash := backend.Commit()

	logs, err := backend.FilterLogs(ctx, ethereum.FilterQuery{
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

type capability struct {
	donCapabilityConfig *pb.CapabilityConfig
	registryConfig      kcr.CapabilitiesRegistryCapability
}

// SetupDON sets up a new DON with the given capabilities and returns the DON ID
func (r *CapabilitiesRegistry) setupDON(donInfo DonConfiguration, capabilities []capability) int {
	var hashedCapabilityIDs [][32]byte

	for _, c := range capabilities {
		id, err := r.contract.GetHashedCapabilityId(&bind.CallOpts{}, c.registryConfig.LabelledName, c.registryConfig.Version)
		require.NoError(r.t, err)
		hashedCapabilityIDs = append(hashedCapabilityIDs, id)
	}

	var registryCapabilities []kcr.CapabilitiesRegistryCapability
	for _, c := range capabilities {
		registryCapabilities = append(registryCapabilities, c.registryConfig)
	}

	_, err := r.contract.AddCapabilities(r.backend.transactionOpts, registryCapabilities)
	require.NoError(r.t, err)

	r.backend.Commit()

	nodes := []kcr.CapabilitiesRegistryNodeParams{}
	for _, peerID := range donInfo.peerIDs {
		n, innerErr := peerToNode(r.nodeOperatorID, peerID)
		require.NoError(r.t, innerErr)

		n.HashedCapabilityIds = hashedCapabilityIDs
		nodes = append(nodes, n)
	}

	_, err = r.contract.AddNodes(r.backend.transactionOpts, nodes)
	require.NoError(r.t, err)
	r.backend.Commit()

	ps, err := peers(donInfo.peerIDs)
	require.NoError(r.t, err)

	var capabilityConfigurations []kcr.CapabilitiesRegistryCapabilityConfiguration
	for i, c := range capabilities {
		configBinary, err2 := proto.Marshal(c.donCapabilityConfig)
		require.NoError(r.t, err2)

		capabilityConfigurations = append(capabilityConfigurations, kcr.CapabilitiesRegistryCapabilityConfiguration{
			CapabilityId: hashedCapabilityIDs[i],
			Config:       configBinary,
		})
	}

	_, err = r.contract.AddDON(r.backend.transactionOpts, ps, capabilityConfigurations, true, donInfo.AcceptsWorkflows, donInfo.F)
	require.NoError(r.t, err)
	r.backend.Commit()

	r.nextDonID++
	return r.nextDonID
}

func newCapabilityConfig() *pb.CapabilityConfig {
	return &pb.CapabilityConfig{
		DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
	}
}
