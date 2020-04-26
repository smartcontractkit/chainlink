package whiteboard

type BoardKey interface{}

type BoardQuery interface {
	PublicKey() PublicKey
	Match(BoardKey) bool
}

type BoardValue interface{}

type BoardError interface {
	Error() string
}

// PublicKey is node public identity against which signatures are verified.
type PublicKey interface{}

type Signature interface{}

// SecretKey is node secret identity, used to sign all messages.
type SecretKey interface{}

// GetHandler receives any values or errors related to query
type GetHandler func(BoardQuery, []BoardValue, []BoardError)

type Board interface {
	// Publish pushes a BoardValue note under (SecretKey.PublicKey(), BoardKey)
	Publish(SecretKey, BoardKey, BoardValue)
	// Get sends GetHandler matches to (PublicKey, BoardQuery)
	Get(PublicKey, BoardQuery, GetHandler)
}
