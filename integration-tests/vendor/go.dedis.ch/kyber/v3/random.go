package kyber

import (
	"crypto/cipher"
)

// Random is an interface that can be mixed in to local suite definitions.
type Random interface {
	// RandomStream returns a cipher.Stream that produces a
	// cryptographically random key stream. The stream must
	// tolerate being used in multiple goroutines.
	RandomStream() cipher.Stream
}
