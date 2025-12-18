package operation

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
)

var RMNRemoteOperations = []*operations.Operation[any, any, any]{
	CurseMultipleOp.AsUntyped(),
	UncurseMultipleOp.AsUntyped(),
}

// CurseMultipleInput is the input for cursing multiple subjects via RMN Remote.
type CurseMultipleInput struct {
	CCIPAddress aptos.AccountAddress
	Subjects    [][]byte
}

// CurseMultipleOp generates an MCMS transaction to curse multiple subjects.
var CurseMultipleOp = operations.NewOperation(
	"curse-multiple-op",
	Version1_0_0,
	"Generates MCMS transaction to curse multiple subjects on RMN Remote",
	curseMultiple,
)

func curseMultiple(b operations.Bundle, deps AptosDeps, in CurseMultipleInput) (mcmstypes.Transaction, error) {
	ccipBind := ccip.Bind(in.CCIPAddress, deps.AptosChain.Client)

	moduleInfo, function, _, args, err := ccipBind.RMNRemote().Encoder().CurseMultiple(in.Subjects)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to encode CurseMultiple: %w", err)
	}

	tx, err := utils.GenerateMCMSTx(in.CCIPAddress, moduleInfo, function, args)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to generate MCMS transaction: %w", err)
	}

	return tx, nil
}

// UncurseMultipleInput is the input for uncursing multiple subjects via RMN Remote.
type UncurseMultipleInput struct {
	CCIPAddress aptos.AccountAddress
	Subjects    [][]byte
}

// UncurseMultipleOp generates an MCMS transaction to uncurse multiple subjects.
var UncurseMultipleOp = operations.NewOperation(
	"uncurse-multiple-op",
	Version1_0_0,
	"Generates MCMS transaction to uncurse multiple subjects on RMN Remote",
	uncurseMultiple,
)

func uncurseMultiple(b operations.Bundle, deps AptosDeps, in UncurseMultipleInput) (mcmstypes.Transaction, error) {
	ccipBind := ccip.Bind(in.CCIPAddress, deps.AptosChain.Client)

	moduleInfo, function, _, args, err := ccipBind.RMNRemote().Encoder().UncurseMultiple(in.Subjects)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to encode UncurseMultiple: %w", err)
	}

	tx, err := utils.GenerateMCMSTx(in.CCIPAddress, moduleInfo, function, args)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to generate MCMS transaction: %w", err)
	}

	return tx, nil
}
