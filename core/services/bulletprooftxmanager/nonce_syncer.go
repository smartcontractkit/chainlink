package bulletprooftxmanager

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
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
		db        *gorm.DB
		ethClient eth.Client
	}
	// NSinserttx represents an EthTx and Attempt to be inserted together
	NSinserttx struct {
		Etx     models.EthTx
		Attempt models.EthTxAttempt
	}
)

// NewNonceSyncer returns a new syncer
func NewNonceSyncer(db *gorm.DB, ethClient eth.Client) *NonceSyncer {
	return &NonceSyncer{
		db,
		ethClient,
	}
}

// SyncAll syncs nonces for all keys in parallel
//
// This should only be called once, before the EthBroadcaster has started.
// Calling it later is not safe and could lead to races.
func (s NonceSyncer) SyncAll(ctx context.Context, keys []ethkey.Key) (merr error) {
	var wg sync.WaitGroup
	var errMu sync.Mutex

	wg.Add(len(keys))
	for _, key := range keys {
		go func(k ethkey.Key) {
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

	selectCtx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	keyNextNonce, err := GetNextNonce(s.db.WithContext(selectCtx), address)
	if err != nil {
		return err
	}

	localNonce := keyNextNonce
	hasInProgressTransaction, err := s.hasInProgressTransaction(address)
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
	logger.Warnw(fmt.Sprintf("NonceSyncer: address %s has been used before, either by an external wallet or a different Chainlink node. "+
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
	return postgres.DBWithDefaultContext(s.db, func(db *gorm.DB) error {
		res := db.Exec(`UPDATE keys SET next_nonce = ?, updated_at = ? WHERE address = ? AND next_nonce = ?`, newNextNonce, time.Now(), address, keyNextNonce)
		if res.Error != nil {
			return errors.Wrap(res.Error, "NonceSyncer#fastForwardNonceIfNecessary failed to update keys.next_nonce")
		}
		if res.RowsAffected == 0 {
			return errors.Errorf("NonceSyncer#fastForwardNonceIfNecessary optimistic lock failure fastforwarding nonce %v to %v for key %s", localNonce, chainNonce, address.Hex())
		}
		return nil
	})
}

func (s NonceSyncer) pendingNonceFromEthClient(ctx context.Context, account common.Address) (nextNonce uint64, err error) {
	ctx, cancel := eth.DefaultQueryCtx(ctx)
	defer cancel()
	nextNonce, err = s.ethClient.PendingNonceAt(ctx, account)
	return nextNonce, errors.WithStack(err)
}

func (s NonceSyncer) hasInProgressTransaction(account common.Address) (exists bool, err error) {
	err = postgres.DBWithDefaultContext(s.db, func(db *gorm.DB) error {
		return db.Raw(`SELECT EXISTS(SELECT 1 FROM eth_txes WHERE state = 'in_progress' AND from_address = ?)`, account).Scan(&exists).Error
	})
	return
}
