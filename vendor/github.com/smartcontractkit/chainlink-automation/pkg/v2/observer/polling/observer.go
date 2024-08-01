/*
A polling observer is an Observer that continuously polls every block, samples
available upkeeps, and surfaces upkeeps that should be agreed on by a network.
*/
package polling

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-automation/internal/util"
	ocr2keepers "github.com/smartcontractkit/chainlink-automation/pkg/v2"
	"github.com/smartcontractkit/chainlink-automation/pkg/v2/observer"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

var ErrTooManyErrors = fmt.Errorf("too many errors in parallel worker process")

// UpkeepProvider is a dependency used by the observer to get upkeeps
// available for polling
type UpkeepProvider interface {
	GetActiveUpkeepIDs(context.Context) ([]ocr2keepers.UpkeepIdentifier, error)
}

// Encoder is a dependency that provides helper functions for results and keys
type Encoder interface {
	// Eligible returns whether or not the upkeep result should be performed
	Eligible(ocr2keepers.UpkeepResult) (bool, error)
	// MakeUpkeepKey combines a block and upkeep id into an upkeep key. This
	// will probably go away with a more structured static upkeep type.
	MakeUpkeepKey(ocr2keepers.BlockKey, ocr2keepers.UpkeepIdentifier) ocr2keepers.UpkeepKey
	// SplitUpkeepKey ...
	SplitUpkeepKey(ocr2keepers.UpkeepKey) (ocr2keepers.BlockKey, ocr2keepers.UpkeepIdentifier, error)
	// Detail is a temporary value that provides upkeep key and gas to perform.
	// A better approach might be needed here.
	Detail(ocr2keepers.UpkeepResult) (ocr2keepers.UpkeepKey, uint32, error)
}

// Runner is a dependency that provides the caching and parallelization
// for checking eligibility of upkeeps
type Runner interface {
	CheckUpkeep(context.Context, bool, ...ocr2keepers.UpkeepKey) ([]ocr2keepers.UpkeepResult, error)
}

// HeadProvider is a dependency for a block data source channel
type HeadProvider interface {
	HeadTicker() chan ocr2keepers.BlockKey
}

// Service ...
type Service interface {
	Start()
	Stop()
}

// Shuffler ...
type Shuffler[T any] interface {
	Shuffle([]T) []T
}

// Ratio is an interface that provides functions to calculate a ratio of a given
// input
type Ratio interface {
	// OfInt should return n out of x such that n/x ~ r (ratio)
	OfInt(int) int
	fmt.Stringer
}

// NewPollingObserver ...
func NewPollingObserver(
	logger *log.Logger,
	src UpkeepProvider,
	heads HeadProvider,
	runner Runner,
	encoder Encoder,
	ratio Ratio,
	maxSamplingDuration time.Duration, // maximum amount of time allowed for RPC calls per head
	coord ocr2keepers.Coordinator,
	mercuryLookup bool,
) *PollingObserver {
	ob := &PollingObserver{
		stopCh:           make(chan struct{}),
		logger:           logger,
		samplingDuration: maxSamplingDuration,
		shuffler:         util.Shuffler[ocr2keepers.UpkeepKey]{Source: util.NewCryptoRandSource()}, // use crypto/rand shuffling for true random
		ratio:            ratio,
		stager:           &stager{},
		coordinator:      coord,
		src:              src,
		runner:           runner,
		encoder:          encoder,
		heads:            heads,
		mercuryLookup:    mercuryLookup,
	}

	// make all go-routines started by this entity automatically recoverable
	// on panics
	ob.services = []Service{
		util.NewRecoverableService(&observer.SimpleService{F: ob.runHeadTasks, C: func() { close(ob.stopCh) }}, logger),
	}

	return ob
}

type PollingObserver struct {
	stopCh    services.StopChan
	startOnce sync.Once
	stopOnce  sync.Once

	// static values provided by constructor
	samplingDuration time.Duration // limits time spent processing a single block
	ratio            Ratio         // ratio for limiting sample size

	// initialized components inside a constructor
	services []Service
	stager   *stager

	// dependency interfaces required by the polling observer
	logger      *log.Logger
	heads       HeadProvider                    // provides new blocks to be operated on
	coordinator ocr2keepers.Coordinator         // key status coordinator tracks in-flight status
	shuffler    Shuffler[ocr2keepers.UpkeepKey] // provides shuffling logic for upkeep keys

	src           UpkeepProvider
	runner        Runner
	encoder       Encoder
	mercuryLookup bool
}

// Observe implements the Observer interface and provides a slice of identifiers
// that were observed to be performable along with the block at which they were
// observed. All ids that are pending are filtered out.
func (o *PollingObserver) Observe() (ocr2keepers.BlockKey, []ocr2keepers.UpkeepIdentifier, error) {
	bl, ids := o.stager.get()
	filteredIDs := make([]ocr2keepers.UpkeepIdentifier, 0, len(ids))

	for _, id := range ids {
		key := o.encoder.MakeUpkeepKey(bl, id)

		if pending, err := o.coordinator.IsPending(key); pending || err != nil {
			if err != nil {
				o.logger.Printf("error checking pending state for '%s': %s", key, err)
			} else {
				o.logger.Printf("filtered out key '%s'", key)
			}

			continue
		}

		filteredIDs = append(filteredIDs, id)
	}

	return bl, filteredIDs, nil
}

// Start will start all required internal services. Calling this function again
// after the first is a noop.
func (o *PollingObserver) Start() {
	o.startOnce.Do(func() {
		for _, svc := range o.services {
			o.logger.Printf("PollingObserver service started")

			svc.Start()
		}
	})
}

// Stop will stop all internal services allowing the observer to exit cleanly.
func (o *PollingObserver) Close() error {
	o.stopOnce.Do(func() {
		for _, svc := range o.services {
			o.logger.Printf("PollingObserver service stopped")

			svc.Stop()
		}
	})

	return nil
}

func (o *PollingObserver) runHeadTasks() error {
	ctx, cancel := o.stopCh.NewCtx()
	defer cancel()
	ch := o.heads.HeadTicker()
	for {
		select {
		case bl := <-ch:
			// run sampling with latest head, the head ticker will drop heads
			// if the following process blocks for an extended period of time
			o.processLatestHead(ctx, bl)
		case <-ctx.Done():
			o.logger.Printf("PollingObserver.runHeadTasks ctx done")

			return ctx.Err()
		}
	}
}

// processLatestHead performs checking upkeep logic for all eligible keys of the given head
func (o *PollingObserver) processLatestHead(ctx context.Context, blockKey ocr2keepers.BlockKey) {
	// limit the context timeout to configured value
	ctx, cancel := context.WithTimeout(ctx, o.samplingDuration)
	defer cancel()
	var (
		keys []ocr2keepers.UpkeepKey
		ids  []ocr2keepers.UpkeepIdentifier
		err  error
	)
	o.logger.Printf("PollingObserver.processLatestHead")

	// Get only the active upkeeps from the id provider. This should not include
	// any cancelled upkeeps.
	if ids, err = o.src.GetActiveUpkeepIDs(ctx); err != nil {
		o.logger.Printf("%s: failed to get active upkeep ids from registry for sampling", err)
		return
	}

	o.logger.Printf("%d active upkeep ids found in registry", len(ids))

	keys = make([]ocr2keepers.UpkeepKey, len(ids))
	for i, id := range ids {
		keys[i] = o.encoder.MakeUpkeepKey(blockKey, id)
	}

	// reduce keys to ratio size and shuffle. this can return a nil array.
	// in that case we have no keys so return.
	if keys = o.shuffleAndSliceKeysToRatio(keys); keys == nil {
		o.logger.Printf("PollingObserver.processLatestHead shuffleAndSliceKeysToRatio returned nil keys")

		return
	}

	o.stager.prepareBlock(blockKey)
	o.logger.Printf("PollingObserver.processLatestHead prepared block")

	// run checkupkeep on all keys. an error from this function should
	// bubble up.
	results, err := o.runner.CheckUpkeep(ctx, o.mercuryLookup, keys...)
	if err != nil {
		o.logger.Printf("%s: failed to parallel check upkeeps", err)
		return
	}

	for _, res := range results {
		eligible, err := o.encoder.Eligible(res)
		if err != nil {
			o.logger.Printf("error testing result eligibility: %s", err)
			continue
		}

		if !eligible {
			continue
		}

		key, _, err := o.encoder.Detail(res)
		if err != nil {
			o.logger.Printf("error getting result detail: %s", err)
			continue
		}

		_, id, err := o.encoder.SplitUpkeepKey(key)
		if err != nil {
			o.logger.Printf("error splitting upkeep key: %s", err)
		}

		o.stager.prepareIdentifier(id)
	}

	// advance the staged block/upkeep id list to the next in line
	o.stager.advance()
	o.logger.Printf("PollingObserver.processLatestHead advanced stager")
}

func (o *PollingObserver) shuffleAndSliceKeysToRatio(keys []ocr2keepers.UpkeepKey) []ocr2keepers.UpkeepKey {
	keys = o.shuffler.Shuffle(keys)
	size := o.ratio.OfInt(len(keys))

	if len(keys) == 0 || size <= 0 {
		o.logger.Printf("PollingObserver.shuffleAndSliceKeysToRatio returning nil")
		return nil
	}

	o.logger.Printf("PollingObserver.shuffleAndSliceKeysToRatio returning %d keys", len(keys[:size]))

	return keys[:size]
}

type stager struct {
	currentIDs   []ocr2keepers.UpkeepIdentifier
	currentBlock ocr2keepers.BlockKey
	nextIDs      []ocr2keepers.UpkeepIdentifier
	nextBlock    ocr2keepers.BlockKey
	sync.RWMutex
}

func (s *stager) prepareBlock(block ocr2keepers.BlockKey) {
	s.Lock()
	defer s.Unlock()

	s.nextBlock = block
}

func (s *stager) prepareIdentifier(id ocr2keepers.UpkeepIdentifier) {
	s.Lock()
	defer s.Unlock()

	if s.nextIDs == nil {
		s.nextIDs = []ocr2keepers.UpkeepIdentifier{}
	}

	s.nextIDs = append(s.nextIDs, id)
}

func (s *stager) advance() {
	s.Lock()
	defer s.Unlock()

	s.currentBlock = s.nextBlock
	s.currentIDs = make([]ocr2keepers.UpkeepIdentifier, len(s.nextIDs))

	copy(s.currentIDs, s.nextIDs)

	s.nextIDs = make([]ocr2keepers.UpkeepIdentifier, 0)
}

func (s *stager) get() (ocr2keepers.BlockKey, []ocr2keepers.UpkeepIdentifier) {
	s.RLock()
	defer s.RUnlock()

	return s.currentBlock, s.currentIDs
}
