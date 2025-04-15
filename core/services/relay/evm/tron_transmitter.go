package evm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	tronsdk "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	tronkeystore "github.com/smartcontractkit/chainlink-tron/relayer/keystore"
	tron "github.com/smartcontractkit/chainlink-tron/relayer/ocr2"
	tronclient "github.com/smartcontractkit/chainlink-tron/relayer/sdk"
	trontxm "github.com/smartcontractkit/chainlink-tron/relayer/txm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type tronEVMTransmitterWrapper struct {
	evmTransmitter ContractTransmitter
}

func (t *tronEVMTransmitterWrapper) LatestTransmissionDetails(ctx context.Context) (types.ConfigDigest, uint32, uint8, *big.Int, time.Time, error) {
	configDigest, epoch, err := t.evmTransmitter.LatestConfigDigestAndEpoch(ctx)
	if err != nil {
		return types.ConfigDigest{}, 0, 0, nil, time.Time{}, err
	}
	return configDigest, epoch, 0, nil, time.Time{}, nil
}

func NewTronContractTransmitter(ctx context.Context, lggr logger.Logger, evmTransmitter ContractTransmitter, ethKeystore keys.Store, configWatcher *configWatcher, opts configTransmitterOpts, transmissionContractABI abi.ABI, ocrTransmitterOpts ...OCRTransmitterOption) (ContractTransmitter, error) {
	// On TRON, get the (extra) nodes information from the chain
	chain, ok := configWatcher.chain.(legacyevm.ChainTronSupport)
	if !ok {
		return nil, fmt.Errorf("chain %s does not support TRON", configWatcher.chain.ID())
	}

	// Check there is at least one node
	if len(chain.Nodes()) == 0 {
		return nil, fmt.Errorf("no nodes found for chain %s", configWatcher.chain.ID())
	}

	// Use an (extra) write-specific node URL for the Tron client
	// Notice: TronClient is not multinode aware, so we need to use the first node
	nodeURL := (*chain.Nodes()[0]).HTTPURLExtraWrite.URL()
	tronClient, err := tronclient.CreateFullNodeClient(nodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create Tron client: %w", err)
	}

	// Start the Tron TXM
	tronTXM := trontxm.New(lggr, tronkeystore.NewLoopKeystoreAdapter(ethKeystore), tronClient, trontxm.TronTxmConfig{
		EnergyMultiplier: 3,
		StatusChecker:    true,
	})
	tronTXM.Start(ctx)

	tronEVMTransmitterWrapper := &tronEVMTransmitterWrapper{
		evmTransmitter: evmTransmitter,
	}

	senderAddress, err := ethKeystore.GetNextAddress(ctx)
	if err != nil {
		return nil, err
	}

	// Construct the Tron contract transmitter, it's slightly different from the EVM contract transmitter and due to mismatching types we have to apply the transmitter options manually
	transmitter := tron.NewOCRContractTransmitter(ctx, tronEVMTransmitterWrapper, tronsdk.EVMAddressToAddress(configWatcher.contractAddress), tronsdk.EVMAddressToAddress(senderAddress), tronTXM, lggr).WithEthereumKeystore()
	transmitterOptions := &transmitterOps{
		reportToEvmTxMeta: nil,
		excludeSigs:       false,
		retention:         0,
		maxLogsKept:       0,
	}

	for _, opt := range ocrTransmitterOpts {
		opt(transmitterOptions)
	}

	if transmitterOptions.excludeSigs {
		lggr.Info("Excluding signatures from transmissions")
		transmitter = transmitter.WithExcludeSignatures()
	}
	if transmitterOptions.reportToEvmTxMeta != nil {
		lggr.Info("Using EVM TxMeta for Tron Transmissions")
		transmitter = transmitter.WithReportToEthMetadata(transmitterOptions.reportToEvmTxMeta)
	}

	return transmitter, nil
}
