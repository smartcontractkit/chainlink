package headreporter_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	mocks2 "github.com/smartcontractkit/chainlink/v2/common/types/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/headreporter"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
)

func Test_TelemetryReporter_NewHead(t *testing.T) {
	head := evmtypes.Head{
		Number:     42,
		EVMChainID: ubig.NewI(100),
		Hash:       common.HexToHash("0x1010"),
		Timestamp:  time.UnixMilli(1000),
	}
	h41 := &evmtypes.Head{
		Number:    41,
		Hash:      common.HexToHash("0x1009"),
		Timestamp: time.UnixMilli(999),
	}
	h41.IsFinalized.Store(true)
	head.Parent.Store(h41)
	requestBytes, err := proto.Marshal(&telem.HeadReportRequest{
		ChainID: "100",
		Latest: &telem.Block{
			Timestamp: uint64(head.Timestamp.UTC().Unix()),
			Number:    42,
			Hash:      head.Hash.Hex(),
		},
		Finalized: &telem.Block{
			Timestamp: uint64(head.Parent.Load().Timestamp.UTC().Unix()),
			Number:    41,
			Hash:      head.Parent.Load().Hash.Hex(),
		},
	})
	assert.NoError(t, err)

	monitoringEndpoint := mocks2.NewMonitoringEndpoint(t)
	monitoringEndpoint.On("SendLog", requestBytes).Return()

	monitoringEndpointGen := telemetry.NewMockMonitoringEndpointGenerator(t)
	monitoringEndpointGen.
		On("GenMonitoringEndpoint", "EVM", "100", "", synchronization.HeadReport).
		Return(monitoringEndpoint)
	reporter := headreporter.NewTelemetryReporter(monitoringEndpointGen, logger.TestLogger(t), big.NewInt(100))

	err = reporter.ReportNewHead(testutils.Context(t), &head)
	assert.NoError(t, err)
}

func Test_TelemetryReporter_NewHeadMissingFinalized(t *testing.T) {
	head := evmtypes.Head{
		Number:     42,
		EVMChainID: ubig.NewI(100),
		Hash:       common.HexToHash("0x1010"),
		Timestamp:  time.UnixMilli(1000),
	}
	requestBytes, err := proto.Marshal(&telem.HeadReportRequest{
		ChainID: "100",
		Latest: &telem.Block{
			Timestamp: uint64(head.Timestamp.UTC().Unix()),
			Number:    42,
			Hash:      head.Hash.Hex(),
		},
	})
	assert.NoError(t, err)

	monitoringEndpoint := mocks2.NewMonitoringEndpoint(t)
	monitoringEndpoint.On("SendLog", requestBytes).Return()

	monitoringEndpointGen := telemetry.NewMockMonitoringEndpointGenerator(t)
	monitoringEndpointGen.
		On("GenMonitoringEndpoint", "EVM", "100", "", synchronization.HeadReport).
		Return(monitoringEndpoint)
	reporter := headreporter.NewTelemetryReporter(monitoringEndpointGen, logger.TestLogger(t), big.NewInt(100))

	err = reporter.ReportNewHead(testutils.Context(t), &head)
	assert.NoError(t, err)
}

func Test_TelemetryReporter_NewHead_MissingEndpoint(t *testing.T) {
	monitoringEndpointGen := telemetry.NewMockMonitoringEndpointGenerator(t)
	monitoringEndpointGen.
		On("GenMonitoringEndpoint", "EVM", "100", "", synchronization.HeadReport).
		Return(nil)

	reporter := headreporter.NewTelemetryReporter(monitoringEndpointGen, logger.TestLogger(t), big.NewInt(100))

	head := evmtypes.Head{Number: 42, EVMChainID: ubig.NewI(100)}

	err := reporter.ReportNewHead(testutils.Context(t), &head)
	assert.Errorf(t, err, "No monitoring endpoint provided chain_id=100")
}
