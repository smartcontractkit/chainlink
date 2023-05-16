package vrf_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/sqlx"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	v2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
	evmlogger "github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_consumer_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_consumer_v2_upgradeable_example"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_external_sub_owner_example"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_malicious_consumer_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_single_consumer_example"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2_proxy_admin"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2_reverting_example"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2_transparent_upgradeable_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2_wrapper_consumer_example"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg/datatypes"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/proof"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// vrfConsumerContract is the common interface implemented by
// the example contracts used for the integration tests.
type vrfConsumerContract interface {
	CreateSubscriptionAndFund(opts *bind.TransactOpts, fundingJuels *big.Int) (*gethtypes.Transaction, error)
	SSubId(opts *bind.CallOpts) (uint64, error)
	SRequestId(opts *bind.CallOpts) (*big.Int, error)
	RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*gethtypes.Transaction, error)
	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)
}

type coordinatorV2Universe struct {
	// Golang wrappers of solidity contracts
	consumerContracts         []*vrf_consumer_v2.VRFConsumerV2
	consumerContractAddresses []common.Address

	rootContract                     *vrf_coordinator_v2.VRFCoordinatorV2
	rootContractAddress              common.Address
	batchCoordinatorContract         *batch_vrf_coordinator_v2.BatchVRFCoordinatorV2
	batchCoordinatorContractAddress  common.Address
	linkContract                     *link_token_interface.LinkToken
	linkContractAddress              common.Address
	linkEthFeedAddress               common.Address
	bhsContract                      *blockhash_store.BlockhashStore
	bhsContractAddress               common.Address
	batchBHSContract                 *batch_blockhash_store.BatchBlockhashStore
	batchBHSContractAddress          common.Address
	maliciousConsumerContract        *vrf_malicious_consumer_v2.VRFMaliciousConsumerV2
	maliciousConsumerContractAddress common.Address
	revertingConsumerContract        *vrfv2_reverting_example.VRFV2RevertingExample
	revertingConsumerContractAddress common.Address
	// This is a VRFConsumerV2Upgradeable wrapper that points to the proxy address.
	consumerProxyContract        *vrf_consumer_v2_upgradeable_example.VRFConsumerV2UpgradeableExample
	consumerProxyContractAddress common.Address
	proxyAdminAddress            common.Address

	// Abstract representation of the ethereum blockchain
	backend        *backends.SimulatedBackend
	coordinatorABI *abi.ABI
	consumerABI    *abi.ABI

	// Cast of participants
	vrfConsumers []*bind.TransactOpts // Authors of consuming contracts that request randomness
	sergey       *bind.TransactOpts   // Owns all the LINK initially
	neil         *bind.TransactOpts   // Node operator running VRF service
	ned          *bind.TransactOpts   // Secondary node operator
	nallory      *bind.TransactOpts   // Oracle transactor
	evil         *bind.TransactOpts   // Author of a malicious consumer contract
	reverter     *bind.TransactOpts   // Author of always reverting contract
}

var (
	weiPerUnitLink = decimal.RequireFromString("10000000000000000")
)

func newVRFCoordinatorV2Universe(t *testing.T, key ethkey.KeyV2, numConsumers int) coordinatorV2Universe {
	testutils.SkipShort(t, "VRFCoordinatorV2Universe")
	oracleTransactor, err := bind.NewKeyedTransactorWithChainID(key.ToEcdsaPrivKey(), testutils.SimulatedChainID)
	require.NoError(t, err)
	var (
		sergey       = testutils.MustNewSimTransactor(t)
		neil         = testutils.MustNewSimTransactor(t)
		ned          = testutils.MustNewSimTransactor(t)
		evil         = testutils.MustNewSimTransactor(t)
		reverter     = testutils.MustNewSimTransactor(t)
		nallory      = oracleTransactor
		vrfConsumers []*bind.TransactOpts
	)

	// Create consumer contract deployer identities
	for i := 0; i < numConsumers; i++ {
		vrfConsumers = append(vrfConsumers, testutils.MustNewSimTransactor(t))
	}

	genesisData := core.GenesisAlloc{
		sergey.From:   {Balance: assets.Ether(1000).ToInt()},
		neil.From:     {Balance: assets.Ether(1000).ToInt()},
		ned.From:      {Balance: assets.Ether(1000).ToInt()},
		nallory.From:  {Balance: assets.Ether(1000).ToInt()},
		evil.From:     {Balance: assets.Ether(1000).ToInt()},
		reverter.From: {Balance: assets.Ether(1000).ToInt()},
	}
	for _, consumer := range vrfConsumers {
		genesisData[consumer.From] = core.GenesisAccount{
			Balance: assets.Ether(1000).ToInt(),
		}
	}

	gasLimit := uint32(ethconfig.Defaults.Miner.GasCeil)
	consumerABI, err := abi.JSON(strings.NewReader(
		vrf_consumer_v2.VRFConsumerV2ABI))
	require.NoError(t, err)
	coordinatorABI, err := abi.JSON(strings.NewReader(
		vrf_coordinator_v2.VRFCoordinatorV2ABI))
	require.NoError(t, err)
	backend := cltest.NewSimulatedBackend(t, genesisData, gasLimit)
	// Deploy link
	linkAddress, _, linkContract, err := link_token_interface.DeployLinkToken(
		sergey, backend)
	require.NoError(t, err, "failed to deploy link contract to simulated ethereum blockchain")

	// Deploy feed
	linkEthFeed, _, _, err :=
		mock_v3_aggregator_contract.DeployMockV3AggregatorContract(
			evil, backend, 18, weiPerUnitLink.BigInt()) // 0.01 eth per link
	require.NoError(t, err)

	// Deploy blockhash store
	bhsAddress, _, bhsContract, err := blockhash_store.DeployBlockhashStore(neil, backend)
	require.NoError(t, err, "failed to deploy BlockhashStore contract to simulated ethereum blockchain")

	// Deploy batch blockhash store
	batchBHSAddress, _, batchBHSContract, err := batch_blockhash_store.DeployBatchBlockhashStore(neil, backend, bhsAddress)
	require.NoError(t, err, "failed to deploy BatchBlockhashStore contract to simulated ethereum blockchain")

	// Deploy VRF V2 coordinator
	coordinatorAddress, _, coordinatorContract, err :=
		vrf_coordinator_v2.DeployVRFCoordinatorV2(
			neil, backend, linkAddress, bhsAddress, linkEthFeed /* linkEth*/)
	require.NoError(t, err, "failed to deploy VRFCoordinatorV2 contract to simulated ethereum blockchain")
	backend.Commit()

	// Deploy batch VRF V2 coordinator
	batchCoordinatorAddress, _, batchCoordinatorContract, err :=
		batch_vrf_coordinator_v2.DeployBatchVRFCoordinatorV2(
			neil, backend, coordinatorAddress,
		)
	require.NoError(t, err, "failed to deploy BatchVRFCoordinatorV2 contract to simulated ethereum blockchain")
	backend.Commit()

	// Create the VRF consumers.
	var (
		consumerContracts         []*vrf_consumer_v2.VRFConsumerV2
		consumerContractAddresses []common.Address
	)
	for _, author := range vrfConsumers {
		// Deploy a VRF consumer. It has a starting balance of 500 LINK.
		consumerContractAddress, _, consumerContract, err :=
			vrf_consumer_v2.DeployVRFConsumerV2(
				author, backend, coordinatorAddress, linkAddress)
		require.NoError(t, err, "failed to deploy VRFConsumer contract to simulated ethereum blockchain")
		_, err = linkContract.Transfer(sergey, consumerContractAddress, assets.Ether(500).ToInt()) // Actually, LINK
		require.NoError(t, err, "failed to send LINK to VRFConsumer contract on simulated ethereum blockchain")

		consumerContracts = append(consumerContracts, consumerContract)
		consumerContractAddresses = append(consumerContractAddresses, consumerContractAddress)

		backend.Commit()
	}

	// Deploy malicious consumer with 1 link
	maliciousConsumerContractAddress, _, maliciousConsumerContract, err :=
		vrf_malicious_consumer_v2.DeployVRFMaliciousConsumerV2(
			evil, backend, coordinatorAddress, linkAddress)
	require.NoError(t, err, "failed to deploy VRFMaliciousConsumer contract to simulated ethereum blockchain")
	_, err = linkContract.Transfer(sergey, maliciousConsumerContractAddress, assets.Ether(1).ToInt()) // Actually, LINK
	require.NoError(t, err, "failed to send LINK to VRFMaliciousConsumer contract on simulated ethereum blockchain")
	backend.Commit()

	// Deploy upgradeable consumer, proxy, and proxy admin
	upgradeableConsumerAddress, _, _, err := vrf_consumer_v2_upgradeable_example.DeployVRFConsumerV2UpgradeableExample(neil, backend)
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

	proxiedConsumer, err := vrf_consumer_v2_upgradeable_example.NewVRFConsumerV2UpgradeableExample(
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
	revertingConsumerContractAddress, _, revertingConsumerContract, err := vrfv2_reverting_example.DeployVRFV2RevertingExample(
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
		uint32(vrf.GasAfterPaymentCalculation), // gasAfterPaymentCalculation
		big.NewInt(1e16),                       // 0.01 eth per link fallbackLinkPrice
		vrf_coordinator_v2.VRFCoordinatorV2FeeConfig{
			FulfillmentFlatFeeLinkPPMTier1: uint32(1000),
			FulfillmentFlatFeeLinkPPMTier2: uint32(1000),
			FulfillmentFlatFeeLinkPPMTier3: uint32(100),
			FulfillmentFlatFeeLinkPPMTier4: uint32(10),
			FulfillmentFlatFeeLinkPPMTier5: uint32(1),
			ReqsForTier2:                   big.NewInt(10),
			ReqsForTier3:                   big.NewInt(20),
			ReqsForTier4:                   big.NewInt(30),
			ReqsForTier5:                   big.NewInt(40),
		},
	)
	require.NoError(t, err, "failed to set coordinator configuration")
	backend.Commit()

	return coordinatorV2Universe{
		vrfConsumers:              vrfConsumers,
		consumerContracts:         consumerContracts,
		consumerContractAddresses: consumerContractAddresses,

		batchCoordinatorContract:        batchCoordinatorContract,
		batchCoordinatorContractAddress: batchCoordinatorAddress,

		revertingConsumerContract:        revertingConsumerContract,
		revertingConsumerContractAddress: revertingConsumerContractAddress,

		consumerProxyContract:        proxiedConsumer,
		consumerProxyContractAddress: proxiedConsumer.Address(),
		proxyAdminAddress:            proxyAdminAddress,

		rootContract:                     coordinatorContract,
		rootContractAddress:              coordinatorAddress,
		linkContract:                     linkContract,
		linkContractAddress:              linkAddress,
		linkEthFeedAddress:               linkEthFeed,
		bhsContract:                      bhsContract,
		bhsContractAddress:               bhsAddress,
		batchBHSContract:                 batchBHSContract,
		batchBHSContractAddress:          batchBHSAddress,
		maliciousConsumerContract:        maliciousConsumerContract,
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
	}
}

// Send eth from prefunded account.
// Amount is number of ETH not wei.
func sendEth(t *testing.T, key ethkey.KeyV2, ec *backends.SimulatedBackend, to common.Address, eth int) {
	nonce, err := ec.PendingNonceAt(testutils.Context(t), key.Address)
	require.NoError(t, err)
	tx := gethtypes.NewTx(&gethtypes.DynamicFeeTx{
		ChainID:   big.NewInt(1337),
		Nonce:     nonce,
		GasTipCap: big.NewInt(1),
		GasFeeCap: assets.GWei(10).ToInt(), // block base fee in sim
		Gas:       uint64(21_000),
		To:        &to,
		Value:     big.NewInt(0).Mul(big.NewInt(int64(eth)), big.NewInt(1e18)),
		Data:      nil,
	})
	signedTx, err := gethtypes.SignTx(tx, gethtypes.NewLondonSigner(big.NewInt(1337)), key.ToEcdsaPrivKey())
	require.NoError(t, err)
	err = ec.SendTransaction(testutils.Context(t), signedTx)
	require.NoError(t, err)
	ec.Commit()
}

func subscribeVRF(
	t *testing.T,
	author *bind.TransactOpts,
	consumerContract vrfConsumerContract,
	coordinatorContract vrf_coordinator_v2.VRFCoordinatorV2Interface,
	backend *backends.SimulatedBackend,
	fundingJuels *big.Int,
) (vrf_coordinator_v2.GetSubscription, uint64) {
	_, err := consumerContract.CreateSubscriptionAndFund(author, fundingJuels)
	require.NoError(t, err)
	backend.Commit()

	subID, err := consumerContract.SSubId(nil)
	require.NoError(t, err)

	sub, err := coordinatorContract.GetSubscription(nil, subID)
	require.NoError(t, err)
	return sub, subID
}

func createVRFJobs(
	t *testing.T,
	fromKeys [][]ethkey.KeyV2,
	app *cltest.TestApplication,
	coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
	coordinatorAddress common.Address,
	batchCoordinatorAddress common.Address,
	uni coordinatorV2Universe,
	batchEnabled bool,
	gasLanePrices ...*assets.Wei,
) (jobs []job.Job) {
	if len(gasLanePrices) != len(fromKeys) {
		t.Fatalf("must provide one gas lane price for each set of from addresses. len(gasLanePrices) != len(fromKeys) [%d != %d]",
			len(gasLanePrices), len(fromKeys))
	}
	// Create separate jobs for each gas lane and register their keys
	for i, keys := range fromKeys {
		var keyStrs []string
		for _, k := range keys {
			keyStrs = append(keyStrs, k.Address.String())
		}

		vrfkey, err := app.GetKeyStore().VRF().Create()
		require.NoError(t, err)

		jid := uuid.NewV4()
		incomingConfs := 2
		s := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
			JobID:                    jid.String(),
			Name:                     fmt.Sprintf("vrf-primary-%d", i),
			CoordinatorAddress:       coordinatorAddress.Hex(),
			BatchCoordinatorAddress:  batchCoordinatorAddress.Hex(),
			BatchFulfillmentEnabled:  batchEnabled,
			MinIncomingConfirmations: incomingConfs,
			PublicKey:                vrfkey.PublicKey.String(),
			FromAddresses:            keyStrs,
			BackoffInitialDelay:      10 * time.Millisecond,
			BackoffMaxDelay:          time.Second,
			V2:                       true,
			GasLanePrice:             gasLanePrices[i],
		}).Toml()
		jb, err := vrf.ValidatedVRFSpec(s)
		t.Log(jb.VRFSpec.PublicKey.MustHash(), vrfkey.PublicKey.MustHash())
		require.NoError(t, err)
		err = app.JobSpawner().CreateJob(&jb)
		require.NoError(t, err)
		registerProvingKeyHelper(t, uni, coordinator, vrfkey)
		jobs = append(jobs, jb)
	}
	// Wait until all jobs are active and listening for logs
	gomega.NewWithT(t).Eventually(func() bool {
		jbs := app.JobSpawner().ActiveJobs()
		var count int
		for _, jb := range jbs {
			if jb.Type == job.VRF {
				count++
			}
		}
		return count == len(fromKeys)
	}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())
	// Unfortunately the lb needs heads to be able to backfill logs to new subscribers.
	// To avoid confirming
	// TODO: it could just backfill immediately upon receiving a new subscriber? (though would
	// only be useful for tests, probably a more robust way is to have the job spawner accept a signal that a
	// job is fully up and running and not add it to the active jobs list before then)
	time.Sleep(2 * time.Second)

	return
}

func requestRandomnessForWrapper(
	t *testing.T,
	vrfWrapperConsumer vrfv2_wrapper_consumer_example.VRFV2WrapperConsumerExample,
	consumerOwner *bind.TransactOpts,
	keyHash common.Hash,
	subID uint64,
	numWords uint32,
	cbGasLimit uint32,
	coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
	uni coordinatorV2Universe,
	wrapperOverhead uint32,
) (*big.Int, uint64) {
	minRequestConfirmations := uint16(3)
	_, err := vrfWrapperConsumer.MakeRequest(
		consumerOwner,
		cbGasLimit,
		minRequestConfirmations,
		numWords,
	)
	require.NoError(t, err)
	uni.backend.Commit()

	iter, err := coordinator.FilterRandomWordsRequested(nil, nil, []uint64{subID}, nil)
	require.NoError(t, err, "could not filter RandomWordsRequested events")

	var events []*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested
	for iter.Next() {
		events = append(events, iter.Event)
	}

	wrapperIter, err := vrfWrapperConsumer.FilterWrapperRequestMade(nil, nil)
	require.NoError(t, err, "could not filter WrapperRequestMade events")

	wrapperConsumerEvents := []*vrfv2_wrapper_consumer_example.VRFV2WrapperConsumerExampleWrapperRequestMade{}
	for wrapperIter.Next() {
		wrapperConsumerEvents = append(wrapperConsumerEvents, wrapperIter.Event)
	}

	event := events[len(events)-1]
	wrapperConsumerEvent := wrapperConsumerEvents[len(wrapperConsumerEvents)-1]
	require.Equal(t, event.RequestId, wrapperConsumerEvent.RequestId, "request ID in consumer log does not match request ID in coordinator log")
	require.Equal(t, keyHash.Bytes(), event.KeyHash[:], "key hash of event (%s) and of request not equal (%s)", hex.EncodeToString(event.KeyHash[:]), keyHash.String())
	require.Equal(t, cbGasLimit+(cbGasLimit/63+1)+wrapperOverhead, event.CallbackGasLimit, "callback gas limit of event and of request not equal")
	require.Equal(t, minRequestConfirmations, event.MinimumRequestConfirmations, "min request confirmations of event and of request not equal")
	require.Equal(t, numWords, event.NumWords, "num words of event and of request not equal")

	return event.RequestId, event.Raw.BlockNumber
}

// requestRandomness requests randomness from the given vrf consumer contract
// and asserts that the request ID logged by the RandomWordsRequested event
// matches the request ID that is returned and set by the consumer contract.
// The request ID and request block number are then returned to the caller.
func requestRandomnessAndAssertRandomWordsRequestedEvent(
	t *testing.T,
	vrfConsumerHandle vrfConsumerContract,
	consumerOwner *bind.TransactOpts,
	keyHash common.Hash,
	subID uint64,
	numWords uint32,
	cbGasLimit uint32,
	coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
	uni coordinatorV2Universe,
) (requestID *big.Int, requestBlockNumber uint64) {
	minRequestConfirmations := uint16(2)
	_, err := vrfConsumerHandle.RequestRandomness(
		consumerOwner,
		keyHash,
		subID,
		minRequestConfirmations,
		cbGasLimit,
		numWords,
	)
	require.NoError(t, err)

	uni.backend.Commit()

	iter, err := coordinator.FilterRandomWordsRequested(nil, nil, []uint64{subID}, nil)
	require.NoError(t, err, "could not filter RandomWordsRequested events")

	var events []*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested
	for iter.Next() {
		events = append(events, iter.Event)
	}

	requestID, err = vrfConsumerHandle.SRequestId(nil)
	require.NoError(t, err)

	event := events[len(events)-1]
	require.Equal(t, event.RequestId, requestID, "request ID in contract does not match request ID in log")
	require.Equal(t, keyHash.Bytes(), event.KeyHash[:], "key hash of event (%s) and of request not equal (%s)", hex.EncodeToString(event.KeyHash[:]), keyHash.String())
	require.Equal(t, cbGasLimit, event.CallbackGasLimit, "callback gas limit of event and of request not equal")
	require.Equal(t, minRequestConfirmations, event.MinimumRequestConfirmations, "min request confirmations of event and of request not equal")
	require.Equal(t, numWords, event.NumWords, "num words of event and of request not equal")

	return requestID, event.Raw.BlockNumber
}

// subscribeAndAssertSubscriptionCreatedEvent subscribes the given consumer contract
// to VRF and funds the subscription with the given fundingJuels amount. It returns the
// subscription ID of the resulting subscription.
func subscribeAndAssertSubscriptionCreatedEvent(
	t *testing.T,
	vrfConsumerHandle vrfConsumerContract,
	consumerOwner *bind.TransactOpts,
	consumerContractAddress common.Address,
	fundingJuels *big.Int,
	coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
	uni coordinatorV2Universe,
) uint64 {
	// Create a subscription and fund with LINK.
	sub, subID := subscribeVRF(t, consumerOwner, vrfConsumerHandle, coordinator, uni.backend, fundingJuels)
	require.Equal(t, uint64(1), subID)
	require.Equal(t, fundingJuels.String(), sub.Balance.String())

	// Assert the subscription event in the coordinator contract.
	iter, err := coordinator.FilterSubscriptionCreated(nil, []uint64{subID})
	require.NoError(t, err)
	found := false
	for iter.Next() {
		if iter.Event.Owner != consumerContractAddress {
			require.FailNowf(t, "SubscriptionCreated event contains wrong owner address", "expected: %+v, actual: %+v", consumerContractAddress, iter.Event.Owner)
		} else {
			found = true
		}
	}
	require.True(t, found, "could not find SubscriptionCreated event for subID %d", subID)

	return subID
}

func assertRandomWordsFulfilled(
	t *testing.T,
	requestID *big.Int,
	expectedSuccess bool,
	coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
) (rwfe *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled) {
	// Check many times in case there are delays processing the event
	// this could happen occasionally and cause flaky tests.
	numChecks := 3
	found := false
	for i := 0; i < numChecks; i++ {
		filter, err := coordinator.FilterRandomWordsFulfilled(nil, []*big.Int{requestID})
		require.NoError(t, err)

		for filter.Next() {
			require.Equal(t, expectedSuccess, filter.Event.Success, "fulfillment event success not correct, expected: %+v, actual: %+v", expectedSuccess, filter.Event.Success)
			require.Equal(t, requestID, filter.Event.RequestId)
			found = true
			rwfe = filter.Event
		}

		if found {
			break
		}

		// Wait a bit and try again.
		time.Sleep(time.Second)
	}
	require.True(t, found, "RandomWordsFulfilled event not found")
	return
}

func assertNumRandomWords(
	t *testing.T,
	contract vrfConsumerContract,
	numWords uint32,
) {
	var err error
	for i := uint32(0); i < numWords; i++ {
		_, err = contract.SRandomWords(nil, big.NewInt(int64(i)))
		require.NoError(t, err)
	}
}

func mine(t *testing.T, requestID *big.Int, subID uint64, uni coordinatorV2Universe, db *sqlx.DB) bool {
	return gomega.NewWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		var txs []txmgr.DbEthTx
		err := db.Select(&txs, `
		SELECT * FROM eth_txes
		WHERE eth_txes.state = 'confirmed'
			AND eth_txes.meta->>'RequestID' = $1
			AND CAST(eth_txes.meta->>'SubId' AS NUMERIC) = $2 LIMIT 1
		`, common.BytesToHash(requestID.Bytes()).String(), subID)
		require.NoError(t, err)
		t.Log("num txs", len(txs))
		return len(txs) == 1
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
}

func mineBatch(t *testing.T, requestIDs []*big.Int, subID uint64, uni coordinatorV2Universe, db *sqlx.DB) bool {
	requestIDMap := map[string]bool{}
	for _, requestID := range requestIDs {
		requestIDMap[common.BytesToHash(requestID.Bytes()).String()] = false
	}
	return gomega.NewWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		var txs []txmgr.DbEthTx
		err := db.Select(&txs, `
		SELECT * FROM eth_txes
		WHERE eth_txes.state = 'confirmed'
			AND CAST(eth_txes.meta->>'SubId' AS NUMERIC) = $1
		`, subID)
		require.NoError(t, err)
		for _, tx := range txs {
			var evmTx txmgr.EvmTx
			txmgr.DbEthTxToEthTx(tx, &evmTx)
			meta, err := evmTx.GetMeta()
			require.NoError(t, err)
			t.Log("meta:", meta)
			for _, requestID := range meta.RequestIDs {
				if _, ok := requestIDMap[requestID.String()]; ok {
					requestIDMap[requestID.String()] = true
				}
			}
		}
		foundAll := true
		for _, found := range requestIDMap {
			foundAll = foundAll && found
		}
		t.Log("requestIDMap:", requestIDMap)
		return foundAll
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
}

func TestVRFV2Integration_SingleConsumer_HappyPath_BatchFulfillment(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	testSingleConsumerHappyPathBatchFulfillment(
		t,
		ownerKey,
		uni,
		uni.vrfConsumers[0],
		uni.consumerContracts[0],
		uni.consumerContractAddresses[0],
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		5,     // number of requests to send
		false, // don't send big callback
	)
}

func TestVRFV2Integration_SingleConsumer_HappyPath_BatchFulfillment_BigGasCallback(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	testSingleConsumerHappyPathBatchFulfillment(
		t,
		ownerKey,
		uni,
		uni.vrfConsumers[0],
		uni.consumerContracts[0],
		uni.consumerContractAddresses[0],
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		5,    // number of requests to send
		true, // send big callback
	)
}

func TestVRFV2Integration_SingleConsumer_HappyPath(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	testSingleConsumerHappyPath(
		t,
		ownerKey,
		uni,
		uni.vrfConsumers[0],
		uni.consumerContracts[0],
		uni.consumerContractAddresses[0],
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress)
}

func TestVRFV2Integration_SingleConsumer_EOA_Request(t *testing.T) {
	t.Parallel()
	testEoa(t, false)
}

func TestVRFV2Integration_SingleConsumer_EOA_Request_Batching_Enabled(t *testing.T) {
	t.Parallel()
	testEoa(t, true)
}

func testEoa(t *testing.T, batchingEnabled bool) {
	gasLimit := int64(2_500_000)

	finalityDepth := uint32(50)

	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, _ := heavyweight.FullTestDBV2(t, "vrfv2_singleconsumer_eoa_request", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = ptr(uint32(gasLimit))
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.EVM[0].FinalityDepth = ptr(finalityDepth)
	})
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]
	// Create a subscription and fund with 500 LINK.
	subAmount := big.NewInt(1).Mul(big.NewInt(5e18), big.NewInt(100))
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, subAmount, uni.rootContract, uni)

	// Createa a new subscription.
	_, err := uni.rootContract.CreateSubscription(consumer)
	require.NoError(t, err)
	uni.backend.Commit()

	// Add the EOA as a consumer.
	_, err = uni.rootContract.AddConsumer(consumer, subID+1, consumer.From)
	require.NoError(t, err)
	uni.backend.Commit()

	// Fund the subscription with 1 LINK.
	b, err := utils.ABIEncode(`[{"type":"uint64"}]`, subID+1)
	require.NoError(t, err)
	_, err = uni.linkContract.TransferAndCall(uni.sergey, uni.rootContractAddress, big.NewInt(1e18), b)
	require.NoError(t, err)
	uni.backend.Commit()

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
		uni,
		batchingEnabled,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make a randomness request with the EOA. This request is impossible to fulfill.
	numWords := uint32(1)
	minRequestConfirmations := uint16(2)
	_, err = uni.rootContract.RequestRandomWords(consumer, keyHash, subID+1, minRequestConfirmations, uint32(200_000), numWords)
	require.NoError(t, err)
	uni.backend.Commit()

	// Ensure request is not fulfilled.
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 0
	}, 5*time.Second, time.Second).Should(gomega.BeTrue())

	// Create query to fetch the application's log broadcasts.
	var broadcastsBeforeFinality []evmlogger.LogBroadcast
	var broadcastsAfterFinality []evmlogger.LogBroadcast
	query := `SELECT block_hash, consumed, log_index, job_id FROM log_broadcasts`
	q := pg.NewQ(app.GetSqlxDB(), app.Logger, app.Config)

	// Execute the query.
	err = q.Select(&broadcastsBeforeFinality, query)
	require.NoError(t, err)

	// Ensure there is only one log broadcast (our EOA request), and that
	// it hasn't been marked as consumed yet.
	require.Equal(t, 1, len(broadcastsBeforeFinality))
	require.Equal(t, false, broadcastsBeforeFinality[0].Consumed)

	// Create new blocks until the finality depth has elapsed.
	for i := 0; i < int(finalityDepth); i++ {
		uni.backend.Commit()
	}

	// Ensure the request is still not fulfilled.
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 0
	}, 5*time.Second, time.Second).Should(gomega.BeTrue())

	// Execute the query for log broadcasts again after finality depth has elapsed.
	err = q.Select(&broadcastsAfterFinality, query)
	require.NoError(t, err)

	// Ensure that there is still only one log broadcast (our EOA request), but that
	// it has been marked as "consumed," such that it won't be retried.
	require.Equal(t, 1, len(broadcastsAfterFinality))
	require.Equal(t, true, broadcastsAfterFinality[0].Consumed)

	t.Log("Done!")
}

func TestVRFV2Integration_SingleConsumer_EIP150_HappyPath(t *testing.T) {
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
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]
	// Create a subscription and fund with 500 LINK.
	subAmount := big.NewInt(1).Mul(big.NewInt(5e18), big.NewInt(100))
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, subAmount, uni.rootContract, uni)

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
		uni,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the first randomness request.
	numWords := uint32(1)
	requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, uint32(callBackGasLimit), uni.rootContract, uni)

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

func TestVRFV2Integration_SingleConsumer_EIP150_Revert(t *testing.T) {
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
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]
	// Create a subscription and fund with 500 LINK.
	subAmount := big.NewInt(1).Mul(big.NewInt(5e18), big.NewInt(100))
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, subAmount, uni.rootContract, uni)

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
		uni,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the first randomness request.
	numWords := uint32(1)
	requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, uint32(callBackGasLimit), uni.rootContract, uni)

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

func deployWrapper(t *testing.T, uni coordinatorV2Universe, wrapperOverhead uint32, coordinatorOverhead uint32, keyHash common.Hash) (
	wrapper *vrfv2_wrapper.VRFV2Wrapper,
	wrapperAddress common.Address,
	wrapperConsumer *vrfv2_wrapper_consumer_example.VRFV2WrapperConsumerExample,
	wrapperConsumerAddress common.Address,
) {
	wrapperAddress, _, wrapper, err := vrfv2_wrapper.DeployVRFV2Wrapper(uni.neil, uni.backend, uni.linkContractAddress, uni.linkEthFeedAddress, uni.rootContractAddress)
	require.NoError(t, err)
	uni.backend.Commit()

	_, err = wrapper.SetConfig(uni.neil, wrapperOverhead, coordinatorOverhead, 0, keyHash, 10)
	require.NoError(t, err)
	uni.backend.Commit()

	wrapperConsumerAddress, _, wrapperConsumer, err = vrfv2_wrapper_consumer_example.DeployVRFV2WrapperConsumerExample(uni.neil, uni.backend, uni.linkContractAddress, wrapperAddress)
	require.NoError(t, err)
	uni.backend.Commit()

	return
}

func TestVRFV2Integration_SingleConsumer_Wrapper(t *testing.T) {
	t.Parallel()
	wrapperOverhead := uint32(30_000)
	coordinatorOverhead := uint32(90_000)

	callBackGasLimit := int64(100_000) // base callback gas.
	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, "vrfv2_singleconsumer_wrapper", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint32](3_500_000)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)

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
		uni,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	wrapper, _, consumer, consumerAddress := deployWrapper(t, uni, wrapperOverhead, coordinatorOverhead, keyHash)

	// Fetch Subscription ID for Wrapper.
	wrapperSubID, err := wrapper.SUBSCRIPTIONID(nil)
	require.NoError(t, err)

	// Fund Subscription.
	b, err := utils.ABIEncode(`[{"type":"uint64"}]`, wrapperSubID)
	require.NoError(t, err)
	_, err = uni.linkContract.TransferAndCall(uni.sergey, uni.rootContractAddress, assets.Ether(100).ToInt(), b)
	require.NoError(t, err)
	uni.backend.Commit()

	// Fund Consumer Contract.
	_, err = uni.linkContract.Transfer(uni.sergey, consumerAddress, assets.Ether(100).ToInt())
	require.NoError(t, err)
	uni.backend.Commit()

	// Make the first randomness request.
	numWords := uint32(1)
	requestID, _ := requestRandomnessForWrapper(t, *consumer, uni.neil, keyHash, wrapperSubID, numWords, uint32(callBackGasLimit), uni.rootContract, uni, wrapperOverhead)

	// Wait for simulation to pass.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, requestID, wrapperSubID, uni, db)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, requestID, true, uni.rootContract)

	t.Log("Done!")
}

func TestVRFV2Integration_Wrapper_High_Gas(t *testing.T) {
	t.Parallel()
	wrapperOverhead := uint32(30_000)
	coordinatorOverhead := uint32(90_000)

	key1 := cltest.MustGenerateRandomKey(t)
	callBackGasLimit := int64(2_000_000) // base callback gas.
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, "vrfv2_wrapper_high_gas_revert", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint32](3_500_000)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)

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
		uni,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	wrapper, _, consumer, consumerAddress := deployWrapper(t, uni, wrapperOverhead, coordinatorOverhead, keyHash)

	// Fetch Subscription ID for Wrapper.
	wrapperSubID, err := wrapper.SUBSCRIPTIONID(nil)
	require.NoError(t, err)

	// Fund Subscription.
	b, err := utils.ABIEncode(`[{"type":"uint64"}]`, wrapperSubID)
	require.NoError(t, err)
	_, err = uni.linkContract.TransferAndCall(uni.sergey, uni.rootContractAddress, assets.Ether(100).ToInt(), b)
	require.NoError(t, err)
	uni.backend.Commit()

	// Fund Consumer Contract.
	_, err = uni.linkContract.Transfer(uni.sergey, consumerAddress, assets.Ether(100).ToInt())
	require.NoError(t, err)
	uni.backend.Commit()

	// Make the first randomness request.
	numWords := uint32(1)
	requestID, _ := requestRandomnessForWrapper(t, *consumer, uni.neil, keyHash, wrapperSubID, numWords, uint32(callBackGasLimit), uni.rootContract, uni, wrapperOverhead)

	// Wait for simulation to pass.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, requestID, wrapperSubID, uni, db)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, requestID, true, uni.rootContract)

	t.Log("Done!")
}

func TestVRFV2Integration_SingleConsumer_NeedsBlockhashStore(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 2)
	testMultipleConsumersNeedBHS(
		t,
		ownerKey,
		uni,
		uni.vrfConsumers,
		uni.consumerContracts,
		uni.consumerContractAddresses,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress)
}

func TestVRFV2Integration_SingleConsumer_BlockHeaderFeeder(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	testBlockHeaderFeeder(
		t,
		ownerKey,
		uni,
		uni.vrfConsumers,
		uni.consumerContracts,
		uni.consumerContractAddresses,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress)
}

func TestVRFV2Integration_SingleConsumer_NeedsTopUp(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	testSingleConsumerNeedsTopUp(
		t,
		ownerKey,
		uni,
		uni.vrfConsumers[0],
		uni.consumerContracts[0],
		uni.consumerContractAddresses[0],
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		assets.Ether(1).ToInt(),   // initial funding of 1 LINK
		assets.Ether(100).ToInt(), // top up of 100 LINK
	)
}

func TestVRFV2Integration_SingleConsumer_BigGasCallback_Sandwich(t *testing.T) {
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
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]

	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, assets.Ether(2).ToInt(), uni.rootContract, uni)

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
		uni,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make some randomness requests, each one block apart, which contain a single low-gas request sandwiched between two high-gas requests.
	numWords := uint32(2)
	reqIDs := []*big.Int{}
	callbackGasLimits := []uint32{2_500_000, 50_000, 1_500_000}
	for _, limit := range callbackGasLimits {
		requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, limit, uni.rootContract, uni)
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
	mine(t, reqIDs[1], subID, uni, db)

	// Assert the random word was fulfilled
	assertRandomWordsFulfilled(t, reqIDs[1], false, uni.rootContract)

	// Assert that we've still only completed 1 run before adding new requests.
	runs, err = app.PipelineORM().GetAllRuns()
	require.NoError(t, err)
	assert.Equal(t, 1, len(runs))

	// Make some randomness requests, each one block apart, this time without a low-gas request present in the callbackGasLimit slice.
	callbackGasLimits = []uint32{2_500_000, 2_500_000, 2_500_000}
	for _, limit := range callbackGasLimits {
		_, _ = requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, limit, uni.rootContract, uni)
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

func TestVRFV2Integration_SingleConsumer_MultipleGasLanes(t *testing.T) {
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
	})
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, cheapKey, expensiveKey)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), uni.rootContract, uni)

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
		uni,
		false,
		cheapGasLane, expensiveGasLane)
	cheapHash := jbs[0].VRFSpec.PublicKey.MustHash()
	expensiveHash := jbs[1].VRFSpec.PublicKey.MustHash()

	numWords := uint32(20)
	cheapRequestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, cheapHash, subID, numWords, 500_000, uni.rootContract, uni)

	// Wait for fulfillment to be queued for cheap key hash.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("assert 1", "runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, cheapRequestID, subID, uni, db)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, cheapRequestID, true, uni.rootContract)

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)

	expensiveRequestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, expensiveHash, subID, numWords, 500_000, uni.rootContract, uni)

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
	mine(t, expensiveRequestID, subID, uni, db)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, expensiveRequestID, true, uni.rootContract)

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)
}

func TestVRFV2Integration_SingleConsumer_AlwaysRevertingCallback_StillFulfilled(t *testing.T) {
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
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 0)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key)
	consumer := uni.reverter
	consumerContract := uni.revertingConsumerContract
	consumerContractAddress := uni.revertingConsumerContractAddress

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), uni.rootContract, uni)

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
		uni,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the randomness request.
	numWords := uint32(20)
	requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, uni.rootContract, uni)

	// Wait for fulfillment to be queued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, requestID, subID, uni, db)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, requestID, false, uni.rootContract)
	t.Log("Done!")
}

func TestVRFV2Integration_ConsumerProxy_HappyPath(t *testing.T) {
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
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 0)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1, key2)
	consumerOwner := uni.neil
	consumerContract := uni.consumerProxyContract
	consumerContractAddress := uni.consumerProxyContractAddress

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(
		t, consumerContract, consumerOwner, consumerContractAddress,
		assets.Ether(5).ToInt(), uni.rootContract, uni)

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
		uni,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the first randomness request.
	numWords := uint32(20)
	requestID1, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(
		t, consumerContract, consumerOwner, keyHash, subID, numWords, 750_000, uni.rootContract, uni)

	// Wait for fulfillment to be queued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, requestID1, subID, uni, db)

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
		t, consumerContract, consumerOwner, keyHash, subID, numWords, 750_000, uni.rootContract, uni)
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 2
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
	mine(t, requestID2, subID, uni, db)
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

func TestVRFV2Integration_ConsumerProxy_CoordinatorZeroAddress(t *testing.T) {
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 0)

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

func simulatedOverrides(t *testing.T, defaultGasPrice *assets.Wei, ks ...v2.KeySpecific) func(*chainlink.Config, *chainlink.Secrets) {
	return func(c *chainlink.Config, s *chainlink.Secrets) {
		require.Zero(t, testutils.SimulatedChainID.Cmp(c.EVM[0].ChainID.ToInt()))
		c.EVM[0].GasEstimator.Mode = ptr("FixedPrice")
		if defaultGasPrice != nil {
			c.EVM[0].GasEstimator.PriceDefault = defaultGasPrice
		}
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint32](2_000_000)

		c.EVM[0].HeadTracker.MaxBufferSize = ptr[uint32](100)
		c.EVM[0].HeadTracker.SamplingInterval = models.MustNewDuration(0) // Head sampling disabled

		c.EVM[0].Transactions.ResendAfterThreshold = models.MustNewDuration(0)
		c.EVM[0].Transactions.ReaperThreshold = models.MustNewDuration(100 * time.Millisecond)

		c.EVM[0].FinalityDepth = ptr[uint32](15)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](1)
		c.EVM[0].MinContractPayment = assets.NewLinkFromJuels(100)
		c.EVM[0].KeySpecific = ks
	}
}

func registerProvingKeyHelper(t *testing.T, uni coordinatorV2Universe, coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface, vrfkey vrfkey.KeyV2) {
	// Register a proving key associated with the VRF job.
	p, err := vrfkey.PublicKey.Point()
	require.NoError(t, err)
	_, err = coordinator.RegisterProvingKey(
		uni.neil, uni.nallory.From, pair(secp256k1.Coordinates(p)))
	require.NoError(t, err)
	uni.backend.Commit()
}

func TestExternalOwnerConsumerExample(t *testing.T) {
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
		vrf_coordinator_v2.DeployVRFCoordinatorV2(
			owner, backend, linkAddress, common.Address{}, common.Address{})
	require.NoError(t, err)
	_, err = coordinator.SetConfig(owner, uint16(1), uint32(10000), 1, 1, big.NewInt(10), vrf_coordinator_v2.VRFCoordinatorV2FeeConfig{
		FulfillmentFlatFeeLinkPPMTier1: 0,
		FulfillmentFlatFeeLinkPPMTier2: 0,
		FulfillmentFlatFeeLinkPPMTier3: 0,
		FulfillmentFlatFeeLinkPPMTier4: 0,
		FulfillmentFlatFeeLinkPPMTier5: 0,
		ReqsForTier2:                   big.NewInt(0),
		ReqsForTier3:                   big.NewInt(0),
		ReqsForTier4:                   big.NewInt(0),
		ReqsForTier5:                   big.NewInt(0),
	})
	require.NoError(t, err)
	backend.Commit()
	consumerAddress, _, consumer, err := vrf_external_sub_owner_example.DeployVRFExternalSubOwnerExample(owner, backend, coordinatorAddress, linkAddress)
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
	_, err = consumer.RequestRandomWords(random, 1, 1, 1, 1, [32]byte{})
	require.Error(t, err)
	_, err = consumer.RequestRandomWords(owner, 1, 1, 1, 1, [32]byte{})
	require.NoError(t, err)

	// Reassign ownership, check that only new owner can request
	_, err = consumer.TransferOwnership(owner, random.From)
	require.NoError(t, err)
	_, err = consumer.RequestRandomWords(owner, 1, 1, 1, 1, [32]byte{})
	require.Error(t, err)
	_, err = consumer.RequestRandomWords(random, 1, 1, 1, 1, [32]byte{})
	require.NoError(t, err)
}

func TestSimpleConsumerExample(t *testing.T) {
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
	coordinatorAddress, _, _, err :=
		vrf_coordinator_v2.DeployVRFCoordinatorV2(
			owner, backend, linkAddress, common.Address{}, common.Address{})
	require.NoError(t, err)
	backend.Commit()
	consumerAddress, _, consumer, err := vrf_single_consumer_example.DeployVRFSingleConsumerExample(owner, backend, coordinatorAddress, linkAddress, 1, 1, 1, [32]byte{})
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

func TestIntegrationVRFV2(t *testing.T) {
	t.Parallel()
	// Reconfigure the sim chain with a default gas price of 1 gwei,
	// max gas limit of 2M and a key specific max 10 gwei price.
	// Keep the prices low so we can operate with small link balance subscriptions.
	gasPrice := assets.GWei(1)
	key := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, _ := heavyweight.FullTestDBV2(t, "vrf_v2_integration", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, gasPrice, v2.KeySpecific{
			Key:          &key.EIP55Address,
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	uni := newVRFCoordinatorV2Universe(t, key, 1)
	carol := uni.vrfConsumers[0]
	carolContract := uni.consumerContracts[0]
	carolContractAddress := uni.consumerContractAddresses[0]

	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, key)
	keys, err := app.KeyStore.Eth().EnabledKeysForChain(testutils.SimulatedChainID)
	require.NoError(t, err)
	require.Zero(t, key.Cmp(keys[0]))

	require.NoError(t, app.Start(testutils.Context(t)))

	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		uni,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Create and fund a subscription.
	// We should see that our subscription has 1 link.
	AssertLinkBalances(t, uni.linkContract, []common.Address{
		carolContractAddress,
		uni.rootContractAddress,
	}, []*big.Int{
		assets.Ether(500).ToInt(), // 500 link
		big.NewInt(0),             // 0 link
	})
	subFunding := decimal.RequireFromString("1000000000000000000")
	_, err = carolContract.CreateSubscriptionAndFund(carol,
		subFunding.BigInt())
	require.NoError(t, err)
	uni.backend.Commit()
	AssertLinkBalances(t, uni.linkContract, []common.Address{
		carolContractAddress,
		uni.rootContractAddress,
		uni.nallory.From, // Oracle's own address should have nothing
	}, []*big.Int{
		assets.Ether(499).ToInt(),
		assets.Ether(1).ToInt(),
		big.NewInt(0),
	})
	subId, err := carolContract.SSubId(nil)
	require.NoError(t, err)
	subStart, err := uni.rootContract.GetSubscription(nil, subId)
	require.NoError(t, err)

	// Make a request for random words.
	// By requesting 500k callback with a configured eth gas limit default of 500k,
	// we ensure that the job is indeed adjusting the gaslimit to suit the users request.
	gasRequested := 500_000
	nw := 10
	requestedIncomingConfs := 3
	_, err = carolContract.RequestRandomness(carol, keyHash, subId, uint16(requestedIncomingConfs), uint32(gasRequested), uint32(nw))
	require.NoError(t, err)

	// Oracle tries to withdraw before its fulfilled should fail
	_, err = uni.rootContract.OracleWithdraw(uni.nallory, uni.nallory.From, big.NewInt(1000))
	require.Error(t, err)

	for i := 0; i < requestedIncomingConfs; i++ {
		uni.backend.Commit()
	}

	// We expect the request to be serviced
	// by the node.
	var runs []pipeline.Run
	gomega.NewWithT(t).Eventually(func() bool {
		runs, err = app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		// It is possible that we send the test request
		// before the job spawner has started the vrf services, which is fine
		// the lb will backfill the logs. However, we need to
		// keep blocks coming in for the lb to send the backfilled logs.
		uni.backend.Commit()
		return len(runs) == 1 && runs[0].State == pipeline.RunStatusCompleted
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

	// Wait for the request to be fulfilled on-chain.
	var rf []*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled
	gomega.NewWithT(t).Eventually(func() bool {
		rfIterator, err2 := uni.rootContract.FilterRandomWordsFulfilled(nil, nil)
		require.NoError(t, err2, "failed to logs")
		uni.backend.Commit()
		for rfIterator.Next() {
			rf = append(rf, rfIterator.Event)
		}
		return len(rf) == 1
	}, testutils.WaitTimeout(t), 500*time.Millisecond).Should(gomega.BeTrue())
	assert.True(t, rf[0].Success, "expected callback to succeed")
	fulfillReceipt, err := uni.backend.TransactionReceipt(testutils.Context(t), rf[0].Raw.TxHash)
	require.NoError(t, err)

	// Assert all the random words received by the consumer are different and non-zero.
	seen := make(map[string]struct{})
	var rw *big.Int
	for i := 0; i < nw; i++ {
		rw, err = carolContract.SRandomWords(nil, big.NewInt(int64(i)))
		require.NoError(t, err)
		_, ok := seen[rw.String()]
		assert.False(t, ok)
		seen[rw.String()] = struct{}{}
	}

	// We should have exactly as much gas as we requested
	// after accounting for function look up code, argument decoding etc.
	// which should be fixed in this test.
	ga, err := carolContract.SGasAvailable(nil)
	require.NoError(t, err)
	gaDecoding := big.NewInt(0).Add(ga, big.NewInt(3701))
	assert.Equal(t, 0, gaDecoding.Cmp(big.NewInt(int64(gasRequested))), "expected gas available %v to exceed gas requested %v", gaDecoding, gasRequested)
	t.Log("gas available", ga.String())

	// Assert that we were only charged for how much gas we actually used.
	// We should be charged for the verification + our callbacks execution in link.
	subEnd, err := uni.rootContract.GetSubscription(nil, subId)
	require.NoError(t, err)
	var (
		end   = decimal.RequireFromString(subEnd.Balance.String())
		start = decimal.RequireFromString(subStart.Balance.String())
		wei   = decimal.RequireFromString("1000000000000000000")
		gwei  = decimal.RequireFromString("1000000000")
	)
	t.Log("end balance", end)
	linkWeiCharged := start.Sub(end)
	// Remove flat fee of 0.001 to get fee for just gas.
	linkCharged := linkWeiCharged.Sub(decimal.RequireFromString("1000000000000000")).Div(wei)
	gasPriceD := decimal.NewFromBigInt(gasPrice.ToInt(), 0)
	t.Logf("subscription charged %s with gas prices of %s gwei and %s ETH per LINK\n", linkCharged, gasPriceD.Div(gwei), weiPerUnitLink.Div(wei))
	expected := decimal.RequireFromString(strconv.Itoa(int(fulfillReceipt.GasUsed))).Mul(gasPriceD).Div(weiPerUnitLink)
	t.Logf("expected sub charge gas use %v %v off by %v", fulfillReceipt.GasUsed, expected, expected.Sub(linkCharged))
	// The expected sub charge should be within 200 gas of the actual gas usage.
	// wei/link * link / wei/gas = wei / (wei/gas) = gas
	gasDiff := linkCharged.Sub(expected).Mul(weiPerUnitLink).Div(gasPriceD).Abs().IntPart()
	t.Log("gasDiff", gasDiff)
	assert.Less(t, gasDiff, int64(200))

	// If the oracle tries to withdraw more than it was paid it should fail.
	_, err = uni.rootContract.OracleWithdraw(uni.nallory, uni.nallory.From, linkWeiCharged.Add(decimal.NewFromInt(1)).BigInt())
	require.Error(t, err)

	// Assert the oracle can withdraw its payment.
	_, err = uni.rootContract.OracleWithdraw(uni.nallory, uni.nallory.From, linkWeiCharged.BigInt())
	require.NoError(t, err)
	uni.backend.Commit()
	AssertLinkBalances(t, uni.linkContract, []common.Address{
		carolContractAddress,
		uni.rootContractAddress,
		uni.nallory.From, // Oracle's own address should have nothing
	}, []*big.Int{
		assets.Ether(499).ToInt(),
		subFunding.Sub(linkWeiCharged).BigInt(),
		linkWeiCharged.BigInt(),
	})

	// We should see the response count present
	chain, err := app.Chains.EVM.Get(big.NewInt(1337))
	require.NoError(t, err)

	q := pg.NewQ(app.GetSqlxDB(), app.Logger, app.Config)
	counts := vrf.GetStartingResponseCountsV2(q, app.Logger, chain.Client().ConfiguredChainID().Uint64(), chain.Config().EvmFinalityDepth())
	t.Log(counts, rf[0].RequestId.String())
	assert.Equal(t, uint64(1), counts[rf[0].RequestId.String()])
}

func TestMaliciousConsumer(t *testing.T) {
	t.Parallel()
	config, _ := heavyweight.FullTestDBV2(t, "vrf_v2_integration_malicious", func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint32](2_000_000)
		c.EVM[0].GasEstimator.PriceMax = assets.GWei(1)
		c.EVM[0].GasEstimator.PriceDefault = assets.GWei(1)
		c.EVM[0].GasEstimator.FeeCapDefault = assets.GWei(1)
	})
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, key, 1)
	carol := uni.vrfConsumers[0]

	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, key)
	require.NoError(t, app.Start(testutils.Context(t)))

	err := app.GetKeyStore().Unlock(cltest.Password)
	require.NoError(t, err)
	vrfkey, err := app.GetKeyStore().VRF().Create()
	require.NoError(t, err)

	jid := uuid.NewV4()
	incomingConfs := 2
	s := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		JobID:                    jid.String(),
		Name:                     "vrf-primary",
		FromAddresses:            []string{key.Address.String()},
		CoordinatorAddress:       uni.rootContractAddress.String(),
		BatchCoordinatorAddress:  uni.batchCoordinatorContractAddress.String(),
		MinIncomingConfirmations: incomingConfs,
		GasLanePrice:             assets.GWei(1),
		PublicKey:                vrfkey.PublicKey.String(),
		V2:                       true,
	}).Toml()
	jb, err := vrf.ValidatedVRFSpec(s)
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

	_, err = uni.maliciousConsumerContract.SetKeyHash(carol,
		vrfkey.PublicKey.MustHash())
	require.NoError(t, err)
	subFunding := decimal.RequireFromString("1000000000000000000")
	_, err = uni.maliciousConsumerContract.CreateSubscriptionAndFund(carol,
		subFunding.BigInt())
	require.NoError(t, err)
	uni.backend.Commit()

	// Send a re-entrant request
	_, err = uni.maliciousConsumerContract.RequestRandomness(carol)
	require.NoError(t, err)

	// We expect the request to be serviced
	// by the node.
	var attempts []txmgr.EvmTxAttempt
	gomega.NewWithT(t).Eventually(func() bool {
		//runs, err = app.PipelineORM().GetAllRuns()
		attempts, _, err = app.TxmStorageService().EthTxAttempts(0, 1000)
		require.NoError(t, err)
		// It possible that we send the test request
		// before the job spawner has started the vrf services, which is fine
		// the lb will backfill the logs. However we need to
		// keep blocks coming in for the lb to send the backfilled logs.
		t.Log("attempts", attempts)
		uni.backend.Commit()
		return len(attempts) == 1 && attempts[0].EthTx.State == txmgr.EthTxConfirmed
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
	var fulfillments []*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled
	for it.Next() {
		fulfillments = append(fulfillments, it.Event)
	}
	require.Equal(t, 1, len(fulfillments))
	require.Equal(t, false, fulfillments[0].Success)

	// It should not have succeeded in placing another request.
	it2, err2 := uni.rootContract.FilterRandomWordsRequested(nil, nil, nil, nil)
	require.NoError(t, err2)
	var requests []*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested
	for it2.Next() {
		requests = append(requests, it2.Event)
	}
	require.Equal(t, 1, len(requests))
}

func TestRequestCost(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, key, 1)

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
		subId, err := carolContract.SSubId(nil)
		require.NoError(tt, err)
		// Ensure even with large number of consumers its still cheap
		var addrs []common.Address
		for i := 0; i < 99; i++ {
			addrs = append(addrs, testutils.NewAddress())
		}
		_, err = carolContract.UpdateSubscription(carol, addrs)
		require.NoError(tt, err)
		estimate := estimateGas(tt, uni.backend, common.Address{},
			carolContractAddress, uni.consumerABI,
			"requestRandomness", vrfkey.PublicKey.MustHash(), subId, uint16(2), uint32(10000), uint32(1))
		tt.Log("gas estimate of non-proxied testRequestRandomness:", estimate)
		// V2 should be at least (87000-134000)/134000 = 35% cheaper
		// Note that a second call drops further to 68998 gas, but would also drop in V1.
		assert.Less(tt, estimate, uint64(90_000),
			"requestRandomness tx gas cost more than expected")
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

		theAbi := evmtypes.MustGetABI(vrf_consumer_v2_upgradeable_example.VRFConsumerV2UpgradeableExampleMetaData.ABI)
		estimate := estimateGas(tt, uni.backend, common.Address{},
			consumerContractAddress, &theAbi,
			"requestRandomness", vrfkey.PublicKey.MustHash(), subId, uint16(2), uint32(10000), uint32(1))
		tt.Log("gas estimate of proxied requestRandomness:", estimate)
		// There is some gas overhead of the delegatecall that is made by the proxy
		// to the logic contract. See https://www.evm.codes/#f4?fork=grayGlacier for a detailed
		// breakdown of the gas costs of a delegatecall.
		assert.Less(tt, estimate, uint64(96_000),
			"proxied testRequestRandomness tx gas cost more than expected")
	})
}

func TestMaxConsumersCost(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, key, 1)
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
	assert.Less(t, estimate, uint64(310000))
	estimate = estimateGas(t, uni.backend, carolContractAddress,
		uni.rootContractAddress, uni.coordinatorABI,
		"addConsumer", subId, testutils.NewAddress())
	t.Log(estimate)
	assert.Less(t, estimate, uint64(100000))
}

func TestFulfillmentCost(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, key, 1)

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
		_, err = carolContract.RequestRandomness(carol, vrfkey.PublicKey.MustHash(), subId, uint16(requestedIncomingConfs), uint32(gasRequested), uint32(nw))
		require.NoError(t, err)
		for i := 0; i < requestedIncomingConfs; i++ {
			uni.backend.Commit()
		}

		requestLog := FindLatestRandomnessRequestedLog(tt, uni.rootContract, vrfkey.PublicKey.MustHash())
		s, err := proof.BigToSeed(requestLog.PreSeed)
		require.NoError(t, err)
		proof, rc, err := proof.GenerateProofResponseV2(app.GetKeyStore().VRF(), vrfkey.ID(), proof.PreSeedDataV2{
			PreSeed:          s,
			BlockHash:        requestLog.Raw.BlockHash,
			BlockNum:         requestLog.Raw.BlockNumber,
			SubId:            subId,
			CallbackGasLimit: uint32(gasRequested),
			NumWords:         uint32(nw),
			Sender:           carolContractAddress,
		})
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
		_, err = consumerContract.RequestRandomness(consumerOwner, vrfkey.PublicKey.MustHash(), subId, uint16(requestedIncomingConfs), uint32(gasRequested), uint32(nw))
		require.NoError(t, err)
		for i := 0; i < requestedIncomingConfs; i++ {
			uni.backend.Commit()
		}

		requestLog := FindLatestRandomnessRequestedLog(t, uni.rootContract, vrfkey.PublicKey.MustHash())
		require.Equal(tt, subId, requestLog.SubId)
		s, err := proof.BigToSeed(requestLog.PreSeed)
		require.NoError(t, err)
		proof, rc, err := proof.GenerateProofResponseV2(app.GetKeyStore().VRF(), vrfkey.ID(), proof.PreSeedDataV2{
			PreSeed:          s,
			BlockHash:        requestLog.Raw.BlockHash,
			BlockNum:         requestLog.Raw.BlockNumber,
			SubId:            subId,
			CallbackGasLimit: uint32(gasRequested),
			NumWords:         uint32(nw),
			Sender:           consumerContractAddress,
		})
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

func TestStartingCountsV1(t *testing.T) {
	cfg, db := heavyweight.FullTestDBNoFixturesV2(t, "vrf_test_starting_counts", nil)
	_, err := db.Exec(`INSERT INTO evm_chains (id, created_at, updated_at) VALUES (1337, NOW(), NOW())`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO evm_heads (hash, number, parent_hash, created_at, timestamp, evm_chain_id)
	VALUES ($1, 4, $2, NOW(), NOW(), 1337)`, utils.NewHash(), utils.NewHash())
	require.NoError(t, err)

	lggr := logger.TestLogger(t)
	q := pg.NewQ(db, lggr, cfg)
	finalityDepth := 3
	counts := vrf.GetStartingResponseCountsV1(q, lggr, 1337, uint32(finalityDepth))
	assert.Equal(t, 0, len(counts))
	ks := keystore.New(db, utils.FastScryptParams, lggr, cfg)
	err = ks.Unlock(testutils.Password)
	require.NoError(t, err)
	k, err := ks.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)
	b := time.Now()
	n1, n2, n3, n4 := int64(0), int64(1), int64(2), int64(3)
	reqID := utils.PadByteToHash(0x10)
	m1 := txmgr.EthTxMeta{
		RequestID: &reqID,
	}
	md1, err := json.Marshal(&m1)
	require.NoError(t, err)
	md1_ := datatypes.JSON(md1)
	reqID2 := utils.PadByteToHash(0x11)
	m2 := txmgr.EthTxMeta{
		RequestID: &reqID2,
	}
	md2, err := json.Marshal(&m2)
	md2_ := datatypes.JSON(md2)
	require.NoError(t, err)
	chainID := utils.NewBig(big.NewInt(1337))
	confirmedTxes := []txmgr.EvmTx{
		{
			Nonce:              &n1,
			FromAddress:        k.Address,
			Error:              null.String{},
			BroadcastAt:        &b,
			InitialBroadcastAt: &b,
			CreatedAt:          b,
			State:              txmgr.EthTxConfirmed,
			Meta:               &datatypes.JSON{},
			EncodedPayload:     []byte{},
			EVMChainID:         *chainID,
		},
		{
			Nonce:              &n2,
			FromAddress:        k.Address,
			Error:              null.String{},
			BroadcastAt:        &b,
			InitialBroadcastAt: &b,
			CreatedAt:          b,
			State:              txmgr.EthTxConfirmed,
			Meta:               &md1_,
			EncodedPayload:     []byte{},
			EVMChainID:         *chainID,
		},
		{
			Nonce:              &n3,
			FromAddress:        k.Address,
			Error:              null.String{},
			BroadcastAt:        &b,
			InitialBroadcastAt: &b,
			CreatedAt:          b,
			State:              txmgr.EthTxConfirmed,
			Meta:               &md2_,
			EncodedPayload:     []byte{},
			EVMChainID:         *chainID,
		},
		{
			Nonce:              &n4,
			FromAddress:        k.Address,
			Error:              null.String{},
			BroadcastAt:        &b,
			InitialBroadcastAt: &b,
			CreatedAt:          b,
			State:              txmgr.EthTxConfirmed,
			Meta:               &md2_,
			EncodedPayload:     []byte{},
			EVMChainID:         *chainID,
		},
	}
	// add unconfirmed txes
	unconfirmedTxes := []txmgr.EvmTx{}
	for i := int64(4); i < 6; i++ {
		reqID3 := utils.PadByteToHash(0x12)
		md, err := json.Marshal(&txmgr.EthTxMeta{
			RequestID: &reqID3,
		})
		require.NoError(t, err)
		md1 := datatypes.JSON(md)
		newNonce := i + 1
		unconfirmedTxes = append(unconfirmedTxes, txmgr.EvmTx{
			Nonce:              &newNonce,
			FromAddress:        k.Address,
			Error:              null.String{},
			CreatedAt:          b,
			State:              txmgr.EthTxUnconfirmed,
			BroadcastAt:        &b,
			InitialBroadcastAt: &b,
			Meta:               &md1,
			EncodedPayload:     []byte{},
			EVMChainID:         *chainID,
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
	txAttempts := []txmgr.EvmTxAttempt{}
	for i := range confirmedTxes {
		txAttempts = append(txAttempts, txmgr.EvmTxAttempt{
			EthTxID:                 int64(i + 1),
			GasPrice:                assets.NewWeiI(100),
			SignedRawTx:             []byte(`blah`),
			Hash:                    utils.NewHash(),
			BroadcastBeforeBlockNum: &broadcastBlock,
			State:                   txmgrtypes.TxAttemptBroadcast,
			CreatedAt:               time.Now(),
			ChainSpecificGasLimit:   uint32(100),
		})
	}
	// add eth_tx_attempts for unconfirmed
	for i := range unconfirmedTxes {
		txAttempts = append(txAttempts, txmgr.EvmTxAttempt{
			EthTxID:               int64(i + 1 + len(confirmedTxes)),
			GasPrice:              assets.NewWeiI(100),
			SignedRawTx:           []byte(`blah`),
			Hash:                  utils.NewHash(),
			State:                 txmgrtypes.TxAttemptInProgress,
			CreatedAt:             time.Now(),
			ChainSpecificGasLimit: uint32(100),
		})
	}
	for _, txAttempt := range txAttempts {
		t.Log("tx attempt eth tx id: ", txAttempt.EthTxID)
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
	receipts := []txmgr.EvmReceipt{}
	for i := 0; i < 4; i++ {
		receipts = append(receipts, txmgr.EvmReceipt{
			BlockHash:        utils.NewHash(),
			TxHash:           txAttempts[i].Hash,
			BlockNumber:      broadcastBlock,
			TransactionIndex: 1,
			Receipt:          &evmtypes.Receipt{},
			CreatedAt:        time.Now(),
		})
	}
	sql = `INSERT INTO eth_receipts (block_hash, tx_hash, block_number, transaction_index, receipt, created_at)
		VALUES (:block_hash, :tx_hash, :block_number, :transaction_index, :receipt, :created_at)`
	for _, r := range receipts {
		dbReceipt := txmgr.DbReceiptFromEvmReceipt(&r)
		_, err := db.NamedExec(sql, &dbReceipt)
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

func FindLatestRandomnessRequestedLog(t *testing.T,
	coordContract *vrf_coordinator_v2.VRFCoordinatorV2,
	keyHash [32]byte) *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested {
	var rf []*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested
	gomega.NewWithT(t).Eventually(func() bool {
		rfIterator, err2 := coordContract.FilterRandomWordsRequested(nil, [][32]byte{keyHash}, nil, []common.Address{})
		require.NoError(t, err2, "failed to logs")
		for rfIterator.Next() {
			rf = append(rf, rfIterator.Event)
		}
		return len(rf) >= 1
	}, testutils.WaitTimeout(t), 500*time.Millisecond).Should(gomega.BeTrue())
	latest := len(rf) - 1
	return rf[latest]
}

func AssertLinkBalances(t *testing.T, linkContract *link_token_interface.LinkToken, addresses []common.Address, balances []*big.Int) {
	require.Equal(t, len(addresses), len(balances))
	for i, a := range addresses {
		b, err := linkContract.BalanceOf(nil, a)
		require.NoError(t, err)
		assert.Equal(t, balances[i].String(), b.String(), "invalid balance for %v", a)
	}
}

func ptr[T any](t T) *T { return &t }
