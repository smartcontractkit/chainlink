package main

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/link_token_interface"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type commonSetConfigArgs struct {
	onchainPubKeys         string
	offchainPubKeys        string
	configPubKeys          string
	peerIDs                string
	transmitters           string
	schedule               string
	f                      uint
	deltaProgress          time.Duration
	deltaResend            time.Duration
	deltaRound             time.Duration
	deltaGrace             time.Duration
	deltaStage             time.Duration
	maxRounds              uint8
	maxDurationQuery       time.Duration
	maxDurationObservation time.Duration
	maxDurationReport      time.Duration
	maxDurationAccept      time.Duration
	maxDurationTransmit    time.Duration
}

type dkgSetConfigArgs struct {
	commonSetConfigArgs
	dkgEncryptionPubKeys string
	dkgSigningPubKeys    string
	keyID                string
}

type vrfBeaconSetConfigArgs struct {
	commonSetConfigArgs
	confDelays        string
	coordinatorConfig ocr2vrftypes.CoordinatorConfig
}

func main() {
	e := helpers.SetupEnv(false)

	switch os.Args[1] {
	case "dkg-deploy":
		deployDKG(e)
	case "router-deploy":
		deployVRFRouter(e)
	case "coordinator-deploy":
		cmd := flag.NewFlagSet("coordinator-deploy", flag.ExitOnError)
		beaconPeriodBlocks := cmd.Int64("beacon-period-blocks", 1, "beacon period in number of blocks")
		linkAddress := cmd.String("link-address", "", "link contract address")
		linkEthFeed := cmd.String("link-eth-feed", "", "link/eth feed address")
		routerAddress := cmd.String("router-address", "", "router address")
		helpers.ParseArgs(cmd, os.Args[2:], "beacon-period-blocks", "link-address", "link-eth-feed")
		deployVRFCoordinator(e, big.NewInt(*beaconPeriodBlocks), *linkAddress, *linkEthFeed, *routerAddress)
	case "beacon-deploy":
		cmd := flag.NewFlagSet("beacon-deploy", flag.ExitOnError)
		coordinatorAddress := cmd.String("coordinator-address", "", "coordinator contract address")
		linkAddress := cmd.String("link-address", "", "link contract address")
		dkgAddress := cmd.String("dkg-address", "", "dkg contract address")
		keyID := cmd.String("key-id", "", "key ID")
		helpers.ParseArgs(cmd, os.Args[2:], "beacon-deploy", "coordinator-address", "link-address", "dkg-address", "key-id")
		deployVRFBeacon(e, *coordinatorAddress, *linkAddress, *dkgAddress, *keyID)
	case "dkg-add-client":
		cmd := flag.NewFlagSet("dkg-add-client", flag.ExitOnError)
		dkgAddress := cmd.String("dkg-address", "", "DKG contract address")
		keyID := cmd.String("key-id", "", "key ID")
		clientAddress := cmd.String("client-address", "", "client address")
		helpers.ParseArgs(cmd, os.Args[2:], "dkg-address", "key-id", "client-address")
		addClientToDKG(e, *dkgAddress, *keyID, *clientAddress)
	case "dkg-remove-client":
		cmd := flag.NewFlagSet("dkg-add-client", flag.ExitOnError)
		dkgAddress := cmd.String("dkg-address", "", "DKG contract address")
		keyID := cmd.String("key-id", "", "key ID")
		clientAddress := cmd.String("client-address", "", "client address")
		helpers.ParseArgs(cmd, os.Args[2:], "dkg-address", "key-id", "client-address")
		removeClientFromDKG(e, *dkgAddress, *keyID, *clientAddress)
	case "dkg-set-config":
		cmd := flag.NewFlagSet("dkg-set-config", flag.ExitOnError)
		dkgAddress := cmd.String("dkg-address", "", "DKG contract address")
		keyID := cmd.String("key-id", "", "key ID")
		onchainPubKeys := cmd.String("onchain-pub-keys", "", "comma-separated list of OCR on-chain pubkeys")
		offchainPubKeys := cmd.String("offchain-pub-keys", "", "comma-separated list of OCR off-chain pubkeys")
		configPubKeys := cmd.String("config-pub-keys", "", "comma-separated list of OCR config pubkeys")
		peerIDs := cmd.String("peer-ids", "", "comma-separated list of peer IDs")
		transmitters := cmd.String("transmitters", "", "comma-separated list transmitters")
		dkgEncryptionPubKeys := cmd.String("dkg-encryption-pub-keys", "", "comma-separated list of DKG encryption pubkeys")
		dkgSigningPubKeys := cmd.String("dkg-signing-pub-keys", "", "comma-separated list of DKG signing pubkeys")
		schedule := cmd.String("schedule", "", "comma-separted list of transmission schedule")
		f := cmd.Uint("f", 1, "number of faulty oracles")
		deltaProgress := cmd.Duration("delta-progress", 30*time.Second, "duration of delta progress")
		deltaResend := cmd.Duration("delta-resend", 10*time.Second, "duration of delta resend")
		deltaRound := cmd.Duration("delta-round", 10*time.Second, "duration of delta round")
		deltaGrace := cmd.Duration("delta-grace", 20*time.Second, "duration of delta grace")
		deltaStage := cmd.Duration("delta-stage", 20*time.Second, "duration of delta stage")
		maxRounds := cmd.Uint("max-rounds", 3, "maximum number of rounds")
		maxDurationQuery := cmd.Duration("max-duration-query", 10*time.Millisecond, "maximum duration of query")
		maxDurationObservation := cmd.Duration("max-duration-observation", 10*time.Second, "maximum duration of observation method")
		maxDurationReport := cmd.Duration("max-duration-report", 10*time.Second, "maximum duration of report method")
		maxDurationAccept := cmd.Duration("max-duration-accept", 10*time.Millisecond, "maximum duration of shouldAcceptFinalizedReport method")
		maxDurationTransmit := cmd.Duration("max-duration-transmit", 1*time.Second, "maximum duration of shouldTransmitAcceptedReport method")

		helpers.ParseArgs(cmd,
			os.Args[2:],
			"dkg-address",
			"key-id",
			"onchain-pub-keys",
			"offchain-pub-keys",
			"config-pub-keys",
			"peer-ids",
			"transmitters",
			"dkg-encryption-pub-keys",
			"dkg-signing-pub-keys",
			"schedule")

		commands := dkgSetConfigArgs{
			commonSetConfigArgs: commonSetConfigArgs{
				onchainPubKeys:         *onchainPubKeys,
				offchainPubKeys:        *offchainPubKeys,
				configPubKeys:          *configPubKeys,
				peerIDs:                *peerIDs,
				transmitters:           *transmitters,
				schedule:               *schedule,
				f:                      *f,
				deltaProgress:          *deltaProgress,
				deltaResend:            *deltaResend,
				deltaRound:             *deltaRound,
				deltaGrace:             *deltaGrace,
				deltaStage:             *deltaStage,
				maxRounds:              uint8(*maxRounds),
				maxDurationQuery:       *maxDurationQuery,
				maxDurationObservation: *maxDurationObservation,
				maxDurationReport:      *maxDurationReport,
				maxDurationAccept:      *maxDurationAccept,
				maxDurationTransmit:    *maxDurationTransmit,
			},
			dkgEncryptionPubKeys: *dkgEncryptionPubKeys,
			dkgSigningPubKeys:    *dkgSigningPubKeys,
			keyID:                *keyID,
		}

		setDKGConfig(e, *dkgAddress, commands)
	case "beacon-set-config":
		cmd := flag.NewFlagSet("beacon-set-config", flag.ExitOnError)
		beaconAddress := cmd.String("beacon-address", "", "VRF beacon contract address")
		confDelays := cmd.String("conf-delays", "1,2,3,4,5,6,7,8", "comma-separted list of 8 confirmation delays")
		onchainPubKeys := cmd.String("onchain-pub-keys", "", "comma-separated list of OCR on-chain pubkeys")
		offchainPubKeys := cmd.String("offchain-pub-keys", "", "comma-separated list of OCR off-chain pubkeys")
		configPubKeys := cmd.String("config-pub-keys", "", "comma-separated list of OCR config pubkeys")
		peerIDs := cmd.String("peer-ids", "", "comma-separated list of peer IDs")
		transmitters := cmd.String("transmitters", "", "comma-separated list transmitters")
		schedule := cmd.String("schedule", "", "comma-separted list of transmission schedule")
		f := cmd.Uint("f", 1, "number of faulty oracles")
		// TODO: Adjust default delta* and maxDuration* values below after benchmarking latency
		deltaProgress := cmd.Duration("delta-progress", 30*time.Second, "duration of delta progress")
		deltaResend := cmd.Duration("delta-resend", 10*time.Second, "duration of delta resend")
		deltaRound := cmd.Duration("delta-round", 10*time.Second, "duration of delta round")
		deltaGrace := cmd.Duration("delta-grace", 20*time.Second, "duration of delta grace")
		deltaStage := cmd.Duration("delta-stage", 20*time.Second, "duration of delta stage")
		cacheEvictionWindowSeconds := cmd.Int64("cache-eviction-window", 60, "cache eviction window, in seconds")
		batchGasLimit := cmd.Int64("batch-gas-limit", 5_000_000, "batch gas limit")
		coordinatorOverhead := cmd.Int64("coordinator-overhead", 50_000, "coordinator overhead")
		callbackOverhead := cmd.Int64("callback-overhead", 50_000, "callback overhead")
		blockGasOverhead := cmd.Int64("block-gas-overhead", 50_000, "block gas overhead")
		lookbackBlocks := cmd.Int64("lookback-blocks", 1000, "lookback blocks")
		maxRounds := cmd.Uint("max-rounds", 3, "maximum number of rounds")
		maxDurationQuery := cmd.Duration("max-duration-query", 10*time.Millisecond, "maximum duration of query")
		maxDurationObservation := cmd.Duration("max-duration-observation", 10*time.Second, "maximum duration of observation method")
		maxDurationReport := cmd.Duration("max-duration-report", 10*time.Second, "maximum duration of report method")
		maxDurationAccept := cmd.Duration("max-duration-accept", 5*time.Second, "maximum duration of shouldAcceptFinalizedReport method")
		maxDurationTransmit := cmd.Duration("max-duration-transmit", 1*time.Second, "maximum duration of shouldTransmitAcceptedReport method")

		helpers.ParseArgs(cmd,
			os.Args[2:],
			"beacon-address",
			"onchain-pub-keys",
			"offchain-pub-keys",
			"config-pub-keys",
			"peer-ids",
			"transmitters",
			"schedule")

		commands := vrfBeaconSetConfigArgs{
			commonSetConfigArgs: commonSetConfigArgs{
				onchainPubKeys:         *onchainPubKeys,
				offchainPubKeys:        *offchainPubKeys,
				configPubKeys:          *configPubKeys,
				peerIDs:                *peerIDs,
				transmitters:           *transmitters,
				schedule:               *schedule,
				f:                      *f,
				deltaProgress:          *deltaProgress,
				deltaResend:            *deltaResend,
				deltaRound:             *deltaRound,
				deltaGrace:             *deltaGrace,
				deltaStage:             *deltaStage,
				maxRounds:              uint8(*maxRounds),
				maxDurationQuery:       *maxDurationQuery,
				maxDurationObservation: *maxDurationObservation,
				maxDurationReport:      *maxDurationReport,
				maxDurationAccept:      *maxDurationAccept,
				maxDurationTransmit:    *maxDurationTransmit,
			},
			confDelays: *confDelays,
			coordinatorConfig: ocr2vrftypes.CoordinatorConfig{
				CacheEvictionWindowSeconds: *cacheEvictionWindowSeconds,
				BatchGasLimit:              *batchGasLimit,
				CoordinatorOverhead:        *coordinatorOverhead,
				CallbackOverhead:           *callbackOverhead,
				BlockGasOverhead:           *blockGasOverhead,
				LookbackBlocks:             *lookbackBlocks,
			},
		}

		commands.setVRFBeaconConfig(e, *beaconAddress)
	case "coordinator-set-producer":
		cmd := flag.NewFlagSet("coordinator-set-producer", flag.ExitOnError)
		coordinatorAddress := cmd.String("coordinator-address", "", "VRF coordinator contract address")
		beaconAddress := cmd.String("beacon-address", "", "VRF beacon contract address")
		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address", "beacon-address")
		setProducer(e, *coordinatorAddress, *beaconAddress)
	case "router-request-randomness":
		cmd := flag.NewFlagSet("router-request-randomness", flag.ExitOnError)
		routerAddress := cmd.String("router-address", "", "VRF coordinator contract address")
		numWords := cmd.Uint("num-words", 1, "number of words to request")
		subID := cmd.String("sub-id", "", "subscription ID")
		confDelay := cmd.Int64("conf-delay", 1, "confirmation delay")
		helpers.ParseArgs(cmd, os.Args[2:], "router-address", "sub-id")
		requestRandomness(
			e,
			*routerAddress,
			uint16(*numWords),
			decimal.RequireFromString(*subID).BigInt(),
			big.NewInt(*confDelay))
	case "router-redeem-randomness":
		cmd := flag.NewFlagSet("router-redeem-randomness", flag.ExitOnError)
		routerAddress := cmd.String("router-address", "", "VRF coordinator contract address")
		subID := cmd.String("sub-id", "", "subscription ID")
		requestID := cmd.Int64("request-id", 0, "request ID")
		helpers.ParseArgs(cmd, os.Args[2:], "router-address", "sub-id", "request-id")
		redeemRandomness(
			e,
			*routerAddress,
			decimal.RequireFromString(*subID).BigInt(),
			big.NewInt(*requestID))
	case "beacon-info":
		cmd := flag.NewFlagSet("beacon-info", flag.ExitOnError)
		beaconAddress := cmd.String("beacon-address", "", "VRF beacon contract address")
		helpers.ParseArgs(cmd, os.Args[2:], "beacon-address")
		beacon := newVRFBeacon(common.HexToAddress(*beaconAddress), e.Ec)
		keyID, err := beacon.SKeyID(nil)
		helpers.PanicErr(err)
		fmt.Println("beacon key id:", hexutil.Encode(keyID[:]))
		keyHash, err := beacon.SProvingKeyHash(nil)
		helpers.PanicErr(err)
		fmt.Println("beacon proving key hash:", hexutil.Encode(keyHash[:]))
	case "coordinator-create-sub":
		cmd := flag.NewFlagSet("coordinator-create-sub", flag.ExitOnError)
		coordinatorAddress := cmd.String("coordinator-address", "", "VRF coordinator contract address")
		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address")
		createSubscription(e, *coordinatorAddress)
	case "coordinator-add-consumer":
		cmd := flag.NewFlagSet("coordinator-add-consumer", flag.ExitOnError)
		coordinatorAddress := cmd.String("coordinator-address", "", "VRF coordinator contract address")
		consumerAddress := cmd.String("consumer-address", "", "VRF consumer contract address")
		subId := cmd.String("sub-id", "", "subscription ID")
		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address", "consumer-address")
		addConsumer(e, *coordinatorAddress, *consumerAddress, decimal.RequireFromString(*subId).BigInt())
	case "coordinator-get-sub":
		cmd := flag.NewFlagSet("coordinator-get-sub", flag.ExitOnError)
		coordinatorAddress := cmd.String("coordinator-address", "", "VRF coordinator contract address")
		subId := cmd.String("sub-id", "", "subscription ID")
		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address")
		sub := getSubscription(e, *coordinatorAddress, decimal.RequireFromString(*subId).BigInt())
		fmt.Println("subscription ID:", *subId)
		fmt.Println("balance:", sub.Balance)
		fmt.Println("consumers:", sub.Consumers)
		fmt.Println("owner:", sub.Owner)
		fmt.Println("request count:", sub.ReqCount)
	case "link-balance":
		cmd := flag.NewFlagSet("link-balance", flag.ExitOnError)
		linkAddress := cmd.String("link-address", "", "link address")
		helpers.ParseArgs(cmd, os.Args[2:], "link-address")

		l, err := link_token_interface.NewLinkToken(common.HexToAddress(*linkAddress), e.Ec)
		helpers.PanicErr(err)

		bal, err := l.BalanceOf(nil, e.Owner.From)
		helpers.PanicErr(err)
		fmt.Println("link balance of", e.Owner.From, "is", assets.NewWei(bal).String())
	case "get-balances":
		cmd := flag.NewFlagSet("get-balances", flag.ExitOnError)
		addresses := cmd.String("addresses", "", "comma-separated list of addresses")
		helpers.ParseArgs(cmd, os.Args[2:], "addresses")

		for _, account := range strings.Split(*addresses, ",") {
			bal, err := e.Ec.BalanceAt(context.Background(), common.HexToAddress(account), nil)
			helpers.PanicErr(err)

			fmt.Println("ETH balance of", account, "is", assets.NewWei(bal).String())
		}
	case "coordinator-fund-sub":
		cmd := flag.NewFlagSet("coordinator-fund-sub", flag.ExitOnError)
		coordinatorAddress := cmd.String("coordinator-address", "", "VRF coordinator contract address")
		linkAddress := cmd.String("link-address", "", "link-address")
		fundingAmount := cmd.String("funding-amount", "5e18", "funding amount in juels. can use scientific notation, e.g 10e18 for 10 LINK") // 5 LINK
		subID := cmd.String("sub-id", "", "subscription ID")
		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address", "link-address")
		eoaFundSubscription(e, *coordinatorAddress, *linkAddress, decimal.RequireFromString(*fundingAmount).BigInt(), decimal.RequireFromString(*subID).BigInt())
	case "beacon-set-payees":
		cmd := flag.NewFlagSet("beacon-set-payees", flag.ExitOnError)
		beaconAddress := cmd.String("beacon-address", "", "VRF beacon contract address")
		transmitters := cmd.String("transmitters", "", "comma-separated list of transmitters")
		payees := cmd.String("payees", "", "comma-separated list of payees")
		helpers.ParseArgs(cmd, os.Args[2:], "beacon-address", "transmitters", "payees")
		setPayees(e, *beaconAddress, helpers.ParseAddressSlice(*transmitters), helpers.ParseAddressSlice(*payees))
	case "consumer-deploy":
		cmd := flag.NewFlagSet("consumer-deploy", flag.ExitOnError)
		coordinatorAddress := cmd.String("coordinator-address", "", "VRF coordinator address")
		shouldFail := cmd.Bool("should-fail", false, "shouldFail flag")
		beaconPeriodBlocks := cmd.Int64("beacon-period-blocks", 1, "beacon period in number of blocks")
		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address", "beacon-period-blocks")
		deployVRFBeaconCoordinatorConsumer(e, *coordinatorAddress, *shouldFail, big.NewInt(*beaconPeriodBlocks))
	case "consumer-request-randomness":
		cmd := flag.NewFlagSet("consumer-request-randomness", flag.ExitOnError)
		consumerAddress := cmd.String("consumer-address", "", "VRF coordinator consumer address")
		numWords := cmd.Uint("num-words", 1, "number of words to request")
		subID := cmd.String("sub-id", "", "subscription ID")
		confDelay := cmd.Int64("conf-delay", 1, "confirmation delay")
		helpers.ParseArgs(cmd, os.Args[2:], "consumer-address", "sub-id")
		requestRandomnessFromConsumer(e, *consumerAddress, uint16(*numWords), decimal.RequireFromString(*subID).BigInt(), big.NewInt(*confDelay))
	case "consumer-redeem-randomness":
		cmd := flag.NewFlagSet("consumer-redeem-randomness", flag.ExitOnError)
		consumerAddress := cmd.String("consumer-address", "", "VRF coordinator consumer address")
		subID := cmd.String("sub-id", "", "subscription ID")
		requestID := cmd.Int64("request-id", 0, "request ID")
		numWords := cmd.Int64("num-words", 1, "number of words to print after redeeming")
		helpers.ParseArgs(cmd, os.Args[2:], "consumer-address", "request-id")
		redeemRandomnessFromConsumer(e, *consumerAddress, decimal.RequireFromString(*subID).BigInt(), big.NewInt(*requestID), *numWords)
	case "consumer-request-callback":
		cmd := flag.NewFlagSet("consumer-request-callback", flag.ExitOnError)
		consumerAddress := cmd.String("consumer-address", "", "VRF coordinator consumer address")
		numWords := cmd.Uint("num-words", 1, "number of words to request")
		subID := cmd.String("sub-id", "", "subscription ID")
		confDelay := cmd.Int64("conf-delay", 1, "confirmation delay")
		callbackGasLimit := cmd.Uint("cb-gas-limit", 100_000, "callback gas limit")
		helpers.ParseArgs(cmd, os.Args[2:], "consumer-address")
		requestRandomnessCallback(
			e,
			*consumerAddress,
			uint16(*numWords),
			decimal.RequireFromString(*subID).BigInt(),
			big.NewInt(int64(*confDelay)),
			uint32(*callbackGasLimit),
			nil, // test consumer doesn't use any args
		)
	case "consumer-read-randomness":
		cmd := flag.NewFlagSet("consumer-read-randomness", flag.ExitOnError)
		consumerAddress := cmd.String("consumer-address", "", "VRF coordinator consumer address")
		requestID := cmd.String("request-id", "", "VRF request ID")
		numWords := cmd.Int("num-words", 1, "number of words to fetch")
		helpers.ParseArgs(cmd, os.Args[2:], "consumer-address")
		readRandomness(e, *consumerAddress, decimal.RequireFromString(*requestID).BigInt(), *numWords)
	case "consumer-request-callback-batch":
		cmd := flag.NewFlagSet("consumer-request-callback", flag.ExitOnError)
		consumerAddress := cmd.String("consumer-address", "", "VRF beacon consumer address")
		numWords := cmd.Uint("num-words", 1, "number of words to request")
		subID := cmd.String("sub-id", "", "subscription ID")
		confDelay := cmd.Int64("conf-delay", 1, "confirmation delay")
		batchSize := cmd.Int64("batch-size", 1, "batch size")
		callbackGasLimit := cmd.Uint("cb-gas-limit", 200_000, "callback gas limit")
		helpers.ParseArgs(cmd, os.Args[2:], "consumer-address")
		requestRandomnessCallbackBatch(
			e,
			*consumerAddress,
			uint16(*numWords),
			decimal.RequireFromString(*subID).BigInt(),
			big.NewInt(int64(*confDelay)),
			uint32(*callbackGasLimit),
			nil, // test consumer doesn't use any args,
			big.NewInt(*batchSize),
		)
	case "consumer-request-callback-batch-load-test":
		cmd := flag.NewFlagSet("consumer-request-callback-load-test", flag.ExitOnError)
		consumerAddress := cmd.String("consumer-address", "", "VRF beacon batch consumer address")
		numWords := cmd.Uint("num-words", 1, "number of words to request")
		subID := cmd.String("sub-id", "", "subscription ID")
		confDelay := cmd.Int64("conf-delay", 1, "confirmation delay")
		batchSize := cmd.Int64("batch-size", 1, "batch size")
		batchCount := cmd.Int64("batch-count", 1, "number of batches to run")
		callbackGasLimit := cmd.Uint("cb-gas-limit", 200_000, "callback gas limit")
		helpers.ParseArgs(cmd, os.Args[2:], "consumer-address")

		for i := int64(0); i < *batchCount; i++ {
			requestRandomnessCallbackBatch(
				e,
				*consumerAddress,
				uint16(*numWords),
				decimal.RequireFromString(*subID).BigInt(),
				big.NewInt(int64(*confDelay)),
				uint32(*callbackGasLimit),
				nil, // test consumer doesn't use any args,
				big.NewInt(*batchSize),
			)
		}
	case "deploy-load-test-consumer":
		cmd := flag.NewFlagSet("deploy-load-test-consumer", flag.ExitOnError)
		coordinatorAddress := cmd.String("coordinator-address", "", "coordinator address")
		beaconPeriodBlocks := cmd.Int64("beacon-period-blocks", 1, "beacon period in number of blocks")
		helpers.ParseArgs(cmd, os.Args[2:], "coordinator-address", "beacon-period-blocks")

		deployLoadTestVRFBeaconCoordinatorConsumer(e, *coordinatorAddress, false, big.NewInt(*beaconPeriodBlocks))
	case "get-load-test-results":
		cmd := flag.NewFlagSet("get-load-test-results", flag.ExitOnError)
		consumerAddress := cmd.String("consumer-address", "", "Load test contract address")
		helpers.ParseArgs(cmd, os.Args[2:], "consumer-address")

		printLoadtestResults(e, *consumerAddress)
	case "verify-beacon-randomness":
		cmd := flag.NewFlagSet("verify-randomness", flag.ExitOnError)
		dkgAddress := cmd.String("dkg-address", "", "DKG contract address")
		beaconAddress := cmd.String("beacon-address", "", "VRF beacon contract address")
		coordinatorAddress := cmd.String("coordinator-address", "", "VRF coordinator contract address")
		height := cmd.Uint64("height", 0, "block height of VRF beacon output")
		confDelay := cmd.Uint64("conf-delay", 1, "confirmation delay of VRF beacon output")
		searchWindow := cmd.Uint64("search-window", 200, "search space size for beacon transmission. Number of blocks after beacon height")
		helpers.ParseArgs(cmd, os.Args[2:], "dkg-address", "coordinator-address", "beacon-address", "height", "conf-delay")

		verifyBeaconRandomness(e, *dkgAddress, *beaconAddress, *coordinatorAddress, *height, *confDelay, *searchWindow)
	case "dkg-setup":
		setupDKGNodes(e)
	case "ocr2vrf-setup":
		setupOCR2VRFNodes(e)
	case "ocr2vrf-setup-infra-forwarder":
		setupOCR2VRFNodesForInfraWithForwarder(e)
	case "ocr2vrf-fund-nodes":
		fundOCR2VRFNodes(e)
	default:
		panic("unrecognized subcommand: " + os.Args[1])
	}
}
