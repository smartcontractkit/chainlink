package cosmostxm_test

import (
	"fmt"
	"testing"
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	tmservicetypes "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	tmtypes "github.com/tendermint/tendermint/proto/tendermint/types"
	"go.uber.org/zap/zapcore"

	relayutils "github.com/smartcontractkit/chainlink-relay/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos"
	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos/cosmostxm"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/cosmostest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	cosmosclient "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/client"
	tcmocks "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/client/mocks"
	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
	. "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"
)

func generateExecuteMsg(t *testing.T, msg []byte, from, to cosmostypes.AccAddress) cosmostypes.Msg {
	return &wasmtypes.MsgExecuteContract{
		Sender:   from.String(),
		Contract: to.String(),
		Msg:      msg,
		Funds:    cosmostypes.Coins{},
	}
}

func newReaderWriterMock(t *testing.T) *tcmocks.ReaderWriter {
	tc := new(tcmocks.ReaderWriter)
	tc.Test(t)
	t.Cleanup(func() { tc.AssertExpectations(t) })
	return tc
}

func TestTxm(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	lggr := testutils.LoggerAssertMaxLevel(t, zapcore.ErrorLevel)
	ks := keystore.New(db, utils.FastScryptParams, lggr, pgtest.NewQConfig(true))
	require.NoError(t, ks.Unlock("blah"))
	k1, err := ks.Cosmos().Create()
	require.NoError(t, err)
	sender1, err := cosmostypes.AccAddressFromBech32(k1.PublicKeyStr())
	require.NoError(t, err)
	k2, err := ks.Cosmos().Create()
	require.NoError(t, err)
	sender2, err := cosmostypes.AccAddressFromBech32(k2.PublicKeyStr())
	require.NoError(t, err)
	contract, err := cosmostypes.AccAddressFromBech32("cosmos1z94322r480rhye2atp8z7v0wm37pk36ghzkdnd")
	require.NoError(t, err)
	contract2, err := cosmostypes.AccAddressFromBech32("cosmos1pe6e59rzm5upts599wl8hrvh95afy859yrcva8")
	require.NoError(t, err)
	logCfg := pgtest.NewQConfig(true)
	chainID := cosmostest.RandomChainID()
	two := int64(2)
	cfg := &cosmos.CosmosConfig{Chain: coscfg.Chain{
		MaxMsgsPerBatch: &two,
	}}
	cfg.SetDefaults()
	gpe := cosmosclient.NewMustGasPriceEstimator([]cosmosclient.GasPricesEstimator{
		cosmosclient.NewFixedGasPriceEstimator(map[string]cosmostypes.DecCoin{
			"uatom": cosmostypes.NewDecCoinFromDec("uatom", cosmostypes.MustNewDecFromStr("0.01")),
		}),
	}, lggr)

	t.Run("single msg", func(t *testing.T) {
		tc := newReaderWriterMock(t)
		tcFn := func() (cosmosclient.ReaderWriter, error) { return tc, nil }
		txm := cosmostxm.NewTxm(db, tcFn, *gpe, chainID, cfg, ks.Cosmos(), lggr, logCfg, nil)

		// Enqueue a single msg, then send it in a batch
		id1, err := txm.Enqueue(contract.String(), generateExecuteMsg(t, []byte(`1`), sender1, contract))
		require.NoError(t, err)
		tc.On("Account", mock.Anything).Return(uint64(0), uint64(0), nil)
		tc.On("BatchSimulateUnsigned", mock.Anything, mock.Anything).Return(&cosmosclient.BatchSimResults{
			Failed: nil,
			Succeeded: cosmosclient.SimMsgs{{ID: id1, Msg: &wasmtypes.MsgExecuteContract{
				Sender: sender1.String(),
				Msg:    []byte(`1`),
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
		txm.SendMsgBatch(testutils.Context(t))

		// Should be in completed state
		completed, err := txm.ORM().GetMsgs(id1)
		require.NoError(t, err)
		require.Equal(t, 1, len(completed))
		assert.Equal(t, completed[0].State, Confirmed)
	})

	t.Run("two msgs different accounts", func(t *testing.T) {
		tc := newReaderWriterMock(t)
		tcFn := func() (cosmosclient.ReaderWriter, error) { return tc, nil }
		txm := cosmostxm.NewTxm(db, tcFn, *gpe, chainID, cfg, ks.Cosmos(), lggr, pgtest.NewQConfig(true), nil)

		id1, err := txm.Enqueue(contract.String(), generateExecuteMsg(t, []byte(`0`), sender1, contract))
		require.NoError(t, err)
		id2, err := txm.Enqueue(contract.String(), generateExecuteMsg(t, []byte(`1`), sender2, contract))
		require.NoError(t, err)

		tc.On("Account", mock.Anything).Return(uint64(0), uint64(0), nil).Once()
		// Note this must be arg dependent, we don't know which order
		// the procesing will happen in (map iteration by from address).
		tc.On("BatchSimulateUnsigned", cosmosclient.SimMsgs{
			{
				ID: id2,
				Msg: &wasmtypes.MsgExecuteContract{
					Sender:   sender2.String(),
					Msg:      []byte(`1`),
					Contract: contract.String(),
				},
			},
		}, mock.Anything).Return(&cosmosclient.BatchSimResults{
			Failed: nil,
			Succeeded: cosmosclient.SimMsgs{
				{
					ID: id2,
					Msg: &wasmtypes.MsgExecuteContract{
						Sender:   sender2.String(),
						Msg:      []byte(`1`),
						Contract: contract.String(),
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
		txm.SendMsgBatch(testutils.Context(t))

		// Should be in completed state
		completed, err := txm.ORM().GetMsgs(id1, id2)
		require.NoError(t, err)
		require.Equal(t, 2, len(completed))
		assert.Equal(t, Errored, completed[0].State) // cancelled
		assert.Equal(t, Confirmed, completed[1].State)
	})

	t.Run("two msgs different contracts", func(t *testing.T) {
		tc := newReaderWriterMock(t)
		tcFn := func() (cosmosclient.ReaderWriter, error) { return tc, nil }
		txm := cosmostxm.NewTxm(db, tcFn, *gpe, chainID, cfg, ks.Cosmos(), lggr, pgtest.NewQConfig(true), nil)

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
			tc.On("BatchSimulateUnsigned", cosmosclient.SimMsgs{
				{
					ID: ids[i],
					Msg: &wasmtypes.MsgExecuteContract{
						Sender:   senders[i],
						Msg:      []byte(fmt.Sprintf(`%d`, i)),
						Contract: contracts[i],
					},
				},
			}, mock.Anything).Return(&cosmosclient.BatchSimResults{
				Failed: nil,
				Succeeded: cosmosclient.SimMsgs{
					{
						ID: ids[i],
						Msg: &wasmtypes.MsgExecuteContract{
							Sender:   senders[i],
							Msg:      []byte(fmt.Sprintf(`%d`, i)),
							Contract: contracts[i],
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
		txm.SendMsgBatch(testutils.Context(t))

		// Should be in completed state
		completed, err := txm.ORM().GetMsgs(id1, id2)
		require.NoError(t, err)
		require.Equal(t, 2, len(completed))
		assert.Equal(t, Confirmed, completed[0].State)
		assert.Equal(t, Confirmed, completed[1].State)
	})

	t.Run("failed to confirm", func(t *testing.T) {
		tc := newReaderWriterMock(t)
		tc.On("Tx", mock.Anything).Return(&txtypes.GetTxResponse{
			Tx:         &txtypes.Tx{},
			TxResponse: &cosmostypes.TxResponse{TxHash: "0x123"},
		}, errors.New("not found")).Twice()
		tcFn := func() (cosmosclient.ReaderWriter, error) { return tc, nil }
		txm := cosmostxm.NewTxm(db, tcFn, *gpe, chainID, cfg, ks.Cosmos(), lggr, pgtest.NewQConfig(true), nil)
		i, err := txm.ORM().InsertMsg("blah", "", []byte{0x01})
		require.NoError(t, err)
		txh := "0x123"
		require.NoError(t, txm.ORM().UpdateMsgs([]int64{i}, Started, &txh))
		require.NoError(t, txm.ORM().UpdateMsgs([]int64{i}, Broadcasted, &txh))
		err = txm.ConfirmTx(testutils.Context(t), tc, txh, []int64{i}, 2, 1*time.Millisecond)
		require.NoError(t, err)
		m, err := txm.ORM().GetMsgs(i)
		require.NoError(t, err)
		require.Equal(t, 1, len(m))
		assert.Equal(t, Errored, m[0].State)
	})

	t.Run("confirm any unconfirmed", func(t *testing.T) {
		require.Equal(t, int64(2), cfg.MaxMsgsPerBatch())
		txHash1 := "0x1234"
		txHash2 := "0x1235"
		txHash3 := "0xabcd"
		tc := newReaderWriterMock(t)
		tc.On("Tx", txHash1).Return(&txtypes.GetTxResponse{
			TxResponse: &cosmostypes.TxResponse{TxHash: txHash1},
		}, nil).Once()
		tc.On("Tx", txHash2).Return(&txtypes.GetTxResponse{
			TxResponse: &cosmostypes.TxResponse{TxHash: txHash2},
		}, nil).Once()
		tc.On("Tx", txHash3).Return(&txtypes.GetTxResponse{
			TxResponse: &cosmostypes.TxResponse{TxHash: txHash3},
		}, nil).Once()
		tcFn := func() (cosmosclient.ReaderWriter, error) { return tc, nil }
		txm := cosmostxm.NewTxm(db, tcFn, *gpe, chainID, cfg, ks.Cosmos(), lggr, pgtest.NewQConfig(true), nil)

		// Insert and broadcast 3 msgs with different txhashes.
		id1, err := txm.ORM().InsertMsg("blah", "", []byte{0x01})
		require.NoError(t, err)
		id2, err := txm.ORM().InsertMsg("blah", "", []byte{0x02})
		require.NoError(t, err)
		id3, err := txm.ORM().InsertMsg("blah", "", []byte{0x03})
		require.NoError(t, err)
		err = txm.ORM().UpdateMsgs([]int64{id1}, Started, &txHash1)
		require.NoError(t, err)
		err = txm.ORM().UpdateMsgs([]int64{id2}, Started, &txHash2)
		require.NoError(t, err)
		err = txm.ORM().UpdateMsgs([]int64{id3}, Started, &txHash3)
		require.NoError(t, err)
		err = txm.ORM().UpdateMsgs([]int64{id1}, Broadcasted, &txHash1)
		require.NoError(t, err)
		err = txm.ORM().UpdateMsgs([]int64{id2}, Broadcasted, &txHash2)
		require.NoError(t, err)
		err = txm.ORM().UpdateMsgs([]int64{id3}, Broadcasted, &txHash3)
		require.NoError(t, err)

		// Confirm them as in a restart while confirming scenario
		txm.ConfirmAnyUnconfirmed(testutils.Context(t))
		msgs, err := txm.ORM().GetMsgs(id1, id2, id3)
		require.NoError(t, err)
		require.Equal(t, 3, len(msgs))
		assert.Equal(t, Confirmed, msgs[0].State)
		assert.Equal(t, Confirmed, msgs[1].State)
		assert.Equal(t, Confirmed, msgs[2].State)
	})

	t.Run("expired msgs", func(t *testing.T) {
		tc := new(tcmocks.ReaderWriter)
		timeout, err := relayutils.NewDuration(1 * time.Millisecond)
		require.NoError(t, err)
		tcFn := func() (cosmosclient.ReaderWriter, error) { return tc, nil }
		two := int64(2)
		cfgShortExpiry := &cosmos.CosmosConfig{Chain: coscfg.Chain{
			MaxMsgsPerBatch: &two,
			TxMsgTimeout:    &timeout,
		}}
		cfgShortExpiry.SetDefaults()
		txm := cosmostxm.NewTxm(db, tcFn, *gpe, chainID, cfgShortExpiry, ks.Cosmos(), lggr, pgtest.NewQConfig(true), nil)

		// Send a single one expired
		id1, err := txm.ORM().InsertMsg("blah", "", []byte{0x03})
		require.NoError(t, err)
		time.Sleep(1 * time.Millisecond)
		txm.SendMsgBatch(testutils.Context(t))
		// Should be marked errored
		m, err := txm.ORM().GetMsgs(id1)
		require.NoError(t, err)
		assert.Equal(t, Errored, m[0].State)

		// Send a batch which is all expired
		id2, err := txm.ORM().InsertMsg("blah", "", []byte{0x03})
		require.NoError(t, err)
		id3, err := txm.ORM().InsertMsg("blah", "", []byte{0x03})
		require.NoError(t, err)
		time.Sleep(1 * time.Millisecond)
		txm.SendMsgBatch(testutils.Context(t))
		require.NoError(t, err)
		ms, err := txm.ORM().GetMsgs(id2, id3)
		require.NoError(t, err)
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
		tcFn := func() (cosmosclient.ReaderWriter, error) { return tc, nil }
		two := int64(2)
		cfgMaxMsgs := &cosmos.CosmosConfig{Chain: coscfg.Chain{
			MaxMsgsPerBatch: &two,
		}}
		cfgMaxMsgs.SetDefaults()
		txm := cosmostxm.NewTxm(db, tcFn, *gpe, chainID, cfgMaxMsgs, ks.Cosmos(), lggr, pgtest.NewQConfig(true), nil)

		// Leftover started is processed
		msg1 := generateExecuteMsg(t, []byte{0x03}, sender1, contract)
		id1 := mustInsertMsg(t, txm, contract.String(), msg1)
		require.NoError(t, txm.ORM().UpdateMsgs([]int64{id1}, Started, nil))
		msgs := cosmosclient.SimMsgs{{ID: id1, Msg: &wasmtypes.MsgExecuteContract{
			Sender:   sender1.String(),
			Msg:      []byte{0x03},
			Contract: contract.String(),
		}}}
		tc.On("BatchSimulateUnsigned", msgs, mock.Anything).
			Return(&cosmosclient.BatchSimResults{Failed: nil, Succeeded: msgs}, nil).Once()
		time.Sleep(1 * time.Millisecond)
		txm.SendMsgBatch(testutils.Context(t))
		m, err := txm.ORM().GetMsgs(id1)
		require.NoError(t, err)
		assert.Equal(t, Confirmed, m[0].State)

		// Leftover started is not cancelled
		msg2 := generateExecuteMsg(t, []byte{0x04}, sender1, contract)
		msg3 := generateExecuteMsg(t, []byte{0x05}, sender1, contract)
		id2 := mustInsertMsg(t, txm, contract.String(), msg2)
		require.NoError(t, txm.ORM().UpdateMsgs([]int64{id2}, Started, nil))
		time.Sleep(time.Millisecond) // ensure != CreatedAt
		id3 := mustInsertMsg(t, txm, contract.String(), msg3)
		msgs = cosmosclient.SimMsgs{{ID: id2, Msg: &wasmtypes.MsgExecuteContract{
			Sender:   sender1.String(),
			Msg:      []byte{0x04},
			Contract: contract.String(),
		}}, {ID: id3, Msg: &wasmtypes.MsgExecuteContract{
			Sender:   sender1.String(),
			Msg:      []byte{0x05},
			Contract: contract.String(),
		}}}
		tc.On("BatchSimulateUnsigned", msgs, mock.Anything).
			Return(&cosmosclient.BatchSimResults{Failed: nil, Succeeded: msgs}, nil).Once()
		time.Sleep(1 * time.Millisecond)
		txm.SendMsgBatch(testutils.Context(t))
		require.NoError(t, err)
		ms, err := txm.ORM().GetMsgs(id2, id3)
		require.NoError(t, err)
		assert.Equal(t, Confirmed, ms[0].State)
		assert.Equal(t, Confirmed, ms[1].State)
	})
}

func mustInsertMsg(t *testing.T, txm *cosmostxm.Txm, contractID string, msg cosmostypes.Msg) int64 {
	typeURL, raw, err := txm.MarshalMsg(msg)
	require.NoError(t, err)
	id, err := txm.ORM().InsertMsg(contractID, typeURL, raw)
	require.NoError(t, err)
	return id
}
