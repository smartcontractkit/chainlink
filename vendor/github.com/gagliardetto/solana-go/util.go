// Copyright 2020 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package solana

import (
	"math/big"
)

var _10b = big.NewInt(10)
var decimalsBigInt = []*big.Int{
	new(big.Int).Exp(_10b, big.NewInt(1), nil),
	new(big.Int).Exp(_10b, big.NewInt(2), nil),
	new(big.Int).Exp(_10b, big.NewInt(3), nil),
	new(big.Int).Exp(_10b, big.NewInt(4), nil),
	new(big.Int).Exp(_10b, big.NewInt(5), nil),
	new(big.Int).Exp(_10b, big.NewInt(6), nil),
	new(big.Int).Exp(_10b, big.NewInt(7), nil),
	new(big.Int).Exp(_10b, big.NewInt(8), nil),
	new(big.Int).Exp(_10b, big.NewInt(9), nil),
	new(big.Int).Exp(_10b, big.NewInt(10), nil),
	new(big.Int).Exp(_10b, big.NewInt(11), nil),
	new(big.Int).Exp(_10b, big.NewInt(12), nil),
	new(big.Int).Exp(_10b, big.NewInt(13), nil),
	new(big.Int).Exp(_10b, big.NewInt(14), nil),
	new(big.Int).Exp(_10b, big.NewInt(15), nil),
	new(big.Int).Exp(_10b, big.NewInt(16), nil),
	new(big.Int).Exp(_10b, big.NewInt(17), nil),
	new(big.Int).Exp(_10b, big.NewInt(18), nil),
}

func DecimalsInBigInt(decimal uint32) *big.Int {
	if decimal == 0 {
		return big.NewInt(1)
	}
	var decimalsBig *big.Int
	if decimal <= uint32(len(decimalsBigInt)) {
		decimalsBig = decimalsBigInt[decimal-1]
	} else {
		decimalsBig = new(big.Int).Exp(_10b, big.NewInt(int64(decimal)), nil)
	}
	return decimalsBig
}

//
//func foo(numerator, denomiator *big.Int) {
//	quotient := new(big.Int).Quo(numerator, denomiator)
//	remainder := new(big.Int).Rem(numerator, denomiator)
//	gcd := new(big.Int).GCD(nil, nil, remainder, denomiator)
//
//}
