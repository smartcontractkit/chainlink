package types

import (
	"encoding/hex"
	"fmt"

	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
)

// NOTE: this is larger then the limit in the contract (64)
const maxWorkflowNameLength = 256

type WorkflowName interface {
	// A 10-byte hash of the name, hex-encoded to a 20-byte string.
	// Used in the metadata we send onchain and for authorizing
	// the workflow with the consumer contract.
	Hex() string

	// User-defined workflow name has can be of any length between
	// 1 and maxWorkflowNameLength. Used for logging and metrics.
	String() string
}

type workflowName struct {
	userDefinedName string
}

func (n workflowName) Hex() string {
	truncatedName := pkgworkflows.HashTruncateName(n.userDefinedName)
	hexName := hex.EncodeToString([]byte(truncatedName))
	return hexName
}

func (n workflowName) String() string {
	return n.userDefinedName
}

func NewWorkflowName(userDefinedName string) (WorkflowName, error) {
	if len(userDefinedName) == 0 || len(userDefinedName) > maxWorkflowNameLength {
		return nil, fmt.Errorf("invalid workflow name length: %d", len(userDefinedName))
	}
	return workflowName{userDefinedName: userDefinedName}, nil
}

type WorkflowID [32]byte

func (w WorkflowID) Hex() string {
	return hex.EncodeToString(w[:])
}

func (w WorkflowID) Equal(o WorkflowID) bool {
	return w.Hex() == o.Hex()
}

func WorkflowIDFromHex(h string) (WorkflowID, error) {
	b, err := hex.DecodeString(h)
	if err != nil {
		return [32]byte{}, err
	}

	if len(b) != 32 {
		return [32]byte{}, fmt.Errorf("invalid workflow id: incorrect length, expected 32, got %d", len(b))
	}

	return WorkflowID([32]byte(b)), nil
}

// expects a hex-encoded [20]byte string, no "0x" prefix
func ValidateWorkflowOwner(owner string) error {
	if len(owner) != 40 {
		return fmt.Errorf("invalid workflow owner: incorrect length, expected 40, got %d", len(owner))
	}
	if _, err := hex.DecodeString(owner); err != nil {
		return fmt.Errorf("invalid workflow owner: %w", err)
	}
	return nil
}
