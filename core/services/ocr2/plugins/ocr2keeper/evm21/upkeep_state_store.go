package evm

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
)

const Separator = "|"

type upkeepState struct {
	payload *ocr2keepers.UpkeepPayload
	state   *UpkeepState
}

type UpkeepState uint8

const Performed UpkeepState = iota

type UpkeepStateReader interface {
	// SelectByID retrieves a single upkeep state
	SelectByID(ID string) (*ocr2keepers.UpkeepPayload, *UpkeepState, error)
	// SelectByBlock retrieves upkeep states at a specific block
	SelectByBlock(block int64) ([]*ocr2keepers.UpkeepPayload, []*UpkeepState, error)
	// SelectByBlockRange retrieves upkeep states within block range from start (inclusive) to end (exclusive)
	SelectByBlockRange(start, end int64) ([]*ocr2keepers.UpkeepPayload, []*UpkeepState, error)
	// SelectByUpkeepID retrieves upkeep states for an upkeep
	SelectByUpkeepID(upkeepId *big.Int) ([]*ocr2keepers.UpkeepPayload, []*UpkeepState, error)
	// SelectByUpkeepIDs retrieves upkeep states for provided upkeeps
	SelectByUpkeepIDs([]*big.Int) ([]*ocr2keepers.UpkeepPayload, []*UpkeepState, error)
}

type UpkeepStateUpdater interface {
	SetUpkeepState(ocr2keepers.UpkeepPayload, UpkeepState) error
}

type UpkeepStateStore struct {
	mu               sync.RWMutex
	statesByID       map[string]*upkeepState
	statesByBlock    map[int64][]*upkeepState
	statesByUpkeepID map[string][]*upkeepState
}

func NewUpkeepStateStore() *UpkeepStateStore {
	return &UpkeepStateStore{
		statesByID:       map[string]*upkeepState{},
		statesByBlock:    map[int64][]*upkeepState{},
		statesByUpkeepID: map[string][]*upkeepState{},
	}
}

func (u *UpkeepStateStore) SelectByID(ID string) (*ocr2keepers.UpkeepPayload, *UpkeepState, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	if s, ok := u.statesByID[ID]; ok {
		return s.payload, s.state, nil
	}
	return nil, nil, nil
}

func (u *UpkeepStateStore) SelectByBlock(block int64) ([]*ocr2keepers.UpkeepPayload, []*UpkeepState, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	var pl []*ocr2keepers.UpkeepPayload
	var us []*UpkeepState
	if state, ok := u.statesByBlock[block]; ok {
		for _, s := range state {
			pl = append(pl, s.payload)
			us = append(us, s.state)
		}
	}
	return pl, us, nil
}

func (u *UpkeepStateStore) SelectByBlockRange(start, end int64) ([]*ocr2keepers.UpkeepPayload, []*UpkeepState, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	var pl []*ocr2keepers.UpkeepPayload
	var us []*UpkeepState
	for i := start; i < end; i++ {
		pl1, us1, err := u.SelectByBlock(i)
		if err != nil {
			return nil, nil, err
		}
		pl = append(pl, pl1...)
		us = append(us, us1...)
	}
	return pl, us, nil
}

func (u *UpkeepStateStore) SelectByUpkeepID(upkeepId *big.Int) ([]*ocr2keepers.UpkeepPayload, []*UpkeepState, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	var pl []*ocr2keepers.UpkeepPayload
	var us []*UpkeepState
	if state, ok := u.statesByUpkeepID[upkeepId.String()]; ok {
		for _, s := range state {
			pl = append(pl, s.payload)
			us = append(us, s.state)
		}
	}
	return pl, us, nil
}

func (u *UpkeepStateStore) SelectByUpkeepIDs(upkeepIds []*big.Int) ([]*ocr2keepers.UpkeepPayload, []*UpkeepState, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	var pl []*ocr2keepers.UpkeepPayload
	var us []*UpkeepState
	for _, id := range upkeepIds {
		pl1, us1, err := u.SelectByUpkeepID(id)
		if err != nil {
			return nil, nil, err
		}
		pl = append(pl, pl1...)
		us = append(us, us1...)
	}
	return pl, us, nil
}

func (u *UpkeepStateStore) SetUpkeepState(pl ocr2keepers.UpkeepPayload, us UpkeepState) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	state := &upkeepState{
		payload: &pl,
		state:   &us,
	}
	s, ok := u.statesByID[pl.ID]
	if ok {
		s.payload = &pl
		s.state = &us
		return nil
	}
	u.statesByID[pl.ID] = state

	upkeepId := big.NewInt(0).SetBytes(pl.Upkeep.ID)
	res1, _ := u.statesByUpkeepID[upkeepId.String()]
	res1 = append(res1, state)
	u.statesByUpkeepID[upkeepId.String()] = res1

	arrs := strings.Split(string(pl.CheckBlock), Separator)
	if len(arrs) != 2 {
		return fmt.Errorf("check block %s is invalid for upkeep %s", pl.CheckBlock, upkeepId)
	}
	block, err := strconv.ParseInt(arrs[0], 10, 64)
	if err != nil {
		return err
	}
	res2, _ := u.statesByBlock[block]
	res2 = append(res2, state)
	u.statesByBlock[block] = res2
	return nil
}
