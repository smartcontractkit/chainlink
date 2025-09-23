package ccip_attestation

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	signer_registry "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/signer_registry_solana"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cs_solana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
)

// executeOrBuildMCMSProposal handles the decision to execute instructions directly or build an MCMS proposal.
// If mcmsConfig is nil, it executes the instructions directly on the chain.
// If mcmsConfig is provided, it builds MCMS transactions and creates a proposal.
func executeOrBuildMCMSProposal(
	e cldf.Environment,
	chain *cldf_solana.Chain,
	instructions []solana.Instruction,
	programID string,
	contractType cldf.ContractType,
	mcmsConfig *proposalutils.TimelockConfig,
	proposalDescription string,
) (cldf.ChangesetOutput, error) {
	if mcmsConfig == nil {
		// Direct execution - confirm each instruction individually to avoid Tx size limits
		for i, ixn := range instructions {
			if err := chain.Confirm([]solana.Instruction{ixn}); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instruction %d: %w", i, err)
			}
		}
		return cldf.ChangesetOutput{}, nil
	}

	mcmsTxns := make([]mcmsTypes.Transaction, 0, len(instructions))
	for _, ixn := range instructions {
		tx, err := cs_solana.BuildMCMSTxn(ixn, programID, contractType)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create MCMS transaction: %w", err)
		}
		mcmsTxns = append(mcmsTxns, *tx)
	}

	proposal, err := cs_solana.BuildProposalsForTxns(
		e, chain.Selector, proposalDescription, mcmsConfig.MinDelay, mcmsTxns)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
	}, nil
}

// use this changeset to rotate NOPs (entirely remove addresses or add new ones)
var _ cldf.ChangeSet[RotateBaseSignerNopsConfig] = RotateBaseSignerNopsChangeset

// use this changeset to begin a key rotation (add green keys to blue ones)
var _ cldf.ChangeSet[AddGreenKeysConfig] = AddGreenKeysChangeset

// use this changeset to finalize a key rotation (promote green keys to blue ones)
var _ cldf.ChangeSet[PromoteKeysConfig] = PromoteKeysChangeset

// use this changeset to change upgrade authority
var _ cldf.ChangeSet[SetUpgradeAuthorityConfig] = SetUpgradeAuthorityChangeset

type RotateBaseSignerNopsConfig struct {
	ChainSelector   uint64
	NopKeysToAdd    []string
	NopKeysToRemove []string
	// if set, assumes current upgrade authority is the timelock
	MCMS *proposalutils.TimelockConfig
}

type AddGreenKeysConfig struct {
	ChainSelector uint64
	// Pairs of blue key (existing on the account) and new green key for that NOP
	BlueGreenKeys [][2]string
	// if set, assumes current upgrade authority is the timelock
	MCMS *proposalutils.TimelockConfig
}

type PromoteKeysConfig struct {
	ChainSelector uint64
	// Keys to promote (nops can be identified by blue or green indistinctly)
	KeysToPromote []string
	// if set, assumes current upgrade authority is the timelock
	MCMS *proposalutils.TimelockConfig
}

type SetUpgradeAuthorityConfig struct {
	ChainSelector       uint64
	NewUpgradeAuthority solana.PublicKey
	// if set, assumes current upgrade authority is the timelock
	MCMS *proposalutils.TimelockConfig
}

func RotateBaseSignerNopsChangeset(e cldf.Environment, c RotateBaseSignerNopsConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Rotating Base signer nops", "chain_selector", c.ChainSelector, "removing", c.NopKeysToRemove, "adding", c.NopKeysToAdd)
	chain := e.BlockChains.SolanaChains()[c.ChainSelector]
	err := c.Validate(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to rotate signer nop: %w", err)
	}

	configPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("config")}, signer_registry.ProgramID)
	signersPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("signers")}, signer_registry.ProgramID)
	eventAuthorityPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("__event_authority")}, signer_registry.ProgramID)
	currentAuthority := chain.DeployerKey.PublicKey()
	if c.MCMS != nil {
		timelockSignerPDA, err := cs_solana.FetchTimelockSigner(e, chain.Selector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get timelock signer: %w", err)
		}
		currentAuthority = timelockSignerPDA
	}

	var instructions []solana.Instruction

	for _, hexKey := range c.NopKeysToRemove {
		key, _ := parseEVMAddress(hexKey)

		ix, err := signer_registry.NewRemoveSignerInstruction(key, currentAuthority, configPda, signersPda, eventAuthorityPda, signer_registry.ProgramID)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to remove signer: %w", err)
		}
		instructions = append(instructions, ix)
	}
	for _, hexKey := range c.NopKeysToAdd {
		key, _ := parseEVMAddress(hexKey)
		ix, err := signer_registry.NewAddSignerInstruction(key, currentAuthority, configPda, signersPda, eventAuthorityPda, signer_registry.ProgramID)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to add signer: %w", err)
		}
		instructions = append(instructions, ix)
	}

	return executeOrBuildMCMSProposal(
		e,
		&chain,
		instructions,
		signer_registry.ProgramID.String(),
		cldf.ContractType("BaseSignerRegistry"),
		c.MCMS,
		"proposal to rotate attestation signer NOPs in Solana",
	)
}

func (c RotateBaseSignerNopsConfig) Validate(e cldf.Environment) error {
	keysToAddParsed := make([][20]uint8, len(c.NopKeysToAdd))
	for i, key := range c.NopKeysToAdd {
		parsed, err := parseEVMAddress(key)
		if err != nil {
			return fmt.Errorf("invalid NopKeysToAdd[%d]: %w", i, err)
		}
		keysToAddParsed[i] = parsed
	}

	keysToRemoveParsed := make([][20]uint8, len(c.NopKeysToRemove))
	for i, key := range c.NopKeysToRemove {
		parsed, err := parseEVMAddress(key)
		if err != nil {
			return fmt.Errorf("invalid NopKeysToRemove[%d]: %w", i, err)
		}
		keysToRemoveParsed[i] = parsed
	}

	chain := e.BlockChains.SolanaChains()[c.ChainSelector]
	if len(c.NopKeysToRemove) > 0 {
		signersPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("signers")}, signer_registry.ProgramID)

		data, err := solanastateview.GetAccountData(e, &chain, signersPda)
		if err != nil {
			return fmt.Errorf("failed to get signers: %w", err)
		}

		signersAccount, err := signer_registry.ParseAccount_Signers(data)

		if err != nil {
			return fmt.Errorf("failed to get signers: %w", err)
		}
		// Check that all NopKeysToRemove exist in signersAccount
		for i, keyBytes := range keysToRemoveParsed {
			if !keyExistsInSigners(keyBytes, signersAccount.Signers) {
				return fmt.Errorf("NopKeysToRemove[%d] (%s) does not exist in current signers", i, c.NopKeysToRemove[i])
			}
		}
	}

	// Check that there are no keys in common between ToAdd and ToRemove
	for i, addKey := range keysToAddParsed {
		for j, removeKey := range keysToRemoveParsed {
			if addKey == removeKey {
				return fmt.Errorf("key %s appears in both NopKeysToAdd[%d] and NopKeysToRemove[%d]", c.NopKeysToAdd[i], i, j)
			}
		}
	}

	return solanastateview.ValidateOwnershipSolana(&e, chain, c.MCMS != nil, signer_registry.ProgramID, shared.SVMSignerRegistry, solana.PublicKey{})
}

func AddGreenKeysChangeset(e cldf.Environment, c AddGreenKeysConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Adding green keys to begin rotation", "chain_selector", c.ChainSelector)
	chain := e.BlockChains.SolanaChains()[c.ChainSelector]
	err := c.Validate(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to add green keys: %w", err)
	}
	configPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("config")}, signer_registry.ProgramID)
	signersPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("signers")}, signer_registry.ProgramID)
	eventAuthorityPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("__event_authority")}, signer_registry.ProgramID)
	currentAuthority := chain.DeployerKey.PublicKey()
	if c.MCMS != nil {
		timelockSignerPDA, err := cs_solana.FetchTimelockSigner(e, chain.Selector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get timelock signer: %w", err)
		}
		currentAuthority = timelockSignerPDA
	}

	var instructions []solana.Instruction

	for _, keyPair := range c.BlueGreenKeys {
		blue, _ := parseEVMAddress(keyPair[0])
		green, _ := parseEVMAddress(keyPair[1])

		ix, err := signer_registry.NewSetSignerNewAddressInstruction(blue, green, currentAuthority, configPda, signersPda, eventAuthorityPda, signer_registry.ProgramID)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to add green key: %w", err)
		}
		instructions = append(instructions, ix)
	}

	return executeOrBuildMCMSProposal(
		e,
		&chain,
		instructions,
		signer_registry.ProgramID.String(),
		cldf.ContractType("BaseSignerRegistry"),
		c.MCMS,
		"proposal to add green keys for rotation in Solana",
	)
}

func (c AddGreenKeysConfig) Validate(e cldf.Environment) error {
	chain := e.BlockChains.SolanaChains()[c.ChainSelector]

	// Parse and validate all blue and green keys
	blueKeysParsed := make([][20]uint8, len(c.BlueGreenKeys))
	greenKeysParsed := make([][20]uint8, len(c.BlueGreenKeys))

	for i, keyPair := range c.BlueGreenKeys {
		blueParsed, err := parseEVMAddress(keyPair[0])
		if err != nil {
			return fmt.Errorf("invalid BlueGreenKeys[%d] blue key: %w", i, err)
		}
		blueKeysParsed[i] = blueParsed

		greenParsed, err := parseEVMAddress(keyPair[1])
		if err != nil {
			return fmt.Errorf("invalid BlueGreenKeys[%d] green key: %w", i, err)
		}
		greenKeysParsed[i] = greenParsed
	}

	// Get current signers account
	signersPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("signers")}, signer_registry.ProgramID)

	data, err := solanastateview.GetAccountData(e, &chain, signersPda)
	if err != nil {
		return fmt.Errorf("failed to get signers: %w", err)
	}

	signersAccount, err := signer_registry.ParseAccount_Signers(data)
	if err != nil {
		return fmt.Errorf("failed to get signers: %w", err)
	}

	// Check that all blue keys exist in signersAccount (either as EvmAddress or NewEvmAddress)
	for i, blueKey := range blueKeysParsed {
		if !keyExistsInSigners(blueKey, signersAccount.Signers) {
			return fmt.Errorf("BlueGreenKeys[%d] blue key (%s) does not exist in current signers", i, c.BlueGreenKeys[i][0])
		}
	}

	// Check that none of the green keys already exist in signersAccount
	for i, greenKey := range greenKeysParsed {
		if keyExistsInSigners(greenKey, signersAccount.Signers) {
			return fmt.Errorf("BlueGreenKeys[%d] green key (%s) already exists in current signers", i, c.BlueGreenKeys[i][1])
		}
	}

	// Check that no green key appears as a blue key in the same config (no circular references)
	for i, greenKey := range greenKeysParsed {
		for j, blueKey := range blueKeysParsed {
			if greenKey == blueKey {
				return fmt.Errorf("green key %s appears as both green key in BlueGreenKeys[%d] and blue key in BlueGreenKeys[%d]", c.BlueGreenKeys[i][1], i, j)
			}
		}
	}

	return solanastateview.ValidateOwnershipSolana(&e, chain, c.MCMS != nil, signer_registry.ProgramID, shared.SVMSignerRegistry, solana.PublicKey{})
}

func PromoteKeysChangeset(e cldf.Environment, c PromoteKeysConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Promoting green keys to finalize rotation", "chain_selector", c.ChainSelector, "keys", c.KeysToPromote)
	chain := e.BlockChains.SolanaChains()[c.ChainSelector]
	err := c.Validate(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy base signer registry contract: %w", err)
	}
	configPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("config")}, signer_registry.ProgramID)
	signersPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("signers")}, signer_registry.ProgramID)
	eventAuthorityPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("__event_authority")}, signer_registry.ProgramID)
	currentAuthority := chain.DeployerKey.PublicKey()
	if c.MCMS != nil {
		timelockSignerPDA, err := cs_solana.FetchTimelockSigner(e, chain.Selector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get timelock signer: %w", err)
		}
		currentAuthority = timelockSignerPDA
	}

	var instructions []solana.Instruction

	for _, keyHex := range c.KeysToPromote {
		key, _ := parseEVMAddress(keyHex)

		ix, err := signer_registry.NewPromoteSignerAddressInstruction(key, currentAuthority, configPda, signersPda, eventAuthorityPda, signer_registry.ProgramID)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to promote key: %w", err)
		}
		instructions = append(instructions, ix)
	}

	return executeOrBuildMCMSProposal(
		e,
		&chain,
		instructions,
		signer_registry.ProgramID.String(),
		cldf.ContractType("BaseSignerRegistry"),
		c.MCMS,
		"proposal to promote green keys in Solana",
	)
}

func (c PromoteKeysConfig) Validate(e cldf.Environment) error {
	chain := e.BlockChains.SolanaChains()[c.ChainSelector]

	// Parse and validate all keys to promote
	keysParsed := make([][20]uint8, len(c.KeysToPromote))
	for i, key := range c.KeysToPromote {
		parsed, err := parseEVMAddress(key)
		if err != nil {
			return fmt.Errorf("invalid KeysToPromote[%d]: %w", i, err)
		}
		keysParsed[i] = parsed
	}

	// Get current signers account
	signersPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("signers")}, signer_registry.ProgramID)

	data, err := solanastateview.GetAccountData(e, &chain, signersPda)
	if err != nil {
		return fmt.Errorf("failed to get signers: %w", err)
	}

	signersAccount, err := signer_registry.ParseAccount_Signers(data)
	if err != nil {
		return fmt.Errorf("failed to get signers: %w", err)
	}

	// Check that each key exists and has an active blue/green pair
	for i, keyBytes := range keysParsed {
		signer := findSignerWithKey(keyBytes, signersAccount.Signers)
		if signer == nil {
			return fmt.Errorf("KeysToPromote[%d] (%s) does not exist in current signers", i, c.KeysToPromote[i])
		}

		// Check that this signer has an active blue/green pair (NewEvmAddress is non-zero)
		if signer.NewEvmAddress == nil {
			return fmt.Errorf("KeysToPromote[%d] (%s) does not have a green key to promote", i, c.KeysToPromote[i])
		}
	}

	return solanastateview.ValidateOwnershipSolana(&e, chain, c.MCMS != nil, signer_registry.ProgramID, shared.SVMSignerRegistry, solana.PublicKey{})
}

func (c SetUpgradeAuthorityConfig) Validate(e cldf.Environment) error {
	chain := e.BlockChains.SolanaChains()[c.ChainSelector]
	return solanastateview.ValidateOwnershipSolana(&e, chain, c.MCMS != nil, signer_registry.ProgramID, shared.SVMSignerRegistry, solana.PublicKey{})
}

func SetUpgradeAuthorityChangeset(
	e cldf.Environment,
	config SetUpgradeAuthorityConfig,
) (cldf.ChangesetOutput, error) {
	chain := e.BlockChains.SolanaChains()[config.ChainSelector]
	currentAuthority := chain.DeployerKey.PublicKey()
	if config.MCMS != nil {
		timelockSignerPDA, err := cs_solana.FetchTimelockSigner(e, chain.Selector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get timelock signer: %w", err)
		}
		currentAuthority = timelockSignerPDA
	}
	e.Logger.Infow("Setting upgrade authority", "newUpgradeAuthority", config.NewUpgradeAuthority.String())

	ixn := cs_solana.SetUpgradeAuthority(&e, &chain, signer_registry.ProgramID, currentAuthority, config.NewUpgradeAuthority, false)

	return executeOrBuildMCMSProposal(
		e,
		&chain,
		[]solana.Instruction{ixn},
		solana.BPFLoaderUpgradeableProgramID.String(),
		cldf.ContractType(solana.BPFLoaderUpgradeableProgramID.String()),
		config.MCMS,
		"proposal to SetUpgradeAuthority in Solana",
	)
}

func parseEVMAddress(addr string) ([20]uint8, error) {
	if strings.HasPrefix(addr, "0x") || strings.HasPrefix(addr, "0X") {
		addr = addr[2:]
	}

	decoded, err := hex.DecodeString(addr)
	if err != nil {
		return [20]uint8{}, err
	}

	if len(decoded) != 20 {
		return [20]uint8{}, fmt.Errorf("expected 20 bytes, got %d", len(decoded))
	}

	var result [20]uint8
	copy(result[:], decoded)
	return result, nil
}

// keyExistsInSigners checks if a key exists in the signers list (either as EvmAddress or NewEvmAddress)
func keyExistsInSigners(key [20]uint8, signers []signer_registry.Signer) bool {
	// Special case: the zero-key is never considered to be in the signers list (as it's an alias for "removal")
	var zeroKey [20]uint8
	if key == zeroKey {
		return false
	}

	for _, signer := range signers {
		// Check current EvmAddress
		if signer.EvmAddress == key {
			return true
		}
		// Check NewEvmAddress if it exists
		if signer.NewEvmAddress != nil && *signer.NewEvmAddress == key {
			return true
		}
	}
	return false
}

// findSignerWithKey finds and returns the signer that contains the given key (either as EvmAddress or NewEvmAddress)
func findSignerWithKey(key [20]uint8, signers []signer_registry.Signer) *signer_registry.Signer {
	// Return nil for all-zero keys
	var zeroKey [20]uint8
	if key == zeroKey {
		return nil
	}

	for i := range signers {
		signer := &signers[i]
		// Check current EvmAddress
		if signer.EvmAddress == key {
			return signer
		}
		// Check NewEvmAddress if it exists
		if signer.NewEvmAddress != nil && *signer.NewEvmAddress == key {
			return signer
		}
	}
	return nil
}
