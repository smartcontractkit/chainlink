package keystore

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	evmkeystore "github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// Eth is the external interface for EthKeyStore
type Eth interface {
	Get(ctx context.Context, id string) (ethkey.KeyV2, error)
	GetAll(ctx context.Context) ([]ethkey.KeyV2, error)
	Create(ctx context.Context, chainIDs ...*big.Int) (ethkey.KeyV2, error)
	Delete(ctx context.Context, id string) (ethkey.KeyV2, error)
	Import(ctx context.Context, keyJSON []byte, password string, chainIDs ...*big.Int) (ethkey.KeyV2, error)
	Export(ctx context.Context, id string, password string) ([]byte, error)

	Enable(ctx context.Context, address common.Address, chainID *big.Int) error
	Disable(ctx context.Context, address common.Address, chainID *big.Int) error
	Add(ctx context.Context, address common.Address, chainID *big.Int) error

	EnsureKeys(ctx context.Context, chainIDs ...*big.Int) error

	EnabledKeysForChain(ctx context.Context, chainID *big.Int) (keys []ethkey.KeyV2, err error)
	GetRoundRobinAddress(ctx context.Context, chainID *big.Int, addresses ...common.Address) (address common.Address, err error)
	CheckEnabled(ctx context.Context, address common.Address, chainID *big.Int) error

	GetState(ctx context.Context, id string, chainID *big.Int) (ethkey.State, error)
	GetStatesForKeys(ctx context.Context, keys []ethkey.KeyV2) ([]ethkey.State, error)
	GetStateForKey(ctx context.Context, key ethkey.KeyV2) (ethkey.State, error)
	EnabledAddressesForChain(ctx context.Context, chainID *big.Int) (addresses []common.Address, err error)
	GetResourceMutex(ctx context.Context, address common.Address) *evmkeystore.Mutex

	XXXTestingOnlySetState(ctx context.Context, keyState ethkey.State)
	XXXTestingOnlyAdd(ctx context.Context, key ethkey.KeyV2)
}

var _ loop.Keystore = &EthSigner{}

type EthSigner struct {
	Eth
	chainID *big.Int
}

func NewEthSigner(eth Eth, chainID *big.Int) *EthSigner {
	return &EthSigner{Eth: eth, chainID: chainID}
}

func (e *EthSigner) Accounts(ctx context.Context) (accounts []string, err error) {
	as, err := e.EnabledAddressesForChain(ctx, e.chainID)
	if err != nil {
		return nil, err
	}
	for _, a := range as {
		accounts = append(accounts, a.String())
	}
	return
}

func (e *EthSigner) Sign(ctx context.Context, account string, data []byte) (signed []byte, err error) {
	k, err := e.Get(ctx, account)
	if err != nil {
		return nil, err
	}
	// loopp spec requires passing nil hash to check existence of id
	if data == nil {
		return nil, nil
	}
	return k.Sign(data)
}

type eth struct {
	*keyManager
	keystateORM
	ds            sqlutil.DataSource
	resourceMutex map[common.Address]*evmkeystore.Mutex // ResourceMutex is an internal field and ought not be persisted to the database. Its main usage is to verify that the same key is not used for both TXMv1 and TXMv2 (usage in both TXMs will cause nonce drift and will lead to missing transactions). This functionality should be removed after we completely switch to TXMv2
}

// GetResourceMutex gets the resource mutex associates with the address if no resource mutex is found a new one is created
func (ks *eth) GetResourceMutex(ctx context.Context, address common.Address) *evmkeystore.Mutex {
	ks.lock.Lock()
	defer ks.lock.Unlock()

	if ks.resourceMutex == nil {
		ks.resourceMutex = make(map[common.Address]*evmkeystore.Mutex)
	}

	_, exists := ks.resourceMutex[address]
	if !exists {
		ks.resourceMutex[address] = &evmkeystore.Mutex{}
	}
	return ks.resourceMutex[address]
}

var _ Eth = &eth{}

func newEthKeyStore(km *keyManager, orm keystateORM, ds sqlutil.DataSource) *eth {
	return &eth{
		keystateORM: orm,
		keyManager:  km,
		ds:          ds,
	}
}

func (ks *eth) Get(ctx context.Context, id string) (ethkey.KeyV2, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return ethkey.KeyV2{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *eth) GetAll(ctx context.Context) (keys []ethkey.KeyV2, _ error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	return ks.getAll(ctx), nil
}

// caller must hold lock!
func (ks *eth) getAll(ctx context.Context) (keys []ethkey.KeyV2) {
	for _, key := range ks.keyRing.Eth {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].Cmp(keys[j]) < 0 })
	return
}

// Create generates a fresh new key and enables it for the given chain IDs
func (ks *eth) Create(ctx context.Context, chainIDs ...*big.Int) (ethkey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ethkey.KeyV2{}, ErrLocked
	}
	key, err := ethkey.NewV2()
	if err != nil {
		return ethkey.KeyV2{}, err
	}
	err = ks.add(ctx, key, chainIDs...)
	if err != nil {
		return ethkey.KeyV2{}, errors.Wrap(err, "unable to add eth key")
	}
	ks.logger.Infow("Created EVM key with ID "+key.Address.Hex(), "address", key.Address.Hex(), "evmChainIDs", chainIDs)
	return key, err
}

// EnsureKeys ensures that each chain has at least one key with a state
// linked to that chain. If a key and state exists for a chain but it is
// disabled, we do not enable it automatically here.
func (ks *eth) EnsureKeys(ctx context.Context, chainIDs ...*big.Int) (err error) {
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
		err = ks.add(ctx, newKey, chainID)
		if err != nil {
			return fmt.Errorf("failed to add key %s for chain %s: %w", newKey.Address, chainID, err)
		}
		ks.logger.Infow("Created EVM key with ID "+newKey.Address.Hex(), "address", newKey.Address.Hex(), "evmChainID", chainID)
	}

	return nil
}

func (ks *eth) Import(ctx context.Context, keyJSON []byte, password string, chainIDs ...*big.Int) (ethkey.KeyV2, error) {
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
		return ethkey.KeyV2{}, ErrKeyExists
	}
	err = ks.add(ctx, key, chainIDs...)
	if err != nil {
		return ethkey.KeyV2{}, errors.Wrap(err, "unable to add eth key")
	}
	return key, nil
}

func (ks *eth) Export(ctx context.Context, id string, password string) ([]byte, error) {
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

func (ks *eth) Add(ctx context.Context, address common.Address, chainID *big.Int) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	_, found := ks.keyRing.Eth[address.Hex()]
	if !found {
		return ErrKeyNotFound
	}
	return ks.addKey(ctx, nil, address, chainID)
}

// caller must hold lock!
// ds is optional, for transactions
func (ks *eth) addKey(ctx context.Context, ds sqlutil.DataSource, address common.Address, chainID *big.Int) error {
	if ds == nil {
		ds = ks.ds
	}
	state := new(ethkey.State)
	sql := `INSERT INTO evm.key_states (address, disabled, evm_chain_id, created_at, updated_at)
			VALUES ($1, false, $2, NOW(), NOW()) 
			RETURNING *;`

	if err := ds.GetContext(ctx, state, sql, address, chainID.String()); err != nil {
		return errors.Wrap(err, "failed to insert key_state")
	}
	// consider: do we really need a cache of the keystates?
	ks.keyStates.add(state)
	return nil
}

func (ks *eth) Enable(ctx context.Context, address common.Address, chainID *big.Int) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	_, found := ks.keyRing.Eth[address.Hex()]
	if !found {
		return ErrKeyNotFound
	}
	return ks.enable(ctx, address, chainID)
}

// caller must hold lock!
func (ks *eth) enable(ctx context.Context, address common.Address, chainID *big.Int) error {
	state := new(ethkey.State)
	sql := `INSERT INTO evm.key_states as key_states ("address", "evm_chain_id", "disabled", "created_at", "updated_at") VALUES ($1, $2, false, NOW(), NOW())
			ON CONFLICT ("address", "evm_chain_id") DO UPDATE SET "disabled" = false, "updated_at" = NOW() WHERE key_states."address" = $1 AND key_states."evm_chain_id" = $2
    		RETURNING *;`
	if err := ks.ds.GetContext(ctx, state, sql, address, chainID.String()); err != nil {
		return errors.Wrap(err, "failed to enable state")
	}

	if state.CreatedAt.Equal(state.UpdatedAt) {
		ks.keyStates.add(state)
	} else {
		ks.keyStates.enable(address, chainID, state.UpdatedAt)
	}
	return nil
}

func (ks *eth) Disable(ctx context.Context, address common.Address, chainID *big.Int) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	_, found := ks.keyRing.Eth[address.Hex()]
	if !found {
		return errors.Errorf("no key exists with ID %s", address.Hex())
	}
	return ks.disable(ctx, address, chainID)
}

func (ks *eth) disable(ctx context.Context, address common.Address, chainID *big.Int) error {
	state := new(ethkey.State)
	sql := `INSERT INTO evm.key_states as key_states ("address", "evm_chain_id", "disabled", "created_at", "updated_at") VALUES ($1, $2, true, NOW(), NOW())
			ON CONFLICT ("address", "evm_chain_id") DO UPDATE SET "disabled" = true, "updated_at" = NOW() WHERE key_states."address" = $1 AND key_states."evm_chain_id" = $2
			RETURNING *;`
	if err := ks.ds.GetContext(ctx, state, sql, address, chainID.String()); err != nil {
		return errors.Wrap(err, "failed to disable state")
	}

	if state.CreatedAt.Equal(state.UpdatedAt) {
		ks.keyStates.add(state)
	} else {
		ks.keyStates.disable(address, chainID, state.UpdatedAt)
	}
	return nil
}

func (ks *eth) Delete(ctx context.Context, id string) (ethkey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ethkey.KeyV2{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return ethkey.KeyV2{}, err
	}
	err = ks.safeRemoveKey(ctx, key, func(ds sqlutil.DataSource) error {
		_, err2 := ds.ExecContext(ctx, `DELETE FROM evm.key_states WHERE address = $1`, key.Address)
		return err2
	})
	if err != nil {
		return ethkey.KeyV2{}, errors.Wrap(err, "unable to remove eth key")
	}
	ks.keyStates.delete(key.Address)
	return key, nil
}

// EnabledKeysForChain returns all keys that are enabled for the given chain
func (ks *eth) EnabledKeysForChain(ctx context.Context, chainID *big.Int) (sendingKeys []ethkey.KeyV2, err error) {
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

func (ks *eth) GetRoundRobinAddress(ctx context.Context, chainID *big.Int, whitelist ...common.Address) (common.Address, error) {
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
			err = errors.Errorf("no sending keys available for chain %s that match whitelist: %v", chainID, whitelist)
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
func (ks *eth) CheckEnabled(ctx context.Context, address common.Address, chainID *big.Int) error {
	if utils.IsZero(address) {
		return errors.Errorf("empty address provided as input")
	}
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
		return errors.Errorf("no eth key exists with address %s", address.String())
	}
	states := ks.keyStates.KeyIDChainID[address.String()]
	state, exists := states[chainID.String()]
	if !exists {
		var chainIDs []string
		for cid, state := range states {
			if !state.Disabled {
				chainIDs = append(chainIDs, cid)
			}
		}
		return errors.Errorf("eth key with address %s exists but is has not been enabled for chain %s (enabled only for chain IDs: %s)", address, chainID.String(), strings.Join(chainIDs, ","))
	}
	if state.Disabled {
		var chainIDs []string
		for cid, state := range states {
			if !state.Disabled {
				chainIDs = append(chainIDs, cid)
			}
		}
		return errors.Errorf("eth key with address %s exists but is disabled for chain %s (enabled only for chain IDs: %s)", address.String(), chainID.String(), strings.Join(chainIDs, ","))
	}
	return nil
}

func (ks *eth) GetState(ctx context.Context, id string, chainID *big.Int) (ethkey.State, error) {
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

func (ks *eth) GetStatesForKeys(ctx context.Context, keys []ethkey.KeyV2) (states []ethkey.State, err error) {
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

// Useful to fetch the ChainID for a given key
func (ks *eth) GetStateForKey(ctx context.Context, key ethkey.KeyV2) (state ethkey.State, err error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	for _, state := range ks.keyStates.All {
		if state.KeyID() == key.ID() {
			return *state, err
		}
	}
	err = fmt.Errorf("no state found for key with id %s", key.ID())
	return
}
func (ks *eth) EnabledAddressesForChain(ctx context.Context, chainID *big.Int) (addresses []common.Address, err error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if chainID == nil {
		return nil, errors.New("chainID must be non-nil")
	}
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, s := range ks.keyStates.ChainIDKeyID[chainID.String()] {
		if !s.Disabled {
			evmAddress := s.Address.Address()
			addresses = append(addresses, evmAddress)
		}
	}
	return
}

// XXXTestingOnlySetState is only used in tests to manually update a key's state
func (ks *eth) XXXTestingOnlySetState(ctx context.Context, state ethkey.State) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		panic(ErrLocked)
	}
	existingState, exists := ks.keyStates.ChainIDKeyID[state.EVMChainID.String()][state.KeyID()]
	if !exists {
		panic("key not found with ID " + state.KeyID())
	}
	*existingState = state
	sql := `UPDATE evm.key_states SET address = :address, is_disabled = :is_disabled, evm_chain_id = :evm_chain_id, updated_at = NOW()
	WHERE address = :address;`
	_, err := ks.ds.NamedExecContext(ctx, sql, state)
	if err != nil {
		panic(err.Error())
	}
}

// XXXTestingOnlyAdd is only used in tests to manually add a key
func (ks *eth) XXXTestingOnlyAdd(ctx context.Context, key ethkey.KeyV2) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		panic(ErrLocked)
	}
	if _, found := ks.keyRing.Eth[key.ID()]; found {
		panic(fmt.Sprintf("key with ID %s already exists", key.ID()))
	}
	err := ks.add(ctx, key)
	if err != nil {
		panic(err.Error())
	}
}

// caller must hold lock!
func (ks *eth) getByID(id string) (ethkey.KeyV2, error) {
	key, found := ks.keyRing.Eth[id]
	if !found {
		return ethkey.KeyV2{}, ErrKeyNotFound
	}
	return key, nil
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
func (ks *eth) add(ctx context.Context, key ethkey.KeyV2, chainIDs ...*big.Int) (err error) {
	return ks.safeAddKey(ctx, key, func(tx sqlutil.DataSource) (serr error) {
		for _, chainID := range chainIDs {
			if serr = ks.addKey(ctx, tx, key.Address, chainID); serr != nil {
				return serr
			}
		}
		return nil
	})
}
