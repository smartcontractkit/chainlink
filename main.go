package main

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	cborv2 "github.com/fxamacker/cbor/v2"

	"github.com/smartcontractkit/chainlink/v2/core/cbor"
)

type UpkeepOffchainConfig struct {
	MaxGasPrice *big.Int `json:"maxGasPrice" cbor:"maxGasPrice"`
}

type UpkeepOffchainConfig1 struct {
	MaxGasPrice *big.Int `json:"maxGasPrice" cbor:"maxGasPrice"`
	MinGasPrice *big.Int `json:"minGasPrice" cbor:"minGasPrice"`
}

//go:generate make modgraph
func main() {
	data, _ := cborv2.Marshal(UpkeepOffchainConfig{MaxGasPrice: big.NewInt(100_000_000)})

	fmt.Println(data)
	fmt.Println(hexutil.Encode(data))
	fmt.Println(hex.EncodeToString(data))

	var uoc1 UpkeepOffchainConfig
	_ = cbor.ParseDietCBORToStruct(data, &uoc1)
	fmt.Println(uoc1.MaxGasPrice) // 1000000000

	var uoc2 UpkeepOffchainConfig
	s := "6b6d61784761735072696365c258200000000000000000000000000000000000000000000000000000000005f5e100"
	newD, _ := hex.DecodeString(s)
	fmt.Println(newD)
	_ = cbor.ParseDietCBORToStruct(newD, &uoc2)
	fmt.Println(uoc2.MaxGasPrice) // 1000000000
}
