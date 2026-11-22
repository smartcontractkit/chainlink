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
)

// blobCache is a short-lived cache for coordinating concurrent fetches of event blobs.
type blobCache struct {
	fetcher ocr3_1types.BlobFetcher

	mu    sync.RWMutex
	cache map[string]func() ([]byte, error) // blob handle -> blob
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

// fetchCoordinatedEventFromBlob fetches and extracts the idx-th event from the blob, without unmarshaling any other events.
func fetchCoordinatedEventFromBlob(ctx context.Context, fetcher *blobCache, handle []byte, idx int) (*pb.CoordinatedEvent, []byte, error) {
	bb, err := fetcher.fetchBlob(ctx, handle)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch blob: %w", err)
	}

	coordEvent, eb, err := extractEvent(bb, idx)
	if err != nil {
		return nil, nil, err
	}
	return coordEvent, eb, nil
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

// extractEvent extracts the idx-th event from the blob, without unmarshaling any other events.
func extractEvent(b []byte, idx int) (*pb.CoordinatedEvent, []byte, error) {
	i := 0
	for len(b) > 0 {
		num, typ, n := protowire.ConsumeTag(b)
		if n < 0 {
			return nil, nil, protowire.ParseError(n)
		}
		b = b[n:]
		if num == coordinatedEventsFieldNumber {
			if typ != protowire.BytesType {
				return nil, nil, fmt.Errorf("expected bytes type for field %d, got %d", coordinatedEventsFieldNumber, typ)
			}
			v, m := protowire.ConsumeBytes(b)
			if m < 0 {
				return nil, nil, protowire.ParseError(m)
			}
			if i == idx {
				var event pb.CoordinatedEvent
				if err := proto.Unmarshal(v, &event); err != nil {
					return nil, v, fmt.Errorf("failed to unmarshal event: %w", err)
				}
				return &event, v, nil
			}
			// advance
			i++
			b = b[m:]
			continue
		} else {
			n := protowire.ConsumeFieldValue(num, typ, b)
			if n < 0 {
				return nil, nil, protowire.ParseError(n)
			}
			b = b[n:]
		}
	}
	return nil, nil, fmt.Errorf("index %d out of range", idx)
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

// broadcastAsBlob broadcasts a blob with as many coordinated events as possible, returning the blob handle and count of events included.
func (c *consensusQueuePlugin) broadcastAsBlob(ctx context.Context, broadcaster ocr3_1types.BlobBroadcaster, events []*pb.CoordinatedEvent,
	seqNr uint64) (blobHandle []byte, n int, err error) {
	n = c.blobPrefix(events)
	if n == 0 {
		return nil, 0, errors.New("failed to package any events in blob")
	}
	blobHandle, err = c.broadcastBlob(ctx, broadcaster, events[:n], seqNr)
	if err != nil {
		return nil, -1, err
	}
	return blobHandle, n, nil
}
