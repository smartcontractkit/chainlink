package config

import (
	"strconv"

	"github.com/pkg/errors"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func GetChainFromSpec(spec *job.OCR2OracleSpec, chainSet evm.LegacyChainContainer) (evm.Chain, int64, error) {
	chainIDInterface, ok := spec.RelayConfig["chainID"]
	if !ok {
		return nil, 0, errors.New("chainID must be provided in relay config")
	}
	destChainID := uint64(chainIDInterface.(float64))
	return GetChainByChainID(chainSet, destChainID)
}

func GetChainByChainSelector(chainSet evm.LegacyChainContainer, chainSelector uint64) (evm.Chain, int64, error) {
	chainID, err := chainselectors.ChainIdFromSelector(chainSelector)
	if err != nil {
		return nil, 0, err
	}
	return GetChainByChainID(chainSet, chainID)
}

func GetChainByChainID(chainSet evm.LegacyChainContainer, chainID uint64) (evm.Chain, int64, error) {
	chain, err := chainSet.Get(strconv.FormatUint(chainID, 10))
	if err != nil {
		return nil, 0, errors.Wrap(err, "chain not found in chainset")
	}
	return chain, chain.ID().Int64(), nil
}
