package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/scripts/vrfv2/testnet/constants"
	"github.com/smartcontractkit/chainlink/core/scripts/vrfv2/testnet/scripts"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"io"
	"math/big"
	"os"
	"strings"

	"github.com/urfave/cli"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	clcmd "github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func newApp(remoteNodeURL string, writer io.Writer) (*clcmd.Shell, *cli.App) {
	prompter := clcmd.NewTerminalPrompter()
	client := &clcmd.Shell{
		Renderer:                       clcmd.RendererJSON{Writer: writer},
		AppFactory:                     clcmd.ChainlinkAppFactory{},
		KeyStoreAuthenticator:          clcmd.TerminalKeyStoreAuthenticator{Prompter: prompter},
		FallbackAPIInitializer:         clcmd.NewPromptingAPIInitializer(prompter),
		Runner:                         clcmd.ChainlinkRunner{},
		PromptingSessionRequestBuilder: clcmd.NewPromptingSessionRequestBuilder(prompter),
		ChangePasswordPrompter:         clcmd.NewChangePasswordPrompter(),
		PasswordPrompter:               clcmd.NewPasswordPrompter(),
	}
	app := clcmd.NewApp(client)
	fs := flag.NewFlagSet("blah", flag.ContinueOnError)
	fs.Bool("json", true, "")
	fs.String("remote-node-url", remoteNodeURL, "")
	helpers.PanicErr(app.Before(cli.NewContext(nil, fs, nil)))
	// overwrite renderer since it's set to stdout after Before() is called
	client.Renderer = clcmd.RendererJSON{Writer: writer}
	return client, app
}

var (
	checkMarkEmoji = "✅"
	xEmoji         = "❌"
	infoEmoji      = "ℹ️"
)

func main() {

	cmd := flag.NewFlagSet("setup-envs", flag.ExitOnError)
	vrfPrimaryNodeURL := cmd.String("vrf-primary-node-url", "http://localhost:6610", "remote node URL")
	credsFile := cmd.String("creds-file", "/Users/iljapavlovs/Desktop/Chainlink/chainlink/core/scripts/vrfv2/testnet/docker/secrets/apicredentials", "Creds to authenticate to the node")

	numEthKeys := cmd.Int("num-eth-keys", 5, "Number of eth keys to create")
	maxGasPriceGwei := cmd.Int("max-gas-price-gwei", -1, "Max gas price gwei of the eth keys")
	numVRFKeys := cmd.Int("num-vrf-keys", 1, "Number of vrf keys to create")
	//helpers.ParseArgs(cmd, os.Args[1:], "remote-node-urls", "creds-file")

	file := *credsFile

	var (
		allETHKeysVRFPrimaryNode []presenters.ETHKeyResource
		allVRFKeysVRFPrimaryNode []presenters.VRFKeyResource
	)

	output := &bytes.Buffer{}
	client, app := newApp(*vrfPrimaryNodeURL, output)

	// login first to establish the session
	fmt.Println("logging in to:", *vrfPrimaryNodeURL)
	loginFs := flag.NewFlagSet("test", flag.ContinueOnError)
	loginFs.String("file", file, "")
	loginFs.Bool("bypass-version-check", true, "")
	loginCtx := cli.NewContext(app, loginFs, nil)
	err := client.RemoteLogin(loginCtx)
	helpers.PanicErr(err)
	output.Reset()
	fmt.Println()

	{
		// check for ETH keys
		err = client.ListETHKeys(&cli.Context{
			App: app,
		})
		helpers.PanicErr(err)
		var ethKeys []presenters.ETHKeyResource
		var newKeys []presenters.ETHKeyResource
		helpers.PanicErr(json.Unmarshal(output.Bytes(), &ethKeys))
		switch {
		case len(ethKeys) >= *numEthKeys:
			fmt.Println(checkMarkEmoji, "found", len(ethKeys), "eth keys on", vrfPrimaryNodeURL)
		case len(ethKeys) < *numEthKeys:
			fmt.Println(xEmoji, "found only", len(ethKeys), "eth keys on", vrfPrimaryNodeURL,
				"; creating", *numEthKeys-len(ethKeys), "more")
			toCreate := *numEthKeys - len(ethKeys)
			for i := 0; i < toCreate; i++ {
				output.Reset()
				var newKey presenters.ETHKeyResource

				flagSet := flag.NewFlagSet("blah", flag.ExitOnError)
				if *maxGasPriceGwei > 0 {
					helpers.PanicErr(flagSet.Set("max-gas-price-gwei", fmt.Sprintf("%d", *maxGasPriceGwei)))
				}
				err = client.CreateETHKey(cli.NewContext(app, flagSet, nil))
				helpers.PanicErr(err)
				helpers.PanicErr(json.Unmarshal(output.Bytes(), &newKey))
				newKeys = append(newKeys, newKey)
			}

			fmt.Println("NEW ETH KEYS:", strings.Join(func() (r []string) {
				for _, k := range newKeys {
					r = append(r, k.Address)
				}
				return
			}(), ", "))
		}
		output.Reset()
		fmt.Println()

		for _, ethKey := range ethKeys {
			allETHKeysVRFPrimaryNode = append(allETHKeysVRFPrimaryNode, ethKey)
		}
		for _, nk := range newKeys {
			allETHKeysVRFPrimaryNode = append(allETHKeysVRFPrimaryNode, nk)
		}
	}

	{
		// check for VRF keys
		err = client.ListVRFKeys(&cli.Context{
			App: app,
		})
		helpers.PanicErr(err)
		var vrfKeys []presenters.VRFKeyResource
		var newKeys []presenters.VRFKeyResource
		helpers.PanicErr(json.Unmarshal(output.Bytes(), &vrfKeys))
		switch {
		case len(vrfKeys) >= *numVRFKeys:
			fmt.Println(checkMarkEmoji, "found", len(vrfKeys), "vrf keys on", vrfPrimaryNodeURL)
		default:
			fmt.Println(xEmoji, "found only", len(vrfKeys), "vrf keys on", vrfPrimaryNodeURL, ", creating",
				*numVRFKeys-len(vrfKeys), "more")
			toCreate := *numVRFKeys - len(vrfKeys)
			for i := 0; i < toCreate; i++ {
				output.Reset()
				var newKey presenters.VRFKeyResource

				err = client.CreateVRFKey(
					cli.NewContext(app, flag.NewFlagSet("blah", flag.ExitOnError), nil))
				helpers.PanicErr(err)
				helpers.PanicErr(json.Unmarshal(output.Bytes(), &newKey))
				newKeys = append(newKeys, newKey)

				fmt.Println("NEW VRF KEYS:", strings.Join(func() (r []string) {
					for _, k := range newKeys {
						r = append(r, k.Uncompressed)
					}
					return
				}(), ", "))
			}
		}

		output.Reset()
		fmt.Println()

		for _, vrfKey := range vrfKeys {
			allVRFKeysVRFPrimaryNode = append(allVRFKeysVRFPrimaryNode, vrfKey)
		}
		for _, nk := range newKeys {
			allVRFKeysVRFPrimaryNode = append(allVRFKeysVRFPrimaryNode, nk)
		}
	}

	var allETHKeysVRFPrimaryNodeString []string
	fmt.Println("------------- NODE INFORMATION -------------")

	for _, ethKey := range allETHKeysVRFPrimaryNode {
		fmt.Println("-----------ETH Key-----------")
		fmt.Println("Address: ", ethKey.Address)
		fmt.Println("MaxGasPriceWei: ", ethKey.MaxGasPriceWei)
		fmt.Println("EthBalance: ", ethKey.EthBalance)
		fmt.Println("NextNonce: ", ethKey.NextNonce)
		fmt.Println("-----------------------------")

		allETHKeysVRFPrimaryNodeString = append(allETHKeysVRFPrimaryNodeString, ethKey.Address)
	}

	fmt.Println()

	for _, vrfKey := range allVRFKeysVRFPrimaryNode {
		fmt.Println("-----------VRF Key-----------")
		fmt.Println("Compressed: ", vrfKey.Compressed)
		fmt.Println("Uncompressed: ", vrfKey.Uncompressed)
		fmt.Println("Hash: ", vrfKey.Hash)
		fmt.Println("-----------------------------")
	}

	e := helpers.SetupEnv(false)

	feeConfig := vrf_coordinator_v2.VRFCoordinatorV2FeeConfig{
		FulfillmentFlatFeeLinkPPMTier1: uint32(constants.FlatFeeTier1),
		FulfillmentFlatFeeLinkPPMTier2: uint32(constants.FlatFeeTier2),
		FulfillmentFlatFeeLinkPPMTier3: uint32(constants.FlatFeeTier3),
		FulfillmentFlatFeeLinkPPMTier4: uint32(constants.FlatFeeTier4),
		FulfillmentFlatFeeLinkPPMTier5: uint32(constants.FlatFeeTier5),
		ReqsForTier2:                   big.NewInt(constants.ReqsForTier2),
		ReqsForTier3:                   big.NewInt(constants.ReqsForTier3),
		ReqsForTier4:                   big.NewInt(constants.ReqsForTier4),
		ReqsForTier5:                   big.NewInt(constants.ReqsForTier5),
	}

	vrfv2PrimaryNodeJob := scripts.VRFV2DeployUniverse(
		e,
		decimal.RequireFromString(constants.FallbackWeiPerUnitLinkString).BigInt(),
		decimal.RequireFromString(constants.SubscriptionBalanceString).BigInt(),
		&allVRFKeysVRFPrimaryNode[0].Uncompressed,
		new(string),
		new(string),
		&constants.MinConfs,
		&constants.MaxGasLimit,
		&constants.StalenessSeconds,
		&constants.GasAfterPayment,
		feeConfig,
		allETHKeysVRFPrimaryNodeString,
		&constants.OracleFundingAmount,
	)

	if err := os.WriteFile("vrf-job-spec.toml", []byte(vrfv2PrimaryNodeJob), 0666); err != nil {
		helpers.PanicErr(err)
	}

	job := presenters.JobResource{}

	flagSet := flag.NewFlagSet("blah", flag.ExitOnError)
	err = flagSet.Parse([]string{"./vrf-job-spec.toml"})
	helpers.PanicErr(err)

	err = client.CreateJob(cli.NewContext(app, flagSet, nil))
	helpers.PanicErr(err)
	helpers.PanicErr(json.Unmarshal(output.Bytes(), &job))
	fmt.Println(job)
}
