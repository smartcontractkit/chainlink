package ccipbase

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/gagliardetto/solana-go"
	signerRegistry "github.com/smartcontractkit/ccip-base/chains/solana/go_bindings"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// use this changeset to rotate NOPs (entirely remove addresses or add new ones)
var _ cldf.ChangeSet[RotateBaseSignerNopsConfig] = RotateBaseSignerNopsChangeset

// use this changeset to begin a key rotation (add green keys to blue ones)
var _ cldf.ChangeSet[AddGreenKeysConfig] = AddGreenKeysChangeset

// use this changeset to finalize a key rotation (promote green keys to blue ones)
var _ cldf.ChangeSet[PromoteKeysConfig] = PromoteKeysChangeset

type RotateBaseSignerNopsConfig struct {
	ChainSelector   uint64
	NopKeysToAdd    []string
	NopKeysToRemove []string
}

type AddGreenKeysConfig struct {
	ChainSelector uint64
	// Pairs of blue key (existing on the account) and new green key for that NOP
	BlueGreenKeys [][2]string
}

type PromoteKeysConfig struct {
	ChainSelector uint64
	// Keys to promote (nops can be identified by blue or green indistinctly)
	KeysToPromote []string
}

func RotateBaseSignerNopsChangeset(e cldf.Environment, c RotateBaseSignerNopsConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Rotating Base signer nops", "chain_selector", c.ChainSelector, "removing", c.NopKeysToAdd, "adding", c.NopKeysToAdd)
	chain := e.BlockChains.SolanaChains()[c.ChainSelector]
	c.Validate(e)

	configPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("config")}, signerRegistry.ProgramID)
	signersPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("signers")}, signerRegistry.ProgramID)
	eventAuthorityPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("__event_authority")}, signerRegistry.ProgramID)

	for _, hexKey := range c.NopKeysToRemove {
		key, _ := parseEVMAddress(hexKey)

		ix, err := signerRegistry.NewRemoveSignerInstruction(key, chain.DeployerKey.PublicKey(), configPda, signersPda, eventAuthorityPda, signerRegistry.ProgramID)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("Failed to remove signer: %w", err)
		}
		if err := chain.Confirm([]solana.Instruction{ix}); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("Failed to remove signer: %w", err)
		}
	}
	for _, hexKey := range c.NopKeysToAdd {
		key, _ := parseEVMAddress(hexKey)
		ix, err := signerRegistry.NewAddSignerInstruction(key, chain.DeployerKey.PublicKey(), configPda, signersPda, eventAuthorityPda, signerRegistry.ProgramID)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("Failed to add signer: %w", err)
		}
		if err := chain.Confirm([]solana.Instruction{ix}); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("Failed to add signer: %w", err)
		}
	}

	return cldf.ChangesetOutput{}, nil
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

	var signersAccount signerRegistry.Signers
	signersPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("signers")}, signerRegistry.ProgramID)

	chain := e.BlockChains.SolanaChains()[c.ChainSelector]
	if err := chain.GetAccountDataBorshInto(e.GetContext(), signersPda, &signersAccount); err != nil {
		return fmt.Errorf("failed to get signers: %w", err)
	}

	// Check that all NopKeysToRemove exist in signersAccount
	for i, keyBytes := range keysToRemoveParsed {
		if !keyExistsInSigners(keyBytes, signersAccount.Signers) {
			return fmt.Errorf("NopKeysToRemove[%d] (%s) does not exist in current signers", i, c.NopKeysToRemove[i])
		}
	}

	// Check that none of NopKeysToAdd already exist in signersAccount
	for i, keyBytes := range keysToAddParsed {
		if keyExistsInSigners(keyBytes, signersAccount.Signers) {
			return fmt.Errorf("NopKeysToAdd[%d] (%s) already exists in current signers", i, c.NopKeysToAdd[i])
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

	// Check that the final number of signers doesn't exceed 20
	currentSignerCount := len(signersAccount.Signers)
	finalSignerCount := currentSignerCount - len(keysToRemoveParsed) + len(keysToAddParsed)
	if finalSignerCount > 20 {
		return fmt.Errorf("final signer count would be %d, which exceeds the maximum of 20 (current: %d, removing: %d, adding: %d)",
			finalSignerCount, currentSignerCount, len(keysToRemoveParsed), len(keysToAddParsed))
	}

	return nil
}

func AddGreenKeysChangeset(e cldf.Environment, c AddGreenKeysConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Adding green keys to begin rotation", "chain_selector", c.ChainSelector)
	chain := e.BlockChains.SolanaChains()[c.ChainSelector]
	c.Validate(e)
	configPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("config")}, signerRegistry.ProgramID)
	signersPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("signers")}, signerRegistry.ProgramID)
	eventAuthorityPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("__event_authority")}, signerRegistry.ProgramID)

	for _, keyPair := range c.BlueGreenKeys {
		blue, _ := parseEVMAddress(keyPair[0])
		green, _ := parseEVMAddress(keyPair[1])

		ix, err := signerRegistry.NewSetSignerNewAddressInstruction(blue, green, chain.DeployerKey.PublicKey(), configPda, signersPda, eventAuthorityPda, signerRegistry.ProgramID)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("Failed to add green key: %w", err)
		}
		if err := chain.Confirm([]solana.Instruction{ix}); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("Failed to add green key: %w", err)
		}
	}

	return cldf.ChangesetOutput{}, nil
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
	var signersAccount signerRegistry.Signers
	signersPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("signers")}, signerRegistry.ProgramID)

	if err := chain.GetAccountDataBorshInto(e.GetContext(), signersPda, &signersAccount); err != nil {
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

	return nil
}

func PromoteKeysChangeset(e cldf.Environment, c PromoteKeysConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Promoting green keys to finalize rotation", "chain_selector", c.ChainSelector, "keys", c.KeysToPromote)
	chain := e.BlockChains.SolanaChains()[c.ChainSelector]
	c.Validate(e)
	configPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("config")}, signerRegistry.ProgramID)
	signersPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("signers")}, signerRegistry.ProgramID)
	eventAuthorityPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("__event_authority")}, signerRegistry.ProgramID)

	for _, keyHex := range c.KeysToPromote {
		key, _ := parseEVMAddress(keyHex)

		ix, err := signerRegistry.NewPromoteSignerAddressInstruction(key, chain.DeployerKey.PublicKey(), configPda, signersPda, eventAuthorityPda, signerRegistry.ProgramID)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("Failed to promote key: %w", err)
		}
		if err := chain.Confirm([]solana.Instruction{ix}); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("Failed to promote key: %w", err)
		}
	}

	return cldf.ChangesetOutput{}, nil
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
	var signersAccount signerRegistry.Signers
	signersPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("signers")}, signerRegistry.ProgramID)

	if err := chain.GetAccountDataBorshInto(e.GetContext(), signersPda, &signersAccount); err != nil {
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

	return nil
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
func keyExistsInSigners(key [20]uint8, signers []signerRegistry.Signer) bool {
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
func findSignerWithKey(key [20]uint8, signers []signerRegistry.Signer) *signerRegistry.Signer {
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
