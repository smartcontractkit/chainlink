package autotelemetry21

import (
	"context"
	"encoding/hex"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/static"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type AutomationCustomTelemetryService struct {
	services.StateMachine
	monitoringEndpoint    commontypes.MonitoringEndpoint
	blockSubscriber       ocr2keepers.BlockSubscriber
	blockSubChanID        int
	threadCtrl            utils.ThreadControl
	lggr                  logger.Logger
	configDigest          [32]byte
	contractConfigTracker types.ContractConfigTracker
	mu                    sync.RWMutex
}

// NewAutomationCustomTelemetryService creates a telemetry service for new blocks and node version
func NewAutomationCustomTelemetryService(me commontypes.MonitoringEndpoint,
	lggr logger.Logger, blocksub ocr2keepers.BlockSubscriber, configTracker types.ContractConfigTracker) (*AutomationCustomTelemetryService, error) {
	return &AutomationCustomTelemetryService{
		monitoringEndpoint:    me,
		threadCtrl:            utils.NewThreadControl(),
		lggr:                  lggr.Named("AutomationCustomTelem"),
		contractConfigTracker: configTracker,
		blockSubscriber:       blocksub,
	}, nil
}

// Start starts Custom Telemetry Service, sends 1 NodeVersion message to endpoint at start and sends new BlockNumber messages
func (e *AutomationCustomTelemetryService) Start(ctx context.Context) error {
	return e.StartOnce("AutomationCustomTelemetryService", func() error {
		e.lggr.Infof("Starting: Custom Telemetry Service")
		_, configDetails, err := e.contractConfigTracker.LatestConfigDetails(ctx)
		if err != nil {
			e.lggr.Errorf("Error occurred while getting newestConfigDetails for initialization %s", err)
		} else {
			e.configDigest = configDetails
			e.sendNodeVersionMsg()
		}
		e.threadCtrl.Go(func(ctx context.Context) {
			minuteTicker := time.NewTicker(1 * time.Minute)
			hourTicker := time.NewTicker(1 * time.Hour)
			defer minuteTicker.Stop()
			defer hourTicker.Stop()
			for {
				select {
				case <-minuteTicker.C:
					_, newConfigDigest, err := e.contractConfigTracker.LatestConfigDetails(ctx)
					if err != nil {
						e.lggr.Errorf("Error occurred while getting newestConfigDetails in configDigest loop %s", err)
					}
					configChanged := false
					e.mu.Lock()
					if newConfigDigest != e.configDigest {
						e.configDigest = newConfigDigest
						configChanged = true
					}
					e.mu.Unlock()

					if configChanged {
						e.sendNodeVersionMsg()
					}
				case <-hourTicker.C:
					e.sendNodeVersionMsg()
				case <-ctx.Done():
					return
				}
			}
		})

		chanID, blockSubscriberChan, blockSubErr := e.blockSubscriber.Subscribe()
		if blockSubErr != nil {
			e.lggr.Errorf("Block Subscriber Error: Subscribe(): %s", blockSubErr)
			return blockSubErr
		}
		e.blockSubChanID = chanID
		e.threadCtrl.Go(func(ctx context.Context) {
			e.lggr.Debug("Started: Sending BlockNumber Messages")
			for {
				select {
				case blockHistory := <-blockSubscriberChan:
					// Exploratory: Debounce blocks to avoid overflow in case of re-org
					latestBlockKey, err := blockHistory.Latest()
					if err != nil {
						e.lggr.Errorf("BlockSubscriber BlockHistory.Latest() failed: %s", err)
						continue
					}
					e.sendBlockNumberMsg(latestBlockKey)
				case <-ctx.Done():
					return
				}
			}
		})
		return nil
	})
}

// Close stops go routines and closes channels
func (e *AutomationCustomTelemetryService) Close() error {
	return e.StopOnce("AutomationCustomTelemetryService", func() error {
		e.lggr.Debug("Stopping: custom telemetry service")
		e.threadCtrl.Close()
		err := e.blockSubscriber.Unsubscribe(e.blockSubChanID)
		if err != nil {
			e.lggr.Errorf("Custom telemetry service encounters error %v when stopping", err)
			return err
		}
		e.lggr.Infof("Stopped: Custom telemetry service")
		return nil
	})
}

func (e *AutomationCustomTelemetryService) sendNodeVersionMsg() {
	e.mu.RLock()
	configDigest := e.configDigest
	e.mu.RUnlock()

	vMsg := &telem.NodeVersion{
		Timestamp:    uint64(time.Now().UTC().UnixMilli()),
		NodeVersion:  static.Version,
		ConfigDigest: configDigest[:],
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
		e.lggr.Debugf("NodeVersion Message Sent to Endpoint: %d", vMsg.Timestamp)
	}
}

func (e *AutomationCustomTelemetryService) sendBlockNumberMsg(blockKey ocr2keepers.BlockKey) {
	e.mu.RLock()
	configDigest := e.configDigest
	e.mu.RUnlock()

	blockNumMsg := &telem.BlockNumber{
		Timestamp:    uint64(time.Now().UTC().UnixMilli()),
		BlockNumber:  uint64(blockKey.Number),
		BlockHash:    hex.EncodeToString(blockKey.Hash[:]),
		ConfigDigest: configDigest[:],
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
		e.lggr.Debugf("BlockNumber Message Sent to Endpoint: %d", blockNumMsg.Timestamp)
	}
}
