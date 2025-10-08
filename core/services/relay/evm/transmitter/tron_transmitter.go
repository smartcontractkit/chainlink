// TODO: Move this to chainlink-tron once chainlink-evm is fully extracted
package transmitter

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	tronsdk "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	tron "github.com/smartcontractkit/chainlink-tron/relayer/ocr2"
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
		return types.ConfigDigest{}, 0, 0, nil, time.Time{}, fmt.Errorf("failed to proxy the call to the EVM transmitter: %w", err)
	}
	return configDigest, epoch, 0, nil, time.Time{}, nil
}

// TronContractTransmitterOpts contains the configuration options for creating a Tron contract transmitter
type TronContractTransmitterOpts struct {
	Logger             logger.Logger
	TransmissionsCache tron.TransmissionsCache
	Keystore           keys.Store
	Chain              legacyevm.Chain
	ContractAddress    common.Address
	OCRTransmitterOpts []OCRTransmitterOption
}

// NewTronContractTransmitter creates a new ContractTransmitter for Tron chains
func NewTronContractTransmitter(ctx context.Context, opts TronContractTransmitterOpts) (ContractTransmitter, error) {
	// On TRON, get the chain specific txm
	chain, ok := opts.Chain.(legacyevm.ChainTronSupport)
	if !ok {
		return nil, fmt.Errorf("chain %s does not support TRON", opts.Chain.ID())
	}

	senderAddress, err := opts.Keystore.GetNextAddress(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender address: %w", err)
	}

	// Construct the Tron contract transmitter, it's slightly different from the EVM contract transmitter and due to mismatching types we have to apply the transmitter options manually
	transmitterTron := tron.NewOCRContractTransmitter(ctx, opts.TransmissionsCache, tronsdk.EVMAddressToAddress(opts.ContractAddress), tronsdk.EVMAddressToAddress(senderAddress), chain.GetTronTXM(), opts.Logger)

	// Use the EVM keystore for the transmitter
	transmitterTron.WithEthereumKeystore()

	if HasExcludeSignatures(opts.OCRTransmitterOpts) {
		opts.Logger.Info("Excluding signatures from transmissions")
		transmitterTron.WithExcludeSignatures()
	}

	return transmitterTron, nil
}
