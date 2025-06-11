package ocrimpls

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

type NoopContractTransmitterFactory struct {
	extraDataCodec ccipcommon.ExtraDataCodec
}

// NewNoopContractTransmitterFactory constructs a Noop transmitter.
func NewNoopContractTransmitterFactory(extraDataCodec ccipcommon.ExtraDataCodec) *NoopContractTransmitterFactory {
	return &NoopContractTransmitterFactory{
		extraDataCodec: extraDataCodec,
	}
}

// NewCommitTransmitter constructs a Noop commit transmitter.
func (f *NoopContractTransmitterFactory) NewCommitTransmitter(
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
		extraDataCodec: f.extraDataCodec,
	}
}

// NewExecTransmitter constructs a Noop execute transmitter.
func (f *NoopContractTransmitterFactory) NewExecTransmitter(
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
		extraDataCodec: f.extraDataCodec,
	}
}

// NewNoopCommitCalldataFunc returns a ToCalldataFunc for noop commits that omits any Info object.
func NewNoopCommitCalldataFunc(commitMethod string) ToCalldataFunc {
	return func(
		rawReportCtx [2][32]byte,
		report ocr3types.ReportWithInfo[[]byte],
		rs, ss [][32]byte,
		vs [32]byte,
		_ ccipcommon.ExtraDataCodec,
	) (string, string, any, error) {
		return consts.ContractNameOffRamp,
			commitMethod,
			nil,
			nil
	}
}

// NoopExecCallDataFunc builds the noop execute call data.
var NoopExecCallDataFunc = func(
	rawReportCtx [2][32]byte,
	report ocr3types.ReportWithInfo[[]byte],
	_, _ [][32]byte,
	_ [32]byte,
	_ ccipcommon.ExtraDataCodec,
) (contract string, method string, args any, err error) {
	return consts.ContractNameOffRamp,
		consts.MethodExecute,
		nil, nil
}
