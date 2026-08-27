package headreporter

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	testutils2 "github.com/smartcontractkit/chainlink/v2/core/web/testutils"
)

type metricsMockRelayer struct {
	testutils2.MockRelayer
	latestHead    types.Head
	finalizedHead *types.Head
	finalizedErr  error
}

func (m *metricsMockRelayer) LatestHead(_ context.Context) (types.Head, error) {
	return m.latestHead, nil
}

func (m *metricsMockRelayer) FinalizedHead(_ context.Context) (types.Head, error) {
	if m.finalizedErr != nil {
		return types.Head{}, m.finalizedErr
	}
	if m.finalizedHead != nil {
		return *m.finalizedHead, nil
	}
	return types.Head{}, status.Errorf(codes.Unimplemented, "method FinalizedHead not implemented")
}

func Test_EVMMetricsReporter_ReportNewHead_WithFinalized(t *testing.T) {
	t.Parallel()
	head := evmtypes.Head{
		Number:     42,
		EVMChainID: sqlutil.NewI(100),
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

	metrics := NewMockHeadMetrics(t)
	metrics.EXPECT().RecordHeadReport(mock.Anything, headReport{
		chainID:       "100",
		network:       "evm",
		chainSelector: 465200170687744372, // gnosis_chain-mainnet, chain ID 100 is registered
		hasSelector:   true,
		latestNumber:  42,
		latestTs:      utils.NonNegativeInt64ToUint64(head.Timestamp.UTC().Unix()),
		finalized: &finalizedBlock{
			number: 41,
			ts:     utils.NonNegativeInt64ToUint64(h41.Timestamp.UTC().Unix()),
		},
	}).Return()

	reporter := NewEVMMetricsReporter(metrics, logger.TestLogger(t), big.NewInt(100))
	err := reporter.ReportNewHead(t.Context(), &head)
	assert.NoError(t, err)
}

func Test_EVMMetricsReporter_ReportNewHead_MissingFinalized(t *testing.T) {
	t.Parallel()
	head := evmtypes.Head{
		Number:     42,
		EVMChainID: sqlutil.NewI(100),
		Hash:       common.HexToHash("0x1010"),
		Timestamp:  time.UnixMilli(1000),
	}

	metrics := NewMockHeadMetrics(t)
	metrics.EXPECT().RecordHeadReport(mock.Anything, headReport{
		chainID:       "100",
		network:       "evm",
		chainSelector: 465200170687744372,
		hasSelector:   true,
		latestNumber:  42,
		latestTs:      utils.NonNegativeInt64ToUint64(head.Timestamp.UTC().Unix()),
		finalized:     nil,
	}).Return()

	reporter := NewEVMMetricsReporter(metrics, logger.TestLogger(t), big.NewInt(100))
	err := reporter.ReportNewHead(t.Context(), &head)
	assert.NoError(t, err)
}

func Test_EVMMetricsReporter_ReportNewHead_UnknownChainSelector(t *testing.T) {
	t.Parallel()
	unregisteredChainID := big.NewInt(123456789012345) // not present in the chain-selectors registry
	head := evmtypes.Head{
		Number:     42,
		EVMChainID: sqlutil.NewI(123456789012345),
		Hash:       common.HexToHash("0x1010"),
		Timestamp:  time.UnixMilli(1000),
	}

	metrics := NewMockHeadMetrics(t)
	metrics.EXPECT().RecordHeadReport(mock.Anything, headReport{
		chainID:      "123456789012345",
		network:      "evm",
		hasSelector:  false,
		latestNumber: 42,
		latestTs:     utils.NonNegativeInt64ToUint64(head.Timestamp.UTC().Unix()),
	}).Return()

	reporter := NewEVMMetricsReporter(metrics, logger.TestLogger(t), unregisteredChainID)
	err := reporter.ReportNewHead(t.Context(), &head)
	assert.NoError(t, err)
}

func Test_EVMMetricsReporter_ReportPeriodic_NoOp(t *testing.T) {
	t.Parallel()
	metrics := NewMockHeadMetrics(t)
	reporter := NewEVMMetricsReporter(metrics, logger.TestLogger(t), big.NewInt(100))
	assert.NoError(t, reporter.ReportPeriodic(t.Context()))
}

func Test_RelayerMetricsReporter_NilRelayers(t *testing.T) {
	t.Parallel()
	metrics := NewMockHeadMetrics(t)
	reporter := NewRelayerMetricsReporter(metrics, logger.TestLogger(t), nil)
	assert.Nil(t, reporter)
}

func Test_RelayerMetricsReporter_ReportNewHead_NoOp(t *testing.T) {
	t.Parallel()
	metrics := NewMockHeadMetrics(t)
	relays := map[types.RelayID]loop.Relayer{
		{Network: "Solana", ChainID: "testchain"}: &metricsMockRelayer{},
	}
	reporter := NewRelayerMetricsReporter(metrics, logger.TestLogger(t), relays)
	assert.NoError(t, reporter.ReportNewHead(t.Context(), &evmtypes.Head{}))
}

func Test_RelayerMetricsReporter_ReportPeriodic(t *testing.T) {
	t.Parallel()
	privKey, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)
	blockHash := [32]byte(privKey.PublicKey())

	head := types.Head{
		Height:    "42",
		Hash:      blockHash[:],
		Timestamp: 1000,
	}
	relay := &metricsMockRelayer{latestHead: head}
	relays := map[types.RelayID]loop.Relayer{
		{Network: "solana", ChainID: "testchain"}: relay,
	}

	metrics := NewMockHeadMetrics(t)
	metrics.EXPECT().RecordHeadReport(mock.Anything, headReport{
		chainID:      "testchain",
		network:      "solana",
		hasSelector:  false, // "testchain" is not a registered chain ID
		latestNumber: 42,
		latestTs:     1000,
	}).Return()

	reporter := NewRelayerMetricsReporter(metrics, logger.TestLogger(t), relays)
	err = reporter.ReportPeriodic(t.Context())
	assert.NoError(t, err)
}

func Test_RelayerMetricsReporter_ReportPeriodic_WithFinalizedHead(t *testing.T) {
	t.Parallel()
	privKey, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)
	blockHash := [32]byte(privKey.PublicKey())

	head := types.Head{
		Height:    "42",
		Hash:      blockHash[:],
		Timestamp: 1000,
	}
	fHead := &types.Head{Height: "40", Timestamp: 900}
	relay := &metricsMockRelayer{latestHead: head, finalizedHead: fHead}
	relays := map[types.RelayID]loop.Relayer{
		{Network: "solana", ChainID: "testchain"}: relay,
	}

	metrics := NewMockHeadMetrics(t)
	metrics.EXPECT().RecordHeadReport(mock.Anything, headReport{
		chainID:      "testchain",
		network:      "solana",
		hasSelector:  false,
		latestNumber: 42,
		latestTs:     1000,
		finalized: &finalizedBlock{
			number: 40,
			ts:     900,
		},
	}).Return()

	reporter := NewRelayerMetricsReporter(metrics, logger.TestLogger(t), relays)
	err = reporter.ReportPeriodic(t.Context())
	assert.NoError(t, err)
}

func Test_RelayerMetricsReporter_ReportPeriodic_FinalizedHeadError(t *testing.T) {
	t.Parallel()
	privKey, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)
	blockHash := [32]byte(privKey.PublicKey())

	head := types.Head{
		Height:    "42",
		Hash:      blockHash[:],
		Timestamp: 1000,
	}
	relay := &metricsMockRelayer{latestHead: head, finalizedErr: errors.New("rpc error")}
	relays := map[types.RelayID]loop.Relayer{
		{Network: "solana", ChainID: "testchain"}: relay,
	}

	metrics := NewMockHeadMetrics(t)
	metrics.EXPECT().RecordHeadReport(mock.Anything, headReport{
		chainID:      "testchain",
		network:      "solana",
		hasSelector:  false,
		latestNumber: 42,
		latestTs:     1000,
	}).Return()

	reporter := NewRelayerMetricsReporter(metrics, logger.TestLogger(t), relays)
	err = reporter.ReportPeriodic(t.Context())
	assert.NoError(t, err)
}

func Test_RelayerMetricsReporter_ReportPeriodic_EmptyBlockHeight(t *testing.T) {
	t.Parallel()
	relay := &metricsMockRelayer{latestHead: types.Head{Height: ""}}
	relays := map[types.RelayID]loop.Relayer{
		{Network: "solana", ChainID: "testchain"}: relay,
	}

	metrics := NewMockHeadMetrics(t)
	reporter := NewRelayerMetricsReporter(metrics, logger.TestLogger(t), relays)
	err := reporter.ReportPeriodic(t.Context())
	assert.ErrorContains(t, err, "latest block height returned by relayer is empty for {solana testchain}")
}
