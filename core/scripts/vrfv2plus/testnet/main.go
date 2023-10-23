package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/smartcontractkit/chainlink/core/scripts/vrfv2plus/testnet/v2plusscripts"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/chain_specific_util_helper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2plus_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_load_test_with_metrics"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/sqlx"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_vrf_coordinator_v2plus"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/trusted_blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_load_test_external_sub_owner"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_single_consumer"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_sub_owner"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2plus_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2plus_wrapper_consumer_example"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/extraargs"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/proof"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	batchCoordinatorV2PlusABI = evmtypes.MustGetABI(batch_vrf_coordinator_v2plus.BatchVRFCoordinatorV2PlusABI)
)

func main() {
	e := helpers.SetupEnv(false)

	switch os.Args[1] {
	case "csu-deploy":
		addr, tx, _, err := chain_specific_util_helper.DeployChainSpecificUtilHelper(e.Owner, e.Ec)
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "deploying chain specific util helper")
		fmt.Println("deployed chain specific util helper at:", addr)
	case "csu-block-number":
		cmd := flag.NewFlagSet("csu-block-number", flag.ExitOnError)
		csuAddress := cmd.String("csu-address", "", "address of the chain specific util helper contract")
		helpers.ParseArgs(cmd, os.Args[2:], "csu-address")
		csu, err := chain_specific_util_helper.NewChainSpecificUtilHelper(common.HexToAddress(*csuAddress), e.Ec)
		helpers.PanicErr(err)
		blockNumber, err := csu.GetBlockNumber(nil)
		helpers.PanicErr(err)
		fmt.Println("block number:", blockNumber)
	case "csu-block-hash":
		cmd := flag.NewFlagSet("csu-block-hash", flag.ExitOnError)
		csuAddress := cmd.String("csu-address", "", "address of the chain specific util helper contract")
		blockNumber := cmd.Uint64("block-number", 0, "block number to get the hash of")
		helpers.ParseArgs(cmd, os.Args[2:], "csu-address")
		csu, err := chain_specific_util_helper.NewChainSpecificUtilHelper(common.HexToAddress(*csuAddress), e.Ec)
		helpers.PanicErr(err)
		blockHash, err := csu.GetBlockhash(nil, *blockNumber)
		helpers.PanicErr(err)
		fmt.Println("block hash:", hexutil.Encode(blockHash[:]))
	case "csu-current-tx-l1-gas-fees":
		cmd := flag.NewFlagSet("csu-current-tx-l1-gas-fees", flag.ExitOnError)
		csuAddress := cmd.String("csu-address", "", "address of the chain specific util helper contract")
		calldata := cmd.String("calldata", "", "calldata to estimate gas fees for")
		helpers.ParseArgs(cmd, os.Args[2:], "csu-address", "calldata")
		csu, err := chain_specific_util_helper.NewChainSpecificUtilHelper(common.HexToAddress(*csuAddress), e.Ec)
		helpers.PanicErr(err)
		gasFees, err := csu.GetCurrentTxL1GasFees(nil, *calldata)
		helpers.PanicErr(err)
		fmt.Println("gas fees:", gasFees)
	case "csu-l1-calldata-gas-cost":
		cmd := flag.NewFlagSet("csu-l1-calldata-gas-cost", flag.ExitOnError)
		csuAddress := cmd.String("csu-address", "", "address of the chain specific util helper contract")
		calldataSize := cmd.String("calldata-size", "", "size of the calldata to estimate gas fees for")
		helpers.ParseArgs(cmd, os.Args[2:], "csu-address", "calldata-size")
		csu, err := chain_specific_util_helper.NewChainSpecificUtilHelper(common.HexToAddress(*csuAddress), e.Ec)
		helpers.PanicErr(err)
		gasCost, err := csu.GetL1CalldataGasCost(nil, decimal.RequireFromString(*calldataSize).BigInt())
		helpers.PanicErr(err)
		fmt.Println("gas cost:", gasCost)
	case "smoke":
		v2plusscripts.SmokeTestVRF(e)
	case "smoke-bhs":
		v2plusscripts.SmokeTestBHS(e)
	case "manual-fulfill":
		cmd := flag.NewFlagSet("manual-fulfill", flag.ExitOnError)
		// In order to get the tx data for a fulfillment transaction, you can grep the
		// chainlink node logs for the VRF v2 request ID in hex. You will find a log for
		// the vrf task in the VRF pipeline, specifically the "output" log field.
		// Sample Loki query:
		// {app="app-name"} | json | taskType="vrfv2plus" |~ "39f2d812c04e07cb9c71e93ce6547e48b7dd23ed4cc02616dfef5ef063a58bde"
		txdatas := cmd.String("txdatas", "", "hex encoded tx data")
		coordinatorAddress := cmd.String("coordinator-address", "", "coordinator address")
		gasMultiplier := cmd.Float64("gas-multiplier", 1.1, "gas multiplier")
		helpers.ParseArgs(cmd, os.Args[2:], "txdatas", "coordinator-address")
		txdatasParsed := helpers.ParseHexSlice(*txdatas)
		coordinatorAddr := common.HexToAddress(*coordinatorAddress)
		for i, txdata := range txdatasParsed {
			nonce, err := e.Ec.PendingNonceAt(context.Background(), e.Owner.From)
			helpers.PanicErr(err)
			estimate, err := e.Ec.EstimateGas(context.Background(), ethereum.CallMsg{
				From: common.HexToAddress("0x0"),
				To:   &coordinatorAddr,
				Data: txdata,
			})
			helpers.PanicErr(err)
			finalEstimate := uint64(*gasMultiplier * float64(estimate))
			tx := types.NewTx(&types.LegacyTx{
				Nonce:    nonce,
				GasPrice: e.Owner.GasPrice,
				Gas:      finalEstimate,
				To:       &coordinatorAddr,
				Data:     txdata,
			})
			signedTx, err := e.Owner.Signer(e.Owner.From, tx)
			helpers.PanicErr(err)
			err = e.Ec.SendTransaction(context.Background(), signedTx)
			helpers.PanicErr(err)
			helpers.ConfirmTXMined(context.Background(), e.Ec, signedTx, e.ChainID, fmt.Sprintf("manual fulfillment %d", i+1))
		}
	case "topics":
		randomWordsRequested := vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsRequested{}.Topic()
		randomWordsFulfilled := vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled{}.Topic()
		fmt.Println("RandomWordsRequested:", randomWordsRequested.String(),
			"RandomWordsFulfilled:", randomWordsFulfilled.String())
	case "batch-coordinatorv2plus-deploy":
		cmd := flag.NewFlagSet("batch-coordinatorv2plus-deploy", flag.ExitOnError)
		coordinatorAddr := cmd.String("coordinator-address", "", "address of the vrf coordinator v2 contract")
		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address")
		_, tx, _, err := batch_vrf_coordinator_v2plus.DeployBatchVRFCoordinatorV2Plus(
			e.Owner, e.Ec, common.HexToAddress(*coordinatorAddr))
		helpers.PanicErr(err)
		helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
	case "batch-coordinatorv2plus-fulfill":
		cmd := flag.NewFlagSet("batch-coordinatorv2plus-fulfill", flag.ExitOnError)
		batchCoordinatorAddr := cmd.String("batch-coordinator-address", "", "address of the batch vrf coordinator v2 contract")
		pubKeyHex := cmd.String("pubkeyhex", "", "compressed pubkey hex")
		dbURL := cmd.String("db-url", "", "postgres database url")
		keystorePassword := cmd.String("keystore-pw", "", "password to the keystore")
		submit := cmd.Bool("submit", false, "whether to submit the fulfillments or not")
		estimateGas := cmd.Bool("estimate-gas", false, "whether to estimate gas or not")
		nativePayment := cmd.Bool("native-payment", false, "whether to use native payment or not")

		// NOTE: it is assumed all of these are of the same length and that
		// elements correspond to each other index-wise. this property is not checked.
		preSeeds := cmd.String("preseeds", "", "comma-separated request preSeeds")
		blockHashes := cmd.String("blockhashes", "", "comma-separated request blockhashes")
		blockNums := cmd.String("blocknums", "", "comma-separated request blocknumbers")
		subIDs := cmd.String("subids", "", "comma-separated request subids")
		cbGasLimits := cmd.String("cbgaslimits", "", "comma-separated request callback gas limits")
		numWordses := cmd.String("numwordses", "", "comma-separated request num words")
		senders := cmd.String("senders", "", "comma-separated request senders")

		helpers.ParseArgs(cmd, os.Args[2:],
			"batch-coordinator-address", "pubkeyhex", "db-url",
			"keystore-pw", "preseeds", "blockhashes", "blocknums",
			"subids", "cbgaslimits", "numwordses", "senders", "submit",
		)

		preSeedSlice := helpers.ParseBigIntSlice(*preSeeds)
		bhSlice := helpers.ParseHashSlice(*blockHashes)
		blockNumSlice := helpers.ParseBigIntSlice(*blockNums)
		subIDSlice := helpers.ParseBigIntSlice(*subIDs)
		cbLimitsSlice := helpers.ParseBigIntSlice(*cbGasLimits)
		numWordsSlice := helpers.ParseBigIntSlice(*numWordses)
		senderSlice := helpers.ParseAddressSlice(*senders)

		batchCoordinator, err := batch_vrf_coordinator_v2plus.NewBatchVRFCoordinatorV2Plus(common.HexToAddress(*batchCoordinatorAddr), e.Ec)
		helpers.PanicErr(err)

		db := sqlx.MustOpen("postgres", *dbURL)
		lggr, _ := logger.NewLogger()

		keyStore := keystore.New(db, utils.DefaultScryptParams, lggr, pg.NewQConfig(false))
		err = keyStore.Unlock(*keystorePassword)
		helpers.PanicErr(err)

		k, err := keyStore.VRF().Get(*pubKeyHex)
		helpers.PanicErr(err)

		fmt.Println("vrf key found:", k)

		proofs := []batch_vrf_coordinator_v2plus.VRFTypesProof{}
		reqCommits := []batch_vrf_coordinator_v2plus.VRFTypesRequestCommitmentV2Plus{}
		for i := range preSeedSlice {
			ps, err := proof.BigToSeed(preSeedSlice[i])
			helpers.PanicErr(err)
			extraArgs, err := extraargs.ExtraArgsV1(*nativePayment)
			helpers.PanicErr(err)
			preSeedData := proof.PreSeedDataV2Plus{
				PreSeed:          ps,
				BlockHash:        bhSlice[i],
				BlockNum:         blockNumSlice[i].Uint64(),
				SubId:            subIDSlice[i],
				CallbackGasLimit: uint32(cbLimitsSlice[i].Uint64()),
				NumWords:         uint32(numWordsSlice[i].Uint64()),
				Sender:           senderSlice[i],
				ExtraArgs:        extraArgs,
			}
			fmt.Printf("preseed data iteration %d: %+v\n", i, preSeedData)
			finalSeed := proof.FinalSeedV2Plus(preSeedData)

			p, err := keyStore.VRF().GenerateProof(*pubKeyHex, finalSeed)
			helpers.PanicErr(err)

			onChainProof, rc, err := proof.GenerateProofResponseFromProofV2Plus(p, preSeedData)
			helpers.PanicErr(err)

			proofs = append(proofs, batch_vrf_coordinator_v2plus.VRFTypesProof(onChainProof))
			reqCommits = append(reqCommits, batch_vrf_coordinator_v2plus.VRFTypesRequestCommitmentV2Plus(rc))
		}

		fmt.Printf("proofs: %+v\n\n", proofs)
		fmt.Printf("request commitments: %+v\n\n", reqCommits)

		if *submit {
			fmt.Println("submitting fulfillments...")
			tx, err := batchCoordinator.FulfillRandomWords(e.Owner, proofs, reqCommits)
			helpers.PanicErr(err)

			fmt.Println("waiting for it to mine:", helpers.ExplorerLink(e.ChainID, tx.Hash()))
			_, err = bind.WaitMined(context.Background(), e.Ec, tx)
			helpers.PanicErr(err)
			fmt.Println("done")
		}

		if *estimateGas {
			fmt.Println("estimating gas")
			payload, err := batchCoordinatorV2PlusABI.Pack("fulfillRandomWords", proofs, reqCommits)
			helpers.PanicErr(err)

			a := batchCoordinator.Address()
			gasEstimate, err := e.Ec.EstimateGas(context.Background(), ethereum.CallMsg{
				From: e.Owner.From,
				To:   &a,
				Data: payload,
			})
			helpers.PanicErr(err)

			fmt.Println("gas estimate:", gasEstimate)
		}
	case "coordinatorv2-fulfill":
		cmd := flag.NewFlagSet("coordinatorv2-fulfill", flag.ExitOnError)
		coordinatorAddr := cmd.String("coordinator-address", "", "address of the vrf coordinator v2 contract")
		pubKeyHex := cmd.String("pubkeyhex", "", "compressed pubkey hex")
		dbURL := cmd.String("db-url", "", "postgres database url")
		keystorePassword := cmd.String("keystore-pw", "", "password to the keystore")
		nativePayment := cmd.Bool("native-payment", false, "whether to use native payment or not")

		preSeed := cmd.String("preseed", "", "request preSeed")
		blockHash := cmd.String("blockhash", "", "request blockhash")
		blockNum := cmd.Uint64("blocknum", 0, "request blocknumber")
		subID := cmd.String("subid", "", "request subid")
		cbGasLimit := cmd.Uint("cbgaslimit", 0, "request callback gas limit")
		numWords := cmd.Uint("numwords", 0, "request num words")
		sender := cmd.String("sender", "", "request sender")

		helpers.ParseArgs(cmd, os.Args[2:],
			"coordinator-address", "pubkeyhex", "db-url",
			"keystore-pw", "preseed", "blockhash", "blocknum",
			"subid", "cbgaslimit", "numwords", "sender",
		)

		coordinator, err := vrf_coordinator_v2plus_interface.NewIVRFCoordinatorV2PlusInternal(common.HexToAddress(*coordinatorAddr), e.Ec)
		helpers.PanicErr(err)

		db := sqlx.MustOpen("postgres", *dbURL)
		lggr, _ := logger.NewLogger()

		keyStore := keystore.New(db, utils.DefaultScryptParams, lggr, pg.NewQConfig(false))
		err = keyStore.Unlock(*keystorePassword)
		helpers.PanicErr(err)

		k, err := keyStore.VRF().Get(*pubKeyHex)
		helpers.PanicErr(err)

		fmt.Println("vrf key found:", k)

		ps, err := proof.BigToSeed(decimal.RequireFromString(*preSeed).BigInt())
		helpers.PanicErr(err)

		parsedSubID := parseSubID(*subID)
		extraArgs, err := extraargs.ExtraArgsV1(*nativePayment)
		helpers.PanicErr(err)
		preSeedData := proof.PreSeedDataV2Plus{
			PreSeed:          ps,
			BlockHash:        common.HexToHash(*blockHash),
			BlockNum:         *blockNum,
			SubId:            parsedSubID,
			CallbackGasLimit: uint32(*cbGasLimit),
			NumWords:         uint32(*numWords),
			Sender:           common.HexToAddress(*sender),
			ExtraArgs:        extraArgs,
		}
		fmt.Printf("preseed data: %+v\n", preSeedData)
		finalSeed := proof.FinalSeedV2Plus(preSeedData)

		p, err := keyStore.VRF().GenerateProof(*pubKeyHex, finalSeed)
		helpers.PanicErr(err)

		onChainProof, rc, err := proof.GenerateProofResponseFromProofV2Plus(p, preSeedData)
		helpers.PanicErr(err)

		fmt.Printf("Proof: %+v, commitment: %+v\nSending fulfillment!", onChainProof, rc)

		tx, err := coordinator.FulfillRandomWords(e.Owner, onChainProof, rc)
		helpers.PanicErr(err)

		fmt.Println("waiting for it to mine:", helpers.ExplorerLink(e.ChainID, tx.Hash()))
		_, err = bind.WaitMined(context.Background(), e.Ec, tx)
		helpers.PanicErr(err)
		fmt.Println("done")
	case "batch-bhs-deploy":
		cmd := flag.NewFlagSet("batch-bhs-deploy", flag.ExitOnError)
		bhsAddr := cmd.String("bhs-address", "", "address of the blockhash store contract")
		helpers.ParseArgs(cmd, os.Args[2:], "bhs-address")
		v2plusscripts.DeployBatchBHS(e, common.HexToAddress(*bhsAddr))
	case "batch-bhs-store":
		cmd := flag.NewFlagSet("batch-bhs-store", flag.ExitOnError)
		batchAddr := cmd.String("batch-bhs-address", "", "address of the batch bhs contract")
		blockNumbersArg := cmd.String("block-numbers", "", "block numbers to store in a single transaction")
		helpers.ParseArgs(cmd, os.Args[2:], "batch-bhs-address", "block-numbers")
		batchBHS, err := batch_blockhash_store.NewBatchBlockhashStore(common.HexToAddress(*batchAddr), e.Ec)
		helpers.PanicErr(err)
		blockNumbers := helpers.ParseBigIntSlice(*blockNumbersArg)
		helpers.PanicErr(err)
		tx, err := batchBHS.Store(e.Owner, blockNumbers)
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
	case "batch-bhs-get":
		cmd := flag.NewFlagSet("batch-bhs-get", flag.ExitOnError)
		batchAddr := cmd.String("batch-bhs-address", "", "address of the batch bhs contract")
		blockNumbersArg := cmd.String("block-numbers", "", "block numbers to store in a single transaction")
		helpers.ParseArgs(cmd, os.Args[2:], "batch-bhs-address", "block-numbers")
		batchBHS, err := batch_blockhash_store.NewBatchBlockhashStore(common.HexToAddress(*batchAddr), e.Ec)
		helpers.PanicErr(err)
		blockNumbers := helpers.ParseBigIntSlice(*blockNumbersArg)
		helpers.PanicErr(err)
		blockhashes, err := batchBHS.GetBlockhashes(nil, blockNumbers)
		helpers.PanicErr(err)
		for i, bh := range blockhashes {
			fmt.Println("blockhash(", blockNumbers[i], ") = ", common.Bytes2Hex(bh[:]))
		}
	case "batch-bhs-storeVerify":
		cmd := flag.NewFlagSet("batch-bhs-storeVerify", flag.ExitOnError)
		batchAddr := cmd.String("batch-bhs-address", "", "address of the batch bhs contract")
		startBlock := cmd.Int64("start-block", -1, "block number to start from. Must be in the BHS already.")
		numBlocks := cmd.Int64("num-blocks", -1, "number of blockhashes to store. will be stored in a single tx, can't be > 150")
		helpers.ParseArgs(cmd, os.Args[2:], "batch-bhs-address", "start-block", "num-blocks")
		batchBHS, err := batch_blockhash_store.NewBatchBlockhashStore(common.HexToAddress(*batchAddr), e.Ec)
		helpers.PanicErr(err)
		blockRange, err := blockhashstore.DecreasingBlockRange(big.NewInt(*startBlock-1), big.NewInt(*startBlock-*numBlocks-1))
		helpers.PanicErr(err)
		rlpHeaders, _, err := helpers.GetRlpHeaders(e, blockRange, true)
		helpers.PanicErr(err)
		tx, err := batchBHS.StoreVerifyHeader(e.Owner, blockRange, rlpHeaders)
		helpers.PanicErr(err)
		fmt.Println("storeVerifyHeader(", blockRange, ", ...) tx:", helpers.ExplorerLink(e.ChainID, tx.Hash()))
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, fmt.Sprintf("blockRange: %d", blockRange))
	case "batch-bhs-backwards":
		cmd := flag.NewFlagSet("batch-bhs-backwards", flag.ExitOnError)
		batchAddr := cmd.String("batch-bhs-address", "", "address of the batch bhs contract")
		bhsAddr := cmd.String("bhs-address", "", "address of the bhs contract")
		startBlock := cmd.Int64("start-block", -1, "block number to start from. Must be in the BHS already.")
		endBlock := cmd.Int64("end-block", -1, "block number to end at. Must be less than startBlock")
		batchSize := cmd.Int64("batch-size", -1, "batch size")
		gasMultiplier := cmd.Int64("gas-price-multiplier", 1, "gas price multiplier to use, defaults to 1 (no multiplication)")
		helpers.ParseArgs(cmd, os.Args[2:], "batch-bhs-address", "bhs-address", "end-block", "batch-size")

		batchBHS, err := batch_blockhash_store.NewBatchBlockhashStore(common.HexToAddress(*batchAddr), e.Ec)
		helpers.PanicErr(err)

		bhs, err := blockhash_store.NewBlockhashStore(common.HexToAddress(*bhsAddr), e.Ec)
		helpers.PanicErr(err)

		// Sanity check BHS address in the Batch BHS.
		bhsAddressBatchBHS, err := batchBHS.BHS(nil)
		helpers.PanicErr(err)

		if bhsAddressBatchBHS != common.HexToAddress(*bhsAddr) {
			log.Panicf("Mismatch in bhs addresses: batch bhs has %s while given %s", bhsAddressBatchBHS.String(), *bhsAddr)
		}

		if *startBlock == -1 {
			tx, err2 := bhs.StoreEarliest(e.Owner)
			helpers.PanicErr(err2)
			receipt := helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "Store Earliest")
			// storeEarliest will store receipt block number minus 256 which is the earliest block
			// the blockhash() instruction will work on.
			*startBlock = receipt.BlockNumber.Int64() - 256
		}

		// Check if the provided start block is in the BHS. If it's not, print out an appropriate
		// helpful error message. Otherwise users would get the cryptic "header has unknown blockhash"
		// error which is a bit more difficult to diagnose.
		// The Batch BHS returns a zero'd [32]byte array in the event the provided block number doesn't
		// have it's blockhash in the BHS.
		var notFound [32]byte
		hsh, err := batchBHS.GetBlockhashes(nil, []*big.Int{big.NewInt(*startBlock)})
		helpers.PanicErr(err)

		if len(hsh) != 1 {
			helpers.PanicErr(fmt.Errorf("expected 1 item in returned array from BHS store, got: %d", len(hsh)))
		}

		if bytes.Equal(hsh[0][:], notFound[:]) {
			helpers.PanicErr(fmt.Errorf("expected block number %d (start-block argument) to be in the BHS already, did not find it there", *startBlock))
		}

		blockRange, err := blockhashstore.DecreasingBlockRange(big.NewInt(*startBlock-1), big.NewInt(*endBlock))
		helpers.PanicErr(err)

		for i := 0; i < len(blockRange); i += int(*batchSize) {
			j := i + int(*batchSize)
			if j > len(blockRange) {
				j = len(blockRange)
			}

			// Get suggested gas price and multiply by multiplier on every iteration
			// so we don't have our transaction getting stuck. Need to be as fast as
			// possible.
			gp, err := e.Ec.SuggestGasPrice(context.Background())
			helpers.PanicErr(err)
			e.Owner.GasPrice = new(big.Int).Mul(gp, big.NewInt(*gasMultiplier))

			fmt.Println("using gas price", e.Owner.GasPrice, "wei")

			blockNumbers := blockRange[i:j]
			blockHeaders, _, err := helpers.GetRlpHeaders(e, blockNumbers, true)
			fmt.Println("storing blockNumbers:", blockNumbers)
			helpers.PanicErr(err)

			tx, err := batchBHS.StoreVerifyHeader(e.Owner, blockNumbers, blockHeaders)
			helpers.PanicErr(err)

			fmt.Println("sent tx:", helpers.ExplorerLink(e.ChainID, tx.Hash()))

			fmt.Println("waiting for it to mine...")
			_, err = bind.WaitMined(context.Background(), e.Ec, tx)
			helpers.PanicErr(err)

			fmt.Println("received receipt, continuing")
		}
		fmt.Println("done")
	case "trusted-bhs-store":
		cmd := flag.NewFlagSet("trusted-bhs-backwards", flag.ExitOnError)
		trustedBHSAddr := cmd.String("trusted-bhs-address", "", "address of the trusted bhs contract")
		blockNumbersString := cmd.String("block-numbers", "", "comma-separated list of block numbers e.g 123,456 ")
		batchSizePtr := cmd.Int64("batch-size", -1, "batch size")
		helpers.ParseArgs(cmd, os.Args[2:], "trusted-bhs-address", "batch-size", "block-numbers")

		// Parse batch size.
		batchSize := int(*batchSizePtr)

		// Parse block numbers.
		blockNumbers := helpers.ParseBigIntSlice(*blockNumbersString)

		// Instantiate trusted bhs.
		trustedBHS, err := trusted_blockhash_store.NewTrustedBlockhashStore(common.HexToAddress(*trustedBHSAddr), e.Ec)
		helpers.PanicErr(err)

		for i := 0; i < len(blockNumbers); i += batchSize {
			// Get recent blockhash and block number anew each iteration. We do this so they do not get stale.
			recentBlockNumber, err := e.Ec.BlockNumber(context.Background())
			helpers.PanicErr(err)
			recentBlock, err := e.Ec.HeaderByNumber(context.Background(), big.NewInt(int64(recentBlockNumber)))
			helpers.PanicErr(err)
			recentBlockhash := recentBlock.Hash()

			// Get blockhashes to store.
			blockNumbersSlice := blockNumbers[i : i+batchSize]
			_, blockhashesStrings, err := helpers.GetRlpHeaders(e, blockNumbersSlice, false)
			helpers.PanicErr(err)
			fmt.Println("storing blockNumbers:", blockNumbers)
			var blockhashes [][32]byte
			for _, h := range blockhashesStrings {
				blockhashes = append(blockhashes, common.HexToHash(h))
			}

			// Execute storage tx.
			tx, err := trustedBHS.StoreTrusted(e.Owner, blockNumbersSlice, blockhashes, big.NewInt(int64(recentBlockNumber)), recentBlockhash)
			helpers.PanicErr(err)
			fmt.Println("waiting for it to mine...")
			helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
			fmt.Println("received receipt, continuing")
		}
		fmt.Println("done")
	case "latest-head":
		h, err := e.Ec.HeaderByNumber(context.Background(), nil)
		helpers.PanicErr(err)
		fmt.Println("latest head number:", h.Number.String())
	case "bhs-deploy":
		v2plusscripts.DeployBHS(e)
	case "coordinator-deploy":
		coordinatorDeployCmd := flag.NewFlagSet("coordinator-deploy", flag.ExitOnError)
		coordinatorDeployLinkAddress := coordinatorDeployCmd.String("link-address", "", "address of link token")
		coordinatorDeployBHSAddress := coordinatorDeployCmd.String("bhs-address", "", "address of bhs")
		coordinatorDeployLinkEthFeedAddress := coordinatorDeployCmd.String("link-eth-feed", "", "address of link-eth-feed")
		helpers.ParseArgs(coordinatorDeployCmd, os.Args[2:], "link-address", "bhs-address", "link-eth-feed")
		v2plusscripts.DeployCoordinator(e, *coordinatorDeployLinkAddress, *coordinatorDeployBHSAddress, *coordinatorDeployLinkEthFeedAddress)
	case "coordinator-get-config":
		cmd := flag.NewFlagSet("coordinator-get-config", flag.ExitOnError)
		coordinatorAddress := cmd.String("coordinator-address", "", "coordinator address")
		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address")

		coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(common.HexToAddress(*coordinatorAddress), e.Ec)
		helpers.PanicErr(err)

		v2plusscripts.PrintCoordinatorConfig(coordinator)
	case "coordinator-set-config":
		cmd := flag.NewFlagSet("coordinator-set-config", flag.ExitOnError)
		setConfigAddress := cmd.String("coordinator-address", "", "coordinator address")
		minConfs := cmd.Int("min-confs", 3, "min confs")
		maxGasLimit := cmd.Int64("max-gas-limit", 2.5e6, "max gas limit")
		stalenessSeconds := cmd.Int64("staleness-seconds", 86400, "staleness in seconds")
		gasAfterPayment := cmd.Int64("gas-after-payment", 33285, "gas after payment calculation")
		fallbackWeiPerUnitLink := cmd.String("fallback-wei-per-unit-link", "", "fallback wei per unit link")
		flatFeeLinkPPM := cmd.Int64("flat-fee-link-ppm", 500, "fulfillment flat fee LINK ppm")
		flatFeeEthPPM := cmd.Int64("flat-fee-eth-ppm", 500, "fulfillment flat fee ETH ppm")

		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address", "fallback-wei-per-unit-link")
		coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(common.HexToAddress(*setConfigAddress), e.Ec)
		helpers.PanicErr(err)

		v2plusscripts.SetCoordinatorConfig(
			e,
			*coordinator,
			uint16(*minConfs),
			uint32(*maxGasLimit),
			uint32(*stalenessSeconds),
			uint32(*gasAfterPayment),
			decimal.RequireFromString(*fallbackWeiPerUnitLink).BigInt(),
			vrf_coordinator_v2_5.VRFCoordinatorV25FeeConfig{
				FulfillmentFlatFeeLinkPPM:   uint32(*flatFeeLinkPPM),
				FulfillmentFlatFeeNativePPM: uint32(*flatFeeEthPPM),
			},
		)
	case "coordinator-register-key":
		coordinatorRegisterKey := flag.NewFlagSet("coordinator-register-key", flag.ExitOnError)
		registerKeyAddress := coordinatorRegisterKey.String("address", "", "coordinator address")
		registerKeyUncompressedPubKey := coordinatorRegisterKey.String("pubkey", "", "uncompressed pubkey")
		registerKeyOracleAddress := coordinatorRegisterKey.String("oracle-address", "", "oracle address")
		helpers.ParseArgs(coordinatorRegisterKey, os.Args[2:], "address", "pubkey", "oracle-address")
		coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(common.HexToAddress(*registerKeyAddress), e.Ec)
		helpers.PanicErr(err)

		// Put key in ECDSA format
		if strings.HasPrefix(*registerKeyUncompressedPubKey, "0x") {
			*registerKeyUncompressedPubKey = strings.Replace(*registerKeyUncompressedPubKey, "0x", "04", 1)
		}

		v2plusscripts.RegisterCoordinatorProvingKey(e, *coordinator, *registerKeyUncompressedPubKey, *registerKeyOracleAddress)
	case "coordinator-deregister-key":
		coordinatorDeregisterKey := flag.NewFlagSet("coordinator-deregister-key", flag.ExitOnError)
		deregisterKeyAddress := coordinatorDeregisterKey.String("address", "", "coordinator address")
		deregisterKeyUncompressedPubKey := coordinatorDeregisterKey.String("pubkey", "", "uncompressed pubkey")
		helpers.ParseArgs(coordinatorDeregisterKey, os.Args[2:], "address", "pubkey")
		coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(common.HexToAddress(*deregisterKeyAddress), e.Ec)
		helpers.PanicErr(err)

		// Put key in ECDSA format
		if strings.HasPrefix(*deregisterKeyUncompressedPubKey, "0x") {
			*deregisterKeyUncompressedPubKey = strings.Replace(*deregisterKeyUncompressedPubKey, "0x", "04", 1)
		}
		pubBytes, err := hex.DecodeString(*deregisterKeyUncompressedPubKey)
		helpers.PanicErr(err)
		pk, err := crypto.UnmarshalPubkey(pubBytes)
		helpers.PanicErr(err)
		tx, err := coordinator.DeregisterProvingKey(e.Owner, [2]*big.Int{pk.X, pk.Y})
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
	case "coordinator-subscription":
		coordinatorSub := flag.NewFlagSet("coordinator-subscription", flag.ExitOnError)
		address := coordinatorSub.String("address", "", "coordinator address")
		subID := coordinatorSub.String("sub-id", "", "sub-id")
		helpers.ParseArgs(coordinatorSub, os.Args[2:], "address", "sub-id")
		coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(common.HexToAddress(*address), e.Ec)
		helpers.PanicErr(err)
		fmt.Println("sub-id", *subID, "address", *address, coordinator.Address())
		parsedSubID := parseSubID(*subID)
		s, err := coordinator.GetSubscription(nil, parsedSubID)
		helpers.PanicErr(err)
		fmt.Printf("Subscription %+v\n", s)
	case "consumer-deploy":
		consumerDeployCmd := flag.NewFlagSet("consumer-deploy", flag.ExitOnError)
		consumerCoordinator := consumerDeployCmd.String("coordinator-address", "", "coordinator address")
		keyHash := consumerDeployCmd.String("key-hash", "", "key hash")
		consumerLinkAddress := consumerDeployCmd.String("link-address", "", "link-address")
		nativePayment := consumerDeployCmd.Bool("native-payment", false, "whether to use native payment or not")

		// TODO: add other params
		helpers.ParseArgs(consumerDeployCmd, os.Args[2:], "coordinator-address", "key-hash", "link-address")
		keyHashBytes := common.HexToHash(*keyHash)
		_, tx, _, err := vrf_v2plus_single_consumer.DeployVRFV2PlusSingleConsumerExample(
			e.Owner,
			e.Ec,
			common.HexToAddress(*consumerCoordinator),
			common.HexToAddress(*consumerLinkAddress),
			uint32(1000000), // gas callback
			uint16(5),       // confs
			uint32(1),       // words
			keyHashBytes,
			*nativePayment)
		helpers.PanicErr(err)
		helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
	case "consumer-subscribe":
		consumerSubscribeCmd := flag.NewFlagSet("consumer-subscribe", flag.ExitOnError)
		consumerSubscribeAddress := consumerSubscribeCmd.String("address", "", "consumer address")
		helpers.ParseArgs(consumerSubscribeCmd, os.Args[2:], "address")
		consumer, err := vrf_v2plus_single_consumer.NewVRFV2PlusSingleConsumerExample(common.HexToAddress(*consumerSubscribeAddress), e.Ec)
		helpers.PanicErr(err)
		tx, err := consumer.Subscribe(e.Owner)
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
	case "link-balance":
		linkBalanceCmd := flag.NewFlagSet("link-balance", flag.ExitOnError)
		linkAddress := linkBalanceCmd.String("link-address", "", "link-address")
		address := linkBalanceCmd.String("address", "", "address")
		helpers.ParseArgs(linkBalanceCmd, os.Args[2:], "link-address", "address")
		lt, err := link_token_interface.NewLinkToken(common.HexToAddress(*linkAddress), e.Ec)
		helpers.PanicErr(err)
		b, err := lt.BalanceOf(nil, common.HexToAddress(*address))
		helpers.PanicErr(err)
		fmt.Println(b)
	case "consumer-cancel":
		consumerCancelCmd := flag.NewFlagSet("consumer-cancel", flag.ExitOnError)
		consumerCancelAddress := consumerCancelCmd.String("address", "", "consumer address")
		helpers.ParseArgs(consumerCancelCmd, os.Args[2:], "address")
		consumer, err := vrf_v2plus_single_consumer.NewVRFV2PlusSingleConsumerExample(common.HexToAddress(*consumerCancelAddress), e.Ec)
		helpers.PanicErr(err)
		tx, err := consumer.Unsubscribe(e.Owner, e.Owner.From)
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
	case "consumer-topup":
		// NOTE NEED TO FUND CONSUMER WITH LINK FIRST
		consumerTopupCmd := flag.NewFlagSet("consumer-topup", flag.ExitOnError)
		consumerTopupAmount := consumerTopupCmd.String("amount", "", "amount in juels")
		consumerTopupAddress := consumerTopupCmd.String("address", "", "consumer address")
		helpers.ParseArgs(consumerTopupCmd, os.Args[2:], "amount", "address")
		consumer, err := vrf_v2plus_single_consumer.NewVRFV2PlusSingleConsumerExample(common.HexToAddress(*consumerTopupAddress), e.Ec)
		helpers.PanicErr(err)
		amount, s := big.NewInt(0).SetString(*consumerTopupAmount, 10)
		if !s {
			panic(fmt.Sprintf("failed to parse top up amount '%s'", *consumerTopupAmount))
		}
		tx, err := consumer.TopUpSubscription(e.Owner, amount)
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
	case "consumer-request":
		consumerRequestCmd := flag.NewFlagSet("consumer-request", flag.ExitOnError)
		consumerRequestAddress := consumerRequestCmd.String("address", "", "consumer address")
		helpers.ParseArgs(consumerRequestCmd, os.Args[2:], "address")
		consumer, err := vrf_v2plus_single_consumer.NewVRFV2PlusSingleConsumerExample(common.HexToAddress(*consumerRequestAddress), e.Ec)
		helpers.PanicErr(err)
		tx, err := consumer.RequestRandomWords(e.Owner)
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
	case "consumer-fund-and-request":
		consumerRequestCmd := flag.NewFlagSet("consumer-request", flag.ExitOnError)
		consumerRequestAddress := consumerRequestCmd.String("address", "", "consumer address")
		helpers.ParseArgs(consumerRequestCmd, os.Args[2:], "address")
		consumer, err := vrf_v2plus_single_consumer.NewVRFV2PlusSingleConsumerExample(common.HexToAddress(*consumerRequestAddress), e.Ec)
		helpers.PanicErr(err)
		// Fund and request 3 link
		tx, err := consumer.FundAndRequestRandomWords(e.Owner, big.NewInt(3000000000000000000))
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
	case "consumer-print":
		consumerPrint := flag.NewFlagSet("consumer-print", flag.ExitOnError)
		address := consumerPrint.String("address", "", "consumer address")
		helpers.ParseArgs(consumerPrint, os.Args[2:], "address")
		consumer, err := vrf_v2plus_single_consumer.NewVRFV2PlusSingleConsumerExample(common.HexToAddress(*address), e.Ec)
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
	case "deploy-universe":
		v2plusscripts.DeployUniverseViaCLI(e)
	case "generate-proof-v2-plus":
		generateProofForV2Plus(e)
	case "eoa-consumer-deploy":
		consumerDeployCmd := flag.NewFlagSet("eoa-consumer-deploy", flag.ExitOnError)
		consumerCoordinator := consumerDeployCmd.String("coordinator-address", "", "coordinator address")
		consumerLinkAddress := consumerDeployCmd.String("link-address", "", "link-address")
		helpers.ParseArgs(consumerDeployCmd, os.Args[2:], "coordinator-address", "link-address", "key-hash")

		v2plusscripts.EoaDeployConsumer(e, *consumerCoordinator, *consumerLinkAddress)
	case "eoa-load-test-consumer-deploy":
		loadTestConsumerDeployCmd := flag.NewFlagSet("eoa-load-test-consumer-deploy", flag.ExitOnError)
		consumerCoordinator := loadTestConsumerDeployCmd.String("coordinator-address", "", "coordinator address")
		consumerLinkAddress := loadTestConsumerDeployCmd.String("link-address", "", "link-address")
		helpers.ParseArgs(loadTestConsumerDeployCmd, os.Args[2:], "coordinator-address", "link-address")
		_, tx, _, err := vrf_load_test_external_sub_owner.DeployVRFLoadTestExternalSubOwner(
			e.Owner,
			e.Ec,
			common.HexToAddress(*consumerCoordinator),
			common.HexToAddress(*consumerLinkAddress))
		helpers.PanicErr(err)
		helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
	case "eoa-load-test-consumer-with-metrics-deploy":
		loadTestConsumerDeployCmd := flag.NewFlagSet("eoa-load-test-consumer-with-metrics-deploy", flag.ExitOnError)
		consumerCoordinator := loadTestConsumerDeployCmd.String("coordinator-address", "", "coordinator address")
		helpers.ParseArgs(loadTestConsumerDeployCmd, os.Args[2:], "coordinator-address")
		_, tx, _, err := vrf_v2plus_load_test_with_metrics.DeployVRFV2PlusLoadTestWithMetrics(
			e.Owner,
			e.Ec,
			common.HexToAddress(*consumerCoordinator),
		)
		helpers.PanicErr(err)
		helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
	case "eoa-create-sub":
		createSubCmd := flag.NewFlagSet("eoa-create-sub", flag.ExitOnError)
		coordinatorAddress := createSubCmd.String("coordinator-address", "", "coordinator address")
		helpers.ParseArgs(createSubCmd, os.Args[2:], "coordinator-address")
		coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(common.HexToAddress(*coordinatorAddress), e.Ec)
		helpers.PanicErr(err)
		v2plusscripts.EoaCreateSub(e, *coordinator)
	case "eoa-add-sub-consumer":
		addSubConsCmd := flag.NewFlagSet("eoa-add-sub-consumer", flag.ExitOnError)
		coordinatorAddress := addSubConsCmd.String("coordinator-address", "", "coordinator address")
		subID := addSubConsCmd.String("sub-id", "", "sub-id")
		consumerAddress := addSubConsCmd.String("consumer-address", "", "consumer address")
		helpers.ParseArgs(addSubConsCmd, os.Args[2:], "coordinator-address", "sub-id", "consumer-address")
		coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(common.HexToAddress(*coordinatorAddress), e.Ec)
		helpers.PanicErr(err)
		parsedSubID := parseSubID(*subID)
		v2plusscripts.EoaAddConsumerToSub(e, *coordinator, parsedSubID, *consumerAddress)
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
		coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(common.HexToAddress(*coordinatorAddress), e.Ec)
		helpers.PanicErr(err)
		fmt.Println(amount, consumerLinkAddress)
		txcreate, err := coordinator.CreateSubscription(e.Owner)
		helpers.PanicErr(err)
		fmt.Println("Create sub", "TX", helpers.ExplorerLink(e.ChainID, txcreate.Hash()))
		helpers.ConfirmTXMined(context.Background(), e.Ec, txcreate, e.ChainID)
		sub := make(chan *vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCreated)
		subscription, err := coordinator.WatchSubscriptionCreated(nil, sub, nil)
		helpers.PanicErr(err)
		defer subscription.Unsubscribe()
		created := <-sub
		linkToken, err := link_token_interface.NewLinkToken(common.HexToAddress(*consumerLinkAddress), e.Ec)
		helpers.PanicErr(err)
		bal, err := linkToken.BalanceOf(nil, e.Owner.From)
		helpers.PanicErr(err)
		fmt.Println("OWNER BALANCE", bal, e.Owner.From.String(), amount.String())
		b, err := utils.ABIEncode(`[{"type":"uint64"}]`, created.SubId)
		helpers.PanicErr(err)
		e.Owner.GasLimit = 500000
		tx, err := linkToken.TransferAndCall(e.Owner, coordinator.Address(), amount, b)
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, fmt.Sprintf("Sub id: %d", created.SubId))
		subFunded := make(chan *vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionFunded)
		fundSub, err := coordinator.WatchSubscriptionFunded(nil, subFunded, []*big.Int{created.SubId})
		helpers.PanicErr(err)
		defer fundSub.Unsubscribe()
		<-subFunded // Add a consumer once its funded
		txadd, err := coordinator.AddConsumer(e.Owner, created.SubId, common.HexToAddress(*consumerAddress))
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, txadd, e.ChainID)
	case "eoa-request":
		request := flag.NewFlagSet("eoa-request", flag.ExitOnError)
		consumerAddress := request.String("consumer-address", "", "consumer address")
		subID := request.String("sub-id", "", "subscription ID")
		cbGasLimit := request.Uint("cb-gas-limit", 1_000_000, "callback gas limit")
		requestConfirmations := request.Uint("request-confirmations", 3, "minimum request confirmations")
		numWords := request.Uint("num-words", 3, "number of words to request")
		keyHash := request.String("key-hash", "", "key hash")
		nativePayment := request.Bool("native-payment", false, "whether to use native payment or not")
		helpers.ParseArgs(request, os.Args[2:], "consumer-address")
		keyHashBytes := common.HexToHash(*keyHash)
		consumer, err := vrf_v2plus_sub_owner.NewVRFV2PlusExternalSubOwnerExample(
			common.HexToAddress(*consumerAddress),
			e.Ec)
		helpers.PanicErr(err)
		tx, err := consumer.RequestRandomWords(e.Owner, parseSubID(*subID), uint32(*cbGasLimit), uint16(*requestConfirmations), uint32(*numWords), keyHashBytes, *nativePayment)
		helpers.PanicErr(err)
		fmt.Println("TX", helpers.ExplorerLink(e.ChainID, tx.Hash()))
		r, err := bind.WaitMined(context.Background(), e.Ec, tx)
		helpers.PanicErr(err)
		fmt.Println("Receipt blocknumber:", r.BlockNumber)
	case "eoa-load-test-read":
		cmd := flag.NewFlagSet("eoa-load-test-read", flag.ExitOnError)
		consumerAddress := cmd.String("consumer-address", "", "consumer address")
		helpers.ParseArgs(cmd, os.Args[2:], "consumer-address")
		consumer, err := vrf_load_test_external_sub_owner.NewVRFLoadTestExternalSubOwner(
			common.HexToAddress(*consumerAddress),
			e.Ec)
		helpers.PanicErr(err)
		rc, err := consumer.SResponseCount(nil)
		helpers.PanicErr(err)
		fmt.Println("load tester", *consumerAddress, "response count:", rc)
	case "eoa-load-test-request":
		request := flag.NewFlagSet("eoa-load-test-request", flag.ExitOnError)
		consumerAddress := request.String("consumer-address", "", "consumer address")
		subID := request.Uint64("sub-id", 0, "subscription ID")
		requestConfirmations := request.Uint("request-confirmations", 3, "minimum request confirmations")
		keyHash := request.String("key-hash", "", "key hash")
		requests := request.Uint("requests", 10, "number of randomness requests to make per run")
		runs := request.Uint("runs", 1, "number of runs to do. total randomness requests will be (requests * runs).")
		helpers.ParseArgs(request, os.Args[2:], "consumer-address", "sub-id", "key-hash")
		keyHashBytes := common.HexToHash(*keyHash)
		consumer, err := vrf_load_test_external_sub_owner.NewVRFLoadTestExternalSubOwner(
			common.HexToAddress(*consumerAddress),
			e.Ec)
		helpers.PanicErr(err)
		var txes []*types.Transaction
		for i := 0; i < int(*runs); i++ {
			tx, err := consumer.RequestRandomWords(e.Owner, *subID, uint16(*requestConfirmations),
				keyHashBytes, uint16(*requests))
			helpers.PanicErr(err)
			fmt.Printf("TX %d: %s\n", i+1, helpers.ExplorerLink(e.ChainID, tx.Hash()))
			txes = append(txes, tx)
		}
		fmt.Println("Total number of requests sent:", (*requests)*(*runs))
		fmt.Println("fetching receipts for all transactions")
		for i, tx := range txes {
			helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, fmt.Sprintf("load test %d", i+1))
		}
	case "eoa-load-test-request-with-metrics":
		request := flag.NewFlagSet("eoa-load-test-request-with-metrics", flag.ExitOnError)
		consumerAddress := request.String("consumer-address", "", "consumer address")
		subID := request.String("sub-id", "", "subscription ID")
		requestConfirmations := request.Uint("request-confirmations", 3, "minimum request confirmations")
		keyHash := request.String("key-hash", "", "key hash")
		cbGasLimit := request.Uint("cb-gas-limit", 100_000, "request callback gas limit")
		nativePaymentEnabled := request.Bool("native-payment-enabled", false, "native payment enabled")
		numWords := request.Uint("num-words", 1, "num words to request")
		requests := request.Uint("requests", 10, "number of randomness requests to make per run")
		runs := request.Uint("runs", 1, "number of runs to do. total randomness requests will be (requests * runs).")
		helpers.ParseArgs(request, os.Args[2:], "consumer-address", "sub-id", "key-hash")
		keyHashBytes := common.HexToHash(*keyHash)
		consumer, err := vrf_v2plus_load_test_with_metrics.NewVRFV2PlusLoadTestWithMetrics(
			common.HexToAddress(*consumerAddress),
			e.Ec)
		helpers.PanicErr(err)
		var txes []*types.Transaction
		for i := 0; i < int(*runs); i++ {
			tx, err := consumer.RequestRandomWords(
				e.Owner,
				decimal.RequireFromString(*subID).BigInt(),
				uint16(*requestConfirmations),
				keyHashBytes,
				uint32(*cbGasLimit),
				*nativePaymentEnabled,
				uint32(*numWords),
				uint16(*requests),
			)
			helpers.PanicErr(err)
			fmt.Printf("TX %d: %s\n", i+1, helpers.ExplorerLink(e.ChainID, tx.Hash()))
			txes = append(txes, tx)
		}
		fmt.Println("Total number of requests sent:", (*requests)*(*runs))
		fmt.Println("fetching receipts for all transactions")
		for i, tx := range txes {
			helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, fmt.Sprintf("load test %d", i+1))
		}
	case "eoa-load-test-read-metrics":
		request := flag.NewFlagSet("eoa-load-test-read-metrics", flag.ExitOnError)
		consumerAddress := request.String("consumer-address", "", "consumer address")
		helpers.ParseArgs(request, os.Args[2:], "consumer-address")
		consumer, err := vrf_v2plus_load_test_with_metrics.NewVRFV2PlusLoadTestWithMetrics(
			common.HexToAddress(*consumerAddress),
			e.Ec)
		helpers.PanicErr(err)
		responseCount, err := consumer.SResponseCount(nil)
		helpers.PanicErr(err)
		fmt.Println("Response Count: ", responseCount)
		requestCount, err := consumer.SRequestCount(nil)
		helpers.PanicErr(err)
		fmt.Println("Request Count: ", requestCount)
		averageFulfillmentInMillions, err := consumer.SAverageFulfillmentInMillions(nil)
		helpers.PanicErr(err)
		fmt.Println("Average Fulfillment In Millions: ", averageFulfillmentInMillions)
		slowestFulfillment, err := consumer.SSlowestFulfillment(nil)
		helpers.PanicErr(err)
		fmt.Println("Slowest Fulfillment: ", slowestFulfillment)
		fastestFulfillment, err := consumer.SFastestFulfillment(nil)
		helpers.PanicErr(err)
		fmt.Println("Fastest Fulfillment: ", fastestFulfillment)
	case "eoa-load-test-reset-metrics":
		request := flag.NewFlagSet("eoa-load-test-reset-metrics", flag.ExitOnError)
		consumerAddress := request.String("consumer-address", "", "consumer address")
		helpers.ParseArgs(request, os.Args[2:], "consumer-address")
		consumer, err := vrf_v2plus_load_test_with_metrics.NewVRFV2PlusLoadTestWithMetrics(
			common.HexToAddress(*consumerAddress),
			e.Ec)
		helpers.PanicErr(err)
		_, err = consumer.Reset(e.Owner)
		helpers.PanicErr(err)
		fmt.Println("Load Test Consumer With Metrics was reset ")
	case "eoa-transfer-sub":
		trans := flag.NewFlagSet("eoa-transfer-sub", flag.ExitOnError)
		coordinatorAddress := trans.String("coordinator-address", "", "coordinator address")
		subID := trans.String("sub-id", "", "sub-id")
		to := trans.String("to", "", "to")
		helpers.ParseArgs(trans, os.Args[2:], "coordinator-address", "sub-id", "to")
		coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(common.HexToAddress(*coordinatorAddress), e.Ec)
		helpers.PanicErr(err)
		tx, err := coordinator.RequestSubscriptionOwnerTransfer(e.Owner, parseSubID(*subID), common.HexToAddress(*to))
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
	case "eoa-accept-sub":
		accept := flag.NewFlagSet("eoa-accept-sub", flag.ExitOnError)
		coordinatorAddress := accept.String("coordinator-address", "", "coordinator address")
		subID := accept.String("sub-id", "", "sub-id")
		helpers.ParseArgs(accept, os.Args[2:], "coordinator-address", "sub-id")
		coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(common.HexToAddress(*coordinatorAddress), e.Ec)
		helpers.PanicErr(err)
		tx, err := coordinator.AcceptSubscriptionOwnerTransfer(e.Owner, parseSubID(*subID))
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
	case "eoa-cancel-sub":
		cancel := flag.NewFlagSet("eoa-cancel-sub", flag.ExitOnError)
		coordinatorAddress := cancel.String("coordinator-address", "", "coordinator address")
		subID := cancel.String("sub-id", "", "sub-id")
		helpers.ParseArgs(cancel, os.Args[2:], "coordinator-address", "sub-id")
		coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(common.HexToAddress(*coordinatorAddress), e.Ec)
		helpers.PanicErr(err)
		tx, err := coordinator.CancelSubscription(e.Owner, parseSubID(*subID), e.Owner.From)
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
	case "eoa-fund-sub-with-native-token":
		fund := flag.NewFlagSet("eoa-fund-sub-with-native-token", flag.ExitOnError)
		coordinatorAddress := fund.String("coordinator-address", "", "coordinator address")
		amountStr := fund.String("amount", "", "amount to fund in wei")
		subID := fund.String("sub-id", "", "sub-id")
		helpers.ParseArgs(fund, os.Args[2:], "coordinator-address", "amount", "sub-id")
		amount, s := big.NewInt(0).SetString(*amountStr, 10)
		if !s {
			panic(fmt.Sprintf("failed to parse top up amount '%s'", *amountStr))
		}
		parsedSubID := parseSubID(*subID)

		v2plusscripts.EoaFundSubWithNative(e, common.HexToAddress(*coordinatorAddress), parsedSubID, amount)
	case "eoa-fund-sub":
		fund := flag.NewFlagSet("eoa-fund-sub", flag.ExitOnError)
		coordinatorAddress := fund.String("coordinator-address", "", "coordinator address")
		amountStr := fund.String("amount", "", "amount to fund in juels")
		subID := fund.String("sub-id", "", "sub-id")
		consumerLinkAddress := fund.String("link-address", "", "link-address")
		helpers.ParseArgs(fund, os.Args[2:], "coordinator-address", "amount", "sub-id", "link-address")
		amount, s := big.NewInt(0).SetString(*amountStr, 10)
		if !s {
			panic(fmt.Sprintf("failed to parse top up amount '%s'", *amountStr))
		}
		coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(common.HexToAddress(*coordinatorAddress), e.Ec)
		helpers.PanicErr(err)

		v2plusscripts.EoaFundSubWithLink(e, *coordinator, *consumerLinkAddress, amount, parseSubID(*subID))
	case "eoa-read":
		cmd := flag.NewFlagSet("eoa-read", flag.ExitOnError)
		consumerAddress := cmd.String("consumer", "", "consumer address")
		helpers.ParseArgs(cmd, os.Args[2:], "consumer")
		consumer, err := vrf_v2plus_single_consumer.NewVRFV2PlusSingleConsumerExample(common.HexToAddress(*consumerAddress), e.Ec)
		helpers.PanicErr(err)
		word, err := consumer.SRandomWords(nil, big.NewInt(0))
		if err != nil {
			fmt.Println("no words (yet?)")
		}
		reqID, err := consumer.SRequestId(nil)
		helpers.PanicErr(err)
		fmt.Println("request id:", reqID.String(), "1st random word:", word)
	case "owner-cancel-sub":
		cancel := flag.NewFlagSet("owner-cancel-sub", flag.ExitOnError)
		coordinatorAddress := cancel.String("coordinator-address", "", "coordinator address")
		subID := cancel.String("sub-id", "", "sub-id")
		helpers.ParseArgs(cancel, os.Args[2:], "coordinator-address", "sub-id")
		coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(common.HexToAddress(*coordinatorAddress), e.Ec)
		helpers.PanicErr(err)
		tx, err := coordinator.OwnerCancelSubscription(e.Owner, parseSubID(*subID))
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
	case "sub-balance":
		consumerBalanceCmd := flag.NewFlagSet("sub-balance", flag.ExitOnError)
		coordinatorAddress := consumerBalanceCmd.String("coordinator-address", "", "coordinator address")
		subID := consumerBalanceCmd.String("sub-id", "", "subscription id")
		helpers.ParseArgs(consumerBalanceCmd, os.Args[2:], "coordinator-address", "sub-id")
		coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(common.HexToAddress(*coordinatorAddress), e.Ec)
		helpers.PanicErr(err)
		resp, err := coordinator.GetSubscription(nil, parseSubID(*subID))
		helpers.PanicErr(err)
		fmt.Println("sub id", *subID, "balance:", resp.Balance)
	case "coordinator-withdrawable-tokens":
		withdrawableTokensCmd := flag.NewFlagSet("coordinator-withdrawable-tokens", flag.ExitOnError)
		coordinator := withdrawableTokensCmd.String("coordinator-address", "", "coordinator address")
		oracle := withdrawableTokensCmd.String("oracle-address", "", "oracle address")
		start := withdrawableTokensCmd.Int("start-link", 10_000, "the starting amount of LINK to check")
		helpers.ParseArgs(withdrawableTokensCmd, os.Args[2:], "coordinator-address", "oracle-address")

		coordinatorAddress := common.HexToAddress(*coordinator)
		oracleAddress := common.HexToAddress(*oracle)
		abi, err := vrf_coordinator_v2_5.VRFCoordinatorV25MetaData.GetAbi()
		helpers.PanicErr(err)

		isWithdrawable := func(amount *big.Int) bool {
			data, err := abi.Pack("oracleWithdraw", oracleAddress /* this can be any address */, amount)
			helpers.PanicErr(err)

			_, err = e.Ec.CallContract(context.Background(), ethereum.CallMsg{
				From: oracleAddress,
				To:   &coordinatorAddress,
				Data: data,
			}, nil)
			if err == nil {
				return true
			} else if strings.Contains(err.Error(), "execution reverted") {
				return false
			} else {
				panic(err)
			}
		}

		result := helpers.BinarySearch(assets.Ether(int64(*start*2)).ToInt(), big.NewInt(0), isWithdrawable)

		fmt.Printf("Withdrawable amount for oracle %s is %s\n", oracleAddress.String(), result.String())
	case "coordinator-transfer-ownership":
		cmd := flag.NewFlagSet("coordinator-transfer-ownership", flag.ExitOnError)
		coordinatorAddress := cmd.String("coordinator-address", "", "v2 coordinator address")
		newOwner := cmd.String("new-owner", "", "new owner address")
		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address", "new-owner")

		coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(common.HexToAddress(*coordinatorAddress), e.Ec)
		helpers.PanicErr(err)

		tx, err := coordinator.TransferOwnership(e.Owner, common.HexToAddress(*newOwner))
		helpers.PanicErr(err)

		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "transfer ownership to", *newOwner)
	case "coordinator-reregister-proving-key":
		coordinatorReregisterKey := flag.NewFlagSet("coordinator-register-key", flag.ExitOnError)
		coordinatorAddress := coordinatorReregisterKey.String("coordinator-address", "", "coordinator address")
		uncompressedPubKey := coordinatorReregisterKey.String("pubkey", "", "uncompressed pubkey")
		newOracleAddress := coordinatorReregisterKey.String("new-oracle-address", "", "oracle address")
		skipDeregister := coordinatorReregisterKey.Bool("skip-deregister", false, "if true, key will not be deregistered")
		helpers.ParseArgs(coordinatorReregisterKey, os.Args[2:], "coordinator-address", "pubkey", "new-oracle-address")

		coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(common.HexToAddress(*coordinatorAddress), e.Ec)
		helpers.PanicErr(err)

		// Put key in ECDSA format
		if strings.HasPrefix(*uncompressedPubKey, "0x") {
			*uncompressedPubKey = strings.Replace(*uncompressedPubKey, "0x", "04", 1)
		}
		pubBytes, err := hex.DecodeString(*uncompressedPubKey)
		helpers.PanicErr(err)
		pk, err := crypto.UnmarshalPubkey(pubBytes)
		helpers.PanicErr(err)

		var deregisterTx *types.Transaction
		if !*skipDeregister {
			deregisterTx, err = coordinator.DeregisterProvingKey(e.Owner, [2]*big.Int{pk.X, pk.Y})
			helpers.PanicErr(err)
			fmt.Println("Deregister transaction", helpers.ExplorerLink(e.ChainID, deregisterTx.Hash()))
		}

		// Use a higher gas price for the register call
		e.Owner.GasPrice.Mul(e.Owner.GasPrice, big.NewInt(2))
		registerTx, err := coordinator.RegisterProvingKey(e.Owner,
			common.HexToAddress(*newOracleAddress),
			[2]*big.Int{pk.X, pk.Y})
		helpers.PanicErr(err)
		fmt.Println("Register transaction", helpers.ExplorerLink(e.ChainID, registerTx.Hash()))

		if !*skipDeregister {
			fmt.Println("Waiting for deregister transaction to be mined...")
			var deregisterReceipt *types.Receipt
			deregisterReceipt, err = bind.WaitMined(context.Background(), e.Ec, deregisterTx)
			helpers.PanicErr(err)
			fmt.Printf("Deregister transaction included in block %s\n", deregisterReceipt.BlockNumber.String())
		}

		fmt.Println("Waiting for register transaction to be mined...")
		registerReceipt, err := bind.WaitMined(context.Background(), e.Ec, registerTx)
		helpers.PanicErr(err)
		fmt.Printf("Register transaction included in block %s\n", registerReceipt.BlockNumber.String())
	case "wrapper-deploy":
		cmd := flag.NewFlagSet("wrapper-deploy", flag.ExitOnError)
		linkAddress := cmd.String("link-address", "", "address of link token")
		linkETHFeedAddress := cmd.String("link-eth-feed", "", "address of link-eth-feed")
		coordinatorAddress := cmd.String("coordinator-address", "", "address of the vrf coordinator v2 contract")
		helpers.ParseArgs(cmd, os.Args[2:], "link-address", "link-eth-feed", "coordinator-address")
		v2plusscripts.WrapperDeploy(e,
			common.HexToAddress(*linkAddress),
			common.HexToAddress(*linkETHFeedAddress),
			common.HexToAddress(*coordinatorAddress))
	case "wrapper-withdraw":
		cmd := flag.NewFlagSet("wrapper-withdraw", flag.ExitOnError)
		wrapperAddress := cmd.String("wrapper-address", "", "address of the VRFV2Wrapper contract")
		recipientAddress := cmd.String("recipient-address", "", "address to withdraw to")
		linkAddress := cmd.String("link-address", "", "address of link token")
		helpers.ParseArgs(cmd, os.Args[2:], "wrapper-address", "recipient-address", "link-address")
		wrapper, err := vrfv2plus_wrapper.NewVRFV2PlusWrapper(common.HexToAddress(*wrapperAddress), e.Ec)
		helpers.PanicErr(err)
		link, err := link_token_interface.NewLinkToken(common.HexToAddress(*linkAddress), e.Ec)
		helpers.PanicErr(err)
		balance, err := link.BalanceOf(nil, common.HexToAddress(*wrapperAddress))
		helpers.PanicErr(err)
		tx, err := wrapper.Withdraw(e.Owner, common.HexToAddress(*recipientAddress), balance)
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "withdrawing", balance.String(), "Juels from", *wrapperAddress, "to", *recipientAddress)
	case "wrapper-get-subscription-id":
		cmd := flag.NewFlagSet("wrapper-get-subscription-id", flag.ExitOnError)
		wrapperAddress := cmd.String("wrapper-address", "", "address of the VRFV2Wrapper contract")
		helpers.ParseArgs(cmd, os.Args[2:], "wrapper-address")
		wrapper, err := vrfv2plus_wrapper.NewVRFV2PlusWrapper(common.HexToAddress(*wrapperAddress), e.Ec)
		helpers.PanicErr(err)
		subID, err := wrapper.SUBSCRIPTIONID(nil)
		helpers.PanicErr(err)
		fmt.Println("subscription id of wrapper", *wrapperAddress, "is:", subID)
	case "wrapper-configure":
		cmd := flag.NewFlagSet("wrapper-configure", flag.ExitOnError)
		wrapperAddress := cmd.String("wrapper-address", "", "address of the VRFV2Wrapper contract")
		wrapperGasOverhead := cmd.Uint("wrapper-gas-overhead", 50_000, "amount of gas overhead in wrapper fulfillment")
		coordinatorGasOverhead := cmd.Uint("coordinator-gas-overhead", 52_000, "amount of gas overhead in coordinator fulfillment")
		wrapperPremiumPercentage := cmd.Uint("wrapper-premium-percentage", 25, "gas premium charged by wrapper")
		keyHash := cmd.String("key-hash", "", "the keyhash that wrapper requests should use")
		maxNumWords := cmd.Uint("max-num-words", 10, "the keyhash that wrapper requests should use")
		fallbackWeiPerUnitLink := cmd.String("fallback-wei-per-unit-link", "", "the fallback wei per unit link")
		stalenessSeconds := cmd.Uint("staleness-seconds", 86400, "the number of seconds of staleness to allow")
		fulfillmentFlatFeeLinkPPM := cmd.Uint("fulfillment-flat-fee-link-ppm", 500, "the link flat fee in ppm to charge for fulfillment")
		fulfillmentFlatFeeNativePPM := cmd.Uint("fulfillment-flat-fee-native-ppm", 500, "the native flat fee in ppm to charge for fulfillment")
		helpers.ParseArgs(cmd, os.Args[2:], "wrapper-address", "key-hash", "fallback-wei-per-unit-link")

		v2plusscripts.WrapperConfigure(e,
			common.HexToAddress(*wrapperAddress),
			*wrapperGasOverhead,
			*coordinatorGasOverhead,
			*wrapperPremiumPercentage,
			*keyHash,
			*maxNumWords,
			decimal.RequireFromString(*fallbackWeiPerUnitLink).BigInt(),
			uint32(*stalenessSeconds),
			uint32(*fulfillmentFlatFeeLinkPPM),
			uint32(*fulfillmentFlatFeeNativePPM))
	case "wrapper-get-fulfillment-tx-size":
		cmd := flag.NewFlagSet("wrapper-get-fulfillment-tx-size", flag.ExitOnError)
		wrapperAddress := cmd.String("wrapper-address", "", "address of the VRFV2Wrapper contract")
		helpers.ParseArgs(cmd, os.Args[2:], "wrapper-address")
		wrapper, err := vrfv2plus_wrapper.NewVRFV2PlusWrapper(common.HexToAddress(*wrapperAddress), e.Ec)
		helpers.PanicErr(err)
		size, err := wrapper.SFulfillmentTxSizeBytes(nil)
		helpers.PanicErr(err)
		fmt.Println("fulfillment tx size of wrapper", *wrapperAddress, "is:", size)
	case "wrapper-set-fulfillment-tx-size":
		cmd := flag.NewFlagSet("wrapper-set-fulfillment-tx-size", flag.ExitOnError)
		wrapperAddress := cmd.String("wrapper-address", "", "address of the VRFV2Wrapper contract")
		size := cmd.Uint("size", 0, "size of the fulfillment transaction")
		helpers.ParseArgs(cmd, os.Args[2:], "wrapper-address", "size")
		wrapper, err := vrfv2plus_wrapper.NewVRFV2PlusWrapper(common.HexToAddress(*wrapperAddress), e.Ec)
		helpers.PanicErr(err)
		tx, err := wrapper.SetFulfillmentTxSize(e.Owner, uint32(*size))
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "set fulfillment tx size")
	case "wrapper-consumer-deploy":
		cmd := flag.NewFlagSet("wrapper-consumer-deploy", flag.ExitOnError)
		linkAddress := cmd.String("link-address", "", "address of link token")
		wrapperAddress := cmd.String("wrapper-address", "", "address of the VRFV2Wrapper contract")
		helpers.ParseArgs(cmd, os.Args[2:], "link-address", "wrapper-address")

		v2plusscripts.WrapperConsumerDeploy(e,
			common.HexToAddress(*linkAddress),
			common.HexToAddress(*wrapperAddress))
	case "wrapper-consumer-request":
		cmd := flag.NewFlagSet("wrapper-consumer-request", flag.ExitOnError)
		consumerAddress := cmd.String("consumer-address", "", "address of wrapper consumer")
		cbGasLimit := cmd.Uint("cb-gas-limit", 100_000, "request callback gas limit")
		confirmations := cmd.Uint("request-confirmations", 3, "request confirmations")
		numWords := cmd.Uint("num-words", 1, "num words to request")
		helpers.ParseArgs(cmd, os.Args[2:], "consumer-address")

		consumer, err := vrfv2plus_wrapper_consumer_example.NewVRFV2PlusWrapperConsumerExample(
			common.HexToAddress(*consumerAddress), e.Ec)
		helpers.PanicErr(err)

		tx, err := consumer.MakeRequest(e.Owner, uint32(*cbGasLimit), uint16(*confirmations), uint32(*numWords))
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
	case "wrapper-consumer-request-status":
		cmd := flag.NewFlagSet("wrapper-consumer-request-status", flag.ExitOnError)
		consumerAddress := cmd.String("consumer-address", "", "address of wrapper consumer")
		requestID := cmd.String("request-id", "", "request id of vrf request")
		helpers.ParseArgs(cmd, os.Args[2:], "consumer-address", "request-id")

		consumer, err := vrfv2plus_wrapper_consumer_example.NewVRFV2PlusWrapperConsumerExample(
			common.HexToAddress(*consumerAddress), e.Ec)
		helpers.PanicErr(err)

		status, err := consumer.GetRequestStatus(nil, decimal.RequireFromString(*requestID).BigInt())
		helpers.PanicErr(err)

		statusStringer := func(status vrfv2plus_wrapper_consumer_example.GetRequestStatus) string {
			return fmt.Sprint("paid (juels):", status.Paid.String(),
				", fulfilled?:", status.Fulfilled,
				", random words:", status.RandomWords)
		}

		fmt.Println("status for request", *requestID, "is:")
		fmt.Println(statusStringer(status))
	case "wrapper-consumer-withdraw-link":
		cmd := flag.NewFlagSet("wrapper-consumer-withdraw-link", flag.ExitOnError)
		consumerAddress := cmd.String("consumer-address", "", "address of wrapper consumer")
		linkAddress := cmd.String("link-address", "", "address of link token")
		helpers.ParseArgs(cmd, os.Args[2:], "consumer-address")
		consumer, err := vrfv2plus_wrapper_consumer_example.NewVRFV2PlusWrapperConsumerExample(
			common.HexToAddress(*consumerAddress), e.Ec)
		helpers.PanicErr(err)
		link, err := link_token_interface.NewLinkToken(common.HexToAddress(*linkAddress), e.Ec)
		helpers.PanicErr(err)
		balance, err := link.BalanceOf(nil, common.HexToAddress(*consumerAddress))
		helpers.PanicErr(err)
		tx, err := consumer.WithdrawLink(e.Owner, balance)
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID,
			"withdrawing", balance.String(), "juels from", *consumerAddress, "to", e.Owner.From.Hex())
	case "transfer-link":
		cmd := flag.NewFlagSet("transfer-link", flag.ExitOnError)
		linkAddress := cmd.String("link-address", "", "address of link token")
		amountJuels := cmd.String("amount-juels", "0", "amount in juels to fund")
		receiverAddress := cmd.String("receiver-address", "", "address of receiver (contract or eoa)")
		helpers.ParseArgs(cmd, os.Args[2:], "amount-juels", "link-address", "receiver-address")
		link, err := link_token_interface.NewLinkToken(common.HexToAddress(*linkAddress), e.Ec)
		helpers.PanicErr(err)
		tx, err := link.Transfer(e.Owner, common.HexToAddress(*receiverAddress), decimal.RequireFromString(*amountJuels).BigInt())
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "transfer", *amountJuels, "juels to", *receiverAddress)
	case "latest-block-header":
		cmd := flag.NewFlagSet("latest-block-header", flag.ExitOnError)
		blockNumber := cmd.Int("block-number", -1, "block number")
		helpers.ParseArgs(cmd, os.Args[2:])
		_ = helpers.CalculateLatestBlockHeader(e, *blockNumber)
	case "wrapper-universe-deploy":
		v2plusscripts.DeployWrapperUniverse(e)
	default:
		panic("unrecognized subcommand: " + os.Args[1])
	}
}

func parseSubID(subID string) *big.Int {
	parsedSubID, ok := new(big.Int).SetString(subID, 10)
	if !ok {
		helpers.PanicErr(fmt.Errorf("sub ID %s cannot be parsed", subID))
	}
	return parsedSubID
}
