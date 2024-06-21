package headreporter

import (
	"context"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
	"github.com/smartcontractkit/libocr/commontypes"
	"google.golang.org/protobuf/proto"
	"math/big"
)

type (
	telemetryReporter struct {
		logger    logger.Logger
		endpoints map[*big.Int]commontypes.MonitoringEndpoint
	}
)

func NewTelemetryReporter(chainContainer legacyevm.LegacyChainContainer, lggr logger.Logger, monitoringEndpointGen telemetry.MonitoringEndpointGenerator) *telemetryReporter {
	endpoints := make(map[*big.Int]commontypes.MonitoringEndpoint)
	for _, chain := range chainContainer.Slice() {
		endpoints[chain.ID()] = monitoringEndpointGen.GenMonitoringEndpoint("EVM", chain.ID().String(), "", synchronization.HeadReport)
	}
	return &telemetryReporter{
		logger:    lggr.Named("TelemetryReporter"),
		endpoints: endpoints,
	}
}

func (t *telemetryReporter) ReportNewHead(ctx context.Context, head *evmtypes.Head) {
	monitoringEndpoint := t.endpoints[head.EVMChainID.ToInt()]
	request := &telem.HeadReportRequest{
		ChainId:     head.EVMChainID.String(),
		Timestamp:   uint64(head.Timestamp.UTC().Unix()),
		BlockNumber: uint64(head.Number),
		BlockHash:   head.Hash.Hex(),
		Finalized:   head.IsFinalized,
	}
	bytes, err := proto.Marshal(request)
	if err != nil {
		t.logger.Warnw("telem.HeadReportRequest marshal error", "err", err)
		return
	}
	monitoringEndpoint.SendLog(bytes)
}

func (t *telemetryReporter) ReportPeriodic(ctx context.Context) {
	//do nothing
}
