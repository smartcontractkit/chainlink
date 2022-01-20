package chaintype

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type ChainType string

const (
	EVM    ChainType = "evm"
	Solana ChainType = "solana"
	Terra  ChainType = "terra"
)

var SupportedChainTypes = []ChainType{EVM, Solana, Terra}
var ErrInvalidChainType error

func init() {
	supported := make([]string, 0, len(SupportedChainTypes))
	for _, chainType := range SupportedChainTypes {
		supported = append(supported, fmt.Sprintf(`"%s"`, chainType))
	}
	ErrInvalidChainType = fmt.Errorf("valid types include: [%s]", strings.Join(supported, ", "))
}

func IsSupportedChainType(chainType ChainType) bool {
	for _, v := range SupportedChainTypes {
		if v == chainType {
			return true
		}
	}
	return false
}

func NewErrInvalidChainType(chainType ChainType) error {
	return errors.Wrapf(ErrInvalidChainType, `unknown chain type "%s"`, chainType)
}
