package solana

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	csState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
)

var _ deployment.ChangeSet[VerifyBuildConfig] = VerifyBuild

func runCommandStreaming(name string, args []string, dir string, onLine func(string)) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error getting stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error getting stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %w", err)
	}

	scanner := bufio.NewScanner(io.MultiReader(stdoutPipe, stderrPipe))
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			onLine(line)
		}
	}()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}
	return nil
}

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

	cmdArgs = []string{
		"verify-from-repo",
		"-u", "https://solana-mainnet.core.chainstack.com/2faf655d61d30edc66df431358f2d7cf",
		"--program-id", "Ccip842gzYHhvdDkSyi2YVCoAWPbYJoApMFzSxQroE9C",
		"--library-name", "ccip_router",
		strings.TrimSuffix(repoURL, ".git"),
		"--commit-hash", commitHash,
		"--mount-path", mountPath,
	}

	// Add --remote flag if remote verification is enabled
	if remote {
		cmdArgs = append(cmdArgs, "--remote")
	}

	done := make(chan struct{})
	var injectionPath string

	err := runCommandStreaming("solana-verify", cmdArgs, ".", func(line string) {
		fmt.Println(line)

		// Watch for the path creation
		if strings.HasPrefix(line, "Build path: ") {
			// Extract the path
			trimmed := strings.TrimPrefix(line, "Build path: ")
			trimmed = strings.Trim(trimmed, `"`)
			injectionPath = trimmed

			go func() {
				// Wait for path to exist
				for {
					if _, err := os.Stat(injectionPath); err == nil {
						break
					}
					time.Sleep(10 * time.Millisecond)
				}

				// Inject your custom file
				err := os.WriteFile(filepath.Join(injectionPath, "CCIP_BUILD_GIT_HASH"), []byte(commitHash), 0644)
				if err != nil {
					fmt.Println("Error injecting file:", err)
				}
				close(done)
			}()
		}
	})

	// Wait for injection to complete before returning
	<-done
	if err != nil {
		return fmt.Errorf("solana program verification failed: %w", err)
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
