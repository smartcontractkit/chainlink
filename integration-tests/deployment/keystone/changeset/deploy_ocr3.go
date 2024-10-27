package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo/models"
	kslib "github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
)

func DeployOCR3(lggr logger.Logger, env deployment.Environment, ab deployment.AddressBook, registryChainSel uint64) (deployment.ChangesetOutput, error) {
	// ocr3 only deployed on registry chain
	c, ok := env.Chains[registryChainSel]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain not found in environment")
	}
	err := kslib.DeployOCR3(lggr, c, ab)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy OCR3Capability: %w", err)
	}
	return deployment.ChangesetOutput{AddressBook: ab}, nil
}

func ConfigureOCR3Contract(lggr logger.Logger, env deployment.Environment, ab deployment.AddressBook, registryChainSel uint64, nodes []*models.Node, cfg kslib.OracleConfigSource) (deployment.ChangesetOutput, error) {

	err := kslib.ConfigureOCR3ContractFromCLO(&env, registryChainSel, nodes, ab, &cfg)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to configure OCR3Capability: %w", err)
	}

	return deployment.ChangesetOutput{AddressBook: ab}, nil
}
