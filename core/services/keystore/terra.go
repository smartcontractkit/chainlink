package keystore

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/terrakey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

//go:generate mockery --name Terra --output mocks/ --case=underscore

// Terra is the external interface for EthKeyStore
type Terra interface {
	Get(id string) (terrakey.KeyV2, error)
	GetAll() ([]terrakey.KeyV2, error)
	Create(chainID *big.Int) (terrakey.KeyV2, error)
	Add(key terrakey.KeyV2, chainID *big.Int) error
	Delete(id string) (terrakey.KeyV2, error)
	Import(keyJSON []byte, password string, chainID *big.Int) (terrakey.KeyV2, error)
	Export(id string, password string) ([]byte, error)

	EnsureKeys(chainID *big.Int) (terrakey.KeyV2, bool, terrakey.KeyV2, bool, error)
	SubscribeToKeyChanges() (ch chan struct{}, unsub func())

	SignTx(fromAddress common.Address, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)

	SendingKeys() (keys []terrakey.KeyV2, err error)
	FundingKeys() (keys []terrakey.KeyV2, err error)
	GetRoundRobinAddress(addresses ...common.Address) (address common.Address, err error)

	GetState(id string) (terrakey.State, error)
	SetState(terrakey.State) error
	GetStatesForKeys([]terrakey.KeyV2) ([]terrakey.State, error)
	GetStatesForChain(chainID *big.Int) ([]terrakey.State, error)

	GetV1KeysAsV2() ([]terrakey.KeyV2, []terrakey.State, error)
}

type terra struct {
	*keyManager
}

var _ Terra = terra{}

func newTerraKeyStore(km *keyManager) terra {
	return terra{
		km,
	}
}

func (ks terra) Get(id string) (terrakey.KeyV2, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return terrakey.KeyV2{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks terra) GetAll() (keys []terrakey.KeyV2, _ error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, key := range ks.keyRing.Terra {
		keys = append(keys, key)
	}
	return keys, nil
}

func (ks terra) Create(chainID *big.Int) (terrakey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return terrakey.KeyV2{}, ErrLocked
	}
	key, err := terrakey.NewV2()
	if err != nil {
		return terrakey.KeyV2{}, err
	}
	return key, ks.add(key, chainID)
}

func (ks terra) Add(key terrakey.KeyV2, chainID *big.Int) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if _, found := ks.keyRing.Terra[key.ID()]; found {
		return fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return ks.add(key, chainID)
}

func (ks terra) EnsureKeys(chainID *big.Int) (
	sendingKey terrakey.KeyV2,
	sendDidExist bool,
	fundingKey terrakey.KeyV2,
	fundDidExist bool,
	err error,
) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return terrakey.KeyV2{}, false, terrakey.KeyV2{}, false, ErrLocked
	}
	// check & setup sending key
	sendingKeys := ks.sendingKeys()
	if len(sendingKeys) > 0 {
		sendingKey = sendingKeys[0]
		sendDidExist = true
	} else {
		sendingKey, err = terrakey.NewV2()
		if err != nil {
			return terrakey.KeyV2{}, false, terrakey.KeyV2{}, false, err
		}
		err = ks.addEthKeyWithState(sendingKey, terrakey.State{EVMChainID: *utils.NewBig(chainID), IsFunding: false})
		if err != nil {
			return terrakey.KeyV2{}, false, terrakey.KeyV2{}, false, err
		}
	}
	// check & setup funding key
	fundingKeys := ks.fundingKeys()
	if err != nil {
		return terrakey.KeyV2{}, false, terrakey.KeyV2{}, false, err
	}
	if len(fundingKeys) > 0 {
		fundingKey = fundingKeys[0]
		fundDidExist = true
	} else {
		fundingKey, err = terrakey.NewV2()
		if err != nil {
			return terrakey.KeyV2{}, false, terrakey.KeyV2{}, false, err
		}
		err = ks.addEthKeyWithState(fundingKey, terrakey.State{EVMChainID: *utils.NewBig(chainID), IsFunding: true})
		if err != nil {
			return terrakey.KeyV2{}, false, terrakey.KeyV2{}, false, err
		}
	}
	return sendingKey, sendDidExist, fundingKey, fundDidExist, nil
}

func (ks terra) Import(keyJSON []byte, password string, chainID *big.Int) (terrakey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return terrakey.KeyV2{}, ErrLocked
	}
	key, err := terrakey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return terrakey.KeyV2{}, errors.Wrap(err, "EthKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.Terra[key.ID()]; found {
		return terrakey.KeyV2{}, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return key, ks.add(key, chainID)
}

func (ks terra) Export(id string, password string) ([]byte, error) {
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

func (ks terra) Delete(id string) (terrakey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return terrakey.KeyV2{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return terrakey.KeyV2{}, err
	}
	err = ks.safeRemoveKey(key, func(db *gorm.DB) error {
		return db.Where("address = ?", key.Address).Delete(terrakey.State{}).Error
	})
	return key, err
}

func (ks terra) SubscribeToKeyChanges() (ch chan struct{}, unsub func()) {
	return nil, func() {}
}

func (ks terra) SignTx(address common.Address, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
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

func (ks terra) SendingKeys() (sendingKeys []terrakey.KeyV2, err error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	return ks.sendingKeys(), nil
}

func (ks terra) FundingKeys() (fundingKeys []terrakey.KeyV2, err error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	return ks.fundingKeys(), nil
}

func (ks terra) GetRoundRobinAddress(whitelist ...common.Address) (common.Address, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return common.Address{}, ErrLocked
	}

	var keys []terrakey.KeyV2
	if len(whitelist) == 0 {
		keys = ks.sendingKeys()
	} else if len(whitelist) > 0 {
		for _, k := range ks.sendingKeys() {
			for _, addr := range whitelist {
				if addr == k.Address.Address() {
					keys = append(keys, k)
				}
			}
		}
	}

	if len(keys) == 0 {
		return common.Address{}, errors.New("no keys available")
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return ks.keyStates.Terra[keys[i].ID()].LastUsed().Before(ks.keyStates.Terra[keys[j].ID()].LastUsed())
	})

	leastRecentlyUsed := keys[0]
	ks.keyStates.Terra[leastRecentlyUsed.ID()].WasUsed()
	return leastRecentlyUsed.Address.Address(), nil
}

func (ks terra) GetState(id string) (terrakey.State, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return terrakey.State{}, ErrLocked
	}
	state, exists := ks.keyStates.Terra[id]
	if !exists {
		return terrakey.State{}, errors.Errorf("state not found for terra key ID %s", id)
	}
	return *state, nil
}

// SetState is only used in tests to manually update a key's state
func (ks terra) SetState(state terrakey.State) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	_, exists := ks.keyStates.Terra[state.KeyID()]
	if !exists {
		return errors.Errorf("key not found with ID %s", state.KeyID())
	}
	ks.keyStates.Terra[state.KeyID()] = &state
	return ks.orm.db.
		Model(terrakey.State{}).
		Where("address = ?", state.Address).
		Updates(state).Error
}

func (ks terra) GetStatesForKeys(keys []terrakey.KeyV2) (states []terrakey.State, err error) {
	for _, k := range keys {
		state, err := ks.GetState(k.ID())
		if err != nil {
			return nil, err
		}
		states = append(states, state)
	}
	return
}

func (ks terra) GetStatesForChain(chainID *big.Int) (states []terrakey.State, err error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, s := range ks.keyStates.Terra {
		if s.EVMChainID.Equal(utils.NewBig(chainID)) {
			states = append(states, *s)
		}
	}
	return
}

func (ks terra) GetV1KeysAsV2() (keys []terrakey.KeyV2, states []terrakey.State, _ error) {
	v1Keys, err := ks.orm.GetEncryptedV1TerraKeys()
	if err != nil {
		return keys, states, err
	}
	for _, keyV1 := range v1Keys {
		dKey, err := keystore.DecryptKey(keyV1.JSON, ks.password)
		if err != nil {
			return keys, states, err
		}
		keyV2 := terrakey.FromPrivateKey(dKey.PrivateKey)
		keys = append(keys, keyV2)
		state := terrakey.State{
			Address:   keyV1.Address,
			NextNonce: keyV1.NextNonce,
			IsFunding: keyV1.IsFunding,
		}
		states = append(states, state)
	}
	return keys, states, nil
}

// caller must hold lock!
func (ks terra) getByID(id string) (terrakey.KeyV2, error) {
	key, found := ks.keyRing.Terra[id]
	if !found {
		return terrakey.KeyV2{}, fmt.Errorf("unable to find terra key with id %s", id)
	}
	return key, nil
}

// caller must hold lock!
func (ks terra) fundingKeys() (fundingKeys []terrakey.KeyV2) {
	for _, k := range ks.keyRing.Terra {
		if ks.keyStates.Terra[k.ID()].IsFunding {
			fundingKeys = append(fundingKeys, k)
		}
	}
	sort.Slice(fundingKeys, func(i, j int) bool { return fundingKeys[i].Cmp(fundingKeys[j]) < 0 })
	return fundingKeys
}

// caller must hold lock!
func (ks terra) sendingKeys() (sendingKeys []terrakey.KeyV2) {
	for _, k := range ks.keyRing.Terra {
		if !ks.keyStates.Terra[k.ID()].IsFunding {
			sendingKeys = append(sendingKeys, k)
		}
	}
	sort.Slice(sendingKeys, func(i, j int) bool { return sendingKeys[i].Cmp(sendingKeys[j]) < 0 })
	return sendingKeys
}

// caller must hold lock!
func (ks terra) add(key terrakey.KeyV2, chainID *big.Int) error {
	return ks.addEthKeyWithState(key, terrakey.State{EVMChainID: *utils.NewBig(chainID)})
}

// caller must hold lock!
func (ks terra) addEthKeyWithState(key terrakey.KeyV2, state terrakey.State) error {
	state.Address = key.Address
	return ks.safeAddKey(key, func(db *gorm.DB) error {
		if err := db.Create(&state).Error; err != nil {
			return err
		}
		ks.keyStates.Terra[key.ID()] = &state
		return nil
	})
}
