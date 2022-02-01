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
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_external_sub_owner_example"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_single_consumer_example"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
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
	owner, err := bind.NewKeyedTransactorWithChainID(&privateKey, big.NewInt(chainID))
	helpers.PanicErr(err)
	// Explicitly set gas price to ensure non-eip 1559
	gp, err := ec.SuggestGasPrice(context.Background())
	helpers.PanicErr(err)
	owner.GasPrice = gp
	switch os.Args[1] {
	case "bhs-deploy":
		bhsAddress, tx, _, err := blockhash_store.DeployBlockhashStore(owner, ec)
		helpers.PanicErr(err)
		fmt.Println("BlockhashStore", bhsAddress.String(), "TX", helpers.ExplorerLink(chainID, tx.Hash()))
	case "coordinator-deploy":
		coordinatorDeployCmd := flag.NewFlagSet("coordinator-deploy", flag.ExitOnError)
		coordinatorDeployLinkAddress := coordinatorDeployCmd.String("link-address", "", "address of link token")
		coordinatorDeployBHSAddress := coordinatorDeployCmd.String("bhs-address", "", "address of bhs")
		coordinatorDeployLinkEthFeedAddress := coordinatorDeployCmd.String("link-eth-feed", "", "address of link-eth-feed")
		helpers.ParseArgs(coordinatorDeployCmd, os.Args[2:], "link-address", "bhs-address", "link-eth-feed")
		coordinatorAddress, tx, _, err := vrf_coordinator_v2.DeployVRFCoordinatorV2(
			owner,
			ec,
			common.HexToAddress(*coordinatorDeployLinkAddress),
			common.HexToAddress(*coordinatorDeployBHSAddress),
			common.HexToAddress(*coordinatorDeployLinkEthFeedAddress))
		helpers.PanicErr(err)
		fmt.Println("Coordinator", coordinatorAddress.String(), "TX", helpers.ExplorerLink(chainID, tx.Hash()))
	case "coordinator-set-config":
		coordinatorSetConfigCmd := flag.NewFlagSet("coordinator-set-config", flag.ExitOnError)
		setConfigAddress := coordinatorSetConfigCmd.String("address", "", "coordinator address")
		// TODO: add config parameters as cli args here
		helpers.PanicErr(coordinatorSetConfigCmd.Parse(os.Args[2:]))
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*setConfigAddress), ec)
		helpers.PanicErr(err)
		helpers.ParseArgs(coordinatorSetConfigCmd, os.Args[2:], "address")
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
		helpers.PanicErr(err)
		fmt.Println("TX", helpers.ExplorerLink(chainID, tx.Hash()))
	case "coordinator-register-key":
		coordinatorRegisterKey := flag.NewFlagSet("coordinator-register-key", flag.ExitOnError)
		registerKeyAddress := coordinatorRegisterKey.String("address", "", "coordinator address")
		registerKeyUncompressedPubKey := coordinatorRegisterKey.String("pubkey", "", "uncompressed pubkey")
		registerKeyOracleAddress := coordinatorRegisterKey.String("oracle-address", "", "oracle address")
		helpers.ParseArgs(coordinatorRegisterKey, os.Args[2:], "address", "pubkey", "oracle-address")
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*registerKeyAddress), ec)
		helpers.PanicErr(err)

		// Put key in ECDSA format
		if strings.HasPrefix(*registerKeyUncompressedPubKey, "0x") {
			*registerKeyUncompressedPubKey = strings.Replace(*registerKeyUncompressedPubKey, "0x", "04", 1)
		}
		pubBytes, err := hex.DecodeString(*registerKeyUncompressedPubKey)
		helpers.PanicErr(err)
		pk, err := crypto.UnmarshalPubkey(pubBytes)
		helpers.PanicErr(err)
		tx, err := coordinator.RegisterProvingKey(owner,
			common.HexToAddress(*registerKeyOracleAddress),
			[2]*big.Int{pk.X, pk.Y})
		helpers.PanicErr(err)
		fmt.Println("TX", helpers.ExplorerLink(chainID, tx.Hash()))
	case "coordinator-deregister-key":
		coordinatorDeregisterKey := flag.NewFlagSet("coordinator-deregister-key", flag.ExitOnError)
		deregisterKeyAddress := coordinatorDeregisterKey.String("address", "", "coordinator address")
		deregisterKeyUncompressedPubKey := coordinatorDeregisterKey.String("pubkey", "", "uncompressed pubkey")
		helpers.ParseArgs(coordinatorDeregisterKey, os.Args[2:], "address", "pubkey")
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*deregisterKeyAddress), ec)
		helpers.PanicErr(err)

		// Put key in ECDSA format
		if strings.HasPrefix(*deregisterKeyUncompressedPubKey, "0x") {
			*deregisterKeyUncompressedPubKey = strings.Replace(*deregisterKeyUncompressedPubKey, "0x", "04", 1)
		}
		pubBytes, err := hex.DecodeString(*deregisterKeyUncompressedPubKey)
		helpers.PanicErr(err)
		pk, err := crypto.UnmarshalPubkey(pubBytes)
		helpers.PanicErr(err)
		tx, err := coordinator.DeregisterProvingKey(owner, [2]*big.Int{pk.X, pk.Y})
		helpers.PanicErr(err)
		fmt.Println("TX", helpers.ExplorerLink(chainID, tx.Hash()))
	case "coordinator-subscription":
		coordinatorSub := flag.NewFlagSet("coordinator-subscription", flag.ExitOnError)
		address := coordinatorSub.String("address", "", "coordinator address")
		subID := coordinatorSub.Int64("sub-id", 0, "sub-id")
		helpers.ParseArgs(coordinatorSub, os.Args[2:], "address", "sub-id")
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*address), ec)
		helpers.PanicErr(err)
		fmt.Println("sub-id", *subID, "address", *address, coordinator.Address())
		s, err := coordinator.GetSubscription(nil, uint64(*subID))
		helpers.PanicErr(err)
		fmt.Printf("Subscription %+v\n", s)
	case "consumer-deploy":
		consumerDeployCmd := flag.NewFlagSet("consumer-deploy", flag.ExitOnError)
		consumerCoordinator := consumerDeployCmd.String("coordinator-address", "", "coordinator address")
		keyHash := consumerDeployCmd.String("key-hash", "", "key hash")
		consumerLinkAddress := consumerDeployCmd.String("link-address", "", "link-address")
		// TODO: add other params
		helpers.ParseArgs(consumerDeployCmd, os.Args[2:], "coordinator-address", "key-hash", "link-address")
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
		helpers.PanicErr(err)
		fmt.Println("Consumer address", consumerAddress, "TX", helpers.ExplorerLink(chainID, tx.Hash()))
	case "consumer-subscribe":
		consumerSubscribeCmd := flag.NewFlagSet("consumer-subscribe", flag.ExitOnError)
		consumerSubscribeAddress := consumerSubscribeCmd.String("address", "", "consumer address")
		helpers.ParseArgs(consumerSubscribeCmd, os.Args[2:], "address")
		consumer, err := vrf_single_consumer_example.NewVRFSingleConsumerExample(common.HexToAddress(*consumerSubscribeAddress), ec)
		helpers.PanicErr(err)
		tx, err := consumer.Subscribe(owner)
		helpers.PanicErr(err)
		fmt.Println("hash", tx.Hash())
	case "link-balance":
		linkBalanceCmd := flag.NewFlagSet("link-balance", flag.ExitOnError)
		linkAddress := linkBalanceCmd.String("link-address", "", "link-address")
		address := linkBalanceCmd.String("address", "", "address")
		helpers.ParseArgs(linkBalanceCmd, os.Args[2:], "link-address", "address")
		lt, err := link_token_interface.NewLinkToken(common.HexToAddress(*linkAddress), ec)
		helpers.PanicErr(err)
		b, err := lt.BalanceOf(nil, common.HexToAddress(*address))
		helpers.PanicErr(err)
		fmt.Println(b)
	case "consumer-cancel":
		consumerCancelCmd := flag.NewFlagSet("consumer-cancel", flag.ExitOnError)
		consumerCancelAddress := consumerCancelCmd.String("address", "", "consumer address")
		helpers.ParseArgs(consumerCancelCmd, os.Args[2:], "address")
		consumer, err := vrf_single_consumer_example.NewVRFSingleConsumerExample(common.HexToAddress(*consumerCancelAddress), ec)
		helpers.PanicErr(err)
		tx, err := consumer.Unsubscribe(owner, owner.From)
		helpers.PanicErr(err)
		fmt.Println("TX", helpers.ExplorerLink(chainID, tx.Hash()))
	case "consumer-topup":
		// NOTE NEED TO FUND CONSUMER WITH LINK FIRST
		consumerTopupCmd := flag.NewFlagSet("consumer-topup", flag.ExitOnError)
		consumerTopupAmount := consumerTopupCmd.String("amount", "", "amount in juels")
		consumerTopupAddress := consumerTopupCmd.String("address", "", "consumer address")
		helpers.ParseArgs(consumerTopupCmd, os.Args[2:], "amount", "address")
		consumer, err := vrf_single_consumer_example.NewVRFSingleConsumerExample(common.HexToAddress(*consumerTopupAddress), ec)
		helpers.PanicErr(err)
		amount, s := big.NewInt(0).SetString(*consumerTopupAmount, 10)
		if !s {
			panic(fmt.Sprintf("failed to parse top up amount '%s'", *consumerTopupAmount))
		}
		tx, err := consumer.TopUpSubscription(owner, amount)
		helpers.PanicErr(err)
		fmt.Println("TX", helpers.ExplorerLink(chainID, tx.Hash()))
	case "consumer-request":
		consumerRequestCmd := flag.NewFlagSet("consumer-request", flag.ExitOnError)
		consumerRequestAddress := consumerRequestCmd.String("address", "", "consumer address")
		helpers.ParseArgs(consumerRequestCmd, os.Args[2:], "address")
		consumer, err := vrf_single_consumer_example.NewVRFSingleConsumerExample(common.HexToAddress(*consumerRequestAddress), ec)
		helpers.PanicErr(err)
		tx, err := consumer.RequestRandomWords(owner)
		helpers.PanicErr(err)
		fmt.Println("TX", helpers.ExplorerLink(chainID, tx.Hash()))
	case "consumer-fund-and-request":
		consumerRequestCmd := flag.NewFlagSet("consumer-request", flag.ExitOnError)
		consumerRequestAddress := consumerRequestCmd.String("address", "", "consumer address")
		helpers.ParseArgs(consumerRequestCmd, os.Args[2:], "address")
		consumer, err := vrf_single_consumer_example.NewVRFSingleConsumerExample(common.HexToAddress(*consumerRequestAddress), ec)
		helpers.PanicErr(err)
		// Fund and request 3 link
		tx, err := consumer.FundAndRequestRandomWords(owner, big.NewInt(3000000000000000000))
		helpers.PanicErr(err)
		fmt.Println("TX", helpers.ExplorerLink(chainID, tx.Hash()))
	case "consumer-print":
		consumerPrint := flag.NewFlagSet("consumer-print", flag.ExitOnError)
		address := consumerPrint.String("address", "", "consumer address")
		helpers.ParseArgs(consumerPrint, os.Args[2:], "address")
		consumer, err := vrf_single_consumer_example.NewVRFSingleConsumerExample(common.HexToAddress(*address), ec)
		helpers.PanicErr(err)
		rc, err := consumer.SRequestConfig(nil)
		helpers.PanicErr(err)
		rw, err := consumer.SRandomWords(nil, big.NewInt(0))
		if err != nil {
			fmt.Println("no words")
		}
		rid, err := consumer.SRequestId(nil)
		helpers.PanicErr(err)
		fmt.Printf("Request config %+v Rw %+v Rid %+v\n", rc, rw, rid)
	case "eoa-consumer-deploy":
		consumerDeployCmd := flag.NewFlagSet("eoa-consumer-deploy", flag.ExitOnError)
		consumerCoordinator := consumerDeployCmd.String("coordinator-address", "", "coordinator address")
		consumerLinkAddress := consumerDeployCmd.String("link-address", "", "link-address")
		helpers.ParseArgs(consumerDeployCmd, os.Args[2:], "coordinator-address", "link-address")
		consumerAddress, tx, _, err := vrf_external_sub_owner_example.DeployVRFExternalSubOwnerExample(
			owner,
			ec,
			common.HexToAddress(*consumerCoordinator),
			common.HexToAddress(*consumerLinkAddress))
		helpers.PanicErr(err)
		fmt.Println("Consumer address", consumerAddress, "TX", helpers.ExplorerLink(chainID, tx.Hash()))
	case "eoa-create-sub":
		createSubCmd := flag.NewFlagSet("eoa-create-sub", flag.ExitOnError)
		coordinatorAddress := createSubCmd.String("coordinator-address", "", "coordinator address")
		helpers.ParseArgs(createSubCmd, os.Args[2:], "coordinator-address")
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddress), ec)
		helpers.PanicErr(err)
		tx, err := coordinator.CreateSubscription(owner)
		helpers.PanicErr(err)
		fmt.Println("Create subscription", "TX", helpers.ExplorerLink(chainID, tx.Hash()))
	case "eoa-add-sub-consumer":
		addSubConsCmd := flag.NewFlagSet("eoa-add-sub-consumer", flag.ExitOnError)
		coordinatorAddress := addSubConsCmd.String("coordinator-address", "", "coordinator address")
		subID := addSubConsCmd.Uint64("sub-id", 0, "sub-id")
		consumerAddress := addSubConsCmd.String("consumer-address", "", "consumer address")
		helpers.ParseArgs(addSubConsCmd, os.Args[2:], "coordinator-address", "sub-id", "consumer-address")
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddress), ec)
		helpers.PanicErr(err)
		txadd, err := coordinator.AddConsumer(owner, *subID, common.HexToAddress(*consumerAddress))
		helpers.PanicErr(err)
		fmt.Println("Adding consumer", "TX hash", txadd.Hash())
	case "eoa-create-fund-authorize-sub":
		// Lets just treat the owner key as the EOA controlling the sub
		cfaSubCmd := flag.NewFlagSet("eoa-create-fund-authorize-sub", flag.ExitOnError)
		coordinatorAddress := cfaSubCmd.String("coordinator-address", "", "coordinator address")
		amountStr := cfaSubCmd.String("amount", "", "amount to fund in juels")
		consumerAddress := cfaSubCmd.String("consumer-address", "", "consumer address")
		consumerLinkAddress := cfaSubCmd.String("link-address", "", "link-address")
		helpers.ParseArgs(cfaSubCmd, os.Args[2:], "coordinator-address", "amount", "consumer-address", "link-address")
		amount, s := big.NewInt(0).SetString(*amountStr, 10)
		if !s {
			panic(fmt.Sprintf("failed to parse top up amount '%s'", *amountStr))
		}
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddress), ec)
		helpers.PanicErr(err)
		fmt.Println(amount, consumerLinkAddress)
		txcreate, err := coordinator.CreateSubscription(owner)
		helpers.PanicErr(err)
		fmt.Println("Create sub", "TX", helpers.ExplorerLink(chainID, txcreate.Hash()))
		sub := make(chan *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCreated)
		subscription, err := coordinator.WatchSubscriptionCreated(nil, sub, nil)
		helpers.PanicErr(err)
		defer subscription.Unsubscribe()
		created := <-sub
		linkToken, err := link_token_interface.NewLinkToken(common.HexToAddress(*consumerLinkAddress), ec)
		helpers.PanicErr(err)
		bal, err := linkToken.BalanceOf(nil, owner.From)
		helpers.PanicErr(err)
		fmt.Println("OWNER BALANCE", bal, owner.From.String(), amount.String())
		b, err := utils.GenericEncode([]string{"uint64"}, created.SubId)
		helpers.PanicErr(err)
		owner.GasLimit = 500000
		tx, err := linkToken.TransferAndCall(owner, coordinator.Address(), amount, b)
		helpers.PanicErr(err)
		fmt.Println("Funding sub", created.SubId, "TX", helpers.ExplorerLink(chainID, tx.Hash()))
		subFunded := make(chan *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionFunded)
		fundSub, err := coordinator.WatchSubscriptionFunded(nil, subFunded, []uint64{created.SubId})
		helpers.PanicErr(err)
		defer fundSub.Unsubscribe()
		<-subFunded // Add a consumer once its funded
		txadd, err := coordinator.AddConsumer(owner, created.SubId, common.HexToAddress(*consumerAddress))
		helpers.PanicErr(err)
		fmt.Println("adding consumer", "TX", helpers.ExplorerLink(chainID, txadd.Hash()))
	case "eoa-request":
		request := flag.NewFlagSet("eoa-request", flag.ExitOnError)
		consumerAddress := request.String("consumer-address", "", "consumer address")
		subID := request.Uint64("sub-id", 0, "subscription ID")
		cbGasLimit := request.Uint("cb-gas-limit", 1_000_000, "callback gas limit")
		requestConfirmations := request.Uint("request-confirmations", 3, "minimum request confirmations")
		numWords := request.Uint("num-words", 3, "number of words to request")
		keyHash := request.String("key-hash", "", "key hash")
		helpers.ParseArgs(request, os.Args[2:], "consumer-address", "sub-id", "key-hash")
		keyHashBytes := common.HexToHash(*keyHash)
		consumer, err := vrf_external_sub_owner_example.NewVRFExternalSubOwnerExample(
			common.HexToAddress(*consumerAddress),
			ec)
		helpers.PanicErr(err)
		tx, err := consumer.RequestRandomWords(owner, *subID, uint32(*cbGasLimit), uint16(*requestConfirmations), uint32(*numWords), keyHashBytes)
		helpers.PanicErr(err)
		fmt.Println("TX", helpers.ExplorerLink(chainID, tx.Hash()))
	case "eoa-transfer-sub":
		trans := flag.NewFlagSet("eoa-transfer-sub", flag.ExitOnError)
		coordinatorAddress := trans.String("coordinator-address", "", "coordinator address")
		subID := trans.Int64("sub-id", 0, "sub-id")
		to := trans.String("to", "", "to")
		helpers.ParseArgs(trans, os.Args[2:], "coordinator-address", "sub-id", "to")
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddress), ec)
		helpers.PanicErr(err)
		tx, err := coordinator.RequestSubscriptionOwnerTransfer(owner, uint64(*subID), common.HexToAddress(*to))
		helpers.PanicErr(err)
		fmt.Println("ownership transfer requested TX", helpers.ExplorerLink(chainID, tx.Hash()))
	case "eoa-accept-sub":
		accept := flag.NewFlagSet("eoa-accept-sub", flag.ExitOnError)
		coordinatorAddress := accept.String("coordinator-address", "", "coordinator address")
		subID := accept.Int64("sub-id", 0, "sub-id")
		helpers.ParseArgs(accept, os.Args[2:], "coordinator-address", "sub-id")
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddress), ec)
		helpers.PanicErr(err)
		tx, err := coordinator.AcceptSubscriptionOwnerTransfer(owner, uint64(*subID))
		helpers.PanicErr(err)
		fmt.Println("ownership transfer accepted TX", helpers.ExplorerLink(chainID, tx.Hash()))
	case "eoa-cancel-sub":
		cancel := flag.NewFlagSet("eoa-cancel-sub", flag.ExitOnError)
		coordinatorAddress := cancel.String("coordinator-address", "", "coordinator address")
		subID := cancel.Int64("sub-id", 0, "sub-id")
		helpers.ParseArgs(cancel, os.Args[2:], "coordinator-address", "sub-id")
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddress), ec)
		helpers.PanicErr(err)
		tx, err := coordinator.CancelSubscription(owner, uint64(*subID), owner.From)
		helpers.PanicErr(err)
		fmt.Println("sub cancelled TX", helpers.ExplorerLink(chainID, tx.Hash()))
	case "eoa-fund-sub":
		fund := flag.NewFlagSet("eoa-fund-sub", flag.ExitOnError)
		coordinatorAddress := fund.String("coordinator-address", "", "coordinator address")
		amountStr := fund.String("amount", "", "amount to fund in juels")
		subID := fund.Int64("sub-id", 0, "sub-id")
		consumerLinkAddress := fund.String("link-address", "", "link-address")
		helpers.ParseArgs(fund, os.Args[2:], "coordinator-address", "amount", "sub-id", "link-address")
		amount, s := big.NewInt(0).SetString(*amountStr, 10)
		if !s {
			panic(fmt.Sprintf("failed to parse top up amount '%s'", *amountStr))
		}
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddress), ec)
		helpers.PanicErr(err)
		linkToken, err := link_token_interface.NewLinkToken(common.HexToAddress(*consumerLinkAddress), ec)
		helpers.PanicErr(err)
		bal, err := linkToken.BalanceOf(nil, owner.From)
		helpers.PanicErr(err)
		fmt.Println("Initial account balance:", bal, owner.From.String(), "Funding amount:", amount.String())
		b, err := utils.GenericEncode([]string{"uint64"}, uint64(*subID))
		helpers.PanicErr(err)
		owner.GasLimit = 500000
		tx, err := linkToken.TransferAndCall(owner, coordinator.Address(), amount, b)
		helpers.PanicErr(err)
		fmt.Println("Funding sub", *subID, "TX", helpers.ExplorerLink(chainID, tx.Hash()))
		helpers.PanicErr(err)
	case "owner-cancel-sub":
		cancel := flag.NewFlagSet("owner-cancel-sub", flag.ExitOnError)
		coordinatorAddress := cancel.String("coordinator-address", "", "coordinator address")
		subID := cancel.Int64("sub-id", 0, "sub-id")
		helpers.ParseArgs(cancel, os.Args[2:], "coordinator-address", "sub-id")
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddress), ec)
		helpers.PanicErr(err)
		tx, err := coordinator.OwnerCancelSubscription(owner, uint64(*subID))
		helpers.PanicErr(err)
		fmt.Println("sub cancelled TX", helpers.ExplorerLink(chainID, tx.Hash()))
	case "sub-balance":
		consumerBalanceCmd := flag.NewFlagSet("sub-balance", flag.ExitOnError)
		coordinatorAddress := consumerBalanceCmd.String("coordinator-address", "", "coordinator address")
		subID := consumerBalanceCmd.Uint64("sub-id", 0, "subscription id")
		helpers.ParseArgs(consumerBalanceCmd, os.Args[2:], "coordinator-address", "sub-id")
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddress), ec)
		helpers.PanicErr(err)
		resp, err := coordinator.GetSubscription(nil, *subID)
		helpers.PanicErr(err)
		fmt.Println("sub id", *subID, "balance:", resp.Balance)
	default:
		panic("unrecognized subcommand: " + os.Args[1])
	}
}
