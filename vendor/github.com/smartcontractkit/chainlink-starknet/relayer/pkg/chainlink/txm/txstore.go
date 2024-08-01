package txm

import (
	"fmt"
	"sort"
	"sync"

	"github.com/NethermindEth/juno/core/felt"
	starknetrpc "github.com/NethermindEth/starknet.go/rpc"
	"golang.org/x/exp/maps"
)

type UnconfirmedTx struct {
	Hash      string
	PublicKey *felt.Felt
	Nonce     *felt.Felt
	Call      starknetrpc.FunctionCall
}

// TxStore tracks broadcast & unconfirmed txs per account address per chain id
type TxStore struct {
	lock sync.RWMutex

	nextNonce         *felt.Felt
	unconfirmedNonces map[string]*UnconfirmedTx
}

func NewTxStore(initialNonce *felt.Felt) *TxStore {
	return &TxStore{
		nextNonce:         new(felt.Felt).Set(initialNonce),
		unconfirmedNonces: map[string]*UnconfirmedTx{},
	}
}

func (s *TxStore) SetNextNonce(newNextNonce *felt.Felt) []*UnconfirmedTx {
	s.lock.Lock()
	defer s.lock.Unlock()

	staleTxs := []*UnconfirmedTx{}
	s.nextNonce = new(felt.Felt).Set(newNextNonce)

	// Remove any stale transactions with nonces greater than the new next nonce.
	for nonceStr, tx := range s.unconfirmedNonces {
		if tx.Nonce.Cmp(s.nextNonce) >= 0 {
			staleTxs = append(staleTxs, tx)
			delete(s.unconfirmedNonces, nonceStr)
		}
	}

	sort.Slice(staleTxs, func(i, j int) bool {
		a := staleTxs[i]
		b := staleTxs[j]
		return a.Nonce.Cmp(b.Nonce) < 0
	})

	return staleTxs
}

func (s *TxStore) GetNextNonce() *felt.Felt {
	s.lock.Lock()
	defer s.lock.Unlock()
	return new(felt.Felt).Set(s.nextNonce)
}

func (s *TxStore) AddUnconfirmed(nonce *felt.Felt, hash string, call starknetrpc.FunctionCall, publicKey *felt.Felt) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if nonce.Cmp(s.nextNonce) < 0 {
		return fmt.Errorf("tried to add an unconfirmed tx at an old nonce: expected %s, got %s", s.nextNonce, nonce)
	}
	if nonce.Cmp(s.nextNonce) > 0 {
		return fmt.Errorf("tried to add an unconfirmed tx at a future nonce: expected %s, got %s", s.nextNonce, nonce)
	}

	nonceStr := nonce.String()
	if h, exists := s.unconfirmedNonces[nonceStr]; exists {
		return fmt.Errorf("nonce used: tried to use nonce (%s) for tx (%s), already used by (%s)", nonce, h.Hash, h)
	}

	s.unconfirmedNonces[nonceStr] = &UnconfirmedTx{
		Nonce:     new(felt.Felt).Set(nonce),
		PublicKey: new(felt.Felt).Set(publicKey),
		Hash:      hash,
		Call:      call,
	}

	s.nextNonce = new(felt.Felt).Add(s.nextNonce, new(felt.Felt).SetUint64(1))
	return nil
}

func (s *TxStore) Confirm(nonce *felt.Felt, hash string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	nonceStr := nonce.String()
	unconfirmed, exists := s.unconfirmedNonces[nonceStr]
	if !exists {
		return fmt.Errorf("no such unconfirmed nonce: %s", nonce)
	}
	// sanity check that the hash matches
	if unconfirmed.Hash != hash {
		return fmt.Errorf("unexpected tx hash: expected %s, got %s", unconfirmed.Hash, hash)
	}
	delete(s.unconfirmedNonces, nonceStr)
	return nil
}

func (s *TxStore) GetUnconfirmed() []*UnconfirmedTx {
	s.lock.RLock()
	defer s.lock.RUnlock()

	unconfirmed := maps.Values(s.unconfirmedNonces)
	sort.Slice(unconfirmed, func(i, j int) bool {
		a := unconfirmed[i]
		b := unconfirmed[j]
		return a.Nonce.Cmp(b.Nonce) < 0
	})

	return unconfirmed
}

func (s *TxStore) InflightCount() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.unconfirmedNonces)
}

type AccountStore struct {
	store map[string]*TxStore // map account address to txstore
	lock  sync.RWMutex
}

func NewAccountStore() *AccountStore {
	return &AccountStore{
		store: map[string]*TxStore{},
	}
}

func (c *AccountStore) CreateTxStore(accountAddress *felt.Felt, initialNonce *felt.Felt) (*TxStore, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	addressStr := accountAddress.String()
	_, ok := c.store[addressStr]
	if ok {
		return nil, fmt.Errorf("TxStore already exists: %s", accountAddress)
	}
	store := NewTxStore(initialNonce)
	c.store[addressStr] = store
	return store, nil
}

// GetTxStore returns the TxStore for the provided account.
func (c *AccountStore) GetTxStore(accountAddress *felt.Felt) *TxStore {
	c.lock.Lock()
	defer c.lock.Unlock()
	store, ok := c.store[accountAddress.String()]
	if !ok {
		return nil
	}
	return store
}

func (c *AccountStore) GetTotalInflightCount() int {
	// use read lock for methods that read underlying data
	c.lock.RLock()
	defer c.lock.RUnlock()

	count := 0
	for _, store := range c.store {
		count += store.InflightCount()
	}

	return count
}

func (c *AccountStore) GetAllUnconfirmed() map[string][]*UnconfirmedTx {
	// use read lock for methods that read underlying data
	c.lock.RLock()
	defer c.lock.RUnlock()

	allUnconfirmed := map[string][]*UnconfirmedTx{}
	for accountAddressStr, store := range c.store {
		allUnconfirmed[accountAddressStr] = store.GetUnconfirmed()
	}
	return allUnconfirmed
}
