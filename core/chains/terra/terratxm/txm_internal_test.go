package terratxm

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"

	tmservicetypes "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	tmtypes "github.com/tendermint/tendermint/proto/tendermint/types"
	wasmtypes "github.com/terra-money/core/x/wasm/types"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra"
	terraclient "github.com/smartcontractkit/chainlink-terra/pkg/terra/client"
	tcmocks "github.com/smartcontractkit/chainlink-terra/pkg/terra/client/mocks"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/terratest"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	. "github.com/smartcontractkit/chainlink-terra/pkg/terra/db"
)

func generateExecuteMsg(t *testing.T, msg []byte, from, to cosmostypes.AccAddress) cosmostypes.Msg {
	return wasmtypes.NewMsgExecuteContract(from, to, msg, cosmostypes.Coins{})
}

func TestTxm(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	lggr := testutils.LoggerAssertMaxLevel(t, zapcore.ErrorLevel)
	ks := keystore.New(db, utils.FastScryptParams, lggr, pgtest.NewPGCfg(true))
	require.NoError(t, ks.Unlock("blah"))
	k1, err := ks.Terra().Create()
	require.NoError(t, err)
	sender1, err := cosmostypes.AccAddressFromBech32(k1.PublicKeyStr())
	require.NoError(t, err)
	k2, err := ks.Terra().Create()
	require.NoError(t, err)
	sender2, err := cosmostypes.AccAddressFromBech32(k2.PublicKeyStr())
	require.NoError(t, err)
	contract, err := cosmostypes.AccAddressFromBech32("terra1pp76d50yv2ldaahsdxdv8mmzqfjr2ax97gmue8")
	require.NoError(t, err)
	contract2, err := cosmostypes.AccAddressFromBech32("terra1mx72uukvzqtzhc6gde7shrjqfu5srk22v7gmww")
	require.NoError(t, err)
	logCfg := pgtest.NewPGCfg(true)
	chainID := fmt.Sprintf("Chainlinktest-%d", rand.Int31n(999999))
	terratest.MustInsertChain(t, db, &Chain{ID: chainID})
	require.NoError(t, err)
	cfg := terra.NewConfig(ChainCfg{
		MaxMsgsPerBatch: null.IntFrom(2),
	}, lggr)
	gpe := terraclient.NewMustGasPriceEstimator([]terraclient.GasPricesEstimator{
		terraclient.NewFixedGasPriceEstimator(map[string]cosmostypes.DecCoin{
			"uluna": cosmostypes.NewDecCoinFromDec("uluna", cosmostypes.MustNewDecFromStr("0.01")),
		}),
	}, lggr)

	t.Run("single msg", func(t *testing.T) {
		tc := new(tcmocks.ReaderWriter)
		tcFn := func() (terraclient.ReaderWriter, error) { return tc, nil }
		txm := NewTxm(db, tcFn, *gpe, chainID, cfg, ks.Terra(), lggr, logCfg, nil)

		// Enqueue a single msg, then send it in a batch
		id1, err := txm.Enqueue(contract.String(), generateExecuteMsg(t, []byte(`1`), sender1, contract))
		require.NoError(t, err)
		tc.On("Account", mock.Anything).Return(uint64(0), uint64(0), nil)
		tc.On("BatchSimulateUnsigned", mock.Anything, mock.Anything).Return(&terraclient.BatchSimResults{
			Failed: nil,
			Succeeded: terraclient.SimMsgs{{ID: id1, Msg: &wasmtypes.MsgExecuteContract{
				Sender:     sender1.String(),
				ExecuteMsg: []byte(`1`),
			}}},
		}, nil)
		tc.On("SimulateUnsigned", mock.Anything, mock.Anything).Return(&txtypes.SimulateResponse{GasInfo: &cosmostypes.GasInfo{
			GasUsed: 1_000_000,
		}}, nil)
		tc.On("LatestBlock").Return(&tmservicetypes.GetLatestBlockResponse{Block: &tmtypes.Block{
			Header: tmtypes.Header{Height: 1},
		}}, nil)
		tc.On("CreateAndSign", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]byte{0x01}, nil)

		txResp := &cosmostypes.TxResponse{TxHash: "4BF5122F344554C53BDE2EBB8CD2B7E3D1600AD631C385A5D7CCE23C7785459A"}
		tc.On("Broadcast", mock.Anything, mock.Anything).Return(&txtypes.BroadcastTxResponse{TxResponse: txResp}, nil)
		tc.On("Tx", mock.Anything).Return(&txtypes.GetTxResponse{Tx: &txtypes.Tx{}, TxResponse: txResp}, nil)
		txm.sendMsgBatch(testutils.Context(t))

		// Should be in completed state
		completed, err := txm.orm.GetMsgs(id1)
		require.NoError(t, err)
		require.Equal(t, 1, len(completed))
		assert.Equal(t, completed[0].State, Confirmed)
		tc.AssertExpectations(t)
	})

	t.Run("two msgs different accounts", func(t *testing.T) {
		tc := new(tcmocks.ReaderWriter)
		tcFn := func() (terraclient.ReaderWriter, error) { return tc, nil }
		txm := NewTxm(db, tcFn, *gpe, chainID, cfg, ks.Terra(), lggr, pgtest.NewPGCfg(true), nil)

		id1, err := txm.Enqueue(contract.String(), generateExecuteMsg(t, []byte(`0`), sender1, contract))
		require.NoError(t, err)
		id2, err := txm.Enqueue(contract.String(), generateExecuteMsg(t, []byte(`1`), sender2, contract))
		require.NoError(t, err)

		tc.On("Account", mock.Anything).Return(uint64(0), uint64(0), nil).Once()
		// Note this must be arg dependent, we don't know which order
		// the procesing will happen in (map iteration by from address).
		tc.On("BatchSimulateUnsigned", terraclient.SimMsgs{
			{
				ID: id2,
				Msg: &wasmtypes.MsgExecuteContract{
					Sender:     sender2.String(),
					ExecuteMsg: []byte(`1`),
					Contract:   contract.String(),
				},
			},
		}, mock.Anything).Return(&terraclient.BatchSimResults{
			Failed: nil,
			Succeeded: terraclient.SimMsgs{
				{
					ID: id2,
					Msg: &wasmtypes.MsgExecuteContract{
						Sender:     sender2.String(),
						ExecuteMsg: []byte(`1`),
						Contract:   contract.String(),
					},
				},
			},
		}, nil).Once()
		tc.On("SimulateUnsigned", mock.Anything, mock.Anything).Return(&txtypes.SimulateResponse{GasInfo: &cosmostypes.GasInfo{
			GasUsed: 1_000_000,
		}}, nil).Once()
		tc.On("LatestBlock").Return(&tmservicetypes.GetLatestBlockResponse{Block: &tmtypes.Block{
			Header: tmtypes.Header{Height: 1},
		}}, nil).Once()
		tc.On("CreateAndSign", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]byte{0x01}, nil).Once()
		txResp := &cosmostypes.TxResponse{TxHash: "4BF5122F344554C53BDE2EBB8CD2B7E3D1600AD631C385A5D7CCE23C7785459A"}
		tc.On("Broadcast", mock.Anything, mock.Anything).Return(&txtypes.BroadcastTxResponse{TxResponse: txResp}, nil).Once()
		tc.On("Tx", mock.Anything).Return(&txtypes.GetTxResponse{Tx: &txtypes.Tx{}, TxResponse: txResp}, nil).Once()
		txm.sendMsgBatch(testutils.Context(t))

		// Should be in completed state
		completed, err := txm.orm.GetMsgs(id1, id2)
		require.NoError(t, err)
		require.Equal(t, 2, len(completed))
		assert.Equal(t, Errored, completed[0].State) // cancelled
		assert.Equal(t, Confirmed, completed[1].State)
		tc.AssertExpectations(t)
	})

	t.Run("two msgs different contracts", func(t *testing.T) {
		tc := new(tcmocks.ReaderWriter)
		tcFn := func() (terraclient.ReaderWriter, error) { return tc, nil }
		txm := NewTxm(db, tcFn, *gpe, chainID, cfg, ks.Terra(), lggr, pgtest.NewPGCfg(true), nil)

		id1, err := txm.Enqueue(contract.String(), generateExecuteMsg(t, []byte(`0`), sender1, contract))
		require.NoError(t, err)
		id2, err := txm.Enqueue(contract2.String(), generateExecuteMsg(t, []byte(`1`), sender2, contract2))
		require.NoError(t, err)
		ids := []int64{id1, id2}
		senders := []string{sender1.String(), sender2.String()}
		contracts := []string{contract.String(), contract2.String()}
		for i := 0; i < 2; i++ {
			tc.On("Account", mock.Anything).Return(uint64(0), uint64(0), nil).Once()
			// Note this must be arg dependent, we don't know which order
			// the procesing will happen in (map iteration by from address).
			tc.On("BatchSimulateUnsigned", terraclient.SimMsgs{
				{
					ID: ids[i],
					Msg: &wasmtypes.MsgExecuteContract{
						Sender:     senders[i],
						ExecuteMsg: []byte(fmt.Sprintf(`%d`, i)),
						Contract:   contracts[i],
					},
				},
			}, mock.Anything).Return(&terraclient.BatchSimResults{
				Failed: nil,
				Succeeded: terraclient.SimMsgs{
					{
						ID: ids[i],
						Msg: &wasmtypes.MsgExecuteContract{
							Sender:     senders[i],
							ExecuteMsg: []byte(fmt.Sprintf(`%d`, i)),
							Contract:   contracts[i],
						},
					},
				},
			}, nil).Once()
			tc.On("SimulateUnsigned", mock.Anything, mock.Anything).Return(&txtypes.SimulateResponse{GasInfo: &cosmostypes.GasInfo{
				GasUsed: 1_000_000,
			}}, nil).Once()
			tc.On("LatestBlock").Return(&tmservicetypes.GetLatestBlockResponse{Block: &tmtypes.Block{
				Header: tmtypes.Header{Height: 1},
			}}, nil).Once()
			tc.On("CreateAndSign", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]byte{0x01}, nil).Once()
		}
		txResp := &cosmostypes.TxResponse{TxHash: "4BF5122F344554C53BDE2EBB8CD2B7E3D1600AD631C385A5D7CCE23C7785459A"}
		tc.On("Broadcast", mock.Anything, mock.Anything).Return(&txtypes.BroadcastTxResponse{TxResponse: txResp}, nil).Twice()
		tc.On("Tx", mock.Anything).Return(&txtypes.GetTxResponse{Tx: &txtypes.Tx{}, TxResponse: txResp}, nil).Twice()
		txm.sendMsgBatch(testutils.Context(t))

		// Should be in completed state
		completed, err := txm.orm.GetMsgs(id1, id2)
		require.NoError(t, err)
		require.Equal(t, 2, len(completed))
		assert.Equal(t, Confirmed, completed[0].State)
		assert.Equal(t, Confirmed, completed[1].State)
		tc.AssertExpectations(t)
	})

	t.Run("failed to confirm", func(t *testing.T) {
		tc := new(tcmocks.ReaderWriter)
		tc.On("Tx", mock.Anything).Return(&txtypes.GetTxResponse{
			Tx:         &txtypes.Tx{},
			TxResponse: &cosmostypes.TxResponse{TxHash: "0x123"},
		}, errors.New("not found")).Twice()
		cfg := terra.NewConfig(ChainCfg{}, lggr)
		tcFn := func() (terraclient.ReaderWriter, error) { return tc, nil }
		txm := NewTxm(db, tcFn, *gpe, chainID, cfg, ks.Terra(), lggr, pgtest.NewPGCfg(true), nil)
		i, err := txm.orm.InsertMsg("blah", "", []byte{0x01})
		require.NoError(t, err)
		txh := "0x123"
		require.NoError(t, txm.orm.UpdateMsgs([]int64{i}, Started, &txh))
		require.NoError(t, txm.orm.UpdateMsgs([]int64{i}, Broadcasted, &txh))
		err = txm.confirmTx(testutils.Context(t), tc, txh, []int64{i}, 2, 1*time.Millisecond)
		require.NoError(t, err)
		m, err := txm.orm.GetMsgs(i)
		require.NoError(t, err)
		require.Equal(t, 1, len(m))
		assert.Equal(t, Errored, m[0].State)
		tc.AssertExpectations(t)
	})

	t.Run("confirm any unconfirmed", func(t *testing.T) {
		require.Equal(t, int64(2), cfg.MaxMsgsPerBatch())
		txHash1 := "0x1234"
		txHash2 := "0x1235"
		txHash3 := "0xabcd"
		tc := new(tcmocks.ReaderWriter)
		tc.On("Tx", txHash1).Return(&txtypes.GetTxResponse{
			TxResponse: &cosmostypes.TxResponse{TxHash: txHash1},
		}, nil).Once()
		tc.On("Tx", txHash2).Return(&txtypes.GetTxResponse{
			TxResponse: &cosmostypes.TxResponse{TxHash: txHash2},
		}, nil).Once()
		tc.On("Tx", txHash3).Return(&txtypes.GetTxResponse{
			TxResponse: &cosmostypes.TxResponse{TxHash: txHash3},
		}, nil).Once()
		tcFn := func() (terraclient.ReaderWriter, error) { return tc, nil }
		txm := NewTxm(db, tcFn, *gpe, chainID, cfg, ks.Terra(), lggr, pgtest.NewPGCfg(true), nil)

		// Insert and broadcast 3 msgs with different txhashes.
		id1, err := txm.orm.InsertMsg("blah", "", []byte{0x01})
		require.NoError(t, err)
		id2, err := txm.orm.InsertMsg("blah", "", []byte{0x02})
		require.NoError(t, err)
		id3, err := txm.orm.InsertMsg("blah", "", []byte{0x03})
		require.NoError(t, err)
		err = txm.orm.UpdateMsgs([]int64{id1}, Started, &txHash1)
		require.NoError(t, err)
		err = txm.orm.UpdateMsgs([]int64{id2}, Started, &txHash2)
		require.NoError(t, err)
		err = txm.orm.UpdateMsgs([]int64{id3}, Started, &txHash3)
		require.NoError(t, err)
		err = txm.orm.UpdateMsgs([]int64{id1}, Broadcasted, &txHash1)
		require.NoError(t, err)
		err = txm.orm.UpdateMsgs([]int64{id2}, Broadcasted, &txHash2)
		require.NoError(t, err)
		err = txm.orm.UpdateMsgs([]int64{id3}, Broadcasted, &txHash3)
		require.NoError(t, err)

		// Confirm them as in a restart while confirming scenario
		txm.confirmAnyUnconfirmed(testutils.Context(t))
		msgs, err := txm.orm.GetMsgs(id1, id2, id3)
		require.NoError(t, err)
		require.Equal(t, 3, len(msgs))
		assert.Equal(t, Confirmed, msgs[0].State)
		assert.Equal(t, Confirmed, msgs[1].State)
		assert.Equal(t, Confirmed, msgs[2].State)
		tc.AssertExpectations(t)
	})

	t.Run("expired msgs", func(t *testing.T) {
		tc := new(tcmocks.ReaderWriter)
		timeout := models.MustMakeDuration(1 * time.Millisecond)
		tcFn := func() (terraclient.ReaderWriter, error) { return tc, nil }
		cfgShortExpiry := terra.NewConfig(ChainCfg{
			MaxMsgsPerBatch: null.IntFrom(2),
			TxMsgTimeout:    &timeout,
		}, lggr)
		txm := NewTxm(db, tcFn, *gpe, chainID, cfgShortExpiry, ks.Terra(), lggr, pgtest.NewPGCfg(true), nil)

		// Send a single one expired
		id1, err := txm.orm.InsertMsg("blah", "", []byte{0x03})
		require.NoError(t, err)
		time.Sleep(1 * time.Millisecond)
		txm.sendMsgBatch(context.Background())
		// Should be marked errored
		m, err := txm.orm.GetMsgs(id1)
		require.NoError(t, err)
		assert.Equal(t, Errored, m[0].State)

		// Send a batch which is all expired
		id2, err := txm.orm.InsertMsg("blah", "", []byte{0x03})
		require.NoError(t, err)
		id3, err := txm.orm.InsertMsg("blah", "", []byte{0x03})
		require.NoError(t, err)
		time.Sleep(1 * time.Millisecond)
		txm.sendMsgBatch(context.Background())
		require.NoError(t, err)
		ms, err := txm.orm.GetMsgs(id2, id3)
		assert.Equal(t, Errored, ms[0].State)
		assert.Equal(t, Errored, ms[1].State)
	})

	t.Run("started msgs", func(t *testing.T) {
		tc := new(tcmocks.ReaderWriter)
		tc.On("Account", mock.Anything).Return(uint64(0), uint64(0), nil)
		tc.On("SimulateUnsigned", mock.Anything, mock.Anything).Return(&txtypes.SimulateResponse{GasInfo: &cosmostypes.GasInfo{
			GasUsed: 1_000_000,
		}}, nil)
		tc.On("LatestBlock").Return(&tmservicetypes.GetLatestBlockResponse{Block: &tmtypes.Block{
			Header: tmtypes.Header{Height: 1},
		}}, nil)
		tc.On("CreateAndSign", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]byte{0x01}, nil)
		txResp := &cosmostypes.TxResponse{TxHash: "4BF5122F344554C53BDE2EBB8CD2B7E3D1600AD631C385A5D7CCE23C7785459A"}
		tc.On("Broadcast", mock.Anything, mock.Anything).Return(&txtypes.BroadcastTxResponse{TxResponse: txResp}, nil)
		tc.On("Tx", mock.Anything).Return(&txtypes.GetTxResponse{Tx: &txtypes.Tx{}, TxResponse: txResp}, nil)
		tcFn := func() (terraclient.ReaderWriter, error) { return tc, nil }
		cfg := terra.NewConfig(ChainCfg{
			MaxMsgsPerBatch: null.IntFrom(2),
		}, lggr)
		txm := NewTxm(db, tcFn, *gpe, chainID, cfg, ks.Terra(), lggr, pgtest.NewPGCfg(true), nil)

		// Leftover started is processed
		msg1 := generateExecuteMsg(t, []byte{0x03}, sender1, contract)
		id1 := mustInsertMsg(t, txm, contract.String(), msg1)
		require.NoError(t, txm.orm.UpdateMsgs([]int64{id1}, Started, nil))
		msgs := terraclient.SimMsgs{{ID: id1, Msg: &wasmtypes.MsgExecuteContract{
			Sender:     sender1.String(),
			ExecuteMsg: []byte{0x03},
			Contract:   contract.String(),
		}}}
		tc.On("BatchSimulateUnsigned", msgs, mock.Anything).
			Return(&terraclient.BatchSimResults{Failed: nil, Succeeded: msgs}, nil).Once()
		time.Sleep(1 * time.Millisecond)
		txm.sendMsgBatch(context.Background())
		m, err := txm.orm.GetMsgs(id1)
		require.NoError(t, err)
		assert.Equal(t, Confirmed, m[0].State)

		// Leftover started is not cancelled
		msg2 := generateExecuteMsg(t, []byte{0x04}, sender1, contract)
		msg3 := generateExecuteMsg(t, []byte{0x05}, sender1, contract)
		id2 := mustInsertMsg(t, txm, contract.String(), msg2)
		require.NoError(t, txm.orm.UpdateMsgs([]int64{id2}, Started, nil))
		time.Sleep(time.Millisecond) // ensure != CreatedAt
		id3 := mustInsertMsg(t, txm, contract.String(), msg3)
		msgs = terraclient.SimMsgs{{ID: id2, Msg: &wasmtypes.MsgExecuteContract{
			Sender:     sender1.String(),
			ExecuteMsg: []byte{0x04},
			Contract:   contract.String(),
		}}, {ID: id3, Msg: &wasmtypes.MsgExecuteContract{
			Sender:     sender1.String(),
			ExecuteMsg: []byte{0x05},
			Contract:   contract.String(),
		}}}
		tc.On("BatchSimulateUnsigned", msgs, mock.Anything).
			Return(&terraclient.BatchSimResults{Failed: nil, Succeeded: msgs}, nil).Once()
		time.Sleep(1 * time.Millisecond)
		txm.sendMsgBatch(context.Background())
		require.NoError(t, err)
		ms, err := txm.orm.GetMsgs(id2, id3)
		assert.Equal(t, Confirmed, ms[0].State)
		assert.Equal(t, Confirmed, ms[1].State)
	})
}

func mustInsertMsg(t *testing.T, txm *Txm, contractID string, msg cosmostypes.Msg) int64 {
	typeURL, raw, err := txm.marshalMsg(msg)
	require.NoError(t, err)
	id, err := txm.orm.InsertMsg(contractID, typeURL, raw)
	require.NoError(t, err)
	return id
}
