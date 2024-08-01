package flows

import (
	"context"
	"fmt"
	"log"
	"time"

	ocr2keepersv3 "github.com/smartcontractkit/chainlink-automation/pkg/v3"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/postprocessors"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/service"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/telemetry"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/tickers"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"
	common "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
)

const (
	// These are the max number of payloads dequeued on every tick from the retry queue in the retry flow
	RetryBatchSize = 10
	// This is the ticker interval for retry flow
	RetryCheckInterval = 5 * time.Second
)

func NewRetryFlow(
	coord ocr2keepersv3.PreProcessor[common.UpkeepPayload],
	resultStore types.ResultStore,
	runner ocr2keepersv3.Runner,
	retryQ types.RetryQueue,
	retryTickerInterval time.Duration,
	stateUpdater common.UpkeepStateUpdater,
	logger *log.Logger,
) service.Recoverable {
	preprocessors := []ocr2keepersv3.PreProcessor[common.UpkeepPayload]{coord}
	post := postprocessors.NewCombinedPostprocessor(
		postprocessors.NewEligiblePostProcessor(resultStore, telemetry.WrapLogger(logger, "retry-eligible-postprocessor")),
		postprocessors.NewRetryablePostProcessor(retryQ, telemetry.WrapLogger(logger, "retry-retryable-postprocessor")),
		postprocessors.NewIneligiblePostProcessor(stateUpdater, telemetry.WrapLogger(logger, "retry-ineligible-postprocessor")),
	)

	obs := ocr2keepersv3.NewRunnableObserver(
		preprocessors,
		post,
		runner,
		ObservationProcessLimit,
		log.New(logger.Writer(), fmt.Sprintf("[%s | retry-observer]", telemetry.ServiceName), telemetry.LogPkgStdFlags),
	)

	timeTick := tickers.NewTimeTicker[[]common.UpkeepPayload](retryTickerInterval, obs, func(ctx context.Context, _ time.Time) (tickers.Tick[[]common.UpkeepPayload], error) {
		return retryTick{logger: logger, q: retryQ, batchSize: RetryBatchSize}, nil
	}, log.New(logger.Writer(), fmt.Sprintf("[%s | retry-ticker]", telemetry.ServiceName), telemetry.LogPkgStdFlags))

	return timeTick
}

type retryTick struct {
	logger    *log.Logger
	q         types.RetryQueue
	batchSize int
}

func (t retryTick) Value(ctx context.Context) ([]common.UpkeepPayload, error) {
	if t.q == nil {
		return nil, nil
	}

	payloads, err := t.q.Dequeue(t.batchSize)
	if err != nil {
		return nil, fmt.Errorf("failed to dequeue from retry queue: %w", err)
	}
	t.logger.Printf("%d payloads returned by retry queue", len(payloads))

	return payloads, err
}
