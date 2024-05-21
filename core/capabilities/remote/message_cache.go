package remote

// MessageCache is a simple store for messages, grouped by event ID and peer ID.
// It is used to collect messages from multiple peers until they are ready for aggregation
// based on quantity and freshness.
type messageCache[EventID comparable, PeerID comparable] struct {
	events map[EventID]*eventState[PeerID]
}

type eventState[PeerID comparable] struct {
	peerMsgs          map[PeerID]*msgState
	creationTimestamp int64
	wasReady          bool
}

type msgState struct {
	timestamp int64
	payload   []byte
}

func NewMessageCache[EventID comparable, PeerID comparable]() *messageCache[EventID, PeerID] {
	return &messageCache[EventID, PeerID]{
		events: make(map[EventID]*eventState[PeerID]),
	}
}

// Insert or overwrite a message for <eventID>. Return creation timestamp of the event.
func (c *messageCache[EventID, PeerID]) Insert(eventID EventID, peerID PeerID, timestamp int64, payload []byte) int64 {
	if _, ok := c.events[eventID]; !ok {
		c.events[eventID] = &eventState[PeerID]{
			peerMsgs:          make(map[PeerID]*msgState),
			creationTimestamp: timestamp,
		}
	}
	c.events[eventID].peerMsgs[peerID] = &msgState{
		timestamp: timestamp,
		payload:   payload,
	}
	return c.events[eventID].creationTimestamp
}

// Return true if there are messages from at least <minCount> peers,
// received more recently than <minTimestamp>.
// Return all messages that satisfy the above condition.
// Ready() will return true at most once per event if <once> is true.
func (c *messageCache[EventID, PeerID]) Ready(eventID EventID, minCount uint32, minTimestamp int64, once bool) (bool, [][]byte) {
	ev, ok := c.events[eventID]
	if !ok {
		return false, nil
	}
	if ev.wasReady && once {
		return false, nil
	}
	if uint32(len(ev.peerMsgs)) < minCount {
		return false, nil
	}
	countAboveMinTimestamp := uint32(0)
	accPayloads := [][]byte{}
	for _, msg := range ev.peerMsgs {
		if msg.timestamp >= minTimestamp {
			countAboveMinTimestamp++
			accPayloads = append(accPayloads, msg.payload)
			if countAboveMinTimestamp >= minCount {
				ev.wasReady = true
				return true, accPayloads
			}
		}
	}
	return false, nil
}

func (c *messageCache[EventID, PeerID]) Delete(eventID EventID) {
	delete(c.events, eventID)
}

// Return the number of events deleted.
// Scans all keys, which might be slow for large caches.
func (c *messageCache[EventID, PeerID]) DeleteOlderThan(cutoffTimestamp int64) int {
	nDeleted := 0
	for id, event := range c.events {
		if event.creationTimestamp < cutoffTimestamp {
			c.Delete(id)
			nDeleted++
		}
	}
	return nDeleted
}
