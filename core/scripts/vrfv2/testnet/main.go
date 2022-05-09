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

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/sqlx"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/batch_vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keepers_vrf_consumer"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_external_sub_owner_example"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_load_test_external_sub_owner"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_single_consumer_example"
	"github.com/smartcontractkit/chainlink/core/logger"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/vrf/proof"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	batchCoordinatorV2ABI = evmtypes.MustGetABI(batch_vrf_coordinator_v2.BatchVRFCoordinatorV2ABI)
)

type logconfig struct{}

func (c logconfig) LogSQL() bool {
	return false
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

	// Uncomment the block below if transactions are not getting picked up due to nonce issues:
	//
	//block, err := ec.BlockNumber(context.Background())
	//helpers.PanicErr(err)
	//
	//nonce, err := ec.NonceAt(context.Background(), owner.From, big.NewInt(int64(block)))
	//helpers.PanicErr(err)
	//
	//owner.Nonce = big.NewInt(int64(nonce))
	//owner.GasPrice = gp.Mul(gp, big.NewInt(2))

	switch os.Args[1] {
	case "keepers-vrf-consumer-deploy":
		cmd := flag.NewFlagSet("keepers-vrf-consumer-deploy", flag.ExitOnError)
		coordinatorAddress := cmd.String("coordinator-address", "", "vrf coordinator v2 address")
		subID := cmd.Uint64("sub-id", 0, "subscription id")
		keyHash := cmd.String("key-hash", "", "vrf v2 key hash")
		requestConfs := cmd.Uint("request-confs", 3, "request confirmations")
		upkeepIntervalSeconds := cmd.Int64("upkeep-interval-seconds", 600, "upkeep interval in seconds")
		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address", "sub-id", "key-hash")
		_, tx, _, err := keepers_vrf_consumer.DeployKeepersVRFConsumer(
			owner, ec,
			common.HexToAddress(*coordinatorAddress), // vrf coordinator address
			*subID,                                   // subscription id
			common.HexToHash(*keyHash),               // key hash
			uint16(*requestConfs),                    // request confirmations
			big.NewInt(*upkeepIntervalSeconds),       // upkeep interval seconds
		)
		helpers.PanicErr(err)
		keepersVrfConsumer, err := bind.WaitDeployed(context.Background(), ec, tx)
		helpers.PanicErr(err)
		fmt.Println("Deploy tx:", helpers.ExplorerLink(chainID, tx.Hash()))
		fmt.Println("Keepers vrf consumer:", keepersVrfConsumer.Hex())
	case "batch-coordinatorv2-deploy":
		cmd := flag.NewFlagSet("batch-coordinatorv2-deploy", flag.ExitOnError)
		coordinatorAddr := cmd.String("coordinator-address", "", "address of the vrf coordinator v2 contract")
		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address")
		batchCoordinatorAddress, tx, _, err := batch_vrf_coordinator_v2.DeployBatchVRFCoordinatorV2(owner, ec, common.HexToAddress(*coordinatorAddr))
		helpers.PanicErr(err)
		fmt.Println("BatchVRFCoordinatorV2:", batchCoordinatorAddress.Hex(), "tx:", helpers.ExplorerLink(chainID, tx.Hash()))
	case "batch-coordinatorv2-fulfill":
		cmd := flag.NewFlagSet("batch-coordinatorv2-fulfill", flag.ExitOnError)
		batchCoordinatorAddr := cmd.String("batch-coordinator-address", "", "address of the batch vrf coordinator v2 contract")
		pubKeyHex := cmd.String("pubkeyhex", "", "compressed pubkey hex")
		dbURL := cmd.String("db-url", "", "postgres database url")
		keystorePassword := cmd.String("keystore-pw", "", "password to the keystore")
		submit := cmd.Bool("submit", false, "whether to submit the fulfillments or not")
		estimateGas := cmd.Bool("estimate-gas", false, "whether to estimate gas or not")

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

		preSeedSlice := parseIntSlice(*preSeeds)
		bhSlice := parseHashSlice(*blockHashes)
		blockNumSlice := parseIntSlice(*blockNums)
		subIDSlice := parseIntSlice(*subIDs)
		cbLimitsSlice := parseIntSlice(*cbGasLimits)
		numWordsSlice := parseIntSlice(*numWordses)
		senderSlice := parseAddressSlice(*senders)

		batchCoordinator, err := batch_vrf_coordinator_v2.NewBatchVRFCoordinatorV2(common.HexToAddress(*batchCoordinatorAddr), ec)
		helpers.PanicErr(err)

		db := sqlx.MustOpen("postgres", *dbURL)
		lggr, _ := logger.NewLogger()

		keyStore := keystore.New(db, utils.DefaultScryptParams, lggr, logconfig{})
		err = keyStore.Unlock(*keystorePassword)
		helpers.PanicErr(err)

		k, err := keyStore.VRF().Get(*pubKeyHex)
		helpers.PanicErr(err)

		fmt.Println("vrf key found:", k)

		proofs := []batch_vrf_coordinator_v2.VRFTypesProof{}
		reqCommits := []batch_vrf_coordinator_v2.VRFTypesRequestCommitment{}
		for i := range preSeedSlice {
			ps, err := proof.BigToSeed(preSeedSlice[i])
			helpers.PanicErr(err)
			preSeedData := proof.PreSeedDataV2{
				PreSeed:          ps,
				BlockHash:        bhSlice[i],
				BlockNum:         blockNumSlice[i].Uint64(),
				SubId:            subIDSlice[i].Uint64(),
				CallbackGasLimit: uint32(cbLimitsSlice[i].Uint64()),
				NumWords:         uint32(numWordsSlice[i].Uint64()),
				Sender:           senderSlice[i],
			}
			fmt.Printf("preseed data iteration %d: %+v\n", i, preSeedData)
			finalSeed := proof.FinalSeedV2(preSeedData)

			p, err := keyStore.VRF().GenerateProof(*pubKeyHex, finalSeed)
			helpers.PanicErr(err)

			onChainProof, rc, err := proof.GenerateProofResponseFromProofV2(p, preSeedData)
			helpers.PanicErr(err)

			proofs = append(proofs, batch_vrf_coordinator_v2.VRFTypesProof(onChainProof))
			reqCommits = append(reqCommits, batch_vrf_coordinator_v2.VRFTypesRequestCommitment(rc))
		}

		fmt.Printf("proofs: %+v\n\n", proofs)
		fmt.Printf("request commitments: %+v\n\n", reqCommits)

		if *submit {
			fmt.Println("submitting fulfillments...")
			tx, err := batchCoordinator.FulfillRandomWords(owner, proofs, reqCommits)
			helpers.PanicErr(err)

			fmt.Println("waiting for it to mine:", helpers.ExplorerLink(chainID, tx.Hash()))
			_, err = bind.WaitMined(context.Background(), ec, tx)
			helpers.PanicErr(err)
			fmt.Println("done")
		}

		if *estimateGas {
			fmt.Println("estimating gas")
			payload, err := batchCoordinatorV2ABI.Pack("fulfillRandomWords", proofs, reqCommits)
			helpers.PanicErr(err)

			a := batchCoordinator.Address()
			gasEstimate, err := ec.EstimateGas(context.Background(), ethereum.CallMsg{
				From: owner.From,
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

		preSeed := cmd.String("preseed", "", "request preSeed")
		blockHash := cmd.String("blockhash", "", "request blockhash")
		blockNum := cmd.Uint64("blocknum", 0, "request blocknumber")
		subID := cmd.Uint64("subid", 0, "request subid")
		cbGasLimit := cmd.Uint("cbgaslimit", 0, "request callback gas limit")
		numWords := cmd.Uint("numwords", 0, "request num words")
		sender := cmd.String("sender", "", "request sender")

		helpers.ParseArgs(cmd, os.Args[2:],
			"coordinator-address", "pubkeyhex", "db-url",
			"keystore-pw", "preseed", "blockhash", "blocknum",
			"subid", "cbgaslimit", "numwords", "sender",
		)

		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*coordinatorAddr), ec)
		helpers.PanicErr(err)

		db := sqlx.MustOpen("postgres", *dbURL)
		lggr, _ := logger.NewLogger()

		keyStore := keystore.New(db, utils.DefaultScryptParams, lggr, logconfig{})
		err = keyStore.Unlock(*keystorePassword)
		helpers.PanicErr(err)

		k, err := keyStore.VRF().Get(*pubKeyHex)
		helpers.PanicErr(err)

		fmt.Println("vrf key found:", k)

		ps, err := proof.BigToSeed(decimal.RequireFromString(*preSeed).BigInt())
		helpers.PanicErr(err)
		preSeedData := proof.PreSeedDataV2{
			PreSeed:          ps,
			BlockHash:        common.HexToHash(*blockHash),
			BlockNum:         *blockNum,
			SubId:            *subID,
			CallbackGasLimit: uint32(*cbGasLimit),
			NumWords:         uint32(*numWords),
			Sender:           common.HexToAddress(*sender),
		}
		fmt.Printf("preseed data: %+v\n", preSeedData)
		finalSeed := proof.FinalSeedV2(preSeedData)

		p, err := keyStore.VRF().GenerateProof(*pubKeyHex, finalSeed)
		helpers.PanicErr(err)

		onChainProof, rc, err := proof.GenerateProofResponseFromProofV2(p, preSeedData)
		helpers.PanicErr(err)

		fmt.Printf("Proof: %+v, commitment: %+v\nSending fulfillment!", onChainProof, rc)

		tx, err := coordinator.FulfillRandomWords(owner, onChainProof, rc)
		helpers.PanicErr(err)

		fmt.Println("waiting for it to mine:", helpers.ExplorerLink(chainID, tx.Hash()))
		_, err = bind.WaitMined(context.Background(), ec, tx)
		helpers.PanicErr(err)
		fmt.Println("done")
	case "batch-bhs-deploy":
		cmd := flag.NewFlagSet("batch-bhs-deploy", flag.ExitOnError)
		bhsAddr := cmd.String("bhs-address", "", "address of the blockhash store contract")
		helpers.ParseArgs(cmd, os.Args[2:], "bhs-address")
		batchBHSAddress, tx, _, err := batch_blockhash_store.DeployBatchBlockhashStore(owner, ec, common.HexToAddress(*bhsAddr))
		helpers.PanicErr(err)
		fmt.Println("BatchBlockhashStore:", batchBHSAddress.Hex(), "tx:", helpers.ExplorerLink(chainID, tx.Hash()))
	case "batch-bhs-store":
		cmd := flag.NewFlagSet("batch-bhs-store", flag.ExitOnError)
		batchAddr := cmd.String("batch-bhs-address", "", "address of the batch bhs contract")
		blockNumbersArg := cmd.String("block-numbers", "", "block numbers to store in a single transaction")
		helpers.ParseArgs(cmd, os.Args[2:], "batch-bhs-address", "block-numbers")
		batchBHS, err := batch_blockhash_store.NewBatchBlockhashStore(common.HexToAddress(*batchAddr), ec)
		helpers.PanicErr(err)
		blockNumbers := parseIntSlice(*blockNumbersArg)
		helpers.PanicErr(err)
		tx, err := batchBHS.Store(owner, blockNumbers)
		helpers.PanicErr(err)
		fmt.Println("Store tx:", helpers.ExplorerLink(chainID, tx.Hash()))
	case "batch-bhs-get":
		cmd := flag.NewFlagSet("batch-bhs-get", flag.ExitOnError)
		batchAddr := cmd.String("batch-bhs-address", "", "address of the batch bhs contract")
		blockNumbersArg := cmd.String("block-numbers", "", "block numbers to store in a single transaction")
		helpers.ParseArgs(cmd, os.Args[2:], "batch-bhs-address", "block-numbers")
		batchBHS, err := batch_blockhash_store.NewBatchBlockhashStore(common.HexToAddress(*batchAddr), ec)
		helpers.PanicErr(err)
		blockNumbers := parseIntSlice(*blockNumbersArg)
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
		batchBHS, err := batch_blockhash_store.NewBatchBlockhashStore(common.HexToAddress(*batchAddr), ec)
		helpers.PanicErr(err)
		blockRange, err := decreasingBlockRange(big.NewInt(*startBlock-1), big.NewInt(*startBlock-*numBlocks-1))
		helpers.PanicErr(err)
		rlpHeaders, err := getRlpHeaders(ec, blockRange)
		helpers.PanicErr(err)
		tx, err := batchBHS.StoreVerifyHeader(owner, blockRange, rlpHeaders)
		helpers.PanicErr(err)
		fmt.Println("storeVerifyHeader(", blockRange, ", ...) tx:", helpers.ExplorerLink(chainID, tx.Hash()))
	case "batch-bhs-backwards":
		cmd := flag.NewFlagSet("batch-bhs-backwards", flag.ExitOnError)
		batchAddr := cmd.String("batch-bhs-address", "", "address of the batch bhs contract")
		startBlock := cmd.Int64("start-block", -1, "block number to start from. Must be in the BHS already.")
		endBlock := cmd.Int64("end-block", -1, "block number to end at. Must be less than startBlock")
		batchSize := cmd.Int64("batch-size", -1, "batch size")
		gasMultiplier := cmd.Int64("gas-price-multiplier", 1, "gas price multiplier to use, defaults to 1 (no multiplication)")
		helpers.ParseArgs(cmd, os.Args[2:], "batch-bhs-address", "start-block", "end-block", "batch-size")

		batchBHS, err := batch_blockhash_store.NewBatchBlockhashStore(common.HexToAddress(*batchAddr), ec)
		helpers.PanicErr(err)

		blockRange, err := decreasingBlockRange(big.NewInt(*startBlock-1), big.NewInt(*endBlock))
		helpers.PanicErr(err)

		for i := 0; i < len(blockRange); i += int(*batchSize) {
			j := i + int(*batchSize)
			if j > len(blockRange) {
				j = len(blockRange)
			}

			// Get suggested gas price and multiply by multiplier on every iteration
			// so we don't have our transaction getting stuck. Need to be as fast as
			// possible.
			gp, err := ec.SuggestGasPrice(context.Background())
			helpers.PanicErr(err)
			owner.GasPrice = new(big.Int).Mul(gp, big.NewInt(*gasMultiplier))

			fmt.Println("using gas price", owner.GasPrice, "wei")

			blockNumbers := blockRange[i:j]
			blockHeaders, err := getRlpHeaders(ec, blockNumbers)
			fmt.Println("storing blockNumbers:", blockNumbers)
			helpers.PanicErr(err)

			tx, err := batchBHS.StoreVerifyHeader(owner, blockNumbers, blockHeaders)
			helpers.PanicErr(err)

			fmt.Println("sent tx:", helpers.ExplorerLink(chainID, tx.Hash()))

			fmt.Println("waiting for it to mine...")
			_, err = bind.WaitMined(context.Background(), ec, tx)
			helpers.PanicErr(err)

			fmt.Println("received receipt, continuing")
		}
		fmt.Println("done")
	case "latest-head":
		h, err := ec.HeaderByNumber(context.Background(), nil)
		helpers.PanicErr(err)
		fmt.Println("latest head number:", h.Number.String())
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
	case "coordinator-get-config":
		cmd := flag.NewFlagSet("coordinator-get-config", flag.ExitOnError)
		setConfigAddress := cmd.String("coordinator-address", "", "coordinator address")
		helpers.ParseArgs(cmd, os.Args[2:], "address")

		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*setConfigAddress), ec)
		helpers.PanicErr(err)

		cfg, err := coordinator.GetConfig(nil)
		helpers.PanicErr(err)

		feeConfig, err := coordinator.GetFeeConfig(nil)
		helpers.PanicErr(err)

		fmt.Printf("config: %+v\n", cfg)
		fmt.Printf("fee config: %+v\n", feeConfig)
	case "coordinator-set-config":
		cmd := flag.NewFlagSet("coordinator-set-config", flag.ExitOnError)
		setConfigAddress := cmd.String("coordinator-address", "", "coordinator address")
		minConfs := cmd.Int("min-confs", 3, "min confs")
		maxGasLimit := cmd.Int64("max-gas-limit", 2.5e6, "max gas limit")
		stalenessSeconds := cmd.Int64("staleness-seconds", 86400, "staleness in seconds")
		gasAfterPayment := cmd.Int64("gas-after-payment", 33285, "gas after payment calculation")
		fallbackWeiPerUnitLink := cmd.String("fallback-wei-per-unit-link", "", "fallback wei per unit link")
		flatFeeTier1 := cmd.Int64("flat-fee-tier-1", 500, "flat fee tier 1")
		flatFeeTier2 := cmd.Int64("flat-fee-tier-2", 500, "flat fee tier 2")
		flatFeeTier3 := cmd.Int64("flat-fee-tier-3", 500, "flat fee tier 3")
		flatFeeTier4 := cmd.Int64("flat-fee-tier-4", 500, "flat fee tier 4")
		flatFeeTier5 := cmd.Int64("flat-fee-tier-5", 500, "flat fee tier 5")
		reqsForTier2 := cmd.Int64("reqs-for-tier-2", 0, "requests for tier 2")
		reqsForTier3 := cmd.Int64("reqs-for-tier-3", 0, "requests for tier 3")
		reqsForTier4 := cmd.Int64("reqs-for-tier-4", 0, "requests for tier 4")
		reqsForTier5 := cmd.Int64("reqs-for-tier-5", 0, "requests for tier 5")

		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address", "fallback-wei-per-unit-link")

		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(*setConfigAddress), ec)
		helpers.PanicErr(err)

		tx, err := coordinator.SetConfig(owner,
			uint16(*minConfs),         // minRequestConfirmations
			uint32(*maxGasLimit),      // max gas limit
			uint32(*stalenessSeconds), // stalenessSeconds
			uint32(*gasAfterPayment),  // gasAfterPaymentCalculation
			decimal.RequireFromString(*fallbackWeiPerUnitLink).BigInt(), // 0.01 eth per link fallbackLinkPrice
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
	case "eoa-load-test-consumer-deploy":
		loadTestConsumerDeployCmd := flag.NewFlagSet("eoa-load-test-consumer-deploy", flag.ExitOnError)
		consumerCoordinator := loadTestConsumerDeployCmd.String("coordinator-address", "", "coordinator address")
		consumerLinkAddress := loadTestConsumerDeployCmd.String("link-address", "", "link-address")
		helpers.ParseArgs(loadTestConsumerDeployCmd, os.Args[2:], "coordinator-address", "link-address")
		consumerAddress, tx, _, err := vrf_load_test_external_sub_owner.DeployVRFLoadTestExternalSubOwner(
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
		fmt.Println("Adding consumer", "TX hash", helpers.ExplorerLink(chainID, txadd.Hash()))
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
		r, err := bind.WaitMined(context.Background(), ec, tx)
		helpers.PanicErr(err)
		fmt.Println("Receipt blocknumber:", r.BlockNumber)
	case "eoa-load-test-read":
		cmd := flag.NewFlagSet("eoa-load-test-read", flag.ExitOnError)
		consumerAddress := cmd.String("consumer-address", "", "consumer address")
		helpers.ParseArgs(cmd, os.Args[2:], "consumer-address")
		consumer, err := vrf_load_test_external_sub_owner.NewVRFLoadTestExternalSubOwner(
			common.HexToAddress(*consumerAddress),
			ec)
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
			ec)
		helpers.PanicErr(err)
		for i := 0; i < int(*runs); i++ {
			tx, err := consumer.RequestRandomWords(owner, *subID, uint16(*requestConfirmations),
				keyHashBytes, uint16(*requests))
			helpers.PanicErr(err)
			fmt.Printf("TX %d: %s\n", i+1, helpers.ExplorerLink(chainID, tx.Hash()))
		}
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
	case "eoa-read":
		cmd := flag.NewFlagSet("eoa-read", flag.ExitOnError)
		consumerAddress := cmd.String("consumer", "", "consumer address")
		helpers.ParseArgs(cmd, os.Args[2:], "consumer")
		consumer, err := vrf_external_sub_owner_example.NewVRFExternalSubOwnerExample(common.HexToAddress(*consumerAddress), ec)
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

func parseIntSlice(arg string) (ret []*big.Int) {
	parts := strings.Split(arg, ",")
	ret = []*big.Int{}
	for _, part := range parts {
		ret = append(ret, decimal.RequireFromString(part).BigInt())
	}
	return ret
}

func parseAddressSlice(arg string) (ret []common.Address) {
	parts := strings.Split(arg, ",")
	ret = []common.Address{}
	for _, part := range parts {
		ret = append(ret, common.HexToAddress(part))
	}
	return
}

func parseHashSlice(arg string) (ret []common.Hash) {
	parts := strings.Split(arg, ",")
	ret = []common.Hash{}
	for _, part := range parts {
		ret = append(ret, common.HexToHash(part))
	}
	return
}

// decreasingBlockRange creates a continugous block range starting with
// block `start` and ending at block `end`.
func decreasingBlockRange(start, end *big.Int) (ret []*big.Int, err error) {
	if start.Cmp(end) == -1 {
		return nil, fmt.Errorf("start (%s) must be greater than end (%s)", start.String(), end.String())
	}
	ret = []*big.Int{}
	for i := new(big.Int).Set(start); i.Cmp(end) >= 0; i.Sub(i, big.NewInt(1)) {
		ret = append(ret, new(big.Int).Set(i))
	}
	return
}

func getRlpHeaders(ec *ethclient.Client, blockNumbers []*big.Int) (headers [][]byte, err error) {
	headers = [][]byte{}
	for _, blockNum := range blockNumbers {
		// Get child block since it's the one that has the parent hash in it's header.
		h, err := ec.HeaderByNumber(
			context.Background(),
			new(big.Int).Set(blockNum).Add(blockNum, big.NewInt(1)),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get header: %+v", err)
		}
		rlpHeader, err := rlp.EncodeToBytes(h)
		if err != nil {
			return nil, fmt.Errorf("failed to encode rlp: %+v", err)
		}
		// Uncomment in case storeVerifyHeader calls are reverting, there may be an issue with the RLP
		// encoding.
		// h2, err := ec.HeaderByNumber(context.Background(), blockNum)
		// if err != nil {
		// 	return nil, fmt.Errorf("failed to get header: %v", err)
		// }
		// fmt.Println("block number:", blockNum, "blockhash:", h2.Hash(), "encoded header of next block:", common.Bytes2Hex(rlpHeader))
		headers = append(headers, rlpHeader)
	}
	return
}
