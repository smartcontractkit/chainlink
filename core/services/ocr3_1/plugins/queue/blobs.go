package queue

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
)

// blobCache is a short-lived cache for coordinating concurrent fetches of event blobs, and of splitting them up.
type blobCache struct {
	fetcher ocr3_1types.BlobFetcher

	mu       sync.RWMutex
	cache    map[string]func() ([]byte, error)   // blob handle -> blob
	elements map[string]func() ([][]byte, error) // blob handle -> its events, split up
}

func (c *blobCache) fetchBlob(ctx context.Context, handle []byte) ([]byte, error) {
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

const coordinatedEventsFieldNumber = 1 // pb.CoordinatedEventList.Events

// maxEventFramingOverhead is the most framing (field tag + length prefix) an event adds as an element of a blob.
// The length prefix is at most as long as the varint of the blob payload limit, which is capped by MaxMaxBlobPayloadBytes.
var maxEventFramingOverhead = protowire.SizeTag(coordinatedEventsFieldNumber) + protowire.SizeVarint(uint64(ocr3_1types.MaxMaxBlobPayloadBytes))

// maxEventBytes returns the largest encoded event (without timestamps) guaranteed to fit alone in a blob.
// It is conservative by at most a few bytes.
func (c *consensusQueuePlugin) maxEventBytes() int {
	return max(0, c.limits.MaxBlobPayloadBytes-maxEventFramingOverhead)
}

// maxEventsPerBlob is the most events a blob may hold. An observation has at most maxObservedEvents events in all of
// its blobs, so an honest blob has no more. Anything else only makes reading it costly.
const maxEventsPerBlob = maxObservedEvents

// splitEvents splits a blob into its marshaled events, without unmarshaling any of them. It rejects a blob with more
// than max events.
func splitEvents(b []byte, max int) ([][]byte, error) {
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

		if len(events) == max {
			return nil, fmt.Errorf("blob has more than %d events", max)
		}

		events = append(events, v)
		b = b[m:]
	}
	return events, nil
}

const (
	// minBlobLifetime is the least time a blob is kept past the events in it expiring. Peers reap expired blobs
	// asynchronously, and nodes that are behind, or that never received the blob, need time to fetch it.
	minBlobLifetime = 30 * time.Second
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
		rounds = max(rounds, uint64((lifetime+roundDur-1)/roundDur))
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

// broadcastBlob broadcasts a blob with all of the events, which must fit, and returns its marshaled handle.
func (c *consensusQueuePlugin) broadcastBlob(ctx context.Context, broadcaster ocr3_1types.BlobBroadcaster, events []*pb.CoordinatedEvent, seqNr uint64) ([]byte, error) {
	var maxDeadline time.Time
	for _, event := range events {
		if t := event.Deadline.AsTime(); t.After(maxDeadline) {
			maxDeadline = t
		}
	}

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
