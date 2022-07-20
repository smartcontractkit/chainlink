package main

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_consumer_v2"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
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
	coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress("0x569620752AbB8a31AC2832b6680efAFcB399E37e"), ec)
	panicErr(err)
	sub, err := coordinator.GetSubscription(nil, 1)
	fmt.Println(sub, err)
	_, _, kh, err := coordinator.GetRequestConfig(nil)
	fmt.Println(kh, err)

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
