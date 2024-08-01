// To fuzz-test this package:
//
//   $ go get -u github.com/dvyukov/go-fuzz/go-fuzz
//   $ go get -u github.com/dvyukov/go-fuzz/go-fuzz-build
//   $ cd `go env GOPATH`
//   $ go-fuzz-build go.dedis.ch/protobuf
//   $ go-fuzz -workdir=workdir -bin protobuf-fuzz.zip

// +build gofuzz

package protobuf

import (
	"fmt"
	"reflect"
)

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

func Fuzz(data []byte) int {
	var it1, it2 t2
	var err error
	if err = Decode(data, &it1); err != nil {
		return 0
	}
	var buf []byte
	if buf, err = Encode(it1); err != nil {
		return 0
	}
	if err = Decode(buf, &it2); err != nil {
		return 0
	}
	if !reflect.DeepEqual(it1, it2) {
		panic(fmt.Sprintf("round trip not equal %#v %#v", it1, it2))
	}
	return 1
}
