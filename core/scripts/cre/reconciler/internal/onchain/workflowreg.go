package onchain

import (
	"context"
	"slices"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldflogger "github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	cre "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
)

func (d *Deployer) configureWorkflowReg(
	ctx context.Context,
	env *cldf.Environment,
	chainSelector uint64,
	desired *domain.DesiredState,
	state *domain.StateFile,
) error {
	d.log.Info().Msg("P6: Configuring Workflow Registry")

	wfRegAddr := state.GetAddress(keystone_changeset.WorkflowRegistry.String())
	if wfRegAddr == "" {
		return errors.New("workflow registry address not found in state")
	}

	var allowedDonIDs []uint64
	for _, don := range desired.DONs {
		if !slices.Contains(don.DONTypes, "workflow") {
			continue
		}
		if id, ok := state.DONIDs[don.Name]; ok {
			allowedDonIDs = append(allowedDonIDs, id)
		}
	}
	if len(allowedDonIDs) == 0 {
		return errors.New("no workflow DON IDs found in state")
	}

	cldfLogger, err := cldflogger.New()
	if err != nil {
		return errors.Wrap(err, "failed to create cldf logger")
	}

	workflowOwner, err := deployerAddress(d.deployerKey)
	if err != nil {
		return errors.Wrap(err, "failed to resolve deployer workflow owner address")
	}

	out, err := workflow.ConfigureWorkflowRegistry(ctx, d.log, cldfLogger, &cre.WorkflowRegistryInput{
		ContractAddress: common.HexToAddress(wfRegAddr),
		ContractVersion: cldf.TypeAndVersion{Version: *crecontracts.V2Version},
		ChainSelector:   chainSelector,
		CldEnv:          env,
		AllowedDonIDs:   allowedDonIDs,
		WorkflowOwners:  []common.Address{workflowOwner},
	})
	if err != nil {
		return errors.Wrap(err, "failed to configure workflow registry")
	}

	owners := make([]string, len(out.WorkflowOwners))
	for i, addr := range out.WorkflowOwners {
		owners[i] = addr.Hex()
	}
	state.WorkflowReg = &domain.WorkflowRegState{
		ChainSelector:  domain.ChainSelector(chainSelector),
		AllowedDonIDs:  out.AllowedDonIDs,
		WorkflowOwners: owners,
	}

	d.log.Info().
		Str("wfReg", wfRegAddr).
		Int("allowedDons", len(allowedDonIDs)).
		Msg("Workflow Registry configured")

	return nil
}
