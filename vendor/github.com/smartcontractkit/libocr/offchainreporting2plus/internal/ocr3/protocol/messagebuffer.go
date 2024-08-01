package protocol

import "github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3/protocol/ringbuffer"

// We have this wrapper to deal with what appears to be a bug in the Go compiler
// that prevents us from using ringbuffer.RingBuffer in the outcome generation
// protocol:
// offchainreporting2plus/internal/ocr3/protocol/outcome_generation.go:241:21: internal compiler error: (*ringbuffer.RingBuffer[go.shape.interface { github.com/smartcontractkit/offchain-reporting/lib/offchainreporting2plus/internal/ocr3/protocol.epoch() uint64; github.com/smartcontractkit/offchain-reporting/lib/offchainreporting2plus/internal/ocr3/protocol.processOutcomeGeneration(*github.com/smartcontractkit/offchain-reporting/lib/offchainreporting2plus/internal/ocr3/protocol.outcomeGenerationState[go.shape.struct {}], github.com/smartcontractkit/offchain-reporting/lib/commontypes.OracleID) }]).Peek(buffer, (*[9]uintptr)(.dict[3])) (type go.shape.interface { github.com/smartcontractkit/offchain-reporting/lib/offchainreporting2plus/internal/ocr3/protocol.epoch() uint64; github.com/smartcontractkit/offchain-reporting/lib/offchainreporting2plus/internal/ocr3/protocol.processOutcomeGeneration(*github.com/smartcontractkit/offchain-reporting/lib/offchainreporting2plus/internal/ocr3/protocol.outcomeGenerationState[go.shape.struct {}], github.com/smartcontractkit/offchain-reporting/lib/commontypes.OracleID) }) is not shape-identical to MessageToOutcomeGeneration[go.shape.struct {}]
// Consider removing it in a future release.
type MessageBuffer[RI any] ringbuffer.RingBuffer[MessageToOutcomeGeneration[RI]]

func NewMessageBuffer[RI any](cap int) *MessageBuffer[RI] {
	return (*MessageBuffer[RI])(ringbuffer.NewRingBuffer[MessageToOutcomeGeneration[RI]](cap))
}

func (rb *MessageBuffer[RI]) Length() int {
	return (*ringbuffer.RingBuffer[MessageToOutcomeGeneration[RI]])(rb).Length()
}

func (rb *MessageBuffer[RI]) Peek() MessageToOutcomeGeneration[RI] {
	return (*ringbuffer.RingBuffer[MessageToOutcomeGeneration[RI]])(rb).Peek()
}

func (rb *MessageBuffer[RI]) Pop() MessageToOutcomeGeneration[RI] {
	return (*ringbuffer.RingBuffer[MessageToOutcomeGeneration[RI]])(rb).Pop()
}

func (rb *MessageBuffer[RI]) Push(msg MessageToOutcomeGeneration[RI]) {
	(*ringbuffer.RingBuffer[MessageToOutcomeGeneration[RI]])(rb).Push(msg)
}
