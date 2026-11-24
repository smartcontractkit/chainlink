// Package queuelimits holds the consensus-queue plugin's runtime sizing constants and the formula for its blob quotas.
// It is shared by the plugin, config validation and deployment tooling. Configuration validation itself lives in
// queuepb, next to the config it validates.
package queuelimits

import (
	"math"
	"time"

	"google.golang.org/protobuf/encoding/protowire"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
)

// Plugin sizing constants.
const (
	// MaxBlobHandlesPerObservation is the most blobs a node broadcasts in one Observation.
	MaxBlobHandlesPerObservation = 10
	// MaxObservedEvents is the most events a node reports in one Observation.
	MaxObservedEvents = 1000
	// MinBlobLifetime is the least time a blob is kept past the events in it expiring.
	MinBlobLifetime = 30 * time.Second
	// BlobReapingSlack is the extra time a blob may remain after it expires before it is reaped. libocr reaps
	// asynchronously, and asks for a margin of tens of seconds when sizing the per-oracle blob quotas.
	BlobReapingSlack = time.Minute
)

// The KV keys StateTransition writes every round.
const (
	ObservedEventsKey  = "observed_events::index"
	QueuedEventsKey    = "queued_events::index"
	ConsensusEventsKey = "consensus_events::index"
	BlobHandlesKey     = "blob_handles::index"

	// NumKVKeys is the number of KV keys StateTransition writes every round.
	NumKVKeys = 4
)

// Blob framing constants.
const (
	// CoordinatedEventsFieldNumber is the protobuf field number of CoordinatedEventList.Events.
	CoordinatedEventsFieldNumber = 1
)

// MaxEventFramingOverhead returns the most framing an event adds as an element of a blob: its field tag and the
// length prefix, whose varint is at most as long as the blob payload limit.
func MaxEventFramingOverhead(fieldNumber protowire.Number) int {
	return protowire.SizeTag(fieldNumber) + protowire.SizeVarint(uint64(ocr3_1types.MaxMaxBlobPayloadBytes))
}

// BlobQuota returns the per-oracle unexpired blob quota, in blobs and cumulative payload bytes, needed to hold a
// node's broadcasts for the lifetime of their events plus a margin for asynchronous reaping.
//
// A blob is kept until its event's deadline plus the grace period, and at least MinBlobLifetime. A node broadcasts at
// most MaxBlobHandlesPerObservation blobs per round. eventTimeout and grace are the workflow's trigger event timeout
// and the consensus-queue expiration grace; roundDur is the DON's DeltaRound.
//
// payloadBytesPerRound is the sustained blob payload per round the quota should cover. Pass the reference the
// deployment targets, not necessarily MaxBlobPayloadBytes: the config field is a uint32, and covering the worst case
// (MaxBlobHandlesPerObservation full blobs every round) would need far more than it can express.
//
// Both results saturate at math.MaxUint32, the widest value the config proto can carry.
func BlobQuota(roundDur, eventTimeout, grace time.Duration, payloadBytesPerRound int) (count, payloadBytes uint32) {
	if roundDur <= 0 {
		roundDur = time.Second
	}
	lifetime := max(0, eventTimeout) + max(0, grace) + MinBlobLifetime
	rounds := (lifetime + roundDur - 1) / roundDur // round up: every partial round is a full one
	rounds = max(rounds, 3)                        // blobs are kept for at least three rounds
	reapingRounds := (BlobReapingSlack + roundDur - 1) / roundDur
	// Both are non-negative: roundDur is positive (defaulted above), and each is a positive duration divided by it.
	totalRounds := uint64(rounds + reapingRounds) //nolint:gosec // G115: non-negative rounds, see above

	count = saturatingMul(uint64(MaxBlobHandlesPerObservation), totalRounds)
	payloadBytes = saturatingMul(uint64(max(0, payloadBytesPerRound)), totalRounds)
	return
}

// saturatingMul multiplies without overflowing, saturating at math.MaxUint32.
func saturatingMul(a, b uint64) uint32 {
	if a == 0 || b == 0 {
		return 0
	}
	if a > math.MaxUint32/b {
		return math.MaxUint32
	}
	return uint32(a * b) //nolint:gosec // G115 the guard above bounds the product by math.MaxUint32
}
