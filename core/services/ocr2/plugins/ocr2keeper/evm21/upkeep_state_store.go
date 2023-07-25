package evm

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"

	mapset "github.com/deckarep/golang-set/v2"
)

const Separator = "|"

type upkeepState struct {
	payload  *ocr2keepers.UpkeepPayload
	state    *UpkeepState
	block    int64
	upkeepId string
}

type UpkeepState uint8

const Performed UpkeepState = iota

type UpkeepStateReader interface {
	// SelectByUpkeepIDsAndBlockRange retrieves upkeep states for provided upkeep ids and block range, the result is currently not in particular order
	SelectByUpkeepIDsAndBlockRange(upkeepIds []*big.Int, start, end int64) ([]*ocr2keepers.UpkeepPayload, []*UpkeepState, error)
}

type UpkeepStateUpdater interface {
	SetUpkeepState(ocr2keepers.UpkeepPayload, UpkeepState) error
}

type UpkeepStateStore struct {
	mu         sync.RWMutex
	statesByID map[string]*upkeepState
	states     []*upkeepState
}

// NewUpkeepStateStore creates a new state store. This is an initial version of this store. More improvements to come:
// TODO: AUTO-4027
func NewUpkeepStateStore() *UpkeepStateStore {
	return &UpkeepStateStore{
		statesByID: map[string]*upkeepState{},
	}
}

func (u *UpkeepStateStore) SelectByUpkeepIDsAndBlockRange(upkeepIds []*big.Int, start, end int64) ([]*ocr2keepers.UpkeepPayload, []*UpkeepState, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	var pl []*ocr2keepers.UpkeepPayload
	var us []*UpkeepState

	uids := mapset.NewSet[string]()
	for _, uid := range upkeepIds {
		uids.Add(uid.String())
	}

	for _, s := range u.states {
		if s.block < start || s.block >= end || !uids.Contains(s.upkeepId) {
			continue
		}
		pl = append(pl, s.payload)
		us = append(us, s.state)
	}
	return pl, us, nil
}

func (u *UpkeepStateStore) SetUpkeepState(pl ocr2keepers.UpkeepPayload, us UpkeepState) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	upkeepId := big.NewInt(0).SetBytes(pl.Upkeep.ID)
	arrs := strings.Split(string(pl.CheckBlock), Separator)
	if len(arrs) != 2 {
		return fmt.Errorf("check block %s is invalid for upkeep %s", pl.CheckBlock, upkeepId)
	}
	block, err := strconv.ParseInt(arrs[0], 10, 64)
	if err != nil {
		return err
	}
	state := &upkeepState{
		payload:  &pl,
		state:    &us,
		block:    block,
		upkeepId: upkeepId.String(),
	}

	s, ok := u.statesByID[pl.ID]
	if ok {
		s.payload = state.payload
		s.state = state.state
		s.block = state.block
		s.upkeepId = state.upkeepId
		u.statesByID[pl.ID] = s
		return nil
	}

	u.statesByID[pl.ID] = state
	u.states = append(u.states, state)
	return nil
}
