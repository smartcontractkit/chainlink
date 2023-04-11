//go:build integration && wasmd

package cosmostxm_test

import (
	"testing"
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	cosmosclient "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/client"
	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"

	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos"
	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos/cosmostxm"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	. "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"
)

func TestTxm_Integration(t *testing.T) {
	// TODO(BCI-978): cleanup once SetupLocalCosmosNode is updated
	chainID := cosmotest.RandomChainID()
	fallbackGasPrice := sdk.NewDecCoinFromDec("uatom", sdk.MustNewDecFromStr("0.01"))
	chain := cosmos.CosmosConfig{ChainID: &chainID, Enabled: ptr(true), Chain: coscfg.Chain{
		FallbackGasPriceUAtom: ptr(decimal.RequireFromString("0.01")),
		GasLimitMultiplier:    ptr(decimal.RequireFromString("1.5")),
	}}
	cfg, db := heavyweight.FullTestDBNoFixturesV2(t, "cosmos_txm", func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Cosmos = cosmos.CosmosConfigs{&chain}
	})
	lggr := logger.TestLogger(t)
	logCfg := pgtest.NewQConfig(true)
	gpe := cosmosclient.NewMustGasPriceEstimator([]cosmosclient.GasPricesEstimator{
		cosmosclient.NewFixedGasPriceEstimator(map[string]sdk.DecCoin{
			"uatom": fallbackGasPrice,
		}),
	}, lggr)
	orm := cosmostxm.NewORM(chainID, db, lggr, logCfg)
	eb := pg.NewEventBroadcaster(cfg.DatabaseURL(), 0, 0, lggr, uuid.NewV4())
	require.NoError(t, eb.Start(testutils.Context(t)))
	t.Cleanup(func() { require.NoError(t, eb.Close()) })
	ks := keystore.New(db, utils.FastScryptParams, lggr, pgtest.NewQConfig(true))
	accounts, testdir, tendermintURL := cosmosclient.SetupLocalCosmosNode(t, "42")
	tc, err := cosmosclient.NewClient("42", tendermintURL, cosmos.DefaultRequestTimeout, lggr)
	require.NoError(t, err)

	// First create a transmitter key and fund it with 1k uatom
	require.NoError(t, ks.Unlock("blah"))
	transmitterKey, err := ks.Cosmos().Create()
	require.NoError(t, err)
	transmitterID, err := sdk.AccAddressFromBech32(transmitterKey.PublicKeyStr())
	require.NoError(t, err)
	an, sn, err := tc.Account(accounts[0].Address)
	require.NoError(t, err)
	_, err = tc.SignAndBroadcast([]sdk.Msg{banktypes.NewMsgSend(accounts[0].Address, transmitterID, sdk.NewCoins(sdk.NewInt64Coin("uatom", 100000)))},
		an, sn, gpe.GasPrices()["uatom"], accounts[0].PrivateKey, txtypes.BroadcastMode_BROADCAST_MODE_BLOCK)
	require.NoError(t, err)

	// TODO: find a way to pull this test artifact from
	// the chainlink-cosmos repo instead of copying it to cores testdata
	contractID := cosmosclient.DeployTestContract(t, tendermintURL, accounts[0], cosmosclient.Account{
		Name:       "transmitter",
		PrivateKey: cosmostxm.NewKeyWrapper(transmitterKey),
		Address:    transmitterID,
	}, tc, testdir, "../../../testdata/cosmos/my_first_contract.wasm")

	tcFn := func() (cosmosclient.ReaderWriter, error) { return tc, nil }
	// Start txm
	txm := cosmostxm.NewTxm(db, tcFn, *gpe, chainID, &chain, ks.Cosmos(), lggr, pgtest.NewQConfig(true), eb)
	require.NoError(t, txm.Start(testutils.Context(t)))

	// Change the contract state
	setMsg := &wasmtypes.MsgExecuteContract{
		Sender:   transmitterID.String(),
		Contract: contractID.String(),
		Msg:      []byte(`{"reset":{"count":5}}`),
		Funds:    sdk.Coins{},
	}
	_, err = txm.Enqueue(contractID.String(), setMsg)
	require.NoError(t, err)

	// Observe the counter gets set eventually
	gomega.NewWithT(t).Eventually(func() bool {
		d, err := tc.ContractState(contractID, []byte(`{"get_count":{}}`))
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
	invalidMsg := &wasmtypes.MsgExecuteContract{
		Sender:   transmitterID.String(),
		Contract: contractID.String(),
		Msg:      []byte(`{"blah":{"blah":5}}`),
		Funds:    sdk.Coins{},
	}
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

func ptr[T any](t T) *T { return &t }
