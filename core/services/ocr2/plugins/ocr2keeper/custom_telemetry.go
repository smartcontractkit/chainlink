package ocr2keeper

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/smartcontractkit/libocr/commontypes"
	"google.golang.org/protobuf/proto"

	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/static"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	customTelemChanSize = 100
)

type AutomationCustomTelemetryService struct {
	utils.StartStopOnce
	monitoringEndpoint commontypes.MonitoringEndpoint
	headBroadcaster    httypes.HeadBroadcaster
	headCh             chan blockKey
	unsubscribe        func()
	chDone             chan struct{}
	lggr               logger.Logger
	configDigest       [32]byte
	latestConfigDigest latestConfigDigestGetter
}

type latestConfigDigestGetter interface {
	LatestConfigDetails(opts *bind.CallOpts) (keeper_registry_wrapper2_0.LatestConfigDetails, error)
}

// NewAutomationCustomTelemetryService creates a telemetry service for new blocks and node version
func NewAutomationCustomTelemetryService(me commontypes.MonitoringEndpoint, hb httypes.HeadBroadcaster,
	lggr logger.Logger, latestcd latestConfigDigestGetter) *AutomationCustomTelemetryService {
	return &AutomationCustomTelemetryService{
		monitoringEndpoint: me,
		headBroadcaster:    hb,
		headCh:             make(chan blockKey, customTelemChanSize),
		chDone:             make(chan struct{}),
		lggr:               lggr.Named("Automation Custom Telem"),
		latestConfigDigest: latestcd,
	}
}

// Start starts Custom Telemetry Service, sends 1 NodeVersion message to endpoint at start and sends new BlockNumber messages
func (e *AutomationCustomTelemetryService) Start(ctx context.Context) error {
	return e.StartOnce("AutomationCustomTelemetryService", func() error {
		e.lggr.Infof("Starting: Custom Telemetry Service")
		callOpt := &bind.CallOpts{Context: ctx}
		configDetails, cdErr0 := e.latestConfigDigest.LatestConfigDetails(callOpt)
		if cdErr0 != nil {
			e.lggr.Errorf("Error occurred while getting newestConfigDetails for initialization %v", cdErr0)
		} else {
			e.configDigest = configDetails.ConfigDigest
		}
		go func() {
			ticker := time.NewTicker(1 * time.Minute)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					callOpt := &bind.CallOpts{Context: ctx}
					newConfigDetails, cdErr := e.latestConfigDigest.LatestConfigDetails(callOpt)
					if cdErr != nil {
						e.lggr.Errorf("Error occurred while getting newestConfigDetails  %v", cdErr)
						continue
					}
					newConfigDigest := newConfigDetails.ConfigDigest
					if newConfigDigest != e.configDigest {
						e.configDigest = newConfigDigest
					}
				case <-ctx.Done():
					return
				}
			}
		}()
		vMsg := &telem.NodeVersion{
			Timestamp:    uint64(time.Now().UTC().UnixMilli()),
			NodeVersion:  static.Version,
			ConfigDigest: e.configDigest[:],
		}
		wrappedVMsg := &telem.AutomationTelemWrapper{
			Msg: &telem.AutomationTelemWrapper_NodeVersion{
				NodeVersion: vMsg,
			},
		}
		bytes, err := proto.Marshal(wrappedVMsg)
		if err != nil {
			e.lggr.Errorf("Error occurred while marshalling the Node Version Message %s: %v", wrappedVMsg.String(), err)
		} else {
			e.monitoringEndpoint.SendLog(bytes)
			e.lggr.Infof("NodeVersion Message Sent to Endpoint: %d", vMsg.Timestamp)
		}
		_, e.unsubscribe = e.headBroadcaster.Subscribe(&headWrapper{e.headCh})
		go func() {
			e.lggr.Infof("Started: Sending BlockNumber Messages")
			for {
				select {
				case blockInfo := <-e.headCh:
					blockNumMsg := &telem.BlockNumber{
						Timestamp:    uint64(time.Now().UTC().UnixMilli()),
						BlockNumber:  uint64(blockInfo.block),
						BlockHash:    blockInfo.hash,
						ConfigDigest: e.configDigest[:],
					}
					wrappedBlockNumMsg := &telem.AutomationTelemWrapper{
						Msg: &telem.AutomationTelemWrapper_BlockNumber{
							BlockNumber: blockNumMsg,
						},
					}
					b, err := proto.Marshal(wrappedBlockNumMsg)
					if err != nil {
						e.lggr.Errorf("Error occurred while marshalling the Block Num Message %s: %v", wrappedBlockNumMsg.String(), err)
					} else {
						e.monitoringEndpoint.SendLog(b)
						e.lggr.Infof("BlockNumber Message Sent to Endpoint: %d", blockNumMsg.Timestamp)
					}
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
