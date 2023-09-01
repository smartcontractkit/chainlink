package web

import (
	"math/big"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
)

var (
	ErrMissingChainID = errors.New("evmChainID does not match any local chains")
	ErrInvalidChainID = errors.New("invalid evmChainID")
	ErrMultipleChains = errors.New("more than one chain available, you must specify evmChainID parameter")
)

func getChain(legacyChains evm.LegacyChainContainer, chainIDstr string) (chain evm.Chain, err error) {

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

	chain, err = legacyChains.Default()
	if err != nil {
		return nil, err
	}
	return chain, nil
}
