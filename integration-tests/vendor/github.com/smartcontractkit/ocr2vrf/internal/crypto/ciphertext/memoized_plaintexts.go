package ciphertext

import (
	"go.dedis.ch/kyber/v3"
)

var memoizedPlainTexts map[string][]kyber.Point

func rawPlainTexts(g kyber.Group) (rv []kyber.Point) {
	gen := g.Point().Base()
	rv = append(rv, g.Point().Null())
	for bi := 1; bi < 4; bi++ {
		rv = append(rv, g.Point().Add(rv[len(rv)-1], gen))
	}
	return rv
}

func memPlainTexts(g kyber.Group) []kyber.Point {
	rv, ok := memoizedPlainTexts[g.String()]
	if !ok {
		rv = rawPlainTexts(g)
	}
	return rv
}
