//go:build integration

package terratxm_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	wasmtypes "github.com/terra-money/core/x/wasm/types"
	"gopkg.in/guregu/null.v4"

	pkgterra "github.com/smartcontractkit/chainlink-terra/pkg/terra"
	terraclient "github.com/smartcontractkit/chainlink-terra/pkg/terra/client"
	"github.com/smartcontractkit/terra.go/msg"

	"github.com/smartcontractkit/chainlink/core/chains/terra"
	"github.com/smartcontractkit/chainlink/core/chains/terra/terratxm"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"

	. "github.com/smartcontractkit/chainlink-terra/pkg/terra/db"
)

func TestTxm_Integration(t *testing.T) {
	cfg, db := heavyweight.FullTestDBNoFixtures(t, "terra_txm")
	lggr := logger.TestLogger(t)
	chainID := fmt.Sprintf("Chainlinktest-%d", rand.Int31n(999999))
	logCfg := pgtest.NewPGCfg(true)
	fallbackGasPrice := sdk.NewDecCoinFromDec("uluna", sdk.MustNewDecFromStr("0.01"))
	gpe := terraclient.NewMustGasPriceEstimator([]terraclient.GasPricesEstimator{
		terraclient.NewFixedGasPriceEstimator(map[string]sdk.DecCoin{
			"uluna": fallbackGasPrice,
		}),
	}, lggr)
	dbChain, err := terra.NewORM(db, lggr, logCfg).CreateChain(chainID, &ChainCfg{
		FallbackGasPriceULuna: null.StringFrom(fallbackGasPrice.Amount.String()),
		GasLimitMultiplier:    null.FloatFrom(1.5),
	})
	require.NoError(t, err)
	chainCfg := pkgterra.NewConfig(*dbChain.Cfg, lggr)
	orm := terratxm.NewORM(chainID, db, lggr, logCfg)
	eb := pg.NewEventBroadcaster(cfg.DatabaseURL(), 0, 0, lggr, uuid.NewV4())
	require.NoError(t, eb.Start(testutils.Context(t)))
	t.Cleanup(func() { require.NoError(t, eb.Close()) })
	ks := keystore.New(db, utils.FastScryptParams, lggr, pgtest.NewPGCfg(true))
	accounts, testdir, tendermintURL := terraclient.SetupLocalTerraNode(t, "42")
	tc, err := terraclient.NewClient("42", tendermintURL, terra.DefaultRequestTimeout, lggr)
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
		an, sn, gpe.GasPrices()["uluna"], accounts[0].PrivateKey, txtypes.BroadcastMode_BROADCAST_MODE_BLOCK)
	require.NoError(t, err)

	// TODO: find a way to pull this test artifact from
	// the chainlink-terra repo instead of copying it to cores testdata
	contractID := terraclient.DeployTestContract(t, tendermintURL, accounts[0], terraclient.Account{
		Name:       "transmitter",
		PrivateKey: terratxm.NewKeyWrapper(transmitterKey),
		Address:    transmitterID,
	}, tc, testdir, "../../../testdata/terra/my_first_contract.wasm")

	tcFn := func() (terraclient.ReaderWriter, error) { return tc, nil }
	// Start txm
	txm := terratxm.NewTxm(db, tcFn, *gpe, chainID, chainCfg, ks.Terra(), lggr, pgtest.NewPGCfg(true), eb)
	require.NoError(t, txm.Start(testutils.Context(t)))

	// Change the contract state
	setMsg := wasmtypes.NewMsgExecuteContract(transmitterID, contractID, []byte(`{"reset":{"count":5}}`), sdk.Coins{})
	_, err = txm.Enqueue(contractID.String(), setMsg)
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
		msgs, err := orm.GetMsgsState(Confirmed, 5)
		require.NoError(t, err)
		return 1 == len(msgs)
	}, 5*time.Second, time.Second).Should(gomega.BeTrue())

	// Ensure invalid msgs are marked as errored
	invalidMsg := wasmtypes.NewMsgExecuteContract(transmitterID, contractID, []byte(`{"blah":{"blah":5}}`), sdk.Coins{})
	_, err = txm.Enqueue(contractID.String(), invalidMsg)
	require.NoError(t, err)
	_, err = txm.Enqueue(contractID.String(), invalidMsg)
	require.NoError(t, err)
	_, err = txm.Enqueue(contractID.String(), setMsg)
	require.NoError(t, err)

	// Ensure messages are completed
	gomega.NewWithT(t).Eventually(func() bool {
		succeeded, err := orm.GetMsgsState(Confirmed, 5)
		require.NoError(t, err)
		errored, err := orm.GetMsgsState(Errored, 5)
		require.NoError(t, err)
		t.Log("errored", len(errored), "succeeded", len(succeeded))
		return 2 == len(succeeded) && 2 == len(errored)
	}, 10*time.Second, time.Second).Should(gomega.BeTrue())

	// Observe the messages have been marked as completed
	require.NoError(t, txm.Close())
}
