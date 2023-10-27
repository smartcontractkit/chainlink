package v2plusscripts

import (
	"bytes"
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/smartcontractkit/chainlink/core/scripts/common/vrf/constants"
	"github.com/smartcontractkit/chainlink/core/scripts/common/vrf/jobs"
	"github.com/smartcontractkit/chainlink/core/scripts/common/vrf/model"
	"github.com/smartcontractkit/chainlink/core/scripts/common/vrf/util"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2plus_interface"
	"math/big"
	"os"
	"strings"

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
	MinConfs               *int
	MaxGasLimit            *int64
	StalenessSeconds       *int64
	GasAfterPayment        *int64
	FallbackWeiPerUnitLink *big.Int
	FeeConfig              vrf_coordinator_v2_5.VRFCoordinatorV25FeeConfig
}

func SmokeTestVRF(e helpers.Environment) {
	smokeCmd := flag.NewFlagSet("smoke", flag.ExitOnError)

	// required flags
	linkAddress := smokeCmd.String("link-address", "", "address of link token")
	linkEthAddress := smokeCmd.String("link-eth-feed", "", "address of link eth feed")
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
	flatFeeLinkPPM := smokeCmd.Int64("flat-fee-link-ppm", 500, "fulfillment flat fee LINK ppm")
	flatFeeEthPPM := smokeCmd.Int64("flat-fee-eth-ppm", 500, "fulfillment flat fee ETH ppm")

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

	if len(*linkEthAddress) == 0 {
		fmt.Println("\nDeploying LINK/ETH Feed...")
		address := helpers.DeployLinkEthFeed(e, *linkAddress, fallbackWeiPerUnitLink).String()
		linkEthAddress = &address
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
		fmt.Println("\nDeploying Coordinator...")
		coordinatorAddress = DeployCoordinator(e, *linkAddress, bhsContractAddress.String(), *linkEthAddress)
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
			vrf_coordinator_v2_5.VRFCoordinatorV25FeeConfig{
				FulfillmentFlatFeeLinkPPM:   uint32(*flatFeeLinkPPM),
				FulfillmentFlatFeeNativePPM: uint32(*flatFeeEthPPM),
			},
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
	tx, err := coordinator.RegisterProvingKey(e.Owner, e.Owner.From, [2]*big.Int{x, y})
	helpers.PanicErr(err)
	registerReceipt := helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "register proving key on", coordinatorAddress.String())
	var provingKeyRegisteredLog *vrf_coordinator_v2_5.VRFCoordinatorV25ProvingKeyRegistered
	for _, log := range registerReceipt.Logs {
		if log.Address == coordinatorAddress {
			var err error
			provingKeyRegisteredLog, err = coordinator.ParseProvingKeyRegistered(*log)
			if err != nil {
				continue
			}
		}
	}
	if provingKeyRegisteredLog == nil {
		panic("no proving key registered log found")
	}
	if !bytes.Equal(provingKeyRegisteredLog.KeyHash[:], keyHash[:]) {
		panic(fmt.Sprintf("unexpected key hash registered %s, expected %s", hexutil.Encode(provingKeyRegisteredLog.KeyHash[:]), hexutil.Encode(keyHash[:])))
	} else {
		fmt.Println("key hash registered:", hexutil.Encode(provingKeyRegisteredLog.KeyHash[:]))
	}

	fmt.Println("\nProving key registered, getting proving key hashes from deployed contract...")
	_, _, provingKeyHashes, configErr := coordinator.GetRequestConfig(nil)
	helpers.PanicErr(configErr)
	fmt.Println("Key hash registered:", hexutil.Encode(provingKeyHashes[len(provingKeyHashes)-1][:]))
	ourKeyHash := key.PublicKey.MustHash()
	if !bytes.Equal(provingKeyHashes[len(provingKeyHashes)-1][:], ourKeyHash[:]) {
		panic(fmt.Sprintf("unexpected key hash %s, expected %s", hexutil.Encode(provingKeyHashes[len(provingKeyHashes)-1][:]), hexutil.Encode(ourKeyHash[:])))
	}

	fmt.Println("\nDeploying consumer...")
	consumerAddress := EoaDeployConsumer(e, coordinatorAddress.String(), *linkAddress)

	fmt.Println("\nAdding subscription...")
	EoaCreateSub(e, *coordinator)

	subID := FindSubscriptionID(e, coordinator)
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
		"\nLINK/ETH Feed contract address:", *linkEthAddress,
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
	receipt := helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "request random words from", consumerAddress.String())
	fmt.Println("request blockhash:", receipt.BlockHash)

	// extract the RandomWordsRequested log from the receipt logs
	var rwrLog *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsRequested
	for _, log := range receipt.Logs {
		if log.Address == coordinatorAddress {
			var err error
			rwrLog, err = coordinator.ParseRandomWordsRequested(*log)
			if err != nil {
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
	} else {
		fmt.Println("storeEarliest succeeded, checking BH is there")
		bh, err := bhs.GetBlockhash(nil, seReceipt.BlockNumber.Sub(seReceipt.BlockNumber, big.NewInt(256)))
		helpers.PanicErr(err)
		fmt.Println("blockhash stored by storeEarliest:", hexutil.Encode(bh[:]))
		anchorBlockNumber = seReceipt.BlockNumber
	}
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
	} else {
		fmt.Println("store succeeded, checking BH is there")
		bh, err := bhs.GetBlockhash(nil, toStore)
		helpers.PanicErr(err)
		fmt.Println("blockhash stored by store:", hexutil.Encode(bh[:]))
	}

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
	} else {
		fmt.Println("storeVerifyHeader succeeded, checking BH is there")
		bh, err := bhs.GetBlockhash(nil, toStore)
		helpers.PanicErr(err)
		fmt.Println("blockhash stored by storeVerifyHeader:", hexutil.Encode(bh[:]))
	}
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

func DeployUniverseViaCLI(e helpers.Environment) {
	deployCmd := flag.NewFlagSet("deploy-universe", flag.ExitOnError)

	// required flags
	linkAddress := *deployCmd.String("link-address", "", "address of link token")
	linkEthAddress := *deployCmd.String("link-eth-feed", "", "address of link eth feed")
	bhsContractAddressString := *deployCmd.String("bhs-address", "", "address of BHS contract")
	batchBHSAddressString := *deployCmd.String("batch-bhs-address", "", "address of Batch BHS contract")
	coordinatorAddressString := *deployCmd.String("coordinator-address", "", "address of VRF Coordinator contract")
	batchCoordinatorAddressString := *deployCmd.String("batch-coordinator-address", "", "address Batch VRF Coordinator contract")
	subscriptionBalanceJuelsString := deployCmd.String("subscription-balance", "1e19", "amount to fund subscription with Link token (Juels)")
	subscriptionBalanceNativeWeiString := deployCmd.String("subscription-balance-native", "1e18", "amount to fund subscription with native token (Wei)")

	batchFulfillmentEnabled := deployCmd.Bool("batch-fulfillment-enabled", constants.BatchFulfillmentEnabled, "whether send randomness fulfillments in batches inside one tx from CL node")

	// optional flags
	fallbackWeiPerUnitLinkString := deployCmd.String("fallback-wei-per-unit-link", "6e16", "fallback wei/link ratio")
	registerKeyUncompressedPubKey := deployCmd.String("uncompressed-pub-key", "", "uncompressed public key")
	vrfPrimaryNodeSendingKeysString := deployCmd.String("vrf-primary-node-sending-keys", "", "VRF Primary Node sending keys")
	minConfs := deployCmd.Int("min-confs", constants.MinConfs, "min confs")
	nodeSendingKeyFundingAmount := deployCmd.String("sending-key-funding-amount", constants.NodeSendingKeyFundingAmount, "CL node sending key funding amount")
	maxGasLimit := deployCmd.Int64("max-gas-limit", constants.MaxGasLimit, "max gas limit")
	stalenessSeconds := deployCmd.Int64("staleness-seconds", constants.StalenessSeconds, "staleness in seconds")
	gasAfterPayment := deployCmd.Int64("gas-after-payment", constants.GasAfterPayment, "gas after payment calculation")
	flatFeeLinkPPM := deployCmd.Int64("flat-fee-link-ppm", constants.FlatFeeLinkPPM, "fulfillment flat fee LINK ppm")
	flatFeeEthPPM := deployCmd.Int64("flat-fee-eth-ppm", constants.FlatFeeNativePPM, "fulfillment flat fee ETH ppm")

	helpers.ParseArgs(
		deployCmd, os.Args[2:],
	)

	fallbackWeiPerUnitLink := decimal.RequireFromString(*fallbackWeiPerUnitLinkString).BigInt()
	subscriptionBalanceJuels := decimal.RequireFromString(*subscriptionBalanceJuelsString).BigInt()
	subscriptionBalanceNativeWei := decimal.RequireFromString(*subscriptionBalanceNativeWeiString).BigInt()
	fundingAmount := decimal.RequireFromString(*nodeSendingKeyFundingAmount).BigInt()

	feeConfig := vrf_coordinator_v2_5.VRFCoordinatorV25FeeConfig{
		FulfillmentFlatFeeLinkPPM:   uint32(*flatFeeLinkPPM),
		FulfillmentFlatFeeNativePPM: uint32(*flatFeeEthPPM),
	}

	vrfPrimaryNodeSendingKeys := strings.Split(*vrfPrimaryNodeSendingKeysString, ",")

	nodesMap := make(map[string]model.Node)

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

	coordinatorConfig := CoordinatorConfigV2Plus{
		MinConfs:               minConfs,
		MaxGasLimit:            maxGasLimit,
		StalenessSeconds:       stalenessSeconds,
		GasAfterPayment:        gasAfterPayment,
		FallbackWeiPerUnitLink: fallbackWeiPerUnitLink,
		FeeConfig:              feeConfig,
	}

	VRFV2PlusDeployUniverse(
		e,
		subscriptionBalanceJuels,
		subscriptionBalanceNativeWei,
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

func VRFV2PlusDeployUniverse(e helpers.Environment,
	subscriptionBalanceJuels *big.Int,
	subscriptionBalanceNativeWei *big.Int,
	registerKeyUncompressedPubKey *string,
	contractAddresses model.ContractAddresses,
	coordinatorConfig CoordinatorConfigV2Plus,
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
	consumerAddress := EoaV2PlusLoadTestConsumerWithMetricsDeploy(e, contractAddresses.CoordinatorAddress.String())

	fmt.Println("\nAdding subscription...")
	EoaCreateSub(e, *coordinator)

	subID := FindSubscriptionID(e, coordinator)
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

	formattedVrfV2PlusBackupJobSpec := fmt.Sprintf(
		jobs.VRFV2PlusJobFormatted,
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
		"\nVRF Subscription LINK Balance:", *subscriptionBalanceJuels,
		"\nVRF Subscription Native Balance:", *subscriptionBalanceNativeWei,
		"\nPossible VRF Request command: ",
		fmt.Sprintf("go run . eoa-load-test-request-with-metrics --consumer-address=%s --sub-id=%d --key-hash=%s --request-confirmations %d --requests 1 --runs 1 --cb-gas-limit 1_000_000", consumerAddress, subID, keyHash, *coordinatorConfig.MinConfs),
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
	fallbackWeiPerUnitLink := cmd.String("fallback-wei-per-unit-link", "", "the fallback wei per unit link")
	stalenessSeconds := cmd.Uint("staleness-seconds", 86400, "the number of seconds of staleness to allow")
	fulfillmentFlatFeeLinkPPM := cmd.Uint("fulfillment-flat-fee-link-ppm", 500, "the link flat fee in ppm to charge for fulfillment")
	fulfillmentFlatFeeNativePPM := cmd.Uint("fulfillment-flat-fee-native-ppm", 500, "the native flat fee in ppm to charge for fulfillment")
	helpers.ParseArgs(cmd, os.Args[2:], "link-address", "link-eth-feed", "coordinator-address", "key-hash", "fallback-wei-per-unit-link")

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
		*maxNumWords,
		decimal.RequireFromString(*fallbackWeiPerUnitLink).BigInt(),
		uint32(*stalenessSeconds),
		uint32(*fulfillmentFlatFeeLinkPPM),
		uint32(*fulfillmentFlatFeeNativePPM),
	)

	consumer := WrapperConsumerDeploy(e,
		common.HexToAddress(*linkAddress),
		wrapper)

	coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(common.HexToAddress(*coordinatorAddress), e.Ec)
	helpers.PanicErr(err)

	EoaFundSubWithLink(e, *coordinator, *linkAddress, amount, subID)

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
