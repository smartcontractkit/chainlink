package txmgr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/common/config"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
)

type stuckTxDetectorGasEstimator interface {
	GetFee(ctx context.Context, calldata []byte, feeLimit uint64, maxFeePrice *assets.Wei, opts ...feetypes.Opt) (fee gas.EvmFee, chainSpecificFeeLimit uint64, err error)
}

type stuckTxDetectorClient interface {
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error
}

type stuckTxDetectorTxStore interface {
	FindUnconfirmedTxsByFromAddresses(ctx context.Context, addresses []common.Address, chainID *big.Int) (txs []Tx, err error)
}

type stuckTxDetectorConfig interface {
	AutoPurgeStuckTxs() bool
	AutoPurgeThreshold() uint32
	AutoPurgeMinAttempts() uint32
	AutoPurgeDetectionApiUrl() *url.URL
}

type stuckTxDetector struct {
	lggr      logger.Logger
	chainID   *big.Int
	chainType config.ChainType
	cfg       stuckTxDetectorConfig

	gasEstimator stuckTxDetectorGasEstimator
	txStore      stuckTxDetectorTxStore
	chainClient  stuckTxDetectorClient
	httpClient   *http.Client

	purgeBlockNumLock sync.RWMutex
	purgeBlockNumMap  map[common.Address]int64 // Tracks the last block num a tx was purged for each from address if the PurgeOverflowTxs feature is enabled
}

func NewStuckTxDetector(lggr logger.Logger, chainID *big.Int, chainType config.ChainType, cfg stuckTxDetectorConfig, gasEstimator stuckTxDetectorGasEstimator, txStore stuckTxDetectorTxStore, chainClient stuckTxDetectorClient) *stuckTxDetector {
	// TODO: ensure to initialize client with the usual security standards
	// TODO: Load purgeBlockNumMap with some DB state or confirm rate limit is not needed on first purge after restart
	return &stuckTxDetector{
		lggr:             lggr,
		chainID:          chainID,
		chainType:        chainType,
		cfg:              cfg,
		gasEstimator:     gasEstimator,
		txStore:          txStore,
		chainClient:      chainClient,
		httpClient:       &http.Client{},
		purgeBlockNumMap: make(map[common.Address]int64),
	}
}

// If the AutoPurgeStuckTxs feature is enabled, finds terminally stuck transactions
// Uses a chain specific method for detection, or if one does not exist, applies a general heuristic
func (d *stuckTxDetector) DetectStuckTransactions(ctx context.Context, enabledAddresses []common.Address, blockNum int64) ([]Tx, error) {
	if !d.cfg.AutoPurgeStuckTxs() {
		return nil, nil
	}
	txs, err := d.FindPotentialStuckTxs(ctx, enabledAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to find potential stuck transactions: %w", err)
	}
	// No transactions found
	if len(txs) == 0 {
		return nil, nil
	}

	var stuckTxs []Tx
	switch d.chainType {
	case config.ChainScroll:
		stuckTxs, err = d.detectStuckTransactionsScroll(ctx, txs)
	case config.ChainZkEvm:
		stuckTxs, err = d.detectStuckTransactionsZkEVM(ctx, txs)
	default:
		stuckTxs, err = d.detectStuckTransactionsHeuristic(ctx, txs, blockNum)
	}

	for _, stuckTx := range stuckTxs {
		lggr := stuckTx.GetLogger(d.lggr)
		lggr.Errorw("marking transaction as terminally stuck", "etx", stuckTx)
	}

	return stuckTxs, err
}

// Finds the lowest nonce Unconfirmed transaction for each enabled address
// Only the earliest transaction can be considered terminally stuck. All others may be valid and just stuck behind the nonce
func (d *stuckTxDetector) FindPotentialStuckTxs(ctx context.Context, enabledAddresses []common.Address) ([]Tx, error) {
	// Loads attempts within tx
	txs, err := d.txStore.FindUnconfirmedTxsByFromAddresses(ctx, enabledAddresses, d.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve unconfirmed transactions for enabled addresses: %w", err)
	}
	// Stores the lowest nonce tx found in the query results for each from address
	lowestNonceTxMap := make(map[common.Address]Tx)
	for _, tx := range txs {
		if _, ok := lowestNonceTxMap[tx.FromAddress]; !ok {
			lowestNonceTxMap[tx.FromAddress] = tx
		} else if lowestNonceTx := lowestNonceTxMap[tx.FromAddress]; *lowestNonceTx.Sequence > *tx.Sequence {
			lowestNonceTxMap[tx.FromAddress] = tx
		}
	}

	// Build list of potentially stuck tx but exclude any that are already marked for purge
	var stuckTxs []Tx
	for _, tx := range lowestNonceTxMap {
		// Attempts are loaded newest to oldest so one marked for purge will always be first
		if len(tx.TxAttempts) > 0 && !tx.TxAttempts[0].IsPurgeAttempt {
			stuckTxs = append(stuckTxs, tx)
		}
	}

	return stuckTxs, nil
}

// Uses a heuristic to determine a stuck transaction potentially due to overflow
// This method can be unreliable and may result in false positives but it is best effort to keep the TXM from getting blocked
// 1. Check if AutoPurgeThreshold amount of blocks have passed since the last purge of a tx for the same fromAddress
// 2. If 1 is true, check if AutoPurgeThreshold amount of blocks have passed since the initial broadcast
// 3. If 2 is true, check if the transaction has at least AutoPurgeMinAttempts amount of broadcasted attempts
// 4. If 3 is true, check if the latest attempt's gas price is higher than what our gas estimator's GetFee method returns
// 5. If 4 is true, the transaction is likely stuck due to overflow
func (d *stuckTxDetector) detectStuckTransactionsHeuristic(ctx context.Context, txs []Tx, blockNum int64) ([]Tx, error) {
	d.purgeBlockNumLock.RLock()
	defer d.purgeBlockNumLock.RUnlock()
	// Get gas price from internal gas estimator
	// Send with arbitrarily high max gas price to prevent the results from being capped. Need the market gas price here.
	marketGasPrice, _, err := d.gasEstimator.GetFee(ctx, []byte{}, 0, assets.Ether(100))
	if err != nil {
		return txs, fmt.Errorf("failed to get market gas price for overflow detection: %w", err)
	}
	var stuckTxs []Tx
	for _, tx := range txs {
		// 1. Check if AutoPurgeThreshold amount of blocks have passed since the last purge of a tx for the same fromAddress
		// Used to rate limit purging to prevent a potential valid tx that was stuck behind an overflow tx from also getting purged without having enough time to be confirmed
		d.purgeBlockNumLock.RLock()
		lastPurgeBlockNum := d.purgeBlockNumMap[tx.FromAddress]
		d.purgeBlockNumLock.RUnlock()
		if lastPurgeBlockNum > blockNum-int64(d.cfg.AutoPurgeThreshold()) {
			continue
		}
		// Tx attempts are loaded from newest to oldest
		oldestBroadcastAttempt, newestBroadcastAttempt, broadcastedAttemptsCount := findBroadcastedAttempts(tx)
		// 2. Check if AutoPurgeThreshold amount of blocks have passed since the oldest attempt's broadcast block num
		if *oldestBroadcastAttempt.BroadcastBeforeBlockNum > blockNum-int64(d.cfg.AutoPurgeThreshold()) {
			continue
		}
		// 3. Check if the transaction has at least AutoPurgeMinAttempts amount of broadcasted attempts
		if broadcastedAttemptsCount < d.cfg.AutoPurgeMinAttempts() {
			continue
		}
		// 4. Check if the newest broadcasted attempt's gas price is higher than what our gas estimator's GetFee method returns
		if compareGasFees(newestBroadcastAttempt.TxFee, marketGasPrice) <= 0 {
			continue
		}
		// 5. Return the transaction since it is likely stuck due to overflow
		stuckTxs = append(stuckTxs, tx)
	}
	return stuckTxs, nil
}

func compareGasFees(attemptGas gas.EvmFee, marketGas gas.EvmFee) int {
	if attemptGas.Legacy != nil && marketGas.Legacy != nil {
		return attemptGas.Legacy.Cmp(marketGas.Legacy)
	}
	if attemptGas.DynamicFeeCap.Cmp(marketGas.DynamicFeeCap) == 0 {
		return attemptGas.DynamicTipCap.Cmp(marketGas.DynamicTipCap)
	}
	return attemptGas.DynamicFeeCap.Cmp(marketGas.DynamicFeeCap)

}

// Assumes tx attempts are loaded newest to oldest
func findBroadcastedAttempts(tx Tx) (oldestAttempt TxAttempt, newestAttempt TxAttempt, broadcastedCount uint32) {
	foundNewest := false
	for _, attempt := range tx.TxAttempts {
		if attempt.State != types.TxAttemptBroadcast {
			continue
		}
		if !foundNewest {
			newestAttempt = attempt
			foundNewest = true
		}
		oldestAttempt = attempt
		broadcastedCount++
	}
	return
}

type scrollRequest struct {
	Txs []string `json:"txs"`
}

type scrollResponse struct {
	Errcode int            `json:"errcode"`
	Errmsg  string         `json:"errmsg"`
	Data    map[string]int `json:"data"`
}

// Uses the custom Scroll skipped endpoint to determine an overflow transaction
func (d *stuckTxDetector) detectStuckTransactionsScroll(ctx context.Context, txs []Tx) ([]Tx, error) {
	if d.cfg.AutoPurgeDetectionApiUrl() == nil {
		return nil, fmt.Errorf("expected AutoPurgeDetectionApiUrl config to be set for chain type: %s", d.chainType)
	}

	attemptHashMap := make(map[string]Tx)

	request := new(scrollRequest)
	// Populate the request with the tx hash of the latest broadcast attempt from every tx
	for _, tx := range txs {
		for _, attempt := range tx.TxAttempts {
			if attempt.State == types.TxAttemptBroadcast {
				request.Txs = append(request.Txs, attempt.Hash.String())
				attemptHashMap[attempt.Hash.String()] = tx
				break
			}
		}
	}
	jsonReq, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json request %v for custom endpoint: %w", request, err)
	}

	// Build http post request
	url := fmt.Sprintf("%s/v1/sequencer/tx/skipped", d.cfg.AutoPurgeDetectionApiUrl())
	bodyReader := bytes.NewReader(jsonReq)
	postReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to make new request with context: %w", err)
	}
	// Send request
	resp, err := d.httpClient.Do(postReq)
	if err != nil {
		return nil, fmt.Errorf("request to scroll's custom endpoint failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	// Decode the response into expected type
	scrollResp := new(scrollResponse)
	err = json.NewDecoder(resp.Body).Decode(scrollResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response into struct: %w", err)
	}
	if scrollResp.Errcode != 0 || scrollResp.Errmsg != "" {
		return nil, fmt.Errorf("scroll's custom endpoint returned an error with code: %d, message: %s", scrollResp.Errcode, scrollResp.Errmsg)
	}

	// Return all transactions marked with status 1 signaling they have been skipped due to overflow
	var stuckTx []Tx
	for hash, status := range scrollResp.Data {
		if status == 1 {
			stuckTx = append(stuckTx, attemptHashMap[hash])
		}
	}

	return stuckTx, nil
}

// Uses eth_getTransactionByHash to detect that a transaction has been discarded due to overflow
// Currently only used by zkEVM but if other chains follow the same behavior in the future
func (d *stuckTxDetector) detectStuckTransactionsZkEVM(ctx context.Context, txs []Tx) ([]Tx, error) {
	txReqs := make([]rpc.BatchElem, len(txs))
	txHashMap := make(map[common.Hash]Tx)
	txRes := make([]*map[string]interface{}, len(txs))

	// Build batch request elems to perform
	// Does not need to be separated out into smaller batches
	// Max number of transactions to check is equal to the number of enabled addresses which is a relatively small amount
	for i, tx := range txs {
		latestAttemptHash := tx.TxAttempts[0].Hash
		var result map[string]interface{}
		txReqs[i] = rpc.BatchElem{
			Method: "eth_getTransactionByHash",
			Args: []interface{}{
				latestAttemptHash,
			},
			Result: &result,
		}
		txHashMap[latestAttemptHash] = tx
		txRes[i] = &result
	}

	// Send batch request
	err := d.chainClient.BatchCallContext(ctx, txReqs)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by hash in batch: %w", err)
	}

	// Parse results to find tx skipped due to zk overflow
	// If the result is nil, the transaction was discarded due to overflow
	var stuckTxs []Tx
	for i, req := range txReqs {
		txHash := req.Args[0].(common.Hash)
		if req.Error != nil {
			d.lggr.Debugf("failed to get transaction by hash (%s): %w", txHash.String(), req.Error)
			continue
		}
		result := *txRes[i]
		if result == nil {
			tx := txHashMap[txHash]
			stuckTxs = append(stuckTxs, tx)
		}
	}
	return stuckTxs, nil
}

// Once a purged tx's empty attempt is confirmed, this method is used to set at which block num the tx was purged at for the fromAddress
func (d *stuckTxDetector) SetPurgeBlockNum(fromAddress common.Address, blockNum int64) {
	d.purgeBlockNumLock.Lock()
	defer d.purgeBlockNumLock.Unlock()
	d.purgeBlockNumMap[fromAddress] = blockNum
}

// TODO: Find better place to store error messages so they can be served through tx status function for product teams
func (d *stuckTxDetector) StuckTxFatalError() *string {
	var errorMsg string
	switch d.chainType {
	case config.ChainScroll, config.ChainZkEvm:
		errorMsg = "transaction skipped due to ZK overflow"
	default:
		errorMsg = "purged terminally stuck transaction"
	}

	return &errorMsg
}
