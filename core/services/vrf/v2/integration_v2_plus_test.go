package v2_test

import (
	"math"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/batch_vrf_coordinator_v2plus"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/trusted_blockhash_store"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_consumer_v2_plus_upgradeable_example"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_consumer_v2_upgradeable_example"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2_plus_v2_example"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2plus_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_malicious_consumer_v2_plus"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_v2plus_single_consumer"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_v2plus_sub_owner"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrfv2_proxy_admin"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrfv2_transparent_upgradeable_proxy"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrfv2plus_consumer_example"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrfv2plus_reverting_example"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/extraargs"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/proof"
	v22 "github.com/smartcontractkit/chainlink/v2/core/services/vrf/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrftesthelpers"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

type coordinatorV2PlusUniverse struct {
	coordinatorV2UniverseCommon
	submanager                      *bind.TransactOpts // Subscription owner
	batchCoordinatorContract        *batch_vrf_coordinator_v2plus.BatchVRFCoordinatorV2Plus
	batchCoordinatorContractAddress common.Address
	migrationTestCoordinator        *vrf_coordinator_v2_plus_v2_example.VRFCoordinatorV2PlusV2Example
	migrationTestCoordinatorAddress common.Address
	trustedBhsContract              *trusted_blockhash_store.TrustedBlockhashStore
	trustedBhsContractAddress       common.Address
}

func newVRFCoordinatorV2PlusUniverse(t *testing.T, key ethkey.KeyV2, numConsumers int, trusting bool) coordinatorV2PlusUniverse {
	tests.SkipShort(t, "VRFCoordinatorV2Universe")
	oracleTransactor := &bind.TransactOpts{
		From:   key.Address,
		Signer: key.SignerFn(testutils.SimulatedChainID),
	}
	var (
		sergey       = evmtestutils.MustNewSimTransactor(t)
		neil         = evmtestutils.MustNewSimTransactor(t)
		ned          = evmtestutils.MustNewSimTransactor(t)
		evil         = evmtestutils.MustNewSimTransactor(t)
		reverter     = evmtestutils.MustNewSimTransactor(t)
		submanager   = evmtestutils.MustNewSimTransactor(t)
		nallory      = oracleTransactor
		vrfConsumers []*bind.TransactOpts
	)

	// Create consumer contract deployer identities
	for i := 0; i < numConsumers; i++ {
		vrfConsumers = append(vrfConsumers, evmtestutils.MustNewSimTransactor(t))
	}

	genesisData := gethtypes.GenesisAlloc{
		sergey.From:     {Balance: assets.Ether(1000).ToInt()},
		neil.From:       {Balance: assets.Ether(1000).ToInt()},
		ned.From:        {Balance: assets.Ether(1000).ToInt()},
		nallory.From:    {Balance: assets.Ether(1000).ToInt()},
		evil.From:       {Balance: assets.Ether(1000).ToInt()},
		reverter.From:   {Balance: assets.Ether(1000).ToInt()},
		submanager.From: {Balance: assets.Ether(1000).ToInt()},
		// ATTENTION this is needed because VRF simulations
		// simulate from the 0x0 address, i.e. with From unspecified.
		// On real chains, that seems to not require a balance, but
		// the sim does. You'll see this otherwise
		// insufficient funds for gas * price + value: address 0x0000000000000000000000000000000000000000
		common.HexToAddress("0x0"): {Balance: assets.Ether(1000).ToInt()},
	}
	for _, consumer := range vrfConsumers {
		genesisData[consumer.From] = gethtypes.Account{
			Balance: assets.Ether(1000).ToInt(),
		}
	}

	consumerABI, err := abi.JSON(strings.NewReader(
		vrfv2plus_consumer_example.VRFV2PlusConsumerExampleABI))
	require.NoError(t, err)
	coordinatorABI, err := abi.JSON(strings.NewReader(
		vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalABI))
	require.NoError(t, err)
	backend := cltest.NewSimulatedBackend(t, genesisData, ethconfig.Defaults.Miner.GasCeil)
	h, err := backend.Client().HeaderByNumber(testutils.Context(t), nil)
	require.NoError(t, err)
	require.LessOrEqual(t, h.Time, uint64(math.MaxInt64))
	blockTime := time.Unix(int64(h.Time), 0) //nolint:gosec // G115 false positive
	// Move the clock closer to the current time. We set first block to be 24 hours ago.
	err = backend.AdjustTime(time.Since(blockTime) - 24*time.Hour)
	require.NoError(t, err)
	backend.Commit()

	// Deploy link
	linkAddress, _, linkContract, err := link_token_interface.DeployLinkToken(
		sergey, backend.Client())
	require.NoError(t, err, "failed to deploy link contract to simulated ethereum blockchain")
	backend.Commit()

	// Deploy feed
	linkEthFeed, _, _, err :=
		mock_v3_aggregator_contract.DeployMockV3AggregatorContract(
			evil, backend.Client(), 18, vrftesthelpers.WeiPerUnitLink.BigInt()) // 0.01 eth per link
	require.NoError(t, err)
	backend.Commit()

	// Deploy blockhash store
	bhsAddress, _, bhsContract, err := blockhash_store.DeployBlockhashStore(neil, backend.Client())
	require.NoError(t, err, "failed to deploy BlockhashStore contract to simulated ethereum blockchain")
	backend.Commit()

	// Deploy trusted BHS
	trustedBHSAddress, _, trustedBhsContract, err := trusted_blockhash_store.DeployTrustedBlockhashStore(neil, backend.Client(), []common.Address{})
	require.NoError(t, err, "failed to deploy trusted BlockhashStore contract to simulated ethereum blockchain")
	backend.Commit()

	// Deploy batch blockhash store
	batchBHSAddress, _, batchBHSContract, err := batch_blockhash_store.DeployBatchBlockhashStore(neil, backend.Client(), bhsAddress)
	require.NoError(t, err, "failed to deploy BatchBlockhashStore contract to simulated ethereum blockchain")
	backend.Commit()

	// Deploy VRF V2plus coordinator
	var bhsAddr = bhsAddress
	if trusting {
		bhsAddr = trustedBHSAddress
	}
	coordinatorAddress, _, coordinatorContract, err :=
		vrf_coordinator_v2_5.DeployVRFCoordinatorV25(
			neil, backend.Client(), bhsAddr)
	require.NoError(t, err, "failed to deploy VRFCoordinatorV2 contract to simulated ethereum blockchain")
	backend.Commit()

	_, err = coordinatorContract.SetLINKAndLINKNativeFeed(neil, linkAddress, linkEthFeed)
	require.NoError(t, err)
	backend.Commit()

	migrationTestCoordinatorAddress, _, migrationTestCoordinator, err := vrf_coordinator_v2_plus_v2_example.DeployVRFCoordinatorV2PlusV2Example(
		neil, backend.Client(), linkAddress, coordinatorAddress)
	require.NoError(t, err)
	backend.Commit()

	_, err = coordinatorContract.RegisterMigratableCoordinator(neil, migrationTestCoordinatorAddress)
	require.NoError(t, err)
	backend.Commit()

	// Deploy batch VRF V2 coordinator
	batchCoordinatorAddress, _, batchCoordinatorContract, err :=
		batch_vrf_coordinator_v2plus.DeployBatchVRFCoordinatorV2Plus(
			neil, backend.Client(), coordinatorAddress,
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
		consumerContractAddress, _, consumerContract, err2 :=
			vrfv2plus_consumer_example.DeployVRFV2PlusConsumerExample(
				author, backend.Client(), coordinatorAddress, linkAddress)
		require.NoError(t, err2, "failed to deploy VRFConsumer contract to simulated ethereum blockchain")
		backend.Commit()
		_, err2 = linkContract.Transfer(sergey, consumerContractAddress, assets.Ether(500).ToInt()) // Actually, LINK
		require.NoError(t, err2, "failed to send LINK to VRFConsumer contract on simulated ethereum blockchain")
		backend.Commit()

		consumerContracts = append(consumerContracts, vrftesthelpers.NewVRFV2PlusConsumer(consumerContract))
		consumerContractAddresses = append(consumerContractAddresses, consumerContractAddress)

		backend.Commit()
	}

	// Deploy malicious consumer with 1 link
	maliciousConsumerContractAddress, _, maliciousConsumerContract, err :=
		vrf_malicious_consumer_v2_plus.DeployVRFMaliciousConsumerV2Plus(
			evil, backend.Client(), coordinatorAddress, linkAddress)
	require.NoError(t, err, "failed to deploy VRFMaliciousConsumer contract to simulated ethereum blockchain")
	backend.Commit()
	_, err = linkContract.Transfer(sergey, maliciousConsumerContractAddress, assets.Ether(1).ToInt()) // Actually, LINK
	require.NoError(t, err, "failed to send LINK to VRFMaliciousConsumer contract on simulated ethereum blockchain")
	backend.Commit()

	// Deploy upgradeable consumer, proxy, and proxy admin
	upgradeableConsumerAddress, _, _, err := vrf_consumer_v2_plus_upgradeable_example.DeployVRFConsumerV2PlusUpgradeableExample(neil, backend.Client())
	require.NoError(t, err, "failed to deploy upgradeable consumer to simulated ethereum blockchain")
	backend.Commit()

	proxyAdminAddress, _, proxyAdmin, err := vrfv2_proxy_admin.DeployVRFV2ProxyAdmin(neil, backend.Client())
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
		neil, backend.Client(), upgradeableConsumerAddress, proxyAdminAddress, initializeCalldata)
	require.NoError(t, err)
	backend.Commit()

	_, err = linkContract.Transfer(sergey, proxyAddress, assets.Ether(500).ToInt()) // Actually, LINK
	require.NoError(t, err)
	backend.Commit()

	implAddress, err := proxyAdmin.GetProxyImplementation(nil, proxyAddress)
	require.NoError(t, err)
	t.Log("impl address:", implAddress.String())
	require.Equal(t, upgradeableConsumerAddress, implAddress)

	proxiedConsumer, err := vrf_consumer_v2_plus_upgradeable_example.NewVRFConsumerV2PlusUpgradeableExample(
		proxyAddress, backend.Client())
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
		reverter, backend.Client(), coordinatorAddress, linkAddress,
	)
	require.NoError(t, err, "failed to deploy VRFRevertingExample contract to simulated eth blockchain")
	backend.Commit()
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
		uint32(5),                              // 0.000005 ETH premium
		uint32(1),                              // 0.000001 LINK premium discount denominated in ETH
		uint8(10),                              // 10% native payment percentage
		uint8(5),                               // 5% LINK payment percentage
	)
	require.NoError(t, err, "failed to set coordinator configuration")
	backend.Commit()

	for i := 0; i < 200; i++ {
		backend.Commit()
	}

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

			rootContract:                     v22.NewCoordinatorV2_5(coordinatorContract),
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
		migrationTestCoordinator:        migrationTestCoordinator,
		migrationTestCoordinatorAddress: migrationTestCoordinatorAddress,
		trustedBhsContract:              trustedBhsContract,
		trustedBhsContractAddress:       trustedBHSAddress,
	}
}

func TestVRFV2PlusIntegration_SingleConsumer_HappyPath_BatchFulfillment(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1, false)
	t.Run("link payment", func(tt *testing.T) {
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
			false,
			func(t *testing.T, coordinator v22.CoordinatorV2_X, rwfe v22.RandomWordsFulfilled, expectedSubID *big.Int) {
				_, err := coordinator.GetSubscription(nil, rwfe.SubID())
				require.NoError(t, err)
				require.Equal(t, expectedSubID, rwfe.SubID())
			},
		)
	})

	t.Run("native payment", func(tt *testing.T) {
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
			true,
			func(t *testing.T, coordinator v22.CoordinatorV2_X, rwfe v22.RandomWordsFulfilled, expectedSubID *big.Int) {
				_, err := coordinator.GetSubscription(nil, rwfe.SubID())
				require.NoError(t, err)
				require.Equal(t, expectedSubID, rwfe.SubID())
			},
		)
	})
}

func TestVRFV2PlusIntegration_SingleConsumer_HappyPath_BatchFulfillment_BigGasCallback(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1, false)
	t.Run("link payment", func(tt *testing.T) {
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
			false,
			func(t *testing.T, coordinator v22.CoordinatorV2_X, rwfe v22.RandomWordsFulfilled, expectedSubID *big.Int) {
				_, err := coordinator.GetSubscription(nil, rwfe.SubID())
				require.NoError(t, err)
				require.Equal(t, expectedSubID, rwfe.SubID())
			},
		)
	})

	t.Run("native payment", func(tt *testing.T) {
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
			true,
			func(t *testing.T, coordinator v22.CoordinatorV2_X, rwfe v22.RandomWordsFulfilled, expectedSubID *big.Int) {
				_, err := coordinator.GetSubscription(nil, rwfe.SubID())
				require.NoError(t, err)
				require.Equal(t, expectedSubID, rwfe.SubID())
			},
		)
	})
}

func TestVRFV2PlusIntegration_SingleConsumer_HappyPath(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1, false)
	t.Run("link payment", func(tt *testing.T) {
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
			vrfcommon.V2Plus,
			false,
			func(t *testing.T, coordinator v22.CoordinatorV2_X, rwfe v22.RandomWordsFulfilled, expectedSubID *big.Int) {
				_, err := coordinator.GetSubscription(nil, rwfe.SubID())
				require.NoError(t, err)
				require.Equal(t, expectedSubID, rwfe.SubID())
			})
	})
	t.Run("native payment", func(tt *testing.T) {
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
			vrfcommon.V2Plus,
			true,
			func(t *testing.T, coordinator v22.CoordinatorV2_X, rwfe v22.RandomWordsFulfilled, expectedSubID *big.Int) {
				_, err := coordinator.GetSubscription(nil, rwfe.SubID())
				require.NoError(t, err)
				require.Equal(t, expectedSubID, rwfe.SubID())
			})
	})
}

func TestVRFV2PlusIntegration_SingleConsumer_EOA_Request(t *testing.T) {
	t.Skip("questionable value of this test")
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1, false)
	testEoa(
		t,
		ownerKey,
		uni.coordinatorV2UniverseCommon,
		false,
		uni.batchBHSContractAddress,
		nil,
		vrfcommon.V2Plus,
	)
}

func TestVRFV2PlusIntegration_SingleConsumer_EOA_Request_Batching_Enabled(t *testing.T) {
	t.Skip("questionable value of this test")
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1, false)
	testEoa(
		t,
		ownerKey,
		uni.coordinatorV2UniverseCommon,
		true,
		uni.batchBHSContractAddress,
		nil,
		vrfcommon.V2Plus,
	)
}

func TestVRFV2PlusIntegration_SingleConsumer_EIP150_HappyPath(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1, false)
	testSingleConsumerEIP150(
		t,
		ownerKey,
		uni.coordinatorV2UniverseCommon,
		uni.batchCoordinatorContractAddress,
		false,
		vrfcommon.V2Plus,
		false,
	)
}

func TestVRFV2PlusIntegration_SingleConsumer_EIP150_Revert(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1, false)
	testSingleConsumerEIP150Revert(
		t,
		ownerKey,
		uni.coordinatorV2UniverseCommon,
		uni.batchCoordinatorContractAddress,
		false,
		vrfcommon.V2Plus,
		false,
	)
}

func TestVRFV2PlusIntegration_SingleConsumer_NeedsBlockhashStore(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 2, false)
	t.Run("link payment", func(tt *testing.T) {
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
			vrfcommon.V2Plus,
			false,
		)
	})
	t.Run("native payment", func(tt *testing.T) {
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
			vrfcommon.V2Plus,
			true,
		)
	})
}

func TestVRFV2PlusIntegration_SingleConsumer_BlockHeaderFeeder(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1, false)
	t.Run("link payment", func(tt *testing.T) {
		tt.Skip("fails after geth upgrade https://github.com/smartcontractkit/chainlink/pull/11809")
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
			vrfcommon.V2Plus,
			false,
		)
	})
	t.Run("native payment", func(tt *testing.T) {
		tt.Skip("fails after geth upgrade https://github.com/smartcontractkit/chainlink/pull/11809")
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
			vrfcommon.V2Plus,
			true,
		)
	})
}

func TestVRFV2PlusIntegration_SingleConsumer_NeedsTopUp(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1, false)
	t.Run("link payment", func(tt *testing.T) {
		tt.Skip("fails after geth upgrade https://github.com/smartcontractkit/chainlink/pull/11809")
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
			false,
		)
	})
	t.Run("native payment", func(tt *testing.T) {
		tt.Skip("fails after geth upgrade https://github.com/smartcontractkit/chainlink/pull/11809")
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
			big.NewInt(1e17),          // initial funding of 0.1 ETH
			assets.Ether(100).ToInt(), // top up of 100 ETH
			vrfcommon.V2Plus,
			true,
		)
	})
}

func TestVRFV2PlusIntegration_SingleConsumer_BigGasCallback_Sandwich(t *testing.T) {
	t.Skip("fails after geth upgrade https://github.com/smartcontractkit/chainlink/pull/11809")
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1, false)
	testSingleConsumerBigGasCallbackSandwich(t, ownerKey, uni.coordinatorV2UniverseCommon, uni.batchCoordinatorContractAddress, vrfcommon.V2Plus, false)
}

func TestVRFV2PlusIntegration_SingleConsumer_MultipleGasLanes(t *testing.T) {
	t.Skip("fails after geth upgrade https://github.com/smartcontractkit/chainlink/pull/11809")
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1, false)
	testSingleConsumerMultipleGasLanes(t, ownerKey, uni.coordinatorV2UniverseCommon, uni.batchCoordinatorContractAddress, vrfcommon.V2Plus, false)
}

func TestVRFV2PlusIntegration_SingleConsumer_AlwaysRevertingCallback_StillFulfilled(t *testing.T) {
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 0, false)
	testSingleConsumerAlwaysRevertingCallbackStillFulfilled(
		t,
		ownerKey,
		uni.coordinatorV2UniverseCommon,
		uni.batchCoordinatorContractAddress,
		false,
		vrfcommon.V2Plus,
		false,
	)
}

func TestVRFV2PlusIntegration_ConsumerProxy_HappyPath(t *testing.T) {
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 0, false)
	testConsumerProxyHappyPath(
		t,
		ownerKey,
		uni.coordinatorV2UniverseCommon,
		uni.batchCoordinatorContractAddress,
		vrfcommon.V2Plus,
		false,
	)
}

func TestVRFV2PlusIntegration_ConsumerProxy_CoordinatorZeroAddress(t *testing.T) {
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 0, false)
	testConsumerProxyCoordinatorZeroAddress(t, uni.coordinatorV2UniverseCommon)
}

func TestVRFV2PlusIntegration_ExternalOwnerConsumerExample(t *testing.T) {
	owner := evmtestutils.MustNewSimTransactor(t)
	random := evmtestutils.MustNewSimTransactor(t)
	genesisData := gethtypes.GenesisAlloc{
		owner.From:  {Balance: assets.Ether(10).ToInt()},
		random.From: {Balance: assets.Ether(10).ToInt()},
	}
	backend := cltest.NewSimulatedBackend(t, genesisData, ethconfig.Defaults.Miner.GasCeil)
	linkAddress, _, linkContract, err := link_token_interface.DeployLinkToken(
		owner, backend.Client())
	require.NoError(t, err)
	backend.Commit()
	// Deploy feed
	linkEthFeed, _, _, err :=
		mock_v3_aggregator_contract.DeployMockV3AggregatorContract(
			owner, backend.Client(), 18, vrftesthelpers.WeiPerUnitLink.BigInt()) // 0.01 eth per link
	require.NoError(t, err)
	backend.Commit()
	coordinatorAddress, _, coordinator, err :=
		vrf_coordinator_v2_5.DeployVRFCoordinatorV25(
			owner, backend.Client(), common.Address{}) // bhs not needed for this test
	require.NoError(t, err)
	backend.Commit()
	_, err = coordinator.SetConfig(owner,
		uint16(1),      // minimumRequestConfirmations
		uint32(10000),  // maxGasLimit
		1,              // stalenessSeconds
		1,              // gasAfterPaymentCalculation
		big.NewInt(10), // fallbackWeiPerUnitLink
		0,              // fulfillmentFlatFeeNativePPM
		0,              // fulfillmentFlatFeeLinkDiscountPPM
		0,              // nativePremiumPercentage
		0,              // linkPremiumPercentage
	)
	require.NoError(t, err)
	backend.Commit()
	_, err = coordinator.SetLINKAndLINKNativeFeed(owner, linkAddress, linkEthFeed)
	require.NoError(t, err)
	backend.Commit()
	consumerAddress, _, consumer, err := vrf_v2plus_sub_owner.DeployVRFV2PlusExternalSubOwnerExample(owner, backend.Client(), coordinatorAddress, linkAddress)
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

	iter, err := coordinator.FilterSubscriptionCreated(nil, nil)
	require.NoError(t, err)
	require.True(t, iter.Next(), "could not find SubscriptionCreated event for subID")
	subID := iter.Event.SubId

	b, err := utils.ABIEncode(`[{"type":"uint256"}]`, subID)
	require.NoError(t, err)
	_, err = linkContract.TransferAndCall(owner, coordinatorAddress, big.NewInt(0), b)
	require.NoError(t, err)
	_, err = coordinator.AddConsumer(owner, subID, consumerAddress)
	require.NoError(t, err)
	backend.Commit()
	_, err = consumer.RequestRandomWords(random, subID, 1, 1, 1, [32]byte{}, false)
	require.Error(t, err)
	_, err = consumer.RequestRandomWords(owner, subID, 1, 1, 1, [32]byte{}, false)
	require.NoError(t, err)

	// Reassign ownership, check that only new owner can request
	_, err = consumer.TransferOwnership(owner, random.From)
	require.NoError(t, err)
	backend.Commit()
	_, err = consumer.AcceptOwnership(random)
	require.NoError(t, err)
	backend.Commit()
	_, err = consumer.RequestRandomWords(owner, subID, 1, 1, 1, [32]byte{}, false)
	require.Error(t, err)
	_, err = consumer.RequestRandomWords(random, subID, 1, 1, 1, [32]byte{}, false)
	require.NoError(t, err)
}

func TestVRFV2PlusIntegration_SimpleConsumerExample(t *testing.T) {
	owner := evmtestutils.MustNewSimTransactor(t)
	random := evmtestutils.MustNewSimTransactor(t)
	genesisData := gethtypes.GenesisAlloc{
		owner.From: {Balance: assets.Ether(10).ToInt()},
	}
	backend := cltest.NewSimulatedBackend(t, genesisData, ethconfig.Defaults.Miner.GasCeil)
	linkAddress, _, linkContract, err := link_token_interface.DeployLinkToken(
		owner, backend.Client())
	require.NoError(t, err)
	backend.Commit()
	// Deploy feed
	linkEthFeed, _, _, err :=
		mock_v3_aggregator_contract.DeployMockV3AggregatorContract(
			owner, backend.Client(), 18, vrftesthelpers.WeiPerUnitLink.BigInt()) // 0.01 eth per link
	require.NoError(t, err)
	backend.Commit()
	coordinatorAddress, _, coordinator, err :=
		vrf_coordinator_v2_5.DeployVRFCoordinatorV25(
			owner, backend.Client(), common.Address{}) // bhs not needed for this test
	require.NoError(t, err)
	backend.Commit()
	_, err = coordinator.SetLINKAndLINKNativeFeed(owner, linkAddress, linkEthFeed)
	require.NoError(t, err)
	backend.Commit()
	consumerAddress, _, consumer, err := vrf_v2plus_single_consumer.DeployVRFV2PlusSingleConsumerExample(owner, backend.Client(), coordinatorAddress, linkAddress, 1, 1, 1, [32]byte{}, false)
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
	uni := newVRFCoordinatorV2PlusUniverse(t, key, 1, false)
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
	ctx := testutils.Context(t)
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, key, 1, false)

	cfg := configtest.NewGeneralConfigSimulated(t, nil)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, cfg, uni.backend, key)
	require.NoError(t, app.Start(testutils.Context(t)))

	vrfkey, err := app.GetKeyStore().VRF().Create(ctx)
	require.NoError(t, err)
	registerProvingKeyHelper(t, uni.coordinatorV2UniverseCommon, uni.rootContract, vrfkey, &defaultMaxGasPrice)
	t.Run("non-proxied consumer", func(tt *testing.T) {
		carol := uni.vrfConsumers[0]
		carolContract := uni.consumerContracts[0]
		carolContractAddress := uni.consumerContractAddresses[0]

		_, err = carolContract.CreateSubscriptionAndFund(carol,
			big.NewInt(1000000000000000000)) // 0.1 LINK
		require.NoError(tt, err)
		uni.backend.Commit()
		_, err = carolContract.TopUpSubscriptionNative(carol,
			big.NewInt(2000000000000000000)) // 0.2 ETH
		uni.backend.Commit()
		// Ensure even with large number of consumers its still cheap
		var addrs []common.Address
		for i := 0; i < 99; i++ {
			addrs = append(addrs, testutils.NewAddress())
		}
		_, err = carolContract.UpdateSubscription(carol, addrs)
		require.NoError(tt, err)
		linkEstimate := estimateGas(tt, uni.backend, common.Address{},
			carolContractAddress, uni.consumerABI,
			"requestRandomWords", uint32(10000), uint16(2), uint32(1),
			vrfkey.PublicKey.MustHash(), false)
		tt.Log("gas estimate of non-proxied requestRandomWords with LINK payment:", linkEstimate)
		nativeEstimate := estimateGas(tt, uni.backend, common.Address{},
			carolContractAddress, uni.consumerABI,
			"requestRandomWords", uint32(10000), uint16(2), uint32(1),
			vrfkey.PublicKey.MustHash(), false)
		tt.Log("gas estimate of non-proxied requestRandomWords with Native payment:", nativeEstimate)
		assert.Less(tt, nativeEstimate, uint64(127_000),
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
		r, err := uni.backend.Client().TransactionReceipt(testutils.Context(t), tx.Hash())
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
		assert.Less(tt, estimate, uint64(106_000),
			"proxied testRequestRandomness tx gas cost more than expected")
	})
}

func TestVRFV2PlusIntegration_MaxConsumersCost(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, key, 1, false)
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
	assert.Less(t, estimate, uint64(540000))
	estimate = estimateGas(t, uni.backend, carolContractAddress,
		uni.rootContractAddress, uni.coordinatorABI,
		"addConsumer", subId, testutils.NewAddress())
	t.Log(estimate)
	assert.Less(t, estimate, uint64(100000))
}

func requestAndEstimateFulfillmentCost(
	t *testing.T,
	subID *big.Int,
	consumer *bind.TransactOpts,
	vrfkey vrfkey.KeyV2,
	minConfs uint16,
	gas uint32,
	numWords uint32,
	consumerContract vrftesthelpers.VRFConsumerContract,
	consumerContractAddress common.Address,
	uni coordinatorV2UniverseCommon,
	app *cltest.TestApplication,
	nativePayment bool,
	lowerBound, upperBound uint64,
) {
	_, err := consumerContract.RequestRandomness(consumer, vrfkey.PublicKey.MustHash(), subID, minConfs, gas, numWords, nativePayment)
	require.NoError(t, err)
	for i := 0; i < int(minConfs); i++ {
		uni.backend.Commit()
	}

	requestLog := FindLatestRandomnessRequestedLog(t, uni.rootContract, vrfkey.PublicKey.MustHash(), nil)
	s, err := proof.BigToSeed(requestLog.PreSeed())
	require.NoError(t, err)
	extraArgs, err := extraargs.EncodeV1(nativePayment)
	require.NoError(t, err)
	proof, rc, err := proof.GenerateProofResponseV2Plus(app.GetKeyStore().VRF(), vrfkey.ID(), proof.PreSeedDataV2Plus{
		PreSeed:          s,
		BlockHash:        requestLog.Raw().BlockHash,
		BlockNum:         requestLog.Raw().BlockNumber,
		SubId:            subID,
		CallbackGasLimit: gas,
		NumWords:         numWords,
		Sender:           consumerContractAddress,
		ExtraArgs:        extraArgs,
	})
	require.NoError(t, err)
	gasEstimate := estimateGas(t, uni.backend, common.Address{},
		uni.rootContractAddress, uni.coordinatorABI,
		"fulfillRandomWords", proof, rc, false)
	t.Log("consumer fulfillment gas estimate:", gasEstimate)
	assert.Greater(t, gasEstimate, lowerBound)
	assert.Less(t, gasEstimate, upperBound)
}

func TestVRFV2PlusIntegration_FulfillmentCost(t *testing.T) {
	ctx := testutils.Context(t)
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, key, 1, false)

	cfg := configtest.NewGeneralConfigSimulated(t, nil)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, cfg, uni.backend, key)
	require.NoError(t, app.Start(testutils.Context(t)))

	vrfkey, err := app.GetKeyStore().VRF().Create(ctx)
	require.NoError(t, err)
	registerProvingKeyHelper(t, uni.coordinatorV2UniverseCommon, uni.rootContract, vrfkey, &defaultMaxGasPrice)

	t.Run("non-proxied consumer", func(tt *testing.T) {
		carol := uni.vrfConsumers[0]
		carolContract := uni.consumerContracts[0]
		carolContractAddress := uni.consumerContractAddresses[0]

		_, err = carolContract.CreateSubscriptionAndFund(carol,
			big.NewInt(1000000000000000000)) // 0.1 LINK
		require.NoError(tt, err)
		uni.backend.Commit()
		subID, err2 := carolContract.SSubId(nil)
		require.NoError(tt, err2)
		_, err2 = carolContract.TopUpSubscriptionNative(carol,
			big.NewInt(2000000000000000000)) // 0.2 ETH
		require.NoError(tt, err2)
		gasRequested := 50_000
		nw := 1
		requestedIncomingConfs := 3
		t.Run("native payment", func(tt *testing.T) {
			requestAndEstimateFulfillmentCost(
				t,
				subID,
				carol,
				vrfkey,
				uint16(requestedIncomingConfs),
				uint32(gasRequested),
				uint32(nw),
				carolContract,
				carolContractAddress,
				uni.coordinatorV2UniverseCommon,
				app,
				true,
				120_000,
				500_000,
			)
		})

		t.Run("link payment", func(tt *testing.T) {
			requestAndEstimateFulfillmentCost(
				t,
				subID,
				carol,
				vrfkey,
				uint16(requestedIncomingConfs),
				uint32(gasRequested),
				uint32(nw),
				carolContract,
				carolContractAddress,
				uni.coordinatorV2UniverseCommon,
				app,
				false,
				120_000,
				500_000,
			)
		})
	})

	t.Run("proxied consumer", func(tt *testing.T) {
		consumerOwner := uni.neil
		consumerContract := uni.consumerProxyContract
		consumerContractAddress := uni.consumerProxyContractAddress

		_, err2 := consumerContract.CreateSubscriptionAndFund(consumerOwner, assets.Ether(5).ToInt())
		require.NoError(t, err2)
		uni.backend.Commit()
		subID, err2 := consumerContract.SSubId(nil)
		require.NoError(t, err2)
		gasRequested := 50_000
		nw := 1
		requestedIncomingConfs := 3
		requestAndEstimateFulfillmentCost(
			t,
			subID,
			consumerOwner,
			vrfkey,
			uint16(requestedIncomingConfs),
			uint32(gasRequested),
			uint32(nw),
			consumerContract,
			consumerContractAddress,
			uni.coordinatorV2UniverseCommon,
			app,
			false,
			120_000,
			500_000,
		)
	})
}

func setupSubscriptionAndFund(
	t *testing.T,
	uni coordinatorV2UniverseCommon,
	consumer *bind.TransactOpts,
	consumerContract vrftesthelpers.VRFConsumerContract,
	consumerAddress common.Address,
	linkAmount *big.Int,
	nativeAmount *big.Int) *big.Int {
	_, err := uni.rootContract.CreateSubscription(consumer)
	require.NoError(t, err)
	uni.backend.Commit()

	iter, err := uni.rootContract.FilterSubscriptionCreated(nil, nil)
	require.NoError(t, err)
	require.True(t, iter.Next(), "could not find SubscriptionCreated event for subID")
	subID := iter.Event().SubID()

	_, err = consumerContract.SetSubID(consumer, subID)
	require.NoError(t, err)

	_, err = uni.rootContract.AddConsumer(consumer, subID, consumerAddress)
	require.NoError(t, err, "failed to add consumer")
	uni.backend.Commit()

	b, err := utils.ABIEncode(`[{"type":"uint256"}]`, subID)
	require.NoError(t, err)
	_, err = uni.linkContract.TransferAndCall(
		uni.sergey, uni.rootContractAddress, linkAmount, b)
	require.NoError(t, err, "failed to fund sub")
	uni.backend.Commit()

	_, err = uni.rootContract.FundSubscriptionWithNative(consumer, subID, nativeAmount)
	require.NoError(t, err, "failed to fund sub with native")
	uni.backend.Commit()

	return subID
}

func TestVRFV2PlusIntegration_Migration(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1, false)
	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), toml.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: toml.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint64](5_000_000)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(1 * time.Second)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)

	// Create a subscription and fund with 5 LINK.
	consumerContract := uni.consumerContracts[0]
	consumer := uni.vrfConsumers[0]
	consumerAddress := uni.consumerContractAddresses[0]

	subID := setupSubscriptionAndFund(
		t,
		uni.coordinatorV2UniverseCommon,
		consumer,
		consumerContract,
		consumerAddress,
		new(big.Int).SetUint64(5e18),
		new(big.Int).SetUint64(3e18),
	)

	// Fund gas lane.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	require.NoError(t, app.Start(ctx))

	// Create VRF job using key1 and key2 on the same gas lane.
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

	// Make some randomness requests.
	numWords := uint32(2)

	requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, uni.rootContract, uni.backend, false)

	// Wait for fulfillment to be queued.
	require.Eventually(t, func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), time.Second)

	mine(t, requestID, subID, uni.backend, db, vrfcommon.V2Plus, testutils.SimulatedChainID)
	assertRandomWordsFulfilled(t, requestID, true, uni.rootContract, false)

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)

	rw, err := consumerContract.SRandomWords(nil, big.NewInt(0))
	require.NoError(t, err)

	subV1, err := uni.rootContract.GetSubscription(nil, subID)
	require.NoError(t, err)

	_, err = uni.rootContract.Migrate(consumer, subID, uni.migrationTestCoordinatorAddress)
	require.NoError(t, err)
	uni.backend.Commit()

	subV2, err := uni.migrationTestCoordinator.GetSubscription(nil, subID)
	require.NoError(t, err)

	totalLinkBalance, err := uni.migrationTestCoordinator.STotalLinkBalance(nil)
	require.NoError(t, err)
	totalNativeBalance, err := uni.migrationTestCoordinator.STotalNativeBalance(nil)
	require.NoError(t, err)
	linkContractBalance, err := uni.linkContract.BalanceOf(nil, uni.migrationTestCoordinatorAddress)
	require.NoError(t, err)
	balance, err := uni.backend.Client().BalanceAt(ctx, uni.migrationTestCoordinatorAddress, nil)
	require.NoError(t, err)

	require.Equal(t, subV1.Balance(), totalLinkBalance)
	require.Equal(t, subV1.NativeBalance(), totalNativeBalance)
	require.Equal(t, subV1.Balance(), linkContractBalance)
	require.Equal(t, subV1.NativeBalance(), balance)

	require.Equal(t, subV1.Balance(), subV2.LinkBalance)
	require.Equal(t, subV1.NativeBalance(), subV2.NativeBalance)
	require.Equal(t, subV1.Owner(), subV2.Owner)
	require.Equal(t, len(subV1.Consumers()), len(subV2.Consumers))
	for i, c := range subV1.Consumers() {
		require.Equal(t, c, subV2.Consumers[i])
	}

	minRequestConfirmations := uint16(2)
	requestID2, rw2 := requestRandomnessAndValidate(
		t,
		consumer,
		consumerContract,
		keyHash,
		subID,
		minRequestConfirmations,
		50_000,
		numWords,
		uni,
		true,
	)
	require.NotEqual(t, requestID, requestID2)
	require.NotEqual(t, rw, rw2)
	requestID3, rw3 := requestRandomnessAndValidate(
		t,
		consumer,
		consumerContract,
		keyHash,
		subID,
		minRequestConfirmations,
		50_000,
		numWords,
		uni,
		false,
	)
	require.NotEqual(t, requestID2, requestID3)
	require.NotEqual(t, rw2, rw3)
}

func requestRandomnessAndValidate(t *testing.T,
	consumer *bind.TransactOpts,
	consumerContract vrftesthelpers.VRFConsumerContract,
	keyHash common.Hash,
	subID *big.Int,
	minConfs uint16,
	gas, numWords uint32,
	uni coordinatorV2PlusUniverse,
	nativePayment bool) (*big.Int, *big.Int) {
	_, err := consumerContract.RequestRandomness(
		consumer,
		keyHash,
		subID,
		minConfs,
		gas,
		numWords,
		nativePayment, // test link payment works after migration
	)
	require.NoError(t, err)
	uni.backend.Commit()

	requestID, err := consumerContract.SRequestId(nil)
	require.NoError(t, err)

	_, err = uni.migrationTestCoordinator.FulfillRandomWords(uni.neil, requestID)
	require.NoError(t, err)
	uni.backend.Commit()

	rw, err := consumerContract.SRandomWords(nil, big.NewInt(0))
	require.NoError(t, err)

	return requestID, rw
}

func TestVRFV2PlusIntegration_CancelSubscription(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, key, 1, false)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]
	linkAmount := new(big.Int).SetUint64(5e18)
	nativeAmount := new(big.Int).SetUint64(3e18)
	subID := setupSubscriptionAndFund(
		t,
		uni.coordinatorV2UniverseCommon,
		consumer,
		consumerContract,
		consumerContractAddress,
		linkAmount,
		nativeAmount,
	)

	linkBalanceBeforeCancel, err := uni.linkContract.BalanceOf(nil, uni.neil.From)
	require.NoError(t, err)
	nativeBalanceBeforeCancel, err := uni.backend.Client().BalanceAt(testutils.Context(t), uni.neil.From, nil)
	require.NoError(t, err)

	// non-owner cannot cancel subscription
	_, err = uni.rootContract.CancelSubscription(uni.neil, subID, consumer.From)
	require.Error(t, err)

	_, err = uni.rootContract.CancelSubscription(consumer, subID, uni.neil.From)
	require.NoError(t, err)
	uni.backend.Commit()

	AssertLinkBalance(t, uni.linkContract, uni.neil.From, linkBalanceBeforeCancel.Add(linkBalanceBeforeCancel, linkAmount))
	AssertNativeBalance(t, uni.backend, uni.neil.From, nativeBalanceBeforeCancel.Add(nativeBalanceBeforeCancel, nativeAmount))
}

func TestVRFV2PlusIntegration_ReplayOldRequestsOnStartUp(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 1, false)

	testReplayOldRequestsOnStartUp(
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
		vrfcommon.V2Plus,
		false,
	)
}
