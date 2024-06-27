package headreporter_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/headreporter"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	mocks2 "github.com/smartcontractkit/chainlink/v2/core/services/telemetry/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
)

type IngressAgent struct {
	mock.Mock
}

func (t *IngressAgent) SendLog(telemetry []byte) {
	_ = t.Called(telemetry)
}

func NewIngressAgent(t interface {
	mock.TestingT
	Cleanup(func())
}) *IngressAgent {
	m := &IngressAgent{}
	m.Mock.Test(t)

	t.Cleanup(func() { m.AssertExpectations(t) })

	return m
}

func Test_TelemetryReporter_NewHead(t *testing.T) {
	chain := mocks.NewChain(t)
	chain.On("ID").Return(big.NewInt(100))

	chains := legacyevm.NewLegacyChains(map[string]legacyevm.Chain{"100": chain}, nil)

	ingressAgent := NewIngressAgent(t)

	monitoringEndpointGen := mocks2.NewMonitoringEndpointGenerator(t)
	monitoringEndpointGen.
		On("GenMonitoringEndpoint", "EVM", "100", "", synchronization.HeadReport).
		Return(ingressAgent)
	reporter := headreporter.NewTelemetryReporter(chains, logger.TestLogger(t), monitoringEndpointGen)

	head := evmtypes.Head{
		Number:      42,
		EVMChainID:  ubig.NewI(100),
		Hash:        common.HexToHash("0x1010"),
		Timestamp:   time.UnixMilli(1000),
		IsFinalized: false,
		Parent: &evmtypes.Head{
			Number:      41,
			Hash:        common.HexToHash("0x1009"),
			Timestamp:   time.UnixMilli(999),
			IsFinalized: true,
		},
	}
	requestBytes, err := proto.Marshal(&telem.HeadReportRequest{
		Latest: &telem.Block{
			Timestamp: uint64(head.Timestamp.UTC().Unix()),
			Number:    42,
			Hash:      head.Hash.Hex(),
		},
		Finalized: &telem.Block{
			Timestamp: uint64(head.Parent.Timestamp.UTC().Unix()),
			Number:    41,
			Hash:      head.Parent.Hash.Hex(),
		},
	})
	assert.NoError(t, err)

	ingressAgent.On("SendLog", requestBytes).Return()

	err = reporter.ReportNewHead(testutils.Context(t), &head)
	assert.NoError(t, err)
}

func Test_TelemetryReporter_NewHeadMissingFinalized(t *testing.T) {
	chain := mocks.NewChain(t)
	chain.On("ID").Return(big.NewInt(100))

	chains := legacyevm.NewLegacyChains(map[string]legacyevm.Chain{"100": chain}, nil)

	ingressAgent := NewIngressAgent(t)

	monitoringEndpointGen := mocks2.NewMonitoringEndpointGenerator(t)
	monitoringEndpointGen.
		On("GenMonitoringEndpoint", "EVM", "100", "", synchronization.HeadReport).
		Return(ingressAgent)
	reporter := headreporter.NewTelemetryReporter(chains, logger.TestLogger(t), monitoringEndpointGen)

	head := evmtypes.Head{
		Number:      42,
		EVMChainID:  ubig.NewI(100),
		Hash:        common.HexToHash("0x1010"),
		Timestamp:   time.UnixMilli(1000),
		IsFinalized: false,
	}
	requestBytes, err := proto.Marshal(&telem.HeadReportRequest{
		Latest: &telem.Block{
			Timestamp: uint64(head.Timestamp.UTC().Unix()),
			Number:    42,
			Hash:      head.Hash.Hex(),
		},
	})
	assert.NoError(t, err)

	ingressAgent.On("SendLog", requestBytes).Return()

	err = reporter.ReportNewHead(testutils.Context(t), &head)
	assert.Errorf(t, err, "No finalized block was found for chain_id=100")
}

func Test_TelemetryReporter_NewHead_MissingEndpoint(t *testing.T) {
	chain := mocks.NewChain(t)
	chain.On("ID").Return(big.NewInt(100))

	chains := legacyevm.NewLegacyChains(map[string]legacyevm.Chain{"100": chain}, nil)

	monitoringEndpointGen := mocks2.NewMonitoringEndpointGenerator(t)
	monitoringEndpointGen.
		On("GenMonitoringEndpoint", "EVM", "100", "", synchronization.HeadReport).
		Return(nil)

	reporter := headreporter.NewTelemetryReporter(chains, logger.TestLogger(t), monitoringEndpointGen)

	head := evmtypes.Head{Number: 42, EVMChainID: ubig.NewI(100)}

	err := reporter.ReportNewHead(testutils.Context(t), &head)
	assert.Errorf(t, err, "No monitoring endpoint provided chain_id=100")
}
