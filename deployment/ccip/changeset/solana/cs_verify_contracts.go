package solana

import (
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	csState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
)

// https://solana.com/developers/guides/advanced/verified-builds
type VerifyBuildConfig struct {
	GitCommitSha               string
	ChainSelector              uint64
	VerifyFeeQuoter            bool
	VerifyRouter               bool
	VerifyOffRamp              bool
	VerifyRMNRemote            bool
	VerifyBurnMintTokenPool    bool
	VerifyLockReleaseTokenPool bool
	VerifyAccessController     bool
	VerifyMCM                  bool
	VerifyTimelock             bool
	RemoteVerification         bool
	MCMSSolana                 *MCMSConfigSolana
}

func runSolanaVerify(networkURL, programID, libraryName, commitHash, mountPath string, remote bool) error {
	cmdArgs := []string{
		"verify-from-repo",
		"-u", networkURL,
		"--program-id", programID,
		"--library-name", libraryName,
		strings.TrimSuffix(repoURL, ".git"),
		"--commit-hash", commitHash,
		"--mount-path", mountPath,
	}

	// Add --remote flag if remote verification is enabled
	if remote {
		cmdArgs = append(cmdArgs, "--remote")
	}

	output, err := runCommand("solana-verify", cmdArgs, ".")
	fmt.Println(output)
	if err != nil {
		return fmt.Errorf("solana program verification failed: %s %w", output, err)
	}
	return nil
}

func VerifyBuild(e deployment.Environment, cfg VerifyBuildConfig) (deployment.ChangesetOutput, error) {
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]

	addresses, err := e.ExistingAddresses.AddressesForChain(cfg.ChainSelector)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get existing addresses: %w", err)
	}
	mcmState, err := csState.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	verifications := []struct {
		name       string
		programID  string
		programLib string
		enabled    bool
	}{
		{"Fee Quoter", chainState.FeeQuoter.String(), deployment.FeeQuoterProgramName, cfg.VerifyFeeQuoter},
		{"Router", chainState.Router.String(), deployment.RouterProgramName, cfg.VerifyRouter},
		{"OffRamp", chainState.OffRamp.String(), deployment.OffRampProgramName, cfg.VerifyOffRamp},
		{"RMN Remote", chainState.RMNRemote.String(), deployment.RMNRemoteProgramName, cfg.VerifyRMNRemote},
		{"Burn Mint Token Pool", chainState.BurnMintTokenPool.String(), deployment.BurnMintTokenPoolProgramName, cfg.VerifyBurnMintTokenPool},
		{"Lock Release Token Pool", chainState.LockReleaseTokenPool.String(), deployment.LockReleaseTokenPoolProgramName, cfg.VerifyLockReleaseTokenPool},
		{"Access Controller", mcmState.AccessControllerProgram.String(), deployment.AccessControllerProgramName, cfg.VerifyAccessController},
		{"MCM", mcmState.McmProgram.String(), deployment.McmProgramName, cfg.VerifyMCM},
		{"Timelock", mcmState.TimelockProgram.String(), deployment.TimelockProgramName, cfg.VerifyTimelock},
	}

	for _, v := range verifications {
		if !v.enabled {
			continue
		}

		e.Logger.Debugw("Verifying program", "name", v.name, "programID", v.programID, "programLib", v.programLib)
		err := runSolanaVerify(
			chain.URL,
			v.programID,
			v.programLib,
			cfg.GitCommitSha,
			anchorDir,
			cfg.RemoteVerification,
		)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error verifying %s: %w", v.name, err)
		}
	}

	return deployment.ChangesetOutput{}, nil
}
