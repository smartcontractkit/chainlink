package dkg

import (
	"fmt"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"
	"github.com/smartcontractkit/ocr2vrf/internal/dkg/contract"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"
)

type KeyConsumer interface {
	KeyInvalidated(contract.KeyID)

	NewKey(contract.KeyID, *KeyData)
}

type KeyData struct {
	PublicKey   kyber.Point
	Shares      []share.PubShare
	SecretShare *SecretShare

	T       player_idx.Int
	Present bool
}

func (kd KeyData) Clone() *KeyData {
	shares := make([]share.PubShare, len(kd.Shares))
	for i := 0; i < len(kd.Shares); i++ {
		if kd.Shares[i].V == nil {
			panic(fmt.Errorf("%dth public share is nil", i))
		}
		s := share.PubShare{kd.Shares[i].I, kd.Shares[i].V.Clone()}
		shares[i] = s
	}
	if kd.PublicKey == nil {
		panic("nil public key")
	}
	if kd.SecretShare == nil {
		panic("nil secret share")
	}
	return &KeyData{
		kd.PublicKey.Clone(),
		shares, kd.SecretShare.Clone(),
		kd.T,
		kd.Present,
	}
}
