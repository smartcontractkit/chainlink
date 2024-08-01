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

var (
	ErrNotRetryable = fmt.Errorf("payload is not retryable")
)

const (
	// This is the ticker interval for log trigger flow
	LogCheckInterval = 1 * time.Second
	// Limit for processing a whole observer flow given a payload. The main component of this
	// is the checkPipeline which involves some RPC, DB and Mercury calls, this is limited
	// to 20 seconds for now
	ObservationProcessLimit = 20 * time.Second
)

// log trigger flow is the happy path entry point for log triggered upkeeps
func newLogTriggerFlow(
	preprocessors []ocr2keepersv3.PreProcessor[common.UpkeepPayload],
	rs types.ResultStore,
	rn ocr2keepersv3.Runner,
	logProvider common.LogEventProvider,
	logInterval time.Duration,
	retryQ types.RetryQueue,
	stateUpdater common.UpkeepStateUpdater,
	logger *log.Logger,
) service.Recoverable {
	post := postprocessors.NewCombinedPostprocessor(
		postprocessors.NewEligiblePostProcessor(rs, telemetry.WrapLogger(logger, "log-trigger-eligible-postprocessor")),
		postprocessors.NewRetryablePostProcessor(retryQ, telemetry.WrapLogger(logger, "log-trigger-retryable-postprocessor")),
		postprocessors.NewIneligiblePostProcessor(stateUpdater, telemetry.WrapLogger(logger, "retry-ineligible-postprocessor")),
	)

	obs := ocr2keepersv3.NewRunnableObserver(
		preprocessors,
		post,
		rn,
		ObservationProcessLimit,
		log.New(logger.Writer(), fmt.Sprintf("[%s | log-trigger-observer]", telemetry.ServiceName), telemetry.LogPkgStdFlags),
	)

	timeTick := tickers.NewTimeTicker[[]common.UpkeepPayload](logInterval, obs, func(ctx context.Context, _ time.Time) (tickers.Tick[[]common.UpkeepPayload], error) {
		return logTick{logger: logger, logProvider: logProvider}, nil
	}, log.New(logger.Writer(), fmt.Sprintf("[%s | log-trigger-ticker]", telemetry.ServiceName), telemetry.LogPkgStdFlags))

	return timeTick
}

type logTick struct {
	logProvider common.LogEventProvider
	logger      *log.Logger
}

func (et logTick) Value(ctx context.Context) ([]common.UpkeepPayload, error) {
	if et.logProvider == nil {
		return nil, nil
	}

	logs, err := et.logProvider.GetLatestPayloads(ctx)

	et.logger.Printf("%d logs returned by log provider", len(logs))

	return logs, err
}
