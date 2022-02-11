package keeper_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/basic_upkeep_contract"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	webpresenters "github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/smartcontractkit/libocr/gethwrappers/link_token_interface"
)

var (
	oneEth    = big.NewInt(1000000000000000000)
	tenEth    = big.NewInt(0).Mul(oneEth, big.NewInt(10))
	oneHunEth = big.NewInt(0).Mul(oneEth, big.NewInt(100))

	payload1 = common.Hex2Bytes("1234")
	payload2 = common.Hex2Bytes("ABCD")
	payload3 = common.Hex2Bytes("6789")
)

func TestKeeperEthIntegration(t *testing.T) {
	tests := []struct {
		name    string
		eip1559 bool
	}{
		{"legacy mode", false},
		{"eip1559 mode", true},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			g := gomega.NewWithT(t)

			// setup node key
			nodeKey := cltest.MustGenerateRandomKey(t)
			nodeAddress := nodeKey.Address.Address()
			nodeAddressEIP55 := ethkey.EIP55AddressFromAddress(nodeAddress)

			// setup blockchain
			sergey := cltest.NewSimulatedBackendIdentity(t) // owns all the link
			steve := cltest.NewSimulatedBackendIdentity(t)  // registry owner
			carrol := cltest.NewSimulatedBackendIdentity(t) // client
			nelly := cltest.NewSimulatedBackendIdentity(t)  // other keeper operator 1
			nick := cltest.NewSimulatedBackendIdentity(t)   // other keeper operator 2
			genesisData := core.GenesisAlloc{
				sergey.From: {Balance: assets.Ether(1000)},
				steve.From:  {Balance: assets.Ether(1000)},
				carrol.From: {Balance: assets.Ether(1000)},
				nelly.From:  {Balance: assets.Ether(1000)},
				nick.From:   {Balance: assets.Ether(1000)},
				nodeAddress: {Balance: assets.Ether(1000)},
			}

			gasLimit := ethconfig.Defaults.Miner.GasCeil * 2
			backend := cltest.NewSimulatedBackend(t, genesisData, gasLimit)

			stopMining := cltest.Mine(backend, 1*time.Second) // >> 2 seconds and the test gets slow, << 1 second and the app may miss heads
			defer stopMining()

			linkAddr, _, linkToken, err := link_token_interface.DeployLinkToken(sergey, backend)
			require.NoError(t, err)
			gasFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(steve, backend, 18, big.NewInt(60000000000))
			require.NoError(t, err)
			linkFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(steve, backend, 18, big.NewInt(20000000000000000))
			require.NoError(t, err)
			regAddr, _, registryContract, err := keeper_registry_wrapper.DeployKeeperRegistry(steve, backend, linkAddr, linkFeedAddr, gasFeedAddr, 250_000_000, 0, big.NewInt(1), 20_000_000, big.NewInt(3600), 1, big.NewInt(60000000000), big.NewInt(20000000000000000))
			require.NoError(t, err)
			upkeepAddr, _, upkeepContract, err := basic_upkeep_contract.DeployBasicUpkeepContract(carrol, backend)
			require.NoError(t, err)
			_, err = linkToken.Transfer(sergey, carrol.From, oneHunEth)
			require.NoError(t, err)
			_, err = linkToken.Approve(carrol, regAddr, oneHunEth)
			require.NoError(t, err)
			_, err = registryContract.SetKeepers(steve, []common.Address{nodeAddress, nelly.From}, []common.Address{nodeAddress, nelly.From})
			require.NoError(t, err)
			_, err = registryContract.RegisterUpkeep(steve, upkeepAddr, 2_500_000, carrol.From, []byte{})
			require.NoError(t, err)
			_, err = upkeepContract.SetBytesToSend(carrol, payload1)
			require.NoError(t, err)
			_, err = upkeepContract.SetShouldPerformUpkeep(carrol, true)
			require.NoError(t, err)
			_, err = registryContract.AddFunds(carrol, big.NewInt(0), tenEth)
			require.NoError(t, err)
			backend.Commit()

			// setup app
			config, db := heavyweight.FullTestDB(t, fmt.Sprintf("keeper_eth_integration_%v", test.eip1559), true, true)
			korm := keeper.NewORM(db, logger.TestLogger(t), nil, nil)
			config.Overrides.GlobalEvmEIP1559DynamicFees = null.BoolFrom(test.eip1559)
			d := 24 * time.Hour
			// disable full sync ticker for test
			config.Overrides.KeeperRegistrySyncInterval = &d
			// backfill will trigger sync on startup
			config.Overrides.BlockBackfillDepth = null.IntFrom(0)
			// disable reorg protection for this test
			config.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(1)
			// avoid waiting to re-submit for upkeeps
			config.Overrides.KeeperMaximumGracePeriod = null.IntFrom(0)
			// test with gas price feature enabled
			config.Overrides.KeeperCheckUpkeepGasPriceFeatureEnabled = null.BoolFrom(true)
			// helps prevent missed heads
			config.Overrides.GlobalEvmHeadTrackerMaxBufferSize = null.IntFrom(100)
			app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, backend, nodeKey)
			require.NoError(t, app.Start(context.TODO()))

			// create job
			regAddrEIP55 := ethkey.EIP55AddressFromAddress(regAddr)
			job := cltest.MustInsertKeeperJob(t, db, korm, nodeAddressEIP55, regAddrEIP55)
			err = app.JobSpawner().StartService(job)
			require.NoError(t, err)

			// keeper job is triggered and payload is received
			receivedBytes := func() []byte {
				received, err2 := upkeepContract.ReceivedBytes(nil)
				require.NoError(t, err2)
				return received
			}
			g.Eventually(receivedBytes, 20*time.Second, cltest.DBPollingInterval).Should(gomega.Equal(payload1))

			// submit from other keeper (because keepers must alternate)
			_, err = registryContract.PerformUpkeep(nelly, big.NewInt(0), []byte{})
			require.NoError(t, err)

			// change payload
			_, err = upkeepContract.SetBytesToSend(carrol, payload2)
			require.NoError(t, err)
			_, err = upkeepContract.SetShouldPerformUpkeep(carrol, true)
			require.NoError(t, err)

			// observe 2nd job run and received payload changes
			g.Eventually(receivedBytes, 20*time.Second, cltest.DBPollingInterval).Should(gomega.Equal(payload2))

			// cancel upkeep
			_, err = registryContract.CancelUpkeep(carrol, big.NewInt(0))
			require.NoError(t, err)
			backend.Commit()

			cltest.WaitForCount(t, app.GetSqlxDB(), "upkeep_registrations", 0)

			// add new upkeep (same target contract)
			_, err = registryContract.RegisterUpkeep(steve, upkeepAddr, 2_500_000, carrol.From, []byte{})
			require.NoError(t, err)
			_, err = upkeepContract.SetBytesToSend(carrol, payload3)
			require.NoError(t, err)
			_, err = upkeepContract.SetShouldPerformUpkeep(carrol, true)
			require.NoError(t, err)
			_, err = registryContract.AddFunds(carrol, big.NewInt(1), tenEth)
			require.NoError(t, err)
			backend.Commit()

			// observe update
			g.Eventually(receivedBytes, 20*time.Second, cltest.DBPollingInterval).Should(gomega.Equal(payload3))

			// remove this node from keeper list
			_, err = registryContract.SetKeepers(steve, []common.Address{nick.From, nelly.From}, []common.Address{nick.From, nelly.From})
			require.NoError(t, err)

			var registry keeper.Registry
			require.NoError(t, app.GetSqlxDB().Get(&registry, `SELECT * FROM keeper_registries`))
			cltest.AssertRecordEventually(t, app.GetSqlxDB(), &registry, fmt.Sprintf("SELECT * FROM keeper_registries WHERE id = %d", registry.ID), func() bool {
				return registry.KeeperIndex == -1
			})
			runs, err := app.PipelineORM().GetAllRuns()
			require.NoError(t, err)
			require.Equal(t, 3, len(runs))
			prr := webpresenters.NewPipelineRunResource(runs[0], logger.TestLogger(t))
			require.Equal(t, 1, len(prr.Outputs))
			require.Nil(t, prr.Outputs[0])
		})
	}
}
