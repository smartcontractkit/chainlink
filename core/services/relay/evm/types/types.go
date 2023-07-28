package types

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"

	"gopkg.in/guregu/null.v2"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type RelayConfig struct {
	ChainID                *utils.Big  `json:"chainID"`
	FromBlock              uint64      `json:"fromBlock"`
	EffectiveTransmitterID null.String `json:"effectiveTransmitterID"`

	// Contract-specific
	SendingKeys pq.StringArray `json:"sendingKeys"`

	// Mercury-specific
	FeedID *common.Hash `json:"feedID"`
}

type ConfigPoller interface {
	ocrtypes.ContractConfigTracker

	Start()
	Close() error
	Replay(ctx context.Context, fromBlock int64) error
}

// TODO(FUN-668): Move chain-agnostic types to Relayer
type FunctionsProvider interface {
	relaytypes.PluginProvider
	LogPollerWrapper() LogPollerWrapper
}

// A LogPoller wrapper that understands router proxy contracts
type LogPollerWrapper interface {
	relaytypes.Service
	LatestRoutes() (activeCoordinator common.Address, proposedCoordinator common.Address, err error)
}
