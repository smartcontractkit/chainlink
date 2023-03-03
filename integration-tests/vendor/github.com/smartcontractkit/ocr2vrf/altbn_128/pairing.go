package altbn_128

import (
	"crypto/cipher"
	"hash"
	"io"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing"
	"go.dedis.ch/kyber/v3/util/random"
)

type PairingSuite struct{}

var _ pairing.Suite = (*PairingSuite)(nil)

func (p *PairingSuite) G1() kyber.Group {
	return newG1()
}

func (p *PairingSuite) G2() kyber.Group {
	return &G2{}
}

func (p *PairingSuite) GT() kyber.Group {
	return &GT{}
}

func (p *PairingSuite) Pair(p1 kyber.Point, p2 kyber.Point) kyber.Point {
	pg1, pg2 := p1.(*g1Point), p2.(*g2Point)
	return &gTPoint{bn256.Pair(pg1.G1.(*bn256.G1), pg2.G2.(*bn256.G2))}
}

func (p *PairingSuite) Write(w io.Writer, objs ...interface{}) error {
	panic("not implemented")
}

func (p *PairingSuite) Read(r io.Reader, objs ...interface{}) error {
	panic("not implemented")
}

func (p *PairingSuite) Hash() hash.Hash {
	panic("not implemented")
}

func (p *PairingSuite) XOF(seed []byte) kyber.XOF {
	panic("not implemented")
}

func (p *PairingSuite) RandomStream() cipher.Stream {
	return random.New()
}
