package queue

import (
	"context"
	"fmt"
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

// broadcastAsBlob broadcasts a blob with as many coordinated events as possible, returning the blob handle and count of events included.
func (c *consensusQueuePlugin) broadcastAsBlob(ctx context.Context, broadcaster ocr3_1types.BlobBroadcastFetcher, events []*pb.CoordinatedEvent,
	seqNr uint64) (blobHandle []byte, n int, err error) {
	protoOpts := proto.MarshalOptions{Deterministic: true}
	for ; protoOpts.Size(&pb.CoordinatedEventList{Events: events}) > c.limits.MaxBlobPayloadBytes; events = events[:len(events)-1] {
		//TODO opt: binary search for largest fit instead....
	}
	n = len(events)
	var maxDeadline time.Time
	for _, event := range events {
		if t := event.Deadline.AsTime(); t.After(maxDeadline) {
			maxDeadline = t
		}
	}
	roundDur := c.onchainCfg.EstimatedRoundInterval //TODO OK? pad?
	//TODO pad extra time and/or rounds?
	expSeqNr := seqNr + uint64(time.Until(maxDeadline)/roundDur)

	data, err := protoOpts.Marshal(&pb.CoordinatedEventList{Events: events})
	if err != nil {
		return nil, -1, fmt.Errorf("failed to marshal blob: %w", err)
	}
	bh, err := broadcaster.BroadcastBlob(ctx, data, ocr3_1types.BlobExpirationHintSequenceNumber{SeqNr: expSeqNr})
	if err != nil {
		return nil, -1, fmt.Errorf("failed to broadcast blob: %w", err)
	}
	blobHandle, err = bh.MarshalBinary()
	return
}
