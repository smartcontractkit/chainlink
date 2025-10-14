package solana

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"

	"github.com/gagliardetto/solana-go"
	"github.com/pelletier/go-toml"
	chainsel "github.com/smartcontractkit/chain-selectors"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	cldfsolana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

// use this changeset to upload the IDL for a program
var _ cldf.ChangeSet[IDLConfig] = UploadIDL

// use this changeset to set the authority for the IDL of a program (timelock)
var _ cldf.ChangeSet[IDLConfig] = SetAuthorityIDL

// use this changeset to upgrade the IDL of a program via timelock
var _ cldf.ChangeSet[IDLConfig] = UpgradeIDL

// use this changeset to set the authority for the IDL of a program (timelock) via timelock
var _ cldf.ChangeSet[IDLConfig] = SetAuthorityIDLByMCMs

// changeset to upgrade idl for a program via timelock
var _ cldf.ChangeSet[IDLConfig] = UpgradeIDL

// use this changeset to close the IDL account for a program via timelock
var _ cldf.ChangeSet[IDLConfig] = CloseIDLs

type IDLConfig struct {
	ChainSelector                uint64
	SolanaContractVersion        string                        // Get the commit sha with VersionToShortCommitSHA[VersionSolanaV0_1_2] this will be used to download the correct artifacts (idls) -> best if same as what was used to deploy the programs
	Router                       bool                          // whether to upload the IDL for the router
	FeeQuoter                    bool                          // whether to upload the IDL for the fee quoter
	OffRamp                      bool                          // whether to upload the IDL for the off ramp
	RMNRemote                    bool                          // whether to upload the IDL for the rmn remote
	AccessController             bool                          // whether to upload the IDL for the access controller
	MCM                          bool                          // whether to upload the IDL for the mcm
	Timelock                     bool                          // whether to upload the IDL for the timelock
	BurnMintTokenPoolMetadata    []string                      // whether to upload the IDL for the token pool (keyed my client identifier (metadata))
	LockReleaseTokenPoolMetadata []string                      // metadata for the lock release token pool (keyed my client identifier (metadata))
	MCMS                         *proposalutils.TimelockConfig // timelock config for mcms
	CCTPTokenPool                bool
}

func (c IDLConfig) Validate(e cldf.Environment) error {
	if err := cldf.IsValidChainSelector(c.ChainSelector); err != nil {
		return fmt.Errorf("invalid chain selector: %d - %w", c.ChainSelector, err)
	}
	family, _ := chainsel.GetSelectorFamily(c.ChainSelector)
	if family != chainsel.FamilySolana {
		return fmt.Errorf("chain %d is not a solana chain", c.ChainSelector)
	}
	existingState, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load existing onchain state: %w", err)
	}
	if _, exists := existingState.SupportedChains()[c.ChainSelector]; !exists {
		return fmt.Errorf("chain %d not supported", c.ChainSelector)
	}
	chainState := existingState.SolChains[c.ChainSelector]
	chain := e.BlockChains.SolanaChains()[c.ChainSelector]
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
	for _, bnmMetadata := range c.BurnMintTokenPoolMetadata {
		bnmTokenPool := chainState.GetActiveTokenPool(shared.BurnMintTokenPool, bnmMetadata)
		if bnmTokenPool.IsZero() {
			return fmt.Errorf("burnMintTokenPool not deployed for chain %d, cannot upload idl", c.ChainSelector)
		}
	}
	for _, lrMetadata := range c.LockReleaseTokenPoolMetadata {
		lrTokenPool := chainState.GetActiveTokenPool(shared.LockReleaseTokenPool, lrMetadata)
		if lrTokenPool.IsZero() {
			return fmt.Errorf("lockReleaseTokenPool not deployed for chain %d, cannot upload idl", c.ChainSelector)
		}
	}
	addresses, err := e.ExistingAddresses.AddressesForChain(c.ChainSelector)
	if err != nil {
		return fmt.Errorf("failed to get existing addresses: %w", err)
	}
	mcmState, err := commonstate.MaybeLoadMCMSWithTimelockChainStateSolana(e.BlockChains.SolanaChains()[c.ChainSelector], addresses)
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
	commitSha := VersionToShortCommitSHA[c.SolanaContractVersion]

	return RepoSetup(e, chain, commitSha)
}

// ANCHOR CLI OPERATIONS

// changeset to set idl authority for a program to timelock
func SetAuthorityIDL(e cldf.Environment, c IDLConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error validating idl config: %w", err)
	}
	state, _ := stateview.LoadOnchainState(e)
	chainState := state.SolChains[c.ChainSelector]
	chain := e.BlockChains.SolanaChains()[c.ChainSelector]

	timelockSignerPDA, err := FetchTimelockSigner(e, c.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error loading timelockSignerPDA: %w", err)
	}

	// set idl authority
	if c.Router {
		err = SetAuthorityIDLByCLI(e, timelockSignerPDA.String(), chain.ProgramsPath, chainState.Router.String(), deployment.RouterProgramName, "")
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}
	if c.FeeQuoter {
		err = SetAuthorityIDLByCLI(e, timelockSignerPDA.String(), chain.ProgramsPath, chainState.FeeQuoter.String(), deployment.FeeQuoterProgramName, "")
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}
	if c.OffRamp {
		err = SetAuthorityIDLByCLI(e, timelockSignerPDA.String(), chain.ProgramsPath, chainState.OffRamp.String(), deployment.OffRampProgramName, "")
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}
	if c.RMNRemote {
		err = SetAuthorityIDLByCLI(e, timelockSignerPDA.String(), chain.ProgramsPath, chainState.RMNRemote.String(), deployment.RMNRemoteProgramName, "")
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}
	for _, bnmMetadata := range c.BurnMintTokenPoolMetadata {
		tokenPool := chainState.GetActiveTokenPool(shared.BurnMintTokenPool, bnmMetadata)
		err = SetAuthorityIDLByCLI(e, timelockSignerPDA.String(), chain.ProgramsPath, tokenPool.String(), deployment.BurnMintTokenPoolProgramName, "")
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}
	for _, lrMetadata := range c.LockReleaseTokenPoolMetadata {
		tokenPool := chainState.GetActiveTokenPool(shared.LockReleaseTokenPool, lrMetadata)
		err = SetAuthorityIDLByCLI(e, timelockSignerPDA.String(), chain.ProgramsPath, tokenPool.String(), deployment.LockReleaseTokenPoolProgramName, "")
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}
	if c.CCTPTokenPool {
		err = SetAuthorityIDLByCLI(e, timelockSignerPDA.String(), chain.ProgramsPath, chainState.CCTPTokenPool.String(), deployment.CCTPTokenPoolProgramName, "")
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}

	addresses, err := e.ExistingAddresses.AddressesForChain(chain.Selector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get existing addresses: %w", err)
	}
	mcmState, err := commonstate.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load MCMS with timelock chain state: %w", err)
	}

	if c.AccessController {
		err = SetAuthorityIDLByCLI(e, timelockSignerPDA.String(), chain.ProgramsPath, mcmState.AccessControllerProgram.String(), types.AccessControllerProgram.String(), "")
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}
	if c.Timelock {
		err = SetAuthorityIDLByCLI(e, timelockSignerPDA.String(), chain.ProgramsPath, mcmState.TimelockProgram.String(), types.RBACTimelockProgram.String(), "")
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}
	if c.MCM {
		err = SetAuthorityIDLByCLI(e, timelockSignerPDA.String(), chain.ProgramsPath, mcmState.McmProgram.String(), types.ManyChainMultisigProgram.String(), "")
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}

	return cldf.ChangesetOutput{}, nil
}

// parse anchor version from running anchor --version
func ParseAnchorVersion(output string) (string, error) {
	const prefix = "anchor-cli "
	if after, ok := strings.CutPrefix(output, prefix); ok {
		return strings.TrimSpace(after), nil
	}
	return "", fmt.Errorf("unexpected version output: %q", output)
}

// create Anchor.toml file to simulate anchor workspace
func WriteAnchorToml(e cldf.Environment, filename, anchorVersion, cluster, wallet string) error {
	e.Logger.Debugw("Writing Anchor.toml", "filename", filename, "anchorVersion", anchorVersion, "cluster", cluster, "wallet", wallet)
	config := map[string]any{
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
func RepoSetup(e cldf.Environment, chain cldfsolana.Chain, gitCommitSha string) error {
	e.Logger.Debug("Downloading Solana CCIP program artifacts...")
	err := memory.DownloadSolanaCCIPProgramArtifacts(e.GetContext(), chain.ProgramsPath, e.Logger, gitCommitSha)
	if err != nil {
		return fmt.Errorf("error downloading solana ccip program artifacts: %w", err)
	}

	// get anchor version
	output, err := RunCommand("anchor", []string{"--version"}, ".")
	if err != nil {
		return errors.New("anchor-cli not installed in path")
	}
	e.Logger.Debugw("Anchor version command output", "output", output)
	anchorVersion, err := ParseAnchorVersion(output)
	if err != nil {
		return fmt.Errorf("error parsing anchor version: %w", err)
	}
	// create Anchor.toml
	// this creates anchor workspace with cluster and wallet configured
	if err := WriteAnchorToml(e, filepath.Join(chain.ProgramsPath, "Anchor.toml"), anchorVersion, chain.URL, chain.KeypairPath); err != nil {
		return fmt.Errorf("error writing Anchor.toml: %w", err)
	}

	return nil
}

// update IDL with program ID
func updateIDL(e cldf.Environment, idlFile string, programID string) error {
	e.Logger.Debug("Reading IDL")
	idlBytes, err := os.ReadFile(idlFile)
	if err != nil {
		return fmt.Errorf("failed to read IDL: %w", err)
	}
	e.Logger.Debug("Parsing IDL")
	var idl map[string]any
	if err := json.Unmarshal(idlBytes, &idl); err != nil {
		return fmt.Errorf("failed to parse legacy IDL: %w", err)
	}
	e.Logger.Debugw("Updating IDL with programID", "programID", programID)
	idl["metadata"] = map[string]any{
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
func getIDL(e cldf.Environment, programsPath, programID string, programName string) (string, error) {
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
func IdlInit(e cldf.Environment, programsPath, programID, programName string) error {
	idlFile, err := getIDL(e, programsPath, programID, programName)
	if err != nil {
		return fmt.Errorf("error getting IDL: %w", err)
	}
	e.Logger.Infow("Uploading IDL", "programName", programName)
	args := []string{"idl", "init", "--filepath", idlFile, programID}
	e.Logger.Info(args)
	output, err := RunCommand("anchor", args, programsPath)
	e.Logger.Debugw("IDL init output", "output", output)
	if err != nil {
		e.Logger.Debugw("IDL init error", "error", err)
		return fmt.Errorf("error uploading idl: %w", err)
	}
	e.Logger.Infow("IDL uploaded", "programName", programName)
	return nil
}

// get IDL address for a program
func getIDLAddress(e cldf.Environment, programID solana.PublicKey) (solana.PublicKey, error) {
	base, _, _ := solana.FindProgramAddress([][]byte{}, programID)
	idlAddress, _ := solana.CreateWithSeed(base, "anchor:idl", programID)
	e.Logger.Infof("IDL Address:  %s", idlAddress.String())
	return idlAddress, nil
}

// parse IDL buffer from `anchor idl write-buffer` output
func parseIdlBuffer(output string) (string, error) {
	const prefix = "Idl buffer created: "
	for line := range strings.SplitSeq(output, "\n") {
		if after, ok := strings.CutPrefix(line, prefix); ok {
			return strings.TrimSpace(after), nil
		}
	}
	return "", errors.New("failed to find IDL buffer in output")
}

// write IDL buffer for a program
func writeBuffer(e cldf.Environment, programsPath, programID, programName string) (solana.PublicKey, error) {
	idlFile, err := getIDL(e, programsPath, programID, programName)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("error getting IDL: %w", err)
	}
	e.Logger.Infow("Writing IDL buffer", "programID", programID)
	args := []string{"idl", "write-buffer", "--filepath", idlFile, programID}
	e.Logger.Info(args)
	output, err := RunCommand("anchor", args, programsPath)
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

func SetAuthorityIDLByCLI(e cldf.Environment, newAuthority, programsPath, programID, programName, bufferAccount string) error {
	e.Logger.Infow("Setting IDL authority", "programName", programName, "newAuthority", newAuthority)
	args := []string{"idl", "set-authority", "-n", newAuthority, "-p", programID}
	if bufferAccount != "" {
		e.Logger.Infow("Setting IDL authority for buffer", "bufferAccount", bufferAccount)
		args = append(args, bufferAccount)
	}
	e.Logger.Info(args)
	_, err := RunCommand("anchor", args, programsPath)
	if err != nil {
		return fmt.Errorf("error setting idl authority: %w", err)
	}
	return nil
}

// SOLANA INSTRUCTIONS

// Discriminator to invoke IDL operations
const IdlIxTag uint64 = 0x0a69e9a778bcf440

const DefaultIDLMaxSize = 512 * 1024 // 512 KB

// Number ids of the operations: copied from https://github.com/solana-foundation/anchor/blob/v0.29.0/lang/src/idl.rs#L36
const (
	IdlInstructionCreate       int = iota // One time initializer for creating the program's idl account.
	IdlInstructionCreateBuffer            // Creates a new IDL account buffer. Can be called several times.
	IdlInstructionWrite                   // Appends the given data to the end of the idl account buffer.
	IdlInstructionSetBuffer               // Sets a new data buffer for the IdlAccount.
	IdlInstructionSetAuthority            // Sets a new authority on the IdlAccount.
	IdlInstructionClose                   // Closes the IDL pda Account
	IdlInstructionResize                  // Increases account size for accounts that need over 10kb.
)

// changeset to set idl authority for a program to timelock
func SetAuthorityIDLByMCMs(e cldf.Environment, c IDLConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error validating idl config: %w", err)
	}
	state, _ := stateview.LoadOnchainState(e)
	chainState := state.SolChains[c.ChainSelector]
	chain := e.BlockChains.SolanaChains()[c.ChainSelector]

	mcmsTxs := make([]mcmsTypes.Transaction, 0)
	newAuthority := e.BlockChains.SolanaChains()[c.ChainSelector].DeployerKey.PublicKey()

	programs, err := getAffectedPrograms(e, c, chainState, chain)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	for programID, programName := range programs {
		setAuthorityTx, err := setAuthorityIDLIx(e, programID, programName, newAuthority, c)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		if setAuthorityTx != nil {
			mcmsTxs = append(mcmsTxs, *setAuthorityTx)
		}
	}

	return generateProposalIfMCMS(e, c.ChainSelector, c.MCMS, mcmsTxs)
}

func UploadIDL(e cldf.Environment, c IDLConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error validating idl config: %w", err)
	}
	chain := e.BlockChains.SolanaChains()[c.ChainSelector]
	state, _ := stateview.LoadOnchainState(e)
	chainState := state.SolChains[c.ChainSelector]

	mcmsTxs := make([]mcmsTypes.Transaction, 0)

	programs, err := getAffectedPrograms(e, c, chainState, chain)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	for programID, programName := range programs {
		upgradeTx, err := IdlInitIx(e, chain.ProgramsPath, programID.String(), programName, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating idl init tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}

	return generateProposalIfMCMS(e, c.ChainSelector, c.MCMS, mcmsTxs)
}

// changeset to upgrade idl for a program via timelock
// write buffer using anchor cli
// set buffer authority to timelock using anchor cli
// generate set buffer ix using solana-go sdk
// build mcms txn to upgrade idl
func UpgradeIDL(e cldf.Environment, c IDLConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error validating idl config: %w", err)
	}
	chain := e.BlockChains.SolanaChains()[c.ChainSelector]
	state, _ := stateview.LoadOnchainState(e)
	chainState := state.SolChains[c.ChainSelector]

	mcmsTxs := make([]mcmsTypes.Transaction, 0)

	programs, err := getAffectedPrograms(e, c, chainState, chain)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	for programID, programName := range programs {
		upgradeTx, err := upgradeIDLIx(e, chain.ProgramsPath, programID.String(), programName, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}

	return generateProposalIfMCMS(e, c.ChainSelector, c.MCMS, mcmsTxs)
}

// changeset to close idl account for a program - this is needed when the idl increased so much in size that it no longer fits in the account
// and the idl account can not be resized when it already contains data, so you need to close the account and reinitialize it
func CloseIDLs(e cldf.Environment, c IDLConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error validating idl config: %w", err)
	}
	chain := e.BlockChains.SolanaChains()[c.ChainSelector]
	state, _ := stateview.LoadOnchainState(e)
	chainState := state.SolChains[c.ChainSelector]

	mcmsTxs := make([]mcmsTypes.Transaction, 0)

	programs, err := getAffectedPrograms(e, c, chainState, chain)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	for programID, programName := range programs {
		closeIdlIx, err := closeIdlInstruction(e, programID, programName, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating close IDL tx: %w", err)
		}
		if closeIdlIx != nil {
			mcmsTxs = append(mcmsTxs, *closeIdlIx)
		}
	}

	return generateProposalIfMCMS(e, c.ChainSelector, c.MCMS, mcmsTxs)
}

func getAffectedPrograms(e cldf.Environment, c IDLConfig, chainState solanastateview.CCIPChainState, chain cldfsolana.Chain) (map[solana.PublicKey]string, error) {
	programs := make(map[solana.PublicKey]string)
	if c.Router {
		programs[chainState.Router] = deployment.RouterProgramName
	}
	if c.FeeQuoter {
		programs[chainState.FeeQuoter] = deployment.FeeQuoterProgramName
	}
	if c.OffRamp {
		programs[chainState.OffRamp] = deployment.OffRampProgramName
	}
	if c.RMNRemote {
		programs[chainState.RMNRemote] = deployment.RMNRemoteProgramName
	}
	for _, bnmMetadata := range c.BurnMintTokenPoolMetadata {
		tokenPool := chainState.GetActiveTokenPool(shared.BurnMintTokenPool, bnmMetadata)
		programs[tokenPool] = deployment.BurnMintTokenPoolProgramName
	}
	for _, lrMetadata := range c.LockReleaseTokenPoolMetadata {
		tokenPool := chainState.GetActiveTokenPool(shared.LockReleaseTokenPool, lrMetadata)
		programs[tokenPool] = deployment.LockReleaseTokenPoolProgramName
	}
	if c.CCTPTokenPool {
		tokenPool := chainState.GetActiveTokenPool(shared.CCTPTokenPool, shared.CLLMetadata)
		programs[tokenPool] = deployment.CCTPTokenPoolProgramName
	}
	addresses, err := e.ExistingAddresses.AddressesForChain(chain.Selector)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing addresses: %w", err)
	}
	mcmState, err := commonstate.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to load MCMS with timelock chain state: %w", err)
	}
	if c.AccessController {
		programs[mcmState.AccessControllerProgram] = deployment.AccessControllerProgramName
	}
	if c.Timelock {
		programs[mcmState.TimelockProgram] = deployment.TimelockProgramName
	}
	if c.MCM {
		programs[mcmState.McmProgram] = deployment.McmProgramName
	}
	return programs, nil
}

// Build instruction to interact with Anchor IDL using the list of ids above for each message
func buildIdlInstruction(programID solana.PublicKey, accountsForIx solana.AccountMetaSlice, idlInstruction int, params []byte) (solana.GenericInstruction, error) {
	data := binary.LittleEndian.AppendUint64([]byte{}, IdlIxTag) // 8-byte Extend instruction identifier
	data = append(data, byte(idlInstruction))                    // Append the numeric ID of the operation
	data = append(data, params...)                               // Append any additional parameters

	instruction := solana.NewInstruction(
		programID,
		accountsForIx,
		data,
	)
	return *instruction, nil
}

func calculateAuthority(e cldf.Environment, c IDLConfig) (solana.PublicKey, error) {
	timelockSignerPDA, err := FetchTimelockSigner(e, c.ChainSelector)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("error loading timelockSignerPDA: %w", err)
	}
	authority := e.BlockChains.SolanaChains()[c.ChainSelector].DeployerKey.PublicKey()
	if c.MCMS != nil {
		authority = timelockSignerPDA
	}
	return authority, err
}

func getTxIfMCMSExecuteIfNot(e cldf.Environment, programID string, programName string, c IDLConfig, instruction solana.GenericInstruction) (*mcmsTypes.Transaction, error) {
	if c.MCMS != nil {
		upgradeTx, err := BuildMCMSTxn(&instruction, programID, cldf.ContractType(programName))
		if err != nil {
			return nil, fmt.Errorf("failed to create close IDL transaction: %w", err)
		}
		return upgradeTx, nil
	}
	if err := e.BlockChains.SolanaChains()[c.ChainSelector].Confirm([]solana.Instruction{&instruction}); err != nil {
		return nil, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return nil, nil
}

// generate set buffer ix using solana-go sdk
func setBufferIdlInstruction(e cldf.Environment, programID, buffer, authority solana.PublicKey) (solana.GenericInstruction, error) {
	accounts, instruction, err := getAccountsFoSetBufferIdlInstruction(e, programID, buffer, authority)
	if err != nil {
		return instruction, err
	}
	return buildIdlInstruction(programID, accounts, IdlInstructionSetBuffer, []byte{})
}

func createIdlInstruction(e cldf.Environment, programID, authority solana.PublicKey) (solana.GenericInstruction, error) {
	accounts, instruction, err := getAccountsFoCreateIdlInstruction(e, programID, authority)
	if err != nil {
		return instruction, err
	}

	params := idlCreateParams(uint64(DefaultIDLMaxSize))

	return buildIdlInstruction(programID, accounts, IdlInstructionCreate, params)
}

func idlCreateParams(dataLen uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, dataLen)
	return b
}

// Following the Anchor 0.29.0 Implementation: https://github.com/solana-foundation/anchor/blob/2a050757609a3c59bd77084a259f5ea64fcebfa6/lang/syn/src/codegen/program/idl.rs#L38
func getAccountsFoCreateIdlInstruction(
	e cldf.Environment,
	programID solana.PublicKey,
	authority solana.PublicKey,
) (solana.AccountMetaSlice, solana.GenericInstruction, error) {
	// Derive the IDL account PDA (the "to" account in IdlCreateAccounts)
	idlAddress, err := getIDLAddress(e, programID)
	if err != nil {
		return nil, solana.GenericInstruction{}, fmt.Errorf("error getting idl address for %s: %w", programID.String(), err)
	}

	// Derive the base PDA with empty seeds ([]), used by the IDL program as signer
	base, _, err := solana.FindProgramAddress([][]byte{}, programID)
	if err != nil {
		return nil, solana.GenericInstruction{}, fmt.Errorf("error deriving base PDA for %s: %w", programID.String(), err)
	}

	// Build the accounts in the exact order expected by Anchor 0.29's IdlCreateAccounts:
	//   1. from            -> signer + writable (payer)
	//   2. to              -> writable (IDL account being created)
	//   3. base            -> readonly (PDA with seeds=[])
	//   4. system_program  -> readonly
	//   5. program         -> readonly (target program)
	accounts := solana.AccountMetaSlice{
		solana.Meta(authority).SIGNER().WRITE(), // from (payer)
		solana.Meta(idlAddress).WRITE(),         // to (IDL account)
		solana.Meta(base),                       // base PDA (readonly)
		solana.Meta(solana.SystemProgramID),     // system_program
		solana.Meta(programID),                  // program (target)
	}

	return accounts, solana.GenericInstruction{}, nil
}

func getAccountsFoSetBufferIdlInstruction(e cldf.Environment, programID solana.PublicKey, buffer solana.PublicKey, authority solana.PublicKey) (solana.AccountMetaSlice, solana.GenericInstruction, error) {
	idlAddress, err := getIDLAddress(e, programID)
	if err != nil {
		return nil, solana.GenericInstruction{}, fmt.Errorf("error getting idl address for %s: %w", programID.String(), err)
	}
	accounts := solana.AccountMetaSlice{
		solana.Meta(buffer).WRITE(),
		solana.Meta(idlAddress).WRITE(),
		solana.Meta(authority).SIGNER(),
	}
	return accounts, solana.GenericInstruction{}, nil
}

func IdlInitIx(e cldf.Environment, programsPath, programID, programName string, c IDLConfig) (*mcmsTypes.Transaction, error) {
	timelockSignerPDA, err := FetchTimelockSigner(e, c.ChainSelector)
	if err != nil {
		return nil, fmt.Errorf("error loading timelockSignerPDA: %w", err)
	}
	buffer, err := writeBuffer(e, programsPath, programID, programName)
	if err != nil {
		return nil, fmt.Errorf("error writing buffer: %w", err)
	}
	authority := e.BlockChains.SolanaChains()[c.ChainSelector].DeployerKey.PublicKey()
	if c.MCMS != nil {
		authority = timelockSignerPDA
		err = SetAuthorityIDLByCLI(e, timelockSignerPDA.String(), programsPath, programID, programName, buffer.String())
		if err != nil {
			return nil, fmt.Errorf("error setting buffer authority: %w", err)
		}
	}
	instruction, err := createIdlInstruction(e, solana.MustPublicKeyFromBase58(programID), authority)
	if err != nil {
		return nil, fmt.Errorf("error generating set buffer ix: %w", err)
	}
	if c.MCMS != nil {
		createIdlIx, err := BuildMCMSTxn(&instruction, programID, cldf.ContractType(programName))
		if err != nil {
			return nil, fmt.Errorf("failed to create upgrade transaction: %w", err)
		}
		return createIdlIx, nil
	}
	if err := e.BlockChains.SolanaChains()[c.ChainSelector].Confirm([]solana.Instruction{&instruction}); err != nil {
		return nil, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return nil, nil
}

// generate upgrade IDL ix for a program via timelock
func upgradeIDLIx(e cldf.Environment, programsPath, programID, programName string, c IDLConfig) (*mcmsTypes.Transaction, error) {
	timelockSignerPDA, err := FetchTimelockSigner(e, c.ChainSelector)
	if err != nil {
		return nil, fmt.Errorf("error loading timelockSignerPDA: %w", err)
	}
	buffer, err := writeBuffer(e, programsPath, programID, programName)
	if err != nil {
		return nil, fmt.Errorf("error writing buffer: %w", err)
	}
	authority := e.BlockChains.SolanaChains()[c.ChainSelector].DeployerKey.PublicKey()
	if c.MCMS != nil {
		authority = timelockSignerPDA
		err = SetAuthorityIDLByCLI(e, timelockSignerPDA.String(), programsPath, programID, programName, buffer.String())
		if err != nil {
			return nil, fmt.Errorf("error setting buffer authority: %w", err)
		}
	}
	instruction, err := setBufferIdlInstruction(e, solana.MustPublicKeyFromBase58(programID), buffer, authority)
	if err != nil {
		return nil, fmt.Errorf("error generating set buffer ix: %w", err)
	}
	if c.MCMS != nil {
		upgradeTx, err := BuildMCMSTxn(&instruction, programID, cldf.ContractType(programName))
		if err != nil {
			return nil, fmt.Errorf("failed to create upgrade transaction: %w", err)
		}
		return upgradeTx, nil
	}
	if err := e.BlockChains.SolanaChains()[c.ChainSelector].Confirm([]solana.Instruction{&instruction}); err != nil {
		return nil, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return nil, nil
}

// generate close IDL ix for a program (via timelock or via prod deployer)
func closeIdlInstruction(e cldf.Environment, programID solana.PublicKey, programName string, c IDLConfig) (*mcmsTypes.Transaction, error) {
	authority, err := calculateAuthority(e, c)
	if err != nil {
		return nil, err
	}
	// The spill Address should always be the deployer key, even if using timelock to close the account
	// since the deployer key is the one that pays for the account rent
	// and is therefore the one that should receive the remaining lamports
	// in the account when it is closed.
	spillAddress := e.BlockChains.SolanaChains()[c.ChainSelector].DeployerKey.PublicKey()
	accounts, err := getAccountsForCloseIdlInstruction(e, programID, authority, spillAddress)
	if err != nil {
		return nil, fmt.Errorf("error getting accounts for close idl instruction %s: %w", programID.String(), err)
	}
	instruction, err := buildIdlInstruction(programID, accounts, IdlInstructionClose, []byte{})
	if err != nil {
		return nil, fmt.Errorf("error closing IDL account ix: %w", err)
	}
	return getTxIfMCMSExecuteIfNot(e, programID.String(), programName, c, instruction)
}

func getAccountsForCloseIdlInstruction(e cldf.Environment, programID solana.PublicKey, authority solana.PublicKey, spillAddress solana.PublicKey) (solana.AccountMetaSlice, error) {
	idlAddress, err := getIDLAddress(e, programID)
	accounts := solana.AccountMetaSlice{
		solana.Meta(idlAddress).WRITE(),
		solana.Meta(authority).SIGNER(),
		solana.Meta(spillAddress).WRITE(), // SOL destination to close funds
	}
	return accounts, err
}

// Set Authority IDL
func setAuthorityIDLIx(e cldf.Environment, programID solana.PublicKey, programName string, newAuthority solana.PublicKey, c IDLConfig) (*mcmsTypes.Transaction, error) {
	authority, err := calculateAuthority(e, c)
	if err != nil {
		return nil, err
	}
	accounts, err := getAccountsForSetAuthorityIdlInstruction(e, programID, authority)
	if err != nil {
		return nil, err
	}
	instruction, err := buildIdlInstruction(programID, accounts, IdlInstructionSetAuthority, newAuthority.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error setting authority IDL ix: %w", err)
	}
	return getTxIfMCMSExecuteIfNot(e, programID.String(), programName, c, instruction)
}

func getAccountsForSetAuthorityIdlInstruction(e cldf.Environment, programID solana.PublicKey, authority solana.PublicKey) (solana.AccountMetaSlice, error) {
	idlAddress, err := getIDLAddress(e, programID)
	if err != nil {
		return nil, fmt.Errorf("error getting idl address for %s: %w", programID.String(), err)
	}
	accounts := solana.AccountMetaSlice{
		solana.Meta(idlAddress).WRITE(),
		solana.Meta(authority).SIGNER(),
	}
	return accounts, nil
}
