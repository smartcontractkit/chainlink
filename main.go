package main

import (
	// "os"
	"context"
	"fmt"
	"math/big"
	"time"

	// "github.com/smartcontractkit/chainlink/v2/core"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/ocr2"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/starknet"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// //go:generate make modgraph
// func main() {
// 	os.Exit(core.Main())
// }

func main() {
	logger := logger.NullLogger
	lggr := logger.With(logger, "starknetChainID")
	timeout := time.Second * 10
	client, err := starknet.NewClient("1234", "https://starknet-sepolia.core.chainstack.com/304b1c13eca4621063f0e75222433570", "", lggr, &timeout)
	if err != nil {
		fmt.Println("error")
		return
	}
	oClient, err := ocr2.NewClient(client, lggr)
	if err != nil {
		fmt.Println("error")
		return
	}

	bigKey, _ := new(big.Int).SetString("79c0bc2a03570241c27235a2dca7696a658cbdaae0bad5762e30204b2791aba", 16)
	addr := new(felt.Felt).SetBytes(bigKey.Bytes())

	events, err := oClient.NewTransmissionsFromEventsAt(context.Background(), addr, 56202)
	if err != nil {
		fmt.Println("event error")
	}
	fmt.Println(events)
	if len(events) == 0 {
		// NOTE This shouldn't happen! LatestRound says this block should have a transmission and we didn't find any!
		fmt.Errorf("no transmissions found in the block")
	}

}
