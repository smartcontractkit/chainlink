package operation

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
)

type ProposeAdministratorInput struct {
	TokenAddress       aptos.AccountAddress
	TokenAdministrator aptos.AccountAddress
}

var ProposeAdministratorOp = operations.NewOperation(
	"propose-administrator-op",
	Version1_0_0,
	"Proposes a new administrator for a given token",
	proposeAdministrator,
)

func proposeAdministrator(b operations.Bundle, deps AptosDeps, in ProposeAdministratorInput) (types.Transaction, error) {
	// Bind CCIP Package
	ccipAddress := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].CCIPAddress
	ccipBind := ccip.Bind(ccipAddress, deps.AptosChain.Client)

	moduleInfo, function, _, args, err := ccipBind.TokenAdminRegistry().Encoder().ProposeAdministrator(in.TokenAddress, in.TokenAdministrator)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode ProposeAdministrator: %w", err)
	}
	tx, err := utils.GenerateMCMSTx(ccipAddress, moduleInfo, function, args)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	return tx, nil
}

var AcceptAdminRoleOp = operations.NewOperation(
	"accept-admin-role-op",
	Version1_0_0,
	"Accepts a new administrator for a given token",
	acceptAdminRole,
)

func acceptAdminRole(b operations.Bundle, deps AptosDeps, tokenAddress aptos.AccountAddress) (types.Transaction, error) {
	// Bind CCIP Package
	ccipAddress := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].CCIPAddress
	ccipBind := ccip.Bind(ccipAddress, deps.AptosChain.Client)

	moduleInfo, function, _, args, err := ccipBind.TokenAdminRegistry().Encoder().AcceptAdminRole(tokenAddress)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode AcceptAdminRole: %w", err)
	}
	tx, err := utils.GenerateMCMSTx(ccipAddress, moduleInfo, function, args)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	return tx, nil
}

type SetPoolInput struct {
	TokenAddress     aptos.AccountAddress
	TokenPoolAddress aptos.AccountAddress
}

var SetPoolOp = operations.NewOperation(
	"set-pool-op",
	Version1_0_0,
	"Sets the pool for a given token",
	setPool,
)

func setPool(b operations.Bundle, deps AptosDeps, in SetPoolInput) (types.Transaction, error) {
	// Bind CCIP Package
	ccipAddress := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].CCIPAddress
	ccipBind := ccip.Bind(ccipAddress, deps.AptosChain.Client)

	moduleInfo, function, _, args, err := ccipBind.TokenAdminRegistry().Encoder().SetPool(in.TokenAddress, in.TokenPoolAddress)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode SetPool: %w", err)
	}
	tx, err := utils.GenerateMCMSTx(ccipAddress, moduleInfo, function, args)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	return tx, nil
}
