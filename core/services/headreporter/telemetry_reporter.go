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

type evmTelemetryReporter struct {
	lggr      logger.Logger
	endpoints map[uint64]commontypes.MonitoringEndpoint
}

func NewEVMTelemetryReporter(monitoringEndpointGen telemetry.MonitoringEndpointGenerator, lggr logger.Logger, chainIDs ...*big.Int) HeadReporter {
	endpoints := make(map[uint64]commontypes.MonitoringEndpoint)
	for _, chainID := range chainIDs {
		endpoints[chainID.Uint64()] = monitoringEndpointGen.GenMonitoringEndpoint("EVM", chainID.String(), "", synchronization.HeadReport)
	}
	return &evmTelemetryReporter{lggr: lggr.Named("TelemetryReporter"), endpoints: endpoints}
}

func (t *evmTelemetryReporter) ReportNewHead(ctx context.Context, head *evmtypes.Head) error {
	monitoringEndpoint := t.endpoints[head.EVMChainID.ToInt().Uint64()]
	if monitoringEndpoint == nil {
		return fmt.Errorf("No monitoring endpoint provided chain_id=%d", head.EVMChainID.Int64())
	}
	var finalized *telem.Block
	latestFinalizedHead := head.LatestFinalizedHead()
	if latestFinalizedHead != nil {
		finalized = &telem.Block{
			Timestamp: utils.NonNegativeInt64ToUint64(latestFinalizedHead.GetTimestamp().UTC().Unix()), // golint:gosec
			Number:    utils.NonNegativeInt64ToUint64(latestFinalizedHead.BlockNumber()),               // golint:gosec
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

func (t *evmTelemetryReporter) ReportPeriodic(_ context.Context) error {
	return nil
}

type solanaTelemetryReporter struct {
	lggr      logger.Logger
	endpoints map[string]commontypes.MonitoringEndpoint
	relays    map[string]loop.Relayer
}

// NewSolanaTelemetryReporter creates a new telemetry reporter for Solana.
func NewSolanaTelemetryReporter(monitoringEndpointGen telemetry.MonitoringEndpointGenerator, lggr logger.Logger, solanaRelays map[types.RelayID]loop.Relayer) HeadReporter {
	if solanaRelays == nil {
		return nil
	}
	endpoints := make(map[string]commontypes.MonitoringEndpoint)
	relays := make(map[string]loop.Relayer)
	for relayID, relay := range solanaRelays {
		chainID := relayID.ChainID
		endpoints[relayID.ChainID] = monitoringEndpointGen.GenMonitoringEndpoint("Solana", chainID, "", synchronization.HeadReport)
		relays[relayID.ChainID] = relay
	}
	return &solanaTelemetryReporter{lggr: lggr.Named("TelemetryReporter"), endpoints: endpoints, relays: relays}
}

// ReportNewHead is unimplemented on Solana because there is no Headtracker to subscribe to
func (t *solanaTelemetryReporter) ReportNewHead(_ context.Context, _ *evmtypes.Head) error {
	return nil
}

// ReportPeriodic is used on Solana to report the latest head
func (t *solanaTelemetryReporter) ReportPeriodic(ctx context.Context) error {
	for chainID, endpoint := range t.endpoints {
		relay, ok := t.relays[chainID]
		if !ok {
			return fmt.Errorf("no relay found for Solana chain_id=%s", chainID)
		}
		err := reportLatestHead(ctx, endpoint, chainID, relay)
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
