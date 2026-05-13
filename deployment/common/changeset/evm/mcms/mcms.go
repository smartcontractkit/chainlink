package mcmsnew

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cast"

	bindings "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset/evm/mcms/seqs"
	"github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
)

// TODO: delete this function after it is available in timelock Inspector
func getAdminAddresses(ctx context.Context, timelock *bindings.RBACTimelock) ([]string, error) {
	numAddresses, err := timelock.GetRoleMemberCount(&bind.CallOpts{
		Context: ctx,
	}, v1_0.ADMIN_ROLE.ID)
	if err != nil {
		return nil, err
	}
	adminAddresses := make([]string, 0, numAddresses.Uint64())
	for i := range numAddresses.Uint64() {
		if i > math.MaxUint32 {
			return nil, fmt.Errorf("value %d exceeds uint32 range", i)
		}
		idx, err := cast.ToInt64E(i)
		if err != nil {
			return nil, err
		}
		address, err := timelock.GetRoleMember(&bind.CallOpts{
			Context: ctx,
		}, v1_0.ADMIN_ROLE.ID, big.NewInt(idx))
		if err != nil {
			return nil, err
		}
		adminAddresses = append(adminAddresses, address.String())
	}
	return adminAddresses, nil
}

func GrantRolesForTimelock(
	env cldf.Environment,
	chain cldf_evm.Chain,
	timelockContracts *cldfproposalutils.MCMSWithTimelockContracts,
	skipIfDeployerKeyNotAdmin bool, // If true, skip role grants if the deployer key is not an admin.
	gasBoostConfig *cldfproposalutils.GasBoostConfig,
) (operations.SequenceReport[seqs.SeqGrantRolesTimelockInput, map[uint64][]opsutils.EVMCallOutput], error) {
	lggr := env.Logger
	ctx := env.GetContext()

	if timelockContracts == nil {
		lggr.Errorw("Timelock contracts not found", "chain", chain.String())
		return operations.SequenceReport[seqs.SeqGrantRolesTimelockInput, map[uint64][]opsutils.EVMCallOutput]{}, fmt.Errorf("timelock contracts not found for chain %s", chain.String())
	}

	timelock := timelockContracts.Timelock
	proposer := timelockContracts.ProposerMcm
	canceller := timelockContracts.CancellerMcm
	bypasser := timelockContracts.BypasserMcm
	callProxy := timelockContracts.CallProxy

	// get admin addresses
	adminAddresses, err := getAdminAddresses(ctx, timelock)
	if err != nil {
		return operations.SequenceReport[seqs.SeqGrantRolesTimelockInput, map[uint64][]opsutils.EVMCallOutput]{}, fmt.Errorf("failed to get admin addresses: %w", err)
	}
	isDeployerKeyAdmin := slices.Contains(adminAddresses, chain.DeployerKey.From.String())
	isTimelockAdmin := slices.Contains(adminAddresses, timelock.Address().String())
	if !isDeployerKeyAdmin && skipIfDeployerKeyNotAdmin {
		lggr.Infow("Deployer key is not admin, skipping role grants", "chain", chain.String())
		return operations.SequenceReport[seqs.SeqGrantRolesTimelockInput, map[uint64][]opsutils.EVMCallOutput]{}, nil
	}
	if !isDeployerKeyAdmin && !isTimelockAdmin {
		return operations.SequenceReport[seqs.SeqGrantRolesTimelockInput, map[uint64][]opsutils.EVMCallOutput]{}, errors.New("neither deployer key nor timelock is admin, cannot grant roles")
	}

	seqDeps := seqs.SeqGrantRolesTimelockDeps{
		Chain: chain,
	}

	seqInput := seqs.SeqGrantRolesTimelockInput{
		ContractType:       commontypes.RBACTimelock,
		ChainSelector:      chain.Selector,
		Timelock:           timelock.Address(),
		IsDeployerKeyAdmin: isDeployerKeyAdmin,
		RolesAndAddresses: []seqs.RolesAndAddresses{
			{
				Role:      v1_0.PROPOSER_ROLE.ID,
				Name:      v1_0.PROPOSER_ROLE.Name,
				Addresses: []common.Address{proposer.Address()},
			},
			{
				Role:      v1_0.CANCELLER_ROLE.ID,
				Name:      v1_0.CANCELLER_ROLE.Name,
				Addresses: []common.Address{proposer.Address(), canceller.Address(), bypasser.Address()},
			},
			{
				Role:      v1_0.BYPASSER_ROLE.ID,
				Name:      v1_0.BYPASSER_ROLE.Name,
				Addresses: []common.Address{bypasser.Address()},
			},
			{
				Role:      v1_0.EXECUTOR_ROLE.ID,
				Name:      v1_0.EXECUTOR_ROLE.Name,
				Addresses: []common.Address{callProxy.Address()},
			},
		},
		GasBoostConfig: gasBoostConfig,
	}

	if !isTimelockAdmin {
		// We grant the timelock the admin role on the MCMS contracts.
		seqInput.RolesAndAddresses = append(seqInput.RolesAndAddresses, seqs.RolesAndAddresses{
			Role:      v1_0.ADMIN_ROLE.ID,
			Name:      v1_0.ADMIN_ROLE.Name,
			Addresses: []common.Address{timelock.Address()},
		})
	}

	report, err := operations.ExecuteSequence(
		env.OperationsBundle,
		seqs.SeqGrantRolesTimelock,
		seqDeps,
		seqInput,
	)
	if err != nil {
		lggr.Errorw("Failed to grant roles for timelock", "chain", chain.String(), "err", err)
		return operations.SequenceReport[seqs.SeqGrantRolesTimelockInput, map[uint64][]opsutils.EVMCallOutput]{}, err
	}

	return report, nil
}
