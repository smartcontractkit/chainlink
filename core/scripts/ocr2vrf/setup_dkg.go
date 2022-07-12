package main

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/urfave/cli"
)

func setupDKGNodes(e helpers.Environment) {
	client := NewProductionClient()
	app := cmd.NewApp(client)

	// Parse os arguments.
	cmd := flag.NewFlagSet("dkg-setup", flag.ExitOnError)
	keyID := cmd.String("key-id", "aee00d81f822f882b6fe28489822f59ebb21ea95c0ae21d9f67c0239461148fc", "key ID")
	apiFile := cmd.String("api", "../../../tools/secrets/apicredentials", "api file")
	passwordFile := cmd.String("password", "../../../tools/secrets/password.txt", "password file")
	databasePrefix := cmd.String("database-prefix", "postgres://postgres:postgres@localhost:5432/dkg-test", "database prefix")
	databaseSuffixes := cmd.String("database-suffixes", "sslmode=disable", "database parameters to be added")
	nodeCount := cmd.Int("node-cout", 5, "number of nodes")
	fundingAmount := cmd.Int64("funding-amount", 10000000000000000, "amount to fund nodes") // .1 ETH
	helpers.ParseArgs(cmd, os.Args[2:])

	if *nodeCount < 5 {
		fmt.Println("Node count too low for DKG job.")
		os.Exit(1)
	}

	// Set environment variables needed to set up DKG jobs.
	configureEnvironmentVariables()

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
		flagSet.String("keyID", *keyID, "")
		flagSet.String("contractID", dkgAddress, "the contract address of the DKG")
		flagSet.Int64("chainID", e.ChainID, "the chain ID")
		flagSet.Bool("dangerWillRobinson", true, "for resetting databases")
		flagSet.Bool("isBootstrapper", i == 0, "is first node")
		bootstrapperPeerID := ""
		if len(peerIDs) != 0 {
			bootstrapperPeerID = peerIDs[0]
		}
		flagSet.String("bootstrapperPeerID", bootstrapperPeerID, "is first node")

		// Create context from flags.
		context := cli.NewContext(app, flagSet, nil)

		// Reset DKG node database.
		resetDatbase(client, context, i, *databasePrefix, *databaseSuffixes)

		// Setup DKG node.
		payload := setupDKGNodeFromClient(client, context)

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
	fundNodes(e, transmitters, *fundingAmount)

	// Construct and print dkg-set-config command.
	fmt.Println("Generated setConfig Command:")
	command := fmt.Sprintf(
		"go run . dkg-set-config --dkg-address %s -key-id %s -onchain-pub-keys %s -offchain-pub-keys %s -config-pub-keys %s -peer-ids %s -transmitters %s -dkg-encryption-pub-keys %s -dkg-signing-pub-keys %s -schedule 1,1,1,1,1",
		dkgAddress,
		*keyID,
		strings.Join(onChainPublicKeys, ","),
		strings.Join(offChainPublicKeys, ","),
		strings.Join(configPublicKeys, ","),
		strings.Join(peerIDs, ","),
		strings.Join(transmitters, ","),
		strings.Join(dkgEncrypters, ","),
		strings.Join(dkgSigners, ","),
	)

	fmt.Println(command)
}

func fundNodes(e helpers.Environment, transmitters []string, fundingAmount int64) {
	fmt.Println("Funding Nodes:")

	block, err := e.Ec.BlockNumber(context.Background())
	helpers.PanicErr(err)

	nonce, err := e.Ec.NonceAt(context.Background(), e.Owner.From, big.NewInt(int64(block)))
	helpers.PanicErr(err)

	for i := 0; i < len(transmitters); i++ {
		tx := types.NewTransaction(
			nonce+uint64(i),
			common.HexToAddress(transmitters[i]),
			big.NewInt(fundingAmount),
			uint64(21000),
			e.Owner.GasPrice,
			nil,
		)
		signedTx, err := e.Owner.Signer(e.Owner.From, tx)
		helpers.PanicErr(err)
		err = e.Ec.SendTransaction(context.Background(), signedTx)
		helpers.PanicErr(err)

		fmt.Printf("Sending to %s: %s\n", transmitters[i], helpers.ExplorerLink(e.ChainID, signedTx.Hash()))
	}
}

func setupDKGNodeFromClient(client *cmd.Client, context *cli.Context) *cmd.SetupDKGNodePayload {
	payload, err := client.ConfigureDKGNode(context)
	helpers.PanicErr(err)

	return payload
}

func configureEnvironmentVariables() {
	os.Setenv("FEATURE_OFFCHAIN_REPORTING2", "true")
}

func resetDatbase(client *cmd.Client, context *cli.Context, index int, databasePrefix string, databaseSuffixes string) {
	os.Setenv("DATABASE_URL", fmt.Sprintf("%s-%d?%s", databasePrefix, index, databaseSuffixes))
	client.ResetDatabase(context)
}

func NewProductionClient() *cmd.Client {
	lggr, closeLggr := logger.NewLogger()
	cfg := config.NewGeneralConfig(lggr)

	prompter := cmd.NewTerminalPrompter()
	return &cmd.Client{
		Renderer:                       cmd.RendererTable{Writer: os.Stdout},
		Config:                         cfg,
		Logger:                         lggr,
		CloseLogger:                    closeLggr,
		AppFactory:                     cmd.ChainlinkAppFactory{},
		KeyStoreAuthenticator:          cmd.TerminalKeyStoreAuthenticator{Prompter: prompter},
		FallbackAPIInitializer:         cmd.NewPromptingAPIInitializer(prompter),
		Runner:                         cmd.ChainlinkRunner{},
		PromptingSessionRequestBuilder: cmd.NewPromptingSessionRequestBuilder(prompter),
		ChangePasswordPrompter:         cmd.NewChangePasswordPrompter(),
		PasswordPrompter:               cmd.NewPasswordPrompter(),
	}
}
