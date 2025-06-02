package por

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// Delegate represents a POR (Proof of Reserve) delegate
type Delegate struct {
	lggr logger.Logger
}

var _ job.Delegate = (*Delegate)(nil)

// NewDelegate creates a new POR delegate
func NewDelegate(
	lggr logger.Logger,
) *Delegate {
	return &Delegate{
		lggr: lggr.Named("POR"),
	}
}

// JobType satisfies the job.Delegate interface.
func (d *Delegate) JobType() job.Type {
	return job.POR
}

// BeforeJobCreated satisfies the job.Delegate interface.
func (d *Delegate) BeforeJobCreated(spec job.Job) {}

// AfterJobCreated satisfies the job.Delegate interface.
func (d *Delegate) AfterJobCreated(spec job.Job) {}

// BeforeJobDeleted satisfies the job.Delegate interface.
func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

// OnDeleteJob satisfies the job.Delegate interface.
func (d *Delegate) OnDeleteJob(context.Context, job.Job) error { return nil }

// ServicesForSpec satisfies the job.Delegate interface.
func (d *Delegate) ServicesForSpec(ctx context.Context, spec job.Job) ([]job.ServiceCtx, error) {
	if spec.PORSpec == nil {
		return nil, fmt.Errorf("POR Delegate expects a *job.PORSpec to be present, got %v", spec)
	}

	// Instantiate the POR plugin factory (from porplugin_simple.go)
	// In a real implementation, you would wire up the actual dependencies here.
	// For this example, we use the mock implementations as in porplugin_simple.go.
	factory := &PorReportingPluginFactory{
		Logger:          nil, // TODO: Provide a real logger if needed
		ExternalAdapter: NewMockExternalAdapterImpl(),
		ContractReader:  NewMockContractReader(nil),
		ReportMarshaler: NewMockReportMarshaler(),
	}

	// Wrap the plugin factory in a service (implement job.ServiceCtx)
	service := NewPORPluginService(factory, spec.PORSpec)

	return []job.ServiceCtx{service}, nil
}
