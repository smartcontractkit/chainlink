package txmgr

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/common/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

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
// within the past EVM.FinalityDepth blocks and find the ones sent by our
// address(es).
//
// This gives us re-org protection up to EVM.FinalityDepth deep in the
// worst case, which is in line with our other guarantees.
type NonceSyncer[ADDR types.Hashable, TX_HASH types.Hashable, BLOCK_HASH types.Hashable] interface {
	Sync(ctx context.Context, addr ADDR) (err error)
}

var _ NonceSyncer[common.Address, common.Hash, common.Hash] = &nonceSyncerImpl{}

type nonceSyncerImpl struct {
	txStore   EvmTxStore
	ethClient evmclient.Client
	chainID   *big.Int
	logger    logger.Logger
	kst       EvmKeyStore
}

// NewNonceSyncer returns a new syncer
func NewNonceSyncer(
	txStore EvmTxStore,
	lggr logger.Logger,
	ethClient evmclient.Client,
	kst EvmKeyStore,
) EvmNonceSyncer {
	lggr = lggr.Named("NonceSyncer")
	return &nonceSyncerImpl{
		txStore:   txStore,
		ethClient: ethClient,
		chainID:   ethClient.ConfiguredChainID(),
		logger:    lggr,
		kst:       kst,
	}
}

// SyncAll syncs nonces for all enabled keys in parallel
//
// This should only be called once, before the EthBroadcaster has started.
// Calling it later is not safe and could lead to races.
func (s nonceSyncerImpl) Sync(ctx context.Context, addr common.Address) (err error) {
	err = s.fastForwardNonceIfNecessary(ctx, addr)
	return errors.Wrap(err, "NonceSyncer#fastForwardNoncesIfNecessary failed")
}

func (s nonceSyncerImpl) fastForwardNonceIfNecessary(ctx context.Context, address common.Address) error {
	chainNonce, err := s.pendingNonceFromEthClient(ctx, address)
	if err != nil {
		return errors.Wrap(err, "GetNextNonce failed to loadInitialNonceFromEthClient")
	}
	if chainNonce == 0 {
		return nil
	}

	keyNextNonce, err := s.kst.NextSequence(address, s.chainID, pg.WithParentCtx(ctx))
	if err != nil {
		return err
	}

	localNonce := keyNextNonce
	hasInProgressTransaction, err := s.txStore.HasInProgressTransaction(address, s.chainID, pg.WithParentCtx(ctx))

	if err != nil {
		return errors.Wrapf(err, "failed to query for in_progress transaction for address %s", address.String())
	} else if hasInProgressTransaction {
		// If we have an 'in_progress' transaction, our keys.next_nonce will be
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
		address, localNonce, chainNonce),
		"address", address.String(), "keyNextNonce", keyNextNonce, "localNonce", localNonce, "chainNonce", chainNonce)

	// Need to remember to decrement the chain nonce by one to account for in_progress transaction
	newNextNonce := int64(chainNonce)
	if hasInProgressTransaction {
		newNextNonce--
	}

	err = s.txStore.UpdateEthKeyNextNonce(evmtypes.Nonce(newNextNonce), keyNextNonce, address, s.chainID, pg.WithParentCtx(ctx))

	if errors.Is(err, ErrKeyNotUpdated) {
		return errors.Errorf("NonceSyncer#fastForwardNonceIfNecessary optimistic lock failure fastforwarding nonce %v to %v for key %s", localNonce, chainNonce, address.String())
	} else if err == nil {
		s.logger.Infow("Fast-forwarded nonce", "address", address, "newNextNonce", newNextNonce, "oldNextNonce", keyNextNonce)
	}
	return err
}

func (s nonceSyncerImpl) pendingNonceFromEthClient(ctx context.Context, account common.Address) (nextNonce uint64, err error) {
	nextNonce, err = s.ethClient.PendingNonceAt(ctx, account)
	return nextNonce, errors.WithStack(err)
}
