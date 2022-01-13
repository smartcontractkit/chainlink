package terratxm_test

import (
	"testing"
	"time"

	"github.com/onsi/gomega"

	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	terraclient "github.com/smartcontractkit/chainlink-terra/pkg/terra/client"
	"github.com/smartcontractkit/terra.go/msg"
	wasmtypes "github.com/terra-money/core/x/wasm/types"

	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/utils"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/chains/terra/terratxm"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/stretchr/testify/require"
)

func TestTxmStartStop(t *testing.T) {
	t.Skip() // Local only unless we want to add terrad to CI env
	cfg, db := heavyweight.FullTestDB(t, "terra_txm", true, false)
	lggr := logger.TestLogger(t)
	orm := terratxm.NewORM(db, lggr, pgtest.NewPGCfg(true))
	eb := pg.NewEventBroadcaster(cfg.DatabaseURL(), 0, 0, lggr, uuid.NewV4())
	require.NoError(t, eb.Start())
	t.Cleanup(func() { require.NoError(t, eb.Close()) })
	ks := keystore.New(db, utils.FastScryptParams, lggr, pgtest.NewPGCfg(true))
	accounts, testdir := terraclient.SetupLocalTerraNode(t, "42")
	fallbackGasPrice := sdk.NewDecCoinFromDec("ulunua", sdk.MustNewDecFromStr("0.01"))
	gasLimitMultiplier := 1.5
	tc, err := terraclient.NewClient("42", "http://127.0.0.1:26657", "https://fcd.terra.dev/", 10, lggr)
	require.NoError(t, err)

	// First create a transmitter key and fund it with 1k uluna
	require.NoError(t, ks.Unlock("blah"))
	transmitterKey, err := ks.Terra().Create()
	require.NoError(t, err)
	transmitterID, err := msg.AccAddressFromBech32(transmitterKey.PublicKeyStr())
	require.NoError(t, err)
	an, sn, err := tc.Account(accounts[0].Address)
	require.NoError(t, err)
	_, err = tc.SignAndBroadcast([]msg.Msg{msg.NewMsgSend(accounts[0].Address, transmitterID, msg.NewCoins(msg.NewInt64Coin("uluna", 100000)))},
		an, sn, tc.GasPrice(fallbackGasPrice), accounts[0].PrivateKey, txtypes.BroadcastMode_BROADCAST_MODE_BLOCK)
	require.NoError(t, err)

	contractID := terraclient.DeployTestContract(t, accounts[0], terraclient.Account{
		Name:       "transmitter",
		PrivateKey: terratxm.NewKeyWrapper(transmitterKey),
		Address:    transmitterID,
	}, tc, testdir, "../../../testdata/my_first_contract.wasm")

	// Start txm
	txm, err := terratxm.NewTxm(db, tc, fallbackGasPrice.Amount.String(), gasLimitMultiplier, ks.Terra(), lggr, pgtest.NewPGCfg(true), eb, 5*time.Second)
	require.NoError(t, err)
	require.NoError(t, txm.Start())

	// Change the contract state
	setMsg := wasmtypes.NewMsgExecuteContract(transmitterID, contractID, []byte(`{"reset":{"count":5}}`), sdk.Coins{})
	validBytes, err := setMsg.Marshal()
	require.NoError(t, err)
	_, err = txm.Enqueue(contractID.String(), validBytes)
	require.NoError(t, err)

	// Observe the counter gets set eventually
	gomega.NewWithT(t).Eventually(func() bool {
		d, err := tc.ContractStore(contractID, []byte(`{"get_count":{}}`))
		require.NoError(t, err)
		t.Log("contract value", string(d))
		return string(d) == `{"count":5}`
	}, 10*time.Second, time.Second).Should(gomega.BeTrue())
	// Ensure messages are completed
	gomega.NewWithT(t).Eventually(func() bool {
		msgs, err := orm.SelectMsgsWithState(terratxm.Confirmed)
		require.NoError(t, err)
		return 1 == len(msgs)
	}, 5*time.Second, time.Second).Should(gomega.BeTrue())

	// Ensure invalid msgs are marked as errored
	invalidMsg := wasmtypes.NewMsgExecuteContract(transmitterID, contractID, []byte(`{"blah":{"blah":5}}`), sdk.Coins{})
	invalidBytes, err := invalidMsg.Marshal()
	require.NoError(t, err)
	_, err = txm.Enqueue(contractID.String(), invalidBytes)
	_, err = txm.Enqueue(contractID.String(), invalidBytes)
	_, err = txm.Enqueue(contractID.String(), validBytes)
	require.NoError(t, err)

	// Ensure messages are completed
	gomega.NewWithT(t).Eventually(func() bool {
		succeeded, err := orm.SelectMsgsWithState(terratxm.Confirmed)
		require.NoError(t, err)
		errored, err := orm.SelectMsgsWithState(terratxm.Errored)
		require.NoError(t, err)
		t.Log("errored", len(errored), "succeeded", len(succeeded))
		return 2 == len(succeeded) && 2 == len(errored)
	}, 10*time.Second, time.Second).Should(gomega.BeTrue())

	// Observe the messages have been marked as completed
	require.NoError(t, txm.Close())
}
