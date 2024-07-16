package keeper_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/gethwrappers/link_token_interface"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/basic_upkeep_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic1_3"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_3"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keeper"
	webpresenters "github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

var (
	oneEth    = big.NewInt(1000000000000000000)
	tenEth    = big.NewInt(0).Mul(oneEth, big.NewInt(10))
	oneHunEth = big.NewInt(0).Mul(oneEth, big.NewInt(100))

	payload1 = common.Hex2Bytes("1234")
	payload2 = common.Hex2Bytes("ABCD")
	payload3 = common.Hex2Bytes("6789")
)

func deployKeeperRegistry(
	t *testing.T,
	version keeper.RegistryVersion,
	auth *bind.TransactOpts,
	backend *client.SimulatedBackendClient,
	linkAddr, linkFeedAddr, gasFeedAddr common.Address,
) (common.Address, *keeper.RegistryWrapper) {
	switch version {
	case keeper.RegistryVersion_1_1:
		regAddr, _, _, err := keeper_registry_wrapper1_1.DeployKeeperRegistry(
			auth,
			backend,
			linkAddr,
			linkFeedAddr,
			gasFeedAddr,
			250_000_000,
			0,
			big.NewInt(1),
			20_000_000,
			big.NewInt(3600),
			1,
			big.NewInt(60000000000),
			big.NewInt(20000000000000000),
		)
		require.NoError(t, err)
		backend.Commit()

		wrapper, err := keeper.NewRegistryWrapper(evmtypes.EIP55AddressFromAddress(regAddr), backend)
		require.NoError(t, err)
		return regAddr, wrapper
	case keeper.RegistryVersion_1_2:
		regAddr, _, _, err := keeper_registry_wrapper1_2.DeployKeeperRegistry(
			auth,
			backend,
			linkAddr,
			linkFeedAddr,
			gasFeedAddr,
			keeper_registry_wrapper1_2.Config{
				PaymentPremiumPPB:    250_000_000,
				FlatFeeMicroLink:     0,
				BlockCountPerTurn:    big.NewInt(1),
				CheckGasLimit:        20_000_000,
				StalenessSeconds:     big.NewInt(3600),
				GasCeilingMultiplier: 1,
				MinUpkeepSpend:       big.NewInt(0),
				MaxPerformGas:        5_000_000,
				FallbackGasPrice:     big.NewInt(60000000000),
				FallbackLinkPrice:    big.NewInt(20000000000000000),
				Transcoder:           testutils.NewAddress(),
				Registrar:            testutils.NewAddress(),
			},
		)
		require.NoError(t, err)
		backend.Commit()
		wrapper, err := keeper.NewRegistryWrapper(evmtypes.EIP55AddressFromAddress(regAddr), backend)
		require.NoError(t, err)
		return regAddr, wrapper
	case keeper.RegistryVersion_1_3:
		logicAddr, _, _, err := keeper_registry_logic1_3.DeployKeeperRegistryLogic(
			auth,
			backend,
			0,
			big.NewInt(80000),
			linkAddr,
			linkFeedAddr,
			gasFeedAddr)
		require.NoError(t, err)
		backend.Commit()

		regAddr, _, _, err := keeper_registry_wrapper1_3.DeployKeeperRegistry(
			auth,
			backend,
			logicAddr,
			keeper_registry_wrapper1_3.Config{
				PaymentPremiumPPB:    250_000_000,
				FlatFeeMicroLink:     0,
				BlockCountPerTurn:    big.NewInt(1),
				CheckGasLimit:        20_000_000,
				StalenessSeconds:     big.NewInt(3600),
				GasCeilingMultiplier: 1,
				MinUpkeepSpend:       big.NewInt(0),
				MaxPerformGas:        5_000_000,
				FallbackGasPrice:     big.NewInt(60000000000),
				FallbackLinkPrice:    big.NewInt(20000000000000000),
				Transcoder:           testutils.NewAddress(),
				Registrar:            testutils.NewAddress(),
			},
		)
		require.NoError(t, err)
		backend.Commit()
		wrapper, err := keeper.NewRegistryWrapper(evmtypes.EIP55AddressFromAddress(regAddr), backend)
		require.NoError(t, err)
		return regAddr, wrapper
	default:
		panic(errors.Errorf("Deployment of registry verdion %d not defined", version))
	}
}

func getUpkeepIdFromTx(t *testing.T, registryWrapper *keeper.RegistryWrapper, registrationTx *types.Transaction, backend *client.SimulatedBackendClient) *big.Int {
	receipt, err := backend.TransactionReceipt(testutils.Context(t), registrationTx.Hash())
	require.NoError(t, err)
	upkeepId, err := registryWrapper.GetUpkeepIdFromRawRegistrationLog(*receipt.Logs[0])
	require.NoError(t, err)
	return upkeepId
}

func TestKeeperEthIntegration(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		eip1559         bool
		registryVersion keeper.RegistryVersion
	}{
		// name should be a valid ORM name, only containing alphanumeric/underscore
		{"legacy_registry1_1", false, keeper.RegistryVersion_1_1},
		{"eip1559_registry1_1", true, keeper.RegistryVersion_1_1},
		{"legacy_registry1_2", false, keeper.RegistryVersion_1_2},
		{"eip1559_registry1_2", true, keeper.RegistryVersion_1_2},
		{"legacy_registry1_3", false, keeper.RegistryVersion_1_3},
		{"eip1559_registry1_3", true, keeper.RegistryVersion_1_3},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx := testutils.Context(t)
			g := gomega.NewWithT(t)

			// setup node key
			nodeKey := cltest.MustGenerateRandomKey(t)
			nodeAddress := nodeKey.Address
			nodeAddressEIP55 := evmtypes.EIP55AddressFromAddress(nodeAddress)

			// setup blockchain
			sergey := testutils.MustNewSimTransactor(t) // owns all the link
			steve := testutils.MustNewSimTransactor(t)  // registry owner
			carrol := testutils.MustNewSimTransactor(t) // client
			nelly := testutils.MustNewSimTransactor(t)  // other keeper operator 1
			nick := testutils.MustNewSimTransactor(t)   // other keeper operator 2
			genesisData := core.GenesisAlloc{
				sergey.From: {Balance: assets.Ether(1000).ToInt()},
				steve.From:  {Balance: assets.Ether(1000).ToInt()},
				carrol.From: {Balance: assets.Ether(1000).ToInt()},
				nelly.From:  {Balance: assets.Ether(1000).ToInt()},
				nick.From:   {Balance: assets.Ether(1000).ToInt()},
				nodeAddress: {Balance: assets.Ether(1000).ToInt()},
			}

			gasLimit := uint32(ethconfig.Defaults.Miner.GasCeil * 2)
			b := cltest.NewSimulatedBackend(t, genesisData, gasLimit)
			backend := client.NewSimulatedBackendClient(t, b, testutils.SimulatedChainID)

			stopMining := cltest.Mine(backend.Backend(), 1*time.Second) // >> 2 seconds and the test gets slow, << 1 second and the app may miss heads
			defer stopMining()

			linkAddr, _, linkToken, err := link_token_interface.DeployLinkToken(sergey, backend)
			require.NoError(t, err)
			gasFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(steve, backend, 18, big.NewInt(60000000000))
			require.NoError(t, err)
			linkFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(steve, backend, 18, big.NewInt(20000000000000000))
			require.NoError(t, err)

			regAddr, registryWrapper := deployKeeperRegistry(t, test.registryVersion, steve, backend, linkAddr, linkFeedAddr, gasFeedAddr)

			upkeepAddr, _, upkeepContract, err := basic_upkeep_contract.DeployBasicUpkeepContract(carrol, backend)
			require.NoError(t, err)
			_, err = linkToken.Transfer(sergey, carrol.From, oneHunEth)
			require.NoError(t, err)
			_, err = linkToken.Approve(carrol, regAddr, oneHunEth)
			require.NoError(t, err)
			_, err = registryWrapper.SetKeepers(steve, []common.Address{nodeAddress, nelly.From}, []common.Address{nodeAddress, nelly.From})
			require.NoError(t, err)
			registrationTx, err := registryWrapper.RegisterUpkeep(steve, upkeepAddr, 2_500_000, carrol.From, []byte{})
			require.NoError(t, err)
			backend.Commit()
			upkeepID := getUpkeepIdFromTx(t, registryWrapper, registrationTx, backend)

			_, err = upkeepContract.SetBytesToSend(carrol, payload1)
			require.NoError(t, err)
			_, err = upkeepContract.SetShouldPerformUpkeep(carrol, true)
			require.NoError(t, err)
			_, err = registryWrapper.AddFunds(carrol, upkeepID, tenEth)
			require.NoError(t, err)
			backend.Commit()

			// setup app
			config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.EIP1559DynamicFees = &test.eip1559
				c.Keeper.MaxGracePeriod = ptr[int64](0)                                       // avoid waiting to re-submit for upkeeps
				c.Keeper.Registry.SyncInterval = commonconfig.MustNewDuration(24 * time.Hour) // disable full sync ticker for test

				c.Keeper.TurnLookBack = ptr[int64](0) // testing doesn't need to do far look back

				c.EVM[0].BlockBackfillDepth = ptr[uint32](0)          // backfill will trigger sync on startup
				c.EVM[0].MinIncomingConfirmations = ptr[uint32](1)    // disable reorg protection for this test
				c.EVM[0].HeadTracker.MaxBufferSize = ptr[uint32](100) // helps prevent missed heads
			})
			korm := keeper.NewORM(db, logger.TestLogger(t))

			app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, backend.Backend(), nodeKey)
			require.NoError(t, app.Start(ctx))

			// create job
			regAddrEIP55 := evmtypes.EIP55AddressFromAddress(regAddr)
			job := cltest.MustInsertKeeperJob(t, db, korm, nodeAddressEIP55, regAddrEIP55)
			err = app.JobSpawner().StartService(ctx, job)
			require.NoError(t, err)

			// keeper job is triggered and payload is received
			receivedBytes := func() []byte {
				received, err2 := upkeepContract.ReceivedBytes(nil)
				require.NoError(t, err2)
				return received
			}
			g.Eventually(receivedBytes, 20*time.Second, cltest.DBPollingInterval).Should(gomega.Equal(payload1))

			// submit from other keeper (because keepers must alternate)
			_, err = registryWrapper.PerformUpkeep(nelly, upkeepID, []byte{})
			require.NoError(t, err)

			// change payload
			_, err = upkeepContract.SetBytesToSend(carrol, payload2)
			require.NoError(t, err)
			_, err = upkeepContract.SetShouldPerformUpkeep(carrol, true)
			require.NoError(t, err)

			// observe 2nd job run and received payload changes
			g.Eventually(receivedBytes, 20*time.Second, cltest.DBPollingInterval).Should(gomega.Equal(payload2))

			// cancel upkeep
			_, err = registryWrapper.CancelUpkeep(carrol, upkeepID)
			require.NoError(t, err)
			backend.Commit()

			cltest.WaitForCount(t, app.GetDB(), "upkeep_registrations", 0)

			// add new upkeep (same target contract)
			registrationTx, err = registryWrapper.RegisterUpkeep(steve, upkeepAddr, 2_500_000, carrol.From, []byte{})
			require.NoError(t, err)
			backend.Commit()

			upkeepID = getUpkeepIdFromTx(t, registryWrapper, registrationTx, backend)
			_, err = upkeepContract.SetBytesToSend(carrol, payload3)
			require.NoError(t, err)
			_, err = upkeepContract.SetShouldPerformUpkeep(carrol, true)
			require.NoError(t, err)
			_, err = registryWrapper.AddFunds(carrol, upkeepID, tenEth)
			require.NoError(t, err)
			backend.Commit()

			// observe update
			g.Eventually(receivedBytes, 20*time.Second, cltest.DBPollingInterval).Should(gomega.Equal(payload3))

			// remove this node from keeper list
			_, err = registryWrapper.SetKeepers(steve, []common.Address{nick.From, nelly.From}, []common.Address{nick.From, nelly.From})
			require.NoError(t, err)

			var registry keeper.Registry
			require.NoError(t, app.GetDB().GetContext(ctx, &registry, `SELECT * FROM keeper_registries`))
			cltest.AssertRecordEventually(t, app.GetDB(), &registry, fmt.Sprintf("SELECT * FROM keeper_registries WHERE id = %d", registry.ID), func() bool {
				return registry.KeeperIndex == -1
			})
			runs, err := app.PipelineORM().GetAllRuns(ctx)
			require.NoError(t, err)
			// Since we set grace period to 0, we can have more than 1 pipeline run per perform
			// This happens in case we start a pipeline run before previous perform tx is committed to chain
			require.GreaterOrEqual(t, len(runs), 3)
			prr := webpresenters.NewPipelineRunResource(runs[0], logger.TestLogger(t))
			require.Equal(t, 1, len(prr.Outputs))
			require.Nil(t, prr.Outputs[0])
		})
	}
}

func TestKeeperForwarderEthIntegration(t *testing.T) {
	t.Parallel()
	t.Run("keeper_forwarder_flow", func(t *testing.T) {
		ctx := testutils.Context(t)
		g := gomega.NewWithT(t)

		// setup node key
		nodeKey := cltest.MustGenerateRandomKey(t)
		nodeAddress := nodeKey.Address
		nodeAddressEIP55 := evmtypes.EIP55AddressFromAddress(nodeAddress)

		// setup blockchain
		sergey := testutils.MustNewSimTransactor(t) // owns all the link
		steve := testutils.MustNewSimTransactor(t)  // registry owner
		carrol := testutils.MustNewSimTransactor(t) // client
		nelly := testutils.MustNewSimTransactor(t)  // other keeper operator 1
		nick := testutils.MustNewSimTransactor(t)   // other keeper operator 2
		genesisData := core.GenesisAlloc{
			sergey.From: {Balance: assets.Ether(1000).ToInt()},
			steve.From:  {Balance: assets.Ether(1000).ToInt()},
			carrol.From: {Balance: assets.Ether(1000).ToInt()},
			nelly.From:  {Balance: assets.Ether(1000).ToInt()},
			nick.From:   {Balance: assets.Ether(1000).ToInt()},
			nodeAddress: {Balance: assets.Ether(1000).ToInt()},
		}

		gasLimit := uint32(ethconfig.Defaults.Miner.GasCeil * 2)
		b := cltest.NewSimulatedBackend(t, genesisData, gasLimit)
		backend := client.NewSimulatedBackendClient(t, b, testutils.SimulatedChainID)

		stopMining := cltest.Mine(backend.Backend(), 1*time.Second) // >> 2 seconds and the test gets slow, << 1 second and the app may miss heads
		defer stopMining()

		linkAddr, _, linkToken, err := link_token_interface.DeployLinkToken(sergey, backend)
		require.NoError(t, err)
		gasFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(steve, backend, 18, big.NewInt(60000000000))
		require.NoError(t, err)
		linkFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(steve, backend, 18, big.NewInt(20000000000000000))
		require.NoError(t, err)

		regAddr, registryWrapper := deployKeeperRegistry(t, keeper.RegistryVersion_1_3, steve, backend, linkAddr, linkFeedAddr, gasFeedAddr)

		fwdrAddress, _, authorizedForwarder, err := authorized_forwarder.DeployAuthorizedForwarder(sergey, backend, linkAddr, sergey.From, steve.From, []byte{})
		require.NoError(t, err)
		_, err = authorizedForwarder.SetAuthorizedSenders(sergey, []common.Address{nodeAddress})
		require.NoError(t, err)

		upkeepAddr, _, upkeepContract, err := basic_upkeep_contract.DeployBasicUpkeepContract(carrol, backend)
		require.NoError(t, err)
		_, err = linkToken.Transfer(sergey, carrol.From, oneHunEth)
		require.NoError(t, err)
		_, err = linkToken.Approve(carrol, regAddr, oneHunEth)
		require.NoError(t, err)
		_, err = registryWrapper.SetKeepers(steve, []common.Address{fwdrAddress, nelly.From}, []common.Address{nodeAddress, nelly.From})
		require.NoError(t, err)
		registrationTx, err := registryWrapper.RegisterUpkeep(steve, upkeepAddr, 2_500_000, carrol.From, []byte{})
		require.NoError(t, err)
		backend.Commit()
		upkeepID := getUpkeepIdFromTx(t, registryWrapper, registrationTx, backend)

		_, err = upkeepContract.SetBytesToSend(carrol, payload1)
		require.NoError(t, err)
		_, err = upkeepContract.SetShouldPerformUpkeep(carrol, true)
		require.NoError(t, err)
		_, err = registryWrapper.AddFunds(carrol, upkeepID, tenEth)
		require.NoError(t, err)
		backend.Commit()

		// setup app
		config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.Feature.LogPoller = ptr(true)
			c.EVM[0].GasEstimator.EIP1559DynamicFees = ptr(true)
			c.Keeper.MaxGracePeriod = ptr[int64](0)                                       // avoid waiting to re-submit for upkeeps
			c.Keeper.Registry.SyncInterval = commonconfig.MustNewDuration(24 * time.Hour) // disable full sync ticker for test

			c.Keeper.TurnLookBack = ptr[int64](0) // testing doesn't need to do far look back

			c.EVM[0].BlockBackfillDepth = ptr[uint32](0)          // backfill will trigger sync on startup
			c.EVM[0].MinIncomingConfirmations = ptr[uint32](1)    // disable reorg protection for this test
			c.EVM[0].HeadTracker.MaxBufferSize = ptr[uint32](100) // helps prevent missed heads
			c.EVM[0].Transactions.ForwardersEnabled = ptr(true)   // Enable Operator Forwarder flow
			c.EVM[0].ChainID = (*ubig.Big)(testutils.SimulatedChainID)
		})
		korm := keeper.NewORM(db, logger.TestLogger(t))

		app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, backend.Backend(), nodeKey)
		require.NoError(t, app.Start(ctx))

		forwarderORM := forwarders.NewORM(db)
		chainID := ubig.Big(*backend.ConfiguredChainID())
		_, err = forwarderORM.CreateForwarder(ctx, fwdrAddress, chainID)
		require.NoError(t, err)

		addr, err := app.GetRelayers().LegacyEVMChains().Slice()[0].TxManager().GetForwarderForEOA(ctx, nodeAddress)
		require.NoError(t, err)
		require.Equal(t, addr, fwdrAddress)

		// create job
		regAddrEIP55 := evmtypes.EIP55AddressFromAddress(regAddr)

		jb := job.Job{
			ID:   1,
			Type: job.Keeper,
			KeeperSpec: &job.KeeperSpec{
				FromAddress:     nodeAddressEIP55,
				ContractAddress: regAddrEIP55,
				EVMChainID:      (*ubig.Big)(testutils.SimulatedChainID),
			},
			SchemaVersion:     1,
			ForwardingAllowed: true,
		}
		err = app.JobORM().CreateJob(testutils.Context(t), &jb)
		require.NoError(t, err)

		registry := keeper.Registry{
			ContractAddress:   regAddrEIP55,
			BlockCountPerTurn: 1,
			CheckGas:          150_000,
			FromAddress:       nodeAddressEIP55,
			JobID:             jb.ID,
			KeeperIndex:       0,
			NumKeepers:        2,
			KeeperIndexMap: map[evmtypes.EIP55Address]int32{
				nodeAddressEIP55: 0,
				evmtypes.EIP55AddressFromAddress(nelly.From): 1,
			},
		}
		err = korm.UpsertRegistry(ctx, &registry)
		require.NoError(t, err)

		callOpts := bind.CallOpts{From: nodeAddress}
		// Read last keeper on the upkeep contract
		lastKeeper := func() common.Address {
			upkeepCfg, err2 := registryWrapper.GetUpkeep(&callOpts, upkeepID)
			require.NoError(t, err2)
			return upkeepCfg.LastKeeper
		}
		require.Equal(t, lastKeeper(), common.Address{})

		err = app.JobSpawner().StartService(ctx, jb)
		require.NoError(t, err)

		// keeper job is triggered and payload is received
		receivedBytes := func() []byte {
			received, err2 := upkeepContract.ReceivedBytes(nil)
			require.NoError(t, err2)
			return received
		}
		g.Eventually(receivedBytes, 20*time.Second, cltest.DBPollingInterval).Should(gomega.Equal(payload1))

		// Upkeep performed by the node through the forwarder
		g.Eventually(lastKeeper, 20*time.Second, cltest.DBPollingInterval).Should(gomega.Equal(fwdrAddress))
	})
}

func TestMaxPerformDataSize(t *testing.T) {
	t.Parallel()
	t.Run("max_perform_data_size_test", func(t *testing.T) {
		ctx := testutils.Context(t)
		maxPerformDataSize := 1000 // Will be set as config override
		g := gomega.NewWithT(t)

		// setup node key
		nodeKey := cltest.MustGenerateRandomKey(t)
		nodeAddress := nodeKey.Address
		nodeAddressEIP55 := evmtypes.EIP55AddressFromAddress(nodeAddress)

		// setup blockchain
		sergey := testutils.MustNewSimTransactor(t) // owns all the link
		steve := testutils.MustNewSimTransactor(t)  // registry owner
		carrol := testutils.MustNewSimTransactor(t) // client
		nelly := testutils.MustNewSimTransactor(t)  // other keeper operator 1
		nick := testutils.MustNewSimTransactor(t)   // other keeper operator 2
		genesisData := core.GenesisAlloc{
			sergey.From: {Balance: assets.Ether(1000).ToInt()},
			steve.From:  {Balance: assets.Ether(1000).ToInt()},
			carrol.From: {Balance: assets.Ether(1000).ToInt()},
			nelly.From:  {Balance: assets.Ether(1000).ToInt()},
			nick.From:   {Balance: assets.Ether(1000).ToInt()},
			nodeAddress: {Balance: assets.Ether(1000).ToInt()},
		}

		gasLimit := uint32(ethconfig.Defaults.Miner.GasCeil * 2)
		b := cltest.NewSimulatedBackend(t, genesisData, gasLimit)
		backend := client.NewSimulatedBackendClient(t, b, testutils.SimulatedChainID)

		stopMining := cltest.Mine(backend.Backend(), 1*time.Second) // >> 2 seconds and the test gets slow, << 1 second and the app may miss heads
		defer stopMining()

		linkAddr, _, linkToken, err := link_token_interface.DeployLinkToken(sergey, backend)
		require.NoError(t, err)
		gasFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(steve, backend, 18, big.NewInt(60000000000))
		require.NoError(t, err)
		linkFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(steve, backend, 18, big.NewInt(20000000000000000))
		require.NoError(t, err)

		regAddr, registryWrapper := deployKeeperRegistry(t, keeper.RegistryVersion_1_3, steve, backend, linkAddr, linkFeedAddr, gasFeedAddr)

		upkeepAddr, _, upkeepContract, err := basic_upkeep_contract.DeployBasicUpkeepContract(carrol, backend)
		require.NoError(t, err)
		_, err = linkToken.Transfer(sergey, carrol.From, oneHunEth)
		require.NoError(t, err)
		_, err = linkToken.Approve(carrol, regAddr, oneHunEth)
		require.NoError(t, err)
		_, err = registryWrapper.SetKeepers(steve, []common.Address{nodeAddress, nelly.From}, []common.Address{nodeAddress, nelly.From})
		require.NoError(t, err)
		registrationTx, err := registryWrapper.RegisterUpkeep(steve, upkeepAddr, 2_500_000, carrol.From, []byte{})
		require.NoError(t, err)
		backend.Commit()
		upkeepID := getUpkeepIdFromTx(t, registryWrapper, registrationTx, backend)

		_, err = registryWrapper.AddFunds(carrol, upkeepID, tenEth)
		require.NoError(t, err)
		backend.Commit()

		// setup app
		config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.Keeper.MaxGracePeriod = ptr[int64](0)                                       // avoid waiting to re-submit for upkeeps
			c.Keeper.Registry.SyncInterval = commonconfig.MustNewDuration(24 * time.Hour) // disable full sync ticker for test
			c.Keeper.Registry.MaxPerformDataSize = ptr(uint32(maxPerformDataSize))        // set the max perform data size

			c.Keeper.TurnLookBack = ptr[int64](0) // testing doesn't need to do far look back

			c.EVM[0].BlockBackfillDepth = ptr[uint32](0)          // backfill will trigger sync on startup
			c.EVM[0].MinIncomingConfirmations = ptr[uint32](1)    // disable reorg protection for this test
			c.EVM[0].HeadTracker.MaxBufferSize = ptr[uint32](100) // helps prevent missed heads
		})
		korm := keeper.NewORM(db, logger.TestLogger(t))

		app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, backend.Backend(), nodeKey)
		require.NoError(t, app.Start(ctx))

		// create job
		regAddrEIP55 := evmtypes.EIP55AddressFromAddress(regAddr)
		job := cltest.MustInsertKeeperJob(t, db, korm, nodeAddressEIP55, regAddrEIP55)
		err = app.JobSpawner().StartService(ctx, job)
		require.NoError(t, err)

		// keeper job is triggered
		receivedBytes := func() []byte {
			received, err2 := upkeepContract.ReceivedBytes(nil)
			require.NoError(t, err2)
			return received
		}

		hugePayload := make([]byte, maxPerformDataSize)
		_, err = upkeepContract.SetBytesToSend(carrol, hugePayload)
		require.NoError(t, err)
		_, err = upkeepContract.SetShouldPerformUpkeep(carrol, true)
		require.NoError(t, err)

		// Huge payload should not result in a perform
		g.Consistently(receivedBytes, 20*time.Second, cltest.DBPollingInterval).Should(gomega.Equal([]byte{}))

		// Set payload to be small and it should get received
		smallPayload := make([]byte, maxPerformDataSize-1)
		_, err = upkeepContract.SetBytesToSend(carrol, smallPayload)
		require.NoError(t, err)
		g.Eventually(receivedBytes, 20*time.Second, cltest.DBPollingInterval).Should(gomega.Equal(smallPayload))
	})
}
