package operation

import (
	"fmt"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
)

// export operations like other places
var RMNRemoteOperations = []*operations.Operation[any, any, any]{
	CurseSubjectsOp.AsUntyped(),
	UncurseSubjectsOp.AsUntyped(),
	GlobalCurseOp.AsUntyped(),
	GlobalUncurseOp.AsUntyped(),
}

// CurseSubjectsInput contains configuration for cursing subjects on RMN Remote
type CurseSubjectsInput struct {
	Subjects [][]byte
}

// CurseSubjectsOp operation to curse subjects on RMN Remote
var CurseSubjectsOp = operations.NewOperation(
	"curse-subjects-op",
	Version1_0_0,
	"Curses subjects on RMN Remote",
	curseSubjects,
)

func curseSubjects(b operations.Bundle, deps AptosDeps, in CurseSubjectsInput) ([]mcmstypes.Transaction, error) {
	// Bind CCIP Package
	ccipAddress := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].CCIPAddress
	ccipBind := ccip.Bind(ccipAddress, deps.AptosChain.Client)

	var txs []mcmstypes.Transaction

	if len(in.Subjects) == 0 {
		b.Logger.Infow("No subjects to curse")
		return nil, nil
	}

	// Encode the curse operation
	moduleInfo, function, _, args, err := ccipBind.RMNRemote().Encoder().CurseMultiple(in.Subjects)
	if err != nil {
		return nil, fmt.Errorf("failed to encode CurseMultiple: %w", err)
	}

	tx, err := utils.GenerateMCMSTx(ccipAddress, moduleInfo, function, args)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	txs = append(txs, tx)

	b.Logger.Infow("Adding RMN Remote curse subjects operation",
		"subjectCount", len(in.Subjects))

	return txs, nil
}

// UncurseSubjectsInput contains configuration for uncursing subjects on RMN Remote
type UncurseSubjectsInput struct {
	Subjects [][]byte
}

// UncurseSubjectsOp operation to uncurse subjects on RMN Remote
var UncurseSubjectsOp = operations.NewOperation(
	"uncurse-subjects-op",
	Version1_0_0,
	"Uncurses subjects on RMN Remote",
	uncurseSubjects,
)

func uncurseSubjects(b operations.Bundle, deps AptosDeps, in UncurseSubjectsInput) ([]mcmstypes.Transaction, error) {
	// Bind CCIP Package
	ccipAddress := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].CCIPAddress
	ccipBind := ccip.Bind(ccipAddress, deps.AptosChain.Client)

	var txs []mcmstypes.Transaction

	if len(in.Subjects) == 0 {
		b.Logger.Infow("No subjects to uncurse")
		return nil, nil
	}

	// Encode the uncurse operation
	moduleInfo, function, _, args, err := ccipBind.RMNRemote().Encoder().UncurseMultiple(in.Subjects)
	if err != nil {
		return nil, fmt.Errorf("failed to encode UncurseMultiple: %w", err)
	}

	tx, err := utils.GenerateMCMSTx(ccipAddress, moduleInfo, function, args)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	txs = append(txs, tx)

	b.Logger.Infow("Adding RMN Remote uncurse subjects operation",
		"subjectCount", len(in.Subjects))

	return txs, nil
}

// GlobalCurseInput contains configuration for global curse on RMN Remote
type GlobalCurseInput struct {
	// No additional parameters needed for global curse
}

// GlobalCurseOp operation to perform global curse on RMN Remote
var GlobalCurseOp = operations.NewOperation(
	"global-curse-op",
	Version1_0_0,
	"Performs global curse on RMN Remote",
	globalCurse,
)

func globalCurse(b operations.Bundle, deps AptosDeps, in GlobalCurseInput) ([]mcmstypes.Transaction, error) {
	// Bind CCIP Package
	ccipAddress := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].CCIPAddress
	ccipBind := ccip.Bind(ccipAddress, deps.AptosChain.Client)

	var txs []mcmstypes.Transaction

	// Global curse uses the global curse subject from globals package
	globalSubject := globals.GlobalCurseSubject()
	subjects := [][]byte{globalSubject[:]}

	// Encode the global curse operation
	moduleInfo, function, _, args, err := ccipBind.RMNRemote().Encoder().CurseMultiple(subjects)
	if err != nil {
		return nil, fmt.Errorf("failed to encode global curse: %w", err)
	}

	tx, err := utils.GenerateMCMSTx(ccipAddress, moduleInfo, function, args)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	txs = append(txs, tx)

	b.Logger.Infow("Adding RMN Remote global curse operation")

	return txs, nil
}

// GlobalUncurseInput contains configuration for global uncurse on RMN Remote
type GlobalUncurseInput struct {
	// No additional parameters needed for global uncurse
}

// GlobalUncurseOp operation to perform global uncurse on RMN Remote
var GlobalUncurseOp = operations.NewOperation(
	"global-uncurse-op",
	Version1_0_0,
	"Performs global uncurse on RMN Remote",
	globalUncurse,
)

func globalUncurse(b operations.Bundle, deps AptosDeps, in GlobalUncurseInput) ([]mcmstypes.Transaction, error) {
	// Bind CCIP Package
	ccipAddress := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].CCIPAddress
	ccipBind := ccip.Bind(ccipAddress, deps.AptosChain.Client)

	var txs []mcmstypes.Transaction

	// Global uncurse uses the same global curse subject from globals package
	globalSubject := globals.GlobalCurseSubject()
	subjects := [][]byte{globalSubject[:]}

	// Encode the global uncurse operation
	moduleInfo, function, _, args, err := ccipBind.RMNRemote().Encoder().UncurseMultiple(subjects)
	if err != nil {
		return nil, fmt.Errorf("failed to encode global uncurse: %w", err)
	}

	tx, err := utils.GenerateMCMSTx(ccipAddress, moduleInfo, function, args)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	txs = append(txs, tx)

	b.Logger.Infow("Adding RMN Remote global uncurse operation")

	return txs, nil
}
