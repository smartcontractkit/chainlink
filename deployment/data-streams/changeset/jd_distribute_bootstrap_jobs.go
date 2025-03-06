package changeset

import (
	"context"
	"fmt"
	"time"

	chainsel "github.com/smartcontractkit/chain-selectors"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jd"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jobs"
)

const (
	defaultBootstrapJobSpecsTimeout = 120 * time.Second
)

var _ deployment.ChangeSetV2[CsDistributeBootstrapJobSpecsConfig] = CsDistributeBootstrapJobSpecs{}

type CsDistributeBootstrapJobSpecsConfig struct {
	ChainSelectorEVM uint64
	Filter           *jd.ListFilter
}

func findConfiguratorAddressByDON(addresses map[string]deployment.TypeAndVersion, donID uint64) (string, error) {
	for address, contract := range addresses {
		if contract.Type == "Configurator" && contract.Labels.Contains(fmt.Sprintf("don-%d", donID)) {
			return address, nil
		}
	}
	return "", fmt.Errorf("Configurator contract not found for DON %d", donID)
}

type CsDistributeBootstrapJobSpecs struct{}

func (CsDistributeBootstrapJobSpecs) Apply(e deployment.Environment, cfg CsDistributeBootstrapJobSpecsConfig) (deployment.ChangesetOutput, error) {
	ctx, cancel := context.WithTimeout(e.GetContext(), defaultBootstrapJobSpecsTimeout)
	defer cancel()

	bootstrapNodes, err := jd.FetchDONBootstrappersFromJD(ctx, e.Offchain, cfg.Filter)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get workflow don nodes: %w", err)
	}

	chainID, err := chainsel.GetChainIDFromSelector(cfg.ChainSelectorEVM)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get chain ID from selector: %w", err)
	}

	addresses, err := e.ExistingAddresses.AddressesForChain(cfg.ChainSelectorEVM)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get existing addresses: %w", err)
	}

	// Search for the Configurator address in the address book by DON ID
	configuratorAddress, err := findConfiguratorAddressByDON(addresses, cfg.Filter.DONID)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	bootstrapSpec := jobs.NewBootstrapSpec(
		configuratorAddress,
		cfg.Filter.DONID,
		jobs.RelayTypeEVM,
		jobs.RelayConfig{
			ChainID: chainID,
		},
	)

	renderedSpec, err := bootstrapSpec.MarshalTOML()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to marshal bootstrap spec: %w", err)
	}

	// Add a label to the job spec to identify the related DON
	labels := append([]*ptypes.Label(nil),
		&ptypes.Label{
			Key: fmt.Sprintf("don-%d-%s", cfg.Filter.DONID, cfg.Filter.DONName),
		})

	// TODO: For now the implementation uses a very simple approach, in case of partial failures there is a risk
	// of sending the same job spec multiple times to the same node. We need to understand the implications of this
	// and decide if we need to implement a more complex approach.
	for _, node := range bootstrapNodes {
		_, err = e.Offchain.ProposeJob(ctx, &jobv1.ProposeJobRequest{
			NodeId: node.Id,
			Spec:   string(renderedSpec),
			Labels: labels,
		})

		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to propose job: %w", err)
		}
	}

	return deployment.ChangesetOutput{}, nil
}

func (f CsDistributeBootstrapJobSpecs) VerifyPreconditions(e deployment.Environment, config CsDistributeBootstrapJobSpecsConfig) error {
	return nil
}
