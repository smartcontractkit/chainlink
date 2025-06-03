package por

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

// Delegate represents a POR (Proof of Reserve) delegate
type Delegate struct {
	lggr logger.Logger
}

// PORPluginService wraps the POR plugin factory and implements job.ServiceCtx
type PORPluginService struct {
	factory *por.PorReportingPluginFactory
	spec    *job.PORSpec
	lggr    logger.Logger
	
	// Plugin lifecycle state
	started bool
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewPORPluginService creates a new POR plugin service
func NewPORPluginService(factory *por.PorReportingPluginFactory, spec *job.PORSpec, lggr logger.Logger) *PORPluginService {
	ctx, cancel := context.WithCancel(context.Background())
	return &PORPluginService{
		factory: factory,
		spec:    spec,
		lggr:    lggr.Named("PORPluginService"),
		started: false,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start implements job.ServiceCtx
func (s *PORPluginService) Start(ctx context.Context) error {
	if s.started {
		return fmt.Errorf("POR plugin service already started")
	}
	
	s.lggr.Info("Starting POR plugin service")
	
	// In a real implementation, you would:
	// 1. Initialize the plugin with proper configuration
	// 2. Start any background processes for data collection
	// 3. Set up communication channels for reporting
	// 4. Begin the OCR protocol participation
	
	// For now, we just validate that the factory is properly configured
	if s.factory == nil {
		return fmt.Errorf("POR plugin factory is nil")
	}
	
	if s.factory.ExternalAdapter == nil {
		return fmt.Errorf("POR plugin external adapter is nil")
	}
	
	if s.factory.ContractReader == nil {
		return fmt.Errorf("POR plugin contract reader is nil")
	}
	
	if s.factory.ReportMarshaler == nil {
		return fmt.Errorf("POR plugin report marshaler is nil")
	}
	
	s.started = true
	s.lggr.Info("POR plugin service started successfully")
	return nil
}

// Close implements job.ServiceCtx  
func (s *PORPluginService) Close() error {
	if !s.started {
		return nil
	}
	
	s.lggr.Info("Closing POR plugin service")
	
	// Cancel any background operations
	s.cancel()
	
	// In a real implementation, you would:
	// 1. Gracefully shutdown any running OCR processes
	// 2. Clean up database connections
	// 3. Close network connections
	// 4. Persist any final state
	
	s.started = false
	s.lggr.Info("POR plugin service closed successfully")
	return nil
}

var _ job.ServiceCtx = (*PORPluginService)(nil)

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

	d.lggr.Info("Creating POR plugin services for job", "jobID", spec.ID, "jobName", spec.Name)

	// Instantiate the POR plugin factory (from porplugin_simple.go)
	// In a real implementation, you would wire up the actual dependencies here.
	// For this example, we use the mock implementations as in porplugin_simple.go.
	
	// Create a dummy config digest for the mock contract reader
	// In production, this would come from the actual OCR configuration
	dummyConfigDigest := types.ConfigDigest{}
	copy(dummyConfigDigest[:], []byte("dummy_config_digest_for_por_test"))
	
	// Initialize mock components
	externalAdapter := por.NewMockExternalAdapterImpl()
	contractReader := por.NewMockContractReader(dummyConfigDigest)
	reportMarshaler := por.NewMockReportMarshaler()
	
	// Create the plugin factory with proper configuration
	factory := &por.PorReportingPluginFactory{
		Logger:          nil, // TODO: Wire up proper libOCR logger
		ExternalAdapter: externalAdapter,
		ContractReader:  contractReader,
		ReportMarshaler: reportMarshaler,
	}

	// Wrap the plugin factory in a service (implement job.ServiceCtx)
	service := NewPORPluginService(factory, spec.PORSpec, d.lggr)

	d.lggr.Info("Successfully created POR plugin service", "jobID", spec.ID)
	return []job.ServiceCtx{service}, nil
}
