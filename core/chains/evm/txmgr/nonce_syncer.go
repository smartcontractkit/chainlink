package txmgr

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type (
	NonceSyncerKeyStore interface {
		GetNextNonce(address common.Address, chainID *big.Int, qopts ...pg.QOpt) (int64, error)
	}
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
		q         pg.Q
		ethClient evmclient.Client
		chainID   *big.Int
		logger    logger.Logger
		kst       NonceSyncerKeyStore
	}
	// NSinserttx represents an EthTx and Attempt to be inserted together
	NSinserttx struct {
		Etx     EthTx
		Attempt EthTxAttempt
	}
)

// NewNonceSyncer returns a new syncer
func NewNonceSyncer(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig, ethClient evmclient.Client, kst NonceSyncerKeyStore) *NonceSyncer {
	lggr = lggr.Named("NonceSyncer")
	q := pg.NewQ(db, lggr, cfg)
	return &NonceSyncer{
		q,
		ethClient,
		ethClient.ChainID(),
		lggr,
		kst,
	}
}

// SyncAll syncs nonces for all enabled keys in parallel
//
// This should only be called once, before the EthBroadcaster has started.
// Calling it later is not safe and could lead to races.
func (s NonceSyncer) Sync(ctx context.Context, keyState ethkey.State) (err error) {
	if keyState.Disabled {
		return errors.Errorf("cannot sync disabled key state: %s", keyState.Address)
	}
	err = s.fastForwardNonceIfNecessary(ctx, keyState.Address.Address())
	return errors.Wrap(err, "NonceSyncer#fastForwardNoncesIfNecessary failed")
}

func (s NonceSyncer) fastForwardNonceIfNecessary(ctx context.Context, address common.Address) error {
	chainNonce, err := s.pendingNonceFromEthClient(ctx, address)
	if err != nil {
		return errors.Wrap(err, "GetNextNonce failed to loadInitialNonceFromEthClient")
	}
	if chainNonce == 0 {
		return nil
	}

	keyNextNonce, err := s.kst.GetNextNonce(address, s.chainID, pg.WithParentCtx(ctx))
	if err != nil {
		return err
	}

	q := s.q.WithOpts(pg.WithParentCtx(ctx))

	localNonce := keyNextNonce
	hasInProgressTransaction, err := s.hasInProgressTransaction(q, address)
	if err != nil {
		return errors.Wrapf(err, "failed to query for in_progress transaction for address %s", address.Hex())
	} else if hasInProgressTransaction {
		// If we have an 'in_progress' transacion, our keys.next_nonce will be
		// one lower than it should because we must have crashed mid-execution.
		// The EthBroadcaster will automatically take care of this and
		// increment it by one later, for now we just increment by one here.
		localNonce++
	}
	if chainNonce <= uint64(localNonce) {
		return nil
	}
	s.logger.Warnw(fmt.Sprintf("address %s has been used before, either by an external wallet or a different Chainlink node. "+
		"Local nonce is %v but the on-chain nonce for this account was %v. "+
		"It's possible that this node was restored from a backup. If so, transactions sent by the previous node will NOT be re-org protected and in rare cases may need to be manually bumped/resubmitted. "+
		"Please note that using the chainlink keys with an external wallet is NOT SUPPORTED and can lead to missed or stuck transactions. ",
		address.Hex(), localNonce, chainNonce),
		"address", address.Hex(), "keyNextNonce", keyNextNonce, "localNonce", localNonce, "chainNonce", chainNonce)

	// Need to remember to decrement the chain nonce by one to account for in_progress transaction
	newNextNonce := chainNonce
	if hasInProgressTransaction {
		newNextNonce--
	}
	//  We pass in next_nonce here as an optimistic lock to make sure it
	//  didn't get changed out from under us. Shouldn't happen but can't hurt.
	err = q.Transaction(func(tx pg.Queryer) error {
		res, err := tx.Exec(`UPDATE evm_key_states SET next_nonce = $1, updated_at = $2 WHERE address = $3 AND next_nonce = $4 AND evm_chain_id = $5`, newNextNonce, time.Now(), address, keyNextNonce, s.chainID.String())
		if err != nil {
			return errors.Wrap(err, "NonceSyncer#fastForwardNonceIfNecessary failed to update keys.next_nonce")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return errors.Wrap(err, "NonceSyncer#fastForwardNonceIfNecessary failed to get RowsAffected")
		}
		if rowsAffected == 0 {
			return errors.Errorf("NonceSyncer#fastForwardNonceIfNecessary optimistic lock failure fastforwarding nonce %v to %v for key %s", localNonce, chainNonce, address.Hex())
		}
		return nil
	})
	if err == nil {
		s.logger.Infow("Fast-forwarded nonce", "address", address, "newNextNonce", newNextNonce, "oldNextNonce", keyNextNonce)
	}
	return err
}

func (s NonceSyncer) pendingNonceFromEthClient(ctx context.Context, account common.Address) (nextNonce uint64, err error) {
	nextNonce, err = s.ethClient.PendingNonceAt(ctx, account)
	return nextNonce, errors.WithStack(err)
}

func (s NonceSyncer) hasInProgressTransaction(q pg.Queryer, account common.Address) (exists bool, err error) {
	err = q.Get(&exists, `SELECT EXISTS(SELECT 1 FROM eth_txes WHERE state = 'in_progress' AND from_address = $1 AND evm_chain_id = $2)`, account, s.chainID.String())
	return exists, errors.Wrap(err, "hasInProgressTransaction failed")
}
