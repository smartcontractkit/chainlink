package legacygasstation_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	geth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/test-go/testify/mock"
	"github.com/test-go/testify/require"

	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/smartcontractkit/capital-markets-projects/core/gethwrappers/ccip/generated/evm_2_evm_off_ramp"
	mock_contracts "github.com/smartcontractkit/capital-markets-projects/core/gethwrappers/ccip/mocks"
	forwarder_wrapper "github.com/smartcontractkit/capital-markets-projects/core/gethwrappers/legacygasstation/generated/legacy_gas_station_forwarder"
	forwarder_mocks "github.com/smartcontractkit/capital-markets-projects/core/gethwrappers/legacygasstation/mocks"
	"github.com/smartcontractkit/capital-markets-projects/lib/services/legacygasstation"
	lgsmocks "github.com/smartcontractkit/capital-markets-projects/lib/services/legacygasstation/mocks"
	"github.com/smartcontractkit/capital-markets-projects/lib/services/legacygasstation/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	lgsservice "github.com/smartcontractkit/chainlink/v2/core/services/legacygasstation"
)

type request struct {
	tx        legacygasstation.TestLegacyGaslessTx
	confirmed bool
	failed    bool
}

type testcase struct {
	name                 string
	latestBlock          int64
	lookbackBlock        int64
	chainID              uint64
	requestData          []request
	forwardSucceededLogs []*forwarder_wrapper.LegacyGasStationForwarderForwardSucceeded
	offrampExecutionLogs []*evm_2_evm_off_ramp.EVM2EVMOffRampExecutionStateChanged
	resultData           []legacygasstation.TestLegacyGaslessTx
}

type testStatusUpdater struct {
	statusCounter map[string]int
}

func newTestStatusUpdater() *testStatusUpdater {
	return &testStatusUpdater{
		statusCounter: make(map[string]int),
	}
}

func (s *testStatusUpdater) Update(tx types.LegacyGaslessTx) error {
	s.statusCounter[tx.Status.String()]++
	return nil
}

var (
	tests = []testcase{
		{
			name:          "submitted transaction confirmed",
			latestBlock:   100,
			lookbackBlock: 50,
			chainID:       testutils.SimulatedChainID.Uint64(),
			requestData: []request{
				{
					tx: legacygasstation.TestLegacyGaslessTx{
						ID:                 "4877f0a6-4b05-49d9-8776-4c50c24bed03",
						Nonce:              big.NewInt(0),
						Amount:             big.NewInt(1e18),
						SourceChainID:      testutils.SimulatedChainID.Uint64(),
						DestinationChainID: testutils.SimulatedChainID.Uint64(),
					},
					confirmed: true,
				},
			},
			resultData: []legacygasstation.TestLegacyGaslessTx{
				{
					ID:                 "4877f0a6-4b05-49d9-8776-4c50c24bed03",
					Nonce:              big.NewInt(0),
					Amount:             big.NewInt(1e18),
					SourceChainID:      testutils.SimulatedChainID.Uint64(),
					DestinationChainID: testutils.SimulatedChainID.Uint64(),
					Status:             types.Confirmed,
				},
			},
		},
		{
			name:          "submitted transaction failed",
			latestBlock:   100,
			lookbackBlock: 50,
			chainID:       testutils.SimulatedChainID.Uint64(),
			requestData: []request{
				{
					tx: legacygasstation.TestLegacyGaslessTx{
						ID:                 "4877f0a6-4b05-49d9-8776-4c50c24bed03",
						Nonce:              big.NewInt(0),
						Amount:             big.NewInt(1e18),
						SourceChainID:      testutils.SimulatedChainID.Uint64(),
						DestinationChainID: testutils.SimulatedChainID.Uint64(),
					},
					failed: true,
				},
			},
			resultData: []legacygasstation.TestLegacyGaslessTx{
				{
					ID:                 "4877f0a6-4b05-49d9-8776-4c50c24bed03",
					Nonce:              big.NewInt(0),
					Amount:             big.NewInt(1e18),
					SourceChainID:      testutils.SimulatedChainID.Uint64(),
					DestinationChainID: testutils.SimulatedChainID.Uint64(),
					Status:             types.Failure,
				},
			},
		},
		{
			name:          "confirmed transaction finalized",
			latestBlock:   100,
			lookbackBlock: 50,
			chainID:       testutils.SimulatedChainID.Uint64(),
			requestData: []request{
				{
					tx: legacygasstation.TestLegacyGaslessTx{
						ID:                 "4877f0a6-4b05-49d9-8776-4c50c24bed03",
						Nonce:              big.NewInt(0),
						Amount:             big.NewInt(1e18),
						SourceChainID:      testutils.SimulatedChainID.Uint64(),
						DestinationChainID: testutils.SimulatedChainID.Uint64(),
						Status:             types.Confirmed,
					},
				},
			},
			forwardSucceededLogs: []*forwarder_wrapper.LegacyGasStationForwarderForwardSucceeded{
				{
					From:  legacygasstation.FromAddress,
					Nonce: big.NewInt(0),
					Raw: geth_types.Log{
						Address: legacygasstation.ForwarderAddress,
					},
				},
			},
			resultData: []legacygasstation.TestLegacyGaslessTx{
				{
					ID:                 "4877f0a6-4b05-49d9-8776-4c50c24bed03",
					Nonce:              big.NewInt(0),
					Amount:             big.NewInt(1e18),
					SourceChainID:      testutils.SimulatedChainID.Uint64(),
					DestinationChainID: testutils.SimulatedChainID.Uint64(),
					Status:             types.Finalized,
				},
			},
		},
		{
			name:          "confirmed transaction failed",
			latestBlock:   100,
			lookbackBlock: 50,
			chainID:       testutils.SimulatedChainID.Uint64(),
			requestData: []request{
				{
					tx: legacygasstation.TestLegacyGaslessTx{
						ID:                 "4877f0a6-4b05-49d9-8776-4c50c24bed03",
						Nonce:              big.NewInt(0),
						Amount:             big.NewInt(1e18),
						SourceChainID:      testutils.SimulatedChainID.Uint64(),
						DestinationChainID: testutils.SimulatedChainID.Uint64(),
						Status:             types.Confirmed,
					},
					failed: true,
				},
			},
			resultData: []legacygasstation.TestLegacyGaslessTx{
				{
					ID:                 "4877f0a6-4b05-49d9-8776-4c50c24bed03",
					Nonce:              big.NewInt(0),
					Amount:             big.NewInt(1e18),
					SourceChainID:      testutils.SimulatedChainID.Uint64(),
					DestinationChainID: testutils.SimulatedChainID.Uint64(),
					Status:             types.Failure,
				},
			},
		},
		{
			name:          "multiple submitted txs finalized",
			latestBlock:   100,
			lookbackBlock: 50,
			chainID:       testutils.SimulatedChainID.Uint64(),
			requestData: []request{
				{
					tx: legacygasstation.TestLegacyGaslessTx{
						ID:                 "4877f0a6-4b05-49d9-8776-4c50c24bed03",
						Nonce:              big.NewInt(0),
						Amount:             big.NewInt(1e18),
						SourceChainID:      testutils.SimulatedChainID.Uint64(),
						DestinationChainID: testutils.SimulatedChainID.Uint64(),
						Status:             types.Confirmed,
					},
				},
				{
					tx: legacygasstation.TestLegacyGaslessTx{
						ID:                 "a4efbb8b-ac67-46fb-8ded-c883f7f5fcab",
						From:               common.HexToAddress("0x780b3102c62d5DfDCc658B3480B93041Ba46F499"),
						Nonce:              big.NewInt(0),
						Amount:             big.NewInt(1e18),
						SourceChainID:      testutils.SimulatedChainID.Uint64(),
						DestinationChainID: testutils.SimulatedChainID.Uint64(),
						Status:             types.Confirmed,
					},
				},
			},
			forwardSucceededLogs: []*forwarder_wrapper.LegacyGasStationForwarderForwardSucceeded{
				{
					From:  legacygasstation.FromAddress,
					Nonce: big.NewInt(0),
					Raw: geth_types.Log{
						Address: legacygasstation.ForwarderAddress,
					},
				},
				{
					From:  common.HexToAddress("0x780b3102c62d5DfDCc658B3480B93041Ba46F499"),
					Nonce: big.NewInt(0),
					Raw: geth_types.Log{
						Address: legacygasstation.ForwarderAddress,
					},
				},
			},
			resultData: []legacygasstation.TestLegacyGaslessTx{
				{
					ID:                 "4877f0a6-4b05-49d9-8776-4c50c24bed03",
					Nonce:              big.NewInt(0),
					Amount:             big.NewInt(1e18),
					SourceChainID:      testutils.SimulatedChainID.Uint64(),
					DestinationChainID: testutils.SimulatedChainID.Uint64(),
					Status:             types.Finalized,
				},
				{
					ID:                 "a4efbb8b-ac67-46fb-8ded-c883f7f5fcab",
					From:               common.HexToAddress("0x780b3102c62d5DfDCc658B3480B93041Ba46F499"),
					Nonce:              big.NewInt(0),
					Amount:             big.NewInt(1e18),
					SourceChainID:      testutils.SimulatedChainID.Uint64(),
					DestinationChainID: testutils.SimulatedChainID.Uint64(),
					Status:             types.Finalized,
				},
			},
		},
		{
			name:          "no forwarder logs",
			latestBlock:   100,
			lookbackBlock: 50,
			chainID:       testutils.SimulatedChainID.Uint64(),
			requestData: []request{
				{
					tx: legacygasstation.TestLegacyGaslessTx{
						ID:                 "4877f0a6-4b05-49d9-8776-4c50c24bed03",
						Nonce:              big.NewInt(0),
						Amount:             big.NewInt(1e18),
						SourceChainID:      testutils.SimulatedChainID.Uint64(),
						DestinationChainID: testutils.SimulatedChainID.Uint64(),
					},
				},
				{
					tx: legacygasstation.TestLegacyGaslessTx{
						ID:                 "a4efbb8b-ac67-46fb-8ded-c883f7f5fcab",
						From:               common.HexToAddress("0x780b3102c62d5DfDCc658B3480B93041Ba46F499"),
						Nonce:              big.NewInt(0),
						Amount:             big.NewInt(1e18),
						SourceChainID:      testutils.SimulatedChainID.Uint64(),
						DestinationChainID: testutils.SimulatedChainID.Uint64(),
					},
				},
			},
			forwardSucceededLogs: []*forwarder_wrapper.LegacyGasStationForwarderForwardSucceeded{},
			resultData: []legacygasstation.TestLegacyGaslessTx{
				{
					ID:                 "4877f0a6-4b05-49d9-8776-4c50c24bed03",
					Nonce:              big.NewInt(0),
					Amount:             big.NewInt(1e18),
					SourceChainID:      testutils.SimulatedChainID.Uint64(),
					DestinationChainID: testutils.SimulatedChainID.Uint64(),
				},
				{
					ID:                 "a4efbb8b-ac67-46fb-8ded-c883f7f5fcab",
					From:               common.HexToAddress("0x780b3102c62d5DfDCc658B3480B93041Ba46F499"),
					Nonce:              big.NewInt(0),
					Amount:             big.NewInt(1e18),
					SourceChainID:      testutils.SimulatedChainID.Uint64(),
					DestinationChainID: testutils.SimulatedChainID.Uint64(),
				},
			},
		},
		{
			name:          "cross chain submitted to source finalized log",
			latestBlock:   100,
			lookbackBlock: 50,
			chainID:       testutils.SimulatedChainID.Uint64(),
			requestData: []request{
				{
					tx: legacygasstation.TestLegacyGaslessTx{
						ID:                 "4877f0a6-4b05-49d9-8776-4c50c24bed03",
						Nonce:              big.NewInt(0),
						Amount:             big.NewInt(1e18),
						SourceChainID:      testutils.SimulatedChainID.Uint64(),
						DestinationChainID: 1000,
						Status:             types.Confirmed,
					},
				},
			},
			forwardSucceededLogs: []*forwarder_wrapper.LegacyGasStationForwarderForwardSucceeded{
				{
					From:  legacygasstation.FromAddress,
					Nonce: big.NewInt(0),
					Raw: geth_types.Log{
						Address: legacygasstation.ForwarderAddress,
					},
					ReturnValue: common.HexToHash("0x30").Bytes(),
				},
			},
			resultData: []legacygasstation.TestLegacyGaslessTx{
				{
					ID:                 "4877f0a6-4b05-49d9-8776-4c50c24bed03",
					Nonce:              big.NewInt(0),
					Amount:             big.NewInt(1e18),
					SourceChainID:      testutils.SimulatedChainID.Uint64(),
					DestinationChainID: 1000,
					Status:             types.SourceFinalized,
				},
			},
		},
		{
			name:          "cross chain source finalized to finalized",
			latestBlock:   100,
			lookbackBlock: 50,
			chainID:       1000,
			requestData: []request{
				{
					tx: legacygasstation.TestLegacyGaslessTx{
						ID:                 "4877f0a6-4b05-49d9-8776-4c50c24bed03",
						Nonce:              big.NewInt(0),
						Amount:             big.NewInt(1e18),
						SourceChainID:      1000,
						DestinationChainID: testutils.SimulatedChainID.Uint64(),
						CCIPMessageID:      ptr(common.HexToHash("0x30")),
						Status:             types.SourceFinalized,
					},
				},
			},
			offrampExecutionLogs: []*evm_2_evm_off_ramp.EVM2EVMOffRampExecutionStateChanged{
				{
					MessageId: common.HexToHash("0x30"),
					Raw: geth_types.Log{
						Address: legacygasstation.OfframpAddress,
					},
				},
			},
			resultData: []legacygasstation.TestLegacyGaslessTx{
				{
					ID:                 "4877f0a6-4b05-49d9-8776-4c50c24bed03",
					Nonce:              big.NewInt(0),
					Amount:             big.NewInt(1e18),
					SourceChainID:      1000,
					DestinationChainID: testutils.SimulatedChainID.Uint64(),
					Status:             types.Finalized,
				},
			},
		},
	}
)

func TestSidecar(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sc, orm, su := setUp(t, test)
			err := sc.Run(testutils.Context(t))
			require.NoError(t, err)
			assertAfterSidecarRun(t, test, orm, su)
		})
	}
}

func setUp(t *testing.T, test testcase) (*legacygasstation.Sidecar, legacygasstation.ORM, *testStatusUpdater) {
	cfg, db := heavyweight.FullTestDBV2(t, "legacy_gas_station_sidecar_test", func(c *chainlink.Config, s *chainlink.Secrets) {
		require.Zero(t, testutils.SimulatedChainID.Cmp(c.EVM[0].ChainID.ToInt()))
		c.Feature.LogPoller = ptr(true)
	})
	backend := cltest.NewSimulatedBackend(t, core.GenesisAlloc{}, uint32(ethconfig.Defaults.Miner.GasCeil))
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, cfg, backend)
	forwarder := forwarder_mocks.NewLegacyGasStationForwarderInterface(t)
	lggr := logger.TestLogger(t)
	offramp := mock_contracts.NewEVM2EVMOffRampInterface(t)
	orm := lgsservice.NewORM(db, lggr, cfg.Database())
	chain, err := app.GetRelayers().LegacyEVMChains().Get(testutils.SimulatedChainID.String())
	require.NoError(t, err)
	lp := lgsmocks.NewLogPoller(t)
	lp.On("FilterName", mock.Anything, mock.Anything, mock.Anything).Return("filterName")
	lp.On("RegisterFilter", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	lp.On("LatestBlock", mock.Anything).Return(test.latestBlock, nil)
	forwarder.On("Address").Return(legacygasstation.ForwarderAddress)
	offramp.On("Address").Return(legacygasstation.OfframpAddress)
	var (
		fsLpLogs  []gethtypes.Log
		oelLpLogs []gethtypes.Log
	)

	for _, fl := range test.forwardSucceededLogs {
		forwarder.On("ParseLog", mock.Anything).Return(fl, nil).Once()
		fsLpLogs = append(fsLpLogs, gethtypes.Log{
			Topics: []common.Hash{forwarder_wrapper.LegacyGasStationForwarderForwardSucceeded{}.Topic()},
		})
	}
	for _, oel := range test.offrampExecutionLogs {
		offramp.On("ParseLog", mock.Anything).Return(oel, nil).Once()
		oelLpLogs = append(oelLpLogs, gethtypes.Log{
			Topics: []common.Hash{evm_2_evm_off_ramp.EVM2EVMOffRampExecutionStateChanged{}.Topic()},
		})
	}
	lp.On("IndexedLogsByBlockRange",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		forwarder_wrapper.LegacyGasStationForwarderForwardSucceeded{}.Topic(),
		legacygasstation.ForwarderAddress,
		1,
		mock.Anything,
	).Return(fsLpLogs, nil).Maybe()
	lp.On("IndexedLogsByBlockRange",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		evm_2_evm_off_ramp.EVM2EVMOffRampExecutionStateChanged{}.Topic(),
		legacygasstation.OfframpAddress,
		2,
		mock.Anything,
	).Return(oelLpLogs, nil).Maybe()

	su := newTestStatusUpdater()
	sc, err := legacygasstation.NewSidecar(
		lggr,
		lp,
		forwarder,
		offramp,
		testutils.SimulatedChainID.Uint64(),
		chain.Config().EVM().FinalityDepth(),
		uint32(test.lookbackBlock),
		orm,
		su,
	)
	require.NoError(t, err)
	for i, r := range test.requestData {
		chainID := cltest.FixtureChainID
		blockNumber := int64(75)
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, app.KeyStore.Eth(), chainID)
		txStore := cltest.NewTestTxStore(t, app.GetSqlxDB(), app.Config.Database())
		var ethTx txmgr.Tx
		if r.confirmed {
			ethTx = cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, int64(i), blockNumber, fromAddress)
			blockHash := utils.NewHash()
			receipt := evmtypes.Receipt{
				TxHash:           ethTx.TxAttempts[0].Hash,
				BlockHash:        blockHash,
				BlockNumber:      big.NewInt(int64(i)),
				TransactionIndex: uint(1),
				Status:           uint64(1), // reverted txs have 0 as status. non-zero other wise
			}
			err := app.TxmStorageService().SaveFetchedReceipts([]*evmtypes.Receipt{&receipt}, &chainID)
			require.NoError(t, err)
		} else if r.failed {
			ethTx = cltest.MustInsertFatalErrorEthTx(t, txStore, fromAddress)
		} else {
			ethTx = cltest.MustInsertInProgressEthTxWithAttempt(t, txStore, evmtypes.Nonce(int64(i)), fromAddress)
		}
		r.tx.EthTxID = ethTx.GetID()
		tx := legacygasstation.LegacyGaslessTx(t, r.tx)
		err = orm.InsertLegacyGaslessTx(testutils.Context(t), tx)
		require.NoError(t, err)
		err = orm.UpdateLegacyGaslessTx(testutils.Context(t), tx) // update populates ccipMessageID and failureReason
		require.NoError(t, err)
	}
	return sc, orm, su
}

func assertAfterSidecarRun(t *testing.T, test testcase, orm legacygasstation.ORM, su *testStatusUpdater) {
	confirmedTxs, submittedTxs, finalizedTxs, sourceFinalizedTxs, failedTxs := categorizeTestTxs(t, test.resultData)

	txs, err := orm.SelectBySourceChainIDAndStatus(testutils.Context(t), test.chainID, types.Confirmed)
	require.NoError(t, err)
	require.Equal(t, len(confirmedTxs), len(txs))

	txs, err = orm.SelectBySourceChainIDAndStatus(testutils.Context(t), test.chainID, types.Submitted)
	require.NoError(t, err)
	require.Equal(t, len(submittedTxs), len(txs))

	txs, err = orm.SelectBySourceChainIDAndStatus(testutils.Context(t), test.chainID, types.Finalized)
	require.NoError(t, err)
	require.Equal(t, len(finalizedTxs), len(txs))

	txs, err = orm.SelectBySourceChainIDAndStatus(testutils.Context(t), test.chainID, types.SourceFinalized)
	require.NoError(t, err)
	require.Equal(t, len(sourceFinalizedTxs), len(txs))

	txs, err = orm.SelectBySourceChainIDAndStatus(testutils.Context(t), test.chainID, types.Failure)
	require.NoError(t, err)
	require.Equal(t, len(failedTxs), len(txs))

	expectedStatusUpdates := make(map[string]int)
	for i, req := range test.requestData {
		resultStatus := test.resultData[i].Status.String()
		if req.tx.Status != test.resultData[i].Status {
			expectedStatusUpdates[resultStatus]++
		}
	}
	require.Equal(t, expectedStatusUpdates[types.Confirmed.String()], su.statusCounter[types.Confirmed.String()])
	require.Equal(t, expectedStatusUpdates[types.Failure.String()], su.statusCounter[types.Failure.String()])
	require.Equal(t, expectedStatusUpdates[types.SourceFinalized.String()], su.statusCounter[types.SourceFinalized.String()])
	require.Equal(t, expectedStatusUpdates[types.Finalized.String()], su.statusCounter[types.Finalized.String()])
}

func categorizeTestTxs(t *testing.T, testTxs []legacygasstation.TestLegacyGaslessTx) (
	confirmedTxs,
	submittedTxs,
	finalizedTxs,
	sourceFinalizedTxs,
	failedTxs []types.LegacyGaslessTx,
) {
	for _, testTx := range testTxs {
		tx := legacygasstation.LegacyGaslessTx(t, testTx)
		switch tx.Status {
		case types.Confirmed:
			confirmedTxs = append(confirmedTxs, tx)
		case types.Submitted:
			submittedTxs = append(submittedTxs, tx)
		case types.SourceFinalized:
			sourceFinalizedTxs = append(sourceFinalizedTxs, tx)
		case types.Finalized:
			finalizedTxs = append(finalizedTxs, tx)
		case types.Failure:
			failedTxs = append(failedTxs, tx)
		default:
			t.Errorf("unexpected status: %s", tx.Status)
		}
	}
	return
}
