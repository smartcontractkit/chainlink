package types

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/types/automation"
)

//go:generate mockery --name Encoder --structname MockEncoder --srcpkg "github.com/smartcontractkit/chainlink-common/pkg/types/automation" --case underscore --filename encoder.generated.go

//go:generate mockery --name LogEventProvider --structname MockLogEventProvider --srcpkg "github.com/smartcontractkit/chainlink-common/pkg/types/automation" --case underscore --filename logeventprovider.generated.go

//go:generate mockery --name RecoverableProvider --structname MockRecoverableProvider --srcpkg "github.com/smartcontractkit/chainlink-common/pkg/types/automation" --case underscore --filename recoverableprovider.generated.go

//go:generate mockery --name ConditionalUpkeepProvider --structname MockConditionalUpkeepProvider --srcpkg "github.com/smartcontractkit/chainlink-common/pkg/types/automation" --case underscore --filename conditionalupkeepprovider.generated.go

//go:generate mockery --name PayloadBuilder --structname MockPayloadBuilder --srcpkg "github.com/smartcontractkit/chainlink-common/pkg/types/automation" --case underscore --filename payloadbuilder.generated.go

//go:generate mockery --name BlockSubscriber --structname MockBlockSubscriber --srcpkg "github.com/smartcontractkit/chainlink-common/pkg/types/automation" --case underscore --filename block_subscriber.generated.go

//go:generate mockery --name UpkeepStateUpdater --structname MockUpkeepStateUpdater --srcpkg "github.com/smartcontractkit/chainlink-common/pkg/types/automation" --case underscore --filename upkeep_state_updater.generated.go

type UpkeepTypeGetter func(automation.UpkeepIdentifier) UpkeepType
type WorkIDGenerator func(automation.UpkeepIdentifier, automation.Trigger) string

type RetryQueue interface {
	// Enqueue adds new items to the queue
	Enqueue(items ...RetryRecord) error
	// Dequeue returns the next n items in the queue, considering retry time schedules
	Dequeue(n int) ([]automation.UpkeepPayload, error)
}

type ProposalQueue interface {
	// Enqueue adds new items to the queue
	Enqueue(items ...automation.CoordinatedBlockProposal) error
	// Dequeue returns the next n items in the queue, considering retry time schedules
	Dequeue(t UpkeepType, n int) ([]automation.CoordinatedBlockProposal, error)
}

//go:generate mockery --name TransmitEventProvider --srcpkg "github.com/smartcontractkit/chainlink-automation/pkg/v3/types" --case underscore --filename transmit_event_provider.generated.go
type TransmitEventProvider interface {
	GetLatestEvents(context.Context) ([]automation.TransmitEvent, error)
}

//go:generate mockery --name Runnable --structname MockRunnable --srcpkg "github.com/smartcontractkit/chainlink-automation/pkg/v3/types" --case underscore --filename runnable.generated.go
type Runnable interface {
	// Can get results for a subset of payloads along with an error
	CheckUpkeeps(context.Context, ...automation.UpkeepPayload) ([]automation.CheckResult, error)
}

//go:generate mockery --name ResultStore --structname MockResultStore --srcpkg "github.com/smartcontractkit/chainlink-automation/pkg/v3/types" --case underscore --filename result_store.generated.go
type ResultStore interface {
	Add(...automation.CheckResult)
	Remove(...string)
	View() ([]automation.CheckResult, error)
}

//go:generate mockery --name Coordinator --structname MockCoordinator --srcpkg "github.com/smartcontractkit/chainlink-automation/pkg/v3/types" --case underscore --filename coordinator.generated.go
type Coordinator interface {
	PreProcess(_ context.Context, payloads []automation.UpkeepPayload) ([]automation.UpkeepPayload, error)

	Accept(automation.ReportedUpkeep) bool
	ShouldTransmit(automation.ReportedUpkeep) bool
	FilterResults([]automation.CheckResult) ([]automation.CheckResult, error)
	FilterProposals([]automation.CoordinatedBlockProposal) ([]automation.CoordinatedBlockProposal, error)
}

//go:generate mockery --name MetadataStore --structname MockMetadataStore --srcpkg "github.com/smartcontractkit/chainlink-automation/pkg/v3/types" --case underscore --filename metadatastore.generated.go
type MetadataStore interface {
	SetBlockHistory(automation.BlockHistory)
	GetBlockHistory() automation.BlockHistory

	AddProposals(proposals ...automation.CoordinatedBlockProposal)
	ViewProposals(utype UpkeepType) []automation.CoordinatedBlockProposal
	RemoveProposals(proposals ...automation.CoordinatedBlockProposal)

	Start(context.Context) error
	Close() error
}

//go:generate mockery --name Ratio --structname MockRatio --srcpkg "github.com/smartcontractkit/chainlink-automation/pkg/v3/types" --case underscore --filename ratio.generated.go
type Ratio interface {
	// OfInt should return n out of x such that n/x ~ r (ratio)
	OfInt(int) int
}
