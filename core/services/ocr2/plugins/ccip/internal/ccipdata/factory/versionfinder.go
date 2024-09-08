package factory

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
)

// VersionFinder accepts a contract address and a client and performs an on-chain call to
// determine the contract type.
type VersionFinder interface {
	TypeAndVersion(addr cciptypes.Address, client bind.ContractBackend) (config.ContractType, semver.Version, error)
}

type EvmVersionFinder struct{}

func NewEvmVersionFinder() EvmVersionFinder {
	return EvmVersionFinder{}
}

func (e EvmVersionFinder) TypeAndVersion(addr cciptypes.Address, client bind.ContractBackend) (config.ContractType, semver.Version, error) {
	evmAddr, err := ccipcalc.GenericAddrToEvm(addr)
	if err != nil {
		return "", semver.Version{}, err
	}
	return config.TypeAndVersion(evmAddr, client)
}

type mockVersionFinder struct {
	typ     config.ContractType
	version semver.Version
	err     error
}

func newMockVersionFinder(typ config.ContractType, version semver.Version, err error) *mockVersionFinder {
	return &mockVersionFinder{typ: typ, version: version, err: err}
}

func (m mockVersionFinder) TypeAndVersion(addr cciptypes.Address, client bind.ContractBackend) (config.ContractType, semver.Version, error) {
	return m.typ, m.version, m.err
}
