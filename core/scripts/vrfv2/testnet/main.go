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
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_external_sub_owner_example"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_single_consumer_example"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
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
		coordinatorDeployBHSAddress := coordinatorDeployCmd.String("bhs-address", "", "address of bhs")
		coordinatorDeployLinkEthFeedAddress := coordinatorDeployCmd.String("link-eth-feed", "", "address of link-eth-feed")
		panicErr(coordinatorDeployCmd.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"link-address", "bhs-address", "link-eth-feed"})
		coordinatorAddress, tx, _, err := vrf_coordinator_v2.DeployVRFCoordinatorV2(
			owner,
			ec,
			common.HexToAddress(*coordinatorDeployLinkAddress),
			common.HexToAddress(*coordinatorDeployBHSAddress),
			common.HexToAddress(*coordinatorDeployLinkEthFeedAddress))
		panicErr(err)
		fmt.Println("Coordinator", coordinatorAddress.String(), "hash", tx.Hash())
	case "coordinator-set-config":
		coordinatorSetConfigCmd := flag.NewFlagSet("coordinator-set-config", flag.ExitOnError)
		setConfigAddress := coordinatorSetConfigCmd.String("address", "", "coordinator address")
		// TODO: add config parameters as cli args here
		panicErr(coordinatorSetConfigCmd.Parse(os.Args[2:]))
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*setConfigAddress), ec)
		panicErr(err)
		failIfRequiredArgumentsAreEmpty([]string{"address"})
		tx, err := coordinator.SetConfig(owner,
			uint16(1),                              // minRequestConfirmations
			uint32(1000000),                        // max gas limit
			uint32(60*60*24),                       // stalenessSeconds
			uint32(vrf.GasAfterPaymentCalculation), // gasAfterPaymentCalculation
			big.NewInt(10000000000000000),          // 0.01 eth per link fallbackLinkPrice
			vrf_coordinator_v2.VRFCoordinatorV2FeeConfig{
				FulfillmentFlatFeeLinkPPMTier1: uint32(10000),
				FulfillmentFlatFeeLinkPPMTier2: uint32(1000),
				FulfillmentFlatFeeLinkPPMTier3: uint32(100),
				FulfillmentFlatFeeLinkPPMTier4: uint32(10),
				FulfillmentFlatFeeLinkPPMTier5: uint32(1),
				ReqsForTier2:                   big.NewInt(10),
				ReqsForTier3:                   big.NewInt(20),
				ReqsForTier4:                   big.NewInt(30),
				ReqsForTier5:                   big.NewInt(40),
			},
		)
		panicErr(err)
		fmt.Println("hash", tx.Hash())
	case "coordinator-register-key":
		coordinatorRegisterKey := flag.NewFlagSet("coordinator-register-key", flag.ExitOnError)
		registerKeyAddress := coordinatorRegisterKey.String("address", "", "coordinator address")
		registerKeyUncompressedPubKey := coordinatorRegisterKey.String("pubkey", "", "uncompressed pubkey")
		registerKeyOracleAddress := coordinatorRegisterKey.String("oracle-address", "", "oracle address")
		panicErr(coordinatorRegisterKey.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"address", "pubkey", "oracle-address"})
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*registerKeyAddress), ec)
		panicErr(err)
		pubBytes, err := hex.DecodeString(*registerKeyUncompressedPubKey)
		panicErr(err)
		pk, err := crypto.UnmarshalPubkey(pubBytes)
		panicErr(err)
		tx, err := coordinator.RegisterProvingKey(owner,
			common.HexToAddress(*registerKeyOracleAddress),
			[2]*big.Int{pk.X, pk.Y})
		panicErr(err)
		fmt.Println("hash", tx.Hash())
	case "coordinator-subscription":
		coordinatorSub := flag.NewFlagSet("coordinator-subscription", flag.ExitOnError)
		address := coordinatorSub.String("address", "", "coordinator address")
		subID := coordinatorSub.Int64("sub", 0, "subID")
		panicErr(coordinatorSub.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"address", "pubkey"})
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*address), ec)
		panicErr(err)
		fmt.Println("subID", *subID, "address", *address, coordinator.Address())
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
		failIfRequiredArgumentsAreEmpty([]string{"coordinator-address", "key-hash", "link-address"})
		keyHashBytes := common.HexToHash(*keyHash)
		consumerAddress, tx, _, err := vrf_single_consumer_example.DeployVRFSingleConsumerExample(
			owner,
			ec,
			common.HexToAddress(*consumerCoordinator),
			common.HexToAddress(*consumerLinkAddress),
			uint32(1000000), // gas callback
			uint16(5),       // confs
			uint32(1),       // words
			keyHashBytes)
		panicErr(err)
		fmt.Println("Consumer address", consumerAddress, "hash", tx.Hash())
	case "consumer-subscribe":
		consumerSubscribeCmd := flag.NewFlagSet("consumer-subscribe", flag.ExitOnError)
		consumerSubscribeAddress := consumerSubscribeCmd.String("address", "", "consumer address")
		panicErr(consumerSubscribeCmd.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"address"})
		consumer, err := vrf_single_consumer_example.NewVRFSingleConsumerExample(common.HexToAddress(*consumerSubscribeAddress), ec)
		panicErr(err)
		tx, err := consumer.Subscribe(owner)
		panicErr(err)
		fmt.Println("hash", tx.Hash())
	case "link-balance":
		linkBalanceCmd := flag.NewFlagSet("link-balance", flag.ExitOnError)
		linkAddress := linkBalanceCmd.String("link-address", "", "link-address")
		address := linkBalanceCmd.String("address", "", "address")
		panicErr(linkBalanceCmd.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"link-address", "address"})
		lt, err := link_token_interface.NewLinkToken(common.HexToAddress(*linkAddress), ec)
		panicErr(err)
		b, err := lt.BalanceOf(nil, common.HexToAddress(*address))
		panicErr(err)
		fmt.Println(b)
	case "consumer-cancel":
		consumerCancelCmd := flag.NewFlagSet("consumer-cancel", flag.ExitOnError)
		consumerCancelAddress := consumerCancelCmd.String("address", "", "consumer address")
		panicErr(consumerCancelCmd.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"address"})
		consumer, err := vrf_single_consumer_example.NewVRFSingleConsumerExample(common.HexToAddress(*consumerCancelAddress), ec)
		panicErr(err)
		tx, err := consumer.Unsubscribe(owner, owner.From)
		panicErr(err)
		fmt.Println("hash", tx.Hash())
	case "consumer-topup":
		// NOTE NEED TO FUND CONSUMER WITH LINK FIRST
		consumerTopupCmd := flag.NewFlagSet("consumer-topup", flag.ExitOnError)
		consumerTopupAmount := consumerTopupCmd.String("amount", "", "amount")
		consumerTopupAddress := consumerTopupCmd.String("address", "", "consumer address")
		panicErr(consumerTopupCmd.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"amount", "address"})
		consumer, err := vrf_single_consumer_example.NewVRFSingleConsumerExample(common.HexToAddress(*consumerTopupAddress), ec)
		panicErr(err)
		amount, s := big.NewInt(0).SetString(*consumerTopupAmount, 10)
		if !s {
			panic(fmt.Sprintf("failed to parse top up amount '%s'", *consumerTopupAmount))
		}
		tx, err := consumer.TopUpSubscription(owner, amount)
		panicErr(err)
		fmt.Println("hash", tx.Hash())
	case "consumer-request":
		consumerRequestCmd := flag.NewFlagSet("consumer-request", flag.ExitOnError)
		consumerRequestAddress := consumerRequestCmd.String("address", "", "consumer address")
		panicErr(consumerRequestCmd.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"address"})
		consumer, err := vrf_single_consumer_example.NewVRFSingleConsumerExample(common.HexToAddress(*consumerRequestAddress), ec)
		panicErr(err)
		// Fund and request 1 link
		tx, err := consumer.RequestRandomWords(owner)
		panicErr(err)
		fmt.Println("tx", tx.Hash())
	case "consumer-fund-and-request":
		consumerRequestCmd := flag.NewFlagSet("consumer-request", flag.ExitOnError)
		consumerRequestAddress := consumerRequestCmd.String("address", "", "consumer address")
		panicErr(consumerRequestCmd.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"address"})
		consumer, err := vrf_single_consumer_example.NewVRFSingleConsumerExample(common.HexToAddress(*consumerRequestAddress), ec)
		panicErr(err)
		// Fund and request 3 link
		tx, err := consumer.FundAndRequestRandomWords(owner, big.NewInt(3000000000000000000))
		panicErr(err)
		fmt.Println("tx", tx.Hash())
	case "consumer-print":
		consumerPrint := flag.NewFlagSet("consumer-print", flag.ExitOnError)
		address := consumerPrint.String("address", "", "consumer address")
		panicErr(consumerPrint.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"address"})
		consumer, err := vrf_single_consumer_example.NewVRFSingleConsumerExample(common.HexToAddress(*address), ec)
		panicErr(err)
		rc, err := consumer.SRequestConfig(nil)
		panicErr(err)
		rw, err := consumer.SRandomWords(nil, big.NewInt(0))
		if err != nil {
			fmt.Println("no words")
		}
		rid, err := consumer.SRequestId(nil)
		panicErr(err)
		fmt.Printf("Request config %+v Rw %+v Rid %+v\n", rc, rw, rid)
	case "eoa-consumer-deploy":
		consumerDeployCmd := flag.NewFlagSet("eoa-consumer-deploy", flag.ExitOnError)
		consumerCoordinator := consumerDeployCmd.String("coordinator-address", "", "coordinator address")
		consumerLinkAddress := consumerDeployCmd.String("link-address", "", "link-address")
		panicErr(consumerDeployCmd.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"coordinator-address", "link-address"})
		consumerAddress, tx, _, err := vrf_external_sub_owner_example.DeployVRFExternalSubOwnerExample(
			owner,
			ec,
			common.HexToAddress(*consumerCoordinator),
			common.HexToAddress(*consumerLinkAddress))
		panicErr(err)
		fmt.Println("Consumer address", consumerAddress, "hash", tx.Hash())
	case "eoa-create-sub":
		createSubCmd := flag.NewFlagSet("eoa-create-sub", flag.ExitOnError)
		coordinatorAddress := createSubCmd.String("coordinator-address", "", "coordinator address")
		panicErr(createSubCmd.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"coordinator-address"})
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddress), ec)
		panicErr(err)
		tx, err := coordinator.CreateSubscription(owner)
		panicErr(err)
		fmt.Println("Create subscription", "TX hash", tx.Hash())
	case "eoa-add-sub-consumer":
		addSubConsCmd := flag.NewFlagSet("eoa-add-sub-consumer", flag.ExitOnError)
		coordinatorAddress := addSubConsCmd.String("coordinator-address", "", "coordinator address")
		subID := addSubConsCmd.Uint64("sub-id", 0, "subID")
		consumerAddress := addSubConsCmd.String("consumer-address", "", "consumer address")
		panicErr(addSubConsCmd.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"coordinator-address", "consumer-address"})
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddress), ec)
		panicErr(err)
		txadd, err := coordinator.AddConsumer(owner, *subID, common.HexToAddress(*consumerAddress))
		panicErr(err)
		fmt.Println("Adding consumer", "TX hash", txadd.Hash())
	case "eoa-create-fund-authorize-sub":
		// Lets just treat the owner key as the EOA controlling the sub
		cfaSubCmd := flag.NewFlagSet("eoa-create-fund-authorize-sub", flag.ExitOnError)
		coordinatorAddress := cfaSubCmd.String("coordinator-address", "", "coordinator address")
		amountStr := cfaSubCmd.String("amount", "", "amount to fund")
		consumerAddress := cfaSubCmd.String("consumer-address", "", "consumer address")
		consumerLinkAddress := cfaSubCmd.String("link-address", "", "link-address")
		panicErr(cfaSubCmd.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"coordinator-address", "amount", "consumer-address", "link-address"})
		amount, s := big.NewInt(0).SetString(*amountStr, 10)
		if !s {
			panic(fmt.Sprintf("failed to parse top up amount '%s'", *amountStr))
		}
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddress), ec)
		panicErr(err)
		fmt.Println(amount, consumerLinkAddress)
		txcreate, err := coordinator.CreateSubscription(owner)
		panicErr(err)
		fmt.Println("Create sub", "hash", txcreate.Hash())
		sub := make(chan *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCreated)
		subscription, err := coordinator.WatchSubscriptionCreated(nil, sub, nil)
		panicErr(err)
		defer subscription.Unsubscribe()
		created := <-sub
		linkToken, err := link_token_interface.NewLinkToken(common.HexToAddress(*consumerLinkAddress), ec)
		panicErr(err)
		bal, err := linkToken.BalanceOf(nil, owner.From)
		panicErr(err)
		fmt.Println("OWNER BALANCE", bal, owner.From.String(), amount.String())
		b, err := utils.GenericEncode([]string{"uint64"}, created.SubId)
		panicErr(err)
		owner.GasLimit = 500000
		tx, err := linkToken.TransferAndCall(owner, coordinator.Address(), amount, b)
		panicErr(err)
		fmt.Println("Funding sub", created.SubId, "hash", tx.Hash())
		subFunded := make(chan *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionFunded)
		fundSub, err := coordinator.WatchSubscriptionFunded(nil, subFunded, []uint64{created.SubId})
		panicErr(err)
		defer fundSub.Unsubscribe()
		<-subFunded // Add a consumer once its funded
		txadd, err := coordinator.AddConsumer(owner, created.SubId, common.HexToAddress(*consumerAddress))
		panicErr(err)
		fmt.Println("adding consumer", "hash", txadd.Hash())
	case "eoa-request":
		request := flag.NewFlagSet("eoa-request", flag.ExitOnError)
		consumerAddress := request.String("consumer-address", "", "consumer address")
		subID := request.Uint64("sub-id", 0, "subscription ID")
		cbGasLimit := request.Uint("cb-gas-limit", 1_000_000, "callback gas limit")
		requestConfirmations := request.Uint("request-confirmations", 3, "minimum request confirmations")
		numWords := request.Uint("num-words", 3, "number of words to request")
		keyHash := request.String("key-hash", "", "key hash")
		panicErr(request.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"consumer-address", "key-hash"})
		keyHashBytes := common.HexToHash(*keyHash)
		consumer, err := vrf_external_sub_owner_example.NewVRFExternalSubOwnerExample(
			common.HexToAddress(*consumerAddress),
			ec)
		panicErr(err)
		tx, err := consumer.RequestRandomWords(owner, *subID, uint32(*cbGasLimit), uint16(*requestConfirmations), uint32(*numWords), keyHashBytes)
		panicErr(err)
		fmt.Println("TX hash:", tx.Hash())
	case "eoa-transfer-sub":
		trans := flag.NewFlagSet("eoa-transfer-sub", flag.ExitOnError)
		coordinatorAddress := trans.String("coordinator-address", "", "coordinator address")
		subID := trans.Int64("subID", 0, "subID")
		to := trans.String("to", "", "to")
		panicErr(trans.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"coordinator-address", "to"})
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddress), ec)
		panicErr(err)
		tx, err := coordinator.RequestSubscriptionOwnerTransfer(owner, uint64(*subID), common.HexToAddress(*to))
		panicErr(err)
		fmt.Println("ownership transfer requested", tx.Hash())
	case "eoa-accept-sub":
		accept := flag.NewFlagSet("eoa-accept-sub", flag.ExitOnError)
		coordinatorAddress := accept.String("coordinator-address", "", "coordinator address")
		subID := accept.Int64("subID", 0, "subID")
		panicErr(accept.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"coordinator-address"})
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddress), ec)
		panicErr(err)
		tx, err := coordinator.AcceptSubscriptionOwnerTransfer(owner, uint64(*subID))
		panicErr(err)
		fmt.Println("ownership transfer accepted", tx.Hash())
	case "eoa-cancel-sub":
		cancel := flag.NewFlagSet("eoa-cancel-sub", flag.ExitOnError)
		coordinatorAddress := cancel.String("coordinator-address", "", "coordinator address")
		subID := cancel.Int64("subID", 0, "subID")
		panicErr(cancel.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"coordinator-address"})
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddress), ec)
		panicErr(err)
		tx, err := coordinator.CancelSubscription(owner, uint64(*subID), owner.From)
		panicErr(err)
		fmt.Println("sub cancelled", tx.Hash())
	case "eoa-fund-sub":
		fund := flag.NewFlagSet("eoa-fund-sub", flag.ExitOnError)
		coordinatorAddress := fund.String("coordinator-address", "", "coordinator address")
		amountStr := fund.String("amount", "", "amount to fund")
		subID := fund.Int64("sub-id", 0, "subID")
		consumerLinkAddress := fund.String("link-address", "", "link-address")
		panicErr(fund.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"coordinator-address", "amount", "link-address"})
		amount, s := big.NewInt(0).SetString(*amountStr, 10)
		if !s {
			panic(fmt.Sprintf("failed to parse top up amount '%s'", *amountStr))
		}
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddress), ec)
		panicErr(err)
		linkToken, err := link_token_interface.NewLinkToken(common.HexToAddress(*consumerLinkAddress), ec)
		panicErr(err)
		bal, err := linkToken.BalanceOf(nil, owner.From)
		panicErr(err)
		fmt.Println("Initial account balance:", bal, owner.From.String(), "Funding amount:", amount.String())
		b, err := utils.GenericEncode([]string{"uint64"}, uint64(*subID))
		panicErr(err)
		owner.GasLimit = 500000
		tx, err := linkToken.TransferAndCall(owner, coordinator.Address(), amount, b)
		panicErr(err)
		fmt.Println("Funding sub", *subID, "hash", tx.Hash())
		panicErr(err)
	case "owner-cancel-sub":
		cancel := flag.NewFlagSet("owner-cancel-sub", flag.ExitOnError)
		coordinatorAddress := cancel.String("coordinator-address", "", "coordinator address")
		subID := cancel.Int64("subID", 0, "subID")
		panicErr(cancel.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"coordinator-address"})
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddress), ec)
		panicErr(err)
		tx, err := coordinator.OwnerCancelSubscription(owner, uint64(*subID))
		panicErr(err)
		fmt.Println("sub cancelled", tx.Hash())
	case "sub-balance":
		consumerBalanceCmd := flag.NewFlagSet("sub-balance", flag.ExitOnError)
		coordinatorAddress := consumerBalanceCmd.String("coordinator-address", "", "coordinator address")
		subID := consumerBalanceCmd.Uint64("sub-id", 0, "subscription id")
		panicErr(consumerBalanceCmd.Parse(os.Args[2:]))
		failIfRequiredArgumentsAreEmpty([]string{"coordinator-address"})
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddress), ec)
		panicErr(err)
		resp, err := coordinator.GetSubscription(nil, *subID)
		panicErr(err)
		fmt.Println("sub id", *subID, "balance:", resp.Balance)
	default:
		panic("unrecognized subcommand: " + os.Args[1])
	}
}
