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
	linktoken "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	vrfoc "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_ownerless_consumer_example"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func failIfRequiredArgumentsAreEmpty(required []string) {
	flag.Parse()
	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			panicErr(fmt.Errorf("missing required -%s argument/flag", req))
		}
	}
}

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
	panicErr(err)

	chainID, err := strconv.ParseInt(chainIDEnv, 10, 64)
	panicErr(err)

	// Owner key. Make sure it has eth
	b, err := hex.DecodeString(accountKey)
	panicErr(err)
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
	panicErr(err)

	// Explicitly set gas price to ensure non-eip 1559
	gp, err := ec.SuggestGasPrice(context.Background())
	panicErr(err)
	account.GasPrice = gp

	switch os.Args[1] {
	case "ownerless-consumer-deploy":
		cmd := flag.NewFlagSet("ownerless-consumer-deploy", flag.ExitOnError)
		coordAddr := cmd.String("coordinator-address", "", "address of VRF coordinator")
		linkAddr := cmd.String("link-address", "", "address of link token")
		panicErr(cmd.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"coordinator-address", "link-address"})
		consumerAddr, tx, _, err := vrfoc.DeployVRFOwnerlessConsumerExample(
			account,
			ec,
			common.HexToAddress(*coordAddr),
			common.HexToAddress(*linkAddr))
		panicErr(err)
		fmt.Printf("Ownerless Consumer: %s TX Hash: %s\n", consumerAddr, tx.Hash())
	case "ownerless-consumer-request":
		cmd := flag.NewFlagSet("ownerless-consumer-deploy", flag.ExitOnError)
		linkAddr := cmd.String("link-address", "", "address of link token")
		consumerAddr := cmd.String("consumer-address", "", "address of the deployed ownerless consumer")
		paymentStr := cmd.String("payment", "100000000000000000" /* 0.1 LINK */, "the payment amount in LINK")
		keyHash := cmd.String("key-hash", "", "key hash")
		panicErr(cmd.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"link-address", "consumer-address", "key-hash"})
		payment, ok := big.NewInt(0).SetString(*paymentStr, 10)
		if !ok {
			panic(fmt.Sprintf("failed to parse payment amount: %s", *paymentStr))
		}
		link, err := linktoken.NewLinkToken(common.HexToAddress(*linkAddr), ec)
		panicErr(err)
		data, err := utils.GenericEncode([]string{"bytes32"}, common.HexToHash(*keyHash))
		panicErr(err)
		tx, err := link.TransferAndCall(account, common.HexToAddress(*consumerAddr), payment, data)
		panicErr(err)
		fmt.Printf("TX Hash: %s\n", tx.Hash())
	}
}
