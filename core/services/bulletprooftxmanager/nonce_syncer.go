package bulletprooftxmanager

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
	"go.uber.org/multierr"
	"gorm.io/gorm"
)

type (
	// NonceSyncer manages the delicate task of syncing the local nonce with the
	// chain nonce in case of divergence.
	//
	// On startup, we check each key for the nonce value on chain and compare
	// it to our local value.
	//
	// Usually the on-chain nonce will be the same as (or lower than) the
	// next_nonce in the DB, in which case we do nothing.
	//
	// If we are restoring from a backup however, or another wallet has used the
	// account, the chain nonce might be higher than our local one. In this
	// scenario, we must fastforward the local nonce to match the chain nonce.
	//
	// The problem with doing this is that now Chainlink does not have any
	// ownership or control over potentially pending transactions with nonces
	// between our local highest nonce and the chain nonce. If one of those
	// transactions is pushed out of the mempool or re-org'd out of the chain,
	// we run the risk of being stuck with a gap in the nonce sequence that
	// will never be filled.
	//
	// The solution is to query the chain for our own transactions and take
	// ownership of them by writing them to the database and letting the
	// EthConfirmer handle them as it would any other transaction.
	//
	// This is not quite as straightforward as one might expect. We cannot
	// query transactions from our account to infinite depth (geth does not
	// support this). The best we can do is to query for all transactions sent
	// within the past ETH_FINALITY_DEPTH blocks and find the ones sent by our
	// address(es).
	//
	// This gives us re-org protection up to ETH_FINALITY_DEPTH deep in the
	// worst case, which is in line with our other guarantees.
	NonceSyncer struct {
		store     *store.Store
		config    orm.ConfigReader
		ethClient eth.Client
	}
	// NSinserttx represents an EthTx and Attempt to be inserted together
	NSinserttx struct {
		Etx     models.EthTx
		Attempt models.EthTxAttempt
	}
)

// NewNonceSyncer returns a new syncer
func NewNonceSyncer(store *store.Store, config orm.ConfigReader, ethClient eth.Client) *NonceSyncer {
	return &NonceSyncer{
		store,
		config,
		ethClient,
	}
}

// SyncAll syncs nonces for all keys in parallel
//
// This should only be called once, before the EthBroadcaster has started.
// Calling it later is not safe and could lead to races.
func (s NonceSyncer) SyncAll(ctx context.Context) (merr error) {
	keys, err := s.store.SendKeys()
	if err != nil {
		return errors.Wrap(err, "NonceSyncer#fastForwardNoncesIfNecessary failed to get keys")
	}

	var wg sync.WaitGroup
	var errMu sync.Mutex

	wg.Add(len(keys))
	for _, key := range keys {
		go func(k models.Key) {
			defer wg.Done()
			if err := s.fastForwardNonceIfNecessary(ctx, k.Address.Address()); err != nil {
				errMu.Lock()
				defer errMu.Unlock()
				merr = multierr.Combine(merr, err)
			}
		}(key)
	}

	wg.Wait()

	return errors.Wrap(merr, "NonceSyncer#fastForwardNoncesIfNecessary failed")
}

func (s NonceSyncer) fastForwardNonceIfNecessary(ctx context.Context, address common.Address) error {
	chainNonce, err := s.pendingNonceFromEthClient(ctx, address)
	if err != nil {
		return errors.Wrap(err, "GetNextNonce failed to loadInitialNonceFromEthClient")
	}
	if chainNonce == 0 {
		return nil
	}

	localNonce, err := GetNextNonce(s.store.DB, address)
	if err != nil {
		return err
	}
	if chainNonce <= uint64(localNonce) {
		return nil
	}
	logger.Warnw(fmt.Sprintf("NonceSyncer: address %s has been used before, either by an external wallet or a different Chainlink node. "+
		"Local nonce is %v but the on-chain nonce for this account was %v. "+
		"Please note that using the chainlink keys with an external wallet is NOT SUPPORTED and can lead to missed or stuck transactions. "+
		"This Chainlink node will now take ownership of this address and may overwrite currently pending transactions",
		address.Hex(), localNonce, chainNonce),
		"address", address.Hex(), "localNonce", localNonce, "chainNonce", chainNonce)

	// First fetch
	//
	// We have to get the latest block first and then fetch deeper blocks in a
	// subsequent step, because otherwise we don't know what block numbers to
	// query for.
	// pending can be fetched in the initial request because it doesn't rely on
	// any block number.
	reqs := []rpc.BatchElem{
		rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{"pending", true},
			Result: &models.Block{},
		},
		rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{"latest", true},
			Result: &models.Block{},
		},
	}

	err = s.ethClient.BatchCallContext(ctx, reqs)
	if err != nil {
		return err
	}

	if reqs[1].Error != nil {
		return errors.Wrap(reqs[1].Error, "latest block request returned error")
	}
	latestBlock, is := reqs[1].Result.(*models.Block)
	if !is {
		panic(fmt.Sprintf("invariant violation, expected %T but got %T", &models.Block{}, latestBlock))
	}
	latestBlockNum := latestBlock.Number

	floor := latestBlockNum - int64(s.config.EthFinalityDepth())
	if floor < 0 {
		floor = 0
	}

	// Second fetch
	//
	// OPTIMISATION NOTE:
	// The astute observer will note that if multiple keys are behind, we fetch
	// the same blocks multiple times, doing redundant extra work. This does
	// put unnecessary load on the eth node, but the most common use-case is
	// with a single key and given how rare this scenario is, I think we can
	// live with it.
	var reqs2 []rpc.BatchElem
	for i := floor; i < latestBlockNum; i++ {
		req := rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{models.Int64ToHex(i), true},
			Result: &models.Block{},
		}
		reqs2 = append(reqs2, req)
	}

	err = s.ethClient.BatchCallContext(ctx, reqs2)
	if err != nil {
		return err
	}

	reqs = append(reqs, reqs2...)
	var txes []types.Transaction
	signer := types.NewEIP155Signer(s.config.ChainID())

	// Rip through all transactions in all blocks and keep only the ones sent
	// from our key
	for _, req := range reqs {
		if req.Error != nil {
			logger.Warnw("NonceSyncer: got error querying for block", "blockNum", req.Args[0], "err", req.Error)
			continue
		}
		block, is := req.Result.(*models.Block)
		if !is {
			panic(fmt.Sprintf("invariant violation, expected %T but got %T", &models.Block{}, block))
		}
		for _, tx := range block.Transactions {
			from, err2 := types.Sender(signer, &tx)
			if err2 != nil {
				logger.Warnw("NonceSyncer#fastForwardNonceIfNecessary failed to extract 'from' from transaction", "tx", tx, "err", err2)
				continue
			}
			if from == address {
				txes = append(txes, tx)
			}
		}
	}

	account, err := s.store.KeyStore.GetAccountByAddress(address)
	if err != nil {
		return errors.Wrap(err, "NonceSyncer#fastForwardNonceIfNecessary could not get account from keystore")
	}
	sort.Slice(txes, func(i, j int) bool { return txes[i].Nonce() < txes[j].Nonce() })

	inserts, err := s.makeInserts(account, latestBlock.Number, txes, chainNonce-1)
	if err != nil {
		return errors.Wrap(err, "NonceSyncer#fastForwardNonceIfNecessary error generating transactions for backfill")
	}

	now := time.Now()
	return postgres.GormTransaction(ctx, s.store.DB, func(dbtx *gorm.DB) error {
		//  We pass in next_nonce here as an optimistic lock to make sure it
		//  didn't get changed out from under us
		res := dbtx.Exec(`UPDATE keys SET next_nonce = ?, updated_at = ? WHERE address = ? AND next_nonce = ?`, chainNonce, now, address, localNonce)
		if res.Error != nil {
			return errors.Wrap(res.Error, "NonceSyncer#fastForwardNonceIfNecessary failed to update keys.next_nonce")
		}
		if res.RowsAffected == 0 {
			return errors.Errorf("NonceSyncer#fastForwardNonceIfNecessary optimistic lock failure fastforwarding nonce %v to %v for key %s", localNonce, chainNonce, address.Hex())
		}

		for _, ins := range inserts {
			// Setting broadcast_at here is a bit of a misnomer since this node
			// didn't actually broadcast the transaction, but including it
			// allows us to avoid changing the state machine limitations and
			// represents roughly the time we read the tx from the blockchain
			// so we can pretty much assume it was "broadcast" at this time.
			ins.Etx.BroadcastAt = &now
			if err := dbtx.Create(&ins.Etx).Error; err != nil {
				return errors.Wrap(err, "NonceSyncer#fastForwardNonceIfNecessary failed to create eth_tx")
			}
			ins.Attempt.EthTxID = ins.Etx.ID
			if err := dbtx.Create(&ins.Attempt).Error; err != nil {
				return errors.Wrap(err, "NonceSyncer#fastForwardNonceIfNecessary failed to create eth_tx_attempt")
			}
		}
		return nil
	})
}

func (s NonceSyncer) pendingNonceFromEthClient(ctx context.Context, account common.Address) (nextNonce uint64, err error) {
	ctx, cancel := context.WithTimeout(ctx, maxEthNodeRequestTime)
	defer cancel()
	nextNonce, err = s.ethClient.PendingNonceAt(ctx, account)
	return nextNonce, errors.WithStack(err)
}

func (s NonceSyncer) makeInserts(acct accounts.Account, blockNum int64, txes []types.Transaction, toNonce uint64) (inserts []NSinserttx, err error) {
	if len(txes) == 0 {
		return
	}
	if !sort.SliceIsSorted(txes, func(i, j int) bool { return txes[i].Nonce() < txes[j].Nonce() }) {
		return nil, errors.New("expected txes to be sorted in nonce ascending order")
	}
	fromNonce := txes[0].Nonce()
	if fromNonce > toNonce {
		// I don't know how this could ever happen but we should handle the case anyway
		return nil, errors.Errorf("fromNonce of %v was greater than toNonce of %v", fromNonce, toNonce)
	}
	txMap := make(map[uint64]types.Transaction)
	for _, tx := range txes {
		txMap[tx.Nonce()] = tx
	}
	for n := fromNonce; n <= toNonce; n++ {
		nonce := int64(n)
		tx, exists := txMap[n]
		if exists {
			ins, err := s.MakeInsert(tx, acct, blockNum, nonce)
			if err != nil {
				logger.Errorw("NonceSyncer: failed to generate transaction, this nonce will not be re-org protected", "address", acct.Address.Hex(), "err", err, "nonce", nonce)
				continue
			}
			inserts = append(inserts, ins)
		} else {
			// Use a zero-transaction if its missing for whatever reason
			// Should really never happen but you never know with geth
			logger.Warnw(fmt.Sprintf("NonceSyncer: missing transaction for nonce %d, falling back to zero transaction", n), "address", acct.Address.Hex(), "blockNum", blockNum, "nonce", nonce)
			ins, err := s.MakeZeroInsert(acct, blockNum, nonce)
			if err != nil {
				logger.Errorw("NonceSyncer: failed to generate empty transaction, this nonce will not be re-org protected", "address", acct.Address.Hex(), "err", err, "nonce", nonce)
				continue
			}
			inserts = append(inserts, ins)
		}
	}

	return inserts, nil
}

// MakeInsert generates a NSinserttx that perfectly mirrors the on-chain transaction
//
// This can be handed off to the EthConfirmer and used to query for receipts
// and bump gas etc exactly like any other transaction we might have sent.
func (s NonceSyncer) MakeInsert(tx types.Transaction, acct accounts.Account, blockNum, nonce int64) (ins NSinserttx, err error) {
	v, _, _ := tx.RawSignatureValues()
	if v.BitLen() == 0 {
		// Believe it or not, this is the only way to determine if the tx
		// is a zero struct without panicking. Thank you, geth.
		logger.Warnw("NonceSyncer: tx was empty/unsigned. Falling back to zero transaction", "err", err, "txHash", tx.Hash(), "nonce", nonce, "address", acct.Address.Hex())
		return s.MakeZeroInsert(acct, blockNum, int64(nonce))
	}
	// NOTE: We set all transactions to unconfirmed even if they are mined.
	//
	// This works out, because the first round of the EthConfirmer will check
	// for receipts. All these transactions should get receipts if they are confirmed.
	//
	// Any transaction not yet confirmed will go into the regular gas bumping
	// cycle as if it were any normal transaction we had sent ourselves.
	ins.Etx = models.EthTx{
		Nonce:          &nonce,
		FromAddress:    acct.Address,
		ToAddress:      *tx.To(),
		EncodedPayload: tx.Data(),
		Value:          assets.Eth(*tx.Value()),
		GasLimit:       tx.Gas(),
		State:          models.EthTxUnconfirmed,
	}
	rlp := new(bytes.Buffer)
	if err := tx.EncodeRLP(rlp); err != nil {
		logger.Warnw("NonceSyncer: could not encode RLP. Falling back to zero transaction", "err", err, "txHash", tx.Hash(), "nonce", nonce, "address", acct.Address.Hex())
		return s.MakeZeroInsert(acct, blockNum, int64(nonce))
	}
	signedRawTx := rlp.Bytes()
	ins.Attempt = models.EthTxAttempt{
		GasPrice:                utils.Big(*tx.GasPrice()),
		SignedRawTx:             signedRawTx,
		Hash:                    tx.Hash(),
		BroadcastBeforeBlockNum: &blockNum,
		State:                   models.EthTxAttemptBroadcast,
	}

	return ins, nil
}

// MakeZeroInsert generates a NSinserttx that represents a zero transaction
//
// This transaction will never get a receipt and does not match anything
// on-chain, but it does serve the purpose of a placeholder for this nonce in
// case the on-chain version is ejected from the mempool or re-org'd out of the
// main chain.
func (s NonceSyncer) MakeZeroInsert(acct accounts.Account, blockNum, nonce int64) (ins NSinserttx, err error) {
	gasLimit := s.config.EthGasLimitDefault()
	gasPrice := s.config.EthGasPriceDefault()

	tx, err := makeEmptyTransaction(s.store.KeyStore, uint64(nonce), gasLimit, gasPrice, acct, s.config.ChainID())
	if err != nil {
		return ins, errors.Wrap(err, "NonceSyncer#MakeZeroInsert failed to makeEmptyTransaction")
	}
	rlp := new(bytes.Buffer)
	if err := tx.EncodeRLP(rlp); err != nil {
		return ins, err
	}
	// NOTE: These transactions will never get a receipt, but setting them to
	// unconfirmed still works out.
	//
	// If there is a transaction on-chain with the same nonce that is pending,
	// then this zero transaction will go into the bumping cycle and may or may
	// not replace the on-chain version.
	//
	// If the on-chain transaction is confirmed, this one will eventually be
	// marked confirmed_missing_receipt when a new transaction is confirmed on
	// top of it, and will exit the bumping cycle once it's deeper than
	// ETH_FINALITY_DEPTH.
	ins.Etx = models.EthTx{
		Nonce:          &nonce,
		FromAddress:    acct.Address,
		ToAddress:      *tx.To(),
		EncodedPayload: tx.Data(),
		Value:          assets.Eth(*tx.Value()),
		GasLimit:       tx.Gas(),
		State:          models.EthTxUnconfirmed,
	}
	ins.Attempt = models.EthTxAttempt{
		GasPrice:                utils.Big(*gasPrice),
		SignedRawTx:             rlp.Bytes(),
		Hash:                    tx.Hash(),
		BroadcastBeforeBlockNum: &blockNum,
		State:                   models.EthTxAttemptBroadcast,
	}
	return ins, nil
}
