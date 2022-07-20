package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

// TransactionConfirmer is an implementation of HeaderEventSubscription that checks whether tx are confirmed
type TransactionConfirmer struct {
	minConfirmations int
	confirmations    int
	client           EVMClient
	tx               *types.Transaction
	doneChan         chan struct{}
	context          context.Context
	cancel           context.CancelFunc
	networkConfig    *EVMNetwork
}

// NewTransactionConfirmer returns a new instance of the transaction confirmer that waits for on-chain minimum
// confirmations
func NewTransactionConfirmer(client EVMClient, tx *types.Transaction, minConfirmations int) *TransactionConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), client.GetNetworkConfig().Timeout)
	tc := &TransactionConfirmer{
		minConfirmations: minConfirmations,
		confirmations:    0,
		client:           client,
		tx:               tx,
		doneChan:         make(chan struct{}, 1),
		context:          ctx,
		cancel:           ctxCancel,
		networkConfig:    client.GetNetworkConfig(),
	}
	return tc
}

// ReceiveBlock the implementation of the HeaderEventSubscription that receives each block and checks
// tx confirmation
func (t *TransactionConfirmer) ReceiveBlock(block NodeBlock) error {
	if block.Block == nil {
		// Strange case that happens in some EVM testnets
		log.Info().Msg("Received nil block")
		return nil
	}
	confirmationLog := log.Debug().Str("Network Name", t.networkConfig.Name).
		Str("Block Hash", block.Hash().Hex()).
		Str("Block Number", block.Number().String()).
		Str("Tx Hash", t.tx.Hash().String()).
		Uint64("Nonce", t.tx.Nonce()).
		Int("Minimum Confirmations", t.minConfirmations)
	isConfirmed, err := t.client.IsTxConfirmed(t.tx.Hash())
	if err != nil {
		return err
	} else if isConfirmed {
		t.confirmations++
	}
	if t.confirmations == t.minConfirmations {
		confirmationLog.Int("Current Confirmations", t.confirmations).
			Msg("Transaction confirmations met")
		t.doneChan <- struct{}{}
	} else if t.confirmations <= t.minConfirmations {
		confirmationLog.Int("Current Confirmations", t.confirmations).
			Msg("Waiting on minimum confirmations")
	}
	return nil
}

// Wait is a blocking function that waits until the transaction is complete
func (t *TransactionConfirmer) Wait() error {
	for {
		select {
		case <-t.doneChan:
			t.cancel()
			return nil
		case <-t.context.Done():
			return fmt.Errorf("timeout waiting for transaction to confirm: %s", t.tx.Hash())
		}
	}
}

// InstantConfirmations is a no-op confirmer as all transactions are instantly mined so no confirmations are needed
type InstantConfirmations struct{}

// ReceiveBlock is a no-op
func (i *InstantConfirmations) ReceiveBlock(block NodeBlock) error {
	return nil
}

// Wait is a no-op
func (i *InstantConfirmations) Wait() error {
	return nil
}

// GetNonce keep tracking of nonces per address, add last nonce for addr if the map is empty
func (e *EthereumClient) GetNonce(ctx context.Context, addr common.Address) (uint64, error) {
	if e.BorrowNonces {
		e.NonceMu.Lock()
		defer e.NonceMu.Unlock()
		if _, ok := e.Nonces[addr.Hex()]; !ok {
			lastNonce, err := e.Client.PendingNonceAt(ctx, addr)
			if err != nil {
				return 0, err
			}
			e.Nonces[addr.Hex()] = lastNonce
			return lastNonce, nil
		}
		e.Nonces[addr.Hex()]++
		return e.Nonces[addr.Hex()], nil
	}
	lastNonce, err := e.Client.PendingNonceAt(ctx, addr)
	if err != nil {
		return 0, err
	}
	return lastNonce, nil
}

// GetHeaderSubscriptions returns a duplicate map of the queued transactions
func (e *EthereumClient) GetHeaderSubscriptions() map[string]HeaderEventSubscription {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	newMap := map[string]HeaderEventSubscription{}
	for k, v := range e.headerSubscriptions {
		newMap[k] = v
	}
	return newMap
}

// subscribeToNewHeaders
func (e *EthereumClient) subscribeToNewHeaders() error {
	headerChannel := make(chan *types.Header)
	subscription, err := e.Client.SubscribeNewHead(context.Background(), headerChannel)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	log.Info().Str("Network", e.NetworkConfig.Name).Msg("Subscribed to new block headers")

	for {
		select {
		case err := <-subscription.Err():
			return err
		case header := <-headerChannel:
			e.receiveHeader(header)
		case <-e.doneChan:
			log.Debug().Str("Network", e.NetworkConfig.Name).Msg("Subscription cancelled")
			return nil
		}
	}
}

// receiveHeader
func (e *EthereumClient) receiveHeader(header *types.Header) {
	suggestedPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		suggestedPrice = big.NewInt(0)
		log.Err(err).
			Str("Block Hash", header.Hash().String()).
			Msg("Error retrieving Suggested Gas Price for new block header")
	}
	log.Debug().
		Str("NetworkName", e.NetworkConfig.Name).
		Int("Node", e.ID).
		Str("Hash", header.Hash().String()).
		Str("Number", header.Number.String()).
		Str("Gas Price", suggestedPrice.String()).
		Msg("Received block header")

	subs := e.GetHeaderSubscriptions()
	block, err := e.Client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		log.Err(fmt.Errorf("error fetching block by number: %v", err))
	}

	g := errgroup.Group{}
	for _, sub := range subs {
		sub := sub
		g.Go(func() error {
			return sub.ReceiveBlock(NodeBlock{NodeID: e.ID, Block: block})
		})
	}
	if err := g.Wait(); err != nil {
		log.Err(fmt.Errorf("error on sending block to receivers: %v", err))
	}
}

// errorReason decodes tx revert reason
func (e *EthereumClient) errorReason(
	b ethereum.ContractCaller,
	tx *types.Transaction,
	receipt *types.Receipt,
) (string, error) {
	chID, err := e.Client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}
	msg, err := tx.AsMessage(types.NewEIP155Signer(chID), nil)
	if err != nil {
		return "", err
	}
	callMsg := ethereum.CallMsg{
		From:     msg.From(),
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}
	res, err := b.CallContract(context.Background(), callMsg, receipt.BlockNumber)
	if err != nil {
		return "", errors.Wrap(err, "CallContract")
	}
	return abi.UnpackRevert(res)
}
