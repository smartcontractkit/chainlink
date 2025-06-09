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
	// Cosmos for the Cosmos chain
	Cosmos ChainType = "cosmos"
	// Solana for the Solana chain
	Solana ChainType = "solana"
	// StarkNet for the StarkNet chain
	StarkNet ChainType = "starknet"
	// Aptos for the Aptos chain
	Aptos ChainType = "aptos"
	// Tron for the Tron chain
	Tron ChainType = "tron"
	// TON for the TON chain
	TON ChainType = "ton"
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

func NewChainType(typ uint8) (ChainType, error) {
	switch typ {
	case 1:
		return EVM, nil
	case 2:
		return Solana, nil
	case 3:
		return Cosmos, nil
	case 4:
		return StarkNet, nil
	case 5:
		return Aptos, nil
	case 6:
		return Tron, nil
	case 7:
		return TON, nil
	default:
		return "", fmt.Errorf("unexpected chaintype.ChainType: %#v", typ)
	}
}

func (c ChainType) Type() (uint8, error) {
	switch c {
	case EVM:
		return 1, nil
	case Solana:
		return 2, nil
	case Cosmos:
		return 3, nil
	case StarkNet:
		return 4, nil
	case Aptos:
		return 5, nil
	case Tron:
		return 6, nil
	case TON:
		return 7, nil
	default:
		return 0, fmt.Errorf("unexpected chaintype.ChainType: %#v", c)
	}
}

// SupportedChainTypes contain all chains that are supported
var SupportedChainTypes = ChainTypes{EVM, Cosmos, Solana, StarkNet, Aptos, Tron, TON}

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
