package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/cmd"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type jobType string

const (
	jobTypeDKG     jobType = "DKG"
	jobTypeOCR2VRF jobType = "OCR2VRF"
)

func setupOCR2VRFNodes(e helpers.Environment) {
	fs := flag.NewFlagSet("ocr2vrf-setup", flag.ExitOnError)

	keyID := fs.String("key-id", "aee00d81f822f882b6fe28489822f59ebb21ea95c0ae21d9f67c0239461148fc", "key ID")
	linkAddress := fs.String("link-address", "", "LINK token address")
	linkEthFeed := fs.String("link-eth-feed", "", "LINK-ETH feed address")
	useForwarder := fs.Bool("use-forwarder", false, "boolean to use the forwarder")
	confDelays := fs.String("conf-delays", "1,2,3,4,5,6,7,8", "8 confirmation delays")
	lookbackBlocks := fs.Int64("lookback-blocks", 1000, "lookback blocks")
	weiPerUnitLink := fs.String("wei-per-unit-link", "6e16", "wei per unit link price for feed")
	beaconPeriodBlocks := fs.Int64("beacon-period-blocks", 3, "beacon period in blocks")
	subscriptionBalanceString := fs.String("subscription-balance", "1e19", "amount to fund subscription")

	apiFile := fs.String("api", "../../../tools/secrets/apicredentials", "api credentials file")
	passwordFile := fs.String("password", "../../../tools/secrets/password.txt", "password file")
	databasePrefix := fs.String("database-prefix", "postgres://postgres:postgres@localhost:5432/ocr2vrf-test", "database prefix")
	databaseSuffixes := fs.String("database-suffixes", "sslmode=disable", "database parameters to be added")
	nodeCount := fs.Int("node-count", 6, "number of nodes")
	fundingAmount := fs.Int64("funding-amount", 1e17, "amount to fund nodes") // .1 ETH

	helpers.ParseArgs(fs, os.Args[2:])

	if *nodeCount < 6 {
		fmt.Println("Node count too low for OCR2VRF job, need at least 6.")
		os.Exit(1)
	}

	delays := helpers.ParseIntSlice(*confDelays)
	if len(delays) != 8 {
		fmt.Println("confDelays must have a length of 8")
		os.Exit(1)
	}

	var link common.Address
	if *linkAddress == "" {
		link = helpers.DeployLinkToken(e)
	} else {
		link = common.HexToAddress(*linkAddress)
	}

	// Deploy DKG and VRF contracts, and add VRF
	// as a consumer of DKG events.
	fmt.Println("Deploying DKG contract...")
	dkgAddress := deployDKG(e)

	// Deploy a new feed if needed
	var feedAddress common.Address
	if *linkEthFeed == "" {
		fmt.Println("Deploying LINK-ETH feed...")
		feedAddress = helpers.DeployLinkEthFeed(e, *linkAddress, decimal.RequireFromString(*weiPerUnitLink).BigInt())
	} else {
		feedAddress = common.HexToAddress(*linkEthFeed)
	}

	fmt.Println("Deploying VRF coordinator...")
	vrfCoordinatorAddress := deployVRFCoordinator(e, big.NewInt(*beaconPeriodBlocks), link.String(), feedAddress.String())

	fmt.Println("Deploying VRF beacon...")
	vrfBeaconAddress := deployVRFBeacon(e, vrfCoordinatorAddress.String(), link.String(), dkgAddress.String(), *keyID)

	fmt.Println("Adding VRF Beacon as DKG client...")
	addClientToDKG(e, dkgAddress.String(), *keyID, vrfBeaconAddress.String())

	fmt.Println("Adding VRF Beacon as producer in VRF Coordinator")
	setProducer(e, vrfCoordinatorAddress.String(), vrfBeaconAddress.String())

	fmt.Println("Deploying beacon consumer...")
	consumerAddress := deployVRFBeaconCoordinatorConsumer(e, vrfCoordinatorAddress.String(), false, big.NewInt(*beaconPeriodBlocks))

	fmt.Println("Creating subscription...")
	createSubscription(e, vrfCoordinatorAddress.String())
	subID := 1

	fmt.Println("Adding consumer to subscription...")
	addConsumer(e, vrfCoordinatorAddress.String(), consumerAddress.String(), big.NewInt(int64(subID)))

	subscriptionBalance := decimal.RequireFromString(*subscriptionBalanceString).BigInt()
	if subscriptionBalance.Cmp(big.NewInt(0)) > 0 {
		fmt.Println("\nFunding subscription with", subscriptionBalance, "juels...")
		eoaFundSubscription(e, vrfCoordinatorAddress.String(), link.String(), subscriptionBalance, uint64(subID))
	} else {
		fmt.Println("Subscription", subID, "NOT getting funded. You must fund the subscription in order to use it!")
	}

	var forwarderAddresses []common.Address
	var forwarderAddressesStrings []string
	// If using the forwarder, set up a forwarder for each node.
	if *useForwarder {
		fmt.Println("Deploying transaction forwarders...")
		for i := 0; i < *nodeCount-1; i++ {
			// Deploy an authorized forwarder, and add it to the list of forwarders.
			f := deployAuthorizedForwarder(e, link, e.Owner.From)
			forwarderAddresses = append(forwarderAddresses, f)
			forwarderAddressesStrings = append(forwarderAddressesStrings, f.String())
		}
		fmt.Printf("ForwarderAddresses : %v", forwarderAddressesStrings)
	}

	fmt.Println("Deploying batch beacon consumer...")
	loadTestConsumerAddress := deployLoadTestVRFBeaconCoordinatorConsumer(e, vrfCoordinatorAddress.String(), false, big.NewInt(*beaconPeriodBlocks))
	addConsumer(e, vrfCoordinatorAddress.String(), loadTestConsumerAddress.String(), big.NewInt(int64(subID)))

	fmt.Println("Configuring nodes with OCR2VRF jobs...")
	var (
		onChainPublicKeys  []string
		offChainPublicKeys []string
		configPublicKeys   []string
		peerIDs            []string
		transmitters       []string
		dkgEncrypters      []string
		dkgSigners         []string
		sendingKeys        [][]string
	)

	for i := 0; i < *nodeCount; i++ {
		flagSet := flag.NewFlagSet("run-ocr2vrf-job-creation", flag.ExitOnError)
		flagSet.String("api", *apiFile, "api file")
		flagSet.String("password", *passwordFile, "password file")
		flagSet.String("bootstrapPort", fmt.Sprintf("%d", 8000), "port of bootstrap")
		flagSet.Int64("chainID", e.ChainID, "the chain ID")

		flagSet.String("job-type", string(jobTypeOCR2VRF), "the job type")

		// used by bootstrap template instantiation
		flagSet.String("contractID", dkgAddress.String(), "the contract to get peers from")

		// DKG args
		flagSet.String("keyID", *keyID, "")
		flagSet.String("dkg-address", dkgAddress.String(), "the contract address of the DKG")

		// VRF args
		flagSet.String("vrf-beacon-address", vrfBeaconAddress.String(), "the contract address of the VRF Beacon")
		flagSet.String("vrf-coordinator-address", vrfCoordinatorAddress.String(), "the contract address of the VRF Coordinator")
		flagSet.String("link-eth-feed-address", feedAddress.Hex(), "link eth feed address")
		flagSet.Int64("lookback-blocks", *lookbackBlocks, "lookback blocks")
		flagSet.String("confirmation-delays", *confDelays, "confirmation delays")

		// Apply forwarder args if using the forwarder.
		if i > 0 && *useForwarder {
			flagSet.Bool("use-forwarder", *useForwarder, "use a transaction forwarder")
			flagSet.String("forwarder-address", forwarderAddressesStrings[i-1], "transaction forwarder address")
		}

		flagSet.Bool("dangerWillRobinson", true, "for resetting databases")
		flagSet.Bool("isBootstrapper", i == 0, "is first node")
		bootstrapperPeerID := ""
		if len(peerIDs) != 0 {
			bootstrapperPeerID = peerIDs[0]
		}
		flagSet.String("bootstrapperPeerID", bootstrapperPeerID, "peerID of first node")

		payload := setupNode(flagSet, i, *databasePrefix, *databaseSuffixes, *useForwarder)

		onChainPublicKeys = append(onChainPublicKeys, payload.OnChainPublicKey)
		offChainPublicKeys = append(offChainPublicKeys, payload.OffChainPublicKey)
		configPublicKeys = append(configPublicKeys, payload.ConfigPublicKey)
		peerIDs = append(peerIDs, payload.PeerID)
		transmitters = append(transmitters, payload.Transmitter)
		dkgEncrypters = append(dkgEncrypters, payload.DkgEncrypt)
		dkgSigners = append(dkgSigners, payload.DkgSign)
		sendingKeys = append(sendingKeys, payload.SendingKeys)
	}

	var nodesToFund []string

	// If using the forwarder, set up a forwarder for each node.
	if *useForwarder {
		fmt.Println("Setting authorized senders...")
		for i, f := range forwarderAddresses {

			// Convert the sending strings for a transmitter to addresses.
			var sendinKeysAddresses []common.Address
			sendingKeysStrings := sendingKeys[i+1]
			for _, s := range sendingKeysStrings {
				sendinKeysAddresses = append(sendinKeysAddresses, common.HexToAddress(s))
			}

			// Set authorized senders for the corresponding forwarder.
			setAuthorizedSenders(e, f, sendinKeysAddresses)

			// Fund the sending keys.
			nodesToFund = append(nodesToFund, sendingKeysStrings...)

			// Set the authorized forwarder as the OCR transmitter.
			transmitters[i+1] = f.String()
		}
	} else {
		for _, t := range transmitters[1:] {
			nodesToFund = append(nodesToFund, t)
		}
	}

	var payees []common.Address
	var reportTransmitters []common.Address // all transmitters excluding bootstrap
	for _, t := range transmitters[1:] {
		payees = append(payees, e.Owner.From)
		reportTransmitters = append(reportTransmitters, common.HexToAddress(t))
	}

	fmt.Printf("Setting EOA: %s as payee for transmitters: %v \n", e.Owner.From, reportTransmitters)
	setPayees(e, vrfBeaconAddress.String(), reportTransmitters, payees)

	fmt.Println("Funding transmitters...")
	helpers.FundNodes(e, nodesToFund, big.NewInt(*fundingAmount))

	fmt.Println("Generated dkg setConfig command:")
	dkgCommand := fmt.Sprintf(
		"go run . dkg-set-config -dkg-address %s -key-id %s -onchain-pub-keys %s -offchain-pub-keys %s -config-pub-keys %s -peer-ids %s -transmitters %s -dkg-encryption-pub-keys %s -dkg-signing-pub-keys %s -schedule 1,1,1,1,1",
		dkgAddress.String(),
		*keyID,
		strings.Join(onChainPublicKeys[1:], ","),
		strings.Join(offChainPublicKeys[1:], ","),
		strings.Join(configPublicKeys[1:], ","),
		strings.Join(peerIDs[1:], ","),
		strings.Join(transmitters[1:], ","),
		strings.Join(dkgEncrypters[1:], ","),
		strings.Join(dkgSigners[1:], ","),
	)
	fmt.Println(dkgCommand)

	fmt.Println()
	fmt.Println("Generated vrf setConfig command:")
	vrfCommand := fmt.Sprintf(
		"go run . beacon-set-config -beacon-address %s -conf-delays %s -onchain-pub-keys %s -offchain-pub-keys %s -config-pub-keys %s -peer-ids %s -transmitters %s -schedule 1,1,1,1,1",
		vrfBeaconAddress.String(),
		*confDelays,
		strings.Join(onChainPublicKeys[1:], ","),
		strings.Join(offChainPublicKeys[1:], ","),
		strings.Join(configPublicKeys[1:], ","),
		strings.Join(peerIDs[1:], ","),
		strings.Join(transmitters[1:], ","),
	)
	fmt.Println(vrfCommand)

	fmt.Println()
	fmt.Println("Consumer address:", consumerAddress.String())
	fmt.Println("Consumer request command:")
	requestCommand := fmt.Sprintf(
		"go run . consumer-request-randomness -consumer-address %s -sub-id <sub-id>",
		consumerAddress.Hex())
	fmt.Println(requestCommand)
	fmt.Println()

	fmt.Println("Consumer callback request command:")
	callbackCommand := fmt.Sprintf(
		"go run . consumer-request-callback -consumer-address %s -sub-id <sub-id>",
		consumerAddress.Hex())
	fmt.Println(callbackCommand)
	fmt.Println()

	fmt.Println("Consumer callback batch request command:")
	callbackCommand = fmt.Sprintf(
		"go run . consumer-request-callback-batch -consumer-address %s -sub-id <sub-id> -batch-size <batch-size>",
		loadTestConsumerAddress.Hex())
	fmt.Println(callbackCommand)
	fmt.Println()

	fmt.Println("Consumer redeem randomness command:")
	redeemCommand := fmt.Sprintf(
		"go run . consumer-redeem-randomness -consumer-address %s -request-id <req-id>",
		consumerAddress.Hex())
	fmt.Println(redeemCommand)
	fmt.Println()
}

func setupNode(flagSet *flag.FlagSet, nodeIdx int, databasePrefix, databaseSuffixes string, useForwarder bool) *cmd.SetupOCR2VRFNodePayload {
	client := newSetupClient()
	app := cmd.NewApp(client)
	ctx := cli.NewContext(app, flagSet, nil)

	defer func() {
		err := app.After(ctx)
		helpers.PanicErr(err)
	}()

	err := app.Before(ctx)
	helpers.PanicErr(err)

	resetDatabase(client, ctx, nodeIdx, databasePrefix, databaseSuffixes)
	configureEnvironmentVariables((useForwarder) && (nodeIdx > 0))

	return setupOCR2VRFNodeFromClient(client, ctx)
}
