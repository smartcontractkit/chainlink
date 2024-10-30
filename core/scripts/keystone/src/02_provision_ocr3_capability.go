package src

import (
	"bytes"
	"text/template"

	"context"
	"flag"
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"

	ksdeploy "github.com/smartcontractkit/chainlink/deployment/keystone"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/ocr3_capability"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func provisionOCR3(
	env helpers.Environment,
	nodeSet NodeSet,
	chainID int64,
	p2pPort int64,
	ocrConfigFile string,
	artefactsDir string,
) (onchainMeta *onchainMeta, cacheHit bool) {
	onchainMeta, cacheHit = deployOCR3Contract(
		nodeSet,
		env,
		ocrConfigFile,
		artefactsDir,
	)

	deployOCR3JobSpecsTo(
		nodeSet,
		chainID,
		p2pPort,
		artefactsDir,
		onchainMeta,
	)

	return
}

func deployOCR3Contract(
	nodeSet NodeSet,
	env helpers.Environment,
	configFile string,
	artefacts string,
) (o *onchainMeta, cacheHit bool) {
	o = LoadOnchainMeta(artefacts, env)
	ocrConf := generateOCR3Config(
		nodeSet,
		configFile,
		env.ChainID,
	)

	if o.OCR3 != nil {
		// types.ConfigDigestPrefixKeystoneOCR3Capability
		fmt.Println("OCR3 Contract already deployed, checking config...")
		latestConfigDigestBytes, err := o.OCR3.LatestConfigDetails(nil)
		PanicErr(err)
		latestConfigDigest, err := types.BytesToConfigDigest(latestConfigDigestBytes.ConfigDigest[:])

		cc := ocrConfToContractConfig(ocrConf, latestConfigDigestBytes.ConfigCount)
		digester := evm.OCR3CapabilityOffchainConfigDigester{
			ChainID:         uint64(env.ChainID),
			ContractAddress: o.OCR3.Address(),
		}
		digest, err := digester.ConfigDigest(context.Background(), cc)
		PanicErr(err)

		if digest.Hex() == latestConfigDigest.Hex() {
			fmt.Printf("OCR3 Contract already deployed with the same config (digest: %s), skipping...\n", digest.Hex())
			return o, false
		}

		fmt.Printf("OCR3 Contract contains a different config, updating...\nOld digest: %s\nNew digest: %s\n", latestConfigDigest.Hex(), digest.Hex())
		setOCRConfig(o, env, ocrConf, artefacts)
		return o, true
	}

	fmt.Println("Deploying keystone ocr3 contract...")
	_, tx, ocrContract, err := ocr3_capability.DeployOCR3Capability(env.Owner, env.Ec)
	PanicErr(err)
	helpers.ConfirmContractDeployed(context.Background(), env.Ec, tx, env.ChainID)
	o.OCR3 = ocrContract
	setOCRConfig(o, env, ocrConf, artefacts)

	return o, true
}

func generateOCR3Config(nodeSet NodeSet, configFile string, chainID int64) ksdeploy.Orc2drOracleConfig {
	topLevelCfg := mustReadOCR3Config(configFile)
	cfg := topLevelCfg.OracleConfig
	cfg.OCRSecrets = deployment.XXXGenerateTestOCRSecrets()
	c, err := ksdeploy.GenerateOCR3Config(cfg, nodeKeysToKsDeployNodeKeys(nodeSet.NodeKeys[1:])) // skip the bootstrap node
	helpers.PanicErr(err)
	return c
}

func setOCRConfig(o *onchainMeta, env helpers.Environment, ocrConf ksdeploy.Orc2drOracleConfig, artefacts string) {
	tx, err := o.OCR3.SetConfig(env.Owner,
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
	artefactsDir string,
	onchainMeta *onchainMeta,
) {
	ocrAddress := onchainMeta.OCR3.Address().Hex()
	nodeKeys := nodeSet.NodeKeys
	nodes := nodeSet.Nodes

	var specName string
	for i, n := range nodes {
		var spec string

		if i == 0 {
			bootstrapSpecConfig := BootstrapJobSpecConfig{
				JobSpecName:              "ocr3_bootstrap",
				OCRConfigContractAddress: ocrAddress,
				ChainID:                  chainID,
			}
			specName = bootstrapSpecConfig.JobSpecName
			spec = createBootstrapJobSpec(bootstrapSpecConfig)
		} else {
			oc := OracleJobSpecConfig{
				JobSpecName:              fmt.Sprintf("ocr3_oracle"),
				OCRConfigContractAddress: ocrAddress,
				OCRKeyBundleID:           nodeKeys[i].OCR2BundleID,
				BootstrapURI:             fmt.Sprintf("%s@%s:%d", nodeKeys[0].P2PPeerID, nodeSet.Nodes[0].ServiceName, p2pPort),
				TransmitterID:            nodeKeys[i].EthAddress,
				ChainID:                  chainID,
				AptosKeyBundleID:         nodeKeys[i].AptosBundleID,
			}
			specName = oc.JobSpecName
			spec = createOracleJobSpec(oc)
		}

		api := newNodeAPI(n)
		upsertJob(api, specName, spec)

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

// BootstrapJobSpecConfig holds configuration for the bootstrap job spec
type BootstrapJobSpecConfig struct {
	JobSpecName              string
	OCRConfigContractAddress string
	ChainID                  int64
}

// OracleJobSpecConfig holds configuration for the oracle job spec
type OracleJobSpecConfig struct {
	JobSpecName              string
	OCRConfigContractAddress string
	OCRKeyBundleID           string
	BootstrapURI             string
	TransmitterID            string
	ChainID                  int64
	AptosKeyBundleID         string
}

func createBootstrapJobSpec(config BootstrapJobSpecConfig) string {
	const bootstrapTemplate = `
type = "bootstrap"
schemaVersion = 1
name = "{{ .JobSpecName }}"
contractID = "{{ .OCRConfigContractAddress }}"
relay = "evm"

[relayConfig]
chainID = "{{ .ChainID }}"
providerType = "ocr3-capability"
`

	tmpl, err := template.New("bootstrap").Parse(bootstrapTemplate)
	if err != nil {
		panic(err)
	}

	var rendered bytes.Buffer
	err = tmpl.Execute(&rendered, config)
	if err != nil {
		panic(err)
	}

	return rendered.String()
}

func createOracleJobSpec(config OracleJobSpecConfig) string {
	const oracleTemplate = `
type = "offchainreporting2"
schemaVersion = 1
name = "{{ .JobSpecName }}"
contractID = "{{ .OCRConfigContractAddress }}"
ocrKeyBundleID = "{{ .OCRKeyBundleID }}"
p2pv2Bootstrappers = [
  "{{ .BootstrapURI }}",
]
relay = "evm"
pluginType = "plugin"
transmitterID = "{{ .TransmitterID }}"

[relayConfig]
chainID = "{{ .ChainID }}"

[pluginConfig]
command = "chainlink-ocr3-capability"
ocrVersion = 3
pluginName = "ocr-capability"
providerType = "ocr3-capability"
telemetryType = "plugin"

[onchainSigningStrategy]
strategyName = 'multi-chain'
[onchainSigningStrategy.config]
evm = "{{ .OCRKeyBundleID }}"
aptos = "{{ .AptosKeyBundleID }}"
`

	tmpl, err := template.New("oracle").Parse(oracleTemplate)
	if err != nil {
		panic(err)
	}

	var rendered bytes.Buffer
	err = tmpl.Execute(&rendered, config)
	if err != nil {
		panic(err)
	}

	return rendered.String()
}
