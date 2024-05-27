package capabilities_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/mock"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	commonMocks "github.com/smartcontractkit/chainlink-common/pkg/types/mocks"

	cap "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	utils "github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/internal"
	remoteMocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
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

	// ==========================================================================================
	// START - Using ChainReaderService - This works, but we want to use a relayer instead.
	// ==========================================================================================

	// chainID, err := simulatedBackendClient.ChainID()
	// require.NoError(t, err)

	// chainConfig := types.ChainReaderConfig{
	// 	Contracts: map[string]types.ChainContractReader{
	// 		"capability_registry": {
	// 			ContractABI: keystone_capability_registry.CapabilityRegistryABI,
	// 			Configs: map[string]*types.ChainReaderDefinition{
	// 				"get_capabilities": {
	// 					ChainSpecificName: "getCapabilities",
	// 					OutputModifications: codec.ModifiersConfig{
	// 						&codec.RenameModifierConfig{Fields: map[string]string{"labelledName": "name"}},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }
	// cr, err := evm.NewChainReaderService(ctx, lggr, lp, simulatedBackendClient, chainConfig)
	// require.NoError(t, err)

	// require.NoError(t, cr.Bind(ctx, []commontypes.BoundContract{
	// 	{
	// 		Name:    "capability_registry",
	// 		Address: capabilityRegistry.Address().String(),
	// 	}}))

	// require.NoError(t, cr.Start(ctx))

	// type Cap struct {
	// 	Name                  string
	// 	Version               string
	// 	ResponseType          int
	// 	ConfigurationContract []byte
	// }

	// var returnedCapabilities []Cap

	// err = cr.GetLatestValue(ctx, "capability_registry", "get_capabilities", nil, &returnedCapabilities)
	// require.NoError(t, err)

	// fmt.Println("Returned capabilities:", returnedCapabilities)

	// ==========================================================================================
	// END - Using ChainReaderService
	// ==========================================================================================

	// SYNCER DEPENDENCIES
	var pid ragetypes.PeerID
	err = pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
	require.NoError(t, err)
	peer := mocks.NewPeer(t)
	peer.On("UpdateConnections", mock.Anything).Return(nil)
	peer.On("ID").Return(pid)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)
	workflowEngineCapabilitiesRegistry := commonMocks.NewCapabilitiesRegistry(t)
	workflowEngineCapabilitiesRegistry.On("Add", mock.Anything, mock.Anything).Return(nil)
	dispatcher := remoteMocks.NewDispatcher(t)
	dispatcher.On("SetReceiver", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// ==========================================================================================
	// Setting up a new relayer that reads from the simulated backend. This part doesn't work.
	// ==========================================================================================

	keyStore := cltest.NewKeyStore(t, db)
	mockChain := &evmmocks.Chain{}
	c := client.NewSimulatedBackendClient(t, simulatedBackend, big.NewInt(1337))
	mockChain.On("Client").Return(c)
	require.NoError(t, lp.Start(ctx))
	mockChain.On("LogPoller").Return(lp)

	relayer, err := evm.NewRelayer(
		lggr,
		mockChain,
		evm.RelayerOpts{
			DS:                   db,
			CSAETHKeystore:       keyStore,
			CapabilitiesRegistry: workflowEngineCapabilitiesRegistry,
		},
	)
	require.NoError(t, err)

	syncer := cap.NewRegistrySyncer(
		wrapper,
		workflowEngineCapabilitiesRegistry,
		dispatcher,
		lggr,
		relayer, // relayer
		capabilityRegistry.Address().String(),
	)

	require.NoError(t, syncer.Start(ctx))

	// Syncer.LocalState().getCapabilities()
	// fmt.Println("Synced capabilities:", returnedCapabilities)

	// // Do assertions here

	// require.NoError(t, syncer.Close())
}
