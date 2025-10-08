package transmitter

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink-evm/pkg/config"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/chaintype"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	evmtxmgr "github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	cciptransmitter "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/transmitter"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
)

type ConfigTransmitterOpts struct {
	// PluginGasLimit overrides the gas limit default provided in the config watcher.
	PluginGasLimit *uint32
	// SubjectID overrides the queueing subject id (the job external id will be used by default).
	SubjectID *uuid.UUID
}

// newOnChainContractTransmitter creates a new contract transmitter.
func newOnChainContractTransmitter(ctx context.Context, lggr logger.Logger, rargs types.RelayArgs, ethKeystore keys.Store, chain legacyevm.Chain, address common.Address, opts ConfigTransmitterOpts, transmissionContractABI abi.ABI, ocrTransmitterOpts ...OCRTransmitterOption) (ContractTransmitter, error) {
	transmitter, err := generateTransmitterFrom(ctx, rargs, ethKeystore, chain, opts)
	if err != nil {
		return nil, err
	}

	evmTransmitter, err := NewOCRContractTransmitter(
		ctx,
		address,
		chain.Client(),
		transmissionContractABI,
		transmitter,
		chain.LogPoller(),
		lggr,
		ethKeystore,
		ocrTransmitterOpts...,
	)
	if err != nil {
		return nil, err
	}

	// This code path should only be called when running CCIP 1.5 jobs on Tron.
	// All other products should use the standard Tron relayer implementation.
	if chain.Config().EVM().ChainType() == chaintype.ChainTron {
		return NewTronContractTransmitter(ctx, TronContractTransmitterOpts{
			Logger:             lggr,
			TransmissionsCache: NewTronTransmissionsCache(evmTransmitter),
			Keystore:           ethKeystore,
			Chain:              chain,
			ContractAddress:    address,
			OCRTransmitterOpts: ocrTransmitterOpts,
		})
	}

	return evmTransmitter, err
}

// newOnChainDualContractTransmitter creates a new dual contract transmitter.
func newOnChainDualContractTransmitter(ctx context.Context, lggr logger.Logger, rargs types.RelayArgs, ethKeystore keys.Store, chain legacyevm.Chain, address common.Address, opts ConfigTransmitterOpts, transmissionContractABI abi.ABI, ocrTransmitterOpts ...OCRTransmitterOption) (*dualContractTransmitter, error) {
	transmitter, err := generateTransmitterFrom(ctx, rargs, ethKeystore, chain, opts)
	if err != nil {
		return nil, err
	}

	return NewOCRDualContractTransmitter(
		ctx,
		address,
		chain.Client(),
		transmissionContractABI,
		transmitter,
		chain.LogPoller(),
		lggr,
		ethKeystore,
		ocrTransmitterOpts...,
	)
}

func NewContractTransmitter(ctx context.Context, lggr logger.Logger, rargs types.RelayArgs, ethKeystore keys.Store, chain legacyevm.Chain, address common.Address, opts ConfigTransmitterOpts, transmissionContractABI abi.ABI, dualTransmission bool, ocrTransmitterOpts ...OCRTransmitterOption) (ContractTransmitter, error) {
	if dualTransmission {
		return newOnChainDualContractTransmitter(ctx, lggr, rargs, ethKeystore, chain, address, opts, transmissionContractABI, ocrTransmitterOpts...)
	}

	return newOnChainContractTransmitter(ctx, lggr, rargs, ethKeystore, chain, address, opts, transmissionContractABI, ocrTransmitterOpts...)
}

func generateTransmitterFrom(ctx context.Context, rargs types.RelayArgs, ethKeystore keys.Store, chain legacyevm.Chain, opts ConfigTransmitterOpts) (Transmitter, error) {
	var relayConfig config.RelayConfig
	if err := json.Unmarshal(rargs.RelayConfig, &relayConfig); err != nil {
		return nil, err
	}
	sendingKeys := relayConfig.SendingKeys
	if !relayConfig.EffectiveTransmitterID.Valid {
		return nil, errors.New("EffectiveTransmitterID must be specified")
	}
	effectiveTransmitterAddress := common.HexToAddress(relayConfig.EffectiveTransmitterID.String)

	sendingKeysLength := len(sendingKeys)
	if sendingKeysLength == 0 {
		return nil, errors.New("no sending keys provided")
	}

	// If we are using multiple sending keys, then a forwarder is needed to rotate transmissions.
	// Ensure that this forwarder is not set to a local sending key, and ensure our sending keys are enabled.
	var fromAddresses = make([]common.Address, 0, sendingKeysLength)
	for _, s := range sendingKeys {
		if sendingKeysLength > 1 && s == effectiveTransmitterAddress.String() {
			return nil, errors.New("the transmitter is a local sending key with transaction forwarding enabled")
		}
		if err := ethKeystore.CheckEnabled(ctx, common.HexToAddress(s)); err != nil {
			return nil, errors.Wrap(err, "one of the sending keys given is not enabled")
		}
		fromAddresses = append(fromAddresses, common.HexToAddress(s))
	}

	subject := rargs.ExternalJobID
	if opts.SubjectID != nil {
		subject = *opts.SubjectID
	}
	strategy := txmgr.NewQueueingTxStrategy(subject, relayConfig.DefaultTransactionQueueDepth)

	var checker evmtxmgr.TransmitCheckerSpec
	if relayConfig.SimulateTransactions {
		checker.CheckerType = evmtxmgr.TransmitCheckerTypeSimulate
	}

	gasLimit := chain.Config().EVM().GasEstimator().LimitDefault()
	ocr2Limit := chain.Config().EVM().GasEstimator().LimitJobType().OCR2()
	if ocr2Limit != nil {
		gasLimit = uint64(*ocr2Limit)
	}
	if opts.PluginGasLimit != nil {
		gasLimit = uint64(*opts.PluginGasLimit)
	}

	var transmitter Transmitter
	var err error

	switch types.OCR2PluginType(rargs.ProviderType) {
	case types.Median:
		transmitter, err = ocrcommon.NewOCR2FeedsTransmitter(
			chain.TxManager(),
			fromAddresses,
			common.HexToAddress(rargs.ContractID),
			gasLimit,
			effectiveTransmitterAddress,
			strategy,
			checker,
			ethKeystore,
			relayConfig.DualTransmissionConfig,
		)
	case types.CCIPExecution:
		transmitter, err = cciptransmitter.NewTransmitterWithStatusChecker(
			chain.TxManager(),
			fromAddresses,
			gasLimit,
			effectiveTransmitterAddress,
			strategy,
			checker,
			chain.ID(),
			ethKeystore,
		)
	default:
		transmitter, err = ocrcommon.NewTransmitter(
			chain.TxManager(),
			fromAddresses,
			gasLimit,
			effectiveTransmitterAddress,
			strategy,
			checker,
			ethKeystore,
		)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transmitter")
	}
	return transmitter, nil
}
