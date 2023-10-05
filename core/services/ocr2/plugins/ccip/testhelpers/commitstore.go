package testhelpers

import (
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	mock_contracts "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type FakeCommitStore struct {
	*mock_contracts.CommitStoreInterface

	isPaused     bool
	blessedRoots map[[32]byte]bool
	staticCfg    commit_store.CommitStoreStaticConfig
	dynamicCfg   commit_store.CommitStoreDynamicConfig
	nextSeqNum   uint64

	mu sync.RWMutex
}

func NewFakeCommitStore(t *testing.T, nextSeqNum uint64) (*FakeCommitStore, common.Address) {
	addr := utils.RandomAddress()
	//mockCommitStore := mock_contracts.NewCommitStoreInterface(t)
	mockCommitStore := mock_contracts.NewCommitStoreInterface(t)
	mockCommitStore.On("Address").Return(addr).Maybe()

	commitStore := &FakeCommitStore{CommitStoreInterface: mockCommitStore}
	commitStore.SetPaused(false)
	commitStore.SetNextSequenceNumber(nextSeqNum)

	return commitStore, addr
}

func (cs *FakeCommitStore) SetPaused(isPaused bool) {
	setCommitStoreVal(cs, func(cs *FakeCommitStore) { cs.isPaused = isPaused })
}

func (cs *FakeCommitStore) IsUnpausedAndARMHealthy(opts *bind.CallOpts) (bool, error) {
	return getCommitStoreVal(cs, func(cs *FakeCommitStore) bool { return !cs.isPaused }), nil
}

func (cs *FakeCommitStore) SetBlessedRoots(roots map[[32]byte]bool) {
	setCommitStoreVal(cs, func(cs *FakeCommitStore) { cs.blessedRoots = roots })
}

func (cs *FakeCommitStore) IsBlessed(opts *bind.CallOpts, root [32]byte) (bool, error) {
	return getCommitStoreVal(cs, func(cs *FakeCommitStore) bool { return cs.blessedRoots[root] }), nil
}

func (cs *FakeCommitStore) SetStaticConfig(cfg commit_store.CommitStoreStaticConfig) {
	setCommitStoreVal(cs, func(cs *FakeCommitStore) { cs.staticCfg = cfg })
}

func (cs *FakeCommitStore) GetStaticConfig(opts *bind.CallOpts) (commit_store.CommitStoreStaticConfig, error) {
	return getCommitStoreVal(cs, func(cs *FakeCommitStore) commit_store.CommitStoreStaticConfig { return cs.staticCfg }), nil
}

func (cs *FakeCommitStore) SetDynamicConfig(cfg commit_store.CommitStoreDynamicConfig) {
	setCommitStoreVal(cs, func(cs *FakeCommitStore) { cs.dynamicCfg = cfg })
}

func (cs *FakeCommitStore) GetDynamicConfig(opts *bind.CallOpts) (commit_store.CommitStoreDynamicConfig, error) {
	return getCommitStoreVal(cs, func(cs *FakeCommitStore) commit_store.CommitStoreDynamicConfig { return cs.dynamicCfg }), nil
}

func (cs *FakeCommitStore) SetNextSequenceNumber(seqNum uint64) {
	setCommitStoreVal(cs, func(cs *FakeCommitStore) { cs.nextSeqNum = seqNum })
}

func (cs *FakeCommitStore) GetExpectedNextSequenceNumber(opts *bind.CallOpts) (uint64, error) {
	return getCommitStoreVal(cs, func(cs *FakeCommitStore) uint64 { return cs.nextSeqNum }), nil
}

func getCommitStoreVal[T any](cs *FakeCommitStore, getter func(cs *FakeCommitStore) T) T {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return getter(cs)
}

func setCommitStoreVal(cs *FakeCommitStore, setter func(cs *FakeCommitStore)) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	setter(cs)
}
