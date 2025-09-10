package sets

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	customcompute "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/compute"
	consensusv1capability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/consensus/v1"
	consensusv2capability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/consensus/v2"
	croncapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/cron"
	evmcapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/evm"
	httpactioncapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/httpaction"
	httptriggercapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/httptrigger"
	logeventtriggercapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/logeventtrigger"
	mockcapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/mock"
	readcontractcapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/readcontract"
	vaultcapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/vault"
	webapitargetcapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/webapitarget"
	webapitriggercapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/webapitrigger"
	writeevmcapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/writeevm"
	writesolanacapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/writesolana"
)

func NewDefaultSet(homeChainID uint64) ([]cre.InstallableCapability, error) {
	capabilities := []cre.InstallableCapability{}

	cron, cErr := croncapability.New()
	if cErr != nil {
		return nil, errors.Wrap(cErr, "failed to create cron capability")
	}
	capabilities = append(capabilities, cron)

	customCompute, customComputeErr := customcompute.New()
	if customComputeErr != nil {
		return nil, errors.Wrap(customComputeErr, "failed to create custom compute capability")
	}
	capabilities = append(capabilities, customCompute)

	c1, c1Err := consensusv1capability.New(homeChainID)
	if c1Err != nil {
		return nil, errors.Wrap(c1Err, "failed to create consensus capability v1")
	}
	capabilities = append(capabilities, c1)

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

	httpaction, httpactionErr := httpactioncapability.New()
	if httpactionErr != nil {
		return nil, errors.Wrap(httpactionErr, "failed to create http action capability")
	}
	capabilities = append(capabilities, httpaction)

	httptrigger, httptriggerErr := httptriggercapability.New()
	if httptriggerErr != nil {
		return nil, errors.Wrap(httptriggerErr, "failed to create http trigger capability")
	}
	capabilities = append(capabilities, httptrigger)

	webapitrigger, webapitriggerErr := webapitriggercapability.New()
	if webapitriggerErr != nil {
		return nil, errors.Wrap(webapitriggerErr, "failed to create web api trigger capability")
	}
	capabilities = append(capabilities, webapitrigger)

	webapitarget, webapitargetErr := webapitargetcapability.New()
	if webapitargetErr != nil {
		return nil, errors.Wrap(webapitargetErr, "failed to create web api target capability")
	}
	capabilities = append(capabilities, webapitarget)

	vault, vaultErr := vaultcapability.New(homeChainID)
	if vaultErr != nil {
		return nil, errors.Wrap(vaultErr, "failed to create vault capability")
	}
	capabilities = append(capabilities, vault)

	mock, mockErr := mockcapability.New()
	if mockErr != nil {
		return nil, errors.Wrap(mockErr, "failed to create mock capability")
	}
	capabilities = append(capabilities, mock)

	writeevm, writeevmErr := writeevmcapability.New()
	if writeevmErr != nil {
		return nil, errors.Wrap(writeevmErr, "failed to create write evm capability")
	}
	capabilities = append(capabilities, writeevm)

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
