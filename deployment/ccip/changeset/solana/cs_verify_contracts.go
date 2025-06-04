package solana

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	csState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	mcmsTypes "github.com/smartcontractkit/mcms/types"
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
	// if program is owned by deployer key, set to true
	// verification and remote job submission will be done in the same call to this changeset
	// if program is owned by timelock signer, set to false in the first call to this changeset
	// get the proposal -> signed -> executed on chain
	// once thats done, call this changeset again, set to true and it will submit the remote job
	RemoteVerification bool
	// set this to timelock signer if the upgrade authority of the program is the timelock signer
	// if not, set to deployer key
	UpgradeAuthority solana.PublicKey
}

func runSolanaVerify(e cldf.Environment,
	cfg VerifyBuildConfig,
	chain cldf_solana.Chain,
	programID, libraryName, mountPath string,
	timelockSignerPDA solana.PublicKey,
	mcmsTxs []mcmsTypes.Transaction,
) error {
	params := map[string]string{
		"Keypair Path": chain.KeypairPath,
		"Network URL":  chain.URL,
		"Program ID":   programID,
		"Lib Name":     libraryName,
		"Commit Hash":  cfg.GitCommitSha,
		"Mount Path":   mountPath,
	}
	log, err := json.MarshalIndent(params, "", "")
	if err != nil {
		return err
	}
	e.Logger.Infow("solana verify params", "params", string(log))

	cmdArgs := []string{
		"config",
		"set",
		"--keypair", chain.KeypairPath,
	}
	output, err := runCommand("solana", cmdArgs, ".")
	e.Logger.Infow("solana config set output", "output", output)
	if err != nil {
		return fmt.Errorf("failed to set keypair during program verification: %s %w", output, err)
	}

	if !timelockSignerPDA.IsZero() && cfg.UpgradeAuthority == timelockSignerPDA {
		// enter here only if mcms tx has been signed and submitted to the chain
		// https://solana.com/developers/guides/advanced/verified-builds#7-submit-remote-verification-job
		if cfg.RemoteVerification {
			cmdArgs = []string{
				"remote",
				"submit-job",
				"--url", chain.URL,
				"--uploader", timelockSignerPDA.String(),
				"--program-id", programID,
			}
			output, err := runCommand("solana-verify", cmdArgs, chain.ProgramsPath)
			e.Logger.Infow("remote submit-job output", "output", output)
			if err != nil {
				return fmt.Errorf("solana program verification failed: %s %w", output, err)
			}
			// only need to submit job this time as we are assuming here that the mcms tx has been signed and submitted to the chain
			return nil
		}
		cmdArgs = []string{
			"export-pda-tx",
			"--url", chain.URL,
			"--program-id", programID,
			"--library-name", libraryName,
			strings.TrimSuffix(repoURL, ".git"),
			"--commit-hash", cfg.GitCommitSha,
			"--mount-path", mountPath,
			"--uploader", timelockSignerPDA.String(),
		}

		output, err = runCommand("solana-verify", cmdArgs, ".")
		e.Logger.Infow("export-pda-tx output", "output", output)
		if err != nil {
			return fmt.Errorf("solana program verification failed: %s %w", output, err)
		}

		resolvedIxn, err := getIxnFromEncodedTx(e, output, timelockSignerPDA)
		if err != nil {
			return fmt.Errorf("failed to get ixn from encoded tx: %w", err)
		}
		if resolvedIxn == nil {
			return fmt.Errorf("failed to get ixn from encoded tx")
		}

		upgradeTx, err := BuildMCMSTxn(resolvedIxn, programID, cldf.ContractType(libraryName))
		if err != nil {
			return fmt.Errorf("failed to build upgrade transaction: %w", err)
		}
		if upgradeTx != nil {
			e.Logger.Infow("upgradeTx", "tx", upgradeTx)
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	} else {
		cmdArgs = []string{
			"verify-from-repo",
			"--url", chain.URL,
			"--program-id", programID,
			"--library-name", libraryName,
			strings.TrimSuffix(repoURL, ".git"),
			"--commit-hash", cfg.GitCommitSha,
			"--mount-path", mountPath,
			"--skip-prompt",
		}

		output, err = runCommand("solana-verify", cmdArgs, ".")
		e.Logger.Infow("verify-from-repo output", "output", output)
		if err != nil {
			return fmt.Errorf("solana program verification failed: %s %w", output, err)
		}
		if cfg.RemoteVerification {
			cmdArgs = []string{
				"remote",
				"submit-job",
				"--url", chain.URL,
				"--uploader", chain.DeployerKey.PublicKey().String(),
				"--program-id", programID,
			}
			output, err := runCommand("solana-verify", cmdArgs, chain.ProgramsPath)
			e.Logger.Infow("remote submit-job output", "output", output)
			if err != nil {
				return fmt.Errorf("solana program verification failed: %s %w", output, err)
			}
		}
	}

	return nil
}

func getIxnFromEncodedTx(e cldf.Environment, output string, timelockSignerPDA solana.PublicKey) (*solana.GenericInstruction, error) {
	lines := strings.Split(string(output), "\n")
	var base58EncodedTx string
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) != "" {
			base58EncodedTx = strings.TrimSpace(lines[i])
			break
		}
	}
	if base58EncodedTx == "" {
		return nil, fmt.Errorf("failed to extract base58-encoded transaction")
	}
	e.Logger.Infow("base58-encoded transaction", "tx", base58EncodedTx)

	txBytes, err := base58.Decode(base58EncodedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base58-encoded transaction: %w", err)
	}
	e.Logger.Infow("txBytes", "txBytes", txBytes)
	tx, err := solana.TransactionFromBytes(txBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction from bytes: %w", err)
	}
	inst := tx.Message.Instructions[0]
	resolved, err := resolveCompiledInstruction(timelockSignerPDA, tx.Message, inst)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve compiled instruction: %w", err)
	}
	return resolved, nil
}

func resolveCompiledInstruction(
	timelockSignerPDA solana.PublicKey,
	msg solana.Message,
	compiled solana.CompiledInstruction,
) (*solana.GenericInstruction, error) {
	accounts := make(solana.AccountMetaSlice, len(compiled.Accounts))
	for i, idx := range compiled.Accounts {
		if int(idx) >= len(msg.AccountKeys) {
			return nil, fmt.Errorf("account index out of range: %d", idx)
		}
		pub := msg.AccountKeys[idx]
		isSigner := msg.IsSigner(pub)

		isWritable, err := msg.IsWritable(pub)
		if err != nil {
			return nil, fmt.Errorf("failed to check if account is writable: %w", err)
		}
		accounts[i] = &solana.AccountMeta{
			PublicKey:  pub,
			IsSigner:   isSigner,
			IsWritable: isWritable,
		}
	}
	if int(compiled.ProgramIDIndex) >= len(msg.AccountKeys) {
		return nil, fmt.Errorf("program ID index out of range: %d", compiled.ProgramIDIndex)
	}

	programID := msg.AccountKeys[compiled.ProgramIDIndex]

	data, err := base58.Decode(compiled.Data.String())
	if err != nil {
		fmt.Errorf("failed to decode instruction data: %w", err)
	}

	return &solana.GenericInstruction{
		ProgID:        programID,
		AccountValues: accounts,
		DataBytes:     data,
	}, nil
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
	var timelockSignerPDA solana.PublicKey
	if mcmState != nil {
		timelockSignerPDA = csState.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
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
	mcmsTxs := make([]mcmsTypes.Transaction, 0)
	for _, v := range verifications {
		if !v.enabled {
			continue
		}

		e.Logger.Debugw("Verifying program", "name", v.name, "programID", v.programID, "programLib", v.programLib)
		err := runSolanaVerify(
			e,
			cfg,
			chain,
			v.programID,
			v.programLib,
			anchorDir,
			timelockSignerPDA,
			mcmsTxs,
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error verifying %s: %w", v.name, err)
		}
	}
	if len(mcmsTxs) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to verify CCIP contracts", cfg.MCMS.MinDelay, mcmsTxs)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}

		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}
