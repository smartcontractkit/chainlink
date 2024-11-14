package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
)

func DeployOCR3(env deployment.Environment, config interface{}) (deployment.ChangesetOutput, error) {
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
	err := kslib.DeployOCR3(env.Logger, c, ab)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy OCR3Capability: %w", err)
	}
	return deployment.ChangesetOutput{AddressBook: ab}, nil
}

func ConfigureOCR3Contract(lggr logger.Logger, env deployment.Environment, ab deployment.AddressBook, registryChainSel uint64, nodes []string, cfg kslib.OracleConfigWithSecrets) (deployment.ChangesetOutput, error) {
	err := kslib.ConfigureOCR3ContractFromJD(&env, registryChainSel, nodes, ab, &cfg)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to configure OCR3Capability: %w", err)
	}

	return deployment.ChangesetOutput{AddressBook: ab}, nil
}
