// To fuzz-test this package:
//
//   $ go get -u github.com/dvyukov/go-fuzz/go-fuzz
//   $ go get -u github.com/dvyukov/go-fuzz/go-fuzz-build
//   $ cd `go env GOPATH`
//   $ go-fuzz-build go.dedis.ch/fixbuf
//   $ go-fuzz -workdir=workdir -bin fixbuf-fuzz.zip
//
// See also: https://medium.com/@dgryski/go-fuzz-github-com-arolek-ase-3c74d5a3150c
// The cd to $GOPATH is required as a workaround for:
// https://github.com/dvyukov/go-fuzz/issues/195

// +build gofuzz

package fixbuf

import (
	"bytes"
	"reflect"
)

type theFactory int
type t1 [32]byte
type t2 struct {
	X, Y t1
	Sl   []bool
	T3   t3
	T3s  [3]t3
}
type t3 struct {
	I int
	F float64
	B bool
}

var aT1 t1
var aT2 t2
var tT1 = reflect.TypeOf(&aT1).Elem()
var tT2 = reflect.TypeOf(&aT2).Elem()

func (theFactory) New(t reflect.Type) interface{} {
	switch t {
	case tT1:
		return new(t1)
	case tT2:
		return new(t2)
	}
	return nil
}

func Fuzz(data []byte) int {
	fac := theFactory(1)
	objs := make([]interface{}, 5)
	objs[0] = new(t1)
	objs[1] = new(t1)
	objs[2] = &t2{Sl: make([]bool, 3)}
	objs[3] = new(t1)
	objs[4] = &t2{Sl: make([]bool, 3)}
	if err := Read(bytes.NewReader(data), fac, objs...); err != nil {
		return 0
	}
	return 1
}
