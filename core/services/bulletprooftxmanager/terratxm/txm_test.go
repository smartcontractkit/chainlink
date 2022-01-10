package terratxm_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	terraclient "github.com/smartcontractkit/chainlink-terra/pkg/terra/client"
	"github.com/smartcontractkit/terra.go/msg"
	"github.com/stretchr/testify/assert"
	wasmtypes "github.com/terra-money/core/x/wasm/types"

	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/utils"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager/terratxm"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/stretchr/testify/require"
)

func TestTxmStartStop(t *testing.T) {
	t.Skip() // Local only unless we want to add terrad to CI env
	cfg, db := heavyweight.FullTestDB(t, "terra_txm", true, false)
	lggr := logger.TestLogger(t)
	eb := pg.NewEventBroadcaster(cfg.DatabaseURL(), 0, 0, lggr, uuid.NewV4())
	require.NoError(t, eb.Start())
	t.Cleanup(func() { require.NoError(t, eb.Close()) })
	ks := keystore.New(db, utils.FastScryptParams, lggr, pgtest.NewPGCfg(true))
	accounts, testdir := terraclient.SetupLocalTerraNode(t, "42")
	tc, err := terraclient.NewClient("42", "0.01", "1.5", "http://127.0.0.1:26657", "https://fcd.terra.dev/", time.Second, lggr)
	require.NoError(t, err)

	// First create a transmitter key and fund it with 1k uluna
	require.NoError(t, ks.Unlock("blah"))
	transmitterKey, err := ks.Terra().Create()
	require.NoError(t, err)
	transmitterID, err := msg.AccAddressFromBech32(transmitterKey.PublicKeyStr())
	require.NoError(t, err)
	an, sn, err := tc.Account(accounts[0].Address)
	require.NoError(t, err)
	resp, err := tc.SignAndBroadcast([]msg.Msg{msg.NewMsgSend(accounts[0].Address, transmitterID, msg.NewCoins(msg.NewInt64Coin("uluna", 100000)))},
		an, sn, tc.GasPrice(), accounts[0].PrivateKey, txtypes.BroadcastMode_BROADCAST_MODE_BLOCK)
	require.NoError(t, err)
	t.Log(resp.TxResponse)

	contractID := terraclient.DeployTestContract(t, accounts[0], terraclient.Account{
		Name:       "transmitter",
		PrivateKey: terratxm.NewPrivKey(transmitterKey),
		Address:    transmitterID,
	}, tc, testdir, "../../../testdata/my_first_contract.wasm")

	// Start txm
	txm := terratxm.NewTxm(db, tc, ks.Terra(), lggr, pgtest.NewPGCfg(true), eb, time.Second)
	require.NoError(t, txm.Start())

	// Change the contract state
	setMsg := wasmtypes.NewMsgExecuteContract(transmitterID, contractID, []byte(`{"reset":{"count":5}}`), sdk.Coins{})
	b, err := setMsg.Marshal()
	require.NoError(t, err)
	id, err := txm.Enqueue(contractID.String(), b)
	require.NoError(t, err)
	t.Log(id)

	// Observe the counter gets set eventually
	set := false
	for i := 0; i < 5; i++ {
		d, err := tc.ContractStore(contractID.String(), []byte(`{"get_count":{}}`))
		require.NoError(t, err)
		t.Log("contract value", string(d))
		if string(d) == `{"count":5}` {
			set = true
			break
		}
		time.Sleep(2 * time.Second)
	}
	assert.True(t, set)

	// Observe the messages have been marked as completed
	require.NoError(t, txm.Close())
}
