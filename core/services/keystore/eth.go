package keystore

import (
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name Eth --output mocks/ --case=underscore

// Eth is the external interface for EthKeyStore
type Eth interface {
	Get(id string) (ethkey.KeyV2, error)
	GetAll() ([]ethkey.KeyV2, error)
	Create(chainID *big.Int) (ethkey.KeyV2, error)
	Add(key ethkey.KeyV2, chainID *big.Int) error
	Delete(id string) (ethkey.KeyV2, error)
	Import(keyJSON []byte, password string, chainID *big.Int) (ethkey.KeyV2, error)
	Export(id string, password string) ([]byte, error)

	EnsureKeys(chainID *big.Int) error
	SubscribeToKeyChanges() (ch chan struct{}, unsub func())

	SignTx(fromAddress common.Address, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)

	SendingKeys(chainID *big.Int) (keys []ethkey.KeyV2, err error)
	FundingKeys() (keys []ethkey.KeyV2, err error)
	GetRoundRobinAddress(chainID *big.Int, addresses ...common.Address) (address common.Address, err error)

	GetState(id string) (ethkey.State, error)
	SetState(ethkey.State) error
	GetStatesForKeys([]ethkey.KeyV2) ([]ethkey.State, error)
	GetStatesForChain(chainID *big.Int) ([]ethkey.State, error)

	GetV1KeysAsV2(f DefaultEVMChainIDFunc) ([]ethkey.KeyV2, []ethkey.State, error)
}

type eth struct {
	*keyManager
	subscribers   [](chan struct{})
	subscribersMu *sync.RWMutex
}

var _ Eth = &eth{}

func newEthKeyStore(km *keyManager) *eth {
	return &eth{
		keyManager:    km,
		subscribers:   make([](chan struct{}), 0),
		subscribersMu: new(sync.RWMutex),
	}
}

func (ks *eth) Get(id string) (ethkey.KeyV2, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return ethkey.KeyV2{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *eth) GetAll() (keys []ethkey.KeyV2, _ error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, key := range ks.keyRing.Eth {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].Cmp(keys[j]) < 0 })
	return keys, nil
}

func (ks *eth) Create(chainID *big.Int) (ethkey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ethkey.KeyV2{}, ErrLocked
	}
	key, err := ethkey.NewV2()
	if err != nil {
		return ethkey.KeyV2{}, err
	}
	err = ks.add(key, chainID)
	if err != nil {
		return ethkey.KeyV2{}, err
	}
	ks.notify()
	return key, nil
}

func (ks *eth) Add(key ethkey.KeyV2, chainID *big.Int) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if _, found := ks.keyRing.Eth[key.ID()]; found {
		return fmt.Errorf("key with ID %s already exists", key.ID())
	}
	err := ks.add(key, chainID)
	if err != nil {
		return err
	}
	ks.notify()
	return nil
}

// EnsureKeys verifies whether the ETH keys have been seeded, if not, it creates them.
func (ks *eth) EnsureKeys(chainID *big.Int) (err error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}

	var sendDidExist bool
	var fundDidExist bool
	var sendingKey ethkey.KeyV2
	var fundingKey ethkey.KeyV2

	// check & setup sending key
	sendingKeys := ks.sendingKeys(chainID)
	if len(sendingKeys) > 0 {
		sendingKey = sendingKeys[0]
		sendDidExist = true
	} else {
		sendingKey, err = ethkey.NewV2()
		if err != nil {
			return err
		}
		err = ks.addEthKeyWithState(sendingKey, ethkey.State{EVMChainID: *utils.NewBig(chainID), IsFunding: false})
		if err != nil {
			return err
		}
		ks.logger.Infow("New sending address created", "address", sendingKey.Address.Hex(), "evmChainID", chainID)
	}

	// check & setup funding key
	fundingKeys := ks.fundingKeys()
	if len(fundingKeys) > 0 {
		fundingKey = fundingKeys[0]
		fundDidExist = true
	} else {
		fundingKey, err = ethkey.NewV2()
		if err != nil {
			return err
		}
		err = ks.addEthKeyWithState(fundingKey, ethkey.State{EVMChainID: *utils.NewBig(chainID), IsFunding: true})
		if err != nil {
			return err
		}
		ks.logger.Infow("New funding address created", "address", fundingKey.Address.Hex(), "evmChainID", chainID)
	}

	if !sendDidExist || !fundDidExist {
		ks.notify()
	}

	return nil
}

func (ks *eth) Import(keyJSON []byte, password string, chainID *big.Int) (ethkey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ethkey.KeyV2{}, ErrLocked
	}
	dKey, err := keystore.DecryptKey(keyJSON, password)
	if err != nil {
		return ethkey.KeyV2{}, errors.Wrap(err, "EthKeyStore#ImportKey failed to decrypt key")
	}
	key := ethkey.FromPrivateKey(dKey.PrivateKey)
	if _, found := ks.keyRing.Eth[key.ID()]; found {
		return ethkey.KeyV2{}, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	err = ks.add(key, chainID)
	if err != nil {
		return ethkey.KeyV2{}, errors.Wrap(err, "unable to add eth key")
	}
	ks.notify()
	return key, nil
}

func (ks *eth) Export(id string, password string) ([]byte, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return nil, err
	}
	return key.ToEncryptedJSON(password, ks.scryptParams)
}

func (ks *eth) Delete(id string) (ethkey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ethkey.KeyV2{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return ethkey.KeyV2{}, err
	}
	err = ks.safeRemoveKey(key, func(tx pg.Queryer) error {
		_, err2 := tx.Exec(`DELETE FROM eth_key_states WHERE address = $1`, key.Address)
		return err2
	})
	if err != nil {
		return ethkey.KeyV2{}, errors.Wrap(err, "unable to remove eth key")
	}
	ks.notify()
	return key, nil
}

func (ks *eth) SubscribeToKeyChanges() (ch chan struct{}, unsub func()) {
	ch = make(chan struct{}, 1)
	ks.subscribersMu.Lock()
	defer ks.subscribersMu.Unlock()
	ks.subscribers = append(ks.subscribers, ch)
	return ch, func() {
		ks.subscribersMu.Lock()
		defer ks.subscribersMu.Unlock()
		for i, sub := range ks.subscribers {
			if sub == ch {
				ks.subscribers = append(ks.subscribers[:i], ks.subscribers[i+1:]...)
				close(ch)
			}
		}
	}
}

func (ks *eth) SignTx(address common.Address, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	key, err := ks.getByID(address.Hex())
	if err != nil {
		return nil, err
	}
	signer := types.LatestSignerForChainID(chainID)
	return types.SignTx(tx, signer, key.ToEcdsaPrivKey())
}

// SendingKeys returns all sending keys for the given chain
// If chainID is nil, returns all sending keys for all chains
func (ks *eth) SendingKeys(chainID *big.Int) (sendingKeys []ethkey.KeyV2, err error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	return ks.sendingKeys(chainID), nil
}

func (ks *eth) FundingKeys() (fundingKeys []ethkey.KeyV2, err error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	return ks.fundingKeys(), nil
}

func (ks *eth) GetRoundRobinAddress(chainID *big.Int, whitelist ...common.Address) (common.Address, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return common.Address{}, ErrLocked
	}

	var keys []ethkey.KeyV2
	if len(whitelist) == 0 {
		keys = ks.sendingKeys(chainID)
	} else if len(whitelist) > 0 {
		for _, k := range ks.sendingKeys(chainID) {
			for _, addr := range whitelist {
				if addr == k.Address.Address() {
					keys = append(keys, k)
				}
			}
		}
	}

	if len(keys) == 0 {
		var err error
		if chainID == nil && len(whitelist) == 0 {
			err = errors.New("no sending keys available")
		} else if chainID == nil {
			err = errors.Errorf("no sending keys available that match whitelist: %v", whitelist)
		} else if len(whitelist) == 0 {
			err = errors.Errorf("no sending keys available for chain %s", chainID.String())
		} else {
			err = errors.Errorf("no sending keys available for chain %s that match whitelist: %v", chainID.String(), whitelist)
		}
		return common.Address{}, err
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return ks.keyStates.Eth[keys[i].ID()].LastUsed().Before(ks.keyStates.Eth[keys[j].ID()].LastUsed())
	})

	leastRecentlyUsed := keys[0]
	ks.keyStates.Eth[leastRecentlyUsed.ID()].WasUsed()
	return leastRecentlyUsed.Address.Address(), nil
}

func (ks *eth) GetState(id string) (ethkey.State, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return ethkey.State{}, ErrLocked
	}
	state, exists := ks.keyStates.Eth[id]
	if !exists {
		return ethkey.State{}, errors.Errorf("state not found for eth key ID %s", id)
	}
	return *state, nil
}

// SetState is only used in tests to manually update a key's state
func (ks *eth) SetState(state ethkey.State) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	_, exists := ks.keyStates.Eth[state.KeyID()]
	if !exists {
		return errors.Errorf("key not found with ID %s", state.KeyID())
	}
	ks.keyStates.Eth[state.KeyID()] = &state
	sql := `UPDATE eth_key_states SET address = :address, next_nonce = :next_nonce, is_funding = :is_funding, evm_chain_id = :evm_chain_id, updated_at = NOW()
	WHERE address = :address;`
	_, err := ks.orm.q.NamedExec(sql, state)
	return errors.Wrap(err, "SetState#Exec failed")
}

func (ks *eth) GetStatesForKeys(keys []ethkey.KeyV2) (states []ethkey.State, err error) {
	for _, k := range keys {
		state, err := ks.GetState(k.ID())
		if err != nil {
			return nil, err
		}
		states = append(states, state)
	}
	return
}

func (ks *eth) GetStatesForChain(chainID *big.Int) (states []ethkey.State, err error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, s := range ks.keyStates.Eth {
		if s.EVMChainID.Equal(utils.NewBig(chainID)) {
			states = append(states, *s)
		}
	}
	return
}

func (ks *eth) GetV1KeysAsV2(f DefaultEVMChainIDFunc) (keys []ethkey.KeyV2, states []ethkey.State, _ error) {
	v1Keys, err := ks.orm.GetEncryptedV1EthKeys()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get encrypted v1 eth keys")
	}
	if len(v1Keys) == 0 {
		return keys, states, nil
	}
	chainID, err := f()
	if err != nil {
		return nil, nil, errors.Wrapf(err, `%d legacy eth keys detected, but no default EVM chain ID was specified

PLEASE READ THIS ADDITIONAL INFO

If you are running Chainlink with EVM_ENABLED=false and don't care about EVM keys at all, you can run the following SQL to remove any lingering eth keys that may have been autogenerated by an older version of Chainlink, and boot the node again:

pqsl> TRUNCATE keys;

WARNING: This will PERMANENTLY AND IRRECOVERABLY delete any legacy eth keys, so please be absolutely sure this is what you want before you run this. Consider taking a database backup first.
`, len(v1Keys))
	}
	for _, keyV1 := range v1Keys {
		dKey, err := keystore.DecryptKey(keyV1.JSON, ks.password)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "could not decrypt eth key %s", keyV1.Address.Hex())
		}
		keyV2 := ethkey.FromPrivateKey(dKey.PrivateKey)
		keys = append(keys, keyV2)
		state := ethkey.State{
			Address:    keyV1.Address,
			NextNonce:  keyV1.NextNonce,
			IsFunding:  keyV1.IsFunding,
			EVMChainID: *utils.NewBig(chainID),
		}
		states = append(states, state)
	}
	return keys, states, nil
}

// caller must hold lock!
func (ks *eth) getByID(id string) (ethkey.KeyV2, error) {
	key, found := ks.keyRing.Eth[id]
	if !found {
		return ethkey.KeyV2{}, fmt.Errorf("unable to find eth key with id %s", id)
	}
	return key, nil
}

// caller must hold lock!
func (ks *eth) fundingKeys() (fundingKeys []ethkey.KeyV2) {
	for _, k := range ks.keyRing.Eth {
		if ks.keyStates.Eth[k.ID()].IsFunding {
			fundingKeys = append(fundingKeys, k)
		}
	}
	sort.Slice(fundingKeys, func(i, j int) bool { return fundingKeys[i].Cmp(fundingKeys[j]) < 0 })
	return fundingKeys
}

// caller must hold lock!
// if chainID is nil, returns keys for all chains
func (ks *eth) sendingKeys(chainID *big.Int) (sendingKeys []ethkey.KeyV2) {
	for _, k := range ks.keyRing.Eth {
		state := ks.keyStates.Eth[k.ID()]
		if !state.IsFunding && (chainID == nil || (((*big.Int)(&state.EVMChainID)).Cmp(chainID) == 0)) {
			sendingKeys = append(sendingKeys, k)
		}
	}
	sort.Slice(sendingKeys, func(i, j int) bool { return sendingKeys[i].Cmp(sendingKeys[j]) < 0 })
	return sendingKeys
}

// caller must hold lock!
func (ks *eth) add(key ethkey.KeyV2, chainID *big.Int) error {
	return ks.addEthKeyWithState(key, ethkey.State{EVMChainID: *utils.NewBig(chainID)})
}

// caller must hold lock!
func (ks *eth) addEthKeyWithState(key ethkey.KeyV2, state ethkey.State) error {
	state.Address = key.Address
	return ks.safeAddKey(key, func(tx pg.Queryer) error {
		sql := `INSERT INTO eth_key_states (address, next_nonce, is_funding, evm_chain_id, created_at, updated_at)
VALUES (:address, :next_nonce, :is_funding, :evm_chain_id, NOW(), NOW())
RETURNING *;`
		if err := ks.orm.q.GetNamed(sql, &state, state); err != nil {
			return errors.Wrap(err, "failed to insert eth_key_state")
		}
		ks.keyStates.Eth[key.ID()] = &state
		return nil
	})
}

// notify notifies subscribers that eth keys have changed
func (ks *eth) notify() {
	ks.subscribersMu.RLock()
	defer ks.subscribersMu.RUnlock()
	for _, ch := range ks.subscribers {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}
