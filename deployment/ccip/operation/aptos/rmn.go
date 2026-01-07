package aptos

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/aptos-labs/aptos-go-sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/dependency"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
)

// allow test seams
var (
	ccipBindFn       = ccip.Bind
	generateMCMSTxFn = utils.GenerateMCMSTx
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

func curseMultiple(b operations.Bundle, deps dependency.AptosDeps, in CurseMultipleInput) (mcmstypes.Transaction, error) {
	// Bind CCIP Package
	ccipAddress := in.CCIPAddress
	ccipBind := ccipBindFn(ccipAddress, deps.AptosChain.Client)

	// Encode curse multiple operation
	moduleInfo, function, _, args, err := ccipBind.RMNRemote().Encoder().CurseMultiple(in.Subjects)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to encode CurseMultiple: %w", err)
	}

	// Generate MCMS transaction
	tx, err := generateMCMSTxFn(ccipAddress, moduleInfo, function, args)
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

func uncurseMultiple(b operations.Bundle, deps dependency.AptosDeps, in UncurseMultipleInput) (mcmstypes.Transaction, error) {
	// Bind CCIP Package
	ccipAddress := in.CCIPAddress
	ccipBind := ccipBindFn(ccipAddress, deps.AptosChain.Client)

	// Encode uncurse multiple operation
	moduleInfo, function, _, args, err := ccipBind.RMNRemote().Encoder().UncurseMultiple(in.Subjects)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to encode UncurseMultiple: %w", err)
	}

	// Generate MCMS transaction
	tx, err := generateMCMSTxFn(ccipAddress, moduleInfo, function, args)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to generate MCMS transaction: %w", err)
	}

	return tx, nil
}

// IsSubjectCursed checks whether the given subject (or a global curse) exists on the RMN Remote.
func IsSubjectCursed(deps dependency.AptosDeps, ccipAddress aptos.AccountAddress, subject []byte) (bool, error) {
	ccipBind := ccipBindFn(ccipAddress, deps.AptosChain.Client)
	callOpts := &bind.CallOpts{}
	cursed, err := ccipBind.RMNRemote().IsCursed(callOpts, subject)
	if err != nil {
		return false, fmt.Errorf("failed to check if subject is cursed: %w", err)
	}
	return cursed, nil
}
