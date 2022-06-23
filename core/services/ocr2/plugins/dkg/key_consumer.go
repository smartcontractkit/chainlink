package dkg

import (
	"encoding/hex"
	"fmt"

	"github.com/smartcontractkit/ocr2vrf/dkg"
	dkgpkg "github.com/smartcontractkit/ocr2vrf/pkg/dkg"
	dkgcontract "github.com/smartcontractkit/ocr2vrf/pkg/dkg/contract"
)

type dummyKeyConsumer struct{}

func (d dummyKeyConsumer) KeyInvalidated(keyID dkgcontract.KeyID) {
	fmt.Println("KEY INVALIDATED:", hex.EncodeToString(keyID[:]))
}

func (d dummyKeyConsumer) NewKey(keyID dkgcontract.KeyID, data *dkgpkg.KeyData) {
	fmt.Println("NEW KEY FOR KEY ID:", hex.EncodeToString(keyID[:]), "KEY:", data)
}

var _ dkg.KeyConsumer = dummyKeyConsumer{}

func newDummyKeyConsumer() dummyKeyConsumer {
	return dummyKeyConsumer{}
}
