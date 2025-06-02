package ocrimpls

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type NoopTransmitter struct {
	extraDataCodec common.ExtraDataCodec
}

// NewNoopTransmitter constructs a Noop transmitter.
func NewNoopTransmitter(extraDataCodec common.ExtraDataCodec) *NoopTransmitter {
	return &NoopTransmitter{
		extraDataCodec: extraDataCodec,
	}
}

// NewCommitTransmitter constructs an EVM commit transmitter.
func (f *NoopTransmitter) NewCommitTransmitter(
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
		toCalldataFn:   NewNoopCommitCalldataFunc(commitMethod),
		extraDataCodec: f.extraDataCodec,
	}
}

// NewExecTransmitter constructs an EVM execute transmitter.
func (f *NoopTransmitter) NewExecTransmitter(
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
		toCalldataFn:   NoopExecCallDataFunc,
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
		_ common.ExtraDataCodec,
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
	_ common.ExtraDataCodec,
) (contract string, method string, args any, err error) {
	return consts.ContractNameOffRamp,
		consts.MethodExecute,
		nil, nil
}
