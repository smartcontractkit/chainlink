package queue

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
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
	//TODO optimize using protowire to avoid unmarshalling the entire blob
	var coordEvents *pb.CoordinatedEventList
	if err := proto.Unmarshal(b, coordEvents); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal blob: %w", err)
	}
	if idx >= len(coordEvents.Events) {
		return nil, nil, fmt.Errorf("blob index %d is out of range for blob with %d events", idx, len(coordEvents.Events))
	}
	coordEvent := coordEvents.Events[idx]
	return coordEvent, nil, nil
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
