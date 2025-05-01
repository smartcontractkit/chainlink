package evm_test

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontestutils "github.com/smartcontractkit/chainlink-common/pkg/loop/testutils"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	clcommontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	clevmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"

	lpMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	evmtxmgr "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"

	. "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/evmtesting" //nolint:revive // dot-imports
)

const commonGasLimitOnEvms = uint64(4712388)
const finalityDepth = 4

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
				"%w: event %q doesn't exist",
				clcommontypes.ErrInvalidConfig, "EventName"),
		},
		{
			name: "Event has a unnecessary data word index override",
			chainContractReaders: map[string]types.ChainContractReader{
				"ContractWithConflict": {
					ContractABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"someDW\",\"type\":\"address\"}],\"name\":\"EventName\",\"type\":\"event\"}]",
					ContractPollingFilter: types.ContractPollingFilter{
						GenericEventNames: []string{"SomeEvent"},
					},
					Configs: map[string]*types.ChainReaderDefinition{
						"SomeEvent": {
							ChainSpecificName: "EventName",
							ReadType:          types.Event,

							EventDefinitions: &types.EventDefinitions{
								GenericDataWordDetails: map[string]types.DataWordDetail{
									"DW": {
										Name:  "someDW",
										Index: ptr(0),
									},
								},
							},
						},
					},
				},
			},
			expectedError: fmt.Errorf("failed to init dw querying for event: %q, err: data word: %q at index: %d details, were calculated automatically and shouldn't be manully overridden by cfg",
				"SomeEvent", "DW", 0),
		},
		{
			name: "Event has a bad type defined in data word detail override config",
			chainContractReaders: map[string]types.ChainContractReader{
				"ContractWithConflict": {
					ContractABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"someDW\",\"type\":\"string\"}],\"name\":\"EventName\",\"type\":\"event\"}]",
					ContractPollingFilter: types.ContractPollingFilter{
						GenericEventNames: []string{"SomeEvent"},
					},
					Configs: map[string]*types.ChainReaderDefinition{
						"SomeEvent": {
							ChainSpecificName: "EventName",
							ReadType:          types.Event,

							EventDefinitions: &types.EventDefinitions{
								GenericDataWordDetails: map[string]types.DataWordDetail{
									"DW": {
										Name:  "someDW",
										Index: ptr(0),
										Type:  "abcdefg",
									},
								},
							},
						},
					},
				},
			},
			expectedError: fmt.Errorf("failed to init dw querying for event: %q, err: bad abi type: \"abcdefg\" provided for data word: %q at index: %d in config",
				"SomeEvent", "DW", 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := evm.NewChainReaderService(testutils.Context(t), logger.Nop(), nil, nil, nil, types.ChainReaderConfig{Contracts: tt.chainContractReaders})
			require.Error(t, err)
			if err != nil {
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			}
		})
	}
}

func TestChainReader_HealthReport(t *testing.T) {
	lp := lpMocks.NewLogPoller(t)
	lp.EXPECT().HealthReport().Return(map[string]error{"lp_name": clcommontypes.ErrFinalityViolated}).Once()
	ht := headstest.NewTracker[*clevmtypes.Head, common.Hash](t)
	htError := errors.New("head tracker error")
	ht.EXPECT().HealthReport().Return(map[string]error{"ht_name": htError}).Once()
	cr, err := evm.NewChainReaderService(testutils.Context(t), logger.Nop(), lp, ht, nil, types.ChainReaderConfig{Contracts: nil})
	require.NoError(t, err)
	healthReport := cr.HealthReport()
	require.True(t, services.ContainsError(healthReport, clcommontypes.ErrFinalityViolated), "expected chain reader to propagate logpoller's error")
	require.True(t, services.ContainsError(healthReport, htError), "expected chain reader to propagate headtracker's error")
}

func TestChainComponents(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-401")

	t.Parallel()
	// shared helper so all tests can run efficiently in parallel
	helper := &helper{}
	helper.Init(t)
	deployLock := sync.Mutex{}
	// add new subtests here so that it can be run on real chains too
	t.Run("RunChainComponentsEvmTests", func(t *testing.T) {
		t.Parallel()
		it := &EVMChainComponentsInterfaceTester[*testing.T]{Helper: helper, DeployLock: &deployLock}
		// These tests are broken in develop as well, so disable them for now
		it.DisableTests([]string{
			interfacetests.ContractReaderQueryKeysReturnsDataTwoEventTypes,
			interfacetests.ContractReaderQueryKeysReturnsDataAsValuesDotValue,
			interfacetests.ContractReaderQueryKeysCanFilterWithValueComparator,
		})
		it.Setup(t)

		RunChainComponentsEvmTests(t, it)
	})

	t.Run("RunChainComponentsInLoopEvmTests", func(t *testing.T) {
		t.Parallel()
		it := &EVMChainComponentsInterfaceTester[*testing.T]{Helper: helper, DeployLock: &deployLock}
		wrapped := commontestutils.WrapContractReaderTesterForLoop(it)
		// These tests are broken in develop as well, so disable them for now
		wrapped.DisableTests([]string{
			interfacetests.ContractReaderQueryKeysReturnsDataTwoEventTypes,
			interfacetests.ContractReaderQueryKeysReturnsDataAsValuesDotValue,
			interfacetests.ContractReaderQueryKeysCanFilterWithValueComparator,
		})
		wrapped.Setup(t)

		RunChainComponentsInLoopEvmTests[*testing.T](t, wrapped, true)
	})

	t.Run("RunChainComponentsInLoopEvmTestsWithBindings", func(t *testing.T) {
		t.Parallel()
		it := &EVMChainComponentsInterfaceTester[*testing.T]{Helper: helper, DeployLock: &deployLock}
		wrapped := WrapContractReaderTesterWithBindings(t, it)
		// TODO, generated binding tests are broken
		wrapped.DisableTests([]string{interfacetests.ContractReaderGetLatestValue})
		wrapped.Setup(t)
		// generated tests are not compatible with parallel running atm
		RunChainComponentsInLoopEvmTests(t, wrapped, false)
	})
}

type helper struct {
	sim         *simulated.Backend
	accounts    []*bind.TransactOpts
	deployerKey *ecdsa.PrivateKey
	senderKey   *ecdsa.PrivateKey
	txm         evmtxmgr.TxManager
	client      client.Client
	db          *sqlx.DB
	lp          logpoller.LogPoller
	ht          logpoller.HeadTracker
}

func getLPOpts() logpoller.Opts {
	return logpoller.Opts{
		PollPeriod:               time.Millisecond,
		FinalityDepth:            finalityDepth,
		BackfillBatchSize:        1,
		RPCBatchSize:             1,
		KeepFinalizedBlocksDepth: 10000,
	}
}

func (h *helper) LogPoller(t *testing.T) logpoller.LogPoller {
	if h.lp != nil {
		return h.lp
	}
	ctx := testutils.Context(t)
	lggr := logger.Nop()
	db := h.Database()

	h.lp = logpoller.NewLogPoller(logpoller.NewORM(h.ChainID(), db, lggr), h.Client(t), lggr, h.HeadTracker(t), getLPOpts())
	require.NoError(t, h.lp.Start(ctx))
	return h.lp
}

func (h *helper) HeadTracker(t *testing.T) logpoller.HeadTracker {
	if h.ht != nil {
		return h.ht
	}
	lpOpts := getLPOpts()
	h.ht = headstest.NewSimulatedHeadTracker(h.Client(t), lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	return h.ht
}

func (h *helper) Init(t *testing.T) {
	h.SetupKeys(t)

	h.accounts = h.Accounts(t)

	h.db = pgtest.NewSqlxDB(t)

	h.Backend()
	h.client = h.Client(t)
	h.LogPoller(t)

	h.txm = h.TXM(t, h.client)
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
		h.sim = simulated.NewBackend(
			evmtypes.GenesisAlloc{h.accounts[0].From: {Balance: big.NewInt(math.MaxInt64)}, h.accounts[1].From: {Balance: big.NewInt(math.MaxInt64)}}, simulated.WithBlockGasLimit(commonGasLimitOnEvms*5000))
		cltest.Mine(h.sim, 1*time.Second)
	}

	return h.sim.Client()
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

func (h *helper) Database() *sqlx.DB {
	return h.db
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

func (h *helper) WrappedChainWriter(cw clcommontypes.ContractWriter, client client.Client) clcommontypes.ContractWriter {
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

	keyStore.XXXTestingOnlyAdd(h.Context(t), ethkey.FromPrivateKey(h.deployerKey))
	require.NoError(t, keyStore.Add(h.Context(t), h.accounts[0].From, h.ChainID()))
	require.NoError(t, keyStore.Enable(h.Context(t), h.accounts[0].From, h.ChainID()))

	keyStore.XXXTestingOnlyAdd(h.Context(t), ethkey.FromPrivateKey(h.senderKey))
	require.NoError(t, keyStore.Add(h.Context(t), h.accounts[1].From, h.ChainID()))
	require.NoError(t, keyStore.Enable(h.Context(t), h.accounts[1].From, h.ChainID()))

	chain, err := app.GetRelayers().LegacyEVMChains().Get((h.ChainID()).String())
	require.NoError(t, err)

	h.txm = chain.TxManager()
	return h.txm
}

func ptr[T any](v T) *T { return &v }
