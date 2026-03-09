package v2

import (
	"sync"
)

// ObservationBuffer holds events from Engine.Put for the plugin's Observation to read.
// Draft: in-memory buffer; full impl would integrate with OCR 3.1 KV store.
type ObservationBuffer struct {
	mu     sync.Mutex
	events []EnqueuedTriggerEvent
}

// Add appends an event to the buffer. Called by OCRQueue.Put.
func (b *ObservationBuffer) Add(event EnqueuedTriggerEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.events = append(b.events, event)
}

// TakeForObservation returns and clears buffered events for this round.
// Called by the plugin's Observation to produce its observation.
func (b *ObservationBuffer) TakeForObservation() []EnqueuedTriggerEvent {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := b.events
	b.events = nil
	return out
}
