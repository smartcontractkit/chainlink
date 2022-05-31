package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

func deployBHSCoordinatorAndConsumer(e environment) {
	coordinatorDeployCmd := flag.NewFlagSet("full-deploy", flag.ExitOnError)
	linkAddress := coordinatorDeployCmd.String("link-address", "", "address of link token")
	linkEthAddress := coordinatorDeployCmd.String("link-eth-feed", "", "address of link eth feed")
	fallbackWeiPerUnitLink := coordinatorDeployCmd.String("fallback-wei-per-unit-link", "", "fallback wei/link ratio")
	registerKeyUncompressedPubKey := coordinatorDeployCmd.String("uncompressed-pub-key", "", "uncompressed public key")
	registerKeyOracleAddress := coordinatorDeployCmd.String("oracle-address", "", "oracle address")
	subscriptionBalanceString := coordinatorDeployCmd.String("subscription-balance", "", "subscription balance")
	helpers.ParseArgs(
		coordinatorDeployCmd, os.Args[2:],
		"link-address",
		"link-eth-feed",
		"fallback-wei-per-unit-link",
		"subscription-balance",
	)

	subscriptionBalance, success := big.NewInt(0).SetString(*subscriptionBalanceString, 10)
	if !success {
		panic(fmt.Sprintf("failed to parse subscriptionBalance '%s'", *subscriptionBalanceString))
	}

	// Put key in ECDSA format
	if strings.HasPrefix(*registerKeyUncompressedPubKey, "0x") {
		*registerKeyUncompressedPubKey = strings.Replace(*registerKeyUncompressedPubKey, "0x", "04", 1)
	}

	fmt.Println("\nDeploying BHS...")
	bhsContractAddress := deployBHS(e)

	fmt.Println("\nDeploying Coordinator...")
	coordinatorAddress := deployCoordinator(e, *linkAddress, bhsContractAddress.String(), *linkEthAddress)
	coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(coordinatorAddress, e.ec)
	helpers.PanicErr(err)

	fmt.Println("\nSetting Config...")
	setCoordinatorConfig(e, *coordinator, 3, 2.5e6, 86400, 33285, decimal.RequireFromString(*fallbackWeiPerUnitLink).BigInt(),
		vrf_coordinator_v2.VRFCoordinatorV2FeeConfig{
			FulfillmentFlatFeeLinkPPMTier1: uint32(500),
			FulfillmentFlatFeeLinkPPMTier2: uint32(500),
			FulfillmentFlatFeeLinkPPMTier3: uint32(500),
			FulfillmentFlatFeeLinkPPMTier4: uint32(500),
			FulfillmentFlatFeeLinkPPMTier5: uint32(500),
			ReqsForTier2:                   big.NewInt(0),
			ReqsForTier3:                   big.NewInt(0),
			ReqsForTier4:                   big.NewInt(0),
			ReqsForTier5:                   big.NewInt(0),
		},
	)

	fmt.Println("\nConfig set, getting current config from deployed contract...")
	printCoordinatorConfig(e, *coordinator)

	if len(*registerKeyUncompressedPubKey) > 0 {
		fmt.Println("\nRegistering proving key...")
		registerCoordinatorProvingKey(e, *coordinator, *registerKeyUncompressedPubKey, *registerKeyOracleAddress)

		fmt.Println("\nProving key registered, getting proving key hashes from deployed contract...")
		_, _, s_provingKeyHashes, err := coordinator.GetRequestConfig(nil)
		helpers.PanicErr(err)
		fmt.Printf("Hashes: %+v\n", s_provingKeyHashes)
	}

	fmt.Println("\nDeploying consumer...")
	consumerAddress := eoaDeployConsumer(e, coordinatorAddress.String(), *linkAddress)

	fmt.Println("\nAdding subscription...")
	eoaCreateSub(e, *coordinator)
	subID := uint64(1)

	fmt.Println("\nAdding consumer to subscription...")
	eoaAddConsumerToSub(e, *coordinator, subID, consumerAddress.String())

	fmt.Println("\nFunding subscription...")
	eoaFundSubscription(e, *coordinator, *linkAddress, subscriptionBalance, subID)

	fmt.Println("\nSubscribed and funded, retrieving subscription from deployed contract...")
	s, err := coordinator.GetSubscription(nil, subID)
	helpers.PanicErr(err)
	fmt.Printf("Subscription %+v\n", s)
	fmt.Println(
		"\nDeployment complete.",
		"\nBlockhash Store contract address:", bhsContractAddress,
		"\nVRF Coordinator Address:", coordinatorAddress,
		"\nVRF Consumer Address:", consumerAddress,
		"\nVRF Subscription Id:", subID,
		"\nVRF Subscription Balance:", *subscriptionBalanceString,
		"\nA node can now be configured to run a VRF job with the above configuration.",
	)
}
