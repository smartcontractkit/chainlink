package types

import ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

type FunctionsProvider interface {
	PluginProvider
	FunctionsEvents() FunctionsEvents
}

type OracleRequest struct {
	RequestID          [32]byte
	RequestingContract ocrtypes.Account
	RequestInitiator   ocrtypes.Account
	//nolint:revive
	SubscriptionId      uint64
	SubscriptionOwner   ocrtypes.Account
	Data                []byte
	DataVersion         uint16
	Flags               [32]byte
	CallbackGasLimit    uint64
	TxHash              []byte
	CoordinatorContract ocrtypes.Account
	OnchainMetadata     []byte
}

type OracleResponse struct {
	RequestID [32]byte
}

// An on-chain event source, which understands router proxy contracts.
type FunctionsEvents interface {
	Service
	LatestEvents() ([]OracleRequest, []OracleResponse, error)
}
