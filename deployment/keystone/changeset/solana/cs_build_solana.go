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

var keystoneBuildParams = helpers.DomainParams{
	RepoURL:          repoURL,
	CloneDir:         cloneDir,
	AnchorDir:        anchorDir,
	DeployDir:        deployDir,
	ProgramFilesView: programToFileMap,
	BuildCmd:         buildCmd,
	ReplaceKeysCmd:   replaceKeysCmd,
}

// Map program names to their Rust file paths (relative to the Anchor project root)
// Needed for upgrades in place
var keystoneForwarder = cldf.ContractType("keystone_forwarder")
var programToFileMap = map[cldf.ContractType]string{
	keystoneForwarder: "programs/keystone-forwarder/src/lib.rs",
}
