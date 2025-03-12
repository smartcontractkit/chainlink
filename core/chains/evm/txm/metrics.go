package txm

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	svrv1 "github.com/smartcontractkit/chainlink-protos/svr/v1"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
)

var (
	promNumBroadcastedTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "txm_num_broadcasted_transactions",
		Help: "Total number of successful broadcasted transactions.",
	}, []string{"chainID"})
	promNumConfirmedTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "txm_num_confirmed_transactions",
		Help: "Total number of confirmed transactions. Note that this can happen multiple times per transaction in the case of re-orgs or when filling the nonce for untracked transactions.",
	}, []string{"chainID"})
	promNumNonceGaps = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "txm_num_nonce_gaps",
		Help: "Total number of nonce gaps created that the transaction manager had to fill.",
	}, []string{"chainID"})
	promTimeUntilTxConfirmed = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "txm_time_until_tx_confirmed",
		Help: "The amount of time elapsed from a transaction being broadcast to being included in a block.",
	}, []string{"chainID"})
)

type txmMetrics struct {
	metrics.Labeler
	chainID              *big.Int
	numBroadcastedTxs    metric.Int64Counter
	numConfirmedTxs      metric.Int64Counter
	numNonceGaps         metric.Int64Counter
	timeUntilTxConfirmed metric.Float64Histogram
}

func NewTxmMetrics(chainID *big.Int) (*txmMetrics, error) {
	numBroadcastedTxs, err := beholder.GetMeter().Int64Counter("txm_num_broadcasted_transactions")
	if err != nil {
		return nil, fmt.Errorf("failed to register broadcasted txs number: %w", err)
	}

	numConfirmedTxs, err := beholder.GetMeter().Int64Counter("txm_num_confirmed_transactions")
	if err != nil {
		return nil, fmt.Errorf("failed to register confirmed txs number: %w", err)
	}

	numNonceGaps, err := beholder.GetMeter().Int64Counter("txm_num_nonce_gaps")
	if err != nil {
		return nil, fmt.Errorf("failed to register nonce gaps number: %w", err)
	}

	timeUntilTxConfirmed, err := beholder.GetMeter().Float64Histogram("txm_time_until_tx_confirmed")
	if err != nil {
		return nil, fmt.Errorf("failed to register time until tx confirmed: %w", err)
	}

	return &txmMetrics{
		chainID:              chainID,
		Labeler:              metrics.NewLabeler().With("chainID", chainID.String()),
		numBroadcastedTxs:    numBroadcastedTxs,
		numConfirmedTxs:      numConfirmedTxs,
		numNonceGaps:         numNonceGaps,
		timeUntilTxConfirmed: timeUntilTxConfirmed,
	}, nil
}

func (m *txmMetrics) IncrementNumBroadcastedTxs(ctx context.Context) {
	promNumBroadcastedTxs.WithLabelValues(m.chainID.String()).Add(float64(1))
	m.numBroadcastedTxs.Add(ctx, 1)
}

func (m *txmMetrics) IncrementNumConfirmedTxs(ctx context.Context, confirmedTransactions int) {
	promNumConfirmedTxs.WithLabelValues(m.chainID.String()).Add(float64(confirmedTransactions))
	m.numConfirmedTxs.Add(ctx, int64(confirmedTransactions))
}

func (m *txmMetrics) IncrementNumNonceGaps(ctx context.Context) {
	promNumNonceGaps.WithLabelValues(m.chainID.String()).Add(float64(1))
	m.numNonceGaps.Add(ctx, 1)
}

func (m *txmMetrics) RecordTimeUntilTxConfirmed(ctx context.Context, duration float64) {
	promTimeUntilTxConfirmed.WithLabelValues(m.chainID.String()).Observe(duration)
	m.timeUntilTxConfirmed.Record(ctx, duration)
}

func (m *txmMetrics) EmitTxMessage(ctx context.Context, txHash common.Hash, fromAddress common.Address, tx *types.Transaction) error {
	meta, err := tx.GetMeta()
	if err != nil {
		return fmt.Errorf("failed to get meta for tx %s: %w", txHash, err)
	}

	var destAddress string
	if meta != nil && meta.FwdrDestAddress != nil {
		destAddress = meta.FwdrDestAddress.String()
	}

	message := &svrv1.TxMessage{
		Hash:        txHash.String(),
		FromAddress: fromAddress.String(),
		ToAddress:   tx.ToAddress.String(),
		Nonce:       strconv.FormatUint(*tx.Nonce, 10),
		CreatedAt:   time.Now().UnixMicro(),
		ChainId:     m.chainID.String(),
		FeedAddress: destAddress,
	}

	messageBytes, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	return beholder.GetEmitter().Emit(
		ctx,
		messageBytes,
		"beholder_domain", "svr",
		"beholder_entity", "svr.v1.TxMessage",
		"beholder_data_schema", "/beholder-tx-message/versions/2",
	)
}
