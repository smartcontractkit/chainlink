package config

import (
	"strconv"

	"github.com/pkg/errors"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func GetChainFromSpec(spec *job.OCR2OracleSpec, chainSet legacyevm.LegacyChainContainer) (legacyevm.Chain, int64, error) {
	chainIDInterface, ok := spec.RelayConfig["chainID"]
	if !ok {
		return nil, 0, errors.New("chainID must be provided in relay config")
	}
	destChainID := uint64(chainIDInterface.(float64))
	return GetChainByChainID(chainSet, destChainID)
}

func GetChainByChainSelector(chainSet legacyevm.LegacyChainContainer, chainSelector uint64) (legacyevm.Chain, int64, error) {
	chainID, err := chainselectors.ChainIdFromSelector(chainSelector)
	if err != nil {
		return nil, 0, err
	}
	return GetChainByChainID(chainSet, chainID)
}

func GetChainByChainID(chainSet legacyevm.LegacyChainContainer, chainID uint64) (legacyevm.Chain, int64, error) {
	chain, err := chainSet.Get(strconv.FormatUint(chainID, 10))
	if err != nil {
		return nil, 0, errors.Wrap(err, "chain not found in chainset")
	}
	return chain, chain.ID().Int64(), nil
}

func ResolveChainNames(sourceChainId int64, destChainId int64) (string, string, error) {
	sourceChainName, err := chainselectors.NameFromChainId(uint64(sourceChainId))
	if err != nil {
		return "", "", err
	}
	destChainName, err := chainselectors.NameFromChainId(uint64(destChainId))
	if err != nil {
		return "", "", err
	}
	return sourceChainName, destChainName, nil
}
