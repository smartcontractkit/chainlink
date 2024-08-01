package flows

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	ocr2keepersv3 "github.com/smartcontractkit/chainlink-automation/pkg/v3"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/postprocessors"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/preprocessors"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/service"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/telemetry"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/tickers"
	common "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
)

const (
	// This is the ticker interval for recovery final flow
	RecoveryFinalInterval = 1 * time.Second
	// These are the maximum number of log upkeeps dequeued on every tick from proposal queue in FinalRecoveryFlow
	// This is kept same as OutcomeSurfacedProposalsLimit as those many can get enqueued by plugin in every round
	FinalRecoveryBatchSize = 50
	// This is the ticker interval for recovery proposal flow
	RecoveryProposalInterval = 1 * time.Second
)

func newFinalRecoveryFlow(
	preprocessors []ocr2keepersv3.PreProcessor[common.UpkeepPayload],
	resultStore types.ResultStore,
	runner ocr2keepersv3.Runner,
	retryQ types.RetryQueue,
	recoveryFinalizationInterval time.Duration,
	proposalQ types.ProposalQueue,
	builder common.PayloadBuilder,
	stateUpdater common.UpkeepStateUpdater,
	logger *log.Logger,
) service.Recoverable {
	post := postprocessors.NewCombinedPostprocessor(
		postprocessors.NewEligiblePostProcessor(resultStore, telemetry.WrapLogger(logger, "recovery-final-eligible-postprocessor")),
		postprocessors.NewRetryablePostProcessor(retryQ, telemetry.WrapLogger(logger, "recovery-final-retryable-postprocessor")),
		postprocessors.NewIneligiblePostProcessor(stateUpdater, telemetry.WrapLogger(logger, "retry-ineligible-postprocessor")),
	)
	// create observer that only pushes results to result stores. everything at
	// this point can be dropped. this process is only responsible for running
	// recovery proposals that originate from network agreements
	recoveryObserver := ocr2keepersv3.NewRunnableObserver(
		preprocessors,
		post,
		runner,
		ObservationProcessLimit,
		log.New(logger.Writer(), fmt.Sprintf("[%s | recovery-final-observer]", telemetry.ServiceName), telemetry.LogPkgStdFlags),
	)

	ticker := tickers.NewTimeTicker[[]common.UpkeepPayload](recoveryFinalizationInterval, recoveryObserver, func(ctx context.Context, _ time.Time) (tickers.Tick[[]common.UpkeepPayload], error) {
		return coordinatedProposalsTick{
			logger:    logger,
			builder:   builder,
			q:         proposalQ,
			utype:     types.LogTrigger,
			batchSize: FinalRecoveryBatchSize,
		}, nil
	}, log.New(logger.Writer(), fmt.Sprintf("[%s | recovery-final-ticker]", telemetry.ServiceName), telemetry.LogPkgStdFlags))

	return ticker
}

// coordinatedProposalsTick is used to push proposals from the proposal queue to some observer
type coordinatedProposalsTick struct {
	logger    *log.Logger
	builder   common.PayloadBuilder
	q         types.ProposalQueue
	utype     types.UpkeepType
	batchSize int
}

func (t coordinatedProposalsTick) Value(ctx context.Context) ([]common.UpkeepPayload, error) {
	if t.q == nil {
		return nil, nil
	}

	proposals, err := t.q.Dequeue(t.utype, t.batchSize)
	if err != nil {
		return nil, fmt.Errorf("failed to dequeue from retry queue: %w", err)
	}
	t.logger.Printf("%d proposals returned from queue", len(proposals))

	builtPayloads, err := t.builder.BuildPayloads(ctx, proposals...)
	if err != nil {
		return nil, fmt.Errorf("failed to build payloads from proposals: %w", err)
	}
	payloads := []common.UpkeepPayload{}
	filtered := 0
	for _, p := range builtPayloads {
		if p.IsEmpty() {
			filtered++
			continue
		}
		payloads = append(payloads, p)
	}
	t.logger.Printf("%d payloads built from %d proposals, %d filtered", len(payloads), len(proposals), filtered)
	return payloads, nil
}

func newRecoveryProposalFlow(
	preProcessors []ocr2keepersv3.PreProcessor[common.UpkeepPayload],
	runner ocr2keepersv3.Runner,
	metadataStore types.MetadataStore,
	recoverableProvider common.RecoverableProvider,
	recoveryInterval time.Duration,
	stateUpdater common.UpkeepStateUpdater,
	logger *log.Logger,
) service.Recoverable {
	preProcessors = append(preProcessors, preprocessors.NewProposalFilterer(metadataStore, types.LogTrigger))
	postprocessors := postprocessors.NewCombinedPostprocessor(
		postprocessors.NewIneligiblePostProcessor(stateUpdater, logger),
		postprocessors.NewAddProposalToMetadataStorePostprocessor(metadataStore),
	)

	observer := ocr2keepersv3.NewRunnableObserver(
		preProcessors,
		postprocessors,
		runner,
		ObservationProcessLimit,
		log.New(logger.Writer(), fmt.Sprintf("[%s | recovery-proposal-observer]", telemetry.ServiceName), telemetry.LogPkgStdFlags),
	)

	return tickers.NewTimeTicker[[]common.UpkeepPayload](recoveryInterval, observer, func(ctx context.Context, _ time.Time) (tickers.Tick[[]common.UpkeepPayload], error) {
		return logRecoveryTick{logger: logger, logRecoverer: recoverableProvider}, nil
	}, log.New(logger.Writer(), fmt.Sprintf("[%s | recovery-proposal-ticker]", telemetry.ServiceName), telemetry.LogPkgStdFlags))
}

type logRecoveryTick struct {
	logRecoverer common.RecoverableProvider
	logger       *log.Logger
}

func (et logRecoveryTick) Value(ctx context.Context) ([]common.UpkeepPayload, error) {
	if et.logRecoverer == nil {
		return nil, nil
	}

	logs, err := et.logRecoverer.GetRecoveryProposals(ctx)

	et.logger.Printf("%d logs returned by log recoverer", len(logs))

	return logs, err
}
