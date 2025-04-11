package solana

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

const IdlIxTag uint64 = 0x0a69e9a778bcf440

// IDL
type IDLConfig struct {
	ChainSelector        uint64
	GitCommitSha         string
	Router               bool
	FeeQuoter            bool
	OffRamp              bool
	RMNRemote            bool
	BurnMintTokenPool    bool
	LockReleaseTokenPool bool
}

func parseAnchorVersion(output string) (string, error) {
	const prefix = "anchor-cli "
	if strings.HasPrefix(output, prefix) {
		return strings.TrimSpace(strings.TrimPrefix(output, prefix)), nil
	}
	return "", fmt.Errorf("unexpected version output: %q", output)
}

func updateToml(e deployment.Environment) error {
	output, err := runCommand("anchor", []string{"--version"}, ".")
	if err != nil {
		return errors.New("anchor-cli not installed in path")
	}
	anchorVersion, err := parseAnchorVersion(output)
	if err != nil {
		return fmt.Errorf("error parsing anchor version: %w", err)
	}

	anchorTomlPath := filepath.Join(cloneDir, anchorDir, "Anchor.toml")

	// Read the whole file as text
	contentBytes, err := os.ReadFile(anchorTomlPath)
	if err != nil {
		return fmt.Errorf("failed to read Anchor.toml: %w", err)
	}
	content := string(contentBytes)

	// Replace exact version string
	replaceText := fmt.Sprintf("anchor_version = \"%s\"", anchorVersion)
	content = strings.Replace(content, `anchor_version = "0.29.0"`, replaceText, 1)

	// Write it back
	if err := os.WriteFile(anchorTomlPath, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write Anchor.toml: %w", err)
	}
	return nil
}

func repoSetup(e deployment.Environment, c IDLConfig) error {
	// setup
	var version string
	var err error
	if c.GitCommitSha == "" {
		version, err = memory.GetSha()
		if err != nil {
			return fmt.Errorf("error getting sha: %w", err)
		}
	} else {
		version = c.GitCommitSha
	}
	if err := cloneRepo(e, version, false); err != nil {
		return fmt.Errorf("error cloning repo: %w", err)
	}

	e.Logger.Debug("Updating Anchor.toml")
	if err := updateToml(e); err != nil {
		return fmt.Errorf("error updating Anchor.toml: %w", err)
	}
	return nil
}

func updateIDL(e deployment.Environment, idlFile string, programID string) error {
	e.Logger.Debug("Reading IDL")
	idlBytes, err := os.ReadFile(idlFile)
	if err != nil {
		return fmt.Errorf("failed to read IDL: %w", err)
	}
	e.Logger.Debug("Parsing IDL")
	var idl map[string]interface{}
	if err := json.Unmarshal(idlBytes, &idl); err != nil {
		return fmt.Errorf("failed to parse legacy IDL: %w", err)
	}
	e.Logger.Debug("Updating IDL with program ID", "programID", programID)
	idl["metadata"] = map[string]interface{}{
		"address": programID,
	}
	// Marshal updated IDL back to JSON
	e.Logger.Debug("Marshalling updated IDL")
	updatedIDLBytes, err := json.MarshalIndent(idl, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated IDL: %w", err)
	}
	e.Logger.Debug("Writing updated IDL")
	// Write updated IDL back to file
	if err := os.WriteFile(idlFile, updatedIDLBytes, 0600); err != nil {
		return fmt.Errorf("failed to write updated IDL: %w", err)
	}
	return nil
}

func getIDL(e deployment.Environment, programID string, programName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting cwd: %w", err)
	}
	idlFile := filepath.Join(cwd, cloneDir, anchorDir, "target", "idl", programName+".json")
	if _, err := os.Stat(idlFile); err != nil {
		return "", fmt.Errorf("idl file not found: %w", err)
	}
	e.Logger.Debug("Updating IDL")
	err = updateIDL(e, idlFile, programID)
	if err != nil {
		return "", fmt.Errorf("error updating IDL: %w", err)
	}
	return idlFile, nil
}

func idlInit(e deployment.Environment, chain deployment.SolChain, programID, programName string) error {
	idlFile, err := getIDL(e, programID, programName)
	if err != nil {
		return fmt.Errorf("error getting IDL: %w", err)
	}
	e.Logger.Infow("Uploading IDL", "programName", programName)
	args := []string{"idl", "init", "--filepath", idlFile, "--provider.wallet", chain.KeypairPath, "--provider.cluster", chain.URL, programID}
	e.Logger.Info(args)
	_, err = runCommand("anchor", args, filepath.Join(cloneDir, anchorDir))
	if err != nil {
		return fmt.Errorf("error uploading idl: %w", err)
	}
	return nil
}

func getTimelockSignerPDA(e deployment.Environment, chainSelector uint64) (solana.PublicKey, error) {
	addresses, _ := e.ExistingAddresses.AddressesForChain(chainSelector)
	mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(e.SolChains[chainSelector], addresses)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("error loading mcm state: %w", err)
	}
	timelockSignerPDA := state.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
	return timelockSignerPDA, nil
}

func setIdlAuthority(e deployment.Environment, chain deployment.SolChain, newAuthority, programID, programName, bufferAccount string) error {
	e.Logger.Infow("Setting IDL authority", "programName", programName, "newAuthority", newAuthority)
	args := []string{"idl", "set-authority", "-n", newAuthority, "-p", programID, "--provider.wallet", chain.KeypairPath, "--provider.cluster", chain.URL}
	if bufferAccount != "" {
		e.Logger.Infow("Setting IDL authority for buffer", "bufferAccount", bufferAccount)
		args = append(args, bufferAccount)
	}
	e.Logger.Info(args)
	_, err := runCommand("anchor", args, filepath.Join(cloneDir, anchorDir))
	if err != nil {
		return fmt.Errorf("error setting idl authority: %w", err)
	}
	return nil
}

func getIDLAddress(e deployment.Environment, programID solana.PublicKey) (solana.PublicKey, error) {
	base, _, _ := solana.FindProgramAddress([][]byte{}, programID)
	idlAddress, _ := solana.CreateWithSeed(base, "anchor:idl", programID)
	e.Logger.Infof("IDL Address:  %s", idlAddress.String())
	return idlAddress, nil
}

func parseIdlBuffer(output string) (string, error) {
	const prefix = "Idl buffer created: "
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix)), nil
		}
	}
	return "", errors.New("failed to find IDL buffer in output")
}

func writeBuffer(e deployment.Environment, chain deployment.SolChain, programID string) (solana.PublicKey, error) {
	idlFile, err := getIDL(e, programID, deployment.RouterProgramName)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("error getting IDL: %w", err)
	}
	e.Logger.Infow("Writing IDL buffer", "programID", programID)
	args := []string{"idl", "write-buffer", "--filepath", idlFile, "--provider.wallet", chain.KeypairPath, "--provider.cluster", chain.URL, programID}
	e.Logger.Info(args)
	output, err := runCommand("anchor", args, filepath.Join(cloneDir, anchorDir))
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("error writing IDL buffer: %w", err)
	}
	e.Logger.Infow("Parsing IDL buffer", "programID", programID)
	buffer, err := parseIdlBuffer(output)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("error parsing IDL buffer: %w", err)
	}
	bufferAddress, err := solana.PublicKeyFromBase58(buffer)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("error parsing IDL buffer: %w", err)
	}
	return bufferAddress, nil
}

func setBufferIx(e deployment.Environment, programID, buffer, authority solana.PublicKey) (solana.GenericInstruction, error) {
	idlAddress, err := getIDLAddress(e, programID)
	if err != nil {
		return solana.GenericInstruction{}, fmt.Errorf("error getting idl address for %s: %w", programID.String(), err)
	}
	data := binary.LittleEndian.AppendUint64([]byte{}, IdlIxTag) // 4-byte Extend instruction identifier
	data = append(data, byte(3))

	instruction := solana.NewInstruction(
		programID,
		solana.AccountMetaSlice{
			solana.NewAccountMeta(buffer, true, false),
			solana.NewAccountMeta(idlAddress, true, false),
			solana.NewAccountMeta(authority, false, true),
		},
		data,
	)
	return *instruction, nil
}

func upgradeIDLIx(e deployment.Environment, chain deployment.SolChain, programID string, timelockSignerPDA solana.PublicKey) (*mcmsTypes.Transaction, error) {
	buffer, err := writeBuffer(e, chain, programID)
	if err != nil {
		return nil, fmt.Errorf("error writing buffer: %w", err)
	}
	err = setIdlAuthority(e, chain, timelockSignerPDA.String(), programID, deployment.RouterProgramName, buffer.String())
	if err != nil {
		return nil, fmt.Errorf("error setting buffer authority: %w", err)
	}
	instruction, err := setBufferIx(e, solana.MustPublicKeyFromBase58(programID), buffer, timelockSignerPDA)
	if err != nil {
		return nil, fmt.Errorf("error generating set buffer ix: %w", err)
	}
	upgradeTx, err := BuildMCMSTxn(&instruction, programID, deployment.RMNRemoteProgramName)
	if err != nil {
		return nil, fmt.Errorf("failed to create upgrade transaction: %w", err)
	}
	return upgradeTx, nil
}

func (c IDLConfig) Validate(e deployment.Environment) error {
	if err := deployment.IsValidChainSelector(c.ChainSelector); err != nil {
		return fmt.Errorf("invalid chain selector: %d - %w", c.ChainSelector, err)
	}
	family, _ := chainsel.GetSelectorFamily(c.ChainSelector)
	if family != chainsel.FamilySolana {
		return fmt.Errorf("chain %d is not a solana chain", c.ChainSelector)
	}
	existingState, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load existing onchain state: %w", err)
	}
	if _, exists := existingState.SupportedChains()[c.ChainSelector]; !exists {
		return fmt.Errorf("chain %d not supported", c.ChainSelector)
	}
	chainState := existingState.SolChains[c.ChainSelector]
	if c.Router && chainState.Router.IsZero() {
		return fmt.Errorf("router not deployed for chain %d, cannot upload idl", c.ChainSelector)
	}
	if c.FeeQuoter && chainState.FeeQuoter.IsZero() {
		return fmt.Errorf("feeQuoter not deployed for chain %d, cannot upload idl", c.ChainSelector)
	}
	if c.OffRamp && chainState.OffRamp.IsZero() {
		return fmt.Errorf("offRamp not deployed for chain %d, cannot upload idl", c.ChainSelector)
	}
	if c.RMNRemote && chainState.Router.IsZero() {
		return fmt.Errorf("rmnRemote not deployed for chain %d, cannot upload idl", c.ChainSelector)
	}
	if c.BurnMintTokenPool && chainState.FeeQuoter.IsZero() {
		return fmt.Errorf("burnMintTokenPool not deployed for chain %d, cannot upload idl", c.ChainSelector)
	}
	if c.LockReleaseTokenPool && chainState.OffRamp.IsZero() {
		return fmt.Errorf("lockReleaseTokenPool not deployed for chain %d, cannot upload idl", c.ChainSelector)
	}
	return nil
}

// changeset to upload idl for a program
func UploadIDL(e deployment.Environment, c IDLConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error validating idl config: %w", err)
	}
	chain := e.SolChains[c.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[c.ChainSelector]

	if err := repoSetup(e, c); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error setting up anchor workspace: %w", err)
	}
	// start uploading
	if c.Router {
		err := idlInit(e, chain, chainState.Router.String(), deployment.RouterProgramName)
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.FeeQuoter {
		err := idlInit(e, chain, chainState.FeeQuoter.String(), deployment.FeeQuoterProgramName)
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.OffRamp {
		err := idlInit(e, chain, chainState.OffRamp.String(), deployment.OffRampProgramName)
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.RMNRemote {
		err := idlInit(e, chain, chainState.RMNRemote.String(), deployment.RMNRemoteProgramName)
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.BurnMintTokenPool {
		err := idlInit(e, chain, chainState.BurnMintTokenPool.String(), deployment.BurnMintTokenPoolProgramName)
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.LockReleaseTokenPool {
		err := idlInit(e, chain, chainState.LockReleaseTokenPool.String(), deployment.LockReleaseTokenPoolProgramName)
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}

	return deployment.ChangesetOutput{}, nil
}

// changeset to set idl authority for a program to timelock
func SetAuthorityIDL(e deployment.Environment, c IDLConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error validating idl config: %w", err)
	}
	chain := e.SolChains[c.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[c.ChainSelector]

	timelockSignerPDA, err := getTimelockSignerPDA(e, c.ChainSelector)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error loading timelockSignerPDA: %w", err)
	}

	// set idl authority
	if c.Router {
		err = setIdlAuthority(e, chain, timelockSignerPDA.String(), chainState.Router.String(), deployment.RouterProgramName, "")
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.FeeQuoter {
		err = setIdlAuthority(e, chain, timelockSignerPDA.String(), chainState.FeeQuoter.String(), deployment.FeeQuoterProgramName, "")
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.OffRamp {
		err = setIdlAuthority(e, chain, timelockSignerPDA.String(), chainState.OffRamp.String(), deployment.OffRampProgramName, "")
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.RMNRemote {
		err = setIdlAuthority(e, chain, timelockSignerPDA.String(), chainState.RMNRemote.String(), deployment.RMNRemoteProgramName, "")
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.BurnMintTokenPool {
		err = setIdlAuthority(e, chain, timelockSignerPDA.String(), chainState.BurnMintTokenPool.String(), deployment.BurnMintTokenPoolProgramName, "")
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.LockReleaseTokenPool {
		err = setIdlAuthority(e, chain, timelockSignerPDA.String(), chainState.LockReleaseTokenPool.String(), deployment.LockReleaseTokenPoolProgramName, "")
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}

	return deployment.ChangesetOutput{}, nil
}

// changeset to upgrade idl for a program via timelock
// write buffer using anchor cli
// set buffer authority to timelock using anchor cli
// generate set buffer ix using solana-go sdk
// build mcms txn to upgrade idl
func UpgradeIDL(e deployment.Environment, c IDLConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error validating idl config: %w", err)
	}
	chain := e.SolChains[c.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[c.ChainSelector]

	timelockSignerPDA, err := getTimelockSignerPDA(e, c.ChainSelector)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error loading timelockSignerPDA: %w", err)
	}
	mcmsTxs := make([]mcmsTypes.Transaction, 0)
	if c.Router {
		upgradeTx, err := upgradeIDLIx(e, chain, chainState.Router.String(), timelockSignerPDA)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		mcmsTxs = append(mcmsTxs, *upgradeTx)
	}
	if c.FeeQuoter {
		upgradeTx, err := upgradeIDLIx(e, chain, chainState.FeeQuoter.String(), timelockSignerPDA)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		mcmsTxs = append(mcmsTxs, *upgradeTx)
	}
	if c.OffRamp {
		upgradeTx, err := upgradeIDLIx(e, chain, chainState.Router.String(), timelockSignerPDA)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		mcmsTxs = append(mcmsTxs, *upgradeTx)
	}
	if c.RMNRemote {
		upgradeTx, err := upgradeIDLIx(e, chain, chainState.Router.String(), timelockSignerPDA)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		mcmsTxs = append(mcmsTxs, *upgradeTx)
	}
	if c.BurnMintTokenPool {
		upgradeTx, err := upgradeIDLIx(e, chain, chainState.Router.String(), timelockSignerPDA)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		mcmsTxs = append(mcmsTxs, *upgradeTx)
	}
	if c.LockReleaseTokenPool {
		upgradeTx, err := upgradeIDLIx(e, chain, chainState.Router.String(), timelockSignerPDA)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		mcmsTxs = append(mcmsTxs, *upgradeTx)
	}

	proposal, err := BuildProposalsForTxns(
		e, c.ChainSelector, "proposal to upgrade CCIP contracts", 1*time.Second, mcmsTxs)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	// do we need to batch this ?
	return deployment.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
	}, nil
}
