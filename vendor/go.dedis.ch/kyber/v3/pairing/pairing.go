package pairing

import "go.dedis.ch/kyber/v3"

// Suite interface represents a triplet of elliptic curve groups (G₁, G₂
// and GT) such that there exists a function e(g₁ˣ,g₂ʸ)=gTˣʸ (where gₓ is a
// generator of the respective group) which is called a pairing.
type Suite interface {
	G1() kyber.Group
	G2() kyber.Group
	GT() kyber.Group
	Pair(p1, p2 kyber.Point) kyber.Point
	kyber.Encoding
	kyber.HashFactory
	kyber.XOFFactory
	kyber.Random
}
