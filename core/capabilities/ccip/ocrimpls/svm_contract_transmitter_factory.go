package ocrimpls

import (
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// SVMContractTransmitterFactory implements ContractTransmitterFactory for SVM-based chains.
type SVMContractTransmitterFactory struct{}

var SVMExecCallDataFunc = func(rawReportCtx [2][32]byte,
	report ocr3types.ReportWithInfo[[]byte],
	_, _ [][32]byte,
	_ [32]byte,
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

	return consts.ContractNameOffRamp,
		consts.MethodExecute,
		struct {
			ReportContext [2][32]byte
			Report        []byte
			Info          ccipocr3.ExecuteReportInfo
		}{
			ReportContext: rawReportCtx,
			Report:        report.Report,
			Info:          info,
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
		toCalldataFn:   SVMExecCallDataFunc,
	}
}
