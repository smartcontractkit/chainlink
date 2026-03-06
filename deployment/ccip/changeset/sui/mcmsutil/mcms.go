package mcmsutil

import (
	"fmt"

	"github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	suistate "github.com/smartcontractkit/chainlink-sui/deployment"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

func GetSuiRoleFromAction(action mcmstypes.TimelockAction) (suisdk.TimelockRole, error) {
	switch action {
	case mcmstypes.TimelockActionSchedule:
		return suisdk.TimelockRoleProposer, nil
	case mcmstypes.TimelockActionBypass:
		return suisdk.TimelockRoleBypasser, nil
	case mcmstypes.TimelockActionCancel:
		return suisdk.TimelockRoleCanceller, nil
	case "":
		return suisdk.TimelockRoleProposer, nil
	default:
		return 0, fmt.Errorf("unsupported timelock action: %v", action)
	}
}

func GenerateProposal(
	env cldf.Environment,
	state suistate.CCIPChainState,
	chainSel uint64,
	operations []mcmstypes.BatchOperation,
	description string,
	mcmsCfg proposalutils.TimelockConfig,
) (*mcms.TimelockProposal, error) {
	role, err := GetSuiRoleFromAction(mcmsCfg.MCMSAction)
	if err != nil {
		return nil, fmt.Errorf("failed to get SUI role from action: %w", err)
	}

	devInspectSigner := signer.NewDevInspectSigner("0x0")
	inspector, err := suisdk.NewInspector(
		env.BlockChains.SuiChains()[chainSel].Client,
		devInspectSigner,
		state.MCMSPackageID,
		role,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create SUI MCMS inspector for chain %d: %w", chainSel, err)
	}

	return proposalutils.BuildProposalFromBatchesV2(
		env,
		map[uint64]string{chainSel: state.MCMSTimelockObjectID},
		map[uint64]string{chainSel: state.MCMSPackageID},
		map[uint64]mcmssdk.Inspector{chainSel: inspector},
		operations,
		description,
		mcmsCfg,
	)
}
