package altbn_128

import (
	"bytes"
	"crypto/cipher"
	"fmt"
	"io"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"

	"go.dedis.ch/kyber/v3"
)

type gTPoint struct{ *bn256.GT }

var _ kyber.Point = (*gTPoint)(nil)

func newGTPoint() *gTPoint { return &gTPoint{new(bn256.GT)} }

func (p *gTPoint) Equal(s2 kyber.Point) bool {
	gtS2, ok := s2.(*gTPoint)
	return ok && bytes.Equal(p.Marshal(), gtS2.Marshal())
}

func (p *gTPoint) String() string { return fmt.Sprintf("&gTPoint{%s}", p.GT.String()) }

func (p *gTPoint) MarshalBinary() (data []byte, err error)        { panic("not implemented") }
func (p *gTPoint) UnmarshalBinary(data []byte) error              { panic("not implemented") }
func (p *gTPoint) MarshalSize() int                               { panic("not implemented") }
func (p *gTPoint) MarshalTo(w io.Writer) (int, error)             { panic("not implemented") }
func (p *gTPoint) UnmarshalFrom(r io.Reader) (int, error)         { panic("not implemented") }
func (p *gTPoint) Null() kyber.Point                              { panic("not implemented") }
func (p *gTPoint) Base() kyber.Point                              { panic("not implemented") }
func (p *gTPoint) Pick(rand cipher.Stream) kyber.Point            { panic("not implemented") }
func (p *gTPoint) Set(p2 kyber.Point) kyber.Point                 { panic("not implemented") }
func (p *gTPoint) Clone() kyber.Point                             { panic("not implemented") }
func (p *gTPoint) EmbedLen() int                                  { panic("not implemented") }
func (p *gTPoint) Embed(data []byte, r cipher.Stream) kyber.Point { panic("not implemented") }
func (p *gTPoint) Data() ([]byte, error)                          { panic("not implemented") }
func (p *gTPoint) Add(a kyber.Point, b kyber.Point) kyber.Point   { panic("not implemented") }
func (p *gTPoint) Sub(a kyber.Point, b kyber.Point) kyber.Point   { panic("not implemented") }
func (p *gTPoint) Neg(a kyber.Point) kyber.Point                  { panic("not implemented") }
func (p *gTPoint) Mul(s kyber.Scalar, p2 kyber.Point) kyber.Point { panic("not implemented") }
