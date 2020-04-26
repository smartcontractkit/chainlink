package bulletin_board

type BoardKey string
type BoardValue []byte

// Boards represents the communication system provided by bulletin_board
type Boards interface {
	// MakeOwnBoard returns publishing mechanism for current node. Synchronous.
	MakeOwnBoard(SecretKey) (OwnBulletinBoard, error)
	// MakeBoard returns the viewing mechanism for the node represented by key
	MakeBoard(PublicKey, TimeoutHandler) BulletinBoard
}

// BulletinBoard represents public information a node presents to others.
//
// It must deal with network failures (except timeouts) transparently to the
// interface user. For this reason, all responses are via callbacks.
//
// It must ignore unauthenticated messages.
type BulletinBoard interface {
	// Get passes the value associated with key, if available, or an error
	Get(BoardKey, BoardUpdateHandler, TimeoutHandler)
	// Subscribe listens for updates on a string / regexp, passing them to handler
	Subscribe(Subscription, BoardUpdateHandler, TimeoutHandler)
}

// OwnBulletinBoard represents the publishing mechanism for a BulletinBoard
type OwnBulletinBoard interface {
	// Publish makes the given key, value pair available to listening Boards
	//
	// It must take responsibility for asynchronously pushing the value out to
	// subscribed participants, so that this is "fire-and-forget" for the
	// interface user (except for timeouts.)
	Publish(BoardKey, BoardValue, TimeoutHandler)
}
