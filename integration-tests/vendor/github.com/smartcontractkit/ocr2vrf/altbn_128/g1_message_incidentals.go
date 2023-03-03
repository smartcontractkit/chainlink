package altbn_128

import (
	"crypto/cipher"
	"io"
	"reflect"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/util/random"
	"go.dedis.ch/kyber/v3/xof/blake2xb"
)

func (g *G1) XOF(seed []byte) kyber.XOF {
	return blake2xb.New(seed)
}

func (g *G1) RandomStream() cipher.Stream {
	if g.r != nil {
		return g.r
	}
	return random.New()
}

func (g *G1) Write(w io.Writer, objs ...interface{}) error { panic("not implemented") }
func (g *G1) Read(r io.Reader, objs ...interface{}) error  { panic("not implemented") }
func (g *G1) New(t reflect.Type) interface{}               { panic("not implemented") }
