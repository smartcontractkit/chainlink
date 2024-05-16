package capabilities_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commonMocks "github.com/smartcontractkit/chainlink-common/pkg/types/mocks"
	coreCapabilities "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	utils "github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/internal"
	remoteMocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types/mocks"
)

func TestIntegration_GetCapabilities(t *testing.T) {
	owner, simulatedBackend := utils.StartNewChain(t)

	capabilityRegistry := utils.DeployCapabilityRegistry(t, owner, simulatedBackend)

	utils.AddCapability(t, owner, simulatedBackend, capabilityRegistry, utils.DataStreamsReportCapability)
	utils.AddCapability(t, owner, simulatedBackend, capabilityRegistry, utils.WriteChainCapability)

	capabilities, err := capabilityRegistry.GetCapabilities(&bind.CallOpts{})
	require.NoError(t, err, "GetCapabilities failed")

	fmt.Println("Capability:", capabilities)

	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	var pid ragetypes.PeerID
	err = pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
	require.NoError(t, err)
	peer := mocks.NewPeer(t)
	peer.On("UpdateConnections", mock.Anything).Return(nil)
	peer.On("ID").Return(pid)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)
	registry := commonMocks.NewCapabilitiesRegistry(t)
	registry.On("Add", mock.Anything, mock.Anything).Return(nil)
	dispatcher := remoteMocks.NewDispatcher(t)
	dispatcher.On("SetReceiver", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	onchainRegistry := coreCapabilities.NewOnchainCapabilityRegistry(types.MustEIP55Address("0x0000000000000000000000000000000000000001").Address(), lggr)

	syncer := coreCapabilities.NewRegistrySyncer(wrapper, registry, dispatcher, lggr, onchainRegistry)
	require.NoError(t, syncer.Start(ctx))
	require.NoError(t, syncer.Close())
}
