package automation

import (
	"context"
	"io"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

// UpkeepStateStore is the interface for managing upkeeps final state in a local store.
type UpkeepStateStore interface {
	UpkeepStateUpdater
	UpkeepStateReader
	Start(context.Context) error
	io.Closer
}

type Registry interface {
	CheckUpkeeps(ctx context.Context, keys ...UpkeepPayload) ([]CheckResult, error)
	Name() string
	Start(ctx context.Context) error
	Close() error
	HealthReport() map[string]error
}

type EventProvider interface {
	services.Service
	GetLatestEvents(ctx context.Context) ([]TransmitEvent, error)
}

type LogRecoverer interface {
	RecoverableProvider
	GetProposalData(context.Context, CoordinatedBlockProposal) ([]byte, error)

	Start(context.Context) error
	io.Closer
}

// UpkeepStateReader is the interface for reading the current state of upkeeps.
type UpkeepStateReader interface {
	SelectByWorkIDs(ctx context.Context, workIDs ...string) ([]UpkeepState, error)
}

type Encoder interface {
	Encode(...CheckResult) ([]byte, error)
	Extract([]byte) ([]ReportedUpkeep, error)
}

type LogEventProvider interface {
	GetLatestPayloads(context.Context) ([]UpkeepPayload, error)
	SetConfig(LogEventProviderConfig)
	Start(context.Context) error
	Close() error
}

type LogEventProviderConfig struct {
	NumLogUpkeeps    uint32
	FastExecLogsHigh uint32
	FastExecLogsLow  uint32
}

type RecoverableProvider interface {
	GetRecoveryProposals(context.Context) ([]UpkeepPayload, error)
}

type ConditionalUpkeepProvider interface {
	GetActiveUpkeeps(context.Context) ([]UpkeepPayload, error)
}

type PayloadBuilder interface {
	// Can get payloads for a subset of proposals along with an error
	BuildPayloads(context.Context, ...CoordinatedBlockProposal) ([]UpkeepPayload, error)
}

type BlockSubscriber interface {
	// Subscribe provides an identifier integer, a new channel, and potentially an error
	Subscribe() (int, chan BlockHistory, error)
	// Unsubscribe requires an identifier integer and indicates the provided channel should be closed
	Unsubscribe(int) error
	Start(context.Context) error
	Close() error
}

type UpkeepStateUpdater interface {
	SetUpkeepState(context.Context, CheckResult, UpkeepState) error
}
