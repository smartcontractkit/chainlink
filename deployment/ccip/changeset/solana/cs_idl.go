package solana

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gagliardetto/solana-go"
	"github.com/pelletier/go-toml"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
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
	AccessController     bool
	MCM                  bool
	Timelock             bool
	MCMS                 *proposalutils.TimelockConfig
}

// parse anchor version from running anchor --version
func parseAnchorVersion(output string) (string, error) {
	const prefix = "anchor-cli "
	if strings.HasPrefix(output, prefix) {
		return strings.TrimSpace(strings.TrimPrefix(output, prefix)), nil
	}
	return "", fmt.Errorf("unexpected version output: %q", output)
}

// create Anchor.toml file to simulate anchor workspace
func writeAnchorToml(e deployment.Environment, filename, anchorVersion, cluster, wallet string) error {
	e.Logger.Debugw("Writing Anchor.toml", "filename", filename, "anchorVersion", anchorVersion, "cluster", cluster, "wallet", wallet)
	config := map[string]interface{}{
		"toolchain": map[string]string{
			"anchor_version": anchorVersion,
		},
		"provider": map[string]string{
			"cluster": cluster,
			"wallet":  wallet,
		},
	}
	e.Logger.Debugw("Anchor.toml config", "config", config)

	tree, err := toml.TreeFromMap(config)
	if err != nil {
		return fmt.Errorf("failed to build TOML tree: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create TOML file: %w", err)
	}
	defer file.Close()

	if _, err := tree.WriteTo(file); err != nil {
		return fmt.Errorf("failed to write TOML to file: %w", err)
	}

	return nil
}

// resolve artifacts based on sha and write anchor.toml file to simulate anchor workspace
func repoSetup(e deployment.Environment, chain deployment.SolChain, gitCommitSha string) error {
	e.Logger.Debug("Downloading Solana CCIP program artifacts...")
	err := memory.DownloadSolanaCCIPProgramArtifacts(e.GetContext(), chain.ProgramsPath, e.Logger, gitCommitSha)
	if err != nil {
		return fmt.Errorf("error downloading solana ccip program artifacts: %w", err)
	}

	// get anchor version
	output, err := runCommand("anchor", []string{"--version"}, ".")
	if err != nil {
		return errors.New("anchor-cli not installed in path")
	}
	e.Logger.Debugw("Anchor version command output", "output", output)
	anchorVersion, err := parseAnchorVersion(output)
	if err != nil {
		return fmt.Errorf("error parsing anchor version: %w", err)
	}
	// create Anchor.toml
	// this creates anchor workspace with cluster and wallet configured
	if err := writeAnchorToml(e, filepath.Join(chain.ProgramsPath, "Anchor.toml"), anchorVersion, chain.URL, chain.KeypairPath); err != nil {
		return fmt.Errorf("error writing Anchor.toml: %w", err)
	}

	return nil
}

// update IDL with program ID
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
	e.Logger.Debugw("Updating IDL with programID", "programID", programID)
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

// get IDL file and update with program ID
func getIDL(e deployment.Environment, programsPath, programID string, programName string) (string, error) {
	idlFile := filepath.Join(programsPath, programName+".json")
	if _, err := os.Stat(idlFile); err != nil {
		return "", fmt.Errorf("idl file not found: %w", err)
	}
	e.Logger.Debug("Updating IDL")
	err := updateIDL(e, idlFile, programID)
	if err != nil {
		return "", fmt.Errorf("error updating IDL: %w", err)
	}
	return idlFile, nil
}

// initialize IDL for a program
func idlInit(e deployment.Environment, programsPath, programID, programName string) error {
	idlFile, err := getIDL(e, programsPath, programID, programName)
	if err != nil {
		return fmt.Errorf("error getting IDL: %w", err)
	}
	e.Logger.Infow("Uploading IDL", "programName", programName)
	args := []string{"idl", "init", "--filepath", idlFile, programID}
	e.Logger.Info(args)
	output, err := runCommand("anchor", args, programsPath)
	e.Logger.Debugw("IDL init output", "output", output)
	if err != nil {
		e.Logger.Debugw("IDL init error", "error", err)
		return fmt.Errorf("error uploading idl: %w", err)
	}
	e.Logger.Infow("IDL uploaded", "programName", programName)
	return nil
}

// set IDL authority for a program
func setIdlAuthority(e deployment.Environment, newAuthority, programsPath, programID, programName, bufferAccount string) error {
	e.Logger.Infow("Setting IDL authority", "programName", programName, "newAuthority", newAuthority)
	args := []string{"idl", "set-authority", "-n", newAuthority, "-p", programID}
	if bufferAccount != "" {
		e.Logger.Infow("Setting IDL authority for buffer", "bufferAccount", bufferAccount)
		args = append(args, bufferAccount)
	}
	e.Logger.Info(args)
	_, err := runCommand("anchor", args, programsPath)
	if err != nil {
		return fmt.Errorf("error setting idl authority: %w", err)
	}
	return nil
}

// get IDL address for a program
func getIDLAddress(e deployment.Environment, programID solana.PublicKey) (solana.PublicKey, error) {
	base, _, _ := solana.FindProgramAddress([][]byte{}, programID)
	idlAddress, _ := solana.CreateWithSeed(base, "anchor:idl", programID)
	e.Logger.Infof("IDL Address:  %s", idlAddress.String())
	return idlAddress, nil
}

// parse IDL buffer from `anchor idl write-buffer` output
func parseIdlBuffer(output string) (string, error) {
	const prefix = "Idl buffer created: "
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix)), nil
		}
	}
	return "", errors.New("failed to find IDL buffer in output")
}

// write IDL buffer for a program
func writeBuffer(e deployment.Environment, programsPath, programID, programName string) (solana.PublicKey, error) {
	idlFile, err := getIDL(e, programsPath, programID, programName)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("error getting IDL: %w", err)
	}
	e.Logger.Infow("Writing IDL buffer", "programID", programID)
	args := []string{"idl", "write-buffer", "--filepath", idlFile, programID}
	e.Logger.Info(args)
	output, err := runCommand("anchor", args, programsPath)
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

// generate set buffer ix using solana-go sdk
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

// generate upgrade IDL ix for a program via timelock
func upgradeIDLIx(e deployment.Environment, programsPath, programID, programName string, c IDLConfig) (*mcmsTypes.Transaction, error) {
	timelockSignerPDA, err := FetchTimelockSigner(e, c.ChainSelector)
	if err != nil {
		return nil, fmt.Errorf("error loading timelockSignerPDA: %w", err)
	}
	buffer, err := writeBuffer(e, programsPath, programID, programName)
	if err != nil {
		return nil, fmt.Errorf("error writing buffer: %w", err)
	}
	authority := e.SolChains[c.ChainSelector].DeployerKey.PublicKey()
	if c.MCMS != nil {
		authority = timelockSignerPDA
		err = setIdlAuthority(e, timelockSignerPDA.String(), programsPath, programID, programName, buffer.String())
		if err != nil {
			return nil, fmt.Errorf("error setting buffer authority: %w", err)
		}
	}
	instruction, err := setBufferIx(e, solana.MustPublicKeyFromBase58(programID), buffer, authority)
	if err != nil {
		return nil, fmt.Errorf("error generating set buffer ix: %w", err)
	}
	if c.MCMS != nil {
		upgradeTx, err := BuildMCMSTxn(&instruction, programID, deployment.ContractType(programName))
		if err != nil {
			return nil, fmt.Errorf("failed to create upgrade transaction: %w", err)
		}
		return upgradeTx, nil
	}
	if err := e.SolChains[c.ChainSelector].Confirm([]solana.Instruction{&instruction}); err != nil {
		return nil, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return nil, nil
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
	if c.RMNRemote && chainState.RMNRemote.IsZero() {
		return fmt.Errorf("rmnRemote not deployed for chain %d, cannot upload idl", c.ChainSelector)
	}
	if c.BurnMintTokenPool && chainState.BurnMintTokenPool.IsZero() {
		return fmt.Errorf("burnMintTokenPool not deployed for chain %d, cannot upload idl", c.ChainSelector)
	}
	if c.LockReleaseTokenPool && chainState.LockReleaseTokenPool.IsZero() {
		return fmt.Errorf("lockReleaseTokenPool not deployed for chain %d, cannot upload idl", c.ChainSelector)
	}
	addresses, err := e.ExistingAddresses.AddressesForChain(c.ChainSelector) //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
	if err != nil {
		return fmt.Errorf("failed to get existing addresses: %w", err)
	}
	mcmState, err := commonstate.MaybeLoadMCMSWithTimelockChainStateSolana(e.SolChains[c.ChainSelector], addresses)
	if err != nil {
		return fmt.Errorf("failed to load MCMS with timelock chain state: %w", err)
	}
	if c.MCM && mcmState.McmProgram.IsZero() {
		return fmt.Errorf("mcm program not deployed for chain %d, cannot upload idl", c.ChainSelector)
	}
	if c.Timelock && mcmState.TimelockProgram.IsZero() {
		return fmt.Errorf("timelock program not deployed for chain %d, cannot upload idl", c.ChainSelector)
	}
	if c.AccessController && mcmState.AccessControllerProgram.IsZero() {
		return fmt.Errorf("access controller program not deployed for chain %d, cannot upload idl", c.ChainSelector)
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

	if err := repoSetup(e, chain, c.GitCommitSha); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error setting up anchor workspace: %w", err)
	}
	// start uploading
	if c.Router {
		err := idlInit(e, chain.ProgramsPath, chainState.Router.String(), deployment.RouterProgramName)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
	}
	if c.FeeQuoter {
		err := idlInit(e, chain.ProgramsPath, chainState.FeeQuoter.String(), deployment.FeeQuoterProgramName)
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.OffRamp {
		err := idlInit(e, chain.ProgramsPath, chainState.OffRamp.String(), deployment.OffRampProgramName)
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.RMNRemote {
		err := idlInit(e, chain.ProgramsPath, chainState.RMNRemote.String(), deployment.RMNRemoteProgramName)
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.BurnMintTokenPool {
		err := idlInit(e, chain.ProgramsPath, chainState.BurnMintTokenPool.String(), deployment.BurnMintTokenPoolProgramName)
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.LockReleaseTokenPool {
		err := idlInit(e, chain.ProgramsPath, chainState.LockReleaseTokenPool.String(), deployment.LockReleaseTokenPoolProgramName)
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	addresses, err := e.ExistingAddresses.AddressesForChain(c.ChainSelector) //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get existing addresses: %w", err)
	}
	mcmState, err := commonstate.MaybeLoadMCMSWithTimelockChainStateSolana(e.SolChains[c.ChainSelector], addresses)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load MCMS with timelock chain state: %w", err)
	}
	if c.MCM {
		err := idlInit(e, chain.ProgramsPath, mcmState.McmProgram.String(), deployment.McmProgramName)
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.Timelock {
		err := idlInit(e, chain.ProgramsPath, mcmState.TimelockProgram.String(), deployment.TimelockProgramName)
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.AccessController {
		err := idlInit(e, chain.ProgramsPath, mcmState.AccessControllerProgram.String(), deployment.AccessControllerProgramName)
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
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[c.ChainSelector]
	chain := e.SolChains[c.ChainSelector]

	timelockSignerPDA, err := FetchTimelockSigner(e, c.ChainSelector)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error loading timelockSignerPDA: %w", err)
	}

	// set idl authority
	if c.Router {
		err = setIdlAuthority(e, timelockSignerPDA.String(), chain.ProgramsPath, chainState.Router.String(), deployment.RouterProgramName, "")
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.FeeQuoter {
		err = setIdlAuthority(e, timelockSignerPDA.String(), chain.ProgramsPath, chainState.FeeQuoter.String(), deployment.FeeQuoterProgramName, "")
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.OffRamp {
		err = setIdlAuthority(e, timelockSignerPDA.String(), chain.ProgramsPath, chainState.OffRamp.String(), deployment.OffRampProgramName, "")
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.RMNRemote {
		err = setIdlAuthority(e, timelockSignerPDA.String(), chain.ProgramsPath, chainState.RMNRemote.String(), deployment.RMNRemoteProgramName, "")
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.BurnMintTokenPool {
		err = setIdlAuthority(e, timelockSignerPDA.String(), chain.ProgramsPath, chainState.BurnMintTokenPool.String(), deployment.BurnMintTokenPoolProgramName, "")
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.LockReleaseTokenPool {
		err = setIdlAuthority(e, timelockSignerPDA.String(), chain.ProgramsPath, chainState.LockReleaseTokenPool.String(), deployment.LockReleaseTokenPoolProgramName, "")
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}

	addresses, err := e.ExistingAddresses.AddressesForChain(chain.Selector) //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get existing addresses: %w", err)
	}
	mcmState, err := commonstate.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load MCMS with timelock chain state: %w", err)
	}

	if c.AccessController {
		err = setIdlAuthority(e, timelockSignerPDA.String(), chain.ProgramsPath, mcmState.AccessControllerProgram.String(), types.AccessControllerProgram.String(), "")
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.Timelock {
		err = setIdlAuthority(e, timelockSignerPDA.String(), chain.ProgramsPath, mcmState.TimelockProgram.String(), types.RBACTimelockProgram.String(), "")
		if err != nil {
			return deployment.ChangesetOutput{}, nil
		}
	}
	if c.MCM {
		err = setIdlAuthority(e, timelockSignerPDA.String(), chain.ProgramsPath, mcmState.McmProgram.String(), types.ManyChainMultisigProgram.String(), "")
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

	if err := repoSetup(e, chain, c.GitCommitSha); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error setting up anchor workspace: %w", err)
	}

	mcmsTxs := make([]mcmsTypes.Transaction, 0)
	if c.Router {
		upgradeTx, err := upgradeIDLIx(e, chain.ProgramsPath, chainState.Router.String(), deployment.RouterProgramName, c)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	if c.FeeQuoter {
		upgradeTx, err := upgradeIDLIx(e, chain.ProgramsPath, chainState.FeeQuoter.String(), deployment.FeeQuoterProgramName, c)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	if c.OffRamp {
		upgradeTx, err := upgradeIDLIx(e, chain.ProgramsPath, chainState.OffRamp.String(), deployment.OffRampProgramName, c)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	if c.RMNRemote {
		upgradeTx, err := upgradeIDLIx(e, chain.ProgramsPath, chainState.RMNRemote.String(), deployment.RMNRemoteProgramName, c)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	if c.BurnMintTokenPool {
		upgradeTx, err := upgradeIDLIx(e, chain.ProgramsPath, chainState.BurnMintTokenPool.String(), deployment.BurnMintTokenPoolProgramName, c)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	if c.LockReleaseTokenPool {
		upgradeTx, err := upgradeIDLIx(e, chain.ProgramsPath, chainState.LockReleaseTokenPool.String(), deployment.LockReleaseTokenPoolProgramName, c)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}

	addresses, err := e.ExistingAddresses.AddressesForChain(chain.Selector) //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get existing addresses: %w", err)
	}
	mcmState, err := commonstate.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load MCMS with timelock chain state: %w", err)
	}

	if c.AccessController {
		upgradeTx, err := upgradeIDLIx(e, chain.ProgramsPath, mcmState.AccessControllerProgram.String(), deployment.AccessControllerProgramName, c)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	if c.Timelock {
		upgradeTx, err := upgradeIDLIx(e, chain.ProgramsPath, mcmState.TimelockProgram.String(), deployment.TimelockProgramName, c)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	if c.MCM {
		upgradeTx, err := upgradeIDLIx(e, chain.ProgramsPath, mcmState.McmProgram.String(), deployment.McmProgramName, c)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}

	if len(mcmsTxs) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, c.ChainSelector, "proposal to upgrade CCIP contracts", c.MCMS.MinDelay, mcmsTxs)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}

		// do we need to batch this ?
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return deployment.ChangesetOutput{}, nil
}
