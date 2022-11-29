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
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/vrf_coordinator_v2"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
)

const formattedVRFJob = `
type = "vrf"
name = "vrf_v2"
schemaVersion = 1
coordinatorAddress = "%s"
batchCoordinatorAddress = "%s"
batchFulfillmentEnabled = %t
batchFulfillmentGasMultiplier = 1.1
publicKey = "%s"
minIncomingConfirmations = 3
evmChainID = "%d"
fromAddresses = ["%s"]
pollPeriod = "5s"
requestTimeout = "24h"
observationSource = """decode_log   [type=ethabidecodelog
              abi="RandomWordsRequested(bytes32 indexed keyHash,uint256 requestId,uint256 preSeed,uint64 indexed subId,uint16 minimumRequestConfirmations,uint32 callbackGasLimit,uint32 numWords,address indexed sender)"
              data="$(jobRun.logData)"
              topics="$(jobRun.logTopics)"]
vrf          [type=vrfv2
              publicKey="$(jobSpec.publicKey)"
              requestBlockHash="$(jobRun.logBlockHash)"
              requestBlockNumber="$(jobRun.logBlockNumber)"
              topics="$(jobRun.logTopics)"]
estimate_gas [type=estimategaslimit
              to="%s"
              multiplier="1.1"
              data="$(vrf.output)"]
simulate     [type=ethcall
              to="%s"
              gas="$(estimate_gas)"
              gasPrice="$(jobSpec.maxGasPrice)"
              extractRevertReason=true
              contract="%s"
              data="$(vrf.output)"]
decode_log->vrf->estimate_gas->simulate
""" 
`

func deployUniverse(e helpers.Environment) {
	deployCmd := flag.NewFlagSet("deploy-universe", flag.ExitOnError)

	// required flags
	linkAddress := deployCmd.String("link-address", "", "address of link token")
	linkEthAddress := deployCmd.String("link-eth-feed", "", "address of link eth feed")
	subscriptionBalanceString := deployCmd.String("subscription-balance", "1e19", "amount to fund subscription")

	// optional flags
	fallbackWeiPerUnitLinkString := deployCmd.String("fallback-wei-per-unit-link", "6e16", "fallback wei/link ratio")
	registerKeyUncompressedPubKey := deployCmd.String("uncompressed-pub-key", "", "uncompressed public key")
	registerKeyOracleAddress := deployCmd.String("oracle-address", "", "oracle sender address")
	minConfs := deployCmd.Int("min-confs", 3, "min confs")
	oracleFundingAmount := deployCmd.Int64("oracle-funding-amount", assets.GWei(100_000_000).Int64(), "amount to fund sending oracle")
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
	noCancelEnabled := deployCmd.Bool("no-cancel-enabled", false, "whether to deploy the no-cancel coordinator")

	helpers.ParseArgs(
		deployCmd, os.Args[2:],
	)

	fallbackWeiPerUnitLink := decimal.RequireFromString(*fallbackWeiPerUnitLinkString).BigInt()
	subscriptionBalance := decimal.RequireFromString(*subscriptionBalanceString).BigInt()

	// Put key in ECDSA format
	if strings.HasPrefix(*registerKeyUncompressedPubKey, "0x") {
		*registerKeyUncompressedPubKey = strings.Replace(*registerKeyUncompressedPubKey, "0x", "04", 1)
	}

	// Generate compressed public key and key hash
	pubBytes, err := hex.DecodeString(*registerKeyUncompressedPubKey)
	helpers.PanicErr(err)
	pk, err := crypto.UnmarshalPubkey(pubBytes)
	helpers.PanicErr(err)
	var pkBytes []byte
	if big.NewInt(0).Mod(pk.Y, big.NewInt(2)).Uint64() != 0 {
		pkBytes = append(pk.X.Bytes(), 1)
	} else {
		pkBytes = append(pk.X.Bytes(), 0)
	}
	var newPK secp256k1.PublicKey
	copy(newPK[:], pkBytes)

	compressedPkHex := hexutil.Encode(pkBytes)
	keyHash, err := newPK.Hash()
	helpers.PanicErr(err)

	if len(*linkAddress) == 0 {
		fmt.Println("\nDeploying LINK Token...")
		address := helpers.DeployLinkToken(e).String()
		linkAddress = &address
	}

	if len(*linkEthAddress) == 0 {
		fmt.Println("\nDeploying LINK/ETH Feed...")
		address := helpers.DeployLinkEthFeed(e, *linkAddress, fallbackWeiPerUnitLink).String()
		linkEthAddress = &address
	}

	fmt.Println("\nDeploying BHS...")
	bhsContractAddress := deployBHS(e)

	fmt.Println("\nDeploying Batch BHS...")
	batchBHSAddress := deployBatchBHS(e, bhsContractAddress)

	var coordinatorAddress common.Address
	if *noCancelEnabled {
		fmt.Println("\nDeploying NoCancelCoordinator...")
		coordinatorAddress = deployNoCancelCoordinator(e, *linkAddress, bhsContractAddress.String(), *linkEthAddress)
	} else {
		fmt.Println("\nDeploying Coordinator...")
		coordinatorAddress = deployCoordinator(e, *linkAddress, bhsContractAddress.String(), *linkEthAddress)
	}

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
		fallbackWeiPerUnitLink,
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

	if len(*registerKeyOracleAddress) > 0 {
		fmt.Println("\nFunding oracle...")
		helpers.FundNodes(e, []string{*registerKeyOracleAddress}, big.NewInt(*oracleFundingAmount))
	}

	formattedJobSpec := fmt.Sprintf(
		formattedVRFJob,
		coordinatorAddress,
		batchCoordinatorAddress,
		false,
		compressedPkHex,
		e.ChainID,
		*registerKeyOracleAddress,
		coordinatorAddress,
		coordinatorAddress,
		coordinatorAddress,
	)

	fmt.Println(
		"\nDeployment complete.",
		"\nLINK Token contract address:", *linkAddress,
		"\nLINK/ETH Feed contract address:", *linkEthAddress,
		"\nBlockhash Store contract address:", bhsContractAddress,
		"\nBatch Blockhash Store contract address:", batchBHSAddress,
		"\nVRF Coordinator Address:", coordinatorAddress,
		"\nBatch VRF Coordinator Address:", batchCoordinatorAddress,
		"\nVRF Consumer Address:", consumerAddress,
		"\nVRF Subscription Id:", subID,
		"\nVRF Subscription Balance:", *subscriptionBalanceString,
		"\nPossible VRF Request command: ",
		fmt.Sprintf("go run . eoa-request --consumer-address %s --sub-id %d --key-hash %s", consumerAddress, subID, keyHash),
		"\nA node can now be configured to run a VRF job with the below job spec :\n",
		formattedJobSpec,
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
