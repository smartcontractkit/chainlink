package ocr

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

// Gun is a gun for the OCR load test
// it triggers new rounds for provided feed(aggregator) contract
type Gun struct {
	roundNum     atomic.Int64
	ocrInstances []contracts.OffchainAggregator
	seth         *seth.Client
	l            zerolog.Logger
}

func NewGun(l zerolog.Logger, seth *seth.Client, ocrInstances []contracts.OffchainAggregator) *Gun {
	return &Gun{
		l:            l,
		seth:         seth,
		ocrInstances: ocrInstances,
	}
}

func (m *Gun) Call(_ *wasp.Generator) *wasp.Response {
	m.roundNum.Add(1)
	requestedRound := m.roundNum.Load()
	m.l.Info().
		Int64("RoundNum", requestedRound).
		Str("FeedID", m.ocrInstances[0].Address()).
		Msg("starting new round")
	err := m.ocrInstances[0].RequestNewRound()
	if err != nil {
		return &wasp.Response{Error: err.Error(), Failed: true}
	}
	for {
		time.Sleep(5 * time.Second)
		lr, err := m.ocrInstances[0].GetLatestRound(context.Background())
		if err != nil {
			return &wasp.Response{Error: err.Error(), Failed: true}
		}
		m.l.Info().Interface("LatestRound", lr).Msg("latest round")
		if lr.RoundId.Int64() >= requestedRound {
			return &wasp.Response{}
		}
	}
}
