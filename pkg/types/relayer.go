package types

import (
	"context"
	"math/big"

	"github.com/google/uuid"
)

// PluginArgs are the args required to create any OCR2 plugin components.
// It's possible that the plugin config might actually be different
// per relay type, so we pass the config directly through.
type PluginArgs struct {
	TransmitterID string
	PluginConfig  []byte
}

type RelayArgs struct {
	ExternalJobID uuid.UUID
	JobID         int32
	ContractID    string
	New           bool // Whether this is a first time job add.
	RelayConfig   []byte
}

type ChainStatus struct {
	ID      string
	Enabled bool
	Config  string // TOML
}

type NodeStatus struct {
	ChainID string
	Name    string
	Config  string // TOML
	State   string
}

// ChainService is a sub-interface of [loop.Relayer] that encapsulates the explicit interactions with a chain
type ChainService interface {
	Service

	GetChainStatus(ctx context.Context) (ChainStatus, error)
	ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) (stats []NodeStatus, nextPageToken string, total int, err error)
	Transact(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error
}

// Deprecated: use loop.Relayer, which includes context.Context.
type Relayer interface {
	Service
	NewConfigProvider(rargs RelayArgs) (ConfigProvider, error)
	NewMedianProvider(rargs RelayArgs, pargs PluginArgs) (MedianProvider, error)
	NewMercuryProvider(rargs RelayArgs, pargs PluginArgs) (MercuryProvider, error)
	NewFunctionsProvider(rargs RelayArgs, pargs PluginArgs) (FunctionsProvider, error)
}
