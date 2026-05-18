package txmgr

import (
	"context"
	"fmt"
	"math/big"
	"slices"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jpillora/backoff"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type NonceTrackerTxStore interface {
	FindLatestSequence(context.Context, common.Address, *big.Int) (evmtypes.Nonce, error)
}

type NonceTrackerClient interface {
	ConfiguredChainID() *big.Int
	PendingSequenceAt(context.Context, common.Address) (evmtypes.Nonce, error)
}

type nonceTracker struct {
	lggr             logger.SugaredLogger
	nextSequenceMap  map[common.Address]evmtypes.Nonce
	txStore          NonceTrackerTxStore
	chainID          *big.Int
	client           NonceTrackerClient
	enabledAddresses []common.Address

	sequenceLock sync.RWMutex
}

func NewNonceTracker(lggr logger.Logger, txStore NonceTrackerTxStore, client NonceTrackerClient) *nonceTracker {
	lggr = logger.Named(lggr, "NonceTracker")
	return &nonceTracker{
		lggr:    logger.Sugared(lggr),
		txStore: txStore,
		chainID: client.ConfiguredChainID(),
		client:  client,
	}
}

func (s *nonceTracker) LoadNextSequences(ctx context.Context, addresses []common.Address) {
	s.sequenceLock.Lock()
	defer s.sequenceLock.Unlock()

	s.enabledAddresses = addresses

	s.nextSequenceMap = make(map[common.Address]evmtypes.Nonce)
	for _, address := range addresses {
		seq, err := s.getSequenceForAddr(ctx, address)
		if err == nil {
			s.nextSequenceMap[address] = seq
		}
	}
}

func (s *nonceTracker) getSequenceForAddr(ctx context.Context, address common.Address) (seq evmtypes.Nonce, err error) {
	// Get the highest sequence from the tx table
	// Will need to be incremented since this sequence is already used
	seq, err = s.txStore.FindLatestSequence(ctx, address, s.chainID)
	if err == nil {
		seq++
		return seq, nil
	}
	// Look for nonce on-chain if no tx found for address in TxStore or if error occurred
	// Returns the nonce that should be used for the next transaction so no need to increment
	nonce, err := s.client.PendingSequenceAt(ctx, address)
	if err == nil {
		return nonce, nil
	}
	s.lggr.Criticalw("failed to retrieve next sequence from on-chain for address: ", "address", address.String())
	return seq, err
}

// syncSequence tries to sync the key sequence, retrying indefinitely until success or stop signal is sent
func (s *nonceTracker) SyncSequence(ctx context.Context, addr common.Address, chStop services.StopChan) {
	sequenceSyncRetryBackoff := backoff.Backoff{
		Min:    100 * time.Millisecond,
		Max:    5 * time.Second,
		Jitter: true,
	}

	localSequence, err := s.GetNextSequence(ctx, addr)
	// Address not found in map so skip sync
	if err != nil {
		s.lggr.Criticalw("Failed to retrieve local next sequence for address", "address", addr.String(), "err", err)
		return
	}

	// Enter loop with retries
	var attempt int
	for {
		select {
		case <-chStop:
			return
		case <-time.After(sequenceSyncRetryBackoff.Duration()):
			attempt++
			err := s.SyncOnChain(ctx, addr, localSequence)
			if err != nil {
				if attempt > 5 {
					s.lggr.Criticalw("Failed to sync with on-chain sequence", "address", addr.String(), "attempt", attempt, "err", err)
				} else {
					s.lggr.Warnw("Failed to sync with on-chain sequence", "address", addr.String(), "attempt", attempt, "err", err)
				}
				continue
			}
			return
		}
	}
}

func (s *nonceTracker) SyncOnChain(ctx context.Context, addr common.Address, localSequence evmtypes.Nonce) error {
	nonce, err := s.client.PendingSequenceAt(ctx, addr)
	if err != nil {
		return err
	}
	if nonce > localSequence {
		s.lggr.Warnw(fmt.Sprintf("address %s has been used before, either by an external wallet or a different Chainlink node. "+
			"Local nonce is %v but the on-chain nonce for this account was %v. "+
			"It's possible that this node was restored from a backup. If so, transactions sent by the previous node will NOT be re-org protected and in rare cases may need to be manually bumped/resubmitted. "+
			"Please note that using the chainlink keys with an external wallet is NOT SUPPORTED and can lead to missed or stuck transactions. ",
			addr, localSequence, nonce),
			"address", addr.String(), "localNonce", localSequence, "chainNonce", nonce)

		s.lggr.Infow("Fast-forward sequence", "address", addr, "newNextSequence", nonce, "oldNextSequence", localSequence)
	}

	s.sequenceLock.Lock()
	defer s.sequenceLock.Unlock()
	s.nextSequenceMap[addr] = max(localSequence, nonce)
	return nil
}

func (s *nonceTracker) GetNextSequence(ctx context.Context, address common.Address) (seq evmtypes.Nonce, err error) {
	s.sequenceLock.Lock()
	defer s.sequenceLock.Unlock()
	// Get next sequence from map
	seq, exists := s.nextSequenceMap[address]
	if exists {
		return seq, nil
	}

	s.lggr.Infow("address not found in local next sequence map. Attempting to search and populate sequence.", "address", address.String())
	// Check if address is in the enabled address list
	if !slices.Contains(s.enabledAddresses, address) {
		return seq, fmt.Errorf("address disabled: %s", address)
	}

	// Try to retrieve next sequence from tx table or on-chain to load the map
	// A scenario could exist where loading the map during startup failed (e.g. All configured RPC's are unreachable at start)
	// The expectation is that the node does not fail startup so sequences need to be loaded during runtime
	foundSeq, err := s.getSequenceForAddr(ctx, address)
	if err != nil {
		return seq, fmt.Errorf("failed to find next sequence for address: %s", address)
	}

	// Set sequence in map
	s.nextSequenceMap[address] = foundSeq
	return foundSeq, nil
}

func (s *nonceTracker) GenerateNextSequence(address common.Address, nonceUsed evmtypes.Nonce) {
	s.sequenceLock.Lock()
	defer s.sequenceLock.Unlock()
	currentNonce := s.nextSequenceMap[address]

	// In most cases, currentNonce would equal nonceUsed
	// There is a chance currentNonce is 1 ahead of nonceUsed if the DB contains an in-progress tx during startup
	// Incrementing currentNonce, which is already set to the next usable nonce, could lead to a nonce gap. Set the map to the incremented nonceUsed instead.
	if currentNonce == nonceUsed || currentNonce == nonceUsed+1 {
		s.nextSequenceMap[address] = nonceUsed + 1
		return
	}

	// If currentNonce is ahead of even the incremented nonceUsed, maintain the unchanged currentNonce in the map
	// This scenario should never occur but logging this discrepancy for visibility
	s.lggr.Warnf("Local nonce map value %d for address %s is ahead of the nonce transmitted %d. Maintaining the existing value in the map without incrementing.", currentNonce, address.String(), nonceUsed)
}
