package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/scripts/vrfv2/testnet/constants"
	"github.com/smartcontractkit/chainlink/core/scripts/vrfv2/testnet/scripts"
	clcmd "github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
	"math/big"
	"os"
	"strings"

	//"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
	"github.com/urfave/cli"
	"io"
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

	vrfPrimaryNodeURL := flag.String("vrf-primary-node-url", "http://localhost:6610", "remote node URL")
	vrfBackupNodeURL := flag.String("vrf-backup-node-url", "", "remote node URL")
	bhsNodeURL := flag.String("bhs-node-url", "", "remote node URL")
	//bhsBackupNodeURL := flag.String("vrf-bhs-backup-node-url", "http://localhost:6613", "remote node URL")
	bhfNodeURL := flag.String("bhf-node-url", "", "remote node URL")
	nodeSendingKeyFundingAmount := flag.Int64("sending-key-funding-amount", constants.NodeSendingKeyFundingAmountGwei, "remote node URL")

	credsFile := flag.String("creds-file", "/Users/iljapavlovs/Desktop/Chainlink/chainlink/core/scripts/vrfv2/testnet/docker/secrets/apicredentials", "Creds to authenticate to the node")

	numEthKeys := flag.Int("num-eth-keys", 5, "Number of eth keys to create")
	maxGasPriceGwei := flag.Int("max-gas-price-gwei", -1, "Max gas price gwei of the eth keys")
	numVRFKeys := flag.Int("num-vrf-keys", 1, "Number of vrf keys to create")

	linkAddress := flag.String("link-address", "", "address of link token")
	linkEthAddress := flag.String("link-eth-feed", "", "address of link eth feed")
	bhsContractAddressString := flag.String("bhs-address", "", "address of BHS contract")
	batchBHSAddressString := flag.String("batch-bhs-address", "", "address of Batch BHS contract")
	coordinatorAddressString := flag.String("coordinator-address", "", "address of VRF Coordinator contract")
	batchCoordinatorAddressString := flag.String("batch-coordinator-address", "", "address Batch VRF Coordinator contract")

	e := helpers.SetupEnv(false)
	flag.Parse()
	file := *credsFile
	nodes := make(map[string]scripts.Node)
	if *vrfPrimaryNodeURL != "" {
		nodes[scripts.VRFPrimaryNodeName] = scripts.Node{
			URL:                     *vrfPrimaryNodeURL,
			SendingKeyFundingAmount: *nodeSendingKeyFundingAmount,
		}
	}

	if *vrfBackupNodeURL != "" {
		nodes[scripts.VRFBackupNodeName] = scripts.Node{
			URL:                     *vrfBackupNodeURL,
			SendingKeyFundingAmount: *nodeSendingKeyFundingAmount,
		}
	}
	if *bhsNodeURL != "" {
		nodes[scripts.BHSNodeName] = scripts.Node{
			URL:                     *bhsNodeURL,
			SendingKeyFundingAmount: *nodeSendingKeyFundingAmount,
		}
	}
	if *bhfNodeURL != "" {
		nodes[scripts.BHFNodeName] = scripts.Node{
			URL:                     *bhfNodeURL,
			SendingKeyFundingAmount: *nodeSendingKeyFundingAmount,
		}
	}

	output := &bytes.Buffer{}

	for key, node := range nodes {

		client, app := connectToNode(&node.URL, output, file)
		ethKeys := createETHKeysIfNeeded(client, app, output, numEthKeys, &node.URL, maxGasPriceGwei)
		if key == scripts.VRFPrimaryNodeName {
			vrfKeys := createVRFKeyIfNeeded(client, app, output, numVRFKeys, &node.URL)
			node.VrfKeys = mapVrfKeysToStringArr(vrfKeys)

			printVRFKeyData(vrfKeys)
		}

		node.SendingKeys = mapEthKeysToStringArr(ethKeys)
		nodes[key] = node
		printETHKeyData(ethKeys)

		if node.SendingKeyFundingAmount > 0 {
			fmt.Println("\nFunding ", key, " Node's Sending Keys...")
			helpers.FundNodes(e, node.SendingKeys, big.NewInt(node.SendingKeyFundingAmount))
		}

	}

	fmt.Println()

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

	vrfv2PrimaryNodeJob, vrfv2BackupNodeJob, bhsNodeJob, bhfNodeJob := scripts.VRFV2DeployUniverse(
		e,
		decimal.RequireFromString(constants.FallbackWeiPerUnitLinkString).BigInt(),
		decimal.RequireFromString(constants.SubscriptionBalanceString).BigInt(),
		&nodes[scripts.VRFPrimaryNodeName].VrfKeys[0],
		*linkAddress,
		*linkEthAddress,
		common.HexToAddress(*bhsContractAddressString),
		common.HexToAddress(*batchBHSAddressString),
		common.HexToAddress(*coordinatorAddressString),
		common.HexToAddress(*batchCoordinatorAddressString),
		&constants.MinConfs,
		&constants.MaxGasLimit,
		&constants.StalenessSeconds,
		&constants.GasAfterPayment,
		feeConfig,
		nodes,
	)

	for key, node := range nodes {
		client, app := connectToNode(&node.URL, output, file)

		//GET ALL JOBS
		jobIDs := getAllJobIDs(client, app, output)

		//DELETE ALL EXISTING JOBS
		for _, jobID := range jobIDs {
			deleteJob(jobID, client, app)
		}
		output.Reset()

		//CREATE JOBS
		switch key {
		case scripts.VRFPrimaryNodeName:
			createJob(vrfv2PrimaryNodeJob, client, app, output)
		case scripts.VRFBackupNodeName:
			createJob(vrfv2BackupNodeJob, client, app, output)
		case scripts.BHSNodeName:
			createJob(bhsNodeJob, client, app, output)
		case scripts.BHFNodeName:
			createJob(bhfNodeJob, client, app, output)
		}

	}
}

func createJob(jobSpec string, client *clcmd.Shell, app *cli.App, output *bytes.Buffer) {
	if err := os.WriteFile("job-spec.toml", []byte(jobSpec), 0666); err != nil {
		helpers.PanicErr(err)
	}
	job := presenters.JobResource{}
	flagSet := flag.NewFlagSet("blah", flag.ExitOnError)
	err := flagSet.Parse([]string{"./job-spec.toml"})
	helpers.PanicErr(err)
	err = client.CreateJob(cli.NewContext(app, flagSet, nil))
	helpers.PanicErr(err)
	helpers.PanicErr(json.Unmarshal(output.Bytes(), &job))
	output.Reset()
}

func deleteJob(jobID string, client *clcmd.Shell, app *cli.App) {
	flagSet := flag.NewFlagSet("blah", flag.ExitOnError)
	err := flagSet.Parse([]string{jobID})
	helpers.PanicErr(err)
	err = client.DeleteJob(cli.NewContext(app, flagSet, nil))
	helpers.PanicErr(err)
}

func getAllJobIDs(client *clcmd.Shell, app *cli.App, output *bytes.Buffer) []string {
	flagSet := flag.NewFlagSet("blah", flag.ExitOnError)
	err := client.ListJobs(cli.NewContext(app, flagSet, nil))
	helpers.PanicErr(err)
	jobs := clcmd.JobPresenters{}
	helpers.PanicErr(json.Unmarshal(output.Bytes(), &jobs))
	var jobIDs []string
	for _, job := range jobs {
		jobIDs = append(jobIDs, job.ID)
	}
	output.Reset()
	return jobIDs
}

func printETHKeyData(ethKeys []presenters.ETHKeyResource) {
	fmt.Println("------------- NODE INFORMATION -------------")
	for _, ethKey := range ethKeys {
		fmt.Println("-----------ETH Key-----------")
		fmt.Println("Address: ", ethKey.Address)
		fmt.Println("MaxGasPriceWei: ", ethKey.MaxGasPriceWei)
		fmt.Println("EthBalance: ", ethKey.EthBalance)
		fmt.Println("NextNonce: ", ethKey.NextNonce)
		fmt.Println("-----------------------------")
	}
}

func mapEthKeysToStringArr(ethKeys []presenters.ETHKeyResource) []string {
	var ethKeysString []string
	for _, ethKey := range ethKeys {
		ethKeysString = append(ethKeysString, ethKey.Address)
	}
	return ethKeysString
}

func mapVrfKeysToStringArr(vrfKeys []presenters.VRFKeyResource) []string {
	var vrfKeysString []string
	for _, vrfKey := range vrfKeys {
		vrfKeysString = append(vrfKeysString, vrfKey.Uncompressed)
	}
	return vrfKeysString
}

func printVRFKeyData(vrfKeys []presenters.VRFKeyResource) {
	fmt.Println("Number of VRF Keys on the node: ", len(vrfKeys))

	fmt.Println("------------- NODE INFORMATION -------------")
	for _, vrfKey := range vrfKeys {
		fmt.Println("-----------VRF Key-----------")
		fmt.Println("Compressed: ", vrfKey.Compressed)
		fmt.Println("Uncompressed: ", vrfKey.Uncompressed)
		fmt.Println("Hash: ", vrfKey.Hash)
		fmt.Println("-----------------------------")
	}
}

func connectToNode(nodeURL *string, output *bytes.Buffer, credFile string) (*clcmd.Shell, *cli.App) {
	client, app := newApp(*nodeURL, output)
	// login first to establish the session
	fmt.Println("logging in to:", *nodeURL)
	loginFs := flag.NewFlagSet("test", flag.ContinueOnError)
	loginFs.String("file", credFile, "")
	loginFs.Bool("bypass-version-check", true, "")
	loginCtx := cli.NewContext(app, loginFs, nil)
	err := client.RemoteLogin(loginCtx)
	helpers.PanicErr(err)
	output.Reset()
	fmt.Println()
	return client, app
}

func createVRFKeyIfNeeded(client *clcmd.Shell, app *cli.App, output *bytes.Buffer, numVRFKeys *int, nodeURL *string) []presenters.VRFKeyResource {
	var allVRFKeys []presenters.VRFKeyResource
	var vrfKeys []presenters.VRFKeyResource
	var newKeys []presenters.VRFKeyResource

	err := client.ListVRFKeys(&cli.Context{
		App: app,
	})
	helpers.PanicErr(err)

	helpers.PanicErr(json.Unmarshal(output.Bytes(), &vrfKeys))
	switch {
	case len(vrfKeys) == *numVRFKeys:
		fmt.Println(checkMarkEmoji, "found", len(vrfKeys), "vrf keys on", nodeURL)
	case len(vrfKeys) > *numVRFKeys:
		fmt.Println(checkMarkEmoji, "found", len(vrfKeys), "vrf keys on", nodeURL, " which is more than expected")
		os.Exit(1)
	default:
		fmt.Println(xEmoji, "found only", len(vrfKeys), "vrf keys on", nodeURL, ", creating",
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
		allVRFKeys = append(allVRFKeys, vrfKey)
	}
	for _, nk := range newKeys {
		allVRFKeys = append(allVRFKeys, nk)
	}
	return allVRFKeys
}

func createETHKeysIfNeeded(client *clcmd.Shell, app *cli.App, output *bytes.Buffer, numEthKeys *int, nodeURL *string, maxGasPriceGwei *int) []presenters.ETHKeyResource {
	var allETHKeysNode []presenters.ETHKeyResource
	var ethKeys []presenters.ETHKeyResource
	var newKeys []presenters.ETHKeyResource
	// check for ETH keys
	err := client.ListETHKeys(&cli.Context{
		App: app,
	})
	helpers.PanicErr(err)

	helpers.PanicErr(json.Unmarshal(output.Bytes(), &ethKeys))
	switch {
	case len(ethKeys) >= *numEthKeys:
		fmt.Println(checkMarkEmoji, "found", len(ethKeys), "eth keys on", nodeURL)
	case len(ethKeys) < *numEthKeys:
		fmt.Println(xEmoji, "found only", len(ethKeys), "eth keys on", nodeURL,
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
		allETHKeysNode = append(allETHKeysNode, ethKey)
	}
	for _, nk := range newKeys {
		allETHKeysNode = append(allETHKeysNode, nk)
	}

	return allETHKeysNode
}
