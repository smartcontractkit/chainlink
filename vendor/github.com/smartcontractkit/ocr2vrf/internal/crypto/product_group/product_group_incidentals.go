package product_group

import (
	"crypto/cipher"
	"io"
	"reflect"

	"go.dedis.ch/fixbuf"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/util/random"
	"go.dedis.ch/kyber/v3/xof/blake2xb"
)

func (pg *ProductGroup) Write(w io.Writer, objs ...interface{}) error {
	return fixbuf.Write(w, objs)
}

func (pg *ProductGroup) Read(r io.Reader, objs ...interface{}) error {
	return fixbuf.Read(r, pg, objs)
}

var aScalar kyber.Scalar
var aPoint kyber.Point

var tScalar = reflect.TypeOf(&aScalar).Elem()
var tPoint = reflect.TypeOf(&aPoint).Elem()

func (pg *ProductGroup) New(t reflect.Type) interface{} {
	switch t {
	case tScalar:
		return pg.Scalar()
	case tPoint:
		return pg.Point()
	}
	return nil
}

func (pg *ProductGroup) XOF(seed []byte) kyber.XOF {
	return blake2xb.New(seed)
}

func (pg *ProductGroup) RandomStream() cipher.Stream {
	if pg.r != nil {
		return pg.r
	}
	return random.New()
}
