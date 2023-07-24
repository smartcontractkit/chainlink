package v2_test

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
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
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/proof"
	v22 "github.com/smartcontractkit/chainlink/v2/core/services/vrf/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrftesthelpers"
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
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1)
	testSingleConsumerEIP150(
		t,
		ownerKey,
		uni.coordinatorV2UniverseCommon,
		uni.batchCoordinatorContractAddress,
		false,
		vrfcommon.V2Plus,
	)
}

func TestVRFV2PlusIntegration_SingleConsumer_EIP150_Revert(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1)
	testSingleConsumerEIP150Revert(
		t,
		ownerKey,
		uni.coordinatorV2UniverseCommon,
		uni.batchCoordinatorContractAddress,
		false,
		vrfcommon.V2Plus,
	)
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
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1)
	testSingleConsumerBigGasCallbackSandwich(
		t,
		ownerKey,
		uni.coordinatorV2UniverseCommon,
		uni.batchCoordinatorContractAddress,
		false,
		vrfcommon.V2Plus,
	)
}

func TestVRFV2PlusIntegration_SingleConsumer_MultipleGasLanes(t *testing.T) {
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1)
	testSingleConsumerMultipleGasLanes(
		t,
		ownerKey,
		uni.coordinatorV2UniverseCommon,
		uni.batchCoordinatorContractAddress,
		false,
		vrfcommon.V2Plus,
	)
}

func TestVRFV2PlusIntegration_SingleConsumer_AlwaysRevertingCallback_StillFulfilled(t *testing.T) {
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 0)
	testSingleConsumerAlwaysRevertingCallbackStillFulfilled(
		t,
		ownerKey,
		uni.coordinatorV2UniverseCommon,
		uni.batchCoordinatorContractAddress,
		false,
		vrfcommon.V2Plus,
	)
}

func TestVRFV2PlusIntegration_ConsumerProxy_HappyPath(t *testing.T) {
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 0)
	testConsumerProxyHappyPath(
		t,
		ownerKey,
		uni.coordinatorV2UniverseCommon,
		uni.batchCoordinatorContractAddress,
		false,
		vrfcommon.V2Plus,
	)
}

func TestVRFV2PlusIntegration_ConsumerProxy_CoordinatorZeroAddress(t *testing.T) {
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 0)
	testConsumerProxyCoordinatorZeroAddress(t, uni.coordinatorV2UniverseCommon)
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
	_, err = consumer.AcceptOwnership(random)
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
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, key, 1)
	testMaliciousConsumer(
		t,
		key,
		uni.coordinatorV2UniverseCommon,
		uni.batchCoordinatorContractAddress,
		false,
		vrfcommon.V2Plus,
	)
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
		assert.Less(tt, estimate, uint64(126_000),
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

		theAbi := evmtypes.MustGetABI(vrf_consumer_v2_plus_upgradeable_example.VRFConsumerV2PlusUpgradeableExampleMetaData.ABI)
		estimate := estimateGas(tt, uni.backend, common.Address{},
			consumerContractAddress, &theAbi,
			"requestRandomness", vrfkey.PublicKey.MustHash(), subId, uint16(2), uint32(10000), uint32(1))
		tt.Log("gas estimate of proxied requestRandomness:", estimate)
		// There is some gas overhead of the delegatecall that is made by the proxy
		// to the logic contract. See https://www.evm.codes/#f4?fork=grayGlacier for a detailed
		// breakdown of the gas costs of a delegatecall.
		assert.Less(tt, estimate, uint64(105_000),
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

func AssertEthBalances(t *testing.T, backend *backends.SimulatedBackend, addresses []common.Address, balances []*big.Int) {
	require.Equal(t, len(addresses), len(balances))
	for i, a := range addresses {
		b, err := backend.BalanceAt(testutils.Context(t), a, nil)
		require.NoError(t, err)
		assert.Equal(t, balances[i].String(), b.String(), "invalid balance for %v", a)
	}
}
