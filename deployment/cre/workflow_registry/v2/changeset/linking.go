package changeset

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type OwnershipProofSignaturePayload struct {
	RequestType              uint8          // should be uint8 in Solidity, 1 byte
	WorkflowOwnerAddress     common.Address // should be 20 bytes in Solidity, address type
	ChainID                  string         // should be uint256 in Solidity, chain-selectors provide it as a string
	WorkflowRegistryContract common.Address // address of the WorkflowRegistry contract, should be 20 bytes in Solidity
	Version                  string         // should be dynamic type in Solidity (string)
	ValidityTimestamp        time.Time      // should be uint256 in Solidity
	OwnershipProofHash       common.Hash    // should be bytes32 in Solidity, 32 bytes hash of the ownership proof
}

// Generates a hash for the ownership proof based on the workflow owner address, organization ID, and nonce.
func GenerateOwnershipProofHash(
	workflowOwnerAddress, organizationID, nonce string,
) string {
	data := workflowOwnerAddress + organizationID + nonce
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Convert payload fields into Solidity-compatible data types and concatenate them in the expected order.
// Use the same hashing algorithm as the Solidity contract (keccak256) to hash the concatenated data.
// Finally, follow the EIP-191 standard to create the final hash for signing.
func PreparePayloadForSigning(payload OwnershipProofSignaturePayload) ([]byte, error) {
	// Prepare a list of ABI arguments in the exact order as expected by the Solidity contract
	arguments, err := prepareABIArguments()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare ABI arguments: %w", err)
	}

	// Convert the payload fields to their respective types
	chainID := new(big.Int)
	chainID.SetString(payload.ChainID, 10)
	validityTimestamp := big.NewInt(payload.ValidityTimestamp.Unix())

	// Concatenate the fields, Solidity contract must follow the same order and use abi.encode()
	packed, err := arguments.Pack(
		payload.RequestType,
		payload.WorkflowOwnerAddress,
		chainID,
		payload.WorkflowRegistryContract,
		payload.Version,
		validityTimestamp,
		payload.OwnershipProofHash,
	)
	if err != nil {
		return nil, fmt.Errorf("abi encoding failed: %w", err)
	}

	// Hash the concatenated result using SHA256, Solidity contract will use keccak256()
	hash := crypto.Keccak256(packed)

	// Prepare a message that can be verified in a Solidity contract.
	// For a signature to be recoverable, it must follow the EIP-191 standard.
	// The message must be prefixed with "\x19Ethereum Signed Message:\n" followed by the length of the message.
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n32%s", hash)
	return crypto.Keccak256([]byte(prefixedMessage)), nil
}

// Prepare the ABI arguments, in the exact order as expected by the Solidity contract.
func prepareABIArguments() (*abi.Arguments, error) {
	arguments := abi.Arguments{}

	uint8Type, err := abi.NewType("uint8", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create uint8 type: %w", err)
	}

	addressType, err := abi.NewType("address", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create address type: %w", err)
	}

	bytes32Type, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create bytes32 type: %w", err)
	}

	uint256Type, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create uint256 type: %w", err)
	}

	stringType, err := abi.NewType("string", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create string type: %w", err)
	}

	arguments = append(arguments, abi.Argument{Type: uint8Type})   // request type
	arguments = append(arguments, abi.Argument{Type: addressType}) // owner address
	arguments = append(arguments, abi.Argument{Type: uint256Type}) // chain ID
	arguments = append(arguments, abi.Argument{Type: addressType}) // address of the contract
	arguments = append(arguments, abi.Argument{Type: stringType})  // version string
	arguments = append(arguments, abi.Argument{Type: uint256Type}) // validity timestamp
	arguments = append(arguments, abi.Argument{Type: bytes32Type}) // ownership proof hash

	return &arguments, nil
}
