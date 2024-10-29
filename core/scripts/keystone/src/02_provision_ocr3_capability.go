package src

import (
	"path/filepath"
	"strconv"
	"strings"

	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	ksdeploy "github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/ocr3_capability"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type provisionOCR3Capability struct{}

func NewProvisionOCR3CapabilityCommand() *provisionOCR3Capability {
	return &provisionOCR3Capability{}
}

func (g *provisionOCR3Capability) Name() string {
	return "provision-ocr3-capability"
}

func (g *provisionOCR3Capability) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ExitOnError)
	ocrConfigFile := fs.String("ocrfile", "ocr_config.json", "path to OCR config file")
	ethUrl := fs.String("ethurl", "", "URL of the Ethereum node")
	chainID := fs.Int64("chainid", 1337, "chain ID of the Ethereum network to deploy to")
	p2pPort := fs.Int64("p2pport", 6690, "p2p port")
	accountKey := fs.String("accountkey", "", "private key of the account to deploy from")
	nodeSetsPath := fs.String("nodesets", defaultNodeSetsPath, "Custom node sets location")
	nodeSetSize := fs.Int("nodesetsize", 5, "number of nodes in a nodeset")
	artefactsDir := fs.String("artefacts", defaultArtefactsDir, "Custom artefacts directory location")
	templatesLocation := fs.String("templates", defaultTemplatesDir, "Custom templates location")

	err := fs.Parse(args)

	if err != nil ||
		*accountKey == "" || accountKey == nil {
		fs.Usage()
		os.Exit(1)
	}

	// use flags for all of the env vars then set the env vars to normalize the interface
	// this is a bit of a hack but it's the easiest way to make this work
	os.Setenv("ETH_URL", *ethUrl)
	os.Setenv("ETH_CHAIN_ID", fmt.Sprintf("%d", *chainID))
	os.Setenv("ACCOUNT_KEY", *accountKey)
	os.Setenv("INSECURE_SKIP_VERIFY", "true")
	env := helpers.SetupEnv(false)

	nodeSet := downloadNodeSets(*chainID, *nodeSetsPath, *nodeSetSize).Workflow

	provisionOCR3(
		env,
		nodeSet,
		*chainID,
		*p2pPort,
		*ocrConfigFile,
		*templatesLocation,
		*artefactsDir,
	)
}

func provisionOCR3(
	env helpers.Environment,
	nodeSet NodeSet,
	chainID int64,
	p2pPort int64,
	ocrConfigFile string,
	templatesLocation string,
	artefactsDir string,
) (onchainMeta onchainMeta, shouldRedeployJobspecs bool) {
	onchainMeta, shouldRedeployJobspecs = deployOCR3Contract(
		nodeSet,
		env,
		ocrConfigFile,
		artefactsDir,
	)

	deployOCR3JobSpecsTo(
		nodeSet,
		chainID,
		p2pPort,
		templatesLocation,
		artefactsDir,
		onchainMeta,
		shouldRedeployJobspecs,
	)

	return
}

func deployOCR3Contract(
	nodeSet NodeSet,
	env helpers.Environment,
	configFile string,
	artefacts string,
) (o onchainMeta, shouldRedeployJobspecs bool) {
	o = LoadOnchainMeta(artefacts, env)
	ocrConf := generateOCR3Config(
		nodeSet,
		configFile,
		env.ChainID,
	)

	if o.OCRContract != nil {
		// types.ConfigDigestPrefixKeystoneOCR3Capability
		fmt.Println("OCR3 Contract already deployed, checking config...")
		latestConfigDigestBytes, err := o.OCRContract.LatestConfigDetails(nil)
		PanicErr(err)
		latestConfigDigest, err := types.BytesToConfigDigest(latestConfigDigestBytes.ConfigDigest[:])

		cc := ocrConfToContractConfig(ocrConf, latestConfigDigestBytes.ConfigCount)
		digester := evmutil.EVMOffchainConfigDigester{
			ChainID:         uint64(env.ChainID),
			ContractAddress: o.OCRContract.Address(),
		}
		digest, err := digester.ConfigDigest(context.Background(), cc)
		PanicErr(err)

		if digest.Hex() == latestConfigDigest.Hex() {
			fmt.Println("OCR3 Contract already deployed with the same config, skipping...")
			return o, false
		}

		fmt.Println("OCR3 Contract contains a different config, updating...")
		setOCRConfig(o, env, ocrConf, artefacts)
		return o, true
	}

	fmt.Println("Deploying keystone ocr3 contract...")
	_, tx, ocrContract, err := ocr3_capability.DeployOCR3Capability(env.Owner, env.Ec)
	PanicErr(err)
	helpers.ConfirmContractDeployed(context.Background(), env.Ec, tx, env.ChainID)
	o.OCRContract = ocrContract
	setOCRConfig(o, env, ocrConf, artefacts)

	return o, true
}

func generateOCR3Config(nodeSet NodeSet, configFile string, chainID int64) ksdeploy.Orc2drOracleConfig {
	topLevelCfg := mustReadOCR3Config(configFile)
	cfg := topLevelCfg.OracleConfig
	c, err := ksdeploy.GenerateOCR3Config(cfg, nodeKeysToKsDeployNodeKeys(nodeSet.NodeKeys[1:])) // skip the bootstrap node
	helpers.PanicErr(err)
	return c
}

func setOCRConfig(o onchainMeta, env helpers.Environment, ocrConf ksdeploy.Orc2drOracleConfig, artefacts string) {
	tx, err := o.OCRContract.SetConfig(env.Owner,
		ocrConf.Signers,
		ocrConf.Transmitters,
		ocrConf.F,
		ocrConf.OnchainConfig,
		ocrConf.OffchainConfigVersion,
		ocrConf.OffchainConfig,
	)
	PanicErr(err)
	receipt := helpers.ConfirmTXMined(context.Background(), env.Ec, tx, env.ChainID)
	o.SetConfigTxBlock = receipt.BlockNumber.Uint64()
	WriteOnchainMeta(o, artefacts)
}

func deployOCR3JobSpecsTo(
	nodeSet NodeSet,
	chainID int64,
	p2pPort int64,
	templatesDir string,
	artefactsDir string,
	onchainMeta onchainMeta,
	replaceJob bool,
) {

	jobspecs := generateOCR3JobSpecs(
		nodeSet,
		templatesDir,
		chainID,
		p2pPort,
		onchainMeta.OCRContract.Address().Hex(),
	)
	flattenedSpecs := []hostSpec{jobspecs.bootstrap}
	flattenedSpecs = append(flattenedSpecs, jobspecs.oracles...)

	if len(nodeSet.Nodes) != len(flattenedSpecs) {
		PanicErr(errors.New("Mismatched node and job spec lengths"))
	}

	for i, n := range nodeSet.Nodes {
		api := newNodeAPI(n)
		specToDeploy := flattenedSpecs[i].spec.ToString()
		maybeUpsertJob(api, specToDeploy, specToDeploy, replaceJob)

		fmt.Printf("Replaying from block: %d\n", onchainMeta.SetConfigTxBlock)
		fmt.Printf("EVM Chain ID: %d\n\n", chainID)
		api.withFlags(api.methods.ReplayFromBlock, func(fs *flag.FlagSet) {
			err := fs.Set("block-number", fmt.Sprint(onchainMeta.SetConfigTxBlock))
			helpers.PanicErr(err)
			err = fs.Set("evm-chain-id", fmt.Sprint(chainID))
			helpers.PanicErr(err)
		}).mustExec()
	}
}

type spec []string

func (s spec) ToString() string {
	return strings.Join(s, "\n")
}

type hostSpec struct {
	spec spec
	host string
}

type donHostSpec struct {
	bootstrap hostSpec
	oracles   []hostSpec
}

func generateOCR3JobSpecs(
	nodeSet NodeSet,
	templatesDir string,
	chainID int64,
	p2pPort int64,
	ocrConfigContractAddress string,
) donHostSpec {
	nodeKeys := nodeKeysToKsDeployNodeKeys(nodeSet.NodeKeys)
	nodes := nodeSet.Nodes
	bootstrapNode := nodeKeys[0]

	bootstrapSpecLines, err := readLines(filepath.Join(templatesDir, bootstrapSpecTemplate))
	helpers.PanicErr(err)
	bootHost := nodes[0].ServiceName
	bootstrapSpecLines = replaceOCR3TemplatePlaceholders(
		bootstrapSpecLines,
		chainID, p2pPort,
		ocrConfigContractAddress, bootHost,
		bootstrapNode, bootstrapNode,
	)
	bootstrap := hostSpec{bootstrapSpecLines, bootHost}

	oracleSpecLinesTemplate, err := readLines(filepath.Join(templatesDir, oracleSpecTemplate))
	helpers.PanicErr(err)
	oracles := []hostSpec{}
	for i := 1; i < len(nodes); i++ {
		oracleSpecLines := oracleSpecLinesTemplate
		oracleSpecLines = replaceOCR3TemplatePlaceholders(
			oracleSpecLines,
			chainID, p2pPort,
			ocrConfigContractAddress, bootHost,
			bootstrapNode, nodeKeys[i],
		)
		oracles = append(oracles, hostSpec{oracleSpecLines, nodes[i].RemoteURL.Host})
	}

	return donHostSpec{
		bootstrap: bootstrap,
		oracles:   oracles,
	}
}

func replaceOCR3TemplatePlaceholders(
	lines []string,

	chainID, p2pPort int64,
	contractAddress, bootHost string,
	boot, node ksdeploy.NodeKeys,
) (output []string) {
	chainIDStr := strconv.FormatInt(chainID, 10)
	bootstrapper := fmt.Sprintf("%s@%s:%d", boot.P2PPeerID, bootHost, p2pPort)
	for _, l := range lines {
		l = strings.Replace(l, "{{ chain_id }}", chainIDStr, 1)
		l = strings.Replace(l, "{{ ocr_config_contract_address }}", contractAddress, 1)
		l = strings.Replace(l, "{{ transmitter_id }}", node.EthAddress, 1)
		l = strings.Replace(l, "{{ ocr_key_bundle_id }}", node.OCR2BundleID, 1)
		l = strings.Replace(l, "{{ aptos_key_bundle_id }}", node.AptosBundleID, 1)
		l = strings.Replace(l, "{{ bootstrapper_p2p_id }}", bootstrapper, 1)
		output = append(output, l)
	}
	return
}

func mustReadOCR3Config(fileName string) (output ksdeploy.TopLevelConfigSource) {
	return mustReadJSON[ksdeploy.TopLevelConfigSource](fileName)
}

func nodeKeysToKsDeployNodeKeys(nks []NodeKeys) []ksdeploy.NodeKeys {
	keys := []ksdeploy.NodeKeys{}
	for _, nk := range nks {
		keys = append(keys, ksdeploy.NodeKeys{
			EthAddress:            nk.EthAddress,
			AptosAccount:          nk.AptosAccount,
			AptosBundleID:         nk.AptosBundleID,
			AptosOnchainPublicKey: nk.AptosOnchainPublicKey,
			P2PPeerID:             nk.P2PPeerID,
			OCR2BundleID:          nk.OCR2BundleID,
			OCR2OnchainPublicKey:  nk.OCR2OnchainPublicKey,
			OCR2OffchainPublicKey: nk.OCR2OffchainPublicKey,
			OCR2ConfigPublicKey:   nk.OCR2ConfigPublicKey,
			CSAPublicKey:          nk.CSAPublicKey,
		})
	}
	return keys
}
