package ocrcommon

import (
	"github.com/smartcontractkit/chainlink/core/services/job"
	ocrtypes "github.com/smartcontractkit/libocr/commontypes"
	"go.uber.org/atomic"
)

// DependentOCRService is an OCR service which is gated on a dependency.
type DependentOCRService struct {
	ready      <-chan struct{}
	stop, done chan struct{}
	service    job.Service
	started    *atomic.Bool
	lggr       ocrtypes.Logger
}

// NewDependentOCRService creates a new service that waits for a signal on the ready chan.
func NewDependentOCRService(ready <-chan struct{}, service job.Service, lggr ocrtypes.Logger) *DependentOCRService {
	return &DependentOCRService{
		ready:   ready,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
		service: service,
		started: atomic.NewBool(false),
		lggr:    lggr,
	}
}

// Start starts the service asynchronously.
func (ds *DependentOCRService) Start() error {
	go ds.run()
	return nil
}

func (ds *DependentOCRService) run() {
	defer close(ds.done)
	select {
	case <-ds.ready:
		if err := ds.service.Start(); err != nil {
			ds.lggr.Error("unable to start service", ocrtypes.LogFields{"err": err})
		} else {
			ds.lggr.Info("started dependent ocr service", ocrtypes.LogFields{})
			ds.started.Store(true)
		}
	case <-ds.stop:
		// In the case we shutdown before detecting config.
	}
}

// Close closes the service synchronously.
func (ds *DependentOCRService) Close() error {
	if ds.started.Load() {
		// Assumes service close is synchronous
		ds.lggr.Info("closed dependent ocr service", ocrtypes.LogFields{})
		return ds.service.Close()
	}
	// If it hasn't started lets stop waiting for the deps
	close(ds.stop)
	<-ds.done
	return nil
}
