package main

import (
	"context"
	"fmt"
	"log"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

// MockLogger simulates a basic logger for demonstration
type MockLogger struct{}

func (m MockLogger) Named(name string) MockLogger { return m }
func (m MockLogger) Info(msg string, args ...interface{}) {
	fmt.Printf("INFO [%s]: %v\n", "POR", args)
}
func (m MockLogger) Error(msg string, args ...interface{}) {
	fmt.Printf("ERROR [%s]: %v\n", "POR", args)
}

// MockPORSpec simulates the PORSpec structure
type MockPORSpec struct {
	ContractAddress string
	ChainID         int64
}

// MockJob simulates a job structure
type MockJob struct {
	ID      int32
	Name    string
	PORSpec *MockPORSpec
}

// PORPluginService demonstrates our service implementation
type PORPluginService struct {
	factory *por.PorReportingPluginFactory
	spec    *MockPORSpec
	lggr    MockLogger
	started bool
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewPORPluginService(factory *por.PorReportingPluginFactory, spec *MockPORSpec, lggr MockLogger) *PORPluginService {
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

func (s *PORPluginService) Start(ctx context.Context) error {
	if s.started {
		return fmt.Errorf("POR plugin service already started")
	}
	
	s.lggr.Info("Starting POR plugin service")
	
	// Validate factory configuration
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

func (s *PORPluginService) Close() error {
	if !s.started {
		return nil
	}
	
	s.lggr.Info("Closing POR plugin service")
	s.cancel()
	s.started = false
	s.lggr.Info("POR plugin service closed successfully")
	return nil
}

// Delegate demonstrates our delegate implementation
type Delegate struct {
	lggr MockLogger
}

func NewDelegate(lggr MockLogger) *Delegate {
	return &Delegate{
		lggr: lggr.Named("POR"),
	}
}

func (d *Delegate) ServicesForSpec(ctx context.Context, spec MockJob) (*PORPluginService, error) {
	if spec.PORSpec == nil {
		return nil, fmt.Errorf("POR Delegate expects a *PORSpec to be present, got %v", spec)
	}

	d.lggr.Info("Creating POR plugin services for job", "jobID", spec.ID, "jobName", spec.Name)

	// Create a dummy config digest for the mock contract reader
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

	// Wrap the plugin factory in a service
	service := NewPORPluginService(factory, spec.PORSpec, d.lggr)

	d.lggr.Info("Successfully created POR plugin service", "jobID", spec.ID)
	return service, nil
}

func main() {
	fmt.Println("=== POR Plugin Service Demo ===")
	
	// Create a mock logger
	lggr := MockLogger{}
	
	// Create the delegate
	delegate := NewDelegate(lggr)
	
	// Create a mock job with POR spec
	job := MockJob{
		ID:   123,
		Name: "test-por-job",
		PORSpec: &MockPORSpec{
			ContractAddress: "0x1234567890123456789012345678901234567890",
			ChainID:         1,
		},
	}
	
	// Create services for the spec
	ctx := context.Background()
	service, err := delegate.ServicesForSpec(ctx, job)
	if err != nil {
		log.Fatalf("Failed to create services: %v", err)
	}
	
	fmt.Println("\n✅ Successfully created POR plugin service")
	
	// Start the service
	err = service.Start(ctx)
	if err != nil {
		log.Fatalf("Failed to start service: %v", err)
	}
	
	fmt.Println("✅ Successfully started POR plugin service")
	
	// Demonstrate that the service is running
	fmt.Println("🔄 POR plugin service is running...")
	
	// Stop the service
	err = service.Close()
	if err != nil {
		log.Fatalf("Failed to close service: %v", err)
	}
	
	fmt.Println("✅ Successfully closed POR plugin service")
	fmt.Println("\n=== Demo Complete ===")
}
