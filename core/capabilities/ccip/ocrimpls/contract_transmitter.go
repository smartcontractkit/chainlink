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
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
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
) (contract string, method string, args any, err error)

// NewToCommitCalldataFunc returns a ToCalldataFunc that is used to generate the calldata for the commit method.
// Multiple methods are accepted in order to allow for different methods to be called based on the report data.
// The Solana on-chain contract has two methods, one for the default commit and one for the price-only commit.
func NewToCommitCalldataFunc(defaultMethod, priceOnlyMethod string) ToCalldataFunc {
	return func(
		rawReportCtx [2][32]byte,
		report ocr3types.ReportWithInfo[[]byte],
		rs, ss [][32]byte,
		vs [32]byte,
	) (contract string, method string, args any, err error) {
		// Note that the name of the struct field is very important, since the encoder used
		// by the chainwriter uses mapstructure, which will use the struct field name to map
		// to the argument name in the function call.
		// If, for whatever reason, we want to change the field name, make sure to add a `mapstructure:"<arg_name>"` tag
		// for that field.
		var info ccipocr3.CommitReportInfo
		if len(report.Info) != 0 {
			var err error
			info, err = ccipocr3.DecodeCommitReportInfo(report.Info)
			if err != nil {
				return "", "", nil, err
			}
		}

		method = defaultMethod
		if priceOnlyMethod != "" && len(info.MerkleRoots) == 0 && len(info.TokenPriceUpdates) > 0 {
			method = priceOnlyMethod
		}

		// WARNING: Be careful if you change the data types.
		// Using a different type e.g. `type Foo [32]byte` instead of `[32]byte`
		// will trigger undefined chainWriter behavior, e.g. transactions submitted with wrong arguments.
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
			}, nil
	}
}

// ToExecCalldata is a ToCalldataFunc that is used to generate the calldata for the execute method.
func ToExecCalldata(
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

var _ ocr3types.ContractTransmitter[[]byte] = &ccipTransmitter{}

type ccipTransmitter struct {
	cw             commontypes.ContractWriter
	fromAccount    ocrtypes.Account
	offrampAddress string
	toCalldataFn   ToCalldataFunc
}

func XXXNewContractTransmitterTestsOnly(
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
		vs [32]byte) (string, string, any, error) {
		_, _, args, err := toCalldataFn(rawReportCtx, report, rs, ss, vs)
		return contractName, method, args, err
	}
	return &ccipTransmitter{
		cw:             cw,
		fromAccount:    fromAccount,
		offrampAddress: offrampAddress,
		toCalldataFn:   wrappedToCalldataFunc,
	}
}

func NewCommitContractTransmitter(
	cw commontypes.ContractWriter,
	fromAccount ocrtypes.Account,
	offrampAddress string,
	defaultMethod, priceOnlyMethod string,
) ocr3types.ContractTransmitter[[]byte] {
	return &ccipTransmitter{
		cw:             cw,
		fromAccount:    fromAccount,
		offrampAddress: offrampAddress,
		toCalldataFn:   NewToCommitCalldataFunc(defaultMethod, priceOnlyMethod),
	}
}

func NewExecContractTransmitter(
	cw commontypes.ContractWriter,
	fromAccount ocrtypes.Account,
	offrampAddress string,
) ocr3types.ContractTransmitter[[]byte] {
	return &ccipTransmitter{
		cw:             cw,
		fromAccount:    fromAccount,
		offrampAddress: offrampAddress,
		toCalldataFn:   ToExecCalldata,
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
	contract, method, args, err := c.toCalldataFn(rawReportCtx, reportWithInfo, rs, ss, vs)
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
	if err := c.cw.SubmitTransaction(ctx, contract, method, args, fmt.Sprintf("%s-%s-%s", contract, c.offrampAddress, txID.String()), c.offrampAddress, &meta, zero); err != nil {
		return fmt.Errorf("failed to submit transaction thru chainwriter: %w", err)
	}

	return nil
}
