package runner

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink-automation/internal/util"
	pkgutil "github.com/smartcontractkit/chainlink-automation/pkg/util"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/telemetry"
)

const WorkerBatchLimit int = 10

var ErrTooManyErrors = fmt.Errorf("too many errors in parallel worker process")

// ensure that the runner implements the same interface it consumes to indicate
// the runner simply wraps the underlying runnable with extra features
var _ types.Runnable = &Runner{}

// Runner is a component that parallelizes calls to the provided runnable both
// by batching tasks to individual calls as well as using parallel threads to
// execute calls to the runnable. All results are cached such that the same
// input job from a previous run will provide a cached response instead of
// calling the runnable.
//
// The Runner is structured as a direct replacement where the runnable is used
// as a dependency.
type Runner struct {
	// injected dependencies
	logger   *log.Logger
	runnable types.Runnable
	// initialized by the constructor
	workers *pkgutil.WorkerGroup[[]ocr2keepers.CheckResult] // parallelizer
	cache   *pkgutil.Cache[ocr2keepers.CheckResult]         // result cache
	// configurations
	workerBatchLimit int // the maximum number of items in RPC batch call
	cacheGcInterval  time.Duration
	// run state data
	running atomic.Bool
	chClose chan struct{}
}

type RunnerConfig struct {
	// Workers is the maximum number of workers in worker group
	Workers int
	// WorkerQueueLength is size of worker queue; set to approximately the number of items expected in workload
	WorkerQueueLength int
	CacheExpire       time.Duration
	CacheClean        time.Duration
}

// NewRunner provides a new configured runner
func NewRunner(
	logger *log.Logger,
	runnable types.Runnable,
	conf RunnerConfig,
) (*Runner, error) {
	return &Runner{
		logger:           log.New(logger.Writer(), fmt.Sprintf("[%s | check-pipeline-runner]", telemetry.ServiceName), telemetry.LogPkgStdFlags),
		runnable:         runnable,
		workers:          pkgutil.NewWorkerGroup[[]ocr2keepers.CheckResult](conf.Workers, conf.WorkerQueueLength),
		cache:            pkgutil.NewCache[ocr2keepers.CheckResult](conf.CacheExpire),
		cacheGcInterval:  conf.CacheClean,
		workerBatchLimit: WorkerBatchLimit,
		chClose:          make(chan struct{}, 1),
	}, nil
}

// CheckUpkeeps accepts an array of payloads, splits the workload into separate
// threads, executes the underlying runnable, and returns all results from all
// threads. If previous runs were already completed for the same one or more
// payloads, results will be pulled from the cache where available.
func (o *Runner) CheckUpkeeps(ctx context.Context, payloads ...ocr2keepers.UpkeepPayload) ([]ocr2keepers.CheckResult, error) {
	r, err := o.parallelCheck(ctx, payloads)
	if err != nil {
		return nil, err
	}

	return r.Values(), nil
}

// Start starts up the cache cleaner
func (o *Runner) Start(_ context.Context) error {
	if o.running.Load() {
		return fmt.Errorf("already running")
	}

	o.running.Swap(true)
	o.logger.Println("starting service")

	go o.cache.Start(o.cacheGcInterval)

	<-o.chClose

	return nil
}

// Close stops the cache cleaner and the parallel worker process
func (o *Runner) Close() error {
	if !o.running.Load() {
		return fmt.Errorf("not running")
	}

	o.cache.Stop()
	o.workers.Stop()
	o.running.Swap(false)

	o.chClose <- struct{}{}

	return nil
}

// parallelCheck should be satisfied by the Runner
func (o *Runner) parallelCheck(ctx context.Context, payloads []ocr2keepers.UpkeepPayload) (*result[ocr2keepers.CheckResult], error) {
	result := newResult[ocr2keepers.CheckResult]()

	if len(payloads) == 0 {
		return result, nil
	}

	toRun := make([]ocr2keepers.UpkeepPayload, 0, len(payloads))
	for _, payload := range payloads {
		// if workID is in cache for the given trigger blocknum/hash, add to result directly
		if res, ok := o.cache.Get(payload.WorkID); ok &&
			(res.Trigger.BlockNumber == payload.Trigger.BlockNumber) &&
			(res.Trigger.BlockHash == payload.Trigger.BlockHash) {
			result.Add(res)
			continue
		}

		// else add to lookup job
		toRun = append(toRun, payload)
	}

	// no more to do
	if len(toRun) == 0 {
		return result, nil
	}

	// Create batches from the given keys.
	// Max keyBatchSize items in the batch.
	pkgutil.RunJobs(
		ctx,
		o.workers,
		util.Unflatten(toRun, o.workerBatchLimit),
		o.wrapWorkerFunc(),
		o.wrapAggregate(result),
	)

	if result.Total() == 0 {
		o.logger.Printf("no network calls were made for this sampling set")
	} else {
		o.logger.Printf("worker call success rate: %.2f; failure rate: %.2f; total calls %d", result.SuccessRate(), result.FailureRate(), result.Total())
	}

	// multiple network calls can result in an error while some can be successful
	// in the case that all workers encounter an error, bubble this up as a hard
	// failure of the process.
	if result.Total() > 0 && result.Total() == result.Failures() && result.Err() != nil {
		return nil, fmt.Errorf("%w: last error encounter by worker was '%s'", ErrTooManyErrors, result.Err())
	}

	return result, nil
}

func (o *Runner) wrapWorkerFunc() func(context.Context, []ocr2keepers.UpkeepPayload) ([]ocr2keepers.CheckResult, error) {
	return func(ctx context.Context, payloads []ocr2keepers.UpkeepPayload) ([]ocr2keepers.CheckResult, error) {
		start := time.Now()

		allPayloadKeys := make([]string, len(payloads))
		for i := range payloads {
			allPayloadKeys[i] = payloads[i].WorkID
		}

		// perform check and update cache with result
		checkResults, err := o.runnable.CheckUpkeeps(ctx, payloads...)
		if err != nil {
			err = fmt.Errorf("%w: failed to check upkeep payloads for ids '%s'", err, strings.Join(allPayloadKeys, ", "))
		} else {
			o.logger.Printf("check %d upkeeps took %dms to perform", len(payloads), time.Since(start)/time.Millisecond)
		}

		return checkResults, err
	}
}

func (o *Runner) wrapAggregate(r *result[ocr2keepers.CheckResult]) func([]ocr2keepers.CheckResult, error) {
	return func(results []ocr2keepers.CheckResult, err error) {
		if err == nil {
			r.AddSuccesses(1)

			for _, result := range results {
				// only add to the cache if pipeline was successful
				if result.PipelineExecutionState == 0 {
					c, ok := o.cache.Get(result.WorkID)
					if !ok || result.Trigger.BlockNumber > c.Trigger.BlockNumber {
						// Add to cache if the workID didn't exist before or if we got a result on a higher checkBlockNumber
						o.cache.Set(result.WorkID, result, pkgutil.DefaultCacheExpiration)
					}
				}

				r.Add(result)
			}
		} else {
			r.SetErr(err)
			o.logger.Printf("error received from worker result: %s", err)
			r.AddFailures(1)
		}
	}
}
