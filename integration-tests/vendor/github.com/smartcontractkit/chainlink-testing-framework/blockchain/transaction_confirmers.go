package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"sync"

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
	minConfirmations      int
	confirmations         int
	client                EVMClient
	tx                    *types.Transaction
	doneChan              chan struct{}
	context               context.Context
	cancel                context.CancelFunc
	networkConfig         *EVMNetwork
	lastReceivedHeaderNum uint64
	complete              bool
	completeMu            sync.Mutex
}

// NewTransactionConfirmer returns a new instance of the transaction confirmer that waits for on-chain minimum
// confirmations
func NewTransactionConfirmer(client EVMClient, tx *types.Transaction, minConfirmations int) *TransactionConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), client.GetNetworkConfig().Timeout.Duration)
	tc := &TransactionConfirmer{
		minConfirmations: minConfirmations,
		confirmations:    0,
		client:           client,
		tx:               tx,
		doneChan:         make(chan struct{}, 1),
		context:          ctx,
		cancel:           ctxCancel,
		networkConfig:    client.GetNetworkConfig(),
		complete:         false,
	}
	return tc
}

// ReceiveHeader the implementation of the HeaderEventSubscription that receives each header and checks
// tx confirmation
func (t *TransactionConfirmer) ReceiveHeader(header NodeHeader) error {
	if header.Number.Uint64() <= t.lastReceivedHeaderNum {
		return nil // Header with same number mined, disregard for confirming
	}
	t.lastReceivedHeaderNum = header.Number.Uint64()
	confirmationLog := log.Debug().
		Str("Network Name", t.networkConfig.Name).
		Str("Header Hash", header.Hash().Hex()).
		Str("Header Number", header.Number.String()).
		Str("Tx Hash", t.tx.Hash().String()).
		Uint64("Nonce", t.tx.Nonce()).
		Int("Minimum Confirmations", t.minConfirmations)
	isConfirmed, err := t.client.IsTxConfirmed(t.tx.Hash())
	if err != nil {
		return err
	} else if isConfirmed {
		t.confirmations++
	}
	if t.confirmations >= t.minConfirmations {
		confirmationLog.Int("Current Confirmations", t.confirmations).Msg("Transaction confirmations met")
		t.complete = true
		t.doneChan <- struct{}{}
	} else {
		confirmationLog.Int("Current Confirmations", t.confirmations).Msg("Waiting on minimum confirmations")
	}
	return nil
}

// Wait is a blocking function that waits until the transaction is complete
func (t *TransactionConfirmer) Wait() error {
	defer func() {
		t.completeMu.Lock()
		t.complete = true
		t.completeMu.Unlock()
	}()

	if t.Complete() {
		t.cancel()
		return nil
	}

	for {
		select {
		case <-t.doneChan:
			t.cancel()
			return nil
		case <-t.context.Done():
			return fmt.Errorf("timeout waiting for transaction to confirm: %s network %s", t.tx.Hash(), t.client.GetNetworkName())
		}
	}
}

// Complete returns if the confirmer has completed or not
func (t *TransactionConfirmer) Complete() bool {
	t.completeMu.Lock()
	defer t.completeMu.Unlock()
	return t.complete
}

// InstantConfirmer is a near-instant confirmation method, primarily for optimistic L2s that have near-instant finalization
type InstantConfirmer struct {
	client       EVMClient
	txHash       common.Hash
	complete     bool // tracks if the subscription is completed or not
	completeChan chan struct{}
	completeMu   sync.Mutex
	context      context.Context
	cancel       context.CancelFunc
	// For events
	confirmed     bool // tracks the confirmation status of the subscription
	confirmedChan chan bool
	errorChan     chan error
}

func NewInstantConfirmer(
	client EVMClient,
	txHash common.Hash,
	confirmedChan chan bool,
	errorChan chan error,
) *InstantConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), client.GetNetworkConfig().Timeout.Duration)
	return &InstantConfirmer{
		client:       client,
		txHash:       txHash,
		completeChan: make(chan struct{}, 1),
		context:      ctx,
		cancel:       ctxCancel,
		// For events
		confirmedChan: confirmedChan,
		errorChan:     errorChan,
	}
}

// ReceiveHeader does a quick check on if the tx is confirmed already
func (l *InstantConfirmer) ReceiveHeader(_ NodeHeader) error {
	var err error
	l.confirmed, err = l.client.IsTxConfirmed(l.txHash)
	if err != nil {
		if err.Error() == "not found" {
			log.Debug().Str("Tx", l.txHash.Hex()).Msg("Transaction not found on chain yet. Waiting to confirm.")
			return err
		}
		log.Error().Str("Tx", l.txHash.Hex()).Err(err).Msg("Error checking tx confirmed")
		if l.errorChan != nil {
			l.errorChan <- err
		}
		return err
	}
	log.Debug().Bool("Confirmed", l.confirmed).Str("Tx", l.txHash.Hex()).Msg("Instant Confirmation")
	if l.confirmed {
		l.completeChan <- struct{}{}
		if l.confirmedChan != nil {
			l.confirmedChan <- l.confirmed
		}
	}
	return nil
}

// Wait checks every header if the tx has been included on chain or not
func (l *InstantConfirmer) Wait() error {
	defer func() {
		l.completeMu.Lock()
		l.complete = true
		l.completeMu.Unlock()
	}()

	for {
		select {
		case <-l.completeChan:
			l.cancel()
			return nil
		case <-l.context.Done():
			return fmt.Errorf("timeout waiting for instant transaction to confirm after %s: %s",
				l.client.GetNetworkConfig().Timeout.String(), l.txHash.Hex())
		}
	}
}

// Complete returns if the transaction is complete or not
func (l *InstantConfirmer) Complete() bool {
	l.completeMu.Lock()
	defer l.completeMu.Unlock()
	return l.complete
}

// EventConfirmer confirms that an event is confirmed by a certain amount of headers
type EventConfirmer struct {
	eventName             string
	minConfirmations      int
	confirmations         int
	client                EVMClient
	event                 *types.Log
	waitChan              chan struct{}
	errorChan             chan error
	confirmedChan         chan bool
	context               context.Context
	cancel                context.CancelFunc
	lastReceivedHeaderNum uint64
	complete              bool
}

// NewEventConfirmer returns a new instance of the event confirmer that waits for on-chain minimum
// confirmations
func NewEventConfirmer(
	eventName string,
	client EVMClient,
	event *types.Log,
	minConfirmations int,
	confirmedChan chan bool,
	errorChan chan error,
) *EventConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), client.GetNetworkConfig().Timeout.Duration)
	tc := &EventConfirmer{
		eventName:        eventName,
		minConfirmations: minConfirmations,
		confirmations:    0,
		client:           client,
		event:            event,
		waitChan:         make(chan struct{}, 1),
		errorChan:        errorChan,
		confirmedChan:    confirmedChan,
		context:          ctx,
		cancel:           ctxCancel,
		complete:         false,
	}
	return tc
}

// ReceiveHeader will attempt to confirm an event for the chain's configured minimum confirmed headers. Errors encountered
// are sent along the eventErrorChan, and the result of confirming the event is sent to eventConfirmedChan.
func (e *EventConfirmer) ReceiveHeader(header NodeHeader) error {
	if header.Number.Uint64() <= e.lastReceivedHeaderNum {
		return nil
	}
	e.lastReceivedHeaderNum = header.Number.Uint64()
	confirmed, removed, err := e.client.IsEventConfirmed(e.event)
	if err != nil {
		e.errorChan <- err
		return err
	}
	if removed {
		e.confirmedChan <- false
		e.complete = true
		return nil
	}
	if confirmed {
		e.confirmations++
	}
	if e.confirmations >= e.minConfirmations {
		e.confirmedChan <- true
		e.complete = true
	}
	return nil
}

// Wait until the event fully presents as complete
func (e *EventConfirmer) Wait() error {
	defer func() { e.complete = true }()
	for {
		select {
		case <-e.waitChan:
			e.cancel()
			return nil
		case <-e.context.Done():
			return fmt.Errorf("timeout waiting for event to confirm after %s: %s",
				e.client.GetNetworkConfig().Timeout.String(), e.event.TxHash.Hex())
		}
	}
}

// Complete returns if the confirmer is done, whether confirmation was successful or not
func (e *EventConfirmer) Complete() bool {
	return e.complete
}

// GetHeaderSubscriptions returns a duplicate map of the queued transactions
func (e *EthereumClient) GetHeaderSubscriptions() map[string]HeaderEventSubscription {
	e.subscriptionMutex.Lock()
	defer e.subscriptionMutex.Unlock()

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
			log.Error().Err(err).Msg("Error while subscribed to new headers, restarting subscription")
			subscription.Unsubscribe()

			subscription, err = e.Client.SubscribeNewHead(context.Background(), headerChannel)
			if err != nil {
				log.Error().Err(err).Msg("Failed to resubscribe to new headers")
				return err
			}
		case header := <-headerChannel:
			e.receiveHeader(header)
		case <-e.doneChan:
			log.Debug().Str("Network", e.NetworkConfig.Name).Msg("Subscription cancelled")
			e.Client.Close()
			return nil
		}
	}
}

// receiveHeader takes in a new header from the chain, and sends the header to all active header subscriptions
func (e *EthereumClient) receiveHeader(header *types.Header) {
	if header == nil {
		log.Debug().Msg("Received Nil Header")
		return
	}
	headerValue := *header

	suggestedPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		suggestedPrice = big.NewInt(0)
		log.Err(err).
			Str("Header Hash", headerValue.Hash().String()).
			Msg("Error retrieving Suggested Gas Price for new block header")
	}
	log.Debug().
		Str("NetworkName", e.NetworkConfig.Name).
		Int("Node", e.ID).
		Str("Hash", headerValue.Hash().String()).
		Str("Number", headerValue.Number.String()).
		Str("Gas Price", suggestedPrice.String()).
		Msg("Received block header")

	subs := e.GetHeaderSubscriptions()

	g := errgroup.Group{}
	for _, sub := range subs {
		sub := sub
		g.Go(func() error {
			return sub.ReceiveHeader(NodeHeader{NodeID: e.ID, Header: headerValue})
		})
	}
	if err := g.Wait(); err != nil {
		log.Err(fmt.Errorf("error on sending block header to receivers: %v", err))
	}
	if len(subs) > 0 {
		var subsRemoved uint
		for key, sub := range subs { // Cleanup subscriptions that might not have Wait called on them
			if sub.Complete() {
				subsRemoved++
				e.DeleteHeaderEventSubscription(key)
			}
		}
		if subsRemoved > 0 {
			log.Trace().
				Uint("Recently Removed", subsRemoved).
				Int("Active", len(e.GetHeaderSubscriptions())).
				Msg("Updated Header Subscriptions")
		}
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
