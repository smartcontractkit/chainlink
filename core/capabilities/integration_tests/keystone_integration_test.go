package capabilities_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	utils "github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/internal"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/keystone_capability_registry"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func TestIntegration_GetCapabilities(t *testing.T) {
	owner, simulatedBackend := utils.StartNewChain(t)

	capabilityRegistry := utils.DeployCapabilityRegistry(t, owner, simulatedBackend)

	utils.AddCapability(t, owner, simulatedBackend, capabilityRegistry, utils.DataStreamsReportCapability)
	utils.AddCapability(t, owner, simulatedBackend, capabilityRegistry, utils.WriteChainCapability)

	capabilities, err := capabilityRegistry.GetCapabilities(&bind.CallOpts{})
	require.NoError(t, err, "GetCapabilities failed")

	fmt.Println("Capabilities:", capabilities)

	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)

	simulatedBackendClient := client.NewSimulatedBackendClient(t, simulatedBackend, testutils.SimulatedChainID)

	// This requires having `CL_DATABASE_URL` set to something. DB does not appear to be used.
	db := pgtest.NewSqlxDB(t)
	lpOpts := logpoller.Opts{
		PollPeriod:               time.Millisecond,
		FinalityDepth:            4,
		BackfillBatchSize:        1,
		RpcBatchSize:             1,
		KeepFinalizedBlocksDepth: 10000,
	}
	lp := logpoller.NewLogPoller(
		logpoller.NewORM(testutils.SimulatedChainID, db, lggr),
		simulatedBackendClient,
		lggr,
		lpOpts,
	)

	require.NoError(t, lp.Start(ctx))

	chainConfig := types.ChainReaderConfig{
		Contracts: map[string]types.ChainContractReader{
			"capability_registry": {
				ContractABI: keystone_capability_registry.CapabilityRegistryABI,
				Configs: map[string]*types.ChainReaderDefinition{
					"get_capabilities": {
						ChainSpecificName: "getCapabilities",
						OutputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{Fields: map[string]string{"labelledName": "name"}},
						},
					},
				},
			},
		},
	}
	cr, err := evm.NewChainReaderService(ctx, lggr, lp, simulatedBackendClient, chainConfig)
	require.NoError(t, err)

	require.NoError(t, cr.Bind(ctx, []commontypes.BoundContract{
		{
			Name:    "capability_registry",
			Address: capabilityRegistry.Address().String(),
		}}))

	require.NoError(t, cr.Start(ctx))

	type Cap struct {
		Name                  string
		Version               string
		ResponseType          int
		ConfigurationContract []byte
	}

	var returnedCapabilities []Cap

	err = cr.GetLatestValue(ctx, "capability_registry", "get_capabilities", nil, &returnedCapabilities)
	require.NoError(t, err)

	fmt.Println("Returned capabilities:", returnedCapabilities)

	// SYNCER DEPENDENCIES
	// var pid ragetypes.PeerID
	// err = pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
	// require.NoError(t, err)
	// peer := mocks.NewPeer(t)
	// peer.On("UpdateConnections", mock.Anything).Return(nil)
	// peer.On("ID").Return(pid)
	// wrapper := mocks.NewPeerWrapper(t)
	// wrapper.On("GetPeer").Return(peer)
	// workflowEngineCapabilitiesRegistry := commonMocks.NewCapabilitiesRegistry(t)
	// workflowEngineCapabilitiesRegistry.On("Add", mock.Anything, mock.Anything).Return(nil)
	// dispatcher := remoteMocks.NewDispatcher(t)
	// dispatcher.On("SetReceiver", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// onchainRegistry := coreCapabilities.NewRemoteRegistry(capabilityRegistry.Address(), lggr)

	// syncer := coreCapabilities.NewRegistrySyncer(
	// 	wrapper,
	// 	workflowEngineCapabilitiesRegistry,
	// 	dispatcher,
	// 	lggr,
	// 	onchainRegistry,
	// 	simulatedBackendClient,
	// )
	// require.NoError(t, syncer.Start(ctx))

	// // Do assertions here

	// require.NoError(t, syncer.Close())
}
