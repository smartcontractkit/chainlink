package evm_test

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	clcommontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	. "github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests" //nolint common practice to import test mods with .
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"

	commontestutils "github.com/smartcontractkit/chainlink-common/pkg/loop/testutils"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtxmgr "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	keytypes "github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	. "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/evmtesting" //nolint common practice to import test mods with .
)

const commonGasLimitOnEvms = uint64(4712388)

func TestChainReaderEventsInitValidation(t *testing.T) {
	tests := []struct {
		name                 string
		chainContractReaders map[string]types.ChainContractReader
		expectedError        error
	}{
		{
			name: "Invalid ABI",
			chainContractReaders: map[string]types.ChainContractReader{
				"InvalidContract": {
					ContractABI: "{invalid json}",
					Configs:     map[string]*types.ChainReaderDefinition{},
				},
			},
			expectedError: fmt.Errorf("failed to parse abi"),
		},
		{
			name: "Conflicting polling filter definitions",
			chainContractReaders: map[string]types.ChainContractReader{
				"ContractWithConflict": {
					ContractABI: "[]",
					Configs: map[string]*types.ChainReaderDefinition{
						"EventWithConflict": {
							ChainSpecificName: "EventName",
							ReadType:          types.Event,
							EventDefinitions: &types.EventDefinitions{
								PollingFilter: &types.PollingFilter{},
							},
						},
					},
					ContractPollingFilter: types.ContractPollingFilter{
						GenericEventNames: []string{"EventWithConflict"},
					},
				},
			},
			expectedError: fmt.Errorf(
				"%w: conflicting chain reader polling filter definitions for contract: %s event: %s, can't have polling filter defined both on contract and event level",
				clcommontypes.ErrInvalidConfig, "ContractWithConflict", "EventWithConflict"),
		},
		{
			name: "No polling filter defined",
			chainContractReaders: map[string]types.ChainContractReader{
				"ContractWithNoFilter": {
					ContractABI: "[]",
					Configs: map[string]*types.ChainReaderDefinition{
						"EventWithNoFilter": {
							ChainSpecificName: "EventName",
							ReadType:          types.Event,
						},
					},
				},
			},
			expectedError: fmt.Errorf(
				"%w: chain reader has no polling filter defined for contract: %s, event: %s",
				clcommontypes.ErrInvalidConfig, "ContractWithNoFilter", "EventWithNoFilter"),
		},
		{
			name: "Invalid chain reader definition read type",
			chainContractReaders: map[string]types.ChainContractReader{
				"ContractWithInvalidReadType": {
					ContractABI: "[]",
					Configs: map[string]*types.ChainReaderDefinition{
						"InvalidReadType": {
							ChainSpecificName: "InvalidName",
							ReadType:          types.ReadType(2),
						},
					},
				},
			},
			expectedError: fmt.Errorf(
				"%w: invalid chain reader definition read type",
				clcommontypes.ErrInvalidConfig),
		},
		{
			name: "Event not present in ABI",
			chainContractReaders: map[string]types.ChainContractReader{
				"ContractWithConflict": {
					ContractABI: "[{\"anonymous\":false,\"inputs\":[],\"name\":\"WrongEvent\",\"type\":\"event\"}]",
					Configs: map[string]*types.ChainReaderDefinition{
						"SomeEvent": {
							ChainSpecificName: "EventName",
							ReadType:          types.Event,
						},
					},
					ContractPollingFilter: types.ContractPollingFilter{
						GenericEventNames: []string{"SomeEvent"},
					},
				},
			},
			expectedError: fmt.Errorf(
				"%w: event %s doesn't exist",
				clcommontypes.ErrInvalidConfig, "EventName"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := evm.NewChainReaderService(testutils.Context(t), logger.NullLogger, nil, nil, nil, types.ChainReaderConfig{Contracts: tt.chainContractReaders})
			require.Error(t, err)
			if err != nil {
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			}
		})
	}
}

func TestChainReader(t *testing.T) {
	t.Parallel()
	it := &EVMChainReaderInterfaceTester[*testing.T]{Helper: &helper{}}
	it.Helper.Init(t)
	// add new subtests here so that it can be run on real chains too
	RunChainReaderEvmTests(t, it)
	RunChainReaderInterfaceTests[*testing.T](t, commontestutils.WrapChainReaderTesterForLoop(it))
}

type helper struct {
	sim    *backends.SimulatedBackend
	auth   *bind.TransactOpts
	pKey   *ecdsa.PrivateKey
	txm    evmtxmgr.TxManager
	client client.Client
}

func (h *helper) Init(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	h.pKey = privateKey

	h.auth, err = bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	require.NoError(t, err)

	h.Backend()
	h.txm = h.TXM(t)
	h.Commit()
}

func (h *helper) MustGenerateRandomKey(t *testing.T) ethkey.KeyV2 {
	return cltest.MustGenerateRandomKey(t)
}

func (h *helper) GasPriceBufferPercent() int64 {
	return 0
}

func (h *helper) SetupAuth(t *testing.T) *bind.TransactOpts {
	return h.auth
}

func (h *helper) Backend() bind.ContractBackend {
	if h.sim == nil {
		h.sim = backends.NewSimulatedBackend(
			core.GenesisAlloc{h.auth.From: {Balance: big.NewInt(math.MaxInt64)}}, commonGasLimitOnEvms*5000)
		cltest.Mine(h.sim, 1*time.Second)
	}

	return h.sim
}

func (h *helper) Commit() {
	h.sim.Commit()
}

func (h *helper) Client(t *testing.T) client.Client {
	if h.client != nil {
		return h.client
	}
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

func (h *helper) FromAddress() common.Address {
	return h.auth.From
}

func (h *helper) MaxWaitTimeForEvents() time.Duration {
	// From trial and error, when running on CI, sometimes the boxes get slow
	maxWaitTime := time.Second * 30
	maxWaitTimeStr, ok := os.LookupEnv("MAX_WAIT_TIME_FOR_EVENTS_S")
	if ok {
		waitS, err := strconv.ParseInt(maxWaitTimeStr, 10, 64)
		if err != nil {
			fmt.Printf("Error parsing MAX_WAIT_TIME_FOR_EVENTS_S: %v, defaulting to %v\n", err, maxWaitTime)
		}
		maxWaitTime = time.Second * time.Duration(waitS)
	}

	return maxWaitTime
}

func (h *helper) TXM(t *testing.T) evmtxmgr.TxManager {
	if h.txm != nil {
		return h.txm
	}

	client := h.Client(t)
	db := h.NewSqlxDB(t)

	clconfig := configtest.NewGeneralConfigSimulated(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Database.Listener.FallbackPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
		c.EVM[0].GasEstimator.EIP1559DynamicFees = ptr(true)
	})

	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, clconfig, h.sim, db, client)
	app.Start(h.Context(t))
	keyStore := app.KeyStore.Eth()

	keyStore.XXXTestingOnlyAdd(h.Context(t), keytypes.FromPrivateKey(h.pKey))
	require.NoError(t, keyStore.Add(h.Context(t), h.FromAddress(), h.ChainID()))
	require.NoError(t, keyStore.Enable(h.Context(t), h.FromAddress(), h.ChainID()))

	chain, err := app.GetRelayers().LegacyEVMChains().Get((*big.Int)(h.ChainID()).String())
	require.NoError(t, err)

	h.txm = chain.TxManager()
	return h.txm
}

func ptr[T any](v T) *T { return &v }
