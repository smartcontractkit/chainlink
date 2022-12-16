package dkg

import (
	"encoding/hex"
	"fmt"

	"github.com/smartcontractkit/ocr2vrf/dkg"
)

type dummyKeyConsumer struct{}

func (d dummyKeyConsumer) KeyInvalidated(keyID dkg.KeyID) {
	fmt.Println("KEY INVALIDATED:", hex.EncodeToString(keyID[:]))
}

func (d dummyKeyConsumer) NewKey(keyID dkg.KeyID, data *dkg.KeyData) {
	fmt.Println("NEW KEY FOR KEY ID:", hex.EncodeToString(keyID[:]), "KEY:", data)
}

var _ dkg.KeyConsumer = dummyKeyConsumer{}

func newDummyKeyConsumer() dummyKeyConsumer {
	return dummyKeyConsumer{}
}
