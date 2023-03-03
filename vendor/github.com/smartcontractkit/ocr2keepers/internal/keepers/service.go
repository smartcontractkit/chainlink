package keepers

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/smartcontractkit/ocr2keepers/pkg/chain"
	"github.com/smartcontractkit/ocr2keepers/pkg/types"
	"github.com/smartcontractkit/ocr2keepers/pkg/util"
)

// maxWorkersBatchSize is the value of max workers batch size
const maxWorkersBatchSize = 10

var (
	ErrTooManyErrors          = fmt.Errorf("too many errors in parallel worker process")
	ErrSamplingNotInitialized = fmt.Errorf("sampling not initialized")
)

type onDemandUpkeepService struct {
	logger           *log.Logger
	ratio            sampleRatio
	headSubscriber   types.HeadSubscriber
	registry         types.Registry
	shuffler         shuffler[types.UpkeepIdentifier]
	cache            *util.Cache[types.UpkeepResult]
	cacheCleaner     *util.IntervalCacheCleaner[types.UpkeepResult]
	samplingResults  samplingUpkeepsResults
	samplingDuration time.Duration
	workers          *util.WorkerGroup[types.UpkeepResults]
	ctx              context.Context
	cancel           context.CancelFunc
}

// newOnDemandUpkeepService provides an object that implements the UpkeepService
// by running a worker pool that makes RPC network calls every time upkeeps
// need to be sampled. This variant has limitations in how quickly large numbers
// of upkeeps can be checked. Be aware that network calls are not rate limited
// from this service.
func newOnDemandUpkeepService(
	ratio sampleRatio,
	headSubscriber types.HeadSubscriber,
	registry types.Registry,
	logger *log.Logger,
	samplingDuration time.Duration,
	cacheExpire time.Duration,
	cacheClean time.Duration,
	workers int,
	workerQueueLength int,
) *onDemandUpkeepService {
	ctx, cancel := context.WithCancel(context.Background())
	s := &onDemandUpkeepService{
		logger:           logger,
		ratio:            ratio,
		headSubscriber:   headSubscriber,
		registry:         registry,
		samplingDuration: samplingDuration,
		shuffler:         new(cryptoShuffler[types.UpkeepIdentifier]),
		cache:            util.NewCache[types.UpkeepResult](cacheExpire),
		cacheCleaner:     util.NewIntervalCacheCleaner[types.UpkeepResult](cacheClean),
		workers:          util.NewWorkerGroup[types.UpkeepResults](workers, workerQueueLength),
		ctx:              ctx,
		cancel:           cancel,
	}

	// stop the cleaner go-routine once the upkeep service is no longer reachable
	runtime.SetFinalizer(s, func(srv *onDemandUpkeepService) { srv.stop() })

	// start background services
	s.start()

	return s
}

var _ upkeepService = (*onDemandUpkeepService)(nil)

func (s *onDemandUpkeepService) SampleUpkeeps(_ context.Context, filters ...func(types.UpkeepKey) bool) (types.BlockKey, types.UpkeepResults, error) {
	blockKey, results, ok := s.samplingResults.get()
	if !ok {
		return nil, nil, ErrSamplingNotInitialized
	}

	filteredResults := make(types.UpkeepResults, 0, len(results))

EachKey:
	for _, result := range results {
		for _, filter := range filters {
			if !filter(result.Key) {
				s.logger.Printf("filtered out key during SampleUpkeeps '%s'", result.Key)
				continue EachKey
			}
		}

		filteredResults = append(filteredResults, result)
	}

	return blockKey, filteredResults, nil
}

func (s *onDemandUpkeepService) CheckUpkeep(ctx context.Context, keys ...types.UpkeepKey) (types.UpkeepResults, error) {
	var (
		wg                sync.WaitGroup
		results           = make([]types.UpkeepResult, len(keys))
		nonCachedKeysLock sync.Mutex
		nonCachedKeysIdxs = make([]int, 0, len(keys))
		nonCachedKeys     = make([]types.UpkeepKey, 0, len(keys))
	)

	for i, key := range keys {
		wg.Add(1)
		go func(i int, key types.UpkeepKey) {
			// the cache is a collection of keys (block & id) that map to cached
			// results. if the same upkeep is checked at a block that has already been
			// checked, return the cached result
			if result, cached := s.cache.Get(key.String()); cached {
				results[i] = result
			} else {
				nonCachedKeysLock.Lock()
				nonCachedKeysIdxs = append(nonCachedKeysIdxs, i)
				nonCachedKeys = append(nonCachedKeys, key)
				nonCachedKeysLock.Unlock()
			}
			wg.Done()
		}(i, key)
	}

	wg.Wait()

	// All keys are cached
	if len(nonCachedKeys) == 0 {
		return results, nil
	}

	// check upkeep at block number in key
	// return result including performData
	checkResults, err := s.registry.CheckUpkeep(ctx, nonCachedKeys...)
	if err != nil {
		return nil, fmt.Errorf("%w: service failed to check upkeep from registry", err)
	}

	// Cache results
	for i, u := range checkResults {
		s.cache.Set(keys[nonCachedKeysIdxs[i]].String(), u, util.DefaultCacheExpiration)
		results[nonCachedKeysIdxs[i]] = u
	}

	return results, nil
}

func (s *onDemandUpkeepService) start() {
	// TODO: if this process panics, restart it
	go s.cacheCleaner.Run(s.cache)
	go func() {
		ch := s.headSubscriber.HeadTicker()
		for {
			select {
			case head := <-ch:
				// run with new head
				s.processLatestHead(s.ctx, head)
			case <-s.ctx.Done():
				return
			}
		}
	}()
}

func (s *onDemandUpkeepService) stop() {
	s.cancel()
	s.workers.Stop()
	s.cacheCleaner.Stop()
}

// processLatestHead performs checking upkeep logic for all eligible keys of the given head
func (s *onDemandUpkeepService) processLatestHead(ctx context.Context, blockKey types.BlockKey) {
	ctx, cancel := context.WithTimeout(ctx, s.samplingDuration)
	defer cancel()

	// Get only the active upkeeps from the contract. This should not include
	// any cancelled upkeeps.
	keys, err := s.registry.GetActiveUpkeepIDs(ctx)
	if err != nil {
		s.logger.Printf("%s: failed to get upkeeps from registry for sampling", err)
		return
	}

	s.logger.Printf("%d active upkeep keys found in registry", len(keys))

	// select x upkeeps at random from set
	keys = s.shuffler.Shuffle(keys)
	sampleSize := s.ratio.OfInt(len(keys))

	s.logger.Printf("%d results selected by provided ratio %s", sampleSize, s.ratio)
	if sampleSize < 0 {
		s.logger.Printf("sample size is too small: %d", sampleSize)
		return
	}

	var upkeepKeys []types.UpkeepKey
	for _, k := range keys {
		upkeepKeys = append(upkeepKeys, chain.NewUpkeepKeyFromBlockAndID(blockKey, k))
	}

	upkeepResults, err := s.parallelCheck(ctx, upkeepKeys[:sampleSize])
	if err != nil {
		s.logger.Printf("%s: failed to parallel check upkeeps", err)
		return
	}

	s.samplingResults.set(blockKey, upkeepResults)
}

func (s *onDemandUpkeepService) parallelCheck(ctx context.Context, keys []types.UpkeepKey) (types.UpkeepResults, error) {
	samples := newSyncedArray[types.UpkeepResult]()

	if len(keys) == 0 {
		return samples.Values(), nil
	}

	var wResults workerResults

	// go through keys and check the cache first
	// if an item doesn't exist on the cache, send the items to the worker threads
	filteredKeys, cacheHits := s.filterFromCache(keys, samples)

	// Create batches from the given keys.
	// Max keyBatchSize items in the batch.
	util.RunJobs(
		ctx,
		s.workers,
		createBatches(filteredKeys, maxWorkersBatchSize),
		s.wrapWorkerFunc(),
		s.wrapAggregate(&wResults, samples),
	)

	if wResults.Total() == 0 {
		s.logger.Printf("no network calls were made for this sampling set")
	} else {
		s.logger.Printf("worker call success rate: %.2f; failure rate: %.2f; total calls %d", wResults.SuccessRate(), wResults.FailureRate(), wResults.Total())
	}

	s.logger.Printf("sampling cache hit ratio %d/%d", cacheHits, len(keys))

	// multiple network calls can result in an error while some can be successful
	// in the case that all workers encounter an error, bubble this up as a hard
	// failure of the process.
	if wResults.Total() > 0 && wResults.Total() == wResults.Failures() && wResults.LastErr() != nil {
		return samples.Values(), fmt.Errorf("%w: last error encounter by worker was '%s'", ErrTooManyErrors, wResults.LastErr())
	}

	return samples.Values(), nil
}

func (s *onDemandUpkeepService) filterFromCache(keys []types.UpkeepKey, samples *syncedArray[types.UpkeepResult]) ([]types.UpkeepKey, int) {
	var keysToSend = make([]types.UpkeepKey, 0, len(keys))
	var cacheHits int

	for _, key := range keys {
		// no RPC lookups need to be done if a result has already been cached
		result, cached := s.cache.Get(key.String())
		if cached {
			cacheHits++
			if result.State == types.Eligible {
				samples.Append(result)
			}
			continue
		}

		// Add key to the slice that is going to be sent to the worker queue
		keysToSend = append(keysToSend, key)
	}

	return keysToSend, cacheHits
}

func (s *onDemandUpkeepService) wrapAggregate(r *workerResults, sa *syncedArray[types.UpkeepResult]) func(types.UpkeepResults, error) {
	return func(result types.UpkeepResults, err error) {
		if err == nil {
			r.AddSuccess(1)

			// Cache results
			for i := range result {
				res := result[i]
				s.cache.Set(string(res.Key.String()), res, util.DefaultCacheExpiration)
				if res.State == types.Eligible {
					sa.Append(res)
				}
			}
		} else {
			r.SetLastErr(err)
			s.logger.Printf("error received from worker result: %s", err)
			r.AddFailure(1)
		}
	}
}

func (s *onDemandUpkeepService) wrapWorkerFunc() func(context.Context, []types.UpkeepKey) (types.UpkeepResults, error) {
	return func(ctx context.Context, keys []types.UpkeepKey) (types.UpkeepResults, error) {
		keysStr := upkeepKeysToString(keys)
		start := time.Now()

		// perform check and update cache with result
		checkResults, err := s.registry.CheckUpkeep(ctx, keys...)
		if err != nil {
			err = fmt.Errorf("%w: failed to check upkeep keys: %s", err, keysStr)
		} else {
			s.logger.Printf("check %d upkeeps took %dms to perform", len(keys), time.Since(start)/time.Millisecond)

			for _, result := range checkResults {
				if result.State == types.Eligible {
					s.logger.Printf("upkeep ready to perform for key %s", result.Key)
				} else {
					s.logger.Printf("upkeep '%s' is not eligible with failure reason: %d", result.Key, result.FailureReason)
				}
			}
		}

		return checkResults, err
	}
}

type workerResults struct {
	success int
	failure int
	lastErr error
}

func (wr *workerResults) AddSuccess(amt int) {
	wr.success = wr.success + amt
}

func (wr *workerResults) Failures() int {
	return wr.failure
}

func (wr *workerResults) LastErr() error {
	return wr.lastErr
}

func (wr *workerResults) AddFailure(amt int) {
	wr.failure = wr.failure + amt
}

func (wr *workerResults) SetLastErr(err error) {
	wr.lastErr = err
}

func (wr *workerResults) Total() int {
	return wr.total()
}

func (wr *workerResults) total() int {
	return wr.success + wr.failure
}

func (wr *workerResults) SuccessRate() float64 {
	return float64(wr.success) / float64(wr.total())
}

func (wr *workerResults) FailureRate() float64 {
	return float64(wr.failure) / float64(wr.total())
}

type samplingUpkeepsResults struct {
	upkeepResults types.UpkeepResults
	blockKey      types.BlockKey
	ok            bool
	sync.Mutex
}

func (sur *samplingUpkeepsResults) set(blockKey types.BlockKey, results types.UpkeepResults) {
	sur.Lock()
	defer sur.Unlock()

	sur.upkeepResults = make(types.UpkeepResults, len(results))
	copy(sur.upkeepResults, results)
	sur.blockKey = blockKey
	sur.ok = true
}

func (sur *samplingUpkeepsResults) get() (types.BlockKey, types.UpkeepResults, bool) {
	sur.Lock()
	defer sur.Unlock()

	return sur.blockKey, sur.upkeepResults, sur.ok
}
