package solana

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	mcmsTypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

const (
	SolanaVerifyProgramID = "verifycLy8mB96wd9wqq3WDXQwM4oU6r42Th37Db9fC"
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
	VerifyCCTPTokenPool          bool
	VerifyAccessController       bool
	VerifyMCM                    bool
	VerifyTimelock               bool
	// if program is owned by deployer key
	// set to true -> verification and remote job submission will be done in the same call to this changeset
	// if program is owned by timelock signer
	// set to false in the first call to this changeset
	// get the proposal -> signed -> executed on chain
	// once thats done, call this changeset again, set to true and it will submit the remote job
	RemoteVerification bool
	// set to the same as upgrade authority of the program
	// timelock signer or deployer key
	UpgradeAuthority solana.PublicKey
	MCMS             *proposalutils.TimelockConfig
}

func runSolanaVerifyMCMS(e cldf.Environment,
	cfg VerifyBuildConfig,
	chain cldf_solana.Chain,
	programID, libraryName, mountPath string,
	timelockSignerPDA solana.PublicKey,
	mcmsTxs *[]mcmsTypes.Transaction,
) error {
	// enter here only if mcms tx has been signed and submitted to the chain
	// https://solana.com/developers/guides/advanced/verified-builds#7-submit-remote-verification-job
	if cfg.RemoteVerification {
		cmdArgs := []string{
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

	// run cli command
	cmdArgs := []string{
		"export-pda-tx",
		"--url", chain.URL,
		"--program-id", programID,
		"--library-name", libraryName,
		strings.TrimSuffix(repoURL, ".git"),
		"--commit-hash", cfg.GitCommitSha,
		"--mount-path", mountPath,
		"--uploader", timelockSignerPDA.String(),
	}
	e.Logger.Infow("export-pda-tx cmdArgs", "cmdArgs", cmdArgs)
	output, err := runCommand("solana-verify", cmdArgs, ".")
	e.Logger.Infow("export-pda-tx output", "output", output)
	if err != nil {
		return fmt.Errorf("solana program verification failed: %s %w", output, err)
	}

	// get ix from output
	resolvedIxn, err := getIxnFromEncodedTx(e, output, timelockSignerPDA)
	if err != nil {
		return fmt.Errorf("failed to get ixn from encoded tx: %w", err)
	}
	if resolvedIxn == nil {
		return errors.New("failed to get ixn from encoded tx")
	}

	// build mcms tx from ix
	upgradeTx, err := BuildMCMSTxn(resolvedIxn, resolvedIxn.ProgID.String(), cldf.ContractType(libraryName))
	if err != nil {
		return fmt.Errorf("failed to build upgrade transaction: %w", err)
	}
	if upgradeTx != nil {
		e.Logger.Infow("upgradeTx", "tx", upgradeTx)
		*mcmsTxs = append(*mcmsTxs, *upgradeTx)
	}
	return nil
}

func runSolanaVerifyWithoutMCMS(e cldf.Environment,
	cfg VerifyBuildConfig,
	chain cldf_solana.Chain,
	programID, libraryName, mountPath string,
	timelockSignerPDA solana.PublicKey,
) error {
	// if timelock signer does not exist
	// or user has set the upgrade authority to the deployer key
	// then we need to run the cli command
	cmdArgs := []string{
		"verify-from-repo",
		"--url", chain.URL,
		"--program-id", programID,
		"--library-name", libraryName,
		strings.TrimSuffix(repoURL, ".git"),
		"--commit-hash", cfg.GitCommitSha,
		"--mount-path", mountPath,
		"--skip-prompt",
	}

	output, err := runCommand("solana-verify", cmdArgs, ".")
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
	return nil
}

func runSolanaVerify(e cldf.Environment,
	cfg VerifyBuildConfig,
	chain cldf_solana.Chain,
	programID, libraryName, mountPath string,
	timelockSignerPDA solana.PublicKey,
	mcmsTxs *[]mcmsTypes.Transaction,
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

	// if timelock signer exists
	// and user has set the upgrade authority to the timelock signer
	// then we need to create mcms txs
	if !timelockSignerPDA.IsZero() && cfg.UpgradeAuthority == timelockSignerPDA {
		return runSolanaVerifyMCMS(e, cfg, chain, programID, libraryName, mountPath, timelockSignerPDA, mcmsTxs)
	}
	return runSolanaVerifyWithoutMCMS(e, cfg, chain, programID, libraryName, mountPath, timelockSignerPDA)
}

// each tx contains 2 things
// 1. message (which contains a list of instructions)
// 2. signatures
// we will extract out the relevant instruction from the message by decoding the different layers
func getIxnFromEncodedTx(e cldf.Environment, output string, timelockSignerPDA solana.PublicKey) (*solana.GenericInstruction, error) {
	// get the base58-encoded transaction from the output
	// this is based on the current tx output format
	// if solana-verify cli changes the output format, this will break
	lines := strings.Split(output, "\n")
	var base58EncodedTx string
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) != "" {
			base58EncodedTx = strings.TrimSpace(lines[i])
			break
		}
	}
	if base58EncodedTx == "" {
		return nil, errors.New("failed to extract base58-encoded transaction")
	}
	e.Logger.Infow("base58-encoded transaction", "tx", base58EncodedTx)

	// create a transaction object from the base58EncodedTx
	tx, err := solana.TransactionFromBase58(base58EncodedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction from bytes: %w", err)
	}

	// we will now find the instruction within the tx that is being executed on the verify program

	// list of all accounts in the tx
	txAccountList := tx.Message.AccountKeys
	for _, inst := range tx.Message.Instructions {
		// this should not happen unless solana-verify cli has a bug
		if int(inst.ProgramIDIndex) >= len(txAccountList) {
			return nil, fmt.Errorf("program ID index out of range: %d", inst.ProgramIDIndex)
		}
		// the programID on which this instruction is being executed
		programID := txAccountList[inst.ProgramIDIndex]
		// if its the verify program, resolve the instruction
		// this is the ix we need to get signed by mcms
		if programID.String() == SolanaVerifyProgramID {
			resolved, err := resolveVerifyInstruction(e, timelockSignerPDA, tx.Message, inst)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve the verify instruction: %w", err)
			}
			return resolved, nil
		}
	}
	return nil, errors.New("failed to find verify instruction")
}

// this function takes in an ix which is part of a tx
// and creates an independent ixn object
// which will then be used to create an independent mcms tx
// it uses the tx.Message to do so
func resolveVerifyInstruction(
	e cldf.Environment,
	timelockSignerPDA solana.PublicKey,
	msg solana.Message,
	verifyIxn solana.CompiledInstruction,
) (*solana.GenericInstruction, error) {
	// the programID on which this instruction is being executed
	// this must be the verify program
	programID := msg.AccountKeys[verifyIxn.ProgramIDIndex]
	// the data of the instruction that we want to run on the verify programs
	data, err := base58.Decode(verifyIxn.Data.String())
	if err != nil {
		return nil, fmt.Errorf("failed to decode instruction data: %w", err)
	}

	// the accounts that are being used in the instruction
	// we need to get the pubkeys of these accounts from the tx.Message
	// and then create an AccountMeta object for each account
	// this is the list of accounts that will be used to create the verify ix
	accounts := make(solana.AccountMetaSlice, len(verifyIxn.Accounts))
	for i, idx := range verifyIxn.Accounts {
		if int(idx) >= len(msg.AccountKeys) {
			return nil, fmt.Errorf("account index out of range: %d", idx)
		}
		accountPubKey := msg.AccountKeys[idx]
		isSigner := msg.IsSigner(accountPubKey)
		isWritable, err := msg.IsWritable(accountPubKey)
		if err != nil {
			return nil, fmt.Errorf("failed to check if account is writable: %w", err)
		}
		accounts[i] = &solana.AccountMeta{
			PublicKey:  accountPubKey,
			IsSigner:   isSigner,
			IsWritable: isWritable,
		}
	}

	return &solana.GenericInstruction{
		ProgID:        programID,
		AccountValues: accounts,
		DataBytes:     data,
	}, nil
}
