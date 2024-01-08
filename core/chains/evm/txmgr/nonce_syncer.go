package txmgr

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/common/txmgr"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// NonceSyncer manages the delicate task of syncing the local nonce with the
// chain nonce in case of divergence.
//
// On startup, we check each key for the nonce value on chain and compare
// it to our local value.
//
// Usually the on-chain nonce will be the same as (or lower than) the
// highest sequence in the DB, in which case we do nothing.
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
var _ txmgr.SequenceSyncer[common.Address, common.Hash, common.Hash, types.Nonce] = &nonceSyncerImpl{}

type nonceSyncerImpl struct {
	txStore EvmTxStore
	client  TxmClient
	chainID *big.Int
	logger  logger.Logger
}

// NewNonceSyncer returns a new syncer
func NewNonceSyncer(
	txStore EvmTxStore,
	lggr logger.Logger,
	ethClient evmclient.Client,
) NonceSyncer {
	lggr = logger.Named(lggr, "NonceSyncer")
	return &nonceSyncerImpl{
		txStore: txStore,
		client:  NewEvmTxmClient(ethClient),
		chainID: ethClient.ConfiguredChainID(),
		logger:  lggr,
	}
}

// SyncAll syncs nonces for all enabled keys in parallel
//
// This should only be called once, before the EthBroadcaster has started.
// Calling it later is not safe and could lead to races.
func (s nonceSyncerImpl) Sync(ctx context.Context, addr common.Address, localNonce types.Nonce) (nonce types.Nonce, err error) {
	nonce, err = s.fastForwardNonceIfNecessary(ctx, addr, localNonce)
	return nonce, errors.Wrap(err, "NonceSyncer#fastForwardNoncesIfNecessary failed")
}

func (s nonceSyncerImpl) fastForwardNonceIfNecessary(ctx context.Context, address common.Address, localNonce types.Nonce) (types.Nonce, error) {
	chainNonce, err := s.pendingNonceFromEthClient(ctx, address)
	if err != nil {
		return localNonce, errors.Wrap(err, "GetNextNonce failed to loadInitialNonceFromEthClient")
	}
	if chainNonce == 0 {
		return localNonce, nil
	}
	if chainNonce <= localNonce {
		return localNonce, nil
	}
	s.logger.Warnw(fmt.Sprintf("address %s has been used before, either by an external wallet or a different Chainlink node. "+
		"Local nonce is %v but the on-chain nonce for this account was %v. "+
		"It's possible that this node was restored from a backup. If so, transactions sent by the previous node will NOT be re-org protected and in rare cases may need to be manually bumped/resubmitted. "+
		"Please note that using the chainlink keys with an external wallet is NOT SUPPORTED and can lead to missed or stuck transactions. ",
		address, localNonce, chainNonce),
		"address", address.String(), "localNonce", localNonce, "chainNonce", chainNonce)

	return chainNonce, nil
}

func (s nonceSyncerImpl) pendingNonceFromEthClient(ctx context.Context, account common.Address) (types.Nonce, error) {
	nextNonce, err := s.client.PendingSequenceAt(ctx, account)
	return nextNonce, errors.WithStack(err)
}
