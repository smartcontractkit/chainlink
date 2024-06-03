package v2_test

import (
	"database/sql"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/jmoiron/sqlx"

	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	evmutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_external_sub_owner_example"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	prooflib "github.com/smartcontractkit/chainlink/v2/core/services/vrf/proof"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/vrf/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	coordinatorV2ABI      = evmtypes.MustGetABI(vrf_coordinator_v2.VRFCoordinatorV2ABI)
	batchCoordinatorV2ABI = evmtypes.MustGetABI(batch_vrf_coordinator_v2.BatchVRFCoordinatorV2ABI)
)

func TestVRFV2Integration_SingleRevertedTxn_ForceFulfillment(t *testing.T) {
	t.Parallel()

	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	th := newRevertTxnTH(t, &uni, ownerKey, false, []uint64{1})

	// Make VRF request without sufficient balance and send fulfillment without simulation
	req := makeVRFReq(t, th, th.subs[0])
	req = fulfillVRFReq(t, th, req, th.subs[0], false, nil)

	waitForForceFulfillment(t, th, req, th.subs[0], true, 1)

	t.Log("Done!")
}

func TestVRFV2Integration_BatchRevertedTxn_ForceFulfillment(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)

	th := newRevertTxnTH(t, &uni, ownerKey, true, []uint64{1})

	numReqs := 2
	reqs := make([]*vrfReq, numReqs)
	for i := 0; i < numReqs; i++ {
		reqs[i] = makeVRFReq(t, th, th.subs[0])
	}
	fulfilBatchVRFReq(t, th, reqs, th.subs[0])

	for i := 0; i < numReqs; i++ {
		// The last request will be the successful one because of the way the example
		// contract is written.
		success := false
		if i == (numReqs - 1) {
			success = true
		}
		waitForForceFulfillment(t, th, reqs[i], th.subs[0], success, 1)
	}
	t.Log("Done!")
}

func TestVRFV2Integration_ForceFulfillmentRevertedTxn_Retry(t *testing.T) {
	t.Parallel()

	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	th := newRevertTxnTH(t, &uni, ownerKey, false, []uint64{1})

	// Make VRF request without sufficient balance and send fulfillment without simulation
	req := makeVRFReq(t, th, th.subs[0])
	req = fulfillVRFReq(t, th, req, th.subs[0], true, ptr(uint64(7)))

	waitForForceFulfillment(t, th, req, th.subs[0], true, 2)

	receipts, err := getTxnReceiptDB(th.db, -1)
	require.Nil(t, err)
	require.Len(t, receipts, 2)
	require.Equal(t, uint64(0), receipts[0].EVMReceipt.Status)
	require.Equal(t, uint64(1), receipts[1].EVMReceipt.Status)
	require.Equal(t, uint64(8), receipts[1].ForceFulfillmentAttempt)

	t.Log("Done!")
}
func TestVRFV2Integration_CanceledSubForceFulfillmentRevertedTxn_Retry(t *testing.T) {
	t.Parallel()

	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	th := newRevertTxnTH(t, &uni, ownerKey, false, []uint64{1})

	// Make VRF request without sufficient balance and send fulfillment without simulation
	req := makeVRFReq(t, th, th.subs[0])
	req = fulfillVRFReq(t, th, req, th.subs[0], true, nil)

	waitForForceFulfillment(t, th, req, th.subs[0], true, 2)

	receipts, err := getTxnReceiptDB(th.db, -1)
	require.Nil(t, err)
	require.Len(t, receipts, 2)
	require.Equal(t, uint64(0), receipts[0].EVMReceipt.Status)
	require.Equal(t, uint64(1), receipts[1].EVMReceipt.Status)
	require.Equal(t, uint64(1), receipts[1].ForceFulfillmentAttempt)

	t.Log("Done!")
}

func TestUniqueReqById_NoPendingReceipts(t *testing.T) {
	revertedForceTxns := []v2.TxnReceiptDB{
		{RequestID: common.BigToHash(big.NewInt(1)).Hex(),
			ForceFulfillmentAttempt: 1, EVMReceipt: evmtypes.Receipt{Status: 0}},
		{RequestID: common.BigToHash(big.NewInt(1)).Hex(),
			ForceFulfillmentAttempt: 2, EVMReceipt: evmtypes.Receipt{Status: 0}},
		{RequestID: common.BigToHash(big.NewInt(2)).Hex(),
			ForceFulfillmentAttempt: 1, EVMReceipt: evmtypes.Receipt{Status: 0}},
		{RequestID: common.BigToHash(big.NewInt(2)).Hex(),
			ForceFulfillmentAttempt: 2, EVMReceipt: evmtypes.Receipt{Status: 0}},
		{RequestID: common.BigToHash(big.NewInt(2)).Hex(),
			ForceFulfillmentAttempt: 3, EVMReceipt: evmtypes.Receipt{Status: 0}},
		{RequestID: common.BigToHash(big.NewInt(2)).Hex(),
			ForceFulfillmentAttempt: 4, EVMReceipt: evmtypes.Receipt{Status: 0}},
	}
	allForceTxns := revertedForceTxns
	res := v2.UniqueByReqID(revertedForceTxns, allForceTxns)
	require.Len(t, res, 2)
	for _, r := range res {
		if r.RequestID == "1" {
			require.Equal(t, r.ForceFulfillmentAttempt, 2)
		}
		if r.RequestID == "2" {
			require.Equal(t, r.ForceFulfillmentAttempt, 4)
		}
	}
}

func TestUniqueReqById_WithPendingReceipts(t *testing.T) {
	revertedForceTxns := []v2.TxnReceiptDB{
		{RequestID: common.BigToHash(big.NewInt(1)).Hex(),
			ForceFulfillmentAttempt: 1, EVMReceipt: evmtypes.Receipt{Status: 0}},
		{RequestID: common.BigToHash(big.NewInt(1)).Hex(),
			ForceFulfillmentAttempt: 2, EVMReceipt: evmtypes.Receipt{Status: 0}},
		{RequestID: common.BigToHash(big.NewInt(2)).Hex(),
			ForceFulfillmentAttempt: 1, EVMReceipt: evmtypes.Receipt{Status: 0}},
		{RequestID: common.BigToHash(big.NewInt(2)).Hex(),
			ForceFulfillmentAttempt: 2, EVMReceipt: evmtypes.Receipt{Status: 0}},
		{RequestID: common.BigToHash(big.NewInt(2)).Hex(),
			ForceFulfillmentAttempt: 3, EVMReceipt: evmtypes.Receipt{Status: 0}},
		{RequestID: common.BigToHash(big.NewInt(2)).Hex(),
			ForceFulfillmentAttempt: 4, EVMReceipt: evmtypes.Receipt{Status: 0}},
	}
	allForceTxns := []v2.TxnReceiptDB{}
	allForceTxns = append(allForceTxns, revertedForceTxns...)
	allForceTxns = append(allForceTxns, v2.TxnReceiptDB{RequestID: common.BigToHash(big.NewInt(2)).Hex(),
		ForceFulfillmentAttempt: 5})
	res := v2.UniqueByReqID(revertedForceTxns, allForceTxns)
	require.Len(t, res, 1)
	for _, r := range res {
		if r.RequestID == "1" {
			require.Equal(t, r.ForceFulfillmentAttempt, 2)
		}
	}
}

// Wait till force fulfillment event fired for the req passed in, till go test timeout
func waitForForceFulfillment(t *testing.T,
	th *revertTxnTH,
	req *vrfReq,
	sub *vrfSub,
	success bool,
	forceFulfilledCount int64) {
	uni := th.uni
	coordinator := th.uni.rootContract
	requestID := req.requestID

	// Wait for force-fulfillment to be queued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		commitment, err := coordinator.GetCommitment(nil, requestID)
		require.NoError(t, err)
		t.Log("commitment is:", hexutil.Encode(commitment[:]), ", requestID: ", common.BigToHash(requestID).Hex())
		checkForForceFulfilledEvent(t, th, req, sub, -1)
		return utils.IsEmpty(commitment[:])
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mineForceFulfilled(t, requestID, sub.subID, forceFulfilledCount, *uni, th.db)

	// Assert correct state of RandomWordsFulfilled event.
	// In this particular case:
	// * success should be true
	// * payment should be zero (forced fulfillment)
	rwfe := assertRandomWordsFulfilled(t, requestID, success, coordinator, false)
	require.Equal(t, "0", rwfe.Payment().String())

	// Check that the RandomWordsForced event is emitted correctly.
	checkForForceFulfilledEvent(t, th, req, sub, 0)
}

// Check if force fulfillment event fired for the req passed in
func checkForForceFulfilledEvent(t *testing.T,
	th *revertTxnTH,
	req *vrfReq,
	sub *vrfSub,
	numForcedLogs int) {
	requestID := req.requestID
	it, err := th.uni.vrfOwnerNew.FilterRandomWordsForced(nil, []*big.Int{requestID},
		[]uint64{sub.subID}, []common.Address{th.eoaConsumerAddr})
	require.NoError(t, err)
	i := 0
	for it.Next() {
		i++
		require.Equal(t, requestID.String(), it.Event.RequestId.String())
		require.Equal(t, sub.subID, it.Event.SubId)
		require.Equal(t, th.eoaConsumerAddr.String(), it.Event.Sender.String())
	}
	t.Log("Number of RandomWordsForced Logs:", i)
	require.Greater(t, i, numForcedLogs)
}

// Make VRF request without sufficient balance and send fulfillment without simulation
func makeVRFReq(t *testing.T, th *revertTxnTH, sub *vrfSub) (req *vrfReq) {
	// Make the randomness request and send fulfillment without simulation
	numWords := uint32(3)
	confs := 10
	callbackGasLimit := uint32(600_000)
	_, err := th.eoaConsumer.RequestRandomWords(th.uni.neil, sub.subID,
		callbackGasLimit, uint16(confs), numWords, th.keyHash)
	require.NoError(t, err, fmt.Sprintf("failed to request randomness from consumer: %v", err))
	th.uni.backend.Commit()

	// Generate VRF proof
	requestID, err := th.eoaConsumer.SRequestId(nil)
	require.NoError(t, err)

	return &vrfReq{requestID: requestID, callbackGasLimit: callbackGasLimit, numWords: numWords}
}

// Fulfill VRF req without prior simulation, after computing req proof and commitment
func fulfillVRFReq(t *testing.T,
	th *revertTxnTH,
	req *vrfReq,
	sub *vrfSub,
	forceFulfill bool,
	forceFulfilmentAttempt *uint64) *vrfReq {
	// Generate VRF proof and commitment
	reqUpdated := genReqProofNCommitment(t, th, *req, sub)
	req = &reqUpdated

	// Send fulfillment TX w/ out simulation to txm, to revert on-chain

	// Construct data payload
	b, err := coordinatorV2ABI.Pack("fulfillRandomWords", req.proof, req.reqCommitment)
	require.NoError(t, err)

	ec := th.uni.backend
	chainID := th.uni.backend.Blockchain().Config().ChainID
	chain, err := th.app.GetRelayers().LegacyEVMChains().Get(chainID.String())
	require.NoError(t, err)

	metadata := &txmgr.TxMeta{
		RequestID:     ptr(common.BytesToHash(req.requestID.Bytes())),
		SubID:         &sub.subID,
		RequestTxHash: req.requestTxHash,
		// No max link since simulation failed
	}
	if forceFulfill {
		metadata.ForceFulfilled = ptr(true)
		if forceFulfilmentAttempt != nil {
			metadata.ForceFulfillmentAttempt = forceFulfilmentAttempt
		}
	}
	etx, err := chain.TxManager().CreateTransaction(testutils.Context(t), txmgr.TxRequest{
		FromAddress:    th.key1.EIP55Address.Address(),
		ToAddress:      th.uni.rootContractAddress,
		EncodedPayload: b,
		FeeLimit:       1e6,
		Strategy:       txmgrcommon.NewSendEveryStrategy(),
		Meta:           metadata,
	})
	require.NoError(t, err)
	ec.Commit()

	// wait for above tx to mine (reach state confirmed)
	mine(t, req.requestID, big.NewInt(int64(sub.subID)), th.uni.backend, th.db, vrfcommon.V2, th.chainID)

	receipts, err := getTxnReceiptDB(th.db, etx.ID)
	require.Nil(t, err)
	require.Len(t, receipts, 1)
	require.Equal(t, uint64(0), receipts[0].EVMReceipt.Status)
	req.txID = etx.ID
	return req
}

// Fulfill VRF req without prior simulation, after computing req proof and commitment
func fulfilBatchVRFReq(t *testing.T,
	th *revertTxnTH,
	reqs []*vrfReq,
	sub *vrfSub) {
	proofs := make([]vrf_coordinator_v2.VRFProof, 0)
	reqCommitments := make([]vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment, 0)
	requestIDs := make([]common.Hash, 0)
	requestIDInts := make([]*big.Int, 0)
	requestTxnHashes := make([]common.Hash, 0)
	// Generate VRF proof and commitment
	for i, req := range reqs {
		reqUpdated := genReqProofNCommitment(t, th, *req, sub)
		reqs[i] = &reqUpdated
		proofs = append(proofs, *reqUpdated.proof)
		reqCommitments = append(reqCommitments, *reqUpdated.reqCommitment)
		requestIDs = append(requestIDs, common.BytesToHash(reqUpdated.requestID.Bytes()))
		requestIDInts = append(requestIDInts, reqUpdated.requestID)
		requestTxnHashes = append(requestTxnHashes, *reqUpdated.requestTxHash)
	}

	// Send fulfillment TX w/ out simulation to txm, to revert on-chain

	// Construct data payload
	b, err := batchCoordinatorV2ABI.Pack("fulfillRandomWords", proofs, reqCommitments)
	require.NoError(t, err)

	ec := th.uni.backend
	chainID := th.uni.backend.Blockchain().Config().ChainID
	chain, err := th.app.GetRelayers().LegacyEVMChains().Get(chainID.String())
	require.NoError(t, err)

	etx, err := chain.TxManager().CreateTransaction(testutils.Context(t), txmgr.TxRequest{
		FromAddress:    th.key1.EIP55Address.Address(),
		ToAddress:      th.uni.batchCoordinatorContractAddress,
		EncodedPayload: b,
		FeeLimit:       1e6,
		Strategy:       txmgrcommon.NewSendEveryStrategy(),
		Meta: &txmgr.TxMeta{
			RequestIDs:      requestIDs,
			RequestTxHashes: requestTxnHashes,
			SubID:           &sub.subID,
			// No max link since simulation failed
		},
	})
	require.NoError(t, err)
	ec.Commit()

	// wait for above tx to mine (reach state confirmed)
	mineBatch(t, requestIDInts, big.NewInt(int64(sub.subID)), th.uni.backend, th.db, vrfcommon.V2, chainID)

	receipts, err := getTxnReceiptDB(th.db, etx.ID)
	require.Nil(t, err)
	require.Len(t, receipts, 1)
	require.Equal(t, uint64(1), receipts[0].EVMReceipt.Status)
}

// Fulfill VRF req without prior simulation, after computing req proof and commitment
func genReqProofNCommitment(t *testing.T,
	th *revertTxnTH,
	req vrfReq,
	sub *vrfSub) vrfReq {
	// Generate VRF proof
	requestLog := FindLatestRandomnessRequestedLog(t, th.uni.rootContract, th.keyHash, req.requestID)
	s, err := prooflib.BigToSeed(requestLog.PreSeed())
	require.NoError(t, err)
	proof, rc, err := prooflib.GenerateProofResponseV2(th.app.GetKeyStore().VRF(), th.vrfKeyID, prooflib.PreSeedDataV2{
		PreSeed:          s,
		BlockHash:        requestLog.Raw().BlockHash,
		BlockNum:         requestLog.Raw().BlockNumber,
		SubId:            sub.subID,
		CallbackGasLimit: req.callbackGasLimit,
		NumWords:         req.numWords,
		Sender:           th.eoaConsumerAddr,
	})
	require.NoError(t, err)
	txHash := requestLog.Raw().TxHash
	req.proof, req.reqCommitment, req.requestTxHash = &proof, &rc, &txHash
	return req
}

// Create VRF jobs in test CL node
func createVRFJobsNew(
	t *testing.T,
	fromKeys [][]ethkey.KeyV2,
	app *cltest.TestApplication,
	coordinator v2.CoordinatorV2_X,
	coordinatorAddress common.Address,
	batchCoordinatorAddress common.Address,
	uni coordinatorV2Universe,
	batchEnabled bool,
	chainID *big.Int,
	gasLanePrices ...*assets.Wei,
) (jobs []job.Job, vrfKeyIDs []string) {
	ctx := testutils.Context(t)
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

		vrfkey, err := app.GetKeyStore().VRF().Create(ctx)
		require.NoError(t, err)

		jid := uuid.New()
		incomingConfs := 2
		s := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
			JobID:                        jid.String(),
			Name:                         fmt.Sprintf("vrf-primary-%d", i),
			CoordinatorAddress:           coordinatorAddress.Hex(),
			BatchCoordinatorAddress:      batchCoordinatorAddress.Hex(),
			BatchFulfillmentEnabled:      batchEnabled,
			MinIncomingConfirmations:     incomingConfs,
			PublicKey:                    vrfkey.PublicKey.String(),
			FromAddresses:                keyStrs,
			BackoffInitialDelay:          10 * time.Millisecond,
			BackoffMaxDelay:              time.Second,
			V2:                           true,
			GasLanePrice:                 gasLanePrices[i],
			VRFOwnerAddress:              uni.vrfOwnerAddressNew.Hex(),
			CustomRevertsPipelineEnabled: true,
			EVMChainID:                   chainID.String(),
		}).Toml()
		jb, err := vrfcommon.ValidatedVRFSpec(s)
		t.Log(jb.VRFSpec.PublicKey.MustHash(), vrfkey.PublicKey.MustHash())
		require.NoError(t, err)
		err = app.JobSpawner().CreateJob(ctx, nil, &jb)
		require.NoError(t, err)
		registerProvingKeyHelper(t, uni.coordinatorV2UniverseCommon, coordinator, vrfkey, ptr(gasLanePrices[i].ToInt().Uint64()))
		jobs = append(jobs, jb)
		vrfKeyIDs = append(vrfKeyIDs, vrfkey.ID())
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

// Get txn receipt from txstore DB for a given txID. Useful to get status
// of a txn on chain, to check if it reverted or not
func getTxnReceiptDB(db *sqlx.DB, txesID int64) ([]v2.TxnReceiptDB, error) {
	sqlQuery := `
		WITH txes AS (
			SELECT *
			FROM evm.txes
			WHERE (state = 'confirmed' OR state = 'unconfirmed')
				AND id = $1
		), attempts AS (
			SELECT *
			FROM evm.tx_attempts
			WHERE eth_tx_id IN (SELECT id FROM txes)
		), receipts AS (
			SELECT *
			FROM evm.receipts
			WHERE tx_hash IN (SELECT hash FROM attempts)
		)
		SELECT r.tx_hash, 
			r.receipt,
			t.from_address,
			t.meta->>'SubId' as sub_id,
			COALESCE(t.meta->>'RequestID', '') as request_id,
			COALESCE(t.meta->>'RequestTxHash', '') as request_tx_hash,
			COALESCE(t.meta->>'ForceFulfillmentAttempt', '0') as force_fulfillment_attempt
		FROM receipts r
		INNER JOIN attempts a ON r.tx_hash = a.hash
		INNER JOIN txes t ON a.eth_tx_id = t.id
	`
	var recentReceipts []v2.TxnReceiptDB
	var err error
	if txesID != -1 {
		err = db.Select(&recentReceipts, sqlQuery, txesID)
	} else {
		sqlQuery = strings.Replace(sqlQuery, "AND id = $1", "AND meta->>'ForceFulfilled' IS NOT NULL", 1)
		err = db.Select(&recentReceipts, sqlQuery)
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "fetch_failed_receipts_txm")
	}

	return recentReceipts, nil
}

// Type to store VRF req details like requestID, proof, reqCommitment
type vrfReq struct {
	requestID        *big.Int
	callbackGasLimit uint32
	numWords         uint32
	txID             int64
	requestTxHash    *common.Hash
	proof            *vrf_coordinator_v2.VRFProof
	reqCommitment    *vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment
}

// Type to store VRF sub details like subID, balance
type vrfSub struct {
	subID   uint64
	balance uint64
}

// Test harness for handling reverted txns functionality
type revertTxnTH struct {
	// VRF Key Details
	key1     ethkey.KeyV2
	key2     ethkey.KeyV2
	vrfKeyID string
	keyHash  [32]byte

	// CL Node Details
	chainID *big.Int
	app     *cltest.TestApplication
	db      *sqlx.DB

	// Contract Details
	uni             *coordinatorV2Universe
	eoaConsumer     *vrf_external_sub_owner_example.VRFExternalSubOwnerExample
	eoaConsumerAddr common.Address

	// VRF Req Details
	subs []*vrfSub
}

// Constructor for handling reverted txns test harness
func newRevertTxnTH(t *testing.T,
	uni *coordinatorV2Universe,
	ownerKey ethkey.KeyV2,
	batchEnabled bool,
	subBalances []uint64) (th *revertTxnTH) {
	key1 := cltest.MustGenerateRandomKey(t)
	key2 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), toml.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: toml.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		}, toml.KeySpecific{
			// Gas lane.
			Key:          ptr(key2.EIP55Address),
			GasEstimator: toml.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1, key2)

	th = &revertTxnTH{
		key1: key1,
		key2: key2,
		app:  app,
		db:   db,
		uni:  uni,
		subs: make([]*vrfSub, len(subBalances)),
	}
	coordinator := uni.rootContract
	coordinatorAddress := uni.rootContractAddress
	th.chainID = th.uni.backend.Blockchain().Config().ChainID
	var err error

	th.eoaConsumerAddr, _, th.eoaConsumer, err = vrf_external_sub_owner_example.DeployVRFExternalSubOwnerExample(
		uni.neil,
		uni.backend,
		coordinatorAddress,
		uni.linkContractAddress,
	)
	require.NoError(t, err, "failed to deploy eoa consumer")
	uni.backend.Commit()

	for i := 0; i < len(subBalances); i++ {
		subID := uint64(i + 1)
		setupSub(t, th, subID, subBalances[i])
		th.subs[i] = &vrfSub{subID: subID, balance: subBalances[i]}
	}

	// Fund gas lanes.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	sendEth(t, ownerKey, uni.backend, key2.Address, 10)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job using key1 and key2 on the same gas lane.
	jbs, vrfKeyIDs := createVRFJobsNew(
		t,
		[][]ethkey.KeyV2{{key1, key2}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		*uni,
		batchEnabled,
		th.chainID,
		gasLanePriceWei)
	vrfKey := jbs[0].VRFSpec.PublicKey

	th.keyHash = vrfKey.MustHash()
	th.vrfKeyID = vrfKeyIDs[0]

	// Transfer ownership of the VRF coordinator to the VRF owner,
	// which is critical for this test.
	t.Log("vrf owner address:", uni.vrfOwnerAddressNew)
	_, err = coordinator.TransferOwnership(uni.neil, uni.vrfOwnerAddressNew)
	require.NoError(t, err, "unable to TransferOwnership of VRF coordinator to VRFOwner")
	uni.backend.Commit()

	_, err = uni.vrfOwnerNew.AcceptVRFOwnership(uni.neil)
	require.NoError(t, err, "unable to Accept VRF Ownership")
	uni.backend.Commit()

	actualCoordinatorAddr, err := uni.vrfOwnerNew.GetVRFCoordinator(nil)
	require.NoError(t, err)
	require.Equal(t, coordinatorAddress, actualCoordinatorAddr)

	// Add allowed callers so that the oracle can call fulfillRandomWords
	// on VRFOwner.
	_, err = uni.vrfOwnerNew.SetAuthorizedSenders(uni.neil, []common.Address{
		key1.EIP55Address.Address(),
		key2.EIP55Address.Address(),
	})
	require.NoError(t, err, "unable to update authorized senders in VRFOwner")
	uni.backend.Commit()

	return th
}

func setupSub(t *testing.T, th *revertTxnTH, subID uint64, balance uint64) {
	uni := th.uni
	coordinator := uni.rootContract
	coordinatorAddress := uni.rootContractAddress
	var err error

	// Create a subscription and fund with amount specified
	_, err = coordinator.CreateSubscription(uni.neil)
	require.NoError(t, err, "failed to create eoa sub")
	uni.backend.Commit()

	// Fund the sub
	b, err := evmutils.ABIEncode(`[{"type":"uint64"}]`, subID)
	require.NoError(t, err)
	_, err = uni.linkContract.TransferAndCall(
		uni.sergey, coordinatorAddress, big.NewInt(int64(balance)), b)
	require.NoError(t, err, "failed to fund sub")
	uni.backend.Commit()

	// Add the consumer to the sub
	subIDBig := big.NewInt(int64(subID))
	_, err = coordinator.AddConsumer(uni.neil, subIDBig, th.eoaConsumerAddr)
	require.NoError(t, err, "failed to add consumer")
	uni.backend.Commit()

	// Check the subscription state
	sub, err := coordinator.GetSubscription(nil, subIDBig)
	consumers := sub.Consumers()
	require.NoError(t, err, "failed to get subscription with id %d", subID)
	require.Equal(t, big.NewInt(int64(balance)), sub.Balance())
	require.Equal(t, 1, len(consumers))
	require.Equal(t, th.eoaConsumerAddr, consumers[0])
	require.Equal(t, uni.neil.From, sub.Owner())
}
