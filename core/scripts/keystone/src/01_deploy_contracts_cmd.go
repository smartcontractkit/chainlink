package src

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/ocr3_capability"
)

type deployedContracts struct {
	OCRContract       common.Address `json:"ocrContract"`
	ForwarderContract common.Address `json:"forwarderContract"`
	// The block number of the transaction that set the config on the OCR3 contract. We use this to replay blocks from this point on
	// when we load the OCR3 job specs on the nodes.
	SetConfigTxBlock uint64 `json:"setConfigTxBlock"`
}

type deployContracts struct{}

func NewDeployContractsCommand() *deployContracts {
	return &deployContracts{}
}

func (g *deployContracts) Name() string {
	return "deploy-contracts"
}

// Run expects the follow environment variables to be set:
//
//  1. Deploys the OCR3 contract
//  2. Deploys the Forwarder contract
//  3. Sets the config on the OCR3 contract
//  4. Writes the deployed contract addresses to a file
//  5. Funds the transmitters
func (g *deployContracts) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ExitOnError)
	ocrConfigFile := fs.String("ocrfile", "config_example.json", "path to OCR config file")
	// create flags for all of the env vars then set the env vars to normalize the interface
	// this is a bit of a hack but it's the easiest way to make this work
	ethUrl := fs.String("ethurl", "", "URL of the Ethereum node")
	chainID := fs.Int64("chainid", 11155111, "chain ID of the Ethereum network to deploy to")
	accountKey := fs.String("accountkey", "", "private key of the account to deploy from")
	skipFunding := fs.Bool("skipfunding", false, "skip funding the transmitters")
	onlySetConfig := fs.Bool("onlysetconfig", false, "set the config on the OCR3 contract without deploying the contracts or funding transmitters")
	dryRun := fs.Bool("dryrun", false, "dry run, don't actually deploy the contracts and do not fund transmitters")

	err := fs.Parse(args)
	if err != nil ||
		*ocrConfigFile == "" || ocrConfigFile == nil ||
		*ethUrl == "" || ethUrl == nil ||
		*chainID == 0 || chainID == nil ||
		*accountKey == "" || accountKey == nil {
		fs.Usage()
		os.Exit(1)
	}

	os.Setenv("ETH_URL", *ethUrl)
	os.Setenv("ETH_CHAIN_ID", fmt.Sprintf("%d", *chainID))
	os.Setenv("ACCOUNT_KEY", *accountKey)

	deploy(*ocrConfigFile, *skipFunding, *dryRun, *onlySetConfig)
}

// deploy does the following:
//  1. Deploys the OCR3 contract
//  2. Deploys the Forwarder contract
//  3. Sets the config on the OCR3 contract
//  4. Writes the deployed contract addresses to a file
//  5. Funds the transmitters
func deploy(
	configFile string,
	skipFunding bool,
	dryRun bool,
	onlySetConfig bool,
) {
	env := helpers.SetupEnv(false)
	ocrConfig := generateOCR3Config(
		configFile,
		env.ChainID,
		".cache/PublicKeys.json",
	)

	if dryRun {
		fmt.Println("Dry run, skipping deployment and funding")
		return
	}

	if onlySetConfig {
		fmt.Println("Skipping deployment of contracts and skipping funding transmitters, only setting config")
		setOCR3Config(env, ocrConfig)
		return
	}

	if ContractsAlreadyDeployed() {
		fmt.Println("Contracts already deployed")
		return
	}

	fmt.Println("Deploying keystone ocr3 contract...")
	ocrContract := DeployKeystoneOCR3Capability(env)
	fmt.Println("Deploying keystone forwarder contract...")
	forwarderContract := DeployForwarder(env)

	fmt.Println("Writing deployed contract addresses to file...")
	contracts := deployedContracts{
		OCRContract:       ocrContract.Address(),
		ForwarderContract: forwarderContract.Address(),
	}
	jsonBytes, err := json.Marshal(contracts)
	PanicErr(err)

	err = os.WriteFile(DeployedContractsFilePath(), jsonBytes, 0600)
	PanicErr(err)

	setOCR3Config(env, ocrConfig)

	if skipFunding {
		fmt.Println("Skipping funding transmitters")
		return
	}
	fmt.Println("Funding transmitters...")
	transmittersStr := []string{}
	for _, t := range ocrConfig.Transmitters {
		transmittersStr = append(transmittersStr, t.String())
	}

	helpers.FundNodes(env, transmittersStr, big.NewInt(50000000000000000)) // 0.05 ETH
}

func setOCR3Config(
	env helpers.Environment,
	ocrConfig orc2drOracleConfig,
) {
	loadedContracts, err := LoadDeployedContracts()
	PanicErr(err)

	ocrContract, err := ocr3_capability.NewOCR3Capability(loadedContracts.OCRContract, env.Ec)
	PanicErr(err)
	fmt.Println("Setting OCR3 contract config...")
	tx, err := ocrContract.SetConfig(env.Owner,
		ocrConfig.Signers,
		ocrConfig.Transmitters,
		ocrConfig.F,
		ocrConfig.OnchainConfig,
		ocrConfig.OffchainConfigVersion,
		ocrConfig.OffchainConfig,
	)
	PanicErr(err)
	receipt := helpers.ConfirmTXMined(context.Background(), env.Ec, tx, env.ChainID)

	// Write blocknumber of the transaction to the deployed contracts file
	loadedContracts.SetConfigTxBlock = receipt.BlockNumber.Uint64()
	jsonBytes, err := json.Marshal(loadedContracts)
	PanicErr(err)
	err = os.WriteFile(DeployedContractsFilePath(), jsonBytes, 0600)
	PanicErr(err)
}

func LoadDeployedContracts() (deployedContracts, error) {
	if !ContractsAlreadyDeployed() {
		return deployedContracts{}, fmt.Errorf("no deployed contracts found, run deploy first")
	}

	jsonBytes, err := os.ReadFile(DeployedContractsFilePath())
	if err != nil {
		return deployedContracts{}, err
	}

	var contracts deployedContracts
	err = json.Unmarshal(jsonBytes, &contracts)
	return contracts, err
}

func ContractsAlreadyDeployed() bool {
	_, err := os.Stat(DeployedContractsFilePath())
	return err == nil
}

func DeployedContractsFilePath() string {
	return filepath.Join(artefactsDir, deployedContractsJSON)
}

func DeployForwarder(e helpers.Environment) *forwarder.KeystoneForwarder {
	_, tx, contract, err := forwarder.DeployKeystoneForwarder(e.Owner, e.Ec)
	PanicErr(err)
	helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)

	return contract
}

func DeployKeystoneOCR3Capability(e helpers.Environment) *ocr3_capability.OCR3Capability {
	_, tx, contract, err := ocr3_capability.DeployOCR3Capability(e.Owner, e.Ec)
	PanicErr(err)
	helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)

	return contract
}
