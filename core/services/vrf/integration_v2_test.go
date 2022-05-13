package vrf_test

import (
	"context"
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

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/batch_vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_consumer_v2"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_external_sub_owner_example"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_malicious_consumer_v2"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_single_consumer_example"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrfv2_reverting_example"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/services/pg/datatypes"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/services/vrf/proof"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// vrfConsumerContract is the common interface implemented by
// the example contracts used for the integration tests.
type vrfConsumerContract interface {
	TestCreateSubscriptionAndFund(opts *bind.TransactOpts, fundingJuels *big.Int) (*gethtypes.Transaction, error)
	SSubId(opts *bind.CallOpts) (uint64, error)
	SRequestId(opts *bind.CallOpts) (*big.Int, error)
	TestRequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*gethtypes.Transaction, error)
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
	bhsContract                      *blockhash_store.BlockhashStore
	bhsContractAddress               common.Address
	maliciousConsumerContract        *vrf_malicious_consumer_v2.VRFMaliciousConsumerV2
	maliciousConsumerContractAddress common.Address
	revertingConsumerContract        *vrfv2_reverting_example.VRFV2RevertingExample
	revertingConsumerContractAddress common.Address

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
	oracleTransactor := cltest.MustNewSimulatedBackendKeyedTransactor(t, key.ToEcdsaPrivKey())
	var (
		sergey       = newIdentity(t)
		neil         = newIdentity(t)
		ned          = newIdentity(t)
		evil         = newIdentity(t)
		reverter     = newIdentity(t)
		nallory      = oracleTransactor
		vrfConsumers = []*bind.TransactOpts{}
	)

	// Create consumer contract deployer identities
	for i := 0; i < numConsumers; i++ {
		vrfConsumers = append(vrfConsumers, newIdentity(t))
	}

	genesisData := core.GenesisAlloc{
		sergey.From:   {Balance: assets.Ether(1000)},
		neil.From:     {Balance: assets.Ether(1000)},
		ned.From:      {Balance: assets.Ether(1000)},
		nallory.From:  {Balance: assets.Ether(1000)},
		evil.From:     {Balance: assets.Ether(1000)},
		reverter.From: {Balance: assets.Ether(1000)},
	}
	for _, consumer := range vrfConsumers {
		genesisData[consumer.From] = core.GenesisAccount{
			Balance: assets.Ether(1000),
		}
	}

	gasLimit := ethconfig.Defaults.Miner.GasCeil
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
	consumerContracts := []*vrf_consumer_v2.VRFConsumerV2{}
	consumerContractAddresses := []common.Address{}
	for _, author := range vrfConsumers {
		// Deploy a VRF consumer. It has a starting balance of 500 LINK.
		consumerContractAddress, _, consumerContract, err :=
			vrf_consumer_v2.DeployVRFConsumerV2(
				author, backend, coordinatorAddress, linkAddress)
		require.NoError(t, err, "failed to deploy VRFConsumer contract to simulated ethereum blockchain")
		_, err = linkContract.Transfer(sergey, consumerContractAddress, assets.Ether(500)) // Actually, LINK
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
	_, err = linkContract.Transfer(sergey, maliciousConsumerContractAddress, assets.Ether(1)) // Actually, LINK
	require.NoError(t, err, "failed to send LINK to VRFMaliciousConsumer contract on simulated ethereum blockchain")
	backend.Commit()

	// Deploy always reverting consumer
	revertingConsumerContractAddress, _, revertingConsumerContract, err := vrfv2_reverting_example.DeployVRFV2RevertingExample(
		reverter, backend, coordinatorAddress, linkAddress,
	)
	require.NoError(t, err, "failed to deploy VRFRevertingExample contract to simulated eth blockchain")
	_, err = linkContract.Transfer(sergey, revertingConsumerContractAddress, assets.Ether(500)) // Actually, LINK
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

		rootContract:                     coordinatorContract,
		rootContractAddress:              coordinatorAddress,
		linkContract:                     linkContract,
		linkContractAddress:              linkAddress,
		bhsContract:                      bhsContract,
		bhsContractAddress:               bhsAddress,
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
	nonce, err := ec.PendingNonceAt(context.Background(), key.Address.Address())
	require.NoError(t, err)
	tx := gethtypes.NewTx(&gethtypes.DynamicFeeTx{
		ChainID:   big.NewInt(1337),
		Nonce:     nonce,
		GasTipCap: big.NewInt(1),
		GasFeeCap: big.NewInt(10e9), // block base fee in sim
		Gas:       uint64(21_000),
		To:        &to,
		Value:     big.NewInt(0).Mul(big.NewInt(int64(eth)), big.NewInt(1e18)),
		Data:      nil,
	})
	signedTx, err := gethtypes.SignTx(tx, gethtypes.NewLondonSigner(big.NewInt(1337)), key.ToEcdsaPrivKey())
	require.NoError(t, err)
	err = ec.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)
	ec.Commit()
}

func subscribeVRF(
	t *testing.T,
	author *bind.TransactOpts,
	consumerContract vrfConsumerContract,
	coordinatorContract *vrf_coordinator_v2.VRFCoordinatorV2,
	backend *backends.SimulatedBackend,
	fundingJuels *big.Int,
) (vrf_coordinator_v2.GetSubscription, uint64) {
	_, err := consumerContract.TestCreateSubscriptionAndFund(author, fundingJuels)
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
	uni coordinatorV2Universe,
	batchEnabled bool,
) (jobs []job.Job) {
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
			CoordinatorAddress:       uni.rootContractAddress.String(),
			BatchCoordinatorAddress:  uni.batchCoordinatorContractAddress.String(),
			BatchFulfillmentEnabled:  batchEnabled,
			MinIncomingConfirmations: incomingConfs,
			PublicKey:                vrfkey.PublicKey.String(),
			FromAddresses:            keyStrs,
			BackoffInitialDelay:      10 * time.Millisecond,
			BackoffMaxDelay:          time.Second,
			V2:                       true,
		}).Toml()
		jb, err := vrf.ValidatedVRFSpec(s)
		t.Log(jb.VRFSpec.PublicKey.MustHash(), vrfkey.PublicKey.MustHash())
		require.NoError(t, err)
		err = app.JobSpawner().CreateJob(&jb)
		require.NoError(t, err)
		registerProvingKeyHelper(t, uni, vrfkey)
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
	}, cltest.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())
	// Unfortunately the lb needs heads to be able to backfill logs to new subscribers.
	// To avoid confirming
	// TODO: it could just backfill immediately upon receiving a new subscriber? (though would
	// only be useful for tests, probably a more robust way is to have the job spawner accept a signal that a
	// job is fully up and running and not add it to the active jobs list before then)
	time.Sleep(2 * time.Second)

	return
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
	uni coordinatorV2Universe,
) (*big.Int, uint64) {
	minRequestConfirmations := uint16(2)
	_, err := vrfConsumerHandle.TestRequestRandomness(
		consumerOwner,
		keyHash,
		subID,
		minRequestConfirmations,
		cbGasLimit,
		numWords,
	)
	require.NoError(t, err)

	uni.backend.Commit()

	iter, err := uni.rootContract.FilterRandomWordsRequested(nil, nil, []uint64{subID}, nil)
	require.NoError(t, err, "could not filter RandomWordsRequested events")

	events := []*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{}
	for iter.Next() {
		events = append(events, iter.Event)
	}

	requestID, err := vrfConsumerHandle.SRequestId(nil)
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
	uni coordinatorV2Universe,
) uint64 {
	// Create a subscription and fund with LINK.
	sub, subID := subscribeVRF(t, consumerOwner, vrfConsumerHandle, uni.rootContract, uni.backend, fundingJuels)
	require.Equal(t, uint64(1), subID)
	require.Equal(t, fundingJuels.String(), sub.Balance.String())

	// Assert the subscription event in the coordinator contract.
	iter, err := uni.rootContract.FilterSubscriptionCreated(nil, []uint64{subID})
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
	uni coordinatorV2Universe,
) {
	// Check many times in case there are delays processing the event
	// this could happen occasionally and cause flaky tests.
	numChecks := 3
	found := false
	for i := 0; i < numChecks; i++ {
		filter, err := uni.rootContract.FilterRandomWordsFulfilled(nil, []*big.Int{requestID})
		require.NoError(t, err)

		for filter.Next() {
			require.Equal(t, expectedSuccess, filter.Event.Success, "fulfillment event success not correct, expected: %+v, actual: %+v", expectedSuccess, filter.Event.Success)
			require.Equal(t, requestID, filter.Event.RequestId)
			found = true
		}

		if found {
			break
		}

		// Wait a bit and try again.
		time.Sleep(time.Second)
	}
	require.True(t, found, "RandomWordsFulfilled event not found")
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
		var txs []txmgr.EthTx
		err := db.Select(&txs, `
		SELECT * FROM eth_txes
		WHERE eth_txes.state = 'confirmed'
			AND eth_txes.meta->>'RequestID' = $1
			AND CAST(eth_txes.meta->>'SubId' AS NUMERIC) = $2 LIMIT 1
		`, common.BytesToHash(requestID.Bytes()).String(), subID)
		require.NoError(t, err)
		t.Log("num txs", len(txs))
		return len(txs) == 1
	}, cltest.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
}

func mineBatch(t *testing.T, requestIDs []*big.Int, subID uint64, uni coordinatorV2Universe, db *sqlx.DB) bool {
	requestIDMap := map[string]bool{}
	for _, requestID := range requestIDs {
		requestIDMap[common.BytesToHash(requestID.Bytes()).String()] = false
	}
	return gomega.NewWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		var txs []txmgr.EthTx
		err := db.Select(&txs, `
		SELECT * FROM eth_txes
		WHERE eth_txes.state = 'confirmed'
			AND CAST(eth_txes.meta->>'SubId' AS NUMERIC) = $1
		`, subID)
		require.NoError(t, err)
		for _, tx := range txs {
			meta, err := tx.GetMeta()
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
	}, cltest.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
}

func TestVRFV2Integration_SingleConsumer_HappyPath_BatchFulfillment(t *testing.T) {
	config, db := heavyweight.FullTestDB(t, "vrfv2_singleconsumer_batch_happypath")
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey)
	config.Overrides.GlobalEvmGasLimitDefault = null.NewInt(5e6, true)
	config.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(2)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), uni)

	// Create gas lane.
	key1, err := app.KeyStore.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)
	sendEth(t, ownerKey, uni.backend, key1.Address.Address(), 10)
	configureSimChain(t, app, map[string]types.ChainCfg{
		key1.Address.String(): {
			EvmMaxGasPriceWei: utils.NewBig(big.NewInt(10e9)), // 10 gwei
		},
	}, big.NewInt(10e9))
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job using key1 and key2 on the same gas lane.
	jbs := createVRFJobs(t, [][]ethkey.KeyV2{{key1}}, app, uni, true)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make some randomness requests.
	numWords := uint32(2)
	reqIDs := []*big.Int{}
	for i := 0; i < 5; i++ {
		requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, uni)
		reqIDs = append(reqIDs, requestID)
	}

	// Wait for fulfillment to be queued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 5
	}, cltest.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	mineBatch(t, reqIDs, subID, uni, db)

	for i, requestID := range reqIDs {
		// Assert correct state of RandomWordsFulfilled event.
		// The last request will be the successful one because of the way the example
		// contract is written.
		if i == (len(reqIDs) - 1) {
			assertRandomWordsFulfilled(t, requestID, true, uni)
		} else {
			assertRandomWordsFulfilled(t, requestID, false, uni)
		}
	}

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)
}

func TestVRFV2Integration_SingleConsumer_HappyPath_BatchFulfillment_BigGasCallback(t *testing.T) {
	config, db := heavyweight.FullTestDB(t, "vrfv2_singleconsumer_batch_bigcallback")
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey)
	config.Overrides.GlobalEvmGasLimitDefault = null.NewInt(5e6, true)
	config.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(2)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), uni)

	// Create gas lane.
	key1, err := app.KeyStore.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)
	sendEth(t, ownerKey, uni.backend, key1.Address.Address(), 10)
	configureSimChain(t, app, map[string]types.ChainCfg{
		key1.Address.String(): {
			EvmMaxGasPriceWei: utils.NewBig(big.NewInt(10e9)), // 10 gwei
		},
	}, big.NewInt(10e9))
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job using key1 and key2 on the same gas lane.
	jbs := createVRFJobs(t, [][]ethkey.KeyV2{{key1}}, app, uni, true)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make some randomness requests with low max gas callback limits.
	// These should all be included in the same batch.
	numWords := uint32(2)
	reqIDs := []*big.Int{}
	for i := 0; i < 5; i++ {
		requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 100_000, uni)
		reqIDs = append(reqIDs, requestID)
	}

	// Make one randomness request with the max callback gas limit.
	// It should live in a batch on it's own.
	requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 2_500_000, uni)
	reqIDs = append(reqIDs, requestID)

	// Wait for fulfillment to be queued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 6
	}, cltest.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	mineBatch(t, reqIDs, subID, uni, db)

	for i, requestID := range reqIDs {
		// Assert correct state of RandomWordsFulfilled event.
		// The last request will be the successful one because of the way the example
		// contract is written.
		if i == (len(reqIDs) - 1) {
			assertRandomWordsFulfilled(t, requestID, true, uni)
		} else {
			assertRandomWordsFulfilled(t, requestID, false, uni)
		}
	}

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)
}

func TestVRFV2Integration_SingleConsumer_HappyPath(t *testing.T) {
	config, db := heavyweight.FullTestDB(t, "vrfv2_singleconsumer_happypath")
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey)
	config.Overrides.GlobalEvmGasLimitDefault = null.NewInt(0, false)
	config.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(2)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), uni)

	// Create gas lane.
	key1, err := app.KeyStore.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)
	key2, err := app.KeyStore.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)
	sendEth(t, ownerKey, uni.backend, key1.Address.Address(), 10)
	sendEth(t, ownerKey, uni.backend, key2.Address.Address(), 10)
	configureSimChain(t, app, map[string]types.ChainCfg{
		key1.Address.String(): {
			EvmMaxGasPriceWei: utils.NewBig(big.NewInt(10e9)), // 10 gwei
		},
		key2.Address.String(): {
			EvmMaxGasPriceWei: utils.NewBig(big.NewInt(10e9)), // 10 gwei
		},
	}, big.NewInt(10e9))
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job using key1 and key2 on the same gas lane.
	jbs := createVRFJobs(t, [][]ethkey.KeyV2{{key1, key2}}, app, uni, false)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the first randomness request.
	numWords := uint32(20)
	requestID1, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, uni)

	// Wait for fulfillment to be queued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, cltest.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, requestID1, subID, uni, db)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, requestID1, true, uni)

	// Make the second randomness request and assert fulfillment is successful
	requestID2, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, uni)
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 2
	}, cltest.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
	mine(t, requestID2, subID, uni, db)
	assertRandomWordsFulfilled(t, requestID2, true, uni)

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)

	// Assert that both send addresses were used to fulfill the requests
	n, err := uni.backend.PendingNonceAt(context.Background(), key1.Address.Address())
	require.NoError(t, err)
	require.EqualValues(t, 1, n)

	n, err = uni.backend.PendingNonceAt(context.Background(), key2.Address.Address())
	require.NoError(t, err)
	require.EqualValues(t, 1, n)

	t.Log("Done!")
}

func TestVRFV2Integration_SingleConsumer_NeedsBlockhashStore(t *testing.T) {
	config, db := heavyweight.FullTestDB(t, "vrfv2_needs_blockhash_store")
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey)
	config.Overrides.GlobalEvmGasLimitDefault = null.NewInt(0, false)
	config.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(2)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]

	// Create a subscription and fund with 0 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, new(big.Int), uni)

	// Create gas lane.
	vrfKey, err := app.KeyStore.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)
	sendEth(t, ownerKey, uni.backend, vrfKey.Address.Address(), 10)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create BHS key
	bhsKey, err := app.KeyStore.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)
	sendEth(t, ownerKey, uni.backend, bhsKey.Address.Address(), 10)

	// Configure VRF and BHS keys
	configureSimChain(t, app, map[string]types.ChainCfg{
		vrfKey.Address.String(): {
			EvmMaxGasPriceWei: utils.NewBig(big.NewInt(10e9)), // 10 gwei
		},
		bhsKey.Address.String(): {
			EvmMaxGasPriceWei: utils.NewBig(big.NewInt(10e9)), // 10 gwei
		},
	}, big.NewInt(10e9))

	// Create VRF job.
	vrfJobs := createVRFJobs(t, [][]ethkey.KeyV2{{vrfKey}}, app, uni, false)
	keyHash := vrfJobs[0].VRFSpec.PublicKey.MustHash()

	_ = createAndStartBHSJob(
		t, vrfKey.Address.String(), app, uni.bhsContractAddress.String(), "",
		uni.rootContractAddress.String())

	// Make the randomness request. It will not yet succeed since it is underfunded.
	numWords := uint32(20)
	requestID, requestBlock := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, uni)

	// Wait 101 blocks.
	for i := 0; i < 100; i++ {
		uni.backend.Commit()
	}

	// Wait for the blockhash to be stored
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		_, err := uni.bhsContract.GetBlockhash(&bind.CallOpts{
			Pending:     false,
			From:        common.Address{},
			BlockNumber: nil,
			Context:     nil,
		}, big.NewInt(int64(requestBlock)))
		if err == nil {
			return true
		} else if strings.Contains(err.Error(), "execution reverted") {
			return false
		} else {
			t.Fatal(err)
			return false
		}
	}, cltest.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	// Wait another 160 blocks so that the request is outside of the 256 block window
	for i := 0; i < 160; i++ {
		uni.backend.Commit()
	}

	// Fund the subscription
	_, err = consumerContract.TopUpSubscription(consumer, big.NewInt(5e18 /* 5 LINK */))
	require.NoError(t, err)

	// Wait for fulfillment to be queued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, cltest.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, requestID, subID, uni, db)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, requestID, true, uni)

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)
}

func TestVRFV2Integration_SingleConsumer_NeedsTopUp(t *testing.T) {
	config, db := heavyweight.FullTestDB(t, "vrfv2_singleconsumer_needstopup")
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey)
	config.Overrides.GlobalEvmGasLimitDefault = null.NewInt(0, false)
	config.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(2)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]

	// Create a subscription and fund with 1 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(1e18), uni)

	// Create expensive gas lane.
	key, err := app.KeyStore.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)
	sendEth(t, ownerKey, uni.backend, key.Address.Address(), 10)
	configureSimChain(t, app, map[string]types.ChainCfg{
		key.Address.String(): {
			EvmMaxGasPriceWei: utils.NewBig(big.NewInt(1000e9)), // 1000 gwei
		},
	}, big.NewInt(1000e9))
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job.
	jbs := createVRFJobs(t, [][]ethkey.KeyV2{{key}}, app, uni, false)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	numWords := uint32(20)
	requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, uni)

	// Fulfillment will not be enqueued because subscriber doesn't have enough LINK.
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("assert 1", "runs", len(runs))
		return len(runs) == 0
	}, 5*time.Second, 1*time.Second).Should(gomega.BeTrue())

	// Top up subscription with enough LINK to see the job through. 100 LINK should do the trick.
	_, err = consumerContract.TopUpSubscription(consumer, decimal.RequireFromString("100e18").BigInt())
	require.NoError(t, err)

	// Wait for fulfillment to go through.
	gomega.NewWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("assert 2", "runs", len(runs))
		return len(runs) == 1
	}, cltest.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment. Need to wait for Txm to mark the tx as confirmed
	// so that we can actually see the event on the simulated chain.
	mine(t, requestID, subID, uni, db)

	// Assert the state of the RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, requestID, true, uni)

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)
}

func TestVRFV2Integration_SingleConsumer_BigGasCallback_Sandwich(t *testing.T) {
	config, db := heavyweight.FullTestDB(t, "vrfv2_singleconsumer_bigcallback_sandwich")
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey)
	config.Overrides.GlobalEvmGasLimitDefault = null.NewInt(5e6, true)
	config.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(2)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]

	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, assets.Ether(3), uni)

	// Create gas lane.
	key1, err := app.KeyStore.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)
	sendEth(t, ownerKey, uni.backend, key1.Address.Address(), 10)
	configureSimChain(t, app, map[string]types.ChainCfg{
		key1.Address.String(): {
			EvmMaxGasPriceWei: utils.NewBig(assets.GWei(100)),
		},
	}, assets.GWei(100))
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job.
	jbs := createVRFJobs(t, [][]ethkey.KeyV2{{key1}}, app, uni, false)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make some randomness requests, each one block apart, which contain a single low-gas request sandwiched between two high-gas requests.
	numWords := uint32(2)
	reqIDs := []*big.Int{}
	callbackGasLimits := []uint32{2_500_000, 50_000, 1_500_000}
	for _, limit := range callbackGasLimits {
		requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, limit, uni)
		reqIDs = append(reqIDs, requestID)
		uni.backend.Commit()
	}

	// Assert that we've completed 0 runs before adding 3 new requests.
	runs, err := app.PipelineORM().GetAllRuns()
	assert.Equal(t, 0, len(runs))
	assert.Equal(t, 3, len(reqIDs))

	// Wait for the 50_000 gas randomness request to be enqueued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, cltest.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

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
	assertRandomWordsFulfilled(t, reqIDs[1], false, uni)

	// Assert that we've still only completed 1 run before adding new requests.
	runs, err = app.PipelineORM().GetAllRuns()
	assert.Equal(t, 1, len(runs))

	// Make some randomness requests, each one block apart, this time without a low-gas request present in the callbackGasLimit slice.
	reqIDs = []*big.Int{}
	callbackGasLimits = []uint32{2_500_000, 2_500_000, 2_500_000}
	for _, limit := range callbackGasLimits {
		requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, limit, uni)
		reqIDs = append(reqIDs, requestID)
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
	config, db := heavyweight.FullTestDB(t, "vrfv2_singleconsumer_multiplegaslanes")
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey)
	config.Overrides.GlobalEvmGasLimitDefault = null.NewInt(0, false)
	config.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(2)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), uni)

	// Create cheap gas lane.
	cheapKey, err := app.KeyStore.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)
	sendEth(t, ownerKey, uni.backend, cheapKey.Address.Address(), 10)
	// Create expensive gas lane.
	expensiveKey, err := app.KeyStore.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)
	sendEth(t, ownerKey, uni.backend, expensiveKey.Address.Address(), 10)
	configureSimChain(t, app, map[string]types.ChainCfg{
		cheapKey.Address.String(): {
			EvmMaxGasPriceWei: utils.NewBig(big.NewInt(10e9)), // 10 gwei
		},
		expensiveKey.Address.String(): {
			EvmMaxGasPriceWei: utils.NewBig(big.NewInt(1000e9)), // 1000 gwei
		},
	}, big.NewInt(10e9))
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF jobs.
	jbs := createVRFJobs(t, [][]ethkey.KeyV2{{cheapKey}, {expensiveKey}}, app, uni, false)
	cheapHash := jbs[0].VRFSpec.PublicKey.MustHash()
	expensiveHash := jbs[1].VRFSpec.PublicKey.MustHash()

	numWords := uint32(20)
	cheapRequestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, cheapHash, subID, numWords, 500_000, uni)

	// Wait for fulfillment to be queued for cheap key hash.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("assert 1", "runs", len(runs))
		return len(runs) == 1
	}, cltest.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, cheapRequestID, subID, uni, db)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, cheapRequestID, true, uni)

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)

	expensiveRequestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, expensiveHash, subID, numWords, 500_000, uni)

	// We should not have any new fulfillments until a top up.
	gomega.NewWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("assert 2", "runs", len(runs))
		return len(runs) == 1
	}, 5*time.Second, 1*time.Second).Should(gomega.BeTrue())

	// Top up subscription with enough LINK to see the job through. 100 LINK should do the trick.
	_, err = consumerContract.TopUpSubscription(consumer, decimal.RequireFromString("100e18").BigInt())
	require.NoError(t, err)

	// Wait for fulfillment to be queued for expensive key hash.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("assert 1", "runs", len(runs))
		return len(runs) == 2
	}, cltest.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, expensiveRequestID, subID, uni, db)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, expensiveRequestID, true, uni)

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)
}

func TestVRFV2Integration_SingleConsumer_AlwaysRevertingCallback_StillFulfilled(t *testing.T) {
	config, db := heavyweight.FullTestDB(t, "vrfv2_singleconsumer_alwaysrevertingcallback")
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 0)
	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey)
	config.Overrides.GlobalEvmGasLimitDefault = null.NewInt(0, false)
	config.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(2)
	consumer := uni.reverter
	consumerContract := uni.revertingConsumerContract
	consumerContractAddress := uni.revertingConsumerContractAddress

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), uni)

	// Create gas lane.
	key, err := app.KeyStore.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)
	sendEth(t, ownerKey, uni.backend, key.Address.Address(), 10)
	configureSimChain(t, app, map[string]types.ChainCfg{
		key.Address.String(): {
			EvmMaxGasPriceWei: utils.NewBig(big.NewInt(10e9)), // 10 gwei
		},
	}, big.NewInt(10e9))
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job.
	jbs := createVRFJobs(t, [][]ethkey.KeyV2{{key}}, app, uni, false)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the randomness request.
	numWords := uint32(20)
	requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, uni)

	// Wait for fulfillment to be queued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, cltest.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, requestID, subID, uni, db)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, requestID, false, uni)
	t.Log("Done!")
}

func configureSimChain(t *testing.T, app *cltest.TestApplication, ks map[string]types.ChainCfg, defaultGasPrice *big.Int) {
	zero := models.MustMakeDuration(0 * time.Millisecond)
	reaperThreshold := models.MustMakeDuration(100 * time.Millisecond)
	app.Chains.EVM.Configure(
		testutils.Context(t),
		*utils.NewBigI(1337),
		true,
		&types.ChainCfg{
			GasEstimatorMode:                 null.StringFrom("FixedPrice"),
			EvmGasPriceDefault:               utils.NewBig(defaultGasPrice),
			EvmHeadTrackerMaxBufferSize:      null.IntFrom(100),
			EvmHeadTrackerSamplingInterval:   &zero, // Head sampling disabled
			EthTxResendAfterThreshold:        &zero,
			EvmFinalityDepth:                 null.IntFrom(15),
			EthTxReaperThreshold:             &reaperThreshold,
			MinIncomingConfirmations:         null.IntFrom(1),
			MinRequiredOutgoingConfirmations: null.IntFrom(1),
			MinimumContractPayment:           assets.NewLinkFromJuels(100),
			EvmGasLimitDefault:               null.NewInt(2000000, true),
			KeySpecific:                      ks,
		},
	)
}

func registerProvingKeyHelper(t *testing.T, uni coordinatorV2Universe, vrfkey vrfkey.KeyV2) {
	// Register a proving key associated with the VRF job.
	p, err := vrfkey.PublicKey.Point()
	require.NoError(t, err)
	_, err = uni.rootContract.RegisterProvingKey(
		uni.neil, uni.nallory.From, pair(secp256k1.Coordinates(p)))
	require.NoError(t, err)
	uni.backend.Commit()
}

func TestExternalOwnerConsumerExample(t *testing.T) {
	owner := newIdentity(t)
	random := newIdentity(t)
	genesisData := core.GenesisAlloc{
		owner.From:  {Balance: assets.Ether(10)},
		random.From: {Balance: assets.Ether(10)},
	}
	backend := cltest.NewSimulatedBackend(t, genesisData, ethconfig.Defaults.Miner.GasCeil)
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
	_, err = linkContract.Transfer(owner, consumerAddress, assets.Ether(2))
	require.NoError(t, err)
	backend.Commit()
	AssertLinkBalances(t, linkContract, []common.Address{owner.From, consumerAddress}, []*big.Int{assets.Ether(999_999_998), assets.Ether(2)})

	// Create sub, fund it and assign consumer
	_, err = coordinator.CreateSubscription(owner)
	require.NoError(t, err)
	backend.Commit()
	b, err := utils.GenericEncode([]string{"uint64"}, uint64(1))
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
	owner := newIdentity(t)
	random := newIdentity(t)
	genesisData := core.GenesisAlloc{
		owner.From: {Balance: assets.Ether(10)},
	}
	backend := cltest.NewSimulatedBackend(t, genesisData, ethconfig.Defaults.Miner.GasCeil)
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
	_, err = linkContract.Transfer(owner, consumerAddress, assets.Ether(2))
	require.NoError(t, err)
	backend.Commit()
	AssertLinkBalances(t, linkContract, []common.Address{owner.From, consumerAddress}, []*big.Int{assets.Ether(999_999_998), assets.Ether(2)})
	_, err = consumer.TopUpSubscription(owner, assets.Ether(1))
	require.NoError(t, err)
	backend.Commit()
	AssertLinkBalances(t, linkContract, []common.Address{owner.From, consumerAddress, coordinatorAddress}, []*big.Int{assets.Ether(999_999_998), assets.Ether(1), assets.Ether(1)})
	// Non-owner cannot withdraw
	_, err = consumer.Withdraw(random, assets.Ether(1), owner.From)
	require.Error(t, err)
	_, err = consumer.Withdraw(owner, assets.Ether(1), owner.From)
	require.NoError(t, err)
	backend.Commit()
	AssertLinkBalances(t, linkContract, []common.Address{owner.From, consumerAddress, coordinatorAddress}, []*big.Int{assets.Ether(999_999_999), assets.Ether(0), assets.Ether(1)})
	_, err = consumer.Unsubscribe(owner, owner.From)
	require.NoError(t, err)
	backend.Commit()
	AssertLinkBalances(t, linkContract, []common.Address{owner.From, consumerAddress, coordinatorAddress}, []*big.Int{assets.Ether(1_000_000_000), assets.Ether(0), assets.Ether(0)})
}

func TestIntegrationVRFV2(t *testing.T) {
	config, _ := heavyweight.FullTestDB(t, "vrf_v2_integration")
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, key, 1)
	carol := uni.vrfConsumers[0]
	carolContract := uni.consumerContracts[0]
	carolContractAddress := uni.consumerContractAddresses[0]

	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, uni.backend, key)
	config.Overrides.GlobalEvmGasLimitDefault = null.NewInt(0, false)
	config.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(2)
	keys, err := app.KeyStore.Eth().SendingKeys(nil)

	// Reconfigure the sim chain with a default gas price of 1 gwei,
	// max gas limit of 2M and a key specific max 10 gwei price.
	// Keep the prices low so we can operate with small link balance subscriptions.
	gasPrice := decimal.NewFromBigInt(big.NewInt(1000000000), 0)
	configureSimChain(t, app, map[string]types.ChainCfg{
		keys[0].Address.String(): {
			EvmMaxGasPriceWei: utils.NewBig(big.NewInt(10000000000)),
		},
	}, gasPrice.BigInt())

	require.NoError(t, app.Start(testutils.Context(t)))
	vrfkey, err := app.GetKeyStore().VRF().Create()
	require.NoError(t, err)

	jid := uuid.NewV4()
	incomingConfs := 2
	s := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		JobID:                    jid.String(),
		Name:                     "vrf-primary",
		CoordinatorAddress:       uni.rootContractAddress.String(),
		BatchCoordinatorAddress:  uni.batchCoordinatorContractAddress.String(),
		MinIncomingConfirmations: incomingConfs,
		PublicKey:                vrfkey.PublicKey.String(),
		FromAddresses:            []string{keys[0].Address.String()},
		V2:                       true,
	}).Toml()
	jb, err := vrf.ValidatedVRFSpec(s)
	require.NoError(t, err)
	err = app.JobSpawner().CreateJob(&jb)
	require.NoError(t, err)

	registerProvingKeyHelper(t, uni, vrfkey)

	// Create and fund a subscription.
	// We should see that our subscription has 1 link.
	AssertLinkBalances(t, uni.linkContract, []common.Address{
		carolContractAddress,
		uni.rootContractAddress,
	}, []*big.Int{
		assets.Ether(500), // 500 link
		big.NewInt(0),     // 0 link
	})
	subFunding := decimal.RequireFromString("1000000000000000000")
	_, err = carolContract.TestCreateSubscriptionAndFund(carol,
		subFunding.BigInt())
	require.NoError(t, err)
	uni.backend.Commit()
	AssertLinkBalances(t, uni.linkContract, []common.Address{
		carolContractAddress,
		uni.rootContractAddress,
		uni.nallory.From, // Oracle's own address should have nothing
	}, []*big.Int{
		assets.Ether(499),
		assets.Ether(1),
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
	_, err = carolContract.TestRequestRandomness(carol, vrfkey.PublicKey.MustHash(), subId, uint16(requestedIncomingConfs), uint32(gasRequested), uint32(nw))
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
	}, cltest.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

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
	}, cltest.WaitTimeout(t), 500*time.Millisecond).Should(gomega.BeTrue())
	assert.True(t, rf[0].Success, "expected callback to succeed")
	fulfillReceipt, err := uni.backend.TransactionReceipt(context.Background(), rf[0].Raw.TxHash)
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
	gaDecoding := big.NewInt(0).Add(ga, big.NewInt(3679))
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
	t.Logf("subscription charged %s with gas prices of %s gwei and %s ETH per LINK\n", linkCharged, gasPrice.Div(gwei), weiPerUnitLink.Div(wei))
	expected := decimal.RequireFromString(strconv.Itoa(int(fulfillReceipt.GasUsed))).Mul(gasPrice).Div(weiPerUnitLink)
	t.Logf("expected sub charge gas use %v %v off by %v", fulfillReceipt.GasUsed, expected, expected.Sub(linkCharged))
	// The expected sub charge should be within 200 gas of the actual gas usage.
	// wei/link * link / wei/gas = wei / (wei/gas) = gas
	gasDiff := linkCharged.Sub(expected).Mul(weiPerUnitLink).Div(gasPrice).Abs().IntPart()
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
		assets.Ether(499),
		subFunding.Sub(linkWeiCharged).BigInt(),
		linkWeiCharged.BigInt(),
	})

	// We should see the response count present
	chain, err := app.Chains.EVM.Get(big.NewInt(1337))
	require.NoError(t, err)

	q := pg.NewQ(app.GetSqlxDB(), app.Logger, app.Config)
	counts := vrf.GetStartingResponseCountsV2(q, app.Logger, chain.Client().ChainID().Uint64(), chain.Config().EvmFinalityDepth())
	t.Log(counts, rf[0].RequestId.String())
	assert.Equal(t, uint64(1), counts[rf[0].RequestId.String()])
}

func TestMaliciousConsumer(t *testing.T) {
	config, _ := heavyweight.FullTestDB(t, "vrf_v2_integration_malicious")
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, key, 1)
	carol := uni.vrfConsumers[0]
	config.Overrides.GlobalEvmGasLimitDefault = null.IntFrom(2000000)
	config.Overrides.GlobalEvmMaxGasPriceWei = assets.GWei(1)
	config.Overrides.GlobalEvmGasPriceDefault = assets.GWei(1)
	config.Overrides.GlobalEvmGasFeeCapDefault = assets.GWei(1)

	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, uni.backend, key)
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
		CoordinatorAddress:       uni.rootContractAddress.String(),
		BatchCoordinatorAddress:  uni.batchCoordinatorContractAddress.String(),
		MinIncomingConfirmations: incomingConfs,
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
	_, err = uni.maliciousConsumerContract.TestCreateSubscriptionAndFund(carol,
		subFunding.BigInt())
	require.NoError(t, err)
	uni.backend.Commit()

	// Send a re-entrant request
	_, err = uni.maliciousConsumerContract.TestRequestRandomness(carol)
	require.NoError(t, err)

	// We expect the request to be serviced
	// by the node.
	var attempts []txmgr.EthTxAttempt
	gomega.NewWithT(t).Eventually(func() bool {
		//runs, err = app.PipelineORM().GetAllRuns()
		attempts, _, err = app.TxmORM().EthTxAttempts(0, 1000)
		require.NoError(t, err)
		// It possible that we send the test request
		// before the job spawner has started the vrf services, which is fine
		// the lb will backfill the logs. However we need to
		// keep blocks coming in for the lb to send the backfilled logs.
		t.Log("attempts", attempts)
		uni.backend.Commit()
		return len(attempts) == 1 && attempts[0].EthTx.State == txmgr.EthTxConfirmed
	}, cltest.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

	// The fulfillment tx should succeed
	ch, err := app.GetChains().EVM.Default()
	require.NoError(t, err)
	r, err := ch.Client().TransactionReceipt(context.Background(), attempts[0].Hash)
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
	carol := uni.vrfConsumers[0]
	carolContract := uni.consumerContracts[0]
	carolContractAddress := uni.consumerContractAddresses[0]

	cfg := cltest.NewTestGeneralConfig(t)
	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, cfg, uni.backend, key)
	require.NoError(t, app.Start(testutils.Context(t)))

	vrfkey, err := app.GetKeyStore().VRF().Create()
	require.NoError(t, err)
	p, err := vrfkey.PublicKey.Point()
	require.NoError(t, err)
	_, err = uni.rootContract.RegisterProvingKey(
		uni.neil, uni.neil.From, pair(secp256k1.Coordinates(p)))
	require.NoError(t, err)
	uni.backend.Commit()
	_, err = carolContract.TestCreateSubscriptionAndFund(carol,
		big.NewInt(1000000000000000000)) // 0.1 LINK
	require.NoError(t, err)
	uni.backend.Commit()
	subId, err := carolContract.SSubId(nil)
	require.NoError(t, err)
	// Ensure even with large number of consumers its still cheap
	var addrs []common.Address
	for i := 0; i < 99; i++ {
		addrs = append(addrs, testutils.NewAddress())
	}
	_, err = carolContract.UpdateSubscription(carol,
		addrs) // 0.1 LINK
	require.NoError(t, err)
	estimate := estimateGas(t, uni.backend, common.Address{},
		carolContractAddress, uni.consumerABI,
		"testRequestRandomness", vrfkey.PublicKey.MustHash(), subId, uint16(2), uint32(10000), uint32(1))
	t.Log(estimate)
	// V2 should be at least (87000-134000)/134000 = 35% cheaper
	// Note that a second call drops further to 68998 gas, but would also drop in V1.
	assert.Less(t, estimate, uint64(90_000),
		"requestRandomness tx gas cost more than expected")
}

func TestMaxConsumersCost(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, key, 1)
	carol := uni.vrfConsumers[0]
	carolContract := uni.consumerContracts[0]
	carolContractAddress := uni.consumerContractAddresses[0]

	cfg := cltest.NewTestGeneralConfig(t)
	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, cfg, uni.backend, key)
	require.NoError(t, app.Start(testutils.Context(t)))
	_, err := carolContract.TestCreateSubscriptionAndFund(carol,
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
	assert.Less(t, estimate, uint64(265000))
	estimate = estimateGas(t, uni.backend, carolContractAddress,
		uni.rootContractAddress, uni.coordinatorABI,
		"addConsumer", subId, testutils.NewAddress())
	t.Log(estimate)
	assert.Less(t, estimate, uint64(100000))
}

func TestFulfillmentCost(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, key, 1)
	carol := uni.vrfConsumers[0]
	carolContract := uni.consumerContracts[0]
	carolContractAddress := uni.consumerContractAddresses[0]

	cfg := cltest.NewTestGeneralConfig(t)
	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, cfg, uni.backend, key)
	require.NoError(t, app.Start(testutils.Context(t)))

	vrfkey, err := app.GetKeyStore().VRF().Create()
	require.NoError(t, err)
	p, err := vrfkey.PublicKey.Point()
	require.NoError(t, err)
	_, err = uni.rootContract.RegisterProvingKey(
		uni.neil, uni.neil.From, pair(secp256k1.Coordinates(p)))
	require.NoError(t, err)
	uni.backend.Commit()
	_, err = carolContract.TestCreateSubscriptionAndFund(carol,
		big.NewInt(1000000000000000000)) // 0.1 LINK
	require.NoError(t, err)
	uni.backend.Commit()
	subId, err := carolContract.SSubId(nil)
	require.NoError(t, err)

	gasRequested := 50000
	nw := 1
	requestedIncomingConfs := 3
	_, err = carolContract.TestRequestRandomness(carol, vrfkey.PublicKey.MustHash(), subId, uint16(requestedIncomingConfs), uint32(gasRequested), uint32(nw))
	require.NoError(t, err)
	for i := 0; i < requestedIncomingConfs; i++ {
		uni.backend.Commit()
	}

	requestLog := FindLatestRandomnessRequestedLog(t, uni.rootContract, vrfkey.PublicKey.MustHash())
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
	require.NoError(t, err)
	estimate := estimateGas(t, uni.backend, common.Address{},
		uni.rootContractAddress, uni.coordinatorABI,
		"fulfillRandomWords", proof, rc)
	t.Log("estimate", estimate)
	// Establish very rough bounds on fulfillment cost
	assert.Greater(t, estimate, uint64(120000))
	assert.Less(t, estimate, uint64(500000))
}

func TestStartingCountsV1(t *testing.T) {
	cfg, db := heavyweight.FullTestDBNoFixtures(t, "vrf_test_starting_counts")
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
	err = ks.Unlock("p4SsW0rD1!@#_")
	require.NoError(t, err)
	k, err := ks.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)
	b := time.Now()
	n1, n2, n3, n4 := int64(0), int64(1), int64(2), int64(3)
	m1 := txmgr.EthTxMeta{
		RequestID: utils.PadByteToHash(0x10),
	}
	md1, err := json.Marshal(&m1)
	require.NoError(t, err)
	md1_ := datatypes.JSON(md1)
	m2 := txmgr.EthTxMeta{
		RequestID: utils.PadByteToHash(0x11),
	}
	md2, err := json.Marshal(&m2)
	md2_ := datatypes.JSON(md2)
	require.NoError(t, err)
	chainID := utils.NewBig(big.NewInt(1337))
	confirmedTxes := []txmgr.EthTx{
		{
			Nonce:              &n1,
			FromAddress:        k.Address.Address(),
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
			FromAddress:        k.Address.Address(),
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
			FromAddress:        k.Address.Address(),
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
			FromAddress:        k.Address.Address(),
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
	unconfirmedTxes := []txmgr.EthTx{}
	for i := int64(4); i < 6; i++ {
		md, err := json.Marshal(&txmgr.EthTxMeta{
			RequestID: utils.PadByteToHash(0x12),
		})
		require.NoError(t, err)
		md1 := datatypes.JSON(md)
		newNonce := i + 1
		unconfirmedTxes = append(unconfirmedTxes, txmgr.EthTx{
			Nonce:              &newNonce,
			FromAddress:        k.Address.Address(),
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
		_, err = db.NamedExec(sql, &tx)
		require.NoError(t, err)
	}

	// add eth_tx_attempts for confirmed
	broadcastBlock := int64(1)
	txAttempts := []txmgr.EthTxAttempt{}
	for i := range confirmedTxes {
		txAttempts = append(txAttempts, txmgr.EthTxAttempt{
			EthTxID:                 int64(i + 1),
			GasPrice:                utils.NewBig(big.NewInt(100)),
			SignedRawTx:             []byte(`blah`),
			Hash:                    utils.NewHash(),
			BroadcastBeforeBlockNum: &broadcastBlock,
			State:                   txmgr.EthTxAttemptBroadcast,
			CreatedAt:               time.Now(),
			ChainSpecificGasLimit:   uint64(100),
		})
	}
	// add eth_tx_attempts for unconfirmed
	for i := range unconfirmedTxes {
		txAttempts = append(txAttempts, txmgr.EthTxAttempt{
			EthTxID:               int64(i + 1 + len(confirmedTxes)),
			GasPrice:              utils.NewBig(big.NewInt(100)),
			SignedRawTx:           []byte(`blah`),
			Hash:                  utils.NewHash(),
			State:                 txmgr.EthTxAttemptInProgress,
			CreatedAt:             time.Now(),
			ChainSpecificGasLimit: uint64(100),
		})
	}
	for _, txAttempt := range txAttempts {
		t.Log("tx attempt eth tx id: ", txAttempt.EthTxID)
	}
	sql = `INSERT INTO eth_tx_attempts (eth_tx_id, gas_price, signed_raw_tx, hash, state, created_at, chain_specific_gas_limit)
		VALUES (:eth_tx_id, :gas_price, :signed_raw_tx, :hash, :state, :created_at, :chain_specific_gas_limit)`
	for _, attempt := range txAttempts {
		_, err = db.NamedExec(sql, &attempt)
		require.NoError(t, err)
	}

	// add eth_receipts
	receipts := []txmgr.EthReceipt{}
	for i := 0; i < 4; i++ {
		receipts = append(receipts, txmgr.EthReceipt{
			BlockHash:        utils.NewHash(),
			TxHash:           txAttempts[i].Hash,
			BlockNumber:      broadcastBlock,
			TransactionIndex: 1,
			Receipt:          []byte(`{}`),
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
	}, cltest.WaitTimeout(t), 500*time.Millisecond).Should(gomega.BeTrue())
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
