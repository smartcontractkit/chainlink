package terratxm

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	tmservicetypes "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	tmtypes "github.com/tendermint/tendermint/proto/tendermint/types"
	wasmtypes "github.com/terra-money/core/x/wasm/types"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra"
	terraclient "github.com/smartcontractkit/chainlink-terra/pkg/terra/client"
	tcmocks "github.com/smartcontractkit/chainlink-terra/pkg/terra/client/mocks"
	terradb "github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/terratest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/utils"

	. "github.com/smartcontractkit/chainlink-terra/pkg/terra/db"
)

func generateExecuteMsg(t *testing.T, msg []byte, from, to cosmostypes.AccAddress) []byte {
	msg1 := wasmtypes.NewMsgExecuteContract(from, to, msg, cosmostypes.Coins{})
	d, err := msg1.Marshal()
	require.NoError(t, err)
	return d
}

func TestTxm(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
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
	logCfg := pgtest.NewPGCfg(true)
	chainID := fmt.Sprintf("Chainlinktest-%d", rand.Int31n(999999))
	terratest.MustInsertChain(t, db, &terradb.Chain{ID: chainID})
	require.NoError(t, err)
	cfg := terra.NewConfig(terradb.ChainCfg{}, terra.DefaultConfigSet, lggr)

	t.Run("single msg", func(t *testing.T) {
		tc := new(tcmocks.ReaderWriter)

		txm, _ := NewTxm(db, tc, chainID, cfg, ks.Terra(), lggr, logCfg, nil)

		// Enqueue a single msg, then send it in a batch
		id1, err := txm.Enqueue(contract.String(), generateExecuteMsg(t, []byte(`1`), sender1, contract))
		require.NoError(t, err)
		tc.On("Account", mock.Anything).Return(uint64(0), uint64(0), nil)
		tc.On("GasPrice", mock.Anything).Return(cosmostypes.NewDecCoinFromDec("uluna", cosmostypes.MustNewDecFromStr("0.01")))
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
		tc.On("Broadcast", mock.Anything, mock.Anything).Return(&txtypes.BroadcastTxResponse{
			TxResponse: &cosmostypes.TxResponse{TxHash: "0x123"},
		}, nil)
		tc.On("Tx", mock.Anything).Return(&txtypes.GetTxResponse{
			Tx:         &txtypes.Tx{},
			TxResponse: &cosmostypes.TxResponse{TxHash: "0x123"},
		}, nil)
		txm.sendMsgBatch()

		// Should be in completed state
		completed, err := txm.orm.SelectMsgsWithIDs([]int64{id1})
		require.NoError(t, err)
		require.Equal(t, 1, len(completed))
		assert.Equal(t, completed[0].State, Confirmed)
		tc.AssertExpectations(t)
	})

	t.Run("two msgs different accounts", func(t *testing.T) {
		tc := new(tcmocks.ReaderWriter)

		txm, _ := NewTxm(db, tc, chainID, cfg, ks.Terra(), lggr, pgtest.NewPGCfg(true), nil)

		id1, err := txm.Enqueue(contract.String(), generateExecuteMsg(t, []byte(`0`), sender1, contract))
		require.NoError(t, err)
		id2, err := txm.Enqueue(contract.String(), generateExecuteMsg(t, []byte(`1`), sender2, contract))
		require.NoError(t, err)
		ids := []int64{id1, id2}
		senders := []string{sender1.String(), sender2.String()}
		tc.On("GasPrice", mock.Anything).Return(cosmostypes.NewDecCoinFromDec("uluna", cosmostypes.MustNewDecFromStr("0.01"))).Once()
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
						Contract:   contract.String(),
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
			tc.On("Broadcast", mock.Anything, mock.Anything).Return(&txtypes.BroadcastTxResponse{
				TxResponse: &cosmostypes.TxResponse{TxHash: "0x123"},
			}, nil).Once()
			tc.On("Tx", mock.Anything).Return(&txtypes.GetTxResponse{
				Tx:         &txtypes.Tx{},
				TxResponse: &cosmostypes.TxResponse{TxHash: "0x123"},
			}, nil).Once()
		}
		txm.sendMsgBatch()

		// Should be in completed state
		completed, err := txm.orm.SelectMsgsWithIDs([]int64{id1, id2})
		require.NoError(t, err)
		require.Equal(t, 2, len(completed))
		assert.Equal(t, completed[0].State, Confirmed)
		assert.Equal(t, completed[1].State, Confirmed)
		tc.AssertExpectations(t)
	})

	t.Run("failed to confirm", func(t *testing.T) {
		tc := new(tcmocks.ReaderWriter)
		tc.On("Tx", mock.Anything).Return(&txtypes.GetTxResponse{
			Tx:         &txtypes.Tx{},
			TxResponse: &cosmostypes.TxResponse{TxHash: "0x123"},
		}, errors.New("not found")).Twice()
		cfg := terra.NewConfig(terradb.ChainCfg{}, terra.DefaultConfigSet, lggr)
		txm, _ := NewTxm(db, tc, chainID, cfg, ks.Terra(), lggr, pgtest.NewPGCfg(true), nil)
		i, err := txm.orm.InsertMsg("blah", []byte{0x01})
		require.NoError(t, err)
		txh := "0x123"
		require.NoError(t, txm.orm.UpdateMsgsWithState([]int64{i}, Broadcasted, &txh))
		err = txm.confirmTx(txh, []int64{i}, 2, 1*time.Millisecond)
		require.NoError(t, err)
		m, err := txm.orm.SelectMsgsWithIDs([]int64{i})
		require.NoError(t, err)
		require.Equal(t, 1, len(m))
		assert.Equal(t, Errored, m[0].State)
		tc.AssertExpectations(t)
	})

	t.Run("confirm any unconfirmed", func(t *testing.T) {
		txHash1 := "0x1234"
		txHash2 := "0x1235"
		tc := new(tcmocks.ReaderWriter)
		tc.On("Tx", txHash1).Return(&txtypes.GetTxResponse{
			TxResponse: &cosmostypes.TxResponse{TxHash: txHash1},
		}, nil).Once()
		tc.On("Tx", txHash2).Return(&txtypes.GetTxResponse{
			TxResponse: &cosmostypes.TxResponse{TxHash: txHash2},
		}, nil).Once()
		txm, _ := NewTxm(db, tc, chainID, cfg, ks.Terra(), lggr, pgtest.NewPGCfg(true), nil)

		// Insert and broadcast 2 msgs with different txhashes.
		id1, err := txm.orm.InsertMsg("blah", []byte{0x01})
		require.NoError(t, err)
		id2, err := txm.orm.InsertMsg("blah", []byte{0x02})
		require.NoError(t, err)
		err = txm.orm.UpdateMsgsWithState([]int64{id1}, Broadcasted, &txHash1)
		require.NoError(t, err)
		err = txm.orm.UpdateMsgsWithState([]int64{id2}, Broadcasted, &txHash2)
		require.NoError(t, err)

		// Confirm them as in a restart while confirming scenario
		txm.confirmAnyUnconfirmed()
		require.NoError(t, err)
		confirmed, err := txm.orm.SelectMsgsWithIDs([]int64{id1, id2})
		require.NoError(t, err)
		require.Equal(t, 2, len(confirmed))
		assert.Equal(t, Confirmed, confirmed[0].State)
		assert.Equal(t, Confirmed, confirmed[1].State)
		tc.AssertExpectations(t)
	})
}
