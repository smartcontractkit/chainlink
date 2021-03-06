package cltest

import (
	"crypto/ecdsa"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

// newIdentity returns a go-ethereum abstraction of an ethereum account for
// interacting with contract golang wrappers
func NewSimulatedBackendIdentity(t *testing.T) *bind.TransactOpts {
	key, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate ethereum identity")
	return MustNewSimulatedBackendKeyedTransactor(t, key)
}

func NewApplicationWithConfigAndKeyOnSimulatedBlockchain(
	t testing.TB,
	tc *TestConfig,
	backend *backends.SimulatedBackend,
	flagsAndDeps ...interface{},
) (app *TestApplication, cleanup func()) {
	chainId := int(backend.Blockchain().Config().ChainID.Int64())
	tc.Config.Set("ETH_CHAIN_ID", chainId)

	client := &SimulatedBackendClient{b: backend, t: t, chainId: chainId}
	flagsAndDeps = append(flagsAndDeps, client)

	app, appCleanup := NewApplicationWithConfigAndKey(t, tc, flagsAndDeps...)
	err := app.Store.KeyStore.Unlock(Password)
	require.NoError(t, err)

	return app, func() { appCleanup(); client.Close() }
}

func MustNewSimulatedBackendKeyedTransactor(t *testing.T, key *ecdsa.PrivateKey) *bind.TransactOpts {
	t.Helper()

	return MustNewKeyedTransactor(t, key, 1337)
}

func MustNewKeyedTransactor(t *testing.T, key *ecdsa.PrivateKey, chainID int64) *bind.TransactOpts {
	t.Helper()

	transactor, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(chainID))
	require.NoError(t, err)

	return transactor
}

// Mine forces the simulated backend to produce a new block every 2 seconds
func Mine(backend *backends.SimulatedBackend) (stopMinning func()) {
	timer := time.NewTicker(2 * time.Second)
	chStop := make(chan struct{})
	go func() {
		for {
			select {
			case <-timer.C:
				backend.Commit()
			case <-chStop:
				return
			}
		}
	}()
	return func() { close(chStop); timer.Stop() }
}
