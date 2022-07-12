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
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

func deployUniverse(e helpers.Environment) {
	deployCmd := flag.NewFlagSet("deploy-universe", flag.ExitOnError)

	// required flags
	linkAddress := deployCmd.String("link-address", "", "address of link token")
	linkEthAddress := deployCmd.String("link-eth-feed", "", "address of link eth feed")
	subscriptionBalanceString := deployCmd.String("subscription-balance", "", "amount to fund subscription")

	// optional flags
	fallbackWeiPerUnitLink := deployCmd.String("fallback-wei-per-unit-link", "60000000000000000", "fallback wei/link ratio")
	registerKeyUncompressedPubKey := deployCmd.String("uncompressed-pub-key", "", "uncompressed public key")
	registerKeyOracleAddress := deployCmd.String("oracle-address", "", "oracle address")
	minConfs := deployCmd.Int("min-confs", 3, "min confs")
	maxGasLimit := deployCmd.Int64("max-gas-limit", 2.5e6, "max gas limit")
	stalenessSeconds := deployCmd.Int64("staleness-seconds", 86400, "staleness in seconds")
	gasAfterPayment := deployCmd.Int64("gas-after-payment", 33285, "gas after payment calculation")
	flatFeeTier1 := deployCmd.Int64("flat-fee-tier-1", 500, "flat fee tier 1")
	flatFeeTier2 := deployCmd.Int64("flat-fee-tier-2", 500, "flat fee tier 2")
	flatFeeTier3 := deployCmd.Int64("flat-fee-tier-3", 500, "flat fee tier 3")
	flatFeeTier4 := deployCmd.Int64("flat-fee-tier-4", 500, "flat fee tier 4")
	flatFeeTier5 := deployCmd.Int64("flat-fee-tier-5", 500, "flat fee tier 5")
	reqsForTier2 := deployCmd.Int64("reqs-for-tier-2", 0, "requests for tier 2")
	reqsForTier3 := deployCmd.Int64("reqs-for-tier-3", 0, "requests for tier 3")
	reqsForTier4 := deployCmd.Int64("reqs-for-tier-4", 0, "requests for tier 4")
	reqsForTier5 := deployCmd.Int64("reqs-for-tier-5", 0, "requests for tier 5")

	helpers.ParseArgs(
		deployCmd, os.Args[2:],
		"link-address",
		"link-eth-feed",
		"subscription-balance",
	)

	subscriptionBalance := decimal.RequireFromString(*subscriptionBalanceString).BigInt()

	// Put key in ECDSA format
	if strings.HasPrefix(*registerKeyUncompressedPubKey, "0x") {
		*registerKeyUncompressedPubKey = strings.Replace(*registerKeyUncompressedPubKey, "0x", "04", 1)
	}

	fmt.Println("\nDeploying BHS...")
	bhsContractAddress := deployBHS(e)

	fmt.Println("\nDeploying Batch BHS...")
	batchBHSAddress := deployBatchBHS(e, bhsContractAddress)

	fmt.Println("\nDeploying Coordinator...")
	coordinatorAddress := deployCoordinator(e, *linkAddress, bhsContractAddress.String(), *linkEthAddress)
	coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(coordinatorAddress, e.Ec)
	helpers.PanicErr(err)

	fmt.Println("\nDeploying Batch Coordinator...")
	batchCoordinatorAddress := deployBatchCoordinatorV2(e, coordinatorAddress)

	fmt.Println("\nSetting Coordinator Config...")
	setCoordinatorConfig(
		e,
		*coordinator,
		uint16(*minConfs),
		uint32(*maxGasLimit),
		uint32(*stalenessSeconds),
		uint32(*gasAfterPayment),
		decimal.RequireFromString(*fallbackWeiPerUnitLink).BigInt(),
		vrf_coordinator_v2.VRFCoordinatorV2FeeConfig{
			FulfillmentFlatFeeLinkPPMTier1: uint32(*flatFeeTier1),
			FulfillmentFlatFeeLinkPPMTier2: uint32(*flatFeeTier2),
			FulfillmentFlatFeeLinkPPMTier3: uint32(*flatFeeTier3),
			FulfillmentFlatFeeLinkPPMTier4: uint32(*flatFeeTier4),
			FulfillmentFlatFeeLinkPPMTier5: uint32(*flatFeeTier5),
			ReqsForTier2:                   big.NewInt(*reqsForTier2),
			ReqsForTier3:                   big.NewInt(*reqsForTier3),
			ReqsForTier4:                   big.NewInt(*reqsForTier4),
			ReqsForTier5:                   big.NewInt(*reqsForTier5),
		},
	)

	fmt.Println("\nConfig set, getting current config from deployed contract...")
	printCoordinatorConfig(coordinator)

	if len(*registerKeyUncompressedPubKey) > 0 && len(*registerKeyOracleAddress) > 0 {
		fmt.Println("\nRegistering proving key...")
		registerCoordinatorProvingKey(e, *coordinator, *registerKeyUncompressedPubKey, *registerKeyOracleAddress)

		fmt.Println("\nProving key registered, getting proving key hashes from deployed contract...")
		_, _, provingKeyHashes, configErr := coordinator.GetRequestConfig(nil)
		helpers.PanicErr(configErr)
		fmt.Println("Key hash registered:", hex.EncodeToString(provingKeyHashes[0][:]))
	} else {
		fmt.Println("NOT registering proving key - you must do this eventually in order to fully deploy VRF!")
	}

	fmt.Println("\nDeploying consumer...")
	consumerAddress := eoaDeployConsumer(e, coordinatorAddress.String(), *linkAddress)

	fmt.Println("\nAdding subscription...")
	eoaCreateSub(e, *coordinator)
	subID := uint64(1)

	fmt.Println("\nAdding consumer to subscription...")
	eoaAddConsumerToSub(e, *coordinator, subID, consumerAddress.String())

	if subscriptionBalance.Cmp(big.NewInt(0)) > 0 {
		fmt.Println("\nFunding subscription with", subscriptionBalance, "juels...")
		eoaFundSubscription(e, *coordinator, *linkAddress, subscriptionBalance, subID)
	} else {
		fmt.Println("Subscription", subID, "NOT getting funded. You must fund the subscription in order to use it!")
	}

	fmt.Println("\nSubscribed and (possibly) funded, retrieving subscription from deployed contract...")
	s, err := coordinator.GetSubscription(nil, subID)
	helpers.PanicErr(err)
	fmt.Printf("Subscription %+v\n", s)
	fmt.Println(
		"\nDeployment complete.",
		"\nBlockhash Store contract address:", bhsContractAddress,
		"\nBatch Blockhash Store contract address:", batchBHSAddress,
		"\nVRF Coordinator Address:", coordinatorAddress,
		"\nBatch VRF Coordinator Address:", batchCoordinatorAddress,
		"\nVRF Consumer Address:", consumerAddress,
		"\nVRF Subscription Id:", subID,
		"\nVRF Subscription Balance:", *subscriptionBalanceString,
		"\nA node can now be configured to run a VRF job with the above configuration.",
	)
}

func deployWrapperUniverse(e helpers.Environment) {
	cmd := flag.NewFlagSet("wrapper-universe-deploy", flag.ExitOnError)
	linkAddress := cmd.String("link-address", "", "address of link token")
	linkETHFeedAddress := cmd.String("link-eth-feed", "", "address of link-eth-feed")
	coordinatorAddress := cmd.String("coordinator-address", "", "address of the vrf coordinator v2 contract")
	wrapperGasOverhead := cmd.Uint("wrapper-gas-overhead", 50_000, "amount of gas overhead in wrapper fulfillment")
	coordinatorGasOverhead := cmd.Uint("coordinator-gas-overhead", 52_000, "amount of gas overhead in coordinator fulfillment")
	wrapperPremiumPercentage := cmd.Uint("wrapper-premium-percentage", 25, "gas premium charged by wrapper")
	keyHash := cmd.String("key-hash", "", "the keyhash that wrapper requests should use")
	maxNumWords := cmd.Uint("max-num-words", 10, "the keyhash that wrapper requests should use")
	subFunding := cmd.String("sub-funding", "10000000000000000000", "amount to fund the subscription with")
	consumerFunding := cmd.String("consumer-funding", "10000000000000000000", "amount to fund the consumer with")
	helpers.ParseArgs(cmd, os.Args[2:], "link-address", "link-eth-feed", "coordinator-address", "key-hash")

	amount, s := big.NewInt(0).SetString(*subFunding, 10)
	if !s {
		panic(fmt.Sprintf("failed to parse top up amount '%s'", *subFunding))
	}

	wrapper, subID := wrapperDeploy(e,
		common.HexToAddress(*linkAddress),
		common.HexToAddress(*linkETHFeedAddress),
		common.HexToAddress(*coordinatorAddress))

	wrapperConfigure(e,
		wrapper,
		*wrapperGasOverhead,
		*coordinatorGasOverhead,
		*wrapperPremiumPercentage,
		*keyHash,
		*maxNumWords)

	consumer := wrapperConsumerDeploy(e,
		common.HexToAddress(*linkAddress),
		wrapper)

	coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddress), e.Ec)
	helpers.PanicErr(err)

	eoaFundSubscription(e, *coordinator, *linkAddress, amount, subID)

	link, err := link_token_interface.NewLinkToken(common.HexToAddress(*linkAddress), e.Ec)
	helpers.PanicErr(err)
	consumerAmount, s := big.NewInt(0).SetString(*consumerFunding, 10)
	if !s {
		panic(fmt.Sprintf("failed to parse top up amount '%s'", *consumerFunding))
	}

	tx, err := link.Transfer(e.Owner, consumer, consumerAmount)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)

}
