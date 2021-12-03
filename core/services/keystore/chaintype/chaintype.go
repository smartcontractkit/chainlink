package chaintype

import "fmt"

type ChainType string
type ErrInvalidChainType error

const (
	EVM    ChainType = "evm"
	Solana ChainType = "solana"
)

var SupportedChainTypes = []ChainType{EVM, Solana}

func IsSupportedChainType(chainType ChainType) bool {
	for _, v := range SupportedChainTypes {
		if v == chainType {
			return true
		}
	}
	return false
}

func NewErrInvalidChainType(chainType ChainType) ErrInvalidChainType {
	return ErrInvalidChainType(
		fmt.Errorf(`invalid chain type "%s", valid types include: "%s" and "%s"`, chainType, EVM, Solana),
	)
}
