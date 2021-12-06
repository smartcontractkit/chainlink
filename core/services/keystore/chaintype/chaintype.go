package chaintype

import (
	"fmt"

	"github.com/pkg/errors"
)

type ChainType string

var ErrInvalidChainType = fmt.Errorf(`invalid chain type: valid types include: "%s" and "%s"`, EVM, Solana)

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

func NewErrInvalidChainType(chainType ChainType) error {
	return errors.Wrapf(ErrInvalidChainType, "unknown chain type: %s", chainType)
}
