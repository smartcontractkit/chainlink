package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"

	linktoken "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_coordinator_interface"
	vrfltoc "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_load_test_ownerless_consumer"
	vrfoc "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_ownerless_consumer_example"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func main() {
	e := helpers.SetupEnv(false)

	switch os.Args[1] {
	case "coordinator-deploy":
		cmd := flag.NewFlagSet("coordinator-deploy", flag.ExitOnError)
		linkAddress := cmd.String("link-address", "", "LINK token contract address")
		bhsAddress := cmd.String("bhs-address", "", "blockhash store contract address")
		helpers.ParseArgs(cmd, os.Args[2:], "link-address", "bhs-address")
		_, tx, _, err := solidity_vrf_coordinator_interface.DeployVRFCoordinator(
			e.Owner, e.Ec, common.HexToAddress(*linkAddress), common.HexToAddress(*bhsAddress))
		helpers.PanicErr(err)
		helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
	case "coordinator-register-key":
		cmd := flag.NewFlagSet("coordinator-register-key", flag.ExitOnError)
		coordinatorAddress := cmd.String("coordinator-address", "", "address of VRF coordinator")
		pubKeyUncompressed := cmd.String("pubkey-uncompressed", "", "uncompressed VRF public key in hex")
		oracleAddress := cmd.String("oracle-address", "", "oracle address")
		fee := cmd.String("fee", "", "VRF fee in juels")
		jobID := cmd.String("job-id", "", "Job UUID on the chainlink node (UUID)")
		helpers.ParseArgs(cmd, os.Args[2:],
			"coordinator-address", "pubkey-uncompressed", "oracle-address", "fee", "job-id")

		coordinator, err := solidity_vrf_coordinator_interface.NewVRFCoordinator(
			common.HexToAddress(*coordinatorAddress),
			e.Ec)
		helpers.PanicErr(err)

		// Put key in ECDSA format
		if strings.HasPrefix(*pubKeyUncompressed, "0x") {
			*pubKeyUncompressed = strings.Replace(*pubKeyUncompressed, "0x", "04", 1)
		}
		pubBytes, err := hex.DecodeString(*pubKeyUncompressed)
		helpers.PanicErr(err)
		pk, err := crypto.UnmarshalPubkey(pubBytes)
		helpers.PanicErr(err)

		uid, err := uuid.FromString(*jobID)
		helpers.PanicErr(err)
		tx, err := coordinator.RegisterProvingKey(
			e.Owner,
			decimal.RequireFromString(*fee).BigInt(),
			common.HexToAddress(*oracleAddress),
			[2]*big.Int{pk.X, pk.Y},
			job.ExternalJobIDEncodeStringToTopic(uid),
		)
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
	case "ownerless-consumer-deploy":
		cmd := flag.NewFlagSet("ownerless-consumer-deploy", flag.ExitOnError)
		coordAddr := cmd.String("coordinator-address", "", "address of VRF coordinator")
		linkAddr := cmd.String("link-address", "", "address of link token")
		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address", "link-address")
		_, tx, _, err := vrfoc.DeployVRFOwnerlessConsumerExample(
			e.Owner,
			e.Ec,
			common.HexToAddress(*coordAddr),
			common.HexToAddress(*linkAddr))
		helpers.PanicErr(err)
		helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
	case "loadtest-ownerless-consumer-deploy":
		cmd := flag.NewFlagSet("loadtest-ownerless-consumer-deploy", flag.ExitOnError)
		coordAddr := cmd.String("coordinator-address", "", "address of VRF coordinator")
		linkAddr := cmd.String("link-address", "", "address of link token")
		priceStr := cmd.String("price", "", "the price of each VRF request in Juels")
		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address", "link-address")
		price := decimal.RequireFromString(*priceStr).BigInt()
		_, tx, _, err := vrfltoc.DeployVRFLoadTestOwnerlessConsumer(
			e.Owner,
			e.Ec,
			common.HexToAddress(*coordAddr),
			common.HexToAddress(*linkAddr),
			price)
		helpers.PanicErr(err)
		helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
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
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
	case "ownerless-consumer-read":
		cmd := flag.NewFlagSet("ownerless-consumer-read", flag.ExitOnError)
		consumerAddr := cmd.String("consumer-address", "", "address of the deployed ownerless consumer")
		helpers.ParseArgs(cmd, os.Args[2:], "consumer-address")
		consumer, err := vrfoc.NewVRFOwnerlessConsumerExample(
			common.HexToAddress(*consumerAddr),
			e.Ec)
		helpers.PanicErr(err)
		requestID, err := consumer.SRequestId(nil)
		helpers.PanicErr(err)
		fmt.Println("request ID:", requestID)
		output, err := consumer.SRandomnessOutput(nil)
		helpers.PanicErr(err)
		fmt.Println("randomness:", output)
	}
}
