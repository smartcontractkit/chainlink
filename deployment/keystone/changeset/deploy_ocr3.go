package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
)

func DeployOCR3(env deployment.Environment, config interface{}) (deployment.ChangesetOutput, error) {
	lggr := env.Logger
	registryChainSel, ok := config.(uint64)
	if !ok {
		return deployment.ChangesetOutput{}, deployment.ErrInvalidConfig
	}
	ab := deployment.NewMemoryAddressBook()
	// ocr3 only deployed on registry chain
	c, ok := env.Chains[registryChainSel]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain not found in environment")
	}
	ocr3Resp, err := kslib.DeployOCR3(c, ab)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy OCR3Capability: %w", err)
	}
	lggr.Infof("Deployed %s chain selector %d addr %s", ocr3Resp.Tv.String(), c.Selector, ocr3Resp.Address.String())
	return deployment.ChangesetOutput{AddressBook: ab}, nil
}

func ConfigureOCR3Contract(lggr logger.Logger, env deployment.Environment, cfg kslib.ConfigureOCR3Config) (deployment.ChangesetOutput, error) {

	_, err := kslib.ConfigureOCR3ContractFromJD(&env, cfg)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to configure OCR3Capability: %w", err)
	}
	// does not create any new addresses
	return deployment.ChangesetOutput{}, nil
}
