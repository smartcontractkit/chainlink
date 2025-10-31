package request

import (
	"bytes"
	"encoding/hex"
	"sync/atomic"

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
		}

		pc.lggr.Debugw("Comparator stored first unique response payload.",
			"responseID", newState.ResponseID,
			"metadata", newState.Metadata,
		)
	}

	existingState := pc.firstState.Load()
	if bytes.Equal(existingState.Hash[:], incomingHash[:]) {
		return
	}

	incomingResp, err := pb.UnmarshalCapabilityResponse(incomingMsg.Payload)
	if err != nil {
		pc.lggr.Warnf("failed to unmarshal incoming message: %s", err)
		return
	}

	expectedResp, err := pb.UnmarshalCapabilityResponse(existingState.Message.Payload)
	if err != nil {
		pc.lggr.Warnf("failed to unmarshal first message: %s", err)
		return
	}

	pc.lggr.Warnw("Divergent response hash detected",
		"expectedResponseID", existingState.ResponseID,
		"expectedMetadata", existingState.Metadata,
		"expectedPayload", existingState.Message.Payload,
		"expectedResp", expectedResp,
		"incomingResponseID", hex.EncodeToString(incomingHash[:]),
		"incomingMetadata", md,
		"incomingPayload", incomingMsg.Payload,
		"incomingResp", incomingResp,
	)
}
