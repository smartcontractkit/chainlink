package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_single_consumer_example"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	ethURL, set := os.LookupEnv("ETH_URL")
	if !set {
		panic("need eth url")
	}
	ownerKey, set := os.LookupEnv("OWNER_KEY")
	if !set {
		panic("need owner key")
	}

	if len(os.Args) < 2 {
		fmt.Println("expected subcommand")
		os.Exit(1)
	}
	ec, err := ethclient.Dial(ethURL)
	panicErr(err)
	chainID := int64(42)

	// Owner key. Make sure it has eth
	d, success := new(big.Int).SetString(ownerKey, 10)
	if !success {
		panic("failed to parse key")
	}
	pkX, pkY := crypto.S256().ScalarBaseMult(d.Bytes())
	privateKey := ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: crypto.S256(),
			X:     pkX,
			Y:     pkY,
		},
		D: d,
	}
	owner, err := bind.NewKeyedTransactorWithChainID(&privateKey, big.NewInt(chainID))
	panicErr(err)
	// Explicitly set gas price to ensure non-eip 1559
	gp, err := ec.SuggestGasPrice(context.Background())
	panicErr(err)
	owner.GasPrice = gp
	switch os.Args[1] {
	case "coordinator-deploy":
		coordinatorDeployCmd := flag.NewFlagSet("coordinator-deploy", flag.ExitOnError)
		coordinatorDeployLinkAddress := coordinatorDeployCmd.String("link-address", "", "address of link token")
		//coordinatorDeployBHSAddress := coordinatorDeployCmd.String("bhs-address", "", "address of bhs")
		coordinatorDeployLinkEthFeedAddress := coordinatorDeployCmd.String("link-eth-feed", "", "address of link-eth-feed")
		panicErr(coordinatorDeployCmd.Parse(os.Args[2:]))
		coordinatorAddress, _, _, err := vrf_coordinator_v2.DeployVRFCoordinatorV2(
			owner,
			ec,
			common.HexToAddress(*coordinatorDeployLinkAddress),
			common.Address{}, // TODO test this
			common.HexToAddress(*coordinatorDeployLinkEthFeedAddress))
		panicErr(err)
		fmt.Println("Coordinator", coordinatorAddress.String())
	case "coordinator-set-config":
		coordinatorSetConfigCmd := flag.NewFlagSet("coordinator-set-config", flag.ExitOnError)
		setConfigAddress := coordinatorSetConfigCmd.String("address", "", "coordinator address")
		// TODO: add config parameters as cli args here
		panicErr(coordinatorSetConfigCmd.Parse(os.Args[2:]))
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*setConfigAddress), ec)
		panicErr(err)
		_, err = coordinator.SetConfig(owner,
			uint16(1),                              // minRequestConfirmations
			uint32(1000),                           // 0.0001 link flat fee
			uint32(1000000),                        // max gas limit
			uint32(60*60*24),                       // stalenessSeconds
			uint32(vrf.GasAfterPaymentCalculation), // gasAfterPaymentCalculation
			big.NewInt(10000000000000000),          // 0.01 eth per link fallbackLinkPrice
		)
		panicErr(err)
	case "coordinator-register-key":
		coordinatorRegisterKey := flag.NewFlagSet("coordinator-register-key", flag.ExitOnError)
		registerKeyAddress := coordinatorRegisterKey.String("address", "", "coordinator address")
		registerKeyUncompressedPubKey := coordinatorRegisterKey.String("pubkey", "", "uncompressed pubkey")
		registerKeyOracleAddress := coordinatorRegisterKey.String("oracle-address", "", "oracle address")
		panicErr(coordinatorRegisterKey.Parse(os.Args[2:]))
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*registerKeyAddress), ec)
		panicErr(err)
		pubBytes, err := hex.DecodeString(*registerKeyUncompressedPubKey)
		panicErr(err)
		pk, err := crypto.UnmarshalPubkey(pubBytes)
		panicErr(err)
		_, err = coordinator.RegisterProvingKey(owner,
			common.HexToAddress(*registerKeyOracleAddress),
			[2]*big.Int{pk.X, pk.Y})
		panicErr(err)
	case "coordinator-subscription":
		coordinatorSub := flag.NewFlagSet("coordinator-subscription", flag.ExitOnError)
		address := coordinatorSub.String("address", "", "coordinator address")
		subID := coordinatorSub.Int64("sub", 0, "subID")
		panicErr(coordinatorSub.Parse(os.Args[2:]))
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*address), ec)
		panicErr(err)
		s, err := coordinator.GetSubscription(nil, uint64(*subID))
		panicErr(err)
		fmt.Printf("Subscription %+v\n", s)
	case "consumer-deploy":
		consumerDeployCmd := flag.NewFlagSet("consumer-deploy", flag.ExitOnError)
		consumerCoordinator := consumerDeployCmd.String("coordinator-address", "", "coordinator address")
		keyHash := consumerDeployCmd.String("key-hash", "", "key hash")
		consumerLinkAddress := consumerDeployCmd.String("link-address", "", "link-address")
		// TODO: add other params
		panicErr(consumerDeployCmd.Parse(os.Args[2:]))
		keyHashBytes := common.HexToHash(*keyHash)
		consumerAddress, _, _, err := vrf_single_consumer_example.DeployVRFSingleConsumerExample(
			owner,
			ec,
			common.HexToAddress(*consumerCoordinator),
			common.HexToAddress(*consumerLinkAddress),
			uint32(300000), // gas callback
			uint16(5),      // confs
			uint32(1),      // words
			keyHashBytes)
		panicErr(err)
		fmt.Println("Consumer address", consumerAddress)
	case "consumer-subscribe":
		consumerSubscribeCmd := flag.NewFlagSet("consumer-subscribe", flag.ExitOnError)
		consumerSubscribeAddress := consumerSubscribeCmd.String("address", "", "consumer address")
		panicErr(consumerSubscribeCmd.Parse(os.Args[2:]))
		consumer, err := vrf_single_consumer_example.NewVRFSingleConsumerExample(common.HexToAddress(*consumerSubscribeAddress), ec)
		panicErr(err)
		_, err = consumer.Subscribe(owner)
		panicErr(err)
	case "consumer-topup":
		consumerTopupCmd := flag.NewFlagSet("consumer-topup", flag.ExitOnError)
		consumerTopupAmount := consumerTopupCmd.String("amount", "", "amount")
		consumerTopupAddress := consumerTopupCmd.String("address", "", "consumer address")
		panicErr(consumerTopupCmd.Parse(os.Args[2:]))
		consumer, err := vrf_single_consumer_example.NewVRFSingleConsumerExample(common.HexToAddress(*consumerTopupAddress), ec)
		panicErr(err)
		amount, s := big.NewInt(0).SetString(*consumerTopupAmount, 10)
		if !s {
			panic("failed to parse top up amount")
		}
		_, err = consumer.TopUpSubscription(owner, amount)
		panicErr(err)
	case "consumer-request":
		consumerRequestCmd := flag.NewFlagSet("consumer-request", flag.ExitOnError)
		consumerRequestAddress := consumerRequestCmd.String("address", "", "consumer address")
		panicErr(consumerRequestCmd.Parse(os.Args[2:]))
		consumer, err := vrf_single_consumer_example.NewVRFSingleConsumerExample(common.HexToAddress(*consumerRequestAddress), ec)
		panicErr(err)
		tx, err := consumer.RequestRandomWords(owner)
		panicErr(err)
		fmt.Println("tx", tx.Hash())
	case "consumer-print":
		consumerPrint := flag.NewFlagSet("consumer-print", flag.ExitOnError)
		address := consumerPrint.String("address", "", "consumer address")
		panicErr(consumerPrint.Parse(os.Args[2:]))
		consumer, err := vrf_single_consumer_example.NewVRFSingleConsumerExample(common.HexToAddress(*address), ec)
		panicErr(err)
		rc, err := consumer.SRequestConfig(nil)
		panicErr(err)
		rw, err := consumer.SRandomWords(nil, big.NewInt(0))
		panicErr(err)
		rid, err := consumer.SRequestId(nil)
		panicErr(err)
		fmt.Printf("Request config %+v Rw %+v Rid %+v\n", rc, rw, rid)
	}
}
