package v2scripts

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/smartcontractkit/chainlink/core/scripts/common/vrf/constants"
	"github.com/smartcontractkit/chainlink/core/scripts/common/vrf/jobs"
	"github.com/smartcontractkit/chainlink/core/scripts/common/vrf/model"
	"github.com/smartcontractkit/chainlink/core/scripts/common/vrf/util"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
)

type CoordinatorConfigV2 struct {
	MinConfs               *int
	MaxGasLimit            *int64
	StalenessSeconds       *int64
	GasAfterPayment        *int64
	FallbackWeiPerUnitLink *big.Int
	FeeConfig              vrf_coordinator_v2.VRFCoordinatorV2FeeConfig
}

func DeployUniverseViaCLI(e helpers.Environment) {
	deployCmd := flag.NewFlagSet("deploy-universe", flag.ExitOnError)

	// required flags
	linkAddress := *deployCmd.String("link-address", "", "address of link token")
	linkEthAddress := *deployCmd.String("link-eth-feed", "", "address of link eth feed")
	bhsContractAddressString := *deployCmd.String("bhs-address", "", "address of BHS contract")
	batchBHSAddressString := *deployCmd.String("batch-bhs-address", "", "address of Batch BHS contract")
	coordinatorAddressString := *deployCmd.String("coordinator-address", "", "address of VRF Coordinator contract")
	batchCoordinatorAddressString := *deployCmd.String("batch-coordinator-address", "", "address Batch VRF Coordinator contract")

	subscriptionBalanceJuelsString := deployCmd.String("subscription-balance", constants.SubscriptionBalanceJuels, "amount to fund subscription")
	nodeSendingKeyFundingAmount := deployCmd.String("sending-key-funding-amount", constants.NodeSendingKeyFundingAmount, "CL node sending key funding amount")

	batchFulfillmentEnabled := deployCmd.Bool("batch-fulfillment-enabled", constants.BatchFulfillmentEnabled, "whether send randomness fulfillments in batches inside one tx from CL node")

	// optional flags
	fallbackWeiPerUnitLinkString := deployCmd.String("fallback-wei-per-unit-link", constants.FallbackWeiPerUnitLink.String(), "fallback wei/link ratio")
	registerKeyUncompressedPubKey := deployCmd.String("uncompressed-pub-key", "", "uncompressed public key")
	vrfPrimaryNodeSendingKeysString := deployCmd.String("vrf-primary-node-sending-keys", "", "VRF Primary Node sending keys")

	minConfs := deployCmd.Int("min-confs", constants.MinConfs, "min confs")
	maxGasLimit := deployCmd.Int64("max-gas-limit", constants.MaxGasLimit, "max gas limit")
	stalenessSeconds := deployCmd.Int64("staleness-seconds", constants.StalenessSeconds, "staleness in seconds")
	gasAfterPayment := deployCmd.Int64("gas-after-payment", constants.GasAfterPayment, "gas after payment calculation")
	flatFeeTier1 := deployCmd.Int64("flat-fee-tier-1", constants.FlatFeeTier1, "flat fee tier 1")
	flatFeeTier2 := deployCmd.Int64("flat-fee-tier-2", constants.FlatFeeTier2, "flat fee tier 2")
	flatFeeTier3 := deployCmd.Int64("flat-fee-tier-3", constants.FlatFeeTier3, "flat fee tier 3")
	flatFeeTier4 := deployCmd.Int64("flat-fee-tier-4", constants.FlatFeeTier4, "flat fee tier 4")
	flatFeeTier5 := deployCmd.Int64("flat-fee-tier-5", constants.FlatFeeTier5, "flat fee tier 5")
	reqsForTier2 := deployCmd.Int64("reqs-for-tier-2", constants.ReqsForTier2, "requests for tier 2")
	reqsForTier3 := deployCmd.Int64("reqs-for-tier-3", constants.ReqsForTier3, "requests for tier 3")
	reqsForTier4 := deployCmd.Int64("reqs-for-tier-4", constants.ReqsForTier4, "requests for tier 4")
	reqsForTier5 := deployCmd.Int64("reqs-for-tier-5", constants.ReqsForTier5, "requests for tier 5")

	helpers.ParseArgs(
		deployCmd, os.Args[2:],
	)

	fallbackWeiPerUnitLink := decimal.RequireFromString(*fallbackWeiPerUnitLinkString).BigInt()
	subscriptionBalanceJuels := decimal.RequireFromString(*subscriptionBalanceJuelsString).BigInt()

	feeConfig := vrf_coordinator_v2.VRFCoordinatorV2FeeConfig{
		FulfillmentFlatFeeLinkPPMTier1: uint32(*flatFeeTier1),
		FulfillmentFlatFeeLinkPPMTier2: uint32(*flatFeeTier2),
		FulfillmentFlatFeeLinkPPMTier3: uint32(*flatFeeTier3),
		FulfillmentFlatFeeLinkPPMTier4: uint32(*flatFeeTier4),
		FulfillmentFlatFeeLinkPPMTier5: uint32(*flatFeeTier5),
		ReqsForTier2:                   big.NewInt(*reqsForTier2),
		ReqsForTier3:                   big.NewInt(*reqsForTier3),
		ReqsForTier4:                   big.NewInt(*reqsForTier4),
		ReqsForTier5:                   big.NewInt(*reqsForTier5),
	}

	vrfPrimaryNodeSendingKeys := strings.Split(*vrfPrimaryNodeSendingKeysString, ",")

	nodesMap := make(map[string]model.Node)

	fundingAmount, ok := new(big.Int).SetString(*nodeSendingKeyFundingAmount, 10)
	if !ok {
		panic(fmt.Sprintf("failed to parse node sending key funding amount '%s'", *nodeSendingKeyFundingAmount))
	}
	nodesMap[model.VRFPrimaryNodeName] = model.Node{
		SendingKeys:             util.MapToSendingKeyArr(vrfPrimaryNodeSendingKeys),
		SendingKeyFundingAmount: fundingAmount,
	}

	bhsContractAddress := common.HexToAddress(bhsContractAddressString)
	batchBHSAddress := common.HexToAddress(batchBHSAddressString)
	coordinatorAddress := common.HexToAddress(coordinatorAddressString)
	batchCoordinatorAddress := common.HexToAddress(batchCoordinatorAddressString)

	contractAddresses := model.ContractAddresses{
		LinkAddress:             linkAddress,
		LinkEthAddress:          linkEthAddress,
		BhsContractAddress:      bhsContractAddress,
		BatchBHSAddress:         batchBHSAddress,
		CoordinatorAddress:      coordinatorAddress,
		BatchCoordinatorAddress: batchCoordinatorAddress,
	}

	coordinatorConfig := CoordinatorConfigV2{
		MinConfs:               minConfs,
		MaxGasLimit:            maxGasLimit,
		StalenessSeconds:       stalenessSeconds,
		GasAfterPayment:        gasAfterPayment,
		FallbackWeiPerUnitLink: fallbackWeiPerUnitLink,
		FeeConfig:              feeConfig,
	}

	VRFV2DeployUniverse(
		e,
		subscriptionBalanceJuels,
		registerKeyUncompressedPubKey,
		contractAddresses,
		coordinatorConfig,
		*batchFulfillmentEnabled,
		nodesMap,
	)

	vrfPrimaryNode := nodesMap[model.VRFPrimaryNodeName]
	fmt.Println("Funding node's sending keys...")
	for _, sendingKey := range vrfPrimaryNode.SendingKeys {
		helpers.FundNode(e, sendingKey.Address, vrfPrimaryNode.SendingKeyFundingAmount)
	}
}

func VRFV2DeployUniverse(
	e helpers.Environment,
	subscriptionBalanceJuels *big.Int,
	registerKeyUncompressedPubKey *string,
	contractAddresses model.ContractAddresses,
	coordinatorConfig CoordinatorConfigV2,
	batchFulfillmentEnabled bool,
	nodesMap map[string]model.Node,
) model.JobSpecs {

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

	if len(contractAddresses.LinkAddress) == 0 {
		fmt.Println("\nDeploying LINK Token...")
		contractAddresses.LinkAddress = helpers.DeployLinkToken(e).String()
	}

	if len(contractAddresses.LinkEthAddress) == 0 {
		fmt.Println("\nDeploying LINK/ETH Feed...")
		contractAddresses.LinkEthAddress = helpers.DeployLinkEthFeed(e, contractAddresses.LinkAddress, coordinatorConfig.FallbackWeiPerUnitLink).String()
	}

	if contractAddresses.BhsContractAddress.String() == "0x0000000000000000000000000000000000000000" {
		fmt.Println("\nDeploying BHS...")
		contractAddresses.BhsContractAddress = DeployBHS(e)
	}

	if contractAddresses.BatchBHSAddress.String() == "0x0000000000000000000000000000000000000000" {
		fmt.Println("\nDeploying Batch BHS...")
		contractAddresses.BatchBHSAddress = DeployBatchBHS(e, contractAddresses.BhsContractAddress)
	}

	if contractAddresses.CoordinatorAddress.String() == "0x0000000000000000000000000000000000000000" {
		fmt.Println("\nDeploying Coordinator...")
		contractAddresses.CoordinatorAddress = DeployCoordinator(e, contractAddresses.LinkAddress, contractAddresses.BhsContractAddress.String(), contractAddresses.LinkEthAddress)
	}

	coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(contractAddresses.CoordinatorAddress, e.Ec)
	helpers.PanicErr(err)

	if contractAddresses.BatchCoordinatorAddress.String() == "0x0000000000000000000000000000000000000000" {
		fmt.Println("\nDeploying Batch Coordinator...")
		contractAddresses.BatchCoordinatorAddress = DeployBatchCoordinatorV2(e, contractAddresses.CoordinatorAddress)
	}

	fmt.Println("\nSetting Coordinator Config...")
	SetCoordinatorConfig(
		e,
		*coordinator,
		uint16(*coordinatorConfig.MinConfs),
		uint32(*coordinatorConfig.MaxGasLimit),
		uint32(*coordinatorConfig.StalenessSeconds),
		uint32(*coordinatorConfig.GasAfterPayment),
		coordinatorConfig.FallbackWeiPerUnitLink,
		coordinatorConfig.FeeConfig,
	)

	fmt.Println("\nConfig set, getting current config from deployed contract...")
	PrintCoordinatorConfig(coordinator)

	if len(*registerKeyUncompressedPubKey) > 0 {
		fmt.Println("\nRegistering proving key...")

		//NOTE - register proving key against EOA account, and not against Oracle's sending address in other to be able
		// easily withdraw funds from Coordinator contract back to EOA account
		RegisterCoordinatorProvingKey(e, *coordinator, *registerKeyUncompressedPubKey, e.Owner.From.String())

		fmt.Println("\nProving key registered, getting proving key hashes from deployed contract...")
		_, _, provingKeyHashes, configErr := coordinator.GetRequestConfig(nil)
		helpers.PanicErr(configErr)
		fmt.Println("Key hash registered:", hex.EncodeToString(provingKeyHashes[0][:]))
	} else {
		fmt.Println("NOT registering proving key - you must do this eventually in order to fully deploy VRF!")
	}

	fmt.Println("\nDeploying consumer...")
	consumerAddress := EoaLoadTestConsumerWithMetricsDeploy(e, contractAddresses.CoordinatorAddress.String())

	fmt.Println("\nAdding subscription...")
	EoaCreateSub(e, *coordinator)
	subID := uint64(1)

	fmt.Println("\nAdding consumer to subscription...")
	EoaAddConsumerToSub(e, *coordinator, subID, consumerAddress.String())

	if subscriptionBalanceJuels.Cmp(big.NewInt(0)) > 0 {
		fmt.Println("\nFunding subscription with", subscriptionBalanceJuels, "juels...")
		EoaFundSubscription(e, *coordinator, contractAddresses.LinkAddress, subscriptionBalanceJuels, subID)
	} else {
		fmt.Println("Subscription", subID, "NOT getting funded. You must fund the subscription in order to use it!")
	}

	fmt.Println("\nSubscribed and (possibly) funded, retrieving subscription from deployed contract...")
	s, err := coordinator.GetSubscription(nil, subID)
	helpers.PanicErr(err)
	fmt.Printf("Subscription %+v\n", s)

	formattedVrfPrimaryJobSpec := fmt.Sprintf(
		jobs.VRFV2JobFormatted,
		contractAddresses.CoordinatorAddress,      //coordinatorAddress
		contractAddresses.BatchCoordinatorAddress, //batchCoordinatorAddress
		batchFulfillmentEnabled,                   //batchFulfillmentEnabled
		compressedPkHex,                           //publicKey
		*coordinatorConfig.MinConfs,               //minIncomingConfirmations
		e.ChainID,                                 //evmChainID
		strings.Join(util.MapToAddressArr(nodesMap[model.VRFPrimaryNodeName].SendingKeys), "\",\""), //fromAddresses
		contractAddresses.CoordinatorAddress,
		nodesMap[model.VRFPrimaryNodeName].SendingKeys[0].Address,
		contractAddresses.CoordinatorAddress,
		contractAddresses.CoordinatorAddress,
	)

	formattedVrfBackupJobSpec := fmt.Sprintf(
		jobs.VRFV2JobFormatted,
		contractAddresses.CoordinatorAddress,      //coordinatorAddress
		contractAddresses.BatchCoordinatorAddress, //batchCoordinatorAddress
		batchFulfillmentEnabled,                   //batchFulfillmentEnabled
		compressedPkHex,                           //publicKey
		100,                                       //minIncomingConfirmations
		e.ChainID,                                 //evmChainID
		strings.Join(util.MapToAddressArr(nodesMap[model.VRFBackupNodeName].SendingKeys), "\",\""), //fromAddresses
		contractAddresses.CoordinatorAddress,
		nodesMap[model.VRFPrimaryNodeName].SendingKeys[0],
		contractAddresses.CoordinatorAddress,
		contractAddresses.CoordinatorAddress,
	)

	formattedBHSJobSpec := fmt.Sprintf(
		jobs.BHSJobFormatted,
		contractAddresses.CoordinatorAddress, //coordinatorAddress
		30,                                   //waitBlocks
		200,                                  //lookbackBlocks
		contractAddresses.BhsContractAddress, //bhs address
		e.ChainID,                            //chain id
		strings.Join(util.MapToAddressArr(nodesMap[model.BHSNodeName].SendingKeys), "\",\""), //sending addresses
	)

	formattedBHSBackupJobSpec := fmt.Sprintf(
		jobs.BHSJobFormatted,
		contractAddresses.CoordinatorAddress, //coordinatorAddress
		100,                                  //waitBlocks
		200,                                  //lookbackBlocks
		contractAddresses.BhsContractAddress, //bhs adreess
		e.ChainID,                            //chain id
		strings.Join(util.MapToAddressArr(nodesMap[model.BHSBackupNodeName].SendingKeys), "\",\""), //sending addresses
	)

	formattedBHFJobSpec := fmt.Sprintf(
		jobs.BHFJobFormatted,
		contractAddresses.CoordinatorAddress, //coordinatorAddress
		contractAddresses.BhsContractAddress, //bhs adreess
		contractAddresses.BatchBHSAddress,    //batchBHS
		e.ChainID,                            //chain id
		strings.Join(util.MapToAddressArr(nodesMap[model.BHFNodeName].SendingKeys), "\",\""), //sending addresses
	)

	fmt.Println(
		"\nDeployment complete.",
		"\nLINK Token contract address:", contractAddresses.LinkAddress,
		"\nLINK/ETH Feed contract address:", contractAddresses.LinkEthAddress,
		"\nBlockhash Store contract address:", contractAddresses.BhsContractAddress,
		"\nBatch Blockhash Store contract address:", contractAddresses.BatchBHSAddress,
		"\nVRF Coordinator Address:", contractAddresses.CoordinatorAddress,
		"\nBatch VRF Coordinator Address:", contractAddresses.BatchCoordinatorAddress,
		"\nVRF Consumer Address:", consumerAddress,
		"\nVRF Subscription Id:", subID,
		"\nVRF Subscription Balance:", *subscriptionBalanceJuels,
		"\nPossible VRF Request command: ",
		fmt.Sprintf("go run . eoa-load-test-request-with-metrics --consumer-address=%s --sub-id=%d --key-hash=%s --request-confirmations %d --requests 1 --runs 1 --cb-gas-limit 1_000_000", consumerAddress, subID, keyHash, *coordinatorConfig.MinConfs),
		"\nRetrieve Request Status: ",
		fmt.Sprintf("go run . eoa-load-test-read-metrics --consumer-address=%s", consumerAddress),
		"\nA node can now be configured to run a VRF job with the below job spec :\n",
		formattedVrfPrimaryJobSpec,
	)

	return model.JobSpecs{
		VRFPrimaryNode: formattedVrfPrimaryJobSpec,
		VRFBackupyNode: formattedVrfBackupJobSpec,
		BHSNode:        formattedBHSJobSpec,
		BHSBackupNode:  formattedBHSBackupJobSpec,
		BHFNode:        formattedBHFJobSpec,
	}
}

func DeployWrapperUniverse(e helpers.Environment) {
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

	wrapper, subID := WrapperDeploy(e,
		common.HexToAddress(*linkAddress),
		common.HexToAddress(*linkETHFeedAddress),
		common.HexToAddress(*coordinatorAddress))

	WrapperConfigure(e,
		wrapper,
		*wrapperGasOverhead,
		*coordinatorGasOverhead,
		*wrapperPremiumPercentage,
		*keyHash,
		*maxNumWords)

	consumer := WrapperConsumerDeploy(e,
		common.HexToAddress(*linkAddress),
		wrapper)

	coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddress), e.Ec)
	helpers.PanicErr(err)

	EoaFundSubscription(e, *coordinator, *linkAddress, amount, subID)

	link, err := link_token_interface.NewLinkToken(common.HexToAddress(*linkAddress), e.Ec)
	helpers.PanicErr(err)
	consumerAmount, s := big.NewInt(0).SetString(*consumerFunding, 10)
	if !s {
		panic(fmt.Sprintf("failed to parse top up amount '%s'", *consumerFunding))
	}

	tx, err := link.Transfer(e.Owner, consumer, consumerAmount)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "link transfer to consumer")

	fmt.Println("wrapper universe deployment complete")
	fmt.Println("wrapper address:", wrapper.String())
	fmt.Println("wrapper consumer address:", consumer.String())
}
