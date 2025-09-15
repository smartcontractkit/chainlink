package ocrimpls

import (
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// SuiCommitCallArgs defines the calldata structure for an Sui commit transaction.
type SuiCommitCallArgs struct {
	ReportContext [2][32]byte `mapstructure:"ReportContext"`
	Report        []byte      `mapstructure:"Report"`
	Signatures    [][96]byte  `mapstructure:"Signatures"`
}

// SuiExecCallArgs defines the calldata structure for an Sui execute transaction.
type SuiExecCallArgs struct {
	ReportContext [2][32]byte                 `mapstructure:"ReportContext"`
	Report        []byte                      `mapstructure:"Report"`
	Info          ccipocr3.ExecuteReportInfo  `mapstructure:"Info"`
	ExtraData     ccipcommon.ExtraDataDecoded `mapstructure:"ExtraData"`
}

// SuiContractTransmitterFactory implements the transmitter factory for Sui chains.
type SuiContractTransmitterFactory struct {
	extraDataCodec ccipocr3.ExtraDataCodecBundle
}

func NewSuiContractTransmitterFactory(extraDataCodec ccipocr3.ExtraDataCodecBundle) *SuiContractTransmitterFactory {
	return &SuiContractTransmitterFactory{
		extraDataCodec: extraDataCodec,
	}
}

// NewSuiCommitCalldataFunc returns a ToCalldataFunc for Sui commits that omits any Info object.
func NewSuiCommitCalldataFunc(commitMethod string) ToEd25519CalldataFunc {
	return func(
		rawReportCtx [2][32]byte,
		report ocr3types.ReportWithInfo[[]byte],
		signatures [][96]byte,
		_ ccipocr3.ExtraDataCodecBundle,
	) (string, string, any, error) {
		return consts.ContractNameOffRamp,
			commitMethod,
			SuiCommitCallArgs{
				ReportContext: rawReportCtx,
				Report:        report.Report,
				Signatures:    signatures,
			},
			nil
	}
}

// NewCommitTransmitter constructs an Sui commit transmitter.
func (f *SuiContractTransmitterFactory) NewCommitTransmitter(
	lggr logger.Logger,
	cw types.ContractWriter,
	fromAccount ocrtypes.Account,
	offrampAddress string,
	commitMethod, _ string, // priceOnlyMethod is ignored for Sui
) ocr3types.ContractTransmitter[[]byte] {
	return &ccipTransmitter{
		lggr:                lggr,
		cw:                  cw,
		fromAccount:         fromAccount,
		offrampAddress:      offrampAddress,
		toEd25519CalldataFn: NewSuiCommitCalldataFunc(commitMethod),
		extraDataCodec:      f.extraDataCodec,
	}
}

// SuiExecCallDataFunc builds the execute call data for Sui
var SuiExecCallDataFunc = func(
	rawReportCtx [2][32]byte,
	report ocr3types.ReportWithInfo[[]byte],
	signatures [][96]byte,
	extraDataCodec ccipocr3.ExtraDataCodecBundle,
) (contract string, method string, args any, err error) {
	var info ccipocr3.ExecuteReportInfo
	var extraDataDecoded ccipcommon.ExtraDataDecoded
	if len(report.Info) != 0 {
		info, err = ccipocr3.DecodeExecuteReportInfo(report.Info)
		if err != nil {
			return "", "", nil, fmt.Errorf("failed to decode execute report info: %w", err)
		}
		if extraDataCodec != nil {
			extraDataDecoded, err = decodeExecDataSui(info, extraDataCodec)
			if err != nil {
				return "", "", nil, fmt.Errorf("failed to decode extra data: %w", err)
			}
		}
	}
	return consts.ContractNameOffRamp,
		consts.MethodExecute,
		SuiExecCallArgs{
			ReportContext: rawReportCtx,
			Report:        report.Report,
			Info:          info,
			ExtraData:     extraDataDecoded,
		}, nil
}

// NewExecTransmitter constructs an Sui execute transmitter.
func (f *SuiContractTransmitterFactory) NewExecTransmitter(
	lggr logger.Logger,
	cw types.ContractWriter,
	fromAccount ocrtypes.Account,
	offrampAddress string,
) ocr3types.ContractTransmitter[[]byte] {
	return &ccipTransmitter{
		lggr:                lggr,
		cw:                  cw,
		fromAccount:         fromAccount,
		offrampAddress:      offrampAddress,
		toEd25519CalldataFn: SuiExecCallDataFunc,
		extraDataCodec:      f.extraDataCodec,
	}
}

// decodeExecData decodes the extra data from an execute report.
func decodeExecDataSui(report ccipocr3.ExecuteReportInfo, codec ccipocr3.ExtraDataCodecBundle) (ccipcommon.ExtraDataDecoded, error) {
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
