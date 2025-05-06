package headreporter

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/libocr/commontypes"
	"google.golang.org/protobuf/proto"

	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type legacyEVMTelemetryReporter struct {
	lggr      logger.Logger
	endpoints map[uint64]commontypes.MonitoringEndpoint
}

func NewLegacyEVMTelemetryReporter(monitoringEndpointGen telemetry.MonitoringEndpointGenerator, lggr logger.Logger, chainIDs ...*big.Int) HeadReporter {
	endpoints := make(map[uint64]commontypes.MonitoringEndpoint)
	for _, chainID := range chainIDs {
		endpoints[chainID.Uint64()] = monitoringEndpointGen.GenMonitoringEndpoint("EVM", chainID.String(), "", synchronization.HeadReport)
	}
	return &legacyEVMTelemetryReporter{lggr: lggr.Named("TelemetryReporter"), endpoints: endpoints}
}

func (t *legacyEVMTelemetryReporter) ReportNewHead(ctx context.Context, head *evmtypes.Head) error {
	monitoringEndpoint := t.endpoints[head.EVMChainID.ToInt().Uint64()]
	if monitoringEndpoint == nil {
		return fmt.Errorf("No monitoring endpoint provided chain_id=%d", head.EVMChainID.Int64())
	}
	var finalized *telem.Block
	latestFinalizedHead := head.LatestFinalizedHead()
	if latestFinalizedHead != nil {
		finalized = &telem.Block{
			Timestamp: utils.NonNegativeInt64ToUint64(latestFinalizedHead.GetTimestamp().UTC().Unix()),
			Number:    utils.NonNegativeInt64ToUint64(latestFinalizedHead.BlockNumber()),
			Hash:      latestFinalizedHead.BlockHash().Hex(),
		}
	}
	request := &telem.HeadReportRequest{
		ChainID: head.EVMChainID.String(),
		Latest: &telem.Block{
			Timestamp: utils.NonNegativeInt64ToUint64(head.Timestamp.UTC().Unix()),
			Number:    utils.NonNegativeInt64ToUint64(head.Number),
			Hash:      head.Hash.Hex(),
		},
		Finalized: finalized,
	}
	bytes, err := proto.Marshal(request)
	if err != nil {
		return fmt.Errorf("telem.HeadReportRequest marshal error: %w", err)
	}
	monitoringEndpoint.SendLog(bytes)
	if finalized == nil {
		t.lggr.Infow("No finalized block was found", "chainID", head.EVMChainID.Int64(),
			"head.number", head.Number, "chainLength", head.ChainLength())
	}
	return nil
}

func (t *legacyEVMTelemetryReporter) ReportPeriodic(_ context.Context) error {
	return nil
}

type loopTelemetryReporter struct {
	lggr      logger.Logger
	endpoints map[types.RelayID]commontypes.MonitoringEndpoint
	relays    map[types.RelayID]loop.Relayer
}

// NewTelemetryReporter creates a new telemetry reporter for each relayer
func NewTelemetryReporter(monitoringEndpointGen telemetry.MonitoringEndpointGenerator, lggr logger.Logger, relayers map[types.RelayID]loop.Relayer) HeadReporter {
	if relayers == nil {
		return nil
	}
	endpoints := make(map[types.RelayID]commontypes.MonitoringEndpoint)
	for relayID := range relayers {
		endpoints[relayID] = monitoringEndpointGen.GenMonitoringEndpoint(relayID.Network, relayID.ChainID, "", synchronization.HeadReport)
	}
	return &loopTelemetryReporter{lggr: lggr.Named("TelemetryReporter"), endpoints: endpoints, relays: relayers}
}

// ReportNewHead is unimplemented on Solana because there is no Headtracker to subscribe to
func (t *loopTelemetryReporter) ReportNewHead(_ context.Context, _ *evmtypes.Head) error {
	return nil
}

// ReportPeriodic is used on Solana to report the latest head
func (t *loopTelemetryReporter) ReportPeriodic(ctx context.Context) error {
	for relayID, endpoint := range t.endpoints {
		relay, ok := t.relays[relayID]
		if !ok {
			return fmt.Errorf("no relay found for Solana chain_id=%s", relayID.ChainID)
		}
		err := reportLatestHead(ctx, endpoint, relayID.ChainID, relay)
		if err != nil {
			return err
		}
	}

	return nil
}

func reportLatestHead(ctx context.Context, endpoint commontypes.MonitoringEndpoint, chainID string, relay loop.Relayer) error {
	head, err := relay.LatestHead(ctx)
	if err != nil {
		return fmt.Errorf("failed getting Solana head for chainID %s: %w", chainID, err)
	}

	blockNum, err := strconv.ParseUint(head.Height, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse Solana block height %s: %w", head.Height, err)
	}

	request := &telem.HeadReportRequest{
		ChainID: chainID,
		Latest: &telem.Block{
			Timestamp: head.Timestamp,
			Number:    blockNum,
			Hash:      hex.EncodeToString(head.Hash),
		},
		Finalized: nil, // latest finalized head retrieval not supported by Solana relayer yet
	}
	bytes, err := proto.Marshal(request)
	if err != nil {
		return fmt.Errorf("telem.HeadReportRequest marshal error: %w", err)
	}
	endpoint.SendLog(bytes)
	return nil
}
