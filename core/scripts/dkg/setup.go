package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/recovery"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	dkgContract "github.com/smartcontractkit/ocr2vrf/gethwrappers/dkg"
	"github.com/urfave/cli"
)

func main() {
	recovery.ReportPanics(func() {
		Run(NewProductionClient(), os.Args...)
	})
}

func Run(client *cmd.Client, args ...string) {

	app := cmd.NewApp(client)
	e := helpers.SetupEnv()

	_, tx, _, err := dkgContract.DeployDKG(e.Owner, e.Ec)
	helpers.PanicErr(err)
	dkgAddress := helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID).String()

	onChainPublicKeys := []string{}
	offChainPublicKeys := []string{}
	configPublicKeys := []string{}
	peerIDs := []string{}
	transmitters := []string{}
	dkgEncrypters := []string{}
	dkgSigners := []string{}

	for i := 0; i < 5; i++ {
		flagSet := flag.NewFlagSet("run-dkg-job-creation", flag.ExitOnError)
		flagSet.String("api", "../../../tools/secrets/apicredentials", "api file")
		flagSet.String("password", "../../../tools/secrets/password.txt", "password file")
		flagSet.String("port", fmt.Sprintf("%d", 8000), "port of bootstrap")
		flagSet.String("keyID", "aee00d81f822f882b6fe28489822f59ebb21ea95c0ae21d9f67c0239461148fc", "")
		flagSet.String("contractID", dkgAddress, "the contract address of the DKG")
		flagSet.Bool("dangerWillRobinson", true, "for resetting databases")
		flagSet.Bool("isBootstrapper", i == 0, "is first node")
		bootstrapperPeerID := ""
		if len(peerIDs) != 0 {
			bootstrapperPeerID = peerIDs[0]
		}
		flagSet.String("bootstrapperPeerID", bootstrapperPeerID, "is first node")

		context := cli.NewContext(app, flagSet, nil)

		os.Setenv("DATABASE_URL", fmt.Sprintf("postgres://postgres:postgres@localhost:5432/dkg-test-%d?sslmode=disable", i))
		client.ResetDatabase(context)

		payload, err := client.RunNodeDKG(context)
		newPayload := cmd.RunNodeDKGPayload(*payload)
		if err != nil {
			client.Logger.Error("Error running dkg", err)
		}

		onChainPublicKeys = append(onChainPublicKeys, newPayload.OnChainPublicKey)
		offChainPublicKeys = append(offChainPublicKeys, newPayload.OffChainPublicKey)
		configPublicKeys = append(configPublicKeys, newPayload.ConfigPublicKey)
		peerIDs = append(peerIDs, newPayload.PeerID)
		transmitters = append(transmitters, newPayload.Transmitter)
		dkgEncrypters = append(dkgEncrypters, newPayload.DkgEncrypt)
		dkgSigners = append(dkgSigners, newPayload.DkgSign)
	}

	command := fmt.Sprintf(
		"go run . dkg-set-config --dkg-address %s -onchain-pub-keys %s -offchain-pub-keys %s -config-pub-keys %s -peer-ids %s -transmitters %s -dkg-encryption-pub-keys %s -dkg-signing-pub-keys %s -schedule 1,1,1,1,1",
		dkgAddress,
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
