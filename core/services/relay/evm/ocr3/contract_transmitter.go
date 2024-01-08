package ocr3

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type Transmitter interface {
	CreateEthTransaction(ctx context.Context, toAddress gethcommon.Address, payload []byte, txMeta *txmgr.TxMeta) error
	FromAddress() gethcommon.Address
}

type ReportToEthMetadata func([]byte) (*txmgr.TxMeta, error)

func reportToEvmTxMetaNoop([]byte) (*txmgr.TxMeta, error) {
	return nil, nil
}

func transmitterFilterName(addr gethcommon.Address) string {
	return logpoller.FilterName("OCR3 ContractTransmitter", addr.String())
}

var _ ocr3types.ContractTransmitter[any] = &contractTransmitterOCR3[any]{}

type contractTransmitterOCR3[RI any] struct {
	contractAddress     gethcommon.Address
	contractABI         abi.ABI
	transmitter         Transmitter
	transmittedEventSig gethcommon.Hash
	lp                  logpoller.LogPoller
	lggr                logger.Logger
	reportToEvmTxMeta   ReportToEthMetadata
}

func NewOCR3ContractTransmitter[RI any](
	address gethcommon.Address,
	contractABI abi.ABI,
	transmitter Transmitter,
	lp logpoller.LogPoller,
	lggr logger.Logger,
	reportToEvmTxMeta ReportToEthMetadata,
) (*contractTransmitterOCR3[RI], error) {
	transmitted, ok := contractABI.Events["Transmitted"]
	if !ok {
		return nil, fmt.Errorf("abi missing Transmitted event")
	}

	err := lp.RegisterFilter(logpoller.Filter{
		Name:      transmitterFilterName(address),
		EventSigs: []gethcommon.Hash{transmitted.ID},
		Addresses: []gethcommon.Address{address},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register filter: %w", err)
	}
	if reportToEvmTxMeta == nil {
		reportToEvmTxMeta = reportToEvmTxMetaNoop
	}
	return &contractTransmitterOCR3[RI]{
		contractAddress:     address,
		contractABI:         contractABI,
		transmitter:         transmitter,
		transmittedEventSig: transmitted.ID,
		lp:                  lp,
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
			panic("eventTransmit(ev): error in SplitSignature")
		}
		rs = append(rs, r)
		ss = append(ss, s)
		vs[i] = v
	}

	// report ctx for OCR3 consists of the following
	// reportContext[0]: ConfigDigest
	// reportContext[1]: 24 byte padding, 8 byte sequence number
	// reportContext[2]: unused
	var rawReportCtx [3][32]byte
	copy(rawReportCtx[0][:], configDigest[:])
	binary.BigEndian.PutUint64(rawReportCtx[1][24:], seqNum)

	txMeta, err := c.reportToEvmTxMeta(rwi.Report)
	if err != nil {
		c.lggr.Warnw("failed to generate tx metadata for report", "err", err)
	}

	c.lggr.Debugw("Transmitting report", "report", hex.EncodeToString(rwi.Report), "rawReportCtx", rawReportCtx, "contractAddress", c.contractAddress, "txMeta", txMeta)

	payload, err := c.contractABI.Pack("transmit", rawReportCtx, []byte(rwi.Report), rs, ss, vs)
	if err != nil {
		return fmt.Errorf("%w: abi.Pack failed with args: (%+v, %s, %+v, %+v, %+v)", err, rawReportCtx, hex.EncodeToString(rwi.Report), rs, ss, vs)
	}

	c.lggr.Debugw("payload", "payload", hex.EncodeToString(payload))

	return errors.Wrap(c.transmitter.CreateEthTransaction(ctx, c.contractAddress, payload, txMeta), "failed to send Eth transaction")
}
