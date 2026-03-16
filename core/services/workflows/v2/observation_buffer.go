package v2

import (
	"sync"
)

// LamportCounter is a thread-safe Lamport clock used for ordering events across nodes.
// Written during StateTransition (when receiving observations); read and incremented during Put.
type LamportCounter struct {
	mu      sync.Mutex
	counter uint64
}

// IncrementAndGet increments the counter and returns the new value.
// Called by ObservationBuffer.Add when enqueuing an event.
func (c *LamportCounter) IncrementAndGet() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counter++
	return c.counter
}

// UpdateFromReceived updates the counter from received observations.
// Called by the plugin during StateTransition: counter = max(counter, maxReceived) + 1.
func (c *LamportCounter) UpdateFromReceived(maxReceived uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if maxReceived >= c.counter {
		c.counter = maxReceived + 1
	}
}

// BufferedEvent exposes an event's ID and Lamport timestamp for observation.
type BufferedEvent interface {
	ID() string
	Lamport() uint64
}

// bufferedEvent implements BufferedEvent.
type bufferedEvent struct {
	event   EnqueuedTriggerEvent
	lamport uint64
}

func (b bufferedEvent) ID() string {
	return b.event.Event().Event.ID
}

func (b bufferedEvent) Lamport() uint64 {
	return b.lamport
}

// ObservationBuffer holds events from Engine.Put for the plugin's Observation to read.
// Draft: in-memory buffer; full impl would integrate with OCR 3.1 KV store.
// Holds a reference to LamportCounter for assigning timestamps on Add.
type ObservationBuffer struct {
	mu      sync.Mutex
	events  []BufferedEvent
	lamport *LamportCounter
}

// NewObservationBuffer creates a buffer with the given Lamport counter.
// Pass nil for lamport to omit Lamport assignment (events will have Lamport 0).
func NewObservationBuffer(lamport *LamportCounter) *ObservationBuffer {
	return &ObservationBuffer{lamport: lamport}
}

// Add appends an event to the buffer. Called by OCRQueue.Put.
// If lamport is set, increments it and assigns the value to the event.
func (b *ObservationBuffer) Add(event EnqueuedTriggerEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()
	var lamport uint64
	if b.lamport != nil {
		lamport = b.lamport.IncrementAndGet()
	}
	b.events = append(b.events, bufferedEvent{event: event, lamport: lamport})
}

// TakeForObservation returns and clears buffered events for this round.
// Called by the plugin's Observation to produce its observation.
func (b *ObservationBuffer) TakeForObservation() []BufferedEvent {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := b.events
	b.events = nil
	return out
}

// UpdateLamportFromReceived updates the Lamport counter from received observations.
// Called by the plugin during StateTransition. No-op if lamport is nil.
func (b *ObservationBuffer) UpdateLamportFromReceived(maxReceived uint64) {
	if b.lamport != nil {
		b.lamport.UpdateFromReceived(maxReceived)
	}
}
