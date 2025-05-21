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

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// ToCalldataFunc is a function that takes in the OCR3 report and signature data and processes them.
// It returns the contract name, method name, and arguments for the on-chain contract call.
// The ReportWithInfo bytes field is also decoded according to the implementation of this function,
// the commit and execute plugins have different representations for this data.
type ToCalldataFunc func(
	rawReportCtx [2][32]byte,
	report ocr3types.ReportWithInfo[[]byte],
	rs, ss [][32]byte,
	vs [32]byte,
	codec ccipcommon.ExtraDataCodec,
) (contract string, method string, args any, err error)

var _ ocr3types.ContractTransmitter[[]byte] = &ccipTransmitter{}

type ccipTransmitter struct {
	cw             commontypes.ContractWriter
	fromAccount    ocrtypes.Account
	offrampAddress string
	toCalldataFn   ToCalldataFunc
	extraDataCodec ccipcommon.ExtraDataCodec
	lggr           logger.Logger
}

func XXXNewContractTransmitterTestsOnly(
	lggr logger.Logger,
	cw commontypes.ContractWriter,
	fromAccount ocrtypes.Account,
	contractName string,
	method string,
	offrampAddress string,
	toCalldataFn ToCalldataFunc,
) ocr3types.ContractTransmitter[[]byte] {
	wrappedToCalldataFunc := func(rawReportCtx [2][32]byte,
		report ocr3types.ReportWithInfo[[]byte],
		rs, ss [][32]byte,
		vs [32]byte,
		extraDataCodec ccipcommon.ExtraDataCodec) (string, string, any, error) {
		_, _, args, err := toCalldataFn(rawReportCtx, report, rs, ss, vs, extraDataCodec)
		return contractName, method, args, err
	}
	return &ccipTransmitter{
		lggr:           lggr,
		cw:             cw,
		fromAccount:    fromAccount,
		offrampAddress: offrampAddress,
		toCalldataFn:   wrappedToCalldataFunc,
	}
}

// FromAccount implements ocr3types.ContractTransmitter.
func (c *ccipTransmitter) FromAccount(context.Context) (ocrtypes.Account, error) {
	return c.fromAccount, nil
}

// Transmit implements ocr3types.ContractTransmitter.
func (c *ccipTransmitter) Transmit(
	ctx context.Context,
	configDigest ocrtypes.ConfigDigest,
	seqNr uint64,
	reportWithInfo ocr3types.ReportWithInfo[[]byte],
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
	rawReportCtx := ocr2key.RawReportContext3(configDigest, seqNr)

	if c.toCalldataFn == nil {
		return errors.New("toCalldataFn is nil")
	}

	// chain writer takes in the raw calldata and packs it on its own.
	contract, method, args, err := c.toCalldataFn(rawReportCtx, reportWithInfo, rs, ss, vs, c.extraDataCodec)
	if err != nil {
		return fmt.Errorf("failed to generate call data: %w", err)
	}

	// TODO: no meta fields yet, what should we add?
	// probably whats in the info part of the report?
	meta := commontypes.TxMeta{}
	txID, err := uuid.NewRandom() // NOTE: CW expects us to generate an ID, rather than return one
	if err != nil {
		return fmt.Errorf("failed to generate UUID: %w", err)
	}
	zero := big.NewInt(0)
	c.lggr.Infow("Submitting transaction", "tx", txID)
	if err := c.cw.SubmitTransaction(ctx, contract, method, args,
		fmt.Sprintf("%s-%s-%s", contract, c.offrampAddress, txID.String()),
		c.offrampAddress, &meta, zero); err != nil {
		return fmt.Errorf("failed to submit transaction via chain writer: %w", err)
	}

	return nil
}
