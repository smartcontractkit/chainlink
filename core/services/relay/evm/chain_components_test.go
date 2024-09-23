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
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	commontestutils "github.com/smartcontractkit/chainlink-common/pkg/loop/testutils"
	clcommontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtxmgr "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	keytypes "github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	. "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/evmtesting" //nolint common practice to import test mods with .
)

const commonGasLimitOnEvms = uint64(4712388)

func TestContractReaderEventsInitValidation(t *testing.T) {
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

func TestChainComponents(t *testing.T) {
	t.Parallel()
	it := &EVMChainComponentsInterfaceTester[*testing.T]{Helper: &helper{}}

	it.Helper.Init(t)

	// add new subtests here so that it can be run on real chains too
	RunChainComponentsEvmTests(t, it)
	RunChainComponentsInLoopEvmTests[*testing.T](t, commontestutils.WrapContractReaderTesterForLoop(it))
}

type helper struct {
	sim         *backends.SimulatedBackend
	accounts    []*bind.TransactOpts
	deployerKey *ecdsa.PrivateKey
	senderKey   *ecdsa.PrivateKey
	txm         evmtxmgr.TxManager
	client      client.Client
	db          *sqlx.DB
}

func (h *helper) Init(t *testing.T) {
	h.SetupKeys(t)

	h.accounts = h.Accounts(t)

	h.db = pgtest.NewSqlxDB(t)

	h.Backend()
	h.client = h.Client(t)

	h.txm = h.TXM(t, h.client)
	h.Commit()
}

func (h *helper) SetupKeys(t *testing.T) {
	deployerPkey, err := crypto.GenerateKey()
	require.NoError(t, err)
	h.deployerKey = deployerPkey

	senderPkey, err := crypto.GenerateKey()
	require.NoError(t, err)
	h.senderKey = senderPkey
}

func (h *helper) Accounts(t *testing.T) []*bind.TransactOpts {
	if h.accounts != nil {
		return h.accounts
	}
	deployer, err := bind.NewKeyedTransactorWithChainID(h.deployerKey, big.NewInt(1337))
	require.NoError(t, err)

	sender, err := bind.NewKeyedTransactorWithChainID(h.senderKey, big.NewInt(1337))
	require.NoError(t, err)

	return []*bind.TransactOpts{deployer, sender}
}

func (h *helper) MustGenerateRandomKey(t *testing.T) ethkey.KeyV2 {
	return cltest.MustGenerateRandomKey(t)
}

func (h *helper) GasPriceBufferPercent() int64 {
	return 0
}

func (h *helper) Backend() bind.ContractBackend {
	if h.sim == nil {
		h.sim = backends.NewSimulatedBackend(
			core.GenesisAlloc{h.accounts[0].From: {Balance: big.NewInt(math.MaxInt64)}, h.accounts[1].From: {Balance: big.NewInt(math.MaxInt64)}}, commonGasLimitOnEvms*5000)
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

func (h *helper) ChainReaderEVMClient(ctx context.Context, t *testing.T, ht logpoller.HeadTracker, conf types.ChainReaderConfig) client.Client {
	// wrap the client so that we can mock historical contract state
	cwh := &evm.ClientWithContractHistory{Client: h.Client(t), HT: ht}
	require.NoError(t, cwh.Init(ctx, conf))
	return cwh
}

func (h *helper) WrappedChainWriter(cw clcommontypes.ChainWriter, client client.Client) clcommontypes.ChainWriter {
	cwhw := evm.NewChainWriterHistoricalWrapper(cw, client.(*evm.ClientWithContractHistory))
	return cwhw
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

func (h *helper) TXM(t *testing.T, client client.Client) evmtxmgr.TxManager {
	if h.txm != nil {
		return h.txm
	}
	db := h.db

	clconfig := configtest.NewGeneralConfigSimulated(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Database.Listener.FallbackPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
		c.EVM[0].GasEstimator.EIP1559DynamicFees = ptr(true)
	})

	clconfig.EVMConfigs()[0].GasEstimator.PriceMax = assets.GWei(100)

	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, clconfig, h.sim, db, client)
	err := app.Start(h.Context(t))
	require.NoError(t, err)

	keyStore := app.KeyStore.Eth()

	keyStore.XXXTestingOnlyAdd(h.Context(t), keytypes.FromPrivateKey(h.deployerKey))
	require.NoError(t, keyStore.Add(h.Context(t), h.accounts[0].From, h.ChainID()))
	require.NoError(t, keyStore.Enable(h.Context(t), h.accounts[0].From, h.ChainID()))

	keyStore.XXXTestingOnlyAdd(h.Context(t), keytypes.FromPrivateKey(h.senderKey))
	require.NoError(t, keyStore.Add(h.Context(t), h.accounts[1].From, h.ChainID()))
	require.NoError(t, keyStore.Enable(h.Context(t), h.accounts[1].From, h.ChainID()))

	chain, err := app.GetRelayers().LegacyEVMChains().Get((h.ChainID()).String())
	require.NoError(t, err)

	h.txm = chain.TxManager()
	return h.txm
}

func ptr[T any](v T) *T { return &v }
