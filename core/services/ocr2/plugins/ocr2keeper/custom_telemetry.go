package ocr2keeper

import (
	"context"
	"fmt"
	"time"

	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/libocr/commontypes"
	"google.golang.org/protobuf/proto"
)

type AutomationCustomTelemetryService struct {
	utils.StartStopOnce
	monitoringEndpoint commontypes.MonitoringEndpoint
	headBroadcaster    httypes.HeadBroadcaster
	headCh             chan BlockKey
	unsubscribe        func()
	// chDone             chan struct{}
	// lggr               logger.Logger
}

func NewAutomationCustomTelemetryService(me commontypes.MonitoringEndpoint, hb httypes.HeadBroadcaster) *AutomationCustomTelemetryService {
	return &AutomationCustomTelemetryService{
		monitoringEndpoint: me,
		headBroadcaster:    hb,
		headCh:             make(chan BlockKey), // Channel Size matter?
		// chDone:             done,
		// lggr:               lggr,
	}
}

// Start starts
func (e *AutomationCustomTelemetryService) Start(context.Context) error {
	return e.StartOnce("AutomationCustomTelemetryService", func() error {
		// read from head channel
		//
		_, e.unsubscribe = e.headBroadcaster.Subscribe(&headWrapper{e.headCh})
		go func() {
			// e.lggr.Infof("Started enhanced telemetry service for job %d", e.job.ID)
			for blockKey := range e.headCh {
				// marshall protobuf message to bytes
				// proto.Marshal takes in a pointer to proto message struct
				blockNumMsg := &telem.BlockNumber{
					NodeId:      "faa",
					Timestamp:   uint64(time.Now().UTC().UnixMilli()),
					BlockNumber: uint64(blockKey.block),
					BlockHash:   blockKey.hash,
				}
				wrappedMessage := &telem.AutomationTelemWrapper{
					Msg: &telem.AutomationTelemWrapper_BlockNumber{
						BlockNumber: blockNumMsg,
					},
				}
				bytes, err := proto.Marshal(wrappedMessage)
				if err != nil {
					fmt.Printf("Error occured: %v", err)
				}
				e.monitoringEndpoint.SendLog(bytes)
				// case <-e.chDone:
				// 	return
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

// Subscribe and unsubscribe functions

type BlockKey struct {
	block int64
	hash  string
}

type headWrapper struct {
	headCh chan BlockKey
}

func (hw *headWrapper) OnNewLongestChain(_ context.Context, head *evmtypes.Head) {
	if head != nil {
		hw.headCh <- BlockKey{
			block: head.Number,
			hash:  head.BlockHash().Hex(),
		}
	}
}
