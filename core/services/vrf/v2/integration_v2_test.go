//go:build integration

package v2_test

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	commonkeystore "github.com/smartcontractkit/chainlink-common/keystore"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ethkey"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/vrfkey"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/vrfkey/secp256k1"
	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	evmlogger "github.com/smartcontractkit/chainlink-evm/pkg/log"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	evmutils "github.com/smartcontractkit/chainlink-evm/pkg/utils"
	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
	mocks2 "github.com/smartcontractkit/chainlink/v2/common/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	v22 "github.com/smartcontractkit/chainlink/v2/core/services/vrf/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrftesthelpers"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

var defaultMaxGasPrice = uint64(1e12)

type coordinatorV2UniverseCommon struct {
	// Golang wrappers of solidity contracts
	consumerContracts                []vrftesthelpers.VRFConsumerContract
	consumerContractAddresses        []common.Address
	rootContract                     v22.CoordinatorV2_X
	rootContractAddress              common.Address
	linkContract                     *link_token_interface.LinkToken
	linkContractAddress              common.Address
	linkEthFeedAddress               common.Address
	bhsContract                      *blockhash_store.BlockhashStore
	bhsContractAddress               common.Address
	batchBHSContract                 *batch_blockhash_store.BatchBlockhashStore
	batchBHSContractAddress          common.Address
	maliciousConsumerContract        vrftesthelpers.VRFConsumerContract
	maliciousConsumerContractAddress common.Address
	revertingConsumerContract        vrftesthelpers.VRFConsumerContract
	revertingConsumerContractAddress common.Address
	// This is a VRFConsumerV2Upgradeable wrapper that points to the proxy address.
	consumerProxyContract        vrftesthelpers.VRFConsumerContract
	consumerProxyContractAddress common.Address
	proxyAdminAddress            common.Address

	// Abstract representation of the ethereum blockchain
	backend        types.Backend
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

func makeTestTxm(t *testing.T, txStore txmgr.TestEvmTxStore, keyStore keystore.Eth, ec *clienttest.Client) txmgrcommon.TxManager[*big.Int, *types.Head, common.Address, common.Hash, common.Hash, types.Nonce, gas.EvmFee] {
	_, _, evmConfig := txmgr.MakeTestConfigs(t)
	txmConfig := txmgr.NewEvmTxmConfig(evmConfig)
	ks := keys.NewStore(keystore.NewEthSigner(keyStore, ec.ConfiguredChainID()))
	builder := mocks.NewTxAttemptBuilder[*big.Int, *types.Head, common.Address, common.Hash, common.Hash, types.Nonce, gas.EvmFee](t)
	servicetest.SetupNoOpMock(builder)
	broadcaster := mocks2.NewBroadcaster[common.Address](t)
	servicetest.SetupNoOpMock(broadcaster)
	confirmer := mocks2.NewConfirmer[*types.Head, common.Address, common.Hash](t)
	servicetest.SetupNoOpMock(confirmer)
	tracker := mocks2.NewTracker[common.Address](t)
	servicetest.SetupNoOpMock(tracker)
	tracker.On("GetEnabledAddresses").Return(nil).Maybe()
	finalizer := mocks.NewFinalizer[common.Hash, *types.Head](t)
	servicetest.SetupNoOpMock(finalizer)
	txm := txmgr.NewEvmTxm(ec.ConfiguredChainID(), txmConfig, evmConfig.Transactions(), ks, logger.TestLogger(t), nil, nil,
		builder, txStore, broadcaster, confirmer, nil, tracker, finalizer, nil, false)

	return txm
}

// Send eth from prefunded account.
// Amount is number of ETH not wei.
func sendEth(t *testing.T, key ethkey.KeyV2, b types.Backend, to common.Address, eth int) {
	ctx := t.Context()
	nonce, err := b.Client().PendingNonceAt(ctx, key.Address)
	require.NoError(t, err)
	tx := gethtypes.NewTx(&gethtypes.DynamicFeeTx{
		ChainID:   testutils.SimulatedChainID,
		Nonce:     nonce,
		GasTipCap: big.NewInt(1000000),    // 1 mwei
		GasFeeCap: assets.GWei(1).ToInt(), // block base fee in sim
		Gas:       uint64(21_000),
		To:        &to,
		Value:     big.NewInt(0).Mul(big.NewInt(int64(eth)), big.NewInt(1e18)),
		Data:      nil,
	})
	balBefore, err := b.Client().BalanceAt(ctx, to, nil)
	require.NoError(t, err)
	signedTx, err := key.SignerFn(testutils.SimulatedChainID)(key.Address, tx)
	require.NoError(t, err)
	err = b.Client().SendTransaction(ctx, signedTx)
	require.NoError(t, err)
	b.Commit()
	balAfter, err := b.Client().BalanceAt(ctx, to, nil)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(0).Sub(balAfter, balBefore).String(), tx.Value().String())
}

func subscribeVRF(
	t *testing.T,
	author *bind.TransactOpts,
	consumerContract vrftesthelpers.VRFConsumerContract,
	coordinator v22.CoordinatorV2_X,
	backend types.Backend,
	fundingAmount *big.Int,
	nativePayment bool,
) (v22.Subscription, *big.Int) {
	var err error
	if nativePayment {
		_, err = consumerContract.CreateSubscriptionAndFundNative(author, fundingAmount)
	} else {
		_, err = consumerContract.CreateSubscriptionAndFund(author, fundingAmount)
	}
	require.NoError(t, err)
	backend.Commit()

	subID, err := consumerContract.SSubID(nil)
	require.NoError(t, err)

	sub, err := coordinator.GetSubscription(nil, subID)
	require.NoError(t, err)

	if nativePayment {
		require.Equal(t, fundingAmount.String(), sub.NativeBalance().String())
	} else {
		require.Equal(t, fundingAmount.String(), sub.Balance().String())
	}

	return sub, subID
}

func createVRFJobs(
	t *testing.T,
	fromKeys [][]ethkey.KeyV2,
	app *cltest.TestApplication,
	coordinator v22.CoordinatorV2_X,
	coordinatorAddress common.Address,
	batchCoordinatorAddress common.Address,
	uni coordinatorV2UniverseCommon,
	vrfOwnerAddress *common.Address,
	vrfVersion vrfcommon.Version,
	batchEnabled bool,
	gasLanePrices ...*assets.Wei,
) (jobs []job.Job) {
	ctx := t.Context()
	require.Len(t, gasLanePrices, len(fromKeys), "must provide one gas lane price for each set of from addresses")
	// Create separate jobs for each gas lane and register their keys
	for i, keys := range fromKeys {
		keyStrs := make([]string, 0, len(keys))
		for _, k := range keys {
			keyStrs = append(keyStrs, k.Address.String())
		}

		//nolint:gosec // we already checked the length of gasLanePrices above
		gasLanePrice := gasLanePrices[i]

		vrfkey, err := app.GetKeyStore().VRF().Create(ctx)
		require.NoError(t, err)

		jid := uuid.New()
		incomingConfs := 2

		spec := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
			JobID:                    jid.String(),
			Name:                     fmt.Sprintf("vrf-primary-%d", i),
			VRFVersion:               vrfVersion,
			CoordinatorAddress:       coordinatorAddress.Hex(),
			BatchCoordinatorAddress:  batchCoordinatorAddress.Hex(),
			BatchFulfillmentEnabled:  batchEnabled,
			MinIncomingConfirmations: incomingConfs,
			PublicKey:                vrfkey.PublicKey.String(),
			FromAddresses:            keyStrs,
			BackoffInitialDelay:      10 * time.Millisecond,
			BackoffMaxDelay:          time.Second,
			V2:                       true,
			GasLanePrice:             gasLanePrice,
			EVMChainID:               testutils.SimulatedChainID.String(),
		}).Toml()

		jb, err := vrfcommon.ValidatedVRFSpec(spec)
		require.NoError(t, err)
		t.Log(jb.VRFSpec.PublicKey.MustHash(), vrfkey.PublicKey.MustHash())
		err = app.JobSpawner().CreateJob(ctx, nil, &jb)
		require.NoError(t, err)
		registerProvingKeyHelper(t, uni, coordinator, vrfkey, new(gasLanePrice.ToInt().Uint64()))
		jobs = append(jobs, jb)
	}
	// Wait until all jobs are active and listening for logs
	require.Eventually(t, func() bool {
		jbs := app.JobSpawner().ActiveJobs()
		var count int
		for _, jb := range jbs {
			if jb.Type == job.VRF {
				count++
			}
		}
		return count == len(fromKeys)
	}, testutils.WaitTimeout(t), 100*time.Millisecond)
	// Unfortunately the lb needs heads to be able to backfill logs to new subscribers.
	// To avoid confirming
	// TODO: it could just backfill immediately upon receiving a new subscriber? (though would
	// only be useful for tests, probably a more robust way is to have the job spawner accept a signal that a
	// job is fully up and running and not add it to the active jobs list before then)

	return jobs
}

// requestRandomness requests randomness from the given vrf consumer contract
// and asserts that the request ID logged by the RandomWordsRequested event
// matches the request ID that is returned and set by the consumer contract.
// The request ID and request block number are then returned to the caller.
func requestRandomnessAndAssertRandomWordsRequestedEvent(
	t *testing.T,
	vrfConsumerHandle vrftesthelpers.VRFConsumerContract,
	consumerOwner *bind.TransactOpts,
	keyHash common.Hash,
	subID *big.Int,
	numWords uint32,
	cbGasLimit uint32,
	coordinator v22.CoordinatorV2_X,
	backend types.Backend,
	nativePayment bool,
) (requestID *big.Int, requestBlockNumber uint64) {
	minRequestConfirmations := uint16(2)
	_, err := vrfConsumerHandle.RequestRandomness(
		consumerOwner,
		keyHash,
		subID,
		minRequestConfirmations,
		cbGasLimit,
		numWords,
		nativePayment,
	)
	require.NoError(t, err)
	filterOpts := commitRequestAndFilterIndexBlock(t, backend)

	// LogPoller indexes asynchronously; retry until the target block is available.
	var iter v22.RandomWordsRequestedIterator
	require.Eventually(t, func() bool {
		var filterErr error
		iter, filterErr = coordinator.FilterRandomWordsRequested(filterOpts, nil, []*big.Int{subID}, nil)
		if filterErr != nil {
			backend.Commit()
			return false
		}
		return true
	}, testutils.WaitTimeout(t), 100*time.Millisecond, "could not filter RandomWordsRequested events")

	var events []v22.RandomWordsRequested
	for iter.Next() {
		events = append(events, iter.Event())
	}

	requestID, err = vrfConsumerHandle.SRequestID(nil)
	require.NoError(t, err)

	event := events[len(events)-1]
	eventKeyHash := event.KeyHash()
	require.Equal(t, event.RequestID(), requestID, "request ID in contract does not match request ID in log")
	require.Equal(t, keyHash.Bytes(), eventKeyHash[:], "key hash of event (%s) and of request not equal (%s)", hex.EncodeToString(eventKeyHash[:]), keyHash.String())
	require.Equal(t, cbGasLimit, event.CallbackGasLimit(), "callback gas limit of event and of request not equal")
	require.Equal(t, minRequestConfirmations, event.MinimumRequestConfirmations(), "min request confirmations of event and of request not equal")
	require.Equal(t, numWords, event.NumWords(), "num words of event and of request not equal")
	require.Equal(t, nativePayment, event.NativePayment())

	return requestID, event.Raw().BlockNumber
}

func commitRequestAndFilterIndexBlock(t *testing.T, backend types.Backend) *bind.FilterOpts {
	ctx := t.Context()
	block, err := backend.Client().BlockByHash(ctx, backend.Commit())
	require.NoError(t, err)
	end := block.NumberU64()
	// Geth's filter index reads end+1, so mine one more block and filter only
	// through the request block.
	backend.Commit()
	return &bind.FilterOpts{Start: 0, End: &end, Context: ctx}
}

// subscribeAndAssertSubscriptionCreatedEvent subscribes the given consumer contract
// to VRF and funds the subscription with the given fundingJuels amount. It returns the
// subscription ID of the resulting subscription.
func subscribeAndAssertSubscriptionCreatedEvent(
	t *testing.T,
	vrfConsumerHandle vrftesthelpers.VRFConsumerContract,
	consumerOwner *bind.TransactOpts,
	consumerContractAddress common.Address,
	fundingAmount *big.Int,
	coordinator v22.CoordinatorV2_X,
	backend types.Backend,
	nativePayment bool,
) *big.Int {
	// Create a subscription and fund with LINK.
	_, subID := subscribeVRF(t, consumerOwner, vrfConsumerHandle, coordinator, backend, fundingAmount, nativePayment)

	// Assert the subscription event in the coordinator contract.
	iter, err := coordinator.FilterSubscriptionCreated(nil, []*big.Int{subID})
	require.NoError(t, err)
	found := false
	for iter.Next() {
		if iter.Event().Owner() != consumerContractAddress {
			require.FailNowf(t, "SubscriptionCreated event contains wrong owner address", "expected: %+v, actual: %+v", consumerContractAddress, iter.Event().Owner())
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
	coordinator v22.CoordinatorV2_X,
	nativePayment bool,
) (rwfe v22.RandomWordsFulfilled) {
	require.Eventually(t, func() bool {
		filter, err := coordinator.FilterRandomWordsFulfilled(nil, []*big.Int{requestID}, nil)
		require.NoError(t, err)
		for filter.Next() {
			require.Equal(t, expectedSuccess, filter.Event().Success(), "fulfillment event success not correct, expected: %+v, actual: %+v", expectedSuccess, filter.Event().Success())
			require.Equal(t, requestID, filter.Event().RequestID())
			require.Equal(t, nativePayment, filter.Event().NativePayment())
			rwfe = filter.Event()
			return true
		}
		return false
	}, 3*time.Second, 100*time.Millisecond, "RandomWordsFulfilled event not found")
	return rwfe
}

func assertNumRandomWords(
	t *testing.T,
	contract vrftesthelpers.VRFConsumerContract,
	numWords uint32,
) {
	var err error
	for i := range numWords {
		_, err = contract.SRandomWords(nil, big.NewInt(int64(i)))
		require.NoError(t, err)
	}
}

func mine(t *testing.T, requestID, subID *big.Int, backend types.Backend, db *sqlx.DB, vrfVersion vrfcommon.Version, chainID *big.Int) bool {
	txstore := txmgr.NewTxStore(db, logger.TestLogger(t))
	var metaField string
	switch vrfVersion {
	case vrfcommon.V2Plus:
		metaField = "GlobalSubId"
	default:
		t.Errorf("unsupported vrf version %s", vrfVersion)
	}

	return assert.Eventually(t, func() bool {
		backend.Commit()
		txes, err := txstore.FindTxesByMetaFieldAndStates(t.Context(), metaField, subID.String(), []txmgrtypes.TxState{txmgrcommon.TxConfirmed, txmgrcommon.TxFinalized}, chainID)
		require.NoError(t, err)
		for _, tx := range txes {
			if !checkForReceipt(t, db, tx.ID) {
				return false
			}
			meta, err := tx.GetMeta()
			require.NoError(t, err)
			if meta.RequestID.String() == common.BytesToHash(requestID.Bytes()).String() {
				return true
			}
		}
		return false
	}, testutils.WaitTimeout(t), 100*time.Millisecond)
}

func mineBatch(t *testing.T, requestIDs []*big.Int, subID *big.Int, backend types.Backend, db *sqlx.DB, vrfVersion vrfcommon.Version, chainID *big.Int) bool {
	requestIDMap := map[string]bool{}
	txstore := txmgr.NewTxStore(db, logger.TestLogger(t))
	var metaField string
	switch vrfVersion {
	case vrfcommon.V2Plus:
		metaField = "GlobalSubId"
	default:
		t.Errorf("unsupported vrf version %s", vrfVersion)
	}
	for _, requestID := range requestIDs {
		requestIDMap[common.BytesToHash(requestID.Bytes()).String()] = false
	}
	return assert.Eventually(t, func() bool {
		backend.Commit()
		txes, err := txstore.FindTxesByMetaFieldAndStates(t.Context(), metaField, subID.String(), []txmgrtypes.TxState{txmgrcommon.TxConfirmed, txmgrcommon.TxFinalized}, chainID)
		require.NoError(t, err)
		for _, tx := range txes {
			if !checkForReceipt(t, db, tx.ID) {
				return false
			}
			meta, err := tx.GetMeta()
			require.NoError(t, err)
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
	}, testutils.WaitTimeout(t), 100*time.Millisecond)
}

func checkForReceipt(t *testing.T, db *sqlx.DB, txID int64) bool {
	// Confirm receipt is fetched and stored for transaction to consider it mined
	var count uint32
	sql := `
	SELECT count(*) FROM evm.receipts
	JOIN evm.tx_attempts ON evm.tx_attempts.hash = evm.receipts.tx_hash
	JOIN evm.txes ON evm.txes.ID = evm.tx_attempts.eth_tx_id
	WHERE evm.txes.ID = $1 AND evm.txes.state IN ('confirmed', 'finalized')`
	if txID != -1 {
		err := db.GetContext(t.Context(), &count, sql, txID)
		require.NoError(t, err)
	} else {
		sql = strings.Replace(sql, "evm.txes.ID = $1", "evm.txes.meta->>'ForceFulfilled' IS NOT NULL", 1)
		err := db.GetContext(t.Context(), &count, sql, txID)
		require.NoError(t, err)
	}
	return count > 0
}

func testEoa(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2UniverseCommon,
	batchingEnabled bool,
	batchCoordinatorAddress common.Address,
	vrfOwnerAddress *common.Address,
	vrfVersion vrfcommon.Version,
) {
	ctx := t.Context()
	gasLimit := uint64(2_500_000)

	finalityDepth := uint32(50)

	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), toml.KeySpecific{
			// Gas lane.
			Key:          new(key1.EIP55Address),
			GasEstimator: toml.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = new(gasLimit)
		c.EVM[0].MinIncomingConfirmations = new(uint32(2))
		c.EVM[0].FinalityDepth = new(finalityDepth)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)
	consumer := uni.vrfConsumers[0]

	// Create a new subscription.
	subID := setupAndFundSubscriptionAndConsumer(
		t,
		uni,
		uni.rootContract,
		uni.rootContractAddress,
		consumer,
		consumer.From,
		vrfVersion,
		assets.Ether(1).ToInt(),
	)

	// Fund gas lane.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	require.NoError(t, app.Start(ctx))

	// Create VRF job.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		batchCoordinatorAddress,
		uni,
		vrfOwnerAddress,
		vrfVersion,
		batchingEnabled,
		gasLanePriceWei,
	)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make a randomness request with the EOA. This request is impossible to fulfill.
	numWords := uint32(1)
	minRequestConfirmations := uint16(2)
	{
		_, err := uni.rootContract.RequestRandomWords(consumer, keyHash, subID, minRequestConfirmations, uint32(200_000), numWords, false)
		require.NoError(t, err)
	}
	uni.backend.Commit()

	// Ensure request is not fulfilled.
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 0
	}, 5*time.Second, time.Second).Should(gomega.BeTrue())

	// Create query to fetch the application's log broadcasts.
	var broadcastsBeforeFinality []evmlogger.LogBroadcast
	var broadcastsAfterFinality []evmlogger.LogBroadcast
	query := `SELECT block_hash, consumed, log_index, job_id FROM log_broadcasts`

	// Execute the query.
	require.NoError(t, app.GetDB().SelectContext(ctx, &broadcastsBeforeFinality, query))

	// Ensure there is only one log broadcast (our EOA request), and that
	// it hasn't been marked as consumed yet.
	require.Len(t, broadcastsBeforeFinality, 1)
	require.False(t, broadcastsBeforeFinality[0].Consumed)

	// Create new blocks until the finality depth has elapsed.
	for range finalityDepth {
		uni.backend.Commit()
	}

	// Ensure the request is still not fulfilled.
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 0
	}, 5*time.Second, time.Second).Should(gomega.BeTrue())

	// Execute the query for log broadcasts again after finality depth has elapsed.
	require.NoError(t, app.GetDB().SelectContext(ctx, &broadcastsAfterFinality, query))

	// Ensure that there is still only one log broadcast (our EOA request), but that
	// it has been marked as "consumed," such that it won't be retried.
	require.Len(t, broadcastsAfterFinality, 1)
	require.True(t, broadcastsAfterFinality[0].Consumed)

	t.Log("Done!")
}

func TestVRFV2Integration_SingleConsumer_NeedsTrustedBlockhashStore(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 2, true)
	testMultipleConsumersNeedTrustedBHS(
		t,
		ownerKey,
		uni,
		uni.vrfConsumers,
		uni.consumerContracts,
		uni.consumerContractAddresses,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		vrfcommon.V2Plus,
		false,
		false,
	)
}

func TestVRFV2Integration_SingleConsumer_NeedsTrustedBlockhashStore_AfterDelay(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2PlusUniverse(t, ownerKey, 2, true)
	testMultipleConsumersNeedTrustedBHS(
		t,
		ownerKey,
		uni,
		uni.vrfConsumers,
		uni.consumerContracts,
		uni.consumerContractAddresses,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		vrfcommon.V2Plus,
		false,
		true,
	)
}

func simulatedOverrides(t *testing.T, defaultGasPrice *assets.Wei, ks ...toml.KeySpecific) func(*chainlink.Config, *chainlink.Secrets) {
	return func(c *chainlink.Config, s *chainlink.Secrets) {
		require.Zero(t, testutils.SimulatedChainID.Cmp(c.EVM[0].ChainID.ToInt()))
		c.EVM[0].GasEstimator.Mode = new("FixedPrice")
		if defaultGasPrice != nil {
			c.EVM[0].GasEstimator.PriceDefault = defaultGasPrice
		}
		c.EVM[0].GasEstimator.LimitDefault = new(uint64(3_500_000))

		c.Feature.LogPoller = new(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)

		c.EVM[0].HeadTracker.MaxBufferSize = new(uint32(100))
		c.EVM[0].HeadTracker.SamplingInterval = commonconfig.MustNewDuration(0) // Head sampling disabled

		c.EVM[0].Transactions.ResendAfterThreshold = commonconfig.MustNewDuration(0)
		c.EVM[0].Transactions.ReaperThreshold = commonconfig.MustNewDuration(100 * time.Millisecond)

		c.EVM[0].FinalityDepth = new(uint32(15))
		c.EVM[0].MinIncomingConfirmations = new(uint32(1))
		c.EVM[0].MinContractPayment = commonassets.NewLinkFromJuels(100)
		c.EVM[0].KeySpecific = ks
	}
}

func registerProvingKeyHelper(t *testing.T, uni coordinatorV2UniverseCommon, coordinator v22.CoordinatorV2_X, vrfkey vrfkey.KeyV2, gasLaneMaxGas *uint64) {
	// Register a proving key associated with the VRF job.
	p, err := vrfkey.PublicKey.Point()
	require.NoError(t, err)
	if gasLaneMaxGas == nil {
		t.Error("gasLaneMaxGas must be non-nil for V2+")
	}
	_, err = coordinator.RegisterProvingKey(
		uni.neil, nil, pair(secp256k1.Coordinates(p)), gasLaneMaxGas,
	)
	require.NoError(t, err)
	uni.backend.Commit()
}

func TestStartingCountsV1(t *testing.T) {
	t.Parallel()
	cfg, db := heavyweight.FullTestDBNoFixturesV2(t, nil)

	ctx := t.Context()
	txStore := txmgr.NewTxStore(db, logger.TestLogger(t))
	lggr := logger.TestLogger(t)
	ks := keystore.NewInMemory(db, commonkeystore.FastScryptParams, lggr.Infof)
	ec := clienttest.NewClient(t)
	ec.On("Dial", mock.Anything).Maybe().Return(nil)
	ec.On("Close").Maybe().Return(nil)
	ec.On("ConfiguredChainID").Return(testutils.SimulatedChainID)
	ec.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(big.NewInt(0), nil)
	ec.On("LatestBlockHeight", mock.Anything).Return(big.NewInt(2), nil).Maybe()
	ec.On("SubscribeToHeads", mock.Anything).Maybe().Return(nil, nil, errors.ErrUnsupported)
	txm := makeTestTxm(t, txStore, ks.Eth(), ec)
	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		KeyStore:       ks.Eth(),
		Client:         ec,
		DB:             db,
		ChainConfigs:   cfg.EVMConfigs(),
		DatabaseConfig: cfg.Database(),
		FeatureConfig:  cfg.Feature(),
		ListenerConfig: cfg.Database().Listener(),
		TxManager:      txm,
	})
	chainService, err := legacyChains.Get(testutils.SimulatedChainID.String())
	require.NoError(t, err)
	chain, ok := chainService.(legacyevm.Chain)
	require.True(t, ok)
	listenerV2 := v22.MakeTestListenerV2(chain)
	var counts map[[32]byte]uint64
	counts, err = vrfcommon.GetStartingResponseCountsV1(t.Context(), chain)
	require.NoError(t, err)
	assert.Empty(t, counts)
	err = ks.Unlock(ctx, testutils.Password)
	require.NoError(t, err)
	k, err := ks.Eth().Create(t.Context(), testutils.SimulatedChainID)
	require.NoError(t, err)
	b := time.Now()
	n1, n2, n3, n4 := types.Nonce(0), types.Nonce(1), types.Nonce(2), types.Nonce(3)
	reqID := evmutils.PadByteToHash(0x10)
	m1 := txmgr.TxMeta{
		RequestID: &reqID,
	}
	md1, err := json.Marshal(&m1)
	require.NoError(t, err)
	md1SQL := sqlutil.JSON(md1)
	reqID2 := evmutils.PadByteToHash(0x11)
	m2 := txmgr.TxMeta{
		RequestID: &reqID2,
	}
	md2, err := json.Marshal(&m2)
	md2SQL := sqlutil.JSON(md2)
	require.NoError(t, err)
	chainID := sqlutil.New(testutils.SimulatedChainID)
	// Build unconfirmed txes first so confirmedTxes can be preallocated for append(confirmed, unconfirmed...).
	unconfirmedTxes := make([]txmgr.Tx, 0, 2)
	for i := int64(4); i < 6; i++ {
		reqID3 := evmutils.PadByteToHash(0x12)
		md, err2 := json.Marshal(&txmgr.TxMeta{
			RequestID: &reqID3,
		})
		require.NoError(t, err2)
		mdSQL := sqlutil.JSON(md)
		newNonce := types.Nonce(i + 1)
		unconfirmedTxes = append(unconfirmedTxes, txmgr.Tx{
			Sequence:           &newNonce,
			FromAddress:        k.Address,
			Error:              null.String{},
			CreatedAt:          b,
			State:              txmgrcommon.TxUnconfirmed,
			BroadcastAt:        &b,
			InitialBroadcastAt: &b,
			Meta:               &mdSQL,
			EncodedPayload:     []byte{},
			ChainID:            chainID.ToInt(),
		})
	}
	confirmedTxes := make([]txmgr.Tx, 0, 4+len(unconfirmedTxes))
	confirmedTxes = append(confirmedTxes,
		txmgr.Tx{
			Sequence:           &n1,
			FromAddress:        k.Address,
			Error:              null.String{},
			BroadcastAt:        &b,
			InitialBroadcastAt: &b,
			CreatedAt:          b,
			State:              txmgrcommon.TxConfirmed,
			Meta:               &sqlutil.JSON{},
			EncodedPayload:     []byte{},
			ChainID:            chainID.ToInt(),
		},
		txmgr.Tx{
			Sequence:           &n2,
			FromAddress:        k.Address,
			Error:              null.String{},
			BroadcastAt:        &b,
			InitialBroadcastAt: &b,
			CreatedAt:          b,
			State:              txmgrcommon.TxConfirmed,
			Meta:               &md1SQL,
			EncodedPayload:     []byte{},
			ChainID:            chainID.ToInt(),
		},
		txmgr.Tx{
			Sequence:           &n3,
			FromAddress:        k.Address,
			Error:              null.String{},
			BroadcastAt:        &b,
			InitialBroadcastAt: &b,
			CreatedAt:          b,
			State:              txmgrcommon.TxConfirmed,
			Meta:               &md2SQL,
			EncodedPayload:     []byte{},
			ChainID:            chainID.ToInt(),
		},
		txmgr.Tx{
			Sequence:           &n4,
			FromAddress:        k.Address,
			Error:              null.String{},
			BroadcastAt:        &b,
			InitialBroadcastAt: &b,
			CreatedAt:          b,
			State:              txmgrcommon.TxConfirmed,
			Meta:               &md2SQL,
			EncodedPayload:     []byte{},
			ChainID:            chainID.ToInt(),
		},
	)
	numConfirmed := len(confirmedTxes)
	confirmedTxes = append(confirmedTxes, unconfirmedTxes...)
	for i := range confirmedTxes {
		err = txStore.InsertTx(ctx, &confirmedTxes[i])
		require.NoError(t, err)
	}

	// add tx attempt for confirmed
	broadcastBlock := int64(1)
	txAttempts := make([]txmgr.TxAttempt, 0, len(confirmedTxes))
	for i := range numConfirmed {
		txAttempts = append(txAttempts, txmgr.TxAttempt{
			TxID:                    int64(i + 1),
			TxFee:                   gas.EvmFee{GasPrice: assets.NewWeiI(100)},
			SignedRawTx:             []byte(`blah`),
			Hash:                    evmutils.NewHash(),
			BroadcastBeforeBlockNum: &broadcastBlock,
			State:                   txmgrtypes.TxAttemptBroadcast,
			CreatedAt:               time.Now(),
			ChainSpecificFeeLimit:   uint64(100),
		})
	}
	// add tx attempt for unconfirmed
	for i := range unconfirmedTxes {
		txAttempts = append(txAttempts, txmgr.TxAttempt{
			TxID:                  int64(i + 1 + numConfirmed),
			TxFee:                 gas.EvmFee{GasPrice: assets.NewWeiI(100)},
			SignedRawTx:           []byte(`blah`),
			Hash:                  evmutils.NewHash(),
			State:                 txmgrtypes.TxAttemptInProgress,
			CreatedAt:             time.Now(),
			ChainSpecificFeeLimit: uint64(100),
		})
	}
	for _, txAttempt := range txAttempts {
		t.Log("tx attempt eth tx id: ", txAttempt.TxID)
	}
	for i := range txAttempts {
		err = txStore.InsertTxAttempt(ctx, &txAttempts[i])
		require.NoError(t, err)
	}

	// add evm.receipts
	receipts := make([]types.Receipt, 0, 4)
	for i := range 4 {
		receipts = append(receipts, types.Receipt{
			BlockHash:        evmutils.NewHash(),
			TxHash:           txAttempts[i].Hash,
			BlockNumber:      big.NewInt(broadcastBlock),
			TransactionIndex: 1,
		})
	}
	for i := range receipts {
		_, err = txStore.InsertReceipt(ctx, &receipts[i])
		require.NoError(t, err)
	}

	counts, err = vrfcommon.GetStartingResponseCountsV1(t.Context(), chain)
	require.NoError(t, err)
	assert.Len(t, counts, 3)
	assert.Equal(t, uint64(1), counts[evmutils.PadByteToHash(0x10)])
	assert.Equal(t, uint64(2), counts[evmutils.PadByteToHash(0x11)])
	assert.Equal(t, uint64(2), counts[evmutils.PadByteToHash(0x12)])

	countsV2, err := listenerV2.GetStartingResponseCountsV2(t.Context())
	require.NoError(t, err)
	t.Log(countsV2)
	assert.Len(t, countsV2, 3)
	assert.Equal(t, uint64(1), countsV2[big.NewInt(0x10).String()])
	assert.Equal(t, uint64(2), countsV2[big.NewInt(0x11).String()])
	assert.Equal(t, uint64(2), countsV2[big.NewInt(0x12).String()])
}

func FindLatestRandomnessRequestedLog(t *testing.T,
	coordContract v22.CoordinatorV2_X,
	keyHash [32]byte,
	requestID *big.Int,
) v22.RandomWordsRequested {
	var rf []v22.RandomWordsRequested
	require.Eventually(t, func() bool {
		rfIterator, err2 := coordContract.FilterRandomWordsRequested(nil, [][32]byte{keyHash}, nil, []common.Address{})
		require.NoError(t, err2, "failed to logs")
		for rfIterator.Next() {
			if requestID == nil || requestID.Cmp(rfIterator.Event().RequestID()) == 0 {
				rf = append(rf, rfIterator.Event())
			}
		}
		return len(rf) >= 1
	}, testutils.WaitTimeout(t), 500*time.Millisecond)
	latest := len(rf) - 1
	return rf[latest]
}

func AssertLinkBalance(t *testing.T, linkContract *link_token_interface.LinkToken, address common.Address, balance *big.Int) {
	b, err := linkContract.BalanceOf(nil, address)
	require.NoError(t, err)
	assert.Equal(t, balance.String(), b.String(), "invalid balance for %v", address)
}

func AssertNativeBalance(t *testing.T, backend types.Backend, address common.Address, balance *big.Int) {
	b, err := backend.Client().BalanceAt(t.Context(), address, nil)
	require.NoError(t, err)
	assert.Equal(t, balance.String(), b.String(), "invalid balance for %v", address)
}

func AssertLinkBalances(t *testing.T, linkContract *link_token_interface.LinkToken, addresses []common.Address, balances []*big.Int) {
	require.Len(t, balances, len(addresses))
	for i, a := range addresses {
		AssertLinkBalance(t, linkContract, a, balances[i])
	}
}

func pair(x, y *big.Int) [2]*big.Int { return [2]*big.Int{x, y} }

// estimateGas returns the estimated gas cost of running the given method on the
// contract at address to, on the given backend, with the given args, and given
// that the transaction is sent from the from address.
func estimateGas(t *testing.T, backend types.Backend,
	from, to common.Address, abi *abi.ABI, method string, args ...any,
) uint64 {
	rawData, err := abi.Pack(method, args...)
	require.NoError(t, err, "failed to construct raw %s transaction with args %s",
		method, args)
	callMsg := ethereum.CallMsg{From: from, To: &to, Data: rawData}
	estimate, err := backend.Client().EstimateGas(t.Context(), callMsg)
	require.NoError(t, err, "failed to estimate gas from %s call with args %s",
		method, args)
	return estimate
}
