package ocrimpls

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// SVMContractTransmitterFactory implements ContractTransmitterFactory for SVM-based chains.
type SVMContractTransmitterFactory struct{}

var SVMExecCalldataFunc = func(rawReportCtx [2][32]byte,
	report ocr3types.ReportWithInfo[[]byte],
	_, _ [][32]byte,
	_ [32]byte,
	extraDataCodec ccipcommon.ExtraDataCodec,
) (contract string, method string, args any, err error) {
	// Note that the name of the struct field is very important, since the encoder used
	// by the chainwriter uses mapstructure, which will use the struct field name to map
	// to the argument name in the function call.
	// If, for whatever reason, we want to change the field name, make sure to add a `mapstructure:"<arg_name>"` tag
	// for that field.

	// WARNING: Be careful if you change the data types.
	// Using a different type e.g. `type Foo [32]byte` instead of `[32]byte`
	// will trigger undefined chainWriter behavior, e.g. transactions submitted with wrong arguments.
	var info ccipocr3.ExecuteReportInfo
	if len(report.Info) != 0 {
		var err error
		info, err = ccipocr3.DecodeExecuteReportInfo(report.Info)
		if err != nil {
			return "", "", nil, err
		}
	}

	extraDataDecoded := ccipcommon.ExtraDataDecoded{}
	if extraDataCodec != nil {
		var err error
		extraDataDecoded, err = decodeExecData(info, extraDataCodec)
		if err != nil {
			return "", "", nil, err
		}
	}

	return consts.ContractNameOffRamp,
		consts.MethodExecute,
		struct {
			ReportContext [2][32]byte
			Report        []byte
			Info          ccipocr3.ExecuteReportInfo
			ExtraData     ccipcommon.ExtraDataDecoded
		}{
			ReportContext: rawReportCtx,
			Report:        report.Report,
			Info:          info,
			ExtraData:     extraDataDecoded,
		}, nil
}

// SVMCommitCalldataFunc Returns a ToCalldataFunc that is used to generate the calldata for the commit method.
// Multiple methods are accepted in order to allow for different methods to be called based on the report data.
// The SVM on-chain contract has two methods, one for the default commit and one for the price-only commit.
func SVMCommitCalldataFunc(defaultMethod, priceOnlyMethod string) ToCalldataFunc {
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
				return "", "", nil, err
			}
		}

		method := defaultMethod
		if priceOnlyMethod != "" && len(info.MerkleRoots) == 0 && len(info.TokenPriceUpdates) > 0 {
			method = priceOnlyMethod
		}

		return consts.ContractNameOffRamp,
			method,
			struct {
				ReportContext [2][32]byte
				Report        []byte
				Rs            [][32]byte
				Ss            [][32]byte
				RawVs         [32]byte
				Info          ccipocr3.CommitReportInfo
			}{
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
		return ccipcommon.ExtraDataDecoded{}, err
	}
	// stopgap solution for missing extra args for Solana. To be replaced in the future.
	destExecDataDecoded := make([]map[string]any, len(message.TokenAmounts))
	for i, tokenAmount := range message.TokenAmounts {
		destExecDataDecoded[i], err = codec.DecodeTokenAmountDestExecData(tokenAmount.DestExecData, report.AbstractReports[0].SourceChainSelector)
		if err != nil {
			return ccipcommon.ExtraDataDecoded{}, err
		}
	}
	extraDataDecoded.DestExecDataDecoded = destExecDataDecoded

	return extraDataDecoded, nil
}

func (f *SVMContractTransmitterFactory) NewCommitTransmitter(
	cw types.ContractWriter,
	fromAccount ocrtypes.Account,
	offrampAddress string,
	defaultMethod, priceOnlyMethod string,
) ocr3types.ContractTransmitter[[]byte] {
	return &ccipTransmitter{
		cw:             cw,
		fromAccount:    fromAccount,
		offrampAddress: offrampAddress,
		toCalldataFn:   SVMCommitCalldataFunc(defaultMethod, priceOnlyMethod),
	}
}

func (f *SVMContractTransmitterFactory) NewExecTransmitter(
	cw types.ContractWriter,
	fromAccount ocrtypes.Account,
	offrampAddress string,
) ocr3types.ContractTransmitter[[]byte] {

	return &ccipTransmitter{
		cw:             cw,
		fromAccount:    fromAccount,
		offrampAddress: offrampAddress,
		toCalldataFn:   SVMExecCalldataFunc,
		extraDataCodec: ccipcommon.NewExtraDataCodec(ccipcommon.NewExtraDataCodecParams(ccipevm.ExtraDataDecoder{}, ccipsolana.ExtraDataDecoder{})),
	}
}
