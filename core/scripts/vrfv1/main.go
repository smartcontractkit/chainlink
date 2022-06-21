package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"

	linktoken "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	vrfltoc "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_load_test_ownerless_consumer"
	vrfoc "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_ownerless_consumer_example"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func main() {
	e := helpers.SetupEnv()

	switch os.Args[1] {
	case "ownerless-consumer-deploy":
		cmd := flag.NewFlagSet("ownerless-consumer-deploy", flag.ExitOnError)
		coordAddr := cmd.String("coordinator-address", "", "address of VRF coordinator")
		linkAddr := cmd.String("link-address", "", "address of link token")
		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address", "link-address")
		consumerAddr, tx, _, err := vrfoc.DeployVRFOwnerlessConsumerExample(
			e.Owner,
			e.Ec,
			common.HexToAddress(*coordAddr),
			common.HexToAddress(*linkAddr))
		helpers.PanicErr(err)
		fmt.Printf("Ownerless Consumer: %s TX: %s\n",
			consumerAddr, helpers.ExplorerLink(e.ChainID, tx.Hash()))
	case "loadtest-ownerless-consumer-deploy":
		cmd := flag.NewFlagSet("loadtest-ownerless-consumer-deploy", flag.ExitOnError)
		coordAddr := cmd.String("coordinator-address", "", "address of VRF coordinator")
		linkAddr := cmd.String("link-address", "", "address of link token")
		priceStr := cmd.String("price", "", "the price of each VRF request in Juels")
		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address", "link-address")
		price := decimal.RequireFromString(*priceStr).BigInt()
		consumerAddr, tx, _, err := vrfltoc.DeployVRFLoadTestOwnerlessConsumer(
			e.Owner,
			e.Ec,
			common.HexToAddress(*coordAddr),
			common.HexToAddress(*linkAddr),
			price)
		helpers.PanicErr(err)
		fmt.Printf("Loadtest Ownerless Consumer: %s TX: %s\n",
			consumerAddr, helpers.ExplorerLink(e.ChainID, tx.Hash()))
	case "ownerless-consumer-request":
		cmd := flag.NewFlagSet("ownerless-consumer-request", flag.ExitOnError)
		linkAddr := cmd.String("link-address", "", "address of link token")
		consumerAddr := cmd.String("consumer-address", "", "address of the deployed ownerless consumer")
		paymentStr := cmd.String("payment", "" /* 0.1 LINK */, "the payment amount in LINK")
		keyHash := cmd.String("key-hash", "", "key hash")
		helpers.ParseArgs(cmd, os.Args[2:], "link-address", "consumer-address", "payment", "key-hash")
		payment, ok := big.NewInt(0).SetString(*paymentStr, 10)
		if !ok {
			panic(fmt.Sprintf("failed to parse payment amount: %s", *paymentStr))
		}
		link, err := linktoken.NewLinkToken(common.HexToAddress(*linkAddr), e.Ec)
		helpers.PanicErr(err)
		data, err := utils.GenericEncode([]string{"bytes32"}, common.HexToHash(*keyHash))
		helpers.PanicErr(err)
		tx, err := link.TransferAndCall(e.Owner, common.HexToAddress(*consumerAddr), payment, data)
		helpers.PanicErr(err)
		fmt.Printf("TX: %s\n", helpers.ExplorerLink(e.ChainID, tx.Hash()))
	}
}
