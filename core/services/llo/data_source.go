package llo

// TODO: llo datasource
import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
)

var (
	promMissingStreamCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "llo_stream_missing_count",
		Help: "Number of times we tried to observe a stream, but it was missing",
	},
		[]string{"streamID"},
	)
	promObservationErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "llo_stream_observation_error_count",
		Help: "Number of times we tried to observe a stream, but it failed with an error",
	},
		[]string{"streamID"},
	)
)

type ErrMissingStream struct {
	id string
}

type Registry interface {
	Get(streamID streams.StreamID) (strm streams.Stream, exists bool)
}

func (e ErrMissingStream) Error() string {
	return fmt.Sprintf("missing stream definition for: %q", e.id)
}

var _ llo.DataSource = &dataSource{}

type dataSource struct {
	lggr     logger.Logger
	registry Registry
}

func NewDataSource(lggr logger.Logger, registry Registry) llo.DataSource {
	// TODO: lggr should include job ID
	return &dataSource{lggr, registry}
}

func (d *dataSource) Observe(ctx context.Context, streamIDs map[commontypes.StreamID]struct{}) (llo.StreamValues, error) {
	// There is no "observationSource" (AKA pipeline)
	// Need a concept of "streams"
	// Streams are referenced by ID from the on-chain config
	// Each stream contains its own pipeline
	// See: https://docs.google.com/document/d/1l1IiDOL1QSteLTnhmiGnJAi6QpcSpyOe0nkqS7D3SvU/edit for stream ID naming

	var wg sync.WaitGroup
	wg.Add(len(streamIDs))
	sv := make(llo.StreamValues)
	var mu sync.Mutex

	for streamID := range streamIDs {
		go func(streamID commontypes.StreamID) {
			defer wg.Done()

			var res llo.ObsResult[*big.Int]

			stream, exists := d.registry.Get(streamID)
			if exists {
				run, trrs, err := stream.Run(ctx)
				// TODO: support types other than *big.Int
				// res.Val, res.Err = stream.Run(ctx)
				if err != nil {
					d.lggr.Debugw("Observation failed for stream", "err", err, "streamID", streamID, "runID", run.ID)
					promObservationErrorCount.WithLabelValues(streamID).Inc()
					res.Err = err
				} else {
					res.Val, res.Err = streams.ExtractBigInt(trrs)
				}
			} else {
				d.lggr.Errorw(fmt.Sprintf("Missing stream: %q", streamID), "streamID", streamID)
				promMissingStreamCount.WithLabelValues(streamID).Inc()
				res.Err = ErrMissingStream{streamID}
			}

			mu.Lock()
			defer mu.Unlock()
			sv[streamID] = res
		}(streamID)
	}

	wg.Wait()

	return sv, nil
}
