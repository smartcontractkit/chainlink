package solana

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
)

// Configuration
const (
	repoURL        = "https://github.com/smartcontractkit/chainlink-solana.git"
	cloneDir       = "./temp-repo"
	anchorDir      = "contracts"               // Path to the Anchor project within the repo
	deployDir      = "contracts/target/deploy" // Path to generated files
	buildCmd       = "build_contracts"
	replaceKeysCmd = "docker_update_contracts"
)

var cacheBuildParams = helpers.DomainParams{
	RepoURL:          repoURL,
	CloneDir:         cloneDir,
	AnchorDir:        anchorDir,
	DeployDir:        deployDir,
	ProgramFilesView: programToFileMap,
	BuildCmd:         buildCmd,
	ReplaceKeysCmd:   replaceKeysCmd,
}

var dataFeedsCache = cldf.ContractType("data_feeds_cache")
var programToFileMap = map[cldf.ContractType]string{
	dataFeedsCache: "programs/data_feeds_cache/src/lib.rs",
}
