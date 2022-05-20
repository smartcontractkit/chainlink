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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/batch_vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_single_consumer_example"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/libocr/gethwrappers/link_token_interface"
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
	case "deploy-bhs-coordinator-consumer":
		coordinatorDeployCmd := flag.NewFlagSet("full-deploy", flag.ExitOnError)
		linkAddress := coordinatorDeployCmd.String("link-address", "", "address of link token")
		linkEthAddress := coordinatorDeployCmd.String("link-eth-feed", "", "address of link eth feed")
		fallbackWeiPerUnitLink := coordinatorDeployCmd.String("fallback-wei-per-unit-link", "", "fallback wei/link ratio")
		registerKeyUncompressedPubKey := coordinatorDeployCmd.String("uncompressed-pub-key", "", "uncompressed public key")
		registerKeyOracleAddress := coordinatorDeployCmd.String("oracle-address", "", "oracle address")
		keyHash := coordinatorDeployCmd.String("key-hash", "", "key hash")
		subscriptionBalanceString := coordinatorDeployCmd.String("subscription-balance", "", "subscription balance")
		helpers.ParseArgs(
			coordinatorDeployCmd, os.Args[2:],
			"link-address",
			"link-eth-feed",
			"fallback-wei-per-unit-link",
			"uncompressed-pub-key",
			"oracle-address",
			"key-hash",
			"subscription-balance",
		)

		subscriptionBalance, success := big.NewInt(0).SetString(*subscriptionBalanceString, 10)
		if !success {
			panic(fmt.Sprintf("failed to parse subscriptionBalance '%s'", *subscriptionBalanceString))
		}

		// Deploy BlockhashStore.
		fmt.Println("\nDeploying BHS...")
		_, tx, _, err := blockhash_store.DeployBlockhashStore(owner, ec)
		helpers.PanicErr(err)
		bhsContractAddress := confirmContractDeployed(context.Background(), ec, tx, chainID)

		// Deploy VRFCoordinatorV2, set config, and register proving key.
		fmt.Println("\nDeploying Coordinator...")
		_, tx, _, err = vrf_coordinator_v2.DeployVRFCoordinatorV2(
			owner,
			ec,
			common.HexToAddress(*linkAddress),
			bhsContractAddress,
			common.HexToAddress(*linkEthAddress))
		helpers.PanicErr(err)
		coordinatorAddress := confirmContractDeployed(context.Background(), ec, tx, chainID)

		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(coordinatorAddress, ec)
		helpers.PanicErr(err)

		fmt.Println("\nSetting Config...")
		tx, err = coordinator.SetConfig(owner,
			uint16(3),     // minRequestConfirmations
			uint32(2.5e6), // max gas limit
			uint32(86400), // stalenessSeconds
			uint32(33285), // gasAfterPaymentCalculation
			decimal.RequireFromString(*fallbackWeiPerUnitLink).BigInt(), // 0.01 eth per link fallbackLinkPrice
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
		helpers.PanicErr(err)
		confirmTXMined(context.Background(), ec, tx, chainID)

		fmt.Println("\nConfig set, getting current config from deployed contract...")
		cfg, err := coordinator.GetConfig(nil)
		helpers.PanicErr(err)
		feeConfig, err := coordinator.GetFeeConfig(nil)
		helpers.PanicErr(err)
		fmt.Printf("Config: %+v\n", cfg)
		fmt.Printf("Fee config: %+v\n", feeConfig)

		fmt.Println("\nRegistering proving key...")
		if strings.HasPrefix(*registerKeyUncompressedPubKey, "0x") {
			*registerKeyUncompressedPubKey = strings.Replace(*registerKeyUncompressedPubKey, "0x", "04", 1)
		}
		pubBytes, err := hex.DecodeString(*registerKeyUncompressedPubKey)
		helpers.PanicErr(err)
		pk, err := crypto.UnmarshalPubkey(pubBytes)
		helpers.PanicErr(err)
		tx, err = coordinator.RegisterProvingKey(owner,
			common.HexToAddress(*registerKeyOracleAddress),
			[2]*big.Int{pk.X, pk.Y})
		helpers.PanicErr(err)
		confirmTXMined(
			context.Background(),
			ec,
			tx,
			chainID,
			fmt.Sprintf("Uncompressed public key: %s,", *registerKeyUncompressedPubKey),
			fmt.Sprintf("Oracle address: %s,", *registerKeyOracleAddress),
		)

		fmt.Println("\nProving key registered, getting proving key hashes from deployed contract...")
		_, _, s_provingKeyHashes, err := coordinator.GetRequestConfig(nil)
		helpers.PanicErr(err)
		fmt.Printf("Hashes: %+v\n", s_provingKeyHashes)

		// Deploy consumer and subscribe.
		fmt.Println("\nDeploying consumer and subscribing...")
		keyHashBytes := common.HexToHash(*keyHash)
		_, tx, _, err = vrf_single_consumer_example.DeployVRFSingleConsumerExample(
			owner,
			ec,
			coordinatorAddress,
			common.HexToAddress(*linkAddress),
			uint32(1000000), // gas callback
			uint16(5),       // confs
			uint32(1),       // words
			keyHashBytes)
		helpers.PanicErr(err)
		consumerAddress := confirmContractDeployed(context.Background(), ec, tx, chainID)
		subId := uint64(1)

		fmt.Println("\nFunding subscription...")
		b, err := utils.GenericEncode([]string{"uint64"}, subId)
		helpers.PanicErr(err)
		linkToken, err := link_token_interface.NewLinkToken(common.HexToAddress(*linkAddress), ec)
		helpers.PanicErr(err)
		tx, err = linkToken.TransferAndCall(owner, coordinator.Address(), subscriptionBalance, b)
		helpers.PanicErr(err)
		confirmTXMined(context.Background(), ec, tx, chainID)

		fmt.Println("\nSubscribed and funded, retrieving subscription from deployed contract...")
		s, err := coordinator.GetSubscription(nil, subId)
		helpers.PanicErr(err)
		fmt.Printf("Subscription %+v\n", s)
		fmt.Println(
			"\nDeployment complete.",
			"\nBlockhash Store contract address:", bhsContractAddress,
			"\nVRF Coordinator Address:", coordinatorAddress,
			"\nVRF Consumer Address:", consumerAddress,
			"\nVRF Subscription Id:", subId,
			"\nVRF Subscription Balance:", *subscriptionBalanceString,
			"\nA node can now be configured to run a VRF job with the above configuration.",
		)
	default:
		panic("unrecognized subcommand: " + os.Args[1])
	}
}

func confirmTXMined(context context.Context, client *ethclient.Client, transaction *types.Transaction, chainID int64, txInfo ...string) {
	fmt.Println("Executing TX", helpers.ExplorerLink(chainID, transaction.Hash()), txInfo)
	receipt, err := bind.WaitMined(context, client, transaction)
	helpers.PanicErr(err)
	fmt.Println("TX", receipt.TxHash, "mined. \nBlock Number:", receipt.BlockNumber, "\nGas Used: ", receipt.GasUsed)
}

func confirmContractDeployed(context context.Context, client *ethclient.Client, transaction *types.Transaction, chainID int64) (address common.Address) {
	fmt.Println("Executing contract deployment, TX:", helpers.ExplorerLink(chainID, transaction.Hash()))
	contractAddress, err := bind.WaitDeployed(context, client, transaction)
	helpers.PanicErr(err)
	fmt.Println("Contract Address:", contractAddress.String())
	return contractAddress
}
