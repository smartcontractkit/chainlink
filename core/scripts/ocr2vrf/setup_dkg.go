package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/urfave/cli"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
)

func setupDKGNodes(e helpers.Environment) {
	client := newSetupClient()
	app := cmd.NewApp(client)

	cmd := flag.NewFlagSet("dkg-setup", flag.ExitOnError)
	keyID := cmd.String("key-id", "aee00d81f822f882b6fe28489822f59ebb21ea95c0ae21d9f67c0239461148fc", "key ID")
	apiFile := cmd.String("api", "../../../tools/secrets/apicredentials", "api credentials file")
	passwordFile := cmd.String("password", "../../../tools/secrets/password.txt", "password file")
	databasePrefix := cmd.String("database-prefix", "postgres://postgres:postgres@localhost:5432/dkg-test", "database prefix")
	databaseSuffixes := cmd.String("database-suffixes", "sslmode=disable", "database parameters to be added")
	nodeCount := cmd.Int("node-count", 6, "number of nodes")
	fundingAmount := cmd.Int64("funding-amount", 10000000000000000, "amount to fund nodes") // .1 ETH
	helpers.ParseArgs(cmd, os.Args[2:])

	if *nodeCount < 6 {
		fmt.Println("Node count too low for DKG job, need at least 6.")
		os.Exit(1)
	}

	//Deploy DKG contract.
	// uncomment for faster txs
	// e.Owner.GasPrice = e.Owner.GasPrice.Mul(e.Owner.GasPrice, big.NewInt(2))
	dkgAddress := deployDKG(e).String()

	// Initialize dkg-set-config arguments.
	onChainPublicKeys := []string{}
	offChainPublicKeys := []string{}
	configPublicKeys := []string{}
	peerIDs := []string{}
	transmitters := []string{}
	dkgEncrypters := []string{}
	dkgSigners := []string{}

	// Iterate through all nodes and create jobs.
	for i := 0; i < *nodeCount; i++ {
		flagSet := flag.NewFlagSet("run-dkg-job-creation", flag.ExitOnError)
		flagSet.String("api", *apiFile, "api file")
		flagSet.String("password", *passwordFile, "password file")
		flagSet.String("bootstrapPort", fmt.Sprintf("%d", 8000), "port of bootstrap")
		flagSet.String("job-type", string(jobTypeDKG), "the job type")
		flagSet.String("keyID", *keyID, "")
		flagSet.String("contractID", dkgAddress, "the contract address of the DKG")
		flagSet.Int64("chainID", e.ChainID, "the chain ID")
		flagSet.Bool("dangerWillRobinson", true, "for resetting databases")
		flagSet.Bool("isBootstrapper", i == 0, "is first node")
		bootstrapperPeerID := ""
		if len(peerIDs) != 0 {
			bootstrapperPeerID = peerIDs[0]
		}
		flagSet.String("bootstrapperPeerID", bootstrapperPeerID, "peerID of first node")

		// Create context from flags.
		context := cli.NewContext(app, flagSet, nil)

		// Set environment variables needed to set up DKG jobs.
		configureEnvironmentVariables(false, i, *databasePrefix, *databaseSuffixes)

		// Reset DKG node database.
		resetDatabase(client, context)

		// Setup DKG node.
		payload := setupOCR2VRFNodeFromClient(client, context, e)

		// Append arguments for dkg-set-config command.
		onChainPublicKeys = append(onChainPublicKeys, payload.OnChainPublicKey)
		offChainPublicKeys = append(offChainPublicKeys, payload.OffChainPublicKey)
		configPublicKeys = append(configPublicKeys, payload.ConfigPublicKey)
		peerIDs = append(peerIDs, payload.PeerID)
		transmitters = append(transmitters, payload.Transmitter)
		dkgEncrypters = append(dkgEncrypters, payload.DkgEncrypt)
		dkgSigners = append(dkgSigners, payload.DkgSign)
	}

	// Fund transmitters with funding amount.
	helpers.FundNodes(e, transmitters, big.NewInt(*fundingAmount))

	// Construct and print dkg-set-config command.
	fmt.Println("Generated setConfig Command:")
	command := fmt.Sprintf(
		"go run . dkg-set-config --dkg-address %s -key-id %s -onchain-pub-keys %s -offchain-pub-keys %s -config-pub-keys %s -peer-ids %s -transmitters %s -dkg-encryption-pub-keys %s -dkg-signing-pub-keys %s -schedule 1,1,1,1,1",
		dkgAddress,
		*keyID,
		strings.Join(onChainPublicKeys[1:], ","),
		strings.Join(offChainPublicKeys[1:], ","),
		strings.Join(configPublicKeys[1:], ","),
		strings.Join(peerIDs[1:], ","),
		strings.Join(transmitters[1:], ","),
		strings.Join(dkgEncrypters[1:], ","),
		strings.Join(dkgSigners[1:], ","),
	)

	fmt.Println(command)
}
