package observation

import (
	"context"
	"math/rand"

	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
)

// testDataSource implements llo.DataSource with random values
type testDataSource struct {
	lggr logger.Logger
}

// NewTestDataSource creates a DataSource that returns random values for testing
func NewTestDataSource(lggr logger.Logger) llo.DataSource {
	return &testDataSource{
		lggr: logger.Named(lggr, "TestDataSource"),
	}
}

// Observe implements llo.DataSource.Observe with random values between 1-100
func (d *testDataSource) Observe(ctx context.Context, streamValues llo.StreamValues, opts llo.DSOpts) error {
	d.lggr.Infow("Using random test values for LLO observation",
		"observationTimestamp", opts.ObservationTimestamp(),
		"seqNr", opts.OutCtx().SeqNr,
		"streamCount", len(streamValues))

	// Generate random values between 1-100 for all streams
	for streamID := range streamValues {
		randomValue := rand.Intn(100) + 1 // Random value between 1-100
		streamValues[streamID] = llo.ToDecimal(decimal.NewFromInt(int64(randomValue)))

		d.lggr.Debugw("Set random value",
			"streamID", streamID,
			"value", randomValue)
	}

	return nil
}
