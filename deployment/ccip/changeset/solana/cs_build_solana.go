package solana

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
)

// Configuration
const (
	repoURL   = "https://github.com/smartcontractkit/chainlink-ccip.git"
	cloneDir  = "./temp-repo"
	anchorDir = "chains/solana/contracts" // Path to the Anchor project within the repo
	deployDir = "chains/solana/contracts/target/deploy"
)

var ccipBuildParams = helpers.DomainParams{
	RepoURL:          repoURL,
	CloneDir:         cloneDir,
	AnchorDir:        anchorDir,
	DeployDir:        deployDir,
	Syncers:          []func() error{syncRouterAndCommon},
	ProgramFilesView: programToFileMap,
}

// Map program names to their Rust file paths (relative to the Anchor project root)
// Needed for upgrades in place
var programToFileMap = map[cldf.ContractType]string{
	shared.Router:                  "programs/ccip-router/src/lib.rs",
	shared.CCIPCommon:              "programs/ccip-common/src/lib.rs",
	shared.FeeQuoter:               "programs/fee-quoter/src/lib.rs",
	shared.OffRamp:                 "programs/ccip-offramp/src/lib.rs",
	shared.BurnMintTokenPool:       "programs/burnmint-token-pool/src/lib.rs",
	shared.LockReleaseTokenPool:    "programs/lockrelease-token-pool/src/lib.rs",
	shared.RMNRemote:               "programs/rmn-remote/src/lib.rs",
	types.AccessControllerProgram:  "programs/access-controller/src/lib.rs",
	types.ManyChainMultisigProgram: "programs/mcm/src/lib.rs",
	types.RBACTimelockProgram:      "programs/timelock/src/lib.rs",
}

func syncRouterAndCommon() error {
	routerFileName := programToFileMap[shared.Router]
	commonFileName := programToFileMap[shared.CCIPCommon]
	routerFile := filepath.Join(cloneDir, anchorDir, routerFileName)
	commonFile := filepath.Join(cloneDir, anchorDir, commonFileName)
	file, err := os.Open(routerFile)
	if err != nil {
		return fmt.Errorf("error opening router file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	declareRegex := regexp.MustCompile(`declare_id!\(\"(.+?)\"\);`)
	var declareID string

	for scanner.Scan() {
		match := declareRegex.FindStringSubmatch(scanner.Text())
		if match != nil {
			declareID = match[0]
			break
		}
	}

	if declareID == "" {
		return errors.New("declare_id not found in router file")
	}

	commonContent, err := os.ReadFile(commonFile)
	if err != nil {
		return fmt.Errorf("error reading common file: %w", err)
	}

	updatedContent := declareRegex.ReplaceAllString(string(commonContent), declareID)

	return os.WriteFile(commonFile, []byte(updatedContent), 0600)
}
