package evm

import (
	"context"
	"database/sql"
	"encoding/hex"
	errors2 "errors"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

// TODO: Remove when new dual transmitter contracts are merged
var dtABI = `[{"inputs":[{"internalType":"bytes32[3]","name":"reportContext","type":"bytes32[3]"},{"internalType":"bytes","name":"report","type":"bytes"},{"internalType":"bytes32[]","name":"rs","type":"bytes32[]"},{"internalType":"bytes32[]","name":"ss","type":"bytes32[]"},{"internalType":"bytes32","name":"rawVs","type":"bytes32"}],"name":"transmitSecondary","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

var _ ContractTransmitter = (*dualContractTransmitter)(nil)

type dualContractTransmitter struct {
	contractAddress     gethcommon.Address
	contractABI         abi.ABI
	dualTransmissionABI abi.ABI
	transmitter         Transmitter
	transmittedEventSig common.Hash
	contractReader      contractReader
	lp                  logpoller.LogPoller
	lggr                logger.Logger
	ks                  keystore.Eth
	// Options
	transmitterOptions *transmitterOps
}

var dualTransmissionABI = sync.OnceValue(func() abi.ABI {
	dualTransmissionABI, err := abi.JSON(strings.NewReader(dtABI))
	if err != nil {
		panic(fmt.Errorf("failed to parse dualTransmission ABI: %w", err))
	}
	return dualTransmissionABI
})

func NewOCRDualContractTransmitter(
	ctx context.Context,
	address gethcommon.Address,
	caller contractReader,
	contractABI abi.ABI,
	transmitter Transmitter,
	lp logpoller.LogPoller,
	lggr logger.Logger,
	ethKeystore keystore.Eth,
	opts ...OCRTransmitterOption,
) (*dualContractTransmitter, error) {
	transmitted, ok := contractABI.Events["Transmitted"]
	if !ok {
		return nil, errors.New("invalid ABI, missing transmitted")
	}

	newContractTransmitter := &dualContractTransmitter{
		contractAddress:     address,
		contractABI:         contractABI,
		dualTransmissionABI: dualTransmissionABI(),
		transmitter:         transmitter,
		transmittedEventSig: transmitted.ID,
		lp:                  lp,
		contractReader:      caller,
		lggr:                logger.Named(lggr, "OCR2DualContractTransmitter"),
		ks:                  ethKeystore,
		transmitterOptions: &transmitterOps{
			reportToEvmTxMeta: reportToEvmTxMetaNoop,
			excludeSigs:       false,
			retention:         0,
			maxLogsKept:       0,
		},
	}

	for _, opt := range opts {
		opt(newContractTransmitter.transmitterOptions)
	}

	err := lp.RegisterFilter(ctx, logpoller.Filter{Name: transmitterFilterName(address), EventSigs: []common.Hash{transmitted.ID}, Addresses: []common.Address{address}, Retention: newContractTransmitter.transmitterOptions.retention, MaxLogsKept: newContractTransmitter.transmitterOptions.maxLogsKept})
	if err != nil {
		return nil, err
	}
	return newContractTransmitter, nil
}

// Transmit sends the report to the on-chain smart contract's Transmit method.
func (oc *dualContractTransmitter) Transmit(ctx context.Context, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signatures []ocrtypes.AttributedOnchainSignature) error {
	var rs [][32]byte
	var ss [][32]byte
	var vs [32]byte
	if len(signatures) > 32 {
		return errors.New("too many signatures, maximum is 32")
	}
	for i, as := range signatures {
		r, s, v, err := evmutil.SplitSignature(as.Signature)
		if err != nil {
			panic("eventTransmit(ev): error in SplitSignature")
		}
		if !oc.transmitterOptions.excludeSigs {
			rs = append(rs, r)
			ss = append(ss, s)
			vs[i] = v
		}
	}
	rawReportCtx := evmutil.RawReportContext(reportCtx)

	txMeta, err := oc.transmitterOptions.reportToEvmTxMeta(report)
	if err != nil {
		oc.lggr.Warnw("failed to generate tx metadata for report", "err", err)
	}

	oc.lggr.Debugw("Transmitting report", "report", hex.EncodeToString(report), "rawReportCtx", rawReportCtx, "contractAddress", oc.contractAddress, "txMeta", txMeta)

	// Primary transmission
	payload, err := oc.contractABI.Pack("transmit", rawReportCtx, []byte(report), rs, ss, vs)
	if err != nil {
		return errors.Wrap(err, "abi.Pack failed")
	}

	transactionErr := errors.Wrap(oc.transmitter.CreateEthTransaction(ctx, oc.contractAddress, payload, txMeta), "failed to send primary Eth transaction")

	oc.lggr.Debugw("Created primary transaction", "error", transactionErr)

	// Secondary transmission
	secondaryPayload, err := oc.dualTransmissionABI.Pack("transmitSecondary", rawReportCtx, []byte(report), rs, ss, vs)
	if err != nil {
		return errors.Wrap(err, "transmitSecondary abi.Pack failed")
	}

	err = errors.Wrap(oc.transmitter.CreateSecondaryEthTransaction(ctx, secondaryPayload, txMeta), "failed to send secondary Eth transaction")
	oc.lggr.Debugw("Created secondary transaction", "error", err)
	return errors2.Join(transactionErr, err)
}

// LatestConfigDigestAndEpoch retrieves the latest config digest and epoch from the OCR2 contract.
// It is plugin independent, in particular avoids use of the plugin specific generated evm wrappers
// by using the evm client Call directly for functions/events that are part of OCR2Abstract.
func (oc *dualContractTransmitter) LatestConfigDigestAndEpoch(ctx context.Context) (ocrtypes.ConfigDigest, uint32, error) {
	latestConfigDigestAndEpoch, err := callContract(ctx, oc.contractAddress, oc.contractABI, "latestConfigDigestAndEpoch", nil, oc.contractReader)
	if err != nil {
		return ocrtypes.ConfigDigest{}, 0, err
	}
	// Panic on these conversions erroring, would mean a broken contract.
	scanLogs := *abi.ConvertType(latestConfigDigestAndEpoch[0], new(bool)).(*bool)
	configDigest := *abi.ConvertType(latestConfigDigestAndEpoch[1], new([32]byte)).(*[32]byte)
	epoch := *abi.ConvertType(latestConfigDigestAndEpoch[2], new(uint32)).(*uint32)
	if !scanLogs {
		return configDigest, epoch, nil
	}

	// Otherwise, we have to scan for the logs.
	latest, err := oc.lp.LatestLogByEventSigWithConfs(ctx, oc.transmittedEventSig, oc.contractAddress, 1)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No transmissions yet
			return configDigest, 0, nil
		}
		return ocrtypes.ConfigDigest{}, 0, err
	}
	return parseTransmitted(latest.Data)
}

// FromAccount returns the account from which the transmitter invokes the contract
func (oc *dualContractTransmitter) FromAccount(ctx context.Context) (ocrtypes.Account, error) {
	return ocrtypes.Account(oc.transmitter.FromAddress(ctx).String()), nil
}

func (oc *dualContractTransmitter) lockTransmitters(ctx context.Context) error {
	err := oc.lockPrimary(ctx)
	if err != nil {
		return err
	}
	err = oc.lockSecondary(ctx)
	if err != nil {
		return multierr.Append(err, oc.unlockPrimary(ctx))
	}
	return nil
}

func (oc *dualContractTransmitter) unlockTransmitters(ctx context.Context) error {
	return multierr.Append(oc.unlockPrimary(ctx), oc.unlockSecondary(ctx))
}

func (oc *dualContractTransmitter) unlockPrimary(ctx context.Context) error {
	primaryAddress := oc.transmitter.FromAddress(ctx)
	rmPrimary, err := oc.ks.GetResourceMutex(ctx, primaryAddress)
	if err != nil {
		return err
	}
	err = rmPrimary.Unlock(keystore.TXMv1)
	if err != nil {
		return err
	}
	oc.lggr.Debugf("Key %s has been unlocked for TXMv1", primaryAddress.String())
	return nil
}
func (oc *dualContractTransmitter) unlockSecondary(ctx context.Context) error {
	secondaryAddress, err := oc.transmitter.SecondaryFromAddress(ctx)
	if err != nil {
		return err
	}
	rmSecondary, err := oc.ks.GetResourceMutex(ctx, secondaryAddress)
	if err != nil {
		return err
	}
	err = rmSecondary.Unlock(keystore.TXMv2)
	if err != nil {
		return err
	}
	oc.lggr.Debugf("Key %s has been unlocked for TXMv2", secondaryAddress.String())
	return nil
}

func (oc *dualContractTransmitter) lockPrimary(ctx context.Context) error {
	primaryAddress := oc.transmitter.FromAddress(ctx)
	rmPrimary, err := oc.ks.GetResourceMutex(ctx, primaryAddress)
	if err != nil {
		return err
	}
	err = rmPrimary.TryLock(keystore.TXMv1)
	if err != nil {
		return err
	}
	oc.lggr.Debugf("Key %s has been locked for TXMv1", primaryAddress.String())
	return nil
}
func (oc *dualContractTransmitter) lockSecondary(ctx context.Context) error {
	secondaryAddress, err := oc.transmitter.SecondaryFromAddress(ctx)
	if err != nil {
		return err
	}
	rmSecondary, err := oc.ks.GetResourceMutex(ctx, secondaryAddress)
	if err != nil {
		return err
	}
	err = rmSecondary.TryLock(keystore.TXMv2)
	if err != nil {
		return err
	}
	oc.lggr.Debugf("Key %s has been locked for TXMv2", secondaryAddress.String())
	return nil
}

func (oc *dualContractTransmitter) Start(ctx context.Context) error {
	return oc.lockTransmitters(ctx)
}
func (oc *dualContractTransmitter) Close() error {
	return oc.unlockTransmitters(context.Background())
}

// Has no state/lifecycle so it's always healthy and ready
func (oc *dualContractTransmitter) Ready() error { return nil }
func (oc *dualContractTransmitter) HealthReport() map[string]error {
	return map[string]error{oc.Name(): nil}
}
func (oc *dualContractTransmitter) Name() string { return oc.lggr.Name() }
