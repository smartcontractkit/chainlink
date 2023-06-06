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

func getChain(cs evm.ChainSet, chainIDstr string) (chain evm.Chain, err error) {
	if chainIDstr != "" && chainIDstr != "<nil>" {
		chainID, ok := big.NewInt(0).SetString(chainIDstr, 10)
		if !ok {
			return nil, ErrInvalidChainID
		}
		chain, err = cs.Get(chainID)
		if err != nil {
			return nil, ErrMissingChainID
		}
		return chain, nil
	}

	if cs.ChainCount() > 1 {
		return nil, ErrMultipleChains
	}
	chain, err = cs.Default()
	if err != nil {
		return nil, err
	}
	return chain, nil
}
