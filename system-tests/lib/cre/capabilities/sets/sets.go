package sets

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	croncapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/cron"
	logeventtriggercapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/logeventtrigger"
	mockcapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/mock"
	readcontractcapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/readcontract"
)

func NewDefaultSet(homeChainID uint64) ([]cre.InstallableCapability, error) {
	capabilities := []cre.InstallableCapability{}

	cron, cErr := croncapability.New()
	if cErr != nil {
		return nil, errors.Wrap(cErr, "failed to create cron capability")
	}
	capabilities = append(capabilities, cron)

	mock, mockErr := mockcapability.New()
	if mockErr != nil {
		return nil, errors.Wrap(mockErr, "failed to create mock capability")
	}
	capabilities = append(capabilities, mock)

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
