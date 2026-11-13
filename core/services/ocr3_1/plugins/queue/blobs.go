package queue

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
)

func fetchFromBlob(ctx context.Context, fetcher ocr3_1types.BlobFetcher, handle []byte, idx int) (*pb.CoordinatedEvent, []byte, error) {
	var blobHandle ocr3_1types.BlobHandle
	if err := blobHandle.UnmarshalBinary(handle); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal blob handle: %w", err)
	}
	//TODO cache to fetch just once
	b, err := fetcher.FetchBlob(ctx, blobHandle)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch blob: %w", err)
	}

	coordEvent, err := extractEvent(b, idx)
	if err != nil {
		return nil, nil, err
	}
	return coordEvent, b, nil
}

// extractEvent extracts the idx-th event from the blob, without unmarshaling any other events.
func extractEvent(b []byte, idx int) (*pb.CoordinatedEvent, error) {
	//TODO can we lock this in with code?
	const fieldNumber = 1 // pb.CoordinatedEventList.Events
	i := 0
	for len(b) > 0 {
		num, typ, n := protowire.ConsumeTag(b)
		if n < 0 {
			return nil, protowire.ParseError(n)
		}
		b = b[n:]
		if num == fieldNumber {
			if typ != protowire.BytesType {
				return nil, fmt.Errorf("expected bytes type for field %d, got %d", fieldNumber, typ)
			}
			v, m := protowire.ConsumeBytes(b)
			if m < 0 {
				return nil, protowire.ParseError(m)
			}
			if i == idx {
				var event pb.CoordinatedEvent
				if err := proto.Unmarshal(v, &event); err != nil {
					return nil, fmt.Errorf("failed to unmarshal event: %w", err)
				}
				return &event, nil
			}
			// advance
			i++
			b = b[m:]
			continue
		}
	}
	return nil, fmt.Errorf("index %d out of range", idx)
}

// broadcastAsBlob broadcasts a list of coordinated events and returns a blob handle.
func broadcastAsBlob(ctx context.Context, broadcaster ocr3_1types.BlobBroadcastFetcher, events []*pb.CoordinatedEvent, expSeq uint64) ([]byte, error) {
	// TODO return a slice of hashes too; encode list with protowire to hash them along the way - w/ ZERO timestamps
	data, err := proto.MarshalOptions{Deterministic: true}.Marshal(&pb.CoordinatedEventList{Events: events})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal blob: %w", err)
	}
	blobHandle, err := broadcaster.BroadcastBlob(ctx, data, ocr3_1types.BlobExpirationHintSequenceNumber{SeqNr: expSeq})
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast blob: %w", err)
	}
	return blobHandle.MarshalBinary()
}
