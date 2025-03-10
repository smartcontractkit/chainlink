package ocrimpls

import (
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// EVMContractTransmitterFactory implements ContractTransmitterFactory for EVM-based chains.
type EVMContractTransmitterFactory struct{}

var EVMExecCallDataFunc = func(
	rawReportCtx [2][32]byte,
	report ocr3types.ReportWithInfo[[]byte],
	_, _ [][32]byte,
	_ [32]byte,
) (contract string, method string, args any, err error) {
	// Note that the name of the struct field is very important, since the encoder used
	// by the chainwriter uses mapstructure, which will use the struct field name to map
	// to the argument name in the function call.
	// If, for whatever reason, we want to change the field name, make sure to add a `mapstructure:"<arg_name>"` tag
	// for that field.

	return consts.ContractNameOffRamp,
		consts.MethodExecute,
		struct {
			ReportContext [2][32]byte
			Report        []byte
		}{
			ReportContext: rawReportCtx,
			Report:        report.Report,
		}, nil
}

// EVMCommitCalldataFunc returns a ToCalldataFunc that omits the Info object for EVM.
func EVMCommitCalldataFunc(defaultMethod string) ToCalldataFunc {
	return func(
		rawReportCtx [2][32]byte,
		report ocr3types.ReportWithInfo[[]byte],
		rs, ss [][32]byte,
		vs [32]byte,
	) (string, string, any, error) {
		// Note that the name of the struct field is very important, since the encoder used
		// by the chainwriter uses mapstructure, which will use the struct field name to map
		// to the argument name in the function call.
		// If, for whatever reason, we want to change the field name, make sure to add a `mapstructure:"<arg_name>"` tag
		// for that field.
		return consts.ContractNameOffRamp,
			defaultMethod,
			struct {
				ReportContext [2][32]byte
				Report        []byte
				Rs            [][32]byte
				Ss            [][32]byte
				RawVs         [32]byte
			}{
				ReportContext: rawReportCtx,
				Report:        report.Report,
				Rs:            rs,
				Ss:            ss,
				RawVs:         vs,
			},
			nil
	}
}

func (f *EVMContractTransmitterFactory) NewCommitTransmitter(
	cw types.ContractWriter,
	fromAccount ocrtypes.Account,
	offrampAddress string,
	defaultMethod, _ string,
) ocr3types.ContractTransmitter[[]byte] {
	return &ccipTransmitter{
		cw:             cw,
		fromAccount:    fromAccount,
		offrampAddress: offrampAddress,
		toCalldataFn:   EVMCommitCalldataFunc(defaultMethod),
	}
}

func (f *EVMContractTransmitterFactory) NewExecTransmitter(
	cw types.ContractWriter,
	fromAccount ocrtypes.Account,
	offrampAddress string,
) ocr3types.ContractTransmitter[[]byte] {

	return &ccipTransmitter{
		cw:             cw,
		fromAccount:    fromAccount,
		offrampAddress: offrampAddress,
		toCalldataFn:   EVMExecCallDataFunc,
	}
}
