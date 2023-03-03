package vrf

import (
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing"
)

func validateSignature(p pairing.Suite, msg, pk, sig kyber.Point) bool {

	return p.Pair(msg, pk).Equal(p.Pair(sig, p.G2().Point().Base()))
}
