package cmd_test

import (
	"flag"
	"testing"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestClient_IndexChains(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	orm := app.EVMORM()

	id := utils.NewBigI(99)
	chain, err := orm.CreateChain(*id, types.ChainCfg{})
	require.NoError(t, err)

	require.Nil(t, client.IndexChains(cltest.EmptyCLIContext()))
	chains := *r.Renders[0].(*cmd.ChainPresenters)
	require.Len(t, chains, 2) // includes default chain
	c := chains[1]            // chains 0 is the default chain (id = 1)
	assert.Equal(t, chain.ID.ToInt().String(), c.ID)
}

func TestClient_CreateChain(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	orm := app.EVMORM()

	set := flag.NewFlagSet("cli", 0)
	set.Int64("id", 99, "")
	set.Parse([]string{`{}`})
	c := cli.NewContext(nil, set, nil)

	err := client.CreateChain(c)
	require.NoError(t, err)

	chains, _, err := orm.Chains(0, 25)
	require.Len(t, chains, 2) // includes default chain
	ch := chains[1]           // chains 0 is the default chain (id = 1)
	assert.Equal(t, int64(99), ch.ID.ToInt().Int64())
}

func TestClient_RemoveChain(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	orm := app.EVMORM()

	id := utils.NewBigI(99)
	_, err := orm.CreateChain(*id, types.ChainCfg{})
	require.NoError(t, err)
	chains, _, err := orm.Chains(0, 25)
	require.Len(t, chains, 2) // includes default chain

	set := flag.NewFlagSet("cli", 0)
	set.Parse([]string{"99"})
	c := cli.NewContext(nil, set, nil)

	err = client.RemoveChain(c)
	require.NoError(t, err)

	chains, _, err = orm.Chains(0, 25)
	require.Len(t, chains, 1) // includes default chain
}
