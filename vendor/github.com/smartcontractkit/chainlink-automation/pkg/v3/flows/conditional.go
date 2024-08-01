package flows

import (
	"context"
	"fmt"
	"log"
	"time"

	common "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	ocr2keepersv3 "github.com/smartcontractkit/chainlink-automation/pkg/v3"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/postprocessors"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/preprocessors"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/random"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/service"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/telemetry"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/tickers"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"
)

const (
	// This is the ticker interval for sampling conditional flow
	SamplingConditionInterval = 3 * time.Second
	// Maximum number of upkeeps to be sampled in every round
	MaxSampledConditionals = 300
	// This is the ticker interval for final conditional flow
	FinalConditionalInterval = 1 * time.Second
	// These are the maximum number of conditional upkeeps dequeued on every tick from proposal queue in FinalConditionalFlow
	// This is kept same as OutcomeSurfacedProposalsLimit as those many can get enqueued by plugin in every round
	FinalConditionalBatchSize = 50
)

func newSampleProposalFlow(
	pre []ocr2keepersv3.PreProcessor[common.UpkeepPayload],
	ratio types.Ratio,
	getter common.ConditionalUpkeepProvider,
	ms types.MetadataStore,
	runner ocr2keepersv3.Runner,
	interval time.Duration,
	logger *log.Logger,
) service.Recoverable {
	pre = append(pre, preprocessors.NewProposalFilterer(ms, types.LogTrigger))
	postprocessors := postprocessors.NewAddProposalToMetadataStorePostprocessor(ms)

	observer := ocr2keepersv3.NewRunnableObserver(
		pre,
		postprocessors,
		runner,
		ObservationProcessLimit,
		log.New(logger.Writer(), fmt.Sprintf("[%s | sample-proposal-observer]", telemetry.ServiceName), telemetry.LogPkgStdFlags),
	)

	return tickers.NewTimeTicker[[]common.UpkeepPayload](interval, observer, func(ctx context.Context, _ time.Time) (tickers.Tick[[]common.UpkeepPayload], error) {
		return NewSampler(ratio, getter, logger), nil
	}, log.New(logger.Writer(), fmt.Sprintf("[%s | sample-proposal-ticker]", telemetry.ServiceName), telemetry.LogPkgStdFlags))
}

func NewSampler(
	ratio types.Ratio,
	getter common.ConditionalUpkeepProvider,
	logger *log.Logger,
) *sampler {
	return &sampler{
		logger:   logger,
		getter:   getter,
		ratio:    ratio,
		shuffler: random.Shuffler[common.UpkeepPayload]{Source: random.NewCryptoRandSource()},
	}
}

type shuffler[T any] interface {
	Shuffle([]T) []T
}

type sampler struct {
	logger *log.Logger

	ratio    types.Ratio
	getter   common.ConditionalUpkeepProvider
	shuffler shuffler[common.UpkeepPayload]
}

func (s *sampler) Value(ctx context.Context) ([]common.UpkeepPayload, error) {
	upkeeps, err := s.getter.GetActiveUpkeeps(ctx)
	if err != nil {
		return nil, err
	}
	if len(upkeeps) == 0 {
		return nil, nil
	}

	upkeeps = s.shuffler.Shuffle(upkeeps)
	size := s.ratio.OfInt(len(upkeeps))

	if size <= 0 {
		return nil, nil
	}
	if size > MaxSampledConditionals {
		s.logger.Printf("Required sample size %d exceeds max allowed conditional samples %d, limiting to max", size, MaxSampledConditionals)
		size = MaxSampledConditionals
	}
	if len(upkeeps) < size {
		size = len(upkeeps)
	}
	s.logger.Printf("sampled %d upkeeps", size)
	return upkeeps[:size], nil
}

func newFinalConditionalFlow(
	preprocessors []ocr2keepersv3.PreProcessor[common.UpkeepPayload],
	resultStore types.ResultStore,
	runner ocr2keepersv3.Runner,
	interval time.Duration,
	proposalQ types.ProposalQueue,
	builder common.PayloadBuilder,
	retryQ types.RetryQueue,
	stateUpdater common.UpkeepStateUpdater,
	logger *log.Logger,
) service.Recoverable {
	post := postprocessors.NewCombinedPostprocessor(
		postprocessors.NewEligiblePostProcessor(resultStore, telemetry.WrapLogger(logger, "conditional-final-eligible-postprocessor")),
		postprocessors.NewRetryablePostProcessor(retryQ, telemetry.WrapLogger(logger, "conditional-final-retryable-postprocessor")),
	)
	// create observer that only pushes results to result stores. everything at
	// this point can be dropped. this process is only responsible for running
	// conditional proposals that originate from network agreements
	observer := ocr2keepersv3.NewRunnableObserver(
		preprocessors,
		post,
		runner,
		ObservationProcessLimit,
		log.New(logger.Writer(), fmt.Sprintf("[%s | conditional-final-observer]", telemetry.ServiceName), telemetry.LogPkgStdFlags),
	)

	ticker := tickers.NewTimeTicker[[]common.UpkeepPayload](interval, observer, func(ctx context.Context, _ time.Time) (tickers.Tick[[]common.UpkeepPayload], error) {
		return coordinatedProposalsTick{
			logger:    logger,
			builder:   builder,
			q:         proposalQ,
			utype:     types.ConditionTrigger,
			batchSize: FinalConditionalBatchSize,
		}, nil
	}, log.New(logger.Writer(), fmt.Sprintf("[%s | conditional-final-ticker]", telemetry.ServiceName), telemetry.LogPkgStdFlags))

	return ticker
}
