package sequence

import (
	"fmt"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

// Set OCR3 Offramp Sequence Input
type SetOCR3OfframpSeqInput struct {
	HomeChainSelector uint64
	ChainSelector     uint64
}

var SetOCR3OfframpSequence = operations.NewSequence(
	"set-aptos-ocr3-offramp-sequence",
	operation.Version1_0_0,
	"Set OCR3 configuration for Aptos CCIP Offramp",
	setOCR3OfframpSequence,
)

func setOCR3OfframpSequence(b operations.Bundle, deps operation.AptosDeps, in SetOCR3OfframpSeqInput) (mcmstypes.BatchOperation, error) {
	var txs []mcmstypes.Transaction

	donID, err := internal.DonIDForChain(
		deps.CCIPOnChainState.Chains[in.HomeChainSelector].CapabilityRegistry,
		deps.CCIPOnChainState.Chains[in.HomeChainSelector].CCIPHome,
		in.ChainSelector,
	)
	if err != nil {
		return mcmstypes.BatchOperation{}, fmt.Errorf("failed to get DON ID: %w", err)
	}

	ocr3Args, err := internal.BuildSetOCR3ConfigArgsAptos(
		donID,
		deps.CCIPOnChainState.Chains[in.HomeChainSelector].CCIPHome,
		in.ChainSelector,
		globals.ConfigTypeActive,
	)
	if err != nil {
		return mcmstypes.BatchOperation{}, fmt.Errorf("failed to build OCR3 config args: %w", err)
	}

	var commitArgs *internal.MultiOCR3BaseOCRConfigArgsAptos
	var execArgs *internal.MultiOCR3BaseOCRConfigArgsAptos
	for _, ocr3Arg := range ocr3Args {
		switch ocr3Arg.OcrPluginType {
		case uint8(types.PluginTypeCCIPCommit):
			commitArgs = &ocr3Arg
		case uint8(types.PluginTypeCCIPExec):
			execArgs = &ocr3Arg
		default:
			return mcmstypes.BatchOperation{}, fmt.Errorf("unknown plugin type %d", ocr3Arg.OcrPluginType)
		}
	}

	// Set commit OCR3 Config
	commitReport, err := operations.ExecuteOperation(
		b,
		operation.SetOcr3ConfigOp,
		deps,
		operation.SetOcr3ConfigInput{
			OcrPluginType: types.PluginTypeCCIPCommit,
			OCRConfigArgs: *commitArgs,
		},
	)
	if err != nil {
		return mcmstypes.BatchOperation{}, err
	}
	txs = append(txs, commitReport.Output)

	// Set exec OCR3 Config
	execReport, err := operations.ExecuteOperation(
		b,
		operation.SetOcr3ConfigOp,
		deps,
		operation.SetOcr3ConfigInput{
			OcrPluginType: types.PluginTypeCCIPExec,
			OCRConfigArgs: *execArgs,
		},
	)
	if err != nil {
		return mcmstypes.BatchOperation{}, err
	}
	txs = append(txs, execReport.Output)

	return mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions:  txs,
	}, nil
}
