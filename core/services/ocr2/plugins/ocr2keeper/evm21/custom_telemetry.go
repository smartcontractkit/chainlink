package evm

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/libocr/commontypes"
	"google.golang.org/protobuf/proto"
)

type AutomationCustomTelemetryService struct {
	utils.StartStopOnce

	chTelem <-chan telem.AutomationTelemWrapper
	// chDone             chan struct{}
	monitoringEndpoint commontypes.MonitoringEndpoint
	// lggr               logger.Logger
}

func NewAutomationCustomTelemetryService(chTelem <-chan telem.AutomationTelemWrapper, me commontypes.MonitoringEndpoint) *AutomationCustomTelemetryService {
	return &AutomationCustomTelemetryService{
		chTelem: chTelem,
		// chDone:             done,
		monitoringEndpoint: me,
		// lggr:               lggr,
	}
}

// Start starts
func (e *AutomationCustomTelemetryService) Start(context.Context) error {
	return e.StartOnce("AutomationCustomTelemetryService", func() error {
		go func() {
			// e.lggr.Infof("Started enhanced telemetry service for job %d", e.job.ID)
			for {
				select {
				case wrappedMessage := <-e.chTelem:
					// marshall protobuf message to bytes
					// proto.Marshal takes in a pointer to proto message struct
					bytes, err := proto.Marshal(&wrappedMessage)
					if err != nil {
						fmt.Printf("Error occured: %v", err)
					}
					e.monitoringEndpoint.SendLog(bytes)
					// case <-e.chDone:
					// 	return
				}
			}
		}()
		return nil
	})
}

func (e *AutomationCustomTelemetryService) Close() error {
	return e.StopOnce("AutomationCustomTelemetryService", func() error {
		// e.chDone <- struct{}{}
		// e.lggr.Infof("Stopping enhanced telemetry service for job %d", e.job.ID)
		return nil
	})
}
