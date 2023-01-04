package evm

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type UpkeepProviderV2_0 struct {
	addr     common.Address
	lookback int64
	ht       httypes.HeadTracker
	poller   logpoller.LogPoller
	chLog    chan logpoller.Log
	mu       sync.RWMutex
	active   map[int64]*managedUpkeep
	inactive map[int64]*managedUpkeep
}

func (up *UpkeepProviderV2_0) Register() error {
	// Add log filters for the log poller so that it can poll and find the logs that
	// we need
	_, err := up.poller.RegisterFilter(logpoller.Filter{
		EventSigs: upkeepStateChangeEventsV2_0,
		Addresses: []common.Address{up.addr},
	})

	return err
}

func (up *UpkeepProviderV2_0) Poll(ctx context.Context) error {
	end := up.calcLatestBlockNumber()

	{
		var logs []logpoller.Log
		var err error

		if logs, err = up.poller.LogsWithSigs(
			end-up.lookback,
			end,
			upkeepStateChangeEventsV2_0,
			up.addr,
			pg.WithParentCtx(ctx),
		); err != nil {
			return fmt.Errorf("%w: %s", ErrLogReadFailure, err)
		}

		for _, log := range logs {
			up.chLog <- log
		}
	}

	return nil
}

func (up *UpkeepProviderV2_0) calcLatestBlockNumber() int64 {
	var pollerEnd int64
	var trackerEnd int64

	pollerEnd, _ = up.poller.LatestBlock()

	ch := up.ht.LatestChain()
	if ch != nil {
		trackerEnd = ch.Number
	}

	if pollerEnd > trackerEnd {
		return pollerEnd
	}

	return trackerEnd
}

var upkeepStateChangeEventsV2_0 = []common.Hash{
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepMigrated{}.Topic(),   // removes upkeep id and detail from registry
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepRegistered{}.Topic(), // adds new upkeep id to registry
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepReceived{}.Topic(),   // adds multiple new upkeep ids to registry
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepCheckDataUpdated{}.Topic(),
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepGasLimitSet{}.Topic(),
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepCanceled{}.Topic(),
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepPaused{}.Topic(),
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepUnpaused{}.Topic(),
	keeper_registry_wrapper2_0.KeeperRegistryFundsAdded{}.Topic(),
}

type managedUpkeep struct {
	ID              *big.Int
	PerformGasLimit uint32
	CheckData       []byte
}
