package request

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"sync/atomic"

	"google.golang.org/protobuf/proto"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
)

type ComparatorState struct {
	Message    *types.MessageBody
	Metadata   commoncap.ResponseMetadata
	Hash       [32]byte
	ResponseID string // The hex-encoded hash string
}

// payloadComparator handles logging and comparing divergent payloads.
type payloadComparator struct {
	lggr       logger.Logger
	firstState atomic.Pointer[ComparatorState]
}

// NewPayloadComparator creates and initializes a new comparator.
func NewPayloadComparator(lggr logger.Logger) *payloadComparator {
	return &payloadComparator{
		lggr: lggr,
	}
}

// Compare accepts an incoming message, stores the first unique instance, and logs
// differences against subsequent unique messages.
func (pc *payloadComparator) Compare(incomingMsg *types.MessageBody) {
	incomingHash, md, err := GetMessageHashAndMetadata(pc.lggr, incomingMsg)
	if err != nil {
		pc.lggr.Errorw("failed to get message data and hash", "error", err)
		return
	}

	if pc.firstState.Load() == nil {
		newState := &ComparatorState{
			Message:    incomingMsg,
			Metadata:   md,
			Hash:       incomingHash,
			ResponseID: hex.EncodeToString(incomingHash[:]),
		}

		if !pc.firstState.CompareAndSwap(nil, newState) {
			pc.lggr.Error("failed to load first state in comparator")
			return
		} else {
			pc.lggr.Debugw("Comparator stored first unique response payload.",
				"responseID", newState.ResponseID,
				"metadata", newState.Metadata,
			)
		}
	}

	existingState := pc.firstState.Load()
	if bytes.Equal(existingState.Hash[:], incomingHash[:]) {
		return
	}

	pc.lggr.Warnw("Divergent response hash detected. Running detailed payload comparison.",
		"firstResponseID", existingState.ResponseID,
		"firstMetadata", existingState.Metadata,
		"incomingResponseID", hex.EncodeToString(incomingHash[:]),
		"incomingMetadata", md,
	)

	incomingResp, err := pb.UnmarshalCapabilityResponse(incomingMsg.Payload)
	if err != nil {
		pc.lggr.Warnw("failed to unmarshal incoming message", err)
		return
	}

	expectedResp, err := pb.UnmarshalCapabilityResponse(existingState.Message.Payload)
	if err != nil {
		pc.lggr.Warnw("failed to unmarshal first message", err)
		return
	}

	logAnyPayloadByteDiffs(pc.lggr, incomingResp, expectedResp, incomingHash, existingState.Hash)
}

func logAnyPayloadByteDiffs(lggr logger.Logger, respA, respB commoncap.CapabilityResponse, hashA, hashB [32]byte) {
	if (respA.Payload == nil) != (respB.Payload == nil) {
		lggr.Errorw("Payload divergence: One response has a Payload and the other doesn't",
			"hashA_has_payload", respA.Payload != nil,
			"hashB_has_payload", respB.Payload != nil,
		)
		return
	}
	if respA.Payload == nil && respB.Payload == nil {
		return
	}

	if proto.Equal(respA.Payload, respB.Payload) {
		return
	}

	valA := respA.Payload.GetValue()
	valB := respB.Payload.GetValue()

	lenA := len(valA)
	lenB := len(valB)
	minLen := lenA
	if lenB < minLen {
		minLen = lenB
	}

	var logBuffer bytes.Buffer

	// Header for the report
	logBuffer.WriteString(
		fmt.Sprintf(
			"\n--- Payload Byte Difference Report (Type: %s vs %s) ---\n",
			respA.Payload.GetTypeUrl(), respB.Payload.GetTypeUrl(),
		),
	)
	logBuffer.WriteString(fmt.Sprintf("Raw Byte Lengths: HashA=%d, HashB=%d\n", lenA, lenB))
	diffsFound := 0

	// Compare byte by byte up to the minimum length
	for i := 0; i < minLen; i++ {
		if valA[i] != valB[i] {
			// Log the index and the hexadecimal value of the differing bytes
			logBuffer.WriteString(fmt.Sprintf("Index %d: HashA=0x%x, HashB=0x%x\n", i, valA[i], valB[i]))
			diffsFound++
		}
	}

	// Report differences in length, if any
	if lenA != lenB {
		diffLen := map[bool]int{true: lenA - lenB, false: lenB - lenA}[lenA != minLen]
		source := map[bool]string{true: "A", false: "B"}[lenA > lenB]
		logBuffer.WriteString(fmt.Sprintf("Payloads differ in length. Remaining bytes in Hash%s: %d\n", source, diffLen))
	}

	// Footer
	logBuffer.WriteString(fmt.Sprintf("Total Byte Differences Found: %d\n", diffsFound))

	lggr.Errorw("Divergence detected in *anypb.Any Payload content.",
		"hashA", hex.EncodeToString(hashA[:]),
		"hashB", hex.EncodeToString(hashB[:]),
		"byte_diff_report", logBuffer.String(),
	)
}
