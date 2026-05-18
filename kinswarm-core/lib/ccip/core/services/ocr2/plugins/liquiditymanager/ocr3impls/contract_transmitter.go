package ocr3impls

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	// TransmittedOCR3 is the name of the Transmitted event in the OCR3 contract.
	TransmittedOCR3 = "Transmitted"
)

type Transmitter interface {
	CreateEthTransaction(ctx context.Context, toAddress gethcommon.Address, payload []byte, txMeta *txmgr.TxMeta) error
	FromAddress() gethcommon.Address
}

type ReportToEthMetadata func([]byte) (*txmgr.TxMeta, error)

func reportToEvmTxMetaNoop([]byte) (*txmgr.TxMeta, error) {
	return nil, nil
}

var _ ocr3types.ContractTransmitter[any] = &contractTransmitterOCR3[any]{}

type contractTransmitterOCR3[RI any] struct {
	contractAddress     gethcommon.Address
	contractABI         abi.ABI
	transmitter         Transmitter
	transmittedEventSig gethcommon.Hash
	lggr                logger.Logger
	reportToEvmTxMeta   ReportToEthMetadata
}

func NewOCR3ContractTransmitter[RI any](
	address gethcommon.Address,
	contractABI abi.ABI,
	transmitter Transmitter,
	lggr logger.Logger,
	reportToEvmTxMeta ReportToEthMetadata,
) (*contractTransmitterOCR3[RI], error) {
	transmitted, ok := contractABI.Events[TransmittedOCR3]
	if !ok {
		return nil, fmt.Errorf("abi missing transmitted event (name: %s)", TransmittedOCR3)
	}

	if reportToEvmTxMeta == nil {
		reportToEvmTxMeta = reportToEvmTxMetaNoop
	}
	return &contractTransmitterOCR3[RI]{
		contractAddress:     address,
		contractABI:         contractABI,
		transmitter:         transmitter,
		transmittedEventSig: transmitted.ID,
		lggr:                lggr,
		reportToEvmTxMeta:   reportToEvmTxMeta,
	}, nil
}

// FromAccount implements ocr3types.ContractTransmitter.
func (c *contractTransmitterOCR3[RI]) FromAccount() (types.Account, error) {
	return types.Account(c.transmitter.FromAddress().Hex()), nil
}

// Transmit implements ocr3types.ContractTransmitter.
func (c *contractTransmitterOCR3[RI]) Transmit(ctx context.Context, configDigest types.ConfigDigest, seqNum uint64, rwi ocr3types.ReportWithInfo[RI], sigs []types.AttributedOnchainSignature) error {
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
	epoch, round := uint64ToUint32AndUint8(seqNum)
	rawReportCtx := evmutil.RawReportContext(types.ReportContext{
		ReportTimestamp: types.ReportTimestamp{
			ConfigDigest: configDigest,
			Epoch:        epoch,
			Round:        round,
		},
		// ExtraData not used in OCR3
	})

	txMeta, err := c.reportToEvmTxMeta(rwi.Report)
	if err != nil {
		c.lggr.Warnw("failed to generate tx metadata for report", "err", err)
	}

	c.lggr.Debugw("Transmitting report", "report", hexutil.Encode(rwi.Report), "rawReportCtx", rawReportCtx, "contractAddress", c.contractAddress, "txMeta", txMeta)

	payload, err := c.contractABI.Pack("transmit", rawReportCtx, []byte(rwi.Report), rs, ss, vs)
	if err != nil {
		return fmt.Errorf("%w: abi.Pack failed with args: (%+v, %s, %+v, %+v, %+v)", err, rawReportCtx, hexutil.Encode(rwi.Report), rs, ss, vs)
	}

	c.lggr.Debugw("transmit payload", "payload", hexutil.Encode(payload))

	return errors.Wrap(c.transmitter.CreateEthTransaction(ctx, c.contractAddress, payload, txMeta), "failed to send Eth transaction")
}
