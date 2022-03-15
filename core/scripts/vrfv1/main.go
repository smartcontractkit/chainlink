package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"

	linktoken "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	vrfltoc "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_load_test_ownerless_consumer"
	vrfoc "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_ownerless_consumer_example"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func main() {
	ethURL, set := os.LookupEnv("ETH_URL")
	if !set {
		panic("need eth url")
	}

	chainIDEnv, set := os.LookupEnv("ETH_CHAIN_ID")
	if !set {
		panic("need chain ID")
	}

	accountKey, set := os.LookupEnv("ACCOUNT_KEY")
	if !set {
		panic("need account key")
	}

	if len(os.Args) < 2 {
		fmt.Println("expected subcommand")
		os.Exit(1)
	}

	ec, err := ethclient.Dial(ethURL)
	helpers.PanicErr(err)

	chainID, err := strconv.ParseInt(chainIDEnv, 10, 64)
	helpers.PanicErr(err)

	// Owner key. Make sure it has eth
	b, err := hex.DecodeString(accountKey)
	helpers.PanicErr(err)
	d := new(big.Int).SetBytes(b)

	pkX, pkY := crypto.S256().ScalarBaseMult(d.Bytes())
	privateKey := ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: crypto.S256(),
			X:     pkX,
			Y:     pkY,
		},
		D: d,
	}
	account, err := bind.NewKeyedTransactorWithChainID(&privateKey, big.NewInt(chainID))
	helpers.PanicErr(err)

	// Explicitly set gas price to ensure non-eip 1559
	gp, err := ec.SuggestGasPrice(context.Background())
	helpers.PanicErr(err)
	account.GasPrice = gp

	// Uncomment the block below if transactions are not getting picked up due to nonce issues:
	//
	//block, err := ec.BlockNumber(context.Background())
	//helpers.PanicErr(err)
	//
	//nonce, err := ec.NonceAt(context.Background(), account.From, big.NewInt(int64(block)))
	//helpers.PanicErr(err)
	//
	//account.Nonce = big.NewInt(int64(nonce))
	//account.GasPrice = gp.Mul(gp, big.NewInt(2))

	switch os.Args[1] {
	case "ownerless-consumer-deploy":
		cmd := flag.NewFlagSet("ownerless-consumer-deploy", flag.ExitOnError)
		coordAddr := cmd.String("coordinator-address", "", "address of VRF coordinator")
		linkAddr := cmd.String("link-address", "", "address of link token")
		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address", "link-address")
		consumerAddr, tx, _, err := vrfoc.DeployVRFOwnerlessConsumerExample(
			account,
			ec,
			common.HexToAddress(*coordAddr),
			common.HexToAddress(*linkAddr))
		helpers.PanicErr(err)
		fmt.Printf("Ownerless Consumer: %s TX: %s\n",
			consumerAddr, helpers.ExplorerLink(chainID, tx.Hash()))
	case "loadtest-ownerless-consumer-deploy":
		cmd := flag.NewFlagSet("loadtest-ownerless-consumer-deploy", flag.ExitOnError)
		coordAddr := cmd.String("coordinator-address", "", "address of VRF coordinator")
		linkAddr := cmd.String("link-address", "", "address of link token")
		priceStr := cmd.String("price", "", "the price of each VRF request in Juels")
		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address", "link-address")
		price := decimal.RequireFromString(*priceStr).BigInt()
		consumerAddr, tx, _, err := vrfltoc.DeployVRFLoadTestOwnerlessConsumer(
			account,
			ec,
			common.HexToAddress(*coordAddr),
			common.HexToAddress(*linkAddr),
			price)
		helpers.PanicErr(err)
		fmt.Printf("Loadtest Ownerless Consumer: %s TX: %s\n",
			consumerAddr, helpers.ExplorerLink(chainID, tx.Hash()))
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
		link, err := linktoken.NewLinkToken(common.HexToAddress(*linkAddr), ec)
		helpers.PanicErr(err)
		data, err := utils.GenericEncode([]string{"bytes32"}, common.HexToHash(*keyHash))
		helpers.PanicErr(err)
		tx, err := link.TransferAndCall(account, common.HexToAddress(*consumerAddr), payment, data)
		helpers.PanicErr(err)
		fmt.Printf("TX: %s\n", helpers.ExplorerLink(chainID, tx.Hash()))
	}
}
