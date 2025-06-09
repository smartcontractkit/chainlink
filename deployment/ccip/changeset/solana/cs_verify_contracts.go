package solana

import (
	"encoding/json"
	"fmt"
	"strings"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	csState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// https://solana.com/developers/guides/advanced/verified-builds
type VerifyBuildConfig struct {
	GitCommitSha                 string
	ChainSelector                uint64
	VerifyFeeQuoter              bool
	VerifyRouter                 bool
	VerifyOffRamp                bool
	VerifyRMNRemote              bool
	BurnMintTokenPoolMetadata    []string
	LockReleaseTokenPoolMetadata []string
	VerifyAccessController       bool
	VerifyMCM                    bool
	VerifyTimelock               bool
	RemoteVerification           bool
	MCMS                         *proposalutils.TimelockConfig
}

func runSolanaVerify(chain cldf_solana.Chain, programID, libraryName, commitHash, mountPath string, remote bool) error {
	params := map[string]string{
		"Keypair Path": chain.KeypairPath,
		"Network URL":  chain.URL,
		"Program ID":   programID,
		"Lib Name":     libraryName,
		"Commit Hash":  commitHash,
		"Mount Path":   mountPath,
	}
	log, err := json.MarshalIndent(params, "", "")
	if err != nil {
		return err
	}
	fmt.Println(string(log))

	cmdArgs := []string{
		"config",
		"set",
		"--keypair", chain.KeypairPath,
	}
	output, err := runCommand("solana", cmdArgs, ".")
	fmt.Println(output)
	if err != nil {
		return fmt.Errorf("solana program verification failed: %s %w", output, err)
	}

	cmdArgs = []string{
		"verify-from-repo",
		"--url", chain.URL,
		"--program-id", programID,
		"--library-name", libraryName,
		strings.TrimSuffix(repoURL, ".git"),
		"--commit-hash", commitHash,
		"--mount-path", mountPath,
		"--skip-prompt",
	}

	output, err = runCommand("solana-verify", cmdArgs, ".")
	fmt.Println(output)
	if err != nil {
		return fmt.Errorf("solana program verification failed: %s %w", output, err)
	}

	// Add --remote flag if remote verification is enabled
	if remote {
		cmdArgs = []string{
			"remote",
			"submit-job",
			"--url", chain.URL,
			"--uploader", chain.DeployerKey.PublicKey().String(),
			"--program-id", programID,
		}
		output, err := runCommand("solana-verify", cmdArgs, chain.ProgramsPath)
		fmt.Println(output)
		if err != nil {
			return fmt.Errorf("solana program verification failed: %s %w", output, err)
		}
	}

	return nil
}

func VerifyBuild(e cldf.Environment, cfg VerifyBuildConfig) (cldf.ChangesetOutput, error) {
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	state, _ := stateview.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]

	addresses, err := e.ExistingAddresses.AddressesForChain(cfg.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get existing addresses: %w", err)
	}
	mcmState, err := csState.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
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
		{"Access Controller", mcmState.AccessControllerProgram.String(), deployment.AccessControllerProgramName, cfg.VerifyAccessController},
		{"MCM", mcmState.McmProgram.String(), deployment.McmProgramName, cfg.VerifyMCM},
		{"Timelock", mcmState.TimelockProgram.String(), deployment.TimelockProgramName, cfg.VerifyTimelock},
	}
	for _, bnmMetadata := range cfg.BurnMintTokenPoolMetadata {
		verifications = append(verifications, struct {
			name       string
			programID  string
			programLib string
			enabled    bool
		}{
			name:       "Burn Mint Token Pool",
			programID:  chainState.BurnMintTokenPools[bnmMetadata].String(),
			programLib: deployment.BurnMintTokenPoolProgramName,
			enabled:    true,
		})
	}

	for _, lnrMetadata := range cfg.LockReleaseTokenPoolMetadata {
		verifications = append(verifications, struct {
			name       string
			programID  string
			programLib string
			enabled    bool
		}{
			name:       "Lock Release Token Pool",
			programID:  chainState.LockReleaseTokenPools[lnrMetadata].String(),
			programLib: deployment.LockReleaseTokenPoolProgramName,
			enabled:    true,
		})
	}

	for _, v := range verifications {
		if !v.enabled {
			continue
		}

		e.Logger.Debugw("Verifying program", "name", v.name, "programID", v.programID, "programLib", v.programLib)
		err := runSolanaVerify(
			chain,
			v.programID,
			v.programLib,
			cfg.GitCommitSha,
			anchorDir,
			cfg.RemoteVerification,
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error verifying %s: %w", v.name, err)
		}
	}

	return cldf.ChangesetOutput{}, nil
}
