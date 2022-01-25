package main

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/blockhash_store"
)

const LOOKBACK_BLOCKS uint64 = 5

type forwardFeeder struct {
	account      *bind.TransactOpts
	ethClient    *ethclient.Client
	chainID      *big.Int
	bhs          *blockhash_store.BlockhashStore
	bhsABI       *abi.ABI
	gasEstimator *gasEstimator
	eip1559      bool
}

func (f *forwardFeeder) feed() error {
	latestNumber, err := f.ethClient.BlockNumber(context.Background())
	if err != nil {
		return errors.Wrap(err, "get latest block number")
	}

	fmt.Println("Latest block number:", latestNumber)

	blockNumber := latestNumber - LOOKBACK_BLOCKS
	fmt.Println("feeding forward, starting at block", blockNumber)
	for {
		latestNumber, err := f.ethClient.BlockNumber(context.Background())
		if err != nil {
			return errors.Wrap(err, "get latest block number")
		}

		batchSize := int(math.Min(float64(BATCH_SIZE), float64(latestNumber)-float64(blockNumber)))
		if batchSize == 0 {
			fmt.Println("caught up with chain tip, going to sleep for 10s")
			time.Sleep(10 * time.Second)
			continue
		}

		txs, nextBlockNumber, err := f.createStoreTxBatch(blockNumber, batchSize)
		if err != nil {
			return errors.Wrap(err, "create store tx batch")
		}

		fmt.Println("submitting batch:", hashes(txs))

		err = f.submitTxBatch(txs)
		if err != nil {
			return errors.Wrap(err, "submit tx batch")
		}

		_, err = bind.WaitMined(context.Background(), f.ethClient, txs[len(txs)-1])
		if err != nil {
			return errors.Wrap(err, "wait mined")
		}

		fmt.Println("Got receipt for batch")

		blockNumber = nextBlockNumber
	}
}

func (f *forwardFeeder) submitTxBatch(txs []*types.Transaction) error {
	for i, tx := range txs {
		if err := f.ethClient.SendTransaction(context.Background(), tx); err != nil {
			fmt.Println("error sending tx, iteration", i, ":", err, "tx:", tx)
			return err
		}
	}
	return nil
}

func (f *forwardFeeder) createStoreTxBatch(startBlockNumber uint64, batchSize int) (txs []*types.Transaction, nextBlockNumber uint64, err error) {
	nonce, err := f.ethClient.PendingNonceAt(context.Background(), f.account.From)
	if err != nil {
		return nil, startBlockNumber, err
	}

	txs = []*types.Transaction{}
	for i := 0; i < batchSize; i++ {
		tx, err := f.createStoreTx(nonce, big.NewInt(int64(startBlockNumber)))
		if err != nil {
			return nil, startBlockNumber, err
		}

		txs = append(txs, tx)
		nonce++
		startBlockNumber++
	}

	return txs, startBlockNumber, nil
}

func (f *forwardFeeder) createStoreTx(nonce uint64, blockNumber *big.Int) (*types.Transaction, error) {
	payload, err := f.bhsABI.Pack("store", blockNumber)
	if err != nil {
		return nil, err
	}

	gasPrice, err := f.gasEstimator.estimate()
	if err != nil {
		return nil, err
	}

	toAddress := f.bhs.Address()
	var txData types.TxData
	if f.eip1559 {
		txData = &types.DynamicFeeTx{
			ChainID:   f.chainID,
			Nonce:     nonce,
			To:        &toAddress,
			Data:      payload,
			Gas:       100_000,
			GasFeeCap: gasPrice,
		}
	} else {
		txData = &types.LegacyTx{
			Nonce:    nonce,
			GasPrice: gasPrice,
			Gas:      100_000,
			To:       &toAddress,
			Data:     payload,
		}
	}
	tx := types.NewTx(txData)
	signedTx, err := f.account.Signer(f.account.From, tx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}
