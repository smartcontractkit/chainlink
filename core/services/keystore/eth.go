package keystore

import (
	"fmt"
	"math/big"
	"sort"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

//go:generate mockery --quiet --name Eth --output mocks/ --case=underscore

// Eth is the external interface for EthKeyStore
type Eth interface {
	Get(id string) (ethkey.KeyV2, error)
	GetAll() ([]ethkey.KeyV2, error)
	Create(chainIDs ...*big.Int) (ethkey.KeyV2, error)
	Delete(id string) (ethkey.KeyV2, error)
	Import(keyJSON []byte, password string, chainIDs ...*big.Int) (ethkey.KeyV2, error)
	Export(id string, password string) ([]byte, error)

	Enable(address common.Address, chainID *big.Int, qopts ...pg.QOpt) error
	Disable(address common.Address, chainID *big.Int, qopts ...pg.QOpt) error
	Reset(address common.Address, chainID *big.Int, nonce int64, qopts ...pg.QOpt) error

	GetNextNonce(address common.Address, chainID *big.Int, qopts ...pg.QOpt) (int64, error)
	IncrementNextNonce(address common.Address, chainID *big.Int, currentNonce int64, qopts ...pg.QOpt) error

	EnsureKeys(chainIDs ...*big.Int) error
	SubscribeToKeyChanges() (ch chan struct{}, unsub func())

	SignTx(fromAddress common.Address, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)

	EnabledKeysForChain(chainID *big.Int) (keys []ethkey.KeyV2, err error)
	GetRoundRobinAddress(chainID *big.Int, addresses ...common.Address) (address common.Address, err error)
	CheckEnabled(address common.Address, chainID *big.Int) error

	GetState(id string, chainID *big.Int) (ethkey.State, error)
	GetStatesForKeys([]ethkey.KeyV2) ([]ethkey.State, error)
	GetStatesForChain(chainID *big.Int) ([]ethkey.State, error)

	XXXTestingOnlySetState(ethkey.State)
	XXXTestingOnlyAdd(key ethkey.KeyV2)
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
	return ks.getAll(), nil
}

// caller must hold lock!
func (ks *eth) getAll() (keys []ethkey.KeyV2) {
	for _, key := range ks.keyRing.Eth {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].Cmp(keys[j]) < 0 })
	return
}

// Create generates a fresh new key and enables it for the given chain IDs
func (ks *eth) Create(chainIDs ...*big.Int) (ethkey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ethkey.KeyV2{}, ErrLocked
	}
	key, err := ethkey.NewV2()
	if err != nil {
		return ethkey.KeyV2{}, err
	}
	err = ks.add(key, chainIDs...)
	if err == nil {
		ks.notify()
	}
	ks.logger.Infow(fmt.Sprintf("Created EVM key with ID %s", key.Address.Hex()), "address", key.Address.Hex(), "evmChainIDs", chainIDs)
	return key, err
}

// EnsureKeys ensures that each chain has at least one key with a state
// linked to that chain. If a key and state exists for a chain but it is
// disabled, we do not enable it automatically here.
func (ks *eth) EnsureKeys(chainIDs ...*big.Int) (err error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}

	for _, chainID := range chainIDs {
		keys := ks.keysForChain(chainID, true)
		if len(keys) > 0 {
			continue
		}
		newKey, err := ethkey.NewV2()
		if err != nil {
			return err
		}
		err = ks.add(newKey, chainID)
		if err != nil {
			return err
		}
		ks.logger.Infow(fmt.Sprintf("Created EVM key with ID %s", newKey.Address.Hex()), "address", newKey.Address.Hex(), "evmChainID", chainID)
	}

	return nil
}

func (ks *eth) Import(keyJSON []byte, password string, chainIDs ...*big.Int) (ethkey.KeyV2, error) {
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
	err = ks.add(key, chainIDs...)
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

// Get the next nonce for the given key and chain. It is safest to always to go the DB for this
func (ks *eth) GetNextNonce(address common.Address, chainID *big.Int, qopts ...pg.QOpt) (nonce int64, err error) {
	if !ks.exists(address) {
		return 0, errors.Errorf("key with address %s does not exist", address.Hex())
	}
	nonce, err = ks.orm.getNextNonce(address, chainID, qopts...)
	if err != nil {
		return 0, errors.Wrap(err, "GetNextNonce failed")
	}
	ks.lock.Lock()
	defer ks.lock.Unlock()
	state, exists := ks.keyStates.KeyIDChainID[address.Hex()][chainID.String()]
	if !exists {
		return 0, errors.Errorf("state not found for address %s, chainID %s", address.Hex(), chainID.String())
	}
	if state.Disabled {
		return 0, errors.Errorf("state is disabled for address %s, chainID %s", address.Hex(), chainID.String())
	}
	// Always clobber the memory nonce with the DB nonce
	state.NextNonce = nonce
	return nonce, nil
}

// IncrementNextNonce increments keys.next_nonce by 1
func (ks *eth) IncrementNextNonce(address common.Address, chainID *big.Int, currentNonce int64, qopts ...pg.QOpt) error {
	if !ks.exists(address) {
		return errors.Errorf("key with address %s does not exist", address.Hex())
	}
	incrementedNonce, err := ks.orm.incrementNextNonce(address, chainID, currentNonce, qopts...)
	if err != nil {
		return errors.Wrap(err, "failed IncrementNextNonce")
	}
	ks.lock.Lock()
	defer ks.lock.Unlock()
	state, exists := ks.keyStates.KeyIDChainID[address.Hex()][chainID.String()]
	if !exists {
		return errors.Errorf("state not found for address %s, chainID %s", address.Hex(), chainID.String())
	}
	if state.Disabled {
		return errors.Errorf("state is disabled for address %s, chainID %s", address.Hex(), chainID.String())
	}
	state.NextNonce = incrementedNonce
	return nil
}

func (ks *eth) Enable(address common.Address, chainID *big.Int, qopts ...pg.QOpt) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	_, found := ks.keyRing.Eth[address.Hex()]
	if !found {
		return errors.Errorf("no key exists with ID %s", address.Hex())
	}
	return ks.enable(address, chainID, qopts...)
}

// caller must hold lock!
func (ks *eth) enable(address common.Address, chainID *big.Int, qopts ...pg.QOpt) error {
	state := new(ethkey.State)
	sql := `INSERT INTO evm_key_states (address, next_nonce, disabled, evm_chain_id, created_at, updated_at)
VALUES ($1, 0, false, $2, NOW(), NOW()) ON CONFLICT (evm_chain_id, address) DO UPDATE SET
disabled=false,
updated_at=NOW()
RETURNING id, next_nonce, address, evm_chain_id, disabled, created_at, updated_at;`
	q := ks.orm.q.WithOpts(qopts...)
	if err := q.Get(state, sql, address, chainID.String()); err != nil {
		return errors.Wrap(err, "failed to insert evm_key_state")
	}
	ks.keyStates.add(state)
	ks.notify()
	return nil
}

func (ks *eth) Disable(address common.Address, chainID *big.Int, qopts ...pg.QOpt) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	_, found := ks.keyRing.Eth[address.Hex()]
	if !found {
		return errors.Errorf("no key exists with ID %s", address.Hex())
	}
	return ks.disable(address, chainID, qopts...)
}

func (ks *eth) disable(address common.Address, chainID *big.Int, qopts ...pg.QOpt) error {
	q := ks.orm.q.WithOpts(qopts...)
	_, err := q.Exec(`UPDATE evm_key_states SET disabled = true WHERE address = $1 AND evm_chain_id = $2`, address, chainID.String())
	if err != nil {
		return errors.Wrap(err, "failed to disable state")
	}

	ks.keyStates.disable(address, chainID)
	ks.notify()
	return nil
}

// Reset the key/chain nonce to the given one
func (ks *eth) Reset(address common.Address, chainID *big.Int, nonce int64, qopts ...pg.QOpt) error {
	q := ks.orm.q.WithOpts(qopts...)
	res, err := q.Exec(`UPDATE evm_key_states SET next_nonce = $1 WHERE address = $2 AND evm_chain_id = $3`, nonce, address, chainID.String())
	if err != nil {
		return errors.Wrap(err, "failed to reset state")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Errorf("key state not found with address %s and chainID %s", address.Hex(), chainID.String())
	}
	ks.lock.Lock()
	defer ks.lock.Unlock()
	state, exists := ks.keyStates.KeyIDChainID[address.Hex()][chainID.String()]
	if !exists {
		return errors.Errorf("state not found for address %s, chainID %s", address.Hex(), chainID.String())
	}
	state.NextNonce = nonce
	return nil
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
		_, err2 := tx.Exec(`DELETE FROM evm_key_states WHERE address = $1`, key.Address)
		return err2
	})
	if err != nil {
		return ethkey.KeyV2{}, errors.Wrap(err, "unable to remove eth key")
	}
	ks.keyStates.delete(key.Address)
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

// EnabledKeysForChain returns all keys that are enabled for the given chain
func (ks *eth) EnabledKeysForChain(chainID *big.Int) (sendingKeys []ethkey.KeyV2, err error) {
	if chainID == nil {
		return nil, errors.New("chainID must be non-nil")
	}
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	return ks.enabledKeysForChain(chainID), nil
}

func (ks *eth) GetRoundRobinAddress(chainID *big.Int, whitelist ...common.Address) (common.Address, error) {
	if chainID == nil {
		return common.Address{}, errors.New("chainID must be non-nil")
	}
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return common.Address{}, ErrLocked
	}

	var keys []ethkey.KeyV2
	if len(whitelist) == 0 {
		keys = ks.enabledKeysForChain(chainID)
	} else if len(whitelist) > 0 {
		for _, k := range ks.enabledKeysForChain(chainID) {
			for _, addr := range whitelist {
				if addr == k.Address {
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

	states := ks.keyStates.ChainIDKeyID[chainID.String()]
	sort.SliceStable(keys, func(i, j int) bool {
		return states[keys[i].ID()].LastUsed().Before(states[keys[j].ID()].LastUsed())
	})

	leastRecentlyUsed := keys[0]
	states[leastRecentlyUsed.ID()].WasUsed()
	return leastRecentlyUsed.Address, nil
}

// CheckEnabled returns nil if state is present and enabled
// The complexity here comes because we want to return nice, useful error messages
func (ks *eth) CheckEnabled(address common.Address, chainID *big.Int) error {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return ErrLocked
	}
	var found bool
	for _, k := range ks.keyRing.Eth {
		if k.Address == address {
			found = true
			break
		}
	}
	if !found {
		return errors.Errorf("no eth key exists with address %s", address.Hex())
	}
	states := ks.keyStates.KeyIDChainID[address.Hex()]
	state, exists := states[chainID.String()]
	if !exists {
		var chainIDs []string
		for cid, state := range states {
			if !state.Disabled {
				chainIDs = append(chainIDs, cid)
			}
		}
		return errors.Errorf("eth key with address %s exists but is has not been enabled for chain %s (enabled only for chain IDs: %s)", address.Hex(), chainID.String(), strings.Join(chainIDs, ","))
	}
	if state.Disabled {
		var chainIDs []string
		for cid, state := range states {
			if !state.Disabled {
				chainIDs = append(chainIDs, cid)
			}
		}
		return errors.Errorf("eth key with address %s exists but is disabled for chain %s (enabled only for chain IDs: %s)", address.Hex(), chainID.String(), strings.Join(chainIDs, ","))
	}
	return nil
}

func (ks *eth) GetState(id string, chainID *big.Int) (ethkey.State, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return ethkey.State{}, ErrLocked
	}
	state, exists := ks.keyStates.KeyIDChainID[id][chainID.String()]
	if !exists {
		return ethkey.State{}, errors.Errorf("state not found for eth key ID %s", id)
	}
	return *state, nil
}

func (ks *eth) GetStatesForKeys(keys []ethkey.KeyV2) (states []ethkey.State, err error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	for _, state := range ks.keyStates.All {
		for _, k := range keys {
			if state.KeyID() == k.ID() {
				states = append(states, *state)
			}
		}
	}
	sort.Slice(states, func(i, j int) bool { return states[i].KeyID() < states[j].KeyID() })
	return
}

func (ks *eth) GetStatesForChain(chainID *big.Int) (states []ethkey.State, err error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, s := range ks.keyStates.ChainIDKeyID[chainID.String()] {
		states = append(states, *s)
	}
	sort.Slice(states, func(i, j int) bool { return states[i].KeyID() < states[j].KeyID() })
	return
}

func (ks *eth) getV1KeysAsV2() (keys []ethkey.KeyV2, nonces []int64, fundings []bool, _ error) {
	v1Keys, err := ks.orm.GetEncryptedV1EthKeys()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to get encrypted v1 eth keys")
	}
	if len(v1Keys) == 0 {
		return nil, nil, nil, nil
	}
	for _, keyV1 := range v1Keys {
		dKey, err := keystore.DecryptKey(keyV1.JSON, ks.password)
		if err != nil {
			return nil, nil, nil, errors.Wrapf(err, "could not decrypt eth key %s", keyV1.Address.Hex())
		}
		keyV2 := ethkey.FromPrivateKey(dKey.PrivateKey)
		keys = append(keys, keyV2)
		nonces = append(nonces, keyV1.NextNonce)
		fundings = append(fundings, keyV1.IsFunding)
	}
	return keys, nonces, fundings, nil
}

// XXXTestingOnlySetState is only used in tests to manually update a key's state
func (ks *eth) XXXTestingOnlySetState(state ethkey.State) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		panic(ErrLocked)
	}
	existingState, exists := ks.keyStates.ChainIDKeyID[state.EVMChainID.String()][state.KeyID()]
	if !exists {
		panic(fmt.Sprintf("key not found with ID %s", state.KeyID()))
	}
	*existingState = state
	sql := `UPDATE evm_key_states SET address = :address, next_nonce = :next_nonce, is_disabled = :is_disabled, evm_chain_id = :evm_chain_id, updated_at = NOW()
	WHERE address = :address;`
	_, err := ks.orm.q.NamedExec(sql, state)
	if err != nil {
		panic(err.Error())
	}
}

// XXXTestingOnlyAdd is only used in tests to manually add a key
func (ks *eth) XXXTestingOnlyAdd(key ethkey.KeyV2) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		panic(ErrLocked)
	}
	if _, found := ks.keyRing.Eth[key.ID()]; found {
		panic(fmt.Sprintf("key with ID %s already exists", key.ID()))
	}
	err := ks.add(key)
	if err != nil {
		panic(err.Error())
	}
}

// caller must hold lock!
func (ks *eth) getByID(id string) (ethkey.KeyV2, error) {
	key, found := ks.keyRing.Eth[id]
	if !found {
		return ethkey.KeyV2{}, fmt.Errorf("unable to find eth key with id %s", id)
	}
	return key, nil
}

func (ks *eth) exists(address common.Address) bool {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	_, found := ks.keyRing.Eth[address.Hex()]
	return found
}

// caller must hold lock!
func (ks *eth) enabledKeysForChain(chainID *big.Int) (keys []ethkey.KeyV2) {
	return ks.keysForChain(chainID, false)
}

// caller must hold lock!
func (ks *eth) keysForChain(chainID *big.Int, includeDisabled bool) (keys []ethkey.KeyV2) {
	states := ks.keyStates.ChainIDKeyID[chainID.String()]
	if states == nil {
		return
	}
	for keyID, state := range states {
		if includeDisabled || !state.Disabled {
			k := ks.keyRing.Eth[keyID]
			keys = append(keys, k)
		}
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].Cmp(keys[j]) < 0 })
	return keys
}

// caller must hold lock!
func (ks *eth) add(key ethkey.KeyV2, chainIDs ...*big.Int) (err error) {
	err = ks.safeAddKey(key, func(tx pg.Queryer) (serr error) {
		for _, chainID := range chainIDs {
			if serr = ks.enable(key.Address, chainID, pg.WithQueryer(tx)); serr != nil {
				return serr
			}
		}
		return nil
	})
	if len(chainIDs) > 0 {
		ks.notify()
	}
	return err
}

func (ks *eth) addWithNonce(key ethkey.KeyV2, chainID *big.Int, nonce int64, isDisabled bool) (err error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	err = ks.safeAddKey(key, func(tx pg.Queryer) (merr error) {
		state := new(ethkey.State)
		sql := `INSERT INTO evm_key_states (address, next_nonce, disabled, evm_chain_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING *;`
		if err = ks.orm.q.Get(state, sql, key.Address, nonce, isDisabled, chainID); err != nil {
			return errors.Wrap(err, "failed to insert evm_key_state")
		}

		ks.keyStates.add(state)
		return nil
	})
	ks.notify()
	return err
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
