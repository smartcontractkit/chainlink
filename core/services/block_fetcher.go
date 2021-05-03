package services

import (
	"context"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type BlockWithReceipts struct {
	block    *types.Block
	receipts []*bulletprooftxmanager.Receipt
}

type BlockFetcher struct {
	store     *strpkg.Store
	logger    *logger.Logger
	addresses map[common.Address]struct{}
}

func NewBlockFetcher(store *strpkg.Store) *BlockFetcher {
	return &BlockFetcher{
		store: store,
		logger: logger.CreateLogger(logger.Default.With(
			"id", "block_fetcher",
		)),
		addresses: make(map[common.Address]struct{}),
	}
}

func (ht *BlockFetcher) fetchBlock(ctx context.Context, head models.Head, blockFetchDuration prometheus.Histogram,
	receiptFetchDuration prometheus.Histogram, receiptLimitedFetchDuration prometheus.Histogram, receiptCount prometheus.Counter) (*BlockWithReceipts, error) {

	start := time.Now()
	block, err := ht.store.EthClient.FastBlockByHash(ctx, head.Hash)
	if err != nil {
		return nil, errors.Wrap(err, "HeadTracker#handleNewHighestHead failed fetching BlockByHash")
	}
	elapsed := time.Since(start)

	blockFetchDuration.Observe(float64(elapsed.Milliseconds()))
	ht.logger.Debugw("========================= HeadTracker: getting whole block took", "elapsedMs", elapsed.Milliseconds())

	logger.Warnf("rpcBlock.Hash: %v", block.Hash())
	logger.Warnf("rpcBlock.num Transactions: %v", block.Transactions().Len())

	start2 := time.Now()

	ht.logger.Debugf("========================= HeadTracker: getting receipts for %v transactions", block.Transactions().Len())

	receipts, err := ht.batchFetchReceipts(ctx, block.Transactions())
	if err != nil {
		return nil, errors.Wrap(err, "HeadTracker#batchFetchReceipts failed ")
	}

	topicCount := 0
	receiptsSizeBytes := 0
	for _, receipt := range receipts {
		for _, log := range receipt.Logs {
			topicCount += len(log.Topics)
			receiptsSizeBytes += len(log.Data)
			receiptsSizeBytes += (len(log.Topics) + 2) * common.HashLength
			receiptsSizeBytes += common.AddressLength
		}
	}

	elapsed2 := time.Since(start2)
	ht.logger.Debugw("========================= HeadTracker: getting tx receipts took",
		"elapsedMs", elapsed2.Milliseconds(), "receiptsSizeBytes", receiptsSizeBytes, "count", len(receipts), "topicCount", topicCount)

	receiptCount.Add(float64(block.Transactions().Len()))
	receiptFetchDuration.Observe(float64(elapsed2.Milliseconds()))

	takeCount := 10
	if takeCount > block.Transactions().Len() {
		takeCount = block.Transactions().Len()
	}
	var preserved []*types.Transaction = block.Transactions()[:takeCount]

	//for _, transaction := range block.Transactions() {
	//	to := transaction.To()
	//	if to == nil {
	//		preserved = append(preserved, transaction)
	//	} else {
	//		for address := range ht.addresses {
	//			if address == *to {
	//				preserved = append(preserved, transaction)
	//			} else {
	//				var signer types.Signer
	//				if transaction.Protected() {
	//					signer = types.LatestSignerForChainID(transaction.ChainId())
	//				} else {
	//					signer = types.HomesteadSigner{}
	//				}
	//
	//				from, _ := types.Sender(signer, transaction)
	//
	//				if address == from {
	//					preserved = append(preserved, transaction)
	//				}
	//			}
	//		}
	//	}
	//}

	ht.logger.Debugf("========================= HeadTracker: getting limited receipts for %v transactions", takeCount)
	start3 := time.Now()
	_, err = ht.batchFetchReceipts(ctx, preserved)
	if err != nil {
		return nil, errors.Wrap(err, "HeadTracker#batchFetchReceipts failed ")
	}

	elapsed3 := time.Since(start3)
	ht.logger.Debugw("========================= HeadTracker: getting limited tx receipts took",
		"elapsedMs", elapsed3.Milliseconds(), "count", len(preserved))

	receiptLimitedFetchDuration.Observe(float64(elapsed3.Milliseconds()))

	return &BlockWithReceipts{
		block,
		receipts,
	}, nil
}

type result struct {
	index int
	res   *bulletprooftxmanager.Receipt
	err   error
}

type resultReceipts struct {
	index int
	res   []*bulletprooftxmanager.Receipt
	err   error
}

// boundedParallelGet sends requests in parallel but only up to a certain
// limit, and furthermore it's only parallel up to the amount of CPUs but
// is always concurrent up to the concurrency limit
func (ht *BlockFetcher) getResultsParallel(ctx context.Context, txs []*types.Transaction) (receipts []*bulletprooftxmanager.Receipt, err error) {

	concurrencyLimit := 10
	// this buffered channel will block at the concurrency limit
	semaphoreChan := make(chan struct{}, concurrencyLimit)

	// this channel will not block and collect the http request results
	resultsChan := make(chan *result)

	// make sure we close these channels when we're done with them
	defer func() {
		close(semaphoreChan)
		close(resultsChan)
	}()

	// keen an index and loop through every url we will send a request to
	for i, tx := range txs {

		// start a go routine with the index and url in a closure
		go func(i int, tx *types.Transaction) {

			// this sends an empty struct into the semaphoreChan which
			// is basically saying add one to the limit, but when the
			// limit has been reached block until there is room
			semaphoreChan <- struct{}{}

			// send the request and put the response in a result struct
			// along with the index so we can sort them later along with
			// any error that might have occoured

			receipt, err := ht.store.EthClient.TransactionReceipt(ctx, tx.Hash())
			//if err != nil {
			//	return nil, err
			//}

			result := &result{i, bulletprooftxmanager.FromGethReceipt(receipt), err}

			// now we can send the result struct through the resultsChan
			resultsChan <- result

			// once we're done it's we read from the semaphoreChan which
			// has the effect of removing one from the limit and allowing
			// another goroutine to start
			<-semaphoreChan

		}(i, tx)
	}

	// make a slice to hold the results we're expecting
	var results []result

	// start listening for any results over the resultsChan
	// once we get a result append it to the result slice
	for {
		result := <-resultsChan
		results = append(results, *result)

		// if we've reached the expected amount of urls then stop
		if len(results) == len(txs) {
			break
		}
	}

	// let's sort these results real quick
	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	for _, res := range results {
		receipts = append(receipts, res.res)
	}
	// now we're done we return the results
	return receipts, nil
}

func (ht *BlockFetcher) getResultsBatchParallel(ctx context.Context, txs []*types.Transaction) (receipts []*bulletprooftxmanager.Receipt, err error) {

	batchLimit := 100
	concurrencyLimit := 10
	// this buffered channel will block at the concurrency limit
	semaphoreChan := make(chan struct{}, concurrencyLimit)

	// this channel will not block and collect the http request results
	resultsChan := make(chan *resultReceipts)

	// make sure we close these channels when we're done with them
	defer func() {
		close(semaphoreChan)
		close(resultsChan)
	}()

	var slices [][]*types.Transaction
	for i := 0; i < len(txs); i += batchLimit {
		end := i + batchLimit
		if end > len(txs) {
			end = len(txs)
		}

		slices = append(slices, txs[i:end])
	}

	// keen an index and loop through every url we will send a request to
	for i, txs := range slices {
		txsInner := txs
		// start a go routine with the index and url in a closure
		go func(i int, txsInner []*types.Transaction) {

			// this sends an empty struct into the semaphoreChan which
			// is basically saying add one to the limit, but when the
			// limit has been reached block until there is room
			semaphoreChan <- struct{}{}
			start := time.Now()
			// send the request and put the response in a result struct
			// along with the index so we can sort them later along with
			// any error that might have occoured

			receipts, err := ht.getResultsBatch(ctx, txsInner)

			elapsed := time.Since(start)
			ht.logger.Debugw("========================= HeadTracker: getting a batch of receipts took", "len", len(txsInner),
				"elapsedMs", elapsed.Milliseconds())

			//if err != nil {
			//	return nil, err
			//}

			result := &resultReceipts{i, receipts, err}

			// now we can send the result struct through the resultsChan
			resultsChan <- result

			// once we're done it's we read from the semaphoreChan which
			// has the effect of removing one from the limit and allowing
			// another goroutine to start
			<-semaphoreChan

		}(i, txsInner)
	}

	// make a slice to hold the results we're expecting
	var results []resultReceipts

	// start listening for any results over the resultsChan
	// once we get a result append it to the result slice
	for {
		result := <-resultsChan
		results = append(results, *result)

		// if we've reached the expected amount of urls then stop
		if len(results) == len(slices) {
			break
		}
	}

	// let's sort these results real quick
	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	for _, res := range results {
		receipts = append(receipts, res.res...)
	}
	// now we're done we return the results
	return receipts, nil
}

func (ht *BlockFetcher) getResultsBatch(ctx context.Context, txs []*types.Transaction) (receipts []*bulletprooftxmanager.Receipt, err error) {
	var reqs []rpc.BatchElem
	for _, tx := range txs { // TODO: how many is too many?
		req := rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{tx.Hash()},
			Result: &bulletprooftxmanager.Receipt{},
		}
		//ht.logger.Debugf("========================= HeadTracker: getting receipt for transaction %v", tx.To())

		reqs = append(reqs, req)
	}

	ctx, cancel := eth.DefaultQueryCtx(ctx)
	defer cancel()

	err = ht.store.EthClient.BatchCallContext(ctx, reqs)
	if err != nil {
		return nil, errors.Wrap(err, "EthConfirmer#batchFetchReceipts error fetching receipts with BatchCallContext")
	}

	for _, req := range reqs {
		result, err := req.Result, req.Error
		if err != nil {
			ht.logger.Errorw("EthConfirmer#batchFetchReceipts: fetchReceipt failed")
			return nil, err
		}
		receipt, is := result.(*bulletprooftxmanager.Receipt)
		if !is {
			return nil, errors.Errorf("expected result to be a %T, got %T", (*bulletprooftxmanager.Receipt)(nil), result)
		}

		receipts = append(receipts, receipt)
	}
	return receipts, nil
}

func (ht *BlockFetcher) batchFetchReceipts(ctx context.Context, txs []*types.Transaction) (receipts []*bulletprooftxmanager.Receipt, err error) {
	if len(txs) == 0 {
		return
	}
	receiptsRaw, err := ht.getResultsBatchParallel(ctx, txs)

	for i, receipt := range receiptsRaw {
		tx := txs[i]

		l := ht.logger.With(
			"txHash", tx.Hash().Hex(), "err", err,
		)

		if err != nil {
			l.Errorw("EthConfirmer#batchFetchReceipts: fetchReceipt failed")
			continue
		}

		if receipt == nil {
			// NOTE: This should never possibly happen, but it seems safer to
			// check regardless to avoid a potential panic
			l.Errorw("EthConfirmer#batchFetchReceipts: invariant violation, got nil receipt")
			continue
		}

		if receipt.IsZero() {
			l.Debugw("EthConfirmer#batchFetchReceipts: still waiting for receipt")
			continue
		}

		l = l.With("receipt", receipt)

		if receipt.IsUnmined() {
			l.Debugw("EthConfirmer#batchFetchReceipts: got receipt for transaction but it's still in the mempool and not included in a block yet")
			continue
		}

		//l.Debugw("EthConfirmer#batchFetchReceipts: got receipt for transaction", "blockNumber", receipt.BlockNumber)

		if receipt.TxHash != tx.Hash() {
			l.Errorf("EthConfirmer#batchFetchReceipts: invariant violation, expected receipt with hash %s to have same hash as attempt with hash %s", receipt.TxHash.Hex(), tx.Hash().Hex())
			continue
		}

		if receipt.BlockNumber == nil {
			l.Error("EthConfirmer#batchFetchReceipts: invariant violation, receipt was missing block number")
			continue
		}

		//if receipt.Status == 0 {
		//	l.Warnf("transaction %s reverted on-chain", receipt.TxHash)
		//	// This is safe to increment here because we save the receipt immediately after
		//	// and once its saved we do not fetch it again.
		//	//promRevertedTxCount.Add(1)
		//}

		receipts = append(receipts, receipt)
	}

	return
}
