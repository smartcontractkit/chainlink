package sets

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	consensusv2capability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/consensus/v2"
	croncapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/cron"
	evmcapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/evm"
	logeventtriggercapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/logeventtrigger"
	mockcapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/mock"
	readcontractcapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/readcontract"
	writesolanacapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/writesolana"
)

func NewDefaultSet(homeChainID uint64) ([]cre.InstallableCapability, error) {
	capabilities := []cre.InstallableCapability{}

	cron, cErr := croncapability.New()
	if cErr != nil {
		return nil, errors.Wrap(cErr, "failed to create cron capability")
	}
	capabilities = append(capabilities, cron)

	c2, c2Err := consensusv2capability.New()
	if c2Err != nil {
		return nil, errors.Wrap(c2Err, "failed to create consensus capability v2")
	}
	capabilities = append(capabilities, c2)

	evm, evmErr := evmcapability.New(homeChainID)
	if evmErr != nil {
		return nil, errors.Wrap(evmErr, "failed to create evm capability")
	}
	capabilities = append(capabilities, evm)

	mock, mockErr := mockcapability.New()
	if mockErr != nil {
		return nil, errors.Wrap(mockErr, "failed to create mock capability")
	}
	capabilities = append(capabilities, mock)

	writesol, writeSolErr := writesolanacapability.New()
	if writeSolErr != nil {
		return nil, errors.Wrap(writeSolErr, "failed to create write solana capability")
	}
	capabilities = append(capabilities, writesol)

	readContract, readContractErr := readcontractcapability.New()
	if readContractErr != nil {
		return nil, errors.Wrap(readContractErr, "failed to create read contract capability")
	}
	capabilities = append(capabilities, readContract)

	logeventtrigger, logeventtriggerErr := logeventtriggercapability.New()
	if logeventtriggerErr != nil {
		return nil, errors.Wrap(logeventtriggerErr, "failed to create log event trigger capability")
	}
	capabilities = append(capabilities, logeventtrigger)

	return capabilities, nil
}
