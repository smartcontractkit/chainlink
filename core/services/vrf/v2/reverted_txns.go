package v2

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	evmutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type (
	TxnReceiptDB struct {
		TxHash                  common.Hash      `db:"tx_hash"`
		EVMReceipt              evmtypes.Receipt `db:"receipt"`
		FromAddress             common.Address   `db:"from_address"`
		ToAddress               common.Address   `db:"to_address"`
		EncodedPayload          hexutil.Bytes    `db:"encoded_payload"`
		GasLimit                uint64           `db:"gas_limit"`
		SubID                   uint64           `db:"sub_id"`
		RequestID               string           `db:"request_id"`
		RequestTxHash           string           `db:"request_tx_hash"`
		ForceFulfillmentAttempt uint64           `db:"force_fulfillment_attempt"`
	}

	RevertedVRFTxn struct {
		DBReceipt  TxnReceiptDB
		IsBatchReq bool
		Proof      vrf_coordinator_v2.VRFProof
		Commitment vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment
	}
)

var ReqScanTimeRangeInDB = "1 hour"

func (lsn *listenerV2) runRevertedTxnsHandler(pollPeriod time.Duration) {
	pollPeriod = pollPeriod + time.Second*3
	tick := time.NewTicker(pollPeriod)
	defer tick.Stop()
	ctx, cancel := lsn.chStop.NewCtx()
	defer cancel()
	for {
		select {
		case <-lsn.chStop:
			return
		case <-tick.C:
			lsn.handleRevertedTxns(ctx, pollPeriod)
		}
	}
}

func (lsn *listenerV2) handleRevertedTxns(ctx context.Context, pollPeriod time.Duration) {
	lsn.l.Infow("Handling reverted txns")

	// Fetch recent single and batch txns, that have not been force-fulfilled
	recentSingleTxns, err := lsn.fetchRecentSingleTxns(ctx, lsn.ds, lsn.chainID.Uint64(), pollPeriod)
	if err != nil {
		lsn.l.Fatalw("Fetch recent txns", "err", err)
	}
	recentBatchTxns, err := lsn.fetchRecentBatchTxns(ctx, lsn.ds, lsn.chainID.Uint64(), pollPeriod)
	if err != nil {
		lsn.l.Fatalw("Fetch recent batch txns", "err", err)
	}
	recentForceFulfillmentTxns, err := lsn.fetchRevertedForceFulfilmentTxns(ctx, lsn.ds, lsn.chainID.Uint64(), pollPeriod)
	if err != nil {
		lsn.l.Fatalw("Fetch recent reverted force-fulfillment txns", "err", err)
	}
	recentTxns := make([]TxnReceiptDB, 0)
	if len(recentSingleTxns) > 0 {
		recentTxns = append(recentTxns, recentSingleTxns...)
	}
	if len(recentBatchTxns) > 0 {
		recentTxns = append(recentTxns, recentBatchTxns...)
	}
	if len(recentForceFulfillmentTxns) > 0 {
		recentTxns = append(recentTxns, recentForceFulfillmentTxns...)
	}

	// Query RPC using TransactionByHash to get the transaction object
	revertedTxns := lsn.filterRevertedTxns(ctx, recentTxns)

	// Extract calldata of function call from transaction object
	for _, revertedTxn := range revertedTxns {
		// Pass that to txm to create a new tx for force fulfillment
		_, err := lsn.enqueueForceFulfillmentForRevertedTxn(ctx, revertedTxn)
		if err != nil {
			lsn.l.Errorw("Enqueue force fulfilment", "err", err)
		}
	}
}

func (lsn *listenerV2) fetchRecentSingleTxns(ctx context.Context,
	ds sqlutil.DataSource,
	chainID uint64,
	pollPeriod time.Duration) ([]TxnReceiptDB, error) {
	// (state = 'confirmed' OR state = 'unconfirmed')
	sqlQuery := fmt.Sprintf(`
		WITH already_ff as (
			SELECT meta->>'RequestID' as request_id
			FROM evm.txes
			WHERE created_at >= NOW() - interval '%s'
				AND evm_chain_id = $1
				AND meta->>'ForceFulfilled' is NOT NULL
		), txes AS (
			SELECT *
			FROM evm.txes
			WHERE created_at >= NOW() - interval '%s'
				AND evm_chain_id = $1
				AND meta->>'SubId' IS NOT NULL
				AND meta->>'RequestID' IS NOT NULL
				AND meta->>'ForceFulfilled' is NULL
				AND meta->>'RequestID' NOT IN (SELECT request_id FROM already_ff)
		), attempts AS (
			SELECT *
			FROM evm.tx_attempts
			WHERE eth_tx_id IN (SELECT id FROM txes)
		), receipts AS (
			SELECT *
			FROM evm.receipts
			WHERE tx_hash IN (SELECT hash FROM attempts)
				AND receipt->>'status' = '0x0'
		)
		SELECT r.tx_hash, 
			r.receipt,
			t.from_address,
			t.to_address,
			t.encoded_payload,
			t.gas_limit,
			t.meta->>'SubId' as sub_id,
			t.meta->>'RequestID' as request_id,
			t.meta->>'RequestTxHash' as request_tx_hash
		FROM receipts r
		INNER JOIN attempts a ON r.tx_hash = a.hash
		INNER JOIN txes t ON a.eth_tx_id = t.id
	`, ReqScanTimeRangeInDB, ReqScanTimeRangeInDB)
	var recentReceipts []TxnReceiptDB

	before := time.Now()
	err := ds.SelectContext(ctx, &recentReceipts, sqlQuery, chainID)
	lsn.postSqlLog(ctx, before, pollPeriod, "FetchRecentSingleTxns")
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "Error fetching recent non-force-fulfilled txns")
	}

	recentReceipts = unique(recentReceipts)
	lsn.l.Infow("finished querying for recently reverting single fulfillments",
		"count", len(recentReceipts),
	)
	for _, r := range recentReceipts {
		lsn.l.Infow("found reverted fulfillment", "requestID", r.RequestID, "fulfillmentTxHash", r.TxHash.String())
	}
	return recentReceipts, nil
}

func (lsn *listenerV2) fetchRecentBatchTxns(ctx context.Context,
	ds sqlutil.DataSource,
	chainID uint64,
	pollPeriod time.Duration) ([]TxnReceiptDB, error) {
	sqlQuery := fmt.Sprintf(`
		WITH already_ff as (
			SELECT meta->>'RequestID' as request_id
			FROM evm.txes
			WHERE created_at >= NOW() - interval '%s'
				AND evm_chain_id = $1
				AND meta->>'ForceFulfilled' is NOT NULL
		), txes AS (
			SELECT *
			FROM (
				SELECT *
				FROM evm.txes
				WHERE created_at >= NOW() - interval '%s'
					AND evm_chain_id = $1
					AND meta->>'SubId' IS NOT NULL
					AND meta->>'RequestIDs' IS NOT NULL
					AND meta->>'ForceFulfilled' IS NULL
			) AS eth_txes1
			WHERE (meta->'RequestIDs' ?| (SELECT ARRAY_AGG(request_id) FROM already_ff)) IS NOT TRUE
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
			t.to_address,
			t.encoded_payload,
			t.gas_limit,
			t.meta->>'SubId' as sub_id
		FROM receipts r
		INNER JOIN attempts a ON r.tx_hash = a.hash
		INNER JOIN txes t ON a.eth_tx_id = t.id
	`, ReqScanTimeRangeInDB, ReqScanTimeRangeInDB)
	var recentReceipts []TxnReceiptDB

	before := time.Now()
	err := ds.SelectContext(ctx, &recentReceipts, sqlQuery, chainID)
	lsn.postSqlLog(ctx, before, pollPeriod, "FetchRecentBatchTxns")
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "Error fetching recent non-force-fulfilled txns")
	}

	recentReceipts = unique(recentReceipts)
	lsn.l.Infow("finished querying for recent batch fulfillments",
		"count", len(recentReceipts),
	)
	return recentReceipts, nil
}

func (lsn *listenerV2) fetchRevertedForceFulfilmentTxns(ctx context.Context,
	ds sqlutil.DataSource,
	chainID uint64,
	pollPeriod time.Duration) ([]TxnReceiptDB, error) {
	sqlQuery := fmt.Sprintf(`
		WITH txes AS (
			SELECT *
			FROM evm.txes
			WHERE created_at >= NOW() - interval '%s'
				AND evm_chain_id = $1
				AND meta->>'SubId' IS NOT NULL
				AND meta->>'RequestID' IS NOT NULL
				AND meta->>'ForceFulfilled' is NOT NULL
		), attempts AS (
			SELECT *
			FROM evm.tx_attempts
			WHERE eth_tx_id IN (SELECT id FROM txes)
		), receipts AS (
			SELECT *
			FROM evm.receipts
			WHERE tx_hash IN (SELECT hash FROM attempts)
				AND receipt->>'status' = '0x0'
		)
		SELECT r.tx_hash, 
			r.receipt,
			t.from_address,
			t.to_address,
			t.encoded_payload,
			t.gas_limit,
			t.meta->>'SubId' as sub_id,
			t.meta->>'RequestID' as request_id,
			t.meta->>'RequestTxHash' as request_tx_hash,
			CAST(COALESCE(t.meta->>'ForceFulfillmentAttempt', '0') AS INT) as force_fulfillment_attempt
		FROM receipts r
		INNER JOIN attempts a ON r.tx_hash = a.hash
		INNER JOIN txes t ON a.eth_tx_id = t.id
	`, ReqScanTimeRangeInDB)
	var recentReceipts []TxnReceiptDB

	before := time.Now()
	err := ds.SelectContext(ctx, &recentReceipts, sqlQuery, chainID)
	lsn.postSqlLog(ctx, before, pollPeriod, "FetchRevertedForceFulfilmentTxns")
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "Error fetching recent reverted force-fulfilled txns")
	}

	sqlQueryAll := fmt.Sprintf(`
		WITH txes AS (
			SELECT *
			FROM evm.txes
			WHERE created_at >= NOW() - interval '%s'
				AND evm_chain_id = $1
				AND meta->>'SubId' IS NOT NULL
				AND meta->>'RequestID' IS NOT NULL
				AND meta->>'ForceFulfilled' is NOT NULL
		), attempts AS (
			SELECT *
			FROM evm.tx_attempts
			WHERE eth_tx_id IN (SELECT id FROM txes)
		)
		SELECT a.hash as tx_hash,
			t.meta->>'SubId' as sub_id,
			t.meta->>'RequestID' as request_id,
			CAST(COALESCE(t.meta->>'ForceFulfillmentAttempt', '0') AS INT) as force_fulfillment_attempt
		FROM attempts a
		INNER JOIN txes t ON a.eth_tx_id = t.id
	`, ReqScanTimeRangeInDB)
	var allReceipts []TxnReceiptDB
	before = time.Now()
	err = ds.SelectContext(ctx, &allReceipts, sqlQueryAll, chainID)
	lsn.postSqlLog(ctx, before, pollPeriod, "Fetch all ForceFulfilment Txns")
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "Error fetching all recent force-fulfilled txns")
	}

	recentReceipts = UniqueByReqID(recentReceipts, allReceipts)

	lsn.l.Infow("finished querying for recently reverting reverted force-fulfillment txns",
		"count", len(recentReceipts),
	)
	for _, r := range recentReceipts {
		lsn.l.Infow("found reverted force-fulfillment txn", "requestID", r.RequestID,
			"fulfillmentTxHash", r.TxHash.String(),
			"ForceFulfillmentAttempt", r.ForceFulfillmentAttempt)
	}
	return unique(recentReceipts), nil
}

func unique(rs []TxnReceiptDB) (res []TxnReceiptDB) {
	if len(rs) == 0 {
		return
	}
	exists := make(map[string]bool)
	res = make([]TxnReceiptDB, 0)
	for _, r := range rs {
		if _, ok := exists[r.TxHash.Hex()]; ok {
			continue
		}
		res = append(res, r)
		exists[r.TxHash.Hex()] = true
	}
	return res
}

func UniqueByReqID(revertedForceTxns []TxnReceiptDB, allForceTxns []TxnReceiptDB) (res []TxnReceiptDB) {
	if len(revertedForceTxns) == 0 {
		return
	}

	// Load all force fulfillment txns into a map
	// allForceTxns would have successful, reverted and pending force fulfillment txns
	allForceTxnsMap := make(map[string]TxnReceiptDB)
	for _, r := range allForceTxns {
		if existingReceipt, ok := allForceTxnsMap[r.RequestID]; ok {
			// Get the latest force fulfillment attempt for a given RequestID
			if existingReceipt.ForceFulfillmentAttempt < r.ForceFulfillmentAttempt {
				allForceTxnsMap[r.RequestID] = r
			}
			continue
		}
		allForceTxnsMap[r.RequestID] = r
	}

	// Deduplicate reverted force fulfillment txns and skip/ignore reverted
	// force-fulfillment txns which have a pending force-fulfillment retry
	revertedForceTxnsMap := make(map[string]TxnReceiptDB)
	res = make([]TxnReceiptDB, 0)
	for _, forceTxn := range revertedForceTxns {
		// If there is a pending force fulfilment without a receipt yet, skip force-fulfilling it now again, until a txn receipt
		// This prevents a race between this Custom-VRF-Reverted-Txns-Pipeline and TransactionManager
		if receipt, ok := allForceTxnsMap[forceTxn.RequestID]; ok && receipt.ForceFulfillmentAttempt > forceTxn.ForceFulfillmentAttempt {
			continue
		}
		if existingReceipt, ok := revertedForceTxnsMap[forceTxn.RequestID]; ok {
			// Get the latest force fulfillment attempt for a given RequestID
			if existingReceipt.ForceFulfillmentAttempt < forceTxn.ForceFulfillmentAttempt {
				revertedForceTxnsMap[forceTxn.RequestID] = forceTxn
			}
			continue
		}
		revertedForceTxnsMap[forceTxn.RequestID] = forceTxn
	}

	// Load the deduplicated map into a list and return
	for _, r := range revertedForceTxnsMap {
		res = append(res, r)
	}
	return res
}

// postSqlLog logs about context cancellation and timing after a query returns.
// Queries which use their full timeout log critical level. More than 50% log error, and 10% warn.
func (lsn *listenerV2) postSqlLog(ctx context.Context, begin time.Time, pollPeriod time.Duration, queryName string) {
	elapsed := time.Since(begin)
	if ctx.Err() != nil {
		lsn.l.Debugw("SQL context canceled", "ms", elapsed.Milliseconds(), "err", ctx.Err(), "sql", queryName)
	}

	timeout := pollPeriod
	deadline, ok := ctx.Deadline()
	if ok {
		timeout = deadline.Sub(begin)
	}

	pct := float64(elapsed) / float64(timeout)
	pct *= 100

	kvs := []any{"ms", elapsed.Milliseconds(),
		"timeout", timeout.Milliseconds(),
		"percent", strconv.FormatFloat(pct, 'f', 1, 64),
		"sql", queryName}

	if elapsed >= timeout {
		lsn.l.Criticalw("ExtremelySlowSQLQuery", kvs...)
	} else if errThreshold := timeout / 5; errThreshold > 0 && elapsed > errThreshold {
		lsn.l.Errorw("VerySlowSQLQuery", kvs...)
	} else if warnThreshold := timeout / 10; warnThreshold > 0 && elapsed > warnThreshold {
		lsn.l.Warnw("SlowSQLQuery", kvs...)
	} else {
		lsn.l.Infow("SQLQueryLatency", kvs...)
	}
}

func (lsn *listenerV2) filterRevertedTxns(ctx context.Context,
	recentReceipts []TxnReceiptDB) []RevertedVRFTxn {
	revertedVRFTxns := make([]RevertedVRFTxn, 0)
	for _, txnReceipt := range recentReceipts {
		switch txnReceipt.ToAddress.Hex() {
		case lsn.vrfOwner.Address().Hex():
			fallthrough
		case lsn.coordinator.Address().Hex():
			// Filter Single VRF Fulfilment
			revertedVRFTxn, err := lsn.filterSingleRevertedTxn(ctx, txnReceipt)
			if err != nil {
				lsn.l.Errorw("Filter reverted single fulfillment txn", "Err", err)
				continue
			}
			// Revert reason is not insufficient balance
			if revertedVRFTxn == nil {
				continue
			}
			revertedVRFTxns = append(revertedVRFTxns, *revertedVRFTxn)
		case lsn.batchCoordinator.Address().Hex():
			// Filter Batch VRF Fulfilment
			revertedBatchVRFTxns, err := lsn.filterBatchRevertedTxn(ctx, txnReceipt)
			if err != nil {
				lsn.l.Errorw("Filter batchfulfilment with reverted txns", "Err", err)
				continue
			}
			// No req in the batch txn with insufficient balance revert reason
			if len(revertedBatchVRFTxns) == 0 {
				continue
			}
			revertedVRFTxns = append(revertedVRFTxns, revertedBatchVRFTxns...)
		default:
			// Unrecognised Txn
			lsn.l.Warnw("Unrecognised txn in VRF-Reverted-Pipeline",
				"ToAddress", txnReceipt.ToAddress.Hex(),
			)
		}
	}

	lsn.l.Infow("Reverted VRF fulfilment txns due to InsufficientBalance",
		"count", len(revertedVRFTxns),
		"reverted_txns", revertedVRFTxns,
	)
	for _, r := range revertedVRFTxns {
		lsn.l.Infow("Reverted VRF fulfilment txns due to InsufficientBalance",
			"RequestID", r.DBReceipt.RequestID,
			"TxnStoreEVMReceipt.BlockHash", r.DBReceipt.EVMReceipt.BlockHash.String(),
			"TxnStoreEVMReceipt.BlockNumber", r.DBReceipt.EVMReceipt.BlockNumber.String(),
			"VRFFulfillmentTxHash", r.DBReceipt.TxHash.String())
	}
	return revertedVRFTxns
}

func (lsn *listenerV2) filterSingleRevertedTxn(ctx context.Context,
	txnReceiptDB TxnReceiptDB) (
	*RevertedVRFTxn, error) {
	requestID := common.HexToHash(txnReceiptDB.RequestID).Big()
	commitment, err := lsn.coordinator.GetCommitment(&bind.CallOpts{Context: ctx}, requestID)
	if err != nil {
		// Not able to get commitment from chain RPC node, continue
		lsn.l.Errorw("Force-fulfilment of single reverted txns: Not able to get commitment from chain RPC node", "err", err)
	} else if utils.IsEmpty(commitment[:]) {
		// VRF request already fulfilled, return
		return nil, nil
	}
	lsn.l.Infow("Single reverted txn: Unfulfilled req", "req", requestID.String())

	// Get txn object from RPC node
	ethClient := lsn.chain.Client()
	tx, err := ethClient.TransactionByHash(ctx, txnReceiptDB.TxHash)
	if err != nil {
		return nil, errors.Wrap(err, "get_txn_by_hash")
	}

	// Simulate txn to get revert error
	call := ethereum.CallMsg{
		From:     txnReceiptDB.FromAddress,
		To:       &txnReceiptDB.ToAddress,
		Data:     tx.Data(), // txnReceiptDB.EncodedPayload,
		Gas:      txnReceiptDB.GasLimit,
		GasPrice: tx.GasPrice(),
	}
	_, rpcError := ethClient.CallContract(ctx, call, txnReceiptDB.EVMReceipt.BlockNumber)
	if rpcError == nil {
		return nil, fmt.Errorf("error fetching revert reason %v: %v", txnReceiptDB.TxHash, err)
	}
	revertErr, err := evmclient.ExtractRPCError(rpcError)
	lsn.l.Infow("InsufficientBalRevertedTxn",
		"RawRevertData", rpcError,
		"ParsedRevertData", revertErr.Data,
		"ParsingErr", err,
	)
	if err != nil {
		return nil, fmt.Errorf("reverted_txn_reason_parse_err: %v", err)
	}
	revertErrDataStr := ""
	revertErrDataBytes := []byte{}
	if revertErr.Data != nil {
		revertErrDataStr = revertErr.Data.(string)
		revertErrDataStr = strings.Replace(revertErrDataStr, "Reverted ", "", 1)
		// If force fulfillment txn reverts on chain due to getFeedData not falling back
		// to MAXINT256 due to stalenessSeconds criteria not satisfying
		revertErrDataBytes = common.FromHex(revertErrDataStr)
	}
	insufficientErr := coordinatorV2ABI.Errors["InsufficientBalance"].ID.Bytes()[0:4]
	// Revert reason may not be accurately determined from all RPC nodes and may
	// not work in some chains
	if len(revertErrDataStr) > 0 && !bytes.Equal(revertErrDataBytes[0:4], insufficientErr) {
		return nil, nil
	}
	// If reached maximum number of retries for force fulfillment
	if txnReceiptDB.ForceFulfillmentAttempt >= 15 {
		return nil, nil
	}

	// Get VRF fulfillment proof and commitment from tx object
	txData := txnReceiptDB.EncodedPayload
	if len(txData) <= 4 {
		return nil, fmt.Errorf("invalid_txn_data_for_tx: %s", tx.Hash().String())
	}
	callData := txData[4:] // Remove first 4 bytes of function signature
	unpacked, err := coordinatorV2ABI.Methods["fulfillRandomWords"].Inputs.Unpack(callData)
	if err != nil {
		return nil, fmt.Errorf("invalid_txn_data_for_tx_pack: %s, err %v", tx.Hash().String(), err)
	}
	proof := abi.ConvertType(unpacked[0], new(vrf_coordinator_v2.VRFProof)).(*vrf_coordinator_v2.VRFProof)
	reqCommitment := abi.ConvertType(unpacked[1], new(vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment)).(*vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment)
	return &RevertedVRFTxn{
		DBReceipt:  txnReceiptDB,
		IsBatchReq: false,
		Proof:      *proof,
		Commitment: *reqCommitment}, nil
}

func (lsn *listenerV2) filterBatchRevertedTxn(ctx context.Context,
	txnReceiptDB TxnReceiptDB) (
	[]RevertedVRFTxn, error) {
	if len(txnReceiptDB.EncodedPayload) <= 4 {
		return nil, fmt.Errorf("invalid encodedPayload: %v", hexutil.Encode(txnReceiptDB.EncodedPayload))
	}
	unpackedInputs, err := batchCoordinatorV2ABI.Methods["fulfillRandomWords"].Inputs.Unpack(txnReceiptDB.EncodedPayload[4:])
	if err != nil {
		return nil, errors.Wrap(err, "cannot_unpack_batch_txn")
	}
	proofs := abi.ConvertType(unpackedInputs[0], new([]vrf_coordinator_v2.VRFProof)).(*[]vrf_coordinator_v2.VRFProof)
	reqCommitments := abi.ConvertType(unpackedInputs[1], new([]vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment)).(*[]vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment)

	proofReqIDs := make([]common.Hash, 0)
	keyHash := lsn.job.VRFSpec.PublicKey.MustHash()
	for _, proof := range *proofs {
		payload, err := evmutils.ABIEncode(`[{"type":"bytes32"},{"type":"uint256"}]`, keyHash, proof.Seed)
		if err != nil {
			return nil, fmt.Errorf("ABI Encode Error: (err %v), (keyHash %v), (prood: %v)", err, keyHash, proof.Seed)
		}
		requestIDOfProof := common.BytesToHash(crypto.Keccak256(payload))
		proofReqIDs = append(proofReqIDs, requestIDOfProof)
	}

	// BatchVRFCoordinatorV2
	revertedTxns := make([]RevertedVRFTxn, 0)
	for _, log := range txnReceiptDB.EVMReceipt.Logs {
		if log.Topics[0] != batchCoordinatorV2ABI.Events["RawErrorReturned"].ID {
			continue
		}

		// Extract revert reason for individual req in batch txn
		unpacked, err := batchCoordinatorV2ABI.Events["RawErrorReturned"].Inputs.Unpack(log.Data)
		if err != nil {
			lsn.l.Errorw("cannot_unpack_batch_coordinator_log", "err", err)
			continue
		}
		lowLevelData := unpacked[0].([]byte)
		if !bytes.Equal(lowLevelData, coordinatorV2ABI.Errors["InsufficientBalance"].ID.Bytes()[0:4]) {
			continue
		}

		// Match current log to a (proof, commitment) pair from rawTxData using requestID
		requestID := log.Topics[1]
		var curProof vrf_coordinator_v2.VRFProof
		var curReqCommitment vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment
		found := false
		for i, proof := range *proofs {
			requestIDOfProof := proofReqIDs[i]
			if requestID == requestIDOfProof {
				found = true
				curProof = proof
				curReqCommitment = (*reqCommitments)[i]
				break
			}
		}

		if found {
			commitment, err := lsn.coordinator.GetCommitment(&bind.CallOpts{Context: ctx}, requestID.Big())
			if err != nil {
				// Not able to get commitment from chain RPC node, continue
				lsn.l.Errorw("Force-fulfilment of batch reverted txns: Not able to get commitment from chain RPC node",
					"err", err,
					"requestID", requestID.Big())
			} else if utils.IsEmpty(commitment[:]) {
				lsn.l.Infow("Batch fulfillment with initial reverted fulfillment txn and later successful fulfillment, Skipping", "req", requestID.String())
				continue
			}
			lsn.l.Infow("Batch fulfillment with reverted fulfillment txn", "req", requestID.String())
			revertedTxn := RevertedVRFTxn{
				DBReceipt: TxnReceiptDB{
					TxHash:      txnReceiptDB.TxHash,
					EVMReceipt:  txnReceiptDB.EVMReceipt,
					FromAddress: txnReceiptDB.FromAddress,
					SubID:       txnReceiptDB.SubID,
					RequestID:   requestID.Hex(),
				},
				IsBatchReq: true,
				Proof:      curProof,
				Commitment: curReqCommitment,
			}
			revertedTxns = append(revertedTxns, revertedTxn)
		} else {
			lsn.l.Criticalw("Reverted Batch fulfilment requestID from log does not have proof in req EncodedPayload",
				"requestIDFromLog", requestID.Big().Int64(),
			)
		}
	}
	return revertedTxns, nil
}

// enqueueForceFulfillment enqueues a forced fulfillment through the
// VRFOwner contract. It estimates gas again on the transaction due
// to the extra steps taken within VRFOwner.fulfillRandomWords.
func (lsn *listenerV2) enqueueForceFulfillmentForRevertedTxn(
	ctx context.Context,
	revertedTxn RevertedVRFTxn,
) (etx txmgr.Tx, err error) {
	if lsn.job.VRFSpec.VRFOwnerAddress == nil {
		return txmgr.Tx{}, errors.New("vrf_owner_not_set_in_job_spec")
	}

	proof := revertedTxn.Proof
	reqCommitment := revertedTxn.Commitment

	fromAddresses := lsn.fromAddresses()
	fromAddress, err := lsn.gethks.GetRoundRobinAddress(ctx, lsn.chainID, fromAddresses...)
	if err != nil {
		return txmgr.Tx{}, errors.Wrap(err, "failed_to_get_vrf_listener_from_address")
	}

	// fulfill the request through the VRF owner
	lsn.l.Infow("VRFOwner.fulfillRandomWords vs. VRFCoordinatorV2.fulfillRandomWords",
		"vrf_owner.fulfillRandomWords", hexutil.Encode(vrfOwnerABI.Methods["fulfillRandomWords"].ID),
		"vrf_coordinator_v2.fulfillRandomWords", hexutil.Encode(coordinatorV2ABI.Methods["fulfillRandomWords"].ID),
	)

	vrfOwnerAddress1 := lsn.vrfOwner.Address()
	vrfOwnerAddressSpec := lsn.job.VRFSpec.VRFOwnerAddress.Address()
	lsn.l.Infow("addresses diff", "wrapper_address", vrfOwnerAddress1, "spec_address", vrfOwnerAddressSpec)

	txData, err := vrfOwnerABI.Pack("fulfillRandomWords", proof, reqCommitment)
	if err != nil {
		return txmgr.Tx{}, errors.Wrap(err, "abi pack VRFOwner.fulfillRandomWords")
	}
	vrfOwnerCoordinator, _ := lsn.vrfOwner.GetVRFCoordinator(nil)
	lsn.l.Infow("RevertedTxnForceFulfilment EstimatingGas",
		"EncodedPayload", hexutil.Encode(txData),
		"VRFOwnerCoordinator", vrfOwnerCoordinator.String(),
	)
	ethClient := lsn.chain.Client()
	estimateGasLimit, err := ethClient.EstimateGas(ctx, ethereum.CallMsg{
		From: fromAddress,
		To:   &vrfOwnerAddressSpec,
		Data: txData,
	})
	if err != nil {
		return txmgr.Tx{}, errors.Wrap(err, "failed to estimate gas on VRFOwner.fulfillRandomWords")
	}
	estimateGasLimit = uint64(1.4 * float64(estimateGasLimit))

	lsn.l.Infow("Estimated gas limit on force fulfillment", "estimateGasLimit", estimateGasLimit)

	reqID := common.BytesToHash(hexutil.MustDecode(revertedTxn.DBReceipt.RequestID))
	var reqTxHash common.Hash
	if revertedTxn.DBReceipt.RequestTxHash != "" {
		reqTxHash = common.BytesToHash(hexutil.MustDecode(revertedTxn.DBReceipt.RequestTxHash))
	}
	lsn.l.Infow("RevertedTxnForceFulfilment CreateTransaction",
		"RequestID", revertedTxn.DBReceipt.RequestID,
		"RequestTxHash", revertedTxn.DBReceipt.RequestTxHash,
	)
	forceFulfiled := true
	forceFulfillmentAttempt := revertedTxn.DBReceipt.ForceFulfillmentAttempt + 1
	etx, err = lsn.chain.TxManager().CreateTransaction(ctx, txmgr.TxRequest{
		FromAddress:    fromAddress,
		ToAddress:      lsn.vrfOwner.Address(),
		EncodedPayload: txData,
		FeeLimit:       estimateGasLimit,
		Strategy:       txmgrcommon.NewSendEveryStrategy(),
		Meta: &txmgr.TxMeta{
			RequestID:               &reqID,
			SubID:                   &revertedTxn.DBReceipt.SubID,
			RequestTxHash:           &reqTxHash,
			ForceFulfilled:          &forceFulfiled,
			ForceFulfillmentAttempt: &forceFulfillmentAttempt,
			// No max link since simulation failed
		},
	})
	return etx, err
}
