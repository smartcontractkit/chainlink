package v2plusscripts

import (
	"bytes"
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/smartcontractkit/chainlink/core/scripts/common/vrf/constants"
	"github.com/smartcontractkit/chainlink/core/scripts/common/vrf/jobs"
	"github.com/smartcontractkit/chainlink/core/scripts/common/vrf/model"
	"github.com/smartcontractkit/chainlink/core/scripts/common/vrf/util"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2plus_interface"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_sub_owner"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/proof"
)

var coordinatorV2PlusABI = evmtypes.MustGetABI(vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalABI)

type CoordinatorConfigV2Plus struct {
	MinConfs                          int
	MaxGasLimit                       int64
	StalenessSeconds                  int64
	GasAfterPayment                   int64
	FallbackWeiPerUnitLink            *big.Int
	FulfillmentFlatFeeNativePPM       uint32
	FulfillmentFlatFeeLinkDiscountPPM uint32
	NativePremiumPercentage           uint8
	LinkPremiumPercentage             uint8
}

func SmokeTestVRF(e helpers.Environment) {
	smokeCmd := flag.NewFlagSet("smoke", flag.ExitOnError)

	// required flags
	coordinatorType := smokeCmd.String("coordinator-type", "", "Specify which coordinator type to use: layer1, arbitrum, optimism")
	linkAddress := smokeCmd.String("link-address", "", "address of link token")
	linkNativeAddress := smokeCmd.String("link-native-feed", "", "address of link native feed")
	bhsAddressStr := smokeCmd.String("bhs-address", "", "address of blockhash store")
	batchBHSAddressStr := smokeCmd.String("batch-bhs-address", "", "address of batch blockhash store")
	coordinatorAddressStr := smokeCmd.String("coordinator-address", "", "address of the vrf coordinator v2 contract")
	batchCoordinatorAddressStr := smokeCmd.String("batch-coordinator-address", "", "address of the batch vrf coordinator v2 contract")
	subscriptionBalanceString := smokeCmd.String("subscription-balance", "1e19", "amount to fund subscription")
	skipConfig := smokeCmd.Bool("skip-config", false, "skip setting coordinator config")

	// optional flags
	fallbackWeiPerUnitLinkString := smokeCmd.String("fallback-wei-per-unit-link", "6e16", "fallback wei/link ratio")
	minConfs := smokeCmd.Int("min-confs", 3, "min confs")
	maxGasLimit := smokeCmd.Int64("max-gas-limit", 2.5e6, "max gas limit")
	stalenessSeconds := smokeCmd.Int64("staleness-seconds", 86400, "staleness in seconds")
	gasAfterPayment := smokeCmd.Int64("gas-after-payment", 33285, "gas after payment calculation")
	flatFeeNativePPM := smokeCmd.Int64("flat-fee-native-ppm", 500, "fulfillment flat fee Native ppm")
	flatFeeLinkDiscountPPM := smokeCmd.Int64("flat-fee-link-discount-ppm", 100, "fulfillment flat fee discount for LINK payment denominated in native ppm")
	nativePremiumPercentage := smokeCmd.Int64("native-premium-percentage", 1, "premium percentage for native payment")
	linkPremiumPercentage := smokeCmd.Int64("link-premium-percentage", 1, "premium percentage for LINK payment")
	gasLaneMaxGas := smokeCmd.Int64("gas-lane-max-gas", 1e12, "gas lane max gas price")

	helpers.ParseArgs(
		smokeCmd, os.Args[2:],
	)

	fallbackWeiPerUnitLink := decimal.RequireFromString(*fallbackWeiPerUnitLinkString).BigInt()
	subscriptionBalance := decimal.RequireFromString(*subscriptionBalanceString).BigInt()

	// generate VRF key
	key, err := vrfkey.NewV2()
	helpers.PanicErr(err)
	fmt.Println("vrf private key:", hexutil.Encode(key.Raw()))
	fmt.Println("vrf public key:", key.PublicKey.String())
	fmt.Println("vrf key hash:", key.PublicKey.MustHash())

	if len(*linkAddress) == 0 {
		fmt.Println("\nDeploying LINK Token...")
		address := helpers.DeployLinkToken(e).String()
		linkAddress = &address
	}

	if len(*linkNativeAddress) == 0 {
		fmt.Println("\nDeploying LINK/Native Feed...")
		address := helpers.DeployLinkEthFeed(e, *linkAddress, fallbackWeiPerUnitLink).String()
		linkNativeAddress = &address
	}

	var bhsContractAddress common.Address
	if len(*bhsAddressStr) == 0 {
		fmt.Println("\nDeploying BHS...")
		bhsContractAddress = DeployBHS(e)
	} else {
		bhsContractAddress = common.HexToAddress(*bhsAddressStr)
	}

	var batchBHSAddress common.Address
	if len(*batchBHSAddressStr) == 0 {
		fmt.Println("\nDeploying Batch BHS...")
		batchBHSAddress = DeployBatchBHS(e, bhsContractAddress)
	} else {
		batchBHSAddress = common.HexToAddress(*batchBHSAddressStr)
	}

	var coordinatorAddress common.Address
	if len(*coordinatorAddressStr) == 0 {
		fmt.Printf("\nDeploying Coordinator [type=%s]...\n", *coordinatorType)
		coordinatorAddress = DeployCoordinator(e, *linkAddress, bhsContractAddress.String(), *linkNativeAddress, *coordinatorType)
	} else {
		coordinatorAddress = common.HexToAddress(*coordinatorAddressStr)
	}

	coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(coordinatorAddress, e.Ec)
	helpers.PanicErr(err)

	var batchCoordinatorAddress common.Address
	if len(*batchCoordinatorAddressStr) == 0 {
		fmt.Println("\nDeploying Batch Coordinator...")
		batchCoordinatorAddress = DeployBatchCoordinatorV2(e, coordinatorAddress)
	} else {
		batchCoordinatorAddress = common.HexToAddress(*batchCoordinatorAddressStr)
	}

	if !*skipConfig {
		fmt.Println("\nSetting Coordinator Config...")
		SetCoordinatorConfig(
			e,
			*coordinator,
			uint16(*minConfs),
			uint32(*maxGasLimit),
			uint32(*stalenessSeconds),
			uint32(*gasAfterPayment),
			fallbackWeiPerUnitLink,
			uint32(*flatFeeNativePPM),
			uint32(*flatFeeLinkDiscountPPM),
			uint8(*nativePremiumPercentage),
			uint8(*linkPremiumPercentage),
		)
	}

	fmt.Println("\nConfig set, getting current config from deployed contract...")
	PrintCoordinatorConfig(coordinator)

	// Generate compressed public key and key hash
	uncompressed, err := key.PublicKey.StringUncompressed()
	helpers.PanicErr(err)
	if strings.HasPrefix(uncompressed, "0x") {
		uncompressed = strings.Replace(uncompressed, "0x", "04", 1)
	}
	pubBytes, err := hex.DecodeString(uncompressed)
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
	fmt.Println("vrf key hash from unmarshal:", hexutil.Encode(keyHash[:]))
	fmt.Println("vrf key hash from key:", key.PublicKey.MustHash())
	if kh := key.PublicKey.MustHash(); !bytes.Equal(keyHash[:], kh[:]) {
		panic(fmt.Sprintf("unexpected key hash %s, expected %s", hexutil.Encode(keyHash[:]), key.PublicKey.MustHash().String()))
	}
	fmt.Println("compressed public key from unmarshal:", compressedPkHex)
	fmt.Println("compressed public key from key:", key.PublicKey.String())
	if compressedPkHex != key.PublicKey.String() {
		panic(fmt.Sprintf("unexpected compressed public key %s, expected %s", compressedPkHex, key.PublicKey.String()))
	}

	kh1, err := coordinator.HashOfKey(nil, [2]*big.Int{pk.X, pk.Y})
	helpers.PanicErr(err)
	fmt.Println("key hash from coordinator:", hexutil.Encode(kh1[:]))
	if !bytes.Equal(kh1[:], keyHash[:]) {
		panic(fmt.Sprintf("unexpected key hash %s, expected %s", hexutil.Encode(kh1[:]), hexutil.Encode(keyHash[:])))
	}

	fmt.Println("\nRegistering proving key...")
	point, err := key.PublicKey.Point()
	helpers.PanicErr(err)
	x, y := secp256k1.Coordinates(point)
	fmt.Println("proving key points x:", x, ", y:", y)
	fmt.Println("proving key points from unmarshal:", pk.X, pk.Y)
	tx, err := coordinator.RegisterProvingKey(e.Owner, [2]*big.Int{x, y}, uint64(*gasLaneMaxGas))
	helpers.PanicErr(err)
	registerReceipt := helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "register proving key on", coordinatorAddress.String())
	var provingKeyRegisteredLog *vrf_coordinator_v2_5.VRFCoordinatorV25ProvingKeyRegistered
	for _, log := range registerReceipt.Logs {
		if log.Address == coordinatorAddress {
			var err2 error
			provingKeyRegisteredLog, err2 = coordinator.ParseProvingKeyRegistered(*log)
			if err2 != nil {
				continue
			}
		}
	}
	if provingKeyRegisteredLog == nil {
		panic("no proving key registered log found")
	}
	if !bytes.Equal(provingKeyRegisteredLog.KeyHash[:], keyHash[:]) {
		panic(fmt.Sprintf("unexpected key hash registered %s, expected %s", hexutil.Encode(provingKeyRegisteredLog.KeyHash[:]), hexutil.Encode(keyHash[:])))
	}
	fmt.Println("key hash registered:", hexutil.Encode(provingKeyRegisteredLog.KeyHash[:]))

	fmt.Println("\nProving key registered, getting proving key hashes from deployed contract...")
	registerdKeyHash, err := coordinator.SProvingKeyHashes(nil, big.NewInt(0))
	helpers.PanicErr(err)
	fmt.Printf("Key hash registered: %x\n", registerdKeyHash)
	ourKeyHash := key.PublicKey.MustHash()
	if !bytes.Equal(registerdKeyHash[:], ourKeyHash[:]) {
		panic(fmt.Sprintf("unexpected key hash %s, expected %s", hexutil.Encode(registerdKeyHash[:]), hexutil.Encode(ourKeyHash[:])))
	}

	fmt.Println("\nDeploying consumer...")
	consumerAddress := EoaDeployConsumer(e, coordinatorAddress.String(), *linkAddress)

	fmt.Println("\nAdding subscription...")
	subID, err := EoaCreateSub(e, *coordinator)
	helpers.PanicErr(err)

	fmt.Println("\nAdding consumer to subscription...")
	EoaAddConsumerToSub(e, *coordinator, subID, consumerAddress.String())

	if subscriptionBalance.Cmp(big.NewInt(0)) > 0 {
		fmt.Println("\nFunding subscription with", subscriptionBalance, "juels...")
		EoaFundSubWithLink(e, *coordinator, *linkAddress, subscriptionBalance, subID)
	} else {
		fmt.Println("Subscription", subID, "NOT getting funded. You must fund the subscription in order to use it!")
	}

	fmt.Println("\nSubscribed and (possibly) funded, retrieving subscription from deployed contract...")
	s, err := coordinator.GetSubscription(nil, subID)
	helpers.PanicErr(err)
	fmt.Printf("Subscription %+v\n", s)

	fmt.Println(
		"\nDeployment complete.",
		"\nLINK Token contract address:", *linkAddress,
		"\nLINK/Native Feed contract address:", *linkNativeAddress,
		"\nBlockhash Store contract address:", bhsContractAddress,
		"\nBatch Blockhash Store contract address:", batchBHSAddress,
		"\nVRF Coordinator Address:", coordinatorAddress,
		"\nBatch VRF Coordinator Address:", batchCoordinatorAddress,
		"\nVRF Consumer Address:", consumerAddress,
		"\nVRF Subscription Id:", subID,
		"\nVRF Subscription Balance:", *subscriptionBalanceString,
	)

	fmt.Println("making a request on consumer", consumerAddress)
	consumer, err := vrf_v2plus_sub_owner.NewVRFV2PlusExternalSubOwnerExample(consumerAddress, e.Ec)
	helpers.PanicErr(err)
	tx, err = consumer.RequestRandomWords(e.Owner, subID, 100_000, 3, 3, provingKeyRegisteredLog.KeyHash, false)
	helpers.PanicErr(err)
	receipt := helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "request random words from", consumerAddress.String())
	fmt.Println("request blockhash:", receipt.BlockHash)

	// extract the RandomWordsRequested log from the receipt logs
	var rwrLog *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsRequested
	for _, log := range receipt.Logs {
		if log.Address == coordinatorAddress {
			var err2 error
			rwrLog, err2 = coordinator.ParseRandomWordsRequested(*log)
			if err2 != nil {
				continue
			}
		}
	}
	if rwrLog == nil {
		panic("no RandomWordsRequested log found")
	}

	fmt.Println("key hash:", hexutil.Encode(rwrLog.KeyHash[:]))
	fmt.Println("request id:", rwrLog.RequestId)
	fmt.Println("preseed:", rwrLog.PreSeed)
	fmt.Println("num words:", rwrLog.NumWords)
	fmt.Println("callback gas limit:", rwrLog.CallbackGasLimit)
	fmt.Println("sender:", rwrLog.Sender)
	fmt.Println("extra args:", hexutil.Encode(rwrLog.ExtraArgs))

	// generate the VRF proof, follow the same process as the node
	// we assume there is enough funds in the subscription to pay for the gas
	preSeed, err := proof.BigToSeed(rwrLog.PreSeed)
	helpers.PanicErr(err)

	preSeedData := proof.PreSeedDataV2Plus{
		PreSeed:          preSeed,
		BlockHash:        rwrLog.Raw.BlockHash,
		BlockNum:         rwrLog.Raw.BlockNumber,
		SubId:            rwrLog.SubId,
		CallbackGasLimit: rwrLog.CallbackGasLimit,
		NumWords:         rwrLog.NumWords,
		Sender:           rwrLog.Sender,
		ExtraArgs:        rwrLog.ExtraArgs,
	}
	finalSeed := proof.FinalSeedV2Plus(preSeedData)
	pf, err := key.GenerateProof(finalSeed)
	helpers.PanicErr(err)
	onChainProof, rc, err := proof.GenerateProofResponseFromProofV2Plus(pf, preSeedData)
	helpers.PanicErr(err)
	b, err := coordinatorV2PlusABI.Pack("fulfillRandomWords", onChainProof, rc)
	helpers.PanicErr(err)
	fmt.Println("calldata for fulfillRandomWords:", hexutil.Encode(b))

	// call fulfillRandomWords with onChainProof and rc appropriately
	fmt.Println("proof c:", onChainProof.C)
	fmt.Println("proof s:", onChainProof.S)
	fmt.Println("proof gamma:", onChainProof.Gamma)
	fmt.Println("proof seed:", onChainProof.Seed)
	fmt.Println("proof pk:", onChainProof.Pk)
	fmt.Println("proof c gamma witness:", onChainProof.CGammaWitness)
	fmt.Println("proof u witness:", onChainProof.UWitness)
	fmt.Println("proof s hash witness:", onChainProof.SHashWitness)
	fmt.Println("proof z inv:", onChainProof.ZInv)
	fmt.Println("request commitment sub id:", rc.SubId)
	fmt.Println("request commitment callback gas limit:", rc.CallbackGasLimit)
	fmt.Println("request commitment num words:", rc.NumWords)
	fmt.Println("request commitment sender:", rc.Sender)
	fmt.Println("request commitment extra args:", hexutil.Encode(rc.ExtraArgs))

	receipt, txHash := sendTx(e, coordinatorAddress, b)
	if receipt.Status != 1 {
		fmt.Println("fulfillment tx failed, extracting revert reason")
		tx, _, err := e.Ec.TransactionByHash(context.Background(), txHash)
		helpers.PanicErr(err)
		call := ethereum.CallMsg{
			From:     e.Owner.From,
			To:       tx.To(),
			Data:     tx.Data(),
			Gas:      tx.Gas(),
			GasPrice: tx.GasPrice(),
		}
		r, err := e.Ec.CallContract(context.Background(), call, receipt.BlockNumber)
		fmt.Println("call contract", "r", r, "err", err)
		rpcError, err := evmclient.ExtractRPCError(err)
		fmt.Println("extracting rpc error", rpcError.String(), err)
		os.Exit(1)
	}

	fmt.Println("\nfulfillment successful")
}

func SmokeTestBHS(e helpers.Environment) {
	smokeCmd := flag.NewFlagSet("smoke-bhs", flag.ExitOnError)

	// optional args
	bhsAddress := smokeCmd.String("bhs-address", "", "address of blockhash store")
	batchBHSAddress := smokeCmd.String("batch-bhs-address", "", "address of batch blockhash store")

	helpers.ParseArgs(smokeCmd, os.Args[2:])

	var bhsContractAddress common.Address
	if len(*bhsAddress) == 0 {
		fmt.Println("\nDeploying BHS...")
		bhsContractAddress = DeployBHS(e)
	} else {
		bhsContractAddress = common.HexToAddress(*bhsAddress)
	}

	var batchBHSContractAddress common.Address
	if len(*batchBHSAddress) == 0 {
		fmt.Println("\nDeploying Batch BHS...")
		batchBHSContractAddress = DeployBatchBHS(e, bhsContractAddress)
	} else {
		batchBHSContractAddress = common.HexToAddress(*batchBHSAddress)
	}

	bhs, err := blockhash_store.NewBlockhashStore(bhsContractAddress, e.Ec)
	helpers.PanicErr(err)

	batchBHS, err := batch_blockhash_store.NewBatchBlockhashStore(batchBHSContractAddress, e.Ec)
	helpers.PanicErr(err)
	batchBHS.Address()

	fmt.Println("\nexecuting storeEarliest")
	tx, err := bhs.StoreEarliest(e.Owner)
	helpers.PanicErr(err)
	seReceipt := helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "storeEarliest on", bhsContractAddress.String())
	var anchorBlockNumber *big.Int
	if seReceipt.Status != 1 {
		fmt.Println("storeEarliest failed")
		os.Exit(1)
	}
	fmt.Println("storeEarliest succeeded, checking BH is there")
	bh, err := bhs.GetBlockhash(nil, seReceipt.BlockNumber.Sub(seReceipt.BlockNumber, big.NewInt(256)))
	helpers.PanicErr(err)
	fmt.Println("blockhash stored by storeEarliest:", hexutil.Encode(bh[:]))
	anchorBlockNumber = seReceipt.BlockNumber

	if anchorBlockNumber == nil {
		panic("no anchor block number")
	}

	fmt.Println("\nexecuting store(n)")
	latestHead, err := e.Ec.HeaderByNumber(context.Background(), nil)
	helpers.PanicErr(err)
	toStore := latestHead.Number.Sub(latestHead.Number, big.NewInt(1))
	tx, err = bhs.Store(e.Owner, toStore)
	helpers.PanicErr(err)
	sReceipt := helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "store on", bhsContractAddress.String())
	if sReceipt.Status != 1 {
		fmt.Println("store failed")
		os.Exit(1)
	}
	fmt.Println("store succeeded, checking BH is there")
	bh, err = bhs.GetBlockhash(nil, toStore)
	helpers.PanicErr(err)
	fmt.Println("blockhash stored by store:", hexutil.Encode(bh[:]))

	fmt.Println("\nexecuting storeVerifyHeader")
	headers, _, err := helpers.GetRlpHeaders(e, []*big.Int{anchorBlockNumber}, false)
	helpers.PanicErr(err)

	toStore = anchorBlockNumber.Sub(anchorBlockNumber, big.NewInt(1))
	tx, err = bhs.StoreVerifyHeader(e.Owner, toStore, headers[0])
	helpers.PanicErr(err)
	svhReceipt := helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "storeVerifyHeader on", bhsContractAddress.String())
	if svhReceipt.Status != 1 {
		fmt.Println("storeVerifyHeader failed")
		os.Exit(1)
	}
	fmt.Println("storeVerifyHeader succeeded, checking BH is there")
	bh, err = bhs.GetBlockhash(nil, toStore)
	helpers.PanicErr(err)
	fmt.Println("blockhash stored by storeVerifyHeader:", hexutil.Encode(bh[:]))
}

func sendTx(e helpers.Environment, to common.Address, data []byte) (*types.Receipt, common.Hash) {
	nonce, err := e.Ec.PendingNonceAt(context.Background(), e.Owner.From)
	helpers.PanicErr(err)
	gasPrice, err := e.Ec.SuggestGasPrice(context.Background())
	helpers.PanicErr(err)
	rawTx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Data:     data,
		Value:    big.NewInt(0),
		Gas:      1_000_000,
		GasPrice: gasPrice,
	})
	signedTx, err := e.Owner.Signer(e.Owner.From, rawTx)
	helpers.PanicErr(err)
	err = e.Ec.SendTransaction(context.Background(), signedTx)
	helpers.PanicErr(err)
	return helpers.ConfirmTXMined(context.Background(), e.Ec, signedTx,
		e.ChainID, "send tx", signedTx.Hash().String(), "to", to.String()), signedTx.Hash()
}

func sendNativeTokens(e helpers.Environment, to common.Address, amount *big.Int) (*types.Receipt, common.Hash) {
	nonce, err := e.Ec.PendingNonceAt(context.Background(), e.Owner.From)
	helpers.PanicErr(err)
	gasPrice, err := e.Ec.SuggestGasPrice(context.Background())
	helpers.PanicErr(err)
	msg := ethereum.CallMsg{
		From:     e.Owner.From,
		To:       &to,
		Value:    amount,
		Gas:      0,
		GasPrice: big.NewInt(0),
		Data:     nil,
	}
	gasLimit, err := e.Ec.EstimateGas(context.Background(), msg)
	helpers.PanicErr(err)
	rawTx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Data:     nil,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: gasPrice,
	})
	signedTx, err := e.Owner.Signer(e.Owner.From, rawTx)
	helpers.PanicErr(err)
	err = e.Ec.SendTransaction(context.Background(), signedTx)
	helpers.PanicErr(err)
	return helpers.ConfirmTXMined(context.Background(), e.Ec, signedTx,
		e.ChainID, "send tx", signedTx.Hash().String(), "to", to.String()), signedTx.Hash()
}

func DeployUniverseViaCLI(e helpers.Environment) {
	deployCmd := flag.NewFlagSet("deploy-universe", flag.ExitOnError)

	// required flags
	nativeOnly := deployCmd.Bool("native-only", false, "if true, link and link feed are not set up")
	linkAddress := deployCmd.String("link-address", "", "address of link token")
	linkNativeAddress := deployCmd.String("link-native-feed", "", "address of link native feed")
	bhsContractAddressString := deployCmd.String("bhs-address", "", "address of BHS contract")
	batchBHSAddressString := deployCmd.String("batch-bhs-address", "", "address of Batch BHS contract")
	coordinatorAddressString := deployCmd.String("coordinator-address", "", "address of VRF Coordinator contract")
	coordinatorType := deployCmd.String("coordinator-type", "", "Specify which coordinator type to use: layer1, arbitrum, optimism")
	batchCoordinatorAddressString := deployCmd.String("batch-coordinator-address", "", "address Batch VRF Coordinator contract")
	subscriptionBalanceJuelsString := deployCmd.String("subscription-balance", "1e19", "amount to fund subscription with Link token (Juels)")
	subscriptionBalanceNativeWeiString := deployCmd.String("subscription-balance-native", "1e18", "amount to fund subscription with native token (Wei)")

	batchFulfillmentEnabled := deployCmd.Bool("batch-fulfillment-enabled", constants.BatchFulfillmentEnabled, "whether send randomness fulfillments in batches inside one tx from CL node")
	batchFulfillmentGasMultiplier := deployCmd.Float64("batch-fulfillment-gas-multiplier", 1.1, "")
	estimateGasMultiplier := deployCmd.Float64("estimate-gas-multiplier", 1.1, "")
	pollPeriod := deployCmd.String("poll-period", "300ms", "")
	requestTimeout := deployCmd.String("request-timeout", "30m0s", "")
	bhsJobWaitBlocks := flag.Int("bhs-job-wait-blocks", 30, "")
	bhsJobLookBackBlocks := flag.Int("bhs-job-look-back-blocks", 200, "")
	bhsJobPollPeriod := flag.String("bhs-job-poll-period", "3s", "")
	bhsJobRunTimeout := flag.String("bhs-job-run-timeout", "1m", "")
	simulationBlock := deployCmd.String("simulation-block", "latest", "simulation block can be 'pending' or 'latest'")

	// optional flags
	fallbackWeiPerUnitLinkString := deployCmd.String("fallback-wei-per-unit-link", "6e16", "fallback wei/link ratio")
	registerVRFKeyUncompressedPubKey := deployCmd.String("uncompressed-pub-key", "", "uncompressed public key")

	vrfPrimaryNodeSendingKeysString := deployCmd.String("vrf-primary-node-sending-keys", "", "VRF Primary Node sending keys")
	minConfs := deployCmd.Int("min-confs", constants.MinConfs, "min confs")
	nodeSendingKeyFundingAmount := deployCmd.String("sending-key-funding-amount", constants.NodeSendingKeyFundingAmount, "CL node sending key funding amount")
	maxGasLimit := deployCmd.Int64("max-gas-limit", constants.MaxGasLimit, "max gas limit")
	stalenessSeconds := deployCmd.Int64("staleness-seconds", constants.StalenessSeconds, "staleness in seconds")
	gasAfterPayment := deployCmd.Int64("gas-after-payment", constants.GasAfterPayment, "gas after payment calculation")
	flatFeeNativePPM := deployCmd.Int64("flat-fee-native-ppm", 500, "fulfillment flat fee Native ppm")
	flatFeeLinkDiscountPPM := deployCmd.Int64("flat-fee-link-discount-ppm", 100, "fulfillment flat fee discount for LINK payment denominated in native ppm")
	nativePremiumPercentage := deployCmd.Int64("native-premium-percentage", 1, "premium percentage for native payment")
	linkPremiumPercentage := deployCmd.Int64("link-premium-percentage", 1, "premium percentage for LINK payment")
	provingKeyMaxGasPriceString := deployCmd.String("proving-key-max-gas-price", "1e12", "gas lane max gas price")

	// only necessary for Optimism coordinator contract
	optimismL1GasFeeCalculationMode := deployCmd.Uint64("optimism-l1-fee-mode", 0, "Choose Optimism coordinator contract L1 fee calculation mode: 0, 1, 2")
	optimismL1GasFeeCoefficient := deployCmd.Uint64("optimism-l1-fee-coefficient", 100, "Choose Optimism coordinator contract L1 fee coefficient percentage [1, 100]")

	helpers.ParseArgs(
		deployCmd, os.Args[2:],
	)

	if *coordinatorType != "layer1" && *coordinatorType != "arbitrum" && *coordinatorType != "optimism" {
		panic(fmt.Sprintf("Invalid Coordinator type `%s`. Only `layer1`, `arbitrum` and `optimism` are supported", *coordinatorType))
	}

	if *nativeOnly {
		if *linkAddress != "" || *linkNativeAddress != "" {
			panic("native-only flag is set, but link address or link native address is provided")
		}
		if *subscriptionBalanceJuelsString != "0" {
			panic("native-only flag is set, but link subscription balance is provided")
		}
	}

	if *simulationBlock != "pending" && *simulationBlock != "latest" {
		helpers.PanicErr(fmt.Errorf("simulation block must be 'pending' or 'latest'"))
	}

	fallbackWeiPerUnitLink := decimal.RequireFromString(*fallbackWeiPerUnitLinkString).BigInt()
	subscriptionBalanceJuels := decimal.RequireFromString(*subscriptionBalanceJuelsString).BigInt()
	subscriptionBalanceNativeWei := decimal.RequireFromString(*subscriptionBalanceNativeWeiString).BigInt()
	fundingAmount := decimal.RequireFromString(*nodeSendingKeyFundingAmount).BigInt()
	provingKeyMaxGasPrice := decimal.RequireFromString(*provingKeyMaxGasPriceString).BigInt()

	var vrfPrimaryNodeSendingKeys []string
	if len(*vrfPrimaryNodeSendingKeysString) > 0 {
		vrfPrimaryNodeSendingKeys = strings.Split(*vrfPrimaryNodeSendingKeysString, ",")
	}

	nodesMap := make(map[string]model.Node)

	nodesMap[model.VRFPrimaryNodeName] = model.Node{
		SendingKeys:             util.MapToSendingKeyArr(vrfPrimaryNodeSendingKeys),
		SendingKeyFundingAmount: fundingAmount,
	}

	bhsContractAddress := common.HexToAddress(*bhsContractAddressString)
	batchBHSAddress := common.HexToAddress(*batchBHSAddressString)
	coordinatorAddress := common.HexToAddress(*coordinatorAddressString)
	batchCoordinatorAddress := common.HexToAddress(*batchCoordinatorAddressString)

	contractAddresses := model.ContractAddresses{
		LinkAddress:             *linkAddress,
		LinkEthAddress:          *linkNativeAddress,
		BhsContractAddress:      bhsContractAddress,
		BatchBHSAddress:         batchBHSAddress,
		CoordinatorAddress:      coordinatorAddress,
		BatchCoordinatorAddress: batchCoordinatorAddress,
	}

	coordinatorConfig := CoordinatorConfigV2Plus{
		MinConfs:                          *minConfs,
		MaxGasLimit:                       *maxGasLimit,
		StalenessSeconds:                  *stalenessSeconds,
		GasAfterPayment:                   *gasAfterPayment,
		FallbackWeiPerUnitLink:            fallbackWeiPerUnitLink,
		FulfillmentFlatFeeNativePPM:       uint32(*flatFeeNativePPM),
		FulfillmentFlatFeeLinkDiscountPPM: uint32(*flatFeeLinkDiscountPPM),
		NativePremiumPercentage:           uint8(*nativePremiumPercentage),
		LinkPremiumPercentage:             uint8(*linkPremiumPercentage),
	}

	vrfKeyRegistrationConfig := model.VRFKeyRegistrationConfig{
		VRFKeyUncompressedPubKey: *registerVRFKeyUncompressedPubKey,
	}

	coordinatorJobSpecConfig := model.CoordinatorJobSpecConfig{
		BatchFulfillmentEnabled:       *batchFulfillmentEnabled,
		BatchFulfillmentGasMultiplier: *batchFulfillmentGasMultiplier,
		EstimateGasMultiplier:         *estimateGasMultiplier,
		PollPeriod:                    *pollPeriod,
		RequestTimeout:                *requestTimeout,
	}

	bhsJobSpecConfig := model.BHSJobSpecConfig{
		RunTimeout:     *bhsJobRunTimeout,
		WaitBlocks:     *bhsJobWaitBlocks,
		LookBackBlocks: *bhsJobLookBackBlocks,
		PollPeriod:     *bhsJobPollPeriod,
	}

	VRFV2PlusDeployUniverse(
		e,
		subscriptionBalanceJuels,
		subscriptionBalanceNativeWei,
		vrfKeyRegistrationConfig,
		contractAddresses,
		coordinatorConfig,
		*nativeOnly,
		nodesMap,
		provingKeyMaxGasPrice.Uint64(),
		coordinatorJobSpecConfig,
		bhsJobSpecConfig,
		*simulationBlock,
		*coordinatorType,
		uint8(*optimismL1GasFeeCalculationMode),
		uint8(*optimismL1GasFeeCoefficient),
	)

	vrfPrimaryNode := nodesMap[model.VRFPrimaryNodeName]
	fmt.Println("Funding node's sending keys...")
	for _, sendingKey := range vrfPrimaryNode.SendingKeys {
		helpers.FundNode(e, sendingKey.Address, vrfPrimaryNode.SendingKeyFundingAmount)
	}
}

func VRFV2PlusDeployUniverse(e helpers.Environment,
	subscriptionBalanceJuels *big.Int,
	subscriptionBalanceNativeWei *big.Int,
	vrfKeyRegistrationConfig model.VRFKeyRegistrationConfig,
	contractAddresses model.ContractAddresses,
	coordinatorConfig CoordinatorConfigV2Plus,
	nativeOnly bool,
	nodesMap map[string]model.Node,
	provingKeyMaxGasPrice uint64,
	coordinatorJobSpecConfig model.CoordinatorJobSpecConfig,
	bhsJobSpecConfig model.BHSJobSpecConfig,
	simulationBlock string,
	coordinatorType string,
	optimismL1FeeMode uint8,
	optimismL1FeeCoefficient uint8,
) model.JobSpecs {
	var compressedPkHex string
	var keyHash common.Hash
	if len(vrfKeyRegistrationConfig.VRFKeyUncompressedPubKey) > 0 {
		// Put key in ECDSA format
		if strings.HasPrefix(vrfKeyRegistrationConfig.VRFKeyUncompressedPubKey, "0x") {
			vrfKeyRegistrationConfig.VRFKeyUncompressedPubKey = strings.Replace(vrfKeyRegistrationConfig.VRFKeyUncompressedPubKey, "0x", "04", 1)
		}

		// Generate compressed public key and key hash
		pubBytes, err := hex.DecodeString(vrfKeyRegistrationConfig.VRFKeyUncompressedPubKey)
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

		compressedPkHex = hexutil.Encode(pkBytes)
		keyHash, err = newPK.Hash()
		helpers.PanicErr(err)
	}

	if !nativeOnly && len(contractAddresses.LinkAddress) == 0 {
		fmt.Println("\nDeploying LINK Token...")
		contractAddresses.LinkAddress = helpers.DeployLinkToken(e).String()
	}

	if !nativeOnly && len(contractAddresses.LinkEthAddress) == 0 {
		fmt.Println("\nDeploying LINK/Native Feed...")
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
		fmt.Printf("\nDeploying Coordinator [type=%s]...\n", coordinatorType)
		contractAddresses.CoordinatorAddress = DeployCoordinator(e, contractAddresses.LinkAddress, contractAddresses.BhsContractAddress.String(), contractAddresses.LinkEthAddress, coordinatorType)
	}

	coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(contractAddresses.CoordinatorAddress, e.Ec)
	helpers.PanicErr(err)

	if contractAddresses.BatchCoordinatorAddress.String() == "0x0000000000000000000000000000000000000000" {
		fmt.Println("\nDeploying Batch Coordinator...")
		contractAddresses.BatchCoordinatorAddress = DeployBatchCoordinatorV2(e, contractAddresses.CoordinatorAddress)
	}

	fmt.Println("\nSetting Coordinator Config...")
	SetCoordinatorConfig(
		e,
		*coordinator,
		uint16(coordinatorConfig.MinConfs),
		uint32(coordinatorConfig.MaxGasLimit),
		uint32(coordinatorConfig.StalenessSeconds),
		uint32(coordinatorConfig.GasAfterPayment),
		coordinatorConfig.FallbackWeiPerUnitLink,
		coordinatorConfig.FulfillmentFlatFeeNativePPM,
		coordinatorConfig.FulfillmentFlatFeeLinkDiscountPPM,
		coordinatorConfig.NativePremiumPercentage,
		coordinatorConfig.LinkPremiumPercentage,
	)

	if coordinatorType == "optimism" {
		fmt.Println("\nSetting L1 gas fee calculation...")
		SetCoordinatorL1FeeCalculation(e, contractAddresses.CoordinatorAddress, optimismL1FeeMode, optimismL1FeeCoefficient)
	}

	fmt.Println("\nConfig set, getting current config from deployed contract...")
	PrintCoordinatorConfig(coordinator)

	if len(vrfKeyRegistrationConfig.VRFKeyUncompressedPubKey) > 0 {
		fmt.Println("\nRegistering proving key...")

		//NOTE - register proving key against EOA account, and not against Oracle's sending address in other to be able
		// easily withdraw funds from Coordinator contract back to EOA account
		RegisterCoordinatorProvingKey(e, *coordinator, vrfKeyRegistrationConfig.VRFKeyUncompressedPubKey, provingKeyMaxGasPrice)

		fmt.Println("\nProving key registered, getting proving key hashes from deployed contract...")
		registerdKeyHash, err2 := coordinator.SProvingKeyHashes(nil, big.NewInt(0))
		helpers.PanicErr(err2)
		fmt.Println("Key hash registered:", hex.EncodeToString(registerdKeyHash[:]))
	} else {
		fmt.Println("NOT registering proving key - you must do this eventually in order to fully deploy VRF!")
	}

	fmt.Println("\nDeploying consumer...")
	consumerAddress := EoaV2PlusLoadTestConsumerWithMetricsDeploy(e, contractAddresses.CoordinatorAddress.String())

	fmt.Println("\nAdding subscription...")
	subID, err := EoaCreateSub(e, *coordinator)
	helpers.PanicErr(err)

	fmt.Println("\nAdding consumer to subscription...")
	EoaAddConsumerToSub(e, *coordinator, subID, consumerAddress.String())

	if subscriptionBalanceJuels.Cmp(big.NewInt(0)) > 0 {
		fmt.Println("\nFunding subscription with Link Token.", subscriptionBalanceJuels, "juels...")
		EoaFundSubWithLink(e, *coordinator, contractAddresses.LinkAddress, subscriptionBalanceJuels, subID)
	} else {
		fmt.Println("Subscription", subID, "NOT getting funded with Link Token. You must fund the subscription in order to use it!")
	}
	if subscriptionBalanceNativeWei.Cmp(big.NewInt(0)) > 0 {
		fmt.Println("\nFunding subscription with Native Token.", subscriptionBalanceNativeWei, "wei...")
		EoaFundSubWithNative(e, coordinator.Address(), subID, subscriptionBalanceNativeWei)
	} else {
		fmt.Println("Subscription", subID, "NOT getting funded with Native Token. You must fund the subscription in order to use it!")
	}

	fmt.Println("\nSubscribed and (possibly) funded, retrieving subscription from deployed contract...")
	s, err := coordinator.GetSubscription(nil, subID)
	helpers.PanicErr(err)
	fmt.Printf("Subscription %+v\n", s)

	formattedVrfV2PlusPrimaryJobSpec := fmt.Sprintf(
		jobs.VRFV2PlusJobFormatted,
		contractAddresses.CoordinatorAddress,                   //coordinatorAddress
		contractAddresses.BatchCoordinatorAddress,              //batchCoordinatorAddress
		coordinatorJobSpecConfig.BatchFulfillmentEnabled,       //batchFulfillmentEnabled
		coordinatorJobSpecConfig.BatchFulfillmentGasMultiplier, //batchFulfillmentGasMultiplier
		compressedPkHex,            //publicKey
		coordinatorConfig.MinConfs, //minIncomingConfirmations
		e.ChainID,                  //evmChainID
		strings.Join(util.MapToAddressArr(nodesMap[model.VRFPrimaryNodeName].SendingKeys), "\",\""), //fromAddresses
		coordinatorJobSpecConfig.PollPeriod,     //pollPeriod
		coordinatorJobSpecConfig.RequestTimeout, //requestTimeout
		contractAddresses.CoordinatorAddress,
		coordinatorJobSpecConfig.EstimateGasMultiplier, //estimateGasMultiplier
		simulationBlock,
		func() string {
			if keys := nodesMap[model.VRFPrimaryNodeName].SendingKeys; len(keys) > 0 {
				return keys[0].Address
			}
			return common.HexToAddress("0x0").String()
		}(),
		contractAddresses.CoordinatorAddress,
		contractAddresses.CoordinatorAddress,
		simulationBlock,
	)

	formattedVrfV2PlusBackupJobSpec := fmt.Sprintf(
		jobs.VRFV2PlusJobFormatted,
		contractAddresses.CoordinatorAddress,                   //coordinatorAddress
		contractAddresses.BatchCoordinatorAddress,              //batchCoordinatorAddress
		coordinatorJobSpecConfig.BatchFulfillmentEnabled,       //batchFulfillmentEnabled
		coordinatorJobSpecConfig.BatchFulfillmentGasMultiplier, //batchFulfillmentGasMultiplier
		compressedPkHex, //publicKey
		100,             //minIncomingConfirmations
		e.ChainID,       //evmChainID
		strings.Join(util.MapToAddressArr(nodesMap[model.VRFBackupNodeName].SendingKeys), "\",\""), //fromAddresses
		coordinatorJobSpecConfig.PollPeriod,     //pollPeriod
		coordinatorJobSpecConfig.RequestTimeout, //requestTimeout
		contractAddresses.CoordinatorAddress,
		coordinatorJobSpecConfig.EstimateGasMultiplier, //estimateGasMultiplier
		simulationBlock,
		func() string {
			if keys := nodesMap[model.VRFPrimaryNodeName].SendingKeys; len(keys) > 0 {
				return keys[0].Address
			}
			return common.HexToAddress("0x0").String()
		}(),
		contractAddresses.CoordinatorAddress,
		contractAddresses.CoordinatorAddress,
		simulationBlock,
	)

	formattedBHSJobSpec := fmt.Sprintf(
		jobs.BHSPlusJobFormatted,
		contractAddresses.CoordinatorAddress, //coordinatorAddress
		contractAddresses.CoordinatorAddress, //coordinatorAddress
		bhsJobSpecConfig.WaitBlocks,          //waitBlocks
		bhsJobSpecConfig.LookBackBlocks,      //lookbackBlocks
		contractAddresses.BhsContractAddress, //bhs address
		bhsJobSpecConfig.PollPeriod,          //pollPeriod
		bhsJobSpecConfig.RunTimeout,          //runTimeout
		e.ChainID,                            //chain id
		strings.Join(util.MapToAddressArr(nodesMap[model.BHSNodeName].SendingKeys), "\",\""), //sending addresses
	)

	formattedBHSBackupJobSpec := fmt.Sprintf(
		jobs.BHSPlusJobFormatted,
		contractAddresses.CoordinatorAddress, //coordinatorAddress
		contractAddresses.CoordinatorAddress, //coordinatorAddress
		100,                                  //waitBlocks
		200,                                  //lookbackBlocks
		contractAddresses.BhsContractAddress, //bhs adreess
		bhsJobSpecConfig.PollPeriod,          //pollPeriod
		bhsJobSpecConfig.RunTimeout,          //runTimeout
		e.ChainID,                            //chain id
		strings.Join(util.MapToAddressArr(nodesMap[model.BHSBackupNodeName].SendingKeys), "\",\""), //sending addresses
	)

	formattedBHFJobSpec := fmt.Sprintf(
		jobs.BHFPlusJobFormatted,
		contractAddresses.CoordinatorAddress, //coordinatorAddress
		contractAddresses.BhsContractAddress, //bhs adreess
		contractAddresses.BatchBHSAddress,    //batchBHS
		e.ChainID,                            //chain id
		strings.Join(util.MapToAddressArr(nodesMap[model.BHFNodeName].SendingKeys), "\",\""), //sending addresses
	)

	fmt.Println(
		"\nDeployment complete.",
		"\nLINK Token contract address:", contractAddresses.LinkAddress,
		"\nLINK/Native Feed contract address:", contractAddresses.LinkEthAddress,
		"\nBlockhash Store contract address:", contractAddresses.BhsContractAddress,
		"\nBatch Blockhash Store contract address:", contractAddresses.BatchBHSAddress,
		"\nVRF Coordinator Address:", contractAddresses.CoordinatorAddress,
		"\nBatch VRF Coordinator Address:", contractAddresses.BatchCoordinatorAddress,
		"\nVRF Consumer Address:", consumerAddress,
		"\nVRF Subscription Id:", subID,
		"\nVRF Subscription LINK Balance:", *subscriptionBalanceJuels,
		"\nVRF Subscription Native Balance:", *subscriptionBalanceNativeWei,
		"\nPossible VRF Request command: ",
		fmt.Sprintf("go run . eoa-load-test-request-with-metrics --consumer-address=%s --sub-id=%d --key-hash=%s --request-confirmations %d --native-payment-enabled=true --requests 1 --runs 1 --cb-gas-limit 1_000_000", consumerAddress, subID, keyHash, coordinatorConfig.MinConfs),
		"\nRetrieve Request Status: ",
		fmt.Sprintf("go run . eoa-load-test-read-metrics --consumer-address=%s", consumerAddress),
		"\nA node can now be configured to run a VRF job with the below job spec :\n",
		formattedVrfV2PlusPrimaryJobSpec,
	)

	return model.JobSpecs{
		VRFPrimaryNode: formattedVrfV2PlusPrimaryJobSpec,
		VRFBackupyNode: formattedVrfV2PlusBackupJobSpec,
		BHSNode:        formattedBHSJobSpec,
		BHSBackupNode:  formattedBHSBackupJobSpec,
		BHFNode:        formattedBHFJobSpec,
	}
}

func DeployWrapperUniverse(e helpers.Environment) {
	cmd := flag.NewFlagSet("wrapper-universe-deploy", flag.ExitOnError)
	wrapperType := cmd.String("wrapper-type", "", "Specify which wrapper type to use: layer1, arbitrum, optimism")
	linkAddress := cmd.String("link-address", "", "address of link token")
	linkNativeFeedAddress := cmd.String("link-native-feed", "", "address of link-native-feed")
	coordinatorAddress := cmd.String("coordinator-address", "", "address of the vrf coordinator v2plus contract")
	subscriptionID := cmd.String("subscription-id", "", "subscription ID for the wrapper")
	wrapperGasOverhead := cmd.Uint("wrapper-gas-overhead", 50_000, "amount of gas overhead in wrapper fulfillment")
	coordinatorGasOverheadNative := cmd.Uint("coordinator-gas-overhead-native", 52_000, "amount of gas overhead in coordinator fulfillment for native payment")
	coordinatorGasOverheadLink := cmd.Uint("coordinator-gas-overhead-link", 74_000, "amount of gas overhead in coordinator fulfillment for link payment")
	coordinatorGasOverheadPerWord := cmd.Uint("coordinator-gas-overhead-per-word", 0, "amount of gas overhead per word in coordinator fulfillment")
	wrapperNativePremiumPercentage := cmd.Uint("wrapper-native-premium-percentage", 25, "gas premium charged by wrapper for native payment")
	wrapperLinkPremiumPercentage := cmd.Uint("wrapper-link-premium-percentage", 25, "gas premium charged by wrapper for link payment")
	keyHash := cmd.String("key-hash", "", "the keyhash that wrapper requests should use")
	maxNumWords := cmd.Uint("max-num-words", 10, "the keyhash that wrapper requests should use")
	subFundingLink := cmd.String("sub-funding-link", "10000000000000000000", "amount in LINK to fund the subscription with")
	subFundingNative := cmd.String("sub-funding-native", "10000000000000000000", "amount in native to fund the subscription with")
	consumerFundingLink := cmd.String("consumer-funding-link", "10000000000000000000", "amount in LINK to fund the consumer with")
	consumerFundingNative := cmd.String("consumer-funding-native", "10000000000000000000", "amount in native to fund the consumer with")
	fallbackWeiPerUnitLink := cmd.String("fallback-wei-per-unit-link", "", "the fallback wei per unit link")
	stalenessSeconds := cmd.Uint("staleness-seconds", 86400, "the number of seconds of staleness to allow")
	fulfillmentFlatFeeNativePPM := cmd.Uint("fulfillment-flat-fee-native-ppm", 500, "the native flat fee in ppm to charge for fulfillment denominated in native")
	fulfillmentFlatFeeLinkDiscountPPM := cmd.Uint("fulfillment-flat-fee-link-discount-ppm", 500, "the link flat fee discount in ppm to charge for fulfillment denominated in native")
	// only necessary for Optimism coordinator contract
	optimismL1GasFeeCalculationMode := cmd.Uint64("optimism-l1-fee-mode", 0, "Choose Optimism coordinator contract L1 fee calculation mode: 0, 1, 2")
	optimismL1GasFeeCoefficient := cmd.Uint64("optimism-l1-fee-coefficient", 100, "Choose Optimism coordinator contract L1 fee coefficient percentage [1, 100]")
	helpers.ParseArgs(cmd, os.Args[2:], "wrapper-type", "link-address", "link-native-feed", "coordinator-address", "key-hash", "fallback-wei-per-unit-link")

	if *wrapperType != "layer1" && *wrapperType != "arbitrum" && *wrapperType != "optimism" {
		panic(fmt.Sprintf("Invalid Wrapper type `%s`. Only `layer1`, `arbitrum` and `optimism` are supported", *wrapperType))
	}

	subAmountLink, s := big.NewInt(0).SetString(*subFundingLink, 10)
	if !s {
		panic(fmt.Sprintf("failed to parse subscription top up amount '%s'", *subFundingLink))
	}

	subAmountNative, s := big.NewInt(0).SetString(*subFundingNative, 10)
	if !s {
		panic(fmt.Sprintf("failed to parse subscription top up amount '%s'", *subFundingNative))
	}

	consumerAmountLink, s := big.NewInt(0).SetString(*consumerFundingLink, 10)
	if !s {
		panic(fmt.Sprintf("failed to parse consumer top up amount '%s'", *consumerFundingLink))
	}

	consumerAmountNative, s := big.NewInt(0).SetString(*consumerFundingNative, 10)
	if !s {
		panic(fmt.Sprintf("failed to parse consumer top up amount '%s'", *consumerFundingNative))
	}

	coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(common.HexToAddress(*coordinatorAddress), e.Ec)
	helpers.PanicErr(err)

	var subId *big.Int
	if *subscriptionID == "" {
		subId, err = EoaCreateSub(e, *coordinator)
		helpers.PanicErr(err)
		fmt.Println("Created subscription ID:", subId)
	} else {
		subId = parseSubID(*subscriptionID)
		fmt.Println("Using existing subscription ID:", subId)
	}

	fmt.Println()

	wrapper := WrapperDeploy(e,
		common.HexToAddress(*linkAddress),
		common.HexToAddress(*linkNativeFeedAddress),
		common.HexToAddress(*coordinatorAddress),
		subId,
		*wrapperType,
	)

	fmt.Println("Deployed wrapper:", wrapper.String())
	fmt.Println("Wrapper type:", *wrapperType)
	fmt.Println()

	WrapperConfigure(e,
		wrapper,
		*wrapperGasOverhead,
		*coordinatorGasOverheadNative,
		*coordinatorGasOverheadLink,
		*coordinatorGasOverheadPerWord,
		*wrapperNativePremiumPercentage,
		*wrapperLinkPremiumPercentage,
		*keyHash,
		*maxNumWords,
		decimal.RequireFromString(*fallbackWeiPerUnitLink).BigInt(),
		uint32(*stalenessSeconds),
		uint32(*fulfillmentFlatFeeNativePPM),
		uint32(*fulfillmentFlatFeeLinkDiscountPPM),
	)

	fmt.Println("Configured wrapper")
	fmt.Println()

	if *wrapperType == "optimism" {
		WrapperSetL1FeeCalculation(e, wrapper, uint8(*optimismL1GasFeeCalculationMode), uint8(*optimismL1GasFeeCoefficient))
		fmt.Println("Set L1 gas fee calculation")
		fmt.Println()
	}

	consumer := WrapperConsumerDeploy(e,
		common.HexToAddress(*linkAddress),
		wrapper)

	fmt.Println("Deployed wrapper consumer:", consumer.String())
	fmt.Println()

	// for v2plus we need to add wrapper as a consumer to the subscription
	EoaAddConsumerToSub(e, *coordinator, subId, wrapper.String())

	fmt.Println("Added wrapper as the subscription consumer")
	fmt.Println()

	link, err := link_token_interface.NewLinkToken(common.HexToAddress(*linkAddress), e.Ec)
	helpers.PanicErr(err)

	tx, err := link.Transfer(e.Owner, consumer, consumerAmountLink)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "link transfer to consumer")

	sendNativeTokens(e, consumer, consumerAmountNative)

	fmt.Println("Funded wrapper consumer")
	fmt.Println()

	EoaFundSubWithLink(e, *coordinator, *linkAddress, subAmountLink, subId)
	// e.Owner.Value is hardcoded inside this helper function, make sure to run it as the last one in the script
	EoaFundSubWithNative(e, common.HexToAddress(*coordinatorAddress), subId, subAmountNative)

	fmt.Println("Funded wrapper subscription")
	fmt.Println()

	fmt.Println("Wrapper universe deployment complete")
	fmt.Println("Wrapper address:", wrapper.String())
	fmt.Println("Wrapper type:", *wrapperType)
	fmt.Println("Wrapper consumer address:", consumer.String())
	fmt.Println("Wrapper subscription ID:", subId)
	fmt.Printf("Send native request example: go run . wrapper-consumer-request --consumer-address=%s --cb-gas-limit=1000000 --native-payment=true\n", consumer.String())
	fmt.Printf("Send LINK request example: go run . wrapper-consumer-request --consumer-address=%s --cb-gas-limit=1000000 --native-payment=false\n", consumer.String())
}

func parseSubID(subID string) *big.Int {
	parsedSubID, ok := new(big.Int).SetString(subID, 10)
	if !ok {
		helpers.PanicErr(fmt.Errorf("sub ID %s cannot be parsed", subID))
	}
	return parsedSubID
}
