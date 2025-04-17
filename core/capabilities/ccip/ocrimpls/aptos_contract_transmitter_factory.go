package ocrimpls

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipaptos"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// AptosCommitCallArgs defines the calldata structure for an Aptos commit transaction.
type AptosCommitCallArgs struct {
	ReportContext [2][32]byte `mapstructure:"ReportContext"`
	Report        []byte      `mapstructure:"Report"`
	Signatures    [][96]byte  `mapstructure:"Signatures"`
}

// AptosExecCallArgs defines the calldata structure for an Aptos execute transaction.
type AptosExecCallArgs struct {
	ReportContext [2][32]byte `mapstructure:"ReportContext"`
	Report        []byte      `mapstructure:"Report"`
}

// AptosContractTransmitterFactory implements the transmitter factory for Aptos chains.
type AptosContractTransmitterFactory struct{}

// NewAptosCommitCalldataFunc returns a ToCalldataFunc for Aptos commits that omits any Info object.
func NewAptosCommitCalldataFunc(commitMethod string) ToEd25519CalldataFunc {
	return func(
		rawReportCtx [2][32]byte,
		report ocr3types.ReportWithInfo[[]byte],
		signatures [][96]byte,
		_ ccipcommon.ExtraDataCodec,
	) (string, string, any, error) {
		return consts.ContractNameOffRamp,
			commitMethod,
			AptosCommitCallArgs{
				ReportContext: rawReportCtx,
				Report:        report.Report,
				Signatures:    signatures,
			},
			nil
	}
}

// NewCommitTransmitter constructs an Aptos commit transmitter.
func (f *AptosContractTransmitterFactory) NewCommitTransmitter(
	cw types.ContractWriter,
	fromAccount ocrtypes.Account,
	offrampAddress string,
	commitMethod, _ string, // priceOnlyMethod is ignored for Aptos
) ocr3types.ContractTransmitter[[]byte] {
	return &ccipTransmitter{
		cw:                  cw,
		fromAccount:         fromAccount,
		offrampAddress:      offrampAddress,
		toEd25519CalldataFn: NewAptosCommitCalldataFunc(commitMethod),
	}
}

// AptosExecCallDataFunc builds the execute call data for Aptos
var AptosExecCallDataFunc = func(
	rawReportCtx [2][32]byte,
	report ocr3types.ReportWithInfo[[]byte],
	signatures [][96]byte,
	_ ccipcommon.ExtraDataCodec,
) (contract string, method string, args any, err error) {
	return consts.ContractNameOffRamp,
		consts.MethodExecute,
		AptosExecCallArgs{
			ReportContext: rawReportCtx,
			Report:        report.Report,
		}, nil
}

// NewExecTransmitter constructs an Aptos execute transmitter.
func (f *AptosContractTransmitterFactory) NewExecTransmitter(
	cw types.ContractWriter,
	fromAccount ocrtypes.Account,
	offrampAddress string,
) ocr3types.ContractTransmitter[[]byte] {
	return &ccipTransmitter{
		cw:                  cw,
		fromAccount:         fromAccount,
		offrampAddress:      offrampAddress,
		toEd25519CalldataFn: AptosExecCallDataFunc,
		extraDataCodec: ccipcommon.NewExtraDataCodec(
			ccipcommon.NewExtraDataCodecParams(ccipevm.ExtraDataDecoder{}, ccipsolana.ExtraDataDecoder{}, ccipaptos.ExtraDataDecoder{}),
		),
	}
}
