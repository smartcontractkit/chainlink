# Fast Implementation of Base58 encoding
[![GoDoc](https://godoc.org/github.com/mr-tron/base58?status.svg)](https://godoc.org/github.com/mr-tron/base58)  [![Go Report Card](https://goreportcard.com/badge/github.com/mr-tron/base58)](https://goreportcard.com/report/github.com/mr-tron/base58)
[![Used By](https://sourcegraph.com/github.com/mr-tron/base58/-/badge.svg)](https://sourcegraph.com/github.com/mr-tron/base58?badge)

Fast implementation of base58 encoding in Go. 

Base algorithm is copied from https://github.com/trezor/trezor-crypto/blob/master/base58.c

## Benchmark
Trivial - encoding via big.Int (over libraries use this implemenation)
Fast - optimized algorythm from trezor

```
BenchmarkTrivialBase58Encoding-4          123063              9568 ns/op
BenchmarkFastBase58Encoding-4             690040              1598 ns/op

BenchmarkTrivialBase58Decoding-4          275216              4301 ns/op
BenchmarkFastBase58Decoding-4            1812105               658 ns/op
```
Encoding - **faster by 6 times**

Decoding - **faster by 6 times**

## Usage example

```go

package main

import (
	"fmt"
	"github.com/mr-tron/base58"
)

func main() {

	encoded := "1QCaxc8hutpdZ62iKZsn1TCG3nh7uPZojq"
	num, err := base58.Decode(encoded)
	if err != nil {
		fmt.Printf("Demo %v, got error %s\n", encoded, err)	
	}
	chk := base58.Encode(num)
	if encoded == string(chk) {
		fmt.Printf ( "Successfully decoded then re-encoded %s\n", encoded )
	} 
}

```
