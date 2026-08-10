package headreporter

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
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
		return fmt.Errorf("No monitoring endpoint provided chain_id=%d", head.EVMChainID.ToInt().Int64())
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
		t.lggr.Infow("No finalized block was found", "chainID", head.EVMChainID.ToInt().Int64(),
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

// ReportNewHead is unimplemented because there is no Headtracker to subscribe to
func (t *loopTelemetryReporter) ReportNewHead(_ context.Context, _ *evmtypes.Head) error {
	return nil
}

func (t *loopTelemetryReporter) ReportPeriodic(ctx context.Context) error {
	for relayID, endpoint := range t.endpoints {
		relay, ok := t.relays[relayID]
		if !ok {
			return fmt.Errorf("no relay found for chain=%s", relayID)
		}
		err := reportLatestHead(ctx, t.lggr, endpoint, relayID, relay)
		if err != nil {
			return err
		}
	}

	return nil
}

func reportLatestHead(ctx context.Context, lggr logger.Logger, endpoint commontypes.MonitoringEndpoint, relayID types.RelayID, relay loop.Relayer) error {
	head, err := relay.LatestHead(ctx)
	if err != nil {
		return fmt.Errorf("failed getting head for chain %s: %w", relayID, err)
	}

	blockNum, err := parseGenericHeadBlockNumber(relayID, head)
	if err != nil {
		return err
	}

	var finalized *telem.Block
	finalizedHead, err := fetchFinalizedGenericHead(ctx, relayID, relay)
	if err != nil {
		lggr.Warnw("Failed to fetch finalized head", "chainID", relayID.ChainID, "network", relayID.Network, "err", err)
	} else if finalizedHead != nil {
		finalizedNum, parseErr := parseGenericHeadBlockNumber(relayID, *finalizedHead)
		if parseErr != nil {
			lggr.Warnw("Failed to parse finalized head block number", "chainID", relayID.ChainID, "network", relayID.Network, "err", parseErr)
		} else {
			finalized = &telem.Block{
				Timestamp: finalizedHead.Timestamp,
				Number:    finalizedNum,
				Hash:      hex.EncodeToString(finalizedHead.Hash),
			}
		}
	}

	request := &telem.HeadReportRequest{
		ChainID: relayID.ChainID,
		Latest: &telem.Block{
			Timestamp: head.Timestamp,
			Number:    blockNum,
			Hash:      hex.EncodeToString(head.Hash),
		},
		Finalized: finalized,
	}
	bytes, err := proto.Marshal(request)
	if err != nil {
		return fmt.Errorf("telem.HeadReportRequest marshal error: %w", err)
	}
	endpoint.SendLog(bytes)
	return nil
}

// parseGenericHeadBlockNumber parses the block height carried by a generic relayer head.
func parseGenericHeadBlockNumber(relayID types.RelayID, head types.Head) (uint64, error) {
	if head.Height == "" {
		return 0, fmt.Errorf("latest block height returned by relayer is empty for %s", relayID)
	}
	blockNum, err := strconv.ParseUint(head.Height, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s block height %s: %w", relayID, head.Height, err)
	}
	return blockNum, nil
}

// fetchFinalizedGenericHead fetches a generic relayer's finalized head, returning nil (no
// error) if the relayer does not implement FinalizedHead (codes.Unimplemented).
func fetchFinalizedGenericHead(ctx context.Context, relayID types.RelayID, relay loop.Relayer) (*types.Head, error) {
	head, err := relay.FinalizedHead(ctx)
	if err != nil {
		if status.Code(err) == codes.Unimplemented {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch finalized head for %s: %w", relayID, err)
	}
	return &head, nil
}
