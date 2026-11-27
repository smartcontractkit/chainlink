package queue

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3_1/plugins/queue/queuelimits"
)

// defaultBlobOpTimeout bounds each blob fetch or broadcast, so a stalled blob store or peer cannot block a plugin
// method indefinitely. It mirrors the vault plugin's blobBroadcastTimeout.
const defaultBlobOpTimeout = 2 * time.Second

// blobCache is a short-lived cache for coordinating concurrent fetches of event blobs, and of splitting them up.
type blobCache struct {
	fetcher ocr3_1types.BlobFetcher
	// timeout bounds each fetch. Zero means defaultBlobOpTimeout.
	timeout time.Duration

	mu       sync.RWMutex
	cache    map[string]func() ([]byte, error)   // blob handle -> blob
	elements map[string]func() ([][]byte, error) // blob handle -> its events, split up
}

func (c *blobCache) opTimeout() time.Duration {
	if c.timeout > 0 {
		return c.timeout
	}
	return defaultBlobOpTimeout
}

// isCancellation reports whether an error is a context ending.
func isCancellation(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

// uncancelled returns what get returns, except that a result which is a context ending is not kept for others. The
// result of a fetch is shared by every caller, but it was made with the context of the first one, so if that context
// ended, its cancellation is not the answer for a caller whose context did not. evict removes the result from the cache,
// and the caller then tries again, once, with its own context. A caller whose own context ended gets the error.
func uncancelled[T any](ctx context.Context, get func() (T, error), evict func()) (T, error) {
	for attempt := 0; ; attempt++ {
		v, err := get()
		if err == nil || !isCancellation(err) {
			return v, err
		}
		evict()
		if ctx.Err() != nil || attempt > 0 {
			return v, err
		}
	}
}

// fetchBlob fetches the blob once however many callers want it, and shares the result, whether a blob or an error. Not
// an error that is a context ending: see uncancelled.
func (c *blobCache) fetchBlob(ctx context.Context, handle []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, c.opTimeout())
	defer cancel()
	return uncancelled(ctx, func() ([]byte, error) { return c.fetchBlobCached(ctx, handle) }, func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.cache, string(handle))
	})
}

func (c *blobCache) fetchBlobCached(ctx context.Context, handle []byte) ([]byte, error) {
	c.mu.RLock()
	if fn, ok := c.cache[string(handle)]; ok {
		c.mu.RUnlock()
		return fn()
	}

	c.mu.RUnlock()

	c.mu.Lock()
	if fn, ok := c.cache[string(handle)]; ok {
		c.mu.Unlock()
		return fn()
	}

	fn := sync.OnceValues(func() ([]byte, error) {
		var blobHandle ocr3_1types.BlobHandle
		if err := blobHandle.UnmarshalBinary(handle); err != nil {
			return nil, fmt.Errorf("failed to unmarshal blob handle: %w", err)
		}

		return c.fetcher.FetchBlob(ctx, blobHandle)
	})

	c.cache[string(handle)] = fn
	c.mu.Unlock()
	return fn()
}

// fetchElements fetches the blob and splits it into its events once, however many events are read from it. Finding an
// event by scanning from the start of a blob each time is quadratic in its events, which a peer could exploit by padding
// a blob with empty events.
func (c *blobCache) fetchElements(ctx context.Context, handle []byte) ([][]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, c.opTimeout())
	defer cancel()
	return uncancelled(ctx, func() ([][]byte, error) { return c.fetchElementsCached(ctx, handle) }, func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.elements, string(handle))
	})
}

func (c *blobCache) fetchElementsCached(ctx context.Context, handle []byte) ([][]byte, error) {
	c.mu.Lock()
	fn, ok := c.elements[string(handle)]
	if !ok {
		fn = sync.OnceValues(func() ([][]byte, error) {
			blob, err := c.fetchBlob(ctx, handle)
			if err != nil {
				return nil, err
			}
			return splitEvents(blob, maxEventsPerBlob)
		})
		if c.elements == nil {
			c.elements = make(map[string]func() ([][]byte, error))
		}
		c.elements[string(handle)] = fn
	}
	c.mu.Unlock()
	return fn()
}

// fetchCoordinatedEventFromBlob fetches the blob and returns the idx-th event in it, and its marshaled bytes, without
// unmarshaling any other events.
func fetchCoordinatedEventFromBlob(ctx context.Context, fetcher *blobCache, handle []byte, idx int) (*pb.CoordinatedEvent, []byte, error) {
	elements, err := fetcher.fetchElements(ctx, handle)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch blob: %w", err)
	}
	if idx < 0 || idx >= len(elements) {
		return nil, nil, fmt.Errorf("index %d out of range", idx)
	}

	var event pb.CoordinatedEvent
	if err := proto.Unmarshal(elements[idx], &event); err != nil {
		return nil, elements[idx], fmt.Errorf("failed to unmarshal event: %w", err)
	}
	return &event, elements[idx], nil
}

const coordinatedEventsFieldNumber = queuelimits.CoordinatedEventsFieldNumber // pb.CoordinatedEventList.Events

// maxEventFramingOverhead is the most framing (field tag + length prefix) an event adds as an element of a blob.
var maxEventFramingOverhead = queuelimits.MaxEventFramingOverhead(coordinatedEventsFieldNumber)

// maxEventBytes returns the largest encoded event (without timestamps) guaranteed to fit alone in a blob.
// It is conservative by at most a few bytes.
func (c *consensusQueuePlugin) maxEventBytes() int {
	return max(0, c.limits.MaxBlobPayloadBytes-maxEventFramingOverhead)
}

// maxEventsPerBlob is the most events a blob may hold. An observation has at most maxObservedEvents events in all of
// its blobs, so an honest blob has no more. Anything else only makes reading it costly.
const maxEventsPerBlob = maxObservedEvents

// splitEvents splits a blob into its marshaled events, without unmarshaling any of them. It rejects a blob with more
// than limit events.
func splitEvents(b []byte, limit int) ([][]byte, error) {
	var events [][]byte
	for len(b) > 0 {
		num, typ, n := protowire.ConsumeTag(b)
		if n < 0 {
			return nil, protowire.ParseError(n)
		}

		b = b[n:]
		if num != coordinatedEventsFieldNumber {
			n := protowire.ConsumeFieldValue(num, typ, b)
			if n < 0 {
				return nil, protowire.ParseError(n)
			}

			b = b[n:]
			continue
		}

		if typ != protowire.BytesType {
			return nil, fmt.Errorf("expected bytes type for field %d, got %d", coordinatedEventsFieldNumber, typ)
		}

		v, m := protowire.ConsumeBytes(b)
		if m < 0 {
			return nil, protowire.ParseError(m)
		}

		if len(events) == limit {
			return nil, fmt.Errorf("blob has more than %d events", limit)
		}

		events = append(events, v)
		b = b[m:]
	}
	return events, nil
}

const (
	// minBlobLifetime is the least time a blob is kept past the events in it expiring. Peers reap expired blobs
	// asynchronously, and nodes that are behind, or that never received the blob, need time to fetch it.
	minBlobLifetime = queuelimits.MinBlobLifetime
	// minBlobLifetimeRounds is the fewest rounds a blob is kept, however the round interval compares to minBlobLifetime.
	// The events in it are executed in Committed for its round and in the Observation of the next, which both fetch it.
	minBlobLifetimeRounds = 3
)

// blobExpirySeqNr returns the sequence number at which a blob broadcast in round seqNr may expire, so that it outlives
// the time remaining until its last event expires, plus minBlobLifetime.
//
// The time is converted to rounds by dividing by the estimated round interval, which in OCR3.1 is DeltaRound: the
// leader waits at least that long between rounds. Rounds are normally longer, so this over-counts the rounds in the
// remaining time and keeps the blob a little longer than needed. That is the safe direction, as an early expiry
// loses events but a late one only holds storage. Partial rounds are rounded up. If the interval is unknown, the blob
// is kept for minBlobLifetimeRounds.
func blobExpirySeqNr(seqNr uint64, remaining, roundDur time.Duration) uint64 {
	rounds := uint64(minBlobLifetimeRounds)
	if roundDur > 0 {
		lifetime := max(0, remaining) + minBlobLifetime
		rounds = max(rounds, uint64((lifetime+roundDur-1)/roundDur)) //nolint:gosec // G115: lifetime and roundDur are positive
	}

	return seqNr + rounds
}

// blobPrefix returns how many of the events (from the first) fit in one blob.
func (c *consensusQueuePlugin) blobPrefix(events []*pb.CoordinatedEvent) int {
	protoOpts := proto.MarshalOptions{Deterministic: true}
	return sort.Search(len(events), func(i int) bool {
		return protoOpts.Size(&pb.CoordinatedEventList{Events: events[:i+1]}) > c.limits.MaxBlobPayloadBytes
	})
}

// latestDeadline returns the latest deadline among the events, or the zero time if none has one.
func latestDeadline(events []*pb.ConsensusQueueEvent) time.Time {
	var latest time.Time
	for _, event := range events {
		if t := event.GetDeadline().AsTime(); t.After(latest) {
			latest = t
		}
	}
	return latest
}

// broadcastBlob broadcasts a blob with all of the events, which must fit, and returns its marshaled handle.
// maxDeadline is the latest deadline of the events, which the caller must take from the observation: the events
// themselves carry no timestamps, as they are stripped for consistent hashing and left out of the blob.
func (c *consensusQueuePlugin) broadcastBlob(ctx context.Context, broadcaster ocr3_1types.BlobBroadcaster, events []*pb.CoordinatedEvent, maxDeadline time.Time, seqNr uint64) ([]byte, error) {
	// An event is kept until its deadline plus the grace period, so its blob must be too.
	grace, err := c.expirationDur(ctx)
	if err != nil {
		c.lggr.Warnw("Failed to get expiration grace duration - using default", "grace", grace, "err", err)
	}

	expSeqNr := blobExpirySeqNr(seqNr, time.Until(maxDeadline.Add(max(0, grace))), c.onchainCfg.EstimatedRoundInterval)

	data, err := proto.MarshalOptions{Deterministic: true}.Marshal(&pb.CoordinatedEventList{Events: events})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal blob: %w", err)
	}

	bh, err := broadcaster.BroadcastBlob(ctx, data, ocr3_1types.BlobExpirationHintSequenceNumber{SeqNr: expSeqNr})
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast blob: %w", err)
	}

	return bh.MarshalBinary()
}
