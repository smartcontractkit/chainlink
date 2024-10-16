package framework

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

type EthBlockchain struct {
	services.StateMachine
	*backends.SimulatedBackend
	transactionOpts *bind.TransactOpts

	blockTimeProcessingTime time.Duration

	stopCh services.StopChan
	wg     sync.WaitGroup
}

func NewEthBlockchain(t *testing.T, initialEth int, blockTimeProcessingTime time.Duration) *EthBlockchain {
	transactOpts := testutils.MustNewSimTransactor(t) // config contract deployer and owner
	genesisData := core.GenesisAlloc{transactOpts.From: {Balance: assets.Ether(initialEth).ToInt()}}
	//nolint:gosec // disable G115
	backend := cltest.NewSimulatedBackend(t, genesisData, uint32(ethconfig.Defaults.Miner.GasCeil))
	gethlog.SetDefault(gethlog.NewLogger(gethlog.NewTerminalHandlerWithLevel(os.Stderr, gethlog.LevelWarn, true)))
	backend.Commit()

	return &EthBlockchain{SimulatedBackend: backend, stopCh: make(services.StopChan),
		blockTimeProcessingTime: blockTimeProcessingTime, transactionOpts: transactOpts}
}

func (b *EthBlockchain) Start(ctx context.Context) error {
	return b.StartOnce("EthBlockchain", func() error {
		b.wg.Add(1)
		go func() {
			defer b.wg.Done()
			ticker := time.NewTicker(b.blockTimeProcessingTime)
			defer ticker.Stop()

			for {
				select {
				case <-b.stopCh:
					return
				case <-ctx.Done():
					return
				case <-ticker.C:
					b.SimulatedBackend.Commit()
				}
			}
		}()

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
