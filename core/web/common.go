package web

import (
	"math/big"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
)

var (
	ErrMissingChainID = errors.New("chain id does not match any local chains")
	ErrEmptyChainID   = errors.New("chainID is empty")
	ErrInvalidChainID = errors.New("invalid chain id")
	ErrMultipleChains = errors.New("more than one chain available, you must specify chain id parameter")
)

func getChain(legacyChains legacyevm.LegacyChainContainer, chainIDstr string) (chain legacyevm.Chain, err error) {
	if chainIDstr != "" && chainIDstr != "<nil>" {
		// evm keys are expected to be parsable as a big int
		_, ok := big.NewInt(0).SetString(chainIDstr, 10)
		if !ok {
			return nil, ErrInvalidChainID
		}
		chain, err = legacyChains.Get(chainIDstr)
		if err != nil {
			return nil, ErrMissingChainID
		}
		return chain, nil
	}

	if legacyChains.Len() > 1 {
		return nil, ErrMultipleChains
	}

	return nil, ErrEmptyChainID
}
