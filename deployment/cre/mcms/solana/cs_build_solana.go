package solana

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
)

const (
	mcmsRepoURL        = "https://github.com/smartcontractkit/chainlink-ccip.git"
	mcmsCloneDir       = "./temp-repo-cre-mcms"
	mcmsAnchorDir      = "chains/solana/contracts"
	mcmsDeployDir      = "chains/solana/contracts/target/deploy"
	mcmsBuildCmd       = "docker-build-contracts"
	mcmsReplaceKeysCmd = "docker-update-contracts"
)

// mcmsBuildParams configures cloning and building MCMS Solana programs from chainlink-ccip.
var mcmsBuildParams = helpers.DomainParams{
	RepoURL:        mcmsRepoURL,
	CloneDir:       mcmsCloneDir,
	AnchorDir:      mcmsAnchorDir,
	DeployDir:      mcmsDeployDir,
	BuildCmd:       mcmsBuildCmd,
	ReplaceKeysCmd: mcmsReplaceKeysCmd,
	ProgramFilesView: map[cldf.ContractType]string{
		types.AccessControllerProgram:  "programs/access-controller/src/lib.rs",
		types.ManyChainMultisigProgram: "programs/mcm/src/lib.rs",
		types.RBACTimelockProgram:      "programs/timelock/src/lib.rs",
	},
}

// BuildMCMSPrograms clones/builds or downloads MCMS program artifacts into destinationDir.
func BuildMCMSPrograms(env cldf.Environment, config helpers.BuildSolanaConfig) error {
	return helpers.BuildSolana(env, config, mcmsBuildParams)
}
