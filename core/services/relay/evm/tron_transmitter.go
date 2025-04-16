// TODO: Move this to chainlink-tron once chainlink-evm is fully extracted
package evm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	tronsdk "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	tron "github.com/smartcontractkit/chainlink-tron/relayer/ocr2"
	trontxm "github.com/smartcontractkit/chainlink-tron/relayer/txm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// We implement the TRON TXM cache API using EVM's contract transmitter
var _ tron.TransmissionsCache = (*tronTransmissionsCache)(nil)

type tronTransmissionsCache struct {
	evmTransmitter ContractTransmitter
}

func NewTronTransmissionsCache(evmTransmitter ContractTransmitter) tron.TransmissionsCache {
	return &tronTransmissionsCache{
		evmTransmitter: evmTransmitter,
	}
}

func (t *tronTransmissionsCache) LatestTransmissionDetails(ctx context.Context) (types.ConfigDigest, uint32, uint8, *big.Int, time.Time, error) {
	configDigest, epoch, err := t.evmTransmitter.LatestConfigDigestAndEpoch(ctx)
	if err != nil {
		return types.ConfigDigest{}, 0, 0, nil, time.Time{}, err
	}
	return configDigest, epoch, 0, nil, time.Time{}, nil
}

// TronContractTransmitterOpts contains the configuration options for creating a Tron contract transmitter
type TronContractTransmitterOpts struct {
	Logger             logger.Logger
	TransmissionsCache tron.TransmissionsCache
	Keystore           keys.Store
	ConfigWatcher      *configWatcher
	OCRTransmitterOpts []OCRTransmitterOption
}

// NewTronContractTransmitter creates a new ContractTransmitter for Tron chains
func NewTronContractTransmitter(ctx context.Context, opts TronContractTransmitterOpts) (ContractTransmitter, error) {
	// On TRON, get the chain specific txm
	chain, ok := opts.ConfigWatcher.chain.(legacyevm.ChainTronSupport)
	if !ok {
		return nil, fmt.Errorf("chain %s does not support TRON", opts.ConfigWatcher.chain.ID())
	}

	senderAddress, err := opts.Keystore.GetNextAddress(ctx)
	if err != nil {
		return nil, err
	}

	// Construct the Tron contract transmitter, it's slightly different from the EVM contract transmitter and due to mismatching types we have to apply the transmitter options manually
	transmitter := tron.NewOCRContractTransmitter(ctx, opts.TransmissionsCache, tronsdk.EVMAddressToAddress(opts.ConfigWatcher.contractAddress), tronsdk.EVMAddressToAddress(senderAddress), chain.GetTronTXM(), opts.Logger).WithEthereumKeystore()
	transmitterOptions := &transmitterOps{
		reportToEvmTxMeta: nil,
		excludeSigs:       false,
		retention:         0,
		maxLogsKept:       0,
	}

	for _, opt := range opts.OCRTransmitterOpts {
		opt(transmitterOptions)
	}

	if transmitterOptions.excludeSigs {
		opts.Logger.Info("Excluding signatures from transmissions")
		transmitter = transmitter.WithExcludeSignatures()
	}
	if transmitterOptions.reportToEvmTxMeta != nil {
		opts.Logger.Info("Using EVM TxMeta for Tron Transmissions")
		transmitter = transmitter.WithReportToEthMetadata(transmitterOptions.reportToEvmTxMeta)
	}

	tronTransmitter := newTronTransmitterWrapper(transmitter, chain.GetTronTXM())

	return tronTransmitter, nil
}

var _ ContractTransmitter = (*tronTransmitterWrapper)(nil)

// Simple wrapper around the tron.ContractTransmitter to provide start / stop hooks to the tron txm
type tronTransmitterWrapper struct {
	transmitter tron.ContractTransmitter
	txm         *trontxm.TronTxm
}

func newTronTransmitterWrapper(transmitter tron.ContractTransmitter, txm *trontxm.TronTxm) ContractTransmitter {
	return &tronTransmitterWrapper{
		transmitter: transmitter,
		txm:         txm,
	}
}

// The Tron Transmitter doesn't close the txm, so we'll close it here
// NOTE: This is called by the Close() method of either the CommitProvider or the ExecProvider
func (t *tronTransmitterWrapper) Close() error {
	err := t.txm.Close()
	if err != nil {
		return err
	}

	return t.transmitter.Close()
}

func (t *tronTransmitterWrapper) FromAccount(ctx context.Context) (types.Account, error) {
	return t.transmitter.FromAccount(ctx)
}

func (t *tronTransmitterWrapper) HealthReport() map[string]error {
	return t.transmitter.HealthReport()
}

func (t *tronTransmitterWrapper) LatestConfigDigestAndEpoch(ctx context.Context) (types.ConfigDigest, uint32, error) {
	return t.transmitter.LatestConfigDigestAndEpoch(ctx)
}

func (t *tronTransmitterWrapper) Name() string {
	return t.transmitter.Name()
}

func (t *tronTransmitterWrapper) Ready() error {
	return t.transmitter.Ready()
}

// As the Tron Transmitter doesn't start the txm, we need to start it within the start hook
func (t *tronTransmitterWrapper) Start(ctx context.Context) error {
	t.txm.Logger.Info("Starting Tron TXM")

	// NOTE: The txm needs to be started before the contract transmitter, so we check if it's ready and start it if it's not already started
	if err := t.txm.Start(ctx); err != nil {
		return err
	}

	return t.transmitter.Start(ctx)
}

func (t *tronTransmitterWrapper) Transmit(ctx context.Context, reportCtx types.ReportContext, report types.Report, signatures []types.AttributedOnchainSignature) error {
	return t.transmitter.Transmit(ctx, reportCtx, report, signatures)
}
