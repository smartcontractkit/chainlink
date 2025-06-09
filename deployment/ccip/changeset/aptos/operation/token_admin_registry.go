package operation

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
)

type SetTokenRegistrarInput struct {
	TokenAddress          aptos.AccountAddress
	TokenPoolOwnerAddress aptos.AccountAddress
}

var SetTokenRegistrarOp = operations.NewOperation(
	"set-token-registrar-op",
	Version1_0_0,
	"Sets the token registrar for a given token",
	setTokenRegistrar,
)

func setTokenRegistrar(b operations.Bundle, deps AptosDeps, in SetTokenRegistrarInput) (types.Transaction, error) {
	// Bind CCIP Package
	ccipAddress := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].CCIPAddress
	ccipBind := ccip.Bind(ccipAddress, deps.AptosChain.Client)

	moduleInfo, function, _, args, err := ccipBind.TokenAdminRegistry().Encoder().SetTokenRegistrar(in.TokenAddress, in.TokenPoolOwnerAddress)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode SetTokenRegistrar: %w", err)
	}
	tx, err := utils.GenerateMCMSTx(ccipAddress, moduleInfo, function, args)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	return tx, nil
}
