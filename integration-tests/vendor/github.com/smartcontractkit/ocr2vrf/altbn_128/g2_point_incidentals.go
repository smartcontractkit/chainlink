package altbn_128

import (
	"crypto/cipher"

	"go.dedis.ch/kyber/v3/util/random"
)

func (g *G2) RandomStream() cipher.Stream {
	if g.r != nil {
		return g.r
	}
	return random.New()
}
