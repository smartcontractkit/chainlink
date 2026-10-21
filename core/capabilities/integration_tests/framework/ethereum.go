package framework

import (
	"context"
	"crypto/ecdsa"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
)

type EthBlockchain struct {
	services.StateMachine
	evmtypes.Backend
	transactionOpts *bind.TransactOpts
	ownerKey        *ecdsa.PrivateKey

	blockTimeProcessingTime time.Duration

	stopCh services.StopChan
	wg     sync.WaitGroup
}

func NewEthBlockchain(t *testing.T, initialEth int, blockTimeProcessingTime time.Duration) *EthBlockchain {
	// The owner key is retained so that contracts requiring off-chain signatures from the deployer,
	// such as the workflow registry's owner link, can be driven from tests.
	ownerKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	transactOpts, err := bind.NewKeyedTransactorWithChainID(ownerKey, testutils.SimulatedChainID) // contract deployer and owner
	require.NoError(t, err)

	genesisData := types.GenesisAlloc{transactOpts.From: {Balance: assets.Ether(initialEth).ToInt()}}
	backend := cltest.NewSimulatedBackend(t, genesisData, ethconfig.Defaults.Miner.GasCeil)
	gethlog.SetDefault(gethlog.NewLogger(gethlog.NewTerminalHandlerWithLevel(os.Stderr, gethlog.LevelWarn, true)))
	backend.Commit()

	return &EthBlockchain{Backend: backend, stopCh: make(services.StopChan),
		blockTimeProcessingTime: blockTimeProcessingTime, transactionOpts: transactOpts, ownerKey: ownerKey}
}

func (b *EthBlockchain) Start(ctx context.Context) error {
	return b.StartOnce("EthBlockchain", func() error {
		b.wg.Go(func() {
			ticker := time.NewTicker(b.blockTimeProcessingTime)
			defer ticker.Stop()

			for {
				select {
				case <-b.stopCh:
					return
				case <-ctx.Done():
					return
				case <-ticker.C:
					b.Commit()
				}
			}
		})

		return nil
	})
}

func (b *EthBlockchain) Close() error {
	return b.StopOnce("EthBlockchain", func() error {
		close(b.stopCh)
		b.wg.Wait()
		return nil
	})
}

func (b *EthBlockchain) TransactionOpts() *bind.TransactOpts {
	return b.transactionOpts
}

// SignHash signs a 32 byte hash with the key owning the contracts deployed against this backend.
// The returned signature has a recovery ID of 0 or 1, callers needing an EIP-191 signature must
// add 27 to it.
func (b *EthBlockchain) SignHash(hash []byte) ([]byte, error) {
	return crypto.Sign(hash, b.ownerKey)
}
