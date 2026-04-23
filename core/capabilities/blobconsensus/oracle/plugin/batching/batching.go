package batching

import (
	"context"

	"google.golang.org/protobuf/proto"
)

type metrics interface {
	IncBatchCapacityExceeded(ctx context.Context, step string)
	IncBatchRequestsTotal(ctx context.Context, step string)
	RecordObservationBatchSize(ctx context.Context, size float64)
	RecordOutcomeBatchSize(ctx context.Context, size float64)
}

// varintSize calculates the size of a varint encoding
func varintSize(x uint64) int {
	switch {
	case x < 1<<7:
		return 1
	case x < 1<<14:
		return 2
	case x < 1<<21:
		return 3
	case x < 1<<28:
		return 4
	case x < 1<<35:
		return 5
	case x < 1<<42:
		return 6
	case x < 1<<49:
		return 7
	case x < 1<<56:
		return 8
	case x < 1<<63:
		return 9
	default:
		return 10
	}
}

// calculateMessageSize calculates the marshalled size of any proto message
func calculateMessageSize(message proto.Message) int {
	// Use proto.Size which gives us the exact marshalled size
	return proto.Size(message)
}

func batchHasCapacityForMessageOnSlice(cachedSize int, message proto.Message, maxSizeBytes int) (bool, int) {
	numBytes := proto.Size(message)
	return batchHasCapacityForSliceBytes(cachedSize, numBytes, maxSizeBytes)
}

func batchHasCapacityForStringOnSlice(cachedSize int, message string, maxSizeBytes int) (bool, int) {
	numBytes := len(message)
	return batchHasCapacityForSliceBytes(cachedSize, numBytes, maxSizeBytes)
}

func batchHasCapacityForSliceBytes(cachedSize int, newMessageSize int, maxSizeBytes int) (bool, int) {
	// Add protobuf field overhead: tag (field number + wire type) + length prefix
	// For repeated fields in protobuf, each element gets:
	// - Tag: field number (1 for the repeated field) << 3 | wire type (2 for length-delimited)
	// - Length: varint encoding of the message size
	tagSize := varintSize(uint64(1<<3 | 2))
	if newMessageSize < 0 {
		return false, cachedSize
	}
	lengthSize := varintSize(uint64(newMessageSize))

	totalSizeWithNewMessage := cachedSize + tagSize + lengthSize + newMessageSize

	// Check against config
	if totalSizeWithNewMessage > maxSizeBytes {
		// Stop adding more messages
		return false, cachedSize
	}

	return true, totalSizeWithNewMessage
}
