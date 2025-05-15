package operation

import (
	"context"
	"fmt"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	"github.com/smartcontractkit/mcms/types"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	mcmsbind "github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// OP: Deploy MCMS Contract
var DeployMCMSOp = operations.NewOperation(
	"deploy-mcms-op",
	Version1_0_0,
	"Deploys MCMS Contract Operation for Aptos Chain",
	deployMCMS,
)

func deployMCMS(b operations.Bundle, deps AptosDeps, _ operations.EmptyInput) (aptos.AccountAddress, error) {
	mcmsSeed := mcmsbind.DefaultSeed + time.Now().String()
	mcmsAddress, mcmsDeployTx, _, err := mcmsbind.DeployToResourceAccount(deps.AptosChain.DeployerSigner, deps.AptosChain.Client, mcmsSeed)
	if err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("failed to deploy MCMS contract: %w", err)
	}
	if err := deps.AptosChain.Confirm(mcmsDeployTx.Hash); err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("failed to confirm MCMS deployment transaction: %w", err)
	}

	return mcmsAddress, nil
}

// OP: Configure MCMS Contract
type ConfigureMCMSInput struct {
	MCMSAddress aptos.AccountAddress
	MCMSConfigs mcmstypes.Config
	MCMSRole    aptosmcms.TimelockRole
}

var ConfigureMCMSOp = operations.NewOperation(
	"configure-mcms-op",
	Version1_0_0,
	"Configure MCMS Contract Operation for Aptos Chain",
	configureMCMS,
)

func configureMCMS(b operations.Bundle, deps AptosDeps, in ConfigureMCMSInput) (any, error) {
	configurer := aptosmcms.NewConfigurer(deps.AptosChain.Client, deps.AptosChain.DeployerSigner, in.MCMSRole)
	setCfgTx, err := configurer.SetConfig(context.Background(), in.MCMSAddress.StringLong(), &in.MCMSConfigs, false)
	if err != nil {
		return nil, fmt.Errorf("failed to setConfig in MCMS contract: %w", err)
	}
	if err := deps.AptosChain.Confirm(setCfgTx.Hash); err != nil {
		return nil, fmt.Errorf("MCMS setConfig transaction failed: %w", err)
	}
	return nil, nil
}

// OP: Transfer Ownership to Self
var TransferOwnershipToSelfOp = operations.NewOperation(
	"transfer-ownership-to-self-op",
	Version1_0_0,
	"Transfer ownership to self",
	transferOwnershipToSelf,
)

func transferOwnershipToSelf(b operations.Bundle, deps AptosDeps, mcmsAddress aptos.AccountAddress) (any, error) {
	opts := &bind.TransactOpts{Signer: deps.AptosChain.DeployerSigner}
	contractMCMS := mcmsbind.Bind(mcmsAddress, deps.AptosChain.Client)
	tx, err := contractMCMS.MCMSAccount().TransferOwnershipToSelf(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to TransferOwnershipToSelf in MCMS contract: %w", err)
	}
	err = deps.AptosChain.Confirm(tx.Hash)
	if err != nil {
		return nil, fmt.Errorf("MCMS TransferOwnershipToSelf transaction failed: %w", err)
	}
	return nil, nil
}

// OP: Accept Ownership
var AcceptOwnershipOp = operations.NewOperation(
	"accept-ownership-op",
	Version1_0_0,
	"Generate Accept Ownership BatchOperations for MCMS Contract",
	acceptOwnership,
)

func acceptOwnership(b operations.Bundle, deps AptosDeps, mcmsAddress aptos.AccountAddress) (mcmstypes.Transaction, error) {
	contractMCMS := mcmsbind.Bind(mcmsAddress, deps.AptosChain.Client)
	moduleInfo, function, _, args, err := contractMCMS.MCMSAccount().Encoder().AcceptOwnership()
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to encode AcceptOwnership: %w", err)
	}

	return utils.GenerateMCMSTx(mcmsAddress, moduleInfo, function, args)
}

// OP: SetMinDelay
type TimelockMinDelayInput struct {
	MCMSAddress      aptos.AccountAddress
	TimelockMinDelay uint64
}

var SetMinDelayOP = operations.NewOperation(
	"set-timelock-min-delay-op",
	Version1_0_0,
	"Generate set timelock min delay MCMS BatchOperations",
	setMinDelay,
)

func setMinDelay(b operations.Bundle, deps AptosDeps, in TimelockMinDelayInput) (mcmstypes.Transaction, error) {
	contractMCMS := mcmsbind.Bind(in.MCMSAddress, deps.AptosChain.Client)
	moduleInfo, function, _, args, err := contractMCMS.MCMS().Encoder().TimelockUpdateMinDelay(in.TimelockMinDelay)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to encode AcceptOwnership: %w", err)
	}

	return utils.GenerateMCMSTx(in.MCMSAddress, moduleInfo, function, args)
}

// OP: CleanupStagingAreaOp generates a batch operation to clean up the staging area
var CleanupStagingAreaOp = operations.NewOperation(
	"cleanup-staging-area-op",
	Version1_0_0,
	"Cleans up MCMS staging area if it's not already clean",
	cleanupStagingArea,
)

func cleanupStagingArea(b operations.Bundle, deps AptosDeps, mcmsAddress aptos.AccountAddress) (types.BatchOperation, error) {
	// Check resources first to see if staging is clean
	IsMCMSStagingAreaClean, err := utils.IsMCMSStagingAreaClean(deps.AptosChain.Client, mcmsAddress)
	if err != nil {
		return types.BatchOperation{}, fmt.Errorf("failed to check if MCMS staging area is clean: %w", err)
	}
	if IsMCMSStagingAreaClean {
		b.Logger.Infow("MCMS Staging Area already clean", "addr", mcmsAddress.String())
		return types.BatchOperation{}, nil
	}

	// Bind MCMS contract
	mcmsContract := mcmsbind.Bind(mcmsAddress, deps.AptosChain.Client)

	// Get cleanup staging operations
	moduleInfo, function, _, args, err := mcmsContract.MCMSDeployer().Encoder().CleanupStagingArea()
	if err != nil {
		return types.BatchOperation{}, fmt.Errorf("failed to EncodeCleanupStagingArea: %w", err)
	}
	mcmsTx, err := utils.GenerateMCMSTx(mcmsAddress, moduleInfo, function, args)
	if err != nil {
		return types.BatchOperation{}, fmt.Errorf("failed to generate MCMS operations for FeeQuoter Initialize: %w", err)
	}

	return types.BatchOperation{
		ChainSelector: types.ChainSelector(deps.AptosChain.Selector),
		Transactions:  []types.Transaction{mcmsTx},
	}, nil
}
