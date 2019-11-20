package bulletin_board

type BoardKey string
type BoardValue []byte
type BoardMap map[BoardKey]BoardValue

// Boards represents the communication system provided by bulletin_board
type Boards interface {
	// MakeOwnBoard returns the publishing mechanism for the current node.
	MakeOwnBoard(key SecretKey) (OwnBulletinBoard, error)
	// MakeBoard returns the viewing mechanism for the node represented by key
	MakeBoard(key PublicKey) (BulletinBoard, error)
}

// BulletinBoard represents the public information a node presents to the
// others.
type BulletinBoard interface {
	// Get returns the value associated with key, if available, or an error
	Get(key BoardKey) BoardValue
	// Subscribe listens for updates on a string / regexp, passing them to handler
	Subscribe(key Subscription, handler SubscriptionHandler)
}

// OwnBulletinBoard represents the publishing mechanism for a BulletinBoard
type OwnBulletinBoard interface {
	// Publish makes the given key, value pair available to listening Boards
	Publish(key BoardKey, value BoardValue)
}
