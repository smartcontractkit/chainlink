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
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	evmutils "github.com/smartcontractkit/chainlink-integrations/evm/utils"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/blockhash_store"
	linktoken "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_coordinator_interface"
	vrfltoc "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_load_test_ownerless_consumer"
	vrfoc "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_ownerless_consumer_example"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func main() {
	e := helpers.SetupEnv(false)

	switch os.Args[1] {
	case "topics":
		randomnessRequestTopic := solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest{}.Topic()
		randomnessFulfilledTopic := solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequestFulfilled{}.Topic()
		fmt.Println("RandomnessRequest:", randomnessRequestTopic.String(),
			"RandomnessRequestFulfilled:", randomnessFulfilledTopic.String())
	case "request-report":
		cmd := flag.NewFlagSet("request-report", flag.ExitOnError)
		txHashes := cmd.String("tx-hashes", "", "comma separated transaction hashes")
		requestIDs := cmd.String("request-ids", "", "comma separated request IDs in hex")
		bhsAddress := cmd.String("bhs-address", "", "BHS contract address")
		coordinatorAddress := cmd.String("coordinator-address", "", "VRF coordinator address")

		helpers.ParseArgs(cmd, os.Args[2:], "tx-hashes", "bhs-address", "request-ids", "coordinator-address")

		hashes := helpers.ParseHashSlice(*txHashes)
		reqIDs := parseRequestIDs(*requestIDs)
		bhs, err := blockhash_store.NewBlockhashStore(
			common.HexToAddress(*bhsAddress),
			e.Ec)
		helpers.PanicErr(err)
		coordinator, err := solidity_vrf_coordinator_interface.NewVRFCoordinator(
			common.HexToAddress(*coordinatorAddress),
			e.Ec)
		helpers.PanicErr(err)

		if len(hashes) != len(reqIDs) {
			panic(fmt.Errorf("len(hashes) [%d] != len(reqIDs) [%d]", len(hashes), len(reqIDs)))
		}

		var bhsMissedBlocks []*big.Int
		for i := range hashes {
			receipt, err := e.Ec.TransactionReceipt(context.Background(), hashes[i])
			helpers.PanicErr(err)

			reqID := reqIDs[i]
			callbacks, err := coordinator.Callbacks(nil, reqID)
			helpers.PanicErr(err)
			fulfilled := utils.IsEmpty(callbacks.SeedAndBlockNum[:])

			_, err = bhs.GetBlockhash(nil, receipt.BlockNumber)
			if err != nil {
				fmt.Println("Blockhash for block", receipt.BlockNumber, "not stored (tx", hashes[i].String(),
					", request ID", hex.EncodeToString(reqID[:]), ", fulfilled:", fulfilled, ")")
				if !fulfilled {
					// not fulfilled and bh not stored means the feeder missed a store
					bhsMissedBlocks = append(bhsMissedBlocks, receipt.BlockNumber)
				}
			} else {
				fmt.Println("Blockhash for block", receipt.BlockNumber, "stored (tx", hashes[i].String(),
					", request ID", hex.EncodeToString(reqID[:]), ", fulfilled:", fulfilled, ")")
			}
		}

		if len(bhsMissedBlocks) == 0 {
			fmt.Println("Didn't miss any bh stores!")
			return
		}
		fmt.Println("Missed stores:")
		for _, blockNumber := range bhsMissedBlocks {
			fmt.Println("\t* ", blockNumber.String())
		}
	case "get-receipt":
		cmd := flag.NewFlagSet("get-tx", flag.ExitOnError)
		txHashes := cmd.String("tx-hashes", "", "comma separated transaction hashes")
		helpers.ParseArgs(cmd, os.Args[2:], "tx-hashes")
		hashes := helpers.ParseHashSlice(*txHashes)

		for _, h := range hashes {
			receipt, err := e.Ec.TransactionReceipt(context.Background(), h)
			helpers.PanicErr(err)
			fmt.Println("Tx", h.String(), "Included in block:", receipt.BlockNumber,
				", blockhash:", receipt.BlockHash.String())
		}
	case "get-callback":
		cmd := flag.NewFlagSet("get-callback", flag.ExitOnError)
		coordinatorAddress := cmd.String("coordinator-address", "", "VRF coordinator address")
		requestIDs := cmd.String("request-ids", "", "comma separated request IDs in hex")
		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address", "request-ids")
		coordinator, err := solidity_vrf_coordinator_interface.NewVRFCoordinator(
			common.HexToAddress(*coordinatorAddress),
			e.Ec)
		helpers.PanicErr(err)
		reqIDs := parseRequestIDs(*requestIDs)
		for _, reqID := range reqIDs {
			callbacks, err := coordinator.Callbacks(nil, reqID)
			helpers.PanicErr(err)
			if utils.IsEmpty(callbacks.SeedAndBlockNum[:]) {
				fmt.Println("request", hex.EncodeToString(reqID[:]), "fulfilled")
			} else {
				fmt.Println("request", hex.EncodeToString(reqID[:]), "not fulfilled")
			}
		}
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

		uid := uuid.MustParse(*jobID)
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
		data, err := evmutils.ABIEncode(`[{"type":"bytes32"}]`, common.HexToHash(*keyHash))
		helpers.PanicErr(err)
		tx, err := link.TransferAndCall(e.Owner, common.HexToAddress(*consumerAddr), payment, data)
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
	case "load-test-read":
		cmd := flag.NewFlagSet("load-test-read", flag.ExitOnError)
		consumerAddress := cmd.String("consumer-address", "", "load test consumer address")
		helpers.ParseArgs(cmd, os.Args[2:], "consumer-address")
		consumer, err := vrfltoc.NewVRFLoadTestOwnerlessConsumer(common.HexToAddress(*consumerAddress), e.Ec)
		helpers.PanicErr(err)
		count, err := consumer.SResponseCount(nil)
		helpers.PanicErr(err)
		fmt.Println("response count:", count.String(), "consumer:", *consumerAddress)
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

func parseRequestIDs(arg string) (ret [][32]byte) {
	split := strings.Split(arg, ",")
	for _, rid := range split {
		if strings.HasPrefix(rid, "0x") {
			rid = strings.Replace(rid, "0x", "", 1)
		}
		reqID, err := hex.DecodeString(rid)
		helpers.PanicErr(err)
		var reqIDFixed [32]byte
		copy(reqIDFixed[:], reqID)
		ret = append(ret, reqIDFixed)
	}
	return
}
