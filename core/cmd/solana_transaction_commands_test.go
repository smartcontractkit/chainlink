//go:build integration

package cmd_test

import (
	"flag"
	"strconv"
	"testing"
	"time"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
	solanaClient "github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"

	"github.com/smartcontractkit/chainlink/core/chains/solana"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
)

func TestClient_SolanaSendSol(t *testing.T) {
	chainID := "localnet"
	url := solanaClient.SetupLocalSolNode(t)
	node := solcfg.Node{
		Name: ptr(t.Name()),
		URL:  utils.MustParseURL(url),
	}
	cfg := solana.SolanaConfig{
		ChainID: &chainID,
		Enabled: ptr(true),
		Nodes:   solana.SolanaNodes{&node},
	}
	app := solanaStartNewApplication(t, &cfg)
	from, err := app.GetKeyStore().Solana().Create()
	require.NoError(t, err)
	to, err := solanago.NewRandomPrivateKey()
	require.NoError(t, err)
	solanaClient.FundTestAccounts(t, []solanago.PublicKey{from.PublicKey()}, url)

	chain, err := app.GetChains().Solana.Chain(testutils.Context(t), chainID)
	require.NoError(t, err)

	reader, err := chain.Reader()
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		coin, err := reader.Balance(from.PublicKey())
		if !assert.NoError(t, err) {
			return false
		}
		return coin == 100*solanago.LAMPORTS_PER_SOL
	}, time.Minute, 5*time.Second)

	client, r := app.NewClientAndRenderer()
	cliapp := cli.NewApp()

	for _, tt := range []struct {
		amount string
		expErr string
	}{
		{amount: "1000000000"},
		{amount: "100000000000", expErr: "is too low for this transaction to be executed:"},
		{amount: "0", expErr: "amount must be greater than zero"},
		{amount: "asdf", expErr: "invalid amount:"},
	} {
		tt := tt
		t.Run(tt.amount, func(t *testing.T) {
			startBal, err := reader.Balance(from.PublicKey())
			require.NoError(t, err)

			set := flag.NewFlagSet("sendsolcoins", 0)
			cltest.CopyFlagSetFromAction(cmd.NewSolanaNodeClient(client).RemoveNode, set, "solana")

			require.NoError(t, set.Set("id", chainID))
			require.NoError(t, set.Parse([]string{tt.amount, from.PublicKey().String(), to.PublicKey().String()}))

			c := cli.NewContext(cliapp, set, nil)
			err = client.SolanaSendSol(c)
			if tt.expErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expErr)
				return
			}

			// Check CLI output
			require.Greater(t, len(r.Renders), 0)
			renderer := r.Renders[len(r.Renders)-1]
			renderedMsg := renderer.(*cmd.SolanaMsgPresenter)
			t.Logf("%+v\n", renderedMsg)
			require.NotEmpty(t, renderedMsg.ID)
			assert.Equal(t, chainID, renderedMsg.ChainID)
			assert.Equal(t, from.PublicKey().String(), renderedMsg.From)
			assert.Equal(t, to.PublicKey().String(), renderedMsg.To)
			assert.Equal(t, tt.amount, strconv.FormatUint(renderedMsg.Amount, 10))

			// wait for updated balance
			updated := false
			endBal := uint64(0)
			for i := 0; i < 5; i++ {
				time.Sleep(time.Second) // wait for tx execution

				// Check balance
				endBal, err = reader.Balance(from.PublicKey())
				require.NoError(t, err)

				// exit if difference found
				if endBal != startBal {
					updated = true
					break
				}
			}
			require.True(t, updated, "end bal == start bal, transaction likely not succeeded")

			// Check balance
			if assert.NotEqual(t, 0, startBal) && assert.NotEqual(t, 0, endBal) {
				diff := startBal - endBal
				receiveBal, err := reader.Balance(to.PublicKey())
				require.NoError(t, err)
				assert.Equal(t, tt.amount, strconv.FormatUint(receiveBal, 10))
				assert.Greater(t, diff, receiveBal)
			}
		})
	}
}
