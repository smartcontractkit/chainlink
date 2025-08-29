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

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

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
var _ cldf.ChangeSet[IDLConfig] = SetAuthorityIDLByDeployerKey

// use this changeset to upgrade the IDL of a program via timelock
var _ cldf.ChangeSet[IDLConfig] = UpgradeIDL

// Discriminator to invoke IDL operations
const IdlIxTag uint64 = 0x0a69e9a778bcf440

// Number ids of the operations:
const (
	IdlInstructionCreate       int = iota // One time initializer for creating the program's idl account.
	IdlInstructionCreateBuffer            // Creates a new IDL account buffer. Can be called several times.
	IdlInstructionWrite                   // Appends the given data to the end of the idl account buffer.
	IdlInstructionSetBuffer               // Sets a new data buffer for the IdlAccount.
	IdlInstructionSetAuthority            // Sets a new authority on the IdlAccount.
	IdlInstructionClose                   // Closes the IDL pda Account
	IdlInstructionResize                  // Increases account size for accounts that need over 10kb.
)

// IDL
type IDLConfig struct {
	ChainSelector                uint64
	GitCommitSha                 string                        // this will be used to download the correct artifacts (idls) -> best if same as what was used to deploy the programs
	SpillAddress                 solana.PublicKey              // used when closing the IDL account
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

// parse anchor version from running anchor --version
func parseAnchorVersion(output string) (string, error) {
	const prefix = "anchor-cli "
	if strings.HasPrefix(output, prefix) {
		return strings.TrimSpace(strings.TrimPrefix(output, prefix)), nil
	}
	return "", fmt.Errorf("unexpected version output: %q", output)
}

// create Anchor.toml file to simulate anchor workspace
func writeAnchorToml(e cldf.Environment, filename, anchorVersion, cluster, wallet string) error {
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
func repoSetup(e cldf.Environment, chain cldf_solana.Chain, gitCommitSha string) error {
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
func updateIDL(e cldf.Environment, idlFile string, programID string) error {
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
func idlInit(e cldf.Environment, programsPath, programID, programName string) error {
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

func setAuthorityIdlInstruction(e cldf.Environment, programID, authority, newAuthority solana.PublicKey) (solana.GenericInstruction, error) {
	accounts, instruction, err2 := getAccountsForSetAuthorityIdlInstruction(e, programID, authority)
	if err2 != nil {
		return instruction, err2
	}
	return buildIdlInstruction(programID, accounts, IdlInstructionSetAuthority, newAuthority.Bytes())
}

func getAccountsForSetAuthorityIdlInstruction(e cldf.Environment, programID solana.PublicKey, authority solana.PublicKey) (solana.AccountMetaSlice, solana.GenericInstruction, error) {
	idlAddress, err := getIDLAddress(e, programID)
	if err != nil {
		return nil, solana.GenericInstruction{}, fmt.Errorf("error getting idl address for %s: %w", programID.String(), err)
	}
	accounts := solana.AccountMetaSlice{
		solana.NewAccountMeta(idlAddress, true, false), // IDL account
		solana.NewAccountMeta(authority, false, true),
	}
	return accounts, solana.GenericInstruction{}, nil
}

// generate close IDL PDA ix for a program via timelock
func setAuthorityIDLIx(e cldf.Environment, programID, programName string, newAuthority solana.PublicKey, c IDLConfig) (*mcmsTypes.Transaction, error) {
	timelockSignerPDA, err := FetchTimelockSigner(e, c.ChainSelector)
	if err != nil {
		return nil, fmt.Errorf("error loading timelockSignerPDA: %w", err)
	}
	authority := e.BlockChains.SolanaChains()[c.ChainSelector].DeployerKey.PublicKey()
	if c.MCMS != nil {
		authority = timelockSignerPDA
	}
	instruction, err := setAuthorityIdlInstruction(e, solana.MustPublicKeyFromBase58(programID), authority, newAuthority)
	if err != nil {
		return nil, fmt.Errorf("error setting authority IDL ix: %w", err)
	}
	if c.MCMS != nil {
		upgradeTx, err := BuildMCMSTxn(&instruction, programID, cldf.ContractType(programName))
		if err != nil {
			return nil, fmt.Errorf("failed to set authority IDL transaction: %w", err)
		}
		return upgradeTx, nil
	}
	if err := e.BlockChains.SolanaChains()[c.ChainSelector].Confirm([]solana.Instruction{&instruction}); err != nil {
		return nil, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return nil, nil
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
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix)), nil
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
func setBufferIdlInstruction(e cldf.Environment, programID, buffer, authority solana.PublicKey) (solana.GenericInstruction, error) {
	accounts, instruction, err2 := getAccountsForCreateBufferIdlInstruction(e, programID, buffer, authority)
	if err2 != nil {
		return instruction, err2
	}
	return buildIdlInstruction(programID, accounts, IdlInstructionSetBuffer, []byte{})
}

func getAccountsForCreateBufferIdlInstruction(e cldf.Environment, programID solana.PublicKey, buffer solana.PublicKey, authority solana.PublicKey) (solana.AccountMetaSlice, solana.GenericInstruction, error) {
	idlAddress, err := getIDLAddress(e, programID)
	if err != nil {
		return nil, solana.GenericInstruction{}, fmt.Errorf("error getting idl address for %s: %w", programID.String(), err)
	}
	accounts := solana.AccountMetaSlice{
		solana.NewAccountMeta(buffer, true, false),
		solana.NewAccountMeta(idlAddress, true, false),
		solana.NewAccountMeta(authority, false, true),
	}
	return accounts, solana.GenericInstruction{}, nil
}

func buildIdlInstruction(programID solana.PublicKey, accountsForIx solana.AccountMetaSlice, idlInstruction int, params []byte) (solana.GenericInstruction, error) {
	data := binary.LittleEndian.AppendUint64([]byte{}, IdlIxTag) // 4-byte Extend instruction identifier
	data = append(data, byte(idlInstruction))                    // Append the numeric ID of the operation
	data = append(data, params...)                               // Append any additional parameters

	instruction := solana.NewInstruction(
		programID,
		accountsForIx,
		data,
	)
	return *instruction, nil
}

func setAuthorityIDLByCLI(e cldf.Environment, newAuthority, programsPath, programID, programName, bufferAccount string) error {
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
		err = setAuthorityIDLByCLI(e, timelockSignerPDA.String(), programsPath, programID, programName, buffer.String())
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

// generate set buffer ix using solana-go sdk
func closeIdlInstruction(e cldf.Environment, programID, authority, spillAddress solana.PublicKey) (solana.GenericInstruction, error) {
	accounts, err := getAccountsForCloseIdlInstruction(e, programID, authority, spillAddress)
	if err != nil {
		return solana.GenericInstruction{}, fmt.Errorf("error getting idl address for %s: %w", programID.String(), err)
	}
	return buildIdlInstruction(programID, accounts, IdlInstructionClose, []byte{})
}

func getAccountsForCloseIdlInstruction(e cldf.Environment, programID solana.PublicKey, authority solana.PublicKey, spillAddress solana.PublicKey) (solana.AccountMetaSlice, error) {
	idlAddress, err := getIDLAddress(e, programID)
	accounts := solana.AccountMetaSlice{
		solana.NewAccountMeta(idlAddress, true, false), // IDL account
		solana.NewAccountMeta(authority, false, true),
		solana.NewAccountMeta(spillAddress, true, false), // sol destination for close funds
	}
	return accounts, err
}

// generate close IDL PDA ix for a program via timelock
func removeIDLIx(e cldf.Environment, programID, programName string, spillAddress solana.PublicKey, c IDLConfig) (*mcmsTypes.Transaction, error) {
	timelockSignerPDA, err := FetchTimelockSigner(e, c.ChainSelector)
	if err != nil {
		return nil, fmt.Errorf("error loading timelockSignerPDA: %w", err)
	}
	authority := e.BlockChains.SolanaChains()[c.ChainSelector].DeployerKey.PublicKey()
	if c.MCMS != nil {
		authority = timelockSignerPDA
	}
	instruction, err := closeIdlInstruction(e, solana.MustPublicKeyFromBase58(programID), authority, spillAddress)
	if err != nil {
		return nil, fmt.Errorf("error closing IDL account ix: %w", err)
	}
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

	return repoSetup(e, chain, c.GitCommitSha)
}

// changeset to upload idl for a program
func UploadIDL(e cldf.Environment, c IDLConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error validating idl config: %w", err)
	}
	chain := e.BlockChains.SolanaChains()[c.ChainSelector]
	state, _ := stateview.LoadOnchainState(e)
	chainState := state.SolChains[c.ChainSelector]

	// start uploading
	if c.Router {
		err := idlInit(e, chain.ProgramsPath, chainState.Router.String(), deployment.RouterProgramName)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}
	if c.FeeQuoter {
		err := idlInit(e, chain.ProgramsPath, chainState.FeeQuoter.String(), deployment.FeeQuoterProgramName)
		if err != nil {
			return cldf.ChangesetOutput{}, nil
		}
	}
	if c.OffRamp {
		err := idlInit(e, chain.ProgramsPath, chainState.OffRamp.String(), deployment.OffRampProgramName)
		if err != nil {
			return cldf.ChangesetOutput{}, nil
		}
	}
	if c.RMNRemote {
		err := idlInit(e, chain.ProgramsPath, chainState.RMNRemote.String(), deployment.RMNRemoteProgramName)
		if err != nil {
			return cldf.ChangesetOutput{}, nil
		}
	}
	for _, bnmMetadata := range c.BurnMintTokenPoolMetadata {
		tokenPool := chainState.GetActiveTokenPool(shared.BurnMintTokenPool, bnmMetadata)
		err := idlInit(e, chain.ProgramsPath, tokenPool.String(), deployment.BurnMintTokenPoolProgramName)
		if err != nil {
			return cldf.ChangesetOutput{}, nil
		}
	}
	for _, lrMetadata := range c.LockReleaseTokenPoolMetadata {
		tokenPool := chainState.GetActiveTokenPool(shared.LockReleaseTokenPool, lrMetadata)
		err := idlInit(e, chain.ProgramsPath, tokenPool.String(), deployment.LockReleaseTokenPoolProgramName)
		if err != nil {
			return cldf.ChangesetOutput{}, nil
		}
	}
	if c.CCTPTokenPool {
		err := idlInit(e, chain.ProgramsPath, chainState.CCTPTokenPool.String(), deployment.CCTPTokenPoolProgramName)
		if err != nil {
			return cldf.ChangesetOutput{}, nil
		}
	}
	addresses, err := e.ExistingAddresses.AddressesForChain(c.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get existing addresses: %w", err)
	}
	mcmState, err := commonstate.MaybeLoadMCMSWithTimelockChainStateSolana(e.BlockChains.SolanaChains()[c.ChainSelector], addresses)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load MCMS with timelock chain state: %w", err)
	}
	if c.MCM {
		err := idlInit(e, chain.ProgramsPath, mcmState.McmProgram.String(), deployment.McmProgramName)
		if err != nil {
			return cldf.ChangesetOutput{}, nil
		}
	}
	if c.Timelock {
		err := idlInit(e, chain.ProgramsPath, mcmState.TimelockProgram.String(), deployment.TimelockProgramName)
		if err != nil {
			return cldf.ChangesetOutput{}, nil
		}
	}
	if c.AccessController {
		err := idlInit(e, chain.ProgramsPath, mcmState.AccessControllerProgram.String(), deployment.AccessControllerProgramName)
		if err != nil {
			return cldf.ChangesetOutput{}, nil
		}
	}

	return cldf.ChangesetOutput{}, nil
}

// changeset to set idl authority for a program to timelock
func SetAuthorityIDLByDeployerKey(e cldf.Environment, c IDLConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error validating idl config: %w", err)
	}
	state, _ := stateview.LoadOnchainState(e)
	chainState := state.SolChains[c.ChainSelector]
	chain := e.BlockChains.SolanaChains()[c.ChainSelector]

	// set idl authority
	mcmsTxs := make([]mcmsTypes.Transaction, 0)
	newAuthority := e.BlockChains.SolanaChains()[c.ChainSelector].DeployerKey.PublicKey()
	if c.Router {
		setAuthorityTx, err := setAuthorityIDLIx(e, chainState.Router.String(), deployment.RouterProgramName, newAuthority, c)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		if setAuthorityTx != nil {
			mcmsTxs = append(mcmsTxs, *setAuthorityTx)
		}
	}
	if c.FeeQuoter {
		setAuthorityTx, err := setAuthorityIDLIx(e, chainState.FeeQuoter.String(), deployment.FeeQuoterProgramName, newAuthority, c)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		if setAuthorityTx != nil {
			mcmsTxs = append(mcmsTxs, *setAuthorityTx)
		}
	}
	if c.OffRamp {
		setAuthorityTx, err := setAuthorityIDLIx(e, chainState.OffRamp.String(), deployment.OffRampProgramName, newAuthority, c)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		if setAuthorityTx != nil {
			mcmsTxs = append(mcmsTxs, *setAuthorityTx)
		}
	}
	if c.RMNRemote {
		setAuthorityTx, err := setAuthorityIDLIx(e, chainState.RMNRemote.String(), deployment.RMNRemoteProgramName, newAuthority, c)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		if setAuthorityTx != nil {
			mcmsTxs = append(mcmsTxs, *setAuthorityTx)
		}
	}
	for _, bnmMetadata := range c.BurnMintTokenPoolMetadata {
		tokenPool := chainState.GetActiveTokenPool(shared.BurnMintTokenPool, bnmMetadata)
		setAuthorityTx, err := setAuthorityIDLIx(e, tokenPool.String(), deployment.BurnMintTokenPoolProgramName, newAuthority, c)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		if setAuthorityTx != nil {
			mcmsTxs = append(mcmsTxs, *setAuthorityTx)
		}
	}
	for _, lrMetadata := range c.LockReleaseTokenPoolMetadata {
		tokenPool := chainState.GetActiveTokenPool(shared.LockReleaseTokenPool, lrMetadata)
		setAuthorityTx, err := setAuthorityIDLIx(e, tokenPool.String(), deployment.LockReleaseTokenPoolProgramName, newAuthority, c)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		if setAuthorityTx != nil {
			mcmsTxs = append(mcmsTxs, *setAuthorityTx)
		}
	}
	if c.CCTPTokenPool {
		setAuthorityTx, err := setAuthorityIDLIx(e, chainState.CCTPTokenPool.String(), deployment.CCTPTokenPoolProgramName, newAuthority, c)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		if setAuthorityTx != nil {
			mcmsTxs = append(mcmsTxs, *setAuthorityTx)
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
		setAuthorityTx, err := setAuthorityIDLIx(e, mcmState.AccessControllerProgram.String(), types.AccessControllerProgram.String(), newAuthority, c)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		if setAuthorityTx != nil {
			mcmsTxs = append(mcmsTxs, *setAuthorityTx)
		}
	}
	if c.Timelock {
		setAuthorityTx, err := setAuthorityIDLIx(e, mcmState.TimelockProgram.String(), types.RBACTimelockProgram.String(), newAuthority, c)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		if setAuthorityTx != nil {
			mcmsTxs = append(mcmsTxs, *setAuthorityTx)
		}
	}
	if c.MCM {
		setAuthorityTx, err := setAuthorityIDLIx(e, mcmState.McmProgram.String(), types.ManyChainMultisigProgram.String(), newAuthority, c)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		if setAuthorityTx != nil {
			mcmsTxs = append(mcmsTxs, *setAuthorityTx)
		}
	}

	if len(mcmsTxs) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, c.ChainSelector, "proposal to upgrade CCIP contracts", c.MCMS.MinDelay, mcmsTxs)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}

		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
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
	if c.Router {
		upgradeTx, err := upgradeIDLIx(e, chain.ProgramsPath, chainState.Router.String(), deployment.RouterProgramName, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	if c.FeeQuoter {
		upgradeTx, err := upgradeIDLIx(e, chain.ProgramsPath, chainState.FeeQuoter.String(), deployment.FeeQuoterProgramName, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	if c.OffRamp {
		upgradeTx, err := upgradeIDLIx(e, chain.ProgramsPath, chainState.OffRamp.String(), deployment.OffRampProgramName, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	if c.RMNRemote {
		upgradeTx, err := upgradeIDLIx(e, chain.ProgramsPath, chainState.RMNRemote.String(), deployment.RMNRemoteProgramName, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	for _, bnmMetadata := range c.BurnMintTokenPoolMetadata {
		tokenPool := chainState.GetActiveTokenPool(shared.BurnMintTokenPool, bnmMetadata)
		upgradeTx, err := upgradeIDLIx(e, chain.ProgramsPath, tokenPool.String(), deployment.BurnMintTokenPoolProgramName, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	for _, lrMetadata := range c.LockReleaseTokenPoolMetadata {
		tokenPool := chainState.GetActiveTokenPool(shared.LockReleaseTokenPool, lrMetadata)
		upgradeTx, err := upgradeIDLIx(e, chain.ProgramsPath, tokenPool.String(), deployment.LockReleaseTokenPoolProgramName, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	if c.CCTPTokenPool {
		tokenPool := chainState.GetActiveTokenPool(shared.CCTPTokenPool, shared.CLLMetadata)
		upgradeTx, err := upgradeIDLIx(e, chain.ProgramsPath, tokenPool.String(), deployment.CCTPTokenPoolProgramName, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
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
		upgradeTx, err := upgradeIDLIx(e, chain.ProgramsPath, mcmState.AccessControllerProgram.String(), deployment.AccessControllerProgramName, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	if c.Timelock {
		upgradeTx, err := upgradeIDLIx(e, chain.ProgramsPath, mcmState.TimelockProgram.String(), deployment.TimelockProgramName, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	if c.MCM {
		upgradeTx, err := upgradeIDLIx(e, chain.ProgramsPath, mcmState.McmProgram.String(), deployment.McmProgramName, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}

	if len(mcmsTxs) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, c.ChainSelector, "proposal to upgrade CCIP contracts", c.MCMS.MinDelay, mcmsTxs)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}

		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}

// changeset to close idl account for a program - this is needed when the idl increased so much in size that it no longer fits in the account
// and the idl account can not be resized when it already contains data, so you need to close the account and reinitialize it
func CloseIDLs(e cldf.Environment, c IDLConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error validating idl config: %w", err)
	}
	if c.SpillAddress.IsZero() {
		return cldf.ChangesetOutput{}, fmt.Errorf("error validating idl config, spill address must be set")
	}
	chain := e.BlockChains.SolanaChains()[c.ChainSelector]
	state, _ := stateview.LoadOnchainState(e)
	chainState := state.SolChains[c.ChainSelector]
	spillAddress := c.SpillAddress

	mcmsTxs := make([]mcmsTypes.Transaction, 0)
	if c.Router {
		upgradeTx, err := removeIDLIx(e, chainState.Router.String(), deployment.RouterProgramName, spillAddress, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	if c.FeeQuoter {
		upgradeTx, err := removeIDLIx(e, chainState.FeeQuoter.String(), deployment.FeeQuoterProgramName, spillAddress, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	if c.OffRamp {
		upgradeTx, err := removeIDLIx(e, chainState.OffRamp.String(), deployment.OffRampProgramName, spillAddress, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	if c.RMNRemote {
		upgradeTx, err := removeIDLIx(e, chainState.RMNRemote.String(), deployment.RMNRemoteProgramName, spillAddress, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	for _, bnmMetadata := range c.BurnMintTokenPoolMetadata {
		tokenPool := chainState.GetActiveTokenPool(shared.BurnMintTokenPool, bnmMetadata)
		upgradeTx, err := removeIDLIx(e, tokenPool.String(), deployment.BurnMintTokenPoolProgramName, spillAddress, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	for _, lrMetadata := range c.LockReleaseTokenPoolMetadata {
		tokenPool := chainState.GetActiveTokenPool(shared.LockReleaseTokenPool, lrMetadata)
		upgradeTx, err := removeIDLIx(e, tokenPool.String(), deployment.LockReleaseTokenPoolProgramName, spillAddress, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	if c.CCTPTokenPool {
		tokenPool := chainState.GetActiveTokenPool(shared.CCTPTokenPool, shared.CLLMetadata)
		upgradeTx, err := removeIDLIx(e, tokenPool.String(), deployment.CCTPTokenPoolProgramName, spillAddress, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
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
		upgradeTx, err := removeIDLIx(e, mcmState.AccessControllerProgram.String(), deployment.AccessControllerProgramName, spillAddress, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	if c.Timelock {
		upgradeTx, err := removeIDLIx(e, mcmState.TimelockProgram.String(), deployment.TimelockProgramName, spillAddress, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}
	if c.MCM {
		upgradeTx, err := removeIDLIx(e, mcmState.McmProgram.String(), deployment.McmProgramName, spillAddress, c)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error generating upgrade tx: %w", err)
		}
		if upgradeTx != nil {
			mcmsTxs = append(mcmsTxs, *upgradeTx)
		}
	}

	if len(mcmsTxs) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, c.ChainSelector, "proposal to upgrade CCIP contracts", c.MCMS.MinDelay, mcmsTxs)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}

		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}
