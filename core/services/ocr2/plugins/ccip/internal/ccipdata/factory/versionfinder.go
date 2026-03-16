package factory

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/version"
)

type mockVersionFinder struct {
	typ     version.ContractType
	version semver.Version
	err     error
}

func newMockVersionFinder(typ version.ContractType, version semver.Version, err error) *mockVersionFinder {
	return &mockVersionFinder{typ: typ, version: version, err: err}
}

func (m mockVersionFinder) TypeAndVersion(addr cciptypes.Address, client bind.ContractBackend) (version.ContractType, semver.Version, error) {
	return m.typ, m.version, m.err
}
