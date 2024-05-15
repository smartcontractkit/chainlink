package evm_test

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/evmtesting"
)

const commonGasLimitOnEvms = uint64(4712388)

func TestChainReader(t *testing.T) {
	t.Parallel()
	evmtesting.RunChainReaderTests(t, &evmtesting.EvmChainReaderInterfaceTester[*testing.T]{Helper: &helper{}}, true)
}

type helper struct {
	sim  *backends.SimulatedBackend
	auth *bind.TransactOpts
}

func (h *helper) SetupAuth(t *testing.T) *bind.TransactOpts {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	h.auth, err = bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	require.NoError(t, err)

	h.Backend()
	h.Commit()
	return h.auth
}

func (h *helper) Backend() bind.ContractBackend {
	if h.sim == nil {
		h.sim = backends.NewSimulatedBackend(
			core.GenesisAlloc{h.auth.From: {Balance: big.NewInt(math.MaxInt64)}}, commonGasLimitOnEvms*5000)
	}

	return h.sim
}

func (h *helper) Commit() {
	h.sim.Commit()
}

func (h *helper) Client(t *testing.T) client.Client {
	return client.NewSimulatedBackendClient(t, h.sim, big.NewInt(1337))
}

func (h *helper) ChainID() *big.Int {
	return testutils.SimulatedChainID
}

func (h *helper) NewSqlxDB(t *testing.T) *sqlx.DB {
	return pgtest.NewSqlxDB(t)
}

func (h *helper) Context(t *testing.T) context.Context {
	return testutils.Context(t)
}

func (h *helper) MaxWaitTimeForEvents() time.Duration {
	// From trial and error, when running on CI, sometimes the boxes get slow
	maxWaitTime := time.Second * 20
	maxWaitTimeStr, ok := os.LookupEnv("MAX_WAIT_TIME_FOR_EVENTS_S")
	if ok {
		wiatS, err := strconv.ParseInt(maxWaitTimeStr, 10, 64)
		if err != nil {
			fmt.Printf("Error parsing MAX_WAIT_TIME_FOR_EVENTS_S: %v, defaulting to %v\n", err, maxWaitTime)
		}
		maxWaitTime = time.Second * time.Duration(wiatS)
	}

	return maxWaitTime
}
