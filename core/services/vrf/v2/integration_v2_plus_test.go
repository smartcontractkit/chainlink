package v2_test

import (
	"encoding/json"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	v2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_vrf_coordinator_v2plus"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_consumer_v2_plus_upgradeable_example"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_consumer_v2_upgradeable_example"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2plus"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_malicious_consumer_v2_plus"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_single_consumer"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_sub_owner"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2_proxy_admin"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2_transparent_upgradeable_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2plus_consumer_example"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2plus_reverting_example"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg/datatypes"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/proof"
	v22 "github.com/smartcontractkit/chainlink/v2/core/services/vrf/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrftesthelpers"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type coordinatorV2PlusUniverse struct {
	coordinatorV2UniverseCommon
	submanager                      *bind.TransactOpts // Subscription owner
	batchCoordinatorContract        *batch_vrf_coordinator_v2plus.BatchVRFCoordinatorV2Plus
	batchCoordinatorContractAddress common.Address
}

func newVRFCoordinatorV2PlusUniverse(t *testing.T, key ethkey.KeyV2, numConsumers int) coordinatorV2PlusUniverse {
	testutils.SkipShort(t, "VRFCoordinatorV2Universe")
	oracleTransactor, err := bind.NewKeyedTransactorWithChainID(key.ToEcdsaPrivKey(), testutils.SimulatedChainID)
	require.NoError(t, err)
	var (
		sergey       = testutils.MustNewSimTransactor(t)
		neil         = testutils.MustNewSimTransactor(t)
		ned          = testutils.MustNewSimTransactor(t)
		evil         = testutils.MustNewSimTransactor(t)
		reverter     = testutils.MustNewSimTransactor(t)
		submanager   = testutils.MustNewSimTransactor(t)
		nallory      = oracleTransactor
		vrfConsumers []*bind.TransactOpts
	)

	// Create consumer contract deployer identities
	for i := 0; i < numConsumers; i++ {
		vrfConsumers = append(vrfConsumers, testutils.MustNewSimTransactor(t))
	}

	genesisData := core.GenesisAlloc{
		sergey.From:     {Balance: assets.Ether(1000).ToInt()},
		neil.From:       {Balance: assets.Ether(1000).ToInt()},
		ned.From:        {Balance: assets.Ether(1000).ToInt()},
		nallory.From:    {Balance: assets.Ether(1000).ToInt()},
		evil.From:       {Balance: assets.Ether(1000).ToInt()},
		reverter.From:   {Balance: assets.Ether(1000).ToInt()},
		submanager.From: {Balance: assets.Ether(1000).ToInt()},
	}
	for _, consumer := range vrfConsumers {
		genesisData[consumer.From] = core.GenesisAccount{
			Balance: assets.Ether(1000).ToInt(),
		}
	}

	gasLimit := uint32(ethconfig.Defaults.Miner.GasCeil)
	consumerABI, err := abi.JSON(strings.NewReader(
		vrfv2plus_consumer_example.VRFV2PlusConsumerExampleABI))
	require.NoError(t, err)
	coordinatorABI, err := abi.JSON(strings.NewReader(
		vrf_coordinator_v2plus.VRFCoordinatorV2PlusABI))
	require.NoError(t, err)
	backend := cltest.NewSimulatedBackend(t, genesisData, gasLimit)
	// Deploy link
	linkAddress, _, linkContract, err := link_token_interface.DeployLinkToken(
		sergey, backend)
	require.NoError(t, err, "failed to deploy link contract to simulated ethereum blockchain")

	// Deploy feed
	linkEthFeed, _, _, err :=
		mock_v3_aggregator_contract.DeployMockV3AggregatorContract(
			evil, backend, 18, vrftesthelpers.WeiPerUnitLink.BigInt()) // 0.01 eth per link
	require.NoError(t, err)

	// Deploy blockhash store
	bhsAddress, _, bhsContract, err := blockhash_store.DeployBlockhashStore(neil, backend)
	require.NoError(t, err, "failed to deploy BlockhashStore contract to simulated ethereum blockchain")

	// Deploy batch blockhash store
	batchBHSAddress, _, batchBHSContract, err := batch_blockhash_store.DeployBatchBlockhashStore(neil, backend, bhsAddress)
	require.NoError(t, err, "failed to deploy BatchBlockhashStore contract to simulated ethereum blockchain")

	// Deploy VRF V2plus coordinator
	coordinatorAddress, _, coordinatorContract, err :=
		vrf_coordinator_v2plus.DeployVRFCoordinatorV2Plus(
			neil, backend, bhsAddress)
	require.NoError(t, err, "failed to deploy VRFCoordinatorV2 contract to simulated ethereum blockchain")
	backend.Commit()

	_, err = coordinatorContract.SetLINK(neil, linkAddress)
	require.NoError(t, err)
	backend.Commit()

	_, err = coordinatorContract.SetLinkEthFeed(neil, linkEthFeed)
	require.NoError(t, err)
	backend.Commit()

	// Deploy batch VRF V2 coordinator
	batchCoordinatorAddress, _, batchCoordinatorContract, err :=
		batch_vrf_coordinator_v2plus.DeployBatchVRFCoordinatorV2Plus(
			neil, backend, coordinatorAddress,
		)
	require.NoError(t, err, "failed to deploy BatchVRFCoordinatorV2 contract to simulated ethereum blockchain")
	backend.Commit()

	// Create the VRF consumers.
	var (
		consumerContracts         []vrftesthelpers.VRFConsumerContract
		consumerContractAddresses []common.Address
	)
	for _, author := range vrfConsumers {
		// Deploy a VRF consumer. It has a starting balance of 500 LINK.
		consumerContractAddress, _, consumerContract, err :=
			vrfv2plus_consumer_example.DeployVRFV2PlusConsumerExample(
				author, backend, coordinatorAddress, linkAddress)
		require.NoError(t, err, "failed to deploy VRFConsumer contract to simulated ethereum blockchain")
		_, err = linkContract.Transfer(sergey, consumerContractAddress, assets.Ether(500).ToInt()) // Actually, LINK
		require.NoError(t, err, "failed to send LINK to VRFConsumer contract on simulated ethereum blockchain")

		consumerContracts = append(consumerContracts, vrftesthelpers.NewVRFV2PlusConsumer(consumerContract))
		consumerContractAddresses = append(consumerContractAddresses, consumerContractAddress)

		backend.Commit()
	}

	// Deploy malicious consumer with 1 link
	maliciousConsumerContractAddress, _, maliciousConsumerContract, err :=
		vrf_malicious_consumer_v2_plus.DeployVRFMaliciousConsumerV2Plus(
			evil, backend, coordinatorAddress, linkAddress)
	require.NoError(t, err, "failed to deploy VRFMaliciousConsumer contract to simulated ethereum blockchain")
	_, err = linkContract.Transfer(sergey, maliciousConsumerContractAddress, assets.Ether(1).ToInt()) // Actually, LINK
	require.NoError(t, err, "failed to send LINK to VRFMaliciousConsumer contract on simulated ethereum blockchain")
	backend.Commit()

	// Deploy upgradeable consumer, proxy, and proxy admin
	upgradeableConsumerAddress, _, _, err := vrf_consumer_v2_plus_upgradeable_example.DeployVRFConsumerV2PlusUpgradeableExample(neil, backend)
	require.NoError(t, err, "failed to deploy upgradeable consumer to simulated ethereum blockchain")
	backend.Commit()

	proxyAdminAddress, _, proxyAdmin, err := vrfv2_proxy_admin.DeployVRFV2ProxyAdmin(neil, backend)
	require.NoError(t, err)
	backend.Commit()

	// provide abi-encoded initialize function call on the implementation contract
	// so that it's called upon the proxy construction, to initialize it.
	upgradeableAbi, err := vrf_consumer_v2_upgradeable_example.VRFConsumerV2UpgradeableExampleMetaData.GetAbi()
	require.NoError(t, err)
	initializeCalldata, err := upgradeableAbi.Pack("initialize", coordinatorAddress, linkAddress)
	hexified := hexutil.Encode(initializeCalldata)
	t.Log("initialize calldata:", hexified, "coordinator:", coordinatorAddress.String(), "link:", linkAddress)
	require.NoError(t, err)
	proxyAddress, _, _, err := vrfv2_transparent_upgradeable_proxy.DeployVRFV2TransparentUpgradeableProxy(
		neil, backend, upgradeableConsumerAddress, proxyAdminAddress, initializeCalldata)
	require.NoError(t, err)

	_, err = linkContract.Transfer(sergey, proxyAddress, assets.Ether(500).ToInt()) // Actually, LINK
	require.NoError(t, err)
	backend.Commit()

	implAddress, err := proxyAdmin.GetProxyImplementation(nil, proxyAddress)
	require.NoError(t, err)
	t.Log("impl address:", implAddress.String())
	require.Equal(t, upgradeableConsumerAddress, implAddress)

	proxiedConsumer, err := vrf_consumer_v2_plus_upgradeable_example.NewVRFConsumerV2PlusUpgradeableExample(
		proxyAddress, backend)
	require.NoError(t, err)

	cAddress, err := proxiedConsumer.COORDINATOR(nil)
	require.NoError(t, err)
	t.Log("coordinator address in proxy to upgradeable consumer:", cAddress.String())
	require.Equal(t, coordinatorAddress, cAddress)

	lAddress, err := proxiedConsumer.LINKTOKEN(nil)
	require.NoError(t, err)
	t.Log("link address in proxy to upgradeable consumer:", lAddress.String())
	require.Equal(t, linkAddress, lAddress)

	// Deploy always reverting consumer
	revertingConsumerContractAddress, _, revertingConsumerContract, err := vrfv2plus_reverting_example.DeployVRFV2PlusRevertingExample(
		reverter, backend, coordinatorAddress, linkAddress,
	)
	require.NoError(t, err, "failed to deploy VRFRevertingExample contract to simulated eth blockchain")
	_, err = linkContract.Transfer(sergey, revertingConsumerContractAddress, assets.Ether(500).ToInt()) // Actually, LINK
	require.NoError(t, err, "failed to send LINK to VRFRevertingExample contract on simulated eth blockchain")
	backend.Commit()

	// Set the configuration on the coordinator.
	_, err = coordinatorContract.SetConfig(neil,
		uint16(1),                              // minRequestConfirmations
		uint32(2.5e6),                          // gas limit
		uint32(60*60*24),                       // stalenessSeconds
		uint32(v22.GasAfterPaymentCalculation), // gasAfterPaymentCalculation
		big.NewInt(1e16),                       // 0.01 eth per link fallbackLinkPrice
		vrf_coordinator_v2plus.VRFCoordinatorV2PlusFeeConfig{
			FulfillmentFlatFeeLinkPPM: uint32(1000), // 0.001 LINK premium
			FulfillmentFlatFeeEthPPM:  uint32(5),    // 0.000005 ETH preimum
		},
	)
	require.NoError(t, err, "failed to set coordinator configuration")
	backend.Commit()

	return coordinatorV2PlusUniverse{
		coordinatorV2UniverseCommon: coordinatorV2UniverseCommon{
			vrfConsumers:              vrfConsumers,
			consumerContracts:         consumerContracts,
			consumerContractAddresses: consumerContractAddresses,

			revertingConsumerContract:        vrftesthelpers.NewRevertingConsumerPlus(revertingConsumerContract),
			revertingConsumerContractAddress: revertingConsumerContractAddress,

			consumerProxyContract:        vrftesthelpers.NewUpgradeableConsumerPlus(proxiedConsumer),
			consumerProxyContractAddress: proxiedConsumer.Address(),
			proxyAdminAddress:            proxyAdminAddress,

			rootContract:                     v22.NewCoordinatorV2Plus(coordinatorContract),
			rootContractAddress:              coordinatorAddress,
			linkContract:                     linkContract,
			linkContractAddress:              linkAddress,
			linkEthFeedAddress:               linkEthFeed,
			bhsContract:                      bhsContract,
			bhsContractAddress:               bhsAddress,
			batchBHSContract:                 batchBHSContract,
			batchBHSContractAddress:          batchBHSAddress,
			maliciousConsumerContract:        vrftesthelpers.NewMaliciousConsumerPlus(maliciousConsumerContract),
			maliciousConsumerContractAddress: maliciousConsumerContractAddress,
			backend:                          backend,
			coordinatorABI:                   &coordinatorABI,
			consumerABI:                      &consumerABI,
			sergey:                           sergey,
			neil:                             neil,
			ned:                              ned,
			nallory:                          nallory,
			evil:                             evil,
			reverter:                         reverter,
		},
		batchCoordinatorContract:        batchCoordinatorContract,
		batchCoordinatorContractAddress: batchCoordinatorAddress,
		submanager:                      submanager,
	}
}

func TestVRFV2PlusIntegration_SingleConsumer_HappyPath_BatchFulfillment(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1)
	testSingleConsumerHappyPathBatchFulfillment(
		t,
		ownerKey,
		uni.coordinatorV2UniverseCommon,
		uni.vrfConsumers[0],
		uni.consumerContracts[0],
		uni.consumerContractAddresses[0],
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		nil,
		5,     // number of requests to send
		false, // don't send big callback
		vrfcommon.V2Plus,
	)
}

func TestVRFV2PlusIntegration_SingleConsumer_HappyPath_BatchFulfillment_BigGasCallback(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1)
	testSingleConsumerHappyPathBatchFulfillment(
		t,
		ownerKey,
		uni.coordinatorV2UniverseCommon,
		uni.vrfConsumers[0],
		uni.consumerContracts[0],
		uni.consumerContractAddresses[0],
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		nil,
		5,    // number of requests to send
		true, // send big callback
		vrfcommon.V2Plus,
	)
}

func TestVRFV2PlusIntegration_SingleConsumer_HappyPath(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1)
	testSingleConsumerHappyPath(
		t,
		ownerKey,
		uni.coordinatorV2UniverseCommon,
		uni.vrfConsumers[0],
		uni.consumerContracts[0],
		uni.consumerContractAddresses[0],
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		nil,
		vrfcommon.V2Plus)
}

func TestVRFV2PlusIntegration_SingleConsumer_EOA_Request(t *testing.T) {
	t.Parallel()
	testEoa(t, false)
}

func TestVRFV2PlusIntegration_SingleConsumer_EOA_Request_Batching_Enabled(t *testing.T) {
	t.Parallel()
	testEoa(t, true)
}

func TestVRFV2PlusIntegration_SingleConsumer_EIP150_HappyPath(t *testing.T) {
	t.Parallel()
	callBackGasLimit := int64(2_500_000)            // base callback gas.
	eip150Fee := callBackGasLimit / 64              // premium needed for callWithExactGas
	coordinatorFulfillmentOverhead := int64(90_000) // fixed gas used in coordinator fulfillment
	gasLimit := callBackGasLimit + eip150Fee + coordinatorFulfillmentOverhead

	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, _ := heavyweight.FullTestDBV2(t, "vrfv2_singleconsumer_eip150_happypath", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = ptr(uint32(gasLimit))
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]
	// Create a subscription and fund with 500 LINK.
	subAmount := big.NewInt(1).Mul(big.NewInt(5e18), big.NewInt(100))
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, subAmount, uni.rootContract, uni.backend)

	// Fund gas lane.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		uni.coordinatorV2UniverseCommon,
		nil,
		vrfcommon.V2Plus,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the first randomness request.
	numWords := uint32(1)
	requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, uint32(callBackGasLimit), uni.rootContract, uni.backend)

	// Wait for simulation to pass.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	t.Log("Done!")
}

func TestVRFV2PlusIntegration_SingleConsumer_EIP150_Revert(t *testing.T) {
	t.Parallel()
	callBackGasLimit := int64(2_500_000)            // base callback gas.
	eip150Fee := int64(0)                           // no premium given for callWithExactGas
	coordinatorFulfillmentOverhead := int64(90_000) // fixed gas used in coordinator fulfillment
	gasLimit := callBackGasLimit + eip150Fee + coordinatorFulfillmentOverhead

	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, _ := heavyweight.FullTestDBV2(t, "vrfv2_singleconsumer_eip150_revert", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = ptr(uint32(gasLimit))
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]
	// Create a subscription and fund with 500 LINK.
	subAmount := big.NewInt(1).Mul(big.NewInt(5e18), big.NewInt(100))
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, subAmount, uni.rootContract, uni.backend)

	// Fund gas lane.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		uni.coordinatorV2UniverseCommon,
		nil,
		vrfcommon.V2Plus,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the first randomness request.
	numWords := uint32(1)
	requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, uint32(callBackGasLimit), uni.rootContract, uni.backend)

	// Simulation should not pass.
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 0
	}, 5*time.Second, time.Second).Should(gomega.BeTrue())

	t.Log("Done!")
}

func TestVRFV2PlusIntegration_SingleConsumer_NeedsBlockhashStore(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 2)
	testMultipleConsumersNeedBHS(
		t,
		ownerKey,
		uni.coordinatorV2UniverseCommon,
		uni.vrfConsumers,
		uni.consumerContracts,
		uni.consumerContractAddresses,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		nil,
		vrfcommon.V2Plus)
}

func TestVRFV2PlusIntegration_SingleConsumer_BlockHeaderFeeder(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1)
	testBlockHeaderFeeder(
		t,
		ownerKey,
		uni.coordinatorV2UniverseCommon,
		uni.vrfConsumers,
		uni.consumerContracts,
		uni.consumerContractAddresses,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		nil,
		vrfcommon.V2Plus)
}

func TestVRFV2PlusIntegration_SingleConsumer_NeedsTopUp(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1)
	testSingleConsumerNeedsTopUp(
		t,
		ownerKey,
		uni.coordinatorV2UniverseCommon,
		uni.vrfConsumers[0],
		uni.consumerContracts[0],
		uni.consumerContractAddresses[0],
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		nil,
		assets.Ether(1).ToInt(),   // initial funding of 1 LINK
		assets.Ether(100).ToInt(), // top up of 100 LINK
		vrfcommon.V2Plus,
	)
}

func TestVRFV2PlusIntegration_SingleConsumer_BigGasCallback_Sandwich(t *testing.T) {
	ownerKey := cltest.MustGenerateRandomKey(t)
	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(100)
	config, db := heavyweight.FullTestDBV2(t, "vrfv2_singleconsumer_bigcallback_sandwich", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(100), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint32](5_000_000)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]

	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, assets.Ether(2).ToInt(), uni.rootContract, uni.backend)

	// Fund gas lane.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		uni.coordinatorV2UniverseCommon,
		nil,
		vrfcommon.V2Plus,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make some randomness requests, each one block apart, which contain a single low-gas request sandwiched between two high-gas requests.
	numWords := uint32(2)
	reqIDs := []*big.Int{}
	callbackGasLimits := []uint32{2_500_000, 50_000, 1_500_000}
	for _, limit := range callbackGasLimits {
		requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, limit, uni.rootContract, uni.backend)
		reqIDs = append(reqIDs, requestID)
		uni.backend.Commit()
	}

	// Assert that we've completed 0 runs before adding 3 new requests.
	runs, err := app.PipelineORM().GetAllRuns()
	require.NoError(t, err)
	assert.Equal(t, 0, len(runs))
	assert.Equal(t, 3, len(reqIDs))

	// Wait for the 50_000 gas randomness request to be enqueued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	// After the first successful request, no more will be enqueued.
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("assert 1", "runs", len(runs))
		return len(runs) == 1
	}, 3*time.Second, 1*time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, reqIDs[1], subID, uni.backend, db)

	// Assert the random word was fulfilled
	assertRandomWordsFulfilled(t, reqIDs[1], false, uni.rootContract)

	// Assert that we've still only completed 1 run before adding new requests.
	runs, err = app.PipelineORM().GetAllRuns()
	require.NoError(t, err)
	assert.Equal(t, 1, len(runs))

	// Make some randomness requests, each one block apart, this time without a low-gas request present in the callbackGasLimit slice.
	callbackGasLimits = []uint32{2_500_000, 2_500_000, 2_500_000}
	for _, limit := range callbackGasLimits {
		_, _ = requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, limit, uni.rootContract, uni.backend)
		uni.backend.Commit()
	}

	// Fulfillment will not be enqueued because subscriber doesn't have enough LINK for any of the requests.
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("assert 1", "runs", len(runs))
		return len(runs) == 1
	}, 5*time.Second, 1*time.Second).Should(gomega.BeTrue())

	t.Log("Done!")
}

func TestVRFV2PlusIntegration_SingleConsumer_MultipleGasLanes(t *testing.T) {
	cheapKey := cltest.MustGenerateRandomKey(t)
	expensiveKey := cltest.MustGenerateRandomKey(t)
	cheapGasLane := assets.GWei(10)
	expensiveGasLane := assets.GWei(1000)
	config, db := heavyweight.FullTestDBV2(t, "vrfv2_singleconsumer_multiplegaslanes", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Cheap gas lane.
			Key:          ptr(cheapKey.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: cheapGasLane},
		}, v2.KeySpecific{
			// Expensive gas lane.
			Key:          ptr(expensiveKey.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: expensiveGasLane},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint32](5_000_000)
	})
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, cheapKey, expensiveKey)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), uni.rootContract, uni.backend)

	// Fund gas lanes.
	sendEth(t, ownerKey, uni.backend, cheapKey.Address, 10)
	sendEth(t, ownerKey, uni.backend, expensiveKey.Address, 10)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF jobs.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{cheapKey}, {expensiveKey}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		uni.coordinatorV2UniverseCommon,
		nil,
		vrfcommon.V2Plus,
		false,
		cheapGasLane, expensiveGasLane)
	cheapHash := jbs[0].VRFSpec.PublicKey.MustHash()
	expensiveHash := jbs[1].VRFSpec.PublicKey.MustHash()

	numWords := uint32(20)
	cheapRequestID, _ :=
		requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, cheapHash, subID, numWords, 500_000, uni.rootContract, uni.backend)

	// Wait for fulfillment to be queued for cheap key hash.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("assert 1", "runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, cheapRequestID, subID, uni.backend, db)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, cheapRequestID, true, uni.rootContract)

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)

	expensiveRequestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, expensiveHash, subID, numWords, 500_000, uni.rootContract, uni.backend)

	// We should not have any new fulfillments until a top up.
	gomega.NewWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("assert 2", "runs", len(runs))
		return len(runs) == 1
	}, 5*time.Second, 1*time.Second).Should(gomega.BeTrue())

	// Top up subscription with enough LINK to see the job through. 100 LINK should do the trick.
	_, err := consumerContract.TopUpSubscription(consumer, decimal.RequireFromString("100e18").BigInt())
	require.NoError(t, err)

	// Wait for fulfillment to be queued for expensive key hash.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("assert 1", "runs", len(runs))
		return len(runs) == 2
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, expensiveRequestID, subID, uni.backend, db)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, expensiveRequestID, true, uni.rootContract)

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)
}

func TestVRFV2PlusIntegration_SingleConsumer_AlwaysRevertingCallback_StillFulfilled(t *testing.T) {
	ownerKey := cltest.MustGenerateRandomKey(t)
	key := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, "vrfv2_singleconsumer_alwaysrevertingcallback", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 0)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key)
	consumer := uni.reverter
	consumerContract := uni.revertingConsumerContract
	consumerContractAddress := uni.revertingConsumerContractAddress

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), uni.rootContract, uni.backend)

	// Fund gas lane.
	sendEth(t, ownerKey, uni.backend, key.Address, 10)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		uni.coordinatorV2UniverseCommon,
		nil,
		vrfcommon.V2Plus,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the randomness request.
	numWords := uint32(20)
	requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, uni.rootContract, uni.backend)

	// Wait for fulfillment to be queued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, requestID, subID, uni.backend, db)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, requestID, false, uni.rootContract)
	t.Log("Done!")
}

func TestVRFV2PlusIntegration_ConsumerProxy_HappyPath(t *testing.T) {
	ownerKey := cltest.MustGenerateRandomKey(t)
	key1 := cltest.MustGenerateRandomKey(t)
	key2 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, "vrfv2_consumerproxy_happypath", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		}, v2.KeySpecific{
			Key:          ptr(key2.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 0)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1, key2)
	consumerOwner := uni.neil
	consumerContract := uni.consumerProxyContract
	consumerContractAddress := uni.consumerProxyContractAddress

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(
		t, consumerContract, consumerOwner, consumerContractAddress,
		assets.Ether(5).ToInt(), uni.rootContract, uni.backend)

	// Create gas lane.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	sendEth(t, ownerKey, uni.backend, key2.Address, 10)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job using key1 and key2 on the same gas lane.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1, key2}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		uni.coordinatorV2UniverseCommon,
		nil,
		vrfcommon.V2Plus,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the first randomness request.
	numWords := uint32(20)
	requestID1, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(
		t, consumerContract, consumerOwner, keyHash, subID, numWords, 750_000, uni.rootContract, uni.backend)

	// Wait for fulfillment to be queued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, requestID1, subID, uni.backend, db)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, requestID1, true, uni.rootContract)

	// Gas available will be around 724,385, which means that 750,000 - 724,385 = 25,615 gas was used.
	// This is ~20k more than what the non-proxied consumer uses.
	// So to be safe, users should probably over-estimate their fulfillment gas by ~25k.
	gasAvailable, err := consumerContract.SGasAvailable(nil)
	require.NoError(t, err)
	t.Log("gas available after proxied callback:", gasAvailable)

	// Make the second randomness request and assert fulfillment is successful
	requestID2, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(
		t, consumerContract, consumerOwner, keyHash, subID, numWords, 750_000, uni.rootContract, uni.backend)
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 2
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
	mine(t, requestID2, subID, uni.backend, db)
	assertRandomWordsFulfilled(t, requestID2, true, uni.rootContract)

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)

	// Assert that both send addresses were used to fulfill the requests
	n, err := uni.backend.PendingNonceAt(testutils.Context(t), key1.Address)
	require.NoError(t, err)
	require.EqualValues(t, 1, n)

	n, err = uni.backend.PendingNonceAt(testutils.Context(t), key2.Address)
	require.NoError(t, err)
	require.EqualValues(t, 1, n)

	t.Log("Done!")
}

func TestVRFV2PlusIntegration_ConsumerProxy_CoordinatorZeroAddress(t *testing.T) {
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 0)

	// Deploy another upgradeable consumer, proxy, and proxy admin
	// to test vrfCoordinator != 0x0 condition.
	upgradeableConsumerAddress, _, _, err := vrf_consumer_v2_upgradeable_example.DeployVRFConsumerV2UpgradeableExample(uni.neil, uni.backend)
	require.NoError(t, err, "failed to deploy upgradeable consumer to simulated ethereum blockchain")
	uni.backend.Commit()

	// Deployment should revert if we give the 0x0 address for the coordinator.
	upgradeableAbi, err := vrf_consumer_v2_upgradeable_example.VRFConsumerV2UpgradeableExampleMetaData.GetAbi()
	require.NoError(t, err)
	initializeCalldata, err := upgradeableAbi.Pack("initialize",
		common.BytesToAddress(common.LeftPadBytes([]byte{}, 20)), // zero address for the coordinator
		uni.linkContractAddress)
	require.NoError(t, err)
	_, _, _, err = vrfv2_transparent_upgradeable_proxy.DeployVRFV2TransparentUpgradeableProxy(
		uni.neil, uni.backend, upgradeableConsumerAddress, uni.proxyAdminAddress, initializeCalldata)
	require.Error(t, err)
}

func TestVRFV2PlusIntegration_ExternalOwnerConsumerExample(t *testing.T) {
	owner := testutils.MustNewSimTransactor(t)
	random := testutils.MustNewSimTransactor(t)
	genesisData := core.GenesisAlloc{
		owner.From:  {Balance: assets.Ether(10).ToInt()},
		random.From: {Balance: assets.Ether(10).ToInt()},
	}
	backend := cltest.NewSimulatedBackend(t, genesisData, uint32(ethconfig.Defaults.Miner.GasCeil))
	linkAddress, _, linkContract, err := link_token_interface.DeployLinkToken(
		owner, backend)
	require.NoError(t, err)
	backend.Commit()
	coordinatorAddress, _, coordinator, err :=
		vrf_coordinator_v2plus.DeployVRFCoordinatorV2Plus(
			owner, backend, common.Address{}) //bhs not needed for this test
	require.NoError(t, err)
	_, err = coordinator.SetConfig(owner, uint16(1), uint32(10000), 1, 1, big.NewInt(10), vrf_coordinator_v2plus.VRFCoordinatorV2PlusFeeConfig{
		FulfillmentFlatFeeLinkPPM: 0,
		FulfillmentFlatFeeEthPPM:  0,
	})
	require.NoError(t, err)
	backend.Commit()
	_, err = coordinator.SetLINK(owner, linkAddress)
	require.NoError(t, err)
	backend.Commit()
	consumerAddress, _, consumer, err := vrf_v2plus_sub_owner.DeployVRFV2PlusExternalSubOwnerExample(owner, backend, coordinatorAddress, linkAddress)
	require.NoError(t, err)
	backend.Commit()
	_, err = linkContract.Transfer(owner, consumerAddress, assets.Ether(2).ToInt())
	require.NoError(t, err)
	backend.Commit()
	AssertLinkBalances(t, linkContract, []common.Address{owner.From, consumerAddress}, []*big.Int{assets.Ether(999_999_998).ToInt(), assets.Ether(2).ToInt()})

	// Create sub, fund it and assign consumer
	_, err = coordinator.CreateSubscription(owner)
	require.NoError(t, err)
	backend.Commit()
	b, err := utils.ABIEncode(`[{"type":"uint64"}]`, uint64(1))
	require.NoError(t, err)
	_, err = linkContract.TransferAndCall(owner, coordinatorAddress, big.NewInt(0), b)
	require.NoError(t, err)
	_, err = coordinator.AddConsumer(owner, 1, consumerAddress)
	require.NoError(t, err)
	_, err = consumer.RequestRandomWords(random, 1, 1, 1, 1, [32]byte{}, false)
	require.Error(t, err)
	_, err = consumer.RequestRandomWords(owner, 1, 1, 1, 1, [32]byte{}, false)
	require.NoError(t, err)

	// Reassign ownership, check that only new owner can request
	_, err = consumer.TransferOwnership(owner, random.From)
	require.NoError(t, err)
	_, err = consumer.RequestRandomWords(owner, 1, 1, 1, 1, [32]byte{}, false)
	require.Error(t, err)
	_, err = consumer.RequestRandomWords(random, 1, 1, 1, 1, [32]byte{}, false)
	require.NoError(t, err)
}

func TestVRFV2PlusIntegration_SimpleConsumerExample(t *testing.T) {
	owner := testutils.MustNewSimTransactor(t)
	random := testutils.MustNewSimTransactor(t)
	genesisData := core.GenesisAlloc{
		owner.From: {Balance: assets.Ether(10).ToInt()},
	}
	backend := cltest.NewSimulatedBackend(t, genesisData, uint32(ethconfig.Defaults.Miner.GasCeil))
	linkAddress, _, linkContract, err := link_token_interface.DeployLinkToken(
		owner, backend)
	require.NoError(t, err)
	backend.Commit()
	coordinatorAddress, _, coordinator, err :=
		vrf_coordinator_v2plus.DeployVRFCoordinatorV2Plus(
			owner, backend, common.Address{}) // bhs not needed for this test
	require.NoError(t, err)
	backend.Commit()
	_, err = coordinator.SetLINK(owner, linkAddress)
	require.NoError(t, err)
	backend.Commit()
	consumerAddress, _, consumer, err := vrf_v2plus_single_consumer.DeployVRFV2PlusSingleConsumerExample(owner, backend, coordinatorAddress, linkAddress, 1, 1, 1, [32]byte{}, false)
	require.NoError(t, err)
	backend.Commit()
	_, err = linkContract.Transfer(owner, consumerAddress, assets.Ether(2).ToInt())
	require.NoError(t, err)
	backend.Commit()
	AssertLinkBalances(t, linkContract, []common.Address{owner.From, consumerAddress}, []*big.Int{assets.Ether(999_999_998).ToInt(), assets.Ether(2).ToInt()})
	_, err = consumer.TopUpSubscription(owner, assets.Ether(1).ToInt())
	require.NoError(t, err)
	backend.Commit()
	AssertLinkBalances(t, linkContract, []common.Address{owner.From, consumerAddress, coordinatorAddress}, []*big.Int{assets.Ether(999_999_998).ToInt(), assets.Ether(1).ToInt(), assets.Ether(1).ToInt()})
	// Non-owner cannot withdraw
	_, err = consumer.Withdraw(random, assets.Ether(1).ToInt(), owner.From)
	require.Error(t, err)
	_, err = consumer.Withdraw(owner, assets.Ether(1).ToInt(), owner.From)
	require.NoError(t, err)
	backend.Commit()
	AssertLinkBalances(t, linkContract, []common.Address{owner.From, consumerAddress, coordinatorAddress}, []*big.Int{assets.Ether(999_999_999).ToInt(), assets.Ether(0).ToInt(), assets.Ether(1).ToInt()})
	_, err = consumer.Unsubscribe(owner, owner.From)
	require.NoError(t, err)
	backend.Commit()
	AssertLinkBalances(t, linkContract, []common.Address{owner.From, consumerAddress, coordinatorAddress}, []*big.Int{assets.Ether(1_000_000_000).ToInt(), assets.Ether(0).ToInt(), assets.Ether(0).ToInt()})
}

func TestVRFV2PlusIntegration_TestMaliciousConsumer(t *testing.T) {
	t.Parallel()
	config, _ := heavyweight.FullTestDBV2(t, "vrf_v2plus_integration_malicious", func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint32](2_000_000)
		c.EVM[0].GasEstimator.PriceMax = assets.GWei(1)
		c.EVM[0].GasEstimator.PriceDefault = assets.GWei(1)
		c.EVM[0].GasEstimator.FeeCapDefault = assets.GWei(1)
	})
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, key, 1)
	carol := uni.vrfConsumers[0]

	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, key)
	require.NoError(t, app.Start(testutils.Context(t)))

	err := app.GetKeyStore().Unlock(cltest.Password)
	require.NoError(t, err)
	vrfkey, err := app.GetKeyStore().VRF().Create()
	require.NoError(t, err)

	jid := uuid.New()
	incomingConfs := 2
	s := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		JobID:                    jid.String(),
		Name:                     "vrf-primary",
		VRFVersion:               vrfcommon.V2Plus,
		FromAddresses:            []string{key.Address.String()},
		CoordinatorAddress:       uni.rootContractAddress.String(),
		BatchCoordinatorAddress:  uni.batchCoordinatorContractAddress.String(),
		MinIncomingConfirmations: incomingConfs,
		GasLanePrice:             assets.GWei(1),
		PublicKey:                vrfkey.PublicKey.String(),
		V2:                       true,
	}).Toml()
	jb, err := vrfcommon.ValidatedVRFSpec(s)
	require.NoError(t, err)
	err = app.JobSpawner().CreateJob(&jb)
	require.NoError(t, err)
	time.Sleep(1 * time.Second)

	// Register a proving key associated with the VRF job.
	p, err := vrfkey.PublicKey.Point()
	require.NoError(t, err)
	_, err = uni.rootContract.RegisterProvingKey(
		uni.neil, uni.nallory.From, pair(secp256k1.Coordinates(p)))
	require.NoError(t, err)

	subFunding := decimal.RequireFromString("1000000000000000000")
	_, err = uni.maliciousConsumerContract.CreateSubscriptionAndFund(carol,
		subFunding.BigInt())
	require.NoError(t, err)
	uni.backend.Commit()

	// Send a re-entrant request
	// subID, nConfs, callbackGas, numWords are hard-coded within the contract, so setting them to 0 here
	_, err = uni.maliciousConsumerContract.RequestRandomness(carol, vrfkey.PublicKey.MustHash(), 0, 0, 0, 0, false)
	require.NoError(t, err)

	// We expect the request to be serviced
	// by the node.
	var attempts []txmgr.TxAttempt
	gomega.NewWithT(t).Eventually(func() bool {
		attempts, _, err = app.TxmStorageService().TxAttempts(0, 1000)
		require.NoError(t, err)
		// It possible that we send the test request
		// before the job spawner has started the vrf services, which is fine
		// the lb will backfill the logs. However we need to
		// keep blocks coming in for the lb to send the backfilled logs.
		t.Log("attempts", attempts)
		uni.backend.Commit()
		return len(attempts) == 1 && attempts[0].Tx.State == txmgrcommon.TxConfirmed
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

	// The fulfillment tx should succeed
	ch, err := app.GetChains().EVM.Default()
	require.NoError(t, err)
	r, err := ch.Client().TransactionReceipt(testutils.Context(t), attempts[0].Hash)
	require.NoError(t, err)
	require.Equal(t, uint64(1), r.Status)

	// The user callback should have errored
	it, err := uni.rootContract.FilterRandomWordsFulfilled(nil, nil)
	require.NoError(t, err)
	var fulfillments []*v22.RandomWordsFulfilled
	for it.Next() {
		fulfillments = append(fulfillments, it.Event())
	}
	require.Equal(t, 1, len(fulfillments))
	require.Equal(t, false, fulfillments[0].Success())

	// It should not have succeeded in placing another request.
	it2, err2 := uni.rootContract.FilterRandomWordsRequested(nil, nil, nil, nil)
	require.NoError(t, err2)
	var requests []*v22.RandomWordsRequested
	for it2.Next() {
		requests = append(requests, it2.Event())
	}
	require.Equal(t, 1, len(requests))
}

func TestVRFV2PlusIntegration_RequestCost(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, key, 1)

	cfg := configtest.NewGeneralConfigSimulated(t, nil)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, cfg, uni.backend, key)
	require.NoError(t, app.Start(testutils.Context(t)))

	vrfkey, err := app.GetKeyStore().VRF().Create()
	require.NoError(t, err)
	p, err := vrfkey.PublicKey.Point()
	require.NoError(t, err)
	_, err = uni.rootContract.RegisterProvingKey(
		uni.neil, uni.neil.From, pair(secp256k1.Coordinates(p)))
	require.NoError(t, err)
	uni.backend.Commit()
	t.Run("non-proxied consumer", func(tt *testing.T) {
		carol := uni.vrfConsumers[0]
		carolContract := uni.consumerContracts[0]
		carolContractAddress := uni.consumerContractAddresses[0]

		_, err = carolContract.CreateSubscriptionAndFund(carol,
			big.NewInt(1000000000000000000)) // 0.1 LINK
		require.NoError(tt, err)
		uni.backend.Commit()
		// Ensure even with large number of consumers its still cheap
		var addrs []common.Address
		for i := 0; i < 99; i++ {
			addrs = append(addrs, testutils.NewAddress())
		}
		_, err = carolContract.UpdateSubscription(carol, addrs)
		require.NoError(tt, err)
		estimate := estimateGas(tt, uni.backend, common.Address{},
			carolContractAddress, uni.consumerABI,
			"requestRandomWords", uint32(10000), uint16(2), uint32(1),
			vrfkey.PublicKey.MustHash(), false)
		tt.Log("gas estimate of non-proxied requestRandomWords:", estimate)
		assert.Less(tt, estimate, uint64(125_000),
			"requestRandomWords tx gas cost more than expected")
	})

	t.Run("proxied consumer", func(tt *testing.T) {
		consumerOwner := uni.neil
		consumerContract := uni.consumerProxyContract
		consumerContractAddress := uni.consumerProxyContractAddress

		// Create a subscription and fund with 5 LINK.
		tx, err := consumerContract.CreateSubscriptionAndFund(consumerOwner, assets.Ether(5).ToInt())
		require.NoError(tt, err)
		uni.backend.Commit()
		r, err := uni.backend.TransactionReceipt(testutils.Context(t), tx.Hash())
		require.NoError(tt, err)
		t.Log("gas used by proxied CreateSubscriptionAndFund:", r.GasUsed)

		subId, err := consumerContract.SSubId(nil)
		require.NoError(tt, err)
		_, err = uni.rootContract.GetSubscription(nil, subId)
		require.NoError(tt, err)

		// Ensure even with large number of consumers it's still cheap
		var addrs []common.Address
		for i := 0; i < 99; i++ {
			addrs = append(addrs, testutils.NewAddress())
		}
		_, err = consumerContract.UpdateSubscription(consumerOwner, addrs)

		theAbi := evmtypes.MustGetABI(vrf_consumer_v2_plus_upgradeable_example.VRFConsumerV2PlusUpgradeableExampleMetaData.ABI)
		estimate := estimateGas(tt, uni.backend, common.Address{},
			consumerContractAddress, &theAbi,
			"requestRandomness", vrfkey.PublicKey.MustHash(), subId, uint16(2), uint32(10000), uint32(1))
		tt.Log("gas estimate of proxied requestRandomness:", estimate)
		// There is some gas overhead of the delegatecall that is made by the proxy
		// to the logic contract. See https://www.evm.codes/#f4?fork=grayGlacier for a detailed
		// breakdown of the gas costs of a delegatecall.
		assert.Less(tt, estimate, uint64(102_000),
			"proxied testRequestRandomness tx gas cost more than expected")
	})
}

func TestVRFV2PlusIntegration_MaxConsumersCost(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, key, 1)
	carol := uni.vrfConsumers[0]
	carolContract := uni.consumerContracts[0]
	carolContractAddress := uni.consumerContractAddresses[0]

	cfg := configtest.NewGeneralConfigSimulated(t, nil)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, cfg, uni.backend, key)
	require.NoError(t, app.Start(testutils.Context(t)))
	_, err := carolContract.CreateSubscriptionAndFund(carol,
		big.NewInt(1000000000000000000)) // 0.1 LINK
	require.NoError(t, err)
	uni.backend.Commit()
	subId, err := carolContract.SSubId(nil)
	require.NoError(t, err)
	var addrs []common.Address
	for i := 0; i < 98; i++ {
		addrs = append(addrs, testutils.NewAddress())
	}
	_, err = carolContract.UpdateSubscription(carol, addrs)
	// Ensure even with max number of consumers its still reasonable gas costs.
	require.NoError(t, err)
	estimate := estimateGas(t, uni.backend, carolContractAddress,
		uni.rootContractAddress, uni.coordinatorABI,
		"removeConsumer", subId, carolContractAddress)
	t.Log(estimate)
	assert.Less(t, estimate, uint64(320000))
	estimate = estimateGas(t, uni.backend, carolContractAddress,
		uni.rootContractAddress, uni.coordinatorABI,
		"addConsumer", subId, testutils.NewAddress())
	t.Log(estimate)
	assert.Less(t, estimate, uint64(100000))
}

func TestVRFV2PlusIntegration_FulfillmentCost(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, key, 1)

	cfg := configtest.NewGeneralConfigSimulated(t, nil)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, cfg, uni.backend, key)
	require.NoError(t, app.Start(testutils.Context(t)))

	vrfkey, err := app.GetKeyStore().VRF().Create()
	require.NoError(t, err)
	p, err := vrfkey.PublicKey.Point()
	require.NoError(t, err)
	_, err = uni.rootContract.RegisterProvingKey(
		uni.neil, uni.neil.From, pair(secp256k1.Coordinates(p)))
	require.NoError(t, err)
	uni.backend.Commit()

	var (
		nonProxiedConsumerGasEstimate uint64
		proxiedConsumerGasEstimate    uint64
	)
	t.Run("non-proxied consumer", func(tt *testing.T) {
		carol := uni.vrfConsumers[0]
		carolContract := uni.consumerContracts[0]
		carolContractAddress := uni.consumerContractAddresses[0]

		_, err = carolContract.CreateSubscriptionAndFund(carol,
			big.NewInt(1000000000000000000)) // 0.1 LINK
		require.NoError(tt, err)
		uni.backend.Commit()
		subId, err := carolContract.SSubId(nil)
		require.NoError(tt, err)

		gasRequested := 50_000
		nw := 1
		requestedIncomingConfs := 3
		_, err = carolContract.RequestRandomness(carol, vrfkey.PublicKey.MustHash(), subId, uint16(requestedIncomingConfs), uint32(gasRequested), uint32(nw), false)
		require.NoError(t, err)
		for i := 0; i < requestedIncomingConfs; i++ {
			uni.backend.Commit()
		}

		requestLog := FindLatestRandomnessRequestedLog(tt, uni.rootContract, vrfkey.PublicKey.MustHash())
		s, err := proof.BigToSeed(requestLog.PreSeed())
		require.NoError(t, err)
		proof, rc, err := proof.GenerateProofResponseV2Plus(app.GetKeyStore().VRF(), vrfkey.ID(), proof.PreSeedDataV2{
			PreSeed:          s,
			BlockHash:        requestLog.Raw().BlockHash,
			BlockNum:         requestLog.Raw().BlockNumber,
			SubId:            subId,
			CallbackGasLimit: uint32(gasRequested),
			NumWords:         uint32(nw),
			Sender:           carolContractAddress,
		}, false)
		require.NoError(tt, err)
		nonProxiedConsumerGasEstimate = estimateGas(tt, uni.backend, common.Address{},
			uni.rootContractAddress, uni.coordinatorABI,
			"fulfillRandomWords", proof, rc)
		t.Log("non-proxied consumer fulfillment gas estimate:", nonProxiedConsumerGasEstimate)
		// Establish very rough bounds on fulfillment cost
		assert.Greater(tt, nonProxiedConsumerGasEstimate, uint64(120_000))
		assert.Less(tt, nonProxiedConsumerGasEstimate, uint64(500_000))
	})

	t.Run("proxied consumer", func(tt *testing.T) {
		consumerOwner := uni.neil
		consumerContract := uni.consumerProxyContract
		consumerContractAddress := uni.consumerProxyContractAddress

		_, err = consumerContract.CreateSubscriptionAndFund(consumerOwner, assets.Ether(5).ToInt())
		require.NoError(t, err)
		uni.backend.Commit()
		subId, err := consumerContract.SSubId(nil)
		require.NoError(t, err)
		gasRequested := 50_000
		nw := 1
		requestedIncomingConfs := 3
		_, err = consumerContract.RequestRandomness(consumerOwner, vrfkey.PublicKey.MustHash(), subId, uint16(requestedIncomingConfs), uint32(gasRequested), uint32(nw), false)
		require.NoError(t, err)
		for i := 0; i < requestedIncomingConfs; i++ {
			uni.backend.Commit()
		}

		requestLog := FindLatestRandomnessRequestedLog(t, uni.rootContract, vrfkey.PublicKey.MustHash())
		require.Equal(tt, subId, requestLog.SubID())
		s, err := proof.BigToSeed(requestLog.PreSeed())
		require.NoError(t, err)
		proof, rc, err := proof.GenerateProofResponseV2Plus(app.GetKeyStore().VRF(), vrfkey.ID(), proof.PreSeedDataV2{
			PreSeed:          s,
			BlockHash:        requestLog.Raw().BlockHash,
			BlockNum:         requestLog.Raw().BlockNumber,
			SubId:            subId,
			CallbackGasLimit: uint32(gasRequested),
			NumWords:         uint32(nw),
			Sender:           consumerContractAddress,
		}, false)
		require.NoError(t, err)
		proxiedConsumerGasEstimate = estimateGas(t, uni.backend, common.Address{},
			uni.rootContractAddress, uni.coordinatorABI,
			"fulfillRandomWords", proof, rc)
		t.Log("proxied consumer fulfillment gas estimate", proxiedConsumerGasEstimate)
		// Establish very rough bounds on fulfillment cost
		assert.Greater(t, proxiedConsumerGasEstimate, uint64(120_000))
		assert.Less(t, proxiedConsumerGasEstimate, uint64(500_000))
	})
}

func TestVRFV2PlusIntegration_StartingCountsV1(t *testing.T) {
	cfg, db := heavyweight.FullTestDBNoFixturesV2(t, "vrf_test_starting_counts", nil)
	_, err := db.Exec(`INSERT INTO evm_chains (id, created_at, updated_at) VALUES (1337, NOW(), NOW())`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO evm_heads (hash, number, parent_hash, created_at, timestamp, evm_chain_id)
	VALUES ($1, 4, $2, NOW(), NOW(), 1337)`, utils.NewHash(), utils.NewHash())
	require.NoError(t, err)

	lggr := logger.TestLogger(t)
	q := pg.NewQ(db, lggr, cfg.Database())
	finalityDepth := 3
	counts := vrf.GetStartingResponseCountsV1(q, lggr, 1337, uint32(finalityDepth))
	assert.Equal(t, 0, len(counts))
	ks := keystore.New(db, utils.FastScryptParams, lggr, cfg.Database())
	err = ks.Unlock(testutils.Password)
	require.NoError(t, err)
	k, err := ks.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)
	b := time.Now()
	n1, n2, n3, n4 := evmtypes.Nonce(0), evmtypes.Nonce(1), evmtypes.Nonce(2), evmtypes.Nonce(3)
	reqID := utils.PadByteToHash(0x10)
	m1 := txmgr.TxMeta{
		RequestID: &reqID,
	}
	md1, err := json.Marshal(&m1)
	require.NoError(t, err)
	md1_ := datatypes.JSON(md1)
	reqID2 := utils.PadByteToHash(0x11)
	m2 := txmgr.TxMeta{
		RequestID: &reqID2,
	}
	md2, err := json.Marshal(&m2)
	md2_ := datatypes.JSON(md2)
	require.NoError(t, err)
	chainID := utils.NewBig(big.NewInt(1337))
	confirmedTxes := []txmgr.Tx{
		{
			Sequence:           &n1,
			FromAddress:        k.Address,
			Error:              null.String{},
			BroadcastAt:        &b,
			InitialBroadcastAt: &b,
			CreatedAt:          b,
			State:              txmgrcommon.TxConfirmed,
			Meta:               &datatypes.JSON{},
			EncodedPayload:     []byte{},
			ChainID:            chainID.ToInt(),
		},
		{
			Sequence:           &n2,
			FromAddress:        k.Address,
			Error:              null.String{},
			BroadcastAt:        &b,
			InitialBroadcastAt: &b,
			CreatedAt:          b,
			State:              txmgrcommon.TxConfirmed,
			Meta:               &md1_,
			EncodedPayload:     []byte{},
			ChainID:            chainID.ToInt(),
		},
		{
			Sequence:           &n3,
			FromAddress:        k.Address,
			Error:              null.String{},
			BroadcastAt:        &b,
			InitialBroadcastAt: &b,
			CreatedAt:          b,
			State:              txmgrcommon.TxConfirmed,
			Meta:               &md2_,
			EncodedPayload:     []byte{},
			ChainID:            chainID.ToInt(),
		},
		{
			Sequence:           &n4,
			FromAddress:        k.Address,
			Error:              null.String{},
			BroadcastAt:        &b,
			InitialBroadcastAt: &b,
			CreatedAt:          b,
			State:              txmgrcommon.TxConfirmed,
			Meta:               &md2_,
			EncodedPayload:     []byte{},
			ChainID:            chainID.ToInt(),
		},
	}
	// add unconfirmed txes
	unconfirmedTxes := []txmgr.Tx{}
	for i := int64(4); i < 6; i++ {
		reqID3 := utils.PadByteToHash(0x12)
		md, err := json.Marshal(&txmgr.TxMeta{
			RequestID: &reqID3,
		})
		require.NoError(t, err)
		md1 := datatypes.JSON(md)
		newNonce := evmtypes.Nonce(i + 1)
		unconfirmedTxes = append(unconfirmedTxes, txmgr.Tx{
			Sequence:           &newNonce,
			FromAddress:        k.Address,
			Error:              null.String{},
			CreatedAt:          b,
			State:              txmgrcommon.TxUnconfirmed,
			BroadcastAt:        &b,
			InitialBroadcastAt: &b,
			Meta:               &md1,
			EncodedPayload:     []byte{},
			ChainID:            chainID.ToInt(),
		})
	}
	txes := append(confirmedTxes, unconfirmedTxes...)
	sql := `INSERT INTO eth_txes (nonce, from_address, to_address, encoded_payload, value, gas_limit, state, created_at, broadcast_at, initial_broadcast_at, meta, subject, evm_chain_id, min_confirmations, pipeline_task_run_id)
VALUES (:nonce, :from_address, :to_address, :encoded_payload, :value, :gas_limit, :state, :created_at, :broadcast_at, :initial_broadcast_at, :meta, :subject, :evm_chain_id, :min_confirmations, :pipeline_task_run_id);`
	for _, tx := range txes {
		dbEtx := txmgr.DbEthTxFromEthTx(&tx)
		_, err = db.NamedExec(sql, &dbEtx)
		txmgr.DbEthTxToEthTx(dbEtx, &tx)
		require.NoError(t, err)
	}

	// add eth_tx_attempts for confirmed
	broadcastBlock := int64(1)
	txAttempts := []txmgr.TxAttempt{}
	for i := range confirmedTxes {
		txAttempts = append(txAttempts, txmgr.TxAttempt{
			TxID:                    int64(i + 1),
			TxFee:                   gas.EvmFee{Legacy: assets.NewWeiI(100)},
			SignedRawTx:             []byte(`blah`),
			Hash:                    utils.NewHash(),
			BroadcastBeforeBlockNum: &broadcastBlock,
			State:                   txmgrtypes.TxAttemptBroadcast,
			CreatedAt:               time.Now(),
			ChainSpecificFeeLimit:   uint32(100),
		})
	}
	// add eth_tx_attempts for unconfirmed
	for i := range unconfirmedTxes {
		txAttempts = append(txAttempts, txmgr.TxAttempt{
			TxID:                  int64(i + 1 + len(confirmedTxes)),
			TxFee:                 gas.EvmFee{Legacy: assets.NewWeiI(100)},
			SignedRawTx:           []byte(`blah`),
			Hash:                  utils.NewHash(),
			State:                 txmgrtypes.TxAttemptInProgress,
			CreatedAt:             time.Now(),
			ChainSpecificFeeLimit: uint32(100),
		})
	}
	for _, txAttempt := range txAttempts {
		t.Log("tx attempt eth tx id: ", txAttempt.TxID)
	}
	sql = `INSERT INTO eth_tx_attempts (eth_tx_id, gas_price, signed_raw_tx, hash, state, created_at, chain_specific_gas_limit)
		VALUES (:eth_tx_id, :gas_price, :signed_raw_tx, :hash, :state, :created_at, :chain_specific_gas_limit)`
	for _, attempt := range txAttempts {
		dbAttempt := txmgr.DbEthTxAttemptFromEthTxAttempt(&attempt)
		_, err = db.NamedExec(sql, &dbAttempt)
		txmgr.DbEthTxAttemptToEthTxAttempt(dbAttempt, &attempt)
		require.NoError(t, err)
	}

	// add eth_receipts
	receipts := []txmgr.Receipt{}
	for i := 0; i < 4; i++ {
		receipts = append(receipts, txmgr.Receipt{
			BlockHash:        utils.NewHash(),
			TxHash:           txAttempts[i].Hash,
			BlockNumber:      broadcastBlock,
			TransactionIndex: 1,
			Receipt:          evmtypes.Receipt{},
			CreatedAt:        time.Now(),
		})
	}
	sql = `INSERT INTO eth_receipts (block_hash, tx_hash, block_number, transaction_index, receipt, created_at)
		VALUES (:block_hash, :tx_hash, :block_number, :transaction_index, :receipt, :created_at)`
	for _, r := range receipts {
		_, err := db.NamedExec(sql, &r)
		require.NoError(t, err)
	}

	counts = vrf.GetStartingResponseCountsV1(q, lggr, 1337, uint32(finalityDepth))
	assert.Equal(t, 3, len(counts))
	assert.Equal(t, uint64(1), counts[utils.PadByteToHash(0x10)])
	assert.Equal(t, uint64(2), counts[utils.PadByteToHash(0x11)])
	assert.Equal(t, uint64(2), counts[utils.PadByteToHash(0x12)])

	countsV2 := vrf.GetStartingResponseCountsV2(q, lggr, 1337, uint32(finalityDepth))
	t.Log(countsV2)
	assert.Equal(t, 3, len(countsV2))
	assert.Equal(t, uint64(1), countsV2[big.NewInt(0x10).String()])
	assert.Equal(t, uint64(2), countsV2[big.NewInt(0x11).String()])
	assert.Equal(t, uint64(2), countsV2[big.NewInt(0x12).String()])
}

func AssertEthBalances(t *testing.T, backend *backends.SimulatedBackend, addresses []common.Address, balances []*big.Int) {
	require.Equal(t, len(addresses), len(balances))
	for i, a := range addresses {
		b, err := backend.BalanceAt(testutils.Context(t), a, nil)
		require.NoError(t, err)
		assert.Equal(t, balances[i].String(), b.String(), "invalid balance for %v", a)
	}
}
