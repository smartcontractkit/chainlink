package ocrimpls

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
)

type ToCalldataFunc func(rawReportCtx [3][32]byte, report []byte, rs, ss [][32]byte, vs [32]byte) any

func ToCommitCalldata(rawReportCtx [3][32]byte, report []byte, rs, ss [][32]byte, vs [32]byte) any {
	// Note that the name of the struct field is very important, since the encoder used
	// by the chainwriter uses mapstructure, which will use the struct field name to map
	// to the argument name in the function call.
	// If, for whatever reason, we want to change the field name, make sure to add a `mapstructure:"<arg_name>"` tag
	// for that field.
	return struct {
		ReportContext [3][32]byte
		Report        []byte
		Rs            [][32]byte
		Ss            [][32]byte
		RawVs         [32]byte
	}{
		ReportContext: rawReportCtx,
		Report:        report,
		Rs:            rs,
		Ss:            ss,
		RawVs:         vs,
	}
}

func ToExecCalldata(rawReportCtx [3][32]byte, report []byte, _, _ [][32]byte, _ [32]byte) any {
	// Note that the name of the struct field is very important, since the encoder used
	// by the chainwriter uses mapstructure, which will use the struct field name to map
	// to the argument name in the function call.
	// If, for whatever reason, we want to change the field name, make sure to add a `mapstructure:"<arg_name>"` tag
	// for that field.
	return struct {
		ReportContext [3][32]byte
		Report        []byte
	}{
		ReportContext: rawReportCtx,
		Report:        report,
	}
}

var _ ocr3types.ContractTransmitter[[]byte] = &commitTransmitter[[]byte]{}

type commitTransmitter[RI any] struct {
	cw             commontypes.ChainWriter
	fromAccount    ocrtypes.Account
	contractName   string
	method         string
	offrampAddress string
	toCalldataFn   ToCalldataFunc
}

func XXXNewContractTransmitterTestsOnly[RI any](
	cw commontypes.ChainWriter,
	fromAccount ocrtypes.Account,
	contractName string,
	method string,
	offrampAddress string,
	toCalldataFn ToCalldataFunc,
) ocr3types.ContractTransmitter[RI] {
	return &commitTransmitter[RI]{
		cw:             cw,
		fromAccount:    fromAccount,
		contractName:   contractName,
		method:         method,
		offrampAddress: offrampAddress,
		toCalldataFn:   toCalldataFn,
	}
}

func NewCommitContractTransmitter[RI any](
	cw commontypes.ChainWriter,
	fromAccount ocrtypes.Account,
	offrampAddress string,
) ocr3types.ContractTransmitter[RI] {
	return &commitTransmitter[RI]{
		cw:             cw,
		fromAccount:    fromAccount,
		contractName:   consts.ContractNameOffRamp,
		method:         consts.MethodCommit,
		offrampAddress: offrampAddress,
		toCalldataFn:   ToCommitCalldata,
	}
}

func NewExecContractTransmitter[RI any](
	cw commontypes.ChainWriter,
	fromAccount ocrtypes.Account,
	offrampAddress string,
) ocr3types.ContractTransmitter[RI] {
	return &commitTransmitter[RI]{
		cw:             cw,
		fromAccount:    fromAccount,
		contractName:   consts.ContractNameOffRamp,
		method:         consts.MethodExecute,
		offrampAddress: offrampAddress,
		toCalldataFn:   ToExecCalldata,
	}
}

// FromAccount implements ocr3types.ContractTransmitter.
func (c *commitTransmitter[RI]) FromAccount(context.Context) (ocrtypes.Account, error) {
	return c.fromAccount, nil
}

// Transmit implements ocr3types.ContractTransmitter.
func (c *commitTransmitter[RI]) Transmit(
	ctx context.Context,
	configDigest ocrtypes.ConfigDigest,
	seqNr uint64,
	reportWithInfo ocr3types.ReportWithInfo[RI],
	sigs []ocrtypes.AttributedOnchainSignature,
) error {
	var rs [][32]byte
	var ss [][32]byte
	var vs [32]byte
	if len(sigs) > 32 {
		return errors.New("too many signatures, maximum is 32")
	}
	for i, as := range sigs {
		r, s, v, err := evmutil.SplitSignature(as.Signature)
		if err != nil {
			return fmt.Errorf("failed to split signature: %w", err)
		}
		rs = append(rs, r)
		ss = append(ss, s)
		vs[i] = v
	}

	// report ctx for OCR3 consists of the following
	// reportContext[0]: ConfigDigest
	// reportContext[1]: 24 byte padding, 8 byte sequence number
	// reportContext[2]: unused
	// convert seqNum, which is a uint64, into a uint32 epoch and uint8 round
	// while this does truncate the sequence number, it is not a problem because
	// it still gives us 2^40 - 1 possible sequence numbers.
	// assuming a sequence number is generated every second, this gives us
	// 1099511627775 seconds, or approximately 34,865 years, before we run out
	// of sequence numbers.
	epoch, round := uint64ToUint32AndUint8(seqNr)
	rawReportCtx := evmutil.RawReportContext(ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: configDigest,
			Epoch:        epoch,
			Round:        round,
		},
		// ExtraData not used in OCR3
	})

	if c.toCalldataFn == nil {
		return errors.New("toCalldataFn is nil")
	}

	// chain writer takes in the raw calldata and packs it on its own.
	args := c.toCalldataFn(rawReportCtx, reportWithInfo.Report, rs, ss, vs)

	// TODO: no meta fields yet, what should we add?
	// probably whats in the info part of the report?
	meta := commontypes.TxMeta{}
	txID, err := uuid.NewRandom() // NOTE: CW expects us to generate an ID, rather than return one
	if err != nil {
		return fmt.Errorf("failed to generate UUID: %w", err)
	}
	zero := big.NewInt(0)
	if err := c.cw.SubmitTransaction(ctx, c.contractName, c.method, args, fmt.Sprintf("%s-%s-%s", c.contractName, c.offrampAddress, txID.String()), c.offrampAddress, &meta, zero); err != nil {
		return fmt.Errorf("failed to submit transaction thru chainwriter: %w", err)
	}

	return nil
}

func uint64ToUint32AndUint8(x uint64) (uint32, uint8) {
	return uint32(x >> 32), uint8(x)
}
