package ocr2keeper

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"google.golang.org/protobuf/proto"

	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/static"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type AutomationCustomTelemetryService struct {
	utils.StartStopOnce
	monitoringEndpoint  commontypes.MonitoringEndpoint
	headBroadcaster     httypes.HeadBroadcaster
	headCh              chan blockKey
	customTelemChanSize uint8
	unsubscribe         func()
	chDone              chan struct{}
	lggr                logger.Logger
}

// NewAutomationCustomTelemetryService creates a telemetry service for new blocks and node version
func NewAutomationCustomTelemetryService(me commontypes.MonitoringEndpoint, hb httypes.HeadBroadcaster, lggr logger.Logger, customTelemChanSize uint8) *AutomationCustomTelemetryService {
	return &AutomationCustomTelemetryService{
		monitoringEndpoint:  me,
		headBroadcaster:     hb,
		headCh:              make(chan blockKey, customTelemChanSize),
		customTelemChanSize: customTelemChanSize,
		chDone:              make(chan struct{}),
		lggr:                lggr.Named("Automation Custom Telem"),
	}
}

// Start starts Custom Telemetry Service, sends 1 NodeVersion message to endpoint at start and sends new BlockNumber messages
func (e *AutomationCustomTelemetryService) Start(context.Context) error {
	return e.StartOnce("AutomationCustomTelemetryService", func() error {
		e.lggr.Infof("Starting: Custom Telemetry Service")
		versionMsg := &telem.NodeVersion{
			Timestamp:   uint64(time.Now().UTC().UnixMilli()),
			NodeVersion: static.Version,
		}
		wrappedMessage := &telem.AutomationTelemWrapper{
			Msg: &telem.AutomationTelemWrapper_NodeVersion{
				NodeVersion: versionMsg,
			},
		}
		bytes, err := proto.Marshal(wrappedMessage)
		if err != nil {
			e.lggr.Errorf("Error occurred while marshalling the message: %v", err)
		}
		e.monitoringEndpoint.SendLog(bytes)
		e.lggr.Infof("NodeVersion Message Sent to Endpoint: %d", versionMsg.Timestamp)
		_, e.unsubscribe = e.headBroadcaster.Subscribe(&headWrapper{e.headCh})
		go func() {
			e.lggr.Infof("Started: Custom Telemetry Service")
			for {
				select {
				case blockInfo := <-e.headCh:
					blockNumMsg := &telem.BlockNumber{
						Timestamp:   uint64(time.Now().UTC().UnixMilli()),
						BlockNumber: uint64(blockInfo.block),
						BlockHash:   blockInfo.hash,
					}
					wrappedMessage := &telem.AutomationTelemWrapper{
						Msg: &telem.AutomationTelemWrapper_BlockNumber{
							BlockNumber: blockNumMsg,
						},
					}
					bytes, err := proto.Marshal(wrappedMessage)
					if err != nil {
						e.lggr.Errorf("Error occurred while marshalling the message: %v", err)
					}
					e.monitoringEndpoint.SendLog(bytes)
					e.lggr.Infof("BlockNumber Message Sent to Endpoint: %d", blockNumMsg.Timestamp)
				case <-e.chDone:
					return
				}
			}
		}()
		return nil
	})
}

// Close stops go routines and closes channels
func (e *AutomationCustomTelemetryService) Close() error {
	return e.StopOnce("AutomationCustomTelemetryService", func() error {
		e.lggr.Infof("Stopping: custom telemetry service")
		e.unsubscribe()
		e.chDone <- struct{}{}
		close(e.headCh)
		close(e.chDone)
		e.lggr.Infof("Stopped: Custom telemetry service")
		return nil
	})
}

// blockKey contains block and hash info for BlockNumber telemetry message
type blockKey struct {
	block int64
	hash  string
}

// headWrapper is passed into HeadBroadcaster's subscribe() function, must implement OnNewLongestChain(_ context.Context, head *evmtypes.Head)
type headWrapper struct {
	headCh chan blockKey
}

// OnNewLongestChain sends block number and hash to head channel where message will be sent to monitoring endpoint
func (hw *headWrapper) OnNewLongestChain(_ context.Context, head *evmtypes.Head) {
	if head != nil {
		hw.headCh <- blockKey{
			block: head.Number,
			hash:  head.BlockHash().Hex(),
		}
	}
}
