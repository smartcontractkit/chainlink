package ccipbase

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/gagliardetto/solana-go"
	signerRegistry "github.com/smartcontractkit/ccip-base/chains/solana/go_bindings"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

type RotateBaseSignerNopsConfig struct {
	ChainSelector   uint64
	NopKeysToAdd    []string
	NopKeysToRemove []string
}

// use these changesets to rotate NOPs (entirely remove addresses or add new ones)
var _ cldf.ChangeSet[RotateBaseSignerNopsConfig] = RotateBaseSignerNopsChangeset

func RotateBaseSignerNopsChangeset(e cldf.Environment, c RotateBaseSignerNopsConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Rotating Base signer keys", "chain_selector", c.ChainSelector)
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
	// Parse and validate NopKeysToAdd
	keysToAddParsed := make([][20]uint8, len(c.NopKeysToAdd))
	for i, key := range c.NopKeysToAdd {
		parsed, err := parseEVMAddress(key)
		if err != nil {
			return fmt.Errorf("invalid NopKeysToAdd[%d]: %w", i, err)
		}
		keysToAddParsed[i] = parsed
	}

	// Parse and validate NopKeysToRemove
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
