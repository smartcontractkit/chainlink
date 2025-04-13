package sequence

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	aptosutils "github.com/smartcontractkit/chainlink-aptos/relayer/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/operations"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	mcmstypes "github.com/smartcontractkit/mcms/types"
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

	offRampBind := ccip_offramp.Bind(deps.OnChainState.CCIPAddress, deps.AptosChain.Client)

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

	var commitArgs *internal.MultiOCR3BaseOCRConfigArgsAptos = nil
	var execArgs *internal.MultiOCR3BaseOCRConfigArgsAptos = nil
	for _, ocr3Arg := range ocr3Args {
		if ocr3Arg.OcrPluginType == uint8(types.PluginTypeCCIPCommit) {
			commitArgs = &ocr3Arg
		} else if ocr3Arg.OcrPluginType == uint8(types.PluginTypeCCIPExec) {
			execArgs = &ocr3Arg
		} else {
			return mcmstypes.BatchOperation{}, fmt.Errorf("unknown plugin type %d", ocr3Arg.OcrPluginType)
		}
	}

	commitSigners := [][]byte{}
	for _, signer := range commitArgs.Signers {
		commitSigners = append(commitSigners, signer)
	}
	commitTransmitters := []aptos.AccountAddress{}
	for _, transmitter := range commitArgs.Transmitters {
		address, err := aptosutils.PublicKeyBytesToAddress(transmitter)
		if err != nil {
			return mcmstypes.BatchOperation{}, fmt.Errorf("failed to convert transmitter to address: %w", err)
		}
		commitTransmitters = append(commitTransmitters, address)
	}
	moduleInfo, function, _, args, err := offRampBind.Offramp().Encoder().SetOcr3Config(
		commitArgs.ConfigDigest[:],
		uint8(types.PluginTypeCCIPCommit),
		commitArgs.F,
		commitArgs.IsSignatureVerificationEnabled,
		commitSigners,
		commitTransmitters,
	)
	if err != nil {
		return mcmstypes.BatchOperation{}, fmt.Errorf("failed to encode SetOcr3Config for commit: %w", err)
	}
	mcmsTx, err := utils.GenerateMCMSTx(deps.OnChainState.CCIPAddress, moduleInfo, function, args)
	if err != nil {
		return mcmstypes.BatchOperation{}, fmt.Errorf("failed to generate MCMS operations for OffRamp Initialize: %w", err)
	}
	txs = append(txs, mcmsTx)

	execSigners := [][]byte{}
	for _, signer := range execArgs.Signers {
		execSigners = append(execSigners, signer)
	}
	execTransmitters := []aptos.AccountAddress{}
	for _, transmitter := range execArgs.Transmitters {
		address, err := aptosutils.PublicKeyBytesToAddress(transmitter)
		if err != nil {
			return mcmstypes.BatchOperation{}, fmt.Errorf("failed to convert transmitter to address: %w", err)
		}
		execTransmitters = append(execTransmitters, address)
	}
	moduleInfo, function, _, args, err = offRampBind.Offramp().Encoder().SetOcr3Config(
		execArgs.ConfigDigest[:],
		uint8(types.PluginTypeCCIPExec),
		execArgs.F,
		execArgs.IsSignatureVerificationEnabled,
		execSigners,
		execTransmitters,
	)
	if err != nil {
		return mcmstypes.BatchOperation{}, fmt.Errorf("failed to encode SetOcr3Config for exec: %w", err)
	}
	mcmsTx, err = utils.GenerateMCMSTx(deps.OnChainState.CCIPAddress, moduleInfo, function, args)
	if err != nil {
		return mcmstypes.BatchOperation{}, fmt.Errorf("failed to generate MCMS operations for OffRamp Initialize: %w", err)
	}
	txs = append(txs, mcmsTx)

	return mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions:  txs,
	}, nil
}
