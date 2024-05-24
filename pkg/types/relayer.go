package types

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/google/uuid"
)

const (
	NetworkEVM      = "evm"
	NetworkCosmos   = "cosmos"
	NetworkSolana   = "solana"
	NetworkStarkNet = "starknet"
)

var SupportedRelays = map[string]struct{}{
	NetworkEVM:      {},
	NetworkCosmos:   {},
	NetworkSolana:   {},
	NetworkStarkNet: {},
}

type RelayID struct {
	Network string
	ChainID string
}

// ID uniquely identifies a relayer by network and chain id
func (i *RelayID) Name() string {
	return fmt.Sprintf("%s.%s", i.Network, i.ChainID)
}

func (i *RelayID) String() string {
	return i.Name()
}
func NewRelayID(n string, c string) RelayID {
	return RelayID{Network: n, ChainID: c}
}

func (i *RelayID) UnmarshalString(s string) error {
	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		return fmt.Errorf("error unmarshaling Identifier. %s does not match expected pattern", s)
	}

	network, chainID := parts[0], parts[1]

	newID := &RelayID{ChainID: chainID}
	for n := range SupportedRelays {
		if network == n {
			newID.Network = n
			break
		}
	}

	if newID.Network == "" {
		return fmt.Errorf("error unmarshaling identifier: did not find network in supported list %q", newID.Network)
	}

	i.ChainID = newID.ChainID
	i.Network = newID.Network
	return nil
}

// PluginArgs are the args required to create any OCR2 plugin components.
// It's possible that the plugin config might actually be different
// per relay type, so we pass the config directly through.
type PluginArgs struct {
	TransmitterID string
	PluginConfig  []byte
}

type RelayArgs struct {
	ExternalJobID      uuid.UUID
	JobID              int32
	ContractID         string
	New                bool // Whether this is a first time job add.
	RelayConfig        []byte
	ProviderType       string
	MercuryCredentials *MercuryCredentials
}

type MercuryCredentials struct {
	LegacyURL string
	URL       string
	Username  string
	Password  string
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
	NewContractReader(contractReaderConfig []byte) (ContractReader, error)
	NewConfigProvider(rargs RelayArgs) (ConfigProvider, error)
	NewMedianProvider(rargs RelayArgs, pargs PluginArgs) (MedianProvider, error)
	NewMercuryProvider(rargs RelayArgs, pargs PluginArgs) (MercuryProvider, error)
	NewFunctionsProvider(rargs RelayArgs, pargs PluginArgs) (FunctionsProvider, error)
	NewAutomationProvider(rargs RelayArgs, pargs PluginArgs) (AutomationProvider, error)
	NewLLOProvider(rargs RelayArgs, pargs PluginArgs) (LLOProvider, error)
	NewPluginProvider(rargs RelayArgs, pargs PluginArgs) (PluginProvider, error)
	NewOCR3CapabilityProvider(rargs RelayArgs, pargs PluginArgs) (OCR3CapabilityProvider, error)
}
