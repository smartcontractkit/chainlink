package scripts

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/smartcontractkit/chainlink/core/scripts/vrfv2/testnet/constants"
	"github.com/smartcontractkit/chainlink/core/scripts/vrfv2/testnet/jobs"
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

var (
	VRFPrimaryNodeName = "vrf-primary-node"
	VRFBackupNodeName  = "vrf-backup-node"
	BHSNodeName        = "bhs-node"
	BHFNodeName        = "bhf-node"
)

type Node struct {
	URL                         string
	CredsFile                   string
	SendingKeys                 []string
	NumberOfSendingKeysToCreate int
	SendingKeyFundingAmount     int64
	VrfKeys                     []string
	jobSpec                     string
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

	subscriptionBalanceString := deployCmd.String("subscription-balance", constants.SubscriptionBalanceString, "amount to fund subscription")

	// optional flags
	fallbackWeiPerUnitLinkString := deployCmd.String("fallback-wei-per-unit-link", constants.FallbackWeiPerUnitLinkString, "fallback wei/link ratio")
	registerKeyUncompressedPubKey := deployCmd.String("uncompressed-pub-key", "", "uncompressed public key")
	vrfPrimaryNodeSendingKeysString := deployCmd.String("vrf-primary-node-sending-keys", "", "VRF Primary Node sending keys")
	vrfBackupNodeSendingKeysString := deployCmd.String("vrf-backup-node-sending-keys", "", "VRF Backup Node sending keys")
	bhsNodeSendingKeysString := deployCmd.String("bhs-node-sending-keys", "", "BHS Node sending keys")
	bhfNodeSendingKeysString := deployCmd.String("bhf-node-sending-keys", "", "BHF Node sending keys")
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
	subscriptionBalance := decimal.RequireFromString(*subscriptionBalanceString).BigInt()

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
	vrfBackupNodeSendingKeys := strings.Split(*vrfBackupNodeSendingKeysString, ",")
	bhsNodeSendingKeys := strings.Split(*bhsNodeSendingKeysString, ",")
	bhfNodeSendingKeys := strings.Split(*bhfNodeSendingKeysString, ",")

	nodes := make(map[string]Node)

	nodes[VRFPrimaryNodeName] = Node{
		SendingKeys: vrfPrimaryNodeSendingKeys,
	}
	nodes[VRFBackupNodeName] = Node{
		SendingKeys: vrfBackupNodeSendingKeys,
	}
	nodes[BHSNodeName] = Node{
		SendingKeys: bhsNodeSendingKeys,
	}
	nodes[BHFNodeName] = Node{
		SendingKeys: bhfNodeSendingKeys,
	}

	bhsContractAddress := common.HexToAddress(bhsContractAddressString)
	batchBHSAddress := common.HexToAddress(batchBHSAddressString)
	coordinatorAddress := common.HexToAddress(coordinatorAddressString)
	batchCoordinatorAddress := common.HexToAddress(batchCoordinatorAddressString)

	VRFV2DeployUniverse(
		e,
		fallbackWeiPerUnitLink,
		subscriptionBalance,
		registerKeyUncompressedPubKey,
		linkAddress,
		linkEthAddress,
		bhsContractAddress,
		batchBHSAddress,
		coordinatorAddress,
		batchCoordinatorAddress,
		minConfs,
		maxGasLimit,
		stalenessSeconds,
		gasAfterPayment,
		feeConfig,
		nodes,
	)
}

func VRFV2DeployUniverse(
	e helpers.Environment,
	fallbackWeiPerUnitLink *big.Int,
	subscriptionBalance *big.Int,
	registerKeyUncompressedPubKey *string,
	linkAddress string,
	linkEthAddress string,
	bhsContractAddress common.Address,
	batchBHSAddress common.Address,
	coordinatorAddress common.Address,
	batchCoordinatorAddress common.Address,
	minConfs *int,
	maxGasLimit *int64,
	stalenessSeconds *int64,
	gasAfterPayment *int64,
	feeConfig vrf_coordinator_v2.VRFCoordinatorV2FeeConfig,
	nodesMap map[string]Node, //todo - is possible to pass a pointer to the node, so that we read data and also update the data and we dont need to return "nodes"
) (string, string, string, string) {

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

	if len(linkAddress) == 0 {
		fmt.Println("\nDeploying LINK Token...")
		linkAddress = helpers.DeployLinkToken(e).String()
	}

	if len(linkEthAddress) == 0 {
		fmt.Println("\nDeploying LINK/ETH Feed...")
		linkEthAddress = helpers.DeployLinkEthFeed(e, linkAddress, fallbackWeiPerUnitLink).String()
	}

	if bhsContractAddress.String() == "0x0000000000000000000000000000000000000000" {
		fmt.Println("\nDeploying BHS...")
		bhsContractAddress = DeployBHS(e)
	}

	if batchBHSAddress.String() == "0x0000000000000000000000000000000000000000" {
		fmt.Println("\nDeploying Batch BHS...")
		batchBHSAddress = DeployBatchBHS(e, bhsContractAddress)
	}

	if coordinatorAddress.String() == "0x0000000000000000000000000000000000000000" {
		fmt.Println("\nDeploying Coordinator...")
		coordinatorAddress = DeployCoordinator(e, linkAddress, bhsContractAddress.String(), linkEthAddress)
	}

	coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(coordinatorAddress, e.Ec)
	helpers.PanicErr(err)

	if batchCoordinatorAddress.String() == "0x0000000000000000000000000000000000000000" {
		fmt.Println("\nDeploying Batch Coordinator...")
		batchCoordinatorAddress = deployBatchCoordinatorV2(e, coordinatorAddress)
	}

	fmt.Println("\nSetting Coordinator Config...")
	SetCoordinatorConfig(
		e,
		*coordinator,
		uint16(*minConfs),
		uint32(*maxGasLimit),
		uint32(*stalenessSeconds),
		uint32(*gasAfterPayment),
		fallbackWeiPerUnitLink,
		feeConfig,
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
	consumerAddress := EoaLoadTestConsumerWithMetricsDeploy(e, coordinatorAddress.String())

	fmt.Println("\nAdding subscription...")
	EoaCreateSub(e, *coordinator)
	subID := uint64(1)

	fmt.Println("\nAdding consumer to subscription...")
	EoaAddConsumerToSub(e, *coordinator, subID, consumerAddress.String())

	if subscriptionBalance.Cmp(big.NewInt(0)) > 0 {
		fmt.Println("\nFunding subscription with", subscriptionBalance, "juels...")
		EoaFundSubscription(e, *coordinator, linkAddress, subscriptionBalance, subID)
	} else {
		fmt.Println("Subscription", subID, "NOT getting funded. You must fund the subscription in order to use it!")
	}

	fmt.Println("\nSubscribed and (possibly) funded, retrieving subscription from deployed contract...")
	s, err := coordinator.GetSubscription(nil, subID)
	helpers.PanicErr(err)
	fmt.Printf("Subscription %+v\n", s)

	formattedVrfPrimaryJobSpec := fmt.Sprintf(
		jobs.VRFJobFormatted,
		coordinatorAddress,      //coordinatorAddress
		batchCoordinatorAddress, //batchCoordinatorAddress
		false,                   //batchFulfillmentEnabled
		compressedPkHex,         //publicKey
		*minConfs,               //minIncomingConfirmations
		e.ChainID,               //evmChainID
		strings.Join(nodesMap[VRFPrimaryNodeName].SendingKeys, "\",\""), //fromAddresses
		coordinatorAddress,
		coordinatorAddress,
		coordinatorAddress,
	)

	formattedVrfBackupJobSpec := fmt.Sprintf(
		jobs.VRFJobFormatted,
		coordinatorAddress,      //coordinatorAddress
		batchCoordinatorAddress, //batchCoordinatorAddress
		false,                   //batchFulfillmentEnabled
		compressedPkHex,         //publicKey
		100,                     //minIncomingConfirmations
		e.ChainID,               //evmChainID
		strings.Join(nodesMap[VRFBackupNodeName].SendingKeys, "\",\""), //fromAddresses
		coordinatorAddress,
		coordinatorAddress,
		coordinatorAddress,
	)

	formattedBHSJobSpec := fmt.Sprintf(
		jobs.BHSJobFormatted,
		coordinatorAddress, //coordinatorAddress
		bhsContractAddress, //bhs adreess
		e.ChainID,          //chain id
		strings.Join(nodesMap[BHSNodeName].SendingKeys, "\",\""), //sending addresses
	)

	formattedBHFJobSpec := fmt.Sprintf(
		jobs.BHFJobFormatted,
		coordinatorAddress, //coordinatorAddress
		bhsContractAddress, //bhs adreess
		batchBHSAddress,    //batchBHS
		strings.Join(nodesMap[BHSNodeName].SendingKeys, "\",\""), //sending addresses
	)

	fmt.Println(
		"\nDeployment complete.",
		"\nLINK Token contract address:", linkAddress,
		"\nLINK/ETH Feed contract address:", linkEthAddress,
		"\nBlockhash Store contract address:", bhsContractAddress,
		"\nBatch Blockhash Store contract address:", batchBHSAddress,
		"\nVRF Coordinator Address:", coordinatorAddress,
		"\nBatch VRF Coordinator Address:", batchCoordinatorAddress,
		"\nVRF Consumer Address:", consumerAddress,
		"\nVRF Subscription Id:", subID,
		"\nVRF Subscription Balance:", *subscriptionBalance,
		"\nPossible VRF Request command: ",
		fmt.Sprintf("go run . eoa-load-test-request-with-metrics --consumer-address=%s --sub-id=%d --key-hash=%s --request-confirmations 1 --requests 1 --runs 1 --cb-gas-limit 1_000_000", consumerAddress, subID, keyHash),
		"\nRetrieve Request Status: ",
		fmt.Sprintf("go run . eoa-load-test-read-metrics --consumer-address=%s", consumerAddress),
		"\nA node can now be configured to run a VRF job with the below job spec :\n",
		formattedVrfPrimaryJobSpec,
	)

	return formattedVrfPrimaryJobSpec, formattedVrfBackupJobSpec, formattedBHSJobSpec, formattedBHFJobSpec
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
