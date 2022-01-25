package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/blockhash_store"
)

// BATCH_SIZE is the size of the batch of storeVerifyHeader transactions to send on-chain
// at a time.
const BATCH_SIZE uint64 = 5

type backwardFeeder struct {
	account      *bind.TransactOpts
	ethClient    *ethclient.Client
	chainID      *big.Int
	bhs          *blockhash_store.BlockhashStore
	bhsABI       *abi.ABI
	gasEstimator *gasEstimator
}

// storeEarliest creates a BlockhashStore.storeEarliest transaction
// and sends it to the eth node.
func (f *backwardFeeder) storeEarliest() (*types.Transaction, error) {
	nonce, err := f.ethClient.PendingNonceAt(context.Background(), f.account.From)
	if err != nil {
		return nil, errors.Wrap(err, "pending nonce at")
	}

	payload, err := f.bhsABI.Pack("storeEarliest")
	if err != nil {
		return nil, errors.Wrap(err, "abi pack storeEarliest")
	}

	gasPrice, err := f.gasEstimator.estimate()
	if err != nil {
		return nil, errors.Wrap(err, "estimating gas price")
	}

	fmt.Println("gas price:", gasPrice)

	toAddress := f.bhs.Address()
	tx, err := f.account.Signer(f.account.From, types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      100_000,
		To:       &toAddress,
		Data:     payload,
	}))
	if err != nil {
		return nil, errors.Wrap(err, "signing tx")
	}

	if err := f.ethClient.SendTransaction(context.Background(), tx); err != nil {
		return nil, errors.Wrap(err, "sending tx")
	}

	fmt.Println("Sent storeEarliest tx:", tx.Hash())

	return tx, nil
}

// store creates a BlockhashStore.store transaction and sends it to the eth node.
func (f *backwardFeeder) store(blockNum *big.Int) (*types.Transaction, error) {
	nonce, err := f.ethClient.PendingNonceAt(context.Background(), f.account.From)
	if err != nil {
		return nil, errors.Wrap(err, "pending nonce at")
	}

	payload, err := f.bhsABI.Pack("store", blockNum)
	if err != nil {
		return nil, errors.Wrapf(err, "abi pack store %s", blockNum.String())
	}

	gasPrice, err := f.gasEstimator.estimate()
	if err != nil {
		return nil, errors.Wrap(err, "estimating gas price")
	}

	toAddress := f.bhs.Address()
	tx, err := f.account.Signer(f.account.From, types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      100_000,
		To:       &toAddress,
		Data:     payload,
	}))
	if err != nil {
		return nil, errors.Wrap(err, "signing tx")
	}

	if err := f.ethClient.SendTransaction(context.Background(), tx); err != nil {
		return nil, errors.Wrap(err, "sending tx")
	}

	fmt.Println("Sent store(", blockNum.String(), ") tx:", tx.Hash())

	return tx, nil
}

// createStoreVerifyHeaderTxBatch creates a batch of BlockhashStore.storeVerifyHeader transactions
// and decrements the given feedFrom integer appropriately.
func (f *backwardFeeder) createStoreVerifyHeaderTxBatch(feedFrom *big.Int) (txs []*types.Transaction, err error) {
	nonce, err := f.ethClient.PendingNonceAt(context.Background(), f.account.From)
	if err != nil {
		return nil, err
	}

	for i := uint64(0); i < BATCH_SIZE; i++ {
		blockNumToStore := new(big.Int).Set(feedFrom).Sub(feedFrom, big.NewInt(1))
		// Call will revert if there is no blockhash for the given block
		blockhash, err := f.bhs.GetBlockhash(nil, blockNumToStore)
		if err != nil {
			tx, err := f.createStoreVerifyHeaderTx(nonce, blockNumToStore, feedFrom)
			if err != nil {
				return nil, err
			}
			fmt.Println("Created (but did not send) storeVerifyHeader(", blockNumToStore.String(), ", blockHeader) tx:", tx.Hash())
			txs = append(txs, tx)
		} else {
			fmt.Println("Blockhash", hex.EncodeToString(blockhash[:]), "already stored for block", blockNumToStore)
		}
		feedFrom.Sub(feedFrom, big.NewInt(1))
		nonce++
	}

	fmt.Println("Created storeVerifyHeader tx batch")

	return txs, nil
}

// createStoreVerifyHeaderTx creates a BlockhashStore.storeVerifyHeader transaction from
// the given toStore and after block numbers.
func (f *backwardFeeder) createStoreVerifyHeaderTx(nonce uint64, toStore, after *big.Int) (*types.Transaction, error) {
	blockHeader, err := serializedBlockHeader(f.ethClient, after, f.chainID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to serialize blockheader:", err)
		return nil, err
	}

	payload, err := f.bhsABI.Pack("storeVerifyHeader", toStore, blockHeader)
	if err != nil {
		return nil, errors.Wrap(err, "packing storeVerifyHeader")
	}

	gasPrice, err := f.gasEstimator.estimate()
	if err != nil {
		return nil, errors.Wrap(err, "estimating gas price")
	}

	toAddress := f.bhs.Address()
	tx, err := f.account.Signer(f.account.From, types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      100_000,
		To:       &toAddress,
		Data:     payload,
	}))
	if err != nil {
		return nil, errors.Wrap(err, "signing tx")
	}

	return tx, nil
}

// sendTxBatch sends the given batch of transactions to the eth node and waits for them to be mined.
func (f *backwardFeeder) sendTxBatch(txs []*types.Transaction) error {
	fmt.Println("Sending batch of", len(txs), "transactions")
	for i, tx := range txs {
		if err := f.ethClient.SendTransaction(context.Background(), tx); err != nil {
			return errors.Wrapf(err, "sending transaction %+v iteration %d", tx, i)
		}
	}

	fmt.Println("Waiting for receipt on last of", len(txs), "transactions")

	// Check that the txs were mined on-chain.
	// Check only the last one for simplicity.
	receipt, err := bind.WaitMined(context.Background(), f.ethClient, txs[len(txs)-1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error waiting for tx to mine:", err)
		return err
	}

	fmt.Println("Got receipt for last tx in batch. Blocknumber:", receipt.BlockNumber)

	return nil
}

func (f *backwardFeeder) feed() error {
	// store earliest first
	_, err := f.storeEarliest()
	if err != nil {
		return errors.Wrap(err, "create and send store earliest tx")
	}

	blockNumber, err := f.ethClient.BlockNumber(context.Background())
	if err != nil {
		return errors.Wrap(err, "get latest block number")
	}

	// store current blocknumber - 1
	txNow, err := f.store(big.NewInt(int64(blockNumber) - 1))
	if err != nil {
		return errors.Wrap(err, "create and send store tx")
	}

	receiptNow, err := bind.WaitMined(context.Background(), f.ethClient, txNow)
	if err != nil {
		return errors.Wrap(err, "wait mined")
	}

	feedFrom := new(big.Int).Set(receiptNow.BlockNumber).
		Sub(receiptNow.BlockNumber, big.NewInt(256))
	fmt.Println("Sleeping for 20 seconds then feeding blockhashes from block", feedFrom, "backwards")

	time.Sleep(20 * time.Second)

	for {
		txs, err := f.createStoreVerifyHeaderTxBatch(feedFrom)
		if err != nil {
			return err
		}

		if err := f.sendTxBatch(txs); err != nil {
			return errors.Wrap(err, "sending tx batch")
		}
	}
}
