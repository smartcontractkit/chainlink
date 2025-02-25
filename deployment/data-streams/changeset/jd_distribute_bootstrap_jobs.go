package changeset

import (
	"context"
	"fmt"
	"time"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jd"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jobs"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
)

const (
	defaultBootstrapJobSpecsTimeout = 120 * time.Second
)

var CsDistributeBootstrapJobSpecs deployment.ChangeLogic[CsDistributeBootstrapJobSpecsConfig] = csDistributeBootstrapJobSpecs

type ContractID string

type CsDistributeBootstrapJobSpecsConfig struct {
	ChainSelectorEVM uint64
	DONFilter        *jd.DONFilter
}

func findConfiguratorAddressByDON(addresses map[ContractID]deployment.TypeAndVersion, DONID uint64) (string, error) {
	for address, contract := range addresses {
		if contract.Type == "Configurator" && contract.Labels.Contains(fmt.Sprintf("don-%d", DONID)) {
			return string(address), nil
		}
	}
	return "", fmt.Errorf("Configurator contract not found for DON %d", DONID)
}

func csDistributeBootstrapJobSpecs(e deployment.Environment, cfg CsDistributeBootstrapJobSpecsConfig) (deployment.ChangesetOutput, error) {
	ctx, cancel := context.WithTimeout(e.GetContext(), defaultBootstrapJobSpecsTimeout)
	defer cancel()

	bootstrapNodes, err := jd.FetchDONBootstrappersFromJD(ctx, e.Offchain, cfg.DONFilter)
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

	configuratorAddress, err := findConfiguratorAddressByDON(addresses, cfg.DONFilter.DONID)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	bootstrapSpec := jobs.NewBootstrapSpec(
		configuratorAddress,
		cfg.DONFilter.DONID,
		jobs.RelayTypeEVM,
		jobs.RelayConfig{
			ChainID: chainID,
		},
	)

	renderedSpec, err := bootstrapSpec.MarshalTOML()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to marshal bootstrap spec: %w", err)
	}

	// TODO: For now the implementation uses a very simple approach, in case of partial failures there is a risk
	// of sending the same job spec multiple times to the same node. We need to understand the implications of this
	// and decide if we need to implement a more complex approach.
	for _, node := range bootstrapNodes {
		err = e.Offchain.ProposeJob(ctx, jobv1.ProposeJobRequest{
			Id:   node.NodeID.String(),
			Spec: renderedSpec,
		})

		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to propose job: %w", err)
		}
	}

	return deployment.ChangesetOutput{}, nil
}
