package runner

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/smartcontractkit/chainlink-automation/internal/util"
	pkgutil "github.com/smartcontractkit/chainlink-automation/pkg/util"
	ocr2keepers "github.com/smartcontractkit/chainlink-automation/pkg/v2"
)

var ErrTooManyErrors = fmt.Errorf("too many errors in parallel worker process")

type Registry interface {
	CheckUpkeep(context.Context, bool, ...ocr2keepers.UpkeepKey) ([]ocr2keepers.UpkeepResult, error)
}

type Encoder interface {
	// Eligible determines if an upkeep is eligible or not. This allows an
	// upkeep result to be abstract and only the encoder is able and responsible
	// for decoding it.
	Eligible(ocr2keepers.UpkeepResult) (bool, error)
	// Detail is a temporary value that provides upkeep key and gas to perform.
	// A better approach might be needed here.
	Detail(ocr2keepers.UpkeepResult) (ocr2keepers.UpkeepKey, uint32, error)
	// SplitUpkeepKey ...
	SplitUpkeepKey(ocr2keepers.UpkeepKey) (ocr2keepers.BlockKey, ocr2keepers.UpkeepIdentifier, error)
}

// Runner ...
type Runner struct {
	// injected dependencies
	logger   *log.Logger
	registry Registry
	encoder  Encoder

	// initialized by the constructor
	workers      *pkgutil.WorkerGroup[[]ocr2keepers.UpkeepResult] // parallelizer for RPC calls
	cache        *pkgutil.Cache[ocr2keepers.UpkeepResult]
	cacheCleaner *pkgutil.IntervalCacheCleaner[ocr2keepers.UpkeepResult]

	// configurations
	workerBatchLimit int // the maximum number of items in RPC batch call

	// run state data
	running atomic.Bool
}

// NewRunner ...
func NewRunner(
	logger *log.Logger,
	registry Registry,
	encoder Encoder,
	workers int, // maximum number of workers in worker group
	workerQueueLength int, // size of worker queue; set to approximately the number of items expected in workload
	cacheExpire time.Duration,
	cacheClean time.Duration,
) (*Runner, error) {
	return &Runner{
		logger:           logger,
		registry:         registry,
		encoder:          encoder,
		workers:          pkgutil.NewWorkerGroup[[]ocr2keepers.UpkeepResult](workers, workerQueueLength),
		cache:            pkgutil.NewCache[ocr2keepers.UpkeepResult](cacheExpire),
		cacheCleaner:     pkgutil.NewIntervalCacheCleaner[ocr2keepers.UpkeepResult](cacheClean),
		workerBatchLimit: 10,
	}, nil
}

func (o *Runner) CheckUpkeep(ctx context.Context, mercuryEnabled bool, keys ...ocr2keepers.UpkeepKey) ([]ocr2keepers.UpkeepResult, error) {
	r, err := o.parallelCheck(ctx, mercuryEnabled, keys)
	if err != nil {
		return nil, err
	}

	return r.Values(), nil
}

func (o *Runner) Start() error {
	if !o.running.Load() {
		go o.cacheCleaner.Run(o.cache)
		o.running.Swap(true)
	}

	return nil
}

func (o *Runner) Close() error {
	if o.running.Load() {
		o.cacheCleaner.Stop()
		o.workers.Stop()
		o.running.Swap(false)
	}

	return nil
}

// parallelCheck should be satisfied by the Runner
func (o *Runner) parallelCheck(ctx context.Context, mercuryEnabled bool, keys []ocr2keepers.UpkeepKey) (*Result, error) {
	result := NewResult()

	if len(keys) == 0 {
		return result, nil
	}

	toRun := make([]ocr2keepers.UpkeepKey, 0, len(keys))
	for _, key := range keys {

		// if in cache, add to result
		if res, ok := o.cache.Get(string(key)); ok {
			result.Add(res)
			continue
		}

		// else add to lookup job
		toRun = append(toRun, key)
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
		o.wrapWorkerFunc(mercuryEnabled),
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

func (o *Runner) wrapWorkerFunc(mercuryEnabled bool) func(context.Context, []ocr2keepers.UpkeepKey) ([]ocr2keepers.UpkeepResult, error) {
	return func(ctx context.Context, keys []ocr2keepers.UpkeepKey) ([]ocr2keepers.UpkeepResult, error) {
		start := time.Now()

		// perform check and update cache with result
		checkResults, err := o.registry.CheckUpkeep(ctx, mercuryEnabled, keys...)
		if err != nil {
			err = fmt.Errorf("%w: failed to check upkeep keys: %s", err, keys)
		} else {
			o.logger.Printf("check %d upkeeps took %dms to perform", len(keys), time.Since(start)/time.Millisecond)

			for _, result := range checkResults {
				ok, err := o.encoder.Eligible(result)
				if err != nil {
					o.logger.Printf("eligibility check error: %s", err)
					continue
				}

				key, _, _ := o.encoder.Detail(result)
				if ok {
					o.logger.Printf("upkeep ready to perform for key %s", key)
				} else {
					o.logger.Printf("upkeep not ready to perform for key %s", key)
				}
			}
		}

		return checkResults, err
	}
}

func (o *Runner) wrapAggregate(r *Result) func([]ocr2keepers.UpkeepResult, error) {
	return func(result []ocr2keepers.UpkeepResult, err error) {
		if err == nil {
			r.AddSuccesses(1)

			for _, res := range result {
				key, _, _ := o.encoder.Detail(res)
				o.cache.Set(string(key), res, pkgutil.DefaultCacheExpiration)
				r.Add(res)
			}
		} else {
			r.SetErr(err)
			o.logger.Printf("error received from worker result: %s", err)
			r.AddFailures(1)
		}
	}
}
