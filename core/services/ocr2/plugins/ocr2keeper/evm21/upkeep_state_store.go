package evm

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const BlockKeySeparator = "|"

type upkeepState struct {
	payload  *ocr2keepers.UpkeepPayload
	state    *UpkeepState
	block    int64
	upkeepId string
}

// TODO: use the same type defined in keeper plugin after a new release is cut
type UpkeepState uint8

const (
	Performed UpkeepState = iota
	Eligible
)

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
	lggr       logger.Logger
}

// NewUpkeepStateStore creates a new state store. This is an initial version of this store. More improvements to come:
// TODO: AUTO-4027
func NewUpkeepStateStore(lggr logger.Logger) *UpkeepStateStore {
	return &UpkeepStateStore{
		statesByID: map[string]*upkeepState{},
		lggr:       lggr.Named("UpkeepStateStore"),
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
	arrs := strings.Split(string(pl.CheckBlock), BlockKeySeparator)
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
		u.lggr.Infof("upkeep %s is overridden, payload ID is %s, block is %d", s.upkeepId, s.payload.ID, s.block)
		return nil
	}

	u.statesByID[pl.ID] = state
	u.states = append(u.states, state)
	u.lggr.Infof("added new state with upkeep %s payload ID %s block %d", state.upkeepId, state.payload.ID, state.block)
	return nil
}
