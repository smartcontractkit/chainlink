package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/scripts/common/vrf/constants"
	"github.com/smartcontractkit/chainlink/core/scripts/common/vrf/model"
	"github.com/smartcontractkit/chainlink/core/scripts/vrfv2/testnet/v2scripts"
	"github.com/smartcontractkit/chainlink/core/scripts/vrfv2plus/testnet/v2plusscripts"
	clcmd "github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
	"github.com/urfave/cli"
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

	vrfPrimaryNodeURL := flag.String("vrf-primary-node-url", "", "remote node URL")
	vrfBackupNodeURL := flag.String("vrf-backup-node-url", "", "remote node URL")
	bhsNodeURL := flag.String("bhs-node-url", "", "remote node URL")
	bhsBackupNodeURL := flag.String("bhs-backup-node-url", "", "remote node URL")
	bhfNodeURL := flag.String("bhf-node-url", "", "remote node URL")
	nodeSendingKeyFundingAmount := flag.String("sending-key-funding-amount", constants.NodeSendingKeyFundingAmount, "sending key funding amount")

	vrfPrimaryCredsFile := flag.String("vrf-primary-creds-file", "", "Creds to authenticate to the node")
	vrfBackupCredsFile := flag.String("vrf-bk-creds-file", "", "Creds to authenticate to the node")
	bhsCredsFile := flag.String("bhs-creds-file", "", "Creds to authenticate to the node")
	bhsBackupCredsFile := flag.String("bhs-bk-creds-file", "", "Creds to authenticate to the node")
	bhfCredsFile := flag.String("bhf-creds-file", "", "Creds to authenticate to the node")

	numEthKeys := flag.Int("num-eth-keys", 5, "Number of eth keys to create")
	maxGasPriceGwei := flag.Int("max-gas-price-gwei", -1, "Max gas price gwei of the eth keys")
	numVRFKeys := flag.Int("num-vrf-keys", 1, "Number of vrf keys to create")
	batchFulfillmentEnabled := flag.Bool("batch-fulfillment-enabled", constants.BatchFulfillmentEnabled, "whether to enable batch fulfillment on Cl node")

	vrfVersion := flag.String("vrf-version", "v2", "VRF version to use")
	deployContractsAndCreateJobs := flag.Bool("deploy-contracts-and-create-jobs", false, "whether to deploy contracts and create jobs")

	subscriptionBalanceJuelsString := flag.String("subscription-balance", constants.SubscriptionBalanceJuels, "amount to fund subscription with Link token (Juels)")
	subscriptionBalanceNativeWeiString := flag.String("subscription-balance-native", constants.SubscriptionBalanceNativeWei, "amount to fund subscription with native token (Wei)")

	minConfs := flag.Int("min-confs", constants.MinConfs, "minimum confirmations")
	linkAddress := flag.String("link-address", "", "address of link token")
	linkEthAddress := flag.String("link-eth-feed", "", "address of link eth feed")
	bhsContractAddressString := flag.String("bhs-address", "", "address of BHS contract")
	batchBHSAddressString := flag.String("batch-bhs-address", "", "address of Batch BHS contract")
	coordinatorAddressString := flag.String("coordinator-address", "", "address of VRF Coordinator contract")
	batchCoordinatorAddressString := flag.String("batch-coordinator-address", "", "address Batch VRF Coordinator contract")

	e := helpers.SetupEnv(false)
	flag.Parse()
	nodesMap := make(map[string]model.Node)

	if *vrfVersion != "v2" && *vrfVersion != "v2plus" {
		panic(fmt.Sprintf("Invalid VRF Version `%s`. Only `v2` and `v2plus` are supported", *vrfVersion))
	}
	fmt.Println("Using VRF Version:", *vrfVersion)

	fundingAmount := decimal.RequireFromString(*nodeSendingKeyFundingAmount).BigInt()
	subscriptionBalanceJuels := decimal.RequireFromString(*subscriptionBalanceJuelsString).BigInt()
	subscriptionBalanceNativeWei := decimal.RequireFromString(*subscriptionBalanceNativeWeiString).BigInt()

	if *vrfPrimaryNodeURL != "" {
		nodesMap[model.VRFPrimaryNodeName] = model.Node{
			URL:                     *vrfPrimaryNodeURL,
			SendingKeyFundingAmount: fundingAmount,
			CredsFile:               *vrfPrimaryCredsFile,
		}
	}
	if *vrfBackupNodeURL != "" {
		nodesMap[model.VRFBackupNodeName] = model.Node{
			URL:                     *vrfBackupNodeURL,
			SendingKeyFundingAmount: fundingAmount,
			CredsFile:               *vrfBackupCredsFile,
		}
	}
	if *bhsNodeURL != "" {
		nodesMap[model.BHSNodeName] = model.Node{
			URL:                     *bhsNodeURL,
			SendingKeyFundingAmount: fundingAmount,
			CredsFile:               *bhsCredsFile,
		}
	}
	if *bhsBackupNodeURL != "" {
		nodesMap[model.BHSBackupNodeName] = model.Node{
			URL:                     *bhsBackupNodeURL,
			SendingKeyFundingAmount: fundingAmount,
			CredsFile:               *bhsBackupCredsFile,
		}
	}

	if *bhfNodeURL != "" {
		nodesMap[model.BHFNodeName] = model.Node{
			URL:                     *bhfNodeURL,
			SendingKeyFundingAmount: fundingAmount,
			CredsFile:               *bhfCredsFile,
		}
	}

	output := &bytes.Buffer{}
	for key, node := range nodesMap {

		client, app := connectToNode(&node.URL, output, node.CredsFile)
		ethKeys := createETHKeysIfNeeded(client, app, output, numEthKeys, &node.URL, maxGasPriceGwei)
		if key == model.VRFPrimaryNodeName {
			vrfKeys := createVRFKeyIfNeeded(client, app, output, numVRFKeys, &node.URL)
			node.VrfKeys = mapVrfKeysToStringArr(vrfKeys)
			printVRFKeyData(vrfKeys)
			exportVRFKey(client, app, vrfKeys[0], output)
		}

		if key == model.VRFBackupNodeName {
			vrfKeys := getVRFKeys(client, app, output)
			node.VrfKeys = mapVrfKeysToStringArr(vrfKeys)
		}

		node.SendingKeys = mapEthKeysToSendingKeyArr(ethKeys)
		printETHKeyData(ethKeys)
		fundNodesIfNeeded(node, key, e)
		nodesMap[key] = node
	}
	importVRFKeyToNodeIfSet(vrfBackupNodeURL, nodesMap, output, nodesMap[model.VRFBackupNodeName].CredsFile)

	if *deployContractsAndCreateJobs {

		contractAddresses := model.ContractAddresses{
			LinkAddress:             *linkAddress,
			LinkEthAddress:          *linkEthAddress,
			BhsContractAddress:      common.HexToAddress(*bhsContractAddressString),
			BatchBHSAddress:         common.HexToAddress(*batchBHSAddressString),
			CoordinatorAddress:      common.HexToAddress(*coordinatorAddressString),
			BatchCoordinatorAddress: common.HexToAddress(*batchCoordinatorAddressString),
		}

		var jobSpecs model.JobSpecs

		switch *vrfVersion {
		case "v2":
			feeConfigV2 := vrf_coordinator_v2.VRFCoordinatorV2FeeConfig{
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

			coordinatorConfigV2 := v2scripts.CoordinatorConfigV2{
				MinConfs:               minConfs,
				MaxGasLimit:            &constants.MaxGasLimit,
				StalenessSeconds:       &constants.StalenessSeconds,
				GasAfterPayment:        &constants.GasAfterPayment,
				FallbackWeiPerUnitLink: constants.FallbackWeiPerUnitLink,
				FeeConfig:              feeConfigV2,
			}

			jobSpecs = v2scripts.VRFV2DeployUniverse(
				e,
				subscriptionBalanceJuels,
				&nodesMap[model.VRFPrimaryNodeName].VrfKeys[0],
				contractAddresses,
				coordinatorConfigV2,
				*batchFulfillmentEnabled,
				nodesMap,
			)
		case "v2plus":
			feeConfigV2Plus := vrf_coordinator_v2_5.VRFCoordinatorV25FeeConfig{
				FulfillmentFlatFeeLinkPPM:   uint32(constants.FlatFeeLinkPPM),
				FulfillmentFlatFeeNativePPM: uint32(constants.FlatFeeNativePPM),
			}
			coordinatorConfigV2Plus := v2plusscripts.CoordinatorConfigV2Plus{
				MinConfs:               minConfs,
				MaxGasLimit:            &constants.MaxGasLimit,
				StalenessSeconds:       &constants.StalenessSeconds,
				GasAfterPayment:        &constants.GasAfterPayment,
				FallbackWeiPerUnitLink: constants.FallbackWeiPerUnitLink,
				FeeConfig:              feeConfigV2Plus,
			}

			jobSpecs = v2plusscripts.VRFV2PlusDeployUniverse(
				e,
				subscriptionBalanceJuels,
				subscriptionBalanceNativeWei,
				&nodesMap[model.VRFPrimaryNodeName].VrfKeys[0],
				contractAddresses,
				coordinatorConfigV2Plus,
				*batchFulfillmentEnabled,
				nodesMap,
			)
		}

		for key, node := range nodesMap {
			client, app := connectToNode(&node.URL, output, node.CredsFile)

			//GET ALL JOBS
			jobIDs := getAllJobIDs(client, app, output)

			//DELETE ALL EXISTING JOBS
			for _, jobID := range jobIDs {
				deleteJob(jobID, client, app, output)
			}
			//CREATE JOBS

			switch key {
			case model.VRFPrimaryNodeName:
				createJob(jobSpecs.VRFPrimaryNode, client, app, output)
			case model.VRFBackupNodeName:
				createJob(jobSpecs.VRFBackupyNode, client, app, output)
			case model.BHSNodeName:
				createJob(jobSpecs.BHSNode, client, app, output)
			case model.BHSBackupNodeName:
				createJob(jobSpecs.BHSBackupNode, client, app, output)
			case model.BHFNodeName:
				createJob(jobSpecs.BHFNode, client, app, output)
			}
		}
	}
}

func fundNodesIfNeeded(node model.Node, key string, e helpers.Environment) {
	if node.SendingKeyFundingAmount.Cmp(big.NewInt(0)) == 1 {
		fmt.Println("\nFunding", key, "Node's Sending Keys. Need to fund each key with", node.SendingKeyFundingAmount, "wei")
		for _, sendingKey := range node.SendingKeys {
			fundingToSendWei := new(big.Int).Sub(node.SendingKeyFundingAmount, sendingKey.BalanceEth)
			if fundingToSendWei.Cmp(big.NewInt(0)) == 1 {
				helpers.FundNode(e, sendingKey.Address, fundingToSendWei)
			} else {
				fmt.Println("\nSkipping Funding", sendingKey.Address, "since it has", sendingKey.BalanceEth.String(), "wei")
			}
		}
	} else {
		fmt.Println("\nSkipping Funding", key, "Node's Sending Keys since funding amount is 0 wei")
	}
}

func importVRFKeyToNodeIfSet(vrfBackupNodeURL *string, nodes map[string]model.Node, output *bytes.Buffer, file string) {
	if *vrfBackupNodeURL != "" {
		vrfBackupNode := nodes[model.VRFBackupNodeName]
		vrfPrimaryNode := nodes[model.VRFBackupNodeName]

		if len(vrfBackupNode.VrfKeys) == 0 || vrfPrimaryNode.VrfKeys[0] != vrfBackupNode.VrfKeys[0] {
			client, app := connectToNode(&vrfBackupNode.URL, output, file)
			importVRFKey(client, app, output)

			vrfKeys := getVRFKeys(client, app, output)

			vrfBackupNode.VrfKeys = mapVrfKeysToStringArr(vrfKeys)
			if len(vrfBackupNode.VrfKeys) == 0 {
				panic("VRF Key was not imported to VRF Backup Node")
			}
			printVRFKeyData(vrfKeys)
		}
	}
}

func getVRFKeys(client *clcmd.Shell, app *cli.App, output *bytes.Buffer) []presenters.VRFKeyResource {
	var vrfKeys []presenters.VRFKeyResource

	err := client.ListVRFKeys(&cli.Context{
		App: app,
	})
	helpers.PanicErr(err)
	helpers.PanicErr(json.Unmarshal(output.Bytes(), &vrfKeys))
	output.Reset()
	return vrfKeys
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

func exportVRFKey(client *clcmd.Shell, app *cli.App, vrfKey presenters.VRFKeyResource, output *bytes.Buffer) {
	if err := os.WriteFile("vrf-key-password.txt", []byte("twochains"), 0666); err != nil {
		helpers.PanicErr(err)
	}
	flagSet := flag.NewFlagSet("blah", flag.ExitOnError)
	flagSet.String("new-password", "./vrf-key-password.txt", "")
	flagSet.String("output", "exportedvrf.json", "")
	err := flagSet.Parse([]string{vrfKey.Compressed})
	helpers.PanicErr(err)
	err = client.ExportVRFKey(cli.NewContext(app, flagSet, nil))
	helpers.PanicErr(err)
	output.Reset()
}

func importVRFKey(client *clcmd.Shell, app *cli.App, output *bytes.Buffer) {
	if err := os.WriteFile("vrf-key-password.txt", []byte("twochains"), 0666); err != nil {
		helpers.PanicErr(err)
	}
	flagSet := flag.NewFlagSet("blah", flag.ExitOnError)
	flagSet.String("old-password", "./vrf-key-password.txt", "")
	err := flagSet.Parse([]string{"exportedvrf.json"})
	helpers.PanicErr(err)
	err = client.ImportVRFKey(cli.NewContext(app, flagSet, nil))
	helpers.PanicErr(err)
	output.Reset()
}

func deleteJob(jobID string, client *clcmd.Shell, app *cli.App, output *bytes.Buffer) {
	flagSet := flag.NewFlagSet("blah", flag.ExitOnError)
	err := flagSet.Parse([]string{jobID})
	helpers.PanicErr(err)
	err = client.DeleteJob(cli.NewContext(app, flagSet, nil))
	helpers.PanicErr(err)
	output.Reset()
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
		fmt.Println("-----------------------------")
	}
}

func mapEthKeysToSendingKeyArr(ethKeys []presenters.ETHKeyResource) []model.SendingKey {
	var sendingKeys []model.SendingKey
	for _, ethKey := range ethKeys {
		sendingKey := model.SendingKey{Address: ethKey.Address, BalanceEth: ethKey.EthBalance.ToInt()}
		sendingKeys = append(sendingKeys, sendingKey)
	}
	return sendingKeys
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
	var newKeys []presenters.VRFKeyResource

	vrfKeys := getVRFKeys(client, app, output)

	switch {
	case len(vrfKeys) == *numVRFKeys:
		fmt.Println(checkMarkEmoji, "found", len(vrfKeys), "vrf keys on", *nodeURL)
	case len(vrfKeys) > *numVRFKeys:
		fmt.Println(xEmoji, "found", len(vrfKeys), "vrf keys on", nodeURL, " which is more than expected")
		os.Exit(1)
	default:
		fmt.Println(xEmoji, "found only", len(vrfKeys), "vrf keys on", nodeURL, ", creating",
			*numVRFKeys-len(vrfKeys), "more")
		toCreate := *numVRFKeys - len(vrfKeys)
		for i := 0; i < toCreate; i++ {
			output.Reset()
			newKey := createVRFKey(client, app, output)
			newKeys = append(newKeys, newKey)
		}
		fmt.Println("NEW VRF KEYS:", strings.Join(func() (r []string) {
			for _, k := range newKeys {
				r = append(r, k.Uncompressed)
			}
			return
		}(), ", "))
	}
	fmt.Println()
	for _, vrfKey := range vrfKeys {
		allVRFKeys = append(allVRFKeys, vrfKey)
	}
	for _, nk := range newKeys {
		allVRFKeys = append(allVRFKeys, nk)
	}
	return allVRFKeys
}

func createVRFKey(client *clcmd.Shell, app *cli.App, output *bytes.Buffer) presenters.VRFKeyResource {
	var newKey presenters.VRFKeyResource
	err := client.CreateVRFKey(
		cli.NewContext(app, flag.NewFlagSet("blah", flag.ExitOnError), nil))
	helpers.PanicErr(err)
	helpers.PanicErr(json.Unmarshal(output.Bytes(), &newKey))
	output.Reset()
	return newKey
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
		fmt.Println(checkMarkEmoji, "found", len(ethKeys), "eth keys on", *nodeURL)
	case len(ethKeys) < *numEthKeys:
		fmt.Println(xEmoji, "found only", len(ethKeys), "eth keys on", *nodeURL,
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
