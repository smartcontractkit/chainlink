package ocrimpls

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// EVMCommitCallArgs defines the calldata structure for an EVM commit transaction.
// IMPORTANT: The names and types of the fields are critical because the chainwriter uses mapstructure
// to map these fields to the contract's parameter names. Changing these names or types (or omitting the
// mapstructure tags) may result in transactions being constructed with incorrect arguments.
type EVMCommitCallArgs struct {
	ReportContext [2][32]byte `mapstructure:"ReportContext"`
	Report        []byte      `mapstructure:"Report"`
	Rs            [][32]byte  `mapstructure:"Rs"`
	Ss            [][32]byte  `mapstructure:"Ss"`
	RawVs         [32]byte    `mapstructure:"RawVs"`
}

// EVMExecCallArgs defines the calldata structure for an EVM execute transaction.
// IMPORTANT: The names and types of the fields are critical because the chainwriter uses mapstructure
// to map these fields to the contract's parameter names. Changing these names or types (or omitting the
// mapstructure tags) may result in transactions being constructed with incorrect arguments.
type EVMExecCallArgs struct {
	ReportContext [2][32]byte `mapstructure:"ReportContext"`
	Report        []byte      `mapstructure:"Report"`
}

// EVMContractTransmitterFactory implements the transmitter factory for EVM chains.
type EVMContractTransmitterFactory struct{}

// EVMExecCallDataFunc builds the execute call data for EVM.
var EVMExecCallDataFunc = func(
	rawReportCtx [2][32]byte,
	report ocr3types.ReportWithInfo[[]byte],
	_, _ [][32]byte,
	_ [32]byte,
	_ ccipcommon.ExtraDataCodec,
) (contract string, method string, args any, err error) {
	return consts.ContractNameOffRamp,
		consts.MethodExecute,
		EVMExecCallArgs{
			ReportContext: rawReportCtx,
			Report:        report.Report,
		}, nil
}

// NewEVMCommitCalldataFunc returns a ToCalldataFunc for EVM commits that omits any Info object.
func NewEVMCommitCalldataFunc(commitMethod string) ToCalldataFunc {
	return func(
		rawReportCtx [2][32]byte,
		report ocr3types.ReportWithInfo[[]byte],
		rs, ss [][32]byte,
		vs [32]byte,
		_ ccipcommon.ExtraDataCodec,
	) (string, string, any, error) {
		return consts.ContractNameOffRamp,
			commitMethod,
			EVMCommitCallArgs{
				ReportContext: rawReportCtx,
				Report:        report.Report,
				Rs:            rs,
				Ss:            ss,
				RawVs:         vs,
			},
			nil
	}
}

// NewCommitTransmitter constructs an EVM commit transmitter.
func (f *EVMContractTransmitterFactory) NewCommitTransmitter(
	lggr logger.Logger,
	cw types.ContractWriter,
	fromAccount ocrtypes.Account,
	offrampAddress string,
	commitMethod, _ string, // priceOnlyMethod is ignored for EVM
) ocr3types.ContractTransmitter[[]byte] {
	return &ccipTransmitter{
		lggr:           lggr,
		cw:             cw,
		fromAccount:    fromAccount,
		offrampAddress: offrampAddress,
		toCalldataFn:   NewEVMCommitCalldataFunc(commitMethod),
	}
}

// NewExecTransmitter constructs an EVM execute transmitter.
func (f *EVMContractTransmitterFactory) NewExecTransmitter(
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
		toCalldataFn:   EVMExecCallDataFunc,
	}
}
