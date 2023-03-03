package types

import "context"

type DiscovererDatabase interface {
	// StoreAnnouncement has key-value-store semantics and stores a peerID (key) and an associated serialized
	//announcement (value).
	StoreAnnouncement(ctx context.Context, peerID string, ann []byte) error

	// ReadAnnouncements returns one serialized announcement (if available) for each of the peerIDs in the form of a map
	// keyed by each announcement's corresponding peer ID.
	ReadAnnouncements(ctx context.Context, peerIDs []string) (map[string][]byte, error)
}
