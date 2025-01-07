package src

import (
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/deployment"
	ksdeploy "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

func mustReadConfig(fileName string) (output ksdeploy.TopLevelConfigSource) {
	return mustParseJSON[ksdeploy.TopLevelConfigSource](fileName)
}

func generateOCR3Config(nodeList string, configFile string, chainID int64, pubKeysPath string) ksdeploy.OCR3OnchainConfig {
	topLevelCfg := mustReadConfig(configFile)
	cfg := topLevelCfg.OracleConfig
	nca := downloadNodePubKeys(nodeList, chainID, pubKeysPath)
	c, err := ksdeploy.GenerateOCR3Config(cfg, nca, deployment.XXXGenerateTestOCRSecrets())
	helpers.PanicErr(err)
	return c
}
