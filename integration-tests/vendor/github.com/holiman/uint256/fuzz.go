// uint256: Fixed size 256-bit math library
// Copyright 2020 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

// +build gofuzz

package uint256

import (
	"fmt"
	"math/big"
	"reflect"
	"runtime"
)

const (
	opUdivrem = 0
	opMul     = 1
	opLsh     = 2
	opAdd     = 4
	opSub     = 5
)

type opFunc func(*Int, *Int, *Int) *Int
type bigFunc func(*big.Int, *big.Int, *big.Int) *big.Int

func crash(op opFunc, x, y Int, msg string) {
	fn := runtime.FuncForPC(reflect.ValueOf(op).Pointer())
	fnName := fn.Name()
	fnFile, fnLine := fn.FileLine(fn.Entry())
	panic(fmt.Sprintf("%s\nfor %s (%s:%d)\nx: %x\ny: %x", msg, fnName, fnFile, fnLine, &x, &y))
}

func checkOp(op opFunc, bigOp bigFunc, x, y Int) {
	origX := x
	origY := y

	var result Int
	ret := op(&result, &x, &y)
	if ret != &result {
		crash(op, x, y, "returned not the pointer receiver")
	}
	if x != origX {
		crash(op, x, y, "first argument modified")
	}
	if y != origY {
		crash(op, x, y, "second argument modified")
	}

	expected, _ := FromBig(bigOp(new(big.Int), x.ToBig(), y.ToBig()))
	if result != *expected {
		crash(op, x, y, "unexpected result")
	}

	// Test again when the receiver is not zero.
	var garbage Int
	garbage.Xor(&x, &y)
	ret = op(&garbage, &x, &y)
	if ret != &garbage {
		crash(op, x, y, "returned not the pointer receiver")
	}
	if garbage != *expected {
		crash(op, x, y, "unexpected result")
	}
	if x != origX {
		crash(op, x, y, "first argument modified")
	}
	if y != origY {
		crash(op, x, y, "second argument modified")
	}

	// Test again with the receiver aliasing arguments.
	ret = op(&x, &x, &y)
	if ret != &x {
		crash(op, x, y, "returned not the pointer receiver")
	}
	if x != *expected {
		crash(op, x, y, "unexpected result")
	}

	ret = op(&y, &origX, &y)
	if ret != &y {
		crash(op, x, y, "returned not the pointer receiver")
	}
	if y != *expected {
		crash(op, x, y, "unexpected result")
	}
}

func Fuzz(data []byte) int {
	if len(data) != 65 {
		return 0
	}

	op := data[0]

	var x, y Int
	x.SetBytes(data[1:33])
	y.SetBytes(data[33:])

	switch op {
	case opUdivrem:
		if y.IsZero() {
			return 0
		}
		checkOp((*Int).Div, (*big.Int).Div, x, y)
		checkOp((*Int).Mod, (*big.Int).Mod, x, y)

	case opMul:
		checkOp((*Int).Mul, (*big.Int).Mul, x, y)

	case opLsh:
		lsh := func(z, x, y *Int) *Int {
			return z.Lsh(x, uint(y[0]))
		}
		bigLsh := func(z, x, y *big.Int) *big.Int {
			n := uint(y.Uint64())
			if n > 256 {
				n = 256
			}
			return z.Lsh(x, n)
		}
		checkOp(lsh, bigLsh, x, y)

	case opAdd:
		checkOp((*Int).Add, (*big.Int).Add, x, y)

	case opSub:
		checkOp((*Int).Sub, (*big.Int).Sub, x, y)
	}

	return 0
}
