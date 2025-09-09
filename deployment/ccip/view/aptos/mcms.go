package aptos

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	aptosCommon "github.com/smartcontractkit/chainlink/deployment/common/view/aptos"
)

type MCMSWithTimelockView struct {
	aptosCommon.ContractMetaData

	Bypasser  mcmstypes.Config
	Proposer  mcmstypes.Config
	Canceller mcmstypes.Config

	TimelockMinDelay         uint64
	TimelockBlockedFunctions []TimelockBlockedFunction
}

type TimelockBlockedFunction struct {
	Target       string
	ModuleName   string
	FunctionName string
}

func GenerateMCMSWithTimelockView(chain cldf_aptos.Chain, mcmsAddress aptos.AccountAddress) (MCMSWithTimelockView, error) {
	boundMCMS := mcms.Bind(mcmsAddress, chain.Client)

	mcmsOwner, err := boundMCMS.MCMSAccount().Owner(nil)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to retrieve owner of MCMS: %w", err)
	}

	// Query configs
	configTransformer := aptosmcms.NewConfigTransformer()
	bypasserCfg, err := boundMCMS.MCMS().GetConfig(nil, aptosmcms.TimelockRoleBypasser.Byte())
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to query bypasser config: %w", err)
	}
	bypasserConfig, err := configTransformer.ToConfig(bypasserCfg)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to convert aptos bypasser config: %w", err)
	}
	proposerCfg, err := boundMCMS.MCMS().GetConfig(nil, aptosmcms.TimelockRoleProposer.Byte())
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to query proposer config: %w", err)
	}
	proposerConfig, err := configTransformer.ToConfig(proposerCfg)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to convert aptos proposer config: %w", err)
	}
	cancellerCfg, err := boundMCMS.MCMS().GetConfig(nil, aptosmcms.TimelockRoleCanceller.Byte())
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to query canceler config: %w", err)
	}
	cancellerConfig, err := configTransformer.ToConfig(cancellerCfg)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to convert aptos canceller config: %w", err)
	}

	// Timelock
	timelockMinDelay, err := boundMCMS.MCMS().TimelockMinDelay(nil)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to query timelock min delay: %w", err)
	}
	blockedFunctions, err := boundMCMS.MCMS().TimelockGetBlockedFunctions(nil)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to query timelock blocked functions: %w", err)
	}
	TimelockBlockedFunctions := make([]TimelockBlockedFunction, len(blockedFunctions))
	for i, function := range blockedFunctions {
		TimelockBlockedFunctions[i] = TimelockBlockedFunction{
			Target:       function.Target.StringLong(),
			ModuleName:   function.ModuleName,
			FunctionName: function.FunctionName,
		}
	}

	return MCMSWithTimelockView{
		ContractMetaData: aptosCommon.ContractMetaData{
			Address: mcmsAddress.StringLong(),
			Owner:   mcmsOwner.StringLong(),
		},
		Bypasser:  *bypasserConfig,
		Proposer:  *proposerConfig,
		Canceller: *cancellerConfig,

		TimelockMinDelay: timelockMinDelay,
	}, nil
}
