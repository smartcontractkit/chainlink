package main

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_consumer_v2"
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	consumerAddress := "0x9E79d9A7F68D136ec4c1C0187B97c271CEa6008B"
	ec, err := ethclient.Dial("ws://127.0.0.1:8546")
	panicErr(err)
	consumer, err := vrf_consumer_v2.NewVRFConsumerV2(common.HexToAddress(consumerAddress), ec)
	panicErr(err)
	var rw []*big.Int
	nw := 3
	var r *big.Int
	for i := 0; i < nw; i++ {
		r, err = consumer.SRandomWords(nil, big.NewInt(int64(i)))
		panicErr(err)
		rw = append(rw, r)
	}
	gasAvail, err := consumer.SGasAvailable(nil)
	panicErr(err)
	fmt.Println("Random words", rw)
	fmt.Println("Gas available", gasAvail)
}
