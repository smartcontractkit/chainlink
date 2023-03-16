package chaintype

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// ChainType denotes the chain or network to work with
type ChainType string

const (
	// EVM for Ethereum or other chains supporting the EVM
	EVM ChainType = "evm"
	// Solana for the Solana chain
	Solana ChainType = "solana"
	// StarkNet for the StarkNet chain
	StarkNet ChainType = "starknet"
)

type ChainTypes []ChainType

func (c ChainTypes) String() (out string) {
	var sb strings.Builder
	for i, chain := range c {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(string(chain))
	}
	return sb.String()
}

// SupportedChainTypes contain all chains that are supported
var SupportedChainTypes = ChainTypes{EVM, Solana, StarkNet}

// ErrInvalidChainType is an error to indicate an unsupported chain type
var ErrInvalidChainType error

func init() {
	supported := make([]string, 0, len(SupportedChainTypes))
	for _, chainType := range SupportedChainTypes {
		supported = append(supported, fmt.Sprintf(`"%s"`, chainType))
	}
	ErrInvalidChainType = fmt.Errorf("valid types include: [%s]", strings.Join(supported, ", "))
}

// IsSupportedChainType checks to see if the chain is supported
func IsSupportedChainType(chainType ChainType) bool {
	for _, v := range SupportedChainTypes {
		if v == chainType {
			return true
		}
	}
	return false
}

// NewErrInvalidChainType returns an error wrapping ErrInvalidChainType for an unsupported chain
func NewErrInvalidChainType(chainType ChainType) error {
	return errors.Wrapf(ErrInvalidChainType, `unknown chain type "%s"`, chainType)
}
