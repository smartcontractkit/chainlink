package ocr2keeper

import (
	"context"
	"encoding/hex"
	"time"

	"cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/commontypes"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	evm21 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/static"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type AutomationCustomTelemetryService struct {
	utils.StartStopOnce
	monitoringEndpoint commontypes.MonitoringEndpoint
	blockSubscriber    *evm21.BlockSubscriber
	blockSubChanID     int
	threadCtrl         utils.ThreadControl
	lggr               logger.Logger
	configDigest       [32]byte
	latestConfigDigest latestConfigDigestGetter
}

type latestConfigDigestGetter interface {
	LatestConfigDetails(opts *bind.CallOpts) (keeper_registry_wrapper2_0.LatestConfigDetails, error)
}

// NewAutomationCustomTelemetryService creates a telemetry service for new blocks and node version
func NewAutomationCustomTelemetryService(me commontypes.MonitoringEndpoint,
	lggr logger.Logger, chain evm.Chain, rAddr common.Address, blocksub *evm21.BlockSubscriber) (*AutomationCustomTelemetryService, error) {
	registry, rErr := keeper_registry_wrapper2_0.NewKeeperRegistry(rAddr, chain.Client())
	if rErr != nil {
		return nil, errors.Wrap(rErr, "error creating new Registry Wrapper for customTelemService")
	}
	return &AutomationCustomTelemetryService{
		monitoringEndpoint: me,
		threadCtrl:         utils.NewThreadControl(),
		lggr:               lggr.Named("Automation Custom Telem"),
		latestConfigDigest: registry,
		blockSubscriber:    blocksub,
	}, nil
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
			e.sendNodeVersionMsg()
		}
		e.threadCtrl.Go(func(ctx context.Context) {
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
						e.sendNodeVersionMsg()
					}
				case <-ctx.Done():
					return
				}
			}
		})

		chanID, blockSubscriberChan, blockSubErr := e.blockSubscriber.Subscribe()
		if blockSubErr != nil {
			e.lggr.Errorf("Block Subscriber Error: Subscribe(): %s", blockSubErr)
		}
		e.blockSubChanID = chanID
		e.threadCtrl.Go(func(ctx context.Context) {
			e.lggr.Infof("Started: Sending BlockNumber Messages")
			for {
				select {
				case blockHistory := <-blockSubscriberChan:
					latestBlockKey, err := blockHistory.Latest()
					if err != nil {
						e.lggr.Errorf("BlockSubscriber BlockHistory.Latest() failed: %s", err)
						continue
					}
					blockNumMsg := &telem.BlockNumber{
						Timestamp:    uint64(time.Now().UTC().UnixMilli()),
						BlockNumber:  uint64(latestBlockKey.Number),
						BlockHash:    hex.EncodeToString(latestBlockKey.Hash[:]),
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
	// use utils
	return e.StopOnce("AutomationCustomTelemetryService", func() error {
		e.lggr.Infof("Stopping: custom telemetry service")
		e.threadCtrl.Close()
		e.blockSubscriber.Unsubscribe(e.blockSubChanID)
		e.lggr.Infof("Stopped: Custom telemetry service")
		return nil
	})
}

func (e *AutomationCustomTelemetryService) sendNodeVersionMsg() {
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
}
