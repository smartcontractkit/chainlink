package keeper_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	goEthereumEth "github.com/ethereum/go-ethereum/eth"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/basic_upkeep_contract"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/libocr/gethwrappers/link_token_interface"
	"github.com/stretchr/testify/require"
)

var oneEth = big.NewInt(1000000000000000000)
var oneHunEth = big.NewInt(0).Mul(oneEth, big.NewInt(100))

func TestKeeperEthIntegration(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	// setup node key
	nodeKey := cltest.MustGenerateRandomKey(t)
	nodeAddress := nodeKey.Address.Address()
	nodeAddressEIP55, err := models.EIP55AddressFromAddress(nodeAddress)
	require.NoError(t, err)

	// setup blockchain
	sergey := cltest.NewSimulatedBackendIdentity(t) // owns all the link
	steve := cltest.NewSimulatedBackendIdentity(t)  // registry owner
	carrol := cltest.NewSimulatedBackendIdentity(t) // client
	nelly := cltest.NewSimulatedBackendIdentity(t)  // other keeper operator
	genesisData := core.GenesisAlloc{
		sergey.From: {Balance: oneEth},
		steve.From:  {Balance: oneEth},
		carrol.From: {Balance: oneEth},
		nelly.From:  {Balance: oneEth},
		nodeAddress: {Balance: oneEth},
	}

	gasLimit := goEthereumEth.DefaultConfig.Miner.GasCeil * 2
	backend := backends.NewSimulatedBackend(genesisData, gasLimit)

	payload1 := common.Hex2Bytes("1234")
	payload2 := common.Hex2Bytes("ABCD")

	linkAddr, _, linkToken, err := link_token_interface.DeployLinkToken(sergey, backend)
	require.NoError(t, err)
	gasFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(steve, backend, 18, big.NewInt(60000000000))
	require.NoError(t, err)
	linkFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(steve, backend, 18, big.NewInt(20000000000000000))
	require.NoError(t, err)
	regAddr, _, registryContract, err := keeper_registry_wrapper.DeployKeeperRegistry(steve, backend, linkAddr, linkFeedAddr, gasFeedAddr, 250_000_000, big.NewInt(1), 20_000_000, big.NewInt(3600), big.NewInt(60000000000), big.NewInt(20000000000000000))
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
	backend.Commit()
	_, err = registryContract.AddFunds(carrol, big.NewInt(0), oneHunEth)
	require.NoError(t, err)
	backend.Commit()

	stopMining := cltest.Mine(backend)
	defer stopMining()

	// setup app
	config, _, cfgCleanup := cltest.BootstrapThrowawayORM(t, "keeper_eth_integration", true, true)
	config.Config.Dialect = dialects.PostgresWithoutLock
	defer cfgCleanup()
	config.Set("KEEPER_REGISTRY_SYNC_INTERVAL", 2*time.Second)
	app, appCleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, backend, nodeKey)
	defer appCleanup()
	require.NoError(t, app.StartAndConnect())

	// create job
	regAddrEIP55, err := models.EIP55AddressFromAddress(regAddr)
	require.NoError(t, err)
	cltest.MustInsertKeeperJob(t, app.Store, nodeAddressEIP55, regAddrEIP55)

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
}
