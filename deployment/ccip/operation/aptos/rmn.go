package aptos

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/aptos-labs/aptos-go-sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
)

// CurseMultipleInput is the input for cursing multiple subjects
type CurseMultipleInput struct {
	CCIPAddress aptos.AccountAddress
	Subjects    [][]byte
}

// OP: CurseMultipleOp generates MCMS transaction to curse multiple subjects
var CurseMultipleOp = operations.NewOperation(
	"curse-multiple-op",
	semver.MustParse("1.0.0"),
	"Generates MCMS transaction to curse multiple subjects on RMN Remote",
	curseMultiple,
)

func curseMultiple(b operations.Bundle, aptosChain cldf_aptos.Chain, in CurseMultipleInput) (mcmstypes.Transaction, error) {
	// Bind CCIP Package
	ccipBind := ccip.Bind(in.CCIPAddress, aptosChain.Client)

	// Encode curse multiple operation
	moduleInfo, function, _, args, err := ccipBind.RMNRemote().Encoder().CurseMultiple(in.Subjects)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to encode CurseMultiple: %w", err)
	}

	// Generate MCMS transaction
	tx, err := utils.GenerateMCMSTx(in.CCIPAddress, moduleInfo, function, args)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to generate MCMS transaction: %w", err)
	}

	return tx, nil
}

// UncurseMultipleInput is the input for uncursing multiple subjects
type UncurseMultipleInput struct {
	CCIPAddress aptos.AccountAddress
	Subjects    [][]byte
}

// OP: UncurseMultipleOp generates MCMS transaction to uncurse multiple subjects
var UncurseMultipleOp = operations.NewOperation(
	"uncurse-multiple-op",
	semver.MustParse("1.0.0"),
	"Generates MCMS transaction to uncurse multiple subjects on RMN Remote",
	uncurseMultiple,
)

func uncurseMultiple(b operations.Bundle, aptosChain cldf_aptos.Chain, in UncurseMultipleInput) (mcmstypes.Transaction, error) {
	// Bind CCIP Package
	ccipBind := ccip.Bind(in.CCIPAddress, aptosChain.Client)

	// Encode uncurse multiple operation
	moduleInfo, function, _, args, err := ccipBind.RMNRemote().Encoder().UncurseMultiple(in.Subjects)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to encode UncurseMultiple: %w", err)
	}

	// Generate MCMS transaction
	tx, err := utils.GenerateMCMSTx(in.CCIPAddress, moduleInfo, function, args)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to generate MCMS transaction: %w", err)
	}

	return tx, nil
}
