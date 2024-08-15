package headreporter

import (
	"context"
	"math/big"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/commontypes"
	"google.golang.org/protobuf/proto"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
)

type telemetryReporter struct {
	endpoints map[uint64]commontypes.MonitoringEndpoint
}

func NewTelemetryReporter(monitoringEndpointGen telemetry.MonitoringEndpointGenerator, chainIDs ...*big.Int) HeadReporter {
	endpoints := make(map[uint64]commontypes.MonitoringEndpoint)
	for _, chainID := range chainIDs {
		endpoints[chainID.Uint64()] = monitoringEndpointGen.GenMonitoringEndpoint("EVM", chainID.String(), "", synchronization.HeadReport)
	}
	return &telemetryReporter{endpoints: endpoints}
}

func (t *telemetryReporter) ReportNewHead(ctx context.Context, head *evmtypes.Head) error {
	monitoringEndpoint := t.endpoints[head.EVMChainID.ToInt().Uint64()]
	if monitoringEndpoint == nil {
		return errors.Errorf("No monitoring endpoint provided chain_id=%d", head.EVMChainID.Int64())
	}
	var finalized *telem.Block
	latestFinalizedHead := head.LatestFinalizedHead()
	if latestFinalizedHead != nil {
		finalized = &telem.Block{
			Timestamp: uint64(latestFinalizedHead.GetTimestamp().UTC().Unix()),
			Number:    uint64(latestFinalizedHead.BlockNumber()),
			Hash:      latestFinalizedHead.BlockHash().Hex(),
		}
	}
	request := &telem.HeadReportRequest{
		Latest: &telem.Block{
			Timestamp: uint64(head.Timestamp.UTC().Unix()),
			Number:    uint64(head.Number),
			Hash:      head.Hash.Hex(),
		},
		Finalized: finalized,
	}
	bytes, err := proto.Marshal(request)
	if err != nil {
		return errors.WithMessage(err, "telem.HeadReportRequest marshal error")
	}
	monitoringEndpoint.SendLog(bytes)
	if finalized == nil {
		return errors.Errorf("No finalized block was found for chain_id=%d", head.EVMChainID.Int64())
	}
	return nil
}

func (t *telemetryReporter) ReportPeriodic(ctx context.Context) error {
	return nil
}
