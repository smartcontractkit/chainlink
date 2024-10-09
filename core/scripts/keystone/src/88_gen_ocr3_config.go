package src

import (
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	ksdeploy "github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
)

func mustReadConfig(fileName string) (output ksdeploy.TopLevelConfigSource) {
	return mustParseJSON[ksdeploy.TopLevelConfigSource](fileName)
}

func generateOCR3Config(nodeList string, configFile string, chainID int64, pubKeysPath string) ksdeploy.Orc2drOracleConfig {
	topLevelCfg := mustReadConfig(configFile)
	cfg := topLevelCfg.OracleConfig
	nca := downloadNodePubKeys(nodeList, chainID, pubKeysPath)
	c, err := ksdeploy.GenerateOCR3Config(cfg, nodeKeysToKsDeployNodeKeys(nca))
	helpers.PanicErr(err)
	return c
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
