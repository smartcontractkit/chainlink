package ocrimpls

import (
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// SVMCommitCallArgs defines the calldata structure for an SVM commit transaction.
// IMPORTANT: The names and types of the fields are critical because the chainwriter uses mapstructure
// to map these fields to the contract's parameter names. Changing these names or types (or omitting the
// mapstructure tags) may result in transactions being constructed with incorrect arguments.
type SVMCommitCallArgs struct {
	ReportContext [2][32]byte               `mapstructure:"ReportContext"`
	Report        []byte                    `mapstructure:"Report"`
	Rs            [][32]byte                `mapstructure:"Rs"`
	Ss            [][32]byte                `mapstructure:"Ss"`
	RawVs         [32]byte                  `mapstructure:"RawVs"`
	Info          ccipocr3.CommitReportInfo `mapstructure:"Info"`
}

// SVMExecCallArgs defines the calldata structure for an SVM execute transaction.
// IMPORTANT: The names and types of the fields are critical because the chainwriter uses mapstructure
// to map these fields to the contract's parameter names. Changing these names or types (or omitting the
// mapstructure tags) may result in transactions being constructed with incorrect arguments.
type SVMExecCallArgs struct {
	ReportContext [2][32]byte                 `mapstructure:"ReportContext"`
	Report        []byte                      `mapstructure:"Report"`
	Info          ccipocr3.ExecuteReportInfo  `mapstructure:"Info"`
	ExtraData     ccipcommon.ExtraDataDecoded `mapstructure:"ExtraData"`
}

// SVMContractTransmitterFactory implements the transmitter factory for SVM chains.
type SVMContractTransmitterFactory struct {
	extraDataCodec ccipcommon.ExtraDataCodec
}

// NewSVMContractTransmitterFactory returns a new SVMContractTransmitterFactory.
func NewSVMContractTransmitterFactory(extraDataCodec ccipcommon.ExtraDataCodec) *SVMContractTransmitterFactory {
	return &SVMContractTransmitterFactory{
		extraDataCodec: extraDataCodec,
	}
}

// SVMExecCalldataFunc builds the execute call data for SVM.
var SVMExecCalldataFunc = func(
	rawReportCtx [2][32]byte,
	report ocr3types.ReportWithInfo[[]byte],
	_, _ [][32]byte,
	_ [32]byte,
	extraDataCodec ccipcommon.ExtraDataCodec,
) (contract string, method string, args any, err error) {
	var info ccipocr3.ExecuteReportInfo
	var extraDataDecoded ccipcommon.ExtraDataDecoded
	if len(report.Info) != 0 {
		info, err = ccipocr3.DecodeExecuteReportInfo(report.Info)
		if err != nil {
			return "", "", nil, fmt.Errorf("failed to decode execute report info: %w", err)
		}
		if extraDataCodec != nil {
			extraDataDecoded, err = decodeExecData(info, extraDataCodec)
			if err != nil {
				return "", "", nil, fmt.Errorf("failed to decode extra data: %w", err)
			}
		}
	}

	return consts.ContractNameOffRamp,
		consts.MethodExecute,
		SVMExecCallArgs{
			ReportContext: rawReportCtx,
			Report:        report.Report,
			Info:          info,
			ExtraData:     extraDataDecoded,
		}, nil
}

// NewSVMCommitCalldataFunc Returns a ToCalldataFunc that is used to generate the calldata for the commit method.
// // Multiple methods are accepted in order to allow for different methods to be called based on the report data.
// // The SVM on-chain contract has two methods, one for the default commit and one for the price-only commit.
func NewSVMCommitCalldataFunc(defaultMethod, priceOnlyMethod string) ToCalldataFunc {
	return func(
		rawReportCtx [2][32]byte,
		report ocr3types.ReportWithInfo[[]byte],
		rs, ss [][32]byte,
		vs [32]byte,
		_ ccipcommon.ExtraDataCodec,
	) (string, string, any, error) {
		var info ccipocr3.CommitReportInfo
		if len(report.Info) != 0 {
			var err error
			info, err = ccipocr3.DecodeCommitReportInfo(report.Info)
			if err != nil {
				return "", "", nil, fmt.Errorf("failed to decode commit report info: %w", err)
			}
		}

		method := defaultMethod
		// Switch to price-only method if no Merkle roots and there are token or gas price updates.
		if priceOnlyMethod != "" && len(info.MerkleRoots) == 0 && (len(info.TokenPriceUpdates) > 0 || len(info.GasPriceUpdates) > 0) {
			method = priceOnlyMethod
		}

		return consts.ContractNameOffRamp,
			method,
			SVMCommitCallArgs{
				ReportContext: rawReportCtx,
				Report:        report.Report,
				Rs:            rs,
				Ss:            ss,
				RawVs:         vs,
				Info:          info,
			},
			nil
	}
}

// decodeExecData decodes the extra data from an execute report.
func decodeExecData(report ccipocr3.ExecuteReportInfo, codec ccipcommon.ExtraDataCodec) (ccipcommon.ExtraDataDecoded, error) {
	// only one report one message, since this is a stop-gap solution for solana
	if len(report.AbstractReports) != 1 {
		return ccipcommon.ExtraDataDecoded{}, fmt.Errorf("unexpected report length, expected 1, got %d", len(report.AbstractReports))
	}
	if len(report.AbstractReports[0].Messages) != 1 {
		return ccipcommon.ExtraDataDecoded{}, fmt.Errorf("unexpected message length, expected 1, got %d", len(report.AbstractReports[0].Messages))
	}
	message := report.AbstractReports[0].Messages[0]
	extraDataDecoded := ccipcommon.ExtraDataDecoded{}

	var err error
	extraDataDecoded.ExtraArgsDecoded, err = codec.DecodeExtraArgs(message.ExtraArgs, report.AbstractReports[0].SourceChainSelector)
	if err != nil {
		return ccipcommon.ExtraDataDecoded{}, fmt.Errorf("failed to decode extra args: %w", err)
	}
	// stopgap solution for missing extra args for Solana. To be replaced in the future.
	destExecDataDecoded := make([]map[string]any, len(message.TokenAmounts))
	for i, tokenAmount := range message.TokenAmounts {
		destExecDataDecoded[i], err = codec.DecodeTokenAmountDestExecData(tokenAmount.DestExecData, report.AbstractReports[0].SourceChainSelector)
		if err != nil {
			return ccipcommon.ExtraDataDecoded{}, fmt.Errorf("failed to decode token amount dest exec data: %w", err)
		}
	}
	extraDataDecoded.DestExecDataDecoded = destExecDataDecoded

	return extraDataDecoded, nil
}

// NewCommitTransmitter constructs an SVM commit transmitter.
func (f *SVMContractTransmitterFactory) NewCommitTransmitter(
	lggr logger.Logger,
	cw types.ContractWriter,
	fromAccount ocrtypes.Account,
	offrampAddress string,
	defaultMethod, priceOnlyMethod string,
) ocr3types.ContractTransmitter[[]byte] {
	return &ccipTransmitter{
		lggr:           lggr,
		cw:             cw,
		fromAccount:    fromAccount,
		offrampAddress: offrampAddress,
		toCalldataFn:   NewSVMCommitCalldataFunc(defaultMethod, priceOnlyMethod),
		extraDataCodec: f.extraDataCodec,
	}
}

// NewExecTransmitter constructs an SVM execute transmitter.
func (f *SVMContractTransmitterFactory) NewExecTransmitter(
	lggr logger.Logger,
	cw types.ContractWriter,
	fromAccount ocrtypes.Account,
	offrampAddress string,
) ocr3types.ContractTransmitter[[]byte] {
	return &ccipTransmitter{
		lggr:           lggr,
		cw:             cw,
		fromAccount:    fromAccount,
		offrampAddress: offrampAddress,
		toCalldataFn:   SVMExecCalldataFunc,
		extraDataCodec: f.extraDataCodec,
	}
}
