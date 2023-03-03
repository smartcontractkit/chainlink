package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

// Used for when running tests on a live test network, so tests can share nonces and run in parallel on the same network
var (
	globalNonceManager = make(map[uint64]*NonceSettings)
	globalNonceLock    sync.Mutex
)

// useGlobalNonceManager for when running tests on a non-simulated network
func useGlobalNonceManager(chainId *big.Int) *NonceSettings {
	globalNonceLock.Lock()
	defer globalNonceLock.Unlock()
	if _, ok := globalNonceManager[chainId.Uint64()]; !ok {
		fmt.Printf("Using a new Global Nonce Manager for chain %d\n%v", chainId.Uint64(), globalNonceManager)
		globalNonceManager[chainId.Uint64()] = newNonceSettings()
		go globalNonceManager[chainId.Uint64()].watchInstantTransactions()
	}

	return globalNonceManager[chainId.Uint64()]
}

// convenience function
func newNonceSettings() *NonceSettings {
	return &NonceSettings{
		NonceMu: &sync.Mutex{},
		Nonces:  make(map[string]uint64),

		doneChan:            make(chan struct{}),
		instantTransactions: make(map[string]map[uint64]chan struct{}),
		instantNonces:       make(map[string]uint64),
		registerChan:        make(chan instantTxRegistration),
		sentChan:            make(chan string),
	}
}

// NonceSettings is a convenient wrapper for holding nonce state
type NonceSettings struct {
	NonceMu *sync.Mutex
	Nonces  map[string]uint64

	// used to properly meter out instant txs on L2s
	doneChan            chan struct{}
	instantTransactions map[string]map[uint64]chan struct{}
	instantNonces       map[string]uint64
	instantNoncesMu     sync.Mutex
	registerChan        chan instantTxRegistration
	sentChan            chan string
}

// watchInstantTransactions should only be called when minConfirmations for the chain is 0, generally an L2 chain.
// This helps meter out transactions to L2 chains, so that nonces only send in order. For most (if not all) L2 chains,
// the mempool is small or non-existent, meaning we can't send nonces out of order, otherwise the tx is instantly
// rejected.
func (ns *NonceSettings) watchInstantTransactions() {
	ns.instantTransactions = make(map[string]map[uint64]chan struct{})

	for {
		select {
		case toRegister := <-ns.registerChan:
			if _, ok := ns.instantTransactions[toRegister.fromAddr]; !ok {
				ns.instantTransactions[toRegister.fromAddr] = make(map[uint64]chan struct{})
			}
			ns.instantTransactions[toRegister.fromAddr][toRegister.nonce] = toRegister.releaseChan
		case sentAddr := <-ns.sentChan:
			ns.instantNoncesMu.Lock()
			ns.instantNonces[sentAddr]++
			ns.instantNoncesMu.Unlock()
		default:
			for addr, releaseChannels := range ns.instantTransactions {
				ns.instantNoncesMu.Lock()
				nonceToSend := ns.instantNonces[addr]
				ns.instantNoncesMu.Unlock()
				if txChannel, ok := releaseChannels[nonceToSend]; ok {
					close(txChannel)
					delete(releaseChannels, nonceToSend)
				}
			}
		}
	}
}

// registerInstantTransaction helps meter out txs for L2 chains. Register, then wait to receive from the returned channel
// to know when your Tx can send. See watchInstantTransactions for a deeper explanation.
func (ns *NonceSettings) registerInstantTransaction(fromAddr string, nonce uint64) chan struct{} {
	releaseChan := make(chan struct{})
	ns.registerChan <- instantTxRegistration{
		fromAddr:    fromAddr,
		nonce:       nonce,
		releaseChan: releaseChan,
	}
	return releaseChan
}

// sentInstantTransaction shows that you have sent this instant transaction, unlocking the next L2 transaction to run.
// See watchInstantTransactions for a deeper explanation.
func (ns *NonceSettings) sentInstantTransaction(fromAddr string) {
	ns.sentChan <- fromAddr
}

// GetNonce keep tracking of nonces per address, add last nonce for addr if the map is empty
func (e *EthereumClient) GetNonce(ctx context.Context, addr common.Address) (uint64, error) {
	e.NonceSettings.NonceMu.Lock()
	defer e.NonceSettings.NonceMu.Unlock()

	// See current state of the nonce manager, handy for debugging
	// fmt.Println("-------Nonce Manager Current State-----------------")
	// for address, nonce := range e.NonceSettings.Nonces {
	// 	fmt.Printf("%s: %d\n", address, nonce)
	// }
	// fmt.Println("---------------------------------------------------")

	if _, ok := e.NonceSettings.Nonces[addr.Hex()]; !ok {
		pendingNonce, err := e.Client.PendingNonceAt(ctx, addr)
		if err != nil {
			return 0, err
		}
		e.NonceSettings.Nonces[addr.Hex()] = pendingNonce

		e.NonceSettings.instantNoncesMu.Lock()
		e.NonceSettings.instantNonces[addr.Hex()] = pendingNonce
		e.NonceSettings.instantNoncesMu.Unlock()

		return pendingNonce, nil
	}
	e.NonceSettings.Nonces[addr.Hex()]++
	return e.NonceSettings.Nonces[addr.Hex()], nil
}

// PeekPendingNonce returns the current pending nonce for the address. Does not change any nonce settings state
func (e *EthereumClient) PeekPendingNonce(addr common.Address) (uint64, error) {
	e.NonceSettings.NonceMu.Lock()
	defer e.NonceSettings.NonceMu.Unlock()
	if _, ok := e.NonceSettings.Nonces[addr.Hex()]; !ok {
		pendingNonce, err := e.Client.PendingNonceAt(context.Background(), addr)
		if err != nil {
			return 0, err
		}
		e.NonceSettings.Nonces[addr.Hex()] = pendingNonce
	}
	return e.NonceSettings.Nonces[addr.Hex()], nil
}

type instantTxRegistration struct {
	fromAddr    string
	nonce       uint64
	releaseChan chan struct{}
}
