package headreporter_test

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/utils"

	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"

	mocks2 "github.com/smartcontractkit/chainlink/v2/common/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/headreporter"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
	utils2 "github.com/smartcontractkit/chainlink/v2/core/utils"
	testutils2 "github.com/smartcontractkit/chainlink/v2/core/web/testutils"
)

func Test_EVMTelemetryReporter_NewHead(t *testing.T) {
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
			Timestamp: utils2.NonNegativeInt64ToUint64(head.Timestamp.UTC().Unix()),
			Number:    42,
			Hash:      head.Hash.Hex(),
		},
		Finalized: &telem.Block{
			Timestamp: utils2.NonNegativeInt64ToUint64(head.Parent.Load().Timestamp.UTC().Unix()),
			Number:    41,
			Hash:      head.Parent.Load().Hash.Hex(),
		},
	})
	require.NoError(t, err)

	monitoringEndpoint := mocks2.NewMonitoringEndpoint(t)
	monitoringEndpoint.On("SendLog", requestBytes).Return()

	monitoringEndpointGen := telemetry.NewMockMonitoringEndpointGenerator(t)
	monitoringEndpointGen.
		On("GenMonitoringEndpoint", "EVM", "100", "", synchronization.HeadReport).
		Return(monitoringEndpoint)
	reporter := headreporter.NewLegacyEVMTelemetryReporter(monitoringEndpointGen, logger.TestLogger(t), big.NewInt(100))

	err = reporter.ReportNewHead(testutils.Context(t), &head)
	assert.NoError(t, err)
}

func Test_EVMTelemetryReporter_NewHeadMissingFinalized(t *testing.T) {
	head := evmtypes.Head{
		Number:     42,
		EVMChainID: ubig.NewI(100),
		Hash:       common.HexToHash("0x1010"),
		Timestamp:  time.UnixMilli(1000),
	}
	requestBytes, err := proto.Marshal(&telem.HeadReportRequest{
		ChainID: "100",
		Latest: &telem.Block{
			Timestamp: utils2.NonNegativeInt64ToUint64(head.Timestamp.UTC().Unix()),
			Number:    42,
			Hash:      head.Hash.Hex(),
		},
	})
	require.NoError(t, err)

	monitoringEndpoint := mocks2.NewMonitoringEndpoint(t)
	monitoringEndpoint.On("SendLog", requestBytes).Return()

	monitoringEndpointGen := telemetry.NewMockMonitoringEndpointGenerator(t)
	monitoringEndpointGen.
		On("GenMonitoringEndpoint", "EVM", "100", "", synchronization.HeadReport).
		Return(monitoringEndpoint)
	reporter := headreporter.NewLegacyEVMTelemetryReporter(monitoringEndpointGen, logger.TestLogger(t), big.NewInt(100))

	err = reporter.ReportNewHead(testutils.Context(t), &head)
	assert.NoError(t, err)
}

func Test_EVMTelemetryReporter_NewHead_MissingEndpoint(t *testing.T) {
	monitoringEndpointGen := telemetry.NewMockMonitoringEndpointGenerator(t)
	monitoringEndpointGen.
		On("GenMonitoringEndpoint", "EVM", "100", "", synchronization.HeadReport).
		Return(nil)

	reporter := headreporter.NewLegacyEVMTelemetryReporter(monitoringEndpointGen, logger.TestLogger(t), big.NewInt(100))

	head := evmtypes.Head{Number: 42, EVMChainID: ubig.NewI(100)}

	err := reporter.ReportNewHead(testutils.Context(t), &head)
	assert.Errorf(t, err, "No monitoring endpoint provided chain_id=100")
}

type mockRelayer struct {
	testutils2.MockRelayer
	latestHead types.Head
}

func (m mockRelayer) LatestHead(_ context.Context) (types.Head, error) {
	return m.latestHead, nil
}

func Test_SolanaTelemetryReporter_ReportPeriodic(t *testing.T) {
	blockHash := [32]byte(utils.GetRandomPubKey(t))

	head := types.Head{
		Height:    "42",
		Hash:      blockHash[:],
		Timestamp: 1000,
	}
	relay := mockRelayer{latestHead: head}
	solanaRelays := map[types.RelayID]loop.Relayer{
		types.RelayID{Network: "Solana", ChainID: "testchain"}: relay,
	}

	request := telem.HeadReportRequest{
		ChainID: "testchain",
		Latest: &telem.Block{
			Timestamp: head.Timestamp,
			Number:    42,
			Hash:      hex.EncodeToString(head.Hash),
		},
	}
	requestBytes, err := proto.Marshal(&request)
	require.NoError(t, err)

	monitoringEndpoint := mocks2.NewMonitoringEndpoint(t)
	monitoringEndpoint.On("SendLog", requestBytes).Return()

	monitoringEndpointGen := telemetry.NewMockMonitoringEndpointGenerator(t)
	monitoringEndpointGen.
		On("GenMonitoringEndpoint", "Solana", "testchain", "", synchronization.HeadReport).
		Return(monitoringEndpoint)

	reporter := headreporter.NewTelemetryReporter(monitoringEndpointGen, logger.TestLogger(t), solanaRelays)

	err = reporter.ReportPeriodic(testutils.Context(t))
	assert.NoError(t, err)
}

func Test_SolanaTelemetryReporter_ReportPeriodic_MissingEndpoint(t *testing.T) {
	monitoringEndpoint := mocks2.NewMonitoringEndpoint(t)

	monitoringEndpointGen := telemetry.NewMockMonitoringEndpointGenerator(t)
	monitoringEndpointGen.
		On("GenMonitoringEndpoint", "Solana", "testchain", "", synchronization.HeadReport).
		Return(monitoringEndpoint)

	solanaRelays := map[types.RelayID]loop.Relayer{
		types.RelayID{Network: "Solana", ChainID: "testchain"}: testutils2.MockRelayer{},
	}

	reporter := headreporter.NewTelemetryReporter(monitoringEndpointGen, logger.TestLogger(t), solanaRelays)

	err := reporter.ReportPeriodic(testutils.Context(t))
	assert.Errorf(t, err, "No monitoring endpoint provided chain_id=testchain")
}
