package versionfinder

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/ccipcalc"
)

type EvmVersionFinder struct{}

func NewEvmVersionFinder() EvmVersionFinder {
	return EvmVersionFinder{}
}

func (e EvmVersionFinder) TypeAndVersion(addr ccip.Address, client bind.ContractBackend) (config.ContractType, semver.Version, error) {
	evmAddr, err := ccipcalc.GenericAddrToEvm(addr)
	if err != nil {
		return "", semver.Version{}, err
	}
	return config.TypeAndVersion(evmAddr, client)
}
