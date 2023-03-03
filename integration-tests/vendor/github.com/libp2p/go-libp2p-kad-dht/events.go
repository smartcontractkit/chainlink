package dht

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/google/uuid"

	"github.com/libp2p/go-libp2p-core/peer"
	kbucket "github.com/libp2p/go-libp2p-kbucket"
)

// KeyKadID contains the Kademlia key in string and binary form.
type KeyKadID struct {
	Key string
	Kad kbucket.ID
}

// NewKeyKadID creates a KeyKadID from a string Kademlia ID.
func NewKeyKadID(k string) *KeyKadID {
	return &KeyKadID{
		Key: k,
		Kad: kbucket.ConvertKey(k),
	}
}

// PeerKadID contains a libp2p Peer ID and a binary Kademlia ID.
type PeerKadID struct {
	Peer peer.ID
	Kad  kbucket.ID
}

// NewPeerKadID creates a PeerKadID from a libp2p Peer ID.
func NewPeerKadID(p peer.ID) *PeerKadID {
	return &PeerKadID{
		Peer: p,
		Kad:  kbucket.ConvertPeerID(p),
	}
}

// NewPeerKadIDSlice creates a slice of PeerKadID from the passed slice of libp2p Peer IDs.
func NewPeerKadIDSlice(p []peer.ID) []*PeerKadID {
	r := make([]*PeerKadID, len(p))
	for i := range p {
		r[i] = NewPeerKadID(p[i])
	}
	return r
}

// OptPeerKadID returns a pointer to a PeerKadID or nil if the passed Peer ID is it's default value.
func OptPeerKadID(p peer.ID) *PeerKadID {
	if p == "" {
		return nil
	}
	return NewPeerKadID(p)
}

// NewLookupEvent creates a LookupEvent automatically converting the node
// libp2p Peer ID to a PeerKadID and the string Kademlia key to a KeyKadID.
func NewLookupEvent(
	node peer.ID,
	id uuid.UUID,
	key string,
	request *LookupUpdateEvent,
	response *LookupUpdateEvent,
	terminate *LookupTerminateEvent,
) *LookupEvent {
	return &LookupEvent{
		Node:      NewPeerKadID(node),
		ID:        id,
		Key:       NewKeyKadID(key),
		Request:   request,
		Response:  response,
		Terminate: terminate,
	}
}

// LookupEvent is emitted for every notable event that happens during a DHT lookup.
// LookupEvent supports JSON marshalling because all of its fields do, recursively.
type LookupEvent struct {
	// Node is the ID of the node performing the lookup.
	Node *PeerKadID
	// ID is a unique identifier for the lookup instance.
	ID uuid.UUID
	// Key is the Kademlia key used as a lookup target.
	Key *KeyKadID
	// Request, if not nil, describes a state update event, associated with an outgoing query request.
	Request *LookupUpdateEvent
	// Response, if not nil, describes a state update event, associated with an outgoing query response.
	Response *LookupUpdateEvent
	// Terminate, if not nil, describe a termination event.
	Terminate *LookupTerminateEvent
}

// NewLookupUpdateEvent creates a new lookup update event, automatically converting the passed peer IDs to peer Kad IDs.
func NewLookupUpdateEvent(
	cause peer.ID,
	source peer.ID,
	heard []peer.ID,
	waiting []peer.ID,
	queried []peer.ID,
	unreachable []peer.ID,
) *LookupUpdateEvent {
	return &LookupUpdateEvent{
		Cause:       OptPeerKadID(cause),
		Source:      OptPeerKadID(source),
		Heard:       NewPeerKadIDSlice(heard),
		Waiting:     NewPeerKadIDSlice(waiting),
		Queried:     NewPeerKadIDSlice(queried),
		Unreachable: NewPeerKadIDSlice(unreachable),
	}
}

// LookupUpdateEvent describes a lookup state update event.
type LookupUpdateEvent struct {
	// Cause is the peer whose response (or lack of response) caused the update event.
	// If Cause is nil, this is the first update event in the lookup, caused by the seeding.
	Cause *PeerKadID
	// Source is the peer who informed us about the peer IDs in this update (below).
	Source *PeerKadID
	// Heard is a set of peers whose state in the lookup's peerset is being set to "heard".
	Heard []*PeerKadID
	// Waiting is a set of peers whose state in the lookup's peerset is being set to "waiting".
	Waiting []*PeerKadID
	// Queried is a set of peers whose state in the lookup's peerset is being set to "queried".
	Queried []*PeerKadID
	// Unreachable is a set of peers whose state in the lookup's peerset is being set to "unreachable".
	Unreachable []*PeerKadID
}

// LookupTerminateEvent describes a lookup termination event.
type LookupTerminateEvent struct {
	// Reason is the reason for lookup termination.
	Reason LookupTerminationReason
}

// NewLookupTerminateEvent creates a new lookup termination event with a given reason.
func NewLookupTerminateEvent(reason LookupTerminationReason) *LookupTerminateEvent {
	return &LookupTerminateEvent{Reason: reason}
}

// LookupTerminationReason captures reasons for terminating a lookup.
type LookupTerminationReason int

// MarshalJSON returns the JSON encoding of the passed lookup termination reason.
func (r LookupTerminationReason) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r LookupTerminationReason) String() string {
	switch r {
	case LookupStopped:
		return "stopped"
	case LookupCancelled:
		return "cancelled"
	case LookupStarvation:
		return "starvation"
	case LookupCompleted:
		return "completed"
	}
	panic("unreachable")
}

const (
	// LookupStopped indicates that the lookup was aborted by the user's stopFn.
	LookupStopped LookupTerminationReason = iota
	// LookupCancelled indicates that the lookup was aborted by the context.
	LookupCancelled
	// LookupStarvation indicates that the lookup terminated due to lack of unqueried peers.
	LookupStarvation
	// LookupCompleted indicates that the lookup terminated successfully, reaching the Kademlia end condition.
	LookupCompleted
)

type routingLookupKey struct{}

// TODO: lookupEventChannel copies the implementation of eventChanel.
// The two should be refactored to use a common event channel implementation.
// A common implementation needs to rethink the signature of RegisterForEvents,
// because returning a typed channel cannot be made polymorphic without creating
// additional "adapter" channels. This will be easier to handle when Go
// introduces generics.
type lookupEventChannel struct {
	mu  sync.Mutex
	ctx context.Context
	ch  chan<- *LookupEvent
}

// waitThenClose is spawned in a goroutine when the channel is registered. This
// safely cleans up the channel when the context has been canceled.
func (e *lookupEventChannel) waitThenClose() {
	<-e.ctx.Done()
	e.mu.Lock()
	close(e.ch)
	// 1. Signals that we're done.
	// 2. Frees memory (in case we end up hanging on to this for a while).
	e.ch = nil
	e.mu.Unlock()
}

// send sends an event on the event channel, aborting if either the passed or
// the internal context expire.
func (e *lookupEventChannel) send(ctx context.Context, ev *LookupEvent) {
	e.mu.Lock()
	// Closed.
	if e.ch == nil {
		e.mu.Unlock()
		return
	}
	// in case the passed context is unrelated, wait on both.
	select {
	case e.ch <- ev:
	case <-e.ctx.Done():
	case <-ctx.Done():
	}
	e.mu.Unlock()
}

// RegisterForLookupEvents registers a lookup event channel with the given context.
// The returned context can be passed to DHT queries to receive lookup events on
// the returned channels.
//
// The passed context MUST be canceled when the caller is no longer interested
// in query events.
func RegisterForLookupEvents(ctx context.Context) (context.Context, <-chan *LookupEvent) {
	ch := make(chan *LookupEvent, LookupEventBufferSize)
	ech := &lookupEventChannel{ch: ch, ctx: ctx}
	go ech.waitThenClose()
	return context.WithValue(ctx, routingLookupKey{}, ech), ch
}

// LookupEventBufferSize is the number of events to buffer.
var LookupEventBufferSize = 16

// PublishLookupEvent publishes a query event to the query event channel
// associated with the given context, if any.
func PublishLookupEvent(ctx context.Context, ev *LookupEvent) {
	ich := ctx.Value(routingLookupKey{})
	if ich == nil {
		return
	}

	// We *want* to panic here.
	ech := ich.(*lookupEventChannel)
	ech.send(ctx, ev)
}
